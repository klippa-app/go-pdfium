package requests

import (
	"io"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/structs"
)

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

type FPDFPage_New struct {
	Document  references.FPDF_DOCUMENT
	PageIndex int     // Suggested 0-based index of the page to create. If it is larger than document's current last index(L), the created page index is the next available index -- L+1.
	Width     float64 // The page width in points.
	Height    float64 // The page height in points.
}

type FPDFPage_Delete struct {
	Document  references.FPDF_DOCUMENT
	PageIndex int // The index of the page to delete.
}

type FPDFPage_InsertObject struct {
	Page       Page
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPage_RemoveObject struct {
	Page       Page
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPage_CountObjects struct {
	Page Page
}

type FPDFPage_GetObject struct {
	Page  Page
	Index int
}

type FPDFPage_GenerateContent struct {
	Page Page
}

type FPDFPageObj_Destroy struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_HasTransparency struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_GetType struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_Transform struct {
	PageObject references.FPDF_PAGEOBJECT
	Transform  structs.FPDF_FS_MATRIX
}

type FPDFPageObj_GetMatrix struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_SetMatrix struct {
	PageObject references.FPDF_PAGEOBJECT
	Transform  structs.FPDF_FS_MATRIX
}

type FPDFPage_TransformAnnots struct {
	Page      Page
	Transform structs.FPDF_FS_MATRIX
}

type FPDFPageObj_NewImageObj struct {
	Document references.FPDF_DOCUMENT
}

type FPDFPageObj_CountMarks struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_GetMark struct {
	PageObject references.FPDF_PAGEOBJECT
	Index      uint64
}

type FPDFPageObj_AddMark struct {
	PageObject references.FPDF_PAGEOBJECT
	Name       string
}

type FPDFPageObj_RemoveMark struct {
	PageObject     references.FPDF_PAGEOBJECT
	PageObjectMark references.FPDF_PAGEOBJECTMARK
}

type FPDFPageObjMark_GetName struct {
	PageObjectMark references.FPDF_PAGEOBJECTMARK
}

type FPDFPageObjMark_CountParams struct {
	PageObjectMark references.FPDF_PAGEOBJECTMARK
}

type FPDFPageObjMark_GetParamKey struct {
	PageObjectMark references.FPDF_PAGEOBJECTMARK
	Index          uint64
}

type FPDFPageObjMark_GetParamValueType struct {
	PageObjectMark references.FPDF_PAGEOBJECTMARK
	Key            string
}

type FPDFPageObjMark_GetParamIntValue struct {
	PageObjectMark references.FPDF_PAGEOBJECTMARK
	Key            string
}

type FPDFPageObjMark_GetParamStringValue struct {
	PageObjectMark references.FPDF_PAGEOBJECTMARK
	Key            string
}

type FPDFPageObjMark_GetParamBlobValue struct {
	PageObjectMark references.FPDF_PAGEOBJECTMARK
	Key            string
}

type FPDFPageObjMark_SetIntParam struct {
	Document       references.FPDF_DOCUMENT
	PageObject     references.FPDF_PAGEOBJECT
	PageObjectMark references.FPDF_PAGEOBJECTMARK
	Key            string
	Value          int
}

type FPDFPageObjMark_SetStringParam struct {
	Document       references.FPDF_DOCUMENT
	PageObject     references.FPDF_PAGEOBJECT
	PageObjectMark references.FPDF_PAGEOBJECTMARK
	Key            string
	Value          string
}

type FPDFPageObjMark_SetBlobParam struct {
	Document       references.FPDF_DOCUMENT
	PageObject     references.FPDF_PAGEOBJECT
	PageObjectMark references.FPDF_PAGEOBJECTMARK
	Key            string
	Value          []byte
}

type FPDFPageObjMark_RemoveParam struct {
	PageObject     references.FPDF_PAGEOBJECT
	PageObjectMark references.FPDF_PAGEOBJECTMARK
	Key            string
}

type FPDFImageObj_LoadJpegFile struct {
	Page           *Page // The start of all loaded pages, may be nil.
	Count          int   // Number of pages, may be 0.
	FileData       []byte
	FileReader     io.ReadSeeker
	FileReaderSize int64 // Size of the file when using a reader.
	FilePath       string
}

type FPDFImageObj_LoadJpegFileInline struct {
	FileData       []byte
	FileReader     io.ReadSeeker
	FileReaderSize int64 // Size of the file when using a reader.
	FilePath       string
}

type FPDFImageObj_SetMatrix struct {
	ImageObject references.FPDF_PAGEOBJECT
	Transform   structs.FPDF_FS_MATRIX
}

type FPDFImageObj_SetBitmap struct {
	Page        *Page
	Count       int
	ImageObject references.FPDF_PAGEOBJECT
	Bitmap      references.FPDF_BITMAP
}

type FPDFImageObj_GetBitmap struct {
	ImageObject references.FPDF_PAGEOBJECT
}

type FPDFImageObj_GetRenderedBitmap struct {
	Document    references.FPDF_DOCUMENT
	Page        Page
	ImageObject references.FPDF_PAGEOBJECT
}

type FPDFImageObj_GetImageDataDecoded struct {
	ImageObject references.FPDF_PAGEOBJECT
}

type FPDFImageObj_GetImageDataRaw struct {
	ImageObject references.FPDF_PAGEOBJECT
}

type FPDFImageObj_GetImageFilterCount struct {
	ImageObject references.FPDF_PAGEOBJECT
}

type FPDFImageObj_GetImageFilter struct {
	ImageObject references.FPDF_PAGEOBJECT
	Index       int
}

type FPDFImageObj_GetImageMetadata struct {
	ImageObject references.FPDF_PAGEOBJECT
	Page        Page
}

type FPDFPageObj_CreateNewPath struct {
	X float32
	Y float32
}

type FPDFPageObj_CreateNewRect struct {
	X float32
	Y float32
	W float32
	H float32
}

type FPDFPageObj_GetBounds struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_SetBlendMode struct {
	PageObject references.FPDF_PAGEOBJECT
	BlendMode  enums.PDF_BLEND_MODE
}

type FPDFPageObj_SetStrokeColor struct {
	PageObject  references.FPDF_PAGEOBJECT
	StrokeColor structs.FPDF_COLOR
}

type FPDFPageObj_GetStrokeColor struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_SetStrokeWidth struct {
	PageObject  references.FPDF_PAGEOBJECT
	StrokeWidth float32
}

type FPDFPageObj_GetStrokeWidth struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_GetLineJoin struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_SetLineJoin struct {
	PageObject references.FPDF_PAGEOBJECT
	LineJoin   enums.FPDF_LINEJOIN
}

type FPDFPageObj_GetLineCap struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_SetLineCap struct {
	PageObject references.FPDF_PAGEOBJECT
	LineCap    enums.FPDF_LINECAP
}

type FPDFPageObj_SetFillColor struct {
	PageObject references.FPDF_PAGEOBJECT
	FillColor  structs.FPDF_COLOR
}

type FPDFPageObj_GetFillColor struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_GetDashPhase struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_SetDashPhase struct {
	PageObject references.FPDF_PAGEOBJECT
	DashPhase  float32
}

type FPDFPageObj_GetDashCount struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_GetDashArray struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_SetDashArray struct {
	PageObject references.FPDF_PAGEOBJECT
	DashArray  []float32
	DashPhase  float32
}

type FPDFPath_CountSegments struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPath_GetPathSegment struct {
	PageObject references.FPDF_PAGEOBJECT
	Index      int
}

type FPDFPathSegment_GetPoint struct {
	PathSegment references.FPDF_PATHSEGMENT
}
type FPDFPathSegment_GetType struct {
	PathSegment references.FPDF_PATHSEGMENT
}

type FPDFPathSegment_GetClose struct {
	PathSegment references.FPDF_PATHSEGMENT
}

type FPDFPath_MoveTo struct {
	PageObject references.FPDF_PAGEOBJECT
	X          float32
	Y          float32
}

type FPDFPath_LineTo struct {
	PageObject references.FPDF_PAGEOBJECT
	X          float32
	Y          float32
}

type FPDFPath_BezierTo struct {
	PageObject references.FPDF_PAGEOBJECT
	X1         float32
	Y1         float32
	X2         float32
	Y2         float32
	X3         float32
	Y3         float32
}

type FPDFPath_Close struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPath_SetDrawMode struct {
	PageObject references.FPDF_PAGEOBJECT
	FillMode   enums.FPDF_FILLMODE
	Stroke     bool
}

type FPDFPath_GetDrawMode struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFPageObj_NewTextObj struct {
	Document references.FPDF_DOCUMENT
	Font     string
	FontSize float32
}

type FPDFText_SetText struct {
	PageObject references.FPDF_PAGEOBJECT
	Text       string
}

type FPDFText_SetCharcodes struct {
	PageObject references.FPDF_PAGEOBJECT
	CharCodes  []uint32
}

type FPDFText_LoadFont struct {
	Document references.FPDF_DOCUMENT
	Data     []byte
	FontType enums.FPDF_FONT
	CID      bool // Whether the font is a CID font or not.
}

type FPDFText_LoadStandardFont struct {
	Document references.FPDF_DOCUMENT
	Font     string
}

type FPDFTextObj_GetFontSize struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFFont_Close struct {
	Font references.FPDF_FONT
}

type FPDFPageObj_CreateTextObj struct {
	Document references.FPDF_DOCUMENT
	Font     references.FPDF_FONT
	FontSize float32
}

type FPDFTextObj_GetTextRenderMode struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFTextObj_SetTextRenderMode struct {
	PageObject     references.FPDF_PAGEOBJECT
	TextRenderMode enums.FPDF_TEXT_RENDERMODE
}

type FPDFTextObj_GetText struct {
	PageObject references.FPDF_PAGEOBJECT
	TextPage   references.FPDF_TEXTPAGE
}

type FPDFTextObj_GetFont struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFFont_GetFontName struct {
	Font references.FPDF_FONT
}

type FPDFFont_GetFlags struct {
	Font references.FPDF_FONT
}

type FPDFFont_GetWeight struct {
	Font references.FPDF_FONT
}

type FPDFFont_GetItalicAngle struct {
	Font references.FPDF_FONT
}

type FPDFFont_GetAscent struct {
	Font     references.FPDF_FONT
	FontSize float32
}

type FPDFFont_GetDescent struct {
	Font     references.FPDF_FONT
	FontSize float32
}

type FPDFFont_GetGlyphWidth struct {
	Font     references.FPDF_FONT
	Glyph    uint32
	FontSize float32
}

type FPDFFont_GetGlyphPath struct {
	Font     references.FPDF_FONT
	Glyph    uint32
	FontSize float32
}

type FPDFGlyphPath_CountGlyphSegments struct {
	GlyphPath references.FPDF_GLYPHPATH
}

type FPDFGlyphPath_GetGlyphPathSegment struct {
	GlyphPath references.FPDF_GLYPHPATH
	Index     int
}

type FPDFFormObj_CountObjects struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFFormObj_GetObject struct {
	PageObject references.FPDF_PAGEOBJECT
	Index      uint64
}
