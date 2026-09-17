package renderutil

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/klippa-app/go-pdfium/requests"
)

var _ = Describe("CalculateCrop", func() {
	DescribeTable("calculates the crop",
		func(crop requests.RenderPageCrop, scale float64, want Crop) {
			got, err := CalculateCrop(crop, scale)
			Expect(err).To(BeNil())
			Expect(*got).To(Equal(want))
		},
		Entry("a crop at the origin has no offset",
			requests.RenderPageCrop{X: 0, Y: 0, Width: 100, Height: 50}, 1.0,
			Crop{OffsetX: 0, OffsetY: 0, Width: 100, Height: 50}),
		Entry("a crop is scaled",
			requests.RenderPageCrop{X: 10, Y: 20, Width: 100, Height: 50}, 2.0,
			Crop{OffsetX: 20, OffsetY: 40, Width: 200, Height: 100}),
		Entry("a crop is measured on the pixel grid of the full render",
			requests.RenderPageCrop{X: 0.5, Y: 0.5, Width: 1, Height: 1}, 1.0,
			Crop{OffsetX: 1, OffsetY: 1, Width: 1, Height: 1}),
		Entry("a crop outside of the page is allowed",
			requests.RenderPageCrop{X: 10000, Y: 10000, Width: 100, Height: 50}, 1.0,
			Crop{OffsetX: 10000, OffsetY: 10000, Width: 100, Height: 50}),
		Entry("a crop before the page is allowed and has a negative offset",
			requests.RenderPageCrop{X: -50, Y: -25, Width: 100, Height: 50}, 1.0,
			Crop{OffsetX: -50, OffsetY: -25, Width: 100, Height: 50}),
	)

	DescribeTable("rejects invalid crops",
		func(crop requests.RenderPageCrop, scale float64, wantErr string) {
			_, err := CalculateCrop(crop, scale)
			Expect(err).To(MatchError(wantErr))
		},
		Entry("a crop without a width is rejected",
			requests.RenderPageCrop{X: 0, Y: 0, Width: 0, Height: 50}, 1.0,
			"crop width and height must be larger than 0"),
		Entry("a crop with a negative height is rejected",
			requests.RenderPageCrop{X: 0, Y: 0, Width: 100, Height: -50}, 1.0,
			"crop width and height must be larger than 0"),
		Entry("a crop that is not a number is rejected",
			requests.RenderPageCrop{X: math.NaN(), Y: 0, Width: 100, Height: 50}, 1.0,
			"crop values must be valid numbers"),
		Entry("an infinite crop is rejected",
			requests.RenderPageCrop{X: 0, Y: 0, Width: math.Inf(1), Height: 50}, 1.0,
			"crop values must be valid numbers"),
		Entry("a crop that rounds away to nothing is rejected",
			requests.RenderPageCrop{X: 0, Y: 0, Width: 0.1, Height: 0.1}, 1.0,
			"crop is too small to render"),
		Entry("a crop that does not fit in an int32 is rejected",
			requests.RenderPageCrop{X: 0, Y: 0, Width: 100, Height: 50}, float64(math.MaxInt32),
			"crop is too large to render"),
	)

	// This is the guarantee that the whole pixel edge calculation exists for:
	// neighbouring tiles have to line up exactly, so that rendering a page as
	// tiles gives the same pixels as rendering it in one go.
	It("produces adjacent tiles", func() {
		// Deliberately awkward values, a tile size and a scale that do not
		// divide into whole pixels.
		const tileSize = 33.3
		const scale = 1.3888888888888888

		fullWidth := math.Ceil(tileSize * 3 * scale)

		totalWidth := 0
		previousRight := 0
		for tile := range 3 {
			crop, err := CalculateCrop(requests.RenderPageCrop{
				X:      float64(tile) * tileSize,
				Y:      0,
				Width:  tileSize,
				Height: tileSize,
			}, scale)
			Expect(err).To(BeNil())

			Expect(crop.OffsetX).To(Equal(previousRight), "tile %d starts at %d but the previous tile ended at %d, tiles must not overlap or leave a gap", tile, crop.OffsetX, previousRight)

			previousRight = crop.OffsetX + crop.Width
			totalWidth += crop.Width
		}

		// The tiles together have to cover exactly the same pixels as a single
		// render of the same area, give or take the ceil of the full page size.
		Expect(math.Abs(float64(totalWidth)-fullWidth)).To(BeNumerically("<=", 1), "the tiles are %d pixels wide together but the full render is %v pixels wide", totalWidth, fullWidth)
	})

	// Covers the pdfium-cli crop-px flag, which converts a pixel region to
	// points and relies on getting the exact pixel size it asked for back.
	It("round trips pixel regions", func() {
		for _, dpi := range []int{72, 96, 150, 200, 300, 400} {
			scale := float64(dpi) / 72.0

			for _, pixels := range [][4]float64{
				{0, 0, 500, 400},
				{1000, 500, 500, 400},
				{333, 777, 101, 97},
			} {
				crop, err := CalculateCrop(requests.RenderPageCrop{
					X:      pixels[0] / scale,
					Y:      pixels[1] / scale,
					Width:  pixels[2] / scale,
					Height: pixels[3] / scale,
				}, scale)
				Expect(err).To(BeNil())

				Expect([2]int{crop.Width, crop.Height}).To(Equal([2]int{int(pixels[2]), int(pixels[3])}), "at %d dpi a region of %vx%v pixels came back as %dx%d pixels", dpi, pixels[2], pixels[3], crop.Width, crop.Height)
				Expect([2]int{crop.OffsetX, crop.OffsetY}).To(Equal([2]int{int(pixels[0]), int(pixels[1])}), "at %d dpi a region at %v,%v pixels came back at %d,%d pixels", dpi, pixels[0], pixels[1], crop.OffsetX, crop.OffsetY)
			}
		}
	})
})

var _ = Describe("RenderSize", func() {
	It("calculates the size of a render", func() {
		width, height, err := RenderSize(595.2755737304688, 841.8897094726562, 1.3888888888888888)
		Expect(err).To(BeNil())
		Expect(width).To(Equal(827))
		Expect(height).To(Equal(1170))
	})

	It("rejects a render that does not fit in an int32", func() {
		_, _, err := RenderSize(595, 841, math.MaxInt32)
		Expect(err).ToNot(BeNil())
	})
})

// Pins down the behaviour that was moved here out of the two implementations,
// so that the non cropped render keeps the exact sizes it had before.
var _ = Describe("CalculateImageSize", func() {
	DescribeTable("calculates the image size",
		func(widthInPoints, heightInPoints float64, width, height, wantWidth, wantHeight int, wantRatio float64) {
			gotWidth, gotHeight, ratio := CalculateImageSize(widthInPoints, heightInPoints, width, height)
			Expect([2]int{gotWidth, gotHeight}).To(Equal([2]int{wantWidth, wantHeight}))
			Expect(ratio).To(Equal(wantRatio))
		},
		Entry("only a maximum width", 200.0, 100.0, 400, 0, 400, 200, 2.0),
		Entry("only a maximum height", 200.0, 100.0, 0, 400, 800, 400, 4.0),
		Entry("both maximums, width is the limit", 200.0, 100.0, 400, 400, 400, 200, 2.0),
		Entry("both maximums, height is the limit", 100.0, 200.0, 400, 400, 200, 400, 2.0),
		Entry("a square box in a square maximum", 100.0, 100.0, 300, 300, 300, 300, 3.0),
	)
})
