package implementation_webassembly

import (
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"unsafe"
)

// FPDFPage_Flatten makes annotations and form fields become part of the page contents itself.
// On success the page is reloaded behind the given reference, so that
// subsequent rendering and text extraction through the same reference (or
// the same index) reflect the flattened content.
func (p *PdfiumImplementation) FPDFPage_Flatten(request *requests.FPDFPage_Flatten) (*responses.FPDFPage_Flatten, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.call("FPDFPage_Flatten", *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Usage)))
	if err != nil {
		return nil, err
	}

	flattenPageResult := *(*int32)(unsafe.Pointer(&res[0]))

	// Flattening rewrites the page dictionary, but the loaded page keeps its
	// already parsed content. Reload it so that the same page reference
	// renders and extracts the flattened result.
	if responses.FPDFPage_FlattenResult(flattenPageResult) == responses.FPDFPage_FlattenResultSuccess {
		if err := p.reloadPage(pageHandle); err != nil {
			return nil, err
		}
	}

	return &responses.FPDFPage_Flatten{
		Page:   pageHandle.index,
		Result: responses.FPDFPage_FlattenResult(flattenPageResult),
	}, nil
}
