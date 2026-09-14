// Package renderutil contains the size calculations that are shared between
// the cgo and the webassembly implementation of the render methods.
package renderutil

import (
	"errors"
	"math"

	"github.com/klippa-app/go-pdfium/requests"
)

// Crop is the pixel representation of a requests.RenderPageCrop for a specific
// render scale.
type Crop struct {
	OffsetX int // The X offset of the region inside the full page render, in pixels.
	OffsetY int // The Y offset of the region inside the full page render, in pixels.
	Width   int // The width of the cropped image in pixels.
	Height  int // The height of the cropped image in pixels.
}

// ValidateCrop checks whether a crop region can be rendered at all. It is
// separate from CalculateCrop because a render that is given a maximum width or
// height derives its scale from the size of the region, so the region has to be
// checked before there is a scale to calculate the pixel bounds with.
func ValidateCrop(crop requests.RenderPageCrop) error {
	if !isFinite(crop.X) || !isFinite(crop.Y) || !isFinite(crop.Width) || !isFinite(crop.Height) {
		return errors.New("crop values must be valid numbers")
	}

	if crop.Width <= 0 || crop.Height <= 0 {
		return errors.New("crop width and height must be larger than 0")
	}

	return nil
}

// CalculateCrop calculates the pixel bounds of a crop region for a render with
// the given scale in pixels per point.
//
// The bounds are calculated as pixel edges on the grid of the full page render
// at the same scale, instead of by scaling the size of the region itself. That
// gives two guarantees that a tiling caller depends on: a cropped render
// contains the same pixels as the matching region of a full page render, and
// two regions that share an edge in points also share that edge in pixels, so
// tiles fit together without a seam, a gap or an overlap.
func CalculateCrop(crop requests.RenderPageCrop, scale float64) (*Crop, error) {
	if err := ValidateCrop(crop); err != nil {
		return nil, err
	}

	if !isFinite(scale) {
		return nil, errors.New("crop values must be valid numbers")
	}

	left := math.Round(crop.X * scale)
	top := math.Round(crop.Y * scale)
	right := math.Round((crop.X + crop.Width) * scale)
	bottom := math.Round((crop.Y + crop.Height) * scale)

	// PDFium uses 32 bit signed integers for all render coordinates.
	if !isInt32(left) || !isInt32(top) || !isInt32(right) || !isInt32(bottom) {
		return nil, errors.New("crop is too large to render")
	}

	width := int(right - left)
	height := int(bottom - top)
	if width < 1 || height < 1 {
		return nil, errors.New("crop is too small to render")
	}

	return &Crop{
		OffsetX: int(left),
		OffsetY: int(top),
		Width:   width,
		Height:  height,
	}, nil
}

// RenderSize returns the size in pixels of a full page render at the given
// scale in pixels per point. This is the size that a crop is positioned
// against, it is not the size of the resulting image when cropping.
func RenderSize(widthInPoints, heightInPoints, scale float64) (int, int, error) {
	width := math.Ceil(widthInPoints * scale)
	height := math.Ceil(heightInPoints * scale)

	if !isInt32(width) || !isInt32(height) {
		return 0, 0, errors.New("crop is too large to render")
	}

	return int(width), int(height), nil
}

// CalculateImageSize calculates the pixel size of a box of the given size in
// points when it has to fit inside the given maximum width and/or height.
// A maximum of 0 means that the value is calculated from the other maximum and
// the aspect ratio of the box. It returns the width, the height and the scale
// in pixels per point.
func CalculateImageSize(widthInPoints, heightInPoints float64, width, height int) (int, int, float64) {
	targetWidth := float64(width)
	targetHeight := float64(height)
	ratio := float64(0)
	if height == 0 {
		// Height not set, add ratio to height.
		ratio = heightInPoints / widthInPoints
		targetHeight = targetWidth * ratio
	} else if width == 0 {
		// Width not set, add ratio to width.
		ratio = widthInPoints / heightInPoints
		targetWidth = targetHeight * ratio
	} else {
		// Both values set, automatically pick the correct ratio.
		ratio = heightInPoints / widthInPoints
		if (targetWidth * ratio) < float64(height) {
			targetHeight = targetWidth * ratio
		} else {
			ratio = widthInPoints / heightInPoints
			if (targetHeight * ratio) < float64(width) {
				targetWidth = targetHeight * ratio
			}
		}
	}

	return int(math.Ceil(targetWidth)), int(math.Ceil(targetHeight)), targetWidth / widthInPoints
}

func isFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func isInt32(value float64) bool {
	return value >= math.MinInt32 && value <= math.MaxInt32
}
