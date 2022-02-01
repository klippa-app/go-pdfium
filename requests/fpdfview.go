package requests

import "github.com/klippa-app/go-pdfium/references"

type FPDF_CloseDocument struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_LoadPage struct {
	Document references.FPDF_DOCUMENT
	Index    int // The page number (0-index based).
}

type FPDF_ClosePage struct {
	Document references.FPDF_DOCUMENT
	Page     references.FPDF_PAGE
}

type FPDF_GetFileVersion struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetDocPermissions struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetSecurityHandlerRevision struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetPageCount struct {
	Document references.FPDF_DOCUMENT
}
