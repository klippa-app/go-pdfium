//go:build pdfium_experimental
// +build pdfium_experimental

package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_doc.h"
// #include <stdlib.h>
import "C"
import (
	"bytes"
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFDest_GetView returns the view (fit type) for a given dest.
// Experimental API.
func (p *PdfiumImplementation) FPDFDest_GetView(request *requests.FPDFDest_GetView) (*responses.FPDFDest_GetView, error) {
	p.Lock()
	defer p.Unlock()

	destHandle, err := p.getDestHandle(request.Dest)
	if err != nil {
		return nil, err
	}

	numParams := C.ulong(0)
	params := make([]C.FS_FLOAT, 4, 4)

	destView := C.FPDFDest_GetView(destHandle.handle, &numParams, (*C.FS_FLOAT)(unsafe.Pointer(&params[0])))
	resParams := make([]float32, int(numParams), int(numParams))
	if int(numParams) > 0 {
		for i := range resParams {
			resParams[i] = float32(params[i])
		}
	}

	return &responses.FPDFDest_GetView{
		DestView: enums.FPDF_PDFDEST_VIEW(destView),
		Params:   resParams,
	}, nil
}

// FPDFLink_GetAnnot returns a FPDF_ANNOTATION object for a link.
// Experimental API.
func (p *PdfiumImplementation) FPDFLink_GetAnnot(request *requests.FPDFLink_GetAnnot) (*responses.FPDFLink_GetAnnot, error) {
	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	linkHandle, err := p.getLinkHandle(request.Link)
	if err != nil {
		return nil, err
	}

	annotation := C.FPDFLink_GetAnnot(pageHandle.handle, linkHandle.handle)
	if annotation == nil {
		return &responses.FPDFLink_GetAnnot{}, nil
	}

	annotationHandle := p.registerAnnotation(annotation)

	return &responses.FPDFLink_GetAnnot{
		Annotation: &annotationHandle.nativeRef,
	}, nil
}

// FPDF_GetPageAAction returns an additional-action from page.
// Experimental API
func (p *PdfiumImplementation) FPDF_GetPageAAction(request *requests.FPDF_GetPageAAction) (*responses.FPDF_GetPageAAction, error) {
	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	action := C.FPDF_GetPageAAction(pageHandle.handle, C.int(request.AAType))
	if action == nil {
		return &responses.FPDF_GetPageAAction{}, nil
	}

	actionHandle := p.registerAction(action)

	return &responses.FPDF_GetPageAAction{
		AAType: &request.AAType,
		Action: &actionHandle.nativeRef,
	}, nil
}

// FPDF_GetFileIdentifier Get the file identifier defined in the trailer of a document.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetFileIdentifier(request *requests.FPDF_GetFileIdentifier) (*responses.FPDF_GetFileIdentifier, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	if request.FileIdType != enums.FPDF_FILEIDTYPE_PERMANENT && request.FileIdType != enums.FPDF_FILEIDTYPE_CHANGING {
		return nil, errors.New("invalid file id type given")
	}

	// First get the identifier length.
	identifierSize := C.FPDF_GetFileIdentifier(documentHandle.handle, C.FPDF_FILEIDTYPE(request.FileIdType), C.NULL, 0)
	if identifierSize == 0 {
		return &responses.FPDF_GetFileIdentifier{
			FileIdType: request.FileIdType,
			Identifier: nil,
		}, nil
	}

	charData := make([]byte, identifierSize)
	C.FPDF_GetFileIdentifier(documentHandle.handle, C.FPDF_FILEIDTYPE(request.FileIdType), unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	// Remove NULL terminator.
	charData = bytes.TrimSuffix(charData, []byte("\x00"))

	return &responses.FPDF_GetFileIdentifier{
		FileIdType: request.FileIdType,
		Identifier: charData,
	}, nil
}
