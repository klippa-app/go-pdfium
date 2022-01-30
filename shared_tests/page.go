package shared_tests

import (
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
			var doc pdfium.Document

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

				doc = newDoc
			})

			AfterEach(func() {
				doc.Close()
			})

			When("is opened", func() {
				Context("when an invalid page is given", func() {
					Context("GetPageRotation()", func() {
						It("returns an error", func() {
							pageRotation, err := doc.GetPageRotation(&requests.GetPageRotation{
								Page: 1,
							})
							Expect(err).To(MatchError(errors.ErrPage.Error()))
							Expect(pageRotation).To(BeNil())
						})
					})

					Context("GetPageTransparency()", func() {
						It("returns an error", func() {
							pageTransparency, err := doc.GetPageTransparency(&requests.GetPageTransparency{
								Page: 1,
							})
							Expect(err).To(MatchError(errors.ErrPage.Error()))
							Expect(pageTransparency).To(BeNil())
						})
					})

					Context("FlattenPage()", func() {
						It("returns an error", func() {
							flattenedPage, err := doc.FlattenPage(&requests.FlattenPage{
								Page:  1,
								Usage: requests.FlattenPageUsageNormalDisplay,
							})
							Expect(err).To(MatchError(errors.ErrPage.Error()))
							Expect(flattenedPage).To(BeNil())
						})
					})
				})

				Context("when the page rotation is requested", func() {
					It("returns the correct rotation", func() {
						rotation, err := doc.GetPageRotation(&requests.GetPageRotation{
							Page: 0,
						})
						Expect(err).To(BeNil())
						Expect(rotation).To(Equal(&responses.GetPageRotation{}))
					})
				})

				Context("when the page transparency is requested", func() {
					It("returns the correct transparency", func() {
						pageTransparency, err := doc.GetPageTransparency(&requests.GetPageTransparency{
							Page: 0,
						})
						Expect(err).To(BeNil())
						Expect(pageTransparency).To(Equal(&responses.GetPageTransparency{
							HasTransparency: false,
						}))
					})
				})

				Context("when the page flattening is requested", func() {
					It("returns that the page does not need to be flattened", func() {
						pageFlattenResult, err := doc.FlattenPage(&requests.FlattenPage{
							Page:  0,
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
			var doc pdfium.Document

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

				doc = newDoc
			})

			AfterEach(func() {
				doc.Close()
			})

			When("the page transparency is requested", func() {
				It("returns the correct transparency", func() {
					pageTransparency, err := doc.GetPageTransparency(&requests.GetPageTransparency{
						Page: 0,
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
