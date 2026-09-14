package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdf_dataavail.h"
import "C"
import (
	"github.com/klippa-app/go-pdfium/references"

	"github.com/google/uuid"
)

type DataAvailHandle struct {
	handle C.FPDF_AVAIL

	// fileAvailHandle is the FX_FILEAVAIL that was handed to
	// FPDFAvail_Create. PDFium keeps the pointer for the lifetime of the
	// availability provider, and it is the key of this provider's entry in
	// dataAvailAbilityCallbacks, so it has to be the same struct rather than
	// a copy of it.
	fileAvailHandle *C.FX_FILEAVAIL

	nativeRef     references.FPDF_AVAIL // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
	fileHandleRef string
	hints         *C.FX_DOWNLOADHINTS

	// openDocuments counts the documents FPDFAvail_GetDocument handed out
	// that have not been closed yet, and destroyRequested records that the
	// caller is done with the provider. PDFium requires the provider to
	// outlive those documents, so the two together decide when it can
	// really be destroyed, see destroyIfUnused.
	openDocuments    int
	destroyRequested bool
}

func (p *PdfiumImplementation) registerDataAvail(dataAvail C.FPDF_AVAIL, fileHandleRef string, fileAvailHandle *C.FX_FILEAVAIL, hints *C.FX_DOWNLOADHINTS) *DataAvailHandle {
	ref := uuid.New()
	handle := &DataAvailHandle{
		handle:          dataAvail,
		nativeRef:       references.FPDF_AVAIL(ref.String()),
		fileHandleRef:   fileHandleRef,
		fileAvailHandle: fileAvailHandle,
		hints:           hints,
	}

	p.dataAvailRefs[handle.nativeRef] = handle

	return handle
}
