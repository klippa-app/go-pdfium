//go:build pdfium_experimental
// +build pdfium_experimental

package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_structtree.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
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
