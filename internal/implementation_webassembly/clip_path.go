package implementation_webassembly

import (
	"github.com/klippa-app/go-pdfium/references"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerClipPath(clipPath *uint64) *ClipPathHandle {
	ref := uuid.New()
	handle := &ClipPathHandle{
		handle:    clipPath,
		nativeRef: references.FPDF_CLIPPATH(ref.String()),
	}

	p.clipPathRefs[handle.nativeRef] = handle

	return handle
}
