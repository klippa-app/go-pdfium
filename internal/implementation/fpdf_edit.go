package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_edit.h"
import "C"
import (
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDF_CreateNewDocument returns a new document.
func (p *PdfiumImplementation) FPDF_CreateNewDocument(request *requests.FPDF_CreateNewDocument) (*responses.FPDF_CreateNewDocument, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle := &DocumentHandle{}
	doc := C.FPDF_CreateNewDocument()
	documentHandle.handle = doc
	documentRef := uuid.New()
	documentHandle.nativeRef = references.FPDF_DOCUMENT(documentRef.String())
	Pdfium.documentRefs[documentHandle.nativeRef] = documentHandle
	p.documentRefs[documentHandle.nativeRef] = documentHandle

	return &responses.FPDF_CreateNewDocument{
		Document: documentHandle.nativeRef,
	}, nil
}

// FPDFPage_GetRotation returns the page rotation.
func (p *PdfiumImplementation) FPDFPage_GetRotation(request *requests.FPDFPage_GetRotation) (*responses.FPDFPage_GetRotation, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	rotation := C.FPDFPage_GetRotation(pageHandle.handle)

	return &responses.FPDFPage_GetRotation{
		Page:         pageHandle.index,
		PageRotation: responses.PageRotation(rotation),
	}, nil
}

// FPDFPage_SetRotation sets the page rotation for a given page.
func (p *PdfiumImplementation) FPDFPage_SetRotation(request *requests.FPDFPage_SetRotation) (*responses.FPDFPage_SetRotation, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	C.FPDFPage_SetRotation(pageHandle.handle, C.int(request.Rotate))

	return &responses.FPDFPage_SetRotation{}, nil
}

// FPDFPage_HasTransparency returns whether the page has transparency.
func (p *PdfiumImplementation) FPDFPage_HasTransparency(request *requests.FPDFPage_HasTransparency) (*responses.FPDFPage_HasTransparency, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	alpha := C.FPDFPage_HasTransparency(pageHandle.handle)

	return &responses.FPDFPage_HasTransparency{
		Page:            pageHandle.index,
		HasTransparency: int(alpha) == 1,
	}, nil
}
