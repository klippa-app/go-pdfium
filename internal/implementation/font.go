package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_edit.h"
import "C"
import (
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
)

func (p *PdfiumImplementation) registerFont(pageObject C.FPDF_FONT) *FontHandle {
	ref := uuid.New()
	handle := &FontHandle{
		handle:    pageObject,
		nativeRef: references.FPDF_FONT(ref.String()),
	}

	p.fontRefs[handle.nativeRef] = handle

	return handle
}
