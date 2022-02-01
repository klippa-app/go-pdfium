package shared_tests

import (
	"github.com/klippa-app/go-pdfium/references"
	"io/ioutil"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func RunPageTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("Page", func() {
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
					Context("GetPageRotation()", func() {
						It("returns an error", func() {
							pageRotation, err := pdfiumContainer.FPDFPage_GetRotation(&requests.FPDFPage_GetRotation{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    1,
									},
								},
							})
							Expect(err).To(MatchError(errors.ErrPage.Error()))
							Expect(pageRotation).To(BeNil())
						})
					})

					Context("GetPageTransparency()", func() {
						It("returns an error", func() {
							pageTransparency, err := pdfiumContainer.FPDFPage_HasTransparency(&requests.FPDFPage_HasTransparency{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    1,
									},
								},
							})
							Expect(err).To(MatchError(errors.ErrPage.Error()))
							Expect(pageTransparency).To(BeNil())
						})
					})

					Context("FlattenPage()", func() {
						It("returns an error", func() {
							flattenedPage, err := pdfiumContainer.FPDFPage_Flatten(&requests.FPDFPage_Flatten{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    1,
									},
								},
								Usage: requests.FPDFPage_FlattenUsageNormalDisplay,
							})
							Expect(err).To(MatchError(errors.ErrPage.Error()))
							Expect(flattenedPage).To(BeNil())
						})
					})
				})

				Context("when the page rotation is requested", func() {
					It("returns the correct rotation", func() {
						rotation, err := pdfiumContainer.FPDFPage_GetRotation(&requests.FPDFPage_GetRotation{
							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    0,
								},
							},
						})
						Expect(err).To(BeNil())
						Expect(rotation).To(Equal(&responses.FPDFPage_GetRotation{}))
					})
				})

				Context("when the page transparency is requested", func() {
					It("returns the correct transparency", func() {
						pageTransparency, err := pdfiumContainer.FPDFPage_HasTransparency(&requests.FPDFPage_HasTransparency{
							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    0,
								},
							},
						})
						Expect(err).To(BeNil())
						Expect(pageTransparency).To(Equal(&responses.FPDFPage_HasTransparency{
							HasTransparency: false,
						}))
					})
				})

				Context("when the page flattening is requested", func() {
					It("returns that the page does not need to be flattened", func() {
						pageFlattenResult, err := pdfiumContainer.FPDFPage_Flatten(&requests.FPDFPage_Flatten{
							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    0,
								},
							},
							Usage: requests.FPDFPage_FlattenUsageNormalDisplay,
						})
						Expect(err).To(BeNil())
						Expect(pageFlattenResult).To(Equal(&responses.FPDFPage_Flatten{
							Result: responses.FPDFPage_FlattenResultNothingToDo,
						}))
					})
				})
			})
		})

		Context("a PDF file that uses an alpha channel", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				pdfData, err := ioutil.ReadFile(testsPath + "/testdata/alpha_channel.pdf")
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

			When("the page transparency is requested", func() {
				It("returns the correct transparency", func() {
					pageTransparency, err := pdfiumContainer.FPDFPage_HasTransparency(&requests.FPDFPage_HasTransparency{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})
					Expect(err).To(BeNil())
					Expect(pageTransparency).To(Equal(&responses.FPDFPage_HasTransparency{
						HasTransparency: true,
					}))
				})
			})
		})
	})
}
