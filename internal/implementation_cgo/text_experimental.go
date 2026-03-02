//go:build pdfium_experimental
// +build pdfium_experimental

package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_text.h"
import "C"

import (
	"bytes"
	"math"
	"unsafe"

	"github.com/klippa-app/go-pdfium/responses"
)

func (p *PdfiumImplementation) getFontInformation(textPage C.FPDF_TEXTPAGE, charIndex int) *responses.FontInformation {
	fontSize := C.FPDFText_GetFontSize(textPage, C.int(charIndex))
	fontWeight := C.FPDFText_GetFontWeight(textPage, C.int(charIndex))
	fontFlags := C.int(0)

	// First get the length of the font name.
	fontNameLength := C.FPDFText_GetFontInfo(textPage, C.int(charIndex), C.NULL, 0, &fontFlags)

	fontName := ""
	if fontNameLength > 0 {
		rawFontName := make([]byte, fontNameLength)

		// Get the actual font name.
		// For some reason, the font name is UTF-8.
		C.FPDFText_GetFontInfo(textPage, C.int(charIndex), unsafe.Pointer(&rawFontName[0]), C.ulong(len(rawFontName)), &fontFlags)

		// Convert byte array to string, remove trailing null.
		fontName = string(bytes.TrimSuffix(rawFontName, []byte("\x00")))
	}

	renderedSize := float64(fontSize)

	matrix := C.FS_MATRIX{}
	success := C.FPDFText_GetMatrix(textPage, C.int(charIndex), &matrix)
	if int(success) != 0 {
		renderedSize = float64(fontSize) * math.Sqrt(float64(matrix.c)*float64(matrix.c)+float64(matrix.d)*float64(matrix.d))
	}

	return &responses.FontInformation{
		Size:         float64(fontSize),
		RenderedSize: renderedSize,
		Weight:       int(fontWeight),
		Name:         fontName,
		Flags:        int(fontFlags),
	}
}
