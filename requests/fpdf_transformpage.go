package requests

import (
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/structs"
)

type FPDFPage_SetMediaBox struct {
	Page   Page
	Left   float32
	Bottom float32
	Right  float32
	Top    float32
}

type FPDFPage_SetCropBox struct {
	Page   Page
	Left   float32
	Bottom float32
	Right  float32
	Top    float32
}

type FPDFPage_SetBleedBox struct {
	Page   Page
	Left   float32
	Bottom float32
	Right  float32
	Top    float32
}

type FPDFPage_SetTrimBox struct {
	Page   Page
	Left   float32
	Bottom float32
	Right  float32
	Top    float32
}

type FPDFPage_SetArtBox struct {
	Page   Page
	Left   float32
	Bottom float32
	Right  float32
	Top    float32
}

type FPDFPage_GetMediaBox struct {
	Page Page
}

type FPDFPage_GetCropBox struct {
	Page Page
}

type FPDFPage_GetBleedBox struct {
	Page Page
}

type FPDFPage_GetTrimBox struct {
	Page Page
}

type FPDFPage_GetArtBox struct {
	Page Page
}

type FPDFPage_TransFormWithClip struct {
	Page     Page
	Matrix   *structs.FPDF_FS_MATRIX
	ClipRect *structs.FPDF_FS_RECTF
}

type FPDFPageObj_TransformClipPath struct {
	PageObject references.FPDF_PAGEOBJECT
	A          float64
	B          float64
	C          float64
	D          float64
	E          float64
	F          float64
}

type FPDFPageObj_GetClipPath struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFClipPath_CountPaths struct {
	ClipPath references.FPDF_CLIPPATH
}

type FPDFClipPath_CountPathSegments struct {
	ClipPath  references.FPDF_CLIPPATH
	PathIndex int
}

type FPDFClipPath_GetPathSegment struct {
	ClipPath     references.FPDF_CLIPPATH
	PathIndex    int
	SegmentIndex int
}

type FPDF_CreateClipPath struct {
	Left   float32
	Bottom float32
	Right  float32
	Top    float32
}

type FPDF_DestroyClipPath struct {
	ClipPath references.FPDF_CLIPPATH
}

type FPDFPage_InsertClipPath struct {
	Page     Page
	ClipPath references.FPDF_CLIPPATH
}
