//go:build pdfium_experimental
// +build pdfium_experimental

package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_edit.h"
import "C"
import (
	"errors"
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

// FPDFPageObj_GetMatrix returns the transform matrix of a page object.
// The matrix is composed as:
//   |a c e|
//   |b d f|
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
//   |a c e|
//   |b d f|
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

	return &responses.FPDFPageObj_GetMark{}, nil
}

// FPDFPageObj_AddMark adds a new content mark to the given page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_AddMark(request *requests.FPDFPageObj_AddMark) (*responses.FPDFPageObj_AddMark, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObj_AddMark{}, nil
}

// FPDFPageObj_RemoveMark removes the given content mark from the given page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_RemoveMark(request *requests.FPDFPageObj_RemoveMark) (*responses.FPDFPageObj_RemoveMark, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObj_RemoveMark{}, nil
}

// FPDFPageObjMark_GetName returns the name of a content mark.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetName(request *requests.FPDFPageObjMark_GetName) (*responses.FPDFPageObjMark_GetName, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObjMark_GetName{}, nil
}

// FPDFPageObjMark_CountParams returns the number of key/value pair parameters in the given mark.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_CountParams(request *requests.FPDFPageObjMark_CountParams) (*responses.FPDFPageObjMark_CountParams, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObjMark_CountParams{}, nil
}

// FPDFPageObjMark_GetParamKey returns the key of a property in a content mark.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetParamKey(request *requests.FPDFPageObjMark_GetParamKey) (*responses.FPDFPageObjMark_GetParamKey, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObjMark_GetParamKey{}, nil
}

// FPDFPageObjMark_GetParamValueType returns the type of the value of a property in a content mark by key.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetParamValueType(request *requests.FPDFPageObjMark_GetParamValueType) (*responses.FPDFPageObjMark_GetParamValueType, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObjMark_GetParamValueType{}, nil
}

// FPDFPageObjMark_GetParamIntValue returns the value of a number property in a content mark by key as int.
// FPDFPageObjMark_GetParamValueType() should have returned FPDF_OBJECT_NUMBER
// for this property.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetParamIntValue(request *requests.FPDFPageObjMark_GetParamIntValue) (*responses.FPDFPageObjMark_GetParamIntValue, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObjMark_GetParamIntValue{}, nil
}

// FPDFPageObjMark_GetParamStringValue returns the value of a string property in a content mark by key.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetParamStringValue(request *requests.FPDFPageObjMark_GetParamStringValue) (*responses.FPDFPageObjMark_GetParamStringValue, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObjMark_GetParamStringValue{}, nil
}

// FPDFPageObjMark_GetParamBlobValue returns the value of a blob property in a content mark by key.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetParamBlobValue(request *requests.FPDFPageObjMark_GetParamBlobValue) (*responses.FPDFPageObjMark_GetParamBlobValue, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObjMark_GetParamBlobValue{}, nil
}

// FPDFPageObjMark_SetIntParam sets the value of an int property in a content mark by key. If a parameter
// with the given key exists, its value is set to the given value. Otherwise, it is added as
// a new parameter.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_SetIntParam(request *requests.FPDFPageObjMark_SetIntParam) (*responses.FPDFPageObjMark_SetIntParam, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObjMark_SetIntParam{}, nil
}

// FPDFPageObjMark_SetStringParam sets the value of a string property in a content mark by key. If a parameter
// with the given key exists, its value is set to the given value. Otherwise, it is added as
// a new parameter.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_SetStringParam(request *requests.FPDFPageObjMark_SetStringParam) (*responses.FPDFPageObjMark_SetStringParam, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObjMark_SetStringParam{}, nil
}

// FPDFPageObjMark_SetBlobParam sets the value of a blob property in a content mark by key. If a parameter
// with the given key exists, its value is set to the given value. Otherwise, it is added as
// a new parameter.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_SetBlobParam(request *requests.FPDFPageObjMark_SetBlobParam) (*responses.FPDFPageObjMark_SetBlobParam, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObjMark_SetBlobParam{}, nil
}

// FPDFPageObjMark_RemoveParam removes a property from a content mark by key.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_RemoveParam(request *requests.FPDFPageObjMark_RemoveParam) (*responses.FPDFPageObjMark_RemoveParam, error) {
	p.Lock()
	defer p.Unlock()

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

	return &responses.FPDFImageObj_GetRenderedBitmap{}, nil
}

// FPDFPageObj_GetDashPhase returns the line dash phase of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetDashPhase(request *requests.FPDFPageObj_GetDashPhase) (*responses.FPDFPageObj_GetDashPhase, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObj_GetDashPhase{}, nil
}

// FPDFPageObj_SetDashPhase sets the line dash phase of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_SetDashPhase(request *requests.FPDFPageObj_SetDashPhase) (*responses.FPDFPageObj_SetDashPhase, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObj_SetDashPhase{}, nil
}

// FPDFPageObj_GetDashCount returns the line dash array size of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetDashCount(request *requests.FPDFPageObj_GetDashCount) (*responses.FPDFPageObj_GetDashCount, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObj_GetDashCount{}, nil
}

// FPDFPageObj_GetDashArray returns the line dash array of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetDashArray(request *requests.FPDFPageObj_GetDashArray) (*responses.FPDFPageObj_GetDashArray, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObj_GetDashArray{}, nil
}

// FPDFPageObj_SetDashArray sets the line dash array of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_SetDashArray(request *requests.FPDFPageObj_SetDashArray) (*responses.FPDFPageObj_SetDashArray, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFPageObj_SetDashArray{}, nil
}

// FPDFText_LoadStandardFont loads one of the standard 14 fonts per PDF spec 1.7 page 416. The preferred
// way of using font style is using a dash to separate the name from the style,
// for example 'Helvetica-BoldItalic'.
// The loaded font can be closed using FPDFFont_Close.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_LoadStandardFont(request *requests.FPDFText_LoadStandardFont) (*responses.FPDFText_LoadStandardFont, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFText_LoadStandardFont{}, nil
}

// FPDFTextObj_SetTextRenderMode sets the text rendering mode of a text object.
// Experimental API.
func (p *PdfiumImplementation) FPDFTextObj_SetTextRenderMode(request *requests.FPDFTextObj_SetTextRenderMode) (*responses.FPDFTextObj_SetTextRenderMode, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFTextObj_SetTextRenderMode{}, nil
}

// FPDFTextObj_GetFont returns the font of a text object.
// Experimental API.
func (p *PdfiumImplementation) FPDFTextObj_GetFont(request *requests.FPDFTextObj_GetFont) (*responses.FPDFTextObj_GetFont, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFTextObj_GetFont{}, nil
}

// FPDFFont_GetFontName returns the font name of a font.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetFontName(request *requests.FPDFFont_GetFontName) (*responses.FPDFFont_GetFontName, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFFont_GetFontName{}, nil
}

// FPDFFont_GetFlags returns the descriptor flags of a font.
// Returns the bit flags specifying various characteristics of the font as
// defined in ISO 32000-1:2008, table 123.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetFlags(request *requests.FPDFFont_GetFlags) (*responses.FPDFFont_GetFlags, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFFont_GetFlags{}, nil
}

// FPDFFont_GetWeight returns the font weight of a font.
// Typical values are 400 (normal) and 700 (bold).
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetWeight(request *requests.FPDFFont_GetWeight) (*responses.FPDFFont_GetWeight, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFFont_GetWeight{}, nil
}

// FPDFFont_GetItalicAngle returns the italic angle of a font.
// The italic angle of a font is defined as degrees counterclockwise
// from vertical. For a font that slopes to the right, this will be negative.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetItalicAngle(request *requests.FPDFFont_GetItalicAngle) (*responses.FPDFFont_GetItalicAngle, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFFont_GetItalicAngle{}, nil
}

// FPDFFont_GetAscent returns ascent distance of a font.
// Ascent is the maximum distance in points above the baseline reached by the
// glyphs of the font. One point is 1/72 inch (around 0.3528 mm).
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetAscent(request *requests.FPDFFont_GetAscent) (*responses.FPDFFont_GetAscent, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFFont_GetAscent{}, nil
}

// FPDFFont_GetDescent returns the descent distance of a font.
// Descent is the maximum distance in points below the baseline reached by the
// glyphs of the font. One point is 1/72 inch (around 0.3528 mm).
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetDescent(request *requests.FPDFFont_GetDescent) (*responses.FPDFFont_GetDescent, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFFont_GetDescent{}, nil
}

// FPDFFont_GetGlyphWidth returns the width of a glyph in a font.
// Glyph width is the distance from the end of the prior glyph to the next
// glyph. This will be the vertical distance for vertical writing.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetGlyphWidth(request *requests.FPDFFont_GetGlyphWidth) (*responses.FPDFFont_GetGlyphWidth, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFFont_GetGlyphWidth{}, nil
}

// FPDFFont_GetGlyphPath returns the glyphpath describing how to draw a font glyph.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetGlyphPath(request *requests.FPDFFont_GetGlyphPath) (*responses.FPDFFont_GetGlyphPath, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFFont_GetGlyphPath{}, nil
}

// FPDFGlyphPath_CountGlyphSegments returns the number of segments inside the given glyphpath.
// Experimental API.
func (p *PdfiumImplementation) FPDFGlyphPath_CountGlyphSegments(request *requests.FPDFGlyphPath_CountGlyphSegments) (*responses.FPDFGlyphPath_CountGlyphSegments, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFGlyphPath_CountGlyphSegments{}, nil
}

// FPDFGlyphPath_GetGlyphPathSegment returns the segment in glyphpath at the given index.
// Experimental API.
func (p *PdfiumImplementation) FPDFGlyphPath_GetGlyphPathSegment(request *requests.FPDFGlyphPath_GetGlyphPathSegment) (*responses.FPDFGlyphPath_GetGlyphPathSegment, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDFGlyphPath_GetGlyphPathSegment{}, nil
}
