package implementation_webassembly

import (
	"github.com/google/uuid"

	"github.com/klippa-app/go-pdfium/references"
)

func (p *PdfiumImplementation) registerPageObject(pageObject *uint64) *PageObjectHandle {
	ref := uuid.New()
	handle := &PageObjectHandle{
		handle:    pageObject,
		nativeRef: references.FPDF_PAGEOBJECT(ref.String()),
	}

	p.pageObjectRefs[handle.nativeRef] = handle

	return handle
}
