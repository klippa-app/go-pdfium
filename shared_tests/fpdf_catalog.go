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
