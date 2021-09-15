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
		Text: string(removeNullTerminator(charData[0 : charsWritten*4])),
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

			char := &responses.GetPageTextStructuredChar{
				Text:  string(removeNullTerminator(charData[0 : charsWritten*4])),
				Angle: float64(angle),
				PointPosition: responses.CharPosition{
					Left:   float64(left),
					Top:    float64(top),
					Right:  float64(right),
					Bottom: float64(bottom),
				},
			}

			if request.CollectFontInformation {
				char.FontInformation = p.getFontInformation(textPage, i)
			}

			if request.PixelPositions.Calculate {
				char.PixelPosition = convertPointPositions(char.PointPosition, pointToPixelRatio)

				if char.FontInformation != nil {
					sizeInPixels := int(math.Round(char.FontInformation.Size * pointToPixelRatio))
					char.FontInformation.SizeInPixels = &sizeInPixels
				}
			}

			resp.Chars = append(resp.Chars, char)
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
			char := &responses.GetPageTextStructuredRect{
				Text: string(removeNullTerminator(charData[0 : charsWritten*4])),
				PointPosition: responses.CharPosition{
					Left:   float64(left),
					Top:    float64(top),
					Right:  float64(right),
					Bottom: float64(bottom),
				},
			}

			if request.CollectFontInformation {
				// Find index of the first letter of the rect.
				// @todo: is 5 a "valid" tolerance?
				tolerance := C.double(5)
				charIndex := C.FPDFText_GetCharIndexAtPos(textPage, C.double(char.PointPosition.Left), C.double(char.PointPosition.Top), tolerance, tolerance)
				char.FontInformation = p.getFontInformation(textPage, int(charIndex))
			}

			if request.PixelPositions.Calculate {
				char.PixelPosition = convertPointPositions(char.PointPosition, pointToPixelRatio)
				if char.FontInformation != nil {
					sizeInPixels := int(math.Round(char.FontInformation.Size * pointToPixelRatio))
					char.FontInformation.SizeInPixels = &sizeInPixels
				}
			}

			resp.Rects = append(resp.Rects, char)
		}
	}

	C.FPDFText_ClosePage(textPage)
	p.Unlock()

	return resp, nil
}

func (p *Pdfium) getFontInformation(textPage C.FPDF_TEXTPAGE, charIndex int) *responses.FontInformation {
	fontSize := C.FPDFText_GetFontSize(textPage, C.int(charIndex))
	fontWeight := C.FPDFText_GetFontWeight(textPage, C.int(charIndex))
	fontName := make([]byte, 255)
	fontFlags := C.int(0)
	fontNameLength := C.FPDFText_GetFontInfo(textPage, C.int(charIndex), unsafe.Pointer(&fontName[0]), C.ulong(len(fontName)), &fontFlags)

	return &responses.FontInformation{
		Size:   float64(fontSize),
		Weight: int(fontWeight),
		Name:   string(removeNullTerminator(fontName[:fontNameLength])),
		Flags:  int(fontFlags),
	}
}

func convertPointPositions(pointPositions responses.CharPosition, ratio float64) *responses.CharPosition {
	return &responses.CharPosition{
		Left:   math.Round(pointPositions.Left * ratio),
		Top:    math.Round(pointPositions.Top * ratio),
		Right:  math.Round(pointPositions.Right * ratio),
		Bottom: math.Round(pointPositions.Bottom * ratio),
	}
}

func removeNullTerminator(input []byte) []byte {
	return bytes.ReplaceAll(input, []byte("\x00"), []byte{})
}
