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
	"github.com/klippa-app/go-pdfium/internal/renderutil"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/tetratelabs/wazero/api"
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

	// When a crop is given we report the size of the region instead of the size
	// of the full page, using the same rounding that rendering the region uses.
	if request.Crop != nil {
		crop, err := renderutil.CalculateCrop(*request.Crop, pointToPixelRatio)
		if err != nil {
			return nil, err
		}

		widthInPixels = crop.Width
		heightInPixels = crop.Height
	}

	return &responses.GetPageSizeInPixels{
		Page:              index,
		Width:             widthInPixels,
		Height:            heightInPixels,
		PointToPixelRatio: pointToPixelRatio,
	}, nil
}

// applyCrop changes a render page to render only the given region instead of
// the full page. We still render the full page at the same scale, but we
// position it so that only the region lands inside the bitmap. PDFium clips the
// render to the bitmap, which leaves us with just the region.
func (p *PdfiumImplementation) applyCrop(pageToRender *renderPage, crop requests.RenderPageCrop, scale float64) error {
	_, widthInPoints, heightInPoints, err := p.getPageSize(pageToRender.Page)
	if err != nil {
		return err
	}

	calculatedCrop, err := renderutil.CalculateCrop(crop, scale)
	if err != nil {
		return err
	}

	renderWidth, renderHeight, err := renderutil.RenderSize(widthInPoints, heightInPoints, scale)
	if err != nil {
		return err
	}

	pageToRender.Width = calculatedCrop.Width
	pageToRender.Height = calculatedCrop.Height
	pageToRender.Crop = &renderCrop{
		RenderWidth:  renderWidth,
		RenderHeight: renderHeight,
		OffsetX:      calculatedCrop.OffsetX,
		OffsetY:      calculatedCrop.OffsetY,
	}

	return nil
}

// buildRenderPageInDPI builds the render page for a DPI based render.
func (p *PdfiumImplementation) buildRenderPageInDPI(request *requests.RenderPageInDPI) (int, *renderPage, error) {
	index, widthInPixels, heightInPixels, pointToPixelRatio, err := p.getPageSizeInPixels(request.Page, request.DPI)
	if err != nil {
		return 0, nil, err
	}

	pageToRender := &renderPage{
		Page:              request.Page,
		Width:             widthInPixels,
		Height:            heightInPixels,
		PointToPixelRatio: pointToPixelRatio,
		Flags:             request.RenderFlags,
		RenderForm:        request.RenderForm,
		Document:          request.Document,
	}

	if request.Crop != nil {
		if err := p.applyCrop(pageToRender, *request.Crop, pointToPixelRatio); err != nil {
			return 0, nil, err
		}
	}

	return index, pageToRender, nil
}

// buildRenderPageInPixels builds the render page for a render with a maximum
// width and/or height. When a crop is given, those maximums apply to the region
// instead of to the full page.
func (p *PdfiumImplementation) buildRenderPageInPixels(request *requests.RenderPageInPixels) (int, *renderPage, error) {
	index, widthInPoints, heightInPoints, err := p.getPageSize(request.Page)
	if err != nil {
		return 0, nil, err
	}

	if request.Crop != nil {
		// The scale is calculated from the size of the region, so the region
		// has to be valid before we can use it.
		if err := renderutil.ValidateCrop(*request.Crop); err != nil {
			return 0, nil, err
		}

		widthInPoints = request.Crop.Width
		heightInPoints = request.Crop.Height
	}

	width, height, ratio := renderutil.CalculateImageSize(widthInPoints, heightInPoints, request.Width, request.Height)

	pageToRender := &renderPage{
		Page:              request.Page,
		Width:             width,
		Height:            height,
		PointToPixelRatio: ratio,
		Flags:             request.RenderFlags,
		RenderForm:        request.RenderForm,
		Document:          request.Document,
	}

	if request.Crop != nil {
		if err := p.applyCrop(pageToRender, *request.Crop, ratio); err != nil {
			return 0, nil, err
		}
	}

	return index, pageToRender, nil
}


// RenderPageInDPI renders a specific page in a specific dpi, the result is an image.
func (p *PdfiumImplementation) RenderPageInDPI(request *requests.RenderPageInDPI) (*responses.RenderPageInDPI, error) {
	p.Lock()
	defer p.Unlock()

	if request.DPI == 0 {
		return nil, errors.New("no DPI given")
	}

	index, pageToRender, err := p.buildRenderPageInDPI(request)
	if err != nil {
		return nil, err
	}

	// Render a single page.
	result, cleanupFunc, err := p.renderPages([]renderPage{*pageToRender}, 0)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPageInDPI{
		CleanupFunc: cleanupFunc,
		Result: responses.RenderPage{
			Page:              index,
			Image:             result.Image,
			PointToPixelRatio: pageToRender.PointToPixelRatio,
			Width:             pageToRender.Width,
			Height:            pageToRender.Height,
			HasTransparency:   result.Pages[0].HasTransparency,
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

		if len(request.Pages) > 1 && request.Pages[i].Crop != nil {
			return nil, fmt.Errorf("crop is not supported for requested page %d when rendering multiple pages", i)
		}

		_, pageToRender, err := p.buildRenderPageInDPI(&request.Pages[i])
		if err != nil {
			return nil, err
		}

		pages[i] = *pageToRender
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

// calculateRenderImageSize calculates the pixel size of a page when it has to
// fit inside the given maximum width and/or height.
func (p *PdfiumImplementation) calculateRenderImageSize(page requests.Page, width, height int) (int, int, int, float64, error) {
	index, widthInPoints, heightInPoints, err := p.getPageSize(page)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	width, height, ratio := renderutil.CalculateImageSize(widthInPoints, heightInPoints, width, height)

	return index, width, height, ratio, nil
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

	index, pageToRender, err := p.buildRenderPageInPixels(request)
	if err != nil {
		return nil, err
	}

	// Render a single page.
	result, cleanupFunc, err := p.renderPages([]renderPage{*pageToRender}, 0)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPageInPixels{
		CleanupFunc: cleanupFunc,
		Result: responses.RenderPage{
			Page:              index,
			Image:             result.Image,
			PointToPixelRatio: pageToRender.PointToPixelRatio,
			Width:             pageToRender.Width,
			Height:            pageToRender.Height,
			HasTransparency:   result.Pages[0].HasTransparency,
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

		if len(request.Pages) > 1 && request.Pages[i].Crop != nil {
			return nil, fmt.Errorf("crop is not supported for requested page %d when rendering multiple pages", i)
		}

		_, pageToRender, err := p.buildRenderPageInPixels(&request.Pages[i])
		if err != nil {
			return nil, err
		}

		pages[i] = *pageToRender
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
	Width             int // The width of this page in the bitmap, the width of the region when cropping.
	Height            int // The height of this page in the bitmap, the height of the region when cropping.
	PointToPixelRatio float64
	RenderForm        bool
	Document          *references.FPDF_DOCUMENT
	Crop              *renderCrop // When given, only the region is rendered instead of the full page.
}

// renderCrop contains the values that are needed to render only a region of a
// page. We render the full page at the render scale, but positioned so that
// only the region lands inside the bitmap.
type renderCrop struct {
	RenderWidth  int // The width of the full page at the render scale.
	RenderHeight int // The height of the full page at the render scale.
	OffsetX      int // The X offset of the region inside the full page render.
	OffsetY      int // The Y offset of the region inside the full page render.
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

	if totalWidth < 1 || totalHeight < 1 {
		return nil, nil, errors.New("could not render an empty image")
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
		index, hasTransparency, err := p.renderPage(bitmap, pages[i], currentOffset)
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
func (p *PdfiumImplementation) renderPage(bitmap uint64, pageToRender renderPage, offset int) (int, bool, error) {
	pageHandle, err := p.loadPage(pageToRender.Page)
	if err != nil {
		return 0, false, err
	}

	width := pageToRender.Width
	height := pageToRender.Height
	flags := pageToRender.Flags

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

	// Fill the area of the bitmap that belongs to this page with the specified
	// color. This is always the area of the page in the bitmap, also when
	// cropping, so that the part of a region that falls outside of the page
	// keeps the background color.
	_, err = p.Module.ExportedFunction("FPDFBitmap_FillRect").Call(p.Context, bitmap, api.EncodeI32(0), api.EncodeI32(int32(offset)), api.EncodeI32(int32(width)), api.EncodeI32(int32(height)), fillColor)
	if err != nil {
		return 0, false, err
	}

	// By default we render the full page onto the area of the bitmap that
	// belongs to this page.
	startX := 0
	startY := offset
	sizeX := width
	sizeY := height

	// When cropping we render the full page at the same scale, but we move it
	// (partly) outside of the bitmap. PDFium clips the render to the bitmap,
	// which leaves us with just the region that we want.
	if pageToRender.Crop != nil {
		startX = -pageToRender.Crop.OffsetX
		startY = offset - pageToRender.Crop.OffsetY
		sizeX = pageToRender.Crop.RenderWidth
		sizeY = pageToRender.Crop.RenderHeight
	}

	// Render the bitmap into the given external bitmap, write the bytes
	// in reverse order so that BGRA becomes RGBA.
	flags = flags | enums.FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER
	_, err = p.Module.ExportedFunction("FPDF_RenderPageBitmap").Call(p.Context, bitmap, *pageHandle.handle, api.EncodeI32(int32(startX)), api.EncodeI32(int32(startY)), api.EncodeI32(int32(sizeX)), api.EncodeI32(int32(sizeY)), api.EncodeI32(0), api.EncodeI32(int32(flags)))
	if err != nil {
		return 0, false, err
	}

	if pageToRender.RenderForm {
		document := pageToRender.Document
		if document == nil && pageToRender.Page.ByIndex != nil {
			document = &pageToRender.Page.ByIndex.Document
		}
		if document == nil {
			return 0, false, errors.New("document is required when rendering forms")
		}

		documentHandle, err := p.getDocumentHandle(*document)
		if err != nil {
			return 0, false, err
		}

		res, err := p.Module.ExportedFunction("FPDF_FORMFILLINFO_Create").Call(p.Context)
		if err != nil {
			return 0, false, err
		}

		formInfoStruct := res[0]
		if formInfoStruct == 0 {
			return 0, false, errors.New("could not init form fill environment")
		}

		res, err = p.Module.ExportedFunction("FPDFDOC_InitFormFillEnvironment").Call(p.Context, *documentHandle.handle, formInfoStruct)
		if err != nil {
			return 0, false, err
		}

		formHandle := res[0]
		if formHandle == 0 {
			return 0, false, errors.New("could not init form fill environment")
		}

		// The form has to be drawn with the exact same position and size as the
		// page render, otherwise the form fields end up somewhere else than the
		// content of the page when cropping.
		_, err = p.Module.ExportedFunction("FPDF_FFLDraw").Call(p.Context, formHandle, bitmap, *pageHandle.handle, api.EncodeI32(int32(startX)), api.EncodeI32(int32(startY)), api.EncodeI32(int32(sizeX)), api.EncodeI32(int32(sizeY)), api.EncodeI32(0), api.EncodeI32(int32(flags)))
		if err != nil {
			return 0, false, err
		}

		_, err = p.Module.ExportedFunction("FPDFDOC_ExitFormFillEnvironment").Call(p.Context, formHandle)
		if err != nil {
			return 0, false, err
		}
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

	// If any of the pages have transparency, place a white background under
	// the image like a PDF viewer would. This is also to fix transparency JPEG
	// rendering, when you render a JPG image in Go, it will make the
	// transparent background black.
	if hasTransparency {
		imageWithWhiteBackground := image.NewRGBA(renderedImage.Bounds())
		draw.Draw(imageWithWhiteBackground, imageWithWhiteBackground.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
		// PDFium's FPDFBitmap_BGRA has straight (non-premultiplied) alpha.
		// Wrap as NRGBA so draw.Over uses the correct straight-alpha compositing formula.
		straightAlphaSrc := &image.NRGBA{
			Pix:    renderedImage.Pix,
			Stride: renderedImage.Stride,
			Rect:   renderedImage.Rect,
		}
		draw.Draw(imageWithWhiteBackground, imageWithWhiteBackground.Bounds(), straightAlphaSrc, straightAlphaSrc.Bounds().Min, draw.Over)
		renderedImage = imageWithWhiteBackground
	}

	if request.OutputFormat == requests.RenderToFileOutputFormatJPG {
		opt := image_jpeg.Options{
			Options: &jpeg.Options{
				Quality: 95,
			},
		}

		if request.OutputQuality > 0 {
			opt.Options.Quality = request.OutputQuality
		}

		for {
			err := image_jpeg.Encode(&imgBuf, renderedImage, opt)
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
