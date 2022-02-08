//go:build !pdfium_experimental
// +build !pdfium_experimental

package implementation

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_text.h"
import "C"

import (
	"github.com/klippa-app/go-pdfium/responses"
)

func (p *PdfiumImplementation) getFontInformation(textPage C.FPDF_TEXTPAGE, charIndex int) *responses.FontInformation {
	fontSize := C.FPDFText_GetFontSize(textPage, C.int(charIndex))
	fontWeight := C.FPDFText_GetFontWeight(textPage, C.int(charIndex))
	fontFlags := C.int(0)

	return &responses.FontInformation{
		Size:   float64(fontSize),
		Weight: int(fontWeight),
		Flags:  int(fontFlags),
	}
}
