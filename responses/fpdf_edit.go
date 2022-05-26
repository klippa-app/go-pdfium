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
	ValueType enums.FPDF_OBJECT_TYPE
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
	FontSize float32 // The font size of the text object, measured in points (about 1/72 inch)
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

type FPDFFont_GetFontData struct {
	// The uncompressed font data. i.e. the raw font data after
	// having all stream filters applied, when the data is embedded.
	// If the font is not embedded, then this API will instead return the data for
	// the substitution font it is using.
	FontData []byte
}

type FPDFFont_GetIsEmbedded struct {
	IsEmbedded bool
}

type FPDFFont_GetFlags struct {
	Flags       uint32
	FixedPitch  bool // Whether all glyphs have the same width (as opposed to proportional or variable-pitch fonts, which have different widths).
	Serif       bool // Whether glyphs have serifs, which are short strokes drawn at an angle on the top and bottom of glyph stems. (Sans serif fonts do not have serifs.)
	Symbolic    bool // Whether the font contains glyphs outside the Adobe standard Latin character set. This flag and the Nonsymbolic flag shall not both be set or both be clear.
	Script      bool // Whether the glyphs resemble cursive handwriting.
	Nonsymbolic bool // Whether the font uses the Adobe standard Latin character set or a subset of it.
	Italic      bool // Whether the glyphs have dominant vertical strokes that are slanted.
	AllCap      bool // Whether the font contains no lowercase letters; typically used for display purposes, such as for titles or headlines.
	SmallCap    bool // Whether the font contains both uppercase and lowercase letters. The uppercase letters are similar to those in the regular version of the same typeface family. The glyphs for the lowercase letters have the same shapes as the corresponding uppercase letters, but they are sized and their proportions adjusted so that they have the same size and stroke weight as lowercase glyphs in the same typeface family.
	ForceBold   bool // Whether bold glyphs shall be painted with extra pixels even at very small text sizes by a conforming reader. If the ForceBold flag is set, features of bold glyphs may be thickened at small text sizes.
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
