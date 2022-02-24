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
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

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

			It("returns an error when calling FPDFPage_New", func() {
				FPDFPage_New, err := PdfiumInstance.FPDFPage_New(&requests.FPDFPage_New{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFPage_New).To(BeNil())
			})

			It("returns an error when calling FPDFPage_Delete", func() {
				FPDFPage_Delete, err := PdfiumInstance.FPDFPage_Delete(&requests.FPDFPage_Delete{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFPage_Delete).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_NewImageObj", func() {
				FPDFPageObj_NewImageObj, err := PdfiumInstance.FPDFPageObj_NewImageObj(&requests.FPDFPageObj_NewImageObj{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFPageObj_NewImageObj).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_NewTextObj", func() {
				FPDFPageObj_NewTextObj, err := PdfiumInstance.FPDFPageObj_NewTextObj(&requests.FPDFPageObj_NewTextObj{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFPageObj_NewTextObj).To(BeNil())
			})

			It("returns an error when calling FPDFText_LoadFont", func() {
				FPDFText_LoadFont, err := PdfiumInstance.FPDFText_LoadFont(&requests.FPDFText_LoadFont{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFText_LoadFont).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_CreateTextObj", func() {
				FPDFPageObj_CreateTextObj, err := PdfiumInstance.FPDFPageObj_CreateTextObj(&requests.FPDFPageObj_CreateTextObj{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFPageObj_CreateTextObj).To(BeNil())
			})
		})
	})

	Context("no page", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFPage_InsertObject", func() {
				FPDFPage_InsertObject, err := PdfiumInstance.FPDFPage_InsertObject(&requests.FPDFPage_InsertObject{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_InsertObject).To(BeNil())
			})

			It("returns an error when calling FPDFPage_CountObjects", func() {
				FPDFPage_CountObjects, err := PdfiumInstance.FPDFPage_CountObjects(&requests.FPDFPage_CountObjects{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_CountObjects).To(BeNil())
			})

			It("returns an error when calling FPDFPage_GetObject", func() {
				FPDFPage_GetObject, err := PdfiumInstance.FPDFPage_GetObject(&requests.FPDFPage_GetObject{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_GetObject).To(BeNil())
			})

			It("returns an error when calling FPDFPage_GenerateContent", func() {
				FPDFPage_GenerateContent, err := PdfiumInstance.FPDFPage_GenerateContent(&requests.FPDFPage_GenerateContent{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_GenerateContent).To(BeNil())
			})

			It("returns an error when calling FPDFPage_TransformAnnots", func() {
				FPDFPage_TransformAnnots, err := PdfiumInstance.FPDFPage_TransformAnnots(&requests.FPDFPage_TransformAnnots{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_TransformAnnots).To(BeNil())
			})
		})
	})

	Context("no page object", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFPageObj_Destroy", func() {
				FPDFPageObj_Destroy, err := PdfiumInstance.FPDFPageObj_Destroy(&requests.FPDFPageObj_Destroy{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_Destroy).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_HasTransparency", func() {
				FPDFPageObj_HasTransparency, err := PdfiumInstance.FPDFPageObj_HasTransparency(&requests.FPDFPageObj_HasTransparency{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_HasTransparency).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_GetType", func() {
				FPDFPageObj_GetType, err := PdfiumInstance.FPDFPageObj_GetType(&requests.FPDFPageObj_GetType{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_GetType).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_Transform", func() {
				FPDFPageObj_Transform, err := PdfiumInstance.FPDFPageObj_Transform(&requests.FPDFPageObj_Transform{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_Transform).To(BeNil())
			})

			It("returns an error when calling FPDFImageObj_SetMatrix", func() {
				FPDFImageObj_SetMatrix, err := PdfiumInstance.FPDFImageObj_SetMatrix(&requests.FPDFImageObj_SetMatrix{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFImageObj_SetMatrix).To(BeNil())
			})

			It("returns an error when calling FPDFImageObj_GetImageDataDecoded", func() {
				FPDFImageObj_GetImageDataDecoded, err := PdfiumInstance.FPDFImageObj_GetImageDataDecoded(&requests.FPDFImageObj_GetImageDataDecoded{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFImageObj_GetImageDataDecoded).To(BeNil())
			})

			It("returns an error when calling FPDFImageObj_GetImageDataRaw", func() {
				FPDFImageObj_GetImageDataRaw, err := PdfiumInstance.FPDFImageObj_GetImageDataRaw(&requests.FPDFImageObj_GetImageDataRaw{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFImageObj_GetImageDataRaw).To(BeNil())
			})

			It("returns an error when calling FPDFImageObj_GetImageFilterCount", func() {
				FPDFImageObj_GetImageFilterCount, err := PdfiumInstance.FPDFImageObj_GetImageFilterCount(&requests.FPDFImageObj_GetImageFilterCount{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFImageObj_GetImageFilterCount).To(BeNil())
			})

			It("returns an error when calling FPDFImageObj_GetImageFilter", func() {
				FPDFImageObj_GetImageFilter, err := PdfiumInstance.FPDFImageObj_GetImageFilter(&requests.FPDFImageObj_GetImageFilter{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFImageObj_GetImageFilter).To(BeNil())
			})

			It("returns an error when calling FPDFImageObj_GetImageMetadata", func() {
				FPDFImageObj_GetImageMetadata, err := PdfiumInstance.FPDFImageObj_GetImageMetadata(&requests.FPDFImageObj_GetImageMetadata{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFImageObj_GetImageMetadata).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_GetBounds", func() {
				FPDFPageObj_GetBounds, err := PdfiumInstance.FPDFPageObj_GetBounds(&requests.FPDFPageObj_GetBounds{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_GetBounds).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_SetBlendMode", func() {
				FPDFPageObj_SetBlendMode, err := PdfiumInstance.FPDFPageObj_SetBlendMode(&requests.FPDFPageObj_SetBlendMode{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_SetBlendMode).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_SetStrokeColor", func() {
				FPDFPageObj_SetStrokeColor, err := PdfiumInstance.FPDFPageObj_SetStrokeColor(&requests.FPDFPageObj_SetStrokeColor{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_SetStrokeColor).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_GetStrokeColor", func() {
				FPDFPageObj_GetStrokeColor, err := PdfiumInstance.FPDFPageObj_GetStrokeColor(&requests.FPDFPageObj_GetStrokeColor{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_GetStrokeColor).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_SetStrokeWidth", func() {
				FPDFPageObj_SetStrokeWidth, err := PdfiumInstance.FPDFPageObj_SetStrokeWidth(&requests.FPDFPageObj_SetStrokeWidth{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_SetStrokeWidth).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_GetStrokeWidth", func() {
				FPDFPageObj_GetStrokeWidth, err := PdfiumInstance.FPDFPageObj_GetStrokeWidth(&requests.FPDFPageObj_GetStrokeWidth{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_GetStrokeWidth).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_GetLineJoin", func() {
				FPDFPageObj_GetLineJoin, err := PdfiumInstance.FPDFPageObj_GetLineJoin(&requests.FPDFPageObj_GetLineJoin{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_GetLineJoin).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_SetLineJoin", func() {
				FPDFPageObj_SetLineJoin, err := PdfiumInstance.FPDFPageObj_SetLineJoin(&requests.FPDFPageObj_SetLineJoin{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_SetLineJoin).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_GetLineCap", func() {
				FPDFPageObj_GetLineCap, err := PdfiumInstance.FPDFPageObj_GetLineCap(&requests.FPDFPageObj_GetLineCap{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_GetLineCap).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_SetLineCap", func() {
				FPDFPageObj_SetLineCap, err := PdfiumInstance.FPDFPageObj_SetLineCap(&requests.FPDFPageObj_SetLineCap{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_SetLineCap).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_SetFillColor", func() {
				FPDFPageObj_SetFillColor, err := PdfiumInstance.FPDFPageObj_SetFillColor(&requests.FPDFPageObj_SetFillColor{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_SetFillColor).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_GetFillColor", func() {
				FPDFPageObj_GetFillColor, err := PdfiumInstance.FPDFPageObj_GetFillColor(&requests.FPDFPageObj_GetFillColor{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_GetFillColor).To(BeNil())
			})

			It("returns an error when calling FPDFPath_CountSegments", func() {
				FPDFPath_CountSegments, err := PdfiumInstance.FPDFPath_CountSegments(&requests.FPDFPath_CountSegments{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPath_CountSegments).To(BeNil())
			})

			It("returns an error when calling FPDFPath_GetPathSegment", func() {
				FPDFPath_GetPathSegment, err := PdfiumInstance.FPDFPath_GetPathSegment(&requests.FPDFPath_GetPathSegment{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPath_GetPathSegment).To(BeNil())
			})

			It("returns an error when calling FPDFPath_MoveTo", func() {
				FPDFPath_MoveTo, err := PdfiumInstance.FPDFPath_MoveTo(&requests.FPDFPath_MoveTo{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPath_MoveTo).To(BeNil())
			})

			It("returns an error when calling FPDFPath_LineTo", func() {
				FPDFPath_LineTo, err := PdfiumInstance.FPDFPath_LineTo(&requests.FPDFPath_LineTo{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPath_LineTo).To(BeNil())
			})

			It("returns an error when calling FPDFPath_BezierTo", func() {
				FPDFPath_BezierTo, err := PdfiumInstance.FPDFPath_BezierTo(&requests.FPDFPath_BezierTo{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPath_BezierTo).To(BeNil())
			})

			It("returns an error when calling FPDFPath_Close", func() {
				FPDFPath_Close, err := PdfiumInstance.FPDFPath_Close(&requests.FPDFPath_Close{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPath_Close).To(BeNil())
			})

			It("returns an error when calling FPDFPath_SetDrawMode", func() {
				FPDFPath_SetDrawMode, err := PdfiumInstance.FPDFPath_SetDrawMode(&requests.FPDFPath_SetDrawMode{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPath_SetDrawMode).To(BeNil())
			})

			It("returns an error when calling FPDFPath_GetDrawMode", func() {
				FPDFPath_GetDrawMode, err := PdfiumInstance.FPDFPath_GetDrawMode(&requests.FPDFPath_GetDrawMode{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPath_GetDrawMode).To(BeNil())
			})

			It("returns an error when calling FPDFText_SetText", func() {
				FPDFText_SetText, err := PdfiumInstance.FPDFText_SetText(&requests.FPDFText_SetText{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFText_SetText).To(BeNil())
			})

			It("returns an error when calling FPDFText_SetCharcodes", func() {
				FPDFText_SetCharcodes, err := PdfiumInstance.FPDFText_SetCharcodes(&requests.FPDFText_SetCharcodes{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFText_SetCharcodes).To(BeNil())
			})

			It("returns an error when calling FPDFTextObj_GetFontSize", func() {
				FPDFTextObj_GetFontSize, err := PdfiumInstance.FPDFTextObj_GetFontSize(&requests.FPDFTextObj_GetFontSize{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFTextObj_GetFontSize).To(BeNil())
			})

			It("returns an error when calling FPDFTextObj_GetTextRenderMode", func() {
				FPDFTextObj_GetTextRenderMode, err := PdfiumInstance.FPDFTextObj_GetTextRenderMode(&requests.FPDFTextObj_GetTextRenderMode{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFTextObj_GetTextRenderMode).To(BeNil())
			})

			It("returns an error when calling FPDFTextObj_GetText", func() {
				FPDFTextObj_GetText, err := PdfiumInstance.FPDFTextObj_GetText(&requests.FPDFTextObj_GetText{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFTextObj_GetText).To(BeNil())
			})

			It("returns an error when calling FPDFFormObj_CountObjects", func() {
				FPDFFormObj_CountObjects, err := PdfiumInstance.FPDFFormObj_CountObjects(&requests.FPDFFormObj_CountObjects{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFFormObj_CountObjects).To(BeNil())
			})

			It("returns an error when calling FPDFFormObj_GetObject", func() {
				FPDFFormObj_GetObject, err := PdfiumInstance.FPDFFormObj_GetObject(&requests.FPDFFormObj_GetObject{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFFormObj_GetObject).To(BeNil())
			})

			It("returns an error when calling FPDFImageObj_SetBitmap", func() {
				FPDFImageObj_SetBitmap, err := PdfiumInstance.FPDFImageObj_SetBitmap(&requests.FPDFImageObj_SetBitmap{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFImageObj_SetBitmap).To(BeNil())
			})

			It("returns an error when calling FPDFImageObj_GetBitmap", func() {
				FPDFImageObj_GetBitmap, err := PdfiumInstance.FPDFImageObj_GetBitmap(&requests.FPDFImageObj_GetBitmap{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFImageObj_GetBitmap).To(BeNil())
			})
		})
	})

	Context("no path segment object", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFPathSegment_GetPoint", func() {
				FPDFPathSegment_GetPoint, err := PdfiumInstance.FPDFPathSegment_GetPoint(&requests.FPDFPathSegment_GetPoint{})
				Expect(err).To(MatchError("pathSegment not given"))
				Expect(FPDFPathSegment_GetPoint).To(BeNil())
			})

			It("returns an error when calling FPDFPathSegment_GetType", func() {
				FPDFPathSegment_GetType, err := PdfiumInstance.FPDFPathSegment_GetType(&requests.FPDFPathSegment_GetType{})
				Expect(err).To(MatchError("pathSegment not given"))
				Expect(FPDFPathSegment_GetType).To(BeNil())
			})

			It("returns an error when calling FPDFPathSegment_GetClose", func() {
				FPDFPathSegment_GetClose, err := PdfiumInstance.FPDFPathSegment_GetClose(&requests.FPDFPathSegment_GetClose{})
				Expect(err).To(MatchError("pathSegment not given"))
				Expect(FPDFPathSegment_GetClose).To(BeNil())
			})
		})
	})

	Context("no font object", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFFont_Close", func() {
				FPDFFont_Close, err := PdfiumInstance.FPDFFont_Close(&requests.FPDFFont_Close{})
				Expect(err).To(MatchError("font not given"))
				Expect(FPDFFont_Close).To(BeNil())
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
