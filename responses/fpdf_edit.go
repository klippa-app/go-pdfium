package responses

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/structs"
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

type FPDFPage_New struct {
	Page references.FPDF_PAGE
}

type FPDFPage_Delete struct{}

type FPDFPage_InsertObject struct{}

type FPDFPage_RemoveObject struct{}

type FPDFPage_CountObjects struct {
	Count int
}

type FPDFPage_GetObject struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPage_GenerateContent struct{}

type FPDFPageObj_Destroy struct{}

type FPDFPageObj_HasTransparency struct {
	HasTransparency bool
}

type FPDFPageObj_GetType struct {
	Type enums.FPDF_PAGEOBJ
}

type FPDFPageObj_Transform struct{}

type FPDFPageObj_GetMatrix struct {
	Matrix structs.FPDF_FS_MATRIX
}

type FPDFPageObj_SetMatrix struct{}

type FPDFPage_TransformAnnots struct{}

type FPDFPageObj_NewImageObj struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_CountMarks struct {
	Count int
}

type FPDFPageObj_GetMark struct {
	Mark references.FPDF_PAGEOBJECTMARK
}

type FPDFPageObj_AddMark struct {
	Mark references.FPDF_PAGEOBJECTMARK
}

type FPDFPageObj_RemoveMark struct{}

type FPDFPageObjMark_GetName struct {
	Name string
}

type FPDFPageObjMark_CountParams struct {
	Count int
}

type FPDFPageObjMark_GetParamKey struct {
	Key string
}

type FPDFPageObjMark_GetParamValueType struct {
	Type enums.FPDF_OBJECT_TYPE
}

type FPDFPageObjMark_GetParamIntValue struct {
	Value int
}

type FPDFPageObjMark_GetParamStringValue struct {
	Value string
}

type FPDFPageObjMark_GetParamBlobValue struct {
	Value []byte
}

type FPDFPageObjMark_SetIntParam struct{}

type FPDFPageObjMark_SetStringParam struct{}

type FPDFPageObjMark_SetBlobParam struct{}

type FPDFPageObjMark_RemoveParam struct{}

type FPDFImageObj_LoadJpegFile struct{}

type FPDFImageObj_LoadJpegFileInline struct{}

type FPDFImageObj_SetMatrix struct{}

type FPDFImageObj_SetBitmap struct{}

type FPDFImageObj_GetBitmap struct {
	Bitmap references.FPDF_BITMAP
}

type FPDFImageObj_GetRenderedBitmap struct {
	Bitmap references.FPDF_BITMAP
}

type FPDFImageObj_GetImageDataDecoded struct {
	Data []byte
}

type FPDFImageObj_GetImageDataRaw struct {
	Data []byte
}

type FPDFImageObj_GetImageFilterCount struct {
	Count int
}

type FPDFImageObj_GetImageFilter struct {
	ImageFilter string
}

type FPDFImageObj_GetImageMetadata struct {
	ImageMetadata structs.FPDF_IMAGEOBJ_METADATA
}

type FPDFPageObj_CreateNewPath struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_CreateNewRect struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_GetBounds struct {
	Left   float32
	Bottom float32
	Right  float32
	Top    float32
}

type FPDFPageObj_SetBlendMode struct{}

type FPDFPageObj_SetStrokeColor struct{}

type FPDFPageObj_GetStrokeColor struct {
	StrokeColor structs.FPDF_COLOR
}

type FPDFPageObj_SetStrokeWidth struct{}

type FPDFPageObj_GetStrokeWidth struct {
	StrokeWidth float32
}

type FPDFPageObj_GetLineJoin struct {
	LineJoin enums.FPDF_LINEJOIN
}

type FPDFPageObj_SetLineJoin struct {
}

type FPDFPageObj_GetLineCap struct {
	LineCap enums.FPDF_LINECAP
}
type FPDFPageObj_SetLineCap struct{}

type FPDFPageObj_SetFillColor struct{}

type FPDFPageObj_GetFillColor struct {
	FillColor structs.FPDF_COLOR
}

type FPDFPageObj_GetDashPhase struct {
	DashPhase float32
}

type FPDFPageObj_SetDashPhase struct{}

type FPDFPageObj_GetDashCount struct {
	DashCount int
}

type FPDFPageObj_GetDashArray struct {
	DashArray []float32
}

type FPDFPageObj_SetDashArray struct{}

type FPDFPath_CountSegments struct {
	Count int
}

type FPDFPath_GetPathSegment struct {
	PathSegment references.FPDF_PATHSEGMENT
}

type FPDFPathSegment_GetPoint struct {
	X float32
	Y float32
}

type FPDFPathSegment_GetType struct {
	Type enums.FPDF_SEGMENT
}

type FPDFPathSegment_GetClose struct {
	IsClose bool
}

type FPDFPath_MoveTo struct{}

type FPDFPath_LineTo struct{}

type FPDFPath_BezierTo struct{}

type FPDFPath_Close struct{}

type FPDFPath_SetDrawMode struct{}

type FPDFPath_GetDrawMode struct {
	FillMode enums.FPDF_FILLMODE
	Stroke   bool
}

type FPDFPageObj_NewTextObj struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFText_SetText struct{}

type FPDFText_SetCharcodes struct{}

type FPDFText_LoadFont struct {
	Font references.FPDF_FONT
}

type FPDFText_LoadStandardFont struct {
	Font references.FPDF_FONT
}

type FPDFTextObj_GetFontSize struct {
	FontSize float32
}

type FPDFFont_Close struct{}

type FPDFPageObj_CreateTextObj struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFTextObj_GetTextRenderMode struct {
	TextRenderMode enums.FPDF_TEXT_RENDERMODE
}
type FPDFTextObj_SetTextRenderMode struct{}

type FPDFTextObj_GetText struct {
	Text string
}

type FPDFTextObj_GetFont struct {
	Font references.FPDF_FONT
}

type FPDFFont_GetFontName struct {
	FontName string
}

type FPDFFont_GetFlags struct {
	Flags int
}

type FPDFFont_GetWeight struct {
	Weight int
}

type FPDFFont_GetItalicAngle struct {
	ItalicAngle int
}

type FPDFFont_GetAscent struct {
	Ascent float32
}

type FPDFFont_GetDescent struct {
	Descent float32
}

type FPDFFont_GetGlyphWidth struct {
	GlyphWidth float32
}

type FPDFFont_GetGlyphPath struct {
	GlyphPath references.FPDF_GLYPHPATH
}

type FPDFGlyphPath_CountGlyphSegments struct {
	Count int
}

type FPDFGlyphPath_GetGlyphPathSegment struct {
	GlyphPathSegment references.FPDF_PATHSEGMENT
}

type FPDFFormObj_CountObjects struct {
	Count int
}

type FPDFFormObj_GetObject struct {
	PageObject references.FPDF_PAGEOBJECT
}
