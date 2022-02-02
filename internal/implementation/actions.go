package implementation

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"
import (
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
)

func (p *PdfiumImplementation) registerAction(dest C.FPDF_ACTION, documentHandle *DocumentHandle) *ActionHandle {
	ref := uuid.New()
	handle := &ActionHandle{
		handle:      dest,
		nativeRef:   references.FPDF_ACTION(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.actionRefs[handle.nativeRef] = handle
	p.actionRefs[handle.nativeRef] = handle

	return handle
}
