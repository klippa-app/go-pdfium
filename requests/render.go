package requests

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
)

// RenderPageCrop defines a region of a page to render, in points.
// One point is 1/72 inch (around 0.3528 mm).
//
// The region is in the coordinate space of the rendered page, which means that
// the origin (0,0) is the top-left corner of the page and that Y grows
// downwards, the same as the resulting image. The page rotation and the crop
// box have already been applied to that space, so the region is relative to the
// page size that GetPageSize returns.
//
// The region is allowed to extend outside of the page. The area that falls
// outside of the page is filled with the background color, which makes it
// possible to cut a page into equally sized tiles without having to handle the
// tiles at the edges of the page differently.
//
// Cropping is only supported when rendering a single page.
type RenderPageCrop struct {
	X      float64 // The offset of the region in points from the left of the page.
	Y      float64 // The offset of the region in points from the top of the page.
	Width  float64 // The width of the region in points, must be larger than 0.
	Height float64 // The height of the region in points, must be larger than 0.
}

type RenderPageInDPI struct {
	Page        Page
	DPI         int                       // The DPI to render the page in.
	RenderFlags enums.FPDF_RENDER_FLAG    // FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER will always be set to render to Go image.
	RenderForm  bool                      // Whether to render form elements.
	Document    *references.FPDF_DOCUMENT // The document to render if not passed through the page by index, required when RenderForm is true.
	Crop        *RenderPageCrop           // When given, only this region of the page is rendered. Only supported when rendering a single page.
}

type RenderPagesInDPI struct {
	Pages   []RenderPageInDPI // The pages
	Padding int               // The amount of padding (in pixels) between the images
}

type RenderPageInPixels struct {
	Page        Page
	Width       int                       // The maximum width of the image.
	Height      int                       // The maximum height of the image.
	RenderFlags enums.FPDF_RENDER_FLAG    // FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER will always be set to render to Go image.
	RenderForm  bool                      // Whether to render form elements.
	Document    *references.FPDF_DOCUMENT // The document to render if not passed through the page by index, required when RenderForm is true.
	Crop        *RenderPageCrop           // When given, only this region of the page is rendered, and Width and Height apply to the region instead of to the full page. Only supported when rendering a single page.
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
	OutputQuality       int                      // Only used when OutputFormat RenderToFileOutputFormatJPG. Ranges from 1 to 100 inclusive, higher is better. The default is 95.
	Progressive         bool                     // Only used when OutputFormat RenderToFileOutputFormatJPG and with build tag pdfium_use_turbojpeg. Will render a progressive jpeg.
	MaxFileSize         int64                    // The maximum file size, when OutputFormat RenderToFileOutputFormatJPG, it will try to lower the quality it until it fits.
	TargetFilePath      string                   // When OutputTarget is file, the path to write it to, if not given, a temp file is created
}
