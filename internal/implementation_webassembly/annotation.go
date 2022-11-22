package implementation_webassembly

import (
	"github.com/klippa-app/go-pdfium/references"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerAnnotation(annotation *uint64) *AnnotationHandle {
	ref := uuid.New()
	handle := &AnnotationHandle{
		handle:    annotation,
		nativeRef: references.FPDF_ANNOTATION(ref.String()),
	}

	p.annotationRefs[handle.nativeRef] = handle

	return handle
}
