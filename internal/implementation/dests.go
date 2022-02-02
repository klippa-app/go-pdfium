package implementation

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"
import (
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
)

func (p *PdfiumImplementation) registerDest(dest C.FPDF_DEST, documentHandle *DocumentHandle) *DestHandle {
	ref := uuid.New()
	handle := &DestHandle{
		handle:      dest,
		nativeRef:   references.FPDF_DEST(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.destRefs[handle.nativeRef] = handle
	p.destRefs[handle.nativeRef] = handle

	return handle
}
