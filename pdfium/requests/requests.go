package requests

type OpenDocument struct {
	File *[]byte // A reference to the file data.
}

type GetPageCount struct {}

type RenderPageInDPI struct {
	Page int // The page number (0-index based).
	DPI  int // The DPI to render the page in.
}

type RenderPageInPixels struct {
	Page   int // The page number (0-index based).
	Width  int // The maximum width of the image.
	Height int // The maximum height of the image.
}

type GetPageSize struct {
	Page int // The page number (0-index based).
}

type GetPageSizeInPixels struct {
	Page int // The page number (0-index based).
	DPI  int // The DPI to calculate the size for.
}

type GetPageText struct {
	Page int // The page number (0-index based).
}

type GetPageTextStructured struct {
	Page           int                                 // The page number (0-index based).
	Mode           GetPageTextStructuredMode           // The mode to get structured text for.
	PixelPositions GetPageTextStructuredPixelPositions // Pixel position calculation settings.
}

type GetPageTextStructuredMode string

const (
	GetPageTextStructuredModeChars GetPageTextStructuredMode = "char" // Only get every separate char
	GetPageTextStructuredModeRects GetPageTextStructuredMode = "rect" // Get char rects, strings on the same line with the same font settings.
	GetPageTextStructuredModeBoth  GetPageTextStructuredMode = "both" // Get both rects and chars.
)

type GetPageTextStructuredPixelPositions struct {
	// Whether to calculate from points to pixel.
	// Useful if you used RenderPageInDPI or RenderPageInPixels.
	Calculate bool

	// If rendered in a specific DPI, give the DPI.
	// Useful if you used RenderPageInDPI.
	DPI int

	// If rendered with a specific resolution, give the resolution.
	// Useful if you used RenderPageInPixels.
	Width  int
	Height int
}
