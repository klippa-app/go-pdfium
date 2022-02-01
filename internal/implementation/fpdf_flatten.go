package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_flatten.h"
import "C"
import (
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFPage_Flatten makes annotations and form fields become part of the page contents itself.
func (p *PdfiumImplementation) FPDFPage_Flatten(request *requests.FPDFPage_Flatten) (*responses.FPDFPage_Flatten, error) {
	p.Lock()
	defer p.Unlock()

	nativePage, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	flattenPageResult := C.FPDFPage_Flatten(nativePage.page, C.int(request.Usage))

	return &responses.FPDFPage_Flatten{
		Page:   nativePage.index,
		Result: responses.FPDFPage_FlattenResult(flattenPageResult),
	}, nil
}
