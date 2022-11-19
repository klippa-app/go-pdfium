package implementation_webassembly

import (
	"github.com/klippa-app/go-pdfium/references"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerBitmap(bitmap *uint64) *BitmapHandle {
	ref := uuid.New()
	handle := &BitmapHandle{
		handle:    bitmap,
		nativeRef: references.FPDF_BITMAP(ref.String()),
	}

	p.bitmapRefs[handle.nativeRef] = handle

	return handle
}
