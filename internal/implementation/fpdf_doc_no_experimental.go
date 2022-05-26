//go:build !pdfium_experimental
// +build !pdfium_experimental

package implementation

import (
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFDest_GetView returns the view (fit type) for a given dest.
// Experimental API.
func (p *PdfiumImplementation) FPDFDest_GetView(request *requests.FPDFDest_GetView) (*responses.FPDFDest_GetView, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFLink_GetAnnot returns a FPDF_ANNOTATION object for a link.
// Experimental API.
func (p *PdfiumImplementation) FPDFLink_GetAnnot(request *requests.FPDFLink_GetAnnot) (*responses.FPDFLink_GetAnnot, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_GetPageAAction returns an additional-action from page.
// Experimental API
func (p *PdfiumImplementation) FPDF_GetPageAAction(request *requests.FPDF_GetPageAAction) (*responses.FPDF_GetPageAAction, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_GetFileIdentifier Get the file identifier defined in the trailer of a document.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetFileIdentifier(request *requests.FPDF_GetFileIdentifier) (*responses.FPDF_GetFileIdentifier, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFBookmark_GetCount returns the number of children of a bookmark.
// Experimental API.
func (p *PdfiumImplementation) FPDFBookmark_GetCount(request *requests.FPDFBookmark_GetCount) (*responses.FPDFBookmark_GetCount, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}
