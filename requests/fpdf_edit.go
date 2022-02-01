package requests

import "github.com/klippa-app/go-pdfium/references"

type FPDFPage_SetRotation struct {
	Document references.FPDF_DOCUMENT
	Page     Page
	Rotate   PageRotation // New value of PDF page rotation.
}

type FPDFPage_GetRotation struct {
	Document references.FPDF_DOCUMENT
	Page     Page
}

type FPDFPage_HasTransparency struct {
	Document references.FPDF_DOCUMENT
	Page     Page
}
