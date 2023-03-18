package implementation_webassembly

import (
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
)

func (p *PdfiumImplementation) registerFormHandle(formHandle *uint64, formInfo *uint64) *FormHandleHandle {
	ref := uuid.New()
	handle := &FormHandleHandle{
		handle:           formHandle,
		nativeRef:        references.FPDF_FORMHANDLE(ref.String()),
		formInfo:         formInfo,
		pagePointers:     map[uint64]references.FPDF_PAGE{},
		documentPointers: map[uint64]references.FPDF_DOCUMENT{},
	}

	p.formHandleRefs[handle.nativeRef] = handle

	return handle
}
