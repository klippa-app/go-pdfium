package responses

import "github.com/klippa-app/go-pdfium/references"

type FPDF_ImportPages struct{}

type FPDF_CopyViewerPreferences struct{}

type FPDF_ImportPagesByIndex struct{}

type FPDF_ImportNPagesToOne struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_NewXObjectFromPage struct {
	XObject references.FPDF_XOBJECT
}

type FPDF_CloseXObject struct{}

type FPDF_NewFormObjectFromXObject struct {
	PageObject references.FPDF_PAGEOBJECT
}
