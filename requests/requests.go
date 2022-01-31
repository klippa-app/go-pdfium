package requests

import (
	"io"

	"github.com/klippa-app/go-pdfium/references"
)

type OpenDocument struct {
	File           *[]byte // A reference to the file data.
	FilePath       *string // A path to a PDF file.
	FileReader     io.ReadSeeker
	FileReaderSize int
	Password       *string // The password of the document.
}

type CloseDocument struct {
	Document references.Document
}

type GetFileVersion struct {
	Document references.Document
}

type GetDocPermissions struct {
	Document references.Document
}

type GetSecurityHandlerRevision struct {
	Document references.Document
}

type GetPageCount struct {
	Document references.Document
}

type GetPageMode struct {
	Document references.Document
}

type GetMetadata struct {
	Document references.Document
	Tag      string // A metadata tag. Title, Author, Subject, Keywords, Creator, Producer, CreationDate, ModDate. For detailed explanation of these tags and their respective values, please refer to section 10.2.1 "Document Information Dictionary" in PDF Reference 1.7.
}

// Page can either be the index of a page or a page reference.
// When you use an index. The library will always cache the last opened page.
type Page struct {
	Index     int             // The page number (0-index based).
	Reference references.Page // A reference to a page. Received by GetPage()
}

type LoadPage struct {
	Document references.Document
	Index    int // The page number (0-index based).
}

type UnloadPage struct {
	Document references.Document
	Page     references.Page
}

type GetPageRotation struct {
	Document references.Document
	Page     Page
}

type GetPageTransparency struct {
	Document references.Document
	Page     Page
}

type FlattenPageUsage int

const (
	FlattenPageUsageNormalDisplay FlattenPageUsage = 0
	FlattenPageUsagePrint         FlattenPageUsage = 1
)

type FlattenPage struct {
	Document references.Document
	Page     Page
	Usage    FlattenPageUsage // The usage flag for the flattening.
}

type RenderPageInDPI struct {
	Document references.Document
	Page     Page
	DPI      int // The DPI to render the page in.
}

type RenderPagesInDPI struct {
	Pages   []RenderPageInDPI // The pages
	Padding int               // The amount of padding (in pixels) between the images
}

type RenderPageInPixels struct {
	Document references.Document
	Page     Page
	Width    int // The maximum width of the image.
	Height   int // The maximum height of the image.
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
	Document references.Document
	Page     Page
}

type GetPageSizeInPixels struct {
	Document references.Document
	Page     Page
	DPI      int // The DPI to calculate the size for.
}

type GetPageText struct {
	Document references.Document
	Page     Page
}

type GetPageTextStructured struct {
	Document               references.Document
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
	Document  references.Document
	Calculate bool // Whether to calculate from points to pixel. Useful if you used RenderPageInDPI or RenderPageInPixels.
	DPI       int  // If rendered in a specific DPI, give the DPI. Useful if you used RenderPageInDPI.
	Width     int  // If rendered with a specific resolution, give the width resolution. Useful if you used RenderPageInPixels.
	Height    int  // If rendered with a specific resolution, give the height resolution. Useful if you used RenderPageInPixels.
}
