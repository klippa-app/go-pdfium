package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdf_transformpage.h"
import "C"
import (
	"github.com/klippa-app/go-pdfium/references"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerClipPath(clipPath C.FPDF_CLIPPATH) *ClipPathHandle {
	ref := uuid.New()
	handle := &ClipPathHandle{
		handle:    clipPath,
		nativeRef: references.FPDF_CLIPPATH(ref.String()),
	}

	p.clipPathRefs[handle.nativeRef] = handle

	return handle
}
