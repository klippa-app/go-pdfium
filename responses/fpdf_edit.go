package responses

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
)

type FPDF_CreateNewDocument struct {
	Document references.FPDF_DOCUMENT
}

type FPDFPage_GetRotation struct {
	Page         int                      // The page number (0-index based).
	PageRotation enums.FPDF_PAGE_ROTATION // The page rotation.
}

type FPDFPage_SetRotation struct{}

type FPDFPage_HasTransparency struct {
	Page            int  // The page number (0-index based).
	HasTransparency bool // Whether the page has transparency.
}
