package implementation_webassembly

import (
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDF_CreateNewDocument returns a new document.
func (p *PdfiumImplementation) FPDF_CreateNewDocument(request *requests.FPDF_CreateNewDocument) (*responses.FPDF_CreateNewDocument, error) {
	p.Lock()
	defer p.Unlock()

	res, err := p.Module.ExportedFunction("FPDF_CreateNewDocument").Call(p.Context)
	if err != nil {
		return nil, err
	}

	doc := &res[0]
	documentHandle := p.registerDocument(doc)

	return &responses.FPDF_CreateNewDocument{
		Document: documentHandle.nativeRef,
	}, nil
}
