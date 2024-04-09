//go:build pdfium_xfa
// +build pdfium_xfa

package implementation_cgo

// #cgo pkg-config: pdfium
// #cgo CFLAGS: -DPDF_ENABLE_XFA
// #include "fpdfview.h"
import "C"
import (
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"log"
)

// FPDF_BStr_Init initializes a FPDF_BSTR.
//
// Only available when using a build that includes XFA and when using the
// build flag pdfium_xfa.
func (p *PdfiumImplementation) FPDF_BStr_Init(request *requests.FPDF_BStr_Init) (*responses.FPDF_BStr_Init, error) {
	p.Lock()
	defer p.Unlock()

	FPDF_BStr := C.FPDF_BSTR{}
	log.Println(FPDF_BStr)
	C.FPDF_BStr_Init(&FPDF_BStr)
	log.Println(FPDF_BStr)

	// @todo: implement me.
	return nil, pdfium_errors.ErrXFAUnsupported
}

// FPDF_BStr_Set copies string data into the FPDF_BSTR.
//
// Only available when using a build that includes XFA and when using the
// build flag pdfium_xfa.
func (p *PdfiumImplementation) FPDF_BStr_Set(request *requests.FPDF_BStr_Set) (*responses.FPDF_BStr_Set, error) {
	p.Lock()
	defer p.Unlock()

	// @todo: implement me.
	return nil, pdfium_errors.ErrXFAUnsupported
}

// FPDF_BStr_Clear clears a FPDF_BSTR.
//
// Only available when using a build that includes XFA and when using the
// build flag pdfium_xfa.
func (p *PdfiumImplementation) FPDF_BStr_Clear(request *requests.FPDF_BStr_Clear) (*responses.FPDF_BStr_Clear, error) {
	p.Lock()
	defer p.Unlock()

	// @todo: implement me.
	return nil, pdfium_errors.ErrXFAUnsupported
}
