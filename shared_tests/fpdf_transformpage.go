package shared_tests

import "C"
import (
	"github.com/klippa-app/go-pdfium/structs"
	"os"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_transformpage", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	Context("no page", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFPage_SetMediaBox", func() {
				FPDFPage_SetMediaBox, err := PdfiumInstance.FPDFPage_SetMediaBox(&requests.FPDFPage_SetMediaBox{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_SetMediaBox).To(BeNil())
			})

			It("returns an error when calling FPDFPage_SetCropBox", func() {
				FPDFPage_SetCropBox, err := PdfiumInstance.FPDFPage_SetCropBox(&requests.FPDFPage_SetCropBox{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_SetCropBox).To(BeNil())
			})

			It("returns an error when calling FPDFPage_SetBleedBox", func() {
				FPDFPage_SetBleedBox, err := PdfiumInstance.FPDFPage_SetBleedBox(&requests.FPDFPage_SetBleedBox{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_SetBleedBox).To(BeNil())
			})

			It("returns an error when calling FPDFPage_SetTrimBox", func() {
				FPDFPage_SetTrimBox, err := PdfiumInstance.FPDFPage_SetTrimBox(&requests.FPDFPage_SetTrimBox{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_SetTrimBox).To(BeNil())
			})

			It("returns an error when calling FPDFPage_SetArtBox", func() {
				FPDFPage_SetArtBox, err := PdfiumInstance.FPDFPage_SetArtBox(&requests.FPDFPage_SetArtBox{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_SetArtBox).To(BeNil())
			})

			It("returns an error when calling FPDFPage_GetMediaBox", func() {
				FPDFPage_GetMediaBox, err := PdfiumInstance.FPDFPage_GetMediaBox(&requests.FPDFPage_GetMediaBox{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_GetMediaBox).To(BeNil())
			})

			It("returns an error when calling FPDFPage_GetCropBox", func() {
				FPDFPage_GetCropBox, err := PdfiumInstance.FPDFPage_GetCropBox(&requests.FPDFPage_GetCropBox{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_GetCropBox).To(BeNil())
			})

			It("returns an error when calling FPDFPage_GetBleedBox", func() {
				FPDFPage_GetBleedBox, err := PdfiumInstance.FPDFPage_GetBleedBox(&requests.FPDFPage_GetBleedBox{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_GetBleedBox).To(BeNil())
			})

			It("returns an error when calling FPDFPage_GetTrimBox", func() {
				FPDFPage_GetTrimBox, err := PdfiumInstance.FPDFPage_GetTrimBox(&requests.FPDFPage_GetTrimBox{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_GetTrimBox).To(BeNil())
			})

			It("returns an error when calling FPDFPage_GetArtBox", func() {
				FPDFPage_GetArtBox, err := PdfiumInstance.FPDFPage_GetArtBox(&requests.FPDFPage_GetArtBox{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_GetArtBox).To(BeNil())
			})

			It("returns an error when calling FPDFPage_TransFormWithClip", func() {
				FPDFPage_TransFormWithClip, err := PdfiumInstance.FPDFPage_TransFormWithClip(&requests.FPDFPage_TransFormWithClip{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_TransFormWithClip).To(BeNil())
			})

			It("returns an error when calling FPDFPage_InsertClipPath", func() {
				FPDFPage_InsertClipPath, err := PdfiumInstance.FPDFPage_InsertClipPath(&requests.FPDFPage_InsertClipPath{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_InsertClipPath).To(BeNil())
			})
		})
	})

	Context("no page object", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFPageObj_TransformClipPath", func() {
				FPDFPageObj_TransformClipPath, err := PdfiumInstance.FPDFPageObj_TransformClipPath(&requests.FPDFPageObj_TransformClipPath{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_TransformClipPath).To(BeNil())
			})
		})
	})

	Context("no clippath", func() {
		When("is opened", func() {
			It("returns an error when calling FPDF_DestroyClipPath", func() {
				FPDF_DestroyClipPath, err := PdfiumInstance.FPDF_DestroyClipPath(&requests.FPDF_DestroyClipPath{})
				Expect(err).To(MatchError("clipPath not given"))
				Expect(FPDF_DestroyClipPath).To(BeNil())
			})
		})
	})

	Context("a normal PDF file with 1 page", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			file, err := os.Open(TestDataPath + "/testdata/test.pdf")
			Expect(err).To(BeNil())

			fileStat, err := file.Stat()
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadCustomDocument(&requests.FPDF_LoadCustomDocument{
				Reader: file,
				Size:   fileStat.Size(),
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
			It("when the mediabox is set, it returns the correct mediabox", func() {
				FPDFPage_SetMediaBox, err := PdfiumInstance.FPDFPage_SetMediaBox(&requests.FPDFPage_SetMediaBox{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Left:   20,
					Bottom: 30,
					Right:  100,
					Top:    150,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_SetMediaBox).To(Equal(&responses.FPDFPage_SetMediaBox{}))

				FPDFPage_GetMediaBox, err := PdfiumInstance.FPDFPage_GetMediaBox(&requests.FPDFPage_GetMediaBox{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetMediaBox).To(Equal(&responses.FPDFPage_GetMediaBox{
					Left:   20,
					Bottom: 30,
					Right:  100,
					Top:    150,
				}))
			})

			It("when the cropbox is set, it returns the correct cropbox", func() {
				FPDFPage_SetCropBox, err := PdfiumInstance.FPDFPage_SetCropBox(&requests.FPDFPage_SetCropBox{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Left:   20,
					Bottom: 30,
					Right:  100,
					Top:    150,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_SetCropBox).To(Equal(&responses.FPDFPage_SetCropBox{}))

				FPDFPage_GetCropBox, err := PdfiumInstance.FPDFPage_GetCropBox(&requests.FPDFPage_GetCropBox{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetCropBox).To(Equal(&responses.FPDFPage_GetCropBox{
					Left:   20,
					Bottom: 30,
					Right:  100,
					Top:    150,
				}))
			})

			It("when the bleedbox is set, it returns the correct bleedbox", func() {
				FPDFPage_SetBleedBox, err := PdfiumInstance.FPDFPage_SetBleedBox(&requests.FPDFPage_SetBleedBox{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Left:   20,
					Bottom: 30,
					Right:  100,
					Top:    150,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_SetBleedBox).To(Equal(&responses.FPDFPage_SetBleedBox{}))

				FPDFPage_GetBleedBox, err := PdfiumInstance.FPDFPage_GetBleedBox(&requests.FPDFPage_GetBleedBox{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetBleedBox).To(Equal(&responses.FPDFPage_GetBleedBox{
					Left:   20,
					Bottom: 30,
					Right:  100,
					Top:    150,
				}))
			})

			It("when the trimbox is set, it returns the correct trimbox", func() {
				FPDFPage_SetTrimBox, err := PdfiumInstance.FPDFPage_SetTrimBox(&requests.FPDFPage_SetTrimBox{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Left:   20,
					Bottom: 30,
					Right:  100,
					Top:    150,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_SetTrimBox).To(Equal(&responses.FPDFPage_SetTrimBox{}))

				FPDFPage_GetTrimBox, err := PdfiumInstance.FPDFPage_GetTrimBox(&requests.FPDFPage_GetTrimBox{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetTrimBox).To(Equal(&responses.FPDFPage_GetTrimBox{
					Left:   20,
					Bottom: 30,
					Right:  100,
					Top:    150,
				}))
			})

			It("when the artbox is set, it returns the correct artbox", func() {
				FPDFPage_SetArtBox, err := PdfiumInstance.FPDFPage_SetArtBox(&requests.FPDFPage_SetArtBox{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Left:   20,
					Bottom: 30,
					Right:  100,
					Top:    150,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_SetArtBox).To(Equal(&responses.FPDFPage_SetArtBox{}))

				FPDFPage_GetArtBox, err := PdfiumInstance.FPDFPage_GetArtBox(&requests.FPDFPage_GetArtBox{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetArtBox).To(Equal(&responses.FPDFPage_GetArtBox{
					Left:   20,
					Bottom: 30,
					Right:  100,
					Top:    150,
				}))
			})

			It("a page can be transformed with a clip", func() {
				FPDFPage_TransFormWithClip, err := PdfiumInstance.FPDFPage_TransFormWithClip(&requests.FPDFPage_TransFormWithClip{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Matrix:   &structs.FPDF_FS_MATRIX{0.5, 0, 0, 0.5, 0, 0},
					ClipRect: &structs.FPDF_FS_RECTF{0.0, 0.0, 20.0, 10.0},
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_TransFormWithClip).To(Equal(&responses.FPDFPage_TransFormWithClip{}))
			})

			It("given an error when trying to transform without a clip and matrix", func() {
				FPDFPage_TransFormWithClip, err := PdfiumInstance.FPDFPage_TransFormWithClip(&requests.FPDFPage_TransFormWithClip{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(MatchError("could not apply clip transform"))
				Expect(FPDFPage_TransFormWithClip).To(BeNil())
			})

			It("allows for a clippath to be created, inserted and destroyed", func() {
				FPDF_CreateClipPath, err := PdfiumInstance.FPDF_CreateClipPath(&requests.FPDF_CreateClipPath{
					Left:   100,
					Bottom: 100,
					Right:  100,
					Top:    100,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_CreateClipPath).To(Not(BeNil()))
				Expect(FPDF_CreateClipPath.ClipPath).To(Not(BeNil()))

				FPDFPage_InsertClipPath, err := PdfiumInstance.FPDFPage_InsertClipPath(&requests.FPDFPage_InsertClipPath{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					ClipPath: FPDF_CreateClipPath.ClipPath,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_InsertClipPath).To(Not(BeNil()))

				FPDF_DestroyClipPath, err := PdfiumInstance.FPDF_DestroyClipPath(&requests.FPDF_DestroyClipPath{
					ClipPath: FPDF_CreateClipPath.ClipPath,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_DestroyClipPath).To(Not(BeNil()))
			})

			It("given an error when no clip path is given in FPDFPage_InsertClipPath", func() {
				FPDFPage_InsertClipPath, err := PdfiumInstance.FPDFPage_InsertClipPath(&requests.FPDFPage_InsertClipPath{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(MatchError("clipPath not given"))
				Expect(FPDFPage_InsertClipPath).To(BeNil())
			})

			// @todo: add extra test for FPDFPageObj_TransformClipPath when FPDFPage_GetObject has been implemented.
		})
	})
})
