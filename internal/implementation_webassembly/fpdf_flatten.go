package implementation_webassembly

import (
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"unsafe"
)

// FPDFPage_Flatten makes annotations and form fields become part of the page contents itself.
func (p *PdfiumImplementation) FPDFPage_Flatten(request *requests.FPDFPage_Flatten) (*responses.FPDFPage_Flatten, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPage_Flatten").Call(p.Context, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Usage)))
	if err != nil {
		return nil, err
	}

	flattenPageResult := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFPage_Flatten{
		Page:   pageHandle.index,
		Result: responses.FPDFPage_FlattenResult(flattenPageResult),
	}, nil
}
