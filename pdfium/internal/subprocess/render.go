package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_edit.h"
import "C"

import (
	"errors"
	"image"
	"image/color"
	"math"
	"unsafe"

	"github.com/klippa-app/go-pdfium/pdfium/requests"
	"github.com/klippa-app/go-pdfium/pdfium/responses"
)

// getPageSize returns the pixel size of a page given the pdfium page object DPI.
func (p *Pdfium) getPageSize() (float64, float64) {
	p.Lock()
	imgWidth := C.FPDF_GetPageWidth(p.currentPage)
	imgHeight := C.FPDF_GetPageHeight(p.currentPage)
	p.Unlock()

	return float64(imgWidth), float64(imgHeight)
}

// GetPageSize returns the page size in points
// One point is 1/72 inch (around 0.3528 mm)
func (p *Pdfium) GetPageSize(request *requests.GetPageSize) (*responses.GetPageSize, error) {
	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	p.loadPage(request.Page)
	widthInPoints, heightInPoints := p.getPageSize()
	return &responses.GetPageSize{
		Width:  widthInPoints,
		Height: heightInPoints,
	}, nil
}

// GetPageSizeInPixels returns the pixel size of a page given the page number and the DPI.
func (p *Pdfium) GetPageSizeInPixels(request *requests.GetPageSizeInPixels) (*responses.GetPageSizeInPixels, error) {
	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	scale := float64(request.DPI) / 72.0
	pageSize, err := p.GetPageSize(&requests.GetPageSize{
		Page: request.Page,
	})

	if err != nil {
		return nil, err
	}

	return &responses.GetPageSizeInPixels{
		Width:  int(math.Ceil(pageSize.Width * scale)),
		Height: int(math.Ceil(pageSize.Height * scale)),
	}, nil
}

// RenderPageInDPI renders a specific page in a specific dpi, the result is an image.
func (p *Pdfium) RenderPageInDPI(request *requests.RenderPageInDPI) (*responses.RenderPage, error) {
	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	pixelSize, err := p.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
		Page: request.Page,
		DPI:  request.DPI,
	})
	if err != nil {
		return nil, err
	}

	return &responses.RenderPage{
		Image: p.renderPage(request.Page, pixelSize.Width, pixelSize.Height),
	}, nil
}

// RenderPageInPixels renders a specific page in a specific pixel size, the result is an image.
// The given resolution is a maximum, we automatically calculate either the width or the height
// to make sure it stays withing the maximum resolution.
func (p *Pdfium) RenderPageInPixels(request *requests.RenderPageInPixels) (*responses.RenderPage, error) {
	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	pageSize, err := p.GetPageSize(&requests.GetPageSize{
		Page: request.Page,
	})

	if err != nil {
		return nil, err
	}

	targetWidth := float64(request.Width)
	targetHeight := float64(request.Height)
	if request.Height == 0 {
		// Height not set, add ratio to height.
		ratio := pageSize.Height / pageSize.Width
		targetHeight = targetWidth * ratio
	} else if request.Width == 0 {
		// Width not set, add ratio to width.
		ratio := pageSize.Width / pageSize.Height
		targetWidth = targetHeight * ratio
	} else {
		// Both values set, automatically pick the correct ratio.
		ratio := pageSize.Height / pageSize.Width
		if (targetWidth * ratio) < float64(request.Height) {
			targetHeight = targetWidth * ratio
		} else {
			ratio := pageSize.Width / pageSize.Height
			if (targetHeight * ratio) < float64(request.Width) {
				targetWidth = targetHeight * ratio
			}
		}
	}

	request.Width = int(math.Ceil(targetWidth))
	request.Height = int(math.Ceil(targetHeight))

	return &responses.RenderPage{
		Image: p.renderPage(request.Page, request.Width, request.Height),
	}, nil
}

// RenderPage renders a specific page in a specific dpi, the result is an image.
func (p *Pdfium) renderPage(page, width, height int) *image.RGBA {
	p.loadPage(page)
	p.Lock()
	alpha := C.FPDFPage_HasTransparency(p.currentPage)
	bitmap := C.FPDFBitmap_Create(C.int(width), C.int(height), alpha)

	fillColor := 4294967295
	if int(alpha) == 1 {
		fillColor = 0
	}

	C.FPDFBitmap_FillRect(bitmap, 0, 0, C.int(width), C.int(height), C.ulong(fillColor))
	C.FPDF_RenderPageBitmap(bitmap, p.currentPage, 0, 0, C.int(width), C.int(height), 0, C.FPDF_ANNOT)

	b := C.FPDFBitmap_GetBuffer(bitmap)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	img.Stride = int(C.FPDFBitmap_GetStride(bitmap))
	p.Unlock()

	// This takes a bit of time and I *think* we can do this without the lock
	// @todo: figure out if we can do this better/faster.
	bgra := make([]byte, 4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			for i := range bgra {
				bgra[i] = *((*byte)(b))
				b = unsafe.Pointer(uintptr(b) + 1)
			}
			pixelColor := color.RGBA{B: bgra[0], G: bgra[1], R: bgra[2], A: bgra[3]}
			img.SetRGBA(x, y, pixelColor)
		}
	}
	p.Lock()
	C.FPDFBitmap_Destroy(bitmap)
	p.Unlock()

	return img
}
