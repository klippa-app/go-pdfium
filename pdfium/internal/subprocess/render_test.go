package subprocess_test

import (
	"bytes"
	"encoding/gob"
	"image"
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/pdfium/internal/subprocess"
	"github.com/klippa-app/go-pdfium/pdfium/pdfium_errors"
	"github.com/klippa-app/go-pdfium/pdfium/requests"
	"github.com/klippa-app/go-pdfium/pdfium/responses"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Render", func() {
	pdfium := subprocess.Pdfium{}

	Context("no document", func() {
		When("is opened", func() {
			Context("GetPageSize()", func() {
				It("returns an error", func() {
					pageSize, err := pdfium.GetPageSize(&requests.GetPageSize{
						Page: 0,
					})
					Expect(err).To(MatchError("no current document"))
					Expect(pageSize).To(BeNil())
				})
			})

			Context("GetPageSizeInPixels()", func() {
				It("returns an error", func() {
					pageSize, err := pdfium.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
						Page: 0,
						DPI:  100,
					})
					Expect(err).To(MatchError("no current document"))
					Expect(pageSize).To(BeNil())
				})
			})

			Context("RenderPageInDPI()", func() {
				It("returns an error", func() {
					renderedPage, err := pdfium.RenderPageInDPI(&requests.RenderPageInDPI{
						Page: 0,
						DPI:  300,
					})
					Expect(err).To(MatchError("no current document"))
					Expect(renderedPage).To(BeNil())
				})
			})

			Context("RenderPageInPixels()", func() {
				It("returns an error", func() {
					renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
						Page:   0,
						Width:  2000,
						Height: 2000,
					})
					Expect(err).To(MatchError("no current document"))
					Expect(renderedPage).To(BeNil())
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
				Context("GetPageSize()", func() {
					It("returns an error", func() {
						pageSize, err := pdfium.GetPageSize(&requests.GetPageSize{
							Page: 1,
						})
						Expect(err).To(MatchError(pdfium_errors.ErrPage))
						Expect(pageSize).To(BeNil())
					})
				})

				Context("GetPageSizeInPixels()", func() {
					It("returns an error", func() {
						pageSize, err := pdfium.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
							Page: 1,
							DPI:  100,
						})
						Expect(err).To(MatchError(pdfium_errors.ErrPage))
						Expect(pageSize).To(BeNil())
					})
				})

				Context("RenderPageInDPI()", func() {
					It("returns an error", func() {
						renderedPage, err := pdfium.RenderPageInDPI(&requests.RenderPageInDPI{
							Page: 1,
							DPI:  300,
						})
						Expect(err).To(MatchError(pdfium_errors.ErrPage))
						Expect(renderedPage).To(BeNil())
					})
				})

				Context("RenderPageInPixels()", func() {
					It("returns an error", func() {
						renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
							Page:   1,
							Width:  2000,
							Height: 2000,
						})
						Expect(err).To(MatchError(pdfium_errors.ErrPage))
						Expect(renderedPage).To(BeNil())
					})
				})
			})

			Context("when the page size is requested", func() {
				Context("in points", func() {
					It("returns the correct amount of points", func() {
						pageSize, err := pdfium.GetPageSize(&requests.GetPageSize{
							Page: 0,
						})
						Expect(err).To(BeNil())
						Expect(pageSize).To(Equal(&responses.GetPageSize{
							Width:  595.2755737304688,
							Height: 841.8897094726562,
						}))
					})
				})

				Context("in pixels", func() {
					Context("with no DPI", func() {
						It("returns an error", func() {
							pageSize, err := pdfium.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
								Page: 0,
							})
							Expect(err).To(MatchError("no DPI given"))
							Expect(pageSize).To(BeNil())
						})
					})

					Context("width DPI 100", func() {
						It("returns the right amount of pixels and point to pixel ratio", func() {
							pageSize, err := pdfium.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
								Page: 0,
								DPI:  100,
							})
							Expect(err).To(BeNil())
							Expect(pageSize).To(Equal(&responses.GetPageSizeInPixels{
								Width:             827,
								Height:            1170,
								PointToPixelRatio: 1.3888888888888888,
							}))
						})
					})

					Context("width DPI 300", func() {
						It("returns the right amount of pixels and point to pixel ratio", func() {
							pageSize, err := pdfium.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
								Page: 0,
								DPI:  300,
							})
							Expect(err).To(BeNil())
							Expect(pageSize).To(Equal(&responses.GetPageSizeInPixels{
								Width:             2481,
								Height:            3508,
								PointToPixelRatio: 4.166666666666667,
							}))
						})
					})
				})
			})

			Context("the page is rendered", func() {
				Context("in points", func() {
					Context("with no DPI", func() {
						It("returns an error", func() {
							renderedPage, err := pdfium.RenderPageInDPI(&requests.RenderPageInDPI{
								Page: 0,
							})
							Expect(err).To(MatchError("no DPI given"))
							Expect(renderedPage).To(BeNil())
						})
					})

					Context("width DPI 100", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPageInDPI(&requests.RenderPageInDPI{
								Page: 0,
								DPI:  100,
							})
							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPage{
								Image:             loadPrerenderedImage("./testdata/render_testpdf_dpi_100.gob"),
								PointToPixelRatio: 1.3888888888888888,
							}))
							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(827))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(1170))
						})
					})

					Context("width DPI 300", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPageInDPI(&requests.RenderPageInDPI{
								Page: 0,
								DPI:  300,
							})
							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPage{
								Image:             loadPrerenderedImage("./testdata/render_testpdf_dpi_300.gob"),
								PointToPixelRatio: 4.166666666666667,
							}))

							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(2481))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(3508))
						})
					})
				})

				Context("in pixels", func() {
					Context("with no width or height given", func() {
						It("returns an error", func() {
							renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
								Page: 0,
							})
							Expect(err).To(MatchError("no width or height given"))
							Expect(renderedPage).To(BeNil())
						})
					})

					Context("with only the width given", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
								Page:  0,
								Width: 2000,
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPage{
								Image:             loadPrerenderedImage("./testdata/render_testpdf_pixels_2000x0.gob"),
								PointToPixelRatio: 3.3597884547259587,
							}))
							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(2000))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(2829))
						})
					})

					Context("with only the height given", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
								Page:   0,
								Height: 2000,
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPage{
								Image:             loadPrerenderedImage("./testdata/render_testpdf_pixels_0x2000.gob"),
								PointToPixelRatio: 2.375608084404265,
							}))
							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(1415))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(2000))
						})
					})

					Context("with both the width and height given", func() {
						Context("and the width and height being equal", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
									Page:   0,
									Width:  2000,
									Height: 2000,
								})

								Expect(err).To(BeNil())
								Expect(renderedPage).To(Equal(&responses.RenderPage{
									Image:             loadPrerenderedImage("./testdata/render_testpdf_pixels_2000x2000.gob"),
									PointToPixelRatio: 2.375608084404265,
								}))

								Expect(renderedPage.Image.Bounds().Size().X).To(Equal(1415))
								Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(2000))
							})
						})
						Context("and the width being larger than the height", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
									Page:   0,
									Width:  4000,
									Height: 2000,
								})

								Expect(err).To(BeNil())
								Expect(renderedPage).To(Equal(&responses.RenderPage{
									Image:             loadPrerenderedImage("./testdata/render_testpdf_pixels_4000x2000.gob"),
									PointToPixelRatio: 2.375608084404265,
								}))

								Expect(renderedPage.Image.Bounds().Size().X).To(Equal(1415))
								Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(2000))
							})
						})

						Context("and the height being larger than the width", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
									Page:   0,
									Width:  2000,
									Height: 4000,
								})

								Expect(err).To(BeNil())
								Expect(renderedPage).To(Equal(&responses.RenderPage{
									Image:             loadPrerenderedImage("./testdata/render_testpdf_pixels_2000x4000.gob"),
									PointToPixelRatio: 3.3597884547259587,
								}))

								Expect(renderedPage.Image.Bounds().Size().X).To(Equal(2000))
								Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(2829))
							})
						})
					})
				})
			})
		})
	})

	// This test is only here to test the closing of an opened page.
	Context("a multipage PDF file", func() {
		BeforeEach(func() {
			pdfData, _ := ioutil.ReadFile("./testdata/test_multipage.pdf")
			pdfium.OpenDocument(&requests.OpenDocument{
				File: &pdfData,
			})
		})

		AfterEach(func() {
			pdfium.Close()
		})

		When("is opened", func() {
			Context("when another page is loaded after the first one", func() {
				Context("GetPageSize()", func() {
					It("returns the correct size for both pages", func() {
						pageSize, err := pdfium.GetPageSize(&requests.GetPageSize{
							Page: 0,
						})
						Expect(err).To(BeNil())
						Expect(pageSize).To(Equal(&responses.GetPageSize{
							Width:  595.2755737304688,
							Height: 841.8897094726562,
						}))

						pageSize, err = pdfium.GetPageSize(&requests.GetPageSize{
							Page: 1,
						})
						Expect(err).To(BeNil())
						Expect(pageSize).To(Equal(&responses.GetPageSize{
							Width:  595.2755737304688,
							Height: 841.8897094726562,
						}))
					})
				})
			})
		})
	})
})

func loadPrerenderedImage(path string) *image.RGBA {
	preRender, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	buf := bytes.NewBuffer(preRender)
	dec := gob.NewDecoder(buf)

	var preRenderImage image.RGBA
	err = dec.Decode(&preRenderImage)
	return &preRenderImage
}

func writePrerenderedImage(path string, image image.RGBA) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(&image); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, buf.Bytes(), 0777); err != nil {
		return err
	}

	return nil
}
