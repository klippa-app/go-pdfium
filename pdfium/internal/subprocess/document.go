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
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	return &responses.GetPageCount{
		PageCount: int(C.FPDF_GetPageCount(p.currentDoc)),
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
