package implementation_webassembly

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

	res, err := p.call("FPDF_GetPageWidth", *pageHandle.handle)
	if err != nil {
		return 0, 0, 0, err
	}

	imgWidth := *(*float64)(unsafe.Pointer(&res[0]))

	res, err = p.call("FPDF_GetPageHeight", *pageHandle.handle)
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
		ImageFormat:       request.ImageFormat,
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
		ImageFormat:       request.ImageFormat,
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

	rendered, err := p.renderPageInDPI(request)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPageInDPI{
		CleanupFunc: rendered.cleanup,
		Result:      rendered.singlePage(),
	}, nil
}

// The internal render methods return the rendered pages, which include the
// guest pixel offset and the cleanup function that releases the bitmap. The
// caller must hold the instance lock from before the call until after cleanup.
func (p *PdfiumImplementation) renderPageInDPI(request *requests.RenderPageInDPI) (*renderedPages, error) {
	if request.DPI == 0 {
		return nil, errors.New("no DPI given")
	}

	err := validateRenderImageFormat(request.ImageFormat)
	if err != nil {
		return nil, err
	}

	_, pageToRender, err := p.buildRenderPageInDPI(request)
	if err != nil {
		return nil, err
	}

	// Render a single page.
	return p.renderPages([]renderPage{*pageToRender}, 0)
}

// RenderPagesInDPI renders a list of pages in a specific dpi, the result is an image.
func (p *PdfiumImplementation) RenderPagesInDPI(request *requests.RenderPagesInDPI) (*responses.RenderPagesInDPI, error) {
	p.Lock()
	defer p.Unlock()

	rendered, err := p.renderPagesInDPI(request)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPagesInDPI{
		CleanupFunc: rendered.cleanup,
		Result:      rendered.RenderPages,
	}, nil
}

func (p *PdfiumImplementation) renderPagesInDPI(request *requests.RenderPagesInDPI) (*renderedPages, error) {
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

		err := validateRenderImageFormat(request.Pages[i].ImageFormat)
		if err != nil {
			return nil, fmt.Errorf("invalid ImageFormat given for requested page %d", i)
		}

		// All pages are rendered into one image, which can only have one
		// pixel format and one output field.
		if i > 0 && request.Pages[i].ImageFormat != pages[0].ImageFormat {
			return nil, errors.New("all pages must have the same ImageFormat when rendering multiple pages into one image")
		}

		_, pageToRender, err := p.buildRenderPageInDPI(&request.Pages[i])
		if err != nil {
			return nil, err
		}

		pages[i] = *pageToRender
	}

	return p.renderPages(pages, request.Padding)
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

	rendered, err := p.renderPageInPixels(request)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPageInPixels{
		CleanupFunc: rendered.cleanup,
		Result:      rendered.singlePage(),
	}, nil
}

func (p *PdfiumImplementation) renderPageInPixels(request *requests.RenderPageInPixels) (*renderedPages, error) {
	if request.Width == 0 && request.Height == 0 {
		return nil, errors.New("no width or height given")
	}

	err := validateRenderImageFormat(request.ImageFormat)
	if err != nil {
		return nil, err
	}

	_, pageToRender, err := p.buildRenderPageInPixels(request)
	if err != nil {
		return nil, err
	}

	// Render a single page.
	return p.renderPages([]renderPage{*pageToRender}, 0)
}

// RenderPagesInPixels renders a list of pages in a specific pixel size, the result is an image.
// The given resolution is a maximum, we automatically calculate either the width or the height
// to make sure it stays withing the maximum resolution.
func (p *PdfiumImplementation) RenderPagesInPixels(request *requests.RenderPagesInPixels) (*responses.RenderPagesInPixels, error) {
	p.Lock()
	defer p.Unlock()

	rendered, err := p.renderPagesInPixels(request)
	if err != nil {
		return nil, err
	}

	return &responses.RenderPagesInPixels{
		CleanupFunc: rendered.cleanup,
		Result:      rendered.RenderPages,
	}, nil
}

func (p *PdfiumImplementation) renderPagesInPixels(request *requests.RenderPagesInPixels) (*renderedPages, error) {
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

		err := validateRenderImageFormat(request.Pages[i].ImageFormat)
		if err != nil {
			return nil, fmt.Errorf("invalid ImageFormat given for requested page %d", i)
		}

		// All pages are rendered into one image, which can only have one
		// pixel format and one output field.
		if i > 0 && request.Pages[i].ImageFormat != pages[0].ImageFormat {
			return nil, errors.New("all pages must have the same ImageFormat when rendering multiple pages into one image")
		}

		_, pageToRender, err := p.buildRenderPageInPixels(&request.Pages[i])
		if err != nil {
			return nil, err
		}

		pages[i] = *pageToRender
	}

	return p.renderPages(pages, request.Padding)
}

type renderPage struct {
	Page              requests.Page
	Flags             enums.FPDF_RENDER_FLAG
	Width             int // The width of this page in the bitmap, the width of the region when cropping.
	Height            int // The height of this page in the bitmap, the height of the region when cropping.
	PointToPixelRatio float64
	RenderForm        bool
	Document          *references.FPDF_DOCUMENT
	Crop              *renderCrop                // When given, only the region is rendered instead of the full page.
	ImageFormat       requests.RenderImageFormat // The pixel format to render in.
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

// validateRenderImageFormat validates the given image format. An empty
// value is valid and renders as RGBA.
func validateRenderImageFormat(imageFormat requests.RenderImageFormat) error {
	switch imageFormat {
	case "", requests.RenderImageFormatRGBA, requests.RenderImageFormatGrayscale:
		return nil
	}

	return errors.New("invalid ImageFormat given")
}

// renderedPages retains the guest pixel offset alongside the public image.
// The offset stays valid across memory growth, until cleanup releases the bitmap.
type renderedPages struct {
	responses.RenderPages
	pixelsPtr uint64
	cleanup   func()
}

// singlePage returns the response of a render that contains exactly one page.
func (r *renderedPages) singlePage() responses.RenderPage {
	page := r.Pages[0]
	return responses.RenderPage{
		Page:              page.Page,
		Image:             r.Image,
		RenderedImage:     r.RenderedImage,
		PointToPixelRatio: page.PointToPixelRatio,
		Width:             page.Width,
		Height:            page.Height,
		HasTransparency:   page.HasTransparency,
	}
}

// hasTransparency reports whether any of the rendered pages has transparency.
func (r *renderedPages) hasTransparency() bool {
	for _, page := range r.Pages {
		if page.HasTransparency {
			return true
		}
	}

	return false
}

// toFileResponse builds the RenderToFile response for the rendered pages. The
// point to pixel ratio is only set for a render of a single page.
func (r *renderedPages) toFileResponse(singlePage bool) *responses.RenderToFile {
	resp := &responses.RenderToFile{
		Width:  r.Width,
		Height: r.Height,
		Pages:  r.Pages,
	}

	if singlePage {
		resp.PointToPixelRatio = r.Pages[0].PointToPixelRatio
	}

	return resp
}

// renderPages renders a list of pages, the result is an image.
func (p *PdfiumImplementation) renderPages(pages []renderPage, padding int) (*renderedPages, error) {
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
		return nil, errors.New("could not render an empty image")
	}

	// The image format has been validated by the caller, all pages are
	// guaranteed to have the same format here. An empty format renders as
	// RGBA.
	imageFormat := requests.RenderImageFormatRGBA
	if len(pages) > 0 && pages[0].ImageFormat != "" {
		imageFormat = pages[0].ImageFormat
	}

	// We use a "fake" image here, we will replace the Pix later.
	rect := image.Rect(0, 0, totalWidth, totalHeight)

	var img *image.RGBA
	var imgGray *image.Gray
	var bitmap uint64
	if imageFormat == requests.RenderImageFormatGrayscale {
		// The stride is calculated by PDFium and fetched with
		// FPDFBitmap_GetStride below, it may be larger than the width due to
		// alignment.
		imgGray = &image.Gray{
			Pix:  nil,
			Rect: rect,
		}

		// Pass a NULL buffer pointer so that PDFium allocates (and zero-fills)
		// the buffer inside the WebAssembly memory itself, it will be released
		// by FPDFBitmap_Destroy.
		res, err := p.call("FPDFBitmap_CreateEx", uint64(totalWidth), uint64(totalHeight), uint64(enums.FPDF_BITMAP_FORMAT_GRAY), 0, 0)
		if err != nil {
			return nil, err
		}

		bitmap = res[0]
	} else {
		img = &image.RGBA{
			Pix:    nil,
			Stride: 4 * rect.Dx(),
			Rect:   rect,
		}

		// PDFium runs in a 32 bit environment here, so a bitmap can never be
		// bigger than what a 32 bit integer can address. We have to check that
		// ourselves before asking PDFium for the bitmap, because the size wraps
		// around when we read the buffer back, which would leave us with a much
		// too small view of the bitmap instead of an error.
		if int64(img.Stride)*int64(totalHeight) > math.MaxUint32 {
			return nil, errors.New("the image to render is too large")
		}

		res, err := p.call("FPDFBitmap_Create", uint64(totalWidth), uint64(totalHeight), uint64(1))
		if err != nil {
			return nil, err
		}

		bitmap = res[0]
	}

	// PDFium returns a null bitmap when it could not allocate the buffer, which
	// in practice means that the instance ran out of WebAssembly memory or that
	// the dimensions overflow. Rendering into a null bitmap would silently give
	// us an image full of garbage.
	if bitmap == 0 {
		return nil, errors.New("could not create bitmap, the image to render is most likely too large")
	}

	releaseFunc := func() {
		// Release bitmap resources and buffers.
		p.call("FPDFBitmap_Destroy", bitmap)
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
		index, hasTransparency, err := p.renderPage(bitmap, pages[i], currentOffset, imageFormat)
		if err != nil {
			releaseFunc()
			return nil, err
		}
		pagesInfo[i].Page = index
		pagesInfo[i].HasTransparency = hasTransparency
		currentOffset += pages[i].Height + padding
	}

	imageSize := int64(0)
	if imgGray != nil {
		// The stride is decided by PDFium, it may be larger than the width
		// due to alignment.
		res, err := p.call("FPDFBitmap_GetStride", bitmap)
		if err != nil {
			releaseFunc()
			return nil, err
		}

		imgGray.Stride = int(*(*int32)(unsafe.Pointer(&res[0])))
		imageSize = int64(imgGray.Stride) * int64(totalHeight)
	} else {
		imageSize = int64(img.Stride) * int64(totalHeight)
	}

	// The same 32 bit reasoning as above, the grayscale stride only becomes
	// known here because PDFium decides it.
	if imageSize > math.MaxUint32 {
		releaseFunc()
		return nil, errors.New("the image to render is too large")
	}

	size := uint32(imageSize)

	// The pointer to the first byte of the bitmap buffer.
	res, err := p.call("FPDFBitmap_GetBuffer", bitmap)
	if err != nil {
		releaseFunc()
		return nil, err
	}

	// Create a view of the underlying memory, not a copy.
	data, success := p.Module.Memory().Read(uint32(res[0]), size)
	if !success {
		releaseFunc()
		return nil, errors.New("could not get bitmap buffer")
	}

	var renderedImage image.Image
	if imgGray != nil {
		imgGray.Pix = data
		renderedImage = imgGray
	} else {
		img.Pix = data
		renderedImage = img
	}

	return &renderedPages{
		RenderPages: responses.RenderPages{
			Image:         img,
			RenderedImage: renderedImage,
			Pages:         pagesInfo,
			Width:         totalWidth,
			Height:        totalHeight,
		},
		pixelsPtr: res[0],
		cleanup:   releaseFunc,
	}, nil
}

// renderPage renders a specific page in a specific size on a bitmap.
func (p *PdfiumImplementation) renderPage(bitmap uint64, pageToRender renderPage, offset int, imageFormat requests.RenderImageFormat) (int, bool, error) {
	pageHandle, err := p.loadPage(pageToRender.Page)
	if err != nil {
		return 0, false, err
	}

	width := pageToRender.Width
	height := pageToRender.Height
	flags := pageToRender.Flags

	res, err := p.call("FPDFPage_HasTransparency", *pageHandle.handle)
	if err != nil {
		return 0, false, err
	}

	alpha := *(*int32)(unsafe.Pointer(&res[0]))

	// White
	fillColor := uint64(0xFFFFFFFF)

	hasTransparency := int(alpha) == 1

	if imageFormat == requests.RenderImageFormatGrayscale {
		// A grayscale bitmap has no alpha channel, so the transparent black
		// fill can't be represented, always fill white like a PDF viewer.
		// Byte order is meaningless for a 1 byte per pixel format, so
		// FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER is not set here.
		flags = flags | enums.FPDF_RENDER_FLAG_GRAYSCALE
	} else {
		// When the page has transparency, fill with black, not white.
		if hasTransparency {
			// Black
			fillColor = uint64(0x00000000)
		}

		// Write the bytes in reverse order so that BGRA becomes RGBA.
		flags = flags | enums.FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER
	}

	// Fill the area of the bitmap that belongs to this page with the specified
	// color. This is always the area of the page in the bitmap, also when
	// cropping, so that the part of a region that falls outside of the page
	// keeps the background color.
	_, err = p.call("FPDFBitmap_FillRect", bitmap, api.EncodeI32(0), api.EncodeI32(int32(offset)), api.EncodeI32(int32(width)), api.EncodeI32(int32(height)), fillColor)
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

	// Render the bitmap into the given external bitmap.
	_, err = p.call("FPDF_RenderPageBitmap", bitmap, *pageHandle.handle, api.EncodeI32(int32(startX)), api.EncodeI32(int32(startY)), api.EncodeI32(int32(sizeX)), api.EncodeI32(int32(sizeY)), api.EncodeI32(0), api.EncodeI32(int32(flags)))
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

		res, err := p.call("FPDF_FORMFILLINFO_Create")
		if err != nil {
			return 0, false, err
		}

		formInfoStruct := res[0]
		if formInfoStruct == 0 {
			return 0, false, errors.New("could not init form fill environment")
		}

		res, err = p.call("FPDFDOC_InitFormFillEnvironment", *documentHandle.handle, formInfoStruct)
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
		_, err = p.call("FPDF_FFLDraw", formHandle, bitmap, *pageHandle.handle, api.EncodeI32(int32(startX)), api.EncodeI32(int32(startY)), api.EncodeI32(int32(sizeX)), api.EncodeI32(int32(sizeY)), api.EncodeI32(0), api.EncodeI32(int32(flags)))
		if err != nil {
			return 0, false, err
		}

		_, err = p.call("FPDFDOC_ExitFormFillEnvironment", formHandle)
		if err != nil {
			return 0, false, err
		}
	}

	return pageHandle.index, hasTransparency, nil
}

// renderForFile runs the render operation of a RenderToFile request and builds
// the response for it. The caller must hold the instance lock.
func (p *PdfiumImplementation) renderForFile(request *requests.RenderToFile) (*renderedPages, *responses.RenderToFile, error) {
	var rendered *renderedPages
	var err error
	singlePage := false

	switch {
	case request.RenderPageInDPI != nil:
		rendered, err = p.renderPageInDPI(request.RenderPageInDPI)
		singlePage = true
	case request.RenderPagesInDPI != nil:
		rendered, err = p.renderPagesInDPI(request.RenderPagesInDPI)
	case request.RenderPageInPixels != nil:
		rendered, err = p.renderPageInPixels(request.RenderPageInPixels)
		singlePage = true
	case request.RenderPagesInPixels != nil:
		rendered, err = p.renderPagesInPixels(request.RenderPagesInPixels)
	default:
		return nil, nil, errors.New("no render operation given")
	}

	if err != nil {
		return nil, nil, err
	}

	return rendered, rendered.toFileResponse(singlePage), nil
}

// renderAndEncode renders the pages of a RenderToFile request and encodes the
// image in the requested output format.
//
// The instance lock is held for the whole operation: the render helpers expect
// it, the JPEG encoder borrows the bitmap in guest memory, and the bitmap is
// released at the end. Writing the encoded image to its output target does not
// touch the instance, so RenderToFile does that after the lock is released.
// This keeps the lock free for other work, like the form fill timer callback
// which skips its tick when the instance is busy.
func (p *PdfiumImplementation) renderAndEncode(request *requests.RenderToFile) (*bytes.Buffer, *responses.RenderToFile, error) {
	p.Lock()
	defer p.Unlock()

	rendered, myResp, err := p.renderForFile(request)
	if err != nil {
		return nil, nil, err
	}
	defer rendered.cleanup()

	renderedImage := rendered.RenderedImage

	// If any of the pages have transparency, flatten the image onto a white
	// background like a PDF viewer would. This is also to fix transparency JPEG
	// rendering, when you render a JPG image in Go, it will make the
	// transparent background black.
	// The blend is done in place on the bitmap's pixel view, which is a live
	// window into WASM memory. That keeps the pixel offset valid so the JPEG
	// encoder can still borrow the bitmap instead of allocating and copying a
	// second full-size buffer. This is safe because no WASM call happens
	// between the render returning and this blend, so the view cannot have
	// gone stale.
	// Grayscale images have no alpha channel and are always rendered on a
	// white background, so they don't need this.
	if renderedImageRGBA, isRGBA := renderedImage.(*image.RGBA); rendered.hasTransparency() && isRGBA {
		renderutil.CompositeOnWhiteInPlace(renderedImageRGBA)
	}

	var imgBuf bytes.Buffer

	if request.OutputFormat == requests.RenderToFileOutputFormatJPG {
		opt := image_jpeg.Options{
			Options: &jpeg.Options{
				Quality: 95,
			},
			Progressive: request.Progressive,
		}

		if request.OutputQuality > 0 {
			opt.Options.Quality = request.OutputQuality
		}

		for {
			err := p.encodeJPEG(&imgBuf, renderedImage, rendered.pixelsPtr, opt)
			if err != nil {
				return nil, nil, err
			}

			if request.MaxFileSize == 0 || int64(imgBuf.Len()) < request.MaxFileSize {
				break
			}

			opt.Quality -= 10

			if opt.Quality <= 45 {
				return nil, nil, errors.New("PDF image would exceed maximum filesize")
			}

			imgBuf.Reset()
		}
	} else if request.OutputFormat == requests.RenderToFileOutputFormatPNG {
		// The zero value of PNGCompressionLevel is png.DefaultCompression, so
		// callers that don't set it keep the encoder's default behaviour.
		encoder := png.Encoder{CompressionLevel: request.PNGCompressionLevel}

		err := encoder.Encode(&imgBuf, renderedImage)
		if err != nil {
			return nil, nil, err
		}

		if request.MaxFileSize != 0 && int64(imgBuf.Len()) > request.MaxFileSize {
			return nil, nil, errors.New("PDF image would exceed maximum filesize")
		}
	} else {
		return nil, nil, errors.New("invalid output format given")
	}

	return &imgBuf, myResp, nil
}

func (p *PdfiumImplementation) RenderToFile(request *requests.RenderToFile) (*responses.RenderToFile, error) {
	imgBuf, myResp, err := p.renderAndEncode(request)
	if err != nil {
		return nil, err
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
