//go:build !pdfium_experimental
// +build !pdfium_experimental

package implementation

import (
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDF_LoadMemDocument64 opens and load a PDF document from memory.
// Loaded document can be closed by FPDF_CloseDocument().
// If this function fails, you can use FPDF_GetLastError() to retrieve
// the reason why it failed.
// Experimental API.
func (p *PdfiumImplementation) FPDF_LoadMemDocument64(request *requests.FPDF_LoadMemDocument64) (*responses.FPDF_LoadMemDocument64, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_DocumentHasValidCrossReferenceTable returns whether the document's cross reference table is valid or not.
// Experimental API.
func (p *PdfiumImplementation) FPDF_DocumentHasValidCrossReferenceTable(request *requests.FPDF_DocumentHasValidCrossReferenceTable) (*responses.FPDF_DocumentHasValidCrossReferenceTable, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_GetTrailerEnds returns the byte offsets of trailer ends.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetTrailerEnds(request *requests.FPDF_GetTrailerEnds) (*responses.FPDF_GetTrailerEnds, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_GetPageWidthF returns the page width in float32.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetPageWidthF(request *requests.FPDF_GetPageWidthF) (*responses.FPDF_GetPageWidthF, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_GetPageHeightF returns the page height in float32.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetPageHeightF(request *requests.FPDF_GetPageHeightF) (*responses.FPDF_GetPageHeightF, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_GetPageBoundingBox returns the bounding box of the page. This is the intersection between
// its media box and its crop box.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetPageBoundingBox(request *requests.FPDF_GetPageBoundingBox) (*responses.FPDF_GetPageBoundingBox, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_GetPageSizeByIndexF returns the size of the page at the given index.
// Prefer FPDF_GetPageSizeByIndexF(). This will be deprecated in the future.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetPageSizeByIndexF(request *requests.FPDF_GetPageSizeByIndexF) (*responses.FPDF_GetPageSizeByIndexF, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_VIEWERREF_GetPrintPageRangeCount returns the number of elements in a FPDF_PAGERANGE.
// Experimental API.
func (p *PdfiumImplementation) FPDF_VIEWERREF_GetPrintPageRangeCount(request *requests.FPDF_VIEWERREF_GetPrintPageRangeCount) (*responses.FPDF_VIEWERREF_GetPrintPageRangeCount, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_VIEWERREF_GetPrintPageRangeElement returns an element from a FPDF_PAGERANGE.
// Experimental API.
func (p *PdfiumImplementation) FPDF_VIEWERREF_GetPrintPageRangeElement(request *requests.FPDF_VIEWERREF_GetPrintPageRangeElement) (*responses.FPDF_VIEWERREF_GetPrintPageRangeElement, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_GetXFAPacketCount returns the number of valid packets in the XFA entry.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetXFAPacketCount(request *requests.FPDF_GetXFAPacketCount) (*responses.FPDF_GetXFAPacketCount, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_GetXFAPacketName returns the name of a packet in the XFA array.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetXFAPacketName(request *requests.FPDF_GetXFAPacketName) (*responses.FPDF_GetXFAPacketName, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDF_GetXFAPacketContent returns the content of a packet in the XFA array.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetXFAPacketContent(request *requests.FPDF_GetXFAPacketContent) (*responses.FPDF_GetXFAPacketContent, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}
