//go:build pdfium_experimental
// +build pdfium_experimental

package shared_tests

import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_thumbnail", func() {
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

	Context("no page is given", func() {
		It("returns an error when getting the decoded thumbnail data", func() {
			FPDFPage_GetDecodedThumbnailData, err := PdfiumInstance.FPDFPage_GetDecodedThumbnailData(&requests.FPDFPage_GetDecodedThumbnailData{})
			Expect(err).To(MatchError("either page reference or index should be given"))
			Expect(FPDFPage_GetDecodedThumbnailData).To(BeNil())
		})

		It("returns an error when getting the raw thumbnail data", func() {
			FPDFPage_GetRawThumbnailData, err := PdfiumInstance.FPDFPage_GetRawThumbnailData(&requests.FPDFPage_GetRawThumbnailData{})
			Expect(err).To(MatchError("either page reference or index should be given"))
			Expect(FPDFPage_GetRawThumbnailData).To(BeNil())
		})

		It("returns an error when getting the decoded thumbnail data", func() {
			FPDFPage_GetThumbnailAsBitmap, err := PdfiumInstance.FPDFPage_GetThumbnailAsBitmap(&requests.FPDFPage_GetThumbnailAsBitmap{})
			Expect(err).To(MatchError("either page reference or index should be given"))
			Expect(FPDFPage_GetThumbnailAsBitmap).To(BeNil())
		})
	})

	Context("a normal PDF file without thumbnails", func() {
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

		It("returns no decoded thumbnail data", func() {
			FPDFPage_GetDecodedThumbnailData, err := PdfiumInstance.FPDFPage_GetDecodedThumbnailData(&requests.FPDFPage_GetDecodedThumbnailData{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(FPDFPage_GetDecodedThumbnailData).To(Equal(&responses.FPDFPage_GetDecodedThumbnailData{}))
		})

		It("returns no raw thumbnail data", func() {
			FPDFPage_GetRawThumbnailData, err := PdfiumInstance.FPDFPage_GetRawThumbnailData(&requests.FPDFPage_GetRawThumbnailData{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(FPDFPage_GetRawThumbnailData).To(Equal(&responses.FPDFPage_GetRawThumbnailData{}))
		})

		It("returns no decoded thumbnail data", func() {
			FPDFPage_GetThumbnailAsBitmap, err := PdfiumInstance.FPDFPage_GetThumbnailAsBitmap(&requests.FPDFPage_GetThumbnailAsBitmap{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(FPDFPage_GetThumbnailAsBitmap).To(Equal(&responses.FPDFPage_GetThumbnailAsBitmap{}))
		})
	})

	Context("a PDF file with a thumbnail", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/simple_thumbnail.pdf")
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

		It("returns decoded thumbnail data", func() {
			FPDFPage_GetDecodedThumbnailData, err := PdfiumInstance.FPDFPage_GetDecodedThumbnailData(&requests.FPDFPage_GetDecodedThumbnailData{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(FPDFPage_GetDecodedThumbnailData).To(Not(BeNil()))
			Expect(FPDFPage_GetDecodedThumbnailData.Thumbnail).To(HaveLen(1138))
		})

		It("returns raw thumbnail data", func() {
			FPDFPage_GetRawThumbnailData, err := PdfiumInstance.FPDFPage_GetRawThumbnailData(&requests.FPDFPage_GetRawThumbnailData{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(FPDFPage_GetRawThumbnailData).To(Not(BeNil()))
			Expect(FPDFPage_GetRawThumbnailData.RawThumbnail).To(HaveLen(1851))
		})

		It("returns thumbnail as bitmap", func() {
			By("when the thumbnail is fetched")
			FPDFPage_GetThumbnailAsBitmap, err := PdfiumInstance.FPDFPage_GetThumbnailAsBitmap(&requests.FPDFPage_GetThumbnailAsBitmap{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(FPDFPage_GetThumbnailAsBitmap).To(Not(BeNil()))
			Expect(FPDFPage_GetThumbnailAsBitmap.Bitmap).To(Not(BeNil()))

			By("the first byte in the bitmap buffer is white")
			FPDFBitmap_GetBuffer, err := PdfiumInstance.FPDFBitmap_GetBuffer(&requests.FPDFBitmap_GetBuffer{
				Bitmap: *FPDFPage_GetThumbnailAsBitmap.Bitmap,
			})
			Expect(err).To(BeNil())
			Expect(FPDFBitmap_GetBuffer).To(Not(BeNil()))
			Expect(FPDFBitmap_GetBuffer.Buffer).To(Not(BeNil()))
			Expect(FPDFBitmap_GetBuffer.Buffer[0]).To(Equal(uint8(255)))
		})
	})
})
