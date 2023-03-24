package implementation_webassembly

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/tetratelabs/wazero/api"
)

type PauseHandle struct {
	StringRef uint64
	Pointer   uint64
	Callback  func() bool
}

var PauseHandles = map[references.FPDF_PAGE]*PauseHandle{}

// FPDF_RenderPageBitmap_Start starts to render page contents to a device independent bitmap progressively.
// Not supported on multi-threaded usage.
func (p *PdfiumImplementation) FPDF_RenderPageBitmap_Start(request *requests.FPDF_RenderPageBitmap_Start) (*responses.FPDF_RenderPageBitmap_Start, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	bitmapHandle, err := p.getBitmapHandle(request.Bitmap)
	if err != nil {
		return nil, err
	}

	if request.NeedToPauseNowCallback == nil {
		return nil, errors.New("NeedToPauseNowCallback can't be nil")
	}

	refPointer, err := p.CString(string(pageHandle.nativeRef))
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("IFSDK_PAUSE_Create").Call(p.Context, refPointer.Pointer)
	if err != nil {
		return nil, err
	}

	pausePointer := res[0]

	pauseHandle := &PauseHandle{
		StringRef: refPointer.Pointer,
		Pointer:   pausePointer,
		Callback:  request.NeedToPauseNowCallback,
	}

	PauseHandles[pageHandle.nativeRef] = pauseHandle

	res, err = p.Module.ExportedFunction("FPDF_RenderPageBitmap_Start").Call(p.Context, *bitmapHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.StartX)), *(*uint64)(unsafe.Pointer(&request.StartY)), *(*uint64)(unsafe.Pointer(&request.SizeX)), *(*uint64)(unsafe.Pointer(&request.SizeY)), *(*uint64)(unsafe.Pointer(&request.Rotate)), *(*uint64)(unsafe.Pointer(&request.Flags)), pausePointer)
	if err != nil {
		return nil, err
	}

	renderStatus := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDF_RenderPageBitmap_Start{
		RenderStatus: enums.FPDF_RENDER_STATUS(renderStatus),
	}, nil
}

// FPDF_RenderPage_Continue continues rendering a PDF page.
// Not supported on multi-threaded usage.
func (p *PdfiumImplementation) FPDF_RenderPage_Continue(request *requests.FPDF_RenderPage_Continue) (*responses.FPDF_RenderPage_Continue, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	// Check if we already have the reference. Clean it up.
	if _, ok := PauseHandles[pageHandle.nativeRef]; ok {
		p.Free(PauseHandles[pageHandle.nativeRef].Pointer)
		p.Free(PauseHandles[pageHandle.nativeRef].StringRef)
		delete(PauseHandles, pageHandle.nativeRef)
	}

	pausePointer := uint64(0)
	if request.NeedToPauseNowCallback != nil {
		refPointer, err := p.CString(string(pageHandle.nativeRef))
		if err != nil {
			return nil, err
		}

		res, err := p.Module.ExportedFunction("IFSDK_PAUSE_Create").Call(p.Context, refPointer.Pointer)
		if err != nil {
			return nil, err
		}

		newPausePointer := res[0]
		pauseHandle := &PauseHandle{
			StringRef: refPointer.Pointer,
			Pointer:   newPausePointer,
			Callback:  request.NeedToPauseNowCallback,
		}

		PauseHandles[pageHandle.nativeRef] = pauseHandle
		pausePointer = newPausePointer
	}

	res, err := p.Module.ExportedFunction("FPDF_RenderPage_Continue").Call(p.Context, *pageHandle.handle, pausePointer)
	if err != nil {
		return nil, err
	}

	renderStatus := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDF_RenderPage_Continue{
		RenderStatus: enums.FPDF_RENDER_STATUS(renderStatus),
	}, nil
}

// FPDF_RenderPage_Close Release the resource allocate during page rendering. Need to be called after finishing rendering or cancel the rendering.
// Not supported on multi-threaded usage.
func (p *PdfiumImplementation) FPDF_RenderPage_Close(request *requests.FPDF_RenderPage_Close) (*responses.FPDF_RenderPage_Close, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDF_RenderPage_Close").Call(p.Context, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	// Check if we have the reference. Clean it up.
	if _, ok := PauseHandles[pageHandle.nativeRef]; ok {
		p.Free(PauseHandles[pageHandle.nativeRef].StringRef)
		p.Free(PauseHandles[pageHandle.nativeRef].Pointer)
		delete(PauseHandles, pageHandle.nativeRef)
	}

	return &responses.FPDF_RenderPage_Close{}, nil
}

// FPDF_RenderPageBitmapWithColorScheme_Start starts to render page contents to a device independent bitmap progressively with a specified color scheme for the content.
// Not supported on multi-threaded usage.
// Experimental API.
func (p *PdfiumImplementation) FPDF_RenderPageBitmapWithColorScheme_Start(request *requests.FPDF_RenderPageBitmapWithColorScheme_Start) (*responses.FPDF_RenderPageBitmapWithColorScheme_Start, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	bitmapHandle, err := p.getBitmapHandle(request.Bitmap)
	if err != nil {
		return nil, err
	}

	if request.NeedToPauseNowCallback == nil {
		return nil, errors.New("NeedToPauseNowCallback can't be nil")
	}

	refPointer, err := p.CString(string(pageHandle.nativeRef))
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("IFSDK_PAUSE_Create").Call(p.Context, refPointer.Pointer)
	if err != nil {
		return nil, err
	}

	pausePointer := res[0]
	pauseHandle := &PauseHandle{
		StringRef: refPointer.Pointer,
		Pointer:   pausePointer,
		Callback:  request.NeedToPauseNowCallback,
	}

	PauseHandles[pageHandle.nativeRef] = pauseHandle

	colorSchemeSize := p.CSizeULong() * 4
	colorScheme := uint64(0)
	if request.ColorScheme != nil {
		colorSchemePointer, err := p.Malloc(colorSchemeSize)
		if err != nil {
			return nil, err
		}

		colorScheme = colorSchemePointer

		p.Module.Memory().WriteUint64Le(uint32(colorScheme), api.EncodeU32(uint32(request.ColorScheme.PathFillColor)))
		p.Module.Memory().WriteUint64Le(uint32(colorScheme+4), api.EncodeU32(uint32(request.ColorScheme.PathStrokeColor)))
		p.Module.Memory().WriteUint64Le(uint32(colorScheme+8), api.EncodeU32(uint32(request.ColorScheme.TextFillColor)))
		p.Module.Memory().WriteUint64Le(uint32(colorScheme+12), api.EncodeU32(uint32(request.ColorScheme.TextStrokeColor)))
	}

	res, err = p.Module.ExportedFunction("FPDF_RenderPageBitmapWithColorScheme_Start").Call(p.Context, *bitmapHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.StartX)), *(*uint64)(unsafe.Pointer(&request.StartY)), *(*uint64)(unsafe.Pointer(&request.SizeX)), *(*uint64)(unsafe.Pointer(&request.SizeY)), *(*uint64)(unsafe.Pointer(&request.Rotate)), *(*uint64)(unsafe.Pointer(&request.Flags)), colorScheme, pausePointer)
	if err != nil {
		return nil, err
	}

	renderStatus := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDF_RenderPageBitmapWithColorScheme_Start{
		RenderStatus: enums.FPDF_RENDER_STATUS(renderStatus),
	}, nil
}
