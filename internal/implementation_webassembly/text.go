package implementation_webassembly

import (
	"bytes"
	"errors"
	"math"
	"unsafe"

	"github.com/klippa-app/go-pdfium/internal/textextract"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/google/uuid"
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

	res, err := p.call("FPDFText_LoadPage", *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	textPage := res[0]

	res, err = p.call("FPDFText_CountChars", textPage)
	if err != nil {
		return nil, err
	}

	charsInPage := *(*int32)(unsafe.Pointer(&res[0]))

	charData := make([]rune, 0, charsInPage)
	for i := 0; i < int(charsInPage); i++ {
		res, err = p.call("FPDFText_GetUnicode", textPage, uint64(i))
		if err != nil {
			return nil, err
		}

		uniChar := *(*int)(unsafe.Pointer(&res[0]))
		if uniChar != 0 {
			charData = append(charData, rune(uniChar))
		}
	}

	res, err = p.call("FPDFText_ClosePage", textPage)
	if err != nil {
		return nil, err
	}

	return &responses.GetPageText{
		Page: pageHandle.index,
		Text: string(charData),
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

	res, err := p.call("FPDFText_LoadPage", *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	textPage := res[0]

	res, err = p.call("FPDFText_CountChars", textPage)
	if err != nil {
		return nil, err
	}

	charsInPage := *(*int32)(unsafe.Pointer(&res[0]))

	collectChars := request.Mode == "" || request.Mode == requests.GetPageTextStructuredModeChars || request.Mode == requests.GetPageTextStructuredModeBoth
	collectRects := request.Mode == "" || request.Mode == requests.GetPageTextStructuredModeRects || request.Mode == requests.GetPageTextStructuredModeBoth

	leftPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer leftPointer.Free()

	topPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer topPointer.Free()

	rightPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer rightPointer.Free()

	bottomPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer bottomPointer.Free()

	// Rect text is computed on the Go side from the per-char data (see
	// internal/textextract): FPDFText_GetBoundedText re-scans every char on
	// the page per call, which makes extracting all rects O(chars × rects).
	var extractChars []textextract.Char
	if collectRects {
		extractChars = make([]textextract.Char, 0, charsInPage)
	}

	if collectChars || collectRects {
		for i := 0; i < int(charsInPage); i++ {
			_, err = p.call("FPDFText_GetCharBox", textPage, uint64(i), leftPointer.Pointer, rightPointer.Pointer, bottomPointer.Pointer, topPointer.Pointer)
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

			res, err = p.call("FPDFText_GetUnicode", textPage, uint64(i))
			if err != nil {
				return nil, err
			}
			uniChar := *(*int)(unsafe.Pointer(&res[0]))

			if collectRects {
				_, err = p.call("FPDFText_GetCharOrigin", textPage, uint64(i), leftPointer.Pointer, topPointer.Pointer)
				if err != nil {
					return nil, err
				}

				originY, err := topPointer.Value()
				if err != nil {
					return nil, err
				}

				extractChars = append(extractChars, textextract.Char{
					Left:    float32(left),
					Bottom:  float32(bottom),
					Right:   float32(right),
					Top:     float32(top),
					OriginY: float32(originY),
					Unicode: rune(uniChar),
				})
			}

			if !collectChars {
				continue
			}

			res, err = p.call("FPDFText_GetCharAngle", textPage, uint64(i))
			if err != nil {
				return nil, err
			}
			angle := *(*float32)(unsafe.Pointer(&res[0]))

			text := ""
			if uniChar != 0 {
				text = string(rune(uniChar))
			}

			char := &responses.GetPageTextStructuredChar{
				Text:  text,
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

	if collectRects {
		res, err = p.call("FPDFText_CountRects", textPage, 0, uint64(charsInPage))
		if err != nil {
			return nil, err
		}

		rectsCount := *(*int32)(unsafe.Pointer(&res[0]))

		extractor := textextract.New(extractChars)

		for i := 0; i < int(rectsCount); i++ {
			_, err = p.call("FPDFText_GetRect", textPage, uint64(i), leftPointer.Pointer, topPointer.Pointer, rightPointer.Pointer, bottomPointer.Pointer)
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

			char := &responses.GetPageTextStructuredRect{
				Text: extractor.TextInRect(float32(left), float32(top), float32(right), float32(bottom)),
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
				res, err = p.call("FPDFText_GetCharIndexAtPos", textPage, *(*uint64)(unsafe.Pointer(&char.PointPosition.Left)), *(*uint64)(unsafe.Pointer(&char.PointPosition.Top)), *(*uint64)(unsafe.Pointer(&tolerance)), *(*uint64)(unsafe.Pointer(&tolerance)))
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

	res, err = p.call("FPDFText_ClosePage", textPage)
	if err != nil {
		return nil, err
	}

	return resp, nil
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
	res, err := p.call("FPDFText_GetFontSize", textPage, *(*uint64)(unsafe.Pointer(&charIndex)))
	if err != nil {
		return nil, err
	}

	fontSize := *(*float64)(unsafe.Pointer(&res[0]))

	res, err = p.call("FPDFText_GetFontWeight", textPage, *(*uint64)(unsafe.Pointer(&charIndex)))
	if err != nil {
		return nil, err
	}

	fontWeight := *(*int32)(unsafe.Pointer(&res[0]))
	fontFlagsPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer fontFlagsPointer.Free()

	// First get the length of the font name.
	res, err = p.call("FPDFText_GetFontInfo", textPage, *(*uint64)(unsafe.Pointer(&charIndex)), 0, 0, fontFlagsPointer.Pointer)
	if err != nil {
		return nil, err
	}

	fontNameLength := *(*int32)(unsafe.Pointer(&res[0]))
	fontName := ""
	if fontNameLength > 0 {
		rawFontNamePointer, err := p.ByteArrayPointer(uint64(fontNameLength), nil)
		if err != nil {
			return nil, err
		}
		defer rawFontNamePointer.Free()

		// Get the actual font name.
		// For some reason, the font name is UTF-8.
		_, err = p.call("FPDFText_GetFontInfo", textPage, *(*uint64)(unsafe.Pointer(&charIndex)), rawFontNamePointer.Pointer, uint64(fontNameLength), fontFlagsPointer.Pointer)
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

	renderedSize := fontSize

	matrixPointer, matrixValue, err := p.CStructFS_MATRIX(nil)
	if err == nil {
		defer p.Free(matrixPointer)

		res, err = p.call("FPDFText_GetMatrix", textPage, *(*uint64)(unsafe.Pointer(&charIndex)), matrixPointer)
		if err == nil {
			success := *(*int32)(unsafe.Pointer(&res[0]))
			if int(success) != 0 {
				matrix, err := matrixValue()
				if err == nil {
					renderedSize = fontSize * math.Sqrt(float64(matrix.C)*float64(matrix.C)+float64(matrix.D)*float64(matrix.D))
				}
			}
		}
	}

	return &responses.FontInformation{
		Size:         float64(fontSize),
		RenderedSize: renderedSize,
		Weight:       int(fontWeight),
		Name:         fontName,
		Flags:        int(fontFlags),
	}, nil
}
