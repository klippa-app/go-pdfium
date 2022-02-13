package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_structtree.h"
import "C"
import (
	"github.com/klippa-app/go-pdfium/references"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerStructTree(structTree C.FPDF_STRUCTTREE, documentHandle *DocumentHandle) *StructTreeHandle {
	ref := uuid.New()
	handle := &StructTreeHandle{
		handle:      structTree,
		nativeRef:   references.FPDF_STRUCTTREE(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.structTreeRefs[handle.nativeRef] = handle
	p.structTreeRefs[handle.nativeRef] = handle

	return handle
}

func (p *PdfiumImplementation) registerStructElement(structElement C.FPDF_STRUCTELEMENT, documentHandle *DocumentHandle) *StructElementHandle {
	ref := uuid.New()
	handle := &StructElementHandle{
		handle:      structElement,
		nativeRef:   references.FPDF_STRUCTELEMENT(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.structElementRefs[handle.nativeRef] = handle
	p.structElementRefs[handle.nativeRef] = handle

	return handle
}
