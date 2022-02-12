package responses

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
)

type FPDFAvail_Create struct {
	AvailabilityProvider references.FPDF_AVAIL
}

type FPDFAvail_Destroy struct{}

type FPDFAvail_IsDocAvail struct {
	IsDocAvail enums.PDF_FILEAVAIL_DATA
}

type FPDFAvail_GetDocument struct {
	Document references.FPDF_DOCUMENT
}

type FPDFAvail_GetFirstPageNum struct {
	FirstPageNum int
}

type FPDFAvail_IsPageAvail struct {
	IsPageAvail enums.PDF_FILEAVAIL_DATA
}

type FPDFAvail_IsFormAvail struct {
	IsFormAvail enums.PDF_FILEAVAIL_FORM
}

type FPDFAvail_IsLinearized struct {
	IsLinearized enums.PDF_FILEAVAIL_LINEARIZATION
}
