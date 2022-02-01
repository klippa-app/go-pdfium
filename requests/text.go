package requests

import "github.com/klippa-app/go-pdfium/references"

type GetPageText struct {
	Document references.FPDF_DOCUMENT
	Page     Page
}

type GetPageTextStructured struct {
	Document               references.FPDF_DOCUMENT
	Page                   Page
	Mode                   GetPageTextStructuredMode           // The mode to get structured text for.
	CollectFontInformation bool                                // Whether to collect font information like name/size/weight.
	PixelPositions         GetPageTextStructuredPixelPositions // Pixel position calculation settings.
}

type GetPageTextStructuredMode string

const (
	GetPageTextStructuredModeChars GetPageTextStructuredMode = "char" // Only get every separate char
	GetPageTextStructuredModeRects GetPageTextStructuredMode = "rect" // Get char rects, strings on the same line with the same font settings.
	GetPageTextStructuredModeBoth  GetPageTextStructuredMode = "both" // Get both rects and chars.
)

type GetPageTextStructuredPixelPositions struct {
	Document  references.FPDF_DOCUMENT
	Calculate bool // Whether to calculate from points to pixel. Useful if you used RenderPageInDPI or RenderPageInPixels.
	DPI       int  // If rendered in a specific DPI, give the DPI. Useful if you used RenderPageInDPI.
	Width     int  // If rendered with a specific resolution, give the width resolution. Useful if you used RenderPageInPixels.
	Height    int  // If rendered with a specific resolution, give the height resolution. Useful if you used RenderPageInPixels.
}
