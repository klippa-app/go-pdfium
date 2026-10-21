# Native vs WebAssembly: where the runtimes stand

This report compares go-pdfium's two backends as they actually ship: the native (CGO) backend, built from stock [bblanchon/pdfium-binaries](https://github.com/bblanchon/pdfium-binaries), and the WebAssembly backend, built from the patched [`jerbob92/pdfium-binaries` `wasm-standalone` branch](https://github.com/jerbob92/pdfium-binaries/tree/wasm-standalone). It also inventories the wasm-specific patch set. **go-pdfium's native backend never applies any of jerbob92's fork patches** — it always uses stock bblanchon builds — so this is a comparison of the two real, shipping configurations, not of hypothetical patch states.

**All patches can be found in the [`wasm-standalone` branch of jerbob92/pdfium-binaries](https://github.com/jerbob92/pdfium-binaries/tree/wasm-standalone/patches):** portable patches in [`patches/`](https://github.com/jerbob92/pdfium-binaries/tree/wasm-standalone/patches), wasm-only patches in [`patches/wasm/`](https://github.com/jerbob92/pdfium-binaries/tree/wasm-standalone/patches/wasm). [`steps/03-patch.sh`](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/steps/03-patch.sh) is the authoritative list: it decides which patches apply to which target.

Both backends are measured at the exact same PDFium source revision (chromium/8111, commit `a50f1a595`): the wasm module is `internal/pdfium_wasm/pdfium.wasm` as shipped in this repository, and native is the stock bblanchon `chromium/8111` release build — a true apples-to-apples comparison of the *runtimes*, not of different PDFium versions.

## The comparison

Benchmark: open + render page 1 at 100 DPI + close, per document, through go-pdfium's high-level API (`single_threaded` for native, the default `webassembly` wazero backend for wasm); Apple M5 (10 cores), macOS 27.0.1, arm64, otherwise idle; nothing else running concurrently.

The per-document sample is the four documents in `shared_tests/testdata` that `experimental/benchmark`'s own micro-benchmarks already use — `test`, `alpha_channel`, `embedded_images` (embeds a JPEG photo) and `rect-wrong` — plus the full 5,000-document corpus result below for a statistically robust aggregate.

### Per-document (median of 21 reps, document held open)

| Document | Native (CGO) | WASM (wazero) | ratio |
|---|---|---|---|
| test (minimal text page) | 0.227 ms | 0.205 ms | 0.90× |
| embedded_images (JPEG photo) | 0.503 ms | 1.224 ms | 2.43× |
| alpha_channel (ICC + transparency) | 14.66 ms | 29.94 ms | 2.04× |
| rect-wrong (PNG-predicted flate) | 2.66 ms | 8.89 ms | 3.34× |

`test.pdf`'s render is sub-millisecond on both backends and right at the measurement floor, so its ratio is noisy; the three heavier documents land in a 2.0–3.3× range.

### Corpus scale (5,000 documents, 3 runs each, alternating order)

| Backend | Run 1 | Run 2 | Run 3 |
|---|---|---|---|
| Native (CGO, stock bblanchon 8111) | 31.59 s | 31.88 s | 31.89 s |
| WASM (wazero, patched 8111) | 66.5 s | 66.4 s | 66.3 s |

WASM is consistently **≈2.08–2.09×** native across all three runs — a tight, reproducible ratio.

### What the ratio tells us

The 2.0–3.3× range is runtime overhead, not an algorithmic gap: wazero's single-pass JIT code quality, wasm-level memory bounds checking, and 128-bit-only SIMD, on top of the same PDFium rendering code on both sides. Closing it further is almost entirely runtime work on the wasm side — guard-page memory in wazero (removing the bounds check from every guest memory access) remains the single biggest lever and is still unimplemented in what go-pdfium actually uses (see [wazero](#wazero)).

## The patch inventory

Native builds from stock bblanchon releases: no patches. WASM builds from the `wasm-standalone` fork, which currently applies:

- `lcms_translate_memo`, `t4_psfunc_memo` — portable exact-match memos (ICC transform evaluation, PostScript function evaluation). These are the only two patches here not specific to wasm; go-pdfium's native build doesn't use them since it doesn't build from this fork.
- `compositor_wasm_simd`, `c3_compositor_wasm_simd` — BGR/BGRA compositor row and clip kernels, 3↔4 Bpp conversion, behind `__wasm_simd128__` guards.
- `stretch_engine_wasm_simd` — channel-SIMD horizontal stretch, same guards.
- `glyph_blend_perf`, `fillrect_memset_perf` — gated to the emscripten target only (`fillrect_memset_perf` regresses cache-resident fills on native, so it's never applied there).
- libjpeg-turbo SIMD decode/encode kernels (bit-identical to the scalar paths) and a JPEG encoder shim — PDFium only decodes JPEG natively, so go-pdfium's wasm `RenderToFile` uses the shim for JPEG output.
- zlib SIMD (adler32, chunked inflate).
- Platform plumbing for the emscripten target: callback constructors, a `memory.fill` memset shim, exnref translation for wazero.
- emscripten-toolchain support for PDFium's own test suite.

### wazero

The comparison above measures the version go-pdfium is pinned to (`v1.12.0`, already the latest release). wazero `main` plus three pending pull requests ([#2529](https://github.com/wazero/wazero/pull/2529), [#2530](https://github.com/wazero/wazero/pull/2530), [#2533](https://github.com/wazero/wazero/pull/2533)) renders the 5,000-document corpus in 65.1 s instead of 69.0 s — a modest default-path improvement — but the same build renders it with `WithCloseOnContextDone` enabled (what makes `Kill()` work on a stuck document) in 66.7 s instead of 306.7 s, since #2533 removes almost all of that cost. See `experimental/BENCHMARKS.md` for the full wazero-variant breakdown. Guard-page memory — removing the bounds check from every guest memory access, the single biggest remaining lever on the wasm side — remains designed and unblocked but not implemented in any release.
