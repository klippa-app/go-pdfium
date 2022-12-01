package implementation_webassembly

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

	var pageRange = uint64(0)
	if request.PageRange != nil {
		pageRangePointer, err := p.CString(*request.PageRange)
		if err != nil {
			return nil, err
		}
		defer pageRangePointer.Free()
		pageRange = pageRangePointer.Pointer
	}

	res, err := p.Module.ExportedFunction("FPDF_ImportPages").Call(p.Context, *destinationDocHandle.handle, *sourceDocHandle.handle, pageRange, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDF_CopyViewerPreferences").Call(p.Context, *destinationDocHandle.handle, *sourceDocHandle.handle)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("copy of viewer preferences failed")
	}

	return &responses.FPDF_CopyViewerPreferences{}, nil
}
