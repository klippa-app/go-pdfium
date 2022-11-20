package shared_tests

import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_searchex", func() {
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

	Context("no text page is given", func() {
		It("returns an error when FPDFText_GetCharIndexFromTextIndex is called", func() {
			FPDFText_GetCharIndexFromTextIndex, err := PdfiumInstance.FPDFText_GetCharIndexFromTextIndex(&requests.FPDFText_GetCharIndexFromTextIndex{})
			Expect(err).To(MatchError("textPage not given"))
			Expect(FPDFText_GetCharIndexFromTextIndex).To(BeNil())
		})

		It("returns an error when FPDFText_GetTextIndexFromCharIndex is called", func() {
			FPDFText_GetTextIndexFromCharIndex, err := PdfiumInstance.FPDFText_GetTextIndexFromCharIndex(&requests.FPDFText_GetTextIndexFromCharIndex{})
			Expect(err).To(MatchError("textPage not given"))
			Expect(FPDFText_GetTextIndexFromCharIndex).To(BeNil())
		})
	})

	Context("a normal PDF file", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/hello_world.pdf")
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

		When("a text page is opened", func() {
			var textPage references.FPDF_TEXTPAGE

			BeforeEach(func() {
				textPageResp, err := PdfiumInstance.FPDFText_LoadPage(&requests.FPDFText_LoadPage{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(textPageResp).To(Not(BeNil()))

				textPage = textPageResp.TextPage
			})

			It("returns the correct char index from text index", func() {
				FPDFText_GetCharIndexFromTextIndex, err := PdfiumInstance.FPDFText_GetCharIndexFromTextIndex(&requests.FPDFText_GetCharIndexFromTextIndex{
					TextPage:   textPage,
					NTextIndex: 29,
				})
				Expect(err).To(BeNil())
				Expect(FPDFText_GetCharIndexFromTextIndex).To(Equal(&responses.FPDFText_GetCharIndexFromTextIndex{
					CharIndex: 29,
				}))
			})

			It("returns an error when calling char index from text index with an invalid index", func() {
				FPDFText_GetCharIndexFromTextIndex, err := PdfiumInstance.FPDFText_GetCharIndexFromTextIndex(&requests.FPDFText_GetCharIndexFromTextIndex{
					TextPage:   textPage,
					NTextIndex: 300,
				})
				Expect(err).To(MatchError("could not get char index"))
				Expect(FPDFText_GetCharIndexFromTextIndex).To(BeNil())
			})

			It("returns the correct text index from char index", func() {
				FPDFText_GetCharIndexFromTextIndex, err := PdfiumInstance.FPDFText_GetTextIndexFromCharIndex(&requests.FPDFText_GetTextIndexFromCharIndex{
					TextPage:   textPage,
					NCharIndex: 29,
				})
				Expect(err).To(BeNil())
				Expect(FPDFText_GetCharIndexFromTextIndex).To(Equal(&responses.FPDFText_GetTextIndexFromCharIndex{
					TextIndex: 29,
				}))
			})

			It("returns an error when calling text index from char index with an invalid index", func() {
				FPDFText_GetTextIndexFromCharIndex, err := PdfiumInstance.FPDFText_GetTextIndexFromCharIndex(&requests.FPDFText_GetTextIndexFromCharIndex{
					TextPage:   textPage,
					NCharIndex: 300,
				})
				Expect(err).To(MatchError("could not get text index"))
				Expect(FPDFText_GetTextIndexFromCharIndex).To(BeNil())
			})
		})
	})
})
