package shared_tests

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"
	"io/ioutil"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func RunfpdfTextTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("fpdf_text", func() {
		Context("no page is given", func() {
			It("returns an error when calling FPDFText_LoadPage", func() {
				FPDFText_LoadPage, err := pdfiumContainer.FPDFText_LoadPage(&requests.FPDFText_LoadPage{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDFText_LoadPage).To(BeNil())
			})
		})

		Context("no text page is given", func() {
			It("returns an error when calling FPDFText_LoadPage", func() {
				FPDFText_ClosePage, err := pdfiumContainer.FPDFText_ClosePage(&requests.FPDFText_ClosePage{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_ClosePage).To(BeNil())
			})

			It("returns an error when calling FPDFText_CountChars", func() {
				FPDFText_CountChars, err := pdfiumContainer.FPDFText_CountChars(&requests.FPDFText_CountChars{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_CountChars).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetUnicode", func() {
				FPDFText_GetUnicode, err := pdfiumContainer.FPDFText_ClosePage(&requests.FPDFText_ClosePage{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetUnicode).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetFontSize", func() {
				FPDFText_GetFontSize, err := pdfiumContainer.FPDFText_GetFontSize(&requests.FPDFText_GetFontSize{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetFontSize).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetFontInfo", func() {
				FPDFText_GetFontInfo, err := pdfiumContainer.FPDFText_GetFontInfo(&requests.FPDFText_GetFontInfo{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetFontInfo).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetFontWeight", func() {
				FPDFText_GetFontWeight, err := pdfiumContainer.FPDFText_GetFontWeight(&requests.FPDFText_GetFontWeight{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetFontWeight).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetTextRenderMode", func() {
				FPDFText_GetTextRenderMode, err := pdfiumContainer.FPDFText_GetTextRenderMode(&requests.FPDFText_GetTextRenderMode{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetTextRenderMode).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetFillColor", func() {
				FPDFText_GetFillColor, err := pdfiumContainer.FPDFText_GetFillColor(&requests.FPDFText_GetFillColor{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetFillColor).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetStrokeColor", func() {
				FPDFText_GetStrokeColor, err := pdfiumContainer.FPDFText_GetStrokeColor(&requests.FPDFText_GetStrokeColor{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetStrokeColor).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetCharAngle", func() {
				FPDFText_GetCharAngle, err := pdfiumContainer.FPDFText_GetCharAngle(&requests.FPDFText_GetCharAngle{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetCharAngle).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetCharBox", func() {
				FPDFText_GetCharBox, err := pdfiumContainer.FPDFText_GetCharBox(&requests.FPDFText_GetCharBox{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetCharBox).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetLooseCharBox", func() {
				FPDFText_GetLooseCharBox, err := pdfiumContainer.FPDFText_GetLooseCharBox(&requests.FPDFText_GetLooseCharBox{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetLooseCharBox).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetMatrix", func() {
				FPDFText_GetMatrix, err := pdfiumContainer.FPDFText_GetMatrix(&requests.FPDFText_GetMatrix{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetMatrix).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetCharOrigin", func() {
				FPDFText_GetCharOrigin, err := pdfiumContainer.FPDFText_GetCharOrigin(&requests.FPDFText_GetCharOrigin{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetCharOrigin).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetCharIndexAtPos", func() {
				FPDFText_GetCharIndexAtPos, err := pdfiumContainer.FPDFText_GetCharIndexAtPos(&requests.FPDFText_GetCharIndexAtPos{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetCharIndexAtPos).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetText", func() {
				FPDFText_GetText, err := pdfiumContainer.FPDFText_GetText(&requests.FPDFText_GetText{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetText).To(BeNil())
			})

			It("returns an error when calling FPDFText_CountRects", func() {
				FPDFText_CountRects, err := pdfiumContainer.FPDFText_CountRects(&requests.FPDFText_CountRects{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_CountRects).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetRect", func() {
				FPDFText_GetRect, err := pdfiumContainer.FPDFText_GetRect(&requests.FPDFText_GetRect{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetRect).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetBoundedText", func() {
				FPDFText_GetBoundedText, err := pdfiumContainer.FPDFText_GetBoundedText(&requests.FPDFText_GetBoundedText{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_GetBoundedText).To(BeNil())
			})

			It("returns an error when calling FPDFText_FindStart", func() {
				FPDFText_FindStart, err := pdfiumContainer.FPDFText_FindStart(&requests.FPDFText_FindStart{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFText_FindStart).To(BeNil())
			})

			It("returns an error when calling FPDFLink_LoadWebLinks", func() {
				FPDFLink_LoadWebLinks, err := pdfiumContainer.FPDFLink_LoadWebLinks(&requests.FPDFLink_LoadWebLinks{})
				Expect(err).To(MatchError("textPage not given"))
				Expect(FPDFLink_LoadWebLinks).To(BeNil())
			})
		})

		Context("no search handle is given", func() {
			It("returns an error when calling FPDFText_FindNext", func() {
				FPDFText_FindNext, err := pdfiumContainer.FPDFText_FindNext(&requests.FPDFText_FindNext{})
				Expect(err).To(MatchError("search not given"))
				Expect(FPDFText_FindNext).To(BeNil())
			})

			It("returns an error when calling FPDFText_FindPrev", func() {
				FPDFText_FindPrev, err := pdfiumContainer.FPDFText_FindPrev(&requests.FPDFText_FindPrev{})
				Expect(err).To(MatchError("search not given"))
				Expect(FPDFText_FindPrev).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetSchResultIndex", func() {
				FPDFText_GetSchResultIndex, err := pdfiumContainer.FPDFText_GetSchResultIndex(&requests.FPDFText_GetSchResultIndex{})
				Expect(err).To(MatchError("search not given"))
				Expect(FPDFText_GetSchResultIndex).To(BeNil())
			})

			It("returns an error when calling FPDFText_GetSchCount", func() {
				FPDFText_GetSchCount, err := pdfiumContainer.FPDFText_GetSchCount(&requests.FPDFText_GetSchCount{})
				Expect(err).To(MatchError("search not given"))
				Expect(FPDFText_GetSchCount).To(BeNil())
			})

			It("returns an error when calling FPDFText_FindClose", func() {
				FPDFText_FindClose, err := pdfiumContainer.FPDFText_FindClose(&requests.FPDFText_FindClose{})
				Expect(err).To(MatchError("search not given"))
				Expect(FPDFText_FindClose).To(BeNil())
			})
		})

		Context("no page link is given", func() {
			It("returns an error when calling FPDFLink_CountWebLinks", func() {
				FPDFLink_CountWebLinks, err := pdfiumContainer.FPDFLink_CountWebLinks(&requests.FPDFLink_CountWebLinks{})
				Expect(err).To(MatchError("pageLink not given"))
				Expect(FPDFLink_CountWebLinks).To(BeNil())
			})

			It("returns an error when calling FPDFLink_GetURL", func() {
				FPDFLink_GetURL, err := pdfiumContainer.FPDFLink_GetURL(&requests.FPDFLink_GetURL{})
				Expect(err).To(MatchError("pageLink not given"))
				Expect(FPDFLink_GetURL).To(BeNil())
			})

			It("returns an error when calling FPDFLink_CountRects", func() {
				FPDFLink_CountRects, err := pdfiumContainer.FPDFLink_CountRects(&requests.FPDFLink_CountRects{})
				Expect(err).To(MatchError("pageLink not given"))
				Expect(FPDFLink_CountRects).To(BeNil())
			})

			It("returns an error when calling FPDFLink_GetRect", func() {
				FPDFLink_GetRect, err := pdfiumContainer.FPDFLink_GetRect(&requests.FPDFLink_GetRect{})
				Expect(err).To(MatchError("pageLink not given"))
				Expect(FPDFLink_GetRect).To(BeNil())
			})

			It("returns an error when calling FPDFLink_GetTextRange", func() {
				FPDFLink_GetTextRange, err := pdfiumContainer.FPDFLink_GetTextRange(&requests.FPDFLink_GetTextRange{})
				Expect(err).To(MatchError("pageLink not given"))
				Expect(FPDFLink_GetTextRange).To(BeNil())
			})

			It("returns an error when calling FPDFLink_CloseWebLinks", func() {
				FPDFLink_CloseWebLinks, err := pdfiumContainer.FPDFLink_CloseWebLinks(&requests.FPDFLink_CloseWebLinks{})
				Expect(err).To(MatchError("pageLink not given"))
				Expect(FPDFLink_CloseWebLinks).To(BeNil())
			})
		})

		Context("a normal PDF file", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				pdfData, err := ioutil.ReadFile(testsPath + "/testdata/test.pdf")
				Expect(err).To(BeNil())

				newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData)
				Expect(err).To(BeNil())

				doc = *newDoc
			})

			AfterEach(func() {
				err := pdfiumContainer.FPDF_CloseDocument(doc)
				Expect(err).To(BeNil())
			})

			When("a text page is opened", func() {
				var textPage references.FPDF_TEXTPAGE

				BeforeEach(func() {
					textPageResp, err := pdfiumContainer.FPDFText_LoadPage(&requests.FPDFText_LoadPage{
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

				AfterEach(func() {
					resp, err := pdfiumContainer.FPDFText_ClosePage(&requests.FPDFText_ClosePage{
						TextPage: textPage,
					})
					Expect(err).To(BeNil())
					Expect(resp).To(Not(BeNil()))
				})

				It("returns the correct character count", func() {
					FPDFText_CountChars, err := pdfiumContainer.FPDFText_CountChars(&requests.FPDFText_CountChars{
						TextPage: textPage,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_CountChars).To(Equal(&responses.FPDFText_CountChars{
						Count: 57,
					}))
				})

				It("returns the correct unicode for char 0", func() {
					FPDFText_GetUnicode, err := pdfiumContainer.FPDFText_GetUnicode(&requests.FPDFText_GetUnicode{
						TextPage: textPage,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetUnicode).To(Equal(&responses.FPDFText_GetUnicode{
						Index:   0,
						Unicode: 70,
					}))
				})

				It("returns the correct font size for char 0", func() {
					FPDFText_GetFontSize, err := pdfiumContainer.FPDFText_GetFontSize(&requests.FPDFText_GetFontSize{
						TextPage: textPage,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetFontSize).To(Equal(&responses.FPDFText_GetFontSize{
						Index:    0,
						FontSize: 1,
					}))
				})

				It("returns the correct font info for char 0", func() {
					FPDFText_GetFontInfo, err := pdfiumContainer.FPDFText_GetFontInfo(&requests.FPDFText_GetFontInfo{
						TextPage: textPage,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetFontInfo).To(Equal(&responses.FPDFText_GetFontInfo{
						Index:    0,
						FontName: "CGKWYO+DejaVuSans",
						Flags:    524320,
					}))
				})

				It("returns an error when getting the font info for an invalid char", func() {
					FPDFText_GetFontInfo, err := pdfiumContainer.FPDFText_GetFontInfo(&requests.FPDFText_GetFontInfo{
						TextPage: textPage,
						Index:    -1,
					})
					Expect(err).To(MatchError("could not get font name"))
					Expect(FPDFText_GetFontInfo).To(BeNil())
				})

				It("returns the correct font weight for char 0", func() {
					FPDFText_GetFontWeight, err := pdfiumContainer.FPDFText_GetFontWeight(&requests.FPDFText_GetFontWeight{
						TextPage: textPage,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetFontWeight).To(Equal(&responses.FPDFText_GetFontWeight{
						Index:      0,
						FontWeight: 400,
					}))
				})

				It("returns an error when getting the font weight for an invalid char", func() {
					FPDFText_GetFontWeight, err := pdfiumContainer.FPDFText_GetFontWeight(&requests.FPDFText_GetFontWeight{
						TextPage: textPage,
						Index:    -1,
					})
					Expect(err).To(MatchError("could not get font weight"))
					Expect(FPDFText_GetFontWeight).To(BeNil())
				})

				It("returns the correct text render mode for char 0", func() {
					FPDFText_GetTextRenderMode, err := pdfiumContainer.FPDFText_GetTextRenderMode(&requests.FPDFText_GetTextRenderMode{
						TextPage: textPage,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetTextRenderMode).To(Equal(&responses.FPDFText_GetTextRenderMode{
						Index:          0,
						TextRenderMode: enums.FPDF_TEXTRENDERMODE_FILL,
					}))
				})

				It("returns the correct text fill color for char 0", func() {
					FPDFText_GetFillColor, err := pdfiumContainer.FPDFText_GetFillColor(&requests.FPDFText_GetFillColor{
						TextPage: textPage,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetFillColor).To(Equal(&responses.FPDFText_GetFillColor{
						Index: 0,
						R:     0,
						G:     0,
						B:     0,
						A:     255,
					}))
				})

				It("returns an error when getting the text fill color for an invalid char", func() {
					FPDFText_GetFillColor, err := pdfiumContainer.FPDFText_GetFillColor(&requests.FPDFText_GetFillColor{
						TextPage: textPage,
						Index:    -1,
					})
					Expect(err).To(MatchError("could not get fill color"))
					Expect(FPDFText_GetFillColor).To(BeNil())
				})

				It("returns the correct text stroke color for char 0", func() {
					FPDFText_GetStrokeColor, err := pdfiumContainer.FPDFText_GetStrokeColor(&requests.FPDFText_GetStrokeColor{
						TextPage: textPage,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetStrokeColor).To(Equal(&responses.FPDFText_GetStrokeColor{
						Index: 0,
						R:     0,
						G:     0,
						B:     0,
						A:     255,
					}))
				})

				It("returns an error when getting the text stroke color for an invalid char", func() {
					FPDFText_GetStrokeColor, err := pdfiumContainer.FPDFText_GetStrokeColor(&requests.FPDFText_GetStrokeColor{
						TextPage: textPage,
						Index:    -1,
					})
					Expect(err).To(MatchError("could not get stroke color"))
					Expect(FPDFText_GetStrokeColor).To(BeNil())
				})

				It("returns the correct char angle for char 0", func() {
					FPDFText_GetCharAngle, err := pdfiumContainer.FPDFText_GetCharAngle(&requests.FPDFText_GetCharAngle{
						TextPage: textPage,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetCharAngle).To(Equal(&responses.FPDFText_GetCharAngle{
						Index:     0,
						CharAngle: 0,
					}))
				})

				It("returns an error when getting the char angle for an invalid char", func() {
					FPDFText_GetCharAngle, err := pdfiumContainer.FPDFText_GetCharAngle(&requests.FPDFText_GetCharAngle{
						TextPage: textPage,
						Index:    -1,
					})
					Expect(err).To(MatchError("could not get char angle"))
					Expect(FPDFText_GetCharAngle).To(BeNil())
				})

				It("returns the correct char box for char 0", func() {
					FPDFText_GetCharBox, err := pdfiumContainer.FPDFText_GetCharBox(&requests.FPDFText_GetCharBox{
						TextPage: textPage,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetCharBox).To(Equal(&responses.FPDFText_GetCharBox{
						Index:  0,
						Left:   71.9451904296875,
						Right:  76.55418395996094,
						Bottom: 789.1592407226562,
						Top:    797.17822265625,
					}))
				})

				It("returns an error when getting the char box for an invalid char", func() {
					FPDFText_GetCharBox, err := pdfiumContainer.FPDFText_GetCharBox(&requests.FPDFText_GetCharBox{
						TextPage: textPage,
						Index:    -1,
					})
					Expect(err).To(MatchError("could not get char box"))
					Expect(FPDFText_GetCharBox).To(BeNil())
				})

				It("returns the correct loose char box for char 0", func() {
					FPDFText_GetLooseCharBox, err := pdfiumContainer.FPDFText_GetLooseCharBox(&requests.FPDFText_GetLooseCharBox{
						TextPage: textPage,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetLooseCharBox).To(Equal(&responses.FPDFText_GetLooseCharBox{
						Index: 0,
						Rect: structs.FPDF_FS_RECTF{
							Left:   70.8671875,
							Top:    797.9365234375,
							Right:  77.19218444824219,
							Bottom: 786.9365234375,
						},
					}))
				})

				It("returns an error when getting the loose char box for an invalid char", func() {
					FPDFText_GetLooseCharBox, err := pdfiumContainer.FPDFText_GetLooseCharBox(&requests.FPDFText_GetLooseCharBox{
						TextPage: textPage,
						Index:    -1,
					})
					Expect(err).To(MatchError("could not get loose char box"))
					Expect(FPDFText_GetLooseCharBox).To(BeNil())
				})

				It("returns the correct char matrix for char 0", func() {
					FPDFText_GetMatrix, err := pdfiumContainer.FPDFText_GetMatrix(&requests.FPDFText_GetMatrix{
						TextPage: textPage,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetMatrix).To(Equal(&responses.FPDFText_GetMatrix{
						Index: 0,
						Matrix: structs.FPDF_FS_MATRIX{
							A: 11,
							B: 0,
							C: 0,
							D: 11,
							E: 70.8671875,
							F: 789.1592407226562,
						},
					}))
				})

				It("returns an error when getting the char matrix for an invalid char", func() {
					FPDFText_GetMatrix, err := pdfiumContainer.FPDFText_GetMatrix(&requests.FPDFText_GetMatrix{
						TextPage: textPage,
						Index:    -1,
					})
					Expect(err).To(MatchError("could not get char matrix"))
					Expect(FPDFText_GetMatrix).To(BeNil())
				})

				It("returns the correct char origin for char 0", func() {
					FPDFText_GetCharOrigin, err := pdfiumContainer.FPDFText_GetCharOrigin(&requests.FPDFText_GetCharOrigin{
						TextPage: textPage,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetCharOrigin).To(Equal(&responses.FPDFText_GetCharOrigin{
						Index: 0,
						X:     70.8671875,
						Y:     789.1592407226562,
					}))
				})

				It("returns an error when getting the char origin for an invalid char", func() {
					FPDFText_GetCharOrigin, err := pdfiumContainer.FPDFText_GetCharOrigin(&requests.FPDFText_GetCharOrigin{
						TextPage: textPage,
						Index:    -1,
					})
					Expect(err).To(MatchError("could not get char origin"))
					Expect(FPDFText_GetCharOrigin).To(BeNil())
				})

				It("returns the correct char index for the given position", func() {
					FPDFText_GetCharIndexAtPos, err := pdfiumContainer.FPDFText_GetCharIndexAtPos(&requests.FPDFText_GetCharIndexAtPos{
						TextPage: textPage,
						X:        73,
						Y:        793,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetCharIndexAtPos).To(Equal(&responses.FPDFText_GetCharIndexAtPos{
						CharIndex: 0,
					}))
				})

				It("returns the correct char index for no text at position", func() {
					FPDFText_GetCharIndexAtPos, err := pdfiumContainer.FPDFText_GetCharIndexAtPos(&requests.FPDFText_GetCharIndexAtPos{
						TextPage: textPage,
						X:        2,
						Y:        2,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetCharIndexAtPos).To(Equal(&responses.FPDFText_GetCharIndexAtPos{
						CharIndex: -1,
					}))
				})

				It("returns the correct page text", func() {
					FPDFText_GetText, err := pdfiumContainer.FPDFText_GetText(&requests.FPDFText_GetText{
						TextPage:   textPage,
						StartIndex: 0,
						Count:      57,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetText).To(Equal(&responses.FPDFText_GetText{
						Text: "File: Untitled Document 2 Page 1 of 1\r\nThis is a test PDF",
					}))
				})

				It("returns the correct rect count", func() {
					FPDFText_CountRects, err := pdfiumContainer.FPDFText_CountRects(&requests.FPDFText_CountRects{
						TextPage:   textPage,
						StartIndex: 0,
						Count:      57,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_CountRects).To(Equal(&responses.FPDFText_CountRects{
						Count: 3,
					}))
				})

				It("returns an error when getting a rect without calling FPDFText_CountRects first", func() {
					FPDFText_GetRect, err := pdfiumContainer.FPDFText_GetRect(&requests.FPDFText_GetRect{
						TextPage: textPage,
						Index:    0,
					})
					Expect(err).To(MatchError("could not get rect"))
					Expect(FPDFText_GetRect).To(BeNil())
				})

				It("returns the correct position for rect 0", func() {
					FPDFText_CountRects, err := pdfiumContainer.FPDFText_CountRects(&requests.FPDFText_CountRects{
						TextPage:   textPage,
						StartIndex: 0,
						Count:      57,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_CountRects).To(Equal(&responses.FPDFText_CountRects{
						Count: 3,
					}))

					// This only works if FPDFText_CountRects is called first.
					FPDFText_GetRect, err := pdfiumContainer.FPDFText_GetRect(&requests.FPDFText_GetRect{
						TextPage: textPage,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetRect).To(Equal(&responses.FPDFText_GetRect{
						Left:   71.9451904296875,
						Top:    797.5192260742188,
						Right:  208.60919189453125,
						Bottom: 789.0162353515625,
					}))
				})

				It("returns the correct text within bounds", func() {
					FPDFText_GetBoundedText, err := pdfiumContainer.FPDFText_GetBoundedText(&requests.FPDFText_GetBoundedText{
						TextPage: textPage,
						Left:     71.9451904296875,
						Top:      797.5192260742188,
						Right:    208.60919189453125,
						Bottom:   789.0162353515625,
					})
					Expect(err).To(BeNil())
					Expect(FPDFText_GetBoundedText).To(Equal(&responses.FPDFText_GetBoundedText{
						Text: "File: Untitled Document 2 ",
					}))
				})

				When("a search handle is opened", func() {
					When("moving forwards", func() {
						var searchHandle references.FPDF_SCHHANDLE

						BeforeEach(func() {
							FPDFText_FindStart, err := pdfiumContainer.FPDFText_FindStart(&requests.FPDFText_FindStart{
								TextPage: textPage,
								Find:     "e",
							})
							Expect(err).To(BeNil())
							Expect(FPDFText_FindStart).To(Not(BeNil()))

							searchHandle = FPDFText_FindStart.Search
						})

						AfterEach(func() {
							resp, err := pdfiumContainer.FPDFText_FindClose(&requests.FPDFText_FindClose{
								Search: searchHandle,
							})
							Expect(err).To(BeNil())
							Expect(resp).To(Not(BeNil()))
						})

						It("returns the correct amount of characters matches", func() {
							// We need to call FindNext first to move to the first result.
							FPDFText_FindNext, err := pdfiumContainer.FPDFText_FindNext(&requests.FPDFText_FindNext{
								Search: searchHandle,
							})

							Expect(err).To(BeNil())
							Expect(FPDFText_FindNext).To(Equal(&responses.FPDFText_FindNext{
								GotMatch: true,
							}))

							FPDFText_GetSchCount, err := pdfiumContainer.FPDFText_GetSchCount(&requests.FPDFText_GetSchCount{
								Search: searchHandle,
							})
							Expect(err).To(BeNil())
							Expect(FPDFText_GetSchCount).To(Equal(&responses.FPDFText_GetSchCount{
								Count: 1,
							}))
						})

						It("returns the correct char position of the first match", func() {
							// We need to call FindNext first to move to the first result.
							FPDFText_FindNext, err := pdfiumContainer.FPDFText_FindNext(&requests.FPDFText_FindNext{
								Search: searchHandle,
							})

							Expect(err).To(BeNil())
							Expect(FPDFText_FindNext).To(Equal(&responses.FPDFText_FindNext{
								GotMatch: true,
							}))

							FPDFText_GetSchResultIndex, err := pdfiumContainer.FPDFText_GetSchResultIndex(&requests.FPDFText_GetSchResultIndex{
								Search: searchHandle,
							})
							Expect(err).To(BeNil())
							Expect(FPDFText_GetSchResultIndex).To(Equal(&responses.FPDFText_GetSchResultIndex{
								Index: 3,
							}))
						})
					})

					When("moving backwards", func() {
						var searchHandle references.FPDF_SCHHANDLE

						BeforeEach(func() {
							FPDFText_FindStart, err := pdfiumContainer.FPDFText_FindStart(&requests.FPDFText_FindStart{
								TextPage:   textPage,
								Find:       "e",
								StartIndex: -1,
							})
							Expect(err).To(BeNil())
							Expect(FPDFText_FindStart).To(Not(BeNil()))

							searchHandle = FPDFText_FindStart.Search
						})

						AfterEach(func() {
							resp, err := pdfiumContainer.FPDFText_FindClose(&requests.FPDFText_FindClose{
								Search: searchHandle,
							})
							Expect(err).To(BeNil())
							Expect(resp).To(Not(BeNil()))
						})

						It("returns the correct amount of characters matches", func() {
							// We need to call FindNext first to move to the first result.
							FPDFText_FindNext, err := pdfiumContainer.FPDFText_FindPrev(&requests.FPDFText_FindPrev{
								Search: searchHandle,
							})

							Expect(err).To(BeNil())
							Expect(FPDFText_FindNext).To(Equal(&responses.FPDFText_FindPrev{
								GotMatch: true,
							}))

							FPDFText_GetSchCount, err := pdfiumContainer.FPDFText_GetSchCount(&requests.FPDFText_GetSchCount{
								Search: searchHandle,
							})
							Expect(err).To(BeNil())
							Expect(FPDFText_GetSchCount).To(Equal(&responses.FPDFText_GetSchCount{
								Count: 1,
							}))
						})

						It("returns the correct char position of the first match", func() {
							// We need to call FindNext first to move to the first result.
							FPDFText_FindNext, err := pdfiumContainer.FPDFText_FindPrev(&requests.FPDFText_FindPrev{
								Search: searchHandle,
							})

							Expect(err).To(BeNil())
							Expect(FPDFText_FindNext).To(Equal(&responses.FPDFText_FindPrev{
								GotMatch: true,
							}))

							FPDFText_GetSchResultIndex, err := pdfiumContainer.FPDFText_GetSchResultIndex(&requests.FPDFText_GetSchResultIndex{
								Search: searchHandle,
							})
							Expect(err).To(BeNil())
							Expect(FPDFText_GetSchResultIndex).To(Equal(&responses.FPDFText_GetSchResultIndex{
								Index: 50,
							}))
						})
					})
				})
			})
		})

		Context("a PDF file with weblinks", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				pdfData, err := ioutil.ReadFile(testsPath + "/testdata/weblinks.pdf")
				Expect(err).To(BeNil())

				newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData)
				Expect(err).To(BeNil())

				doc = *newDoc
			})

			AfterEach(func() {
				err := pdfiumContainer.FPDF_CloseDocument(doc)
				Expect(err).To(BeNil())
			})

			When("a text page is opened", func() {
				var textPage references.FPDF_TEXTPAGE

				BeforeEach(func() {
					textPageResp, err := pdfiumContainer.FPDFText_LoadPage(&requests.FPDFText_LoadPage{
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

				AfterEach(func() {
					resp, err := pdfiumContainer.FPDFText_ClosePage(&requests.FPDFText_ClosePage{
						TextPage: textPage,
					})
					Expect(err).To(BeNil())
					Expect(resp).To(Not(BeNil()))
				})

				When("a web links are loaded", func() {
					var pageLink references.FPDF_PAGELINK

					BeforeEach(func() {
						FPDFLink_LoadWebLinks, err := pdfiumContainer.FPDFLink_LoadWebLinks(&requests.FPDFLink_LoadWebLinks{
							TextPage: textPage,
						})
						Expect(err).To(BeNil())
						Expect(FPDFLink_LoadWebLinks).To(Not(BeNil()))

						pageLink = FPDFLink_LoadWebLinks.PageLink
					})

					AfterEach(func() {
						resp, err := pdfiumContainer.FPDFLink_CloseWebLinks(&requests.FPDFLink_CloseWebLinks{
							PageLink: pageLink,
						})
						Expect(err).To(BeNil())
						Expect(resp).To(Not(BeNil()))
					})

					It("returns the correct web link count", func() {
						FPDFLink_CountWebLinks, err := pdfiumContainer.FPDFLink_CountWebLinks(&requests.FPDFLink_CountWebLinks{
							PageLink: pageLink,
						})
						Expect(err).To(BeNil())
						Expect(FPDFLink_CountWebLinks).To(Equal(&responses.FPDFLink_CountWebLinks{
							Count: 2,
						}))
					})

					It("returns the correct URL for link 0", func() {
						FPDFLink_GetURL, err := pdfiumContainer.FPDFLink_GetURL(&requests.FPDFLink_GetURL{
							PageLink: pageLink,
							Index:    0,
						})
						Expect(err).To(BeNil())
						Expect(FPDFLink_GetURL).To(Equal(&responses.FPDFLink_GetURL{
							Index: 0,
							URL:   "http://example.com?q=foo",
						}))
					})

					It("returns the correct rect count for link 0", func() {
						FPDFLink_CountRects, err := pdfiumContainer.FPDFLink_CountRects(&requests.FPDFLink_CountRects{
							PageLink: pageLink,
							Index:    0,
						})
						Expect(err).To(BeNil())
						Expect(FPDFLink_CountRects).To(Equal(&responses.FPDFLink_CountRects{
							Index: 0,
							Count: 1,
						}))
					})

					It("returns the correct rect for link 0 and rect index 0", func() {
						FPDFLink_GetRect, err := pdfiumContainer.FPDFLink_GetRect(&requests.FPDFLink_GetRect{
							PageLink:  pageLink,
							Index:     0,
							RectIndex: 0,
						})
						Expect(err).To(BeNil())
						Expect(FPDFLink_GetRect).To(Equal(&responses.FPDFLink_GetRect{
							Index:  0,
							Left:   50.779998779296875,
							Top:    108.84400177001953,
							Right:  187.9879913330078,
							Bottom: 97.52799987792969,
						}))
					})

					It("returns the correct text range for link 0 and gets the correct text for it", func() {
						FPDFLink_GetTextRange, err := pdfiumContainer.FPDFLink_GetTextRange(&requests.FPDFLink_GetTextRange{
							PageLink: pageLink,
							Index:    0,
						})
						Expect(err).To(BeNil())
						Expect(FPDFLink_GetTextRange).To(Equal(&responses.FPDFLink_GetTextRange{
							Index:          0,
							StartCharIndex: 35,
							CharCount:      24,
						}))

						FPDFText_GetText, err := pdfiumContainer.FPDFText_GetText(&requests.FPDFText_GetText{
							TextPage:   textPage,
							StartIndex: FPDFLink_GetTextRange.StartCharIndex,
							Count:      FPDFLink_GetTextRange.CharCount,
						})
						Expect(err).To(BeNil())
						Expect(FPDFText_GetText).To(Equal(&responses.FPDFText_GetText{
							Text: "http://example.com?q=foo",
						}))
					})
				})
			})
		})
	})
}
