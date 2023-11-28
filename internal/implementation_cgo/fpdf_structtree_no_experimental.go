//go:build !pdfium_experimental
// +build !pdfium_experimental

package implementation_cgo

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

// FPDF_StructElement_GetActualText returns the actual text for a given element.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetActualText(request *requests.FPDF_StructElement_GetActualText) (*responses.FPDF_StructElement_GetActualText, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_GetObjType returns the object type (/Type) for a given element.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetObjType(request *requests.FPDF_StructElement_GetObjType) (*responses.FPDF_StructElement_GetObjType, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_GetParent returns the parent of the structure element.
// If structure element is StructTreeRoot, then this function will return an error.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetParent(request *requests.FPDF_StructElement_GetParent) (*responses.FPDF_StructElement_GetParent, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_GetAttributeCount returns the number of attributes for the structure element.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetAttributeCount(request *requests.FPDF_StructElement_GetAttributeCount) (*responses.FPDF_StructElement_GetAttributeCount, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_GetAttributeAtIndex returns an attribute object in the structure element.
// If the attribute object exists but is not a dict, then this
// function will return an error. This will also return an error for out of
// bounds indices.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetAttributeAtIndex(request *requests.FPDF_StructElement_GetAttributeAtIndex) (*responses.FPDF_StructElement_GetAttributeAtIndex, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_Attr_GetCount returns the number of attributes in a structure element attribute map.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_Attr_GetCount(request *requests.FPDF_StructElement_Attr_GetCount) (*responses.FPDF_StructElement_Attr_GetCount, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_Attr_GetName returns the name of an attribute in a structure element attribute map.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_Attr_GetName(request *requests.FPDF_StructElement_Attr_GetName) (*responses.FPDF_StructElement_Attr_GetName, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_Attr_GetType returns the type of an attribute in a structure element attribute map.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_Attr_GetType(request *requests.FPDF_StructElement_Attr_GetType) (*responses.FPDF_StructElement_Attr_GetType, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_Attr_GetBooleanValue returns the value of a boolean attribute in an attribute map by name as
// FPDF_BOOL. FPDF_StructElement_Attr_GetType() should have returned
// FPDF_OBJECT_BOOLEAN for this property.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_Attr_GetBooleanValue(request *requests.FPDF_StructElement_Attr_GetBooleanValue) (*responses.FPDF_StructElement_Attr_GetBooleanValue, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_Attr_GetNumberValue returns the value of a number attribute in an attribute map by name as
// float. FPDF_StructElement_Attr_GetType() should have returned
// FPDF_OBJECT_NUMBER for this property.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_Attr_GetNumberValue(request *requests.FPDF_StructElement_Attr_GetNumberValue) (*responses.FPDF_StructElement_Attr_GetNumberValue, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_Attr_GetStringValue returns the value of a string attribute in an attribute map by name as
// string. FPDF_StructElement_Attr_GetType() should have returned
// FPDF_OBJECT_STRING or FPDF_OBJECT_NAME for this property.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_Attr_GetStringValue(request *requests.FPDF_StructElement_Attr_GetStringValue) (*responses.FPDF_StructElement_Attr_GetStringValue, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_Attr_GetBlobValue returns the value of a blob attribute in an attribute map by name as
// string.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_Attr_GetBlobValue(request *requests.FPDF_StructElement_Attr_GetBlobValue) (*responses.FPDF_StructElement_Attr_GetBlobValue, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_GetMarkedContentIdCount returns the count of marked content ids for a given element.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetMarkedContentIdCount(request *requests.FPDF_StructElement_GetMarkedContentIdCount) (*responses.FPDF_StructElement_GetMarkedContentIdCount, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_GetMarkedContentIdAtIndex returns the marked content id at a given index for a given element.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetMarkedContentIdAtIndex(request *requests.FPDF_StructElement_GetMarkedContentIdAtIndex) (*responses.FPDF_StructElement_GetMarkedContentIdAtIndex, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_StructElement_GetChildMarkedContentID returns the child's content id.
// If the child exists but is not a stream or object, then this
// function will return an error. This will also return an error for out of bounds
// indices. Compared to FPDF_StructElement_GetMarkedContentIdAtIndex,
// it is scoped to the current page.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetChildMarkedContentID(request *requests.FPDF_StructElement_GetChildMarkedContentID) (*responses.FPDF_StructElement_GetChildMarkedContentID, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}
