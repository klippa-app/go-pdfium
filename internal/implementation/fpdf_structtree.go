package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_structtree.h"
import "C"
import (
	"errors"
	"unsafe"

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

	structTree := C.FPDF_StructTree_GetForPage(pageHandle.handle)
	if structTree == nil {
		return nil, errors.New("could not load struct tree")
	}

	structTreeHandle := p.registerStructTree(structTree, documentHandle)

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

	C.FPDF_StructTree_Close(structTreeHandle.handle)

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

	count := C.FPDF_StructTree_CountChildren(structTreeHandle.handle)

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

	child := C.FPDF_StructTree_GetChildAtIndex(structTreeHandle.handle, C.int(request.Index))
	if child == nil {
		return nil, errors.New("could not load struct tree child")
	}

	documentHandle, err := p.getDocumentHandle(structTreeHandle.documentRef)
	if err != nil {
		return nil, err
	}

	structElementHandle := p.registerStructElement(child, documentHandle)

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

	altTextLength := C.FPDF_StructElement_GetAltText(structElementHandle.handle, nil, 0)
	if altTextLength == 0 {
		return nil, errors.New("Could not get alt text")
	}

	charData := make([]byte, altTextLength)
	C.FPDF_StructElement_GetAltText(structElementHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

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

	markedContentID := C.FPDF_StructElement_GetMarkedContentID(structElementHandle.handle)

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

	typeLength := C.FPDF_StructElement_GetType(structElementHandle.handle, nil, 0)
	if typeLength == 0 {
		return nil, errors.New("Could not get type")
	}

	charData := make([]byte, typeLength)
	C.FPDF_StructElement_GetType(structElementHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

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

	titleLength := C.FPDF_StructElement_GetTitle(structElementHandle.handle, nil, 0)
	if titleLength == 0 {
		return nil, errors.New("Could not get title")
	}

	charData := make([]byte, titleLength)
	C.FPDF_StructElement_GetTitle(structElementHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

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

	count := C.FPDF_StructElement_CountChildren(structElementHandle.handle)

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

	child := C.FPDF_StructElement_GetChildAtIndex(parentStructElementHandle.handle, C.int(request.Index))
	if child == nil {
		return nil, errors.New("could not load struct element child")
	}

	documentHandle, err := p.getDocumentHandle(parentStructElementHandle.documentRef)
	if err != nil {
		return nil, err
	}

	structElementHandle := p.registerStructElement(child, documentHandle)

	return &responses.FPDF_StructElement_GetChildAtIndex{
		StructElement: structElementHandle.nativeRef,
	}, nil
}
