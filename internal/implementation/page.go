package implementation

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_edit.h"
// #include "fpdf_flatten.h"
import "C"
import (
	"errors"

	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// loadPage changes the active page if it's different from what's currently
// open and closes the page that's currently open if any is open.
func (p *PdfiumImplementation) loadPage(doc *NativeDocument, page int) error {
	// Already loaded this page.
	if doc.currentPageNumber != nil && *doc.currentPageNumber == page {
		return nil
	}

	if doc.currentPageNumber != nil {
		// Unload the current page.
		C.FPDF_ClosePage(doc.currentPage)
		doc.currentPage = nil
		doc.currentPageNumber = nil
	}

	pageObject := C.FPDF_LoadPage(doc.currentDoc, C.int(page))
	if pageObject == nil {
		return pdfium_errors.ErrPage
	}

	doc.currentPage = pageObject
	doc.currentPageNumber = &page

	return nil
}

// GetPageRotation returns the page rotation.
func (p *PdfiumImplementation) GetPageRotation(request *requests.GetPageRotation) (*responses.GetPageRotation, error) {
	p.Lock()
	defer p.Unlock()

	nativeDoc, err := p.getNativeDocument(request.Document)
	if err != nil {
		return nil, err
	}

	if nativeDoc.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	err = p.loadPage(nativeDoc, request.Page)
	if err != nil {
		return nil, err
	}

	rotation := C.FPDFPage_GetRotation(nativeDoc.currentPage)

	return &responses.GetPageRotation{
		Page:         request.Page,
		PageRotation: responses.PageRotation(rotation),
	}, nil
}

// GetPageTransparency returns whether the page has transparency.
func (p *PdfiumImplementation) GetPageTransparency(request *requests.GetPageTransparency) (*responses.GetPageTransparency, error) {
	p.Lock()
	defer p.Unlock()

	nativeDoc, err := p.getNativeDocument(request.Document)
	if err != nil {
		return nil, err
	}

	if nativeDoc.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	err = p.loadPage(nativeDoc, request.Page)
	if err != nil {
		return nil, err
	}

	alpha := C.FPDFPage_HasTransparency(nativeDoc.currentPage)
	if int(alpha) == 1 {
		return &responses.GetPageTransparency{
			Page:            request.Page,
			HasTransparency: true,
		}, nil
	}

	return &responses.GetPageTransparency{
		Page:            request.Page,
		HasTransparency: false,
	}, nil
}

// FlattenPage makes annotations and form fields become part of the page contents itself.
func (p *PdfiumImplementation) FlattenPage(request *requests.FlattenPage) (*responses.FlattenPage, error) {
	p.Lock()
	defer p.Unlock()

	nativeDoc, err := p.getNativeDocument(request.Document)
	if err != nil {
		return nil, err
	}

	if nativeDoc.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	err = p.loadPage(nativeDoc, request.Page)
	if err != nil {
		return nil, err
	}

	flattenPageResult := C.FPDFPage_Flatten(nativeDoc.currentPage, C.int(request.Usage))

	return &responses.FlattenPage{
		Page:   request.Page,
		Result: responses.FlattenPageResult(flattenPageResult),
	}, nil
}
