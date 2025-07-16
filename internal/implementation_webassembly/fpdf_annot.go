package implementation_webassembly

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"
)

// FPDFAnnot_IsSupportedSubtype returns whether an annotation subtype is currently supported for creation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_IsSupportedSubtype(request *requests.FPDFAnnot_IsSupportedSubtype) (*responses.FPDFAnnot_IsSupportedSubtype, error) {
	p.Lock()
	defer p.Unlock()

	res, err := p.Module.ExportedFunction("FPDFAnnot_IsSupportedSubtype").Call(p.Context, *(*uint64)(unsafe.Pointer(&request.Subtype)))
	if err != nil {
		return nil, err
	}

	isSupported := *(*int32)(unsafe.Pointer(&res[0]))

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

	res, err := p.Module.ExportedFunction("FPDFPage_CreateAnnot").Call(p.Context, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Subtype)))
	if err != nil {
		return nil, err
	}

	annotation := res[0]
	if annotation == 0 {
		return nil, errors.New("could not create annotation")
	}

	annotationHandle := p.registerAnnotation(&annotation)

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

	res, err := p.Module.ExportedFunction("FPDFPage_GetAnnotCount").Call(p.Context, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

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

	res, err := p.Module.ExportedFunction("FPDFPage_GetAnnot").Call(p.Context, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	annotation := res[0]
	if annotation == 0 {
		return nil, errors.New("could not get annotation")
	}

	annotationHandle := p.registerAnnotation(&annotation)

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

	res, err := p.Module.ExportedFunction("FPDFPage_GetAnnotIndex").Call(p.Context, *pageHandle.handle, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	index := *(*int32)(unsafe.Pointer(&res[0]))
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

	_, err = p.Module.ExportedFunction("FPDFPage_CloseAnnot").Call(p.Context, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

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

	res, err := p.Module.ExportedFunction("FPDFPage_RemoveAnnot").Call(p.Context, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetSubtype").Call(p.Context, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	subtype := *(*int32)(unsafe.Pointer(&res[0]))

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

	res, err := p.Module.ExportedFunction("FPDFAnnot_IsObjectSupportedSubtype").Call(p.Context, *(*uint64)(unsafe.Pointer(&request.Subtype)))
	if err != nil {
		return nil, err
	}

	isSupported := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFAnnot_IsObjectSupportedSubtype{
		IsObjectSupportedSubtype: int(isSupported) == 1,
	}, nil
}

// FPDFAnnot_UpdateObject updates the given object in the given annotation. The object must be in the annotation already and must have
// been retrieved by FPDFAnnot_GetObject(). Currently, only ink and stamp
// annotations are supported by this API. Also note that only path, image, and
// text objects have APIs for modification; see FPDFPath_*(), FPDFText_*(), and
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_UpdateObject").Call(p.Context, *annotationHandle.handle, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	pointSize := uint64(len(request.Points))
	pointArray, err := p.FS_POINTFArrayPointer(pointSize, request.Points)
	if err != nil {
		return nil, err
	}
	defer pointArray.Free()

	res, err := p.Module.ExportedFunction("FPDFAnnot_AddInkStroke").Call(p.Context, *annotationHandle.handle, pointArray.Pointer, pointSize)
	if err != nil {
		return nil, err
	}

	index := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_RemoveInkList").Call(p.Context, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_AppendObject").Call(p.Context, *annotationHandle.handle, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetObjectCount").Call(p.Context, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetObject").Call(p.Context, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	object := res[0]
	if object == 0 {
		return nil, errors.New("could not get object")
	}

	objectHandle := p.registerPageObject(&object)

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

	res, err := p.Module.ExportedFunction("FPDFAnnot_RemoveObject").Call(p.Context, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_SetColor").Call(p.Context, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.ColorType)), *(*uint64)(unsafe.Pointer(&request.R)), *(*uint64)(unsafe.Pointer(&request.G)), *(*uint64)(unsafe.Pointer(&request.B)), *(*uint64)(unsafe.Pointer(&request.A)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	rPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer rPointer.Free()

	gPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer gPointer.Free()

	bPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer bPointer.Free()

	aPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer aPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetColor").Call(p.Context, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.ColorType)), rPointer.Pointer, gPointer.Pointer, bPointer.Pointer, aPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get annotation color")
	}

	r, err := rPointer.Value()
	if err != nil {
		return nil, err
	}

	g, err := gPointer.Value()
	if err != nil {
		return nil, err
	}

	b, err := bPointer.Value()
	if err != nil {
		return nil, err
	}

	a, err := aPointer.Value()
	if err != nil {
		return nil, err
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_HasAttachmentPoints").Call(p.Context, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	hasAttachmentPoints := *(*int32)(unsafe.Pointer(&res[0]))

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

	attachmentPointsPointer, _, err := p.CStructFS_QUADPOINTSF(&request.AttachmentPoints)
	if err != nil {
		return nil, err
	}
	defer p.Free(attachmentPointsPointer)

	res, err := p.Module.ExportedFunction("FPDFAnnot_SetAttachmentPoints").Call(p.Context, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)), attachmentPointsPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	attachmentPointsPointer, _, err := p.CStructFS_QUADPOINTSF(&request.AttachmentPoints)
	if err != nil {
		return nil, err
	}
	defer p.Free(attachmentPointsPointer)

	res, err := p.Module.ExportedFunction("FPDFAnnot_AppendAttachmentPoints").Call(p.Context, *annotationHandle.handle, attachmentPointsPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_CountAttachmentPoints").Call(p.Context, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

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

	attachmentPointsPointer, attachmentPointsValue, err := p.CStructFS_QUADPOINTSF(nil)
	if err != nil {
		return nil, err
	}
	defer p.Free(attachmentPointsPointer)

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetAttachmentPoints").Call(p.Context, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)), attachmentPointsPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not append attachment points")
	}

	attachmentPoints, err := attachmentPointsValue()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetAttachmentPoints{
		QuadPoints: structs.FPDF_FS_QUADPOINTSF{
			X1: float32(attachmentPoints.X1),
			Y1: float32(attachmentPoints.Y1),
			X2: float32(attachmentPoints.X2),
			Y2: float32(attachmentPoints.Y2),
			X3: float32(attachmentPoints.X3),
			Y3: float32(attachmentPoints.Y3),
			X4: float32(attachmentPoints.X4),
			Y4: float32(attachmentPoints.Y4),
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

	rectPointer, _, err := p.CStructFS_RECTF(&request.Rect)
	if err != nil {
		return nil, err
	}
	defer p.Free(rectPointer)

	res, err := p.Module.ExportedFunction("FPDFAnnot_SetRect").Call(p.Context, *annotationHandle.handle, rectPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	rectPointer, rectValue, err := p.CStructFS_RECTF(nil)
	if err != nil {
		return nil, err
	}
	defer p.Free(rectPointer)

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetRect").Call(p.Context, *annotationHandle.handle, rectPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could net set rect")
	}

	rect, err := rectValue()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetRect{
		Rect: structs.FPDF_FS_RECTF{
			Left:   float32(rect.Left),
			Top:    float32(rect.Top),
			Right:  float32(rect.Right),
			Bottom: float32(rect.Bottom),
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
	res, err := p.Module.ExportedFunction("FPDFAnnot_GetVertices").Call(p.Context, *annotationHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	length := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	cVerticesPointer, err := p.FS_POINTFArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer cVerticesPointer.Free()

	if length > 0 {
		// Actually fill the array
		_, err := p.Module.ExportedFunction("FPDFAnnot_GetVertices").Call(p.Context, *annotationHandle.handle, cVerticesPointer.Pointer, length)
		if err != nil {
			return nil, err
		}
	}

	vertices, err := cVerticesPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetVertices{
		Vertices: vertices,
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetInkListCount").Call(p.Context, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

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
	res, err := p.Module.ExportedFunction("FPDFAnnot_GetInkListPath").Call(p.Context, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)), 0, 0)
	if err != nil {
		return nil, err
	}

	length := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	cPathPointer, err := p.FS_POINTFArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer cPathPointer.Free()

	if length > 0 {
		// Actually fill the array
		_, err = p.Module.ExportedFunction("FPDFAnnot_GetInkListPath").Call(p.Context, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)), cPathPointer.Pointer, length)
		if err != nil {
			return nil, err
		}
	}

	cPath, err := cPathPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetInkListPath{
		Path: cPath,
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

	startPointer, err := p.FS_POINTFPointer(nil)
	if err != nil {
		return nil, err
	}
	defer startPointer.Free()

	endPointer, err := p.FS_POINTFPointer(nil)
	if err != nil {
		return nil, err
	}
	defer endPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetLine").Call(p.Context, *annotationHandle.handle, startPointer.Pointer, endPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get line")
	}

	start, err := startPointer.Value()
	if err != nil {
		return nil, err
	}

	end, err := endPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetLine{
		Start: *start,
		End:   *end,
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_SetBorder").Call(p.Context, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.HorizontalRadius)), *(*uint64)(unsafe.Pointer(&request.VerticalRadius)), *(*uint64)(unsafe.Pointer(&request.BorderWidth)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	horizontalRadiusPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer horizontalRadiusPointer.Free()

	verticalRadiusPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer horizontalRadiusPointer.Free()

	borderWidthPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer horizontalRadiusPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetBorder").Call(p.Context, *annotationHandle.handle, horizontalRadiusPointer.Pointer, verticalRadiusPointer.Pointer, borderWidthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get border")
	}

	horizontalRadius, err := horizontalRadiusPointer.Value()
	if err != nil {
		return nil, err
	}

	verticalRadius, err := verticalRadiusPointer.Value()
	if err != nil {
		return nil, err
	}

	borderWidth, err := borderWidthPointer.Value()
	if err != nil {
		return nil, err
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

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAnnot_HasKey").Call(p.Context, *annotationHandle.handle, keyPointer.Pointer)
	if err != nil {
		return nil, err
	}

	hasKey := *(*int32)(unsafe.Pointer(&res[0]))

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

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetValueType").Call(p.Context, *annotationHandle.handle, keyPointer.Pointer)
	if err != nil {
		return nil, err
	}

	valueType := *(*int32)(unsafe.Pointer(&res[0]))

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

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	valuePointer, err := p.CFPDF_WIDESTRING(request.Value)
	if err != nil {
		return nil, err
	}
	defer valuePointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAnnot_SetStringValue").Call(p.Context, *annotationHandle.handle, keyPointer.Pointer, valuePointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	// First get the value length
	res, err := p.Module.ExportedFunction("FPDFAnnot_GetStringValue").Call(p.Context, *annotationHandle.handle, keyPointer.Pointer, 0, 0)
	if err != nil {
		return nil, err
	}

	length := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if uint64(length) == 0 {
		return nil, errors.New("could not get string value")
	}

	charDataPointer, err := p.ByteArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFAnnot_GetStringValue").Call(p.Context, *annotationHandle.handle, keyPointer.Pointer, charDataPointer.Pointer, length)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	valuePointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer valuePointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetNumberValue").Call(p.Context, *annotationHandle.handle, keyPointer.Pointer, valuePointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get number value")
	}

	value, err := valuePointer.Value()
	if err != nil {
		return nil, err
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
		valuePointer, err := p.CFPDF_WIDESTRING(*request.Value)
		if err != nil {
			return nil, err
		}
		defer valuePointer.Free()

		res, err := p.Module.ExportedFunction("FPDFAnnot_SetAP").Call(p.Context, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.AppearanceMode)), valuePointer.Pointer)
		if err != nil {
			return nil, err
		}

		success := *(*int32)(unsafe.Pointer(&res[0]))
		if int(success) == 0 {
			return nil, errors.New("could not set appearance mode")
		}
	} else {
		res, err := p.Module.ExportedFunction("FPDFAnnot_SetAP").Call(p.Context, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.AppearanceMode)), 0)
		if err != nil {
			return nil, err
		}

		success := *(*int32)(unsafe.Pointer(&res[0]))
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
	res, err := p.Module.ExportedFunction("FPDFAnnot_GetAP").Call(p.Context, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.AppearanceMode)), 0, 0)
	if err != nil {
		return nil, err
	}

	length := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if uint64(length) == 0 {
		return nil, errors.New("could not get appearance mode")
	}

	charDataPointer, err := p.ByteArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFAnnot_GetAP").Call(p.Context, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.AppearanceMode)), charDataPointer.Pointer, length)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetAP{
		Value: transformedText,
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

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetLinkedAnnot").Call(p.Context, *annotationHandle.handle, keyPointer.Pointer)
	if err != nil {
		return nil, err
	}

	linkedAnnotation := res[0]
	if linkedAnnotation == 0 {
		return nil, errors.New("could not get linked annotation")
	}

	linkedAnnotationHandle := p.registerAnnotation(&linkedAnnotation)

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

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFlags").Call(p.Context, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	flags := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFAnnot_GetFlags{
		Flags: enums.FPDF_ANNOT_FLAG(flags),
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_SetFlags").Call(p.Context, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.Flags)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFormFieldFlags").Call(p.Context, *formHandle.handle, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	flags := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFAnnot_GetFormFieldFlags{
		Flags: enums.FPDF_FORMFLAG(flags),
	}, nil
}

// FPDFAnnot_SetFormFieldFlags sets the form field flags for an interactive form annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetFormFieldFlags(request *requests.FPDFAnnot_SetFormFieldFlags) (*responses.FPDFAnnot_SetFormFieldFlags, error) {
	p.Lock()
	defer p.Unlock()

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFAnnot_SetFormFieldFlags").Call(p.Context, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.Flags)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set form field flags")
	}

	return &responses.FPDFAnnot_SetFormFieldFlags{}, nil
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

	pointPointer, err := p.FS_POINTFPointer(&request.Point)
	if err != nil {
		return nil, err
	}
	defer pointPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFormFieldAtPoint").Call(p.Context, *formHandle.handle, *pageHandle.handle, pointPointer.Pointer)
	if err != nil {
		return nil, err
	}

	formField := res[0]
	if formField == 0 {
		return nil, errors.New("could not get form field")
	}

	formFieldHandle := p.registerAnnotation(&formField)

	return &responses.FPDFAnnot_GetFormFieldAtPoint{
		Annotation: formFieldHandle.nativeRef,
	}, nil
}

// FPDFAnnot_GetFormAdditionalActionJavaScript returns the JavaScript of an event of the annotation's additional actions.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormAdditionalActionJavaScript(request *requests.FPDFAnnot_GetFormAdditionalActionJavaScript) (*responses.FPDFAnnot_GetFormAdditionalActionJavaScript, error) {
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
	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFormAdditionalActionJavaScript").Call(p.Context, *formHandle.handle, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.Event)), 0, 0)
	if err != nil {
		return nil, err
	}

	length := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if length == 0 {
		return nil, errors.New("could not get form additional action JavaScript")
	}

	charDataPointer, err := p.ByteArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFAnnot_GetFormAdditionalActionJavaScript").Call(p.Context, *formHandle.handle, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.Event)), charDataPointer.Pointer, length)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetFormAdditionalActionJavaScript{
		FormAdditionalActionJavaScript: transformedText,
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
	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFormFieldName").Call(p.Context, *formHandle.handle, *annotationHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	length := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if length == 0 {
		return nil, errors.New("could not get form field name")
	}

	charDataPointer, err := p.ByteArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFAnnot_GetFormFieldName").Call(p.Context, *formHandle.handle, *annotationHandle.handle, charDataPointer.Pointer, length)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetFormFieldName{
		FormFieldName: transformedText,
	}, nil
}

// FPDFAnnot_GetFormFieldAlternateName returns the alternate name of an annotation, which is an interactive form annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFormFieldAlternateName(request *requests.FPDFAnnot_GetFormFieldAlternateName) (*responses.FPDFAnnot_GetFormFieldAlternateName, error) {
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
	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFormFieldAlternateName").Call(p.Context, *formHandle.handle, *annotationHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	length := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if length == 0 {
		return nil, errors.New("could not get form field name")
	}

	charDataPointer, err := p.ByteArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFAnnot_GetFormFieldAlternateName").Call(p.Context, *formHandle.handle, *annotationHandle.handle, charDataPointer.Pointer, length)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetFormFieldAlternateName{
		FormFieldAlternateName: transformedText,
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFormFieldType").Call(p.Context, *formHandle.handle, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	formFieldType := *(*int32)(unsafe.Pointer(&res[0]))
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
	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFormFieldValue").Call(p.Context, *formHandle.handle, *annotationHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	length := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if length == 0 {
		return nil, errors.New("could not get form field name")
	}

	charDataPointer, err := p.ByteArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFAnnot_GetFormFieldValue").Call(p.Context, *formHandle.handle, *annotationHandle.handle, charDataPointer.Pointer, length)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetOptionCount").Call(p.Context, *formHandle.handle, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	optionCount := *(*int32)(unsafe.Pointer(&res[0]))
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
	res, err := p.Module.ExportedFunction("FPDFAnnot_GetOptionLabel").Call(p.Context, *formHandle.handle, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)), 0, 0)
	if err != nil {
		return nil, err
	}

	length := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if length == 0 {
		return nil, errors.New("could not get form field name")
	}

	charDataPointer, err := p.ByteArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFAnnot_GetOptionLabel").Call(p.Context, *formHandle.handle, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)), charDataPointer.Pointer, length)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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

	res, err := p.Module.ExportedFunction("FPDFAnnot_IsOptionSelected").Call(p.Context, *formHandle.handle, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	isOptionSelected := *(*int32)(unsafe.Pointer(&res[0]))

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

	fontSizePointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer fontSizePointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFontSize").Call(p.Context, *formHandle.handle, *annotationHandle.handle, fontSizePointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get font size")
	}

	fontSize, err := fontSizePointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetFontSize{
		FontSize: float32(fontSize),
	}, nil
}

// FPDFAnnot_SetFontColor Set the text color of an annotation.
// Currently supported subtypes: freetext.
// The range for the color components is 0 to 255.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_SetFontColor(request *requests.FPDFAnnot_SetFontColor) (*responses.FPDFAnnot_SetFontColor, error) {
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_SetFontColor").Call(p.Context, *formHandle.handle, *annotationHandle.handle, *(*uint64)(unsafe.Pointer(&request.R)), *(*uint64)(unsafe.Pointer(&request.G)), *(*uint64)(unsafe.Pointer(&request.B)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set font color")
	}

	return &responses.FPDFAnnot_SetFontColor{}, nil
}

// FPDFAnnot_GetFontColor returns the RGB value of the font color for an annotation with variable text.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFontColor(request *requests.FPDFAnnot_GetFontColor) (*responses.FPDFAnnot_GetFontColor, error) {
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

	rPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer rPointer.Free()

	gPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer gPointer.Free()

	bPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer bPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFontColor").Call(p.Context, *formHandle.handle, *annotationHandle.handle, rPointer.Pointer, gPointer.Pointer, bPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get font color")
	}

	rValue, err := rPointer.Value()
	if err != nil {
		return nil, err
	}

	gValue, err := gPointer.Value()
	if err != nil {
		return nil, err
	}

	bValue, err := bPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAnnot_GetFontColor{
		R: rValue,
		G: gValue,
		B: bValue,
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_IsChecked").Call(p.Context, *formHandle.handle, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	isChecked := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFAnnot_IsChecked{
		IsChecked: int(isChecked) == 1,
	}, nil
}

// FPDFAnnot_SetFocusableSubtypes sets the list of focusable annotation subtypes. Annotations of subtype
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

	subtypesSize := uint64(len(request.Subtypes))
	focusableSubtypesPointer, err := p.IntArrayPointer(subtypesSize)
	if err != nil {
		return nil, err
	}
	defer focusableSubtypesPointer.Free()

	for i := range request.Subtypes {
		success := p.Module.Memory().WriteUint32Le(uint32(focusableSubtypesPointer.Pointer+(p.CSizeInt()*uint64(i))), *(*uint32)(unsafe.Pointer(&request.Subtypes[i])))
		if !success {
			return nil, errors.New("could not write focusable subtype to memory")
		}
	}

	res, err := p.Module.ExportedFunction("FPDFAnnot_SetFocusableSubtypes").Call(p.Context, *formHandle.handle, focusableSubtypesPointer.Pointer, subtypesSize)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFocusableSubtypesCount").Call(p.Context, *formHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))
	if int(count) == -1 {
		return nil, errors.New("could net get focusable subtypes count")
	}

	return &responses.FPDFAnnot_GetFocusableSubtypesCount{
		FocusableSubtypesCount: int(count),
	}, nil
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFocusableSubtypesCount").Call(p.Context, *formHandle.handle)
	if err != nil {
		return nil, err
	}

	count := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if int(count) == -1 {
		return nil, errors.New("could net get focusable subtypes")
	}

	focusableSubtypesPointer, err := p.IntArrayPointer(count)
	if err != nil {
		return nil, err
	}

	goFocusableSubtypes := make([]enums.FPDF_ANNOTATION_SUBTYPE, int(count))

	if int(count) > 0 {
		res, err = p.Module.ExportedFunction("FPDFAnnot_GetFocusableSubtypes").Call(p.Context, *formHandle.handle, focusableSubtypesPointer.Pointer, count)
		if err != nil {
			return nil, err
		}

		success := *(*int32)(unsafe.Pointer(&res[0]))
		if int(success) == 0 {
			return nil, errors.New("could net get focusable subtypes")
		}

		focusableSubtypes, err := focusableSubtypesPointer.Value()
		if err != nil {
			return nil, err
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetLink").Call(p.Context, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	link := res[0]
	if link == 0 {
		return nil, errors.New("could not get link")
	}

	linkHandle := p.registerLink(&link)

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

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFormControlCount").Call(p.Context, *formHandle.handle, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	formControlCount := uint64(*(*int32)(unsafe.Pointer(&res[0])))
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

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFormControlIndex").Call(p.Context, *formHandle.handle, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	formControlIndex := uint64(*(*int32)(unsafe.Pointer(&res[0])))
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
	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFormFieldExportValue").Call(p.Context, *formHandle.handle, *annotationHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	length := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if length == 0 {
		return nil, errors.New("could not get form field export value")
	}

	charDataPointer, err := p.ByteArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFAnnot_GetFormFieldExportValue").Call(p.Context, *formHandle.handle, *annotationHandle.handle, charDataPointer.Pointer, length)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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

	uriStrPointer, err := p.CString(request.URI)
	if err != nil {
		return nil, err
	}
	defer uriStrPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAnnot_SetURI").Call(p.Context, *annotationHandle.handle, uriStrPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could net set uri")
	}

	return &responses.FPDFAnnot_SetURI{}, nil
}

// FPDFAnnot_GetFileAttachment get the attachment from the given annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_GetFileAttachment(request *requests.FPDFAnnot_GetFileAttachment) (*responses.FPDFAnnot_GetFileAttachment, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFAnnot_GetFileAttachment").Call(p.Context, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	handle := res[0]
	if handle == 0 {
		return nil, errors.New("could not get attachment object")
	}

	attachmentHandle := p.registerAttachment(&handle, documentHandle)

	return &responses.FPDFAnnot_GetFileAttachment{
		Attachment: attachmentHandle.nativeRef,
	}, nil
}

// FPDFAnnot_AddFileAttachment Add an embedded file to the given annotation.
// Experimental API.
func (p *PdfiumImplementation) FPDFAnnot_AddFileAttachment(request *requests.FPDFAnnot_AddFileAttachment) (*responses.FPDFAnnot_AddFileAttachment, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	namePointer, err := p.CFPDF_WIDESTRING(request.Name)
	if err != nil {
		return nil, err
	}
	defer namePointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAnnot_AddFileAttachment").Call(p.Context, *annotationHandle.handle, namePointer.Pointer)
	if err != nil {
		return nil, err
	}

	handle := res[0]
	if handle == 0 {
		return nil, errors.New("could not get attachment object")
	}

	attachmentHandle := p.registerAttachment(&handle, documentHandle)

	return &responses.FPDFAnnot_AddFileAttachment{
		Attachment: attachmentHandle.nativeRef,
	}, nil
}
