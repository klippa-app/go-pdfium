package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_formfill.h"
import "C"
import (
	"github.com/klippa-app/go-pdfium/references"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerFormHandle(formHandle C.FPDF_FORMHANDLE) *FormHandleHandle {
	ref := uuid.New()
	handle := &FormHandleHandle{
		handle:    formHandle,
		nativeRef: references.FPDF_FORMHANDLE(ref.String()),
	}

	p.formHandleRefs[handle.nativeRef] = handle

	return handle
}
