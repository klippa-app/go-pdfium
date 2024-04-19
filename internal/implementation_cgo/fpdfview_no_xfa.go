//go:build !pdfium_xfa
// +build !pdfium_xfa

package implementation_cgo

import (
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDF_BStr_Init initializes a FPDF_BSTR.
//
// Only available when using a build that includes XFA and when using the
// build flag pdfium_xfa.
func (p *PdfiumImplementation) FPDF_BStr_Init(request *requests.FPDF_BStr_Init) (*responses.FPDF_BStr_Init, error) {
	return nil, pdfium_errors.ErrXFAUnsupported
}

// FPDF_BStr_Set copies string data into the FPDF_BSTR.
//
// Only available when using a build that includes XFA and when using the
// build flag pdfium_xfa.
func (p *PdfiumImplementation) FPDF_BStr_Set(request *requests.FPDF_BStr_Set) (*responses.FPDF_BStr_Set, error) {
	return nil, pdfium_errors.ErrXFAUnsupported
}

// FPDF_BStr_Clear clears a FPDF_BSTR.
//
// Only available when using a build that includes XFA and when using the
// build flag pdfium_xfa.
func (p *PdfiumImplementation) FPDF_BStr_Clear(request *requests.FPDF_BStr_Clear) (*responses.FPDF_BStr_Clear, error) {
	return nil, pdfium_errors.ErrXFAUnsupported
}