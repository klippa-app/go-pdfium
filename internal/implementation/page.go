package implementation

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_edit.h"
// #include "fpdf_flatten.h"
import "C"
import (
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
)

// loadPage changes the active page if it's different from what's currently
// open and closes the page that's currently open if any is open.
func (p *PdfiumImplementation) loadPage(doc *NativeDocument, page requests.Page) (*NativePage, error) {
	if page.Reference != "" {
		return doc.getNativePage(page.Reference)
	}

	// Already loaded this page.
	if doc.currentPage != nil && doc.currentPage.index == page.Index {
		return doc.currentPage, nil
	}

	if doc.currentPage != nil {
		doc.currentPage.Close()
		doc.currentPage = nil
	}

	pageObject := C.FPDF_LoadPage(doc.currentDoc, C.int(page.Index))
	if pageObject == nil {
		return nil, pdfium_errors.ErrPage
	}

	nativePage := &NativePage{
		page:  pageObject,
		index: page.Index,
	}

	doc.currentPage = nativePage

	return nativePage, nil
}
