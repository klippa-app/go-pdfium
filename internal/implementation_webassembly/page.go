package implementation_webassembly

import (
	"errors"

	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"

	"github.com/google/uuid"
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
		documentHandle.currentPage.Close(p)

		// Cleanup refs.
		delete(documentHandle.pageRefs, documentHandle.currentPage.nativeRef)
		delete(p.pageRefs, documentHandle.currentPage.nativeRef)

		documentHandle.currentPage = nil
	}

	pageObject, err := p.Module.ExportedFunction("FPDF_LoadPage").Call(p.Context, *documentHandle.handle, uint64(page.ByIndex.Index))
	if err != nil {
		return nil, err
	}

	if len(pageObject) == 0 || pageObject[0] == 0 {
		return nil, pdfium_errors.ErrPage
	}

	nativePage := p.registerPage(pageObject[0], page.ByIndex.Index, documentHandle)

	documentHandle.currentPage = nativePage

	return nativePage, nil
}

func (p *PdfiumImplementation) registerPage(page uint64, index int, documentHandle *DocumentHandle) *PageHandle {
	pageRef := uuid.New()
	pageHandle := &PageHandle{
		handle:    &page,
		index:     index,
		nativeRef: references.FPDF_PAGE(pageRef.String()),
	}

	if documentHandle != nil {
		pageHandle.documentRef = documentHandle.nativeRef
		documentHandle.pageRefs[pageHandle.nativeRef] = pageHandle
	}

	p.pageRefs[pageHandle.nativeRef] = pageHandle

	return pageHandle
}
