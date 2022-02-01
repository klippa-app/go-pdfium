package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_doc.h"
// #include <stdlib.h>
import "C"

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDF_GetMetaText returns the requested metadata.
func (p *PdfiumImplementation) FPDF_GetMetaText(request *requests.FPDF_GetMetaText) (*responses.FPDF_GetMetaText, error) {
	p.Lock()
	defer p.Unlock()

	nativeDoc, err := p.getNativeDocument(request.Document)
	if err != nil {
		return nil, err
	}

	if nativeDoc.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	cstr := C.CString(request.Tag)
	defer C.free(unsafe.Pointer(cstr))

	// First get the metadata length.
	metaSize := C.FPDF_GetMetaText(nativeDoc.currentDoc, cstr, C.NULL, 0)
	if metaSize == 0 {
		return nil, errors.New("Could not get metadata")
	}

	charData := make([]byte, metaSize)
	C.FPDF_GetMetaText(nativeDoc.currentDoc, cstr, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEText(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_GetMetaText{
		Tag:   request.Tag,
		Value: transformedText,
	}, nil
}
