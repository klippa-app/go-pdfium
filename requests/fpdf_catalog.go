package requests

import "github.com/klippa-app/go-pdfium/references"

type FPDFCatalog_IsTagged struct {
	Document references.FPDF_DOCUMENT
}

type FPDFCatalog_SetLanguage struct {
	Document references.FPDF_DOCUMENT
	Language string
}
