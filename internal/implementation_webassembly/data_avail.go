package implementation_webassembly

import (
	"github.com/klippa-app/go-pdfium/references"

	"github.com/google/uuid"
)

type DataAvailHandle struct {
	handle                *uint64
	fileAvail             *uint64
	hints                 *uint64
	reader                *uint32
	DataAvailableCallback func(offset, size uint64) bool
	AddSegmentCallback    func(offset, size uint64)
	nativeRef             references.FPDF_AVAIL // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

func (p *PdfiumImplementation) registerDataAvail(dataAvail *uint64, fileAvail *uint64, hints *uint64, reader *uint32, dataAvailableCallback func(offset, size uint64) bool, addSegmentCallback func(offset, size uint64)) *DataAvailHandle {
	ref := uuid.New()
	handle := &DataAvailHandle{
		handle:                dataAvail,
		fileAvail:             fileAvail,
		nativeRef:             references.FPDF_AVAIL(ref.String()),
		hints:                 hints,
		reader:                reader,
		DataAvailableCallback: dataAvailableCallback,
		AddSegmentCallback:    addSegmentCallback,
	}

	p.dataAvailRefs[handle.nativeRef] = handle

	return handle
}
