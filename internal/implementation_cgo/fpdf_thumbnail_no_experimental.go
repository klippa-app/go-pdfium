//go:build !pdfium_experimental
// +build !pdfium_experimental

package implementation_cgo

import (
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFPage_GetDecodedThumbnailData returns the decoded data from the thumbnail of the given page if it exists.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_GetDecodedThumbnailData(request *requests.FPDFPage_GetDecodedThumbnailData) (*responses.FPDFPage_GetDecodedThumbnailData, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFPage_GetRawThumbnailData returns the raw data from the thumbnail of the given page if it exists.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_GetRawThumbnailData(request *requests.FPDFPage_GetRawThumbnailData) (*responses.FPDFPage_GetRawThumbnailData, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFPage_GetThumbnailAsBitmap returns the thumbnail of the given page as a FPDF_BITMAP.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_GetThumbnailAsBitmap(request *requests.FPDFPage_GetThumbnailAsBitmap) (*responses.FPDFPage_GetThumbnailAsBitmap, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}
