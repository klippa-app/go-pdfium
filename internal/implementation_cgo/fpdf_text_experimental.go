//go:build pdfium_experimental
// +build pdfium_experimental

package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdf_text.h"
import "C"
import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"
)

// FPDFText_GetFontInfo returns the font name and flags of a particular character.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetFontInfo(request *requests.FPDFText_GetFontInfo) (*responses.FPDFText_GetFontInfo, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	// First get the font name length.
	fontNameSize := C.FPDFText_GetFontInfo(textPageHandle.handle, C.int(request.Index), nil, 0, nil)
	if fontNameSize == 0 {
		return nil, errors.New("could not get font name")
	}

	fontFlags := C.int(0)
	charData := make([]byte, fontNameSize)
	C.FPDFText_GetFontInfo(textPageHandle.handle, C.int(request.Index), unsafe.Pointer(&charData[0]), C.ulong(len(charData)), &fontFlags)

	return &responses.FPDFText_GetFontInfo{
		Index:    request.Index,
		FontName: string(charData[:fontNameSize-1]), // Remove NULL terminator
		Flags:    int(fontFlags),
	}, nil
}

// FPDFText_GetFontWeight returns the font weight of a particular character.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetFontWeight(request *requests.FPDFText_GetFontWeight) (*responses.FPDFText_GetFontWeight, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	fontWeight := C.FPDFText_GetFontWeight(textPageHandle.handle, C.int(request.Index))
	if int(fontWeight) == -1 {
		return nil, errors.New("could not get font weight")
	}

	return &responses.FPDFText_GetFontWeight{
		Index:      request.Index,
		FontWeight: int(fontWeight),
	}, nil
}

// FPDFText_GetTextRenderMode returns the text rendering mode of character.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetTextRenderMode(request *requests.FPDFText_GetTextRenderMode) (*responses.FPDFText_GetTextRenderMode, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	textRenderMode := C.FPDFText_GetTextRenderMode(textPageHandle.handle, C.int(request.Index))
	return &responses.FPDFText_GetTextRenderMode{
		Index:          request.Index,
		TextRenderMode: enums.FPDF_TEXT_RENDERMODE(textRenderMode),
	}, nil
}

// FPDFText_GetFillColor returns the fill color of a particular character.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetFillColor(request *requests.FPDFText_GetFillColor) (*responses.FPDFText_GetFillColor, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	r := C.uint(0)
	g := C.uint(0)
	b := C.uint(0)
	a := C.uint(0)
	success := C.FPDFText_GetFillColor(textPageHandle.handle, C.int(request.Index), &r, &g, &b, &a)
	if int(success) == 0 {
		return nil, errors.New("could not get fill color")
	}

	return &responses.FPDFText_GetFillColor{
		Index: request.Index,
		R:     uint(r),
		G:     uint(g),
		B:     uint(b),
		A:     uint(a),
	}, nil
}

// FPDFText_GetStrokeColor returns the stroke color of a particular character.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetStrokeColor(request *requests.FPDFText_GetStrokeColor) (*responses.FPDFText_GetStrokeColor, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	r := C.uint(0)
	g := C.uint(0)
	b := C.uint(0)
	a := C.uint(0)
	success := C.FPDFText_GetStrokeColor(textPageHandle.handle, C.int(request.Index), &r, &g, &b, &a)
	if int(success) == 0 {
		return nil, errors.New("could not get stroke color")
	}

	return &responses.FPDFText_GetStrokeColor{
		Index: request.Index,
		R:     uint(r),
		G:     uint(g),
		B:     uint(b),
		A:     uint(a),
	}, nil
}

// FPDFText_GetCharAngle returns the character rotation angle.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetCharAngle(request *requests.FPDFText_GetCharAngle) (*responses.FPDFText_GetCharAngle, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	charAngle := C.FPDFText_GetCharAngle(textPageHandle.handle, C.int(request.Index))
	if float64(charAngle) == -1 {
		return nil, errors.New("could not get char angle")
	}

	return &responses.FPDFText_GetCharAngle{
		Index:     request.Index,
		CharAngle: float32(charAngle),
	}, nil
}

// FPDFText_GetLooseCharBox returns a "loose" bounding box of a particular character, i.e., covering
// the entire glyph bounds, without taking the actual glyph shape into
// account. All positions are measured in PDF "user space".
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetLooseCharBox(request *requests.FPDFText_GetLooseCharBox) (*responses.FPDFText_GetLooseCharBox, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	rect := C.FS_RECTF{}
	success := C.FPDFText_GetLooseCharBox(textPageHandle.handle, C.int(request.Index), &rect)
	if int(success) == 0 {
		return nil, errors.New("could not get loose char box")
	}

	return &responses.FPDFText_GetLooseCharBox{
		Rect: structs.FPDF_FS_RECTF{
			Left:   float32(rect.left),
			Top:    float32(rect.top),
			Right:  float32(rect.right),
			Bottom: float32(rect.bottom),
		},
	}, nil
}

// FPDFText_GetMatrix returns the effective transformation matrix for a particular character.
// All positions are measured in PDF "user space".
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetMatrix(request *requests.FPDFText_GetMatrix) (*responses.FPDFText_GetMatrix, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	matrix := C.FS_MATRIX{}
	success := C.FPDFText_GetMatrix(textPageHandle.handle, C.int(request.Index), &matrix)
	if int(success) == 0 {
		return nil, errors.New("could not get char matrix")
	}

	return &responses.FPDFText_GetMatrix{
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

// FPDFLink_GetTextRange returns the start char index and char count for a link.
// Experimental API.
func (p *PdfiumImplementation) FPDFLink_GetTextRange(request *requests.FPDFLink_GetTextRange) (*responses.FPDFLink_GetTextRange, error) {
	p.Lock()
	defer p.Unlock()

	pageLinkhandle, err := p.getPageLinkHandle(request.PageLink)
	if err != nil {
		return nil, err
	}

	startCharIndex := C.int(0)
	charCount := C.int(0)

	success := C.FPDFLink_GetTextRange(pageLinkhandle.handle, C.int(request.Index), &startCharIndex, &charCount)
	if int(success) == 0 {
		return nil, errors.New("could not get text range")
	}

	return &responses.FPDFLink_GetTextRange{
		Index:          request.Index,
		StartCharIndex: int(startCharIndex),
		CharCount:      int(charCount),
	}, nil
}

// FPDFText_IsGenerated returns whether a character in a page is generated by PDFium.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_IsGenerated(request *requests.FPDFText_IsGenerated) (*responses.FPDFText_IsGenerated, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	isGenerated := C.FPDFText_IsGenerated(textPageHandle.handle, C.int(request.Index))
	if int(isGenerated) == -1 {
		return nil, errors.New("could not get whether text is generated")
	}

	return &responses.FPDFText_IsGenerated{
		Index:       request.Index,
		IsGenerated: int(isGenerated) == 1,
	}, nil
}
