package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdf_ppo.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDF_ImportPages imports some pages from one PDF document to another one.
func (p *PdfiumImplementation) FPDF_ImportPages(request *requests.FPDF_ImportPages) (*responses.FPDF_ImportPages, error) {
	p.Lock()
	defer p.Unlock()

	sourceDocHandle, err := p.getDocumentHandle(request.Source)
	if err != nil {
		return nil, err
	}

	destinationDocHandle, err := p.getDocumentHandle(request.Destination)
	if err != nil {
		return nil, err
	}

	var pageRange *C.char
	if request.PageRange != nil {
		pageRange = C.CString(*request.PageRange)
		defer C.free(unsafe.Pointer(pageRange))
	}

	success := C.FPDF_ImportPages(destinationDocHandle.handle, sourceDocHandle.handle, pageRange, C.int(request.Index))
	if int(success) == 0 {
		return nil, errors.New("import of pages failed")
	}

	return &responses.FPDF_ImportPages{}, nil
}

// FPDF_CopyViewerPreferences copies the viewer preferences from one PDF document to another
func (p *PdfiumImplementation) FPDF_CopyViewerPreferences(request *requests.FPDF_CopyViewerPreferences) (*responses.FPDF_CopyViewerPreferences, error) {
	p.Lock()
	defer p.Unlock()

	sourceDocHandle, err := p.getDocumentHandle(request.Source)
	if err != nil {
		return nil, err
	}

	destinationDocHandle, err := p.getDocumentHandle(request.Destination)
	if err != nil {
		return nil, err
	}

	success := C.FPDF_CopyViewerPreferences(destinationDocHandle.handle, sourceDocHandle.handle)
	if int(success) == 0 {
		return nil, errors.New("copy of viewer preferences failed")
	}

	return &responses.FPDF_CopyViewerPreferences{}, nil
}
