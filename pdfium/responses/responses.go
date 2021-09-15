package responses

import "image"

type GetPageCount struct {
	PageCount int
}

type RenderPage struct {
	PointToPixelRatio float64
	Image             *image.RGBA
}

type GetPageSize struct {
	Width  float64
	Height float64
}

type GetPageSizeInPixels struct {
	Width             int
	Height            int
	PointToPixelRatio float64
}

type GetPageText struct {
	Text string
}

type CharPosition struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

type FontInformation struct {
	Size         float64 // Font size in points (also known as em)
	SizeInPixels *int    // Font size in pixels, only available when PixelPositions is used.
	Weight       int     // The weight of the font, can be negative for spaces and newlines.
	Name         string  // The name of the font, can be empty for spaces and newlines.
	Flags        int     // Font flags, should be interpreted per PDF spec 1.7, Section 5.7.1 Font Descriptor Flags.
}

type GetPageTextStructuredChar struct {
	Text            string
	Angle           float64
	PointPosition   CharPosition
	PixelPosition   *CharPosition
	FontInformation *FontInformation
}

type GetPageTextStructuredRect struct {
	Text            string
	PointPosition   CharPosition
	PixelPosition   *CharPosition
	FontInformation *FontInformation
}

type GetPageTextStructured struct {
	Chars             []*GetPageTextStructuredChar
	Rects             []*GetPageTextStructuredRect
	PointToPixelRatio float64
}
