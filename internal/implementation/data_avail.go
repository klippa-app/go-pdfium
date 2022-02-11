package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_dataavail.h"
import "C"
import (
	"github.com/klippa-app/go-pdfium/references"

	"github.com/google/uuid"
)

type DataAvailHandle struct {
	handle          C.FPDF_AVAIL
	fileAvailHandle C.FX_FILEAVAIL
	nativeRef       references.FPDF_AVAIL // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
	fileHandleRef   string
	hints           *C.FX_DOWNLOADHINTS
}

func (p *PdfiumImplementation) registerDataAvail(dataAvail C.FPDF_AVAIL, fileHandleRef string, fileAvailHandle C.FX_FILEAVAIL, hints *C.FX_DOWNLOADHINTS) *DataAvailHandle {
	ref := uuid.New()
	handle := &DataAvailHandle{
		handle:          dataAvail,
		nativeRef:       references.FPDF_AVAIL(ref.String()),
		fileHandleRef:   fileHandleRef,
		fileAvailHandle: fileAvailHandle,
	}

	p.dataAvailRefs[handle.nativeRef] = handle

	return handle
}
