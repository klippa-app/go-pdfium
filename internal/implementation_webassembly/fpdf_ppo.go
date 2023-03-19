package implementation_webassembly

import (
	"errors"
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
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

	pageIndicesSize := uint64(len(request.PageIndices))
	pageIndicesPointer := uint64(0)

	if request.PageIndices != nil && len(request.PageIndices) > 0 {
		pageIndices, err := p.Malloc(pageIndicesSize * p.CSizeInt())
		if err != nil {
			return nil, err
		}

		for i := range request.PageIndices {
			success := p.Module.Memory().WriteUint32Le(uint32(pageIndices)+(uint32(i)*uint32(p.CSizeInt())), uint32(request.PageIndices[i]))
			if !success {
				return nil, errors.New("could not write page indices data to memory")
			}
		}

		pageIndicesPointer = pageIndices
	}

	res, err := p.Module.ExportedFunction("FPDF_ImportPagesByIndex").Call(p.Context, *destinationDocHandle.handle, *sourceDocHandle.handle, pageIndicesPointer, *(*uint64)(unsafe.Pointer(&pageIndicesSize)), *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDF_ImportNPagesToOne").Call(p.Context, *sourceDocHandle.handle, *(*uint64)(unsafe.Pointer(&request.OutputWidth)), *(*uint64)(unsafe.Pointer(&request.OutputHeight)), *(*uint64)(unsafe.Pointer(&request.NumPagesOnXAxis)), *(*uint64)(unsafe.Pointer(&request.NumPagesOnYAxis)))
	if err != nil {
		return nil, err
	}

	doc := res[0]
	if doc == 0 {
		return nil, errors.New("import of pages failed")
	}

	documentHandle := p.registerDocument(&doc)

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

	res, err := p.Module.ExportedFunction("FPDF_NewXObjectFromPage").Call(p.Context, *sourceDocHandle.handle, *destinationDocHandle.handle, *(*uint64)(unsafe.Pointer(&request.SourcePageIndex)))
	if err != nil {
		return nil, err
	}

	xObject := res[0]
	if xObject == 0 {
		return nil, errors.New("creation of xobject failed")
	}

	xObjectHandle := &XObjectHandle{}
	xObjectHandle.handle = &xObject
	xObjectRef := uuid.New()
	xObjectHandle.nativeRef = references.FPDF_XOBJECT(xObjectRef.String())
	p.xObjectRefs[xObjectHandle.nativeRef] = xObjectHandle

	return &responses.FPDF_NewXObjectFromPage{
		XObject: xObjectHandle.nativeRef,
	}, nil
}

// FPDF_CloseXObject closes an FPDF_XOBJECT handle created by FPDF_NewXObjectFromPage().
// Experimental API.
func (p *PdfiumImplementation) FPDF_CloseXObject(request *requests.FPDF_CloseXObject) (*responses.FPDF_CloseXObject, error) {
	p.Lock()
	defer p.Unlock()

	xObjectHandle, err := p.getXObjectHandle(request.XObject)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDF_CloseXObject").Call(p.Context, *xObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_CloseXObject{}, nil
}

// FPDF_NewFormObjectFromXObject creates a new form object from an FPDF_XOBJECT object.
// Experimental API.
func (p *PdfiumImplementation) FPDF_NewFormObjectFromXObject(request *requests.FPDF_NewFormObjectFromXObject) (*responses.FPDF_NewFormObjectFromXObject, error) {
	p.Lock()
	defer p.Unlock()

	xObjectHandle, err := p.getXObjectHandle(request.XObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDF_NewFormObjectFromXObject").Call(p.Context, *xObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	handle := res[0]

	pageObjectHandle := &PageObjectHandle{}
	pageObjectHandle.handle = &handle
	pageObjectRef := uuid.New()
	pageObjectHandle.nativeRef = references.FPDF_PAGEOBJECT(pageObjectRef.String())
	p.pageObjectRefs[pageObjectHandle.nativeRef] = pageObjectHandle

	return &responses.FPDF_NewFormObjectFromXObject{
		PageObject: pageObjectHandle.nativeRef,
	}, nil
}
