package responses

import "image"

type GetFileVersion struct {
	FileVersion int // The numeric version of the file: 14 for 1.4, 15 for 1.5, ...
}

type GetDocPermissions struct {
	DocPermissions uint32 // A 32-bit integer which indicates the permission flags. Please refer to "TABLE 3.20 User access permissions" in PDF Reference 1.7 P123 for detailed description. If the document is not protected, 0xffffffff (4294967295) will be returned.
}

type GetSecurityHandlerRevision struct {
	SecurityHandlerRevision int // The revision number of security handler. Please refer to key "R" in "TABLE 3.19 Additional encryption dictionary entries for the standard security handler" in PDF Reference 1.7 P122 for detailed description. If the document is not protected, -1 will be returned.
}

type GetPageCount struct {
	PageCount int // The amount of pages of the document.
}

type GetMetadata struct {
	Tag   string // The requested metadata tag.
	Value string // The value of the tag if found, string is empty if the value is not found.
}

type RenderPage struct {
	Page              int         // The rendered page number (0-index based).
	PointToPixelRatio float64     // The point to pixel ratio for the rendered image. How many points is 1 pixel in this image.
	Image             *image.RGBA // The rendered image.
	Width             int         // The width of the rendered image.
	Height            int         // The height of the rendered image.
}

type RenderPagesPage struct {
	Page              int     // The rendered page number (0-index based).
	PointToPixelRatio float64 // The point to pixel ratio for the rendered image. How many points is 1 pixel for this page in this image.
	Width             int     // The width of the rendered page inside the image.
	Height            int     // The height of the rendered page inside the image.
	X                 int     // The X start position of this page inside the image.
	Y                 int     // The Y start position of this page inside the image.
}

type RenderPages struct {
	Pages  []RenderPagesPage // Information about the rendered pages inside this image.
	Image  *image.RGBA       // The rendered image.
	Width  int               // The width of the rendered image.
	Height int               // The height of the rendered image.
}

type RenderToFile struct {
	Pages             []RenderPagesPage // Information about the rendered pages inside this image.
	ImageBytes        *[]byte           // The byte array of the rendered file when OutputTarget is RenderToFileOutputTargetBytes.
	ImagePath         string            // The file path when OutputTarget is RenderToFileOutputTargetFile, is a tmp path when TargetFilePath was empty in the request.
	Width             int               // The width of the rendered image.
	Height            int               // The height of the rendered image.
	PointToPixelRatio float64           // The point to pixel ratio for the rendered image. How many points is 1 pixel in this image. Only set when rendering one page.
}

type GetPageSize struct {
	Page   int     // The page this size came from (0-index based).
	Width  float64 // The width of the page in points. One point is 1/72 inch (around 0.3528 mm).
	Height float64 // The height of the page in points. One point is 1/72 inch (around 0.3528 mm).
}

type GetPageSizeInPixels struct {
	Page              int     // The page this size came from (0-index based).
	Width             int     // The width of the page in pixels.
	Height            int     // The height of the page in pixels.
	PointToPixelRatio float64 // The point to pixel ratio for the rendered image. How many points is 1 pixel in this image.
}

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
	Weight       int     // The weight of the font, can be negative for spaces and newlines.
	Name         string  // The name of the font, can be empty for spaces and newlines.
	Flags        int     // Font flags, should be interpreted per PDF spec 1.7, Section 5.7.1 Font Descriptor Flags.
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
