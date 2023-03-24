//go:build pdfium_experimental
// +build pdfium_experimental

package implementation_cgo

/*
#cgo pkg-config: pdfium
#include "fpdf_formfill.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"unsafe"
)

// FORM_OnMouseWheel
// Call this member function when the user scrolls the mouse wheel.
// For X and Y delta, the caller must normalize
// platform-specific wheel deltas. e.g. On Windows, a delta value of 240
// for a WM_MOUSEWHEEL event normalizes to 2, since Windows defines
// WHEEL_DELTA as 120.
// Experimental API
func (p *PdfiumImplementation) FORM_OnMouseWheel(request *requests.FORM_OnMouseWheel) (*responses.FORM_OnMouseWheel, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	pageCoord := C.FS_POINTF{
		x: C.float(request.PageCoord.X),
		y: C.float(request.PageCoord.Y),
	}

	success := C.FORM_OnMouseWheel(formHandleHandle.handle, pageHandle.handle, C.int(request.Modifier), &pageCoord, C.int(request.DeltaX), C.int(request.DeltaY))
	if int(success) == 0 {
		return nil, errors.New("could not do mouse wheel")
	}

	return &responses.FORM_OnMouseWheel{}, nil
}

// FORM_GetFocusedText
// Call this function to obtain the text within the current focused
// field, if any.
// Experimental API
func (p *PdfiumImplementation) FORM_GetFocusedText(request *requests.FORM_GetFocusedText) (*responses.FORM_GetFocusedText, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	// First get the text length
	length := C.FORM_GetFocusedText(formHandleHandle.handle, pageHandle.handle, nil, 0)
	if uint64(length) == 0 {
		return nil, errors.New("could not get focused text length")
	}

	charData := make([]byte, length)
	C.FORM_GetFocusedText(formHandleHandle.handle, pageHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FORM_GetFocusedText{
		FocusedText: transformedText,
	}, nil
}

// FORM_SelectAllText
// Call this function to select all the text within the currently focused
// form text field or form combobox text field.
// Experimental API
func (p *PdfiumImplementation) FORM_SelectAllText(request *requests.FORM_SelectAllText) (*responses.FORM_SelectAllText, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	success := C.FORM_SelectAllText(formHandleHandle.handle, pageHandle.handle)
	if int(success) == 0 {
		return nil, errors.New("could not select all text")
	}

	return &responses.FORM_SelectAllText{}, nil
}

// FORM_GetFocusedAnnot
// Call this member function to get the currently focused annotation.
// Not currently supported for XFA forms - will report no focused
// annotation. Must call FPDFPage_CloseAnnot() when the annotation returned
// by this function is no longer needed.
// Experimental API.
func (p *PdfiumImplementation) FORM_GetFocusedAnnot(request *requests.FORM_GetFocusedAnnot) (*responses.FORM_GetFocusedAnnot, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageIndex := C.int(0)
	var annotation C.FPDF_ANNOTATION

	success := C.FORM_GetFocusedAnnot(formHandleHandle.handle, &pageIndex, &annotation)
	if int(success) == 0 {
		return nil, errors.New("could not get focused annotation")
	}

	annotationHandle := p.registerAnnotation(annotation)

	return &responses.FORM_GetFocusedAnnot{
		PageIndex:  int(pageIndex),
		Annotation: annotationHandle.nativeRef,
	}, nil
}

// FORM_SetFocusedAnnot
// Call this member function to set the currently focused annotation.
// The annotation can't be nil. To kill focus, use FORM_ForceToKillFocus() instead.
// Experimental API.
func (p *PdfiumImplementation) FORM_SetFocusedAnnot(request *requests.FORM_SetFocusedAnnot) (*responses.FORM_SetFocusedAnnot, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	success := C.FORM_SetFocusedAnnot(formHandleHandle.handle, annotationHandle.handle)
	if int(success) == 0 {
		return nil, errors.New("could not set focused annotation")
	}

	return &responses.FORM_SetFocusedAnnot{}, nil
}

// FPDF_GetFormType returns the type of form contained in the PDF document.
// If document is nil, then the return value is FORMTYPE_NONE.
// Experimental API
func (p *PdfiumImplementation) FPDF_GetFormType(request *requests.FPDF_GetFormType) (*responses.FPDF_GetFormType, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	formType := C.FPDF_GetFormType(documentHandle.handle)

	return &responses.FPDF_GetFormType{
		FormType: enums.FPDF_FORMTYPE(formType),
	}, nil
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
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	selected := C.FPDF_BOOL(0)
	if request.Selected {
		selected = C.FPDF_BOOL(1)
	}

	success := C.FORM_SetIndexSelected(formHandleHandle.handle, pageHandle.handle, C.int(request.Index), selected)
	if int(success) == 0 {
		return nil, errors.New("could not set index selected")
	}

	return &responses.FORM_SetIndexSelected{}, nil
}

// FORM_IsIndexSelected returns whether or not the value at index of the focused
// annotation is currently selected.
// Intended for use with listbox/combobox widget types. Default
// implementation is a no-op that will return false for other types.
// Not currently supported for XFA forms - will return false.
// Experimental API
func (p *PdfiumImplementation) FORM_IsIndexSelected(request *requests.FORM_IsIndexSelected) (*responses.FORM_IsIndexSelected, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	isIndexSelected := C.FORM_IsIndexSelected(formHandleHandle.handle, pageHandle.handle, C.int(request.Index))

	return &responses.FORM_IsIndexSelected{
		IsIndexSelected: int(isIndexSelected) == 1,
	}, nil
}

// FORM_ReplaceAndKeepSelection
// Call this function to replace the selected text in a form text field or
// user-editable form combobox text field with another text string (which
// can be empty or non-empty). If there is no selected text, this function
// will append the replacement text after the current caret position. After
// the insertion, the inserted text will be selected.
// Experimental API
func (p *PdfiumImplementation) FORM_ReplaceAndKeepSelection(request *requests.FORM_ReplaceAndKeepSelection) (*responses.FORM_ReplaceAndKeepSelection, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF8ToUTF16LE(request.Text)
	if err != nil {
		return nil, err
	}

	C.FORM_ReplaceAndKeepSelection(formHandleHandle.handle, pageHandle.handle, (C.FPDF_WIDESTRING)(unsafe.Pointer(&transformedText[0])))

	return &responses.FORM_ReplaceAndKeepSelection{}, nil
}
