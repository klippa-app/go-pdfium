package implementation

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_doc.h"
// #include "fpdf_ext.h"
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

	docPermissions := &responses.GetDocPermissions{
		DocPermissions: uint32(permissions),
	}

	PrintDocument := uint32(1 << 2)
	ModifyContents := uint32(1 << 3)
	CopyOrExtractText := uint32(1 << 4)
	AddOrModifyTextAnnotations := uint32(1 << 5)
	FillInExistingInteractiveFormFields := uint32(1 << 8)
	ExtractTextAndGraphics := uint32(1 << 9)
	AssembleDocument := uint32(1 << 10)
	PrintDocumentAsFaithfulDigitalCopy := uint32(1 << 11)

	hasPermission := func(permission uint32) bool {
		if docPermissions.DocPermissions&permission > 0 {
			return true
		}

		return false
	}

	docPermissions.PrintDocument = hasPermission(PrintDocument)
	docPermissions.ModifyContents = hasPermission(ModifyContents)
	docPermissions.CopyOrExtractText = hasPermission(CopyOrExtractText)
	docPermissions.AddOrModifyTextAnnotations = hasPermission(AddOrModifyTextAnnotations)
	docPermissions.FillInInteractiveFormFields = hasPermission(AddOrModifyTextAnnotations)
	docPermissions.FillInExistingInteractiveFormFields = hasPermission(FillInExistingInteractiveFormFields)
	docPermissions.ExtractTextAndGraphics = hasPermission(ExtractTextAndGraphics)
	docPermissions.AssembleDocument = hasPermission(AssembleDocument)
	docPermissions.PrintDocumentAsFaithfulDigitalCopy = hasPermission(PrintDocumentAsFaithfulDigitalCopy)

	// Calculated permissions
	docPermissions.CreateOrModifyInteractiveFormFields = docPermissions.ModifyContents && docPermissions.AddOrModifyTextAnnotations

	return docPermissions, nil
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

// GetPageMode returns the document's page mode, which describes how the document should be displayed when opened.
func (p *Pdfium) GetPageMode(request *requests.GetPageMode) (*responses.GetPageMode, error) {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	pageMode := C.FPDFDoc_GetPageMode(p.currentDoc)

	return &responses.GetPageMode{
		PageMode: responses.PageMode(pageMode),
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
