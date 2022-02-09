//go:build windows
// +build windows

package implementation

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"
import (
	"errors"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDF_SetPrintMode sets printing mode when printing on Windows.
// Experimental API.
// Windows only!
func (p *PdfiumImplementation) FPDF_SetPrintMode(request *requests.FPDF_SetPrintMode) (*responses.FPDF_SetPrintMode, error) {
	p.Lock()
	defer p.Unlock()

	success := C.FPDF_SetPrintMode(C.int(request.PrintMode))
	if int(success) == 0 {
		return errors.New("could not set print mode")
	}

	return &responses.FPDF_SetPrintMode{}, nil
}

// FPDF_RenderPage renders contents of a page to a device (screen, bitmap, or printer).
// Windows only!
func (p *PdfiumImplementation) FPDF_RenderPage(request *requests.FPDF_RenderPage) (*responses.FPDF_RenderPage, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	C.FPDF_RenderPage((C.HDC)(request.DC), pageHandle.handle, C.int(request.StartX), C.int(request.StartY), C.int(request.SizeX), C.int(request.SizeY), C.int(request.Rotate), C.int(request.Flags))

	return &responses.FPDF_RenderPageBitmap{}, nil
}
