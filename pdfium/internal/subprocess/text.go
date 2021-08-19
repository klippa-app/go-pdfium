package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_text.h"
import "C"

import (
	"bytes"
	"unsafe"

	"github.com/klippa-app/go-pdfium/pdfium/internal/commons"
)

// GetPageText returns the text of a page
func (d *Document) GetPageText(page int) string {
	d.loadPage(page)

	mutex.Lock()
	textPage := C.FPDFText_LoadPage(d.page)
	charsInPage := int(C.FPDFText_CountChars(textPage))
	charData := make([]byte, (charsInPage+1)*4) // UTF-8 = Max 4 bytes per char, add 1 for terminator.
	charsWritten := C.FPDFText_GetText(textPage, C.int(0), C.int(charsInPage), (*C.ushort)(unsafe.Pointer(&charData[0])))
	C.FPDFText_ClosePage(textPage)
	mutex.Unlock()

	return string(bytes.ReplaceAll(charData[0:charsWritten*4], []byte("\x00"), []byte{}))
}

// GetPageText returns the text of a page in a structured way
func (d *Document) GetPageTextStructured(page int, mode commons.GetPageTextStructuredRequestMode) *commons.GetPageTextStructuredResponse {
	d.loadPage(page)

	resp := &commons.GetPageTextStructuredResponse{
		Chars: []*commons.GetPageTextStructuredResponseChar{},
		Rects: []*commons.GetPageTextStructuredResponseRect{},
	}

	mutex.Lock()
	textPage := C.FPDFText_LoadPage(d.page)
	charsInPage := C.FPDFText_CountChars(textPage)

	if mode == "" || mode == commons.GetPageTextStructuredRequestModeChars || mode == commons.GetPageTextStructuredRequestModeBoth {
		for i := 0; i < int(charsInPage); i++ {
			angle := C.FPDFText_GetCharAngle(textPage, C.int(i))
			left := C.double(0)
			top := C.double(0)
			right := C.double(0)
			bottom := C.double(0)
			C.FPDFText_GetCharBox(textPage, C.int(i), &left, &top, &right, &bottom)
			charData := make([]byte, 8) // UTF-8 = Max 4 bytes per char, room for 2 chars.
			charsWritten := C.FPDFText_GetText(textPage, C.int(i), C.int(1), (*C.ushort)(unsafe.Pointer(&charData[0])))
			resp.Chars = append(resp.Chars, &commons.GetPageTextStructuredResponseChar{
				Text:   string(bytes.ReplaceAll(charData[0:charsWritten*4], []byte("\x00"), []byte{})),
				Left:   float64(left),
				Top:    float64(top),
				Right:  float64(right),
				Bottom: float64(bottom),
				Angle:  float64(angle),
			})
		}
	}

	if mode == "" || mode == commons.GetPageTextStructuredRequestModeRects || mode == commons.GetPageTextStructuredRequestModeBoth {
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
			resp.Rects = append(resp.Rects, &commons.GetPageTextStructuredResponseRect{
				Text:   string(bytes.ReplaceAll(charData[0:charsWritten*4], []byte("\x00"), []byte{})),
				Left:   float64(left),
				Top:    float64(top),
				Right:  float64(right),
				Bottom: float64(bottom),
			})
		}
	}

	C.FPDFText_ClosePage(textPage)
	mutex.Unlock()

	return resp
}
