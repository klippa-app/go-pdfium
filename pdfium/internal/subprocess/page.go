package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"
import (
	"github.com/klippa-app/go-pdfium/pdfium/pdfium_errors"
)

// loadPage changes the active page if it's different from what's currently
// open and closes the page that's currently open if any is open.
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
