package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_edit.h"
import "C"

import (
	"errors"
	"image"
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

	err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

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

	if request.DPI == 0 {
		return nil, errors.New("no DPI given")
	}

	scale := float64(request.DPI) / 72.0
	pageSize, err := p.GetPageSize(&requests.GetPageSize{
		Page: request.Page,
	})

	if err != nil {
		return nil, err
	}

	return &responses.GetPageSizeInPixels{
		Width:             int(math.Ceil(pageSize.Width * scale)),
		Height:            int(math.Ceil(pageSize.Height * scale)),
		PointToPixelRatio: (pageSize.Width * scale) / pageSize.Width,
	}, nil
}

// RenderPageInDPI renders a specific page in a specific dpi, the result is an image.
func (p *Pdfium) RenderPageInDPI(request *requests.RenderPageInDPI) (*responses.RenderPage, error) {
	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	if request.DPI == 0 {
		return nil, errors.New("no DPI given")
	}

	pixelSize, err := p.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
		Page: request.Page,
		DPI:  request.DPI,
	})
	if err != nil {
		return nil, err
	}

	renderedPage, err := p.renderPage(request.Page, pixelSize.Width, pixelSize.Height)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPage{
		Image:             renderedPage,
		PointToPixelRatio: pixelSize.PointToPixelRatio,
	}, nil
}

func (p *Pdfium) calculateRenderImageSize(page, width, height int) (int, int, float64, error) {
	pageSize, err := p.GetPageSize(&requests.GetPageSize{
		Page: page,
	})

	if err != nil {
		return 0, 0, 0, err
	}

	targetWidth := float64(width)
	targetHeight := float64(height)
	ratio := float64(0)
	if height == 0 {
		// Height not set, add ratio to height.
		ratio = pageSize.Height / pageSize.Width
		targetHeight = targetWidth * ratio
	} else if width == 0 {
		// Width not set, add ratio to width.
		ratio = pageSize.Width / pageSize.Height
		targetWidth = targetHeight * ratio
	} else {
		// Both values set, automatically pick the correct ratio.
		ratio = pageSize.Height / pageSize.Width
		if (targetWidth * ratio) < float64(height) {
			targetHeight = targetWidth * ratio
		} else {
			ratio = pageSize.Width / pageSize.Height
			if (targetHeight * ratio) < float64(width) {
				targetWidth = targetHeight * ratio
			}
		}
	}

	width = int(math.Ceil(targetWidth))
	height = int(math.Ceil(targetHeight))

	return width, height, targetWidth / pageSize.Width, nil
}

// RenderPageInPixels renders a specific page in a specific pixel size, the result is an image.
// The given resolution is a maximum, we automatically calculate either the width or the height
// to make sure it stays withing the maximum resolution.
func (p *Pdfium) RenderPageInPixels(request *requests.RenderPageInPixels) (*responses.RenderPage, error) {
	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	if request.Width == 0 && request.Height == 0 {
		return nil, errors.New("no width or height given")
	}

	width, height, ratio, err := p.calculateRenderImageSize(request.Page, request.Width, request.Height)
	if err != nil {
		return nil, err
	}

	renderedPage, err := p.renderPage(request.Page, width, height)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPage{
		Image:             renderedPage,
		PointToPixelRatio: ratio,
	}, nil
}

// RenderPage renders a specific page in a specific dpi, the result is an image.
func (p *Pdfium) renderPage(page, width, height int) (*image.RGBA, error) {
	err := p.loadPage(page)
	if err != nil {
		return nil, err
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	p.Lock()
	defer p.Unlock()

	// Check whether the page has transparency, this determines the fill color.
	alpha := C.FPDFPage_HasTransparency(p.currentPage)
	fillColor := 4294967295
	if int(alpha) == 1 {
		fillColor = 0
	}

	// Create a device independent bitmap to the external buffer by passing a
	// pointer to the first pixel, pdfium will do the rest.
	bitmap := C.FPDFBitmap_CreateEx(C.int(width), C.int(height), C.FPDFBitmap_BGRA, unsafe.Pointer(&img.Pix[0]), C.int(img.Stride))

	// Fill the rectangle with the color (transparent or white)
	C.FPDFBitmap_FillRect(bitmap, 0, 0, C.int(width), C.int(height), C.ulong(fillColor))

	// Render the bitmap into the given external bitmap, write the bytes
	// in reverse order so that BGRA becomes RGBA.
	C.FPDF_RenderPageBitmap(bitmap, p.currentPage, 0, 0, C.int(width), C.int(height), 0, C.FPDF_ANNOT|C.FPDF_REVERSE_BYTE_ORDER)

	// Release bitmap resources and buffers.
	// This does not clear the Go image pixel buffer.
	C.FPDFBitmap_Destroy(bitmap)

	return img, nil
}
