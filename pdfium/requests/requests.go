package requests

type OpenDocument struct {
	File     *[]byte // A reference to the file data.
	Password *string // The password of the document.
}

type GetPageCount struct{}

type RenderPageInDPI struct {
	Page int // The page number (0-index based).
	DPI  int // The DPI to render the page in.
}

type RenderPagesInDPI struct {
	Pages   []RenderPageInDPI // The pages
	Padding int               // The amount of padding (in pixels) between the images
}

type RenderPageInPixels struct {
	Page   int // The page number (0-index based).
	Width  int // The maximum width of the image.
	Height int // The maximum height of the image.
}

type RenderPagesInPixels struct {
	Pages   []RenderPageInPixels // The pages
	Padding int                  // The amount of padding (in pixels) between the images
}

type RenderToFileOutputFormat string // The file format to render output as.

const (
	RenderToFileOutputFormatJPG RenderToFileOutputFormat = "jpg" // Render the file as a JPEG file.
	RenderToFileOutputFormatPNG RenderToFileOutputFormat = "png" // Render the file as a PNG file.
)

type RenderToFileOutputTarget string // The file target output.

const (
	RenderToFileOutputTargetBytes RenderToFileOutputTarget = "bytes" // Returns the file as a byte array in the response.
	RenderToFileOutputTargetFile  RenderToFileOutputTarget = "file"  // Writes away the file to a given path or a generated tmp file.
)

type RenderToFile struct {
	RenderPageInDPI     *RenderPageInDPI         // To execute the RenderPageInDPI request
	RenderPagesInDPI    *RenderPagesInDPI        // To execute the RenderPagesInDPI request
	RenderPageInPixels  *RenderPageInPixels      // To execute the RenderPageInPixels request
	RenderPagesInPixels *RenderPagesInPixels     // To execute the RenderPagesInPixels request
	OutputFormat        RenderToFileOutputFormat // The format to output the image as
	OutputTarget        RenderToFileOutputTarget // Where to output the image
	MaxFileSize         int64                    // The maximum filesize, if jpg is chosen as output format, it will try to compress it until it fits
	TargetFilePath      string                   // When OutputTarget is file, the path to write it to, if not given, a temp file is created
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
	Page                   int                                 // The page number (0-index based).
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
	Calculate bool // Whether to calculate from points to pixel. Useful if you used RenderPageInDPI or RenderPageInPixels.
	DPI       int  // If rendered in a specific DPI, give the DPI. Useful if you used RenderPageInDPI.
	Width     int  // If rendered with a specific resolution, give the width resolution. Useful if you used RenderPageInPixels.
	Height    int  // If rendered with a specific resolution, give the height resolution. Useful if you used RenderPageInPixels.
}
