package shared_tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func RunTextTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("Render", func() {
		Context("a normal PDF file", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				pdfData, err := ioutil.ReadFile(testsPath + "/testdata/test.pdf")
				Expect(err).To(BeNil())
				if err != nil {
					return
				}

				newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData)
				if err != nil {
					return
				}

				doc = *newDoc
			})

			AfterEach(func() {
				err := pdfiumContainer.FPDF_CloseDocument(doc)
				Expect(err).To(BeNil())
			})

			When("is opened", func() {
				Context("when an invalid page is given", func() {
					Context("GetPageText()", func() {
						It("returns an error", func() {
							pageText, err := pdfiumContainer.GetPageText(&requests.GetPageText{
								Document: doc,
								Page: requests.Page{
									Index: 1,
								},
							})
							Expect(err).To(MatchError(errors.ErrPage.Error()))
							Expect(pageText).To(BeNil())
						})
					})

					Context("GetPageTextStructured()", func() {
						It("returns an error", func() {
							pageTextStructured, err := pdfiumContainer.GetPageTextStructured(&requests.GetPageTextStructured{
								Document: doc,
								Page: requests.Page{
									Index: 1,
								},
							})
							Expect(err).To(MatchError(errors.ErrPage.Error()))
							Expect(pageTextStructured).To(BeNil())
						})
					})
				})

				Context("when the page text is requested", func() {
					It("returns the correct text", func() {
						pageText, err := pdfiumContainer.GetPageText(&requests.GetPageText{
							Document: doc,
							Page: requests.Page{
								Index: 0,
							},
						})
						Expect(err).To(BeNil())
						Expect(pageText).To(Equal(&responses.GetPageText{
							Text: "File: Untitled Document 2 Page 1 of 1\r\nThis is a test PDF",
						}))
					})
				})

				Context("when the structured page text is requested", func() {
					It("returns the correct structured text", func() {
						pageTextStructured, err := pdfiumContainer.GetPageTextStructured(&requests.GetPageTextStructured{
							Document: doc,
							Page: requests.Page{
								Index: 0,
							},
						})

						Expect(err).To(BeNil())
						Expect(pageTextStructured).To(Equal(loadStructuredText(testsPath+"/testdata/text_"+prefix+"_testpdf_without_pixel_calculations.json", pageTextStructured)))
					})

					Context("when PixelPositions is enabled", func() {
						Context("with no DPI and no pixels", func() {
							It("returns an error", func() {
								pageTextStructured, err := pdfiumContainer.GetPageTextStructured(&requests.GetPageTextStructured{
									Document: doc,
									Page: requests.Page{
										Index: 0,
									},
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
								pageTextStructured, err := pdfiumContainer.GetPageTextStructured(&requests.GetPageTextStructured{
									Document: doc,
									Page: requests.Page{
										Index: 0,
									},
									PixelPositions: requests.GetPageTextStructuredPixelPositions{
										Calculate: true,
										DPI:       300,
									},
								})
								Expect(err).To(BeNil())
								Expect(pageTextStructured).To(Equal(loadStructuredText(testsPath+"/testdata/text_"+prefix+"_testpdf_with_dpi_pixel_calculations.json", pageTextStructured)))
							})
						})

						Context("with pixels", func() {
							It("returns the correct calculations", func() {
								pageTextStructured, err := pdfiumContainer.GetPageTextStructured(&requests.GetPageTextStructured{
									Document: doc,
									Page: requests.Page{
										Index: 0,
									},
									PixelPositions: requests.GetPageTextStructuredPixelPositions{
										Calculate: true,
										Width:     3000,
										Height:    3000,
									},
								})

								Expect(err).To(BeNil())
								Expect(pageTextStructured).To(Equal(loadStructuredText(testsPath+"/testdata/text_"+prefix+"_testpdf_with_resolution_pixel_calculations.json", pageTextStructured)))
							})
						})
					})

					Context("when PixelPositions is enabled", func() {
						It("returns the correct font information", func() {
							pageTextStructured, err := pdfiumContainer.GetPageTextStructured(&requests.GetPageTextStructured{
								Document: doc,
								Page: requests.Page{
									Index: 0,
								},
								CollectFontInformation: true,
							})
							Expect(err).To(BeNil())
							Expect(pageTextStructured).To(Equal(loadStructuredText(testsPath+"/testdata/text_"+prefix+"_testpdf_with_font_information.json", pageTextStructured)))
						})

						Context("and PixelPositions is enabled", func() {
							It("returns the correct font information", func() {
								pageTextStructured, err := pdfiumContainer.GetPageTextStructured(&requests.GetPageTextStructured{
									Document: doc,
									Page: requests.Page{
										Index: 0,
									},
									CollectFontInformation: true,
									PixelPositions: requests.GetPageTextStructuredPixelPositions{
										Calculate: true,
										Width:     3000,
										Height:    3000,
									},
								})
								Expect(err).To(BeNil())
								Expect(pageTextStructured).To(Equal(loadStructuredText(testsPath+"/testdata/text_"+prefix+"_testpdf_with_font_information_and_pixel_positions.json", pageTextStructured)))
							})
						})
					})
				})
			})
		})
	})
}

func loadStructuredText(path string, resp *responses.GetPageTextStructured) *responses.GetPageTextStructured {
	writeStructuredText(path, resp)
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

func writeStructuredText(path string, resp *responses.GetPageTextStructured) error {
	return nil // Comment this in case of updating pdfium versions and output has changed.

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(resp); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, buf.Bytes(), 0777); err != nil {
		return err
	}

	return nil
}
