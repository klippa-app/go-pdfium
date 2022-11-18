package implementation_webassembly

import (
	"github.com/klippa-app/go-pdfium/references"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerStructTree(structTree *uint64, documentHandle *DocumentHandle) *StructTreeHandle {
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

func (p *PdfiumImplementation) registerStructElement(structElement *uint64, documentHandle *DocumentHandle) *StructElementHandle {
	ref := uuid.New()
	handle := &StructElementHandle{
		handle:    structElement,
		nativeRef: references.FPDF_STRUCTELEMENT(ref.String()),
	}

	if documentHandle != nil {
		handle.documentRef = documentHandle.nativeRef
		documentHandle.structElementRefs[handle.nativeRef] = handle
	}

	p.structElementRefs[handle.nativeRef] = handle

	return handle
}

type StructElementAttributeHandle struct {
	handle    *int64
	nativeRef references.FPDF_STRUCTELEMENT_ATTR // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

func (p *PdfiumImplementation) registerStructElementAttribute(structElementAttribute *int64) *StructElementAttributeHandle {
	ref := uuid.New()
	handle := &StructElementAttributeHandle{
		handle:    structElementAttribute,
		nativeRef: references.FPDF_STRUCTELEMENT_ATTR(ref.String()),
	}

	p.structElementAttributeRefs[handle.nativeRef] = handle

	return handle
}
