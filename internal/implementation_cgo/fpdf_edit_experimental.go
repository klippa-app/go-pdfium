//go:build pdfium_experimental
// +build pdfium_experimental

package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdf_edit.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"
)

// FPDFPage_RemoveObject removes an object from a page.
// Ownership is transferred to the caller. Call FPDFPageObj_Destroy() to free
// it.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_RemoveObject(request *requests.FPDFPage_RemoveObject) (*responses.FPDFPage_RemoveObject, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	success := C.FPDFPage_RemoveObject(pageHandle.handle, pageObjectHandle.handle)
	if int(success) == 0 {
		return nil, errors.New("could not remove object")
	}

	return &responses.FPDFPage_RemoveObject{}, nil
}

// FPDFPageObj_TransformF transforms the page object by the given matrix.
// The matrix is composed as:
//
//	|a c e|
//	|b d f|
//
// and can be used to scale, rotate, shear and translate the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_TransformF(request *requests.FPDFPageObj_TransformF) (*responses.FPDFPageObj_TransformF, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	matrix := C.FS_MATRIX{}
	matrix.a = C.float(request.Transform.A)
	matrix.b = C.float(request.Transform.B)
	matrix.c = C.float(request.Transform.C)
	matrix.d = C.float(request.Transform.D)
	matrix.e = C.float(request.Transform.E)
	matrix.f = C.float(request.Transform.F)

	success := C.FPDFPageObj_TransformF(pageObjectHandle.handle, &matrix)
	if int(success) == 0 {
		return nil, errors.New("could not transform object")
	}

	return &responses.FPDFPageObj_TransformF{}, nil
}

// FPDFPageObj_GetMatrix returns the transform matrix of a page object.
// The matrix is composed as:
//
//	|a c e|
//	|b d f|
//
// and can be used to scale, rotate, shear and translate the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetMatrix(request *requests.FPDFPageObj_GetMatrix) (*responses.FPDFPageObj_GetMatrix, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	matrix := C.FS_MATRIX{}

	success := C.FPDFPageObj_GetMatrix(pageObjectHandle.handle, &matrix)
	if int(success) == 0 {
		return nil, errors.New("could not get page object matrix")
	}

	return &responses.FPDFPageObj_GetMatrix{
		Matrix: structs.FPDF_FS_MATRIX{
			A: float32(matrix.a),
			B: float32(matrix.b),
			C: float32(matrix.c),
			D: float32(matrix.d),
			E: float32(matrix.e),
			F: float32(matrix.f),
		},
	}, nil
}

// FPDFPageObj_SetMatrix sets the transform matrix on a page object.
// The matrix is composed as:
//
//	|a c e|
//	|b d f|
//
// and can be used to scale, rotate, shear and translate the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_SetMatrix(request *requests.FPDFPageObj_SetMatrix) (*responses.FPDFPageObj_SetMatrix, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	matrix := C.FS_MATRIX{}
	matrix.a = C.float(request.Transform.A)
	matrix.b = C.float(request.Transform.B)
	matrix.c = C.float(request.Transform.C)
	matrix.d = C.float(request.Transform.D)
	matrix.e = C.float(request.Transform.E)
	matrix.f = C.float(request.Transform.F)

	success := C.FPDFPageObj_SetMatrix(pageObjectHandle.handle, &matrix)
	if int(success) == 0 {
		return nil, errors.New("could not set page object matrix")
	}

	return &responses.FPDFPageObj_SetMatrix{}, nil
}

// FPDFPageObj_GetMarkedContentID returns the marked content ID of a page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetMarkedContentID(request *requests.FPDFPageObj_GetMarkedContentID) (*responses.FPDFPageObj_GetMarkedContentID, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	markedContentID := C.FPDFPageObj_GetMarkedContentID(pageObjectHandle.handle)

	return &responses.FPDFPageObj_GetMarkedContentID{
		MarkedContentID: int(markedContentID),
	}, nil
}

// FPDFPageObj_CountMarks returns the count of content marks in a page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_CountMarks(request *requests.FPDFPageObj_CountMarks) (*responses.FPDFPageObj_CountMarks, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	count := C.FPDFPageObj_CountMarks(pageObjectHandle.handle)

	return &responses.FPDFPageObj_CountMarks{
		Count: int(count),
	}, nil
}

// FPDFPageObj_GetMark returns the content mark of a page object at the given index.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetMark(request *requests.FPDFPageObj_GetMark) (*responses.FPDFPageObj_GetMark, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	mark := C.FPDFPageObj_GetMark(pageObjectHandle.handle, C.ulong(request.Index))
	if mark == nil {
		return nil, errors.New("could not get mark")
	}

	markHandle := p.registerPageObjectMark(mark)

	return &responses.FPDFPageObj_GetMark{
		Mark: markHandle.nativeRef,
	}, nil
}

// FPDFPageObj_AddMark adds a new content mark to the given page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_AddMark(request *requests.FPDFPageObj_AddMark) (*responses.FPDFPageObj_AddMark, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	name := C.CString(request.Name)
	defer C.free(unsafe.Pointer(name))

	mark := C.FPDFPageObj_AddMark(pageObjectHandle.handle, name)
	if mark == nil {
		return nil, errors.New("could not add mark")
	}

	markHandle := p.registerPageObjectMark(mark)

	return &responses.FPDFPageObj_AddMark{
		Mark: markHandle.nativeRef,
	}, nil
}

// FPDFPageObj_RemoveMark removes the given content mark from the given page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_RemoveMark(request *requests.FPDFPageObj_RemoveMark) (*responses.FPDFPageObj_RemoveMark, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	success := C.FPDFPageObj_RemoveMark(pageObjectHandle.handle, pageObjectMarkHandle.handle)
	if int(success) == 0 {
		return nil, errors.New("could not remove mark")
	}

	return &responses.FPDFPageObj_RemoveMark{}, nil
}

// FPDFPageObjMark_GetName returns the name of a content mark.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetName(request *requests.FPDFPageObjMark_GetName) (*responses.FPDFPageObjMark_GetName, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	nameLength := C.ulong(0)

	success := C.FPDFPageObjMark_GetName(pageObjectMarkHandle.handle, nil, 0, &nameLength)
	if int(success) == 0 {
		return nil, errors.New("could not get name")
	}

	if uint64(nameLength) == 0 {
		return nil, errors.New("could not get name")
	}

	charData := make([]byte, uint64(nameLength))
	C.FPDFPageObjMark_GetName(pageObjectMarkHandle.handle, (*C.FPDF_WCHAR)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)), &nameLength)

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObjMark_GetName{
		Name: transformedText,
	}, nil
}

// FPDFPageObjMark_CountParams returns the number of key/value pair parameters in the given mark.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_CountParams(request *requests.FPDFPageObjMark_CountParams) (*responses.FPDFPageObjMark_CountParams, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	count := C.FPDFPageObjMark_CountParams(pageObjectMarkHandle.handle)

	return &responses.FPDFPageObjMark_CountParams{
		Count: int(count),
	}, nil
}

// FPDFPageObjMark_GetParamKey returns the key of a property in a content mark.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetParamKey(request *requests.FPDFPageObjMark_GetParamKey) (*responses.FPDFPageObjMark_GetParamKey, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	keyLength := C.ulong(0)

	success := C.FPDFPageObjMark_GetParamKey(pageObjectMarkHandle.handle, C.ulong(request.Index), nil, 0, &keyLength)
	if int(success) == 0 {
		return nil, errors.New("could not get key")
	}

	if uint64(keyLength) == 0 {
		return nil, errors.New("could not get key")
	}

	charData := make([]byte, uint64(keyLength))
	C.FPDFPageObjMark_GetParamKey(pageObjectMarkHandle.handle, C.ulong(request.Index), (*C.FPDF_WCHAR)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)), &keyLength)

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObjMark_GetParamKey{
		Key: transformedText,
	}, nil
}

// FPDFPageObjMark_GetParamValueType returns the type of the value of a property in a content mark by key.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetParamValueType(request *requests.FPDFPageObjMark_GetParamValueType) (*responses.FPDFPageObjMark_GetParamValueType, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	key := C.CString(request.Key)
	defer C.free(unsafe.Pointer(key))

	valueType := C.FPDFPageObjMark_GetParamValueType(pageObjectMarkHandle.handle, key)

	return &responses.FPDFPageObjMark_GetParamValueType{
		ValueType: enums.FPDF_OBJECT_TYPE(valueType),
	}, nil
}

// FPDFPageObjMark_GetParamIntValue returns the value of a number property in a content mark by key as int.
// FPDFPageObjMark_GetParamValueType() should have returned FPDF_OBJECT_NUMBER
// for this property.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetParamIntValue(request *requests.FPDFPageObjMark_GetParamIntValue) (*responses.FPDFPageObjMark_GetParamIntValue, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	key := C.CString(request.Key)
	defer C.free(unsafe.Pointer(key))

	intValue := C.int(0)

	success := C.FPDFPageObjMark_GetParamIntValue(pageObjectMarkHandle.handle, key, &intValue)
	if int(success) == 0 {
		return nil, errors.New("could not get value")
	}

	return &responses.FPDFPageObjMark_GetParamIntValue{
		Value: int(intValue),
	}, nil
}

// FPDFPageObjMark_GetParamStringValue returns the value of a string property in a content mark by key.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetParamStringValue(request *requests.FPDFPageObjMark_GetParamStringValue) (*responses.FPDFPageObjMark_GetParamStringValue, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	key := C.CString(request.Key)
	defer C.free(unsafe.Pointer(key))

	valueLength := C.ulong(0)

	success := C.FPDFPageObjMark_GetParamStringValue(pageObjectMarkHandle.handle, key, nil, 0, &valueLength)
	if int(success) == 0 {
		return nil, errors.New("could not get value")
	}

	if uint64(valueLength) == 0 {
		return nil, errors.New("could not get value")
	}

	charData := make([]byte, uint64(valueLength))
	C.FPDFPageObjMark_GetParamStringValue(pageObjectMarkHandle.handle, key, (*C.FPDF_WCHAR)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)), &valueLength)

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObjMark_GetParamStringValue{
		Value: transformedText,
	}, nil
}

// FPDFPageObjMark_GetParamBlobValue returns the value of a blob property in a content mark by key.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetParamBlobValue(request *requests.FPDFPageObjMark_GetParamBlobValue) (*responses.FPDFPageObjMark_GetParamBlobValue, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	key := C.CString(request.Key)
	defer C.free(unsafe.Pointer(key))

	valueLength := C.ulong(0)

	success := C.FPDFPageObjMark_GetParamBlobValue(pageObjectMarkHandle.handle, key, nil, 0, &valueLength)
	if int(success) == 0 {
		return nil, errors.New("could not get value")
	}

	if uint64(valueLength) == 0 {
		return nil, errors.New("could not get value")
	}

	valueData := make([]byte, uint64(valueLength))
	C.FPDFPageObjMark_GetParamBlobValue(pageObjectMarkHandle.handle, key, (*C.uchar)(unsafe.Pointer(&valueData[0])), C.ulong(len(valueData)), &valueLength)

	return &responses.FPDFPageObjMark_GetParamBlobValue{
		Value: valueData,
	}, nil
}

// FPDFPageObjMark_SetIntParam sets the value of an int property in a content mark by key. If a parameter
// with the given key exists, its value is set to the given value. Otherwise, it is added as
// a new parameter.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_SetIntParam(request *requests.FPDFPageObjMark_SetIntParam) (*responses.FPDFPageObjMark_SetIntParam, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	key := C.CString(request.Key)
	defer C.free(unsafe.Pointer(key))

	success := C.FPDFPageObjMark_SetIntParam(documentHandle.handle, pageObjectHandle.handle, pageObjectMarkHandle.handle, key, C.int(request.Value))
	if int(success) == 0 {
		return nil, errors.New("could not set value")
	}

	return &responses.FPDFPageObjMark_SetIntParam{}, nil
}

// FPDFPageObjMark_SetStringParam sets the value of a string property in a content mark by key. If a parameter
// with the given key exists, its value is set to the given value. Otherwise, it is added as
// a new parameter.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_SetStringParam(request *requests.FPDFPageObjMark_SetStringParam) (*responses.FPDFPageObjMark_SetStringParam, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	key := C.CString(request.Key)
	defer C.free(unsafe.Pointer(key))

	value := C.CString(request.Value)
	defer C.free(unsafe.Pointer(value))

	success := C.FPDFPageObjMark_SetStringParam(documentHandle.handle, pageObjectHandle.handle, pageObjectMarkHandle.handle, key, value)
	if int(success) == 0 {
		return nil, errors.New("could not set value")
	}

	return &responses.FPDFPageObjMark_SetStringParam{}, nil
}

// FPDFPageObjMark_SetBlobParam sets the value of a blob property in a content mark by key. If a parameter
// with the given key exists, its value is set to the given value. Otherwise, it is added as
// a new parameter.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_SetBlobParam(request *requests.FPDFPageObjMark_SetBlobParam) (*responses.FPDFPageObjMark_SetBlobParam, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	if request.Value == nil || len(request.Value) == 0 {
		return nil, errors.New("blob value cant be empty")
	}

	key := C.CString(request.Key)
	defer C.free(unsafe.Pointer(key))

	success := C.FPDFPageObjMark_SetBlobParam(documentHandle.handle, pageObjectHandle.handle, pageObjectMarkHandle.handle, key, (*C.uchar)(unsafe.Pointer(&request.Value[0])), C.ulong(len(request.Value)))
	if int(success) == 0 {
		return nil, errors.New("could not set value")
	}

	return &responses.FPDFPageObjMark_SetBlobParam{}, nil
}

// FPDFPageObjMark_RemoveParam removes a property from a content mark by key.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_RemoveParam(request *requests.FPDFPageObjMark_RemoveParam) (*responses.FPDFPageObjMark_RemoveParam, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	key := C.CString(request.Key)
	defer C.free(unsafe.Pointer(key))

	success := C.FPDFPageObjMark_RemoveParam(pageObjectHandle.handle, pageObjectMarkHandle.handle, key)
	if int(success) == 0 {
		return nil, errors.New("could not set value")
	}

	return &responses.FPDFPageObjMark_RemoveParam{}, nil
}

// FPDFImageObj_GetRenderedBitmap returns a bitmap rasterization of the given image object that takes the image mask and
// image matrix into account. To render correctly, the caller must provide the
// document associated with the image object. If there is a page associated
// with the image object the caller should provide that as well.
// The returned bitmap will be owned by the caller, and FPDFBitmap_Destroy()
// must be called on the returned bitmap when it is no longer needed.
// Experimental API.
func (p *PdfiumImplementation) FPDFImageObj_GetRenderedBitmap(request *requests.FPDFImageObj_GetRenderedBitmap) (*responses.FPDFImageObj_GetRenderedBitmap, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	imageObjectHandle, err := p.getPageObjectHandle(request.ImageObject)
	if err != nil {
		return nil, err
	}

	bitmap := C.FPDFImageObj_GetRenderedBitmap(documentHandle.handle, pageHandle.handle, imageObjectHandle.handle)
	if bitmap == nil {
		return nil, errors.New("could not get bitmap")
	}

	bitmapHandle := p.registerBitmap(bitmap)

	return &responses.FPDFImageObj_GetRenderedBitmap{
		Bitmap: bitmapHandle.nativeRef,
	}, nil
}

// FPDFPageObj_GetRotatedBounds Get the quad points that bounds the page object.
// Similar to FPDFPageObj_GetBounds(), this returns the bounds of a page
// object. When the object is rotated by a non-multiple of 90 degrees, this API
// returns a tighter bound that cannot be represented with just the 4 sides of
// a rectangle.
//
// Currently only works the following page object types: FPDF_PAGEOBJ_TEXT and
// FPDF_PAGEOBJ_IMAGE.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetRotatedBounds(request *requests.FPDFPageObj_GetRotatedBounds) (*responses.FPDFPageObj_GetRotatedBounds, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	var quadPoints C.FS_QUADPOINTSF

	success := C.FPDFPageObj_GetRotatedBounds(pageObjectHandle.handle, &quadPoints)
	if int(success) == 0 {
		return nil, errors.New("could not get rotated bounds for page object")
	}

	return &responses.FPDFPageObj_GetRotatedBounds{
		QuadPoints: structs.FPDF_FS_QUADPOINTSF{
			X1: float32(quadPoints.x1),
			Y1: float32(quadPoints.y1),
			X2: float32(quadPoints.x2),
			Y2: float32(quadPoints.y2),
			X3: float32(quadPoints.x3),
			Y3: float32(quadPoints.y3),
			X4: float32(quadPoints.x4),
			Y4: float32(quadPoints.y4),
		},
	}, nil
}

// FPDFPageObj_GetDashPhase returns the line dash phase of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetDashPhase(request *requests.FPDFPageObj_GetDashPhase) (*responses.FPDFPageObj_GetDashPhase, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	dashPhase := C.float(0)
	success := C.FPDFPageObj_GetDashPhase(pageObjectHandle.handle, &dashPhase)
	if int(success) == 0 {
		return nil, errors.New("could not get dash phase")
	}

	return &responses.FPDFPageObj_GetDashPhase{
		DashPhase: float32(dashPhase),
	}, nil
}

// FPDFPageObj_SetDashPhase sets the line dash phase of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_SetDashPhase(request *requests.FPDFPageObj_SetDashPhase) (*responses.FPDFPageObj_SetDashPhase, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	success := C.FPDFPageObj_SetDashPhase(pageObjectHandle.handle, C.float(request.DashPhase))
	if int(success) == 0 {
		return nil, errors.New("could not set dash phase")
	}

	return &responses.FPDFPageObj_SetDashPhase{}, nil
}

// FPDFPageObj_GetDashCount returns the line dash array size of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetDashCount(request *requests.FPDFPageObj_GetDashCount) (*responses.FPDFPageObj_GetDashCount, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	dashCount := C.FPDFPageObj_GetDashCount(pageObjectHandle.handle)

	return &responses.FPDFPageObj_GetDashCount{
		DashCount: int(dashCount),
	}, nil
}

// FPDFPageObj_GetDashArray returns the line dash array of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetDashArray(request *requests.FPDFPageObj_GetDashArray) (*responses.FPDFPageObj_GetDashArray, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	// First get the Dash Count.
	dashCount := C.FPDFPageObj_GetDashCount(pageObjectHandle.handle)

	convertedData := make([]float32, 0)
	if int(dashCount) > 0 {
		// Create an array that's big enough.
		valueData := make([]C.float, int(dashCount))
		C.FPDFPageObj_GetDashArray(pageObjectHandle.handle, &valueData[0], C.size_t(len(valueData)))

		convertedData = make([]float32, int(dashCount))
		for i := range valueData {
			convertedData[i] = float32(valueData[i])
		}
	}

	return &responses.FPDFPageObj_GetDashArray{
		DashArray: convertedData,
	}, nil
}

// FPDFPageObj_SetDashArray sets the line dash array of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_SetDashArray(request *requests.FPDFPageObj_SetDashArray) (*responses.FPDFPageObj_SetDashArray, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	// Create an array that's big enough.
	valueData := make([]C.float, len(request.DashArray))
	for i := range request.DashArray {
		valueData[i] = C.float(request.DashArray[i])
	}

	C.FPDFPageObj_SetDashArray(pageObjectHandle.handle, &valueData[0], C.size_t(len(valueData)), C.float(request.DashPhase))

	return &responses.FPDFPageObj_SetDashArray{}, nil
}

// FPDFText_LoadStandardFont loads one of the standard 14 fonts per PDF spec 1.7 page 416. The preferred
// way of using font style is using a dash to separate the name from the style,
// for example 'Helvetica-BoldItalic'.
// The loaded font can be closed using FPDFFont_Close().
// Experimental API.
func (p *PdfiumImplementation) FPDFText_LoadStandardFont(request *requests.FPDFText_LoadStandardFont) (*responses.FPDFText_LoadStandardFont, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	fontName := C.CString(request.Font)
	defer C.free(unsafe.Pointer(fontName))

	font := C.FPDFText_LoadStandardFont(documentHandle.handle, fontName)
	if font == nil {
		return nil, errors.New("could not load standard font")
	}

	fontHandle := p.registerFont(font)

	return &responses.FPDFText_LoadStandardFont{
		Font: fontHandle.nativeRef,
	}, nil
}

// FPDFText_LoadCidType2Font returns a font object loaded from a stream of data for a type 2 CID font.
// The font is loaded into the document. Unlike FPDFText_LoadFont(), the ToUnicode data and the CIDToGIDMap
// data are caller provided, instead of auto-generated.
// The loaded font can be closed using FPDFFont_Close().
// Experimental API.
func (p *PdfiumImplementation) FPDFText_LoadCidType2Font(request *requests.FPDFText_LoadCidType2Font) (*responses.FPDFText_LoadCidType2Font, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	toUnicodeCmap := C.CString(request.ToUnicodeCmap)
	defer C.free(unsafe.Pointer(toUnicodeCmap))

	font := C.FPDFText_LoadCidType2Font(documentHandle.handle, (*C.uchar)(unsafe.Pointer(&request.FontData[0])), C.uint32_t(len(request.FontData)), toUnicodeCmap, (*C.uchar)(unsafe.Pointer(&request.CIDToGIDMapData[0])), C.uint32_t(len(request.CIDToGIDMapData)))
	if font == nil {
		return nil, errors.New("could not load CID Type2 font")
	}

	fontHandle := p.registerFont(font)

	return &responses.FPDFText_LoadCidType2Font{
		Font: fontHandle.nativeRef,
	}, nil
}

// FPDFTextObj_SetTextRenderMode sets the text rendering mode of a text object.
// Experimental API.
func (p *PdfiumImplementation) FPDFTextObj_SetTextRenderMode(request *requests.FPDFTextObj_SetTextRenderMode) (*responses.FPDFTextObj_SetTextRenderMode, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	result := C.FPDFTextObj_SetTextRenderMode(pageObjectHandle.handle, C.FPDF_TEXT_RENDERMODE(request.TextRenderMode))
	if int(result) == 0 {
		return nil, errors.New("could not set text render mode")
	}

	return &responses.FPDFTextObj_SetTextRenderMode{}, nil
}

// FPDFTextObj_GetRenderedBitmap returns a bitmap rasterization of the given text object.
// To render correctly, the caller must provide the document associated with the text object.
// If there is a page associated with text object, the caller should provide that as well.
// The returned bitmap will be owned by the caller, and FPDFBitmap_Destroy()
// must be called on the returned bitmap when it is no longer needed.
// Experimental API.
func (p *PdfiumImplementation) FPDFTextObj_GetRenderedBitmap(request *requests.FPDFTextObj_GetRenderedBitmap) (*responses.FPDFTextObj_GetRenderedBitmap, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	var pageHandle C.FPDF_PAGE
	if request.Page != "" {
		pageHandleReference, err := p.getPageHandle(request.Page)
		if err != nil {
			return nil, err
		}

		pageHandle = pageHandleReference.handle
	}

	bitmap := C.FPDFTextObj_GetRenderedBitmap(documentHandle.handle, pageHandle, pageObjectHandle.handle, C.float(request.Scale))
	if bitmap == nil {
		return nil, errors.New("could not render text object as bitmap")
	}

	bitmapHandle := p.registerBitmap(bitmap)

	return &responses.FPDFTextObj_GetRenderedBitmap{
		Bitmap: bitmapHandle.nativeRef,
	}, nil
}

// FPDFTextObj_GetFont returns the font of a text object.
// Experimental API.
func (p *PdfiumImplementation) FPDFTextObj_GetFont(request *requests.FPDFTextObj_GetFont) (*responses.FPDFTextObj_GetFont, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	font := C.FPDFTextObj_GetFont(pageObjectHandle.handle)
	fontHandle := p.registerFont(font)

	return &responses.FPDFTextObj_GetFont{
		Font: fontHandle.nativeRef,
	}, nil
}

// FPDFFont_GetBaseFontName returns the base font name of a font.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetBaseFontName(request *requests.FPDFFont_GetBaseFontName) (*responses.FPDFFont_GetBaseFontName, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	// First get the text value length.
	nameSize := C.FPDFFont_GetBaseFontName(fontHandle.handle, nil, 0)
	if nameSize == 0 {
		return nil, errors.New("could not get font name")
	}

	charData := make([]byte, nameSize)
	C.FPDFFont_GetBaseFontName(fontHandle.handle, (*C.char)(unsafe.Pointer(&charData[0])), C.size_t(len(charData)))

	return &responses.FPDFFont_GetBaseFontName{
		BaseFontName: string(charData[:len(charData)-1]), // Remove NULL-terminator
	}, nil
}

// FPDFFont_GetFamilyName returns the family name of a font.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetFamilyName(request *requests.FPDFFont_GetFamilyName) (*responses.FPDFFont_GetFamilyName, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	// First get the text value length.
	nameSize := C.FPDFFont_GetFamilyName(fontHandle.handle, nil, 0)
	if nameSize == 0 {
		return nil, errors.New("could not get font name")
	}

	charData := make([]byte, nameSize)
	C.FPDFFont_GetFamilyName(fontHandle.handle, (*C.char)(unsafe.Pointer(&charData[0])), C.size_t(len(charData)))

	return &responses.FPDFFont_GetFamilyName{
		FamilyName: string(charData[:len(charData)-1]), // Remove NULL-terminator
	}, nil
}

// FPDFFont_GetFontData returns the decoded data from the given font.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetFontData(request *requests.FPDFFont_GetFontData) (*responses.FPDFFont_GetFontData, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	outBufLen := C.size_t(0)

	// First get the font data length.
	success := C.FPDFFont_GetFontData(fontHandle.handle, (*C.uchar)(nil), 0, &outBufLen)
	if int(success) != 1 {
		return nil, errors.New("could not get font data")
	}

	if int(outBufLen) == 0 {
		return nil, errors.New("could not get font data")
	}

	fontData := make([]byte, outBufLen)
	success = C.FPDFFont_GetFontData(fontHandle.handle, (*C.uchar)(unsafe.Pointer(&fontData[0])), C.size_t(len(fontData)), &outBufLen)
	if int(success) != 1 {
		return nil, errors.New("could not get font data")
	}

	return &responses.FPDFFont_GetFontData{
		FontData: fontData,
	}, nil
}

// FPDFFont_GetIsEmbedded returns whether the given font is embedded or not.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetIsEmbedded(request *requests.FPDFFont_GetIsEmbedded) (*responses.FPDFFont_GetIsEmbedded, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	isEmbedded := C.FPDFFont_GetIsEmbedded(fontHandle.handle)
	return &responses.FPDFFont_GetIsEmbedded{
		IsEmbedded: int(isEmbedded) == 1,
	}, nil
}

// FPDFFont_GetFlags returns the descriptor flags of a font.
// Returns the bit flags specifying various characteristics of the font as
// defined in ISO 32000-1:2008, table 123.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetFlags(request *requests.FPDFFont_GetFlags) (*responses.FPDFFont_GetFlags, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	flags := C.FPDFFont_GetFlags(fontHandle.handle)
	if int(flags) == -1 {
		return nil, errors.New("could not get font flags")
	}

	fontFlags := &responses.FPDFFont_GetFlags{
		Flags: uint32(flags),
	}

	FixedPitch := uint32(1 << 1)
	Serif := uint32(1 << 2)
	Symbolic := uint32(1 << 3)
	Script := uint32(1 << 4)
	Nonsymbolic := uint32(1 << 6)
	Italic := uint32(1 << 7)
	AllCap := uint32(1 << 17)
	SmallCap := uint32(1 << 18)
	ForceBold := uint32(1 << 19)

	hasFlag := func(flag uint32) bool {
		if fontFlags.Flags&flag > 0 {
			return true
		}

		return false
	}

	fontFlags.FixedPitch = hasFlag(FixedPitch)
	fontFlags.Serif = hasFlag(Serif)
	fontFlags.Symbolic = hasFlag(Symbolic)
	fontFlags.Script = hasFlag(Script)
	fontFlags.Nonsymbolic = hasFlag(Nonsymbolic)
	fontFlags.Italic = hasFlag(Italic)
	fontFlags.AllCap = hasFlag(AllCap)
	fontFlags.SmallCap = hasFlag(SmallCap)
	fontFlags.ForceBold = hasFlag(ForceBold)

	return fontFlags, nil
}

// FPDFFont_GetWeight returns the font weight of a font.
// Typical values are 400 (normal) and 700 (bold).
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetWeight(request *requests.FPDFFont_GetWeight) (*responses.FPDFFont_GetWeight, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	fontWeight := C.FPDFFont_GetWeight(fontHandle.handle)
	if int(fontWeight) == -1 {
		return nil, errors.New("could not get font weight")
	}

	return &responses.FPDFFont_GetWeight{
		Weight: int(fontWeight),
	}, nil
}

// FPDFFont_GetItalicAngle returns the italic angle of a font.
// The italic angle of a font is defined as degrees counterclockwise
// from vertical. For a font that slopes to the right, this will be negative.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetItalicAngle(request *requests.FPDFFont_GetItalicAngle) (*responses.FPDFFont_GetItalicAngle, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	angle := C.int(0)

	result := C.FPDFFont_GetItalicAngle(fontHandle.handle, &angle)
	if int(result) == 0 {
		return nil, errors.New("could not get italic angle")
	}

	return &responses.FPDFFont_GetItalicAngle{
		ItalicAngle: int(angle),
	}, nil
}

// FPDFFont_GetAscent returns ascent distance of a font.
// Ascent is the maximum distance in points above the baseline reached by the
// glyphs of the font. One point is 1/72 inch (around 0.3528 mm).
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetAscent(request *requests.FPDFFont_GetAscent) (*responses.FPDFFont_GetAscent, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	ascent := C.float(0)

	result := C.FPDFFont_GetAscent(fontHandle.handle, C.float(request.FontSize), &ascent)
	if int(result) == 0 {
		return nil, errors.New("could not get ascent")
	}

	return &responses.FPDFFont_GetAscent{
		Ascent: float32(ascent),
	}, nil
}

// FPDFFont_GetDescent returns the descent distance of a font.
// Descent is the maximum distance in points below the baseline reached by the
// glyphs of the font. One point is 1/72 inch (around 0.3528 mm).
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetDescent(request *requests.FPDFFont_GetDescent) (*responses.FPDFFont_GetDescent, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	descent := C.float(0)

	result := C.FPDFFont_GetDescent(fontHandle.handle, C.float(request.FontSize), &descent)
	if int(result) == 0 {
		return nil, errors.New("could not get descent")
	}

	return &responses.FPDFFont_GetDescent{
		Descent: float32(descent),
	}, nil
}

// FPDFFont_GetGlyphWidth returns the width of a glyph in a font.
// Glyph width is the distance from the end of the prior glyph to the next
// glyph. This will be the vertical distance for vertical writing.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetGlyphWidth(request *requests.FPDFFont_GetGlyphWidth) (*responses.FPDFFont_GetGlyphWidth, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	glyphWidth := C.float(0)

	result := C.FPDFFont_GetGlyphWidth(fontHandle.handle, C.uint32_t(request.Glyph), C.float(request.FontSize), &glyphWidth)
	if int(result) == 0 {
		return nil, errors.New("could not get glyph width")
	}

	return &responses.FPDFFont_GetGlyphWidth{
		GlyphWidth: float32(glyphWidth),
	}, nil
}

// FPDFFont_GetGlyphPath returns the glyphpath describing how to draw a font glyph.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetGlyphPath(request *requests.FPDFFont_GetGlyphPath) (*responses.FPDFFont_GetGlyphPath, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	glyphPath := C.FPDFFont_GetGlyphPath(fontHandle.handle, C.uint32_t(request.Glyph), C.float(request.FontSize))
	if glyphPath == nil {
		return nil, errors.New("could not get glyph path")
	}

	glyphPathHandle := p.registerGlyphPath(glyphPath)

	return &responses.FPDFFont_GetGlyphPath{
		GlyphPath: glyphPathHandle.nativeRef,
	}, nil
}

// FPDFGlyphPath_CountGlyphSegments returns the number of segments inside the given glyphpath.
// Experimental API.
func (p *PdfiumImplementation) FPDFGlyphPath_CountGlyphSegments(request *requests.FPDFGlyphPath_CountGlyphSegments) (*responses.FPDFGlyphPath_CountGlyphSegments, error) {
	p.Lock()
	defer p.Unlock()

	glyphPathHandle, err := p.getGlyphPathHandle(request.GlyphPath)
	if err != nil {
		return nil, err
	}

	count := C.FPDFGlyphPath_CountGlyphSegments(glyphPathHandle.handle)
	if int(count) == -1 {
		return nil, errors.New("could not get glyph path segment count")
	}

	return &responses.FPDFGlyphPath_CountGlyphSegments{
		Count: int(count),
	}, nil
}

// FPDFGlyphPath_GetGlyphPathSegment returns the segment in glyphpath at the given index.
// Experimental API.
func (p *PdfiumImplementation) FPDFGlyphPath_GetGlyphPathSegment(request *requests.FPDFGlyphPath_GetGlyphPathSegment) (*responses.FPDFGlyphPath_GetGlyphPathSegment, error) {
	p.Lock()
	defer p.Unlock()

	glyphPathHandle, err := p.getGlyphPathHandle(request.GlyphPath)
	if err != nil {
		return nil, err
	}

	segment := C.FPDFGlyphPath_GetGlyphPathSegment(glyphPathHandle.handle, C.int(request.Index))
	if segment == nil {
		return nil, errors.New("could not get glyph path segment")
	}

	segmentHandle := p.registerPathSegment(segment)

	return &responses.FPDFGlyphPath_GetGlyphPathSegment{
		GlyphPathSegment: segmentHandle.nativeRef,
	}, nil
}

// FPDFImageObj_GetImagePixelSize get the image size in pixels. Faster method to get only image size.
// Experimental API.
func (p *PdfiumImplementation) FPDFImageObj_GetImagePixelSize(request *requests.FPDFImageObj_GetImagePixelSize) (*responses.FPDFImageObj_GetImagePixelSize, error) {
	p.Lock()
	defer p.Unlock()

	imageObjectHandle, err := p.getPageObjectHandle(request.ImageObject)
	if err != nil {
		return nil, err
	}

	width := C.uint(0)
	height := C.uint(0)

	result := C.FPDFImageObj_GetImagePixelSize(imageObjectHandle.handle, &width, &height)
	if int(result) == 0 {
		return nil, errors.New("could not get image pixel size")
	}

	return &responses.FPDFImageObj_GetImagePixelSize{
		Width:  uint(width),
		Height: uint(height),
	}, nil
}

// FPDF_MovePages Move the given pages to a new index position.
// When this call fails, the document may be left in an indeterminate state.
// Experimental API.
func (p *PdfiumImplementation) FPDF_MovePages(request *requests.FPDF_MovePages) (*responses.FPDF_MovePages, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	if len(request.PageIndices) == 0 {
		return nil, errors.New("no page indices were given")
	}

	// Create an array that's big enough.
	valueData := make([]C.int, len(request.PageIndices))
	for i := range request.PageIndices {
		valueData[i] = C.int(request.PageIndices[i])
	}

	result := C.FPDF_MovePages(documentHandle.handle, &valueData[0], C.ulong(len(valueData)), C.int(request.DestPageIndex))
	if int(result) == 0 {
		return nil, errors.New("could not move pages")
	}

	return &responses.FPDF_MovePages{}, nil
}
