package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_ext.h"
import "C"
import (
	"errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFDoc_GetPageMode returns the document's page mode, which describes how the document should be displayed when opened.
func (p *PdfiumImplementation) FPDFDoc_GetPageMode(request *requests.FPDFDoc_GetPageMode) (*responses.FPDFDoc_GetPageMode, error) {
	p.Lock()
	defer p.Unlock()

	nativeDoc, err := p.getNativeDocument(request.Document)
	if err != nil {
		return nil, err
	}

	if nativeDoc.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	pageMode := C.FPDFDoc_GetPageMode(nativeDoc.currentDoc)

	return &responses.FPDFDoc_GetPageMode{
		PageMode: responses.FPDFDoc_GetPageModeMode(pageMode),
	}, nil
}
