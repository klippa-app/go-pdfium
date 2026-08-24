# Native vs WebAssembly: where the runtimes stand

This report answers three questions: how far the WebAssembly backend now is from the native (CGO) backend, how much of that depends on the native side receiving our patches (which still have to be upstreamed to PDFium), and what the WebAssembly-specific patch set consists of and where it lives.

**All patches can be found in the [`wasm-standalone` branch of jerbob92/pdfium-binaries](https://github.com/jerbob92/pdfium-binaries/tree/wasm-standalone/patches):** portable patches in [`patches/`](https://github.com/jerbob92/pdfium-binaries/tree/wasm-standalone/patches), wasm-only patches in [`patches/wasm/`](https://github.com/jerbob92/pdfium-binaries/tree/wasm-standalone/patches/wasm). [`steps/03-patch.sh`](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/steps/03-patch.sh) is the authoritative list: it decides which patches apply to which target, and several of them are now target-gated rather than applied everywhere.

**What this describes.** The patch inventory below is the state of that branch as of `e33a1ef` (2026-08-23); the `webassembly/pdfium.wasm` in this repository was rebuilt hours later the same day. The benchmark tables were measured on 2026-08-19, against the patch set of that day — six patches and two kernel sets landed after them, all of which help both backends or wasm only, so treat the tables as the state of that build rather than as today's ceiling. Every patch's own measured effect is listed in [the inventory](#the-patch-inventory).

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

At corpus scale the same ratio holds: over the 5000-document corpus used for the per-patch work, fully patched native renders in ≈36.4s at 200 DPI and wasm in ≈2.44× that (per-cluster ratios 2.19–2.79). The wasm total there is extrapolated from samples rather than measured document-for-document, unlike the eight documents above.

### What the difference between the scenarios tells us

Against today's stock native builds, the patched WebAssembly backend is within **1.4–1.7×** on exactly the document types the optimization campaign targeted (image decode, PNG-predicted streams, alpha compositing). Once native receives the same portable patches, the gap becomes a near-uniform **~1.9–2.2×** across every document type. That uniformity is the signature of *runtime* overhead rather than algorithmic difference: wazero's single-pass JIT code quality, wasm-level memory bounds checking, and 128-bit-only SIMD.

Further closing of the gap is now mostly runtime work — guard-page memory in wazero would remove the bounds check from every access, which is the single change that would lift every wasm number here. On the PDFium side the large wasm-visible kernels are done; what remains is narrower: the 8bpp/palette and blend-mode compositor rows that do not fit the kernel pattern, and an exact-match ICC scanline memo. (The one big unclaimed item, a NEON/SSE twin of the wasm gather kernels, would speed up *native* and widen this gap rather than close it.)

## The patch inventory

Native effects are from the 5000-document corpus benchmark, measured one patch at a time, null-controlled; wasm effects are from the corpus samples and the eight-document bench above. "corpus" means the whole-corpus geomean, everything else is the affected subset.

### Portable patches under upstream review

Three CLs are live on Gerrit; Scenario A models the state where they have landed. All three were rewritten for review to use spans and `fxcrt::Zip()` instead of raw pointers, so they carry no `UNSAFE_BUFFERS` — and the safe versions measured at least as fast as the pointer versions they replaced.

| Patch | What it does | wasm effect | native effect |
|---|---|---|---|
| [png_predictor_perf](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/png_predictor_perf.patch) ([CL 155510](https://pdfium-review.googlesource.com/c/pdfium/+/155510)) | hoisted row boundaries, Sub filter carried in registers as pixel structs, `Zip()` lead loops | −13.8% on a fully Sub-filtered document | PNG-predicted documents −3.29%, corpus −0.59% |
| [compositor_perf](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/compositor_perf.patch) ([CL 155530](https://pdfium-review.googlesource.com/c/pdfium/+/155530)) | typed-source-span opaque BGR row, opaque fast path skipping the per-pixel alpha division | image_page −13% | corpus −1.73% at 200 DPI, −2.63% at 300 |
| [stretch_engine_perf](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/stretch_engine_perf.patch) ([CL 155550](https://pdfium-review.googlesource.com/c/pdfium/+/155550)) | `GetWeights()` span accessor, horizontal `Zip()` + 1/2-tap unroll, `StretchVert` loop interchange and uninitialised accumulator | image_page −16% | corpus −11.1% / −13.7% (200/300 DPI), scanned and photo documents −46/−47% |

The stretch patch also carries a follow-up that is validated but not uploaded yet (it chains onto CL 155550): 1/2-tap destination columns planned into runs once per engine and replayed per row, with a bail-out that skips planning for the 70.8% of engines doing less than 4096 dest_cols × src_rows of work. Native corpus **−2.38%** on top of the CL, bilinear band −4.87%, 1-tap band −14.5%, bit-identical over 90k renders.

### Portable patches not yet upstreamed

Validated on both backends and living in the fork; each has a local upstream branch that compiles against `origin/main` with clean unit and embedder differentials, waiting to be uploaded.

| Patch | What it does | wasm effect | native effect |
|---|---|---|---|
| [jpeg_decode_bgr](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/jpeg_decode_bgr.patch) | decode 3-component JPEGs straight to `JCS_EXT_BGR`, turning `TranslateScanline24bppDefaultDecode` into a pass-through; extended to ICCBased color spaces whose profile `IsSRGB()` | JPEG-heavy sample −12%, ICC-heavy sample −5.2% | corpus −1.64%, totals −4.5%, heavy pass-through documents −11 to −26%; ICC-heavy subset −1.25% |
| [swap_translate_perf](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/swap_translate_perf.patch) | 24bpp RGB→BGR scanline swap through `fxcodec::ReverseRGB()`, which clang already vectorizes; wasm gets a `__wasm_simd128__` kernel because it has no interleaved-store lowering | swap-heavy sample −16% | swap-heavy −10.2%, corpus −0.52% |
| [lcms_translate_memo](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/lcms_translate_memo.patch) | exact-match memo on the quantized tuple in `IccTransform::Translate` — lcms2 refuses to optimize PDFium's transforms, so parametric tone curves were running libm `pow` per pixel | tail sample −21% geomean, −66.5% best | heaviest documents −68/−30/−25/−22%, corpus totals −0.45% |
| [t4_psfunc_memo](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/t4_psfunc_memo.patch) | exact-match memo in `CPDF_PSFunc::v_Call()` keyed on raw input float bits | worst document −68%, PostScript subset −13% | worst document −65% (340→119ms), corpus totals −0.4% |
| [alpha_unroll_native](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/alpha_unroll_native.patch) | 1/2-tap unroll for the alpha stretch arm, which never got CL 155550's unroll. **Native only** | not applied — the equivalent peel of the wasm SIMD tap loop measured flat | alpha documents −1.5% |

Both memos are exact-match caches: bit-identical output, verified over 2×5000 native renders plus wasm samples.

### WebAssembly-specific patches

These exploit wasm SIMD128 or work around wasm/emscripten platform gaps. They are compiled only for the wasm target, either because `03-patch.sh` applies them only there or because the added code sits behind `__wasm_simd128__` guards inside shared files.

| Patch | What it does | wasm effect | why not native |
|---|---|---|---|
| [stretch_horz_wasm_simd](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/stretch_horz_wasm_simd.patch) | `i8x16.swizzle` 4-column gather kernels for the planned 1/2-tap runs, in 16- and 32-byte window variants, emitted only when ≥75% of a run's columns fit | isolated −18.3% on the bilinear-heavy sample, −24% cumulative; per-column WAT ops 77 → 25 | guarded intrinsics; the native twin is a NEON prototype pending an upstream intrinsics decision |
| [compositor_wasm_simd](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/compositor_wasm_simd.patch) | `CopyRowToOpaqueBgra` in all four BGR/BGRx → BGRA/RGBA instantiations, the opaque BGRA row kernel templated over destination byte order (the shipped BGRA-dest-only version was dead code: go-pdfium sets `FPDF_REVERSE_BYTE_ORDER` on every color render, so the destination is always RGBA), clip-row kernels, and an all-transparent-group fast path | compositor-heavy sample −7.1%, alpha_channel −41.9%, bench8 −3.28% | guarded intrinsics |
| [c3_compositor_wasm_simd](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/c3_compositor_wasm_simd.patch) | clip-masked `Bgra2Bgra` kernel over both destination byte orders, and the `Rgb2Rgb` 3↔4-Bpp conversion shuffles (399M px over 1305 corpus documents) | clip documents ≈−6% (−10% heavy), conversion documents ≈−2.7% | guarded intrinsics |
| [stretch_engine_wasm_simd](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/stretch_engine_wasm_simd.patch) | channel-SIMD horizontal stretch | stretch −2.8% | guarded intrinsics |
| [glyph_blend_perf](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/glyph_blend_perf.patch) | glyph blending helpers on raw pointers | rect-heavy renders −3.4% | measures neutral on native |
| [fillrect_memset_perf](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/fillrect_memset_perf.patch) | byte-uniform `FillRect` as `memset` | background fill ~free | *regresses* small cache-resident fills on native, where libc `memset` switches to non-temporal stores |

The last two only pay off under a runtime that neither vectorizes nor elides its own bounds checks; they were applied to every target until 2026-08-20 and are now gated to emscripten.

**libjpeg-turbo SIMD kernels** — [`patches/wasm/jpeg_simd/`](https://github.com/jerbob92/pdfium-binaries/tree/wasm-standalone/patches/wasm/jpeg_simd) + [jpeg_simd_wasm.patch](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/wasm/jpeg_simd_wasm.patch). Ports of libjpeg-turbo's NEON kernels to wasm SIMD128, written to be **bit-identical to the scalar C paths** (unlike upstream's NEON/SSE2, so render hashes are unchanged):

- Decode: YCbCr→RGB (7 layouts), accurate (islow) IDCT, h2v2 + h2v1 fancy chroma upsampling. Pure JPEG decode of a 2.8MP image: 15.3ms → 6.38ms (**−58%**, ~1.3× faster than Go's `image/jpeg`); JPEG-bearing renders −10–14%. The layout dispatch already covers `JCS_EXT_BGR`, so `jpeg_decode_bgr` needed no wasm-side change.
- Encode: RGB→YCbCr (7 layouts), forward islow DCT, sample conversion + quantization (per-lane variable shift emulated with swizzle-built power-of-two multipliers), h2v1/h2v2 downsampling, baseline Huffman block encoding, progressive AC prepare stages. Encoding a 1415×2000 RGBA page at q75: 20.3ms scalar → **7.0ms** (1.94× faster than `image/jpeg`); progressive 39.3ms → 24.7ms.

**JPEG encoder exposure** — [jpeg_encode_shim.c](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/wasm/jpeg_encode_shim.c). PDFium only decodes JPEG, so libjpeg-turbo's compressor is dead-stripped; this shim exports `pdfium_jpeg_encode`/`pdfium_jpeg_free` (RGB/RGBA/BGRA/GRAY input, quality, optional progressive, `jpeg_mem_dest` output, setjmp error handling). go-pdfium's `RenderToFile` uses it for JPEG output on the wasm backend, matching what the cgo backend does with `pdfium_use_turbojpeg`.

**zlib SIMD** — [adler32_simd_wasm.patch](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/wasm/adler32_simd_wasm.patch), [inflate_chunk_wasm.patch](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/wasm/inflate_chunk_wasm.patch). Enables Chromium zlib's SIMD adler32 and chunked inflate copies for wasm128 (rect-heavy renders −13%).

**Platform plumbing** — [`patches/wasm/pdfium.patch` / `build.patch`](https://github.com/jerbob92/pdfium-binaries/tree/wasm-standalone/patches/wasm) (emscripten target support), [callbacks.patch](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/wasm/callbacks.patch) (the `FPDF_FILEACCESS_Create` / `FX_FILEAVAIL_Create` / `FPDF_FORMFILLINFO_Create` / `FPDF_FILEWRITE_Create` / `FX_DOWNLOADHINTS_Create` constructors, their `_Size` helpers and the `FPDF_FORMFILLINFO_CALL_TIMER` trampoline that the wasm backend needs to hand PDFium structs and receive callbacks), a bulk-memory `memset` shim linked in [`steps/06-build.sh`](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/steps/06-build.sh) (lowers to `memory.fill`), a wasm-opt pass translating the exception encoding to standardized exnref (required by wazero), and an opt-in `PDFium_PROFILING_NAMES` build flag that preserves the wasm name section for profilers.

**Test suite under wasm** — [testsuite_emscripten.patch](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/wasm/testsuite_emscripten.patch) + [README-testsuite.md](https://github.com/jerbob92/pdfium-binaries/blob/wasm-standalone/patches/wasm/README-testsuite.md). Makes `pdfium_unittests`/`pdfium_embeddertests` buildable with the emscripten toolchain and runnable under node. Used to verify that all of the above causes **zero test regressions**: patched and unpatched wasm builds fail the exact same (environment-caused) tests.

### Upstreaming status

| Item | State |
|---|---|
| CL 155510 (PNG predictor) | live, patchset 11, awaiting re-review |
| CL 155530 (compositor) | live, patchset 10, awaiting re-review |
| CL 155550 (stretch engine) | live, patchset 9, awaiting re-review |
| stretch run-planning follow-up | local branch, validated, chains onto 155550 — uploads after it lands |
| `jpeg_decode_bgr` (+ ICC-sRGB) | local branch, independent, ready to upload |
| `swap_translate_perf` | local branch, independent, ready to upload |
| `lcms_translate_memo` | local branch, independent, ready to upload |
| `t4_psfunc_memo` | local branch, independent, ready to upload |
| alpha arm unroll | folded onto the stretch follow-up branch |
| wasm SIMD kernels, glyph/fillrect | fork-only by nature (guarded intrinsics, or native-negative) |

### wazero-side work

The benchmark's wazero includes upstream PRs [#2514](https://github.com/wazero/wazero/pull/2514) (elide bounds checks for constant addresses) and [#2515](https://github.com/wazero/wazero/pull/2515) (shared conditional-trap exits, −3–9% on go-pdfium renders), which we benchmarked and validated. Separately, our wazero contributions from this campaign: exception-handling checkpoint pooling and a perfmap fix (`ClearIndex` instead of `Clear` on arm64 branch-relocation reruns) that makes guest-symbol profiling usable. Guard-page memory — the change that would remove bounds checking from every guest access — is designed and unblocked but not implemented.

## Correctness

Every configuration in the tables above passes the full go-pdfium test suite and renders hash-identically across the golden corpus. The wasm SIMD kernels are bit-identical to PDFium's scalar paths by construction; the portable patches change memory access patterns only, and the two memos are exact-match caches, so they are bit-identical as well (the one arithmetic substitution in the compositor, reciprocal division, was removed again — it was a measured loss on arm64). PDFium's own test suite was run for both native (per patch) and wasm (patched vs clean baseline) with identical results throughout — always as a differential, comparing failing-*name* sets rather than counts: the pinned checkout's baseline is 1033/1033 unit tests, and the embedder and pixel suites' pre-existing environment failures came back as byte-identical name sets under every configuration built.

Each patch was also validated against a null control (a second, byte-identical build of the baseline), so effects smaller than the measurement floor were never counted, and against poison builds (a deliberately wrong constant must change the render hash) to prove the new code is actually reached.
