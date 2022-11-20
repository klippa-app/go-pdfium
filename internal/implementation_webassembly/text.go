package implementation_webassembly

import (
	"bytes"
	"errors"
	"io/ioutil"
	"math"
	"unsafe"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/google/uuid"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func (p *PdfiumImplementation) registerTextPage(attachment *uint64, documentHandle *DocumentHandle) *TextPageHandle {
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

func (p *PdfiumImplementation) registerPageLink(pageLink *uint64, documentHandle *DocumentHandle) *PageLinkHandle {
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

	res, err := p.Module.ExportedFunction("FPDFText_LoadPage").Call(p.Context, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	textPage := res[0]

	res, err = p.Module.ExportedFunction("FPDFText_CountChars").Call(p.Context, textPage)
	if err != nil {
		return nil, err
	}

	charsInPage := *(*int32)(unsafe.Pointer(&res[0]))

	charDataPointer, err := p.ByteArrayPointer(uint64((charsInPage + 1) * 2)) // UTF16-LE max 2 bytes per char, add 1 char for terminator.
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDFText_GetText").Call(p.Context, textPage, 0, uint64(charsInPage), charDataPointer.Pointer)
	if err != nil {
		return nil, err
	}

	charsWritten := *(*int32)(unsafe.Pointer(&res[0]))

	res, err = p.Module.ExportedFunction("FPDFText_ClosePage").Call(p.Context, textPage)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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

	res, err := p.Module.ExportedFunction("FPDFText_LoadPage").Call(p.Context, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	textPage := res[0]

	res, err = p.Module.ExportedFunction("FPDFText_CountChars").Call(p.Context, textPage)
	if err != nil {
		return nil, err
	}

	charsInPage := *(*int32)(unsafe.Pointer(&res[0]))

	if request.Mode == "" || request.Mode == requests.GetPageTextStructuredModeChars || request.Mode == requests.GetPageTextStructuredModeBoth {
		for i := 0; i < int(charsInPage); i++ {
			res, err = p.Module.ExportedFunction("FPDFText_GetCharAngle").Call(p.Context, textPage, uint64(i))
			if err != nil {
				return nil, err
			}
			angle := *(*float32)(unsafe.Pointer(&res[0]))

			leftPointer, err := p.DoublePointer(nil)
			if err != nil {
				return nil, err
			}

			topPointer, err := p.DoublePointer(nil)
			if err != nil {
				return nil, err
			}

			rightPointer, err := p.DoublePointer(nil)
			if err != nil {
				return nil, err
			}

			bottomPointer, err := p.DoublePointer(nil)
			if err != nil {
				return nil, err
			}

			_, err = p.Module.ExportedFunction("FPDFText_GetCharBox").Call(p.Context, textPage, uint64(i), leftPointer.Pointer, topPointer.Pointer, rightPointer.Pointer, bottomPointer.Pointer)
			if err != nil {
				return nil, err
			}

			charDataPointer, err := p.ByteArrayPointer(4) // UTF16-LE max 2 bytes per char, so 1 byte for the char, and 1 char for terminator.
			if err != nil {
				return nil, err
			}
			defer charDataPointer.Free()

			res, err = p.Module.ExportedFunction("FPDFText_GetText").Call(p.Context, textPage, uint64(i), 1, charDataPointer.Pointer)
			if err != nil {
				return nil, err
			}

			charsWritten := *(*int32)(unsafe.Pointer(&res[0]))

			charData, err := charDataPointer.Value(false)
			if err != nil {
				return nil, err
			}

			transformedText, err := p.transformUTF16LEToUTF8(charData[0 : (charsWritten)*2])
			if err != nil {
				return nil, err
			}

			left, err := leftPointer.Value()
			if err != nil {
				return nil, err
			}

			top, err := topPointer.Value()
			if err != nil {
				return nil, err
			}

			right, err := rightPointer.Value()
			if err != nil {
				return nil, err
			}

			bottom, err := bottomPointer.Value()
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
				fontInfo, err := p.getFontInformation(textPage, i)
				if err != nil {
					return nil, err
				}
				char.FontInformation = fontInfo
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
		res, err = p.Module.ExportedFunction("FPDFText_CountRects").Call(p.Context, textPage, 0, uint64(charsInPage))
		if err != nil {
			return nil, err
		}

		rectsCount := *(*int32)(unsafe.Pointer(&res[0]))

		// Create a buffer that has room for all chars in this page, since
		// we don't know the amount of chars in the section.
		// We need to clear this every time, because we don't know how much bytes every char is.
		charDataPointer, err := p.ByteArrayPointer(uint64((charsInPage + 1) * 2)) // UTF16-LE max 2 bytes per char, add 1 char for terminator.
		if err != nil {
			return nil, err
		}

		leftPointer, err := p.DoublePointer(nil)
		if err != nil {
			return nil, err
		}

		topPointer, err := p.DoublePointer(nil)
		if err != nil {
			return nil, err
		}

		rightPointer, err := p.DoublePointer(nil)
		if err != nil {
			return nil, err
		}

		bottomPointer, err := p.DoublePointer(nil)
		if err != nil {
			return nil, err
		}

		for i := 0; i < int(rectsCount); i++ {
			_, err = p.Module.ExportedFunction("FPDFText_GetRect").Call(p.Context, textPage, uint64(i), leftPointer.Pointer, topPointer.Pointer, rightPointer.Pointer, bottomPointer.Pointer)
			if err != nil {
				return nil, err
			}

			left, err := leftPointer.Value()
			if err != nil {
				return nil, err
			}

			top, err := topPointer.Value()
			if err != nil {
				return nil, err
			}

			right, err := rightPointer.Value()
			if err != nil {
				return nil, err
			}

			bottom, err := bottomPointer.Value()
			if err != nil {
				return nil, err
			}

			res, err = p.Module.ExportedFunction("FPDFText_GetBoundedText").Call(p.Context, textPage, *(*uint64)(unsafe.Pointer(&left)), *(*uint64)(unsafe.Pointer(&top)), *(*uint64)(unsafe.Pointer(&right)), *(*uint64)(unsafe.Pointer(&bottom)), charDataPointer.Pointer, uint64((charsInPage+1)*2))
			if err != nil {
				return nil, err
			}

			charsWritten := *(*int32)(unsafe.Pointer(&res[0]))

			charData, err := charDataPointer.Value(false)
			if err != nil {
				return nil, err
			}

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
				tolerance := float64(5)
				res, err = p.Module.ExportedFunction("FPDFText_GetCharIndexAtPos").Call(p.Context, textPage, *(*uint64)(unsafe.Pointer(&char.PointPosition.Left)), *(*uint64)(unsafe.Pointer(&char.PointPosition.Top)), *(*uint64)(unsafe.Pointer(&tolerance)), *(*uint64)(unsafe.Pointer(&tolerance)))
				if err != nil {
					return nil, err
				}

				charIndex := *(*int32)(unsafe.Pointer(&res[0]))
				fontInfo, err := p.getFontInformation(textPage, int(charIndex))
				if err != nil {
					return nil, err
				}

				char.FontInformation = fontInfo
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

	res, err = p.Module.ExportedFunction("FPDFText_ClosePage").Call(p.Context, textPage)
	if err != nil {
		return nil, err
	}

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

func (p *PdfiumImplementation) getFontInformation(textPage uint64, charIndex int) (*responses.FontInformation, error) {
	res, err := p.Module.ExportedFunction("FPDFText_GetFontSize").Call(p.Context, textPage, *(*uint64)(unsafe.Pointer(&charIndex)))
	if err != nil {
		return nil, err
	}

	fontSize := *(*float64)(unsafe.Pointer(&res[0]))

	res, err = p.Module.ExportedFunction("FPDFText_GetFontSize").Call(p.Context, textPage, *(*uint64)(unsafe.Pointer(&charIndex)))
	if err != nil {
		return nil, err
	}

	fontWeight := *(*int32)(unsafe.Pointer(&res[0]))
	fontFlagsPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}

	// First get the length of the font name.
	res, err = p.Module.ExportedFunction("FPDFText_GetFontInfo").Call(p.Context, textPage, *(*uint64)(unsafe.Pointer(&charIndex)), 0, 0, fontFlagsPointer.Pointer)
	if err != nil {
		return nil, err
	}

	fontNameLength := *(*int32)(unsafe.Pointer(&res[0]))
	fontName := ""
	if fontNameLength > 0 {
		rawFontNamePointer, err := p.ByteArrayPointer(uint64(fontNameLength))
		if err != nil {
			return nil, err
		}

		// Get the actual font name.
		// For some reason, the font name is UTF-8.
		_, err = p.Module.ExportedFunction("FPDFText_GetFontInfo").Call(p.Context, textPage, *(*uint64)(unsafe.Pointer(&charIndex)), rawFontNamePointer.Pointer, uint64(fontNameLength), fontFlagsPointer.Pointer)
		if err != nil {
			return nil, err
		}

		rawFontName, err := rawFontNamePointer.Value(false)
		if err != nil {
			return nil, err
		}

		// Convert byte array to string, remove trailing null.
		fontName = string(bytes.TrimSuffix(rawFontName, []byte("\x00")))
	}

	fontFlags, err := fontFlagsPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FontInformation{
		Size:   float64(fontSize),
		Weight: int(fontWeight),
		Name:   fontName,
		Flags:  int(fontFlags),
	}, nil
}
