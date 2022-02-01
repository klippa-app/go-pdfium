package requests

import "github.com/klippa-app/go-pdfium/references"

type FPDF_ImportPages struct {
	Source      references.FPDF_DOCUMENT
	Destination references.FPDF_DOCUMENT
	PageRange   *string // The page ranges, such as "1,3,5-7". If it is nil, it means to import all pages from parameter Source to Destination.
	Index       int     // An integer value which specifies the page index in parameter Destination where the imported pages will be inserted.
}

type FPDF_CopyViewerPreferences struct {
	Source      references.FPDF_DOCUMENT
	Destination references.FPDF_DOCUMENT
}
