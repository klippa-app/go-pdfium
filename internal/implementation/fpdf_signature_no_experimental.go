//go:build !pdfium_experimental
// +build !pdfium_experimental

package implementation

import (
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDF_GetSignatureCount returns the total number of signatures in the document.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetSignatureCount(request *requests.FPDF_GetSignatureCount) (*responses.FPDF_GetSignatureCount, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_GetSignatureObject returns the Nth signature of the document.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetSignatureObject(request *requests.FPDF_GetSignatureObject) (*responses.FPDF_GetSignatureObject, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFSignatureObj_GetContents returns the contents of a signature object.
// Experimental API.
func (p *PdfiumImplementation) FPDFSignatureObj_GetContents(request *requests.FPDFSignatureObj_GetContents) (*responses.FPDFSignatureObj_GetContents, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFSignatureObj_GetByteRange returns the byte range of a signature object.
// Experimental API.
func (p *PdfiumImplementation) FPDFSignatureObj_GetByteRange(request *requests.FPDFSignatureObj_GetByteRange) (*responses.FPDFSignatureObj_GetByteRange, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFSignatureObj_GetSubFilter returns the encoding of the value of a signature object.
// Experimental API.
func (p *PdfiumImplementation) FPDFSignatureObj_GetSubFilter(request *requests.FPDFSignatureObj_GetSubFilter) (*responses.FPDFSignatureObj_GetSubFilter, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFSignatureObj_GetReason returns the reason (comment) of the signature object.
// Experimental API.
func (p *PdfiumImplementation) FPDFSignatureObj_GetReason(request *requests.FPDFSignatureObj_GetReason) (*responses.FPDFSignatureObj_GetReason, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFSignatureObj_GetTime returns the time of signing of a signature object.
// Experimental API.
func (p *PdfiumImplementation) FPDFSignatureObj_GetTime(request *requests.FPDFSignatureObj_GetTime) (*responses.FPDFSignatureObj_GetTime, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFSignatureObj_GetDocMDPPermission returns the DocMDP permission of a signature object.
// Experimental API.
func (p *PdfiumImplementation) FPDFSignatureObj_GetDocMDPPermission(request *requests.FPDFSignatureObj_GetDocMDPPermission) (*responses.FPDFSignatureObj_GetDocMDPPermission, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}
