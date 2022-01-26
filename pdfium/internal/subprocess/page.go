package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"
import (
	"github.com/klippa-app/go-pdfium/pdfium/pdfium_errors"
)

// GetPageSize returns the page size in points
// One point is 1/72 inch (around 0.3528 mm)
func (p *Pdfium) loadPage(page int) error {
	// Already loaded this page.
	if p.currentPageNumber != nil && *p.currentPageNumber == page {
		return nil
	}
	
	if p.currentPageNumber != nil {
		// Unload the current page.
		C.FPDF_ClosePage(p.currentPage)
		p.currentPage = nil
		p.currentPageNumber = nil
	}

	pageObject := C.FPDF_LoadPage(p.currentDoc, C.int(page))
	if pageObject == nil {
		return pdfium_errors.ErrPage
	}

	p.currentPage = pageObject
	p.currentPageNumber = &page

	return nil
}
