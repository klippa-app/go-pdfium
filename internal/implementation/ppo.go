package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_ppo.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"unsafe"
)

// ImportPages imports some pages from one PDF document to another one.
func (p *PdfiumImplementation) ImportPages(request *requests.ImportPages) (*responses.ImportPages, error) {
	p.Lock()
	defer p.Unlock()

	nativeSourceDoc, err := p.getNativeDocument(request.Source)
	if err != nil {
		return nil, err
	}

	if nativeSourceDoc.currentDoc == nil {
		return nil, errors.New("no source document")
	}

	nativeDestinationDoc, err := p.getNativeDocument(request.Destination)
	if err != nil {
		return nil, err
	}

	if nativeDestinationDoc.currentDoc == nil {
		return nil, errors.New("no destination document")
	}

	var pageRange *C.char
	if request.PageRange != nil {
		pageRange = C.CString(*request.PageRange)
		defer C.free(unsafe.Pointer(pageRange))
	}

	success := C.FPDF_ImportPages(nativeDestinationDoc.currentDoc, nativeSourceDoc.currentDoc, pageRange, C.int(request.Index))
	if int(success) == 0 {
		return nil, errors.New("import of pages failed")
	}

	return &responses.ImportPages{}, nil
}

// CopyViewerPreferences copies the viewer preferences from one PDF document to another
func (p *PdfiumImplementation) CopyViewerPreferences(request *requests.CopyViewerPreferences) (*responses.CopyViewerPreferences, error) {
	p.Lock()
	defer p.Unlock()

	nativeSourceDoc, err := p.getNativeDocument(request.Source)
	if err != nil {
		return nil, err
	}

	if nativeSourceDoc.currentDoc == nil {
		return nil, errors.New("no source document")
	}

	nativeDestinationDoc, err := p.getNativeDocument(request.Destination)
	if err != nil {
		return nil, err
	}

	if nativeDestinationDoc.currentDoc == nil {
		return nil, errors.New("no destination document")
	}

	success := C.FPDF_CopyViewerPreferences(nativeDestinationDoc.currentDoc, nativeSourceDoc.currentDoc)
	if int(success) == 0 {
		return nil, errors.New("import of pages failed")
	}

	return &responses.CopyViewerPreferences{}, nil
}
