package implementation

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"
import (
	"errors"
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
)

// loadPage changes the active page if it's different from what's currently
// open and closes the page that's currently open if any is open.
func (p *PdfiumImplementation) loadPage(page requests.Page) (*NativePage, error) {
	if page.ByReference == nil && page.ByIndex == nil {
		return nil, errors.New("either page reference or index should be given")
	}
	if page.ByReference != nil {
		if *page.ByReference == "" {
			return nil, errors.New("page reference can't be empty")
		}
		return p.getNativePage(*page.ByReference)
	}

	doc, err := p.getNativeDocument(page.ByIndex.Document)
	if err != nil {
		return nil, err
	}

	if doc.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	// Already loaded this page.
	if doc.currentPage != nil && doc.currentPage.index == page.ByIndex.Index {
		return doc.currentPage, nil
	}

	if doc.currentPage != nil {
		doc.currentPage.Close()
		doc.currentPage = nil
	}

	pageObject := C.FPDF_LoadPage(doc.currentDoc, C.int(page.ByIndex.Index))
	if pageObject == nil {
		return nil, pdfium_errors.ErrPage
	}

	nativePage := &NativePage{
		page:  pageObject,
		index: page.ByIndex.Index,
	}

	doc.currentPage = nativePage

	return nativePage, nil
}
