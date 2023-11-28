//go:build pdfium_experimental
// +build pdfium_experimental

package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"
)

// FPDF_LoadMemDocument64 opens and load a PDF document from memory.
// Loaded document can be closed by FPDF_CloseDocument().
// If this function fails, you can use FPDF_GetLastError() to retrieve
// the reason why it failed.
// Experimental API.
func (p *PdfiumImplementation) FPDF_LoadMemDocument64(request *requests.FPDF_LoadMemDocument64) (*responses.FPDF_LoadMemDocument64, error) {
	// Don't lock, OpenDocument will do that.
	doc, err := p.OpenDocument(&requests.OpenDocument{
		File:     request.Data,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_LoadMemDocument64{
		Document: doc.Document,
	}, nil
}

// FPDF_DocumentHasValidCrossReferenceTable returns whether the document's cross reference table is valid or not.
// Experimental API.
func (p *PdfiumImplementation) FPDF_DocumentHasValidCrossReferenceTable(request *requests.FPDF_DocumentHasValidCrossReferenceTable) (*responses.FPDF_DocumentHasValidCrossReferenceTable, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	isValid := C.FPDF_DocumentHasValidCrossReferenceTable(documentHandle.handle)
	return &responses.FPDF_DocumentHasValidCrossReferenceTable{
		DocumentHasValidCrossReferenceTable: int(isValid) == 1,
	}, nil
}

// FPDF_GetTrailerEnds returns the byte offsets of trailer ends.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetTrailerEnds(request *requests.FPDF_GetTrailerEnds) (*responses.FPDF_GetTrailerEnds, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	trailerSize := C.FPDF_GetTrailerEnds(documentHandle.handle, nil, 0)
	if int(trailerSize) == 0 {
		return nil, errors.New("could not read trailer ends")
	}

	cTrailerEnds := make([]C.uint, int(trailerSize))
	readTrailers := C.FPDF_GetTrailerEnds(documentHandle.handle, (*C.uint)(unsafe.Pointer(&cTrailerEnds[0])), trailerSize)
	if int(readTrailers) == 0 {
		return nil, errors.New("could not read trailer ends")
	}

	trailerEnds := make([]int, trailerSize)
	for i := range cTrailerEnds {
		trailerEnds[i] = int(cTrailerEnds[i])
	}

	return &responses.FPDF_GetTrailerEnds{
		TrailerEnds: trailerEnds,
	}, nil
}

// FPDF_GetPageWidthF returns the page width in float32.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetPageWidthF(request *requests.FPDF_GetPageWidthF) (*responses.FPDF_GetPageWidthF, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	pageWidth := C.FPDF_GetPageWidthF(pageHandle.handle)
	return &responses.FPDF_GetPageWidthF{
		PageWidth: float32(pageWidth),
	}, nil
}

// FPDF_GetPageHeightF returns the page height in float32.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetPageHeightF(request *requests.FPDF_GetPageHeightF) (*responses.FPDF_GetPageHeightF, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	pageHeight := C.FPDF_GetPageHeightF(pageHandle.handle)
	return &responses.FPDF_GetPageHeightF{
		PageHeight: float32(pageHeight),
	}, nil
}

// FPDF_GetPageBoundingBox returns the bounding box of the page. This is the intersection between
// its media box and its crop box.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetPageBoundingBox(request *requests.FPDF_GetPageBoundingBox) (*responses.FPDF_GetPageBoundingBox, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	rect := C.FS_RECTF{}
	success := C.FPDF_GetPageBoundingBox(pageHandle.handle, &rect)
	if int(success) == 0 {
		return nil, errors.New("could not get page bounding box")
	}

	return &responses.FPDF_GetPageBoundingBox{
		Rect: structs.FPDF_FS_RECTF{
			Left:   float32(rect.left),
			Top:    float32(rect.top),
			Right:  float32(rect.right),
			Bottom: float32(rect.bottom),
		},
	}, nil
}

// FPDF_GetPageSizeByIndexF returns the size of the page at the given index.
// Prefer FPDF_GetPageSizeByIndexF(). This will be deprecated in the future.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetPageSizeByIndexF(request *requests.FPDF_GetPageSizeByIndexF) (*responses.FPDF_GetPageSizeByIndexF, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	size := C.FS_SIZEF{}
	success := C.FPDF_GetPageSizeByIndexF(documentHandle.handle, C.int(request.Index), &size)
	if int(success) == 0 {
		return nil, errors.New("could not get page size by index")
	}

	return &responses.FPDF_GetPageSizeByIndexF{
		Size: structs.FPDF_FS_SIZEF{
			Width:  float32(size.width),
			Height: float32(size.height),
		},
	}, nil
}

// FPDF_VIEWERREF_GetPrintPageRangeCount returns the number of elements in a FPDF_PAGERANGE.
// Experimental API.
func (p *PdfiumImplementation) FPDF_VIEWERREF_GetPrintPageRangeCount(request *requests.FPDF_VIEWERREF_GetPrintPageRangeCount) (*responses.FPDF_VIEWERREF_GetPrintPageRangeCount, error) {
	p.Lock()
	defer p.Unlock()

	pageRangeHandle, err := p.getPageRangeHandle(request.PageRange)
	if err != nil {
		return nil, err
	}

	count := C.FPDF_VIEWERREF_GetPrintPageRangeCount(pageRangeHandle.handle)
	return &responses.FPDF_VIEWERREF_GetPrintPageRangeCount{
		Count: uint64(count),
	}, nil
}

// FPDF_VIEWERREF_GetPrintPageRangeElement returns an element from a FPDF_PAGERANGE.
// Experimental API.
func (p *PdfiumImplementation) FPDF_VIEWERREF_GetPrintPageRangeElement(request *requests.FPDF_VIEWERREF_GetPrintPageRangeElement) (*responses.FPDF_VIEWERREF_GetPrintPageRangeElement, error) {
	p.Lock()
	defer p.Unlock()

	pageRangeHandle, err := p.getPageRangeHandle(request.PageRange)
	if err != nil {
		return nil, err
	}

	value := C.FPDF_VIEWERREF_GetPrintPageRangeElement(pageRangeHandle.handle, C.size_t(request.Index))
	if int(value) == -1 {
		return nil, errors.New("could not load page range element")
	}

	return &responses.FPDF_VIEWERREF_GetPrintPageRangeElement{
		Value: int(value),
	}, nil
}

// FPDF_GetXFAPacketCount returns the number of valid packets in the XFA entry.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetXFAPacketCount(request *requests.FPDF_GetXFAPacketCount) (*responses.FPDF_GetXFAPacketCount, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	count := C.FPDF_GetXFAPacketCount(documentHandle.handle)
	if int(count) == -1 {
		return nil, errors.New("error getting XFA packet count")
	}

	return &responses.FPDF_GetXFAPacketCount{
		Count: int(count),
	}, nil
}

// FPDF_GetXFAPacketName returns the name of a packet in the XFA array.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetXFAPacketName(request *requests.FPDF_GetXFAPacketName) (*responses.FPDF_GetXFAPacketName, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	// First get the name length.
	nameSize := C.FPDF_GetXFAPacketName(documentHandle.handle, C.int(request.Index), nil, 0)
	if uint64(nameSize) == 0 {
		return nil, errors.New("could not get name of the XFA packet")
	}

	charData := make([]byte, uint64(nameSize))
	result := C.FPDF_GetXFAPacketName(documentHandle.handle, C.int(request.Index), unsafe.Pointer(&charData[0]), C.ulong(len(charData)))
	if uint64(result) == 0 {
		return nil, errors.New("could not get name of the XFA packet")
	}

	return &responses.FPDF_GetXFAPacketName{
		Index: request.Index,
		Name:  string(charData[:len(charData)-1]),
	}, nil
}

// FPDF_GetXFAPacketContent returns the content of a packet in the XFA array.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetXFAPacketContent(request *requests.FPDF_GetXFAPacketContent) (*responses.FPDF_GetXFAPacketContent, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	outBufLen := C.ulong(0)

	// First get the name length.
	success := C.FPDF_GetXFAPacketContent(documentHandle.handle, C.int(request.Index), nil, 0, &outBufLen)
	if int(success) == 0 || uint64(outBufLen) == 0 {
		return nil, errors.New("could not get content of the XFA packet")
	}

	contentData := make([]byte, uint64(outBufLen))
	success = C.FPDF_GetXFAPacketContent(documentHandle.handle, C.int(request.Index), unsafe.Pointer(&contentData[0]), C.ulong(len(contentData)), &outBufLen)
	if int(success) == 0 {
		return nil, errors.New("could not get content of the XFA packet")
	}

	// Callers must check both the return value and the input |buflen| is no
	// less than the returned |out_buflen| before using the data in |buffer|.
	if uint64(len(contentData)) < uint64(outBufLen) {
		return nil, errors.New("could not get content of the XFA packet")
	}

	return &responses.FPDF_GetXFAPacketContent{
		Index:   request.Index,
		Content: contentData,
	}, nil
}

// FPDF_GetDocUserPermissions returns the user permissions of the PDF.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetDocUserPermissions(request *requests.FPDF_GetDocUserPermissions) (*responses.FPDF_GetDocUserPermissions, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	permissions := C.FPDF_GetDocUserPermissions(documentHandle.handle)

	docPermissions := &responses.FPDF_GetDocUserPermissions{
		DocUserPermissions: uint32(permissions),
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
		if docPermissions.DocUserPermissions&permission > 0 {
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
