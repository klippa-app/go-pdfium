package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"

import (
	"errors"

	"github.com/klippa-app/go-pdfium/pdfium/requests"
	"github.com/klippa-app/go-pdfium/pdfium/responses"
)

// GetPageCount counts the amount of pages
func (p *Pdfium) GetPageCount(request *requests.GetPageCount) (*responses.GetPageCount, error) {
	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	p.Lock()
	defer p.Unlock()
	return &responses.GetPageCount{
		PageCount: int(C.FPDF_GetPageCount(p.currentDoc.doc)),
	}, nil
}

// Close closes the internal references in FPDF
func (p *Pdfium) Close() error {
	if p.currentDoc == nil {
		return errors.New("no current document")
	}

	p.Lock()
	if p.currentDoc.currentPage != nil {
		C.FPDF_ClosePage(p.currentDoc.page)
		p.currentDoc.page = nil
		p.currentDoc.currentPage = nil
	}
	C.FPDF_CloseDocument(p.currentDoc.doc)
	p.currentDoc.doc = nil
	p.currentDoc = nil
	p.Unlock()
	return nil
}
