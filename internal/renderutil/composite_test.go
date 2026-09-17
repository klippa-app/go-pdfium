package renderutil_test

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/klippa-app/go-pdfium/internal/renderutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// compositeOnWhiteReference is the image/draw based implementation that
// CompositeOnWhiteInPlace replaced. It is the source of truth for the render
// golden hashes, so the in-place version has to match it byte for byte.
func compositeOnWhiteReference(src *image.RGBA) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	straightAlphaSrc := &image.NRGBA{
		Pix:    src.Pix,
		Stride: src.Stride,
		Rect:   src.Rect,
	}
	draw.Draw(dst, dst.Bounds(), straightAlphaSrc, straightAlphaSrc.Bounds().Min, draw.Over)
	return dst
}

// fillAllColorAlphaPairs puts every colour value on the x axis against every
// alpha value on the y axis, so all 256 * 256 (colour, alpha) combinations
// are covered by a single 256x256 image.
func fillAllColorAlphaPairs(img *image.RGBA) {
	for y := range 256 {
		for x := range 256 {
			o := img.PixOffset(img.Rect.Min.X+x, img.Rect.Min.Y+y)
			img.Pix[o+0] = uint8(x)
			img.Pix[o+1] = uint8(255 - x)
			img.Pix[o+2] = uint8((x * 7) & 0xff)
			img.Pix[o+3] = uint8(y)
		}
	}
}

// visibleRows returns the pixel bytes of every row without the stride padding.
func visibleRows(img *image.RGBA) [][]byte {
	rows := make([][]byte, 0, img.Rect.Dy())
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		rows = append(rows, img.Pix[img.PixOffset(img.Rect.Min.X, y):img.PixOffset(img.Rect.Max.X, y)])
	}
	return rows
}

var _ = Describe("CompositeOnWhiteInPlace", func() {
	DescribeTable("matches the image/draw reference for every colour and alpha value",
		func(img *image.RGBA) {
			fillAllColorAlphaPairs(img)
			want := compositeOnWhiteReference(img)

			renderutil.CompositeOnWhiteInPlace(img)

			gotRows := visibleRows(img)
			wantRows := visibleRows(want)
			for y := range gotRows {
				Expect(gotRows[y]).To(Equal(wantRows[y]), "row %d (alpha %d) differs from the image/draw reference", y, y)
			}
		},
		Entry("with a tight stride", image.NewRGBA(image.Rect(0, 0, 256, 256))),
		Entry("with a padded stride", &image.RGBA{
			Pix:    make([]byte, 256*(256*4+32)),
			Stride: 256*4 + 32,
			Rect:   image.Rect(0, 0, 256, 256),
		}),
		Entry("with a non-zero origin", image.NewRGBA(image.Rect(10, 20, 266, 276))),
	)

	It("leaves the stride padding untouched", func() {
		img := &image.RGBA{
			Pix:    make([]byte, 4*(2*4+8)),
			Stride: 2*4 + 8,
			Rect:   image.Rect(0, 0, 2, 4),
		}
		for i := range img.Pix {
			img.Pix[i] = 0x42
		}

		renderutil.CompositeOnWhiteInPlace(img)

		for y := range 4 {
			pad := img.Pix[y*img.Stride+8 : (y+1)*img.Stride]
			for _, b := range pad {
				Expect(b).To(Equal(byte(0x42)), "padding of row %d was modified", y)
			}
		}
	})

	It("leaves fully opaque pixels unchanged", func() {
		img := image.NewRGBA(image.Rect(0, 0, 4, 1))
		for x := range 4 {
			img.SetRGBA(x, 0, color.RGBA{R: uint8(x * 60), G: 10, B: 200, A: 255})
		}
		before := append([]byte(nil), img.Pix...)

		renderutil.CompositeOnWhiteInPlace(img)

		Expect(img.Pix).To(Equal(before))
	})

	It("turns fully transparent pixels white", func() {
		img := image.NewRGBA(image.Rect(0, 0, 4, 1))
		for x := range 4 {
			img.SetRGBA(x, 0, color.RGBA{R: uint8(x * 60), G: 10, B: 200, A: 0})
		}

		renderutil.CompositeOnWhiteInPlace(img)

		for x := range 4 {
			Expect(img.RGBAAt(x, 0)).To(Equal(color.RGBA{R: 255, G: 255, B: 255, A: 255}))
		}
	})
})
