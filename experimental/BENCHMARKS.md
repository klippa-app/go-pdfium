# WebAssembly runtime benchmarks

go-pdfium can run the PDFium WebAssembly module on three runtimes: [wazero](https://wazero.io/) in the `webassembly`
package, and [wazy](https://github.com/samyfodil/wazy) and [wago](https://github.com/wago-org/wago) in the
experimental `experimental/wazy` and `experimental/wago` packages. This document compares them.

Everything here goes through the public `pdfium.Pool` API of each backend, so the numbers include go-pdfium's own
overhead in the same way for all of them. All renders were compared pixel for pixel: in the default configuration the
three runtimes produced identical output for every one of the 5,000 corpus documents.

## Setup

| | |
|---|---|
| go-pdfium | this commit, module `internal/pdfium_wasm/pdfium.wasm` (PDFium 8057, 5,754,045 bytes) |
| wazero | v1.12.0, compiler engine, `CoreFeaturesV2 \| CoreFeaturesExceptionHandling` |
| wazy | main `ed4381b8` (2026-09-20, `v0.3.1-0.20260920105944-ed4381b86521`), compiler engine, `CoreFeaturesV2 \| CoreFeatureExceptionHandling` |
| wago | main `fa26daaf` (2026-09-23, `v0.1.0-beta.9.0.20260923033105-fa26daaf4a4b`), default config (explicit bounds checks), WASI plugin v0.3.1 |
| Go | go1.27.1 |
| Machine | Apple M5 (10 cores), macOS 26 (Darwin 25.6.0), arm64, otherwise idle |

Micro-benchmark numbers are the median of five `go test -bench` runs (`-count 5`). The corpus numbers come from a
single interleaved pass over all documents. All numbers in this document were measured on arm64; the
`.github/workflows/benchmark.yml` workflow runs the micro-benchmarks on GitHub's amd64 and arm64 runners on demand.

## Micro-benchmarks

The benchmarks live in `experimental/benchmark`:

```sh
go test -run '^$' -bench . -benchmem -count 5 ./experimental/benchmark
```

- **Init**: create a pool, which compiles the 5.7 MB module to machine code and starts one worker.
- **NewInstance**: instantiate the compiled module, run its constructors and `FPDF_InitLibrary`. This is what every
  `GetInstance` costs when workers are not reused, the default.
- **OpenDocument**: load `test.pdf` from memory and close it again.
- **Render/...**: render the first page at 100 DPI of a small text page (`test.pdf`), a page with an ICC based colour
  space and transparency (`alpha_channel.pdf`), a page with embedded images and a large PNG predicted flate page
  (`rect-wrong.pdf`).
- **RenderToJPEG**: render `embedded_images.pdf` at 150 DPI and encode it to JPEG inside the module.
- **GetPageText**: structured text extraction of `test.pdf`, many small calls into the module.

### Native arm64 (Apple M5)

| Benchmark | wazero | wazy | wago | wazy vs wazero | wago vs wazero |
|---|---|---|---|---|---|
| Init | 1.18 s | 726.58 ms | 244.97 ms | 0.62x | 0.21x |
| NewInstance | 438.4 µs | 231.3 µs | 308.4 µs | 0.53x | 0.70x |
| OpenDocument | 9.0 µs | 6.8 µs | 8.0 µs | 0.76x | 0.89x |
| Render/test | 194.8 µs | 113.0 µs | 153.2 µs | 0.58x | 0.79x |
| Render/alpha_channel | 27.30 ms | 24.37 ms | 35.98 ms | 0.89x | 1.32x |
| Render/embedded_images | 1.05 ms | 818.9 µs | 1.18 ms | 0.78x | 1.12x |
| Render/rect-wrong | 7.57 ms | 7.05 ms | 9.60 ms | 0.93x | 1.27x |
| RenderToJPEG | 13.23 ms | 7.50 ms | 8.58 ms | 0.57x | 0.65x |
| GetPageText | 29.0 µs | 26.4 µs | 76.2 µs | 0.91x | 2.63x |

Allocations per operation (median B/op, allocs/op):

| Benchmark | wazero | wazy | wago |
|---|---|---|---|
| Init | 345,871,832 B, 659,031 | 98,872,424 B, 85,318 | 57,231,697 B, 60,028 |
| NewInstance | 19,305,130 B, 13,830 | 207,568 B, 105 | 114,331 B, 1,208 |
| OpenDocument | 1,008 B, 22 | 1,008 B, 22 | 1,008 B, 22 |
| Render/test | 518 B, 7 | 497 B, 7 | 489 B, 7 |
| Render/alpha_channel | 5,103,415 B, 62 | 2,464,541 B, 8 | 853 B, 13 |
| Render/embedded_images | 2,513,890 B, 238 | 556 B, 7 | 502 B, 7 |
| Render/rect-wrong | 1,326 B, 7 | 1,053 B, 7 | 627 B, 9 |
| RenderToJPEG | 2,608,276 B, 244 | 83,503 B, 12 | 82,721 B, 13 |
| GetPageText | 12,396 B, 197 | 12,395 B, 197 | 12,392 B, 197 |

wazy is the fastest of the three on every execution benchmark: it compiles the module in two thirds of wazero's
time, creates a worker in half the time with a fraction of the allocations, and executes 7 to 43 percent faster than
wazero depending on how much of the time is runtime overhead rather than PDFium compute. wago has a different profile.
Compiling the module takes a quarter of a second, almost five times faster than wazero with a sixth of the
allocations, and creating a worker is faster than wazero as well. Its execution is faster than wazero on the short
text page render and on the JPEG encode, but 12 to 32 percent slower on the compute bound renders and much slower on
text extraction, a workload of many small calls into the module where its per call overhead shows.

### Compilation workers

`Init` above uses each runtime's default compilation mode. wago compiles functions in parallel by default; wazero and
wazy compile with a single goroutine unless asked otherwise, through a context value: `experimental.WithCompilationWorkers`
for wazero and `api.WithCompilationWorkers` for wazy. go-pdfium passes `Config.Context` to the compiler, so the option
can be set there. Time to create a pool (compile plus first worker, best of 3):

| Compilation workers | wazero | wazy |
|---|---|---|
| 1 (default) | 1.16 s | 0.74 s |
| 2 | 0.57 s | 0.40 s |
| 4 | 0.38 s | 0.25 s |
| 8 | 0.35 s | 0.21 s |
| 10 | 0.33 s | 0.19 s |

With four or more workers both come within reach of wago's parallel compiler, and wazy with eight workers is
slightly faster than it. Nothing else in this document depends on this setting: it only affects `Init`.

## Real world corpus

`experimental/benchmark/cmd/corpus` renders the first page of every PDF in a directory at 100 DPI with each runtime,
interleaved per document (with a rotating order so that no runtime always runs first on a cold cache), records the
time per document and a hash of the rendered pixels, and reports totals, percentiles, the per-document geometric mean
ratio and the number of pixel mismatches:

```sh
go run ./experimental/benchmark/cmd/corpus -dir /path/to/pdfs [-runtimes wazero,wazy,wago] [-dpi 100] [-stride N] [-csv out.csv]
```

The corpus used here is a private set of 5,000 real world PDFs (859 MB, invoices, scans, forms, reports, CAD
drawings), rendered natively on the Apple M5. Per document time includes opening the document from memory, rendering
and closing it. All three runtimes rendered all 5,000 documents and produced identical pixels for every one of them.

| Runtime | Compile + first worker | Rendered | Failed | Timed out | Total render time | Mean | Median | p90 | p99 | Geomean ratio vs wazero | Pixel mismatches vs wazero |
| wazero | 1.165s | 5000 | 0 | 0 | 76.4 s | 15.27 ms | 7.66 ms | 32.8 ms | 89 ms | 1.00x | reference |
| wazy | 731ms | 5000 | 0 | 0 | 67.2 s | 13.43 ms | 6.54 ms | 28.2 ms | 84 ms | 0.861x | 0 of 5000 |
| wago | 246ms | 5000 | 0 | 0 | 82.9 s | 16.58 ms | 7.38 ms | 33.5 ms | 121 ms | 0.986x | 0 of 5000 |

Per document, wazy is faster than wazero on 4,766 of the 5,000 documents (wazy/wazero time ratio: 0.80 at the 10th
percentile, 0.86 at the median, 0.94 at the 90th), wago on 3,276 of them (0.87, 0.95 and 1.17). wago
matches wazero overall and is a little faster on the median document, but has a heavier tail: its 99th percentile is
121 ms against wazero's 89 ms, so the heaviest documents take longer on wago.

### The cost of close-on-context-done

All three runtimes can interrupt a running call when a context is cancelled (`WithCloseOnContextDone` on the wazero
and wazy runtime configuration, `InterruptibleCalls` on wago), which is what makes `Kill()` work on a document that
never finishes. It is off by default in go-pdfium. Measured on every fifth document of the corpus (1,000 files),
first with the option off:

| Runtime | Compile + first worker | Rendered | Failed | Timed out | Total render time | Mean | Median | p90 | p99 | Geomean ratio vs wazero | Pixel mismatches vs wazero |
| wazero | 1.162s | 1000 | 0 | 0 | 14.7 s | 14.66 ms | 6.95 ms | 31.6 ms | 84 ms | 1.00x | reference |
| wazy | 753ms | 1000 | 0 | 0 | 13.1 s | 13.09 ms | 6.02 ms | 27.4 ms | 83 ms | 0.870x | 0 of 1000 |
| wago | 260ms | 1000 | 0 | 0 | 16.2 s | 16.23 ms | 6.71 ms | 32.4 ms | 139 ms | 0.998x | 2 of 1000 |

and then with it on:

| Runtime | Compile + first worker | Rendered | Failed | Timed out | Total render time | Mean | Median | p90 | p99 | Geomean ratio vs wazero | Pixel mismatches vs wazero |
| wazero | 1.174s | 1000 | 0 | 0 | 67.1 s | 67.06 ms | 25.87 ms | 182.8 ms | 579 ms | 1.00x | reference |
| wazy | 796ms | 1000 | 0 | 0 | 16.0 s | 16.00 ms | 7.15 ms | 34.4 ms | 100 ms | 0.266x | 0 of 1000 |
| wago | 249ms | 1000 | 0 | 0 | 16.9 s | 16.88 ms | 6.86 ms | 32.5 ms | 149 ms | 0.261x | 2 of 1000 |

wazero becomes 4.6 times slower with the option on, wazy 1.2 times and wago 1.05 times, so in this mode wazy and wago
are roughly equal and four times faster than wazero. For wazero this is addressed by an upstream pull request, see
the wazero versions section below.

The 2 wago pixel mismatches on this 1,000 document sequence occur in both modes (and not on the 5,000 document
sequence above). Their cause is a wago arm64 bug: `memory.fill` clobbers a register that is still in use, which
PDFium hits in `std::fill_n` on a `std::vector<bool>`, corrupting memory in a heap layout dependent way. It has been
reported to wago with a small reproducer; linux/amd64 is not affected.

## wazero versions

The same benchmarks were run against unreleased wazero builds, to see what the upcoming changes bring for go-pdfium.
Each variant was substituted for wazero v1.12.0 through a Go workspace `replace` directive, with everything else
identical, and all five variants pass go-pdfium's complete wazero test suite. wazy (v0.3.0 at the time) was included in every run as a
control for machine drift: its numbers were flat across the runs except for the `Render/test` micro-benchmark, where
the control moved by about 20% between runs, so treat that one row as noisier than the rest.

| Variant | Commit | What it is |
|---|---|---|
| v1.12.0 | release | the version go-pdfium uses |
| main | `451613ca` (2026-09-08) | wazero `main` |
| main + 2533 | `451613ca` + [#2533](https://github.com/wazero/wazero/pull/2533) | function entry as a termination checkpoint |
| main + 2529 + 2530 | `451613ca` + [#2529](https://github.com/wazero/wazero/pull/2529) + [#2530](https://github.com/wazero/wazero/pull/2530) | wazevo group-ID fix; clone-free `try_table` checkpoints via top-relative offsets and a trampoline |
| main + 2529 + 2530 + 2533 | all three merged (one trivial conflict, both PRs add fields to the same struct) | the combination |

### Micro-benchmarks (native arm64, median of 5)

| Benchmark | v1.12.0 | main | main+2533 | main+2529+2530 | main+all three | main vs v1.12.0 | main+2533 vs v1.12.0 | main+2529+2530 vs v1.12.0 | main+all three vs v1.12.0 |
|---|---|---|---|---|---|---|---|---|---|
| Init | 1.14 s | 1.13 s | 1.14 s | 1.15 s | 1.15 s | 0.99x | 1.01x | 1.02x | 1.01x |
| NewInstance | 437.6 µs | 474.7 µs | 459.2 µs | 456.2 µs | 451.8 µs | 1.08x | 1.05x | 1.04x | 1.03x |
| OpenDocument | 8.8 µs | 8.6 µs | 8.6 µs | 8.5 µs | 8.3 µs | 0.97x | 0.97x | 0.97x | 0.94x |
| Render/test | 193.4 µs | 139.2 µs | 140.2 µs | 134.3 µs | 147.4 µs | 0.72x | 0.72x | 0.69x | 0.76x |
| Render/alpha_channel | 26.97 ms | 25.58 ms | 25.54 ms | 24.67 ms | 25.79 ms | 0.95x | 0.95x | 0.91x | 0.96x |
| Render/embedded_images | 1.05 ms | 973.1 µs | 966.4 µs | 905.0 µs | 923.9 µs | 0.93x | 0.92x | 0.86x | 0.88x |
| Render/rect-wrong | 7.04 ms | 6.87 ms | 6.89 ms | 6.80 ms | 6.65 ms | 0.98x | 0.98x | 0.97x | 0.94x |
| RenderToJPEG | 11.20 ms | 6.87 ms | 6.87 ms | 6.72 ms | 6.82 ms | 0.61x | 0.61x | 0.60x | 0.61x |
| GetPageText | 29.0 µs | 26.8 µs | 27.2 µs | 26.6 µs | 26.9 µs | 0.93x | 0.94x | 0.92x | 0.93x |

### Corpus, default configuration (5,000 documents)

| wazero | Total render time | Mean | Median | p90 | p99 | vs v1.12.0 |
|---|---|---|---|---|---|---|
| v1.12.0 | 69.4 s | 13.88 ms | 7.09 ms | 30.7 ms | 83 ms | 1.00x |
| main | 66.6 s | 13.32 ms | 6.79 ms | 29.2 ms | 80 ms | 0.96x |
| main + 2533 | 66.8 s | 13.37 ms | 6.80 ms | 29.4 ms | 81 ms | 0.96x |
| main + 2529 + 2530 | 63.9 s | 12.79 ms | 6.66 ms | 28.0 ms | 77 ms | 0.92x |
| main + 2529 + 2530 + 2533 | 64.6 s | 12.93 ms | 6.71 ms | 28.3 ms | 79 ms | 0.93x |

### Corpus with close-on-context-done (every fifth document, 1,000 files)

| wazero | Total render time | Mean | Median | p90 | p99 | vs v1.12.0 |
|---|---|---|---|---|---|---|
| v1.12.0 | 64.8 s | 64.79 ms | 25.36 ms | 175.1 ms | 426 ms | 1.00x |
| main | 63.6 s | 63.56 ms | 24.79 ms | 174.2 ms | 392 ms | 0.98x |
| main + 2533 | 14.8 s | 14.82 ms | 6.97 ms | 31.6 ms | 85 ms | 0.23x |
| main + 2529 + 2530 | 62.9 s | 62.94 ms | 24.49 ms | 172.4 ms | 379 ms | 0.97x |
| main + 2529 + 2530 + 2533 | 14.2 s | 14.22 ms | 6.79 ms | 30.2 ms | 81 ms | 0.22x |

### Reading the numbers

- **main** is a few percent faster than v1.12.0 across the board (4% on the corpus) and a lot faster on the JPEG
  encode path (0.61x), which runs libjpeg-turbo's SIMD kernels inside the module.
- **#2533** does not change the default path at all, but it removes the close-on-context-done penalty almost
  completely: the interruptible corpus run goes from 63.6 s to 14.8 s, within 3% of the 14.4 s that v1.12.0 needs
  with the option off. With this PR, enabling `WithCloseOnContextDone` so that `Kill()` can interrupt a stuck render
  becomes essentially free on wazero, and wazero is faster than wazy in that mode (wazy: 17.6 s on the same files).
- **#2529 + #2530** improve default execution a little further (another 4% on the corpus, 0.86x on the image page,
  0.60x on JPEG encode) but do not touch the close-on-context-done cost.
- **All three together** behave as the sum of the parts: the default corpus matches #2529 + #2530 (64.6 s, 0.93x),
  and the interruptible corpus matches #2533 (14.2 s, 0.22x). Nothing regresses when they are combined, so this is
  the configuration to look forward to: a few percent faster than today, and `WithCloseOnContextDone` at no cost.

### Compared with wazy

wazy was the control in every run, so the combined variant can be compared with it directly, from the same run:

| Benchmark | wazero main + 2529 + 2530 + 2533 | wazy v0.3.0 | wazy vs wazero |
|---|---|---|---|
| Init | 1.15 s | 714 ms | 0.62x |
| NewInstance | 451.8 µs | 230.7 µs | 0.51x |
| OpenDocument | 8.3 µs | 6.7 µs | 0.81x |
| Render/test | 147.4 µs | 129.8 µs | 0.88x |
| Render/alpha_channel | 25.79 ms | 24.13 ms | 0.94x |
| Render/embedded_images | 923.9 µs | 811.9 µs | 0.88x |
| Render/rect-wrong | 6.65 ms | 6.28 ms | 0.95x |
| RenderToJPEG | 6.82 ms | 7.43 ms | 1.09x |
| GetPageText | 26.9 µs | 24.3 µs | 0.90x |

| Corpus | wazero main + 2529 + 2530 + 2533 | wazy v0.3.0 | wazy vs wazero |
|---|---|---|---|
| default configuration, 5,000 documents | 64.6 s | 61.2 s | 0.92x |
| close-on-context-done, 1,000 documents | 14.2 s | 17.6 s | 1.16x |

Against v1.12.0 wazy is 14% faster on the corpus; against wazero with the three PRs it is 8% faster in the default
configuration and 16% slower once close-on-context-done is enabled, and wazero wins the in-module JPEG encode. What
wazy keeps is the startup side: it compiles the module in 0.62x the time and creates a worker in half the time with a
small fraction of the allocations. In execution the two runtimes end up close to even.

## All three runtimes against wazero with the pull requests

The comparison above uses the released wazero v1.12.0. This section repeats it with wazero `main` @ `c267b507`
(2026-09-22) plus [#2529](https://github.com/wazero/wazero/pull/2529), [#2530](https://github.com/wazero/wazero/pull/2530)
and [#2533](https://github.com/wazero/wazero/pull/2533) merged (one trivial conflict, both sides kept), substituted
through a Go workspace `replace`; wazy and wago are the same commits as above. That build passes go-pdfium's complete
wazero test suite. The runs happened right after the ones above, on the same machine.

### Micro-benchmarks (native arm64, median of 5)

| Benchmark | wazero | wazy | wago | wazy vs wazero | wago vs wazero |
|---|---|---|---|---|---|
| Init | 1.15 s | 719.02 ms | 252.08 ms | 0.63x | 0.22x |
| NewInstance | 459.4 µs | 235.8 µs | 335.1 µs | 0.51x | 0.73x |
| OpenDocument | 8.5 µs | 6.8 µs | 7.9 µs | 0.80x | 0.93x |
| Render/test | 126.8 µs | 119.3 µs | 147.0 µs | 0.94x | 1.16x |
| Render/alpha_channel | 25.14 ms | 24.43 ms | 34.14 ms | 0.97x | 1.36x |
| Render/embedded_images | 902.1 µs | 827.4 µs | 940.4 µs | 0.92x | 1.04x |
| Render/rect-wrong | 6.73 ms | 6.59 ms | 9.74 ms | 0.98x | 1.45x |
| RenderToJPEG | 6.73 ms | 7.50 ms | 8.55 ms | 1.12x | 1.27x |
| GetPageText | 26.4 µs | 24.5 µs | 77.5 µs | 0.93x | 2.93x |

### Corpus, default configuration (5,000 documents)

| Runtime | Compile + first worker | Rendered | Failed | Timed out | Total render time | Mean | Median | p90 | p99 | Geomean ratio vs wazero | Pixel mismatches vs wazero |
| wazero | 1.218s | 5000 | 0 | 0 | 66.5 s | 13.31 ms | 6.84 ms | 29.1 ms | 79 ms | 1.00x | reference |
| wazy | 732ms | 5000 | 0 | 0 | 62.8 s | 12.57 ms | 6.21 ms | 26.3 ms | 82 ms | 0.924x | 0 of 5000 |
| wago | 286ms | 5000 | 0 | 0 | 78.9 s | 15.79 ms | 7.00 ms | 32.1 ms | 118 ms | 1.066x | 0 of 5000 |

Per document, wazy is faster than this wazero on 4,307 of the 5,000 documents (ratio 0.87 / 0.91 / 1.04 at the 10th /
50th / 90th percentile), wago on 2,160 of them (0.93 / 1.02 / 1.29).

### Every fifth document (1,000 files), close-on-context-done off and on

Off:

| Runtime | Compile + first worker | Rendered | Failed | Timed out | Total render time | Mean | Median | p90 | p99 | Geomean ratio vs wazero | Pixel mismatches vs wazero |
| wazero | 1.154s | 1000 | 0 | 0 | 13.7 s | 13.70 ms | 6.59 ms | 29.2 ms | 76 ms | 1.00x | reference |
| wazy | 725ms | 1000 | 0 | 0 | 13.0 s | 13.03 ms | 6.01 ms | 27.1 ms | 80 ms | 0.923x | 0 of 1000 |
| wago | 282ms | 1000 | 0 | 0 | 16.2 s | 16.15 ms | 6.58 ms | 31.7 ms | 136 ms | 1.057x | 2 of 1000 |

On:

| Runtime | Compile + first worker | Rendered | Failed | Timed out | Total render time | Mean | Median | p90 | p99 | Geomean ratio vs wazero | Pixel mismatches vs wazero |
| wazero | 1.281s | 1000 | 0 | 0 | 14.6 s | 14.62 ms | 6.92 ms | 31.1 ms | 88 ms | 1.00x | reference |
| wazy | 803ms | 1000 | 0 | 0 | 15.7 s | 15.68 ms | 7.06 ms | 34.0 ms | 115 ms | 1.036x | 0 of 1000 |
| wago | 282ms | 1000 | 0 | 0 | 16.2 s | 16.18 ms | 6.70 ms | 32.4 ms | 139 ms | 1.004x | 2 of 1000 |

With the three pull requests wazero closes most of the execution gap. On the corpus wazy's advantage shrinks from
0.86x to 0.92x and wago goes from matching wazero to 1.07x; per render the three are within a few percent of each
other on the compute bound pages, and wazero now wins the in-module JPEG encode (wazy 1.12x, wago 1.27x). The
startup side does not change: wazy still compiles in two thirds of the time and creates a worker in half the time,
and wago still compiles the module several times faster than either with its default parallel compiler (see the
compilation workers table above for what wazero and wazy do with more than one compilation goroutine).

The biggest change is close-on-context-done. With the pull requests it costs wazero 7% on this subset (13.7 s to
14.6 s) instead of 4.6x, so wazero becomes the fastest runtime in that mode (wazy 1.04x, wago 1.00x), and enabling it
to make `Kill()` interrupt stuck renders no longer needs to be weighed against throughput. wago's 2 pixel mismatches
on this sequence are the arm64 issue described above and are unrelated to the wazero build.
