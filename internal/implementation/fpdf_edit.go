package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_edit.h"
import "C"
import (
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFPage_GetRotation returns the page rotation.
func (p *PdfiumImplementation) FPDFPage_GetRotation(request *requests.FPDFPage_GetRotation) (*responses.FPDFPage_GetRotation, error) {
	p.Lock()
	defer p.Unlock()

	nativePage, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	rotation := C.FPDFPage_GetRotation(nativePage.page)

	return &responses.FPDFPage_GetRotation{
		Page:         nativePage.index,
		PageRotation: responses.PageRotation(rotation),
	}, nil
}

// FPDFPage_SetRotation sets the page rotation for a given page.
func (p *PdfiumImplementation) FPDFPage_SetRotation(request *requests.FPDFPage_SetRotation) (*responses.FPDFPage_SetRotation, error) {
	p.Lock()
	defer p.Unlock()

	nativePage, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	C.FPDFPage_SetRotation(nativePage.page, C.int(request.Rotate))

	return &responses.FPDFPage_SetRotation{}, nil
}

// FPDFPage_HasTransparency returns whether the page has transparency.
func (p *PdfiumImplementation) FPDFPage_HasTransparency(request *requests.FPDFPage_HasTransparency) (*responses.FPDFPage_HasTransparency, error) {
	p.Lock()
	defer p.Unlock()

	nativePage, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	alpha := C.FPDFPage_HasTransparency(nativePage.page)

	return &responses.FPDFPage_HasTransparency{
		Page:            nativePage.index,
		HasTransparency: int(alpha) == 1,
	}, nil
}
