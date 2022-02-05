package shared_tests

import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func RunfpdfCatalogTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("fpdf_catalog", func() {
		Context("no document", func() {
			When("is opened", func() {
				It("returns an error when getting the document tagged status", func() {
					isTagged, err := pdfiumContainer.FPDFCatalog_IsTagged(&requests.FPDFCatalog_IsTagged{})
					Expect(err).To(MatchError("document not given"))
					Expect(isTagged).To(BeNil())
				})
			})
		})

		Context("a normal PDF file", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				pdfData, err := ioutil.ReadFile(testsPath + "/testdata/test.pdf")
				Expect(err).To(BeNil())

				newDoc, err := pdfiumContainer.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
					Data: &pdfData,
				})
				Expect(err).To(BeNil())

				doc = newDoc.Document
			})

			AfterEach(func() {
				FPDF_CloseDocument, err := pdfiumContainer.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_CloseDocument).To(Not(BeNil()))
			})

			When("is opened", func() {
				It("returns that its not tagged", func() {
					isTagged, err := pdfiumContainer.FPDFCatalog_IsTagged(&requests.FPDFCatalog_IsTagged{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(isTagged).To(Equal(&responses.FPDFCatalog_IsTagged{
						IsTagged: false,
					}))
				})
			})
		})

		Context("a tagged PDF file", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				pdfData, err := ioutil.ReadFile(testsPath + "/testdata/tagged_table.pdf")
				Expect(err).To(BeNil())

				newDoc, err := pdfiumContainer.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
					Data: &pdfData,
				})
				Expect(err).To(BeNil())

				doc = newDoc.Document
			})

			AfterEach(func() {
				FPDF_CloseDocument, err := pdfiumContainer.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_CloseDocument).To(Not(BeNil()))
			})

			When("is opened", func() {
				It("returns that it is tagged", func() {
					isTagged, err := pdfiumContainer.FPDFCatalog_IsTagged(&requests.FPDFCatalog_IsTagged{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(isTagged).To(Equal(&responses.FPDFCatalog_IsTagged{
						IsTagged: true,
					}))
				})
			})
		})
	})
}
