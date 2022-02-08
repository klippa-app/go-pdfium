package shared_tests

import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_ppo", func() {
	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when calling FPDF_ImportPages", func() {
				FPDF_ImportPages, err := PdfiumInstance.FPDF_ImportPages(&requests.FPDF_ImportPages{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_ImportPages).To(BeNil())
			})

			It("returns an error when calling FPDF_CopyViewerPreferences", func() {
				FPDF_CopyViewerPreferences, err := PdfiumInstance.FPDF_CopyViewerPreferences(&requests.FPDF_CopyViewerPreferences{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_CopyViewerPreferences).To(BeNil())
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
			When("no destination document is given", func() {
				It("returns an error when calling FPDF_ImportPages", func() {
					FPDF_ImportPages, err := PdfiumInstance.FPDF_ImportPages(&requests.FPDF_ImportPages{
						Source: doc,
					})
					Expect(err).To(MatchError("document not given"))
					Expect(FPDF_ImportPages).To(BeNil())
				})

				It("returns an error when calling FPDF_CopyViewerPreferences", func() {
					FPDF_CopyViewerPreferences, err := PdfiumInstance.FPDF_CopyViewerPreferences(&requests.FPDF_CopyViewerPreferences{
						Source: doc,
					})
					Expect(err).To(MatchError("document not given"))
					Expect(FPDF_CopyViewerPreferences).To(BeNil())
				})
			})

			Context("a second PDF file is opened", func() {
				var doc2 references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/viewer_ref.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data: &pdfData,
					})
					Expect(err).To(BeNil())

					doc2 = newDoc.Document
				})

				AfterEach(func() {
					FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
						Document: doc2,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_CloseDocument).To(Not(BeNil()))
				})

				When("is opened", func() {
					It("returns no error when FPDF_ImportPages is called", func() {
						FPDF_ImportPages, err := PdfiumInstance.FPDF_ImportPages(&requests.FPDF_ImportPages{
							Source:      doc2,
							Destination: doc,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_ImportPages).To(Not(BeNil()))
					})

					It("returns no error when FPDF_ImportPages is called with a valid pagerange", func() {
						pageRange := "1"
						FPDF_ImportPages, err := PdfiumInstance.FPDF_ImportPages(&requests.FPDF_ImportPages{
							Source:      doc2,
							Destination: doc,
							PageRange:   &pageRange,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_ImportPages).To(Not(BeNil()))
					})

					It("returns no error when FPDF_ImportPages is called with an invalid pagerange", func() {
						pageRange := "32"
						FPDF_ImportPages, err := PdfiumInstance.FPDF_ImportPages(&requests.FPDF_ImportPages{
							Source:      doc2,
							Destination: doc,
							PageRange:   &pageRange,
						})
						Expect(err).To(MatchError("import of pages failed"))
						Expect(FPDF_ImportPages).To(BeNil())
					})

					It("returns an error when calling FPDF_CopyViewerPreferences with a source document that has no viewer preferences", func() {
						FPDF_CopyViewerPreferences, err := PdfiumInstance.FPDF_CopyViewerPreferences(&requests.FPDF_CopyViewerPreferences{
							Source:      doc,
							Destination: doc2,
						})
						Expect(err).To(MatchError("copy of viewer preferences failed"))
						Expect(FPDF_CopyViewerPreferences).To(BeNil())
					})

					It("returns no error when calling FPDF_CopyViewerPreferences with a source document that has viewer preferences", func() {
						FPDF_CopyViewerPreferences, err := PdfiumInstance.FPDF_CopyViewerPreferences(&requests.FPDF_CopyViewerPreferences{
							Source:      doc2,
							Destination: doc,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_CopyViewerPreferences).To(Not(BeNil()))
					})
				})
			})
		})
	})
})
