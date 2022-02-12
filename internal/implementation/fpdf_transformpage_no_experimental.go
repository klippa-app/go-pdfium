//go:build !pdfium_experimental
// +build !pdfium_experimental

package implementation

import (
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFPageObj_GetClipPath Get the clip path of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetClipPath(request *requests.FPDFPageObj_GetClipPath) (*responses.FPDFPageObj_GetClipPath, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFClipPath_CountPaths returns the number of paths inside the given clip path.
// Experimental API.
func (p *PdfiumImplementation) FPDFClipPath_CountPaths(request *requests.FPDFClipPath_CountPaths) (*responses.FPDFClipPath_CountPaths, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFClipPath_CountPathSegments returns the number of segments inside one path of the given clip path.
// Experimental API.
func (p *PdfiumImplementation) FPDFClipPath_CountPathSegments(request *requests.FPDFClipPath_CountPathSegments) (*responses.FPDFClipPath_CountPathSegments, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFClipPath_GetPathSegment returns the segment in one specific path of the given clip path at index.
// Experimental API.
func (p *PdfiumImplementation) FPDFClipPath_GetPathSegment(request *requests.FPDFClipPath_GetPathSegment) (*responses.FPDFClipPath_GetPathSegment, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}
