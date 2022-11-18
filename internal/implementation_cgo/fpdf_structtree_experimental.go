//go:build pdfium_experimental
// +build pdfium_experimental

package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdf_structtree.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"github.com/klippa-app/go-pdfium/enums"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDF_StructElement_GetID returns the ID for a given element.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetID(request *requests.FPDF_StructElement_GetID) (*responses.FPDF_StructElement_GetID, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	idLength := C.FPDF_StructElement_GetID(structElementHandle.handle, nil, 0)
	if idLength == 0 {
		return nil, errors.New("Could not get ID")
	}

	charData := make([]byte, idLength)
	C.FPDF_StructElement_GetID(structElementHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_StructElement_GetID{
		ID: transformedText,
	}, nil
}

// FPDF_StructElement_GetLang returns the case-insensitive IETF BCP 47 language code for an element.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetLang(request *requests.FPDF_StructElement_GetLang) (*responses.FPDF_StructElement_GetLang, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	langLength := C.FPDF_StructElement_GetLang(structElementHandle.handle, nil, 0)
	if langLength == 0 {
		return nil, errors.New("Could not get lang")
	}

	charData := make([]byte, langLength)
	C.FPDF_StructElement_GetLang(structElementHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_StructElement_GetLang{
		Lang: transformedText,
	}, nil
}

// FPDF_StructElement_GetStringAttribute returns a struct element attribute of type "name" or "string"
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetStringAttribute(request *requests.FPDF_StructElement_GetStringAttribute) (*responses.FPDF_StructElement_GetStringAttribute, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	attributeName := C.CString(request.AttributeName)
	defer C.free(unsafe.Pointer(attributeName))

	attributeLength := C.FPDF_StructElement_GetStringAttribute(structElementHandle.handle, attributeName, nil, 0)
	if attributeLength == 0 {
		return nil, errors.New("could not get attribute")
	}

	charData := make([]byte, attributeLength)
	C.FPDF_StructElement_GetStringAttribute(structElementHandle.handle, attributeName, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_StructElement_GetStringAttribute{
		Attribute: request.AttributeName,
		Value:     transformedText,
	}, nil
}

// FPDF_StructElement_GetActualText returns the actual text for a given element.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetActualText(request *requests.FPDF_StructElement_GetActualText) (*responses.FPDF_StructElement_GetActualText, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	textLength := C.FPDF_StructElement_GetActualText(structElementHandle.handle, nil, 0)
	if textLength == 0 {
		return nil, errors.New("could not get actual text")
	}

	charData := make([]byte, textLength)
	C.FPDF_StructElement_GetActualText(structElementHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_StructElement_GetActualText{
		Actualtext: transformedText,
	}, nil
}

// FPDF_StructElement_GetObjType returns the object type (/Type) for a given element.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetObjType(request *requests.FPDF_StructElement_GetObjType) (*responses.FPDF_StructElement_GetObjType, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	objTypeLength := C.FPDF_StructElement_GetObjType(structElementHandle.handle, nil, 0)
	if objTypeLength == 0 {
		return nil, errors.New("could not get obj type")
	}

	charData := make([]byte, objTypeLength)
	C.FPDF_StructElement_GetObjType(structElementHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_StructElement_GetObjType{
		ObjType: transformedText,
	}, nil
}

// FPDF_StructElement_GetParent returns the parent of the structure element.
// If structure element is StructTreeRoot, then this function will return an error.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetParent(request *requests.FPDF_StructElement_GetParent) (*responses.FPDF_StructElement_GetParent, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	parent := C.FPDF_StructElement_GetParent(structElementHandle.handle)
	if parent == nil {
		return nil, errors.New("could not get struct element parent")
	}

	parentStructElementHandle := p.registerStructElement(parent, nil)

	return &responses.FPDF_StructElement_GetParent{
		StructElement: parentStructElementHandle.nativeRef,
	}, nil
}

// FPDF_StructElement_GetAttributeCount returns the number of attributes for the structure element.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetAttributeCount(request *requests.FPDF_StructElement_GetAttributeCount) (*responses.FPDF_StructElement_GetAttributeCount, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	count := C.FPDF_StructElement_GetAttributeCount(structElementHandle.handle)
	if int(count) == -1 {
		return nil, errors.New("could not get struct element attribute count")
	}

	return &responses.FPDF_StructElement_GetAttributeCount{
		Count: int(count),
	}, nil
}

// FPDF_StructElement_GetAttributeAtIndex returns an attribute object in the structure element.
// If the attribute object exists but is not a dict, then this
// function will return an error. This will also return an error for out of
// bounds indices.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetAttributeAtIndex(request *requests.FPDF_StructElement_GetAttributeAtIndex) (*responses.FPDF_StructElement_GetAttributeAtIndex, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	structElementAttribute := C.FPDF_StructElement_GetAttributeAtIndex(structElementHandle.handle, C.int(request.Index))
	if structElementAttribute == nil {
		return nil, errors.New("could not get struct element attribute")
	}

	structElementAttributeHandle := p.registerStructElementAttribute(structElementAttribute)

	return &responses.FPDF_StructElement_GetAttributeAtIndex{
		StructElementAttribute: structElementAttributeHandle.nativeRef,
	}, nil
}

// FPDF_StructElement_Attr_GetCount returns the number of attributes in a structure element attribute map.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_Attr_GetCount(request *requests.FPDF_StructElement_Attr_GetCount) (*responses.FPDF_StructElement_Attr_GetCount, error) {
	p.Lock()
	defer p.Unlock()

	structElementAttributeHandle, err := p.getStructElementAttributeHandle(request.StructElementAttribute)
	if err != nil {
		return nil, err
	}

	count := C.FPDF_StructElement_Attr_GetCount(structElementAttributeHandle.handle)
	if int(count) == -1 {
		return nil, errors.New("could not get struct element attribute count")
	}

	return &responses.FPDF_StructElement_Attr_GetCount{
		Count: int(count),
	}, nil
}

// FPDF_StructElement_Attr_GetName returns the name of an attribute in a structure element attribute map.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_Attr_GetName(request *requests.FPDF_StructElement_Attr_GetName) (*responses.FPDF_StructElement_Attr_GetName, error) {
	p.Lock()
	defer p.Unlock()

	structElementAttributeHandle, err := p.getStructElementAttributeHandle(request.StructElementAttribute)
	if err != nil {
		return nil, err
	}

	// For some reason this does not work as said in the documentation.
	// You can't pass a nil buffer.
	charData := make([]byte, uint64(1))
	objTypeLength := C.ulong(0)
	success := C.FPDF_StructElement_Attr_GetName(structElementAttributeHandle.handle, C.int(request.Index), unsafe.Pointer(&charData[0]), C.ulong(len(charData)), &objTypeLength)
	if int(success) == 0 {
		return nil, errors.New("could not get attribute name")
	}

	charData = make([]byte, uint64(objTypeLength))
	success = C.FPDF_StructElement_Attr_GetName(structElementAttributeHandle.handle, C.int(request.Index), unsafe.Pointer(&charData[0]), C.ulong(len(charData)), &objTypeLength)
	if int(success) == 0 {
		return nil, errors.New("could not get attribute name")
	}

	return &responses.FPDF_StructElement_Attr_GetName{
		Name: string(charData[:len(charData)-1]), // Remove NULL terminator.
	}, nil
}

// FPDF_StructElement_Attr_GetType returns the type of an attribute in a structure element attribute map.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_Attr_GetType(request *requests.FPDF_StructElement_Attr_GetType) (*responses.FPDF_StructElement_Attr_GetType, error) {
	p.Lock()
	defer p.Unlock()

	structElementAttributeHandle, err := p.getStructElementAttributeHandle(request.StructElementAttribute)
	if err != nil {
		return nil, err
	}

	attributeName := C.CString(request.Name)
	defer C.free(unsafe.Pointer(attributeName))

	attrType := C.FPDF_StructElement_Attr_GetType(structElementAttributeHandle.handle, attributeName)

	return &responses.FPDF_StructElement_Attr_GetType{
		ObjectType: enums.FPDF_OBJECT_TYPE(attrType),
	}, nil
}

// FPDF_StructElement_Attr_GetBooleanValue returns the value of a boolean attribute in an attribute map by name as
// FPDF_BOOL. FPDF_StructElement_Attr_GetType() should have returned
// FPDF_OBJECT_BOOLEAN for this property.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_Attr_GetBooleanValue(request *requests.FPDF_StructElement_Attr_GetBooleanValue) (*responses.FPDF_StructElement_Attr_GetBooleanValue, error) {
	p.Lock()
	defer p.Unlock()

	structElementAttributeHandle, err := p.getStructElementAttributeHandle(request.StructElementAttribute)
	if err != nil {
		return nil, err
	}

	attributeName := C.CString(request.Name)
	defer C.free(unsafe.Pointer(attributeName))

	outValue := C.FPDF_BOOL(0)

	success := C.FPDF_StructElement_Attr_GetBooleanValue(structElementAttributeHandle.handle, attributeName, &outValue)
	if int(success) == 0 {
		return nil, errors.New("could not get boolean value")
	}

	return &responses.FPDF_StructElement_Attr_GetBooleanValue{
		Value: int(outValue) == 1,
	}, nil
}

// FPDF_StructElement_Attr_GetNumberValue returns the value of a number attribute in an attribute map by name as
// float. FPDF_StructElement_Attr_GetType() should have returned
// FPDF_OBJECT_NUMBER for this property.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_Attr_GetNumberValue(request *requests.FPDF_StructElement_Attr_GetNumberValue) (*responses.FPDF_StructElement_Attr_GetNumberValue, error) {
	p.Lock()
	defer p.Unlock()

	structElementAttributeHandle, err := p.getStructElementAttributeHandle(request.StructElementAttribute)
	if err != nil {
		return nil, err
	}

	attributeName := C.CString(request.Name)
	defer C.free(unsafe.Pointer(attributeName))

	outValue := C.float(0)

	success := C.FPDF_StructElement_Attr_GetNumberValue(structElementAttributeHandle.handle, attributeName, &outValue)
	if int(success) == 0 {
		return nil, errors.New("could not get number value")
	}

	return &responses.FPDF_StructElement_Attr_GetNumberValue{
		Value: float32(outValue),
	}, nil
}

// FPDF_StructElement_Attr_GetStringValue returns the value of a string attribute in an attribute map by name as
// string. FPDF_StructElement_Attr_GetType() should have returned
// FPDF_OBJECT_STRING or FPDF_OBJECT_NAME for this property.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_Attr_GetStringValue(request *requests.FPDF_StructElement_Attr_GetStringValue) (*responses.FPDF_StructElement_Attr_GetStringValue, error) {
	p.Lock()
	defer p.Unlock()

	structElementAttributeHandle, err := p.getStructElementAttributeHandle(request.StructElementAttribute)
	if err != nil {
		return nil, err
	}

	attributeName := C.CString(request.Name)
	defer C.free(unsafe.Pointer(attributeName))

	objTypeLength := C.ulong(0)
	success := C.FPDF_StructElement_Attr_GetStringValue(structElementAttributeHandle.handle, attributeName, nil, 0, &objTypeLength)
	if int(success) == 0 {
		return nil, errors.New("could not get string value")
	}

	charData := make([]byte, uint64(objTypeLength))
	success = C.FPDF_StructElement_Attr_GetStringValue(structElementAttributeHandle.handle, attributeName, unsafe.Pointer(&charData[0]), C.ulong(len(charData)), &objTypeLength)
	if int(success) == 0 {
		return nil, errors.New("could not get string value")
	}

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_StructElement_Attr_GetStringValue{
		Value: transformedText,
	}, nil
}

// FPDF_StructElement_Attr_GetBlobValue returns the value of a blob attribute in an attribute map by name as
// string.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_Attr_GetBlobValue(request *requests.FPDF_StructElement_Attr_GetBlobValue) (*responses.FPDF_StructElement_Attr_GetBlobValue, error) {
	p.Lock()
	defer p.Unlock()

	structElementAttributeHandle, err := p.getStructElementAttributeHandle(request.StructElementAttribute)
	if err != nil {
		return nil, err
	}

	attributeName := C.CString(request.Name)
	defer C.free(unsafe.Pointer(attributeName))

	blobLength := C.ulong(0)
	success := C.FPDF_StructElement_Attr_GetBlobValue(structElementAttributeHandle.handle, attributeName, nil, 0, &blobLength)
	if int(success) == 0 {
		return nil, errors.New("could not get blob value")
	}

	blobData := make([]byte, uint64(blobLength))
	success = C.FPDF_StructElement_Attr_GetBlobValue(structElementAttributeHandle.handle, attributeName, unsafe.Pointer(&blobData[0]), C.ulong(len(blobData)), &blobLength)
	if int(success) == 0 {
		return nil, errors.New("could not get blob value")
	}

	return &responses.FPDF_StructElement_Attr_GetBlobValue{
		Value: blobData,
	}, nil
}

// FPDF_StructElement_GetMarkedContentIdCount returns the count of marked content ids for a given element.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetMarkedContentIdCount(request *requests.FPDF_StructElement_GetMarkedContentIdCount) (*responses.FPDF_StructElement_GetMarkedContentIdCount, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	count := C.FPDF_StructElement_GetMarkedContentIdCount(structElementHandle.handle)
	if int(count) == -1 {
		return nil, errors.New("could not get struct element marked content id count")
	}

	return &responses.FPDF_StructElement_GetMarkedContentIdCount{
		Count: int(count),
	}, nil
}

// FPDF_StructElement_GetMarkedContentIdAtIndex returns the marked content id at a given index for a given element.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetMarkedContentIdAtIndex(request *requests.FPDF_StructElement_GetMarkedContentIdAtIndex) (*responses.FPDF_StructElement_GetMarkedContentIdAtIndex, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	markedContentID := C.FPDF_StructElement_GetMarkedContentIdAtIndex(structElementHandle.handle, C.int(request.Index))
	if int(markedContentID) == -1 {
		return nil, errors.New("could not get struct element marked content id")
	}

	return &responses.FPDF_StructElement_GetMarkedContentIdAtIndex{
		MarkedContentID: int(markedContentID),
	}, nil
}
