# WebAssembly runtime benchmarks

go-pdfium can run the PDFium WebAssembly module on three runtimes: [wazero](https://wazero.io/) in the `webassembly`
package, and [wazy](https://github.com/samyfodil/wazy) and [wago](https://github.com/wago-org/wago) in the
experimental `experimental/wazy` and `experimental/wago` packages. This document compares them.

Everything here goes through the public `pdfium.Pool` API of each backend, so the numbers include go-pdfium's own
overhead in the same way for all of them. All renders were compared pixel for pixel: all three runtimes produced
identical output for every one of the 5,000 corpus documents, in both the default and the close-on-context-done
configuration.

## Setup

| | |
|---|---|
| go-pdfium | this commit, module `internal/pdfium_wasm/pdfium.wasm` (PDFium 8111, 5,828,543 bytes) |
| wazero | v1.12.0 (latest release), compiler engine, `CoreFeaturesV2 \| CoreFeaturesExceptionHandling` |
| wazy | `v0.5.0` (2026-10-24), compiler engine, `CoreFeaturesV2 \| CoreFeatureExceptionHandling` |
| wago | `v0.1.0-beta.13.0.20261024181356-4de27acba129` (main, 2026-10-24), default config (explicit bounds checks), WASI plugin v0.3.1 |
| Go | go1.26.4 |
| Machine | Apple M5 (10 cores), macOS 27.0.1 (Darwin 27.0.0), arm64, otherwise idle |

wago is pinned to a commit rather than a tag: it could not compile this module at all on arm64 as of `v0.1.0-beta.13`
(`arm64: native branch displacement exceeds instruction range`, a codegen limit the module outgrew), fixed upstream in
`fix(arm64): compile large functions with long conditional branches` (wago-org/wago#1014). That fix, plus several
other arm64 fixes landed the same day, is not in any tagged release yet.

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
| Init | 1.205 s | 1.184 s | 246.4 ms | 0.98x | 0.20x |
| NewInstance | 456.7 µs | 251.9 µs | 370.6 µs | 0.55x | 0.81x |
| OpenDocument | 9.4 µs | 8.5 µs | 8.4 µs | 0.91x | 0.90x |
| Render/test | 203.6 µs | 122.2 µs | 147.6 µs | 0.60x | 0.72x |
| Render/alpha_channel | 27.68 ms | 47.63 ms | 38.35 ms | 1.72x | 1.39x |
| Render/embedded_images | 1.110 ms | 1.344 ms | 977.4 µs | 1.21x | 0.88x |
| Render/rect-wrong | 7.157 ms | 10.47 ms | 9.870 ms | 1.46x | 1.38x |
| RenderToJPEG | 7.456 ms | 10.14 ms | 8.839 ms | 1.36x | 1.19x |
| GetPageText | 30.1 µs | 34.4 µs | 94.8 µs | 1.14x | 3.15x |

Allocations per operation (median B/op, allocs/op):

| Benchmark | wazero | wazy | wago |
|---|---|---|---|
| Init | 377,229,224 B, 721,264 | 1,133,704,264 B, 385,482 | 55,947,611 B, 20,705 |
| NewInstance | 19,372,043 B, 13,892 | 207,637 B, 107 | 117,689 B, 1,252 |
| OpenDocument | 1,008 B, 22 | 1,008 B, 22 | 1,008 B, 22 |
| Render/test | 519 B, 7 | 526 B, 8 | 489 B, 7 |
| Render/alpha_channel | 6,139,815 B, 66 | 4,743,135 B, 11 | 941 B, 14 |
| Render/embedded_images | 2,513,897 B, 238 | 659 B, 8 | 504 B, 7 |
| Render/rect-wrong | 1,326 B, 7 | 1,657 B, 8 | 652 B, 9 |
| RenderToJPEG | 2,607,783 B, 244 | 84,099 B, 13 | 82,746 B, 13 |
| GetPageText | 12,397 B, 197 | 12,397 B, 197 | 12,393 B, 197 |

wazy and wazero now compile at essentially the same speed, but wazy's `Init` allocates roughly 3x the bytes and 1.9x
the objects it did at v0.4.0 (377 MB / 94K allocs then) — a real regression in compile-time footprint, worth watching
on future wazy updates even though wall-clock time didn't move much. wago compiles in a fifth of wazero's time with a
sixth of the allocations and creates a worker faster too. In execution, wago is faster than wazero on the small text
page, the embedded-image page and JPEG encode, but 19–39 percent slower on the ICC/transparency and PNG-predicted
pages, and over 3x slower on text extraction, a workload of many small calls into the module where its per-call
overhead shows. wazy is faster than wazero only on the tiny `test.pdf` page; it's 21–72 percent slower everywhere
else.

### Compilation workers

`Init` above uses each runtime's default compilation mode. wago compiles functions in parallel by default; wazero and
wazy compile with a single goroutine unless asked otherwise, through a context value — `experimental.WithCompilationWorkers`
for wazero, `api.WithCompilationWorkers` for wazy — and go-pdfium passes `Config.Context` to the compiler, so the
option can be set there. This only affects `Init`.

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
|---|---|---|---|---|---|---|---|---|---|---|---|
| wazero | 1.255s | 5000 | 0 | 0 | 71.6 s | 14.31 ms | 7.40 ms | 30.7 ms | 85 ms | 1.00x | reference |
| wazy | 1.233s | 5000 | 0 | 0 | 89.9 s | 17.98 ms | 8.48 ms | 43.8 ms | 112 ms | 1.178x | 0 of 5000 |
| wago | 260ms | 5000 | 0 | 0 | 82.1 s | 16.41 ms | 7.22 ms | 35.1 ms | 118 ms | 1.017x | 0 of 5000 |

wazy is slower than wazero on the corpus as a whole, though it still compiles about as fast. wago essentially matches
wazero on the whole corpus (1.02x) while compiling in a fifth of the time, and no longer shows the two pixel
mismatches it used to on a 1,000-document sample — one of the other arm64 fixes landed alongside the compile fix
appears to have also fixed that.

### The cost of close-on-context-done

All three runtimes can interrupt a running call when a context is cancelled (`WithCloseOnContextDone` on wazero and
wazy, `InterruptibleCalls` on wago), which is what makes `Kill()` work on a document that never finishes. It is off
by default in go-pdfium. Measured on the same 5,000-document corpus with the option on:

| Runtime | Compile + first worker | Rendered | Failed | Timed out | Total render time | Mean | Median | p90 | p99 | Geomean ratio vs wazero | Pixel mismatches vs wazero |
|---|---|---|---|---|---|---|---|---|---|---|---|
| wazero | 1.256s | 5000 | 0 | 0 | 324.3 s | 64.86 ms | 27.72 ms | 175.6 ms | 459 ms | 1.00x | reference |
| wazy | 1.235s | 5000 | 0 | 0 | 105.0 s | 21.00 ms | 9.76 ms | 51.1 ms | 136 ms | 0.348x | 0 of 5000 |
| wago | 254ms | 5000 | 0 | 0 | 83.6 s | 16.73 ms | 7.37 ms | 36.1 ms | 120 ms | 0.267x | 0 of 5000 |

wazero becomes roughly 4.5 times slower with the option on; wazy and wago barely change. So even though wazy and wago
are both close to or slower than wazero in the default configuration above, they are 3-4 times faster than wazero
once `Kill()` capability is turned on — wago is in fact the fastest runtime in this mode.

## wazero variants

The same benchmarks run against three unreleased wazero builds, each substituted for v1.12.0 through a Go workspace
`replace` directive: plain `main`, `main` + [#2533](https://github.com/wazero/wazero/pull/2533) (function entry as a
termination checkpoint), and `main` + #2533 + [#2529](https://github.com/wazero/wazero/pull/2529) +
[#2530](https://github.com/wazero/wazero/pull/2530) (wazevo group-ID fix; clone-free `try_table` checkpoints via
top-relative offsets and a trampoline). All three pass go-pdfium's complete wazero test suite. wazy is included as a
control in every run; its numbers don't move, confirming the variants aren't affecting machine conditions.

### Micro-benchmarks (native arm64, median of 5)

| Benchmark | v1.12.0 | main | main+2533 | main+2529+2530+2533 |
|---|---|---|---|---|
| Init | 1.145 s | 1.175 s | 1.179 s | 1.193 s |
| NewInstance | 476.2 µs | 429.1 µs | 432.1 µs | 447.7 µs |
| OpenDocument | 9.0 µs | 8.5 µs | 8.4 µs | 8.4 µs |
| Render/test | 199.1 µs | 151.3 µs | 153.2 µs | 135.1 µs |
| Render/alpha_channel | 27.27 ms | 25.84 ms | 25.75 ms | 26.31 ms |
| Render/embedded_images | 1.070 ms | 950.7 µs | 948.5 µs | 928.9 µs |
| Render/rect-wrong | 6.691 ms | 6.645 ms | 6.605 ms | 6.608 ms |
| RenderToJPEG | 7.144 ms | 6.853 ms | 6.818 ms | 6.802 ms |
| GetPageText | 30.1 µs | 28.7 µs | 29.6 µs | 28.2 µs |

### Corpus, default configuration (5,000 documents)

| wazero variant | Total render time | wazy (control) |
|---|---|---|
| v1.12.0 | 69.0 s | 87.2 s |
| main | 66.5 s | 86.3 s |
| main+2533 | 67.3 s | 87.0 s |
| main+2529+2530+2533 | 65.1 s | 87.1 s |

### Corpus with close-on-context-done (5,000 documents)

| wazero variant | Total render time | wazy (control) |
|---|---|---|
| v1.12.0 | 306.7 s | 100.4 s |
| main | 309.1 s | 102.0 s |
| main+2533 | 70.6 s | 101.4 s |
| main+2529+2530+2533 | 66.7 s | 100.3 s |

`main` on its own is a few percent faster than v1.12.0 on most benchmarks (the JPEG encode path moves most, since it
runs libjpeg-turbo's SIMD kernels inside the module) but does **not** touch the close-on-context-done cost — it's
actually slightly slower there. #2533 is what removes that cost almost entirely: 306.7 s down to 70.6 s, within a few
percent of the same build's own non-interruptible time (67.3 s), so enabling `WithCloseOnContextDone` so that `Kill()`
can interrupt a stuck render stops needing to be weighed against throughput. #2529 + #2530 add a further few percent
on top in the default configuration without touching the interruptible cost either way. Combined, `main+2529+2530+2533`
is the fastest configuration in both modes: 65.1 s default, 66.7 s interruptible — close-on-context-done is
essentially free, and wazero also overtakes wazy in that mode (wazy: 100.3 s on the same files). This section
predates wago's arm64 fix, so it only has wazero and wazy numbers; it has not been rerun with wago included.
