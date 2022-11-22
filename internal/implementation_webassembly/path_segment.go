package implementation_webassembly

import (
	"github.com/klippa-app/go-pdfium/references"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerPathSegment(pathSegment *uint64) *PathSegmentHandle {
	ref := uuid.New()
	handle := &PathSegmentHandle{
		handle:    pathSegment,
		nativeRef: references.FPDF_PATHSEGMENT(ref.String()),
	}

	p.pathSegmentRefs[handle.nativeRef] = handle

	return handle
}
