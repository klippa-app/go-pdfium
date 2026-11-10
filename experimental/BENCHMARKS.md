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
| wazy | `v0.7.0` (2026-11-06), compiler engine, `CoreFeaturesV2 \| CoreFeatureExceptionHandling` |
| wago | `v0.1.0-beta.13.0.20261107063457-7e6a95708f67` (main, 2026-11-07), default config (explicit bounds checks), WASI plugin v0.3.1 |
| Go | go1.26.4 |
| Machine | Apple M5 (10 cores), macOS 27.0.1 (Darwin 27.0.0), arm64, otherwise idle |

wago is pinned to a commit rather than a tag: its latest tagged release, `v0.1.0-beta.13`, cannot compile this module
on arm64 at all (`arm64: native branch displacement exceeds instruction range`, a codegen limit the module outgrew);
the fix (wago-org/wago#1014) is on `main` but not yet in a tagged release.

wazy v0.6.0 removed `RuntimeConfig.WithCloseOnContextDone`: compiled code is now unconditionally interruptible by a
done context, with nothing to configure, so go-pdfium's `experimental/wazy` backend no longer sets that option (it
doesn't exist anymore) and go-pdfium's own pool code now checks the worker context explicitly before instantiating,
since wazy itself no longer rejects instantiation on an already-cancelled context the way it used to.

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
| Init | 444.79 ms | 524.81 ms | 86.72 ms | 1.18x | 0.20x |
| NewInstance | 571.7 µs | 293.2 µs | 411.3 µs | 0.51x | 0.72x |
| OpenDocument | 9.99 µs | 7.12 µs | 9.13 µs | 0.71x | 0.91x |
| Render/test | 221.87 µs | 107.44 µs | 153.78 µs | 0.48x | 0.69x |
| Render/alpha_channel | 31.01 ms | 22.34 ms | 39.34 ms | 0.72x | 1.27x |
| Render/embedded_images | 1.171 ms | 706.4 µs | 974.8 µs | 0.60x | 0.83x |
| Render/rect-wrong | 7.303 ms | 5.504 ms | 10.05 ms | 0.75x | 1.38x |
| RenderToJPEG | 7.566 ms | 6.598 ms | 9.091 ms | 0.87x | 1.20x |
| GetPageText | 30.64 µs | 30.36 µs | 94.11 µs | 0.99x | 3.07x |

Allocations per operation (median B/op, allocs/op):

| Benchmark | wazero | wazy | wago |
|---|---|---|---|
| Init | 660,571,989 B, 1,183,218 | 690,030,876 B, 481,659 | 97,858,770 B, 22,975 |
| NewInstance | 19,369,701 B, 13,890 | 207,945 B, 113 | 117,815 B, 1,252 |
| OpenDocument | 1,008 B, 22 | 1,008 B, 22 | 1,008 B, 22 |
| Render/test | 523 B, 7 | 498 B, 7 | 489 B, 7 |
| Render/alpha_channel | 6,462,940 B, 70 | 2,106,152 B, 12 | 969 B, 14 |
| Render/embedded_images | 2,513,903 B, 238 | 550 B, 7 | 504 B, 7 |
| Render/rect-wrong | 1,337 B, 7 | 979 B, 8 | 654 B, 9 |
| RenderToJPEG | 2,607,783 B, 244 | 83,436 B, 13 | 82,761 B, 13 |
| GetPageText | 12,397 B, 197 | 12,395 B, 197 | 12,393 B, 197 |

`Init` configures all three runtimes to use the same number of compilation workers (`runtime.NumCPU()`, 10 on this
machine), set through each runtime's own mechanism — `experimental.WithCompilationWorkers`/`api.WithCompilationWorkers`
via `Config.Context` for wazero and wazy, `RuntimeConfig.WithFunctionWorkers` for wago — so the comparison reflects
equal CPU budget rather than each runtime's own default worker count. With 10 cores compiling in parallel, wago is
still by far the fastest to compile (0.20x of wazero's time) with a seventh of the allocations; wazy takes about 18
percent longer than wazero to compile, allocating a similar number of bytes but well under half the objects. In
execution (unaffected by the compile-worker setting, since the timer starts after the pool is built), wazy is faster
than wazero on every benchmark except the ICC/transparency page, where it's 28 percent slower. wago is faster than
wazero on the small text page and close to it on `OpenDocument`, but slower everywhere else — substantially so on the
ICC/transparency page (1.27x) and the PNG-predicted-flate page (1.38x) — and over 3x slower on text extraction, a
workload of many small calls into the module where its per-call overhead shows.

### Compilation workers

wago compiles functions in parallel; wazero and wazy also support parallel compilation through a context value —
`experimental.WithCompilationWorkers` for wazero, `api.WithCompilationWorkers` for wazy — which go-pdfium exposes via
`Config.Context`. None of the three parallelizes by default (each uses a single worker unless told otherwise), so
`Init` above sets the same worker count explicitly on all three rather than relying on that default. This only
affects `Init`; the real world corpus benchmarks below use each runtime's own default (single-worker) compile, since
they measure one `Init` per runtime and worker-count fairness matters less for a single data point than for a
micro-benchmark run many times.

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
| wazero | 1.284s | 5000 | 0 | 0 | 82.3 s | 16.47 ms | 8.23 ms | 34.7 ms | 98 ms | 1.00x | reference |
| wazy | 1.439s | 5000 | 0 | 0 | 58.9 s | 11.79 ms | 6.16 ms | 24.2 ms | 72 ms | 0.737x | 0 of 5000 |
| wago | 238ms | 5000 | 0 | 0 | 92.7 s | 18.54 ms | 8.01 ms | 39.5 ms | 143 ms | 1.005x | 0 of 5000 |

wazy is faster than wazero on the corpus as a whole (0.737x) while compiling in about the same time. wago essentially
matches wazero on the whole corpus (1.005x) while compiling in a fifth of the time.

### The cost of close-on-context-done

wazero can interrupt a running call when its context is cancelled (`WithCloseOnContextDone`), which is what makes
`Kill()` work on a document that never finishes; it is off by default in go-pdfium. wazy interrupts a running call
unconditionally, with no option to turn on or off. wago has its own `InterruptibleCalls` option. Measured on the same
5,000-document corpus with wazero's and wago's option turned on (wazy is unaffected, since there's nothing to turn
on):

| Runtime | Compile + first worker | Rendered | Failed | Timed out | Total render time | Mean | Median | p90 | p99 | Geomean ratio vs wazero | Pixel mismatches vs wazero |
|---|---|---|---|---|---|---|---|---|---|---|---|
| wazero | 1.625s | 5000 | 0 | 0 | 346.7 s | 69.34 ms | 29.12 ms | 186.2 ms | 495 ms | 1.00x | reference |
| wazy | 1.863s | 5000 | 0 | 0 | 56.8 s | 11.36 ms | 5.97 ms | 22.9 ms | 70 ms | 0.195x | 0 of 5000 |
| wago | 284ms | 5000 | 0 | 0 | 88.8 s | 17.77 ms | 7.75 ms | 38.1 ms | 127 ms | 0.266x | 0 of 5000 |

wazero becomes roughly 4.2 times slower with the option on; wazy and wago barely change. wazy's total render time here
(56.8 s) is within noise of its default-mode run above (58.9 s), confirming its unconditional interruptibility costs
nothing extra. So even though wazy and wago are both close to or faster than wazero in the default configuration
above, they pull further ahead once `Kill()` capability is turned on — wazy is the fastest runtime in this mode.

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
