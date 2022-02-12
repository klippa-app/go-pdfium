//go:build pdfium_experimental
// +build pdfium_experimental

package implementation

/*
#cgo pkg-config: pdfium
#include "fpdf_progressive.h"

extern int go_progressive_render_pause_cb(struct _IFSDK_PAUSE *me);

static inline void IFSDK_PAUSE_SET_CB(IFSDK_PAUSE *p, char *id) {
	p->NeedToPauseNow = &go_progressive_render_pause_cb;
	p->user = id;
}
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

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

	var colorScheme *C.FPDF_COLORSCHEME
	if request.ColorScheme != nil {
		colorScheme = &C.FPDF_COLORSCHEME{}
		colorScheme.path_fill_color = C.FPDF_DWORD(request.ColorScheme.PathFillColor)
		colorScheme.path_stroke_color = C.FPDF_DWORD(request.ColorScheme.PathStrokeColor)
		colorScheme.text_fill_color = C.FPDF_DWORD(request.ColorScheme.TextFillColor)
		colorScheme.text_stroke_color = C.FPDF_DWORD(request.ColorScheme.TextStrokeColor)
	}

	renderStatus := C.FPDF_RenderPageBitmapWithColorScheme_Start(bitmapHandle.handle, pageHandle.handle, C.int(request.StartX), C.int(request.StartY), C.int(request.SizeX), C.int(request.SizeY), C.int(request.Rotate), C.int(request.Flags), colorScheme, pauseStruct)

	return &responses.FPDF_RenderPageBitmapWithColorScheme_Start{
		RenderStatus: enums.FPDF_RENDER_STATUS(renderStatus),
	}, nil
}
