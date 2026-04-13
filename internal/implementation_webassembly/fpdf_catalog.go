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

// FPDFCatalog_GetLanguage gets the language of a document from the catalog's /Lang entry.
// Experimental API.
func (p *PdfiumImplementation) FPDFCatalog_GetLanguage(request *requests.FPDFCatalog_GetLanguage) (*responses.FPDFCatalog_GetLanguage, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFCatalog_GetLanguage").Call(p.Context, *documentHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	langLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if langLength == 0 {
		return nil, errors.New("could not get language")
	}

	charDataPointer, err := p.ByteArrayPointer(langLength, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFCatalog_GetLanguage").Call(p.Context, *documentHandle.handle, charDataPointer.Pointer, langLength)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFCatalog_GetLanguage{
		Language: transformedText,
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

	languagePointer, err := p.CFPDF_WIDESTRING(request.Language)
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
