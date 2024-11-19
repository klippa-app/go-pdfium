package implementation_webassembly

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFCatalog_IsTagged determines if the given document represents a tagged PDF.
// For the definition of tagged PDF, See (see 10.7 "Tagged PDF" in PDF Reference 1.7).
// Experimental API.
func (p *PdfiumImplementation) FPDFCatalog_IsTagged(request *requests.FPDFCatalog_IsTagged) (*responses.FPDFCatalog_IsTagged, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFCatalog_IsTagged").Call(p.Context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	isTagged := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFCatalog_IsTagged{
		IsTagged: int(isTagged) == 1,
	}, nil
}

// FPDFCatalog_SetLanguage sets the language of a document.
// Experimental API.
func (p *PdfiumImplementation) FPDFCatalog_SetLanguage(request *requests.FPDFCatalog_SetLanguage) (*responses.FPDFCatalog_SetLanguage, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	languagePointer, err := p.CString(request.Language)
	if err != nil {
		return nil, err
	}
	defer languagePointer.Free()

	res, err := p.Module.ExportedFunction("FPDFCatalog_SetLanguage").Call(p.Context, *documentHandle.handle, languagePointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set language")
	}

	return &responses.FPDFCatalog_SetLanguage{}, nil
}
