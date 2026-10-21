# WebAssembly runtime benchmarks

go-pdfium can run the PDFium WebAssembly module on three runtimes: [wazero](https://wazero.io/) in the `webassembly`
package, and [wazy](https://github.com/samyfodil/wazy) and [wago](https://github.com/wago-org/wago) in the
experimental `experimental/wazy` and `experimental/wago` packages. This document compares them.

Everything here goes through the public `pdfium.Pool` API of each backend, so the numbers include go-pdfium's own
overhead in the same way for all of them. All renders were compared pixel for pixel: wazero and wazy produced
identical output for every one of the 5,000 corpus documents.

**wago is excluded from every table below**: it can no longer compile the PDFium 8111 module on arm64 at all — see
[wago: arm64 compile failure](#wago-arm64-compile-failure).

## Setup

| | |
|---|---|
| go-pdfium | this commit, module `internal/pdfium_wasm/pdfium.wasm` (PDFium 8111, 5,828,543 bytes) |
| wazero | v1.12.0 (unchanged, already latest), compiler engine, `CoreFeaturesV2 \| CoreFeaturesExceptionHandling` |
| wazy | `v0.4.0` (2026-10-20), compiler engine, `CoreFeaturesV2 \| CoreFeatureExceptionHandling` |
| wago | `v0.1.0-beta.13` (2026-10-15) — excluded, see [wago: arm64 compile failure](#wago-arm64-compile-failure) |
| Go | go1.26.4 |
| Machine | Apple M5 (10 cores), macOS 27.0.1 (Darwin 27.0.0), arm64, otherwise idle |

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

| Benchmark | wazero | wazy | wazy vs wazero |
|---|---|---|---|
| Init | 1.145 s | 1.078 s | 0.94x |
| NewInstance | 476.2 µs | 264.8 µs | 0.56x |
| OpenDocument | 9.0 µs | 7.8 µs | 0.87x |
| Render/test | 199.1 µs | 132.6 µs | 0.67x |
| Render/alpha_channel | 27.27 ms | 46.93 ms | 1.72x |
| Render/embedded_images | 1.070 ms | 1.298 ms | 1.21x |
| Render/rect-wrong | 6.691 ms | 10.44 ms | 1.56x |
| RenderToJPEG | 7.144 ms | 9.836 ms | 1.38x |
| GetPageText | 30.1 µs | 32.1 µs | 1.07x |

Allocations per operation (median B/op, allocs/op):

| Benchmark | wazero | wazy |
|---|---|---|
| Init | 377,230,776 B, 721,268 | 193,354,648 B, 94,191 |
| NewInstance | 19,370,737 B, 13,891 | 207,628 B, 107 |
| OpenDocument | 1,008 B, 22 | 1,008 B, 22 |
| Render/test | 522 B, 7 | 526 B, 8 |
| Render/alpha_channel | 5,847,466 B, 64 | 4,551,376 B, 11 |
| Render/embedded_images | 2,513,899 B, 238 | 655 B, 8 |
| Render/rect-wrong | 1,244 B, 7 | 1,617 B, 8 |
| RenderToJPEG | 2,607,721 B, 244 | 84,051 B, 13 |
| GetPageText | 12,396 B, 197 | 12,396 B, 197 |

wazy still compiles a little faster and creates a worker in about half the time with a fraction of the allocations,
but it is now slower than wazero on every render benchmark except the tiny `test.pdf` page — 7 to 72 percent slower
depending on how much of the time is PDFium compute versus runtime overhead, and text extraction is close to even.

### Compilation workers

`Init` above uses each runtime's default compilation mode (a single goroutine). Both wazero and wazy accept more
through a context value — `experimental.WithCompilationWorkers` for wazero, `api.WithCompilationWorkers` for wazy —
and go-pdfium passes `Config.Context` to the compiler, so the option can be set there. This only affects `Init`.

## Real world corpus

`experimental/benchmark/cmd/corpus` renders the first page of every PDF in a directory at 100 DPI with each runtime,
interleaved per document (with a rotating order so that no runtime always runs first on a cold cache), records the
time per document and a hash of the rendered pixels, and reports totals, percentiles, the per-document geometric mean
ratio and the number of pixel mismatches:

```sh
go run ./experimental/benchmark/cmd/corpus -dir /path/to/pdfs [-runtimes wazero,wazy] [-dpi 100] [-stride N] [-csv out.csv]
```

The corpus used here is a private set of 5,000 real world PDFs (859 MB, invoices, scans, forms, reports, CAD
drawings), rendered natively on the Apple M5. Per document time includes opening the document from memory, rendering
and closing it. Both runtimes rendered all 5,000 documents and produced identical pixels for every one of them.

| Runtime | Compile + first worker | Rendered | Failed | Timed out | Total render time | Mean | Median | p90 | p99 | Geomean ratio vs wazero | Pixel mismatches vs wazero |
|---|---|---|---|---|---|---|---|---|---|---|---|
| wazero | 1.206s | 5000 | 0 | 0 | 69.0 s | 13.80 ms | 7.15 ms | 29.5 ms | 81 ms | 1.00x | reference |
| wazy | 1.156s | 5000 | 0 | 0 | 87.2 s | 17.44 ms | 8.23 ms | 42.5 ms | 114 ms | 1.184x | 0 of 5000 |

wazy is slower than wazero on the corpus as a whole, though it still compiles and starts a worker about as fast.

### The cost of close-on-context-done

Both runtimes can interrupt a running call when a context is cancelled (`WithCloseOnContextDone`), which is what
makes `Kill()` work on a document that never finishes. It is off by default in go-pdfium. Measured on the same
5,000-document corpus with the option on:

| Runtime | Compile + first worker | Rendered | Failed | Timed out | Total render time | Mean | Median | p90 | p99 | Geomean ratio vs wazero | Pixel mismatches vs wazero |
|---|---|---|---|---|---|---|---|---|---|---|---|
| wazero | 1.206s | 5000 | 0 | 0 | 306.7 s | 61.34 ms | 25.99 ms | 167.0 ms | 445 ms | 1.00x | reference |
| wazy | 1.235s | 5000 | 0 | 0 | 100.4 s | 20.09 ms | 9.27 ms | 49.3 ms | 133 ms | 0.352x | 0 of 5000 |

wazero becomes roughly 4.4 times slower with the option on; wazy barely changes. So even though wazy is slower than
wazero in the default configuration above, it is still close to 3 times faster than wazero once `Kill()` capability
is turned on.

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
essentially free, and wazero also overtakes wazy in that mode (wazy: 100.3 s on the same files).

## wago: arm64 compile failure

wago can no longer compile the PDFium 8111 module on arm64 at all:

```
could not compile webassembly module: compile: arm64: function 3986: arm64: native branch displacement exceeds instruction range
```

This reproduces identically on wago `v0.1.0-beta.9` and the current `v0.1.0-beta.13`, so it isn't a regression from
upgrading wago — it's a latent arm64-codegen limit that the module (now 5,828,543 bytes, up from 5,754,045) has grown
into. The module itself compiles fine on wazero and wazy. wago still passes go-pdfium's full test suite under Docker
linux/amd64 (1585/1600, 15 pre-existing skips), but amd64 timings under emulation on this machine aren't trustworthy
for benchmarking, and the arm64 failure blocks native testing entirely, so wago has no numbers in this document until
it's fixed upstream.
