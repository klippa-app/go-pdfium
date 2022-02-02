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

type FPDF_ImportPagesByIndex struct {
	Source      references.FPDF_DOCUMENT
	Destination references.FPDF_DOCUMENT
	PageIndices []int // An array of page indices to be imported. The first page is zero. If PageIndices is nil, all pages from Source are imported.
	Index       int   //  The page index at which to insert the first imported page into Destination. The first page is zero.
}

type FPDF_ImportNPagesToOne struct {
	Source          references.FPDF_DOCUMENT // The document to be imported.
	OutputWidth     float32                  // The output page width in PDF "user space" units.
	OutputHeight    float32                  // The output page height in PDF "user space" units.
	NumPagesOnXAxis int                      // The number of pages on X Axis.
	NumPagesOnYAxis int                      // The number of pages on Y Axis.
}

type FPDF_NewXObjectFromPage struct {
	Source          references.FPDF_DOCUMENT
	Destination     references.FPDF_DOCUMENT
	SourcePageIndex int
}

type FPDF_CloseXObject struct {
	XObject references.FPDF_XOBJECT
}

type FPDF_NewFormObjectFromXObject struct {
	XObject references.FPDF_XOBJECT
}
