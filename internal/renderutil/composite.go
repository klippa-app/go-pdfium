package renderutil

import "image"

// CompositeOnWhiteInPlace flattens straight-alpha (non-premultiplied) RGBA
// pixels onto an opaque white background, writing the result back into
// img.Pix. PDFium's FPDFBitmap_BGRA output has straight alpha.
//
// Blending in place avoids allocating a second full-size image. For the
// webassembly implementation this also keeps the pixel data inside the guest
// bitmap, so the JPEG encoder can borrow it instead of copying it back into
// WASM memory.
//
// The integer arithmetic reproduces image/draw's NRGBA-over-RGBA path with a
// white destination exactly, so the result is bit-identical to drawing the
// image with draw.Over onto a white image.NewRGBA. This matters because the
// render golden hashes depend on it. The equivalence is checked against the
// image/draw implementation in composite_test.go.
func CompositeOnWhiteInPlace(img *image.RGBA) {
	const m = 1<<16 - 1
	rowLen := img.Rect.Dx() * 4
	for y := 0; y < img.Rect.Dy(); y++ {
		row := img.Pix[y*img.Stride : y*img.Stride+rowLen]
		for i := 0; i < len(row); i += 4 {
			px := row[i : i+4 : i+4]
			sa := uint32(px[3]) * 0x101
			if sa == m {
				// Opaque pixel, nothing to blend.
				continue
			}
			a := (m - sa) * 0x101
			white := 0xff * a / m
			px[0] = uint8((white + uint32(px[0])*sa/0xff) >> 8)
			px[1] = uint8((white + uint32(px[1])*sa/0xff) >> 8)
			px[2] = uint8((white + uint32(px[2])*sa/0xff) >> 8)
			px[3] = uint8((white + sa) >> 8)
		}
	}
}
