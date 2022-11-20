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

var _ = Describe("fpdf_edit", func() {
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

	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFPageObjMark_SetIntParam", func() {
				FPDFPageObjMark_SetIntParam, err := PdfiumInstance.FPDFPageObjMark_SetIntParam(&requests.FPDFPageObjMark_SetIntParam{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFPageObjMark_SetIntParam).To(BeNil())
			})

			It("returns an error when calling FPDFPageObjMark_SetIntParam", func() {
				FPDFPageObjMark_SetStringParam, err := PdfiumInstance.FPDFPageObjMark_SetStringParam(&requests.FPDFPageObjMark_SetStringParam{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFPageObjMark_SetStringParam).To(BeNil())
			})

			It("returns an error when calling FPDFPageObjMark_SetBlobParam", func() {
				FPDFPageObjMark_SetBlobParam, err := PdfiumInstance.FPDFPageObjMark_SetBlobParam(&requests.FPDFPageObjMark_SetBlobParam{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFPageObjMark_SetBlobParam).To(BeNil())
			})

			It("returns an error when calling FPDFImageObj_GetRenderedBitmap", func() {
				FPDFImageObj_GetRenderedBitmap, err := PdfiumInstance.FPDFImageObj_GetRenderedBitmap(&requests.FPDFImageObj_GetRenderedBitmap{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFImageObj_GetRenderedBitmap).To(BeNil())
			})

			It("returns an error when calling FPDFText_LoadStandardFont", func() {
				FPDFText_LoadStandardFont, err := PdfiumInstance.FPDFText_LoadStandardFont(&requests.FPDFText_LoadStandardFont{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFText_LoadStandardFont).To(BeNil())
			})

			It("returns an error when calling FPDFTextObj_GetRenderedBitmap", func() {
				FPDFTextObj_GetRenderedBitmap, err := PdfiumInstance.FPDFTextObj_GetRenderedBitmap(&requests.FPDFTextObj_GetRenderedBitmap{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFTextObj_GetRenderedBitmap).To(BeNil())
			})
		})
	})

	Context("no page", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFPage_RemoveObject", func() {
				FPDFPage_RemoveObject, err := PdfiumInstance.FPDFPage_RemoveObject(&requests.FPDFPage_RemoveObject{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFPage_RemoveObject).To(BeNil())
			})
		})
	})

	Context("no page object", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFPageObj_GetMatrix", func() {
				FPDFPageObj_GetMatrix, err := PdfiumInstance.FPDFPageObj_GetMatrix(&requests.FPDFPageObj_GetMatrix{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_GetMatrix).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_SetMatrix", func() {
				FPDFPageObj_SetMatrix, err := PdfiumInstance.FPDFPageObj_SetMatrix(&requests.FPDFPageObj_SetMatrix{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_SetMatrix).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_CountMarks", func() {
				FPDFPageObj_CountMarks, err := PdfiumInstance.FPDFPageObj_CountMarks(&requests.FPDFPageObj_CountMarks{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_CountMarks).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_GetMark", func() {
				FPDFPageObj_GetMark, err := PdfiumInstance.FPDFPageObj_GetMark(&requests.FPDFPageObj_GetMark{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_GetMark).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_AddMark", func() {
				FPDFPageObj_AddMark, err := PdfiumInstance.FPDFPageObj_AddMark(&requests.FPDFPageObj_AddMark{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_AddMark).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_RemoveMark", func() {
				FPDFPageObj_RemoveMark, err := PdfiumInstance.FPDFPageObj_RemoveMark(&requests.FPDFPageObj_RemoveMark{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_RemoveMark).To(BeNil())
			})

			It("returns an error when calling FPDFPageObjMark_RemoveParam", func() {
				FPDFPageObjMark_RemoveParam, err := PdfiumInstance.FPDFPageObjMark_RemoveParam(&requests.FPDFPageObjMark_RemoveParam{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObjMark_RemoveParam).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_GetDashPhase", func() {
				FPDFPageObj_GetDashPhase, err := PdfiumInstance.FPDFPageObj_GetDashPhase(&requests.FPDFPageObj_GetDashPhase{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_GetDashPhase).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_SetDashPhase", func() {
				FPDFPageObj_SetDashPhase, err := PdfiumInstance.FPDFPageObj_SetDashPhase(&requests.FPDFPageObj_SetDashPhase{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_SetDashPhase).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_GetDashCount", func() {
				FPDFPageObj_GetDashCount, err := PdfiumInstance.FPDFPageObj_GetDashCount(&requests.FPDFPageObj_GetDashCount{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_GetDashCount).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_GetDashArray", func() {
				FPDFPageObj_GetDashArray, err := PdfiumInstance.FPDFPageObj_GetDashArray(&requests.FPDFPageObj_GetDashArray{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_GetDashArray).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_GetDashArray", func() {
				FPDFPageObj_GetDashArray, err := PdfiumInstance.FPDFPageObj_GetDashArray(&requests.FPDFPageObj_GetDashArray{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_GetDashArray).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_SetDashArray", func() {
				FPDFPageObj_SetDashArray, err := PdfiumInstance.FPDFPageObj_SetDashArray(&requests.FPDFPageObj_SetDashArray{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_SetDashArray).To(BeNil())
			})

			It("returns an error when calling FPDFTextObj_SetTextRenderMode", func() {
				FPDFTextObj_SetTextRenderMode, err := PdfiumInstance.FPDFTextObj_SetTextRenderMode(&requests.FPDFTextObj_SetTextRenderMode{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFTextObj_SetTextRenderMode).To(BeNil())
			})

			It("returns an error when calling FPDFTextObj_GetFont", func() {
				FPDFTextObj_GetFont, err := PdfiumInstance.FPDFTextObj_GetFont(&requests.FPDFTextObj_GetFont{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFTextObj_GetFont).To(BeNil())
			})

			It("returns an error when calling FPDFPageObj_GetRotatedBounds", func() {
				FPDFPageObj_GetRotatedBounds, err := PdfiumInstance.FPDFPageObj_GetRotatedBounds(&requests.FPDFPageObj_GetRotatedBounds{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_GetRotatedBounds).To(BeNil())
			})
		})
	})

	Context("no path object mark", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFPageObjMark_GetName", func() {
				FPDFPageObjMark_GetName, err := PdfiumInstance.FPDFPageObjMark_GetName(&requests.FPDFPageObjMark_GetName{})
				Expect(err).To(MatchError("pageObjectMark not given"))
				Expect(FPDFPageObjMark_GetName).To(BeNil())
			})

			It("returns an error when calling FPDFPageObjMark_CountParams", func() {
				FPDFPageObjMark_CountParams, err := PdfiumInstance.FPDFPageObjMark_CountParams(&requests.FPDFPageObjMark_CountParams{})
				Expect(err).To(MatchError("pageObjectMark not given"))
				Expect(FPDFPageObjMark_CountParams).To(BeNil())
			})

			It("returns an error when calling FPDFPageObjMark_GetParamKey", func() {
				FPDFPageObjMark_GetParamKey, err := PdfiumInstance.FPDFPageObjMark_GetParamKey(&requests.FPDFPageObjMark_GetParamKey{})
				Expect(err).To(MatchError("pageObjectMark not given"))
				Expect(FPDFPageObjMark_GetParamKey).To(BeNil())
			})

			It("returns an error when calling FPDFPageObjMark_GetParamKey", func() {
				FPDFPageObjMark_GetParamValueType, err := PdfiumInstance.FPDFPageObjMark_GetParamValueType(&requests.FPDFPageObjMark_GetParamValueType{})
				Expect(err).To(MatchError("pageObjectMark not given"))
				Expect(FPDFPageObjMark_GetParamValueType).To(BeNil())
			})

			It("returns an error when calling FPDFPageObjMark_GetParamIntValue", func() {
				FPDFPageObjMark_GetParamIntValue, err := PdfiumInstance.FPDFPageObjMark_GetParamIntValue(&requests.FPDFPageObjMark_GetParamIntValue{})
				Expect(err).To(MatchError("pageObjectMark not given"))
				Expect(FPDFPageObjMark_GetParamIntValue).To(BeNil())
			})

			It("returns an error when calling FPDFPageObjMark_GetParamStringValue", func() {
				FPDFPageObjMark_GetParamStringValue, err := PdfiumInstance.FPDFPageObjMark_GetParamStringValue(&requests.FPDFPageObjMark_GetParamStringValue{})
				Expect(err).To(MatchError("pageObjectMark not given"))
				Expect(FPDFPageObjMark_GetParamStringValue).To(BeNil())
			})

			It("returns an error when calling FPDFPageObjMark_GetParamBlobValue", func() {
				FPDFPageObjMark_GetParamBlobValue, err := PdfiumInstance.FPDFPageObjMark_GetParamBlobValue(&requests.FPDFPageObjMark_GetParamBlobValue{})
				Expect(err).To(MatchError("pageObjectMark not given"))
				Expect(FPDFPageObjMark_GetParamBlobValue).To(BeNil())
			})
		})
	})

	Context("no font object", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFFont_GetFontName", func() {
				FPDFFont_GetFontName, err := PdfiumInstance.FPDFFont_GetFontName(&requests.FPDFFont_GetFontName{})
				Expect(err).To(MatchError("font not given"))
				Expect(FPDFFont_GetFontName).To(BeNil())
			})

			It("returns an error when calling FPDFFont_GetFontData", func() {
				FPDFFont_GetFontData, err := PdfiumInstance.FPDFFont_GetFontData(&requests.FPDFFont_GetFontData{})
				Expect(err).To(MatchError("font not given"))
				Expect(FPDFFont_GetFontData).To(BeNil())
			})

			It("returns an error when calling FPDFFont_GetIsEmbedded", func() {
				FPDFFont_GetIsEmbedded, err := PdfiumInstance.FPDFFont_GetIsEmbedded(&requests.FPDFFont_GetIsEmbedded{})
				Expect(err).To(MatchError("font not given"))
				Expect(FPDFFont_GetIsEmbedded).To(BeNil())
			})

			It("returns an error when calling FPDFFont_GetFlags", func() {
				FPDFFont_GetFlags, err := PdfiumInstance.FPDFFont_GetFlags(&requests.FPDFFont_GetFlags{})
				Expect(err).To(MatchError("font not given"))
				Expect(FPDFFont_GetFlags).To(BeNil())
			})

			It("returns an error when calling FPDFFont_GetWeight", func() {
				FPDFFont_GetWeight, err := PdfiumInstance.FPDFFont_GetWeight(&requests.FPDFFont_GetWeight{})
				Expect(err).To(MatchError("font not given"))
				Expect(FPDFFont_GetWeight).To(BeNil())
			})

			It("returns an error when calling FPDFFont_GetItalicAngle", func() {
				FPDFFont_GetItalicAngle, err := PdfiumInstance.FPDFFont_GetItalicAngle(&requests.FPDFFont_GetItalicAngle{})
				Expect(err).To(MatchError("font not given"))
				Expect(FPDFFont_GetItalicAngle).To(BeNil())
			})

			It("returns an error when calling FPDFFont_GetAscent", func() {
				FPDFFont_GetAscent, err := PdfiumInstance.FPDFFont_GetAscent(&requests.FPDFFont_GetAscent{})
				Expect(err).To(MatchError("font not given"))
				Expect(FPDFFont_GetAscent).To(BeNil())
			})

			It("returns an error when calling FPDFFont_GetDescent", func() {
				FPDFFont_GetDescent, err := PdfiumInstance.FPDFFont_GetDescent(&requests.FPDFFont_GetDescent{})
				Expect(err).To(MatchError("font not given"))
				Expect(FPDFFont_GetDescent).To(BeNil())
			})

			It("returns an error when calling FPDFFont_GetGlyphWidth", func() {
				FPDFFont_GetGlyphWidth, err := PdfiumInstance.FPDFFont_GetGlyphWidth(&requests.FPDFFont_GetGlyphWidth{})
				Expect(err).To(MatchError("font not given"))
				Expect(FPDFFont_GetGlyphWidth).To(BeNil())
			})

			It("returns an error when calling FPDFFont_GetGlyphPath", func() {
				FPDFFont_GetGlyphPath, err := PdfiumInstance.FPDFFont_GetGlyphPath(&requests.FPDFFont_GetGlyphPath{})
				Expect(err).To(MatchError("font not given"))
				Expect(FPDFFont_GetGlyphPath).To(BeNil())
			})
		})
	})

	Context("no glyph path", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFGlyphPath_CountGlyphSegments", func() {
				FPDFGlyphPath_CountGlyphSegments, err := PdfiumInstance.FPDFGlyphPath_CountGlyphSegments(&requests.FPDFGlyphPath_CountGlyphSegments{})
				Expect(err).To(MatchError("glyphPath not given"))
				Expect(FPDFGlyphPath_CountGlyphSegments).To(BeNil())
			})

			It("returns an error when calling FPDFGlyphPath_GetGlyphPathSegment", func() {
				FPDFGlyphPath_GetGlyphPathSegment, err := PdfiumInstance.FPDFGlyphPath_GetGlyphPathSegment(&requests.FPDFGlyphPath_GetGlyphPathSegment{})
				Expect(err).To(MatchError("glyphPath not given"))
				Expect(FPDFGlyphPath_GetGlyphPathSegment).To(BeNil())
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
			Context("a font is loaded", func() {
				var font references.FPDF_FONT

				BeforeEach(func() {
					fontData, err := ioutil.ReadFile(TestDataPath + "/testdata/NotoSansSC-Regular.subset.otf")
					Expect(err).To(BeNil())

					FPDFText_LoadFont, err := PdfiumInstance.FPDFText_LoadFont(&requests.FPDFText_LoadFont{
						Document: doc,
						Data:     fontData,
						FontType: enums.FPDF_FONT_TRUETYPE,
						CID:      true,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_LoadFont).To(Not(BeNil()))
					Expect(FPDFText_LoadFont.Font).To(Not(BeEmpty()))
					font = FPDFText_LoadFont.Font
				})

				AfterEach(func() {
					FPDFFont_Close, err := PdfiumInstance.FPDFFont_Close(&requests.FPDFFont_Close{
						Font: font,
					})
					Expect(err).To(BeNil())
					Expect(FPDFFont_Close).To(Equal(&responses.FPDFFont_Close{}))
				})
			})
		})
	})

	Context("a PDF file with an embedded font", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/text_font.pdf")
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
			It("allows us to get a standard font", func() {
				FPDFText_LoadStandardFont, err := PdfiumInstance.FPDFText_LoadStandardFont(&requests.FPDFText_LoadStandardFont{
					Document: doc,
					Font:     "Helvetica-BoldItalic",
				})
				Expect(err).To(BeNil())
				Expect(FPDFText_LoadStandardFont).ToNot(BeNil())
				Expect(FPDFText_LoadStandardFont.Font).ToNot(BeEmpty())

				FPDFFont_Close, err := PdfiumInstance.FPDFFont_Close(&requests.FPDFFont_Close{
					Font: FPDFText_LoadStandardFont.Font,
				})
				Expect(err).To(BeNil())
				Expect(FPDFFont_Close).To(Equal(&responses.FPDFFont_Close{}))
			})

			It("gives an error when loading a non-existing standard font", func() {
				FPDFText_LoadStandardFont, err := PdfiumInstance.FPDFText_LoadStandardFont(&requests.FPDFText_LoadStandardFont{
					Document: doc,
					Font:     "this-font-isnt-here",
				})
				Expect(err).To(MatchError("could not load standard font"))
				Expect(FPDFText_LoadStandardFont).To(BeNil())
			})

			It("returns an error when calling FPDFTextObj_GetRenderedBitmap without a page object", func() {
				FPDFTextObj_GetRenderedBitmap, err := PdfiumInstance.FPDFTextObj_GetRenderedBitmap(&requests.FPDFTextObj_GetRenderedBitmap{
					Document: doc,
				})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFTextObj_GetRenderedBitmap).To(BeNil())
			})

			Context("when an page object has been loaded", func() {
				var pageObject references.FPDF_PAGEOBJECT

				BeforeEach(func() {
					FPDFPage_GetObject, err := PdfiumInstance.FPDFPage_GetObject(&requests.FPDFPage_GetObject{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
						Index: 0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPage_GetObject).To(Not(BeNil()))
					Expect(FPDFPage_GetObject.PageObject).To(Not(BeEmpty()))
					pageObject = FPDFPage_GetObject.PageObject

					FPDFPageObj_GetType, err := PdfiumInstance.FPDFPageObj_GetType(&requests.FPDFPageObj_GetType{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetType).To(Not(BeNil()))
					Expect(FPDFPageObj_GetType).To(Equal(&responses.FPDFPageObj_GetType{
						Type: enums.FPDF_PAGEOBJ_TEXT,
					}))
				})

				It("allows us getting the font", func() {
					FPDFTextObj_GetFont, err := PdfiumInstance.FPDFTextObj_GetFont(&requests.FPDFTextObj_GetFont{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFTextObj_GetFont).ToNot(BeNil())
					Expect(FPDFTextObj_GetFont.Font).ToNot(BeEmpty())
				})

				It("allows us getting the rotated bounds", func() {
					FPDFPageObj_GetRotatedBounds, err := PdfiumInstance.FPDFPageObj_GetRotatedBounds(&requests.FPDFPageObj_GetRotatedBounds{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetRotatedBounds).ToNot(BeNil())
					Expect(FPDFPageObj_GetRotatedBounds.QuadPoints).To(Equal(structs.FPDF_FS_QUADPOINTSF{
						X1: 57.21999740600586,
						Y1: 723.9920043945312,
						X2: 76.79199981689453,
						Y2: 723.9920043945312,
						X3: 76.79199981689453,
						Y3: 732.5479736328125,
						X4: 57.21999740600586,
						Y4: 732.5479736328125,
					}))
				})

				It("allows us getting the rendered bitmap", func() {
					FPDFTextObj_GetRenderedBitmap, err := PdfiumInstance.FPDFTextObj_GetRenderedBitmap(&requests.FPDFTextObj_GetRenderedBitmap{
						Document:   doc,
						PageObject: pageObject,
						Scale:      0.5,
					})
					Expect(err).To(BeNil())
					Expect(FPDFTextObj_GetRenderedBitmap).ToNot(BeNil())
					Expect(FPDFTextObj_GetRenderedBitmap.Bitmap).ToNot(BeNil())
				})

				It("allows us getting the rendered bitmap when giving the page object", func() {
					FPDF_LoadPage, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{
						Document: doc,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_LoadPage).ToNot(BeNil())
					Expect(FPDF_LoadPage.Page).ToNot(BeNil())

					FPDFTextObj_GetRenderedBitmap, err := PdfiumInstance.FPDFTextObj_GetRenderedBitmap(&requests.FPDFTextObj_GetRenderedBitmap{
						Document:   doc,
						Page:       FPDF_LoadPage.Page,
						PageObject: pageObject,
						Scale:      0.5,
					})
					Expect(err).To(BeNil())
					Expect(FPDFTextObj_GetRenderedBitmap).ToNot(BeNil())
					Expect(FPDFTextObj_GetRenderedBitmap.Bitmap).ToNot(BeNil())
				})

				It("returns an error when getting the rendered bitmap when giving an invalid page object", func() {
					FPDFTextObj_GetRenderedBitmap, err := PdfiumInstance.FPDFTextObj_GetRenderedBitmap(&requests.FPDFTextObj_GetRenderedBitmap{
						Document:   doc,
						Page:       "fake",
						PageObject: pageObject,
						Scale:      0.5,
					})
					Expect(err).To(MatchError("could not find page handle, perhaps the page was already closed or you tried to share pages between instances or documents"))
					Expect(FPDFTextObj_GetRenderedBitmap).To(BeNil())
				})

				It("returns an error when getting the rendered bitmap when giving an invalid scale", func() {
					FPDFTextObj_GetRenderedBitmap, err := PdfiumInstance.FPDFTextObj_GetRenderedBitmap(&requests.FPDFTextObj_GetRenderedBitmap{
						Document:   doc,
						PageObject: pageObject,
						Scale:      0,
					})
					Expect(err).To(MatchError("could not render text object as bitmap"))
					Expect(FPDFTextObj_GetRenderedBitmap).To(BeNil())
				})

				It("allows us setting the text render mode", func() {
					FPDFTextObj_GetFont, err := PdfiumInstance.FPDFTextObj_SetTextRenderMode(&requests.FPDFTextObj_SetTextRenderMode{
						PageObject:     pageObject,
						TextRenderMode: enums.FPDF_TEXTRENDERMODE_STROKE,
					})
					Expect(err).To(BeNil())
					Expect(FPDFTextObj_GetFont).ToNot(BeNil())
				})

				Context("when a font has been loaded", func() {
					var font references.FPDF_FONT

					BeforeEach(func() {
						FPDFTextObj_GetFont, err := PdfiumInstance.FPDFTextObj_GetFont(&requests.FPDFTextObj_GetFont{
							PageObject: pageObject,
						})
						Expect(err).To(BeNil())
						Expect(FPDFTextObj_GetFont).ToNot(BeNil())
						Expect(FPDFTextObj_GetFont.Font).ToNot(BeEmpty())
						font = FPDFTextObj_GetFont.Font
					})

					It("allows us getting the font name", func() {
						FPDFFont_GetFontName, err := PdfiumInstance.FPDFFont_GetFontName(&requests.FPDFFont_GetFontName{
							Font: font,
						})
						Expect(err).To(BeNil())
						Expect(FPDFFont_GetFontName).To(Equal(&responses.FPDFFont_GetFontName{
							FontName: "Liberation Serif",
						}))
					})

					It("allows us getting the font data", func() {
						FPDFFont_GetFontData, err := PdfiumInstance.FPDFFont_GetFontData(&requests.FPDFFont_GetFontData{
							Font: font,
						})
						Expect(err).To(BeNil())
						Expect(FPDFFont_GetFontData).ToNot(BeNil())
						Expect(FPDFFont_GetFontData.FontData).ToNot(BeEmpty())
						Expect(FPDFFont_GetFontData.FontData).To(HaveLen(8268))
					})

					It("allows us getting whether the font is embedded", func() {
						FPDFFont_GetIsEmbedded, err := PdfiumInstance.FPDFFont_GetIsEmbedded(&requests.FPDFFont_GetIsEmbedded{
							Font: font,
						})
						Expect(err).To(BeNil())
						Expect(FPDFFont_GetIsEmbedded).To(Equal(&responses.FPDFFont_GetIsEmbedded{
							IsEmbedded: true,
						}))
					})

					It("allows us getting the font flags", func() {
						FPDFFont_GetFlags, err := PdfiumInstance.FPDFFont_GetFlags(&requests.FPDFFont_GetFlags{
							Font: font,
						})
						Expect(err).To(BeNil())
						Expect(FPDFFont_GetFlags).To(Equal(&responses.FPDFFont_GetFlags{
							Flags:       4,
							FixedPitch:  false,
							Serif:       true,
							Symbolic:    false,
							Script:      false,
							Nonsymbolic: false,
							Italic:      false,
							AllCap:      false,
							SmallCap:    false,
							ForceBold:   false,
						}))
					})

					It("allows us getting the font weight", func() {
						FPDFFont_GetWeight, err := PdfiumInstance.FPDFFont_GetWeight(&requests.FPDFFont_GetWeight{
							Font: font,
						})
						Expect(err).To(BeNil())
						Expect(FPDFFont_GetWeight).To(Equal(&responses.FPDFFont_GetWeight{
							Weight: 400,
						}))
					})

					It("allows us getting the font italic angle", func() {
						FPDFFont_GetItalicAngle, err := PdfiumInstance.FPDFFont_GetItalicAngle(&requests.FPDFFont_GetItalicAngle{
							Font: font,
						})
						Expect(err).To(BeNil())
						Expect(FPDFFont_GetItalicAngle).To(Equal(&responses.FPDFFont_GetItalicAngle{}))
					})

					It("allows us getting the font ascent", func() {
						FPDFFont_GetAscent, err := PdfiumInstance.FPDFFont_GetAscent(&requests.FPDFFont_GetAscent{
							Font:     font,
							FontSize: 16,
						})
						Expect(err).To(BeNil())
						Expect(FPDFFont_GetAscent).To(Equal(&responses.FPDFFont_GetAscent{
							Ascent: 14.255999565124512,
						}))
					})

					It("allows us getting the font descent", func() {
						FPDFFont_GetDescent, err := PdfiumInstance.FPDFFont_GetDescent(&requests.FPDFFont_GetDescent{
							Font:     font,
							FontSize: 16,
						})
						Expect(err).To(BeNil())
						Expect(FPDFFont_GetDescent).To(Equal(&responses.FPDFFont_GetDescent{
							Descent: -3.4560000896453857,
						}))
					})

					It("allows us getting the font glypth width", func() {
						FPDFFont_GetGlyphWidth, err := PdfiumInstance.FPDFFont_GetGlyphWidth(&requests.FPDFFont_GetGlyphWidth{
							Font:     font,
							FontSize: 16,
							Glyph:    's',
						})
						Expect(err).To(BeNil())
						Expect(FPDFFont_GetGlyphWidth).To(Equal(&responses.FPDFFont_GetGlyphWidth{
							GlyphWidth: 6.223999977111816,
						}))
					})

					It("gives an error when getting the glyph path of an unsupported glyp", func() {
						FPDFFont_GetGlyphPath, err := PdfiumInstance.FPDFFont_GetGlyphPath(&requests.FPDFFont_GetGlyphPath{
							Font:     font,
							FontSize: 16,
							Glyph:    1,
						})
						Expect(err).To(MatchError("could not get glyph path"))
						Expect(FPDFFont_GetGlyphPath).To(BeNil())
					})

					It("returns an error when getting an unknown glyph path segment", func() {
						FPDFFont_GetGlyphPath, err := PdfiumInstance.FPDFFont_GetGlyphPath(&requests.FPDFFont_GetGlyphPath{
							Font:     font,
							FontSize: 16,
							Glyph:    's',
						})
						Expect(err).To(BeNil())
						Expect(FPDFFont_GetGlyphPath).ToNot(BeNil())
						Expect(FPDFFont_GetGlyphPath.GlyphPath).ToNot(BeEmpty())

						FPDFGlyphPath_GetGlyphPathSegment, err := PdfiumInstance.FPDFGlyphPath_GetGlyphPathSegment(&requests.FPDFGlyphPath_GetGlyphPathSegment{
							GlyphPath: FPDFFont_GetGlyphPath.GlyphPath,
							Index:     100,
						})
						Expect(err).To(MatchError("could not get glyph path segment"))
						Expect(FPDFGlyphPath_GetGlyphPathSegment).To(BeNil())
					})

					It("allows us getting the font glypth path, we can count the path segments and get the segment", func() {
						FPDFFont_GetGlyphPath, err := PdfiumInstance.FPDFFont_GetGlyphPath(&requests.FPDFFont_GetGlyphPath{
							Font:     font,
							FontSize: 16,
							Glyph:    's',
						})
						Expect(err).To(BeNil())
						Expect(FPDFFont_GetGlyphPath).ToNot(BeNil())
						Expect(FPDFFont_GetGlyphPath.GlyphPath).ToNot(BeEmpty())

						FPDFGlyphPath_CountGlyphSegments, err := PdfiumInstance.FPDFGlyphPath_CountGlyphSegments(&requests.FPDFGlyphPath_CountGlyphSegments{
							GlyphPath: FPDFFont_GetGlyphPath.GlyphPath,
						})
						Expect(err).To(BeNil())
						Expect(FPDFGlyphPath_CountGlyphSegments).To(Equal(&responses.FPDFGlyphPath_CountGlyphSegments{
							Count: 74,
						}))

						FPDFGlyphPath_GetGlyphPathSegment, err := PdfiumInstance.FPDFGlyphPath_GetGlyphPathSegment(&requests.FPDFGlyphPath_GetGlyphPathSegment{
							GlyphPath: FPDFFont_GetGlyphPath.GlyphPath,
							Index:     1,
						})
						Expect(err).To(BeNil())
						Expect(FPDFGlyphPath_GetGlyphPathSegment).ToNot(BeNil())
						Expect(FPDFGlyphPath_GetGlyphPathSegment.GlyphPathSegment).ToNot(BeEmpty())
					})
				})
			})
		})
	})

	Context("a PDF file with dashed lines", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/dashed_lines.pdf")
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
			Context("when an page object has been loaded", func() {
				var pageObject references.FPDF_PAGEOBJECT

				BeforeEach(func() {
					FPDFPage_GetObject, err := PdfiumInstance.FPDFPage_GetObject(&requests.FPDFPage_GetObject{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
						Index: 1,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPage_GetObject).To(Not(BeNil()))
					Expect(FPDFPage_GetObject.PageObject).To(Not(BeEmpty()))
					pageObject = FPDFPage_GetObject.PageObject

					FPDFPageObj_GetType, err := PdfiumInstance.FPDFPageObj_GetType(&requests.FPDFPageObj_GetType{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetType).To(Not(BeNil()))
					Expect(FPDFPageObj_GetType).To(Equal(&responses.FPDFPageObj_GetType{
						Type: enums.FPDF_PAGEOBJ_PATH,
					}))
				})

				It("allows getting the dash phase of the path", func() {
					FPDFPageObj_GetDashPhase, err := PdfiumInstance.FPDFPageObj_GetDashPhase(&requests.FPDFPageObj_GetDashPhase{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetDashPhase).To(Equal(&responses.FPDFPageObj_GetDashPhase{
						DashPhase: 5,
					}))
				})

				It("allows settings the dash phase of the path", func() {
					FPDFPageObj_SetDashPhase, err := PdfiumInstance.FPDFPageObj_SetDashPhase(&requests.FPDFPageObj_SetDashPhase{
						PageObject: pageObject,
						DashPhase:  4,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_SetDashPhase).To(Equal(&responses.FPDFPageObj_SetDashPhase{}))
				})

				It("allows settings the dash phase of the path and then getting it again", func() {
					FPDFPageObj_GetDashPhase, err := PdfiumInstance.FPDFPageObj_GetDashPhase(&requests.FPDFPageObj_GetDashPhase{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetDashPhase).To(Equal(&responses.FPDFPageObj_GetDashPhase{
						DashPhase: 5,
					}))

					FPDFPageObj_SetDashPhase, err := PdfiumInstance.FPDFPageObj_SetDashPhase(&requests.FPDFPageObj_SetDashPhase{
						PageObject: pageObject,
						DashPhase:  4,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_SetDashPhase).To(Equal(&responses.FPDFPageObj_SetDashPhase{}))

					FPDFPageObj_GetDashPhase, err = PdfiumInstance.FPDFPageObj_GetDashPhase(&requests.FPDFPageObj_GetDashPhase{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetDashPhase).To(Equal(&responses.FPDFPageObj_GetDashPhase{
						DashPhase: 4,
					}))
				})

				It("allows settings the dash count of the path", func() {
					FPDFPageObj_GetDashCount, err := PdfiumInstance.FPDFPageObj_GetDashCount(&requests.FPDFPageObj_GetDashCount{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetDashCount).To(Equal(&responses.FPDFPageObj_GetDashCount{
						DashCount: 6,
					}))
				})

				It("allows settings the dash array", func() {
					FPDFPageObj_GetDashArray, err := PdfiumInstance.FPDFPageObj_GetDashArray(&requests.FPDFPageObj_GetDashArray{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetDashArray).To(Equal(&responses.FPDFPageObj_GetDashArray{
						DashArray: []float32{6, 5, 4, 3, 2, 1},
					}))
				})

				It("allows settings the dash array", func() {
					FPDFPageObj_SetDashArray, err := PdfiumInstance.FPDFPageObj_SetDashArray(&requests.FPDFPageObj_SetDashArray{
						PageObject: pageObject,
						DashArray:  []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_SetDashArray).To(Equal(&responses.FPDFPageObj_SetDashArray{}))
				})

				It("allows settings the dash array and then getting it again", func() {
					FPDFPageObj_GetDashArray, err := PdfiumInstance.FPDFPageObj_GetDashArray(&requests.FPDFPageObj_GetDashArray{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetDashArray).To(Equal(&responses.FPDFPageObj_GetDashArray{
						DashArray: []float32{6, 5, 4, 3, 2, 1},
					}))

					FPDFPageObj_SetDashArray, err := PdfiumInstance.FPDFPageObj_SetDashArray(&requests.FPDFPageObj_SetDashArray{
						PageObject: pageObject,
						DashArray:  []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_SetDashArray).To(Equal(&responses.FPDFPageObj_SetDashArray{}))

					FPDFPageObj_GetDashArray, err = PdfiumInstance.FPDFPageObj_GetDashArray(&requests.FPDFPageObj_GetDashArray{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetDashArray).To(Equal(&responses.FPDFPageObj_GetDashArray{
						DashArray: []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					}))
				})
			})
		})
	})

	Context("a PDF file with images", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/embedded_images.pdf")
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
			Context("when an image object has been loaded", func() {
				var imageObject references.FPDF_PAGEOBJECT

				BeforeEach(func() {
					FPDFPage_GetObject, err := PdfiumInstance.FPDFPage_GetObject(&requests.FPDFPage_GetObject{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
						Index: 33,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPage_GetObject).To(Not(BeNil()))
					Expect(FPDFPage_GetObject.PageObject).To(Not(BeEmpty()))
					imageObject = FPDFPage_GetObject.PageObject

					FPDFPageObj_GetType, err := PdfiumInstance.FPDFPageObj_GetType(&requests.FPDFPageObj_GetType{
						PageObject: imageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetType).To(Not(BeNil()))
					Expect(FPDFPageObj_GetType).To(Equal(&responses.FPDFPageObj_GetType{
						Type: enums.FPDF_PAGEOBJ_IMAGE,
					}))
				})

				It("returns an error when trying to get the rendered bitmap without giving the page", func() {
					FPDFImageObj_GetRenderedBitmap, err := PdfiumInstance.FPDFImageObj_GetRenderedBitmap(&requests.FPDFImageObj_GetRenderedBitmap{
						Document: doc,
					})

					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FPDFImageObj_GetRenderedBitmap).To(BeNil())
				})

				It("returns an error when trying to get the rendered bitmap without giving the image object", func() {
					FPDFImageObj_GetRenderedBitmap, err := PdfiumInstance.FPDFImageObj_GetRenderedBitmap(&requests.FPDFImageObj_GetRenderedBitmap{
						Document: doc,
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})

					Expect(err).To(MatchError("pageObject not given"))
					Expect(FPDFImageObj_GetRenderedBitmap).To(BeNil())
				})

				It("returns the correct rendered bitmap", func() {
					FPDFImageObj_GetRenderedBitmap, err := PdfiumInstance.FPDFImageObj_GetRenderedBitmap(&requests.FPDFImageObj_GetRenderedBitmap{
						Document: doc,
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
						ImageObject: imageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFImageObj_GetRenderedBitmap).To(Not(BeNil()))
					Expect(FPDFImageObj_GetRenderedBitmap.Bitmap).To(Not(BeEmpty()))

					FPDFBitmap_Destroy, err := PdfiumInstance.FPDFBitmap_Destroy(&requests.FPDFBitmap_Destroy{
						Bitmap: FPDFImageObj_GetRenderedBitmap.Bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_Destroy).To(Equal(&responses.FPDFBitmap_Destroy{}))
				})

				It("gives an error when trying to remove an object without giving the page object", func() {
					FPDFPage_RemoveObject, err := PdfiumInstance.FPDFPage_RemoveObject(&requests.FPDFPage_RemoveObject{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})
					Expect(err).To(MatchError("pageObject not given"))
					Expect(FPDFPage_RemoveObject).To(BeNil())
				})

				It("allows the object to be removed from the page", func() {
					FPDFPage_RemoveObject, err := PdfiumInstance.FPDFPage_RemoveObject(&requests.FPDFPage_RemoveObject{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
						PageObject: imageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPage_RemoveObject).To(Equal(&responses.FPDFPage_RemoveObject{}))
				})

				It("allows us to get a matrix of an object", func() {
					FPDFPageObj_GetMatrix, err := PdfiumInstance.FPDFPageObj_GetMatrix(&requests.FPDFPageObj_GetMatrix{
						PageObject: imageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetMatrix).To(Equal(&responses.FPDFPageObj_GetMatrix{
						Matrix: structs.FPDF_FS_MATRIX{A: 53, B: 0, C: 0, D: 43, E: 72, F: 646.510009765625},
					}))
				})

				It("allows us to set a matrix of an object", func() {
					FPDFPageObj_SetMatrix, err := PdfiumInstance.FPDFPageObj_SetMatrix(&requests.FPDFPageObj_SetMatrix{
						PageObject: imageObject,
						Transform:  structs.FPDF_FS_MATRIX{A: 40, B: 0, C: 0, D: 50, E: 80, F: 646.510009765625},
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_SetMatrix).To(Equal(&responses.FPDFPageObj_SetMatrix{}))
				})

				It("allows us to set a matrix of an object and then get it again", func() {
					FPDFPageObj_GetMatrix, err := PdfiumInstance.FPDFPageObj_GetMatrix(&requests.FPDFPageObj_GetMatrix{
						PageObject: imageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetMatrix).To(Equal(&responses.FPDFPageObj_GetMatrix{
						Matrix: structs.FPDF_FS_MATRIX{A: 53, B: 0, C: 0, D: 43, E: 72, F: 646.510009765625},
					}))

					FPDFPageObj_SetMatrix, err := PdfiumInstance.FPDFPageObj_SetMatrix(&requests.FPDFPageObj_SetMatrix{
						PageObject: imageObject,
						Transform:  structs.FPDF_FS_MATRIX{A: 40, B: 0, C: 0, D: 50, E: 80, F: 646.510009765625},
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_SetMatrix).To(Equal(&responses.FPDFPageObj_SetMatrix{}))

					FPDFPageObj_GetMatrix, err = PdfiumInstance.FPDFPageObj_GetMatrix(&requests.FPDFPageObj_GetMatrix{
						PageObject: imageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetMatrix).To(Equal(&responses.FPDFPageObj_GetMatrix{
						Matrix: structs.FPDF_FS_MATRIX{A: 40, B: 0, C: 0, D: 50, E: 80, F: 646.510009765625},
					}))
				})
			})
		})
	})

	Context("a PDF file with marks", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/text_in_page_marked_indirect.pdf")
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
			When("a page object is opened", func() {
				var pageObject references.FPDF_PAGEOBJECT
				BeforeEach(func() {
					FPDFPage_GetObject, err := PdfiumInstance.FPDFPage_GetObject(&requests.FPDFPage_GetObject{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
						Index: 0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPage_GetObject).To(Not(BeNil()))
					Expect(FPDFPage_GetObject.PageObject).To(Not(BeEmpty()))
					pageObject = FPDFPage_GetObject.PageObject
				})

				It("returns the correct marks", func() {
					By("loading the page object")

					By("getting the mark count")
					FPDFPageObj_CountMarks, err := PdfiumInstance.FPDFPageObj_CountMarks(&requests.FPDFPageObj_CountMarks{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_CountMarks).To(Equal(&responses.FPDFPageObj_CountMarks{
						Count: 1,
					}))

					By("loading the mark")
					FPDFPageObj_GetMark, err := PdfiumInstance.FPDFPageObj_GetMark(&requests.FPDFPageObj_GetMark{
						PageObject: pageObject,
						Index:      0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetMark).To(Not(BeNil()))
					Expect(FPDFPageObj_GetMark.Mark).To(Not(BeEmpty()))

					By("getting the mark name")
					FPDFPageObjMark_GetName, err := PdfiumInstance.FPDFPageObjMark_GetName(&requests.FPDFPageObjMark_GetName{
						PageObjectMark: FPDFPageObj_GetMark.Mark,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObjMark_GetName).To(Equal(&responses.FPDFPageObjMark_GetName{
						Name: "Square",
					}))

					By("getting the mark params count")
					FPDFPageObjMark_CountParams, err := PdfiumInstance.FPDFPageObjMark_CountParams(&requests.FPDFPageObjMark_CountParams{
						PageObjectMark: FPDFPageObj_GetMark.Mark,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObjMark_CountParams).To(Equal(&responses.FPDFPageObjMark_CountParams{
						Count: 1,
					}))

					By("getting the mark params key")
					FPDFPageObjMark_GetParamKey, err := PdfiumInstance.FPDFPageObjMark_GetParamKey(&requests.FPDFPageObjMark_GetParamKey{
						PageObjectMark: FPDFPageObj_GetMark.Mark,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObjMark_GetParamKey).To(Equal(&responses.FPDFPageObjMark_GetParamKey{
						Key: "Factor",
					}))

					By("getting the mark param value type")
					FPDFPageObjMark_GetParamValueType, err := PdfiumInstance.FPDFPageObjMark_GetParamValueType(&requests.FPDFPageObjMark_GetParamValueType{
						PageObjectMark: FPDFPageObj_GetMark.Mark,
						Key:            "Factor",
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObjMark_GetParamValueType).To(Equal(&responses.FPDFPageObjMark_GetParamValueType{
						ValueType: enums.FPDF_OBJECT_TYPE_NUMBER,
					}))

					By("getting the mark param value")
					FPDFPageObjMark_GetParamIntValue, err := PdfiumInstance.FPDFPageObjMark_GetParamIntValue(&requests.FPDFPageObjMark_GetParamIntValue{
						PageObjectMark: FPDFPageObj_GetMark.Mark,
						Key:            "Factor",
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObjMark_GetParamIntValue).To(Equal(&responses.FPDFPageObjMark_GetParamIntValue{
						Value: 1,
					}))
				})

				It("allows adding a mark", func() {
					FPDFPageObj_AddMark, err := PdfiumInstance.FPDFPageObj_AddMark(&requests.FPDFPageObj_AddMark{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_AddMark).To(Not(BeNil()))
					Expect(FPDFPageObj_AddMark.Mark).To(Not(BeEmpty()))
				})

				It("gives an error when trying to get a mark with an invalid index", func() {
					FPDFPageObj_GetMark, err := PdfiumInstance.FPDFPageObj_GetMark(&requests.FPDFPageObj_GetMark{
						PageObject: pageObject,
						Index:      23,
					})
					Expect(err).To(MatchError("could not get mark"))
					Expect(FPDFPageObj_GetMark).To(BeNil())
				})

				It("gives an error when trying to remove a mark without giving the mark", func() {
					FPDFPageObj_RemoveMark, err := PdfiumInstance.FPDFPageObj_RemoveMark(&requests.FPDFPageObj_RemoveMark{
						PageObject: pageObject,
					})
					Expect(err).To(MatchError("pageObjectMark not given"))
					Expect(FPDFPageObj_RemoveMark).To(BeNil())
				})

				It("allows removing a mark", func() {
					FPDFPageObj_AddMark, err := PdfiumInstance.FPDFPageObj_AddMark(&requests.FPDFPageObj_AddMark{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_AddMark).To(Not(BeNil()))
					Expect(FPDFPageObj_AddMark.Mark).To(Not(BeEmpty()))

					FPDFPageObj_RemoveMark, err := PdfiumInstance.FPDFPageObj_RemoveMark(&requests.FPDFPageObj_RemoveMark{
						PageObject:     pageObject,
						PageObjectMark: FPDFPageObj_AddMark.Mark,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_RemoveMark).To(Equal(&responses.FPDFPageObj_RemoveMark{}))
				})

				It("gives an error when missing/incorrect mark params", func() {
					FPDFPageObj_AddMark, err := PdfiumInstance.FPDFPageObj_AddMark(&requests.FPDFPageObj_AddMark{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_AddMark).To(Not(BeNil()))
					Expect(FPDFPageObj_AddMark.Mark).To(Not(BeEmpty()))

					FPDFPageObjMark_SetIntParam, err := PdfiumInstance.FPDFPageObjMark_SetIntParam(&requests.FPDFPageObjMark_SetIntParam{
						Document: doc,
					})
					Expect(err).To(MatchError("pageObject not given"))
					Expect(FPDFPageObjMark_SetIntParam).To(BeNil())

					FPDFPageObjMark_SetIntParam, err = PdfiumInstance.FPDFPageObjMark_SetIntParam(&requests.FPDFPageObjMark_SetIntParam{
						Document:   doc,
						PageObject: pageObject,
					})
					Expect(err).To(MatchError("pageObjectMark not given"))
					Expect(FPDFPageObjMark_SetIntParam).To(BeNil())

					FPDFPageObjMark_SetStringParam, err := PdfiumInstance.FPDFPageObjMark_SetStringParam(&requests.FPDFPageObjMark_SetStringParam{
						Document: doc,
					})
					Expect(err).To(MatchError("pageObject not given"))
					Expect(FPDFPageObjMark_SetStringParam).To(BeNil())

					FPDFPageObjMark_SetStringParam, err = PdfiumInstance.FPDFPageObjMark_SetStringParam(&requests.FPDFPageObjMark_SetStringParam{
						Document:   doc,
						PageObject: pageObject,
					})
					Expect(err).To(MatchError("pageObjectMark not given"))
					Expect(FPDFPageObjMark_SetStringParam).To(BeNil())

					FPDFPageObjMark_SetBlobParam, err := PdfiumInstance.FPDFPageObjMark_SetBlobParam(&requests.FPDFPageObjMark_SetBlobParam{
						Document: doc,
					})
					Expect(err).To(MatchError("pageObject not given"))
					Expect(FPDFPageObjMark_SetBlobParam).To(BeNil())

					FPDFPageObjMark_SetBlobParam, err = PdfiumInstance.FPDFPageObjMark_SetBlobParam(&requests.FPDFPageObjMark_SetBlobParam{
						Document:   doc,
						PageObject: pageObject,
					})
					Expect(err).To(MatchError("pageObjectMark not given"))
					Expect(FPDFPageObjMark_SetBlobParam).To(BeNil())

					FPDFPageObjMark_SetBlobParam, err = PdfiumInstance.FPDFPageObjMark_SetBlobParam(&requests.FPDFPageObjMark_SetBlobParam{
						Document:       doc,
						PageObject:     pageObject,
						PageObjectMark: FPDFPageObj_AddMark.Mark,
					})
					Expect(err).To(MatchError("blob value cant be empty"))
					Expect(FPDFPageObjMark_SetBlobParam).To(BeNil())

					FPDFPageObjMark_RemoveParam, err := PdfiumInstance.FPDFPageObjMark_RemoveParam(&requests.FPDFPageObjMark_RemoveParam{
						PageObject: pageObject,
					})
					Expect(err).To(MatchError("pageObjectMark not given"))
					Expect(FPDFPageObjMark_RemoveParam).To(BeNil())

					FPDFPageObjMark_GetParamIntValue, err := PdfiumInstance.FPDFPageObjMark_GetParamIntValue(&requests.FPDFPageObjMark_GetParamIntValue{
						PageObjectMark: FPDFPageObj_AddMark.Mark,
						Key:            "not-existing",
					})
					Expect(err).To(MatchError("could not get value"))
					Expect(FPDFPageObjMark_GetParamIntValue).To(BeNil())

					FPDFPageObjMark_GetParamStringValue, err := PdfiumInstance.FPDFPageObjMark_GetParamStringValue(&requests.FPDFPageObjMark_GetParamStringValue{
						PageObjectMark: FPDFPageObj_AddMark.Mark,
						Key:            "not-existing",
					})
					Expect(err).To(MatchError("could not get value"))
					Expect(FPDFPageObjMark_GetParamStringValue).To(BeNil())

					FPDFPageObjMark_GetParamBlobValue, err := PdfiumInstance.FPDFPageObjMark_GetParamBlobValue(&requests.FPDFPageObjMark_GetParamBlobValue{
						PageObjectMark: FPDFPageObj_AddMark.Mark,
						Key:            "not-existing",
					})
					Expect(err).To(MatchError("could not get value"))
					Expect(FPDFPageObjMark_GetParamBlobValue).To(BeNil())

					FPDFPageObjMark_GetParamKey, err := PdfiumInstance.FPDFPageObjMark_GetParamKey(&requests.FPDFPageObjMark_GetParamKey{
						PageObjectMark: FPDFPageObj_AddMark.Mark,
						Index:          35,
					})
					Expect(err).To(MatchError("could not get key"))
					Expect(FPDFPageObjMark_GetParamKey).To(BeNil())
				})

				It("allows manipulation of a mark", func() {
					FPDFPageObj_AddMark, err := PdfiumInstance.FPDFPageObj_AddMark(&requests.FPDFPageObj_AddMark{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_AddMark).To(Not(BeNil()))
					Expect(FPDFPageObj_AddMark.Mark).To(Not(BeEmpty()))

					FPDFPageObjMark_SetIntParam, err := PdfiumInstance.FPDFPageObjMark_SetIntParam(&requests.FPDFPageObjMark_SetIntParam{
						Document:       doc,
						PageObject:     pageObject,
						PageObjectMark: FPDFPageObj_AddMark.Mark,
						Key:            "Test",
						Value:          1,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObjMark_SetIntParam).To(Equal(&responses.FPDFPageObjMark_SetIntParam{}))

					FPDFPageObjMark_GetParamIntValue, err := PdfiumInstance.FPDFPageObjMark_GetParamIntValue(&requests.FPDFPageObjMark_GetParamIntValue{
						PageObjectMark: FPDFPageObj_AddMark.Mark,
						Key:            "Test",
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObjMark_GetParamIntValue).To(Equal(&responses.FPDFPageObjMark_GetParamIntValue{
						Value: 1,
					}))

					FPDFPageObjMark_SetStringParam, err := PdfiumInstance.FPDFPageObjMark_SetStringParam(&requests.FPDFPageObjMark_SetStringParam{
						Document:       doc,
						PageObject:     pageObject,
						PageObjectMark: FPDFPageObj_AddMark.Mark,
						Key:            "Test2",
						Value:          "test",
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObjMark_SetStringParam).To(Equal(&responses.FPDFPageObjMark_SetStringParam{}))

					FPDFPageObjMark_GetParamStringValue, err := PdfiumInstance.FPDFPageObjMark_GetParamStringValue(&requests.FPDFPageObjMark_GetParamStringValue{
						PageObjectMark: FPDFPageObj_AddMark.Mark,
						Key:            "Test2",
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObjMark_GetParamStringValue).To(Equal(&responses.FPDFPageObjMark_GetParamStringValue{
						Value: "test",
					}))

					FPDFPageObjMark_SetBlobParam, err := PdfiumInstance.FPDFPageObjMark_SetBlobParam(&requests.FPDFPageObjMark_SetBlobParam{
						Document:       doc,
						PageObject:     pageObject,
						PageObjectMark: FPDFPageObj_AddMark.Mark,
						Key:            "Test3",
						Value:          []byte{1, 2, 3},
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObjMark_SetBlobParam).To(Equal(&responses.FPDFPageObjMark_SetBlobParam{}))

					FPDFPageObjMark_GetParamBlobValue, err := PdfiumInstance.FPDFPageObjMark_GetParamBlobValue(&requests.FPDFPageObjMark_GetParamBlobValue{
						PageObjectMark: FPDFPageObj_AddMark.Mark,
						Key:            "Test3",
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObjMark_GetParamBlobValue).To(Equal(&responses.FPDFPageObjMark_GetParamBlobValue{
						Value: []byte{1, 2, 3},
					}))

					FPDFPageObjMark_RemoveParam, err := PdfiumInstance.FPDFPageObjMark_RemoveParam(&requests.FPDFPageObjMark_RemoveParam{
						PageObject:     pageObject,
						PageObjectMark: FPDFPageObj_AddMark.Mark,
						Key:            "Test3",
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObjMark_RemoveParam).To(Equal(&responses.FPDFPageObjMark_RemoveParam{}))
				})
			})
		})
	})
})
