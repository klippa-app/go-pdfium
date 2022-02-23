package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_edit.h"
import "C"
import (
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
)

func (p *PdfiumImplementation) registerGlyphPath(glyphPath C.FPDF_GLYPHPATH) *GlyphPathHandle {
	ref := uuid.New()
	handle := &GlyphPathHandle{
		handle:    glyphPath,
		nativeRef: references.FPDF_GLYPHPATH(ref.String()),
	}

	p.glyphPathRefs[handle.nativeRef] = handle

	return handle
}
