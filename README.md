# proxy-benchmark

Measures how much latency an image proxy adds compared to fetching the same
image straight from the origin. Supports [imgproxy](https://imgproxy.net) and
[ipx](https://github.com/unjs/ipx) behind one harness so their numbers are
comparable.

The benchmark starts a local origin HTTP server (serving generated test images,
or your own via `-images`), then for each image runs these scenarios:

| scenario   | imgproxy                                      | ipx                                            |
|------------|-----------------------------------------------|------------------------------------------------|
| `direct`   | `GET origin/<img>` (baseline, no proxy)       | same                                           |
| `raw`      | `/insecure/raw:1/plain/<src>`: no processing  | not supported, skipped (ipx always re-encodes) |
| `reencode` | `/insecure/plain/<src>`: decode + re-encode   | `/_/<src>`                                     |
| `resize`   | `/insecure/rs:fit:300:300/plain/<src>`        | `/s_300x300,fit_inside/<src>`                  |
| `webp`     | `/insecure/rs:fit:300:300/plain/<src>@webp`   | `/s_300x300,fit_inside,f_webp/<src>`           |
| `avif`     | `/insecure/rs:fit:300:300/plain/<src>@avif`   | `/s_300x300,fit_inside,f_avif/<src>`           |

It reports TTFB and total time (including reading the body) at p50/p90/p99/max,
throughput, response size, and p50 overhead vs `direct`.

## Layout

```
main.go           the benchmark (go run . -proxy imgproxy|ipx ...)
imgproxy/run.sh   start imgproxy in Docker, benchmark it, tear it down
imgproxy/Note.md  results from the original imgproxy run (M3, 8GB)
ipx/Dockerfile    ipx image (no official one exists), same shape as imgproxy's
ipx/run.sh        build + start ipx in Docker, benchmark it, tear it down
```

## Run

Requires Docker and Go.

```sh
imgproxy/run.sh                                # imgproxy:latest, 200 req, concurrency 1
ipx/run.sh                                     # ipx 4.0.0, same settings
ipx/run.sh -n 500 -c 8 -json results-ipx.json  # more requests, under load
imgproxy/run.sh -images ~/photos -scenarios direct,raw,webp

IMGPROXY_IMAGE=ghcr.io/imgproxy/imgproxy:v3.28 imgproxy/run.sh
IMGPROXY_DOCKER_ARGS="--cpus=2 -e IMGPROXY_WORKERS=2" imgproxy/run.sh
IPX_VERSION=3.1.2 ipx/run.sh
IPX_DOCKER_ARGS="--cpus=2" ipx/run.sh
```

Both runners publish the container on `127.0.0.1:8080`, give it
`host.docker.internal` for the origin, and allow fetching from it
(imgproxy: `IMGPROXY_ALLOW_PRIVATE_SOURCE_ADDRESSES`; ipx: `IPX_HTTP_DOMAINS`).
Run them one after the other, not concurrently.

To benchmark a proxy you're already running:

```sh
go run . -proxy ipx -proxy-url http://host:port -origin-proxy http://<origin as seen by the proxy>:9000
```

## Comparing imgproxy and ipx

- ipx has no passthrough mode. Every request, including `/_/`, is decoded and
  re-encoded by sharp, so `reencode` is its floor. The `raw` scenario is
  skipped for ipx; asking for it explicitly (`-scenarios raw`) is an error.
- ipx issues a `HEAD` to the origin (for `last-modified`/`cache-control`) and
  then a `GET` for the bytes, so each ipx request costs two origin round-trips
  where imgproxy costs one. That is part of what the benchmark measures.
- `reencode` output sizes differ between the two (encoder defaults differ, e.g.
  sharp's PNG output can be larger than the source), so compare the time
  columns, not `bytes`.
- Both run on glibc images (imgproxy: Ubuntu; ipx: `node:22-bookworm-slim`).
  sharp's bundled libvips performs worse under musl, so the ipx image is
  deliberately not Alpine-based.

## Caveats

- On macOS, Docker Desktop runs containers in a VM, so proxied scenarios
  include two extra VM network hops (client→container, container→origin) that
  `direct` doesn't. Treat `raw - direct` (imgproxy) as an upper bound on proxy
  overhead; run on Linux for numbers close to production.
- The origin is local and serves from the page cache, so `direct` is a
  best case. Real origins (S3, CDN) add network latency that applies to both
  paths, which shrinks the *relative* overhead.
- `-c 1` measures pure latency; increase `-c` to see queueing once requests
  exceed the proxy's worker count.
