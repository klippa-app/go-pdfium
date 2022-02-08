package responses

type GetPageText struct {
	Page int    // The page this text came from (0-index based).
	Text string // The plain text of a page.
}

type CharPosition struct {
	Left   float64 // The position of this char from the left.
	Top    float64 // The position of this char from the top.
	Right  float64 // The position of this char from the right.
	Bottom float64 // The position of this char from the bottom.
}

type FontInformation struct {
	Size         float64 // Font size in points (also known as em).
	SizeInPixels *int    // Font size in pixels, only available when PixelPositions is used.
	Weight       int     // The weight of the font, can be negative for spaces and newlines. Will only be filled when compiled with experimental support.
	Name         string  // The name of the font, can be empty for spaces and newlines. Will only be filled when compiled with experimental support.
	Flags        int     // Font flags, should be interpreted per PDF spec 1.7, Section 5.7.1 Font Descriptor Flags. Will only be filled when compiled with experimental support.
}

type GetPageTextStructuredChar struct {
	Text            string           // The text of this char.
	Angle           float64          // The angle this char is in.
	PointPosition   CharPosition     // The position of this char in points.
	PixelPosition   *CharPosition    // The position of this char in pixels. When PixelPositions are requested.
	FontInformation *FontInformation // The font information of this char. When CollectFontInformation is enabled.
}

type GetPageTextStructuredRect struct {
	Text            string           // The text of this rect.
	PointPosition   CharPosition     // The position of this rect in points.
	PixelPosition   *CharPosition    // The position of this rect in pixels. When PixelPositions are requested.
	FontInformation *FontInformation // The font information of this rect. When CollectFontInformation is enabled.
}

type GetPageTextStructured struct {
	Page              int                          // The page structured this text came from (0-index based).
	Chars             []*GetPageTextStructuredChar // A list of chars in a page. When Mode is GetPageTextStructuredModeChars or GetPageTextStructuredModeBoth.
	Rects             []*GetPageTextStructuredRect // A list of rects in a page. When Mode is GetPageTextStructuredModeRects or GetPageTextStructuredModeBoth.
	PointToPixelRatio float64                      // The point to pixel ratio for the calculated positions.
}
