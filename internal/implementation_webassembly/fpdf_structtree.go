package implementation_webassembly

import "C"
import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDF_StructTree_GetForPage returns the structure tree for a page.
func (p *PdfiumImplementation) FPDF_StructTree_GetForPage(request *requests.FPDF_StructTree_GetForPage) (*responses.FPDF_StructTree_GetForPage, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	documentHandle, err := p.getDocumentHandle(pageHandle.documentRef)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDF_StructTree_GetForPage").Call(p.Context, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	structTree := res[0]
	if structTree == 0 {
		return nil, errors.New("could not load struct tree")
	}

	structTreeHandle := p.registerStructTree(&structTree, documentHandle)

	return &responses.FPDF_StructTree_GetForPage{
		StructTree: structTreeHandle.nativeRef,
	}, nil
}

// FPDF_StructTree_Close releases a resource allocated by FPDF_StructTree_GetForPage().
func (p *PdfiumImplementation) FPDF_StructTree_Close(request *requests.FPDF_StructTree_Close) (*responses.FPDF_StructTree_Close, error) {
	p.Lock()
	defer p.Unlock()

	structTreeHandle, err := p.getStructTreeHandle(request.StructTree)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDF_StructTree_Close").Call(p.Context, *structTreeHandle.handle)
	if err != nil {
		return nil, err
	}

	delete(p.structTreeRefs, structTreeHandle.nativeRef)

	documentHandle, err := p.getDocumentHandle(structTreeHandle.documentRef)
	if err != nil {
		return nil, err
	}

	delete(documentHandle.structTreeRefs, structTreeHandle.nativeRef)

	return &responses.FPDF_StructTree_Close{}, nil
}

// FPDF_StructTree_CountChildren counts the number of children for the structure tree.
func (p *PdfiumImplementation) FPDF_StructTree_CountChildren(request *requests.FPDF_StructTree_CountChildren) (*responses.FPDF_StructTree_CountChildren, error) {
	p.Lock()
	defer p.Unlock()

	structTreeHandle, err := p.getStructTreeHandle(request.StructTree)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDF_StructTree_CountChildren").Call(p.Context, *structTreeHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDF_StructTree_CountChildren{
		Count: int(count),
	}, nil
}

// FPDF_StructTree_GetChildAtIndex returns a child in the structure tree.
func (p *PdfiumImplementation) FPDF_StructTree_GetChildAtIndex(request *requests.FPDF_StructTree_GetChildAtIndex) (*responses.FPDF_StructTree_GetChildAtIndex, error) {
	p.Lock()
	defer p.Unlock()

	structTreeHandle, err := p.getStructTreeHandle(request.StructTree)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDF_StructTree_GetChildAtIndex").Call(p.Context, *structTreeHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	child := res[0]
	if child == 0 {
		return nil, errors.New("could not load struct tree child")
	}

	documentHandle, err := p.getDocumentHandle(structTreeHandle.documentRef)
	if err != nil {
		return nil, err
	}

	structElementHandle := p.registerStructElement(&child, documentHandle)

	return &responses.FPDF_StructTree_GetChildAtIndex{
		StructElement: structElementHandle.nativeRef,
	}, nil
}

// FPDF_StructElement_GetAltText returns the alt text for a given element.
func (p *PdfiumImplementation) FPDF_StructElement_GetAltText(request *requests.FPDF_StructElement_GetAltText) (*responses.FPDF_StructElement_GetAltText, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetAltText").Call(p.Context, *structElementHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	altTextLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if altTextLength == 0 {
		return nil, errors.New("Could not get alt text")
	}

	charDataPointer, err := p.ByteArrayPointer(altTextLength, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDF_StructElement_GetAltText").Call(p.Context, *structElementHandle.handle, charDataPointer.Pointer, altTextLength)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_StructElement_GetAltText{
		AltText: transformedText,
	}, nil
}

// FPDF_StructElement_GetMarkedContentID returns the marked content ID for a given element.
func (p *PdfiumImplementation) FPDF_StructElement_GetMarkedContentID(request *requests.FPDF_StructElement_GetMarkedContentID) (*responses.FPDF_StructElement_GetMarkedContentID, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetMarkedContentID").Call(p.Context, *structElementHandle.handle)
	if err != nil {
		return nil, err
	}

	markedContentID := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDF_StructElement_GetMarkedContentID{
		MarkedContentID: int(markedContentID),
	}, nil
}

// FPDF_StructElement_GetType returns the type (/S) for a given element.
func (p *PdfiumImplementation) FPDF_StructElement_GetType(request *requests.FPDF_StructElement_GetType) (*responses.FPDF_StructElement_GetType, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetType").Call(p.Context, *structElementHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	typeLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if typeLength == 0 {
		return nil, errors.New("Could not get type")
	}

	charDataPointer, err := p.ByteArrayPointer(typeLength, nil)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDF_StructElement_GetType").Call(p.Context, *structElementHandle.handle, charDataPointer.Pointer, typeLength)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_StructElement_GetType{
		Type: transformedText,
	}, nil
}

// FPDF_StructElement_GetTitle returns the title (/T) for a given element.
func (p *PdfiumImplementation) FPDF_StructElement_GetTitle(request *requests.FPDF_StructElement_GetTitle) (*responses.FPDF_StructElement_GetTitle, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetTitle").Call(p.Context, *structElementHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	titleLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if titleLength == 0 {
		return nil, errors.New("Could not get title")
	}

	charDataPointer, err := p.ByteArrayPointer(titleLength, nil)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDF_StructElement_GetTitle").Call(p.Context, *structElementHandle.handle, charDataPointer.Pointer, titleLength)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_StructElement_GetTitle{
		Title: transformedText,
	}, nil
}

// FPDF_StructElement_CountChildren counts the number of children for the structure element.
func (p *PdfiumImplementation) FPDF_StructElement_CountChildren(request *requests.FPDF_StructElement_CountChildren) (*responses.FPDF_StructElement_CountChildren, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDF_StructElement_CountChildren").Call(p.Context, *structElementHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDF_StructElement_CountChildren{
		Count: int(count),
	}, nil
}

// FPDF_StructElement_GetChildAtIndex returns a child in the structure element.
// If the child exists but is not an element, then this function will
// return an error. This will also return an error for out of bounds indices.
func (p *PdfiumImplementation) FPDF_StructElement_GetChildAtIndex(request *requests.FPDF_StructElement_GetChildAtIndex) (*responses.FPDF_StructElement_GetChildAtIndex, error) {
	p.Lock()
	defer p.Unlock()

	parentStructElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetChildAtIndex").Call(p.Context, *parentStructElementHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	child := res[0]
	if child == 0 {
		return nil, errors.New("could not load struct element child")
	}

	documentHandle, err := p.getDocumentHandle(parentStructElementHandle.documentRef)
	if err != nil {
		return nil, err
	}

	structElementHandle := p.registerStructElement(&child, documentHandle)

	return &responses.FPDF_StructElement_GetChildAtIndex{
		StructElement: structElementHandle.nativeRef,
	}, nil
}

// FPDF_StructElement_GetID returns the ID for a given element.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetID(request *requests.FPDF_StructElement_GetID) (*responses.FPDF_StructElement_GetID, error) {
	p.Lock()
	defer p.Unlock()

	structElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetID").Call(p.Context, *structElementHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	idLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if idLength == 0 {
		return nil, errors.New("Could not get ID")
	}

	charDataPointer, err := p.ByteArrayPointer(idLength, nil)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDF_StructElement_GetID").Call(p.Context, *structElementHandle.handle, charDataPointer.Pointer, idLength)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetLang").Call(p.Context, *structElementHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	langLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if langLength == 0 {
		return nil, errors.New("Could not get lang")
	}

	charDataPointer, err := p.ByteArrayPointer(langLength, nil)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDF_StructElement_GetLang").Call(p.Context, *structElementHandle.handle, charDataPointer.Pointer, langLength)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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

	attributeName, err := p.CString(request.AttributeName)
	if err != nil {
		return nil, err
	}
	defer attributeName.Free()

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetStringAttribute").Call(p.Context, *structElementHandle.handle, attributeName.Pointer, 0, 0)
	if err != nil {
		return nil, err
	}

	attributeLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if attributeLength == 0 {
		return nil, errors.New("could not get attribute")
	}

	charDataPointer, err := p.ByteArrayPointer(attributeLength, nil)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDF_StructElement_GetStringAttribute").Call(p.Context, *structElementHandle.handle, attributeName.Pointer, charDataPointer.Pointer, attributeLength)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetActualText").Call(p.Context, *structElementHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	textLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if textLength == 0 {
		return nil, errors.New("could not get actual text")
	}

	charDataPointer, err := p.ByteArrayPointer(textLength, nil)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDF_StructElement_GetActualText").Call(p.Context, *structElementHandle.handle, charDataPointer.Pointer, textLength)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetObjType").Call(p.Context, *structElementHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	objTypeLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if objTypeLength == 0 {
		return nil, errors.New("could not get obj type")
	}

	charDataPointer, err := p.ByteArrayPointer(objTypeLength, nil)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDF_StructElement_GetObjType").Call(p.Context, *structElementHandle.handle, charDataPointer.Pointer, objTypeLength)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetParent").Call(p.Context, *structElementHandle.handle)
	if err != nil {
		return nil, err
	}

	parent := res[0]
	if parent == 0 {
		return nil, errors.New("could not get struct element parent")
	}

	parentStructElementHandle := p.registerStructElement(&parent, nil)

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

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetAttributeCount").Call(p.Context, *structElementHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetAttributeAtIndex").Call(p.Context, *structElementHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	structElementAttribute := res[0]
	if structElementAttribute == 0 {
		return nil, errors.New("could not get struct element attribute")
	}

	structElementAttributeHandle := p.registerStructElementAttribute(&structElementAttribute)

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

	res, err := p.Module.ExportedFunction("FPDF_StructElement_Attr_GetCount").Call(p.Context, *structElementAttributeHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))
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
	charDataPointer, err := p.ByteArrayPointer(1, []byte{1})
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	objTypeLengthPointer, err := p.ULongPointer()
	if err != nil {
		return nil, err
	}
	defer objTypeLengthPointer.Free()

	res, err := p.Module.ExportedFunction("FPDF_StructElement_Attr_GetName").Call(p.Context, *structElementAttributeHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)), charDataPointer.Pointer, 1, objTypeLengthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get attribute name")
	}

	objTypeLength, err := objTypeLengthPointer.Value()
	if err != nil {
		return nil, err
	}

	charDataPointer, err = p.ByteArrayPointer(objTypeLength, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDF_StructElement_Attr_GetName").Call(p.Context, *structElementAttributeHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)), charDataPointer.Pointer, objTypeLength, objTypeLengthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success = *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get attribute name")
	}

	charData, err := charDataPointer.Value(true)
	if err != nil {
		return nil, err
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

	attributeName, err := p.CString(request.Name)
	if err != nil {
		return nil, err
	}
	defer attributeName.Free()

	res, err := p.Module.ExportedFunction("FPDF_StructElement_Attr_GetType").Call(p.Context, *structElementAttributeHandle.handle, attributeName.Pointer)
	if err != nil {
		return nil, err
	}

	attrType := *(*int32)(unsafe.Pointer(&res[0]))

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

	attributeName, err := p.CString(request.Name)
	if err != nil {
		return nil, err
	}
	defer attributeName.Free()

	outValuePointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer outValuePointer.Free()

	res, err := p.Module.ExportedFunction("FPDF_StructElement_Attr_GetBooleanValue").Call(p.Context, *structElementAttributeHandle.handle, attributeName.Pointer, outValuePointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get boolean value")
	}

	outValue, err := outValuePointer.Value()
	if err != nil {
		return nil, err
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

	attributeName, err := p.CString(request.Name)
	if err != nil {
		return nil, err
	}
	defer attributeName.Free()

	outValuePointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer outValuePointer.Free()

	res, err := p.Module.ExportedFunction("FPDF_StructElement_Attr_GetNumberValue").Call(p.Context, *structElementAttributeHandle.handle, attributeName.Pointer, outValuePointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get number value")
	}

	outValue, err := outValuePointer.Value()
	if err != nil {
		return nil, err
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

	attributeName, err := p.CString(request.Name)
	if err != nil {
		return nil, err
	}
	defer attributeName.Free()

	objTypeLengthPointer, err := p.ULongPointer()
	if err != nil {
		return nil, err
	}
	defer objTypeLengthPointer.Free()

	res, err := p.Module.ExportedFunction("FPDF_StructElement_Attr_GetStringValue").Call(p.Context, *structElementAttributeHandle.handle, attributeName.Pointer, 0, 0, objTypeLengthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get string value")
	}

	objTypeLength, err := objTypeLengthPointer.Value()
	if err != nil {
		return nil, err
	}

	charDataPointer, err := p.ByteArrayPointer(objTypeLength, nil)
	if err != nil {
		return nil, err
	}

	res, err = p.Module.ExportedFunction("FPDF_StructElement_Attr_GetStringValue").Call(p.Context, *structElementAttributeHandle.handle, attributeName.Pointer, charDataPointer.Pointer, objTypeLength, objTypeLengthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success = *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get string value")
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
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

	attributeName, err := p.CString(request.Name)
	if err != nil {
		return nil, err
	}
	defer attributeName.Free()

	blobLengthPointer, err := p.ULongPointer()
	if err != nil {
		return nil, err
	}
	defer blobLengthPointer.Free()

	res, err := p.Module.ExportedFunction("FPDF_StructElement_Attr_GetBlobValue").Call(p.Context, *structElementAttributeHandle.handle, attributeName.Pointer, 0, 0, blobLengthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get blob value")
	}

	blobLength, err := blobLengthPointer.Value()
	if err != nil {
		return nil, err
	}

	charDataPointer, err := p.ByteArrayPointer(blobLength, nil)
	if err != nil {
		return nil, err
	}

	res, err = p.Module.ExportedFunction("FPDF_StructElement_Attr_GetBlobValue").Call(p.Context, *structElementAttributeHandle.handle, attributeName.Pointer, charDataPointer.Pointer, blobLength, blobLengthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success = *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get blob value")
	}

	blobData, err := charDataPointer.Value(true)
	if err != nil {
		return nil, err
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

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetMarkedContentIdCount").Call(p.Context, *structElementHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetMarkedContentIdAtIndex").Call(p.Context, *structElementHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	markedContentID := *(*int32)(unsafe.Pointer(&res[0]))
	if int(markedContentID) == -1 {
		return nil, errors.New("could not get struct element marked content id")
	}

	return &responses.FPDF_StructElement_GetMarkedContentIdAtIndex{
		MarkedContentID: int(markedContentID),
	}, nil
}

// FPDF_StructElement_GetChildMarkedContentID returns the child's content id.
// If the child exists but is not a stream or object, then this
// function will return an error. This will also return an error for out of bounds
// indices. Compared to FPDF_StructElement_GetMarkedContentIdAtIndex,
// it is scoped to the current page.
// Experimental API.
func (p *PdfiumImplementation) FPDF_StructElement_GetChildMarkedContentID(request *requests.FPDF_StructElement_GetChildMarkedContentID) (*responses.FPDF_StructElement_GetChildMarkedContentID, error) {
	p.Lock()
	defer p.Unlock()

	parentStructElementHandle, err := p.getStructElementHandle(request.StructElement)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDF_StructElement_GetChildMarkedContentID").Call(p.Context, *parentStructElementHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	childMarkedContentID := *(*int32)(unsafe.Pointer(&res[0]))
	if int(childMarkedContentID) == -1 {
		return nil, errors.New("could not get struct element child marked content id")
	}

	return &responses.FPDF_StructElement_GetChildMarkedContentID{
		ChildMarkedContentID: int(childMarkedContentID),
	}, nil
}
