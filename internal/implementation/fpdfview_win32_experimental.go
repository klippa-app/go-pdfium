//go:build windows && pdfium_experimental
// +build windows,pdfium_experimental

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
