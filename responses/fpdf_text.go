package responses

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/structs"
)

type FPDFText_LoadPage struct {
	TextPage references.FPDF_TEXTPAGE
}

type FPDFText_ClosePage struct{}

type FPDFText_CountChars struct {
	Count int //Number of characters in the page. Generated characters, like additional space characters, new line characters, are also counted.
}

type FPDFText_GetUnicode struct {
	Index   int
	Unicode uint //The Unicode of the particular character. If a character is not encoded in Unicode and PDFium can't convert to Unicode, the return value will be zero.
}

type FPDFText_GetFontSize struct {
	Index    int
	FontSize float64 // The font size of the particular character, measured in points (about 1/72 inch). This is the typographic size of the font (so called "em size").
}

type FPDFText_GetFontInfo struct {
	Index    int
	FontName string
	Flags    int // Font flags. These flags should be interpreted per PDF spec 1.7 Section 5.7.1 Font Descriptor Flags.
}

type FPDFText_GetFontWeight struct {
	Index      int
	FontWeight int
}

type FPDFText_GetTextRenderMode struct {
	Index          int
	TextRenderMode enums.FPDF_TEXT_RENDERMODE
}

type FPDFText_GetFillColor struct {
	Index int
	R     uint
	G     uint
	B     uint
	A     uint
}

type FPDFText_GetStrokeColor struct {
	Index int
	R     uint
	G     uint
	B     uint
	A     uint
}

type FPDFText_GetCharAngle struct {
	Index     int
	CharAngle float32
}

type FPDFText_GetCharBox struct {
	Index  int
	Left   float64
	Right  float64
	Bottom float64
	Top    float64
}

type FPDFText_GetLooseCharBox struct {
	Index int
	Rect  structs.FPDF_FS_RECTF
}

type FPDFText_GetMatrix struct {
	Index  int
	Matrix structs.FPDF_FS_MATRIX
}

type FPDFText_GetCharOrigin struct {
	Index int
	X     float64
	Y     float64
}
type FPDFText_GetCharIndexAtPos struct {
	CharIndex int // -1 when not found
}

type FPDFText_GetText struct {
	Text string
}

type FPDFText_CountRects struct {
	Count int // Number of rectangles, 0 if TextPage is null, or -1 on bad StartIndex.
}

type FPDFText_GetRect struct {
	Left   float64 // Left boundary.
	Top    float64 // Top boundary.
	Right  float64 // Right boundary.
	Bottom float64 // Bottom boundary.
}

type FPDFText_GetBoundedText struct {
	Text string
}

type FPDFText_FindStart struct {
	Search references.FPDF_SCHHANDLE
}

type FPDFText_FindNext struct {
	GotMatch bool
}

type FPDFText_FindPrev struct {
	GotMatch bool
}

type FPDFText_GetSchResultIndex struct {
	Index int
}

type FPDFText_GetSchCount struct {
	Count int
}

type FPDFText_FindClose struct{}

type FPDFLink_LoadWebLinks struct {
	PageLink references.FPDF_PAGELINK
}

type FPDFLink_CountWebLinks struct {
	Count int
}

type FPDFLink_GetURL struct {
	Index int
	URL   string
}

type FPDFLink_CountRects struct {
	Index int
	Count int // Number of rectangular areas for the link.
}

type FPDFLink_GetRect struct {
	Index     int
	RectIndex int
	Left      float64
	Top       float64
	Right     float64
	Bottom    float64
}

type FPDFLink_GetTextRange struct {
	Index          int
	StartCharIndex int
	CharCount      int
}

type FPDFLink_CloseWebLinks struct{}
