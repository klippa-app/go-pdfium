package shared_tests

import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_edit", func() {
	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when getting the pdf page rotation", func() {
				FPDFPage_GetRotation, err := PdfiumInstance.FPDFPage_GetRotation(&requests.FPDFPage_GetRotation{
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
				FPDFPage_SetRotation, err := PdfiumInstance.FPDFPage_SetRotation(&requests.FPDFPage_SetRotation{
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
				FPDFPage_HasTransparency, err := PdfiumInstance.FPDFPage_HasTransparency(&requests.FPDFPage_HasTransparency{
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
			Context("when the rotation is changed", func() {
				It("it is being returned by FPDFPage_GetRotation", func() {
					FPDFPage_SetRotation, err := PdfiumInstance.FPDFPage_SetRotation(&requests.FPDFPage_SetRotation{
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

					FPDFPage_GetRotation, err := PdfiumInstance.FPDFPage_GetRotation(&requests.FPDFPage_GetRotation{
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
