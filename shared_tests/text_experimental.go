//go:build pdfium_experimental
// +build pdfium_experimental

package shared_tests

import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("text", func() {
	BeforeEach(func() {
		Locker.Lock()

		if TestType == "webassembly" {
			// @todo: remove me when implemented.
			Skip("This test is skipped on Webassembly")
		}
	})

	AfterEach(func() {
		Locker.Unlock()

		if TestType == "webassembly" {
			// @todo: remove me when implemented.
			Skip("This test is skipped on Webassembly")
		}
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
			Context("when the structured page text is requested", func() {
				Context("when PixelPositions is enabled", func() {
					It("returns the correct font information", func() {
						pageTextStructured, err := PdfiumInstance.GetPageTextStructured(&requests.GetPageTextStructured{
							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    0,
								},
							},
							CollectFontInformation: true,
						})
						Expect(err).To(BeNil())
						Expect(pageTextStructured).To(Equal(loadStructuredText(TestDataPath+"/testdata/text_experimental_"+TestType+"_testpdf_experimental_with_font_information.json", pageTextStructured)))
					})

					Context("and PixelPositions is enabled", func() {
						It("returns the correct font information", func() {
							pageTextStructured, err := PdfiumInstance.GetPageTextStructured(&requests.GetPageTextStructured{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
								CollectFontInformation: true,
								PixelPositions: requests.GetPageTextStructuredPixelPositions{
									Calculate: true,
									Width:     3000,
									Height:    3000,
								},
							})
							Expect(err).To(BeNil())
							Expect(pageTextStructured).To(Equal(loadStructuredText(TestDataPath+"/testdata/text_experimental_"+TestType+"_testpdf_experimental_with_font_information_and_pixel_positions.json", pageTextStructured)))
						})
					})
				})
			})
		})
	})
})
