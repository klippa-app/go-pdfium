//go:build pdfium_experimental
// +build pdfium_experimental

package implementation_cgo

/*
#cgo pkg-config: pdfium
#include "fpdf_catalog.h"
#include <stdlib.h>
*/
import "C"
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

	isTagged := C.FPDFCatalog_IsTagged(documentHandle.handle)

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

	languageStr := C.CString(request.Language)
	defer C.free(unsafe.Pointer(languageStr))

	success := C.FPDFCatalog_SetLanguage(documentHandle.handle, languageStr)
	if int(success) == 0 {
		return nil, errors.New("could not set language")
	}

	return &responses.FPDFCatalog_SetLanguage{}, nil
}
