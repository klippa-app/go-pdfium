package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"
import (
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
)

func (p *PdfiumImplementation) registerBStr(bstr C.FPDF_BSTR) *BStrHandle {
	ref := uuid.New()
	handle := &BStrHandle{
		handle:    bstr,
		nativeRef: references.FPDF_BSTR(ref.String()),
	}

	p.bStrRefs[handle.nativeRef] = handle

	return handle
}
