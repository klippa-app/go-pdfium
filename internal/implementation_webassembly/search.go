package implementation_webassembly

import (
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
)

func (p *PdfiumImplementation) registerSearch(search *uint64, documentHandle *DocumentHandle) *SearchHandle {
	ref := uuid.New()
	handle := &SearchHandle{
		handle:      search,
		nativeRef:   references.FPDF_SCHHANDLE(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.searchRefs[handle.nativeRef] = handle
	p.searchRefs[handle.nativeRef] = handle

	return handle
}
