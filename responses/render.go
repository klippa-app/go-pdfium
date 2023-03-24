package responses

import (
	"image"
)

type RenderPage struct {
	Page              int         // The rendered page number (0-index based).
	PointToPixelRatio float64     // The point to pixel ratio for the rendered image. How many points is 1 pixel in this image.
	Image             *image.RGBA // The rendered image.
	Width             int         // The width of the rendered image.
	Height            int         // The height of the rendered image.
	HasTransparency   bool        // Whether the page has transparency.
}

type RenderPagesPage struct {
	Page              int     // The rendered page number (0-index based).
	PointToPixelRatio float64 // The point to pixel ratio for the rendered image. How many points is 1 pixel for this page in this image.
	Width             int     // The width of the rendered page inside the image.
	Height            int     // The height of the rendered page inside the image.
	X                 int     // The X start position of this page inside the image.
	Y                 int     // The Y start position of this page inside the image.
	HasTransparency   bool    // Whether the page has transparency.
}

type RenderPages struct {
	Pages  []RenderPagesPage // Information about the rendered pages inside this image.
	Image  *image.RGBA       // The rendered image.
	Width  int               // The width of the rendered image.
	Height int               // The height of the rendered image.
}

type RenderPageInPixels struct {
	Result      RenderPage
	CleanupFunc func() // In WebAssembly you MUST call Cleanup() when you are done with the image object to release resources.
}

// Cleanup should be called when using the WebAssembly runtime and when you're
// done with the Image object to release resources.
func (r *RenderPageInPixels) Cleanup() {
	if r.CleanupFunc != nil {
		r.CleanupFunc()
	}
}

type RenderPagesInPixels struct {
	Result      RenderPages
	CleanupFunc func() // In WebAssembly you MUST call Cleanup() when you are done with the image object to release resources.
}

// Cleanup should be called when using the WebAssembly runtime and when you're
// done with the Image object to release resources.
func (r *RenderPagesInPixels) Cleanup() {
	if r.CleanupFunc != nil {
		r.CleanupFunc()
	}
}

type RenderPageInDPI struct {
	Result      RenderPage
	CleanupFunc func() // In WebAssembly you MUST call Cleanup() when you are done with the image object to release resources.
}

// Cleanup should be called when using the WebAssembly runtime and when you're
// done with the Image object to release resources.
func (r *RenderPageInDPI) Cleanup() {
	if r.CleanupFunc != nil {
		r.CleanupFunc()
	}
}

type RenderPagesInDPI struct {
	Result      RenderPages
	CleanupFunc func() // In WebAssembly you MUST call Cleanup() when you are done with the image object to release resources.
}

// Cleanup should be called when using the WebAssembly runtime and when you're
// done with the Image object to release resources.
func (r *RenderPagesInDPI) Cleanup() {
	if r.CleanupFunc != nil {
		r.CleanupFunc()
	}
}

type RenderToFile struct {
	Pages             []RenderPagesPage // Information about the rendered pages inside this image.
	ImageBytes        *[]byte           // The byte array of the rendered file when OutputTarget is RenderToFileOutputTargetBytes.
	ImagePath         string            // The file path when OutputTarget is RenderToFileOutputTargetFile, is a tmp path when TargetFilePath was empty in the request.
	Width             int               // The width of the rendered image.
	Height            int               // The height of the rendered image.
	PointToPixelRatio float64           // The point to pixel ratio for the rendered image. How many points is 1 pixel in this image. Only set when rendering one page.
}
