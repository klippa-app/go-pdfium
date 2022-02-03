package shared_tests

import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func RunfpdfEditTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("fpdf_edit", func() {
		Context("no document", func() {
			When("is opened", func() {
				It("returns an error when getting the pdf page rotation", func() {
					FPDFPage_GetRotation, err := pdfiumContainer.FPDFPage_GetRotation(&requests.FPDFPage_GetRotation{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Index: 0,
							},
						},
					})
					Expect(err).To(MatchError("document not given"))
					Expect(FPDFPage_GetRotation).To(BeNil())
				})

				It("returns an error when setting the pdf page rotation", func() {
					FPDFPage_SetRotation, err := pdfiumContainer.FPDFPage_SetRotation(&requests.FPDFPage_SetRotation{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Index: 0,
							},
						},
					})
					Expect(err).To(MatchError("document not given"))
					Expect(FPDFPage_SetRotation).To(BeNil())
				})

				It("returns an error when getting the pdf page transparency", func() {
					FPDFPage_HasTransparency, err := pdfiumContainer.FPDFPage_HasTransparency(&requests.FPDFPage_HasTransparency{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Index: 0,
							},
						},
					})
					Expect(err).To(MatchError("document not given"))
					Expect(FPDFPage_HasTransparency).To(BeNil())
				})
			})
		})

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
				Context("when the rotation is changed", func() {
					It("it is being returned by FPDFPage_GetRotation", func() {
						FPDFPage_SetRotation, err := pdfiumContainer.FPDFPage_SetRotation(&requests.FPDFPage_SetRotation{
							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    0,
								},
							},
							Rotate: enums.FPDF_PAGE_ROTATION_180_CW,
						})
						Expect(err).To(BeNil())
						Expect(FPDFPage_SetRotation).To(Equal(&responses.FPDFPage_SetRotation{}))

						FPDFPage_GetRotation, err := pdfiumContainer.FPDFPage_GetRotation(&requests.FPDFPage_GetRotation{
							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    0,
								},
							},
						})
						Expect(err).To(BeNil())
						Expect(FPDFPage_GetRotation).To(Equal(&responses.FPDFPage_GetRotation{
							PageRotation: enums.FPDF_PAGE_ROTATION_180_CW,
						}))
					})
				})
			})
		})
	})
}
