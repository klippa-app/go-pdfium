# Native vs WebAssembly: where the runtimes stand

This report answers three questions: how far the WebAssembly backend now is from the native (CGO) backend, how much of that depends on the native side receiving our patches (which still have to be upstreamed to PDFium), and what the WebAssembly-specific patch set consists of and where it lives.

**All patches can be found in the [`wasm-standalone` branch of jerbob92/pdfium-binaries](https://github.com/jerbob92/pdfium-binaries/tree/wasm-standalone/patches):** portable patches in [`patches/`](https://github.com/jerbob92/pdfium-binaries/tree/wasm-standalone/patches), wasm-only patches in [`patches/wasm/`](https://github.com/jerbob92/pdfium-binaries/tree/wasm-standalone/patches/wasm), applied by [`steps/03-patch.sh`](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/steps/03-patch.sh).

## The comparison

Both scenarios use the same WebAssembly configuration: the fully patched `pdfium.wasm` running under wazero `main` with the two pending performance PRs ([#2514](https://github.com/wazero/wazero/pull/2514), [#2515](https://github.com/wazero/wazero/pull/2515)) merged. What varies is the native side:

- **Scenario A — native fully patched**: the future state, once our PDFium patches are upstreamed.
- **Scenario B — native unpatched**: the present state; stock PDFium as everyone builds it today.

Both native builds come from the same PDFium source revision and build configuration as the wasm module (release, AGG, no V8/XFA), so this is a true apples-to-apples comparison of the *runtimes*, not of different PDFium versions. Benchmark: open + render page 1 at 200 DPI + close, per document, through go-pdfium's high-level API; Apple M5; benchstat n=6, all three configurations interleaved in one session; every result p=0.002.

### Scenario A: both fully patched — wasm is 2.04× native (geomean)

| Document | CGO | WASM | ratio |
|---|---|---|---|
| alpha_channel (alpha-heavy) | 47.8ms | 89.3ms | 1.87× |
| image_page (scanned page, large JPEG) | 53.8ms | 101.2ms | 1.88× |
| mona (photo) | 0.27ms | 0.54ms | 2.00× |
| invoice B | 2.93ms | 6.00ms | 2.05× |
| invoice A | 4.55ms | 9.48ms | 2.08× |
| building_plan (CAD-style) | 1.44ms | 3.06ms | 2.13× |
| test (minimal text page) | 0.62ms | 1.34ms | 2.16× |
| rect-wrong (PNG-predicted flate) | 13.2ms | 28.9ms | 2.19× |
| **geomean** | **4.26ms** | **8.70ms** | **2.04×** |

### Scenario B: native as it is today — wasm is 1.83× native (geomean)

| Document | CGO (stock) | WASM | ratio |
|---|---|---|---|
| **image_page** | 70.1ms | 101.2ms | **1.44×** |
| **rect-wrong** | 18.1ms | 28.9ms | **1.59×** |
| alpha_channel | 52.4ms | 89.3ms | 1.70× |
| mona | 0.31ms | 0.54ms | 1.73× |
| invoice A | 4.67ms | 9.48ms | 2.03× |
| invoice B | 2.92ms | 6.00ms | 2.06× |
| test | 0.64ms | 1.34ms | 2.08× |
| building_plan | 1.44ms | 3.06ms | 2.12× |
| **geomean** | **4.76ms** | **8.70ms** | **1.83×** |

### What the difference between the scenarios tells us

Against today's stock native builds, the patched WebAssembly backend is within **1.4–1.7×** on exactly the document types the optimization campaign targeted (image decode, PNG-predicted streams, alpha compositing). Once native receives the same portable patches, the gap becomes a near-uniform **~1.9–2.2×** across every document type. That uniformity is the signature of *runtime* overhead rather than algorithmic difference: wazero's single-pass JIT code quality, wasm-level memory bounds checking, and 128-bit-only SIMD. Further closing of the gap now has to come from the wazero side (PRs like #2514/#2515) — the PDFium-level algorithmic headroom is spent.

## The patch inventory

### Portable patches (help native too — proposed upstream)

These are the patches whose upstreaming Scenario A models. Native numbers per patch come from the upstream proposal document (measured one patch at a time against PDFium's own test suite); wasm numbers were measured during the campaign:

| Patch | What it does | wasm effect | native effect |
|---|---|---|---|
| [png_predict_line_perf](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/png_predict_line_perf.patch) | PNG predictor on raw pointers, per-row validation | rect-heavy renders −7% | +0.6% (enables the next patch) |
| [png_sub_filter_perf](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/png_sub_filter_perf.patch) | Sub filter channel chains carried in registers | rect −6.4% | rect −23% |
| [stretch_engine_perf](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/stretch_engine_perf.patch) | interpolation weights read through raw pointer | image_page −16% | image_page −24%, mona −14% |
| [compositor_perf](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/compositor_perf.patch) | opaque fast path + exact reciprocal division | image_page −13% | alpha −5% |
| [glyph_blend_perf](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/glyph_blend_perf.patch) | glyph blending helpers on raw pointers | rect −3.4% | neutral |
| [fillrect_memset_perf](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/fillrect_memset_perf.patch) | byte-uniform FillRect as memset | background fill ~free | invoice A −4.7% |

### WebAssembly-specific patches

These exploit wasm SIMD128 or work around wasm/emscripten platform gaps; they are compiled only for the wasm target (GN args `jpeg_wasm_simd`/`zlib_wasm_simd`, or `__wasm_simd128__` guards inside shared files).

**libjpeg-turbo SIMD kernels** — [`patches/wasm/jpeg_simd/`](https://github.com/jerbob92/pdfium-binaries/tree/wasm-standalone/patches/wasm/jpeg_simd) + [jpeg_simd_wasm.patch](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/wasm/jpeg_simd_wasm.patch). Ports of libjpeg-turbo's NEON kernels to wasm SIMD128, written to be **bit-identical to the scalar C paths** (unlike upstream's NEON/SSE2, so render hashes are unchanged):

- Decode: YCbCr→RGB (7 layouts), accurate (islow) IDCT, h2v2 + h2v1 fancy chroma upsampling. Pure JPEG decode of a 2.8MP image: 15.3ms → 6.38ms (**−58%**, ~1.3× faster than Go's `image/jpeg`); JPEG-bearing renders −10–14%.
- Encode: RGB→YCbCr (7 layouts), forward islow DCT, sample conversion + quantization (per-lane variable shift emulated with swizzle-built power-of-two multipliers), h2v1/h2v2 downsampling, baseline Huffman block encoding, progressive AC prepare stages. Encoding a 1415×2000 RGBA page at q75: 20.3ms scalar → **7.0ms** (1.94× faster than `image/jpeg`); progressive 39.3ms → 24.7ms.

**JPEG encoder exposure** — [jpeg_encode_shim.c](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/wasm/jpeg_encode_shim.c). PDFium only decodes JPEG, so libjpeg-turbo's compressor is dead-stripped; this shim exports `pdfium_jpeg_encode`/`pdfium_jpeg_free` (RGB/RGBA/BGRA/GRAY input, quality, optional progressive, `jpeg_mem_dest` output, setjmp error handling). go-pdfium's `RenderToFile` uses it for JPEG output on the wasm backend, matching what the cgo backend does with `pdfium_use_turbojpeg`.

**zlib SIMD** — [adler32_simd_wasm.patch](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/wasm/adler32_simd_wasm.patch), [inflate_chunk_wasm.patch](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/wasm/inflate_chunk_wasm.patch). Enables Chromium zlib's SIMD adler32 and chunked inflate copies for wasm128 (rect-heavy renders −13%).

**Compositor/stretch SIMD** (`__wasm_simd128__`-gated inside shared files) — [compositor_wasm_simd.patch](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/compositor_wasm_simd.patch), [stretch_engine_wasm_simd.patch](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/stretch_engine_wasm_simd.patch), [stretch_horz_wasm_simd.patch](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/stretch_horz_wasm_simd.patch). 4-pixels-per-iteration alpha compositing (branches as bitselect masks, exact f32-division with integer fixups) and channel-SIMD horizontal stretch (alpha renders −5%, stretch −2.8%).

**Platform plumbing** — [`patches/wasm/pdfium.patch` / `build.patch`](https://github.com/jerbob92/pdfium-binaries/tree/wasm-standalone/patches/wasm) (emscripten target support), a bulk-memory `memset` shim linked in [`steps/06-build.sh`](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/steps/06-build.sh) (lowers to `memory.fill`), a wasm-opt pass translating the exception encoding to standardized exnref (required by wazero), and an opt-in `PDFium_PROFILING_NAMES` build flag that preserves the wasm name section for profilers.

**Test suite under wasm** — [testsuite_emscripten.patch](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/wasm/testsuite_emscripten.patch) + [README-testsuite.md](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/wasm/README-testsuite.md). Makes `pdfium_unittests`/`pdfium_embeddertests` buildable with the emscripten toolchain and runnable under node. Used to verify that all of the above causes **zero test regressions**: patched and unpatched wasm builds fail the exact same (environment-caused) tests.

### wazero-side work

The benchmark's wazero includes upstream PRs [#2514](https://github.com/wazero/wazero/pull/2514) (elide bounds checks for constant addresses) and [#2515](https://github.com/wazero/wazero/pull/2515) (shared conditional-trap exits, −3–9% on go-pdfium renders), which we benchmarked and validated. Separately, our wazero contributions from this campaign: exception-handling checkpoint pooling and a perfmap fix (`ClearIndex` instead of `Clear` on arm64 branch-relocation reruns) that makes guest-symbol profiling usable.

## Correctness

Every configuration in the tables above passes the full go-pdfium test suite and renders hash-identically across the golden corpus. The wasm SIMD kernels are bit-identical to PDFium's scalar paths by construction; the portable patches change memory access patterns only (the one arithmetic substitution — reciprocal division — is exhaustively verified exact). PDFium's own test suite was run for both native (per patch) and wasm (patched vs clean baseline) with identical results throughout.
