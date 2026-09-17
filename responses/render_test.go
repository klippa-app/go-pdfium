package responses

import (
	"bytes"
	"encoding/gob"
	"image"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestResponses(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Responses Suite")
}

// The gob round trip is the multi-threaded transport of a render.
var _ = Describe("RenderPage gob round trip", func() {
	// An RGBA render is only transferred once in multi-threaded mode, even
	// though it's referenced by both the deprecated Image field and the
	// RenderedImage field.
	It("transfers an RGBA render once", func() {
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
		Expect(gob.NewEncoder(&buf).Encode(in)).To(Succeed())

		var out RenderPage
		Expect(gob.NewDecoder(&buf).Decode(&out)).To(Succeed())

		Expect(out.Image).ToNot(BeNil(), "Image should not be nil after decoding")

		// Pointer equality proves the image was transferred once and relinked,
		// not transferred twice into two separate copies.
		Expect(out.RenderedImage == image.Image(out.Image)).To(BeTrue(), "RenderedImage should point to the same image as Image after decoding")

		Expect(out.Image.Pix).To(Equal(img.Pix), "image pixels should survive the round trip")
	})

	It("transfers a grayscale render", func() {
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
		Expect(gob.NewEncoder(&buf).Encode(in)).To(Succeed())

		var out RenderPage
		Expect(gob.NewDecoder(&buf).Decode(&out)).To(Succeed())

		Expect(out.Image).To(BeNil(), "Image should be nil for grayscale renders")

		outGray, ok := out.RenderedImage.(*image.Gray)
		Expect(ok).To(BeTrue(), "RenderedImage should be an *image.Gray, got %T", out.RenderedImage)

		Expect(outGray.Pix).To(Equal(imgGray.Pix), "image pixels should survive the round trip")
	})
})

// The RenderPages variant of the RenderPage RGBA round trip.
var _ = Describe("RenderPages gob round trip", func() {
	It("transfers an RGBA render once", func() {
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
		Expect(gob.NewEncoder(&buf).Encode(in)).To(Succeed())

		var out RenderPages
		Expect(gob.NewDecoder(&buf).Decode(&out)).To(Succeed())

		Expect(out.Image).ToNot(BeNil(), "Image should not be nil after decoding")

		Expect(out.RenderedImage == image.Image(out.Image)).To(BeTrue(), "RenderedImage should point to the same image as Image after decoding")

		Expect(out.Pages).To(HaveLen(1), "Pages should survive the round trip")
		Expect(out.Pages[0].Page).To(Equal(1), "Pages should survive the round trip")
	})
})
