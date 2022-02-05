package shared_tests

import (
	"github.com/klippa-app/go-pdfium/responses"
	"io/ioutil"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func RunfpdfThumbnailTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("fpdf_thumbnail", func() {
		Context("no page is given", func() {
			It("returns an error when getting the decoded thumbnail data", func() {
				FPDFPage_GetDecodedThumbnailData, err := pdfiumContainer.FPDFPage_GetDecodedThumbnailData(&requests.FPDFPage_GetDecodedThumbnailData{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_GetDecodedThumbnailData).To(BeNil())
			})

			It("returns an error when getting the raw thumbnail data", func() {
				FPDFPage_GetRawThumbnailData, err := pdfiumContainer.FPDFPage_GetRawThumbnailData(&requests.FPDFPage_GetRawThumbnailData{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_GetRawThumbnailData).To(BeNil())
			})

			It("returns an error when getting the decoded thumbnail data", func() {
				FPDFPage_GetThumbnailAsBitmap, err := pdfiumContainer.FPDFPage_GetThumbnailAsBitmap(&requests.FPDFPage_GetThumbnailAsBitmap{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_GetThumbnailAsBitmap).To(BeNil())
			})
		})

		Context("a normal PDF file without thumbnails", func() {
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

			It("returns no decoded thumbnail data", func() {
				FPDFPage_GetDecodedThumbnailData, err := pdfiumContainer.FPDFPage_GetDecodedThumbnailData(&requests.FPDFPage_GetDecodedThumbnailData{
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
				FPDFPage_GetRawThumbnailData, err := pdfiumContainer.FPDFPage_GetRawThumbnailData(&requests.FPDFPage_GetRawThumbnailData{
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
				FPDFPage_GetThumbnailAsBitmap, err := pdfiumContainer.FPDFPage_GetThumbnailAsBitmap(&requests.FPDFPage_GetThumbnailAsBitmap{
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
				pdfData, err := ioutil.ReadFile(testsPath + "/testdata/simple_thumbnail.pdf")
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

			It("returns no decoded thumbnail data", func() {
				FPDFPage_GetDecodedThumbnailData, err := pdfiumContainer.FPDFPage_GetDecodedThumbnailData(&requests.FPDFPage_GetDecodedThumbnailData{
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

			It("returns no raw thumbnail data", func() {
				FPDFPage_GetRawThumbnailData, err := pdfiumContainer.FPDFPage_GetRawThumbnailData(&requests.FPDFPage_GetRawThumbnailData{
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

			It("returns no decoded thumbnail data", func() {
				FPDFPage_GetThumbnailAsBitmap, err := pdfiumContainer.FPDFPage_GetThumbnailAsBitmap(&requests.FPDFPage_GetThumbnailAsBitmap{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetThumbnailAsBitmap).To(Not(BeNil()))
				// @todo: render thumbnail when FPDFBitmap_* is implemented and compare hash.
			})
		})
	})
}
