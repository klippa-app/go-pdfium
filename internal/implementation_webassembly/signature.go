package implementation_webassembly

import (
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
)

func (p *PdfiumImplementation) registerSignature(signature *uint64, documentHandle *DocumentHandle) *SignatureHandle {
	ref := uuid.New()
	handle := &SignatureHandle{
		handle:      signature,
		nativeRef:   references.FPDF_SIGNATURE(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.signatureRefs[handle.nativeRef] = handle
	p.signatureRefs[handle.nativeRef] = handle

	return handle
}
