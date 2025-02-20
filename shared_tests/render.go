package shared_tests

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("Render", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	Context("no document", func() {
		When("is opened", func() {
			Context("GetPageSize()", func() {
				It("returns an error", func() {
					pageSize, err := PdfiumInstance.GetPageSize(&requests.GetPageSize{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Index: 0,
							},
						},
					})
					Expect(err).To(MatchError("document not given"))
					Expect(pageSize).To(BeNil())
				})
			})

			Context("GetPageSizeInPixels()", func() {
				It("returns an error", func() {
					pageSize, err := PdfiumInstance.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Index: 0,
							},
						},
						DPI: 100,
					})
					Expect(err).To(MatchError("document not given"))
					Expect(pageSize).To(BeNil())
				})
			})

			Context("RenderPageInDPI()", func() {
				It("returns an error", func() {
					renderedPage, err := PdfiumInstance.RenderPageInDPI(&requests.RenderPageInDPI{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Index: 0,
							},
						},
						DPI: 300,
					})
					Expect(err).To(MatchError("document not given"))
					Expect(renderedPage).To(BeNil())
				})
			})

			Context("RenderPageInDPI()", func() {
				It("returns an error", func() {
					renderedPage, err := PdfiumInstance.RenderPagesInDPI(&requests.RenderPagesInDPI{
						Pages: []requests.RenderPageInDPI{
							{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Index: 0,
									},
								},
								DPI: 300,
							},
						},
						Padding: 50,
					})
					Expect(err).To(MatchError("document not given"))
					Expect(renderedPage).To(BeNil())
				})
			})

			Context("RenderPageInPixels()", func() {
				It("returns an error", func() {
					renderedPage, err := PdfiumInstance.RenderPageInPixels(&requests.RenderPageInPixels{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Index: 0,
							},
						},
						Width:  2000,
						Height: 2000,
					})
					Expect(err).To(MatchError("document not given"))
					Expect(renderedPage).To(BeNil())
				})
			})

			Context("RenderPagesInPixels()", func() {
				It("returns an error", func() {
					renderedPage, err := PdfiumInstance.RenderPagesInPixels(&requests.RenderPagesInPixels{
						Pages: []requests.RenderPageInPixels{
							{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Index: 0,
									},
								},
								Width:  2000,
								Height: 2000,
							},
						},
						Padding: 50,
					})
					Expect(err).To(MatchError("document not given"))
					Expect(renderedPage).To(BeNil())
				})
			})
		})
	})

	Context("a normal PDF file", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/test.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		When("is opened", func() {
			Context("when an invalid page is given", func() {
				Context("GetPageSize()", func() {
					It("returns an error", func() {
						pageSize, err := PdfiumInstance.GetPageSize(&requests.GetPageSize{
							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    1,
								},
							},
						})
						Expect(err).To(MatchError(errors.ErrPage.Error()))
						Expect(pageSize).To(BeNil())
					})
				})

				Context("GetPageSizeInPixels()", func() {
					It("returns an error", func() {
						pageSize, err := PdfiumInstance.GetPageSizeInPixels(&requests.GetPageSizeInPixels{

							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    1,
								},
							},
							DPI: 100,
						})
						Expect(err).To(MatchError(errors.ErrPage.Error()))
						Expect(pageSize).To(BeNil())
					})
				})

				Context("RenderPageInDPI()", func() {
					It("returns an error", func() {
						renderedPage, err := PdfiumInstance.RenderPageInDPI(&requests.RenderPageInDPI{

							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    1,
								},
							},
							DPI: 300,
						})
						Expect(err).To(MatchError(errors.ErrPage.Error()))
						Expect(renderedPage).To(BeNil())
					})
				})

				Context("RenderPagesInDPI()", func() {
					It("returns an error", func() {
						renderedPage, err := PdfiumInstance.RenderPagesInDPI(&requests.RenderPagesInDPI{
							Pages: []requests.RenderPageInDPI{
								{

									Page: requests.Page{
										ByIndex: &requests.PageByIndex{
											Document: doc,
											Index:    1,
										},
									},
									DPI: 300,
								},
							},
							Padding: 50,
						})
						Expect(err).To(MatchError(errors.ErrPage.Error()))
						Expect(renderedPage).To(BeNil())
					})
				})

				Context("RenderPageInPixels()", func() {
					It("returns an error", func() {
						renderedPage, err := PdfiumInstance.RenderPageInPixels(&requests.RenderPageInPixels{

							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    1,
								},
							},
							Width:  2000,
							Height: 2000,
						})
						Expect(err).To(MatchError(errors.ErrPage.Error()))
						Expect(renderedPage).To(BeNil())
					})
				})
			})

			Context("RenderPagesInPixels()", func() {
				It("returns an error", func() {
					renderedPage, err := PdfiumInstance.RenderPagesInPixels(&requests.RenderPagesInPixels{
						Pages: []requests.RenderPageInPixels{
							{

								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    1,
									},
								},
								Width:  2000,
								Height: 2000,
							},
						},
						Padding: 50,
					})
					Expect(err).To(MatchError(errors.ErrPage.Error()))
					Expect(renderedPage).To(BeNil())
				})
			})

			Context("when the page size is requested", func() {
				Context("in points", func() {
					It("returns the correct amount of points", func() {
						pageSize, err := PdfiumInstance.GetPageSize(&requests.GetPageSize{

							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    0,
								},
							},
						})
						Expect(err).To(BeNil())
						Expect(pageSize).To(Or(
							Equal(&responses.GetPageSize{
								Width:  595.2755737304688,
								Height: 841.8897094726562,
							}),
							Equal(&responses.GetPageSize{
								Width:  595.2755737304688,
								Height: 841.8897705078125,
							}),
						))
					})
				})

				Context("in pixels", func() {
					Context("with no DPI", func() {
						It("returns an error", func() {
							pageSize, err := PdfiumInstance.GetPageSizeInPixels(&requests.GetPageSizeInPixels{

								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
							})
							Expect(err).To(MatchError("no DPI given"))
							Expect(pageSize).To(BeNil())
						})
					})

					Context("width DPI 100", func() {
						It("returns the right amount of pixels and point to pixel ratio", func() {
							pageSize, err := PdfiumInstance.GetPageSizeInPixels(&requests.GetPageSizeInPixels{

								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
								DPI: 100,
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
							pageSize, err := PdfiumInstance.GetPageSizeInPixels(&requests.GetPageSizeInPixels{

								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
								DPI: 300,
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
							renderedPage, err := PdfiumInstance.RenderPageInDPI(&requests.RenderPageInDPI{

								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
							})
							Expect(err).To(MatchError("no DPI given"))
							Expect(renderedPage).To(BeNil())
						})
					})

					Context("width DPI 100", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := PdfiumInstance.RenderPageInDPI(&requests.RenderPageInDPI{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
								DPI: 100,
							})
							Expect(err).To(BeNil())
							Expect(renderedPage).To(Not(BeNil()))
							compareRenderHash(&renderedPage.Result, Or(Equal(&responses.RenderPage{
								PointToPixelRatio: 1.3888888888888888,
								Width:             827,
								Height:            1170,
							})), TestDataPath+"/testdata/render_"+TestType+"_testpdf_dpi_100")
							Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(827))
							Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(1170))
							renderedPage.Cleanup()
						})
					})

					Context("width DPI 300", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := PdfiumInstance.RenderPageInDPI(&requests.RenderPageInDPI{

								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
								DPI: 300,
							})
							Expect(err).To(BeNil())
							Expect(renderedPage).To(Not(BeNil()))
							compareRenderHash(&renderedPage.Result, Or(Equal(&responses.RenderPage{
								PointToPixelRatio: 4.166666666666667,
								Width:             2481,
								Height:            3508,
							})), TestDataPath+"/testdata/render_"+TestType+"_testpdf_dpi_300")
							Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(2481))
							Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(3508))
							renderedPage.Cleanup()
						})
					})
				})

				Context("in pixels", func() {
					Context("with no width or height given", func() {
						It("returns an error", func() {
							renderedPage, err := PdfiumInstance.RenderPageInPixels(&requests.RenderPageInPixels{

								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
							})
							Expect(err).To(MatchError("no width or height given"))
							Expect(renderedPage).To(BeNil())
						})
					})

					Context("with only the width given", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := PdfiumInstance.RenderPageInPixels(&requests.RenderPageInPixels{

								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
								Width: 2000,
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Not(BeNil()))
							compareRenderHash(&renderedPage.Result, Or(Equal(&responses.RenderPage{
								PointToPixelRatio: 3.3597884547259587,
								Width:             2000,
								Height:            2829,
							})), TestDataPath+"/testdata/render_"+TestType+"_testpdf_pixels_2000x0", TestDataPath+"/testdata/render_"+TestType+"_testpdf_pixels_2000x0_7019")
							Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(2000))
							Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(2829))
							renderedPage.Cleanup()
						})
					})

					Context("with only the height given", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := PdfiumInstance.RenderPageInPixels(&requests.RenderPageInPixels{

								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
								Height: 2000,
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Not(BeNil()))
							compareRenderHash(&renderedPage.Result, Or(
								Equal(&responses.RenderPage{
									PointToPixelRatio: 2.375608084404265,
									Width:             1415,
									Height:            2000,
								}),
								Equal(&responses.RenderPage{
									PointToPixelRatio: 2.375607912177905,
									Width:             1415,
									Height:            2000,
								}),
							), TestDataPath+"/testdata/render_"+TestType+"_testpdf_pixels_0x2000", TestDataPath+"/testdata/render_"+TestType+"_testpdf_pixels_0x2000_7019")
							Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(1415))
							Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(2000))
							renderedPage.Cleanup()
						})
					})

					Context("with both the width and height given", func() {
						Context("and the width and height being equal", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								renderedPage, err := PdfiumInstance.RenderPageInPixels(&requests.RenderPageInPixels{

									Page: requests.Page{
										ByIndex: &requests.PageByIndex{
											Document: doc,
											Index:    0,
										},
									},
									Width:  2000,
									Height: 2000,
								})

								Expect(err).To(BeNil())
								Expect(renderedPage).To(Not(BeNil()))
								Expect(renderedPage).To(Not(BeNil()))
								compareRenderHash(&renderedPage.Result, Or(
									Equal(&responses.RenderPage{
										PointToPixelRatio: 2.375608084404265,
										Width:             1415,
										Height:            2000,
									}),
									Equal(&responses.RenderPage{
										PointToPixelRatio: 2.375607912177905,
										Width:             1415,
										Height:            2000,
									}),
								), TestDataPath+"/testdata/render_"+TestType+"_testpdf_pixels_2000x2000", TestDataPath+"/testdata/render_"+TestType+"_testpdf_pixels_2000x2000_7019")
								Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(1415))
								Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(2000))
								renderedPage.Cleanup()
							})
						})
						Context("and the width being larger than the height", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								renderedPage, err := PdfiumInstance.RenderPageInPixels(&requests.RenderPageInPixels{

									Page: requests.Page{
										ByIndex: &requests.PageByIndex{
											Document: doc,
											Index:    0,
										},
									},
									Width:  4000,
									Height: 2000,
								})

								Expect(err).To(BeNil())
								Expect(renderedPage).To(Not(BeNil()))
								compareRenderHash(&renderedPage.Result, Or(
									Equal(&responses.RenderPage{
										PointToPixelRatio: 2.375608084404265,
										Width:             1415,
										Height:            2000,
									}),
									Equal(&responses.RenderPage{
										PointToPixelRatio: 2.375607912177905,
										Width:             1415,
										Height:            2000,
									}),
								), TestDataPath+"/testdata/render_"+TestType+"_testpdf_pixels_4000x2000", TestDataPath+"/testdata/render_"+TestType+"_testpdf_pixels_4000x2000_7019")
								Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(1415))
								Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(2000))
								renderedPage.Cleanup()
							})
						})

						Context("and the height being larger than the width", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								renderedPage, err := PdfiumInstance.RenderPageInPixels(&requests.RenderPageInPixels{

									Page: requests.Page{
										ByIndex: &requests.PageByIndex{
											Document: doc,
											Index:    0,
										},
									},
									Width:  2000,
									Height: 4000,
								})

								Expect(err).To(BeNil())
								Expect(renderedPage).To(Not(BeNil()))
								compareRenderHash(&renderedPage.Result, Or(Equal(&responses.RenderPage{
									PointToPixelRatio: 3.3597884547259587,
									Width:             2000,
									Height:            2829,
								})), TestDataPath+"/testdata/render_"+TestType+"_testpdf_pixels_2000x4000")
								Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(2000))
								Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(2829))
								renderedPage.Cleanup()
							})
						})
					})
				})
			})

			Context("the pages are rendered", func() {
				Context("in points", func() {
					Context("with no pages given", func() {
						It("returns an error", func() {
							renderedPage, err := PdfiumInstance.RenderPagesInDPI(&requests.RenderPagesInDPI{
								Pages: []requests.RenderPageInDPI{},
							})
							Expect(err).To(MatchError("no pages given"))
							Expect(renderedPage).To(BeNil())
						})
					})

					Context("with no DPI", func() {
						It("returns an error", func() {
							renderedPage, err := PdfiumInstance.RenderPagesInDPI(&requests.RenderPagesInDPI{
								Pages: []requests.RenderPageInDPI{
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
									},
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
									},
								},
							})
							Expect(err).To(MatchError("no DPI given for requested page 0"))
							Expect(renderedPage).To(BeNil())
						})
					})

					Context("with DPI 100", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := PdfiumInstance.RenderPagesInDPI(&requests.RenderPagesInDPI{
								Pages: []requests.RenderPageInDPI{
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										DPI: 100,
									},
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										DPI: 100,
									},
								},
							})
							Expect(err).To(BeNil())
							Expect(renderedPage).To(Not(BeNil()))
							compareRenderHashForPages(&renderedPage.Result, Or(Equal(&responses.RenderPages{
								Width:  827,
								Height: 2340,
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 1.3888888888888888,
										Width:             827,
										Height:            1170,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 1.3888888888888888,
										Width:             827,
										Height:            1170,
										X:                 0,
										Y:                 1170,
									},
								},
							})), TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_dpi_100")
							Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(827))
							Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(2340))
							renderedPage.Cleanup()
						})
					})

					Context("with DPI 300", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := PdfiumInstance.RenderPagesInDPI(&requests.RenderPagesInDPI{
								Pages: []requests.RenderPageInDPI{
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										DPI: 300,
									},
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										DPI: 300,
									},
								},
							})
							Expect(err).To(BeNil())
							Expect(renderedPage).To(Not(BeNil()))
							compareRenderHashForPages(&renderedPage.Result, Or(Equal(&responses.RenderPages{
								Width:  2481,
								Height: 7016,
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 4.166666666666667,
										Width:             2481,
										Height:            3508,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 4.166666666666667,
										Width:             2481,
										Height:            3508,
										X:                 0,
										Y:                 3508,
									},
								},
							})), TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_dpi_300")
							Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(2481))
							Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(7016))
							renderedPage.Cleanup()
						})
					})

					Context("with different DPI per page", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := PdfiumInstance.RenderPagesInDPI(&requests.RenderPagesInDPI{
								Pages: []requests.RenderPageInDPI{
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										DPI: 200,
									},
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										DPI: 300,
									},
								},
							})
							Expect(err).To(BeNil())
							Expect(renderedPage).To(Not(BeNil()))
							compareRenderHashForPages(&renderedPage.Result, Or(Equal(&responses.RenderPages{
								Width:  2481,
								Height: 5847,
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 2.7777777777777777,
										Width:             1654,
										Height:            2339,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 4.166666666666667,
										Width:             2481,
										Height:            3508,
										X:                 0,
										Y:                 2339,
									},
								},
							})), TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_dpi_200_300")
							Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(2481))
							Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(5847))
							renderedPage.Cleanup()
						})
					})

					Context("with padding between pages", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := PdfiumInstance.RenderPagesInDPI(&requests.RenderPagesInDPI{
								Pages: []requests.RenderPageInDPI{
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										DPI: 300,
									},
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										DPI: 300,
									},
								},
								Padding: 50,
							})
							Expect(err).To(BeNil())
							Expect(renderedPage).To(Not(BeNil()))
							compareRenderHashForPages(&renderedPage.Result, Or(Equal(&responses.RenderPages{
								Width:  2481,
								Height: 7066,
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 4.166666666666667,
										Width:             2481,
										Height:            3508,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 4.166666666666667,
										Width:             2481,
										Height:            3508,
										X:                 0,
										Y:                 3558,
									},
								},
							})), TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_dpi_300_padding_50")
							Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(2481))
							Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(7066))
							renderedPage.Cleanup()
						})
					})
				})

				Context("in pixels", func() {
					Context("with no pages given", func() {
						It("returns an error", func() {
							renderedPage, err := PdfiumInstance.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{},
							})
							Expect(err).To(MatchError("no pages given"))
							Expect(renderedPage).To(BeNil())
						})
					})
					Context("with no width or height given", func() {
						It("returns an error", func() {
							renderedPage, err := PdfiumInstance.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
									},
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
									},
								},
							})
							Expect(err).To(MatchError("no width or height given for requested page 0"))
							Expect(renderedPage).To(BeNil())
						})
					})

					Context("with only the width given", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := PdfiumInstance.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Width: 2000,
									},
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Width: 2000,
									},
								},
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Not(BeNil()))
							compareRenderHashForPages(&renderedPage.Result, Or(Equal(&responses.RenderPages{
								Width:  2000,
								Height: 5658,
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 3.3597884547259587,
										Width:             2000,
										Height:            2829,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 3.3597884547259587,
										Width:             2000,
										Height:            2829,
										X:                 0,
										Y:                 2829,
									},
								},
							})), TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_pixels_2000x0")
							Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(2000))
							Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(5658))
							renderedPage.Cleanup()
						})
					})

					Context("with only the height given", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := PdfiumInstance.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Height: 2000,
									},
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Height: 2000,
									},
								},
							})

							Expect(err).To(BeNil())
							compareRenderHashForPages(&renderedPage.Result, Or(
								Equal(&responses.RenderPages{
									Width:  1415,
									Height: 4000,
									Pages: []responses.RenderPagesPage{
										{
											PointToPixelRatio: 2.375608084404265,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
										},
										{
											PointToPixelRatio: 2.375608084404265,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 2000,
										},
									},
								}),
								Equal(&responses.RenderPages{
									Width:  1415,
									Height: 4000,
									Pages: []responses.RenderPagesPage{
										{
											Page:              0,
											PointToPixelRatio: 2.375607912177905,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
											HasTransparency:   false,
										},
										{
											Page:              0,
											PointToPixelRatio: 2.375607912177905,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 2000,
											HasTransparency:   false,
										},
									},
								}),
							), TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_pixels_0x2000", TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_pixels_0x2000_7019")
							Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(1415))
							Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(4000))
							renderedPage.Cleanup()
						})
					})

					Context("with both the width and height given", func() {
						Context("and the width and height being equal", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								renderedPage, err := PdfiumInstance.RenderPagesInPixels(&requests.RenderPagesInPixels{
									Pages: []requests.RenderPageInPixels{
										{

											Page: requests.Page{
												ByIndex: &requests.PageByIndex{
													Document: doc,
													Index:    0,
												},
											},
											Width:  2000,
											Height: 2000,
										},
										{

											Page: requests.Page{
												ByIndex: &requests.PageByIndex{
													Document: doc,
													Index:    0,
												},
											},
											Width:  2000,
											Height: 2000,
										},
									},
								})

								Expect(err).To(BeNil())
								Expect(renderedPage).To(Not(BeNil()))
								compareRenderHashForPages(&renderedPage.Result, Or(
									Equal(&responses.RenderPages{
										Width:  1415,
										Height: 4000,
										Pages: []responses.RenderPagesPage{
											{
												PointToPixelRatio: 2.375608084404265,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 0,
											},
											{
												PointToPixelRatio: 2.375608084404265,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 2000,
											},
										},
									}),
									Equal(&responses.RenderPages{
										Width:  1415,
										Height: 4000,
										Pages: []responses.RenderPagesPage{
											{
												Page:              0,
												PointToPixelRatio: 2.375607912177905,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 0,
												HasTransparency:   false,
											},
											{
												Page:              0,
												PointToPixelRatio: 2.375607912177905,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 2000,
												HasTransparency:   false,
											},
										},
									}),
								), TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_pixels_2000x2000", TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_pixels_2000x2000_7019")
								Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(1415))
								Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(4000))
								renderedPage.Cleanup()
							})
						})
						Context("and the width being larger than the height", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								renderedPage, err := PdfiumInstance.RenderPagesInPixels(&requests.RenderPagesInPixels{
									Pages: []requests.RenderPageInPixels{
										{

											Page: requests.Page{
												ByIndex: &requests.PageByIndex{
													Document: doc,
													Index:    0,
												},
											},
											Width:  4000,
											Height: 2000,
										},
										{

											Page: requests.Page{
												ByIndex: &requests.PageByIndex{
													Document: doc,
													Index:    0,
												},
											},
											Width:  4000,
											Height: 2000,
										},
									},
								})

								Expect(err).To(BeNil())
								Expect(renderedPage).To(Not(BeNil()))
								compareRenderHashForPages(&renderedPage.Result, Or(
									Equal(&responses.RenderPages{
										Width:  1415,
										Height: 4000,
										Pages: []responses.RenderPagesPage{
											{
												PointToPixelRatio: 2.375608084404265,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 0,
											},
											{
												PointToPixelRatio: 2.375608084404265,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 2000,
											},
										},
									}),
									Equal(&responses.RenderPages{
										Width:  1415,
										Height: 4000,
										Pages: []responses.RenderPagesPage{
											{
												Page:              0,
												PointToPixelRatio: 2.375607912177905,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 0,
												HasTransparency:   false,
											},
											{
												Page:              0,
												PointToPixelRatio: 2.375607912177905,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 2000,
												HasTransparency:   false,
											},
										},
									}),
								), TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_pixels_4000x2000", TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_pixels_4000x2000_7019")
								Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(1415))
								Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(4000))
								renderedPage.Cleanup()
							})
						})

						Context("and the height being larger than the width", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								renderedPage, err := PdfiumInstance.RenderPagesInPixels(&requests.RenderPagesInPixels{
									Pages: []requests.RenderPageInPixels{
										{

											Page: requests.Page{
												ByIndex: &requests.PageByIndex{
													Document: doc,
													Index:    0,
												},
											},
											Width:  2000,
											Height: 4000,
										},
										{

											Page: requests.Page{
												ByIndex: &requests.PageByIndex{
													Document: doc,
													Index:    0,
												},
											},
											Width:  2000,
											Height: 4000,
										},
									},
								})

								Expect(err).To(BeNil())
								Expect(renderedPage).To(Not(BeNil()))
								compareRenderHashForPages(&renderedPage.Result, Or(Equal(&responses.RenderPages{
									Width:  2000,
									Height: 5658,
									Pages: []responses.RenderPagesPage{
										{
											PointToPixelRatio: 3.3597884547259587,
											Width:             2000,
											Height:            2829,
											X:                 0,
											Y:                 0,
										},
										{
											PointToPixelRatio: 3.3597884547259587,
											Width:             2000,
											Height:            2829,
											X:                 0,
											Y:                 2829,
										},
									},
								})), TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_pixels_2000x4000")
								Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(2000))
								Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(5658))
								renderedPage.Cleanup()
							})
						})
					})

					Context("with the width being different between pages", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := PdfiumInstance.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Width: 2000,
									},
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Width: 1500,
									},
								},
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Not(BeNil()))
							compareRenderHashForPages(&renderedPage.Result, Or(Equal(&responses.RenderPages{
								Width:  2000,
								Height: 4951,
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 3.3597884547259587,
										Width:             2000,
										Height:            2829,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 2.519841341044469,
										Width:             1500,
										Height:            2122,
										X:                 0,
										Y:                 2829,
									},
								},
							})), TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_pixels_2000x0_1500x0")
							Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(2000))
							Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(4951))
							renderedPage.Cleanup()
						})
					})

					Context("with the height being different between pages", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := PdfiumInstance.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Height: 2000,
									},
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Height: 1500,
									},
								},
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Not(BeNil()))
							compareRenderHashForPages(&renderedPage.Result, Or(
								Equal(&responses.RenderPages{
									Width:  1415,
									Height: 3500,
									Pages: []responses.RenderPagesPage{
										{
											PointToPixelRatio: 2.375608084404265,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
										},
										{
											PointToPixelRatio: 1.7817060633031987,
											Width:             1061,
											Height:            1500,
											X:                 0,
											Y:                 2000,
										},
									},
								}),
								Equal(&responses.RenderPages{
									Width:  1415,
									Height: 3500,
									Pages: []responses.RenderPagesPage{
										{
											Page:              0,
											PointToPixelRatio: 2.375607912177905,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
											HasTransparency:   false,
										},
										{
											Page:              0,
											PointToPixelRatio: 1.7817059341334285,
											Width:             1061,
											Height:            1500,
											X:                 0,
											Y:                 2000,
											HasTransparency:   false,
										},
									},
								}),
							), TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_pixels_0x2000_0x1500", TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_pixels_0x2000_0x1500_7019")
							Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(1415))
							Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(3500))
							renderedPage.Cleanup()
						})
					})

					Context("with the width and height being different between pages", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := PdfiumInstance.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Width:  2000,
										Height: 2000,
									},
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Width:  1500,
										Height: 1500,
									},
								},
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Not(BeNil()))
							compareRenderHashForPages(&renderedPage.Result, Or(
								Equal(&responses.RenderPages{
									Width:  1415,
									Height: 3500,
									Pages: []responses.RenderPagesPage{
										{
											PointToPixelRatio: 2.375608084404265,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
										},
										{
											PointToPixelRatio: 1.7817060633031987,
											Width:             1061,
											Height:            1500,
											X:                 0,
											Y:                 2000,
										},
									},
								}),
								Equal(&responses.RenderPages{
									Width:  1415,
									Height: 3500,
									Pages: []responses.RenderPagesPage{
										{
											Page:              0,
											PointToPixelRatio: 2.375607912177905,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
											HasTransparency:   false,
										},
										{
											Page:              0,
											PointToPixelRatio: 1.7817059341334285,
											Width:             1061,
											Height:            1500,
											X:                 0,
											Y:                 2000,
											HasTransparency:   false,
										},
									},
								}),
							), TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_pixels_2000x2000_1500x1500", TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_pixels_2000x2000_1500x1500_7019")
							Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(1415))
							Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(3500))
							renderedPage.Cleanup()
						})
					})

					Context("with padding between pages", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := PdfiumInstance.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Width:  2000,
										Height: 2000,
									},
									{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Width:  2000,
										Height: 2000,
									},
								},
								Padding: 50,
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Not(BeNil()))
							compareRenderHashForPages(&renderedPage.Result, Or(
								Equal(&responses.RenderPages{
									Width:  1415,
									Height: 4050,
									Pages: []responses.RenderPagesPage{
										{
											PointToPixelRatio: 2.375608084404265,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
										},
										{
											PointToPixelRatio: 2.375608084404265,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 2050,
										},
									},
								}),
								Equal(&responses.RenderPages{
									Image:  nil,
									Width:  1415,
									Height: 4050,
									Pages: []responses.RenderPagesPage{
										{
											Page:              0,
											PointToPixelRatio: 2.375607912177905,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
											HasTransparency:   false,
										},
										{
											Page:              0,
											PointToPixelRatio: 2.375607912177905,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 2050,
											HasTransparency:   false,
										},
									},
								}),
							), TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_pixels_2000x2000_2000x2000_padding_50", TestDataPath+"/testdata/render_"+TestType+"_pages_testpdf_pixels_2000x2000_2000x2000_padding_50_7019")
							Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(1415))
							Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(4050))
							renderedPage.Cleanup()
						})
					})
				})

				Context("and directly rendered to a file", func() {
					Context("with no output target given", func() {
						It("returns an error", func() {
							renderedPage, err := PdfiumInstance.RenderToFile(&requests.RenderToFile{
								OutputFormat: requests.RenderToFileOutputFormatJPG,
								RenderPagesInPixels: &requests.RenderPagesInPixels{
									Pages: []requests.RenderPageInPixels{
										{

											Page: requests.Page{
												ByIndex: &requests.PageByIndex{
													Document: doc,
													Index:    0,
												},
											},
											Width:  2000,
											Height: 2000,
										},
									},
								},
							})
							Expect(err).To(MatchError("invalid output target given"))
							Expect(renderedPage).To(BeNil())
						})
					})

					Context("with no output format given", func() {
						It("returns an error", func() {
							renderedPage, err := PdfiumInstance.RenderToFile(&requests.RenderToFile{
								OutputTarget: requests.RenderToFileOutputTargetBytes,
								RenderPagesInPixels: &requests.RenderPagesInPixels{
									Pages: []requests.RenderPageInPixels{
										{

											Page: requests.Page{
												ByIndex: &requests.PageByIndex{
													Document: doc,
													Index:    0,
												},
											},
											Width:  2000,
											Height: 2000,
										},
									},
								},
							})
							Expect(err).To(MatchError("invalid output format given"))
							Expect(renderedPage).To(BeNil())
						})
					})

					Context("with no render operation given", func() {
						It("returns an error", func() {
							renderedPage, err := PdfiumInstance.RenderToFile(&requests.RenderToFile{
								OutputTarget: requests.RenderToFileOutputTargetBytes,
								OutputFormat: requests.RenderToFileOutputFormatJPG,
							})
							Expect(err).To(MatchError("no render operation given"))
							Expect(renderedPage).To(BeNil())
						})
					})

					Context("in pixels", func() {
						Context("with 1 page", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								request := &requests.RenderToFile{
									OutputTarget: requests.RenderToFileOutputTargetBytes,
									OutputFormat: requests.RenderToFileOutputFormatJPG,
									RenderPageInPixels: &requests.RenderPageInPixels{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Width:  2000,
										Height: 2000,
									},
								}
								renderedFile, err := PdfiumInstance.RenderToFile(request)

								Expect(err).To(BeNil())
								compareFileHash(request, renderedFile, Or(
									Equal(&responses.RenderToFile{
										Pages: []responses.RenderPagesPage{
											{
												Page:              0,
												PointToPixelRatio: 2.375608084404265,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 0,
											},
										},
										Width:             1415,
										Height:            2000,
										PointToPixelRatio: 2.375608084404265,
									}),
									Equal(&responses.RenderToFile{
										Pages: []responses.RenderPagesPage{
											{
												Page:              0,
												PointToPixelRatio: 2.375607912177905,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 0,
												HasTransparency:   false,
											},
										},
										ImageBytes:        nil,
										ImagePath:         "",
										Width:             1415,
										Height:            2000,
										PointToPixelRatio: 2.375607912177905,
									}),
								), TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_1_page_pixels", TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_1_page_pixels_7019")
							})
						})

						Context("with multiple pages", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								request := &requests.RenderToFile{
									OutputTarget: requests.RenderToFileOutputTargetBytes,
									OutputFormat: requests.RenderToFileOutputFormatJPG,
									RenderPagesInPixels: &requests.RenderPagesInPixels{
										Pages: []requests.RenderPageInPixels{
											{

												Page: requests.Page{
													ByIndex: &requests.PageByIndex{
														Document: doc,
														Index:    0,
													},
												},
												Width:  2000,
												Height: 2000,
											},
											{

												Page: requests.Page{
													ByIndex: &requests.PageByIndex{
														Document: doc,
														Index:    0,
													},
												},
												Width:  2000,
												Height: 2000,
											},
										},
									},
								}
								renderedFile, err := PdfiumInstance.RenderToFile(request)

								Expect(err).To(BeNil())
								compareFileHash(request, renderedFile, Or(
									Equal(&responses.RenderToFile{
										Pages: []responses.RenderPagesPage{
											{
												Page:              0,
												PointToPixelRatio: 2.375608084404265,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 0,
											},
											{
												Page:              0,
												PointToPixelRatio: 2.375608084404265,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 2000,
											},
										},
										Width:  1415,
										Height: 4000,
									}),
									Equal(&responses.RenderToFile{
										Pages: []responses.RenderPagesPage{
											{
												Page:              0,
												PointToPixelRatio: 2.375607912177905,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 0,
												HasTransparency:   false,
											},
											{
												Page:              0,
												PointToPixelRatio: 2.375607912177905,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 2000,
												HasTransparency:   false,
											},
										},
										ImageBytes:        nil,
										ImagePath:         "",
										Width:             1415,
										Height:            4000,
										PointToPixelRatio: 0,
									}),
								), TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_multiple_pages_pixels", TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_multiple_pages_pixels_7019")
							})
						})
					})

					Context("in points", func() {
						Context("with 1 page", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								request := &requests.RenderToFile{
									OutputTarget: requests.RenderToFileOutputTargetBytes,
									OutputFormat: requests.RenderToFileOutputFormatJPG,
									RenderPageInDPI: &requests.RenderPageInDPI{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										DPI: 200,
									},
								}
								renderedFile, err := PdfiumInstance.RenderToFile(request)

								Expect(err).To(BeNil())
								compareFileHash(request, renderedFile, Or(Equal(&responses.RenderToFile{
									Pages: []responses.RenderPagesPage{
										{
											Page:              0,
											PointToPixelRatio: 2.7777777777777777,
											Width:             1654,
											Height:            2339,
											X:                 0,
											Y:                 0,
										},
									},
									Width:             1654,
									Height:            2339,
									PointToPixelRatio: 2.7777777777777777,
								})), TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_1_page_dpi")
							})
						})

						Context("with multiple pages", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								request := &requests.RenderToFile{
									OutputTarget: requests.RenderToFileOutputTargetBytes,
									OutputFormat: requests.RenderToFileOutputFormatJPG,
									RenderPagesInDPI: &requests.RenderPagesInDPI{
										Pages: []requests.RenderPageInDPI{
											{

												Page: requests.Page{
													ByIndex: &requests.PageByIndex{
														Document: doc,
														Index:    0,
													},
												},
												DPI: 200,
											},
											{

												Page: requests.Page{
													ByIndex: &requests.PageByIndex{
														Document: doc,
														Index:    0,
													},
												},
												DPI: 200,
											},
										},
									},
								}
								renderedFile, err := PdfiumInstance.RenderToFile(request)

								Expect(err).To(BeNil())
								compareFileHash(request, renderedFile, Or(Equal(&responses.RenderToFile{
									Pages: []responses.RenderPagesPage{
										{
											Page:              0,
											PointToPixelRatio: 2.7777777777777777,
											Width:             1654,
											Height:            2339,

											X: 0,
											Y: 0,
										},
										{
											Page:              0,
											PointToPixelRatio: 2.7777777777777777,
											Width:             1654,
											Height:            2339,

											X: 0,
											Y: 2339,
										},
									},
									Width:  1654,
									Height: 4678,
								})), TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_multiple_pages_dpi")
							})
						})
					})

					Context("to jpeg", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							request := &requests.RenderToFile{
								OutputTarget: requests.RenderToFileOutputTargetBytes,
								OutputFormat: requests.RenderToFileOutputFormatJPG,
								RenderPageInPixels: &requests.RenderPageInPixels{

									Page: requests.Page{
										ByIndex: &requests.PageByIndex{
											Document: doc,
											Index:    0,
										},
									},
									Width:  2000,
									Height: 2000,
								},
							}
							renderedFile, err := PdfiumInstance.RenderToFile(request)

							Expect(err).To(BeNil())
							compareFileHash(request, renderedFile, Or(
								Equal(&responses.RenderToFile{
									Pages: []responses.RenderPagesPage{
										{
											Page:              0,
											PointToPixelRatio: 2.375608084404265,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
										},
									},
									Width:             1415,
									Height:            2000,
									PointToPixelRatio: 2.375608084404265,
								}),
								Equal(&responses.RenderToFile{
									Pages: []responses.RenderPagesPage{
										{
											Page:              0,
											PointToPixelRatio: 2.375607912177905,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
											HasTransparency:   false,
										},
									},
									ImageBytes:        nil,
									ImagePath:         "",
									Width:             1415,
									Height:            2000,
									PointToPixelRatio: 2.375607912177905,
								}),
							), TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_jpg", TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_jpg_7019")
						})
					})

					Context("to png", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							request := &requests.RenderToFile{
								OutputTarget: requests.RenderToFileOutputTargetBytes,
								OutputFormat: requests.RenderToFileOutputFormatPNG,
								RenderPageInPixels: &requests.RenderPageInPixels{

									Page: requests.Page{
										ByIndex: &requests.PageByIndex{
											Document: doc,
											Index:    0,
										},
									},
									Width:  2000,
									Height: 2000,
								},
							}
							renderedFile, err := PdfiumInstance.RenderToFile(request)

							Expect(err).To(BeNil())
							compareFileHash(request, renderedFile, Or(
								Equal(&responses.RenderToFile{
									Pages: []responses.RenderPagesPage{
										{
											Page:              0,
											PointToPixelRatio: 2.375608084404265,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
										},
									},
									Width:             1415,
									Height:            2000,
									PointToPixelRatio: 2.375608084404265,
								}),
								Equal(&responses.RenderToFile{
									Pages: []responses.RenderPagesPage{
										{
											Page:              0,
											PointToPixelRatio: 2.375607912177905,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
											HasTransparency:   false,
										},
									},
									ImageBytes:        nil,
									ImagePath:         "",
									Width:             1415,
									Height:            2000,
									PointToPixelRatio: 2.375607912177905,
								}),
							), TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_png", TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_png_7019")
						})
					})

					Context("to bytes", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							request := &requests.RenderToFile{
								OutputTarget: requests.RenderToFileOutputTargetBytes,
								OutputFormat: requests.RenderToFileOutputFormatJPG,
								RenderPageInPixels: &requests.RenderPageInPixels{

									Page: requests.Page{
										ByIndex: &requests.PageByIndex{
											Document: doc,
											Index:    0,
										},
									},
									Width:  2000,
									Height: 2000,
								},
							}
							renderedFile, err := PdfiumInstance.RenderToFile(request)

							Expect(err).To(BeNil())
							compareFileHash(request, renderedFile, Or(
								Equal(&responses.RenderToFile{
									Pages: []responses.RenderPagesPage{
										{
											Page:              0,
											PointToPixelRatio: 2.375608084404265,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
										},
									},
									Width:             1415,
									Height:            2000,
									PointToPixelRatio: 2.375608084404265,
								}),
								Equal(&responses.RenderToFile{
									Pages: []responses.RenderPagesPage{
										{
											Page:              0,
											PointToPixelRatio: 2.375607912177905,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
											HasTransparency:   false,
										},
									},
									ImageBytes:        nil,
									ImagePath:         "",
									Width:             1415,
									Height:            2000,
									PointToPixelRatio: 2.375607912177905,
								}),
							), TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_bytes", TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_bytes_7019")
						})
					})

					Context("to file", func() {
						Context("with an invalid filepath given", func() {
							It("returns an error", func() {
								request := &requests.RenderToFile{
									OutputTarget: requests.RenderToFileOutputTargetFile,
									OutputFormat: requests.RenderToFileOutputFormatJPG,
									RenderPageInPixels: &requests.RenderPageInPixels{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Width:  2000,
										Height: 2000,
									},
									TargetFilePath: "/file/path/that/is/invalid",
								}
								renderedFile, err := PdfiumInstance.RenderToFile(request)
								Expect(err).To(Not(BeNil()))
								Expect(renderedFile).To(BeNil())
							})
						})

						Context("with a filepath given", func() {
							It("returns the right image, point to pixel ratio and resolution in the given filepath", func() {
								tmpPath, err := os.CreateTemp("", "render_file_testpdf_filepath_*")
								Expect(err).To(BeNil())
								err = tmpPath.Close()
								Expect(err).To(BeNil())
								defer os.Remove(tmpPath.Name())
								request := &requests.RenderToFile{
									OutputTarget: requests.RenderToFileOutputTargetFile,
									OutputFormat: requests.RenderToFileOutputFormatJPG,
									RenderPageInPixels: &requests.RenderPageInPixels{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Width:  2000,
										Height: 2000,
									},
									TargetFilePath: tmpPath.Name(),
								}
								renderedFile, err := PdfiumInstance.RenderToFile(request)

								Expect(err).To(BeNil())
								compareFileHash(request, renderedFile, Or(
									Equal(&responses.RenderToFile{
										Pages: []responses.RenderPagesPage{
											{
												Page:              0,
												PointToPixelRatio: 2.375608084404265,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 0,
											},
										},
										Width:             1415,
										Height:            2000,
										PointToPixelRatio: 2.375608084404265,
									}),
									Equal(&responses.RenderToFile{
										Pages: []responses.RenderPagesPage{
											{
												Page:              0,
												PointToPixelRatio: 2.375607912177905,
												Width:             1415,
												Height:            2000,
												X:                 0,
												Y:                 0,
												HasTransparency:   false,
											},
										},
										ImageBytes:        nil,
										ImagePath:         "",
										Width:             1415,
										Height:            2000,
										PointToPixelRatio: 2.375607912177905,
									}),
								), TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_filepath", TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_filepath_7019")
							})
						})

						Context("with no filepath given", func() {
							It("returns the right image, point to pixel ratio and resolution in a temp filepath", func() {
								request := &requests.RenderToFile{
									OutputTarget: requests.RenderToFileOutputTargetFile,
									OutputFormat: requests.RenderToFileOutputFormatJPG,
									RenderPageInPixels: &requests.RenderPageInPixels{

										Page: requests.Page{
											ByIndex: &requests.PageByIndex{
												Document: doc,
												Index:    0,
											},
										},
										Width:  2000,
										Height: 2000,
									},
								}
								renderedFile, err := PdfiumInstance.RenderToFile(request)

								Expect(err).To(BeNil())
								compareFileHash(request, renderedFile, Or(Equal(&responses.RenderToFile{
									Pages: []responses.RenderPagesPage{
										{
											Page:              0,
											PointToPixelRatio: 2.375608084404265,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
										},
									},
									Width:             1415,
									Height:            2000,
									PointToPixelRatio: 2.375608084404265,
								}), Equal(&responses.RenderToFile{
									Pages: []responses.RenderPagesPage{
										{
											Page:              0,
											PointToPixelRatio: 2.375607912177905,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
											HasTransparency:   false,
										},
									},
									ImageBytes:        nil,
									ImagePath:         "",
									Width:             1415,
									Height:            2000,
									PointToPixelRatio: 2.375607912177905,
								})), TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_no_filepath", TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_no_filepath_7019")
							})
						})

						Context("with a max filesize given", func() {
							Context("while rendering to jpg", func() {
								Context("with a max filesize that is unreasonable", func() {
									It("returns an error", func() {
										renderedPage, err := PdfiumInstance.RenderToFile(&requests.RenderToFile{
											OutputTarget: requests.RenderToFileOutputTargetBytes,
											OutputFormat: requests.RenderToFileOutputFormatJPG,
											RenderPageInPixels: &requests.RenderPageInPixels{

												Page: requests.Page{
													ByIndex: &requests.PageByIndex{
														Document: doc,
														Index:    0,
													},
												},
												Width:  2000,
												Height: 2000,
											},
											MaxFileSize: 1000, // 1000 bytes
										})
										Expect(err).To(MatchError("PDF image would exceed maximum filesize"))
										Expect(renderedPage).To(BeNil())
									})
								})
								Context("with a max filesize that is reasonable", func() {
									It("returns the right image, point to pixel ratio and resolution, filesize under the limit", func() {
										request := &requests.RenderToFile{
											OutputTarget: requests.RenderToFileOutputTargetBytes,
											OutputFormat: requests.RenderToFileOutputFormatJPG,
											RenderPageInPixels: &requests.RenderPageInPixels{

												Page: requests.Page{
													ByIndex: &requests.PageByIndex{
														Document: doc,
														Index:    0,
													},
												},
												Width:  2000,
												Height: 2000,
											},
											MaxFileSize: 60000, // 60 kb
										}
										renderedFile, err := PdfiumInstance.RenderToFile(request)

										Expect(err).To(BeNil())
										compareFileHash(request, renderedFile, Or(Equal(&responses.RenderToFile{
											Pages: []responses.RenderPagesPage{
												{
													Page:              0,
													PointToPixelRatio: 2.375608084404265,
													Width:             1415,
													Height:            2000,
													X:                 0,
													Y:                 0,
												},
											},
											Width:             1415,
											Height:            2000,
											PointToPixelRatio: 2.375608084404265,
										}), Equal(&responses.RenderToFile{
											Pages: []responses.RenderPagesPage{
												{
													Page:              0,
													PointToPixelRatio: 2.375607912177905,
													Width:             1415,
													Height:            2000,
													X:                 0,
													Y:                 0,
													HasTransparency:   false,
												},
											},
											ImageBytes:        nil,
											ImagePath:         "",
											Width:             1415,
											Height:            2000,
											PointToPixelRatio: 2.375607912177905,
										})), TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_max_filesize", TestDataPath+"/testdata/render_"+TestType+"_file_testpdf_max_filesize_7019")
										if renderedFile.ImageBytes != nil {
											Expect(len(*renderedFile.ImageBytes)).To(BeNumerically("<=", 60000))
										}
									})
								})
							})

							Context("while rendering to png", func() {
								Context("with a max filesize that is over the rendered size", func() {
									It("returns an error", func() {
										renderedPage, err := PdfiumInstance.RenderToFile(&requests.RenderToFile{
											OutputTarget: requests.RenderToFileOutputTargetBytes,
											OutputFormat: requests.RenderToFileOutputFormatPNG,
											RenderPageInPixels: &requests.RenderPageInPixels{

												Page: requests.Page{
													ByIndex: &requests.PageByIndex{
														Document: doc,
														Index:    0,
													},
												},
												Width:  2000,
												Height: 2000,
											},
											MaxFileSize: 1000, // 1000 bytes
										})
										Expect(err).To(MatchError("PDF image would exceed maximum filesize"))
										Expect(renderedPage).To(BeNil())
									})
								})
							})
						})
					})
				})
			})
		})
	})

	Context("a PDF file that uses an alpha channel", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/alpha_channel.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		When("it is rendered", func() {
			It("returns the right image", func() {
				renderedPage, err := PdfiumInstance.RenderPageInDPI(&requests.RenderPageInDPI{

					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					DPI: 200,
				})

				Expect(err).To(BeNil())
				compareRenderHash(&renderedPage.Result, Or(Equal(&responses.RenderPage{
					PointToPixelRatio: 2.7777777777777777,
					Width:             1653,
					Height:            2339,
				})), TestDataPath+"/testdata/render_"+TestType+"_page_alpha_channel", TestDataPath+"/testdata/render_"+TestType+"_page_alpha_channel_7019")
				renderedPage.Cleanup()
			})
		})
	})

	// This test is only here to test the closing of an opened page.
	Context("a multipage PDF file", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/test_multipage.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		When("is opened", func() {
			Context("when another page is loaded after the first one", func() {
				Context("GetPageSize()", func() {
					It("returns the correct size for both pages", func() {
						pageSize, err := PdfiumInstance.GetPageSize(&requests.GetPageSize{

							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    0,
								},
							},
						})
						Expect(err).To(BeNil())
						Expect(pageSize).To(Or(
							Equal(&responses.GetPageSize{
								Width:  595.2755737304688,
								Height: 841.8897094726562,
							}),
							Equal(&responses.GetPageSize{
								Width:  595.2755737304688,
								Height: 841.8897705078125,
							}),
						))

						pageSize, err = PdfiumInstance.GetPageSize(&requests.GetPageSize{

							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    1,
								},
							},
						})
						Expect(err).To(BeNil())
						Expect(pageSize).To(Or(
							Equal(&responses.GetPageSize{
								Page:   1,
								Width:  595.2755737304688,
								Height: 841.8897094726562,
							}),
							Equal(&responses.GetPageSize{
								Page:   1,
								Width:  595.2755737304688,
								Height: 841.8897705078125,
							}),
						))
					})
				})
			})
		})
	})

	Context("multiple PDF files", func() {
		var doc references.FPDF_DOCUMENT
		var doc2 references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/test.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document

			pdfData2, err := ioutil.ReadFile(TestDataPath + "/testdata/test_multipage.pdf")
			Expect(err).To(BeNil())

			newDoc2, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData2,
			})
			Expect(err).To(BeNil())

			doc2 = newDoc2.Document
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))

			FPDF_CloseDocument, err = PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc2,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		When("are opened", func() {
			Context("when another document is rendered after the first one", func() {
				It("returns the right image, point to pixel ratio and resolution", func() {
					renderedPage, err := PdfiumInstance.RenderPagesInPixels(&requests.RenderPagesInPixels{
						Pages: []requests.RenderPageInPixels{
							{

								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
								Width:  2000,
								Height: 2000,
							},
							{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc2,
										Index:    0,
									},
								},
								Width:  2000,
								Height: 2000,
							},
							{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc2,
										Index:    1,
									},
								},
								Width:  2000,
								Height: 2000,
							},
						},
						Padding: 50,
					})

					Expect(err).To(BeNil())
					compareRenderHashForPages(&renderedPage.Result, Or(Equal(&responses.RenderPages{
						Width:  1415,
						Height: 6100,
						Pages: []responses.RenderPagesPage{
							{
								PointToPixelRatio: 2.375608084404265,
								Width:             1415,
								Height:            2000,
								X:                 0,
								Y:                 0,
							},
							{
								PointToPixelRatio: 2.375608084404265,
								Width:             1415,
								Height:            2000,
								X:                 0,
								Y:                 2050,
							},
							{
								Page:              1,
								PointToPixelRatio: 2.375608084404265,
								Width:             1415,
								Height:            2000,
								X:                 0,
								Y:                 4100,
							},
						},
					}), Equal(&responses.RenderPages{
						Width:  1415,
						Height: 6100,
						Pages: []responses.RenderPagesPage{
							{
								Page:              0,
								PointToPixelRatio: 2.375607912177905,
								Width:             1415,
								Height:            2000,
								X:                 0,
								Y:                 0,
								HasTransparency:   false,
							},
							{
								Page:              0,
								PointToPixelRatio: 2.375607912177905,
								Width:             1415,
								Height:            2000,
								X:                 0,
								Y:                 2050,
								HasTransparency:   false,
							},
							{
								Page:              1,
								PointToPixelRatio: 2.375607912177905,
								Width:             1415,
								Height:            2000,
								X:                 0,
								Y:                 4100,
								HasTransparency:   false,
							},
						},
					})), TestDataPath+"/testdata/render_"+TestType+"_testpdf_multiple_pdf_combined", TestDataPath+"/testdata/render_"+TestType+"_testpdf_multiple_pdf_combined_7019")
					Expect(renderedPage.Result.Image.Bounds().Size().X).To(Equal(1415))
					Expect(renderedPage.Result.Image.Bounds().Size().Y).To(Equal(6100))
					renderedPage.Cleanup()
				})
			})
		})
	})
})

func compareRenderHash(renderedPage *responses.RenderPage, matcher types.GomegaMatcher, testNames ...string) {
	err := writePrerenderedImage(renderedPage.Image, testNames...)
	Expect(err).To(BeNil())

	// Copy object so we can skip Image.
	// For the image we compare the file hash.
	copiedPage := &responses.RenderPage{
		Page:              renderedPage.Page,
		PointToPixelRatio: renderedPage.PointToPixelRatio,
		Width:             renderedPage.Width,
		Height:            renderedPage.Height,
	}
	Expect(copiedPage).To(matcher)
	existingFileHashes := []types.GomegaMatcher{}

	for _, testName := range testNames {
		existingFileHash, err := ioutil.ReadFile(testName + ".hash")
		Expect(err).To(BeNil())
		existingFileHashes = append(existingFileHashes, Equal(string(existingFileHash)))
	}

	hasher := sha256.New()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(renderedPage.Image.Pix)
	Expect(err).To(BeNil())

	hasher.Write(buf.Bytes())
	currentHash := fmt.Sprintf("%x", hasher.Sum(nil))

	Expect(currentHash).To(Or(existingFileHashes...))
}

func compareRenderHashForPages(renderedPages *responses.RenderPages, matcher types.GomegaMatcher, testNames ...string) {
	err := writePrerenderedImage(renderedPages.Image, testNames...)
	Expect(err).To(BeNil())

	// Copy object so we can skip Image.
	// For the image we compare the file hash.
	copiedPage := &responses.RenderPages{
		Pages:  renderedPages.Pages,
		Width:  renderedPages.Width,
		Height: renderedPages.Height,
	}

	Expect(copiedPage).To(matcher)

	existingFileHashes := []types.GomegaMatcher{}

	for _, testName := range testNames {
		existingFileHash, err := ioutil.ReadFile(testName + ".hash")
		Expect(err).To(BeNil())
		existingFileHashes = append(existingFileHashes, Equal(string(existingFileHash)))
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(renderedPages.Image.Pix)
	Expect(err).To(BeNil())

	hasher := sha256.New()
	hasher.Write(buf.Bytes())
	currentHash := fmt.Sprintf("%x", hasher.Sum(nil))
	Expect(currentHash).To(Or(existingFileHashes...))
}

func compareFileHash(request *requests.RenderToFile, renderedFile *responses.RenderToFile, matcher types.GomegaMatcher, testNames ...string) {
	err := writePrerenderedFile(request, renderedFile, testNames...)
	Expect(err).To(BeNil())

	// Copy object so we can skip Image.
	// For the image we compare the file hash.
	copiedFile := &responses.RenderToFile{
		Pages:             renderedFile.Pages,
		PointToPixelRatio: renderedFile.PointToPixelRatio,
		Width:             renderedFile.Width,
		Height:            renderedFile.Height,
	}
	Expect(copiedFile).To(matcher)

	existingFileHashes := []types.GomegaMatcher{}

	for _, testName := range testNames {
		existingFileHash, err := ioutil.ReadFile(testName + ".hash")
		Expect(err).To(BeNil())
		existingFileHashes = append(existingFileHashes, Equal(string(existingFileHash)))
	}

	hasher := sha256.New()

	if request.OutputTarget == requests.RenderToFileOutputTargetBytes {
		hasher.Write(*renderedFile.ImageBytes)
	} else if request.OutputTarget == requests.RenderToFileOutputTargetFile {
		if request.TargetFilePath != "" {
			Expect(request.TargetFilePath).To(Equal(renderedFile.ImagePath))
		} else {
			// Cleanup tmp file.
			defer os.Remove(renderedFile.ImagePath)
		}
		fileContent, err := ioutil.ReadFile(renderedFile.ImagePath)
		Expect(err).To(BeNil())
		hasher.Write(fileContent)
	}

	currentHash := fmt.Sprintf("%x", hasher.Sum(nil))
	Expect(currentHash).To(Or(existingFileHashes...))

	for _, testName := range testNames {
		existingFileHash, err := ioutil.ReadFile(testName + ".hash")
		Expect(err).To(BeNil())

		if strings.Contains(testName, "_single_") {
			// Compare the single variant to the multi variant.
			existingMultiFileHash, err := ioutil.ReadFile(strings.Replace(testName, "_single_", "_multi_", 1) + ".hash")
			Expect(err).To(BeNil())
			Expect(string(existingMultiFileHash)).To(Equal(string(existingFileHash)))

			// Compare the single variant to the webassembly variant.
			// @todo: figure out why webassembly renders have a different hash.
			//existingWebassemblyFileHash, err := ioutil.ReadFile(strings.Replace(testName, "_single_", "_webassembly_", 1) + ".hash")
			//Expect(err).To(BeNil())
			//Expect(string(existingWebassemblyFileHash)).To(Equal(existingFileHash))

			// Compare the single variant to the internal variant.
			existingInternalFileHash, err := ioutil.ReadFile(strings.Replace(testName, "_single_", "_internal_", 1) + ".hash")
			Expect(err).To(BeNil())
			Expect(string(existingInternalFileHash)).To(Equal(string(existingFileHash)))
		} else if strings.Contains(testName, "_multi_") {
			// Compare the multi variant to the single variant.
			existingSingleFileHash, err := ioutil.ReadFile(strings.Replace(testName, "_multi_", "_single_", 1) + ".hash")
			Expect(err).To(BeNil())
			Expect(string(existingSingleFileHash)).To(Equal(string(existingFileHash)))

			// Compare the multi variant to the webassembly variant.
			// @todo: figure out why webassembly renders have a different hash.
			//existingWebassemblyFileHash, err := ioutil.ReadFile(strings.Replace(testName, "_multi_", "_webassembly_", 1) + ".hash")
			//Expect(err).To(BeNil())
			//Expect(string(existingWebassemblyFileHash)).To(Equal(existingFileHash))

			// Compare the multi variant to the internal variant.
			existingInternalFileHash, err := ioutil.ReadFile(strings.Replace(testName, "_multi_", "_internal_", 1) + ".hash")
			Expect(err).To(BeNil())
			Expect(string(existingInternalFileHash)).To(Equal(string(existingFileHash)))
		} else if strings.Contains(testName, "_internal_") {
			// Compare the internal variant to the single variant.
			existingSingleFileHash, err := ioutil.ReadFile(strings.Replace(testName, "_internal_", "_single_", 1) + ".hash")
			Expect(err).To(BeNil())
			Expect(string(existingSingleFileHash)).To(Equal(string(existingFileHash)))

			// Compare the internal variant to the webassembly variant.
			// @todo: figure out why webassembly renders have a different hash.
			//existingWebassemblyFileHash, err := ioutil.ReadFile(strings.Replace(testName, "_internal_", "_webassembly_", 1) + ".hash")
			//Expect(err).To(BeNil())
			//Expect(string(existingWebassemblyFileHash)).To(Equal(string(existingFileHash)))

			// Compare the internal variant to the multi variant.
			existingMultiFileHash, err := ioutil.ReadFile(strings.Replace(testName, "_internal_", "_multi_", 1) + ".hash")
			Expect(err).To(BeNil())
			Expect(string(existingMultiFileHash)).To(Equal(string(existingFileHash)))
		} else if strings.Contains(testName, "_webassembly_") {
			// @todo: figure out why webassembly renders have a different hash.
			// Compare the webassembly variant to the single variant.
			//existingSingleFileHash, err := ioutil.ReadFile(strings.Replace(testName, "_webassembly_", "_single_", 1) + ".hash")
			//Expect(err).To(BeNil())
			//Expect(string(existingSingleFileHash)).To(Equal(string(existingFileHash)))

			// Compare the webassembly variant to the multi variant.
			//existingMultiFileHash, err := ioutil.ReadFile(strings.Replace(testName, "_webassembly_", "_multi_", 1) + ".hash")
			//Expect(err).To(BeNil())
			//Expect(string(existingMultiFileHash)).To(Equal(string(existingFileHash)))

			// Compare the webassembly variant to the internal variant.
			//existingInternalFileHash, err := ioutil.ReadFile(strings.Replace(testName, "_webassembly_", "internal", 1) + ".hash")
			//Expect(err).To(BeNil())
			//Expect(string(existingInternalFileHash)).To(Equal(string(existingFileHash)))
		}
	}
}

func writePrerenderedImage(renderedImage *image.RGBA, testNames ...string) error {
	filename := testNames[len(testNames)-1]
	if _, err := os.Stat(filename + ".hash"); err == nil {
		return nil // Comment this in case of updating PDFium versions and rendering has changed.
	}

	// Be sure to validate the difference in image to ensure rendering has not been broken.
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(renderedImage.Pix); err != nil {
		return err
	}

	hasher := sha256.New()
	hasher.Write(buf.Bytes())
	currentHash := fmt.Sprintf("%x", hasher.Sum(nil))

	if err := ioutil.WriteFile(filename+".hash", []byte(currentHash), 0777); err != nil {
		return err
	}

	f, err := os.Create(filename + ".png")
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, renderedImage)
	if err != nil {
		return err
	}

	return nil
}

func writePrerenderedFile(request *requests.RenderToFile, renderedFile *responses.RenderToFile, testNames ...string) error {
	filename := testNames[len(testNames)-1]
	if _, err := os.Stat(filename + ".hash"); err == nil {
		return nil // Comment this in case of updating PDFium versions and rendering has changed.
	}

	var fileBytes []byte

	hasher := sha256.New()

	if request.OutputTarget == requests.RenderToFileOutputTargetBytes {
		hasher.Write(*renderedFile.ImageBytes)
		fileBytes = *renderedFile.ImageBytes
	} else if request.OutputTarget == requests.RenderToFileOutputTargetFile {
		fileContent, err := ioutil.ReadFile(renderedFile.ImagePath)
		if err != nil {
			return err
		}

		hasher.Write(fileContent)
		fileBytes = fileContent
	}

	currentHash := fmt.Sprintf("%x", hasher.Sum(nil))

	if err := ioutil.WriteFile(filename+".hash", []byte(currentHash), 0777); err != nil {
		return err
	}

	if request.OutputFormat == requests.RenderToFileOutputFormatPNG {
		filename += ".png"
	} else if request.OutputFormat == requests.RenderToFileOutputFormatJPG {
		filename += ".jpg"
	}

	err := ioutil.WriteFile(filename, fileBytes, 0777)
	if err != nil {
		return err
	}

	return nil
}
