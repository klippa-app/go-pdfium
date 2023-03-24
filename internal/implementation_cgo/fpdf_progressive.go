package implementation_cgo

/*
#cgo pkg-config: pdfium
#include "fpdf_progressive.h"
#include <stdlib.h>

extern int go_progressive_render_pause_cb(struct _IFSDK_PAUSE *me);

static inline void IFSDK_PAUSE_SET_CB(IFSDK_PAUSE *p, char *id) {
	p->NeedToPauseNow = &go_progressive_render_pause_cb;
	p->user = id;
}
*/
import "C"
import (
	"errors"
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"unsafe"
)

// go_progressive_render_pause_cb is the Go implementation of IFSDK_PAUSE::NeedToPauseNow.
// It is exported through cgo so that we can use the reference to it and set
// it on IFSDK_PAUSE structs.
//
//export go_progressive_render_pause_cb
func go_progressive_render_pause_cb(me *C.IFSDK_PAUSE) C.FPDF_BOOL {
	pageRef := C.GoString((*C.char)(me.user))

	// Check if we still have the reference.
	if _, ok := pauseHandles[references.FPDF_PAGE(pageRef)]; !ok {
		return C.FPDF_BOOL(1)
	}

	shouldPause := pauseHandles[references.FPDF_PAGE(pageRef)].Callback()
	if shouldPause {
		return C.FPDF_BOOL(1)
	}

	return C.FPDF_BOOL(0)
}

type PauseHandle struct {
	Struct    *C.IFSDK_PAUSE
	Callback  func() bool
	stringRef unsafe.Pointer
}

var pauseHandles = map[references.FPDF_PAGE]*PauseHandle{}

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

	pauseStruct := &C.IFSDK_PAUSE{}
	pauseStruct.version = 1

	cPageRef := C.CString(string(pageHandle.nativeRef))

	C.IFSDK_PAUSE_SET_CB(pauseStruct, cPageRef)

	pauseHandle := &PauseHandle{
		stringRef: unsafe.Pointer(cPageRef),
		Struct:    pauseStruct,
		Callback:  request.NeedToPauseNowCallback,
	}

	pauseHandles[pageHandle.nativeRef] = pauseHandle

	renderStatus := C.FPDF_RenderPageBitmap_Start(bitmapHandle.handle, pageHandle.handle, C.int(request.StartX), C.int(request.StartY), C.int(request.SizeX), C.int(request.SizeY), C.int(request.Rotate), C.int(request.Flags), pauseStruct)

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
	if _, ok := pauseHandles[pageHandle.nativeRef]; ok {
		C.free(pauseHandles[pageHandle.nativeRef].stringRef)
		delete(pauseHandles, pageHandle.nativeRef)
	}

	var pauseStruct *C.IFSDK_PAUSE

	if request.NeedToPauseNowCallback != nil {
		pauseStruct = &C.IFSDK_PAUSE{}
		pauseStruct.version = 1

		cPageRef := C.CString(string(pageHandle.nativeRef))

		C.IFSDK_PAUSE_SET_CB(pauseStruct, cPageRef)

		pauseHandle := &PauseHandle{
			stringRef: unsafe.Pointer(cPageRef),
			Struct:    pauseStruct,
			Callback:  request.NeedToPauseNowCallback,
		}

		pauseHandles[pageHandle.nativeRef] = pauseHandle
	}

	renderStatus := C.FPDF_RenderPage_Continue(pageHandle.handle, pauseStruct)

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

	C.FPDF_RenderPage_Close(pageHandle.handle)

	// Check if we have the reference. Clean it up.
	if _, ok := pauseHandles[pageHandle.nativeRef]; ok {
		C.free(pauseHandles[pageHandle.nativeRef].stringRef)
		delete(pauseHandles, pageHandle.nativeRef)
	}

	return &responses.FPDF_RenderPage_Close{}, nil
}
