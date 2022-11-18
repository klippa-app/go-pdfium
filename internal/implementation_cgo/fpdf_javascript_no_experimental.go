//go:build !pdfium_experimental
// +build !pdfium_experimental

package implementation_cgo

import (
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFDoc_GetJavaScriptActionCount returns the number of JavaScript actions in the given document.
// Experimental API.
func (p *PdfiumImplementation) FPDFDoc_GetJavaScriptActionCount(request *requests.FPDFDoc_GetJavaScriptActionCount) (*responses.FPDFDoc_GetJavaScriptActionCount, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFDoc_GetJavaScriptAction returns the JavaScript action at the given index in the given document.
// Experimental API.
func (p *PdfiumImplementation) FPDFDoc_GetJavaScriptAction(request *requests.FPDFDoc_GetJavaScriptAction) (*responses.FPDFDoc_GetJavaScriptAction, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFDoc_CloseJavaScriptAction closes a loaded FPDF_JAVASCRIPT_ACTION object.
// Experimental API.
func (p *PdfiumImplementation) FPDFDoc_CloseJavaScriptAction(request *requests.FPDFDoc_CloseJavaScriptAction) (*responses.FPDFDoc_CloseJavaScriptAction, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFJavaScriptAction_GetName returns the name from the javascript handle.
// Experimental API.
func (p *PdfiumImplementation) FPDFJavaScriptAction_GetName(request *requests.FPDFJavaScriptAction_GetName) (*responses.FPDFJavaScriptAction_GetName, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFJavaScriptAction_GetScript returns the script from the javascript handle
// Experimental API.
func (p *PdfiumImplementation) FPDFJavaScriptAction_GetScript(request *requests.FPDFJavaScriptAction_GetScript) (*responses.FPDFJavaScriptAction_GetScript, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}
