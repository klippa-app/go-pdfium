//go:build !pdfium_experimental
// +build !pdfium_experimental

package implementation_cgo

import (
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// GetForm returns the form elements in the given page, including option values and current values.
// Experimental API.
func (p *PdfiumImplementation) GetForm(request *requests.GetForm) (*responses.GetForm, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}
