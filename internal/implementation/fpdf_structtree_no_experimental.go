//go:build !pdfium_experimental
// +build !pdfium_experimental

package implementation

import (
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDF_StructElement_GetID returns the ID for a given element.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetID(request *requests.FPDF_StructElement_GetID) (*responses.FPDF_StructElement_GetID, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_GetLang returns the case-insensitive IETF BCP 47 language code for an element.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetLang(request *requests.FPDF_StructElement_GetLang) (*responses.FPDF_StructElement_GetLang, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_GetStringAttribute returns a struct element attribute of type "name" or "string"
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetStringAttribute(request *requests.FPDF_StructElement_GetStringAttribute) (*responses.FPDF_StructElement_GetStringAttribute, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}
