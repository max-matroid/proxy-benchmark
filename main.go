// Command proxy-benchmark compares the latency of fetching images directly
// from an origin server against fetching them through an image proxy
// (imgproxy or ipx) with various processing options.
//
// It starts a local origin HTTP server that serves a set of test images
// (generated on the fly unless -images is given), then for every image runs
// each scenario: a direct fetch from the origin, and several proxy URLs that
// point back at the same origin.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptrace"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
)

// proxy describes how to talk to one image proxy implementation.
type proxy struct {
	Name string
	// HealthPath is requested once before the run to make sure the proxy is up.
	HealthPath string
	// HealthOK reports whether a response to HealthPath means the proxy is ready.
	HealthOK func(status int) bool
	// BuildURL turns processing options, a source URL and an optional output
	// format into a request URL for this proxy.
	BuildURL func(base, opts, src, ext string) string
}

var proxies = map[string]proxy{
	"imgproxy": {
		Name:       "imgproxy",
		HealthPath: "/health",
		HealthOK:   func(s int) bool { return s == http.StatusOK },
		// /insecure/<opts>/plain/<src>@<ext>; no IMGPROXY_KEY is set, so any
		// signature segment is accepted.
		BuildURL: func(base, opts, src, ext string) string {
			if ext != "" {
				src += "@" + ext
			}
			parts := []string{base, "insecure"}
			if opts != "" {
				parts = append(parts, opts)
			}
			parts = append(parts, "plain", src)
			return strings.Join(parts, "/")
		},
	},
	"ipx": {
		Name:       "ipx",
		HealthPath: "/",
		// ipx has no health endpoint; "/" answers 400 (missing modifiers)
		// once the server is up, which is good enough.
		HealthOK: func(s int) bool { return s < 500 },
		// /<modifiers>/<src>; "_" means no modifiers.
		BuildURL: func(base, opts, src, ext string) string {
			mods := []string{}
			if opts != "" {
				mods = append(mods, opts)
			}
			if ext != "" {
				mods = append(mods, "f_"+ext)
			}
			seg := "_"
			if len(mods) > 0 {
				seg = strings.Join(mods, ",")
			}
			return base + "/" + seg + "/" + src
		},
	},
}

type scenario struct {
	Name string
	// Options maps a proxy name to its processing options for this scenario
	// ("" = no processing). A proxy that is not listed cannot run the scenario.
	Options map[string]string
	// Ext is an optional output format.
	Ext string
}

var allScenarios = []scenario{
	{Name: "direct"},
	{Name: "raw", Options: map[string]string{
		"imgproxy": "raw:1",
		// ipx has no passthrough mode: every request is decoded and re-encoded.
	}},
	{Name: "reencode", Options: map[string]string{
		"imgproxy": "",
		"ipx":      "",
	}},
	{Name: "resize", Options: map[string]string{
		"imgproxy": "rs:fit:300:300",
		"ipx":      "s_300x300,fit_inside",
	}},
	{Name: "webp", Ext: "webp", Options: map[string]string{
		"imgproxy": "rs:fit:300:300",
		"ipx":      "s_300x300,fit_inside",
	}},
	{Name: "avif", Ext: "avif", Options: map[string]string{
		"imgproxy": "rs:fit:300:300",
		"ipx":      "s_300x300,fit_inside",
	}},
}

const defaultScenarios = "direct,raw,reencode,resize,webp,avif"

type sample struct {
	TTFB  time.Duration
	Total time.Duration
	Bytes int64
}

type result struct {
	Proxy     string  `json:"proxy"`
	Image     string  `json:"image"`
	Scenario  string  `json:"scenario"`
	Requests  int     `json:"requests"`
	Errors    int     `json:"errors"`
	RPS       float64 `json:"rps"`
	AvgBytes  int64   `json:"avg_bytes"`
	TTFBP50   float64 `json:"ttfb_p50_ms"`
	TotalMin  float64 `json:"total_min_ms"`
	TotalMean float64 `json:"total_mean_ms"`
	TotalP50  float64 `json:"total_p50_ms"`
	TotalP90  float64 `json:"total_p90_ms"`
	TotalP99  float64 `json:"total_p99_ms"`
	TotalMax  float64 `json:"total_max_ms"`
}

func main() {
	var (
		proxyName   = flag.String("proxy", "imgproxy", "proxy under test: imgproxy or ipx")
		proxyURL    = flag.String("proxy-url", "http://127.0.0.1:8080", "proxy base URL")
		originAddr  = flag.String("origin-listen", "0.0.0.0:9000", "address the local origin server listens on")
		originLocal = flag.String("origin-direct", "http://127.0.0.1:9000", "origin base URL used for direct fetches")
		originProxy = flag.String("origin-proxy", "http://host.docker.internal:9000", "origin base URL as seen from the proxy")
		imagesDir   = flag.String("images", "", "directory with source images (default: generate test images)")
		scenarios   = flag.String("scenarios", defaultScenarios, "comma-separated scenarios to run")
		requests    = flag.Int("n", 200, "measured requests per image/scenario")
		warmup      = flag.Int("warmup", 10, "warm-up requests per image/scenario (not measured)")
		concurrency = flag.Int("c", 1, "concurrent requests")
		timeout     = flag.Duration("timeout", 30*time.Second, "per-request timeout")
		jsonOut     = flag.String("json", "", "write results as JSON to this file")
	)
	flag.Parse()

	px, ok := proxies[*proxyName]
	if !ok {
		log.Fatalf("unknown proxy %q (want imgproxy or ipx)", *proxyName)
	}

	selected, err := selectScenarios(*scenarios, px, *scenarios == defaultScenarios)
	if err != nil {
		log.Fatal(err)
	}

	dir := *imagesDir
	if dir == "" {
		dir, err = os.MkdirTemp("", "proxy-benchmark-*")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(dir)
		if err := generateImages(dir); err != nil {
			log.Fatal(err)
		}
	}

	images, err := listImages(dir)
	if err != nil {
		log.Fatal(err)
	}
	if len(images) == 0 {
		log.Fatalf("no images found in %s", dir)
	}

	ln, err := net.Listen("tcp", *originAddr)
	if err != nil {
		log.Fatal(err)
	}
	go http.Serve(ln, http.FileServer(http.Dir(dir)))

	client := &http.Client{
		Timeout: *timeout,
		Transport: &http.Transport{
			MaxIdleConns:        *concurrency * 2,
			MaxIdleConnsPerHost: *concurrency * 2,
			DisableCompression:  true,
		},
	}

	if err := checkReachable(client, *proxyURL+px.HealthPath, px.HealthOK); err != nil {
		log.Fatalf("%s is not reachable at %s: %v", px.Name, *proxyURL, err)
	}

	var results []result
	for _, img := range images {
		for _, sc := range selected {
			url := buildURL(px, sc, img, *proxyURL, *originLocal, *originProxy)

			// Fail fast on a misconfigured scenario instead of measuring errors.
			if s, err := fetch(client, url); err != nil {
				log.Fatalf("%s/%s: %v (%s)", img, sc.Name, err, url)
			} else if s.Bytes == 0 {
				log.Fatalf("%s/%s: empty response (%s)", img, sc.Name, url)
			}

			run(client, url, *warmup, *concurrency)
			samples, errs, elapsed := run(client, url, *requests, *concurrency)

			r := summarize(px.Name, img, sc.Name, samples, errs, elapsed)
			results = append(results, r)
			fmt.Fprintf(os.Stderr, "%-14s %-9s p50=%7.2fms p99=%7.2fms errors=%d\n",
				img, sc.Name, r.TotalP50, r.TotalP99, r.Errors)
		}
	}

	fmt.Println()
	printResults(px.Name, results, *requests, *concurrency)

	if *jsonOut != "" {
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(*jsonOut, data, 0o644); err != nil {
			log.Fatal(err)
		}
	}
}

// selectScenarios resolves a comma-separated scenario list. Scenarios the
// proxy cannot run are skipped with a notice when the list is the default,
// and are an error when the user asked for them explicitly.
func selectScenarios(list string, px proxy, isDefault bool) ([]scenario, error) {
	var out []scenario
	for _, name := range strings.Split(list, ",") {
		name = strings.TrimSpace(name)
		i := slices.IndexFunc(allScenarios, func(s scenario) bool { return s.Name == name })
		if i < 0 {
			return nil, fmt.Errorf("unknown scenario %q", name)
		}
		sc := allScenarios[i]
		if _, ok := sc.Options[px.Name]; !ok && sc.Name != "direct" {
			if !isDefault {
				return nil, fmt.Errorf("scenario %q is not supported by %s", name, px.Name)
			}
			fmt.Fprintf(os.Stderr, "skipping scenario %q: not supported by %s\n", name, px.Name)
			continue
		}
		out = append(out, sc)
	}
	return out, nil
}

func buildURL(px proxy, sc scenario, img, proxyURL, originLocal, originProxy string) string {
	if sc.Name == "direct" {
		return originLocal + "/" + img
	}
	return px.BuildURL(proxyURL, sc.Options[px.Name], originProxy+"/"+img, sc.Ext)
}

func checkReachable(client *http.Client, url string, ok func(int) bool) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	if !ok(resp.StatusCode) {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}

func fetch(client *http.Client, url string) (sample, error) {
	var ttfb time.Time
	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() { ttfb = time.Now() },
	}

	req, err := http.NewRequestWithContext(
		httptrace.WithClientTrace(context.Background(), trace), http.MethodGet, url, nil,
	)
	if err != nil {
		return sample{}, err
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return sample{}, err
	}
	defer resp.Body.Close()

	n, err := io.Copy(io.Discard, resp.Body)
	total := time.Since(start)
	if err != nil {
		return sample{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return sample{}, fmt.Errorf("status %d", resp.StatusCode)
	}

	return sample{TTFB: ttfb.Sub(start), Total: total, Bytes: n}, nil
}

// run issues n requests to url using c concurrent workers.
func run(client *http.Client, url string, n, c int) ([]sample, int, time.Duration) {
	var (
		mu      sync.Mutex
		samples = make([]sample, 0, n)
		errs    int
		wg      sync.WaitGroup
		jobs    = make(chan struct{})
	)

	start := time.Now()
	for range c {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				s, err := fetch(client, url)
				mu.Lock()
				if err != nil {
					errs++
				} else {
					samples = append(samples, s)
				}
				mu.Unlock()
			}
		}()
	}
	for range n {
		jobs <- struct{}{}
	}
	close(jobs)
	wg.Wait()

	return samples, errs, time.Since(start)
}

func summarize(proxyName, img, name string, samples []sample, errs int, elapsed time.Duration) result {
	r := result{
		Proxy:    proxyName,
		Image:    img,
		Scenario: name,
		Requests: len(samples) + errs,
		Errors:   errs,
		RPS:      float64(len(samples)) / elapsed.Seconds(),
	}
	if len(samples) == 0 {
		return r
	}

	totals := make([]time.Duration, len(samples))
	ttfbs := make([]time.Duration, len(samples))
	var sum time.Duration
	var bytes int64
	for i, s := range samples {
		totals[i] = s.Total
		ttfbs[i] = s.TTFB
		sum += s.Total
		bytes += s.Bytes
	}
	slices.Sort(totals)
	slices.Sort(ttfbs)

	r.AvgBytes = bytes / int64(len(samples))
	r.TTFBP50 = ms(percentile(ttfbs, 50))
	r.TotalMin = ms(totals[0])
	r.TotalMean = ms(sum / time.Duration(len(samples)))
	r.TotalP50 = ms(percentile(totals, 50))
	r.TotalP90 = ms(percentile(totals, 90))
	r.TotalP99 = ms(percentile(totals, 99))
	r.TotalMax = ms(totals[len(totals)-1])
	return r
}

// percentile uses the nearest-rank method on a sorted slice.
func percentile(sorted []time.Duration, p float64) time.Duration {
	rank := int(math.Ceil(p / 100 * float64(len(sorted))))
	return sorted[max(rank-1, 0)]
}

func ms(d time.Duration) float64 {
	return float64(d.Microseconds()) / 1000
}

func printResults(proxyName string, results []result, n, c int) {
	fmt.Printf("%s latency (total time incl. body read), %d requests, concurrency %d\n\n", proxyName, n, c)

	direct := map[string]float64{}
	for _, r := range results {
		if r.Scenario == "direct" {
			direct[r.Image] = r.TotalP50
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "image\tscenario\tbytes\tttfb p50\tp50\tp90\tp99\tmax\trps\tp50 vs direct\terrors\t")
	for _, r := range results {
		overhead := "-"
		if d, ok := direct[r.Image]; ok && r.Scenario != "direct" {
			overhead = fmt.Sprintf("+%.2fms (%.1fx)", r.TotalP50-d, r.TotalP50/d)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.0f\t%s\t%d\t\n",
			r.Image, r.Scenario, humanBytes(r.AvgBytes), r.TTFBP50,
			r.TotalP50, r.TotalP90, r.TotalP99, r.TotalMax, r.RPS, overhead, r.Errors)
	}
	w.Flush()
	fmt.Println("\nAll times in milliseconds.")
}

func humanBytes(b int64) string {
	switch {
	case b >= 1<<20:
		return fmt.Sprintf("%.1fMB", float64(b)/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1fKB", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%dB", b)
	}
}

func listImages(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		switch strings.ToLower(filepath.Ext(e.Name())) {
		case ".jpg", ".jpeg", ".png", ".webp", ".avif", ".gif":
			out = append(out, e.Name())
		}
	}
	return out, nil
}

// generateImages writes deterministic photo-like test images of several
// sizes. Gradients plus noise keep them from compressing unrealistically well.
func generateImages(dir string) error {
	specs := []struct {
		name string
		w, h int
	}{
		{"small.jpg", 640, 480},
		{"medium.jpg", 1920, 1080},
		{"large.jpg", 4032, 3024},
		{"medium.png", 1920, 1080},
	}

	for _, s := range specs {
		img := syntheticImage(s.w, s.h)
		f, err := os.Create(filepath.Join(dir, s.name))
		if err != nil {
			return err
		}
		if strings.HasSuffix(s.name, ".png") {
			err = png.Encode(f, img)
		} else {
			err = jpeg.Encode(f, img, &jpeg.Options{Quality: 85})
		}
		if cerr := f.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func syntheticImage(w, h int) image.Image {
	rng := rand.New(rand.NewSource(42))
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			fx, fy := float64(x)/float64(w), float64(y)/float64(h)
			n := float64(rng.Intn(24) - 12)
			img.SetRGBA(x, y, color.RGBA{
				R: clamp(255*fx + 40*math.Sin(fy*12) + n),
				G: clamp(255*fy + 40*math.Cos(fx*9) + n),
				B: clamp(128 + 100*math.Sin((fx+fy)*7) + n),
				A: 255,
			})
		}
	}
	return img
}

func clamp(v float64) uint8 {
	return uint8(max(0, min(255, v)))
}
