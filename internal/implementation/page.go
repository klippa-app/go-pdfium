package implementation

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_edit.h"
// #include "fpdf_flatten.h"
import "C"
import (
	"errors"

	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/google/uuid"
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

// LoadPage loads a page and returns a reference.
func (p *PdfiumImplementation) LoadPage(request *requests.LoadPage) (*responses.LoadPage, error) {
	p.Lock()
	defer p.Unlock()

	nativeDoc, err := p.getNativeDocument(request.Document)
	if err != nil {
		return nil, err
	}

	if nativeDoc.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	pageObject := C.FPDF_LoadPage(nativeDoc.currentDoc, C.int(request.Index))
	if pageObject == nil {
		return nil, pdfium_errors.ErrPage
	}

	pageRef := uuid.New()
	nativePage := &NativePage{
		page:      pageObject,
		index:     request.Index,
		nativeRef: references.Page(pageRef.String()),
	}

	nativeDoc.pageRefs[nativePage.nativeRef] = nativePage

	return &responses.LoadPage{
		Page: nativePage.nativeRef,
	}, nil
}

// UnloadPage unloads a page by reference.
func (p *PdfiumImplementation) UnloadPage(request *requests.UnloadPage) (*responses.UnloadPage, error) {
	p.Lock()
	defer p.Unlock()

	nativeDoc, err := p.getNativeDocument(request.Document)
	if err != nil {
		return nil, err
	}

	if nativeDoc.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	pageRef, err := nativeDoc.getNativePage(request.Page)
	if err != nil {
		return nil, err
	}

	pageRef.Close()
	delete(nativeDoc.pageRefs, request.Page)

	return &responses.UnloadPage{}, nil
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

	nativePage, err := p.loadPage(nativeDoc, request.Page)
	if err != nil {
		return nil, err
	}

	rotation := C.FPDFPage_GetRotation(nativePage.page)

	return &responses.GetPageRotation{
		Page:         nativePage.index,
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

	nativePage, err := p.loadPage(nativeDoc, request.Page)
	if err != nil {
		return nil, err
	}

	alpha := C.FPDFPage_HasTransparency(nativePage.page)
	if int(alpha) == 1 {
		return &responses.GetPageTransparency{
			Page:            nativePage.index,
			HasTransparency: true,
		}, nil
	}

	return &responses.GetPageTransparency{
		Page:            nativePage.index,
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

	nativePage, err := p.loadPage(nativeDoc, request.Page)
	if err != nil {
		return nil, err
	}

	flattenPageResult := C.FPDFPage_Flatten(nativePage.page, C.int(request.Usage))

	return &responses.FlattenPage{
		Page:   nativePage.index,
		Result: responses.FlattenPageResult(flattenPageResult),
	}, nil
}
