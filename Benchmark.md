
Running on Mac M3 with 8GB RAM

`imgproxy/run.sh`:

```go
large.jpg      direct    p50=   0.28ms p99=   0.78ms errors=0
large.jpg      raw       p50=  12.05ms p99=  21.94ms errors=0
large.jpg      reencode  p50=  69.66ms p99=  87.55ms errors=0
large.jpg      resize    p50=  24.00ms p99=  89.91ms errors=0
large.jpg      webp      p50=  27.40ms p99=  31.24ms errors=0
large.jpg      avif      p50=  25.98ms p99=  48.20ms errors=0
medium.jpg     direct    p50=   0.20ms p99=   0.89ms errors=0
medium.jpg     raw       p50=   2.35ms p99=   4.97ms errors=0
medium.jpg     reencode  p50=  14.22ms p99=  36.65ms errors=0
medium.jpg     resize    p50=   6.16ms p99=   8.70ms errors=0
medium.jpg     webp      p50=   8.91ms p99=  12.08ms errors=0
medium.jpg     avif      p50=   7.72ms p99=  17.41ms errors=0
medium.png     direct    p50=   0.38ms p99=   2.68ms errors=0
medium.png     raw       p50=  15.52ms p99=  22.57ms errors=0
medium.png     reencode  p50= 117.70ms p99= 143.81ms errors=0
medium.png     resize    p50=  27.00ms p99=  33.89ms errors=0
medium.png     webp      p50=  28.79ms p99=  33.14ms errors=0
medium.png     avif      p50=  27.67ms p99=  39.01ms errors=0
small.jpg      direct    p50=   0.16ms p99=   1.20ms errors=0
small.jpg      raw       p50=   0.62ms p99=   2.38ms errors=0
small.jpg      reencode  p50=   3.21ms p99=   5.60ms errors=0
small.jpg      resize    p50=   2.77ms p99=   7.36ms errors=0
small.jpg      webp      p50=   7.25ms p99=  10.41ms errors=0
small.jpg      avif      p50=   4.82ms p99=   6.87ms errors=0

imgproxy latency (total time incl. body read), 200 requests, concurrency 1

       image  scenario    bytes  ttfb p50     p50     p90     p99      max   rps       p50 vs direct  errors
   large.jpg    direct    2.4MB      0.05    0.28    0.58    0.78     0.90  2954                   -       0
   large.jpg       raw    2.4MB      1.19   12.05   14.49   21.94    64.72    79    +11.77ms (42.4x)       0
   large.jpg  reencode    1.9MB     61.75   69.66   73.84   87.55   117.32    14   +69.38ms (245.3x)       0
   large.jpg    resize    3.9KB     23.89   24.00   27.46   89.91  1634.38    27    +23.72ms (84.5x)       0
   large.jpg      webp    2.0KB     27.29   27.40   29.15   31.24    57.37    36    +27.11ms (96.5x)       0
   large.jpg      avif    2.9KB     25.91   25.98   28.90   48.20    75.21    37    +25.70ms (91.5x)       0
  medium.jpg    direct  423.7KB      0.08    0.20    0.30    0.89     5.70  4150                   -       0
  medium.jpg       raw  423.7KB      0.77    2.35    4.20    4.97     5.11   371     +2.15ms (11.6x)       0
  medium.jpg  reencode  346.5KB     13.12   14.22   16.21   36.65    56.15    67    +14.02ms (70.1x)       0
  medium.jpg    resize    3.5KB      6.11    6.16    6.78    8.70    11.14   158     +5.96ms (30.3x)       0
  medium.jpg      webp    1.8KB      8.85    8.91    9.79   12.08    45.29   107     +8.71ms (43.9x)       0
  medium.jpg      avif    2.7KB      7.66    7.72    8.64   17.41    44.36   121     +7.52ms (38.0x)       0
  medium.png    direct    3.2MB      0.05    0.38    0.91    2.68     4.15  1837                   -       0
  medium.png       raw    3.2MB      1.19   15.52   17.92   22.57    26.63    63    +15.14ms (40.8x)       0
  medium.png  reencode    3.0MB    105.17  117.70  121.20  143.81   173.93     8  +117.32ms (309.7x)       0
  medium.png    resize   60.4KB     26.88   27.00   29.07   33.89    53.40    36    +26.62ms (71.0x)       0
  medium.png      webp    1.8KB     28.72   28.79   31.02   33.14    73.86    34    +28.41ms (75.8x)       0
  medium.png      avif    2.6KB     27.59   27.67   30.75   39.01   112.64    34    +27.29ms (72.8x)       0
   small.jpg    direct   65.7KB      0.12    0.16    0.41    1.20     5.67  4054                   -       0
   small.jpg       raw   65.7KB      0.55    0.62    1.46    2.38     3.63  1216      +0.47ms (4.0x)       0
   small.jpg  reencode   54.4KB      3.16    3.21    3.71    5.60     8.42   296     +3.05ms (20.3x)       0
   small.jpg    resize    6.8KB      2.72    2.77    3.74    7.36    10.79   325     +2.62ms (17.6x)       0
   small.jpg      webp    3.3KB      7.20    7.25    7.57   10.41    15.89   135     +7.09ms (45.9x)       0
   small.jpg      avif    4.7KB      4.73    4.82    5.11    6.87    10.66   201     +4.66ms (30.5x)       0

All times in milliseconds.
```

`imgproxy/run.sh -n 500 -c 8`

```go
large.jpg      direct    p50=   1.59ms p99=   4.05ms errors=0
large.jpg      raw       p50=  95.25ms p99= 141.16ms errors=0
large.jpg      reencode  p50= 235.95ms p99= 638.50ms errors=0
large.jpg      resize    p50=  38.02ms p99= 135.23ms errors=0
large.jpg      webp      p50=  43.88ms p99= 426.91ms errors=0
large.jpg      avif      p50=  41.72ms p99= 301.58ms errors=0
medium.jpg     direct    p50=   0.54ms p99=   2.23ms errors=0
medium.jpg     raw       p50=  16.55ms p99=  87.56ms errors=0
medium.jpg     reencode  p50=  37.98ms p99=  59.72ms errors=0
medium.jpg     resize    p50=  17.79ms p99=  26.53ms errors=0
medium.jpg     webp      p50=  20.60ms p99=  81.32ms errors=0
medium.jpg     avif      p50=  20.80ms p99=  39.26ms errors=0
medium.png     direct    p50=   2.14ms p99=   5.20ms errors=0
medium.png     raw       p50= 119.08ms p99= 200.08ms errors=0
medium.png     reencode  p50= 261.43ms p99= 840.78ms errors=0
medium.png     resize    p50=  49.94ms p99= 645.00ms errors=0
medium.png     webp      p50=  53.51ms p99= 109.18ms errors=0
medium.png     avif      p50=  49.46ms p99= 629.24ms errors=0
small.jpg      direct    p50=   0.34ms p99=   2.69ms errors=0
small.jpg      raw       p50=   3.57ms p99=  20.33ms errors=0
small.jpg      reencode  p50=  12.12ms p99=  21.65ms errors=0
small.jpg      resize    p50=  14.44ms p99=  21.27ms errors=0
small.jpg      webp      p50=  18.31ms p99=  27.19ms errors=0
small.jpg      avif      p50=  17.23ms p99=  23.26ms errors=0

Latency (total time incl. body read), 500 requests, concurrency 8

       image  scenario    bytes  ttfb p50     p50     p90     p99      max    rps       p50 vs direct  errors
   large.jpg    direct    2.4MB      0.40    1.59    2.74    4.05     4.65   4605                   -       0
   large.jpg       raw    2.4MB      6.22   95.25  118.70  141.16   172.30     82    +93.66ms (60.0x)       0
   large.jpg  reencode    1.9MB    167.06  235.95  326.23  638.50   995.84     32  +234.36ms (148.7x)       0
   large.jpg    resize    3.9KB     37.69   38.02   58.41  135.23  3075.50    146    +36.44ms (24.0x)       0
   large.jpg      webp    2.0KB     43.83   43.88   62.68  426.91  2084.35    134    +42.29ms (27.6x)       0
   large.jpg      avif    2.9KB     41.69   41.72   67.66  301.58  2721.99     99    +40.14ms (26.3x)       0
  medium.jpg    direct  423.7KB      0.34    0.54    1.01    2.23     3.75  12613                   -       0
  medium.jpg       raw  423.7KB      4.38   16.55   22.05   87.56    92.61    449    +16.01ms (30.8x)       0
  medium.jpg  reencode  346.5KB     32.82   37.98   46.93   59.72    76.79    207    +37.45ms (70.7x)       0
  medium.jpg    resize    3.5KB     17.77   17.79   22.27   26.53    32.43    437    +17.25ms (33.1x)       0
  medium.jpg      webp    1.8KB     20.54   20.60   26.67   81.32   218.82    359    +20.06ms (38.4x)       0
  medium.jpg      avif    2.7KB     20.76   20.80   27.42   39.26   422.87    357    +20.26ms (38.7x)       0
  medium.png    direct    3.2MB      0.59    2.14    3.48    5.20     5.82   3419                   -       0
  medium.png       raw    3.2MB      6.18  119.08  136.49  200.08   224.29     67   +116.94ms (55.6x)       0
  medium.png  reencode    3.0MB    172.93  261.43  328.86  840.78  1992.39     28  +259.29ms (122.2x)       0
  medium.png    resize   60.4KB     49.81   49.94   76.07  645.00  2489.14     83    +47.80ms (23.3x)       0
  medium.png      webp    1.8KB     53.44   53.51   70.12  109.18   889.94    134    +51.37ms (25.0x)       0
  medium.png      avif    2.6KB     49.44   49.46   71.58  629.24  1526.72    103    +47.32ms (23.1x)       0
   small.jpg    direct   65.7KB      0.28    0.34    0.99    2.69     3.35  15692                   -       0
   small.jpg       raw   65.7KB      2.97    3.57    7.92   20.33    26.34   1743     +3.24ms (10.6x)       0
   small.jpg  reencode   54.4KB     11.95   12.12   16.52   21.65    23.04    633    +11.78ms (35.9x)       0
   small.jpg    resize    6.8KB     14.39   14.44   17.64   21.27    22.58    549    +14.11ms (42.9x)       0
   small.jpg      webp    3.3KB     18.28   18.31   22.30   27.19    31.81    431    +17.98ms (54.3x)       0
   small.jpg      avif    4.7KB     17.16   17.23   20.70   23.26    32.45    462    +16.89ms (51.1x)       0

All times in milliseconds.
```

`imgproxy/run.sh -n 1000 -c 32`:

```go
large.jpg      direct    p50=   5.38ms p99=  26.92ms errors=0
large.jpg      raw       p50= 384.35ms p99= 729.05ms errors=0
large.jpg      reencode  p50=1290.14ms p99=3588.02ms errors=0
large.jpg      resize    p50= 185.44ms p99=1405.17ms errors=0
large.jpg      webp      p50= 188.48ms p99=1978.60ms errors=1
large.jpg      avif      p50= 240.08ms p99=1682.58ms errors=0
medium.jpg     direct    p50=   2.20ms p99=  22.57ms errors=0
medium.jpg     raw       p50=  93.97ms p99= 264.96ms errors=0
medium.jpg     reencode  p50= 210.52ms p99= 906.62ms errors=0
medium.jpg     resize    p50=  95.19ms p99= 233.26ms errors=0
medium.jpg     webp      p50=  93.35ms p99= 336.45ms errors=0
medium.jpg     avif      p50=  98.71ms p99= 185.47ms errors=0
medium.png     direct    p50=   8.85ms p99=  70.11ms errors=0
medium.png     raw       p50= 490.23ms p99= 776.53ms errors=0
medium.png     reencode  p50=1185.20ms p99=2494.74ms errors=0
medium.png     resize    p50= 228.92ms p99=2111.06ms errors=1
medium.png     webp      p50= 227.67ms p99=2593.01ms errors=0
medium.png     avif      p50= 257.47ms p99=2049.36ms errors=0
small.jpg      direct    p50=   1.01ms p99=  19.22ms errors=0
small.jpg      raw       p50=  14.96ms p99=  35.38ms errors=0
small.jpg      reencode  p50=  58.89ms p99= 118.01ms errors=0
small.jpg      resize    p50=  69.51ms p99= 111.71ms errors=0
small.jpg      webp      p50=  84.24ms p99= 434.02ms errors=0
small.jpg      avif      p50=  96.25ms p99= 224.33ms errors=0

imgproxy latency (total time incl. body read), 1000 requests, concurrency 32

       image  scenario    bytes  ttfb p50      p50      p90      p99      max    rps        p50 vs direct  errors
   large.jpg    direct    2.4MB      4.11     5.38    10.20    26.92    34.28   4378                    -       0
   large.jpg       raw    2.4MB      9.59   384.35   561.55   729.05  3667.65     76    +378.97ms (71.5x)       0
   large.jpg  reencode    1.9MB    667.13  1290.14  2441.70  3588.02  4330.95     22  +1284.76ms (239.9x)       0
   large.jpg    resize    3.9KB    185.39   185.44   240.25  1405.17  2839.85    134    +180.06ms (34.5x)       0
   large.jpg      webp    2.0KB    188.43   188.48   264.99  1978.60  4258.44    104    +183.10ms (35.0x)       1
   large.jpg      avif    2.9KB    240.03   240.08   404.90  1682.58  3427.57     90    +234.71ms (44.6x)       0
  medium.jpg    direct  423.7KB      1.82     2.20     5.65    22.57    28.79   9000                    -       0
  medium.jpg       raw  423.7KB     10.83    93.97   185.64   264.96   475.44    291     +91.78ms (42.8x)       0
  medium.jpg  reencode  346.5KB    158.82   210.52   424.86   906.62  1251.74    130    +208.33ms (95.8x)       0
  medium.jpg    resize    3.5KB     95.16    95.19   165.58   233.26   504.86    296     +93.00ms (43.3x)       0
  medium.jpg      webp    1.8KB     93.28    93.35   148.26   336.45  1129.11    295     +91.15ms (42.5x)       0
  medium.jpg      avif    2.7KB     98.68    98.71   129.14   185.47   912.90    300     +96.51ms (44.9x)       0
  medium.png    direct    3.2MB      5.84     8.85    25.25    70.11    83.25   2204                    -       0
  medium.png       raw    3.2MB      9.61   490.23   582.74   776.53  2227.72     63    +481.38ms (55.4x)       0
  medium.png  reencode    3.0MB    509.60  1185.20  1816.10  2494.74  5174.65     25  +1176.36ms (134.0x)       0
  medium.png    resize   60.4KB    227.01   228.92   361.98  2111.06  4446.74     92    +220.08ms (25.9x)       1
  medium.png      webp    1.8KB    227.62   227.67   360.32  2593.01  4927.95     89    +218.83ms (25.7x)       0
  medium.png      avif    2.6KB    257.41   257.47   602.44  2049.36  4103.31     86    +248.63ms (29.1x)       0
   small.jpg    direct   65.7KB      0.93     1.01     6.68    19.22    27.22  10925                    -       0
   small.jpg       raw   65.7KB      7.98    14.96    20.80    35.38    38.92   2014     +13.94ms (14.8x)       0
   small.jpg  reencode   54.4KB     58.59    58.89    84.38   118.01   128.37    507     +57.88ms (58.2x)       0
   small.jpg    resize    6.8KB     69.35    69.51    88.77   111.71   125.14    444     +68.50ms (68.7x)       0
   small.jpg      webp    3.3KB     84.14    84.24   106.19   434.02   476.89    328     +83.23ms (83.2x)       0
   small.jpg      avif    4.7KB     96.13    96.25   148.01   224.33   262.40    303     +95.24ms (95.1x)       0

All times in milliseconds.
```

`ipx/run.sh`:

```go
skipping scenario "raw": not supported by ipx
large.jpg      direct    p50=   0.28ms p99=   0.80ms errors=0
large.jpg      reencode  p50= 159.34ms p99= 275.49ms errors=0
large.jpg      resize    p50=  27.96ms p99=  42.94ms errors=0
large.jpg      webp      p50=  31.26ms p99=  62.48ms errors=0
large.jpg      avif      p50=  71.77ms p99= 123.31ms errors=0
medium.jpg     direct    p50=   0.20ms p99=   0.40ms errors=0
medium.jpg     reencode  p50=  30.36ms p99=  38.84ms errors=0
medium.jpg     resize    p50=   8.33ms p99=  16.35ms errors=0
medium.jpg     webp      p50=  10.59ms p99=  14.50ms errors=0
medium.jpg     avif      p50=  45.98ms p99= 113.54ms errors=0
medium.png     direct    p50=   0.38ms p99=   3.76ms errors=0
medium.png     reencode  p50= 102.39ms p99= 172.99ms errors=0
medium.png     resize    p50=  30.50ms p99=  62.62ms errors=0
medium.png     webp      p50=  32.22ms p99=  37.09ms errors=0
medium.png     avif      p50=  68.18ms p99= 121.78ms errors=0
small.jpg      direct    p50=   0.16ms p99=   1.56ms errors=0
small.jpg      reencode  p50=   7.09ms p99=   9.69ms errors=0
small.jpg      resize    p50=   4.96ms p99=  10.89ms errors=0
small.jpg      webp      p50=   8.88ms p99=  14.14ms errors=0
small.jpg      avif      p50=  59.52ms p99=  92.74ms errors=0

ipx latency (total time incl. body read), 200 requests, concurrency 1

       image  scenario    bytes  ttfb p50     p50     p90     p99       max   rps       p50 vs direct  errors
   large.jpg    direct    2.4MB      0.04    0.28    0.60    0.80      0.92  3029                   -       0
   large.jpg  reencode    1.9MB    150.87  159.34  178.61  275.49    420.14     6  +159.06ms (571.1x)       0
   large.jpg    resize    4.0KB     27.87   27.96   31.75   42.94     77.72    34   +27.68ms (100.2x)       0
   large.jpg      webp    1.8KB     31.16   31.26   34.85   62.48   1704.31    25   +30.98ms (112.0x)       0
   large.jpg      avif    1.4KB     71.68   71.77   79.69  123.31   3384.89    10   +71.49ms (257.2x)       0
  medium.jpg    direct  423.7KB      0.08    0.20    0.28    0.40      0.54  4834                   -       0
  medium.jpg  reencode  340.2KB     29.45   30.36   32.34   38.84  10245.42    12   +30.16ms (154.1x)       0
  medium.jpg    resize    3.6KB      8.27    8.33   10.51   16.35     75.00   108     +8.14ms (42.3x)       0
  medium.jpg      webp    1.6KB     10.53   10.59   12.08   14.50     40.81    91    +10.39ms (53.7x)       0
  medium.jpg      avif    1.3KB     45.87   45.98   49.05  113.54  10510.55     7   +45.78ms (233.4x)       0
  medium.png    direct    3.2MB      0.05    0.38    1.11    3.76      6.21  1436                   -       0
  medium.png  reencode    5.5MB     82.07  102.39  105.94  172.99  10761.37     6  +102.02ms (272.3x)       0
  medium.png    resize  140.3KB     30.19   30.50   33.74   62.62  10580.91    12    +30.12ms (81.1x)       0
  medium.png      webp    1.6KB     32.15   32.22   34.42   37.09     43.10    31    +31.84ms (85.7x)       0
  medium.png      avif    1.3KB     68.09   68.18   71.44  121.78    148.37    14   +67.80ms (181.3x)       0
   small.jpg    direct   65.7KB      0.11    0.16    0.29    1.56      2.98  4334                   -       0
   small.jpg  reencode   53.3KB      7.02    7.09    7.95    9.69     12.74   137     +6.94ms (45.2x)       0
   small.jpg    resize    6.7KB      4.89    4.96    5.88   10.89     27.59   186     +4.81ms (31.6x)       0
   small.jpg      webp    3.5KB      8.81    8.88   11.27   14.14     15.71   107     +8.73ms (56.6x)       0
   small.jpg      avif    1.6KB     59.44   59.52   61.26   92.74  10230.21     9   +59.36ms (379.1x)       0

All times in milliseconds.
```


`ipx/run.sh -n 500 -c 8`:

```go
skipping scenario "raw": not supported by ipx
large.jpg      direct    p50=   1.49ms p99=   3.87ms errors=0
large.jpg      reencode  p50= 340.15ms p99=2348.27ms errors=0
large.jpg      resize    p50=  35.43ms p99=10194.29ms errors=1
large.jpg      webp      p50=  40.30ms p99=10534.52ms errors=5
large.jpg      avif      p50= 111.68ms p99=1142.66ms errors=0
medium.jpg     direct    p50=   0.58ms p99=  39.02ms errors=0
medium.jpg     reencode  p50=  68.35ms p99=10439.15ms errors=7
medium.jpg     resize    p50=  18.25ms p99=10294.66ms errors=2
medium.jpg     webp      p50=  21.31ms p99=10541.49ms errors=5
medium.jpg     avif      p50=  48.67ms p99= 144.24ms errors=3
medium.png     direct    p50=   2.17ms p99=  10.19ms errors=0
medium.png     reencode  p50= 295.46ms p99=10399.11ms errors=1
medium.png     resize    p50=  40.63ms p99=10504.01ms errors=4
medium.png     webp      p50=  42.38ms p99=10419.21ms errors=2
medium.png     avif      p50= 120.87ms p99=10593.57ms errors=3
small.jpg      direct    p50=   0.16ms p99=   0.90ms errors=0
small.jpg      reencode  p50=  15.78ms p99=10235.93ms errors=2
small.jpg      resize    p50=  10.01ms p99=10275.01ms errors=3
small.jpg      webp      p50=  19.68ms p99=  45.06ms errors=0
small.jpg      avif      p50= 119.27ms p99=10485.37ms errors=10

ipx latency (total time incl. body read), 500 requests, concurrency 8

       image  scenario    bytes  ttfb p50     p50     p90       p99       max    rps       p50 vs direct  errors
   large.jpg    direct    2.4MB      0.40    1.49    2.43      3.87      5.47   5001                   -       0
   large.jpg  reencode    1.9MB    295.87  340.15  457.96   2348.27   4388.30     21  +338.66ms (228.6x)       0
   large.jpg    resize    4.0KB     35.36   35.43   51.30  10194.29  10488.64     25    +33.94ms (23.8x)       1
   large.jpg      webp    1.8KB     40.26   40.30   67.95  10534.52  10638.06     13    +38.82ms (27.1x)       5
   large.jpg      avif    1.4KB    111.65  111.68  153.15   1142.66   4148.94     50   +110.19ms (75.1x)       0
  medium.jpg    direct  423.7KB      0.37    0.58    1.53     39.02     40.05   5351                   -       0
  medium.jpg  reencode  340.2KB     64.68   68.35   93.70  10439.15  10549.86     20   +67.77ms (117.4x)       7
  medium.jpg    resize    3.6KB     18.22   18.25   25.66  10294.66  10332.80     43    +17.67ms (31.4x)       2
  medium.jpg      webp    1.6KB     21.24   21.31   29.33  10541.49  10594.20     22    +20.72ms (36.6x)       5
  medium.jpg      avif    1.3KB     48.62   48.67   90.87    144.24  10555.31     44    +48.09ms (83.6x)       3
  medium.png    direct    3.2MB      0.64    2.17    4.18     10.19     14.70   3027                   -       0
  medium.png  reencode    5.5MB    126.57  295.46  411.59  10399.11  10730.14     16  +293.28ms (135.8x)       1
  medium.png    resize  140.3KB     39.84   40.63   88.84  10504.01  10656.78     18    +38.46ms (18.7x)       4
  medium.png      webp    1.6KB     42.35   42.38   71.08  10419.21  10536.59     20    +40.21ms (19.5x)       2
  medium.png      avif    1.3KB    120.74  120.87  212.72  10593.57  10693.42     12   +118.70ms (55.6x)       3
   small.jpg    direct   65.7KB      0.13    0.16    0.35      0.90      1.41  36959                   -       0
   small.jpg  reencode   53.3KB     15.74   15.78   22.62  10235.93  10586.85     44    +15.62ms (99.2x)       2
   small.jpg    resize    6.7KB      9.99   10.01   18.67  10275.01  10604.39     28     +9.85ms (62.9x)       3
   small.jpg      webp    3.5KB     19.66   19.68   29.18     45.06     54.38    379   +19.52ms (123.8x)       0
   small.jpg      avif    1.6KB    119.23  119.27  172.59  10485.37  10684.67     13  +119.11ms (750.1x)      10

All times in milliseconds.
```
