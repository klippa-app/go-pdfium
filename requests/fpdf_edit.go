package requests

import "github.com/klippa-app/go-pdfium/enums"

type FPDF_CreateNewDocument struct{}

type FPDFPage_SetRotation struct {
	Page   Page
	Rotate enums.FPDF_PAGE_ROTATION // New value of PDF page rotation.
}

type FPDFPage_GetRotation struct {
	Page Page
}

type FPDFPage_HasTransparency struct {
	Page Page
}
