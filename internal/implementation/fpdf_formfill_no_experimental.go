//go:build !pdfium_experimental
// +build !pdfium_experimental

package implementation

import (
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FORM_OnMouseWheel
// Call this member function when the user scrolls the mouse wheel.
// For X and Y delta, the caller must normalize
// platform-specific wheel deltas. e.g. On Windows, a delta value of 240
// for a WM_MOUSEWHEEL event normalizes to 2, since Windows defines
// WHEEL_DELTA as 120.
// Experimental API
func (p *PdfiumImplementation) FORM_OnMouseWheel(request *requests.FORM_OnMouseWheel) (*responses.FORM_OnMouseWheel, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FORM_GetFocusedText
// Call this function to obtain the text within the current focused
// field, if any.
// Experimental API
func (p *PdfiumImplementation) FORM_GetFocusedText(request *requests.FORM_GetFocusedText) (*responses.FORM_GetFocusedText, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FORM_SelectAllText
// Call this function to select all the text within the currently focused
// form text field or form combobox text field.
// Experimental API
func (p *PdfiumImplementation) FORM_SelectAllText(request *requests.FORM_SelectAllText) (*responses.FORM_SelectAllText, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FORM_GetFocusedAnnot
// Call this member function to get the currently focused annotation.
// Not currently supported for XFA forms - will report no focused
// annotation. Must call FPDFPage_CloseAnnot() when the annotation returned
// by this function is no longer needed.
// Experimental API.
func (p *PdfiumImplementation) FORM_GetFocusedAnnot(request *requests.FORM_GetFocusedAnnot) (*responses.FORM_GetFocusedAnnot, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FORM_SetFocusedAnnot
// Call this member function to set the currently focused annotation.
// The annotation can't be nil. To kill focus, use FORM_ForceToKillFocus() instead.
// Experimental API.
func (p *PdfiumImplementation) FORM_SetFocusedAnnot(request *requests.FORM_SetFocusedAnnot) (*responses.FORM_SetFocusedAnnot, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_GetFormType returns the type of form contained in the PDF document.
// If document is nil, then the return value is FORMTYPE_NONE.
// Experimental API
func (p *PdfiumImplementation) FPDF_GetFormType(request *requests.FPDF_GetFormType) (*responses.FPDF_GetFormType, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FORM_SetIndexSelected selects/deselects the value at the given index of the focused
// annotation.
// Intended for use with listbox/combobox widget types. Comboboxes
// have at most a single value selected at a time which cannot be
// deselected. Deselect on a combobox is a no-op that returns false.
// Default implementation is a no-op that will return false for
// other types.
// Not currently supported for XFA forms - will return false.
// Experimental API
func (p *PdfiumImplementation) FORM_SetIndexSelected(request *requests.FORM_SetIndexSelected) (*responses.FORM_SetIndexSelected, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FORM_IsIndexSelected returns whether or not the value at index of the focused
// annotation is currently selected.
// Intended for use with listbox/combobox widget types. Default
// implementation is a no-op that will return false for other types.
// Not currently supported for XFA forms - will return false.
// Experimental API
func (p *PdfiumImplementation) FORM_IsIndexSelected(request *requests.FORM_IsIndexSelected) (*responses.FORM_IsIndexSelected, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}
