package requests

import "github.com/klippa-app/go-pdfium/references"

type FPDFText_GetCharIndexFromTextIndex struct {
	TextPage   references.FPDF_TEXTPAGE
	NTextIndex int
}

type FPDFText_GetTextIndexFromCharIndex struct {
	TextPage   references.FPDF_TEXTPAGE
	NCharIndex int
}
