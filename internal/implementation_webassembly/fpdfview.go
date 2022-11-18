package implementation_webassembly

import (
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

func (p *PdfiumImplementation) FPDF_GetPageCount(request *requests.FPDF_GetPageCount) (*responses.FPDF_GetPageCount, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	res, err := p.module.ExportedFunction("FPDF_GetPageCount").Call(p.context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_GetPageCount{
		PageCount: int(res[0]),
	}, nil
}
