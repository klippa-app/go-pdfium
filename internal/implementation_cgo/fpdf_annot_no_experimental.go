//go:build !pdfium_experimental
// +build !pdfium_experimental

package implementation_cgo

import (
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFAnnot_IsSupportedSubtype returns whether an annotation subtype is currently supported for creation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_IsSupportedSubtype(request *requests.FPDFAnnot_IsSupportedSubtype) (*responses.FPDFAnnot_IsSupportedSubtype, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFPage_CreateAnnot creates an annotation in the given page of the given subtype. If the specified
// subtype is illegal or unsupported, then a new annotation will not be created.
// Must call FPDFPage_CloseAnnot() when the annotation returned by this
// function is no longer needed.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_CreateAnnot(request *requests.FPDFPage_CreateAnnot) (*responses.FPDFPage_CreateAnnot, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFPage_GetAnnotCount returns the number of annotations in a given page.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_GetAnnotCount(request *requests.FPDFPage_GetAnnotCount) (*responses.FPDFPage_GetAnnotCount, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFPage_GetAnnot returns annotation at the given page and index. Must call FPDFPage_CloseAnnot() when the
// annotation returned by this function is no longer needed.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_GetAnnot(request *requests.FPDFPage_GetAnnot) (*responses.FPDFPage_GetAnnot, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFPage_GetAnnotIndex returns the index of the given annotation in the given page. This is the opposite of
// FPDFPage_GetAnnot().
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_GetAnnotIndex(request *requests.FPDFPage_GetAnnotIndex) (*responses.FPDFPage_GetAnnotIndex, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFPage_CloseAnnot closes an annotation. Must be called when the annotation returned by
// FPDFPage_CreateAnnot() or FPDFPage_GetAnnot() is no longer needed. This
// function does not remove the annotation from the document.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_CloseAnnot(request *requests.FPDFPage_CloseAnnot) (*responses.FPDFPage_CloseAnnot, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFPage_RemoveAnnot removes the annotation in the given page at the given index.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_RemoveAnnot(request *requests.FPDFPage_RemoveAnnot) (*responses.FPDFPage_RemoveAnnot, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetSubtype returns the subtype of an annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetSubtype(request *requests.FPDFAnnot_GetSubtype) (*responses.FPDFAnnot_GetSubtype, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_IsObjectSupportedSubtype checks whether an annotation subtype is currently supported for object extraction,
// update, and removal.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_IsObjectSupportedSubtype(request *requests.FPDFAnnot_IsObjectSupportedSubtype) (*responses.FPDFAnnot_IsObjectSupportedSubtype, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_UpdateObject updates the given object in the given annotation. The object must be in the annotation already and must have
// been retrieved by FPDFAnnot_GetObject(). Currently, only ink and stamp
// annotations are supported by this API. Also note that only path, image, and
// /text objects have APIs for modification; see FPDFPath_*(), FPDFText_*(), and
// FPDFImageObj_*().
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_UpdateObject(request *requests.FPDFAnnot_UpdateObject) (*responses.FPDFAnnot_UpdateObject, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_AddInkStroke adds a new InkStroke, represented by an array of points, to the InkList of
// the annotation. The API creates an InkList if one doesn't already exist in the annotation.
// This API works only for ink annotations. Please refer to ISO 32000-1:2008
// spec, section 12.5.6.13.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_AddInkStroke(request *requests.FPDFAnnot_AddInkStroke) (*responses.FPDFAnnot_AddInkStroke, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_RemoveInkList removes an InkList in the given annotation.
// This API works only for ink annotations.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_RemoveInkList(request *requests.FPDFAnnot_RemoveInkList) (*responses.FPDFAnnot_RemoveInkList, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_AppendObject adds the given object to the given annotation. The object must have been created by
// FPDFPageObj_CreateNew{Path|Rect}() or FPDFPageObj_New{Text|Image}Obj(), and
// will be owned by the annotation. Note that an object cannot belong to more than one
// annotation. Currently, only ink and stamp annotations are supported by this API.
// Also note that only path, image, and text objects have APIs for creation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_AppendObject(request *requests.FPDFAnnot_AppendObject) (*responses.FPDFAnnot_AppendObject, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetObjectCount returns the total number of objects in the given annotation, including path objects, text
// objects, external objects, image objects, and shading objects.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetObjectCount(request *requests.FPDFAnnot_GetObjectCount) (*responses.FPDFAnnot_GetObjectCount, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetObject returns the object in the given annotation at the given index.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetObject(request *requests.FPDFAnnot_GetObject) (*responses.FPDFAnnot_GetObject, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_RemoveObject removes the object in the given annotation at the given index.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_RemoveObject(request *requests.FPDFAnnot_RemoveObject) (*responses.FPDFAnnot_RemoveObject, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_SetColor sets the color of an annotation. Fails when called on annotations with
// appearance streams already defined; instead use
// FPDFPath_Set{Stroke|Fill}Color().
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetColor(request *requests.FPDFAnnot_SetColor) (*responses.FPDFAnnot_SetColor, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetColor returns the color of an annotation. If no color is specified, default to yellow
// for highlight annotation, black for all else. Fails when called on
// annotations with appearance streams already defined; instead use
// FPDFPath_Get{Stroke|Fill}Color().
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetColor(request *requests.FPDFAnnot_GetColor) (*responses.FPDFAnnot_GetColor, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_HasAttachmentPoints returns whether the annotation is of a type that has attachment points
// (i.e. quadpoints). Quadpoints are the vertices of the rectangle that
// encompasses the texts affected by the annotation. They provide the
// coordinates in the page where the annotation is attached. Only text markup
// annotations (i.e. highlight, strikeout, squiggly, and underline) and link
// annotations have quadpoints.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_HasAttachmentPoints(request *requests.FPDFAnnot_HasAttachmentPoints) (*responses.FPDFAnnot_HasAttachmentPoints, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_SetAttachmentPoints replaces the attachment points (i.e. quadpoints) set of an annotation at
// the given quad index. This index needs to be within the result of
// FPDFAnnot_CountAttachmentPoints().
// If the annotation's appearance stream is defined and this annotation is of a
// type with quadpoints, then update the bounding box too if the new quadpoints
// define a bigger one.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetAttachmentPoints(request *requests.FPDFAnnot_SetAttachmentPoints) (*responses.FPDFAnnot_SetAttachmentPoints, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_AppendAttachmentPoints appends to the list of attachment points (i.e. quadpoints) of an annotation.
// If the annotation's appearance stream is defined and this annotation is of a
// type with quadpoints, then update the bounding box too if the new quadpoints
// define a bigger one.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_AppendAttachmentPoints(request *requests.FPDFAnnot_AppendAttachmentPoints) (*responses.FPDFAnnot_AppendAttachmentPoints, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_CountAttachmentPoints returns the number of sets of quadpoints of an annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_CountAttachmentPoints(request *requests.FPDFAnnot_CountAttachmentPoints) (*responses.FPDFAnnot_CountAttachmentPoints, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetAttachmentPoints returns the attachment points (i.e. quadpoints) of an annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetAttachmentPoints(request *requests.FPDFAnnot_GetAttachmentPoints) (*responses.FPDFAnnot_GetAttachmentPoints, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_SetRect sets the annotation rectangle defining the location of the annotation. If the
// annotation's appearance stream is defined and this annotation is of a type
// without quadpoints, then update the bounding box too if the new rectangle
// defines a bigger one.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetRect(request *requests.FPDFAnnot_SetRect) (*responses.FPDFAnnot_SetRect, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetRect returns the annotation rectangle defining the location of the annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetRect(request *requests.FPDFAnnot_GetRect) (*responses.FPDFAnnot_GetRect, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetVertices returns the vertices of a polygon or polyline annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetVertices(request *requests.FPDFAnnot_GetVertices) (*responses.FPDFAnnot_GetVertices, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetInkListCount returns the number of paths in the ink list of an ink annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetInkListCount(request *requests.FPDFAnnot_GetInkListCount) (*responses.FPDFAnnot_GetInkListCount, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetInkListPath returns a path in the ink list of an ink annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetInkListPath(request *requests.FPDFAnnot_GetInkListPath) (*responses.FPDFAnnot_GetInkListPath, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetLine returns the starting and ending coordinates of a line annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetLine(request *requests.FPDFAnnot_GetLine) (*responses.FPDFAnnot_GetLine, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_SetBorder sets the characteristics of the annotation's border (rounded rectangle).
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetBorder(request *requests.FPDFAnnot_SetBorder) (*responses.FPDFAnnot_SetBorder, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetBorder returns the characteristics of the annotation's border (rounded rectangle).
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetBorder(request *requests.FPDFAnnot_GetBorder) (*responses.FPDFAnnot_GetBorder, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_HasKey checks whether the given annotation's dictionary has the given key as a key.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_HasKey(request *requests.FPDFAnnot_HasKey) (*responses.FPDFAnnot_HasKey, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetValueType returns the type of the value corresponding to the given key the annotation's dictionary.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetValueType(request *requests.FPDFAnnot_GetValueType) (*responses.FPDFAnnot_GetValueType, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_SetStringValue sets the string value corresponding to the given key in the annotations's dictionary,
// overwriting the existing value if any. The value type would be
// FPDF_OBJECT_STRING after this function call succeeds.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetStringValue(request *requests.FPDFAnnot_SetStringValue) (*responses.FPDFAnnot_SetStringValue, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetStringValue returns the string value corresponding to the given key in the annotations's dictionary.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetStringValue(request *requests.FPDFAnnot_GetStringValue) (*responses.FPDFAnnot_GetStringValue, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetNumberValue returns the float value corresponding to the given key in the annotations's dictionary.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetNumberValue(request *requests.FPDFAnnot_GetNumberValue) (*responses.FPDFAnnot_GetNumberValue, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_SetAP sets the AP (appearance string) in annotations's dictionary for a given appearance mode.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetAP(request *requests.FPDFAnnot_SetAP) (*responses.FPDFAnnot_SetAP, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetAP returns the AP (appearance string) from annotation's dictionary for a given
// appearance mode.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetAP(request *requests.FPDFAnnot_GetAP) (*responses.FPDFAnnot_GetAP, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetLinkedAnnot returns the annotation corresponding to the given key in the annotations's dictionary. Common
// keys for linking annotations include "IRT" and "Popup". Must call
// FPDFPage_CloseAnnot() when the annotation returned by this function is no
// longer needed.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetLinkedAnnot(request *requests.FPDFAnnot_GetLinkedAnnot) (*responses.FPDFAnnot_GetLinkedAnnot, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFlags returns the annotation flags of the given annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFlags(request *requests.FPDFAnnot_GetFlags) (*responses.FPDFAnnot_GetFlags, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_SetFlags sets the annotation flags of the given annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetFlags(request *requests.FPDFAnnot_SetFlags) (*responses.FPDFAnnot_SetFlags, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFormFieldFlags returns the form field annotation flags of the given annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormFieldFlags(request *requests.FPDFAnnot_GetFormFieldFlags) (*responses.FPDFAnnot_GetFormFieldFlags, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFormFieldAtPoint returns an interactive form annotation whose rectangle contains a given
// point on a page. Must call FPDFPage_CloseAnnot() when the annotation returned
// is no longer needed.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormFieldAtPoint(request *requests.FPDFAnnot_GetFormFieldAtPoint) (*responses.FPDFAnnot_GetFormFieldAtPoint, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFormAdditionalActionJavaScript returns the JavaScript of an event of the annotation's additional actions.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormAdditionalActionJavaScript(request *requests.FPDFAnnot_GetFormAdditionalActionJavaScript) (*responses.FPDFAnnot_GetFormAdditionalActionJavaScript, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFormFieldName returns the name of the given annotation, which is an interactive form annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormFieldName(request *requests.FPDFAnnot_GetFormFieldName) (*responses.FPDFAnnot_GetFormFieldName, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFormFieldAlternateName returns the alternate name of an annotation, which is an interactive form annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormFieldAlternateName(request *requests.FPDFAnnot_GetFormFieldAlternateName) (*responses.FPDFAnnot_GetFormFieldAlternateName, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFormFieldType returns the form field type of the given annotation, which is an interactive form annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormFieldType(request *requests.FPDFAnnot_GetFormFieldType) (*responses.FPDFAnnot_GetFormFieldType, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFormFieldValue returns the value of the given annotation, which is an interactive form annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormFieldValue(request *requests.FPDFAnnot_GetFormFieldValue) (*responses.FPDFAnnot_GetFormFieldValue, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetOptionCount returns the number of options in the annotation's "Opt" dictionary. Intended for
// use with listbox and combobox widget annotations.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetOptionCount(request *requests.FPDFAnnot_GetOptionCount) (*responses.FPDFAnnot_GetOptionCount, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetOptionLabel returns the string value for the label of the option at the given index in annotation's
// "Opt" dictionary. Intended for use with listbox and combobox widget
// annotations.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetOptionLabel(request *requests.FPDFAnnot_GetOptionLabel) (*responses.FPDFAnnot_GetOptionLabel, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_IsOptionSelected returns whether or not the option at the given index in annotation's "Opt" dictionary
// is selected. Intended for use with listbox and combobox widget annotations.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_IsOptionSelected(request *requests.FPDFAnnot_IsOptionSelected) (*responses.FPDFAnnot_IsOptionSelected, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFontSize returns the float value of the font size for an annotation with variable text.
// If 0, the font is to be auto-sized: its size is computed as a function of
// the height of the annotation rectangle.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFontSize(request *requests.FPDFAnnot_GetFontSize) (*responses.FPDFAnnot_GetFontSize, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFontColor returns the RGB value of the font color for an annotation with variable text.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFontColor(request *requests.FPDFAnnot_GetFontColor) (*responses.FPDFAnnot_GetFontColor, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_IsChecked returns whether the given annotation is a form widget that is checked. Intended for use with
// checkbox and radio button widgets.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_IsChecked(request *requests.FPDFAnnot_IsChecked) (*responses.FPDFAnnot_IsChecked, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_SetFocusableSubtypes returns the list of focusable annotation subtypes. Annotations of subtype
// FPDF_ANNOT_WIDGET are by default focusable. New subtypes set using this API
// will override the existing subtypes.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetFocusableSubtypes(request *requests.FPDFAnnot_SetFocusableSubtypes) (*responses.FPDFAnnot_SetFocusableSubtypes, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFocusableSubtypesCount returns the count of focusable annotation subtypes as set by host.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFocusableSubtypesCount(request *requests.FPDFAnnot_GetFocusableSubtypesCount) (*responses.FPDFAnnot_GetFocusableSubtypesCount, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFocusableSubtypes returns the list of focusable annotation subtype as set by host.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFocusableSubtypes(request *requests.FPDFAnnot_GetFocusableSubtypes) (*responses.FPDFAnnot_GetFocusableSubtypes, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetLink returns FPDF_LINK object for the given annotation. Intended to use for link annotations.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetLink(request *requests.FPDFAnnot_GetLink) (*responses.FPDFAnnot_GetLink, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFormControlCount returns the count of annotations in the annotation's control group.
// A group of interactive form annotations is collectively called a form
// control group. Here, annotation, an interactive form annotation, should be
// either a radio button or a checkbox.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormControlCount(request *requests.FPDFAnnot_GetFormControlCount) (*responses.FPDFAnnot_GetFormControlCount, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFormControlIndex returns the index of the given annotation it's control group.
// A group of interactive form annotations is collectively called a form
// control group. Here, the annotation, an interactive form annotation, should be
// either a radio button or a checkbox.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormControlIndex(request *requests.FPDFAnnot_GetFormControlIndex) (*responses.FPDFAnnot_GetFormControlIndex, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFormFieldExportValue returns the export value of the given annotation which is an interactive form annotation.
// Intended for use with radio button and checkbox widget annotations.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormFieldExportValue(request *requests.FPDFAnnot_GetFormFieldExportValue) (*responses.FPDFAnnot_GetFormFieldExportValue, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_SetURI adds a URI action to the given annotation, overwriting the existing action, if any.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetURI(request *requests.FPDFAnnot_SetURI) (*responses.FPDFAnnot_SetURI, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_GetFileAttachment get the attachment from the given annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFileAttachment(request *requests.FPDFAnnot_GetFileAttachment) (*responses.FPDFAnnot_GetFileAttachment, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAnnot_AddFileAttachment Add an embedded file to the given annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_AddFileAttachment(request *requests.FPDFAnnot_AddFileAttachment) (*responses.FPDFAnnot_AddFileAttachment, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}
