package implementation

import (
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFDOC_InitFormFillEnvironment initializes form fill environment
// This function should be called before any form fill operation.
func (p *PdfiumImplementation) FPDFDOC_InitFormFillEnvironment(request *requests.FPDFDOC_InitFormFillEnvironment) (*responses.FPDFDOC_InitFormFillEnvironment, error) {
	return nil, nil
}

// FPDFDOC_ExitFormFillEnvironment takes ownership of the handle and exits form fill environment.
func (p *PdfiumImplementation) FPDFDOC_ExitFormFillEnvironment(request *requests.FPDFDOC_ExitFormFillEnvironment) (*responses.FPDFDOC_ExitFormFillEnvironment, error) {
	return nil, nil
}

// FORM_OnAfterLoadPage
// This method is required for implementing all the form related
// functions. Should be invoked after user successfully loaded a
// PDF page, and FPDFDOC_InitFormFillEnvironment() has been invoked.
func (p *PdfiumImplementation) FORM_OnAfterLoadPage(request *requests.FORM_OnAfterLoadPage) (*responses.FORM_OnAfterLoadPage, error) {
	return nil, nil
}

// FORM_OnBeforeClosePage
// This method is required for implementing all the form related
// functions. Should be invoked before user closes the PDF page.
func (p *PdfiumImplementation) FORM_OnBeforeClosePage(request *requests.FORM_OnBeforeClosePage) (*responses.FORM_OnBeforeClosePage, error) {
	return nil, nil
}

// FORM_DoDocumentJSAction
// This method is required for performing document-level JavaScript
// actions. It should be invoked after the PDF document has been loaded.
// If there is document-level JavaScript action embedded in the
// document, this method will execute the JavaScript action. Otherwise,
// the method will do nothing.
func (p *PdfiumImplementation) FORM_DoDocumentJSAction(request *requests.FORM_DoDocumentJSAction) (*responses.FORM_DoDocumentJSAction, error) {
	return nil, nil
}

// FORM_DoDocumentOpenAction
// This method is required for performing open-action when the document
// is opened.
// This method will do nothing if there are no open-actions embedded
// in the document.
func (p *PdfiumImplementation) FORM_DoDocumentOpenAction(request *requests.FORM_DoDocumentOpenAction) (*responses.FORM_DoDocumentOpenAction, error) {
	return nil, nil
}

// FORM_DoDocumentAAction
// This method is required for performing the document's
// additional-action.
// This method will do nothing if there is no document
// additional-action corresponding to the specified type.
func (p *PdfiumImplementation) FORM_DoDocumentAAction(request *requests.FORM_DoDocumentAAction) (*responses.FORM_DoDocumentAAction, error) {
	return nil, nil
}

// FORM_DoPageAAction
// This method is required for performing the page object's
// additional-action when opened or closed.
// This method will do nothing if no additional-action corresponding
// to the specified type exists.
func (p *PdfiumImplementation) FORM_DoPageAAction(request *requests.FORM_DoPageAAction) (*responses.FORM_DoPageAAction, error) {
	return nil, nil
}

// FORM_OnMouseMove
// Call this member function when the mouse cursor moves.
func (p *PdfiumImplementation) FORM_OnMouseMove(request *requests.FORM_OnMouseMove) (*responses.FORM_OnMouseMove, error) {
	return nil, nil
}

// FORM_OnFocus
// This function focuses the form annotation at a given point. If the
// annotation at the point already has focus, nothing happens. If there
// is no annotation at the point, removes form focus.
func (p *PdfiumImplementation) FORM_OnFocus(request *requests.FORM_OnFocus) (*responses.FORM_OnFocus, error) {
	return nil, nil
}

// FORM_OnLButtonDown
// Call this member function when the user presses the left
// mouse button.
func (p *PdfiumImplementation) FORM_OnLButtonDown(request *requests.FORM_OnLButtonDown) (*responses.FORM_OnLButtonDown, error) {
	return nil, nil
}

// FORM_OnRButtonDown
// Call this member function when the user presses the right
// mouse button.
// At the present time, has no effect except in XFA builds, but is
// included for the sake of symmetry.
func (p *PdfiumImplementation) FORM_OnRButtonDown(request *requests.FORM_OnRButtonDown) (*responses.FORM_OnRButtonDown, error) {
	return nil, nil
}

// FORM_OnLButtonUp
// Call this member function when the user releases the left
// mouse button.
func (p *PdfiumImplementation) FORM_OnLButtonUp(request *requests.FORM_OnLButtonUp) (*responses.FORM_OnLButtonUp, error) {
	return nil, nil
}

// FORM_OnRButtonUp
// Call this member function when the user releases the right
// mouse button.
// At the present time, has no effect except in XFA builds, but is
// included for the sake of symmetry.
func (p *PdfiumImplementation) FORM_OnRButtonUp(request *requests.FORM_OnRButtonUp) (*responses.FORM_OnRButtonUp, error) {
	return nil, nil
}

// FORM_OnLButtonDoubleClick
// Call this member function when the user double clicks the
// left mouse button.
func (p *PdfiumImplementation) FORM_OnLButtonDoubleClick(request *requests.FORM_OnLButtonDoubleClick) (*responses.FORM_OnLButtonDoubleClick, error) {
	return nil, nil
}

// FORM_OnKeyDown
// Call this member function when a nonsystem key is pressed.
func (p *PdfiumImplementation) FORM_OnKeyDown(request *requests.FORM_OnKeyDown) (*responses.FORM_OnKeyDown, error) {
	return nil, nil
}

// FORM_OnKeyUp
// Call this member function when a nonsystem key is released.
func (p *PdfiumImplementation) FORM_OnKeyUp(request *requests.FORM_OnKeyUp) (*responses.FORM_OnKeyUp, error) {
	return nil, nil
}

// FORM_OnChar
// Call this member function when a keystroke translates to a
// nonsystem character.
func (p *PdfiumImplementation) FORM_OnChar(request *requests.FORM_OnChar) (*responses.FORM_OnChar, error) {
	return nil, nil
}

// FORM_GetSelectedText
// Call this function to obtain selected text within a form text
// field or form combobox text field.
func (p *PdfiumImplementation) FORM_GetSelectedText(request *requests.FORM_GetSelectedText) (*responses.FORM_GetSelectedText, error) {
	return nil, nil
}

// FORM_ReplaceSelection
// Call this function to replace the selected text in a form
// text field or user-editable form combobox text field with another
// text string (which can be empty or non-empty). If there is no
// selected text, this function will append the replacement text after
// the current caret position.
func (p *PdfiumImplementation) FORM_ReplaceSelection(request *requests.FORM_ReplaceSelection) (*responses.FORM_ReplaceSelection, error) {
	return nil, nil
}

// FORM_CanUndo
// Find out if it is possible for the current focused widget in a given
// form to perform an undo operation.
func (p *PdfiumImplementation) FORM_CanUndo(request *requests.FORM_CanUndo) (*responses.FORM_CanUndo, error) {
	return nil, nil
}

// FORM_CanRedo
// Find out if it is possible for the current focused widget in a given
// form to perform a redo operation.
func (p *PdfiumImplementation) FORM_CanRedo(request *requests.FORM_CanRedo) (*responses.FORM_CanRedo, error) {
	return nil, nil
}

// FORM_Undo
// Make the current focussed widget perform an undo operation.
func (p *PdfiumImplementation) FORM_Undo(request *requests.FORM_Undo) (*responses.FORM_Undo, error) {
	return nil, nil
}

// FORM_Redo
// Make the current focussed widget perform a redo operation.
func (p *PdfiumImplementation) FORM_Redo(request *requests.FORM_Redo) (*responses.FORM_Redo, error) {
	return nil, nil
}

// FORM_ForceToKillFocus
// Call this member function to force to kill the focus of the form
// field which has focus. If it would kill the focus of a form field,
// save the value of form field if was changed by theuser.
func (p *PdfiumImplementation) FORM_ForceToKillFocus(request *requests.FORM_ForceToKillFocus) (*responses.FORM_ForceToKillFocus, error) {
	return nil, nil
}

// FPDFPage_HasFormFieldAtPoint returns the form field type by point.
func (p *PdfiumImplementation) FPDFPage_HasFormFieldAtPoint(request *requests.FPDFPage_HasFormFieldAtPoint) (*responses.FPDFPage_HasFormFieldAtPoint, error) {
	return nil, nil
}

// FPDFPage_FormFieldZOrderAtPoint returns the form field z-order by point.
func (p *PdfiumImplementation) FPDFPage_FormFieldZOrderAtPoint(request *requests.FPDFPage_FormFieldZOrderAtPoint) (*responses.FPDFPage_FormFieldZOrderAtPoint, error) {
	return nil, nil
}

// FPDF_SetFormFieldHighlightColor sets the highlight color of the specified (or all) form fields
// in the document.
func (p *PdfiumImplementation) FPDF_SetFormFieldHighlightColor(request *requests.FPDF_SetFormFieldHighlightColor) (*responses.FPDF_SetFormFieldHighlightColor, error) {
	return nil, nil
}

// FPDF_SetFormFieldHighlightAlpha sets the transparency of the form field highlight color in the
// document.
func (p *PdfiumImplementation) FPDF_SetFormFieldHighlightAlpha(request *requests.FPDF_SetFormFieldHighlightAlpha) (*responses.FPDF_SetFormFieldHighlightAlpha, error) {
	return nil, nil
}

// FPDF_RemoveFormFieldHighlight removes the form field highlight color in the document.
func (p *PdfiumImplementation) FPDF_RemoveFormFieldHighlight(request *requests.FPDF_RemoveFormFieldHighlight) (*responses.FPDF_RemoveFormFieldHighlight, error) {
	return nil, nil
}

// FPDF_FFLDraw renders FormFields and popup window on a page to a device independent
// bitmap.
// This function is designed to render annotations that are
// user-interactive, which are widget annotations (for FormFields) and
// popup annotations.
// With the FPDF_ANNOT flag, this function will render a popup annotation
// when users mouse-hover on a non-widget annotation. Regardless of
// FPDF_ANNOT flag, this function will always render widget annotations
// for FormFields.
// In order to implement the FormFill functions, implementation should
// call this function after rendering functions, such as
// FPDF_RenderPageBitmap() or FPDF_RenderPageBitmap_Start(), have
// finished rendering the page contents.
func (p *PdfiumImplementation) FPDF_FFLDraw(request *requests.FPDF_FFLDraw) (*responses.FPDF_FFLDraw, error) {
	return nil, nil
}

// FPDF_LoadXFA load XFA fields of the document if it consists of XFA fields.
func (p *PdfiumImplementation) FPDF_LoadXFA(request *requests.FPDF_LoadXFA) (*responses.FPDF_LoadXFA, error) {
	return nil, nil
}
