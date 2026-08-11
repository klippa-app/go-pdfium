package responses

import (
	"bytes"
	"encoding/gob"
	"image"
	"testing"
)

// TestRenderPageGobRoundTripRGBA makes sure that an RGBA render is only
// transferred once in multi-threaded mode, even though it's referenced by
// both the deprecated Image field and the RenderedImage field.
func TestRenderPageGobRoundTripRGBA(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 3, 2))
	for i := range img.Pix {
		img.Pix[i] = uint8(i)
	}

	in := RenderPage{
		Page:          1,
		Image:         img,
		RenderedImage: img,
		Width:         3,
		Height:        2,
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(in); err != nil {
		t.Fatalf("encode: %v", err)
	}

	var out RenderPage
	if err := gob.NewDecoder(&buf).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if out.Image == nil {
		t.Fatal("Image should not be nil after decoding")
	}

	// Pointer equality proves the image was transferred once and relinked,
	// not transferred twice into two separate copies.
	if out.RenderedImage != image.Image(out.Image) {
		t.Fatal("RenderedImage should point to the same image as Image after decoding")
	}

	if !bytes.Equal(out.Image.Pix, img.Pix) {
		t.Fatal("image pixels should survive the round trip")
	}
}

// TestRenderPageGobRoundTripGray makes sure that a grayscale render survives
// the multi-threaded transport.
func TestRenderPageGobRoundTripGray(t *testing.T) {
	imgGray := image.NewGray(image.Rect(0, 0, 3, 2))
	for i := range imgGray.Pix {
		imgGray.Pix[i] = uint8(i)
	}

	in := RenderPage{
		Page:          1,
		RenderedImage: imgGray,
		Width:         3,
		Height:        2,
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(in); err != nil {
		t.Fatalf("encode: %v", err)
	}

	var out RenderPage
	if err := gob.NewDecoder(&buf).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if out.Image != nil {
		t.Fatal("Image should be nil for grayscale renders")
	}

	outGray, ok := out.RenderedImage.(*image.Gray)
	if !ok {
		t.Fatalf("RenderedImage should be an *image.Gray, got %T", out.RenderedImage)
	}

	if !bytes.Equal(outGray.Pix, imgGray.Pix) {
		t.Fatal("image pixels should survive the round trip")
	}
}

// TestRenderPagesGobRoundTripRGBA is the RenderPages variant of
// TestRenderPageGobRoundTripRGBA.
func TestRenderPagesGobRoundTripRGBA(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 3, 2))
	for i := range img.Pix {
		img.Pix[i] = uint8(i)
	}

	in := RenderPages{
		Pages:         []RenderPagesPage{{Page: 1, Width: 3, Height: 2}},
		Image:         img,
		RenderedImage: img,
		Width:         3,
		Height:        2,
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(in); err != nil {
		t.Fatalf("encode: %v", err)
	}

	var out RenderPages
	if err := gob.NewDecoder(&buf).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if out.Image == nil {
		t.Fatal("Image should not be nil after decoding")
	}

	if out.RenderedImage != image.Image(out.Image) {
		t.Fatal("RenderedImage should point to the same image as Image after decoding")
	}

	if len(out.Pages) != 1 || out.Pages[0].Page != 1 {
		t.Fatal("Pages should survive the round trip")
	}
}
