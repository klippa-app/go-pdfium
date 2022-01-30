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

// GetPageRotation returns the page rotation.
func (p *Pdfium) GetPageRotation(request *requests.GetPageRotation) (*responses.GetPageRotation, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	rotation := C.FPDFPage_GetRotation(p.currentPage)

	return &responses.GetPageRotation{
		Page:         request.Page,
		PageRotation: responses.PageRotation((rotation)),
	}, nil
}

// GetPageTransparency returns whether the page has transparency.
func (p *Pdfium) GetPageTransparency(request *requests.GetPageTransparency) (*responses.GetPageTransparency, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	alpha := C.FPDFPage_HasTransparency(p.currentPage)
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
func (p *Pdfium) FlattenPage(request *requests.FlattenPage) (*responses.FlattenPage, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	flattenPageResult := C.FPDFPage_Flatten(p.currentPage, C.int(request.Usage))

	return &responses.FlattenPage{
		Page:   request.Page,
		Result: responses.FlattenPageResult(flattenPageResult),
	}, nil
}
