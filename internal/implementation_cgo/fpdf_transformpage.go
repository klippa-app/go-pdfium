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

// FPDFPage_SetMediaBox sets the "MediaBox" entry to the page dictionary.
func (p *PdfiumImplementation) FPDFPage_SetMediaBox(request *requests.FPDFPage_SetMediaBox) (*responses.FPDFPage_SetMediaBox, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	C.FPDFPage_SetMediaBox(pageHandle.handle, C.float(request.Left), C.float(request.Bottom), C.float(request.Right), C.float(request.Top))

	return &responses.FPDFPage_SetMediaBox{}, nil
}

// FPDFPage_SetCropBox sets the "CropBox" entry to the page dictionary.
func (p *PdfiumImplementation) FPDFPage_SetCropBox(request *requests.FPDFPage_SetCropBox) (*responses.FPDFPage_SetCropBox, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	C.FPDFPage_SetCropBox(pageHandle.handle, C.float(request.Left), C.float(request.Bottom), C.float(request.Right), C.float(request.Top))

	return &responses.FPDFPage_SetCropBox{}, nil
}

// FPDFPage_SetBleedBox sets the "BleedBox" entry to the page dictionary.
func (p *PdfiumImplementation) FPDFPage_SetBleedBox(request *requests.FPDFPage_SetBleedBox) (*responses.FPDFPage_SetBleedBox, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	C.FPDFPage_SetBleedBox(pageHandle.handle, C.float(request.Left), C.float(request.Bottom), C.float(request.Right), C.float(request.Top))

	return &responses.FPDFPage_SetBleedBox{}, nil
}

// FPDFPage_SetTrimBox sets the "TrimBox" entry to the page dictionary.
func (p *PdfiumImplementation) FPDFPage_SetTrimBox(request *requests.FPDFPage_SetTrimBox) (*responses.FPDFPage_SetTrimBox, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	C.FPDFPage_SetTrimBox(pageHandle.handle, C.float(request.Left), C.float(request.Bottom), C.float(request.Right), C.float(request.Top))

	return &responses.FPDFPage_SetTrimBox{}, nil
}

// FPDFPage_SetArtBox sets the "ArtBox" entry to the page dictionary.
func (p *PdfiumImplementation) FPDFPage_SetArtBox(request *requests.FPDFPage_SetArtBox) (*responses.FPDFPage_SetArtBox, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	C.FPDFPage_SetArtBox(pageHandle.handle, C.float(request.Left), C.float(request.Bottom), C.float(request.Right), C.float(request.Top))

	return &responses.FPDFPage_SetArtBox{}, nil
}

// FPDFPage_GetMediaBox gets the "MediaBox" entry from the page dictionary
func (p *PdfiumImplementation) FPDFPage_GetMediaBox(request *requests.FPDFPage_GetMediaBox) (*responses.FPDFPage_GetMediaBox, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	left := C.float(0)
	bottom := C.float(0)
	right := C.float(0)
	top := C.float(0)

	success := C.FPDFPage_GetMediaBox(pageHandle.handle, &left, &bottom, &right, &top)
	if int(success) == 0 {
		return nil, errors.New("could not get media box")
	}

	return &responses.FPDFPage_GetMediaBox{
		Left:   float32(left),
		Bottom: float32(bottom),
		Right:  float32(right),
		Top:    float32(top),
	}, nil
}

// FPDFPage_GetCropBox gets the "CropBox" entry from the page dictionary.
func (p *PdfiumImplementation) FPDFPage_GetCropBox(request *requests.FPDFPage_GetCropBox) (*responses.FPDFPage_GetCropBox, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	left := C.float(0)
	bottom := C.float(0)
	right := C.float(0)
	top := C.float(0)

	success := C.FPDFPage_GetCropBox(pageHandle.handle, &left, &bottom, &right, &top)
	if int(success) == 0 {
		return nil, errors.New("could not get crop box")
	}

	return &responses.FPDFPage_GetCropBox{
		Left:   float32(left),
		Bottom: float32(bottom),
		Right:  float32(right),
		Top:    float32(top),
	}, nil
}

// FPDFPage_GetBleedBox gets the "BleedBox" entry from the page dictionary.
func (p *PdfiumImplementation) FPDFPage_GetBleedBox(request *requests.FPDFPage_GetBleedBox) (*responses.FPDFPage_GetBleedBox, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	left := C.float(0)
	bottom := C.float(0)
	right := C.float(0)
	top := C.float(0)

	success := C.FPDFPage_GetBleedBox(pageHandle.handle, &left, &bottom, &right, &top)
	if int(success) == 0 {
		return nil, errors.New("could not get bleed box")
	}

	return &responses.FPDFPage_GetBleedBox{
		Left:   float32(left),
		Bottom: float32(bottom),
		Right:  float32(right),
		Top:    float32(top),
	}, nil
}

// FPDFPage_GetTrimBox gets the "TrimBox" entry from the page dictionary.
func (p *PdfiumImplementation) FPDFPage_GetTrimBox(request *requests.FPDFPage_GetTrimBox) (*responses.FPDFPage_GetTrimBox, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	left := C.float(0)
	bottom := C.float(0)
	right := C.float(0)
	top := C.float(0)

	success := C.FPDFPage_GetTrimBox(pageHandle.handle, &left, &bottom, &right, &top)
	if int(success) == 0 {
		return nil, errors.New("could not get trim box")
	}

	return &responses.FPDFPage_GetTrimBox{
		Left:   float32(left),
		Bottom: float32(bottom),
		Right:  float32(right),
		Top:    float32(top),
	}, nil
}

// FPDFPage_GetArtBox gets the "ArtBox" entry from the page dictionary.
func (p *PdfiumImplementation) FPDFPage_GetArtBox(request *requests.FPDFPage_GetArtBox) (*responses.FPDFPage_GetArtBox, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	left := C.float(0)
	bottom := C.float(0)
	right := C.float(0)
	top := C.float(0)

	success := C.FPDFPage_GetArtBox(pageHandle.handle, &left, &bottom, &right, &top)
	if int(success) == 0 {
		return nil, errors.New("could not get art box")
	}

	return &responses.FPDFPage_GetArtBox{
		Left:   float32(left),
		Bottom: float32(bottom),
		Right:  float32(right),
		Top:    float32(top),
	}, nil
}

// FPDFPage_TransFormWithClip applies the transforms to the page.
func (p *PdfiumImplementation) FPDFPage_TransFormWithClip(request *requests.FPDFPage_TransFormWithClip) (*responses.FPDFPage_TransFormWithClip, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	var matrix *C.FS_MATRIX
	if request.Matrix != nil {
		matrix = &C.FS_MATRIX{
			a: C.float(request.Matrix.A),
			b: C.float(request.Matrix.B),
			c: C.float(request.Matrix.C),
			d: C.float(request.Matrix.D),
			e: C.float(request.Matrix.E),
			f: C.float(request.Matrix.F),
		}
	}

	var clipRect *C.FS_RECTF
	if request.ClipRect != nil {
		clipRect = &C.FS_RECTF{
			left:   C.float(request.ClipRect.Left),
			top:    C.float(request.ClipRect.Top),
			right:  C.float(request.ClipRect.Right),
			bottom: C.float(request.ClipRect.Bottom),
		}
	}

	success := C.FPDFPage_TransFormWithClip(pageHandle.handle, matrix, clipRect)
	if int(success) == 0 {
		return nil, errors.New("could not apply clip transform")
	}

	return &responses.FPDFPage_TransFormWithClip{}, nil
}

// FPDFPageObj_TransformClipPath transform (scale, rotate, shear, move) the clip path of page object.
func (p *PdfiumImplementation) FPDFPageObj_TransformClipPath(request *requests.FPDFPageObj_TransformClipPath) (*responses.FPDFPageObj_TransformClipPath, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	C.FPDFPageObj_TransformClipPath(pageObjectHandle.handle, C.double(request.A), C.double(request.B), C.double(request.C), C.double(request.D), C.double(request.E), C.double(request.F))

	return &responses.FPDFPageObj_TransformClipPath{}, nil
}

// FPDF_CreateClipPath creates a new clip path, with a rectangle inserted.
func (p *PdfiumImplementation) FPDF_CreateClipPath(request *requests.FPDF_CreateClipPath) (*responses.FPDF_CreateClipPath, error) {
	p.Lock()
	defer p.Unlock()

	clipPath := C.FPDF_CreateClipPath(C.float(request.Left), C.float(request.Bottom), C.float(request.Right), C.float(request.Top))
	if clipPath == nil {
		return nil, errors.New("could not create clip path")
	}

	clipPathHandle := p.registerClipPath(clipPath)
	return &responses.FPDF_CreateClipPath{
		ClipPath: clipPathHandle.nativeRef,
	}, nil
}

// FPDF_DestroyClipPath destroys the clip path.
func (p *PdfiumImplementation) FPDF_DestroyClipPath(request *requests.FPDF_DestroyClipPath) (*responses.FPDF_DestroyClipPath, error) {
	p.Lock()
	defer p.Unlock()

	clipPathHandle, err := p.getClipPathHandle(request.ClipPath)
	if err != nil {
		return nil, err
	}

	C.FPDF_DestroyClipPath(clipPathHandle.handle)
	delete(p.clipPathRefs, clipPathHandle.nativeRef)

	return &responses.FPDF_DestroyClipPath{}, nil
}

// FPDFPage_InsertClipPath Clip the page content, the page content that outside the clipping region become invisible.
func (p *PdfiumImplementation) FPDFPage_InsertClipPath(request *requests.FPDFPage_InsertClipPath) (*responses.FPDFPage_InsertClipPath, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	clipPathHandle, err := p.getClipPathHandle(request.ClipPath)
	if err != nil {
		return nil, err
	}

	C.FPDFPage_InsertClipPath(pageHandle.handle, clipPathHandle.handle)

	return &responses.FPDFPage_InsertClipPath{}, nil
}
