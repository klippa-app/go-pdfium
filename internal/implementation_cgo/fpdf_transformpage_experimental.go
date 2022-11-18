//go:build pdfium_experimental
// +build pdfium_experimental

package implementation_cgo

/*
#cgo pkg-config: pdfium
#include "fpdf_transformpage.h"
*/
import "C"
import (
	"errors"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFPageObj_GetClipPath Get the clip path of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetClipPath(request *requests.FPDFPageObj_GetClipPath) (*responses.FPDFPageObj_GetClipPath, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	clipPath := C.FPDFPageObj_GetClipPath(pageObjectHandle.handle)

	clipPathHandle := p.registerClipPath(clipPath)
	return &responses.FPDFPageObj_GetClipPath{
		ClipPath: clipPathHandle.nativeRef,
	}, nil
}

// FPDFClipPath_CountPaths returns the number of paths inside the given clip path.
// Experimental API.
func (p *PdfiumImplementation) FPDFClipPath_CountPaths(request *requests.FPDFClipPath_CountPaths) (*responses.FPDFClipPath_CountPaths, error) {
	p.Lock()
	defer p.Unlock()

	clipPathHandle, err := p.getClipPathHandle(request.ClipPath)
	if err != nil {
		return nil, err
	}

	count := C.FPDFClipPath_CountPaths(clipPathHandle.handle)
	if int(count) == -1 {
		return nil, errors.New("could not get clip path path count")
	}

	return &responses.FPDFClipPath_CountPaths{
		Count: int(count),
	}, nil
}

// FPDFClipPath_CountPathSegments returns the number of segments inside one path of the given clip path.
// Experimental API.
func (p *PdfiumImplementation) FPDFClipPath_CountPathSegments(request *requests.FPDFClipPath_CountPathSegments) (*responses.FPDFClipPath_CountPathSegments, error) {
	p.Lock()
	defer p.Unlock()

	clipPathHandle, err := p.getClipPathHandle(request.ClipPath)
	if err != nil {
		return nil, err
	}

	count := C.FPDFClipPath_CountPathSegments(clipPathHandle.handle, C.int(request.PathIndex))

	if int(count) == -1 {
		return nil, errors.New("could not get clip path path segment count")
	}

	return &responses.FPDFClipPath_CountPathSegments{
		Count: int(count),
	}, nil
}

// FPDFClipPath_GetPathSegment returns the segment in one specific path of the given clip path at index.
// Experimental API.
func (p *PdfiumImplementation) FPDFClipPath_GetPathSegment(request *requests.FPDFClipPath_GetPathSegment) (*responses.FPDFClipPath_GetPathSegment, error) {
	p.Lock()
	defer p.Unlock()

	clipPathHandle, err := p.getClipPathHandle(request.ClipPath)
	if err != nil {
		return nil, err
	}

	pathSegment := C.FPDFClipPath_GetPathSegment(clipPathHandle.handle, C.int(request.PathIndex), C.int(request.SegmentIndex))
	if pathSegment == nil {
		return nil, errors.New("could not get clip path segment")
	}

	pathSegmentHandle := p.registerPathSegment(pathSegment)

	return &responses.FPDFClipPath_GetPathSegment{
		PathSegment: pathSegmentHandle.nativeRef,
	}, nil
}
