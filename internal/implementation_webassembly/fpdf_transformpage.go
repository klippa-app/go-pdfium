package implementation_webassembly

import (
	"errors"
	"unsafe"

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

	_, err = p.Module.ExportedFunction("FPDFPage_SetMediaBox").Call(p.Context, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Left)), *(*uint64)(unsafe.Pointer(&request.Bottom)), *(*uint64)(unsafe.Pointer(&request.Right)), *(*uint64)(unsafe.Pointer(&request.Top)))
	if err != nil {
		return nil, err
	}

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

	_, err = p.Module.ExportedFunction("FPDFPage_SetCropBox").Call(p.Context, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Left)), *(*uint64)(unsafe.Pointer(&request.Bottom)), *(*uint64)(unsafe.Pointer(&request.Right)), *(*uint64)(unsafe.Pointer(&request.Top)))
	if err != nil {
		return nil, err
	}

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

	_, err = p.Module.ExportedFunction("FPDFPage_SetBleedBox").Call(p.Context, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Left)), *(*uint64)(unsafe.Pointer(&request.Bottom)), *(*uint64)(unsafe.Pointer(&request.Right)), *(*uint64)(unsafe.Pointer(&request.Top)))
	if err != nil {
		return nil, err
	}

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

	_, err = p.Module.ExportedFunction("FPDFPage_SetTrimBox").Call(p.Context, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Left)), *(*uint64)(unsafe.Pointer(&request.Bottom)), *(*uint64)(unsafe.Pointer(&request.Right)), *(*uint64)(unsafe.Pointer(&request.Top)))
	if err != nil {
		return nil, err
	}

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

	_, err = p.Module.ExportedFunction("FPDFPage_SetArtBox").Call(p.Context, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Left)), *(*uint64)(unsafe.Pointer(&request.Bottom)), *(*uint64)(unsafe.Pointer(&request.Right)), *(*uint64)(unsafe.Pointer(&request.Top)))
	if err != nil {
		return nil, err
	}

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

	leftPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer leftPointer.Free()

	bottomPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer bottomPointer.Free()

	rightPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer rightPointer.Free()

	topPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer topPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPage_GetMediaBox").Call(p.Context, *pageHandle.handle, leftPointer.Pointer, bottomPointer.Pointer, rightPointer.Pointer, topPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get media box")
	}

	left, err := leftPointer.Value()
	if err != nil {
		return nil, err
	}

	bottom, err := bottomPointer.Value()
	if err != nil {
		return nil, err
	}

	right, err := rightPointer.Value()
	if err != nil {
		return nil, err
	}

	top, err := topPointer.Value()
	if err != nil {
		return nil, err
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

	leftPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer leftPointer.Free()

	bottomPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer bottomPointer.Free()

	rightPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer rightPointer.Free()

	topPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer topPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPage_GetCropBox").Call(p.Context, *pageHandle.handle, leftPointer.Pointer, bottomPointer.Pointer, rightPointer.Pointer, topPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get crop box")
	}

	left, err := leftPointer.Value()
	if err != nil {
		return nil, err
	}

	bottom, err := bottomPointer.Value()
	if err != nil {
		return nil, err
	}

	right, err := rightPointer.Value()
	if err != nil {
		return nil, err
	}

	top, err := topPointer.Value()
	if err != nil {
		return nil, err
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

	leftPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer leftPointer.Free()

	bottomPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer bottomPointer.Free()

	rightPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer rightPointer.Free()

	topPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer topPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPage_GetBleedBox").Call(p.Context, *pageHandle.handle, leftPointer.Pointer, bottomPointer.Pointer, rightPointer.Pointer, topPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get bleed box")
	}

	left, err := leftPointer.Value()
	if err != nil {
		return nil, err
	}

	bottom, err := bottomPointer.Value()
	if err != nil {
		return nil, err
	}

	right, err := rightPointer.Value()
	if err != nil {
		return nil, err
	}

	top, err := topPointer.Value()
	if err != nil {
		return nil, err
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

	leftPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer leftPointer.Free()

	bottomPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer bottomPointer.Free()

	rightPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer rightPointer.Free()

	topPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer topPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPage_GetTrimBox").Call(p.Context, *pageHandle.handle, leftPointer.Pointer, bottomPointer.Pointer, rightPointer.Pointer, topPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get trim box")
	}

	left, err := leftPointer.Value()
	if err != nil {
		return nil, err
	}

	bottom, err := bottomPointer.Value()
	if err != nil {
		return nil, err
	}

	right, err := rightPointer.Value()
	if err != nil {
		return nil, err
	}

	top, err := topPointer.Value()
	if err != nil {
		return nil, err
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

	leftPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer leftPointer.Free()

	bottomPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer bottomPointer.Free()

	rightPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer rightPointer.Free()

	topPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer topPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPage_GetArtBox").Call(p.Context, *pageHandle.handle, leftPointer.Pointer, bottomPointer.Pointer, rightPointer.Pointer, topPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get art box")
	}

	left, err := leftPointer.Value()
	if err != nil {
		return nil, err
	}

	bottom, err := bottomPointer.Value()
	if err != nil {
		return nil, err
	}

	right, err := rightPointer.Value()
	if err != nil {
		return nil, err
	}

	top, err := topPointer.Value()
	if err != nil {
		return nil, err
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

	matrix := uint64(0)
	if request.Matrix != nil {
		matrixPointer, _, err := p.CStructFS_MATRIX(request.Matrix)
		if err != nil {
			return nil, err
		}
		defer p.Free(matrixPointer)
		matrix = matrixPointer
	}

	clipRect := uint64(0)
	if request.ClipRect != nil {
		clipRectPointer, _, err := p.CStructFS_RECTF(request.ClipRect)
		if err != nil {
			return nil, err
		}
		defer p.Free(clipRectPointer)
		clipRect = clipRectPointer
	}

	res, err := p.Module.ExportedFunction("FPDFPage_TransFormWithClip").Call(p.Context, *pageHandle.handle, matrix, clipRect)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	_, err = p.Module.ExportedFunction("FPDFPageObj_TransformClipPath").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.A)), *(*uint64)(unsafe.Pointer(&request.B)), *(*uint64)(unsafe.Pointer(&request.C)), *(*uint64)(unsafe.Pointer(&request.D)), *(*uint64)(unsafe.Pointer(&request.E)), *(*uint64)(unsafe.Pointer(&request.F)))
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObj_TransformClipPath{}, nil
}

// FPDF_CreateClipPath creates a new clip path, with a rectangle inserted.
func (p *PdfiumImplementation) FPDF_CreateClipPath(request *requests.FPDF_CreateClipPath) (*responses.FPDF_CreateClipPath, error) {
	p.Lock()
	defer p.Unlock()

	res, err := p.Module.ExportedFunction("FPDF_CreateClipPath").Call(p.Context, *(*uint64)(unsafe.Pointer(&request.Left)), *(*uint64)(unsafe.Pointer(&request.Bottom)), *(*uint64)(unsafe.Pointer(&request.Right)), *(*uint64)(unsafe.Pointer(&request.Top)))
	if err != nil {
		return nil, err
	}

	clipPath := res[0]
	if clipPath == 0 {
		return nil, errors.New("could not create clip path")
	}

	clipPathHandle := p.registerClipPath(&clipPath)
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

	_, err = p.Module.ExportedFunction("FPDF_DestroyClipPath").Call(p.Context, *clipPathHandle.handle)
	if err != nil {
		return nil, err
	}

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

	_, err = p.Module.ExportedFunction("FPDFPage_InsertClipPath").Call(p.Context, *pageHandle.handle, *clipPathHandle.handle)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPage_InsertClipPath{}, nil
}

// FPDFPageObj_GetClipPath Get the clip path of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetClipPath(request *requests.FPDFPageObj_GetClipPath) (*responses.FPDFPageObj_GetClipPath, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_GetClipPath").Call(p.Context, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	clipPath := res[0]
	clipPathHandle := p.registerClipPath(&clipPath)
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

	res, err := p.Module.ExportedFunction("FPDFClipPath_CountPaths").Call(p.Context, *clipPathHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDFClipPath_CountPathSegments").Call(p.Context, *clipPathHandle.handle, *(*uint64)(unsafe.Pointer(&request.PathIndex)))
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDFClipPath_GetPathSegment").Call(p.Context, *clipPathHandle.handle, *(*uint64)(unsafe.Pointer(&request.PathIndex)), *(*uint64)(unsafe.Pointer(&request.SegmentIndex)))
	if err != nil {
		return nil, err
	}

	pathSegment := res[0]
	if pathSegment == 0 {
		return nil, errors.New("could not get clip path segment")
	}

	pathSegmentHandle := p.registerPathSegment(&pathSegment)

	return &responses.FPDFClipPath_GetPathSegment{
		PathSegment: pathSegmentHandle.nativeRef,
	}, nil
}
