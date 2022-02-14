package implementation

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"
import (
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
)

func (p *PdfiumImplementation) registerPageObject(pageObject C.FPDF_PAGEOBJECT, documentHandle *DocumentHandle) *PageObjectHandle {
	ref := uuid.New()
	handle := &PageObjectHandle{
		handle:      pageObject,
		nativeRef:   references.FPDF_PAGEOBJECT(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	p.pageObjectRefs[handle.nativeRef] = handle
	documentHandle.pageObjectRefs[handle.nativeRef] = handle

	return handle
}
