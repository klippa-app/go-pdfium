package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_edit.h"
import "C"
import (
	"errors"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// SetRotation sets the page rotation for a given page.
func (p *PdfiumImplementation) SetRotation(request *requests.SetRotation) (*responses.SetRotation, error) {
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

	C.FPDFPage_SetRotation(nativePage.page, C.int(request.Rotate))

	return &responses.SetRotation{}, nil
}
