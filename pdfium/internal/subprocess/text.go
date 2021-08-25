package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_text.h"
import "C"

import (
	"bytes"
	"errors"
	"math"
	"unsafe"

	"github.com/klippa-app/go-pdfium/pdfium/requests"
	"github.com/klippa-app/go-pdfium/pdfium/responses"
)

// GetPageText returns the text of a page
func (p *Pdfium) GetPageText(request *requests.GetPageText) (*responses.GetPageText, error) {
	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	p.Lock()
	textPage := C.FPDFText_LoadPage(p.currentPage)
	charsInPage := int(C.FPDFText_CountChars(textPage))
	charData := make([]byte, (charsInPage+1)*4) // UTF-8 = Max 4 bytes per char, add 1 for terminator.
	charsWritten := C.FPDFText_GetText(textPage, C.int(0), C.int(charsInPage), (*C.ushort)(unsafe.Pointer(&charData[0])))
	C.FPDFText_ClosePage(textPage)
	p.Unlock()

	return &responses.GetPageText{
		Text: string(bytes.ReplaceAll(charData[0:charsWritten*4], []byte("\x00"), []byte{})),
	}, nil
}

// GetPageTextStructured returns the text of a page in a structured way
func (p *Pdfium) GetPageTextStructured(request *requests.GetPageTextStructured) (*responses.GetPageTextStructured, error) {
	if p.currentDoc == nil {
		return nil, errors.New("no current document")
	}

	err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	pointToPixelRatio := float64(0)
	if request.PixelPositions.Calculate {
		if request.PixelPositions.DPI > 0 {
			pagePixelSize, err := p.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
				Page: request.Page,
				DPI:  request.PixelPositions.DPI,
			})
			if err != nil {
				return nil, err
			}
			pointToPixelRatio = pagePixelSize.PointToPixelRatio
		} else if request.PixelPositions.Width == 0 && request.PixelPositions.Height == 0 {
			return nil, errors.New("no DPI or resolution given to calculate pixel positions")
		} else {
			_, _, ratio, err := p.calculateRenderImageSize(request.Page, request.PixelPositions.Width, request.PixelPositions.Height)
			if err != nil {
				return nil, err
			}
			pointToPixelRatio = ratio
		}
	}

	resp := &responses.GetPageTextStructured{
		Chars:             []*responses.GetPageTextStructuredChar{},
		Rects:             []*responses.GetPageTextStructuredRect{},
		PointToPixelRatio: pointToPixelRatio,
	}

	p.Lock()
	textPage := C.FPDFText_LoadPage(p.currentPage)
	charsInPage := C.FPDFText_CountChars(textPage)

	if request.Mode == "" || request.Mode == requests.GetPageTextStructuredModeChars || request.Mode == requests.GetPageTextStructuredModeBoth {
		for i := 0; i < int(charsInPage); i++ {
			angle := C.FPDFText_GetCharAngle(textPage, C.int(i))
			left := C.double(0)
			top := C.double(0)
			right := C.double(0)
			bottom := C.double(0)
			C.FPDFText_GetCharBox(textPage, C.int(i), &left, &top, &right, &bottom)
			charData := make([]byte, 8) // UTF-8 = Max 4 bytes per char, room for 2 chars.
			charsWritten := C.FPDFText_GetText(textPage, C.int(i), C.int(1), (*C.ushort)(unsafe.Pointer(&charData[0])))
			pointPosition := responses.CharPosition{
				Left:   float64(left),
				Top:    float64(top),
				Right:  float64(right),
				Bottom: float64(bottom),
			}
			resp.Chars = append(resp.Chars, &responses.GetPageTextStructuredChar{
				Text:          string(bytes.ReplaceAll(charData[0:charsWritten*4], []byte("\x00"), []byte{})),
				Angle:         float64(angle),
				PointPosition: pointPosition,
				PixelPosition: convertPointPositions(request.PixelPositions, pointPosition, pointToPixelRatio),
			})
		}
	}

	if request.Mode == "" || request.Mode == requests.GetPageTextStructuredModeRects || request.Mode == requests.GetPageTextStructuredModeBoth {
		rectsCount := C.FPDFText_CountRects(textPage, C.int(0), C.int(charsInPage))
		for i := 0; i < int(rectsCount); i++ {
			// Create a buffer that has room for all chars in this page, since
			// we don't know the amount of chars in the section.
			// We need to clear this every time, because we don't know how much bytes every char is.
			charData := make([]byte, (charsInPage+1)*4) // UTF-8 = Max 4 bytes per char, add 1 for terminator.
			left := C.double(0)
			top := C.double(0)
			right := C.double(0)
			bottom := C.double(0)

			C.FPDFText_GetRect(textPage, C.int(i), &left, &top, &right, &bottom)

			charsWritten := C.FPDFText_GetBoundedText(textPage, left, top, right, bottom, (*C.ushort)(unsafe.Pointer(&charData[0])), C.int(len(charData)))
			pointPosition := responses.CharPosition{
				Left:   float64(left),
				Top:    float64(top),
				Right:  float64(right),
				Bottom: float64(bottom),
			}
			resp.Rects = append(resp.Rects, &responses.GetPageTextStructuredRect{
				Text:          string(bytes.ReplaceAll(charData[0:charsWritten*4], []byte("\x00"), []byte{})),
				PointPosition: pointPosition,
				PixelPosition: convertPointPositions(request.PixelPositions, pointPosition, pointToPixelRatio),
			})
		}
	}

	C.FPDFText_ClosePage(textPage)
	p.Unlock()

	return resp, nil
}

func convertPointPositions(pixelPositions requests.GetPageTextStructuredPixelPositions, pointPositions responses.CharPosition, ratio float64) *responses.CharPosition {
	if !pixelPositions.Calculate {
		return nil
	}

	return &responses.CharPosition{
		Left:   math.Round(pointPositions.Left * ratio),
		Top:    math.Round(pointPositions.Top * ratio),
		Right:  math.Round(pointPositions.Right * ratio),
		Bottom: math.Round(pointPositions.Bottom * ratio),
	}
}
