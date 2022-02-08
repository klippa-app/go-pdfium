//go:build !pdfium_experimental
// +build !pdfium_experimental

package implementation

import (
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDF_ImportPagesByIndex imports pages to a FPDF_DOCUMENT.
// Experimental API.
func (p *PdfiumImplementation) FPDF_ImportPagesByIndex(request *requests.FPDF_ImportPagesByIndex) (*responses.FPDF_ImportPagesByIndex, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_ImportNPagesToOne creates a new document from source document. The pages of source document will be
// combined to provide NumPagesOnXAxis x NumPagesOnYAxis pages per page of the output document.
// Experimental API.
func (p *PdfiumImplementation) FPDF_ImportNPagesToOne(request *requests.FPDF_ImportNPagesToOne) (*responses.FPDF_ImportNPagesToOne, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_NewXObjectFromPage creates a template to generate form xobjects from the source document's page at
// the given index, for use in the destination document.
// Experimental API.
func (p *PdfiumImplementation) FPDF_NewXObjectFromPage(request *requests.FPDF_NewXObjectFromPage) (*responses.FPDF_NewXObjectFromPage, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_CloseXObject closes an FPDF_XOBJECT handle created by FPDF_NewXObjectFromPage().
// Experimental API.
func (p *PdfiumImplementation) FPDF_CloseXObject(request *requests.FPDF_CloseXObject) (*responses.FPDF_CloseXObject, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_NewFormObjectFromXObject creates a new form object from an FPDF_XOBJECT object.
// Experimental API.
func (p *PdfiumImplementation) FPDF_NewFormObjectFromXObject(request *requests.FPDF_NewFormObjectFromXObject) (*responses.FPDF_NewFormObjectFromXObject, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}
