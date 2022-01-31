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
			var doc references.Document

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
				err := pdfiumContainer.CloseDocument(doc)
				Expect(err).To(BeNil())
			})

			When("is opened", func() {
				Context("when an invalid page is given", func() {
					Context("GetPageRotation()", func() {
						It("returns an error", func() {
							pageRotation, err := pdfiumContainer.GetPageRotation(&requests.GetPageRotation{
								Document: doc,
								Page: requests.Page{
									Index: 1,
								},
							})
							Expect(err).To(MatchError(errors.ErrPage.Error()))
							Expect(pageRotation).To(BeNil())
						})
					})

					Context("GetPageTransparency()", func() {
						It("returns an error", func() {
							pageTransparency, err := pdfiumContainer.GetPageTransparency(&requests.GetPageTransparency{
								Document: doc,
								Page: requests.Page{
									Index: 1,
								},
							})
							Expect(err).To(MatchError(errors.ErrPage.Error()))
							Expect(pageTransparency).To(BeNil())
						})
					})

					Context("FlattenPage()", func() {
						It("returns an error", func() {
							flattenedPage, err := pdfiumContainer.FlattenPage(&requests.FlattenPage{
								Document: doc,
								Page: requests.Page{
									Index: 1,
								},
								Usage: requests.FlattenPageUsageNormalDisplay,
							})
							Expect(err).To(MatchError(errors.ErrPage.Error()))
							Expect(flattenedPage).To(BeNil())
						})
					})
				})

				Context("when the page rotation is requested", func() {
					It("returns the correct rotation", func() {
						rotation, err := pdfiumContainer.GetPageRotation(&requests.GetPageRotation{
							Document: doc,
							Page: requests.Page{
								Index: 0,
							},
						})
						Expect(err).To(BeNil())
						Expect(rotation).To(Equal(&responses.GetPageRotation{}))
					})
				})

				Context("when the page transparency is requested", func() {
					It("returns the correct transparency", func() {
						pageTransparency, err := pdfiumContainer.GetPageTransparency(&requests.GetPageTransparency{
							Document: doc,
							Page: requests.Page{
								Index: 0,
							},
						})
						Expect(err).To(BeNil())
						Expect(pageTransparency).To(Equal(&responses.GetPageTransparency{
							HasTransparency: false,
						}))
					})
				})

				Context("when the page flattening is requested", func() {
					It("returns that the page does not need to be flattened", func() {
						pageFlattenResult, err := pdfiumContainer.FlattenPage(&requests.FlattenPage{
							Document: doc,
							Page: requests.Page{
								Index: 0,
							},
							Usage: requests.FlattenPageUsageNormalDisplay,
						})
						Expect(err).To(BeNil())
						Expect(pageFlattenResult).To(Equal(&responses.FlattenPage{
							Result: responses.FlattenPageResultNothingToDo,
						}))
					})
				})
			})
		})

		Context("a PDF file that uses an alpha channel", func() {
			var doc references.Document

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
				err := pdfiumContainer.CloseDocument(doc)
				Expect(err).To(BeNil())
			})

			When("the page transparency is requested", func() {
				It("returns the correct transparency", func() {
					pageTransparency, err := pdfiumContainer.GetPageTransparency(&requests.GetPageTransparency{
						Document: doc,
						Page: requests.Page{
							Index: 0,
						},
					})
					Expect(err).To(BeNil())
					Expect(pageTransparency).To(Equal(&responses.GetPageTransparency{
						HasTransparency: true,
					}))
				})
			})
		})
	})
}
