package implementation_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/pdfium/pdfium_errors"
	"github.com/klippa-app/go-pdfium/pdfium/requests"
	"github.com/klippa-app/go-pdfium/pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Text", func() {
	pdfium := implementation.Pdfium{}

	Context("no document", func() {
		When("is opened", func() {
			Context("GetPageText()", func() {
				It("returns an error", func() {
					pageText, err := pdfium.GetPageText(&requests.GetPageText{
						Page: 0,
					})
					Expect(err).To(MatchError("no current document"))
					Expect(pageText).To(BeNil())
				})
			})

			Context("GetPageTextStructured()", func() {
				It("returns an error", func() {
					pageTextStructured, err := pdfium.GetPageTextStructured(&requests.GetPageTextStructured{
						Page: 0,
					})
					Expect(err).To(MatchError("no current document"))
					Expect(pageTextStructured).To(BeNil())
				})
			})
		})
	})

	Context("a normal PDF file", func() {
		BeforeEach(func() {
			pdfData, _ := ioutil.ReadFile("./testdata/test.pdf")
			pdfium.OpenDocument(&requests.OpenDocument{
				File: &pdfData,
			})
		})

		AfterEach(func() {
			pdfium.Close()
		})

		When("is opened", func() {
			Context("when an invalid page is given", func() {
				Context("GetPageText()", func() {
					It("returns an error", func() {
						pageText, err := pdfium.GetPageText(&requests.GetPageText{
							Page: 1,
						})
						Expect(err).To(MatchError(pdfium_errors.ErrPage))
						Expect(pageText).To(BeNil())
					})
				})

				Context("GetPageTextStructured()", func() {
					It("returns an error", func() {
						pageTextStructured, err := pdfium.GetPageTextStructured(&requests.GetPageTextStructured{
							Page: 1,
						})
						Expect(err).To(MatchError(pdfium_errors.ErrPage))
						Expect(pageTextStructured).To(BeNil())
					})
				})
			})

			Context("when the page text is requested", func() {
				It("returns the correct text", func() {
					pageText, err := pdfium.GetPageText(&requests.GetPageText{
						Page: 0,
					})
					Expect(err).To(BeNil())
					Expect(pageText).To(Equal(&responses.GetPageText{
						Text: "File: Untitled Document 2 Page 1 of 1\r\nThis is a test PDF",
					}))
				})
			})

			Context("when the structured page text is requested", func() {
				It("returns the correct structured text", func() {
					pageTextStructured, err := pdfium.GetPageTextStructured(&requests.GetPageTextStructured{
						Page: 0,
					})

					Expect(err).To(BeNil())
					Expect(pageTextStructured).To(Equal(loadStructuredText("./testdata/text_testpdf_without_pixel_calculations.json")))
				})

				Context("when PixelPositions is enabled", func() {
					Context("with no DPI and no pixels", func() {
						It("returns an error", func() {
							pageTextStructured, err := pdfium.GetPageTextStructured(&requests.GetPageTextStructured{
								Page: 0,
								PixelPositions: requests.GetPageTextStructuredPixelPositions{
									Calculate: true,
								},
							})
							Expect(err).To(MatchError("no DPI or resolution given to calculate pixel positions"))
							Expect(pageTextStructured).To(BeNil())
						})
					})

					Context("with DPI", func() {
						It("returns the correct calculations", func() {
							pageTextStructured, err := pdfium.GetPageTextStructured(&requests.GetPageTextStructured{
								Page: 0,
								PixelPositions: requests.GetPageTextStructuredPixelPositions{
									Calculate: true,
									DPI:       300,
								},
							})
							Expect(err).To(BeNil())
							Expect(pageTextStructured).To(Equal(loadStructuredText("./testdata/text_testpdf_with_dpi_pixel_calculations.json")))
						})
					})

					Context("with pixels", func() {
						It("returns the correct calculations", func() {
							pageTextStructured, err := pdfium.GetPageTextStructured(&requests.GetPageTextStructured{
								Page: 0,
								PixelPositions: requests.GetPageTextStructuredPixelPositions{
									Calculate: true,
									Width:     3000,
									Height:    3000,
								},
							})

							Expect(err).To(BeNil())
							Expect(pageTextStructured).To(Equal(loadStructuredText("./testdata/text_testpdf_with_resolution_pixel_calculations.json")))
						})
					})
				})

				Context("when PixelPositions is enabled", func() {
					It("returns the correct font information", func() {
						pageTextStructured, err := pdfium.GetPageTextStructured(&requests.GetPageTextStructured{
							Page:                   0,
							CollectFontInformation: true,
						})
						Expect(err).To(BeNil())
						Expect(pageTextStructured).To(Equal(loadStructuredText("./testdata/text_testpdf_with_font_information.json")))
					})

					Context("and PixelPositions is enabled", func() {
						It("returns the correct font information", func() {
							pageTextStructured, err := pdfium.GetPageTextStructured(&requests.GetPageTextStructured{
								Page:                   0,
								CollectFontInformation: true,
								PixelPositions: requests.GetPageTextStructuredPixelPositions{
									Calculate: true,
									Width:     3000,
									Height:    3000,
								},
							})
							Expect(err).To(BeNil())
							Expect(pageTextStructured).To(Equal(loadStructuredText("./testdata/text_testpdf_with_font_information_and_pixel_positions.json")))
						})
					})
				})
			})
		})
	})
})

func loadStructuredText(path string) *responses.GetPageTextStructured {
	preRender, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	buf := bytes.NewBuffer(preRender)
	dec := json.NewDecoder(buf)

	var text responses.GetPageTextStructured
	err = dec.Decode(&text)
	return &text
}

func writeStructuredText(path string, text responses.GetPageTextStructured) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(&text); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, buf.Bytes(), 0777); err != nil {
		return err
	}

	return nil
}
