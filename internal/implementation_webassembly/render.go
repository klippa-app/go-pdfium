package implementation_webassembly

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/internal/image/image_jpeg"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// getPageSize returns the points size of a page given the PDFium page index.
// One point is 1/72 inch (around 0.3528 mm).
func (p *PdfiumImplementation) getPageSize(page requests.Page) (int, float64, float64, error) {
	pageHandle, err := p.loadPage(page)
	if err != nil {
		return 0, 0, 0, err
	}

	res, err := p.Module.ExportedFunction("FPDF_GetPageWidth").Call(p.Context, *pageHandle.handle)
	if err != nil {
		return 0, 0, 0, err
	}

	imgWidth := *(*float64)(unsafe.Pointer(&res[0]))

	res, err = p.Module.ExportedFunction("FPDF_GetPageHeight").Call(p.Context, *pageHandle.handle)
	if err != nil {
		return 0, 0, 0, err
	}

	imgHeight := *(*float64)(unsafe.Pointer(&res[0]))

	return pageHandle.index, float64(imgWidth), float64(imgHeight), nil
}

// getPageSizeInPixels returns the pixel size of a page given the page index and DPI.
func (p *PdfiumImplementation) getPageSizeInPixels(page requests.Page, dpi int) (int, int, int, float64, error) {
	index, widthInPoints, heightInPoints, err := p.getPageSize(page)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	scale := float64(dpi) / 72.0

	return index, int(math.Ceil(widthInPoints * scale)), int(math.Ceil(heightInPoints * scale)), (widthInPoints * scale) / widthInPoints, nil
}

// GetPageSize returns the page size in points
// One point is 1/72 inch (around 0.3528 mm)
func (p *PdfiumImplementation) GetPageSize(request *requests.GetPageSize) (*responses.GetPageSize, error) {
	p.Lock()
	defer p.Unlock()

	index, widthInPoints, heightInPoints, err := p.getPageSize(request.Page)
	if err != nil {
		return nil, err
	}

	return &responses.GetPageSize{
		Page:   index,
		Width:  widthInPoints,
		Height: heightInPoints,
	}, nil
}

// GetPageSizeInPixels returns the pixel size of a page given the page number and the DPI.
func (p *PdfiumImplementation) GetPageSizeInPixels(request *requests.GetPageSizeInPixels) (*responses.GetPageSizeInPixels, error) {
	p.Lock()
	defer p.Unlock()

	if request.DPI == 0 {
		return nil, errors.New("no DPI given")
	}

	index, widthInPixels, heightInPixels, pointToPixelRatio, err := p.getPageSizeInPixels(request.Page, request.DPI)
	if err != nil {
		return nil, err
	}

	return &responses.GetPageSizeInPixels{
		Page:              index,
		Width:             widthInPixels,
		Height:            heightInPixels,
		PointToPixelRatio: pointToPixelRatio,
	}, nil
}

// RenderPageInDPI renders a specific page in a specific dpi, the result is an image.
func (p *PdfiumImplementation) RenderPageInDPI(request *requests.RenderPageInDPI) (*responses.RenderPageInDPI, error) {
	p.Lock()
	defer p.Unlock()

	if request.DPI == 0 {
		return nil, errors.New("no DPI given")
	}

	index, widthInPixels, heightInPixels, pointToPixelRatio, err := p.getPageSizeInPixels(request.Page, request.DPI)
	if err != nil {
		return nil, err
	}

	// Render a single page.
	result, cleanupFunc, err := p.renderPages([]renderPage{
		{
			Page:              request.Page,
			Width:             widthInPixels,
			Height:            heightInPixels,
			PointToPixelRatio: pointToPixelRatio,
			Flags:             request.RenderFlags,
		},
	}, 0)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPageInDPI{
		CleanupFunc: cleanupFunc,
		Result: responses.RenderPage{
			Page:              index,
			Image:             result.Image,
			PointToPixelRatio: pointToPixelRatio,
			Width:             widthInPixels,
			Height:            heightInPixels,
		},
	}, nil
}

// RenderPagesInDPI renders a list of pages in a specific dpi, the result is an image.
func (p *PdfiumImplementation) RenderPagesInDPI(request *requests.RenderPagesInDPI) (*responses.RenderPagesInDPI, error) {
	p.Lock()
	defer p.Unlock()

	if len(request.Pages) == 0 {
		return nil, errors.New("no pages given")
	}

	pages := make([]renderPage, len(request.Pages))
	for i := range request.Pages {
		if request.Pages[i].DPI == 0 {
			return nil, fmt.Errorf("no DPI given for requested page %d", i)
		}

		_, widthInPixels, heightInPixels, pointToPixelRatio, err := p.getPageSizeInPixels(request.Pages[i].Page, request.Pages[i].DPI)
		if err != nil {
			return nil, err
		}

		pages[i] = renderPage{
			Page:              request.Pages[i].Page,
			Width:             widthInPixels,
			Height:            heightInPixels,
			PointToPixelRatio: pointToPixelRatio,
			Flags:             request.Pages[i].RenderFlags,
		}
	}

	result, cleanupFunc, err := p.renderPages(pages, request.Padding)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPagesInDPI{
		CleanupFunc: cleanupFunc,
		Result:      *result,
	}, nil
}

func (p *PdfiumImplementation) calculateRenderImageSize(page requests.Page, width, height int) (int, int, int, float64, error) {
	index, widthInPoints, heightInPoints, err := p.getPageSize(page)
	if err != nil {
		return 0, 0, 0, 0, err
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

	return index, width, height, targetWidth / widthInPoints, nil
}

// RenderPageInPixels renders a specific page in a specific pixel size, the result is an image.
// The given resolution is a maximum, we automatically calculate either the width or the height
// to make sure it stays withing the maximum resolution.
func (p *PdfiumImplementation) RenderPageInPixels(request *requests.RenderPageInPixels) (*responses.RenderPageInPixels, error) {
	p.Lock()
	defer p.Unlock()

	if request.Width == 0 && request.Height == 0 {
		return nil, errors.New("no width or height given")
	}

	index, width, height, ratio, err := p.calculateRenderImageSize(request.Page, request.Width, request.Height)
	if err != nil {
		return nil, err
	}

	// Render a single page.
	result, cleanupFunc, err := p.renderPages([]renderPage{
		{
			Page:              request.Page,
			Width:             width,
			Height:            height,
			PointToPixelRatio: ratio,
			Flags:             request.RenderFlags,
		},
	}, 0)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPageInPixels{
		CleanupFunc: cleanupFunc,
		Result: responses.RenderPage{
			Page:              index,
			Image:             result.Image,
			PointToPixelRatio: ratio,
			Width:             width,
			Height:            height,
		},
	}, nil
}

// RenderPagesInPixels renders a list of pages in a specific pixel size, the result is an image.
// The given resolution is a maximum, we automatically calculate either the width or the height
// to make sure it stays withing the maximum resolution.
func (p *PdfiumImplementation) RenderPagesInPixels(request *requests.RenderPagesInPixels) (*responses.RenderPagesInPixels, error) {
	p.Lock()
	defer p.Unlock()

	if len(request.Pages) == 0 {
		return nil, errors.New("no pages given")
	}

	pages := make([]renderPage, len(request.Pages))
	for i := range request.Pages {
		if request.Pages[i].Width == 0 && request.Pages[i].Height == 0 {
			return nil, fmt.Errorf("no width or height given for requested page %d", i)
		}

		_, width, height, ratio, err := p.calculateRenderImageSize(request.Pages[i].Page, request.Pages[i].Width, request.Pages[i].Height)
		if err != nil {
			return nil, err
		}

		pages[i] = renderPage{
			Page:              request.Pages[i].Page,
			Width:             width,
			Height:            height,
			PointToPixelRatio: ratio,
			Flags:             request.Pages[i].RenderFlags,
		}
	}

	result, cleanupFunc, err := p.renderPages(pages, request.Padding)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPagesInPixels{
		CleanupFunc: cleanupFunc,
		Result:      *result,
	}, nil
}

type renderPage struct {
	Page              requests.Page
	Flags             enums.FPDF_RENDER_FLAG
	Width             int
	Height            int
	PointToPixelRatio float64
}

// renderPages renders a list of pages, the result is an image.
func (p *PdfiumImplementation) renderPages(pages []renderPage, padding int) (*responses.RenderPages, func(), error) {
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

	// We use a "fake" image here, we will replace the Pix later.
	rect := image.Rect(0, 0, totalWidth, totalHeight)
	img := &image.RGBA{
		Pix:    nil,
		Stride: 4 * rect.Dx(),
		Rect:   rect,
	}

	size := img.Stride * totalHeight

	res, err := p.Module.ExportedFunction("FPDFBitmap_Create").Call(p.Context, uint64(totalWidth), uint64(totalHeight), uint64(1))
	if err != nil {
		return nil, nil, err
	}

	bitmap := res[0]

	releaseFunc := func() {
		// Release bitmap resources and buffers.
		p.Module.ExportedFunction("FPDFBitmap_Destroy").Call(p.Context, bitmap)
	}

	pagesInfo := make([]responses.RenderPagesPage, len(pages))
	currentOffset := 0
	for i := range pages {
		// Keep track of page information in the total image.
		pagesInfo[i] = responses.RenderPagesPage{
			PointToPixelRatio: pages[i].PointToPixelRatio,
			Width:             pages[i].Width,
			Height:            pages[i].Height,
			X:                 0,
			Y:                 currentOffset,
		}
		index, hasTransparency, err := p.renderPage(bitmap, pages[i].Page, pages[i].Width, pages[i].Height, currentOffset, pages[i].Flags)
		if err != nil {
			releaseFunc()
			return nil, nil, err
		}
		pagesInfo[i].Page = index
		pagesInfo[i].HasTransparency = hasTransparency
		currentOffset += pages[i].Height + padding
	}

	// The pointer to the first byte of the bitmap buffer.
	res, err = p.Module.ExportedFunction("FPDFBitmap_GetBuffer").Call(p.Context, bitmap)
	if err != nil {
		releaseFunc()
		return nil, nil, err
	}

	// Create a view of the underlying memory, not a copy.
	data, success := p.Module.Memory().Read(uint32(res[0]), uint32(size))
	if !success {
		releaseFunc()
		return nil, nil, errors.New("could not get bitmap buffer")
	}

	img.Pix = data

	return &responses.RenderPages{
		Image:  img,
		Pages:  pagesInfo,
		Width:  totalWidth,
		Height: totalHeight,
	}, releaseFunc, nil
}

// renderPage renders a specific page in a specific size on a bitmap.
func (p *PdfiumImplementation) renderPage(bitmap uint64, page requests.Page, width, height, offset int, flags enums.FPDF_RENDER_FLAG) (int, bool, error) {
	pageHandle, err := p.loadPage(page)
	if err != nil {
		return 0, false, err
	}

	res, err := p.Module.ExportedFunction("FPDFPage_HasTransparency").Call(p.Context, *pageHandle.handle)
	if err != nil {
		return 0, false, err
	}

	alpha := *(*int32)(unsafe.Pointer(&res[0]))

	// White
	fillColor := uint64(0xFFFFFFFF)

	hasTransparency := int(alpha) == 1

	// When the page has transparency, fill with black, not white.
	if hasTransparency {
		// Black
		fillColor = uint64(0x00000000)
	}

	// Fill the page rect with the specified color.
	_, err = p.Module.ExportedFunction("FPDFBitmap_FillRect").Call(p.Context, bitmap, uint64(0), uint64(offset), uint64(width), uint64(height), fillColor)
	if err != nil {
		return 0, false, err
	}

	// Render the bitmap into the given external bitmap, write the bytes
	// in reverse order so that BGRA becomes RGBA.
	flags = flags | enums.FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER
	_, err = p.Module.ExportedFunction("FPDF_RenderPageBitmap").Call(p.Context, bitmap, *pageHandle.handle, uint64(0), uint64(offset), uint64(width), uint64(height), uint64(0), *(*uint64)(unsafe.Pointer(&flags)))
	if err != nil {
		return 0, false, err
	}

	return pageHandle.index, hasTransparency, nil
}

func (p *PdfiumImplementation) RenderToFile(request *requests.RenderToFile) (*responses.RenderToFile, error) {
	var renderedImage *image.RGBA

	var myResp *responses.RenderToFile
	hasTransparency := false

	if request.RenderPageInDPI != nil {
		resp, err := p.RenderPageInDPI(request.RenderPageInDPI)
		if err != nil {
			return nil, err
		}
		defer resp.Cleanup()

		renderedImage = resp.Result.Image
		hasTransparency = resp.Result.HasTransparency
		myResp = &responses.RenderToFile{
			Width:             resp.Result.Width,
			Height:            resp.Result.Height,
			PointToPixelRatio: resp.Result.PointToPixelRatio,
			Pages: []responses.RenderPagesPage{
				{
					Page:              resp.Result.Page,
					PointToPixelRatio: resp.Result.PointToPixelRatio,
					Width:             resp.Result.Image.Bounds().Max.X,
					Height:            resp.Result.Image.Bounds().Max.Y,
					X:                 0,
					Y:                 0,
					HasTransparency:   resp.Result.HasTransparency,
				},
			},
		}
	} else if request.RenderPagesInDPI != nil {
		resp, err := p.RenderPagesInDPI(request.RenderPagesInDPI)
		if err != nil {
			return nil, err
		}
		defer resp.Cleanup()

		renderedImage = resp.Result.Image

		for _, page := range resp.Result.Pages {
			if page.HasTransparency {
				hasTransparency = true
			}
		}

		myResp = &responses.RenderToFile{
			Width:  resp.Result.Width,
			Height: resp.Result.Height,
			Pages:  resp.Result.Pages,
		}
	} else if request.RenderPageInPixels != nil {
		resp, err := p.RenderPageInPixels(request.RenderPageInPixels)
		if err != nil {
			return nil, err
		}
		defer resp.Cleanup()

		renderedImage = resp.Result.Image
		hasTransparency = resp.Result.HasTransparency
		myResp = &responses.RenderToFile{
			Width:             resp.Result.Width,
			Height:            resp.Result.Height,
			PointToPixelRatio: resp.Result.PointToPixelRatio,
			Pages: []responses.RenderPagesPage{
				{
					Page:              resp.Result.Page,
					PointToPixelRatio: resp.Result.PointToPixelRatio,
					Width:             resp.Result.Image.Bounds().Max.X,
					Height:            resp.Result.Image.Bounds().Max.Y,
					X:                 0,
					Y:                 0,
					HasTransparency:   resp.Result.HasTransparency,
				},
			},
		}
	} else if request.RenderPagesInPixels != nil {
		resp, err := p.RenderPagesInPixels(request.RenderPagesInPixels)
		if err != nil {
			return nil, err
		}
		defer resp.Cleanup()

		renderedImage = resp.Result.Image

		for _, page := range resp.Result.Pages {
			if page.HasTransparency {
				hasTransparency = true
			}
		}

		myResp = &responses.RenderToFile{
			Width:  resp.Result.Width,
			Height: resp.Result.Height,
			Pages:  resp.Result.Pages,
		}
	} else {
		return nil, errors.New("no render operation given")
	}

	var imgBuf bytes.Buffer

	if request.OutputFormat == requests.RenderToFileOutputFormatJPG {
		var opt jpeg.Options
		opt.Quality = 95

		if request.OutputQuality > 0 {
			opt.Quality = request.OutputQuality
		}

		// If any of the pages have transparency, place a white background under
		// the image. When you render a JPG image in Go, it will make the transparent
		// background black. With the added background we make sure that the
		// rendered PDF will look the same as in a PDF viewer, those generally
		// have a white background on the page viewer.
		if hasTransparency {
			imageWithWhiteBackground := image.NewRGBA(renderedImage.Bounds())
			draw.Draw(imageWithWhiteBackground, imageWithWhiteBackground.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
			draw.Draw(imageWithWhiteBackground, imageWithWhiteBackground.Bounds(), renderedImage, renderedImage.Bounds().Min, draw.Over)
			renderedImage = imageWithWhiteBackground
		}

		for {
			err := image_jpeg.Encode(&imgBuf, renderedImage, &opt)
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
