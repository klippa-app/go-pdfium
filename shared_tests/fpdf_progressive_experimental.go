//go:build pdfium_experimental
// +build pdfium_experimental

package shared_tests

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Describe("fpdf_progressive_experimental", func() {
	BeforeEach(func() {
		if TestType == "multi" {
			Skip("Multi-threaded usage does not support setting callbacks")
		}
	})

	AfterEach(func() {
		if TestType == "multi" {
			Skip("Multi-threaded usage does not support setting callbacks")
		}
	})

	Context("no page", func() {
		When("is opened", func() {
			It("returns an error when calling FPDF_RenderPageBitmapWithColorScheme_Start", func() {
				FPDF_RenderPageBitmap_Start, err := PdfiumInstance.FPDF_RenderPageBitmap_Start(&requests.FPDF_RenderPageBitmap_Start{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDF_RenderPageBitmap_Start).To(BeNil())
			})
		})
	})

	Context("a normal PDF file", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/text_form.pdf")
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
			When("no bitmap is given", func() {
				It("returns an error when calling FPDF_RenderPageBitmapWithColorScheme_Start", func() {
					FPDF_RenderPageBitmapWithColorScheme_Start, err := PdfiumInstance.FPDF_RenderPageBitmapWithColorScheme_Start(&requests.FPDF_RenderPageBitmapWithColorScheme_Start{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})
					Expect(err).To(MatchError("bitmap not given"))
					Expect(FPDF_RenderPageBitmapWithColorScheme_Start).To(BeNil())
				})
			})

			When("no callback is given", func() {
				It("returns an error when calling FPDF_RenderPageBitmapWithColorScheme_Start", func() {
					By("creating a bitmap")
					FPDFBitmap_Create, err := PdfiumInstance.FPDFBitmap_Create(&requests.FPDFBitmap_Create{
						Width:  1000,
						Height: 1000,
						Alpha:  1,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_Create).To(Not(BeNil()))
					Expect(FPDFBitmap_Create.Bitmap).To(Not(BeNil()))

					By("calling FPDF_RenderPageBitmap_Start")
					FPDF_RenderPageBitmap_Start, err := PdfiumInstance.FPDF_RenderPageBitmapWithColorScheme_Start(&requests.FPDF_RenderPageBitmapWithColorScheme_Start{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
						Bitmap: FPDFBitmap_Create.Bitmap,
						ColorScheme: &structs.FPDF_COLORSCHEME{
							PathFillColor:   0xFFFF0000,
							PathStrokeColor: 0xFF00FF00,
							TextFillColor:   0xFF0000FF,
							TextStrokeColor: 0xFF00FFFF,
						},
					})
					Expect(err).To(MatchError("NeedToPauseNowCallback can't be nil"))
					Expect(FPDF_RenderPageBitmap_Start).To(BeNil())

					By("destroying the bitmap")
					FPDFBitmap_Destroy, err := PdfiumInstance.FPDFBitmap_Destroy(&requests.FPDFBitmap_Destroy{
						Bitmap: FPDFBitmap_Create.Bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_Destroy).To(Equal(&responses.FPDFBitmap_Destroy{}))
				})
			})

			It("can be rendered progressively", func() {
				By("creating a bitmap")
				FPDFBitmap_Create, err := PdfiumInstance.FPDFBitmap_Create(&requests.FPDFBitmap_Create{
					Width:  1000,
					Height: 1000,
					Alpha:  1,
				})
				Expect(err).To(BeNil())
				Expect(FPDFBitmap_Create).To(Not(BeNil()))
				Expect(FPDFBitmap_Create.Bitmap).To(Not(BeNil()))

				By("starting the progressive rendering and directly pausing it")
				FPDF_RenderPageBitmapWithColorScheme_Start, err := PdfiumInstance.FPDF_RenderPageBitmapWithColorScheme_Start(&requests.FPDF_RenderPageBitmapWithColorScheme_Start{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Bitmap: FPDFBitmap_Create.Bitmap,
					StartX: 0,
					StartY: 0,
					SizeX:  1000,
					SizeY:  1000,
					Rotate: enums.FPDF_PAGE_ROTATION_NONE,
					Flags:  0,
					ColorScheme: &structs.FPDF_COLORSCHEME{
						PathFillColor:   0xFFFF0000,
						PathStrokeColor: 0xFF00FF00,
						TextFillColor:   0xFF0000FF,
						TextStrokeColor: 0xFF00FFFF,
					},
					NeedToPauseNowCallback: func() bool {
						return true
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDF_RenderPageBitmapWithColorScheme_Start).To(Equal(&responses.FPDF_RenderPageBitmapWithColorScheme_Start{
					RenderStatus: enums.FPDF_RENDER_STATUS_TOBECONTINUED,
				}))

				By("starting the continuing the rendering and pausing it again")
				FPDF_RenderPage_Continue, err := PdfiumInstance.FPDF_RenderPage_Continue(&requests.FPDF_RenderPage_Continue{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					NeedToPauseNowCallback: func() bool {
						return true
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDF_RenderPage_Continue).To(Equal(&responses.FPDF_RenderPage_Continue{
					RenderStatus: enums.FPDF_RENDER_STATUS_DONE,
				}))

				By("cleaning up the resources")
				FPDF_RenderPage_Close, err := PdfiumInstance.FPDF_RenderPage_Close(&requests.FPDF_RenderPage_Close{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDF_RenderPage_Close).To(Equal(&responses.FPDF_RenderPage_Close{}))

				By("destroying the bitmap")
				FPDFBitmap_Destroy, err := PdfiumInstance.FPDFBitmap_Destroy(&requests.FPDFBitmap_Destroy{
					Bitmap: FPDFBitmap_Create.Bitmap,
				})
				Expect(err).To(BeNil())
				Expect(FPDFBitmap_Destroy).To(Equal(&responses.FPDFBitmap_Destroy{}))
			})
		})
	})
})
