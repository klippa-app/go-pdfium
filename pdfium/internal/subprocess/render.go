package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_edit.h"
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"unsafe"

	"github.com/klippa-app/go-pdfium/pdfium/requests"
	"github.com/klippa-app/go-pdfium/pdfium/responses"
)

// getPageSize returns the points size of a page given the pdfium page index.
// One point is 1/72 inch (around 0.3528 mm).
func (p *Pdfium) getPageSize(page int) (float64, float64, error) {
	err := p.loadPage(page)
	if err != nil {
		return 0, 0, err
	}

	imgWidth := C.FPDF_GetPageWidth(p.currentPage)
	imgHeight := C.FPDF_GetPageHeight(p.currentPage)

	return float64(imgWidth), float64(imgHeight), nil
}

// getPageSizeInPixels returns the pixel size of a page given the page index and DPI.
func (p *Pdfium) getPageSizeInPixels(page, dpi int) (int, int, float64, error) {
	widthInPoints, heightInPoints, err := p.getPageSize(page)
	if err != nil {
		return 0, 0, 0, err
	}

	scale := float64(dpi) / 72.0

	return int(math.Ceil(widthInPoints * scale)), int(math.Ceil(heightInPoints * scale)), (widthInPoints * scale) / widthInPoints, nil
}

// GetPageSize returns the page size in points
// One point is 1/72 inch (around 0.3528 mm)
func (p *Pdfium) GetPageSize(request *requests.GetPageSize) (*responses.GetPageSize, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	widthInPoints, heightInPoints, err := p.getPageSize(request.Page)
	if err != nil {
		return nil, err
	}

	return &responses.GetPageSize{
		Page:   request.Page,
		Width:  widthInPoints,
		Height: heightInPoints,
	}, nil
}

// GetPageSizeInPixels returns the pixel size of a page given the page number and the DPI.
func (p *Pdfium) GetPageSizeInPixels(request *requests.GetPageSizeInPixels) (*responses.GetPageSizeInPixels, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	if request.DPI == 0 {
		return nil, errors.New("no DPI given")
	}

	widthInPixels, heightInPixels, pointToPixelRatio, err := p.getPageSizeInPixels(request.Page, request.DPI)
	if err != nil {
		return nil, err
	}

	return &responses.GetPageSizeInPixels{
		Page:              request.Page,
		Width:             widthInPixels,
		Height:            heightInPixels,
		PointToPixelRatio: pointToPixelRatio,
	}, nil
}

// RenderPageInDPI renders a specific page in a specific dpi, the result is an image.
func (p *Pdfium) RenderPageInDPI(request *requests.RenderPageInDPI) (*responses.RenderPage, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	if request.DPI == 0 {
		return nil, errors.New("no DPI given")
	}

	widthInPixels, heightInPixels, pointToPixelRatio, err := p.getPageSizeInPixels(request.Page, request.DPI)
	if err != nil {
		return nil, err
	}

	// Render a single page.
	result, err := p.renderPages([]renderPage{
		{
			Page:              request.Page,
			Width:             widthInPixels,
			Height:            heightInPixels,
			PointToPixelRatio: pointToPixelRatio,
		},
	}, 0)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPage{
		Page:              request.Page,
		Image:             result.Image,
		PointToPixelRatio: pointToPixelRatio,
		Width:             widthInPixels,
		Height:            heightInPixels,
	}, nil
}

// RenderPagesInDPI renders a list of pages in a specific dpi, the result is an image.
func (p *Pdfium) RenderPagesInDPI(request *requests.RenderPagesInDPI) (*responses.RenderPages, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	if len(request.Pages) == 0 {
		return nil, errors.New("no pages given")
	}

	pages := make([]renderPage, len(request.Pages))
	for i := range request.Pages {
		if request.Pages[i].DPI == 0 {
			return nil, fmt.Errorf("no DPI given for requested page %d", i)
		}

		widthInPixels, heightInPixels, pointToPixelRatio, err := p.getPageSizeInPixels(request.Pages[i].Page, request.Pages[i].DPI)
		if err != nil {
			return nil, err
		}

		pages[i] = renderPage{
			Page:              request.Pages[i].Page,
			Width:             widthInPixels,
			Height:            heightInPixels,
			PointToPixelRatio: pointToPixelRatio,
		}
	}

	return p.renderPages(pages, request.Padding)
}

func (p *Pdfium) calculateRenderImageSize(page, width, height int) (int, int, float64, error) {
	widthInPoints, heightInPoints, err := p.getPageSize(page)
	if err != nil {
		return 0, 0, 0, err
	}

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

	width = int(math.Ceil(targetWidth))
	height = int(math.Ceil(targetHeight))

	return width, height, targetWidth / widthInPoints, nil
}

// RenderPageInPixels renders a specific page in a specific pixel size, the result is an image.
// The given resolution is a maximum, we automatically calculate either the width or the height
// to make sure it stays withing the maximum resolution.
func (p *Pdfium) RenderPageInPixels(request *requests.RenderPageInPixels) (*responses.RenderPage, error) {
	p.Lock()
	defer p.Unlock()

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
			Page:              request.Page,
			Width:             width,
			Height:            height,
			PointToPixelRatio: ratio,
		},
	}, 0)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPage{
		Page:              request.Page,
		Image:             result.Image,
		PointToPixelRatio: ratio,
		Width:             width,
		Height:            height,
	}, nil
}

// RenderPagesInPixels renders a list of pages in a specific pixel size, the result is an image.
// The given resolution is a maximum, we automatically calculate either the width or the height
// to make sure it stays withing the maximum resolution.
func (p *Pdfium) RenderPagesInPixels(request *requests.RenderPagesInPixels) (*responses.RenderPages, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	if len(request.Pages) == 0 {
		return nil, errors.New("no pages given")
	}

	pages := make([]renderPage, len(request.Pages))
	for i := range request.Pages {
		if request.Pages[i].Width == 0 && request.Pages[i].Height == 0 {
			return nil, fmt.Errorf("no width or height given for requested page %d", i)
		}

		width, height, ratio, err := p.calculateRenderImageSize(request.Pages[i].Page, request.Pages[i].Width, request.Pages[i].Height)
		if err != nil {
			return nil, err
		}

		pages[i] = renderPage{
			Page:              request.Pages[i].Page,
			Width:             width,
			Height:            height,
			PointToPixelRatio: ratio,
		}
	}

	return p.renderPages(pages, request.Padding)
}

type renderPage struct {
	Page              int
	Width             int
	Height            int
	PointToPixelRatio float64
}

// renderPages renders a list of pages, the result is an image.
func (p *Pdfium) renderPages(pages []renderPage, padding int) (*responses.RenderPages, error) {
	totalWidth := 0
	totalHeight := 0

	// First calculate the total image size
	for i := range pages {
		if pages[i].Width > totalWidth {
			totalWidth = pages[i].Width
		}

		totalHeight += pages[i].Height

		// Add padding between the renders
		if i > 0 {
			totalHeight += padding
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))

	// Create a device independent bitmap to the external buffer by passing a
	// pointer to the first pixel, pdfium will do the rest.
	bitmap := C.FPDFBitmap_CreateEx(C.int(totalWidth), C.int(totalHeight), C.FPDFBitmap_BGRA, unsafe.Pointer(&img.Pix[0]), C.int(img.Stride))

	pagesInfo := make([]responses.RenderPagesPage, len(pages))
	currentOffset := 0
	for i := range pages {
		// Keep track of page information in the total image.
		pagesInfo[i] = responses.RenderPagesPage{
			Page:              pages[i].Page,
			PointToPixelRatio: pages[i].PointToPixelRatio,
			Width:             pages[i].Width,
			Height:            pages[i].Height,
			X:                 0,
			Y:                 currentOffset,
		}
		err := p.renderPage(bitmap, pages[i].Page, pages[i].Width, pages[i].Height, currentOffset)
		if err != nil {
			return nil, err
		}
		currentOffset += pages[i].Height + padding
	}

	// Release bitmap resources and buffers.
	// This does not clear the Go image pixel buffer.
	C.FPDFBitmap_Destroy(bitmap)

	return &responses.RenderPages{
		Image:  img,
		Pages:  pagesInfo,
		Width:  totalWidth,
		Height: totalHeight,
	}, nil
}

// renderPage renders a specific page in a specific size on a bitmap.
func (p *Pdfium) renderPage(bitmap C.FPDF_BITMAP, page, width, height, offset int) error {
	err := p.loadPage(page)
	if err != nil {
		return err
	}

	alpha := C.FPDFPage_HasTransparency(p.currentPage)

	// White
	fillColor := 0xFFFFFFFF

	// When the page has transparency, fill with black, not white.
	if int(alpha) == 1 {
		// Black
		fillColor = 0x00000000
	}

	// Fill the page rect with the specified color.
	C.FPDFBitmap_FillRect(bitmap, 0, C.int(offset), C.int(width), C.int(height), C.ulong(fillColor))

	// Render the bitmap into the given external bitmap, write the bytes
	// in reverse order so that BGRA becomes RGBA.
	C.FPDF_RenderPageBitmap(bitmap, p.currentPage, 0, C.int(offset), C.int(width), C.int(height), 0, C.FPDF_ANNOT|C.FPDF_REVERSE_BYTE_ORDER)

	return nil
}

func (p *Pdfium) RenderToFile(request *requests.RenderToFile) (*responses.RenderToFile, error) {
	var renderedImage *image.RGBA

	var myResp *responses.RenderToFile
	if request.RenderPageInDPI != nil {
		resp, err := p.RenderPageInDPI(request.RenderPageInDPI)
		if err != nil {
			return nil, err
		}

		renderedImage = resp.Image
		myResp = &responses.RenderToFile{
			Width:             resp.Width,
			Height:            resp.Height,
			PointToPixelRatio: resp.PointToPixelRatio,
			Pages: []responses.RenderPagesPage{
				{
					Page:              request.RenderPageInDPI.Page,
					PointToPixelRatio: resp.PointToPixelRatio,
					Width:             resp.Image.Bounds().Max.X,
					Height:            resp.Image.Bounds().Max.Y,
					X:                 0,
					Y:                 0,
				},
			},
		}
	} else if request.RenderPagesInDPI != nil {
		resp, err := p.RenderPagesInDPI(request.RenderPagesInDPI)
		if err != nil {
			return nil, err
		}

		renderedImage = resp.Image
		myResp = &responses.RenderToFile{
			Width:  resp.Width,
			Height: resp.Height,
			Pages:  resp.Pages,
		}
	} else if request.RenderPageInPixels != nil {
		resp, err := p.RenderPageInPixels(request.RenderPageInPixels)
		if err != nil {
			return nil, err
		}

		renderedImage = resp.Image
		myResp = &responses.RenderToFile{
			Width:             resp.Width,
			Height:            resp.Height,
			PointToPixelRatio: resp.PointToPixelRatio,
			Pages: []responses.RenderPagesPage{
				{
					Page:              request.RenderPageInPixels.Page,
					PointToPixelRatio: resp.PointToPixelRatio,
					Width:             resp.Image.Bounds().Max.X,
					Height:            resp.Image.Bounds().Max.Y,
					X:                 0,
					Y:                 0,
				},
			},
		}
	} else if request.RenderPagesInPixels != nil {
		resp, err := p.RenderPagesInPixels(request.RenderPagesInPixels)
		if err != nil {
			return nil, err
		}

		renderedImage = resp.Image
		myResp = &responses.RenderToFile{
			Width:  resp.Width,
			Height: resp.Height,
			Pages:  resp.Pages,
		}
	} else {
		return nil, errors.New("no render operation given")
	}

	var imgBuf bytes.Buffer

	if request.OutputFormat == requests.RenderToFileOutputFormatJPG {
		var opt jpeg.Options
		opt.Quality = 95

		for {
			err := jpeg.Encode(&imgBuf, renderedImage, &opt)
			if err != nil {
				return nil, err
			}

			if request.MaxFileSize == 0 || int64(imgBuf.Len()) < request.MaxFileSize {
				break
			}

			opt.Quality -= 10

			if opt.Quality <= 45 {
				return nil, errors.New("PDF image would exceed maximum filesize")
			}

			imgBuf.Reset()
		}
	} else if request.OutputFormat == requests.RenderToFileOutputFormatPNG {
		err := png.Encode(&imgBuf, renderedImage)
		if err != nil {
			return nil, err
		}

		if request.MaxFileSize != 0 && int64(imgBuf.Len()) > request.MaxFileSize {
			return nil, errors.New("PDF image would exceed maximum filesize")
		}
	} else {
		return nil, errors.New("invalid output format given")
	}

	if request.OutputTarget == requests.RenderToFileOutputTargetBytes {
		imageBytes := imgBuf.Bytes()
		myResp.ImageBytes = &imageBytes
	} else if request.OutputTarget == requests.RenderToFileOutputTargetFile {
		var targetFile *os.File
		if request.TargetFilePath != "" {
			existingFile, err := os.Create(request.TargetFilePath)
			if err != nil {
				return nil, err
			}
			targetFile = existingFile
		} else {
			tempFile, err := ioutil.TempFile("", "")
			if err != nil {
				return nil, err
			}
			targetFile = tempFile
		}

		_, err := targetFile.Write(imgBuf.Bytes())
		if err != nil {
			return nil, err
		}

		err = targetFile.Close()
		if err != nil {
			return nil, err
		}

		myResp.ImagePath = targetFile.Name()
	} else {
		return nil, errors.New("invalid output target given")
	}

	return myResp, nil
}
