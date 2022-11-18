package implementation_webassembly

import (
	"github.com/klippa-app/go-pdfium/references"

	"github.com/google/uuid"
)

type DataAvailHandle struct {
	handle          *int64
	fileAvailHandle *int64
	nativeRef       references.FPDF_AVAIL // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
	fileHandleRef   string
	hints           *int64
}

func (p *PdfiumImplementation) registerDataAvail(dataAvail *int64, fileHandleRef string, fileAvailHandle *int64, hints *int64) *DataAvailHandle {
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
