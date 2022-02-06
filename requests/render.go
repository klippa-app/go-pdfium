package requests

import "github.com/klippa-app/go-pdfium/enums"

type RenderPageInDPI struct {
	Page        Page
	DPI         int                    // The DPI to render the page in.
	RenderFlags enums.FPDF_RENDER_FLAG // FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER will always be set to render to Go image.
}

type RenderPagesInDPI struct {
	Pages   []RenderPageInDPI // The pages
	Padding int               // The amount of padding (in pixels) between the images
}

type RenderPageInPixels struct {
	Page        Page
	Width       int                    // The maximum width of the image.
	Height      int                    // The maximum height of the image.
	RenderFlags enums.FPDF_RENDER_FLAG // FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER will always be set to render to Go image.
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
