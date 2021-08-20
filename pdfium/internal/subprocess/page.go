package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"

// GetPageSize returns the page size in points
// One point is 1/72 inch (around 0.3528 mm)
func (p *Pdfium) loadPage(page int) {
	// Already loaded this page.
	if p.currentDoc.currentPage != nil && *p.currentDoc.currentPage == page {
		return
	}

	p.Lock()
	if p.currentDoc.currentPage != nil {
		// Unload the current page.
		C.FPDF_ClosePage(p.currentDoc.page)
		p.currentDoc.page = nil
		p.currentDoc.currentPage = nil
	}

	pageObject := C.FPDF_LoadPage(p.currentDoc.doc, C.int(page))
	p.currentDoc.page = pageObject
	p.currentDoc.currentPage = &page
	p.Unlock()

	return
}
