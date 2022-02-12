package responses

import "github.com/klippa-app/go-pdfium/references"

type FPDFPage_SetMediaBox struct{}

type FPDFPage_SetCropBox struct{}

type FPDFPage_SetBleedBox struct{}

type FPDFPage_SetTrimBox struct{}

type FPDFPage_SetArtBox struct{}

type FPDFPage_GetMediaBox struct {
	Left   float32
	Bottom float32
	Right  float32
	Top    float32
}

type FPDFPage_GetCropBox struct {
	Left   float32
	Bottom float32
	Right  float32
	Top    float32
}

type FPDFPage_GetBleedBox struct {
	Left   float32
	Bottom float32
	Right  float32
	Top    float32
}

type FPDFPage_GetTrimBox struct {
	Left   float32
	Bottom float32
	Right  float32
	Top    float32
}

type FPDFPage_GetArtBox struct {
	Left   float32
	Bottom float32
	Right  float32
	Top    float32
}

type FPDFPage_TransFormWithClip struct{}

type FPDFPageObj_TransformClipPath struct{}

type FPDFPageObj_GetClipPath struct {
	ClipPath references.FPDF_CLIPPATH
}

type FPDFClipPath_CountPaths struct {
	Count int
}

type FPDFClipPath_CountPathSegments struct {
	Count int
}

type FPDFClipPath_GetPathSegment struct {
	PathSegment references.FPDF_PATHSEGMENT
}

type FPDF_CreateClipPath struct {
	ClipPath references.FPDF_CLIPPATH
}

type FPDF_DestroyClipPath struct{}

type FPDFPage_InsertClipPath struct{}
