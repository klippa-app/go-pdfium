package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_edit.h"
import "C"

import (
	"errors"
	"fmt"
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

	// Render a single page.
	result, err := p.renderPages([]renderPage{
		{
			Page:      request.Page,
			PixelSize: *pixelSize,
		},
	}, 0)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPage{
		Image:             result.Image,
		PointToPixelRatio: pixelSize.PointToPixelRatio,
	}, nil
}

// RenderPagesInDPI renders a list of pages in a specific dpi, the result is an image.
func (p *Pdfium) RenderPagesInDPI(request *requests.RenderPagesInDPI) (*responses.RenderPages, error) {
	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	if len(request.Pages) == 0 {
		return nil, errors.New("no pages given")
	}

	pages := []renderPage{}
	for i := range request.Pages {
		if request.Pages[i].DPI == 0 {
			return nil, fmt.Errorf("no DPI given for requested page %d", i)
		}

		pixelSize, err := p.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
			Page: request.Pages[i].Page,
			DPI:  request.Pages[i].DPI,
		})
		if err != nil {
			return nil, err
		}

		pages = append(pages, renderPage{
			Page:      request.Pages[i].Page,
			PixelSize: *pixelSize,
		})
	}

	return p.renderPages(pages, request.Padding)
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

	// Render a single page.
	result, err := p.renderPages([]renderPage{
		{
			Page: request.Page,
			PixelSize: responses.GetPageSizeInPixels{
				Width:             width,
				Height:            height,
				PointToPixelRatio: ratio,
			},
		},
	}, 0)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPage{
		Image:             result.Image,
		PointToPixelRatio: ratio,
	}, nil
}

// RenderPagesInPixels renders a list of pages in a specific pixel size, the result is an image.
// The given resolution is a maximum, we automatically calculate either the width or the height
// to make sure it stays withing the maximum resolution.
func (p *Pdfium) RenderPagesInPixels(request *requests.RenderPagesInPixels) (*responses.RenderPages, error) {
	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	if len(request.Pages) == 0 {
		return nil, errors.New("no pages given")
	}

	pages := []renderPage{}
	for i := range request.Pages {
		if request.Pages[i].Width == 0 && request.Pages[i].Height == 0 {
			return nil, fmt.Errorf("no width or height given for requested page %d", i)
		}

		width, height, ratio, err := p.calculateRenderImageSize(request.Pages[i].Page, request.Pages[i].Width, request.Pages[i].Height)
		if err != nil {
			return nil, err
		}

		pages = append(pages, renderPage{
			Page: request.Pages[i].Page,
			PixelSize: responses.GetPageSizeInPixels{
				Width:             width,
				Height:            height,
				PointToPixelRatio: ratio,
			},
		})
	}

	return p.renderPages(pages, request.Padding)
}

type renderPage struct {
	Page      int
	PixelSize responses.GetPageSizeInPixels
}

// renderPages renders a list of pages, the result is an image.
func (p *Pdfium) renderPages(pages []renderPage, padding int) (*responses.RenderPages, error) {
	totalHeight := 0
	totalWidth := 0

	// First calculate the total image size
	for i := range pages {
		if pages[i].PixelSize.Width > totalWidth {
			totalWidth = pages[i].PixelSize.Width
		}

		totalHeight += pages[i].PixelSize.Height

		// Add padding between the renders
		if i > 0 {
			totalHeight += padding
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))

	// Create a device independent bitmap to the external buffer by passing a
	// pointer to the first pixel, pdfium will do the rest.
	p.Lock()
	bitmap := C.FPDFBitmap_CreateEx(C.int(totalWidth), C.int(totalHeight), C.FPDFBitmap_BGRA, unsafe.Pointer(&img.Pix[0]), C.int(img.Stride))
	p.Unlock()

	pagesInfo := []responses.RenderPagesPage{}

	currentOffset := 0
	for i := range pages {
		// Keep track of page information in the total image.
		pagesInfo = append(pagesInfo, responses.RenderPagesPage{
			PointToPixelRatio: pages[i].PixelSize.PointToPixelRatio,
			Width:             pages[i].PixelSize.Width,
			Height:            pages[i].PixelSize.Height,
			X:                 0,
			Y:                 currentOffset,
		})
		err := p.renderPage(bitmap, pages[i].Page, pages[i].PixelSize.Width, pages[i].PixelSize.Height, currentOffset)
		if err != nil {
			return nil, err
		}
		currentOffset += pages[i].PixelSize.Height + padding
	}

	// Release bitmap resources and buffers.
	// This does not clear the Go image pixel buffer.
	p.Lock()
	C.FPDFBitmap_Destroy(bitmap)
	p.Unlock()

	return &responses.RenderPages{
		Image: img,
		Pages: pagesInfo,
	}, nil
}

// renderPage renders a specific page in a specific size on a bitmap.
func (p *Pdfium) renderPage(bitmap C.FPDF_BITMAP, page, width, height, offset int) error {
	err := p.loadPage(page)
	if err != nil {
		return err
	}

	p.Lock()
	defer p.Unlock()

	// Check whether the page has transparency, this determines the fill color.
	alpha := C.FPDFPage_HasTransparency(p.currentPage)
	fillColor := 4294967295
	if int(alpha) == 1 {
		fillColor = 0
	}

	// Fill the rectangle with the color (transparent or white)
	C.FPDFBitmap_FillRect(bitmap, 0, C.int(offset), C.int(width), C.int(height), C.ulong(fillColor))

	// Render the bitmap into the given external bitmap, write the bytes
	// in reverse order so that BGRA becomes RGBA.
	C.FPDF_RenderPageBitmap(bitmap, p.currentPage, 0, C.int(offset), C.int(width), C.int(height), 0, C.FPDF_ANNOT|C.FPDF_REVERSE_BYTE_ORDER)

	return nil
}
