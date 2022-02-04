package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_javascript.h"
import "C"
import (
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
)

func (p *PdfiumImplementation) registerJavaScriptAction(javaScriptAction C.FPDF_JAVASCRIPT_ACTION, documentHandle *DocumentHandle) *JavaScriptActionHandle {
	ref := uuid.New()
	handle := &JavaScriptActionHandle{
		handle:      javaScriptAction,
		nativeRef:   references.FPDF_JAVASCRIPT_ACTION(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.javaScriptActionRefs[handle.nativeRef] = handle
	p.javaScriptActionRefs[handle.nativeRef] = handle

	return handle
}
