package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_edit.h"
import "C"

import (
	"image"
	"image/color"
	"math"
	"unsafe"
)

// getPageSize returns the pixel size of a page given the pdfium page object DPI.
func (d *Document) getPageSize() (float64, float64) {
	mutex.Lock()
	imgWidth := C.FPDF_GetPageWidth(d.page)
	imgHeight := C.FPDF_GetPageHeight(d.page)
	mutex.Unlock()

	width := C.double(imgWidth)
	height := C.double(imgHeight)

	return float64(width), float64(height)
}

// GetPageSize returns the page size in points
// One point is 1/72 inch (around 0.3528 mm)
func (d *Document) GetPageSize(page int) (float64, float64) {
	d.loadPage(page)
	return d.getPageSize()
}

// GetPageSizeInPixels returns the pixel size of a page given the page number and the DPI.
func (d *Document) GetPageSizeInPixels(page int, dpi int) (int, int) {
	scale := float64(dpi) / 72.0
	widthInPoints, heightInPoints := d.GetPageSize(page)

	return int(math.Ceil(widthInPoints * scale)), int(math.Ceil(heightInPoints * scale))
}

// renderPageInDPI renders a specific page in a specific dpi, the result is an image.
func (d *Document) renderPageInDPI(page, dpi int) *image.RGBA {
	width, height := d.GetPageSizeInPixels(page, dpi)
	return d.renderPage(page, width, height)
}

// renderPageInPixels renders a specific page in a specific pixel size, the result is an image.
// The given resolution is a maximum, we automatically calculate either the width or the height
// to make sure it stays withing the maximum resolution.
func (d *Document) renderPageInPixels(page, width, height int) *image.RGBA {
	widthInPoints, heightInPoints := d.GetPageSize(page)

	targetWidth := float64(width)
	targetHeight := float64(height)
	if height == 0 {
		// Height not set, add ratio to height.
		ratio := heightInPoints / widthInPoints
		targetHeight = targetWidth * ratio
	} else if width == 0 {
		// Width not set, add ratio to width.
		ratio := widthInPoints / heightInPoints
		targetWidth = targetHeight * ratio
	} else {
		// Both values set, automatically pick the correct ratio.
		ratio := heightInPoints / widthInPoints
		if (targetWidth * ratio) < float64(height) {
			targetHeight = targetWidth * ratio
		} else {
			ratio := widthInPoints / heightInPoints
			if (targetHeight * ratio) < float64(width) {
				targetWidth = targetHeight * ratio
			}
		}
	}

	width = int(math.Ceil(targetWidth))
	height = int(math.Ceil(targetHeight))

	return d.renderPage(page, width, height)
}

// RenderPage renders a specific page in a specific dpi, the result is an image.
func (d *Document) renderPage(page, width, height int) *image.RGBA {
	d.loadPage(page)
	mutex.Lock()
	alpha := C.FPDFPage_HasTransparency(d.page)
	bitmap := C.FPDFBitmap_Create(C.int(width), C.int(height), alpha)

	fillColor := 4294967295
	if int(alpha) == 1 {
		fillColor = 0
	}

	C.FPDFBitmap_FillRect(bitmap, 0, 0, C.int(width), C.int(height), C.ulong(fillColor))
	C.FPDF_RenderPageBitmap(bitmap, d.page, 0, 0, C.int(width), C.int(height), 0, C.FPDF_ANNOT)

	p := C.FPDFBitmap_GetBuffer(bitmap)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	img.Stride = int(C.FPDFBitmap_GetStride(bitmap))
	mutex.Unlock()

	// This takes a bit of time and I *think* we can do this without the lock
	// @todo: figure out if we can do this better/faster.
	bgra := make([]byte, 4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			for i := range bgra {
				bgra[i] = *((*byte)(p))
				p = unsafe.Pointer(uintptr(p) + 1)
			}
			pixelColor := color.RGBA{B: bgra[0], G: bgra[1], R: bgra[2], A: bgra[3]}
			img.SetRGBA(x, y, pixelColor)
		}
	}
	mutex.Lock()
	C.FPDFBitmap_Destroy(bitmap)
	mutex.Unlock()

	return img
}
