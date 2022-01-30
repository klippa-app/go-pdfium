package implementation

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_doc.h"
// #include <stdlib.h>
import "C"

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// GetFileVersion returns the version of the PDF file.
func (p *Pdfium) GetFileVersion(request *requests.GetFileVersion) (*responses.GetFileVersion, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	fileVersion := C.int(0)

	success := C.FPDF_GetFileVersion(p.currentDoc, &fileVersion)
	if int(success) == 0 {
		return nil, errors.New("could not get file version")
	}

	return &responses.GetFileVersion{
		FileVersion: int(fileVersion),
	}, nil
}

// GetDocPermissions returns the permissions of the PDF.
func (p *Pdfium) GetDocPermissions(request *requests.GetDocPermissions) (*responses.GetDocPermissions, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	permissions := C.FPDF_GetDocPermissions(p.currentDoc)
	return &responses.GetDocPermissions{
		DocPermissions: uint32(permissions),
	}, nil
}

// GetSecurityHandlerRevision returns the revision number of security handlers of the file.
func (p *Pdfium) GetSecurityHandlerRevision(request *requests.GetSecurityHandlerRevision) (*responses.GetSecurityHandlerRevision, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	securityHandlerRevision := C.FPDF_GetSecurityHandlerRevision(p.currentDoc)

	return &responses.GetSecurityHandlerRevision{
		SecurityHandlerRevision: int(securityHandlerRevision),
	}, nil
}

// GetPageCount counts the amount of pages.
func (p *Pdfium) GetPageCount(request *requests.GetPageCount) (*responses.GetPageCount, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	return &responses.GetPageCount{
		PageCount: int(C.FPDF_GetPageCount(p.currentDoc)),
	}, nil
}

// GetMetadata returns the requested metadata.
func (p *Pdfium) GetMetadata(request *requests.GetMetadata) (*responses.GetMetadata, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	cstr := C.CString(request.Tag)
	defer C.free(unsafe.Pointer(cstr))

	// First get the metadata length.
	metaSize := C.FPDF_GetMetaText(p.currentDoc, cstr, C.NULL, 0)
	if metaSize == 0 {
		return nil, errors.New("Could not get metadata")
	}

	charData := make([]byte, metaSize)
	C.FPDF_GetMetaText(p.currentDoc, cstr, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEText(charData)
	if err != nil {
		return nil, err
	}

	return &responses.GetMetadata{
		Tag:   request.Tag,
		Value: transformedText,
	}, nil
}
