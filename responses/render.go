package responses

import (
	"bytes"
	"encoding/gob"
	"image"
)

func init() {
	// The multi-threaded transport (net/rpc) uses gob encoding. Interface
	// fields like RenderedImage (image.Image) require the concrete types to
	// be registered.
	gob.Register(&image.RGBA{})
	gob.Register(&image.Gray{})
}

type RenderPage struct {
	Page              int     // The rendered page number (0-index based).
	PointToPixelRatio float64 // The point to pixel ratio for the rendered image. How many points is 1 pixel in this image.

	// The rendered image. Nil when the requested ImageFormat was RenderImageFormatGrayscale.
	//
	// Deprecated: use RenderedImage instead, this field will be removed in the next major version.
	Image *image.RGBA

	RenderedImage   image.Image // The rendered image regardless of the requested ImageFormat, the concrete type is *image.RGBA or *image.Gray depending on the request. In WebAssembly mode the pixel buffer is only valid until Cleanup() is called.
	Width           int         // The width of the rendered image.
	Height          int         // The height of the rendered image.
	HasTransparency bool        // Whether the page has transparency.
}

// GobEncode makes sure that the image is not transferred twice in
// multi-threaded mode when RenderedImage contains the same image as the
// deprecated Image field.
func (r RenderPage) GobEncode() ([]byte, error) {
	type renderPageAlias RenderPage
	alias := renderPageAlias(r)
	if renderedImage, ok := r.RenderedImage.(*image.RGBA); ok && renderedImage == r.Image {
		alias.RenderedImage = nil
	}

	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(alias)
	return buf.Bytes(), err
}

// GobDecode restores the RenderedImage field that GobEncode stripped to
// prevent transferring the image twice in multi-threaded mode.
func (r *RenderPage) GobDecode(data []byte) error {
	type renderPageAlias RenderPage
	var alias renderPageAlias
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&alias); err != nil {
		return err
	}

	*r = RenderPage(alias)
	if r.RenderedImage == nil && r.Image != nil {
		r.RenderedImage = r.Image
	}

	return nil
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
	Pages []RenderPagesPage // Information about the rendered pages inside this image.

	// The rendered image. Nil when the requested ImageFormat was RenderImageFormatGrayscale.
	//
	// Deprecated: use RenderedImage instead, this field will be removed in the next major version.
	Image *image.RGBA

	RenderedImage image.Image // The rendered image regardless of the requested ImageFormat, the concrete type is *image.RGBA or *image.Gray depending on the request. In WebAssembly mode the pixel buffer is only valid until Cleanup() is called.
	Width         int         // The width of the rendered image.
	Height        int         // The height of the rendered image.
}

// GobEncode makes sure that the image is not transferred twice in
// multi-threaded mode when RenderedImage contains the same image as the
// deprecated Image field.
func (r RenderPages) GobEncode() ([]byte, error) {
	type renderPagesAlias RenderPages
	alias := renderPagesAlias(r)
	if renderedImage, ok := r.RenderedImage.(*image.RGBA); ok && renderedImage == r.Image {
		alias.RenderedImage = nil
	}

	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(alias)
	return buf.Bytes(), err
}

// GobDecode restores the RenderedImage field that GobEncode stripped to
// prevent transferring the image twice in multi-threaded mode.
func (r *RenderPages) GobDecode(data []byte) error {
	type renderPagesAlias RenderPages
	var alias renderPagesAlias
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&alias); err != nil {
		return err
	}

	*r = RenderPages(alias)
	if r.RenderedImage == nil && r.Image != nil {
		r.RenderedImage = r.Image
	}

	return nil
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
