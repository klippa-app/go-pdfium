//go:build pdfium_experimental
// +build pdfium_experimental

package implementation

/*
#cgo pkg-config: pdfium
#include "fpdf_annot.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"
	"unsafe"
)

// FPDFAnnot_IsSupportedSubtype returns whether an annotation subtype is currently supported for creation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_IsSupportedSubtype(request *requests.FPDFAnnot_IsSupportedSubtype) (*responses.FPDFAnnot_IsSupportedSubtype, error) {
	p.Lock()
	defer p.Unlock()

	isSupported := C.FPDFAnnot_IsSupportedSubtype(C.FPDF_ANNOTATION_SUBTYPE(request.Subtype))

	return &responses.FPDFAnnot_IsSupportedSubtype{
		IsSupported: int(isSupported) == 1,
	}, nil
}

// FPDFPage_CreateAnnot creates an annotation in the given page of the given subtype. If the specified
// subtype is illegal or unsupported, then a new annotation will not be created.
// Must call FPDFPage_CloseAnnot() when the annotation returned by this
// function is no longer needed.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_CreateAnnot(request *requests.FPDFPage_CreateAnnot) (*responses.FPDFPage_CreateAnnot, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	annotation := C.FPDFPage_CreateAnnot(pageHandle.handle, C.FPDF_ANNOTATION_SUBTYPE(request.Subtype))
	if annotation == nil {
		return nil, errors.New("could not create annotation")
	}

	annotationHandle := p.registerAnnotation(annotation)

	return &responses.FPDFPage_CreateAnnot{
		Annotation: annotationHandle.nativeRef,
	}, nil
}

// FPDFPage_GetAnnotCount returns the number of annotations in a given page.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_GetAnnotCount(request *requests.FPDFPage_GetAnnotCount) (*responses.FPDFPage_GetAnnotCount, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	count := C.FPDFPage_GetAnnotCount(pageHandle.handle)

	return &responses.FPDFPage_GetAnnotCount{
		Count: int(count),
	}, nil
}

// FPDFPage_GetAnnot returns annotation at the given page and index. Must call FPDFPage_CloseAnnot() when the
// annotation returned by this function is no longer needed.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_GetAnnot(request *requests.FPDFPage_GetAnnot) (*responses.FPDFPage_GetAnnot, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	annotation := C.FPDFPage_GetAnnot(pageHandle.handle, C.int(request.Index))
	if annotation == nil {
		return nil, errors.New("could not create annotation")
	}

	annotationHandle := p.registerAnnotation(annotation)

	return &responses.FPDFPage_GetAnnot{
		Annotation: annotationHandle.nativeRef,
	}, nil
}

// FPDFPage_GetAnnotIndex returns the index of the given annotation in the given page. This is the opposite of
// FPDFPage_GetAnnot().
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_GetAnnotIndex(request *requests.FPDFPage_GetAnnotIndex) (*responses.FPDFPage_GetAnnotIndex, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	index := C.FPDFPage_GetAnnotIndex(pageHandle.handle, annotationHandle.handle)
	if int(index) == -1 {
		return nil, errors.New("could not get annotation index")
	}

	return &responses.FPDFPage_GetAnnotIndex{
		Index: int(index),
	}, nil
}

// FPDFPage_CloseAnnot closes an annotation. Must be called when the annotation returned by
// FPDFPage_CreateAnnot() or FPDFPage_GetAnnot() is no longer needed. This
// function does not remove the annotation from the document.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_CloseAnnot(request *requests.FPDFPage_CloseAnnot) (*responses.FPDFPage_CloseAnnot, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	C.FPDFPage_CloseAnnot(annotationHandle.handle)

	delete(p.annotationRefs, request.Annotation)

	return &responses.FPDFPage_CloseAnnot{}, nil
}

// FPDFPage_RemoveAnnot removes the annotation in the given page at the given index.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_RemoveAnnot(request *requests.FPDFPage_RemoveAnnot) (*responses.FPDFPage_RemoveAnnot, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	success := C.FPDFPage_RemoveAnnot(pageHandle.handle, C.int(request.Index))
	if int(success) == 0 {
		return nil, errors.New("could not remove annotation")
	}

	return &responses.FPDFPage_RemoveAnnot{}, nil
}

// FPDFAnnot_GetSubtype returns the subtype of an annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetSubtype(request *requests.FPDFAnnot_GetSubtype) (*responses.FPDFAnnot_GetSubtype, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	subtype := C.FPDFAnnot_GetSubtype(annotationHandle.handle)

	return &responses.FPDFAnnot_GetSubtype{
		Subtype: enums.FPDF_ANNOTATION_SUBTYPE(subtype),
	}, nil
}

// FPDFAnnot_IsObjectSupportedSubtype checks whether an annotation subtype is currently supported for object extraction,
// update, and removal.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_IsObjectSupportedSubtype(request *requests.FPDFAnnot_IsObjectSupportedSubtype) (*responses.FPDFAnnot_IsObjectSupportedSubtype, error) {
	p.Lock()
	defer p.Unlock()

	isSupported := C.FPDFAnnot_IsObjectSupportedSubtype(C.FPDF_ANNOTATION_SUBTYPE(request.Subtype))

	return &responses.FPDFAnnot_IsObjectSupportedSubtype{
		IsObjectSupportedSubtype: int(isSupported) == 1,
	}, nil
}

// FPDFAnnot_UpdateObject updates the given object in the given annotation. The object must be in the annotation already and must have
// been retrieved by FPDFAnnot_GetObject(). Currently, only ink and stamp
// annotations are supported by this API. Also note that only path, image, and
///text objects have APIs for modification; see FPDFPath_*(), FPDFText_*(), and
// FPDFImageObj_*().
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_UpdateObject(request *requests.FPDFAnnot_UpdateObject) (*responses.FPDFAnnot_UpdateObject, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	success := C.FPDFAnnot_UpdateObject(annotationHandle.handle, pageObjectHandle.handle)
	if int(success) == 0 {
		return nil, errors.New("could not update object")
	}

	return &responses.FPDFAnnot_UpdateObject{}, nil
}

// FPDFAnnot_AddInkStroke adds a new InkStroke, represented by an array of points, to the InkList of
// the annotation. The API creates an InkList if one doesn't already exist in the annotation.
// This API works only for ink annotations. Please refer to ISO 32000-1:2008
// spec, section 12.5.6.13.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_AddInkStroke(request *requests.FPDFAnnot_AddInkStroke) (*responses.FPDFAnnot_AddInkStroke, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	if len(request.Points) == 0 {
		return nil, errors.New("at least one point is required")
	}

	pointArray := make([]C.FS_POINTF, len(request.Points))
	for i := range request.Points {
		pointArray[i] = C.FS_POINTF{
			x: C.float(request.Points[i].X),
			y: C.float(request.Points[i].Y),
		}
	}

	index := C.FPDFAnnot_AddInkStroke(annotationHandle.handle, &pointArray[0], C.size_t(len(pointArray)))
	if int(index) == -1 {
		return nil, errors.New("could not add ink stroke")
	}

	return &responses.FPDFAnnot_AddInkStroke{
		Index: int(index),
	}, nil
}

// FPDFAnnot_RemoveInkList removes an InkList in the given annotation.
// This API works only for ink annotations.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_RemoveInkList(request *requests.FPDFAnnot_RemoveInkList) (*responses.FPDFAnnot_RemoveInkList, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	success := C.FPDFAnnot_RemoveInkList(annotationHandle.handle)
	if int(success) == 0 {
		return nil, errors.New("could not remove ink list")
	}

	return &responses.FPDFAnnot_RemoveInkList{}, nil
}

// FPDFAnnot_AppendObject adds the given object to the given annotation. The object must have been created by
// FPDFPageObj_CreateNew{Path|Rect}() or FPDFPageObj_New{Text|Image}Obj(), and
// will be owned by the annotation. Note that an object cannot belong to more than one
// annotation. Currently, only ink and stamp annotations are supported by this API.
// Also note that only path, image, and text objects have APIs for creation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_AppendObject(request *requests.FPDFAnnot_AppendObject) (*responses.FPDFAnnot_AppendObject, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	success := C.FPDFAnnot_AppendObject(annotationHandle.handle, pageObjectHandle.handle)
	if int(success) == 0 {
		return nil, errors.New("could not append object")
	}

	return &responses.FPDFAnnot_AppendObject{}, nil
}

// FPDFAnnot_GetObjectCount returns the total number of objects in the given annotation, including path objects, text
// objects, external objects, image objects, and shading objects.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetObjectCount(request *requests.FPDFAnnot_GetObjectCount) (*responses.FPDFAnnot_GetObjectCount, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	count := C.FPDFAnnot_GetObjectCount(annotationHandle.handle)

	return &responses.FPDFAnnot_GetObjectCount{
		Count: int(count),
	}, nil
}

// FPDFAnnot_GetObject returns the object in the given annotation at the given index.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetObject(request *requests.FPDFAnnot_GetObject) (*responses.FPDFAnnot_GetObject, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	object := C.FPDFAnnot_GetObject(annotationHandle.handle, C.int(request.Index))
	if object == nil {
		return nil, errors.New("could not get object")
	}

	objectHandle := p.registerPageObject(object)

	return &responses.FPDFAnnot_GetObject{
		PageObject: objectHandle.nativeRef,
	}, nil
}

// FPDFAnnot_RemoveObject removes the object in the given annotation at the given index.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_RemoveObject(request *requests.FPDFAnnot_RemoveObject) (*responses.FPDFAnnot_RemoveObject, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	success := C.FPDFAnnot_RemoveObject(annotationHandle.handle, C.int(request.Index))
	if int(success) == 0 {
		return nil, errors.New("could not remove object")
	}

	return &responses.FPDFAnnot_RemoveObject{}, nil
}

// FPDFAnnot_SetColor sets the color of an annotation. Fails when called on annotations with
// appearance streams already defined; instead use
// FPDFPath_Set{Stroke|Fill}Color().
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetColor(request *requests.FPDFAnnot_SetColor) (*responses.FPDFAnnot_SetColor, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	success := C.FPDFAnnot_SetColor(annotationHandle.handle, C.FPDFANNOT_COLORTYPE(request.ColorType), C.uint(request.R), C.uint(request.G), C.uint(request.B), C.uint(request.A))
	if int(success) == 0 {
		return nil, errors.New("could not set annotation color")
	}

	return &responses.FPDFAnnot_SetColor{}, nil
}

// FPDFAnnot_GetColor returns the color of an annotation. If no color is specified, default to yellow
// for highlight annotation, black for all else. Fails when called on
// annotations with appearance streams already defined; instead use
// FPDFPath_Get{Stroke|Fill}Color().
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetColor(request *requests.FPDFAnnot_GetColor) (*responses.FPDFAnnot_GetColor, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	r := C.uint(0)
	g := C.uint(0)
	b := C.uint(0)
	a := C.uint(0)

	success := C.FPDFAnnot_GetColor(annotationHandle.handle, C.FPDFANNOT_COLORTYPE(request.ColorType), &r, &g, &b, &a)
	if int(success) == 0 {
		return nil, errors.New("could not get annotation color")
	}

	return &responses.FPDFAnnot_GetColor{
		R: uint(r),
		G: uint(g),
		B: uint(b),
		A: uint(a),
	}, nil
}

// FPDFAnnot_HasAttachmentPoints returns whether the annotation is of a type that has attachment points
// (i.e. quadpoints). Quadpoints are the vertices of the rectangle that
// encompasses the texts affected by the annotation. They provide the
// coordinates in the page where the annotation is attached. Only text markup
// annotations (i.e. highlight, strikeout, squiggly, and underline) and link
// annotations have quadpoints.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_HasAttachmentPoints(request *requests.FPDFAnnot_HasAttachmentPoints) (*responses.FPDFAnnot_HasAttachmentPoints, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	hasAttachmentPoints := C.FPDFAnnot_HasAttachmentPoints(annotationHandle.handle)

	return &responses.FPDFAnnot_HasAttachmentPoints{
		HasAttachmentPoints: int(hasAttachmentPoints) == 1,
	}, nil
}

// FPDFAnnot_SetAttachmentPoints replaces the attachment points (i.e. quadpoints) set of an annotation at
// the given quad index. This index needs to be within the result of
// FPDFAnnot_CountAttachmentPoints().
// If the annotation's appearance stream is defined and this annotation is of a
// type with quadpoints, then update the bounding box too if the new quadpoints
// define a bigger one.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetAttachmentPoints(request *requests.FPDFAnnot_SetAttachmentPoints) (*responses.FPDFAnnot_SetAttachmentPoints, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	attachmentPoints := C.FS_QUADPOINTSF{
		x1: C.float(request.AttachmentPoints.X1),
		y1: C.float(request.AttachmentPoints.Y1),
		x2: C.float(request.AttachmentPoints.X2),
		y2: C.float(request.AttachmentPoints.Y2),
		x3: C.float(request.AttachmentPoints.X3),
		y3: C.float(request.AttachmentPoints.Y3),
		x4: C.float(request.AttachmentPoints.X4),
		y4: C.float(request.AttachmentPoints.Y4),
	}

	success := C.FPDFAnnot_SetAttachmentPoints(annotationHandle.handle, C.size_t(request.Index), &attachmentPoints)
	if int(success) == 0 {
		return nil, errors.New("could not set attachment points")
	}

	return &responses.FPDFAnnot_SetAttachmentPoints{}, nil
}

// FPDFAnnot_AppendAttachmentPoints appends to the list of attachment points (i.e. quadpoints) of an annotation.
// If the annotation's appearance stream is defined and this annotation is of a
// type with quadpoints, then update the bounding box too if the new quadpoints
// define a bigger one.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_AppendAttachmentPoints(request *requests.FPDFAnnot_AppendAttachmentPoints) (*responses.FPDFAnnot_AppendAttachmentPoints, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	attachmentPoints := C.FS_QUADPOINTSF{
		x1: C.float(request.AttachmentPoints.X1),
		y1: C.float(request.AttachmentPoints.Y1),
		x2: C.float(request.AttachmentPoints.X2),
		y2: C.float(request.AttachmentPoints.Y2),
		x3: C.float(request.AttachmentPoints.X3),
		y3: C.float(request.AttachmentPoints.Y3),
		x4: C.float(request.AttachmentPoints.X4),
		y4: C.float(request.AttachmentPoints.Y4),
	}

	success := C.FPDFAnnot_AppendAttachmentPoints(annotationHandle.handle, &attachmentPoints)
	if int(success) == 0 {
		return nil, errors.New("could not append attachment points")
	}

	return &responses.FPDFAnnot_AppendAttachmentPoints{}, nil
}

// FPDFAnnot_CountAttachmentPoints returns the number of sets of quadpoints of an annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_CountAttachmentPoints(request *requests.FPDFAnnot_CountAttachmentPoints) (*responses.FPDFAnnot_CountAttachmentPoints, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	count := C.FPDFAnnot_CountAttachmentPoints(annotationHandle.handle)

	return &responses.FPDFAnnot_CountAttachmentPoints{
		Count: uint64(count),
	}, nil
}

// FPDFAnnot_GetAttachmentPoints returns the attachment points (i.e. quadpoints) of an annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetAttachmentPoints(request *requests.FPDFAnnot_GetAttachmentPoints) (*responses.FPDFAnnot_GetAttachmentPoints, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	attachmentPoints := C.FS_QUADPOINTSF{}

	success := C.FPDFAnnot_GetAttachmentPoints(annotationHandle.handle, C.size_t(request.Index), &attachmentPoints)
	if int(success) == 0 {
		return nil, errors.New("could not append attachment points")
	}

	return &responses.FPDFAnnot_GetAttachmentPoints{
		QuadPoints: structs.FPDF_FS_QUADPOINTSF{
			X1: float32(attachmentPoints.x1),
			Y1: float32(attachmentPoints.y1),
			X2: float32(attachmentPoints.x2),
			Y2: float32(attachmentPoints.y2),
			X3: float32(attachmentPoints.x3),
			Y3: float32(attachmentPoints.y3),
			X4: float32(attachmentPoints.x4),
			Y4: float32(attachmentPoints.y4),
		},
	}, nil
}

// FPDFAnnot_SetRect sets the annotation rectangle defining the location of the annotation. If the
// annotation's appearance stream is defined and this annotation is of a type
// without quadpoints, then update the bounding box too if the new rectangle
// defines a bigger one.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetRect(request *requests.FPDFAnnot_SetRect) (*responses.FPDFAnnot_SetRect, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	rect := C.FS_RECTF{
		left:   C.float(request.Rect.Left),
		top:    C.float(request.Rect.Top),
		right:  C.float(request.Rect.Right),
		bottom: C.float(request.Rect.Bottom),
	}

	success := C.FPDFAnnot_SetRect(annotationHandle.handle, &rect)
	if int(success) == 0 {
		return nil, errors.New("could net set rect")
	}

	return &responses.FPDFAnnot_SetRect{}, nil
}

// FPDFAnnot_GetRect returns the annotation rectangle defining the location of the annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetRect(request *requests.FPDFAnnot_GetRect) (*responses.FPDFAnnot_GetRect, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	rect := C.FS_RECTF{}

	success := C.FPDFAnnot_GetRect(annotationHandle.handle, &rect)
	if int(success) == 0 {
		return nil, errors.New("could net set rect")
	}

	return &responses.FPDFAnnot_GetRect{
		Rect: structs.FPDF_FS_RECTF{
			Left:   float32(rect.left),
			Top:    float32(rect.top),
			Right:  float32(rect.right),
			Bottom: float32(rect.bottom),
		},
	}, nil
}

// FPDFAnnot_GetVertices returns the vertices of a polygon or polyline annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetVertices(request *requests.FPDFAnnot_GetVertices) (*responses.FPDFAnnot_GetVertices, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	// First get the vertices array size.
	length := C.FPDFAnnot_GetVertices(annotationHandle.handle, nil, 0)
	goVertices := make([]structs.FPDF_FS_POINTF, uint64(length))
	cVertices := make([]C.FS_POINTF, uint64(length))
	if length > 0 {
		// Actually fill the array
		C.FPDFAnnot_GetVertices(annotationHandle.handle, &cVertices[0], length)
	}

	for i := range cVertices {
		goVertices[i] = structs.FPDF_FS_POINTF{
			X: float32(cVertices[i].x),
			Y: float32(cVertices[i].y),
		}
	}

	return &responses.FPDFAnnot_GetVertices{
		Vertices: goVertices,
	}, nil
}

// FPDFAnnot_GetInkListCount returns the number of paths in the ink list of an ink annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetInkListCount(request *requests.FPDFAnnot_GetInkListCount) (*responses.FPDFAnnot_GetInkListCount, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	count := C.FPDFAnnot_GetInkListCount(annotationHandle.handle)

	return &responses.FPDFAnnot_GetInkListCount{
		Count: uint64(count),
	}, nil
}

// FPDFAnnot_GetInkListPath returns a path in the ink list of an ink annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetInkListPath(request *requests.FPDFAnnot_GetInkListPath) (*responses.FPDFAnnot_GetInkListPath, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	// First get the path array size.
	length := C.FPDFAnnot_GetInkListPath(annotationHandle.handle, C.ulong(request.Index), nil, 0)
	goPath := make([]structs.FPDF_FS_POINTF, uint64(length))
	cPath := make([]C.FS_POINTF, uint64(length))
	if length > 0 {
		// Actually fill the array
		C.FPDFAnnot_GetInkListPath(annotationHandle.handle, C.ulong(request.Index), &cPath[0], length)
	}

	for i := range cPath {
		goPath[i] = structs.FPDF_FS_POINTF{
			X: float32(cPath[i].x),
			Y: float32(cPath[i].y),
		}
	}
	return &responses.FPDFAnnot_GetInkListPath{
		Path: goPath,
	}, nil
}

// FPDFAnnot_GetLine returns the starting and ending coordinates of a line annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetLine(request *requests.FPDFAnnot_GetLine) (*responses.FPDFAnnot_GetLine, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	start := C.FS_POINTF{}
	end := C.FS_POINTF{}

	success := C.FPDFAnnot_GetLine(annotationHandle.handle, &start, &end)
	if int(success) == 0 {
		return nil, errors.New("could not get line")
	}

	return &responses.FPDFAnnot_GetLine{
		Start: structs.FPDF_FS_POINTF{
			X: float32(start.x),
			Y: float32(start.y),
		},
		End: structs.FPDF_FS_POINTF{
			X: float32(end.x),
			Y: float32(end.y),
		},
	}, nil
}

// FPDFAnnot_SetBorder sets the characteristics of the annotation's border (rounded rectangle).
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetBorder(request *requests.FPDFAnnot_SetBorder) (*responses.FPDFAnnot_SetBorder, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	success := C.FPDFAnnot_SetBorder(annotationHandle.handle, C.float(request.HorizontalRadius), C.float(request.VerticalRadius), C.float(request.BorderWidth))
	if int(success) == 0 {
		return nil, errors.New("could not set border")
	}

	return &responses.FPDFAnnot_SetBorder{}, nil
}

// FPDFAnnot_GetBorder returns the characteristics of the annotation's border (rounded rectangle).
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetBorder(request *requests.FPDFAnnot_GetBorder) (*responses.FPDFAnnot_GetBorder, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	horizontalRadius := C.float(0)
	verticalRadius := C.float(0)
	borderWidth := C.float(0)

	success := C.FPDFAnnot_GetBorder(annotationHandle.handle, &horizontalRadius, &verticalRadius, &borderWidth)
	if int(success) == 0 {
		return nil, errors.New("could not set border")
	}

	return &responses.FPDFAnnot_GetBorder{
		HorizontalRadius: float32(horizontalRadius),
		VerticalRadius:   float32(verticalRadius),
		BorderWidth:      float32(borderWidth),
	}, nil
}

// FPDFAnnot_HasKey checks whether the given annotation's dictionary has the given key as a key.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_HasKey(request *requests.FPDFAnnot_HasKey) (*responses.FPDFAnnot_HasKey, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	keyStr := C.CString(request.Key)
	defer C.free(unsafe.Pointer(keyStr))

	hasKey := C.FPDFAnnot_HasKey(annotationHandle.handle, keyStr)

	return &responses.FPDFAnnot_HasKey{
		HasKey: int(hasKey) == 1,
	}, nil
}

// FPDFAnnot_GetValueType returns the type of the value corresponding to the given key the annotation's dictionary.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetValueType(request *requests.FPDFAnnot_GetValueType) (*responses.FPDFAnnot_GetValueType, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	keyStr := C.CString(request.Key)
	defer C.free(unsafe.Pointer(keyStr))

	valueType := C.FPDFAnnot_GetValueType(annotationHandle.handle, keyStr)

	return &responses.FPDFAnnot_GetValueType{
		ValueType: enums.FPDF_OBJECT_TYPE(valueType),
	}, nil
}

// FPDFAnnot_SetStringValue sets the string value corresponding to the given key in the annotations's dictionary,
// overwriting the existing value if any. The value type would be
// FPDF_OBJECT_STRING after this function call succeeds.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetStringValue(request *requests.FPDFAnnot_SetStringValue) (*responses.FPDFAnnot_SetStringValue, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	keyStr := C.CString(request.Key)
	defer C.free(unsafe.Pointer(keyStr))

	transformedText, err := p.transformUTF8ToUTF16LE(request.Value)
	if err != nil {
		return nil, err
	}

	success := C.FPDFAnnot_SetStringValue(annotationHandle.handle, keyStr, (C.FPDF_WIDESTRING)(unsafe.Pointer(&transformedText[0])))
	if int(success) == 0 {
		return nil, errors.New("could not set string value")
	}

	return &responses.FPDFAnnot_SetStringValue{}, nil
}

// FPDFAnnot_GetStringValue returns the string value corresponding to the given key in the annotations's dictionary.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetStringValue(request *requests.FPDFAnnot_GetStringValue) (*responses.FPDFAnnot_GetStringValue, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	keyStr := C.CString(request.Key)
	defer C.free(unsafe.Pointer(keyStr))

	// First get the value length
	length := C.FPDFAnnot_GetStringValue(annotationHandle.handle, keyStr, nil, 0)
	if uint64(length) == 0 {
		return nil, errors.New("could not get string value")
	}

	charData := make([]byte, length)
	C.FPDFAnnot_GetStringValue(annotationHandle.handle, keyStr, (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetStringValue{
		Value: transformedText,
	}, nil
}

// FPDFAnnot_GetNumberValue returns the float value corresponding to the given key in the annotations's dictionary.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetNumberValue(request *requests.FPDFAnnot_GetNumberValue) (*responses.FPDFAnnot_GetNumberValue, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	keyStr := C.CString(request.Key)
	defer C.free(unsafe.Pointer(keyStr))

	value := C.float(0)
	success := C.FPDFAnnot_GetNumberValue(annotationHandle.handle, keyStr, &value)
	if int(success) == 0 {
		return nil, errors.New("could not get number value")
	}

	return &responses.FPDFAnnot_GetNumberValue{
		Value: float32(value),
	}, nil
}

// FPDFAnnot_SetAP sets the AP (appearance string) in annotations's dictionary for a given appearance mode.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetAP(request *requests.FPDFAnnot_SetAP) (*responses.FPDFAnnot_SetAP, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	if request.Value != nil {
		transformedText, err := p.transformUTF8ToUTF16LE(*request.Value)
		if err != nil {
			return nil, err
		}

		success := C.FPDFAnnot_SetAP(annotationHandle.handle, C.FPDF_ANNOT_APPEARANCEMODE(request.AppearanceMode), (C.FPDF_WIDESTRING)(unsafe.Pointer(&transformedText[0])))
		if int(success) == 0 {
			return nil, errors.New("could not set appearance mode")
		}
	} else {
		success := C.FPDFAnnot_SetAP(annotationHandle.handle, C.FPDF_ANNOT_APPEARANCEMODE(request.AppearanceMode), nil)
		if int(success) == 0 {
			return nil, errors.New("could not set appearance mode")
		}
	}

	return &responses.FPDFAnnot_SetAP{}, nil
}

// FPDFAnnot_GetAP returns the AP (appearance string) from annotation's dictionary for a given
// appearance mode.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetAP(request *requests.FPDFAnnot_GetAP) (*responses.FPDFAnnot_GetAP, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	// First get the value length
	length := C.FPDFAnnot_GetAP(annotationHandle.handle, C.FPDF_ANNOT_APPEARANCEMODE(request.AppearanceMode), nil, 0)
	if uint64(length) == 0 {
		return nil, errors.New("could not get appearance mode")
	}

	charData := make([]byte, length)
	C.FPDFAnnot_GetAP(annotationHandle.handle, C.FPDF_ANNOT_APPEARANCEMODE(request.AppearanceMode), (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetAP{
		AppearanceMode: transformedText,
	}, nil
}

// FPDFAnnot_GetLinkedAnnot returns the annotation corresponding to the given key in the annotations's dictionary. Common
// keys for linking annotations include "IRT" and "Popup". Must call
// FPDFPage_CloseAnnot() when the annotation returned by this function is no
// longer needed.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetLinkedAnnot(request *requests.FPDFAnnot_GetLinkedAnnot) (*responses.FPDFAnnot_GetLinkedAnnot, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	keyStr := C.CString(request.Key)
	defer C.free(unsafe.Pointer(keyStr))

	linkedAnnotation := C.FPDFAnnot_GetLinkedAnnot(annotationHandle.handle, keyStr)
	if linkedAnnotation == nil {
		return nil, errors.New("could not get linked annotation")
	}

	linkedAnnotationHandle := p.registerAnnotation(linkedAnnotation)

	return &responses.FPDFAnnot_GetLinkedAnnot{
		LinkedAnnotation: linkedAnnotationHandle.nativeRef,
	}, nil
}

// FPDFAnnot_GetFlags returns the annotation flags of the given annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFlags(request *requests.FPDFAnnot_GetFlags) (*responses.FPDFAnnot_GetFlags, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	flags := C.FPDFAnnot_GetFlags(annotationHandle.handle)

	return &responses.FPDFAnnot_GetFlags{
		Flags: int(flags),
	}, nil
}

// FPDFAnnot_SetFlags sets the annotation flags of the given annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetFlags(request *requests.FPDFAnnot_SetFlags) (*responses.FPDFAnnot_SetFlags, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	success := C.FPDFAnnot_SetFlags(annotationHandle.handle, C.int(request.Flags))
	if int(success) == 0 {
		return nil, errors.New("could not set flags")
	}

	return &responses.FPDFAnnot_SetFlags{}, nil
}

// FPDFAnnot_GetFormFieldFlags returns the form field annotation flags of the given annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormFieldFlags(request *requests.FPDFAnnot_GetFormFieldFlags) (*responses.FPDFAnnot_GetFormFieldFlags, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	flags := C.FPDFAnnot_GetFormFieldFlags(formHandle.handle, annotationHandle.handle)

	return &responses.FPDFAnnot_GetFormFieldFlags{
		Flags: int(flags),
	}, nil
}

// FPDFAnnot_GetFormFieldAtPoint returns an interactive form annotation whose rectangle contains a given
// point on a page. Must call FPDFPage_CloseAnnot() when the annotation returned
// is no longer needed.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormFieldAtPoint(request *requests.FPDFAnnot_GetFormFieldAtPoint) (*responses.FPDFAnnot_GetFormFieldAtPoint, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	point := C.FS_POINTF{
		x: C.float(request.Point.X),
		y: C.float(request.Point.Y),
	}

	formField := C.FPDFAnnot_GetFormFieldAtPoint(formHandle.handle, pageHandle.handle, &point)
	if formHandle == nil {
		return nil, errors.New("could not get form field")
	}

	formFieldHandle := p.registerAnnotation(formField)

	return &responses.FPDFAnnot_GetFormFieldAtPoint{
		Annotation: formFieldHandle.nativeRef,
	}, nil
}

// FPDFAnnot_GetFormFieldName returns the name of the given annotation, which is an interactive form annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormFieldName(request *requests.FPDFAnnot_GetFormFieldName) (*responses.FPDFAnnot_GetFormFieldName, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	// First get the value length
	length := C.FPDFAnnot_GetFormFieldName(formHandle.handle, annotationHandle.handle, nil, 0)
	if uint64(length) == 0 {
		return nil, errors.New("could not get form field name")
	}

	charData := make([]byte, length)
	C.FPDFAnnot_GetFormFieldName(formHandle.handle, annotationHandle.handle, (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetFormFieldName{
		FormFieldName: transformedText,
	}, nil
}

// FPDFAnnot_GetFormFieldType returns the form field type of the given annotation, which is an interactive form annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormFieldType(request *requests.FPDFAnnot_GetFormFieldType) (*responses.FPDFAnnot_GetFormFieldType, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	formFieldType := C.FPDFAnnot_GetFormFieldType(formHandle.handle, annotationHandle.handle)
	if int(formFieldType) == 0 {
		return nil, errors.New("could not get form field type")
	}

	return &responses.FPDFAnnot_GetFormFieldType{
		FormFieldType: enums.FPDF_FORMFIELD_TYPE(formFieldType),
	}, nil
}

// FPDFAnnot_GetFormFieldValue returns the value of the given annotation, which is an interactive form annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormFieldValue(request *requests.FPDFAnnot_GetFormFieldValue) (*responses.FPDFAnnot_GetFormFieldValue, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	// First get the value length
	length := C.FPDFAnnot_GetFormFieldValue(formHandle.handle, annotationHandle.handle, nil, 0)
	if uint64(length) == 0 {
		return nil, errors.New("could not get form field value")
	}

	charData := make([]byte, length)
	C.FPDFAnnot_GetFormFieldValue(formHandle.handle, annotationHandle.handle, (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetFormFieldValue{
		FormFieldValue: transformedText,
	}, nil
}

// FPDFAnnot_GetOptionCount returns the number of options in the annotation's "Opt" dictionary. Intended for
// use with listbox and combobox widget annotations.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetOptionCount(request *requests.FPDFAnnot_GetOptionCount) (*responses.FPDFAnnot_GetOptionCount, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	optionCount := C.FPDFAnnot_GetOptionCount(formHandle.handle, annotationHandle.handle)
	if int(optionCount) == -1 {
		return nil, errors.New("could not get option count")
	}

	return &responses.FPDFAnnot_GetOptionCount{
		OptionCount: int(optionCount),
	}, nil
}

// FPDFAnnot_GetOptionLabel returns the string value for the label of the option at the given index in annotation's
// "Opt" dictionary. Intended for use with listbox and combobox widget
// annotations.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetOptionLabel(request *requests.FPDFAnnot_GetOptionLabel) (*responses.FPDFAnnot_GetOptionLabel, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	// First get the value length
	length := C.FPDFAnnot_GetOptionLabel(formHandle.handle, annotationHandle.handle, C.int(request.Index), nil, 0)
	if uint64(length) == 0 {
		return nil, errors.New("could not get form field name")
	}

	charData := make([]byte, length)
	C.FPDFAnnot_GetOptionLabel(formHandle.handle, annotationHandle.handle, C.int(request.Index), (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetOptionLabel{
		OptionLabel: transformedText,
	}, nil
}

// FPDFAnnot_IsOptionSelected returns whether or not the option at the given index in annotation's "Opt" dictionary
// is selected. Intended for use with listbox and combobox widget annotations.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_IsOptionSelected(request *requests.FPDFAnnot_IsOptionSelected) (*responses.FPDFAnnot_IsOptionSelected, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	isOptionSelected := C.FPDFAnnot_IsOptionSelected(formHandle.handle, annotationHandle.handle, C.int(request.Index))

	return &responses.FPDFAnnot_IsOptionSelected{
		IsOptionSelected: int(isOptionSelected) == 1,
	}, nil
}

// FPDFAnnot_GetFontSize returns the float value of the font size for an annotation with variable text.
// If 0, the font is to be auto-sized: its size is computed as a function of
// the height of the annotation rectangle.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFontSize(request *requests.FPDFAnnot_GetFontSize) (*responses.FPDFAnnot_GetFontSize, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	fontSize := C.float(0)
	success := C.FPDFAnnot_GetFontSize(formHandle.handle, annotationHandle.handle, &fontSize)
	if int(success) == 0 {
		return nil, errors.New("could not get font size")
	}

	return &responses.FPDFAnnot_GetFontSize{
		FontSize: float32(fontSize),
	}, nil
}

// FPDFAnnot_IsChecked returns whether the given annotation is a form widget that is checked. Intended for use with
// checkbox and radio button widgets.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_IsChecked(request *requests.FPDFAnnot_IsChecked) (*responses.FPDFAnnot_IsChecked, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	isChecked := C.FPDFAnnot_IsChecked(formHandle.handle, annotationHandle.handle)

	return &responses.FPDFAnnot_IsChecked{
		IsChecked: int(isChecked) == 1,
	}, nil
}

// FPDFAnnot_SetFocusableSubtypes returns the list of focusable annotation subtypes. Annotations of subtype
// FPDF_ANNOT_WIDGET are by default focusable. New subtypes set using this API
// will override the existing subtypes.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetFocusableSubtypes(request *requests.FPDFAnnot_SetFocusableSubtypes) (*responses.FPDFAnnot_SetFocusableSubtypes, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	if len(request.Subtypes) == 0 {
		return nil, errors.New("subtypes are required")
	}

	focusableSubtypes := make([]C.FPDF_ANNOTATION_SUBTYPE, len(request.Subtypes))
	for i := range request.Subtypes {
		focusableSubtypes[i] = C.FPDF_ANNOTATION_SUBTYPE(request.Subtypes[i])
	}

	success := C.FPDFAnnot_SetFocusableSubtypes(formHandle.handle, &focusableSubtypes[0], C.size_t(len(focusableSubtypes)))
	if int(success) == 0 {
		return nil, errors.New("could net set focusable subtypes")
	}

	return &responses.FPDFAnnot_SetFocusableSubtypes{}, nil
}

// FPDFAnnot_GetFocusableSubtypesCount returns the count of focusable annotation subtypes as set by host.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFocusableSubtypesCount(request *requests.FPDFAnnot_GetFocusableSubtypesCount) (*responses.FPDFAnnot_GetFocusableSubtypesCount, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	count := C.FPDFAnnot_GetFocusableSubtypesCount(formHandle.handle)
	if int(count) == -1 {
		return nil, errors.New("could net get focusable subtypes count")
	}

	return &responses.FPDFAnnot_GetFocusableSubtypesCount{}, nil
}

// FPDFAnnot_GetFocusableSubtypes returns the list of focusable annotation subtype as set by host.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFocusableSubtypes(request *requests.FPDFAnnot_GetFocusableSubtypes) (*responses.FPDFAnnot_GetFocusableSubtypes, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	count := C.FPDFAnnot_GetFocusableSubtypesCount(formHandle.handle)
	if int(count) == -1 {
		return nil, errors.New("could net get focusable subtypes")
	}

	goFocusableSubtypes := make([]enums.FPDF_ANNOTATION_SUBTYPE, int(count))

	if int(count) > 0 {
		focusableSubtypes := make([]C.FPDF_ANNOTATION_SUBTYPE, int(count))

		success := C.FPDFAnnot_SetFocusableSubtypes(formHandle.handle, &focusableSubtypes[0], C.size_t(len(focusableSubtypes)))
		if int(success) == 0 {
			return nil, errors.New("could net get focusable subtypes")
		}

		for i := range focusableSubtypes {
			goFocusableSubtypes[i] = enums.FPDF_ANNOTATION_SUBTYPE(focusableSubtypes[i])
		}
	}

	return &responses.FPDFAnnot_GetFocusableSubtypes{
		FocusableSubtypes: goFocusableSubtypes,
	}, nil
}

// FPDFAnnot_GetLink returns FPDF_LINK object for the given annotation. Intended to use for link annotations.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetLink(request *requests.FPDFAnnot_GetLink) (*responses.FPDFAnnot_GetLink, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	link := C.FPDFAnnot_GetLink(annotationHandle.handle)
	if link == nil {
		return nil, errors.New("could not get link")
	}

	linkHandle := p.registerLink(link)

	return &responses.FPDFAnnot_GetLink{
		Link: linkHandle.nativeRef,
	}, nil
}

// FPDFAnnot_GetFormControlCount returns the count of annotations in the annotation's control group.
// A group of interactive form annotations is collectively called a form
// control group. Here, annotation, an interactive form annotation, should be
// either a radio button or a checkbox.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormControlCount(request *requests.FPDFAnnot_GetFormControlCount) (*responses.FPDFAnnot_GetFormControlCount, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	formControlCount := C.FPDFAnnot_GetFormControlCount(formHandle.handle, annotationHandle.handle)
	if int(formControlCount) == -1 {
		return nil, errors.New("could net get form control count")
	}

	return &responses.FPDFAnnot_GetFormControlCount{
		FormControlCount: int(formControlCount),
	}, nil
}

// FPDFAnnot_GetFormControlIndex returns the index of the given annotation it's control group.
// A group of interactive form annotations is collectively called a form
// control group. Here, the annotation, an interactive form annotation, should be
// either a radio button or a checkbox.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormControlIndex(request *requests.FPDFAnnot_GetFormControlIndex) (*responses.FPDFAnnot_GetFormControlIndex, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	formControlIndex := C.FPDFAnnot_GetFormControlCount(formHandle.handle, annotationHandle.handle)
	if int(formControlIndex) == -1 {
		return nil, errors.New("could net get form control index")
	}

	return &responses.FPDFAnnot_GetFormControlIndex{
		FormControlIndex: int(formControlIndex),
	}, nil
}

// FPDFAnnot_GetFormFieldExportValue returns the export value of the given annotation which is an interactive form annotation.
// Intended for use with radio button and checkbox widget annotations.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormFieldExportValue(request *requests.FPDFAnnot_GetFormFieldExportValue) (*responses.FPDFAnnot_GetFormFieldExportValue, error) {
	p.Lock()
	defer p.Unlock()

	formHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	// First get the value length
	length := C.FPDFAnnot_GetFormFieldExportValue(formHandle.handle, annotationHandle.handle, nil, 0)
	if uint64(length) == 0 {
		return nil, errors.New("could not get form field export value")
	}

	charData := make([]byte, length)
	C.FPDFAnnot_GetFormFieldExportValue(formHandle.handle, annotationHandle.handle, (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetFormFieldExportValue{
		Value: transformedText,
	}, nil
}

// FPDFAnnot_SetURI adds a URI action to the given annotation, overwriting the existing action, if any.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetURI(request *requests.FPDFAnnot_SetURI) (*responses.FPDFAnnot_SetURI, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	uriStr := C.CString(request.URI)
	defer C.free(unsafe.Pointer(uriStr))

	success := C.FPDFAnnot_SetURI(annotationHandle.handle, uriStr)
	if int(success) == 0 {
		return nil, errors.New("could net set uri")
	}

	return &responses.FPDFAnnot_SetURI{}, nil
}
