package renderutil

import (
	"math"
	"testing"

	"github.com/klippa-app/go-pdfium/requests"
)

func TestCalculateCrop(t *testing.T) {
	tests := []struct {
		name    string
		crop    requests.RenderPageCrop
		scale   float64
		want    Crop
		wantErr string
	}{
		{
			name:  "a crop at the origin has no offset",
			crop:  requests.RenderPageCrop{X: 0, Y: 0, Width: 100, Height: 50},
			scale: 1,
			want:  Crop{OffsetX: 0, OffsetY: 0, Width: 100, Height: 50},
		},
		{
			name:  "a crop is scaled",
			crop:  requests.RenderPageCrop{X: 10, Y: 20, Width: 100, Height: 50},
			scale: 2,
			want:  Crop{OffsetX: 20, OffsetY: 40, Width: 200, Height: 100},
		},
		{
			name:  "a crop is measured on the pixel grid of the full render",
			crop:  requests.RenderPageCrop{X: 0.5, Y: 0.5, Width: 1, Height: 1},
			scale: 1,
			want:  Crop{OffsetX: 1, OffsetY: 1, Width: 1, Height: 1},
		},
		{
			name:  "a crop outside of the page is allowed",
			crop:  requests.RenderPageCrop{X: 10000, Y: 10000, Width: 100, Height: 50},
			scale: 1,
			want:  Crop{OffsetX: 10000, OffsetY: 10000, Width: 100, Height: 50},
		},
		{
			name:  "a crop before the page is allowed and has a negative offset",
			crop:  requests.RenderPageCrop{X: -50, Y: -25, Width: 100, Height: 50},
			scale: 1,
			want:  Crop{OffsetX: -50, OffsetY: -25, Width: 100, Height: 50},
		},
		{
			name:    "a crop without a width is rejected",
			crop:    requests.RenderPageCrop{X: 0, Y: 0, Width: 0, Height: 50},
			scale:   1,
			wantErr: "crop width and height must be larger than 0",
		},
		{
			name:    "a crop with a negative height is rejected",
			crop:    requests.RenderPageCrop{X: 0, Y: 0, Width: 100, Height: -50},
			scale:   1,
			wantErr: "crop width and height must be larger than 0",
		},
		{
			name:    "a crop that is not a number is rejected",
			crop:    requests.RenderPageCrop{X: math.NaN(), Y: 0, Width: 100, Height: 50},
			scale:   1,
			wantErr: "crop values must be valid numbers",
		},
		{
			name:    "an infinite crop is rejected",
			crop:    requests.RenderPageCrop{X: 0, Y: 0, Width: math.Inf(1), Height: 50},
			scale:   1,
			wantErr: "crop values must be valid numbers",
		},
		{
			name:    "a crop that rounds away to nothing is rejected",
			crop:    requests.RenderPageCrop{X: 0, Y: 0, Width: 0.1, Height: 0.1},
			scale:   1,
			wantErr: "crop is too small to render",
		},
		{
			name:    "a crop that does not fit in an int32 is rejected",
			crop:    requests.RenderPageCrop{X: 0, Y: 0, Width: 100, Height: 50},
			scale:   math.MaxInt32,
			wantErr: "crop is too large to render",
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			crop, err := CalculateCrop(tests[i].crop, tests[i].scale)
			if tests[i].wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %s but got no error", tests[i].wantErr)
				}
				if err.Error() != tests[i].wantErr {
					t.Fatalf("expected error %s but got error %s", tests[i].wantErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error but got error %s", err.Error())
			}

			if *crop != tests[i].want {
				t.Errorf("expected %+v but got %+v", tests[i].want, *crop)
			}
		})
	}
}

// TestCalculateCropTilesAreAdjacent is the guarantee that the whole pixel edge
// calculation exists for: neighbouring tiles have to line up exactly, so that
// rendering a page as tiles gives the same pixels as rendering it in one go.
func TestCalculateCropTilesAreAdjacent(t *testing.T) {
	// Deliberately awkward values, a tile size and a scale that do not divide
	// into whole pixels.
	const tileSize = 33.3
	const scale = 1.3888888888888888

	fullWidth := math.Ceil(tileSize * 3 * scale)

	totalWidth := 0
	previousRight := 0
	for tile := 0; tile < 3; tile++ {
		crop, err := CalculateCrop(requests.RenderPageCrop{
			X:      float64(tile) * tileSize,
			Y:      0,
			Width:  tileSize,
			Height: tileSize,
		}, scale)
		if err != nil {
			t.Fatalf("expected no error but got error %s", err.Error())
		}

		if crop.OffsetX != previousRight {
			t.Errorf("tile %d starts at %d but the previous tile ended at %d, tiles must not overlap or leave a gap", tile, crop.OffsetX, previousRight)
		}

		previousRight = crop.OffsetX + crop.Width
		totalWidth += crop.Width
	}

	// The tiles together have to cover exactly the same pixels as a single
	// render of the same area, give or take the ceil of the full page size.
	if math.Abs(float64(totalWidth)-fullWidth) > 1 {
		t.Errorf("the tiles are %d pixels wide together but the full render is %v pixels wide", totalWidth, fullWidth)
	}
}

// TestCalculateCropRoundTripsPixels covers the pdfium-cli crop-px flag, which
// converts a pixel region to points and relies on getting the exact pixel size
// it asked for back.
func TestCalculateCropRoundTripsPixels(t *testing.T) {
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
			if err != nil {
				t.Fatalf("expected no error but got error %s", err.Error())
			}

			if crop.Width != int(pixels[2]) || crop.Height != int(pixels[3]) {
				t.Errorf("at %d dpi a region of %vx%v pixels came back as %dx%d pixels", dpi, pixels[2], pixels[3], crop.Width, crop.Height)
			}

			if crop.OffsetX != int(pixels[0]) || crop.OffsetY != int(pixels[1]) {
				t.Errorf("at %d dpi a region at %v,%v pixels came back at %d,%d pixels", dpi, pixels[0], pixels[1], crop.OffsetX, crop.OffsetY)
			}
		}
	}
}

func TestRenderSize(t *testing.T) {
	width, height, err := RenderSize(595.2755737304688, 841.8897094726562, 1.3888888888888888)
	if err != nil {
		t.Fatalf("expected no error but got error %s", err.Error())
	}

	if width != 827 || height != 1170 {
		t.Errorf("expected 827x1170 but got %dx%d", width, height)
	}

	if _, _, err := RenderSize(595, 841, math.MaxInt32); err == nil {
		t.Errorf("expected an error for a render that does not fit in an int32")
	}
}

// TestCalculateImageSize pins down the behaviour that was moved here out of the
// two implementations, so that the non cropped render keeps the exact sizes it
// had before.
func TestCalculateImageSize(t *testing.T) {
	tests := []struct {
		name          string
		widthInPoints float64
		heightPoints  float64
		width         int
		height        int
		wantWidth     int
		wantHeight    int
		wantRatio     float64
	}{
		{"only a maximum width", 200, 100, 400, 0, 400, 200, 2},
		{"only a maximum height", 200, 100, 0, 400, 800, 400, 4},
		{"both maximums, width is the limit", 200, 100, 400, 400, 400, 200, 2},
		{"both maximums, height is the limit", 100, 200, 400, 400, 200, 400, 2},
		{"a square box in a square maximum", 100, 100, 300, 300, 300, 300, 3},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			width, height, ratio := CalculateImageSize(tests[i].widthInPoints, tests[i].heightPoints, tests[i].width, tests[i].height)
			if width != tests[i].wantWidth || height != tests[i].wantHeight {
				t.Errorf("expected %dx%d but got %dx%d", tests[i].wantWidth, tests[i].wantHeight, width, height)
			}

			if ratio != tests[i].wantRatio {
				t.Errorf("expected ratio %v but got %v", tests[i].wantRatio, ratio)
			}
		})
	}
}
