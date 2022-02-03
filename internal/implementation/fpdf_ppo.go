package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_ppo.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/google/uuid"
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

// FPDF_ImportPagesByIndex imports pages to a FPDF_DOCUMENT.
// Experimental API.
func (p *PdfiumImplementation) FPDF_ImportPagesByIndex(request *requests.FPDF_ImportPagesByIndex) (*responses.FPDF_ImportPagesByIndex, error) {
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

	var pageIndices *C.int
	if request.PageIndices != nil && len(request.PageIndices) > 0 {
		params := make([]C.int, len(request.PageIndices), len(request.PageIndices))
		for i := range params {
			params[i] = C.int(request.PageIndices[i])
		}

		pageIndices = (*C.int)(unsafe.Pointer(&params[0]))
	}

	success := C.FPDF_ImportPagesByIndex(destinationDocHandle.handle, sourceDocHandle.handle, pageIndices, C.ulong(len(request.PageIndices)), C.int(request.Index))
	if int(success) == 0 {
		return nil, errors.New("import of pages failed")
	}

	return &responses.FPDF_ImportPagesByIndex{}, nil
}

// FPDF_ImportNPagesToOne creates a new document from source document. The pages of source document will be
// combined to provide NumPagesOnXAxis x NumPagesOnYAxis pages per page of the output document.
// Experimental API.
func (p *PdfiumImplementation) FPDF_ImportNPagesToOne(request *requests.FPDF_ImportNPagesToOne) (*responses.FPDF_ImportNPagesToOne, error) {
	p.Lock()
	defer p.Unlock()

	sourceDocHandle, err := p.getDocumentHandle(request.Source)
	if err != nil {
		return nil, err
	}

	doc := C.FPDF_ImportNPagesToOne(sourceDocHandle.handle, C.float(request.OutputWidth), C.float(request.OutputHeight), C.size_t(request.NumPagesOnXAxis), C.size_t(request.NumPagesOnYAxis))
	if doc == nil {
		return nil, errors.New("import of pages failed")
	}

	documentHandle := &DocumentHandle{}
	documentHandle.handle = doc
	documentRef := uuid.New()
	documentHandle.nativeRef = references.FPDF_DOCUMENT(documentRef.String())
	Pdfium.documentRefs[documentHandle.nativeRef] = documentHandle
	p.documentRefs[documentHandle.nativeRef] = documentHandle

	return &responses.FPDF_ImportNPagesToOne{
		Document: documentHandle.nativeRef,
	}, nil
}

// FPDF_NewXObjectFromPage creates a template to generate form xobjects from the source document's page at
// the given index, for use in the destination document.
// Experimental API.
func (p *PdfiumImplementation) FPDF_NewXObjectFromPage(request *requests.FPDF_NewXObjectFromPage) (*responses.FPDF_NewXObjectFromPage, error) {
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

	xObject := C.FPDF_NewXObjectFromPage(sourceDocHandle.handle, destinationDocHandle.handle, C.int(request.SourcePageIndex))
	if xObject == nil {
		return nil, errors.New("creation of xobject failed")
	}

	xObjectHandle := &XObjectHandle{}
	xObjectHandle.handle = xObject
	xObjectRef := uuid.New()
	xObjectHandle.nativeRef = references.FPDF_XOBJECT(xObjectRef.String())
	p.xObjectRefs[xObjectHandle.nativeRef] = xObjectHandle

	return &responses.FPDF_NewXObjectFromPage{
		XObject: xObjectHandle.nativeRef,
	}, nil
}

// FPDF_CloseXObject closes an FPDF_XOBJECT handle created by FPDF_NewXObjectFromPage().
func (p *PdfiumImplementation) FPDF_CloseXObject(request *requests.FPDF_CloseXObject) (*responses.FPDF_CloseXObject, error) {
	p.Lock()
	defer p.Unlock()

	xObjectHandle, err := p.getXObjectHandle(request.XObject)
	if err != nil {
		return nil, err
	}

	C.FPDF_CloseXObject(xObjectHandle.handle)

	return &responses.FPDF_CloseXObject{}, nil
}

// FPDF_NewFormObjectFromXObject creates a new form object from an FPDF_XOBJECT object.
func (p *PdfiumImplementation) FPDF_NewFormObjectFromXObject(request *requests.FPDF_NewFormObjectFromXObject) (*responses.FPDF_NewFormObjectFromXObject, error) {
	p.Lock()
	defer p.Unlock()

	xObjectHandle, err := p.getXObjectHandle(request.XObject)
	if err != nil {
		return nil, err
	}

	handle := C.FPDF_NewFormObjectFromXObject(xObjectHandle.handle)

	pageObjectHandle := &PageObjectHandle{}
	pageObjectHandle.handle = handle
	pageObjectRef := uuid.New()
	pageObjectHandle.nativeRef = references.FPDF_PAGEOBJECT(pageObjectRef.String())
	p.pageObjectRefs[pageObjectHandle.nativeRef] = pageObjectHandle

	return &responses.FPDF_NewFormObjectFromXObject{
		PageObject: pageObjectHandle.nativeRef,
	}, nil

}
