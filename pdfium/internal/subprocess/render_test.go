package subprocess_test

import (
	"bytes"
	"encoding/gob"
	"image"
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/pdfium/internal/subprocess"
	"github.com/klippa-app/go-pdfium/pdfium/requests"
	"github.com/klippa-app/go-pdfium/pdfium/responses"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Render", func() {
	pdfium := subprocess.Pdfium{}

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

							preRender, err := ioutil.ReadFile("./testdata/render_testpdf_dpi_100.gob")
							Expect(err).To(BeNil())

							buf := bytes.NewBuffer(preRender)
							dec := gob.NewDecoder(buf)

							var preRenderImage image.RGBA
							err = dec.Decode(&preRenderImage)
							Expect(err).To(BeNil())

							Expect(renderedPage).To(Equal(&responses.RenderPage{
								Image:             &preRenderImage,
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

							preRender, err := ioutil.ReadFile("./testdata/render_testpdf_dpi_300.gob")
							Expect(err).To(BeNil())

							buf := bytes.NewBuffer(preRender)
							dec := gob.NewDecoder(buf)

							var preRenderImage image.RGBA
							err = dec.Decode(&preRenderImage)
							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPage{
								Image:             &preRenderImage,
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
								PointToPixelRatio: 4.166666666666667,
							}))

							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(827))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(1170))
						})
					})

					Context("with only the height given", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							pageSize, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
								Page:   0,
								Height: 2000,
							})
							Expect(err).To(BeNil())
							Expect(pageSize).To(Equal(&responses.RenderPage{
								PointToPixelRatio: 4.166666666666667,
							}))

							Expect(pageSize.Image.Bounds().Size().X).To(Equal(827))
							Expect(pageSize.Image.Bounds().Size().Y).To(Equal(1170))
						})
					})

					Context("with both the width and height given", func() {
						Context("and the width and height being equal", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								pageSize, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
									Page:   0,
									Width:  2000,
									Height: 2000,
								})
								Expect(err).To(BeNil())
								Expect(pageSize).To(Equal(&responses.RenderPage{
									PointToPixelRatio: 4.166666666666667,
								}))

								Expect(pageSize.Image.Bounds().Size().X).To(Equal(827))
								Expect(pageSize.Image.Bounds().Size().Y).To(Equal(1170))
							})
						})
						Context("and the width being larger than the height", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								pageSize, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
									Page:   0,
									Width:  4000,
									Height: 2000,
								})
								Expect(err).To(BeNil())
								Expect(pageSize).To(Equal(&responses.RenderPage{
									PointToPixelRatio: 4.166666666666667,
								}))

								Expect(pageSize.Image.Bounds().Size().X).To(Equal(827))
								Expect(pageSize.Image.Bounds().Size().Y).To(Equal(1170))
							})
						})

						Context("and the height being larger than the width", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								pageSize, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
									Page:   0,
									Width:  2000,
									Height: 4000,
								})
								Expect(err).To(BeNil())
								Expect(pageSize).To(Equal(&responses.RenderPage{
									PointToPixelRatio: 4.166666666666667,
								}))

								Expect(pageSize.Image.Bounds().Size().X).To(Equal(827))
								Expect(pageSize.Image.Bounds().Size().Y).To(Equal(1170))
							})
						})
					})
				})
			})
		})
	})
})
