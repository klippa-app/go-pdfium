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
func (p *PdfiumImplementation) loadPage(page requests.Page) (*PageHandle, error) {
	if page.ByReference == nil && page.ByIndex == nil {
		return nil, errors.New("either page reference or index should be given")
	}
	if page.ByReference != nil {
		if *page.ByReference == "" {
			return nil, errors.New("page reference can't be empty")
		}
		return p.getPageHandle(*page.ByReference)
	}

	documentHandle, err := p.getDocumentHandle(page.ByIndex.Document)
	if err != nil {
		return nil, err
	}

	// Already loaded this page.
	if documentHandle.currentPage != nil && documentHandle.currentPage.index == page.ByIndex.Index {
		return documentHandle.currentPage, nil
	}

	if documentHandle.currentPage != nil {
		documentHandle.currentPage.Close()
		documentHandle.currentPage = nil
	}

	pageObject := C.FPDF_LoadPage(documentHandle.handle, C.int(page.ByIndex.Index))
	if pageObject == nil {
		return nil, pdfium_errors.ErrPage
	}

	nativePage := &PageHandle{
		handle:      pageObject,
		index:       page.ByIndex.Index,
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.currentPage = nativePage

	return nativePage, nil
}
