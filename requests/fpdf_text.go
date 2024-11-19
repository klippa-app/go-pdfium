package requests

import "github.com/klippa-app/go-pdfium/references"

type FPDFText_LoadPage struct {
	Page Page
}

type FPDFText_ClosePage struct {
	TextPage references.FPDF_TEXTPAGE
}

type FPDFText_CountChars struct {
	TextPage references.FPDF_TEXTPAGE
}

type FPDFText_GetUnicode struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int
}

type FPDFText_GetTextObject struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int
}

type FPDFText_IsGenerated struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int
}

type FPDFText_IsHyphen struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int
}

type FPDFText_HasUnicodeMapError struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int
}

type FPDFText_GetFontSize struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int
}

type FPDFText_GetFontInfo struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int
}

type FPDFText_GetFontWeight struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int
}

type FPDFText_GetFillColor struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int
}

type FPDFText_GetStrokeColor struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int
}

type FPDFText_GetCharAngle struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int
}

type FPDFText_GetCharBox struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int
}

type FPDFText_GetLooseCharBox struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int
}

type FPDFText_GetMatrix struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int
}

type FPDFText_GetCharOrigin struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int
}

type FPDFText_GetCharIndexAtPos struct {
	TextPage   references.FPDF_TEXTPAGE
	X          float64
	Y          float64
	XTolerance float64
	YTolerance float64
}

type FPDFText_GetText struct {
	TextPage   references.FPDF_TEXTPAGE
	StartIndex int // Index for the start characters.
	Count      int // Number of characters to be extracted.
}

type FPDFText_CountRects struct {
	TextPage   references.FPDF_TEXTPAGE
	StartIndex int // Index for the start characters.
	Count      int // Number of characters to be extracted, or -1 for all remaining.
}

type FPDFText_GetRect struct {
	TextPage references.FPDF_TEXTPAGE
	Index    int // Zero-based index for the rectangle.
}

type FPDFText_GetBoundedText struct {
	TextPage references.FPDF_TEXTPAGE
	Left     float64 // Left boundary.
	Top      float64 // Top boundary.
	Right    float64 // Right boundary.
	Bottom   float64 // Bottom boundary.
}

type FPDFText_FindStartFlag uint64

const (
	FPDFText_FindStartFlag_MATCHCASE      FPDFText_FindStartFlag = 0x00000001 // If not set, it will not match case by default.
	FPDFText_FindStartFlag_MATCHWHOLEWORD FPDFText_FindStartFlag = 0x00000002 // If not set, it will not match the whole word by default.
	FPDFText_FindStartFlag_CONSECUTIVE    FPDFText_FindStartFlag = 0x00000004 // If not set, it will skip past the current match to look for the next match.
)

type FPDFText_FindStart struct {
	TextPage   references.FPDF_TEXTPAGE
	Find       string                 // A unicode match pattern.
	Flags      FPDFText_FindStartFlag // Option flags.
	StartIndex int                    // Start from this character. -1 for end of the page.
}

type FPDFText_FindNext struct {
	Search references.FPDF_SCHHANDLE
}

type FPDFText_FindPrev struct {
	Search references.FPDF_SCHHANDLE
}

type FPDFText_GetSchResultIndex struct {
	Search references.FPDF_SCHHANDLE
}

type FPDFText_GetSchCount struct {
	Search references.FPDF_SCHHANDLE
}

type FPDFText_FindClose struct {
	Search references.FPDF_SCHHANDLE
}

type FPDFLink_LoadWebLinks struct {
	TextPage references.FPDF_TEXTPAGE
}

type FPDFLink_CountWebLinks struct {
	PageLink references.FPDF_PAGELINK
}

type FPDFLink_GetURL struct {
	PageLink references.FPDF_PAGELINK
	Index    int
}

type FPDFLink_CountRects struct {
	PageLink references.FPDF_PAGELINK
	Index    int
}

type FPDFLink_GetRect struct {
	PageLink  references.FPDF_PAGELINK
	Index     int
	RectIndex int
}

type FPDFLink_GetTextRange struct {
	PageLink references.FPDF_PAGELINK
	Index    int
}

type FPDFLink_CloseWebLinks struct {
	PageLink references.FPDF_PAGELINK
}
