//go:build pdfium_use_turbojpeg

package image_jpeg

import (
	"bytes"
	"image"
	"image/jpeg"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Encode with turbojpeg", func() {
	var img *image.RGBA

	BeforeEach(func() {
		img = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{100, 100}})
	})

	It("encodes with the default options", func() {
		testWriter := bytes.NewBuffer(nil)
		err := Encode(testWriter, img, Options{})
		Expect(err).To(BeNil())
		Expect(testWriter.Len()).To(Equal(823))
	})

	It("encodes with a quality", func() {
		testWriter := bytes.NewBuffer(nil)
		err := Encode(testWriter, img, Options{
			Options: &jpeg.Options{
				Quality: 100,
			},
		})
		Expect(err).To(BeNil())
		Expect(testWriter.Len()).To(Equal(825))
	})

	It("encodes progressive", func() {
		testWriter := bytes.NewBuffer(nil)
		err := Encode(testWriter, img, Options{
			Options: &jpeg.Options{
				Quality: 100,
			},
			Progressive: true,
		})
		Expect(err).To(BeNil())
		Expect(testWriter.Len()).To(Equal(592))
	})
})
