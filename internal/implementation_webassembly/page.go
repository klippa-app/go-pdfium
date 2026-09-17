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

	pageObject, err := p.call("FPDF_LoadPage", *documentHandle.handle, uint64(page.ByIndex.Index))
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

// reloadPage replaces the FPDF_PAGE behind a page handle with a freshly
// loaded one, keeping the reference the same for the caller.
//
// PDFium parses the page content when the page is loaded and never again.
// Operations that rewrite the page dictionary, such as FPDFPage_Flatten,
// therefore do not show up when rendering or extracting text from the same
// FPDF_PAGE; only a newly loaded page sees them. Form fill environments that
// have this page loaded are told about the switch so that their page views
// follow the new page. The old page is kept alive until the handle is closed,
// see PageHandle.stalePages.
func (p *PdfiumImplementation) reloadPage(pageHandle *PageHandle) error {
	if pageHandle.index < 0 || pageHandle.documentRef == "" {
		return errors.New("page can't be reloaded, its index or document is unknown")
	}

	documentHandle, err := p.getDocumentHandle(pageHandle.documentRef)
	if err != nil {
		return err
	}

	res, err := p.call("FPDF_LoadPage", *documentHandle.handle, uint64(pageHandle.index))
	if err != nil {
		return err
	}

	if len(res) == 0 || res[0] == 0 {
		return pdfium_errors.ErrPage
	}

	newPage := res[0]
	oldPage := *pageHandle.handle
	for _, formHandleHandle := range documentHandle.formHandleRefs {
		if _, ok := formHandleHandle.pagePointers[oldPage]; !ok {
			continue
		}

		_, err = p.call("FORM_OnBeforeClosePage", oldPage, *formHandleHandle.handle)
		if err != nil {
			return err
		}
		delete(formHandleHandle.pagePointers, oldPage)

		_, err = p.call("FORM_OnAfterLoadPage", newPage, *formHandleHandle.handle)
		if err != nil {
			return err
		}
		formHandleHandle.pagePointers[newPage] = pageHandle.nativeRef
	}

	pageHandle.stalePages = append(pageHandle.stalePages, oldPage)
	pageHandle.handle = &newPage

	return nil
}
