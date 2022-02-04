package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_attachment.h"
import "C"
import (
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
)

func (p *PdfiumImplementation) registerAttachment(attachment C.FPDF_ATTACHMENT, documentHandle *DocumentHandle) *AttachmentHandle {
	ref := uuid.New()
	handle := &AttachmentHandle{
		handle:      attachment,
		nativeRef:   references.FPDF_ATTACHMENT(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.attachmentRefs[handle.nativeRef] = handle
	p.attachmentRefs[handle.nativeRef] = handle

	return handle
}
