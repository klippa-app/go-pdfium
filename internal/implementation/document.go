package implementation

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_doc.h"
// #include <stdlib.h>
import "C"

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// GetPageCount counts the amount of pages
func (p *Pdfium) GetPageCount(request *requests.GetPageCount) (*responses.GetPageCount, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	return &responses.GetPageCount{
		PageCount: int(C.FPDF_GetPageCount(p.currentDoc)),
	}, nil
}

// GetMetadata returns the requested metadata
func (p *Pdfium) GetMetadata(request *requests.GetMetadata) (*responses.GetMetadata, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	cstr := C.CString(request.Tag)
	defer C.free(unsafe.Pointer(cstr))

	// First get the metadata length.
	metaSize := C.FPDF_GetMetaText(p.currentDoc, cstr, C.NULL, 0)
	if metaSize == 0 {
		return nil, errors.New("Could not get metadata")
	}

	charData := make([]byte, metaSize)
	C.FPDF_GetMetaText(p.currentDoc, cstr, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEText(charData)
	if err != nil {
		return nil, err
	}

	return &responses.GetMetadata{
		Tag:   request.Tag,
		Value: transformedText,
	}, nil
}

// Close closes the internal references in FPDF
func (p *Pdfium) Close() error {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return errors.New("no current document")
	}

	if p.currentPageNumber != nil {
		C.FPDF_ClosePage(p.currentPage)
		p.currentPage = nil
		p.currentPageNumber = nil
	}
	C.FPDF_CloseDocument(p.currentDoc)
	p.currentDoc = nil
	return nil
}
