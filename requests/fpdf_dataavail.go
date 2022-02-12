package requests

import (
	"github.com/klippa-app/go-pdfium/references"
	"io"
)

type FPDFAvail_Create struct {
	Reader                  io.ReadSeeker
	Size                    int64
	IsDataAvailableCallback func(offset, size uint64) bool // Reports if the specified data section is currently available. A section is available if all bytes in the section are available.
	AddSegmentCallback      func(offset, size uint64)      // Add a section to be downloaded. May be nil. The offset and size of the section may not be unique. Part of the section might be already available. The download manager must deal with overlapping sections.
}

type FPDFAvail_Destroy struct {
	AvailabilityProvider references.FPDF_AVAIL
}

type FPDFAvail_IsDocAvail struct {
	AvailabilityProvider references.FPDF_AVAIL
}

type FPDFAvail_GetDocument struct {
	AvailabilityProvider references.FPDF_AVAIL
	Password             *string
}

type FPDFAvail_GetFirstPageNum struct {
	Document references.FPDF_DOCUMENT
}

type FPDFAvail_IsPageAvail struct {
	AvailabilityProvider references.FPDF_AVAIL
	PageIndex            int
}

type FPDFAvail_IsFormAvail struct {
	AvailabilityProvider references.FPDF_AVAIL
}

type FPDFAvail_IsLinearized struct {
	AvailabilityProvider references.FPDF_AVAIL
}
