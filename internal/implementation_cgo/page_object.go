package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"
import (
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
)

func (p *PdfiumImplementation) registerPageObject(pageObject C.FPDF_PAGEOBJECT) *PageObjectHandle {
	ref := uuid.New()
	handle := &PageObjectHandle{
		handle:    pageObject,
		nativeRef: references.FPDF_PAGEOBJECT(ref.String()),
	}

	p.pageObjectRefs[handle.nativeRef] = handle

	return handle
}
