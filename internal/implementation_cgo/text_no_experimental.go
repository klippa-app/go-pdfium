//go:build !pdfium_experimental
// +build !pdfium_experimental

package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_text.h"
import "C"

import (
	"github.com/klippa-app/go-pdfium/responses"
)

func (p *PdfiumImplementation) getFontInformation(textPage C.FPDF_TEXTPAGE, charIndex int) *responses.FontInformation {
	fontSize := C.FPDFText_GetFontSize(textPage, C.int(charIndex))

	return &responses.FontInformation{
		Size: float64(fontSize),
	}
}
