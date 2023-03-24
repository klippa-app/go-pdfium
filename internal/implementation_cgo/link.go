package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"
import (
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
)

func (p *PdfiumImplementation) registerLink(dest C.FPDF_LINK) *LinkHandle {
	ref := uuid.New()
	handle := &LinkHandle{
		handle:    dest,
		nativeRef: references.FPDF_LINK(ref.String()),
	}

	p.linkRefs[handle.nativeRef] = handle

	return handle
}
