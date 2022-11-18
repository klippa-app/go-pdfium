package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_text.h"
import "C"

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
	"io/ioutil"
	"math"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func (p *PdfiumImplementation) registerTextPage(attachment C.FPDF_TEXTPAGE, documentHandle *DocumentHandle) *TextPageHandle {
	ref := uuid.New()
	handle := &TextPageHandle{
		handle:      attachment,
		nativeRef:   references.FPDF_TEXTPAGE(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.textPageRefs[handle.nativeRef] = handle
	p.textPageRefs[handle.nativeRef] = handle

	return handle
}

func (p *PdfiumImplementation) registerPageLink(pageLink C.FPDF_PAGELINK, documentHandle *DocumentHandle) *PageLinkHandle {
	ref := uuid.New()
	handle := &PageLinkHandle{
		handle:      pageLink,
		nativeRef:   references.FPDF_PAGELINK(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.pageLinkRefs[handle.nativeRef] = handle
	p.pageLinkRefs[handle.nativeRef] = handle

	return handle
}

// GetPageText returns the text of a page
func (p *PdfiumImplementation) GetPageText(request *requests.GetPageText) (*responses.GetPageText, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	textPage := C.FPDFText_LoadPage(pageHandle.handle)
	charsInPage := int(C.FPDFText_CountChars(textPage))
	charData := make([]byte, (charsInPage+1)*2) // UTF16-LE max 2 bytes per char, add 1 char for terminator.
	charsWritten := C.FPDFText_GetText(textPage, C.int(0), C.int(charsInPage), (*C.ushort)(unsafe.Pointer(&charData[0])))
	C.FPDFText_ClosePage(textPage)

	transformedText, err := p.transformUTF16LEToUTF8(charData[0 : charsWritten*2])
	if err != nil {
		return nil, err
	}

	return &responses.GetPageText{
		Page: pageHandle.index,
		Text: transformedText,
	}, nil
}

// GetPageTextStructured returns the text of a page in a structured way
func (p *PdfiumImplementation) GetPageTextStructured(request *requests.GetPageTextStructured) (*responses.GetPageTextStructured, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	pointToPixelRatio := float64(0)
	if request.PixelPositions.Calculate {
		if request.PixelPositions.DPI > 0 {
			_, _, _, pointToPixelRatio, err = p.getPageSizeInPixels(request.Page, request.PixelPositions.DPI)
			if err != nil {
				return nil, err
			}
		} else if request.PixelPositions.Width == 0 && request.PixelPositions.Height == 0 {
			return nil, errors.New("no DPI or resolution given to calculate pixel positions")
		} else {
			_, _, _, ratio, err := p.calculateRenderImageSize(request.Page, request.PixelPositions.Width, request.PixelPositions.Height)
			if err != nil {
				return nil, err
			}
			pointToPixelRatio = ratio
		}
	}

	resp := &responses.GetPageTextStructured{
		Page:              pageHandle.index,
		Chars:             []*responses.GetPageTextStructuredChar{},
		Rects:             []*responses.GetPageTextStructuredRect{},
		PointToPixelRatio: pointToPixelRatio,
	}

	textPage := C.FPDFText_LoadPage(pageHandle.handle)
	charsInPage := C.FPDFText_CountChars(textPage)

	if request.Mode == "" || request.Mode == requests.GetPageTextStructuredModeChars || request.Mode == requests.GetPageTextStructuredModeBoth {
		for i := 0; i < int(charsInPage); i++ {
			angle := C.FPDFText_GetCharAngle(textPage, C.int(i))
			left := C.double(0)
			top := C.double(0)
			right := C.double(0)
			bottom := C.double(0)
			C.FPDFText_GetCharBox(textPage, C.int(i), &left, &top, &right, &bottom)
			charData := make([]byte, 4) // UTF16-LE max 2 bytes per char, so 1 byte for the char, and 1 char for terminator.
			charsWritten := C.FPDFText_GetText(textPage, C.int(i), C.int(1), (*C.ushort)(unsafe.Pointer(&charData[0])))

			transformedText, err := p.transformUTF16LEToUTF8(charData[0 : (charsWritten)*2])
			if err != nil {
				return nil, err
			}

			char := &responses.GetPageTextStructuredChar{
				Text:  transformedText,
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
			charData := make([]byte, (charsInPage+1)*2) // UTF16-LE max 2 bytes per char, add 1 char for terminator.
			left := C.double(0)
			top := C.double(0)
			right := C.double(0)
			bottom := C.double(0)

			C.FPDFText_GetRect(textPage, C.int(i), &left, &top, &right, &bottom)

			charsWritten := C.FPDFText_GetBoundedText(textPage, left, top, right, bottom, (*C.ushort)(unsafe.Pointer(&charData[0])), C.int(len(charData)))

			transformedText, err := p.transformUTF16LEToUTF8(charData[0 : charsWritten*2])
			if err != nil {
				return nil, err
			}

			char := &responses.GetPageTextStructuredRect{
				Text: transformedText,
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

	return resp, nil
}

func (p *PdfiumImplementation) transformUTF16LEToUTF8(charData []byte) (string, error) {
	pdf16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	utf16bom := unicode.BOMOverride(pdf16le.NewDecoder())
	unicodeReader := transform.NewReader(bytes.NewReader(charData), utf16bom)

	decoded, err := ioutil.ReadAll(unicodeReader)
	if err != nil {
		return "", err
	}

	// Remove NULL terminator.
	decoded = bytes.TrimSuffix(decoded, []byte("\x00"))

	return string(decoded), nil
}

func (p *PdfiumImplementation) transformUTF8ToUTF16LE(text string) ([]byte, error) {
	pdf16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	utf16bom := unicode.BOMOverride(pdf16le.NewEncoder())

	output := &bytes.Buffer{}
	unicodeWriter := transform.NewWriter(output, utf16bom)
	unicodeWriter.Write([]byte(text))
	unicodeWriter.Close()

	return output.Bytes(), nil
}

func convertPointPositions(pointPositions responses.CharPosition, ratio float64) *responses.CharPosition {
	return &responses.CharPosition{
		Left:   math.Round(pointPositions.Left * ratio),
		Top:    math.Round(pointPositions.Top * ratio),
		Right:  math.Round(pointPositions.Right * ratio),
		Bottom: math.Round(pointPositions.Bottom * ratio),
	}
}
