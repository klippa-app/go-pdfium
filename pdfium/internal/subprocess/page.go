package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"

// GetPageSize returns the page size in points
// One point is 1/72 inch (around 0.3528 mm)
func (d *Document) loadPage(page int) {
	// Already loaded this page.
	if d.currentPage != nil && *d.currentPage == page {
		return
	}

	mutex.Lock()
	if d.currentPage != nil {
		// Unload the current page.
		C.FPDF_ClosePage(d.page)
		d.page = nil
		d.currentPage = nil
	}

	pageObject := C.FPDF_LoadPage(d.doc, C.int(page))
	d.page = pageObject
	d.currentPage = &page
	mutex.Unlock()

	return
}
