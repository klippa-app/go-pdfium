# WebAssembly runtime benchmarks

go-pdfium can run the PDFium WebAssembly module on more than one runtime: [wazero](https://wazero.io/) in the
`webassembly` package and [wazy](https://github.com/samyfodil/wazy) in the experimental `experimental/wazy` package.
This document compares them. The experimental `experimental/wago` backend is left out for now: wago can not compile the
module on arm64 and miscompiles part of it on amd64 (see its package documentation).

Everything here goes through the public `pdfium.Pool` API of each backend, so the numbers include go-pdfium's own
overhead in the same way for both runtimes. All renders were compared pixel for pixel: the two runtimes produced
identical output for every document.

## Setup

| | |
|---|---|
| go-pdfium | this commit, module `internal/pdfium_wasm/pdfium.wasm` (PDFium 8057, 5,754,045 bytes) |
| wazero | v1.12.0, compiler engine, `CoreFeaturesV2 \| CoreFeaturesExceptionHandling` |
| wazy | v0.3.0, compiler engine, `CoreFeaturesV2 \| CoreFeatureExceptionHandling` |
| Go | go1.27.1 |
| Machine | Apple M5 (10 cores), macOS 26 (Darwin 25.6.0), arm64 |

Micro-benchmark numbers are the median of five `go test -bench` runs (`-count 5`). The corpus numbers come from a
single interleaved pass over all documents.

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

All numbers in this document were measured on arm64. The `.github/workflows/benchmark.yml` workflow runs the same
micro-benchmarks on GitHub's amd64 and arm64 runners on demand.

| Benchmark | wazero | wazy | wazy vs wazero |
|---|---|---|---|
| Init | 1.14 s | 719.27 ms | 0.63x |
| NewInstance | 437.6 µs | 234.5 µs | 0.54x |
| OpenDocument | 8.8 µs | 6.7 µs | 0.76x |
| Render/test | 193.4 µs | 106.7 µs | 0.55x |
| Render/alpha_channel | 26.97 ms | 24.28 ms | 0.90x |
| Render/embedded_images | 1.05 ms | 816.4 µs | 0.78x |
| Render/rect-wrong | 7.04 ms | 6.28 ms | 0.89x |
| RenderToJPEG | 11.20 ms | 7.42 ms | 0.66x |
| GetPageText | 29.0 µs | 26.1 µs | 0.90x |

Allocations per operation (median B/op, allocs/op):

| Benchmark | wazero | wazy |
|---|---|---|
| Init | 345,873,464 B, 659,038 | 247,239,932 B, 153,101 |
| NewInstance | 19,305,563 B, 13,830 | 211,003 B, 98 |
| OpenDocument | 1,008 B, 22 | 1,008 B, 22 |
| Render/test | 518 B, 7 | 497 B, 7 |
| Render/alpha_channel | 5,103,415 B, 62 | 2,313,647 B, 8 |
| Render/embedded_images | 2,513,892 B, 238 | 554 B, 7 |
| Render/rect-wrong | 1,286 B, 7 | 994 B, 7 |
| RenderToJPEG | 2,608,264 B, 244 | 83,488 B, 12 |
| GetPageText | 12,396 B, 197 | 12,395 B, 197 |

wazy compiles the module about a third faster and instantiates it twice as fast with a fraction of the allocations,
which matters for the default configuration that creates a fresh worker per `GetInstance`. Execution is faster across
the board, most on short calls where runtime overhead dominates and least on long compute bound renders.

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
| wazero | 1.144s | 5000 | 0 | 0 | 69.4 s | 13.88 ms | 7.09 ms | 30.7 ms | 83 ms | 1.00x | reference |
| wazy | 720ms | 5000 | 0 | 0 | 60.8 s | 12.16 ms | 6.08 ms | 25.6 ms | 74 ms | 0.859x | 0 of 5000 |

Per document, wazy is faster on 4,973 of the 5,000 documents. The distribution of the wazy/wazero time ratio is narrow:
0.82 at the 10th percentile, 0.86 at the median and 0.91 at the 90th percentile, so the gain is a fairly uniform
10 to 20 percent rather than a few outliers. The slowest documents (a second per page) are within a few percent of
each other on both runtimes, which is expected: at that point the time is spent in PDFium's own compute and the
runtime overhead disappears.

### The cost of close-on-context-done

Both runtimes can interrupt a running call when a context is cancelled (`WithCloseOnContextDone` on the runtime
configuration), which is what makes `Kill()` work on a document that never finishes. It is off by default in
go-pdfium. Enabling it changes the picture completely, measured on every fifth document of the corpus (1,000 files):

| Close-on-context-done | Runtime | Total render time | Median | p90 | Geomean ratio vs wazero |
|---|---|---|---|---|---|
| off (default) | wazero | 14.4 s | 6.80 ms | 31.3 ms | 1.00x |
| off (default) | wazy | 12.7 s | 5.87 ms | 26.2 ms | 0.861x |
| on | wazero | 64.8 s | 25.36 ms | 175.1 ms | 1.00x |
| on | wazy | 18.0 s | 8.00 ms | 42.2 ms | 0.304x |

wazero becomes about 4.5 times slower with the option on, wazy about 1.4 times. If you need `Kill()` to interrupt
stuck renders, this is the biggest performance difference between the two runtimes by far, and worth knowing about
even if you stay on wazero.
