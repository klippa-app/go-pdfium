package shared_tests

import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_text", func() {
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
		It("returns an error when calling FPDFText_LoadPage", func() {
			FPDFText_LoadPage, err := PdfiumInstance.FPDFText_LoadPage(&requests.FPDFText_LoadPage{})
			Expect(err).To(MatchError("either page reference or index should be given"))
			Expect(FPDFText_LoadPage).To(BeNil())
		})
	})

	Context("no text page is given", func() {
		It("returns an error when calling FPDFText_LoadPage", func() {
			FPDFText_ClosePage, err := PdfiumInstance.FPDFText_ClosePage(&requests.FPDFText_ClosePage{})
			Expect(err).To(MatchError("textPage not given"))
			Expect(FPDFText_ClosePage).To(BeNil())
		})

		It("returns an error when calling FPDFText_CountChars", func() {
			FPDFText_CountChars, err := PdfiumInstance.FPDFText_CountChars(&requests.FPDFText_CountChars{})
			Expect(err).To(MatchError("textPage not given"))
			Expect(FPDFText_CountChars).To(BeNil())
		})

		It("returns an error when calling FPDFText_GetUnicode", func() {
			FPDFText_GetUnicode, err := PdfiumInstance.FPDFText_ClosePage(&requests.FPDFText_ClosePage{})
			Expect(err).To(MatchError("textPage not given"))
			Expect(FPDFText_GetUnicode).To(BeNil())
		})

		It("returns an error when calling FPDFText_GetFontSize", func() {
			FPDFText_GetFontSize, err := PdfiumInstance.FPDFText_GetFontSize(&requests.FPDFText_GetFontSize{})
			Expect(err).To(MatchError("textPage not given"))
			Expect(FPDFText_GetFontSize).To(BeNil())
		})

		It("returns an error when calling FPDFText_GetCharBox", func() {
			FPDFText_GetCharBox, err := PdfiumInstance.FPDFText_GetCharBox(&requests.FPDFText_GetCharBox{})
			Expect(err).To(MatchError("textPage not given"))
			Expect(FPDFText_GetCharBox).To(BeNil())
		})

		It("returns an error when calling FPDFText_GetCharOrigin", func() {
			FPDFText_GetCharOrigin, err := PdfiumInstance.FPDFText_GetCharOrigin(&requests.FPDFText_GetCharOrigin{})
			Expect(err).To(MatchError("textPage not given"))
			Expect(FPDFText_GetCharOrigin).To(BeNil())
		})

		It("returns an error when calling FPDFText_GetCharIndexAtPos", func() {
			FPDFText_GetCharIndexAtPos, err := PdfiumInstance.FPDFText_GetCharIndexAtPos(&requests.FPDFText_GetCharIndexAtPos{})
			Expect(err).To(MatchError("textPage not given"))
			Expect(FPDFText_GetCharIndexAtPos).To(BeNil())
		})

		It("returns an error when calling FPDFText_GetText", func() {
			FPDFText_GetText, err := PdfiumInstance.FPDFText_GetText(&requests.FPDFText_GetText{})
			Expect(err).To(MatchError("textPage not given"))
			Expect(FPDFText_GetText).To(BeNil())
		})

		It("returns an error when calling FPDFText_CountRects", func() {
			FPDFText_CountRects, err := PdfiumInstance.FPDFText_CountRects(&requests.FPDFText_CountRects{})
			Expect(err).To(MatchError("textPage not given"))
			Expect(FPDFText_CountRects).To(BeNil())
		})

		It("returns an error when calling FPDFText_GetRect", func() {
			FPDFText_GetRect, err := PdfiumInstance.FPDFText_GetRect(&requests.FPDFText_GetRect{})
			Expect(err).To(MatchError("textPage not given"))
			Expect(FPDFText_GetRect).To(BeNil())
		})

		It("returns an error when calling FPDFText_GetBoundedText", func() {
			FPDFText_GetBoundedText, err := PdfiumInstance.FPDFText_GetBoundedText(&requests.FPDFText_GetBoundedText{})
			Expect(err).To(MatchError("textPage not given"))
			Expect(FPDFText_GetBoundedText).To(BeNil())
		})

		It("returns an error when calling FPDFText_FindStart", func() {
			FPDFText_FindStart, err := PdfiumInstance.FPDFText_FindStart(&requests.FPDFText_FindStart{})
			Expect(err).To(MatchError("textPage not given"))
			Expect(FPDFText_FindStart).To(BeNil())
		})

		It("returns an error when calling FPDFLink_LoadWebLinks", func() {
			FPDFLink_LoadWebLinks, err := PdfiumInstance.FPDFLink_LoadWebLinks(&requests.FPDFLink_LoadWebLinks{})
			Expect(err).To(MatchError("textPage not given"))
			Expect(FPDFLink_LoadWebLinks).To(BeNil())
		})
	})

	Context("no search handle is given", func() {
		It("returns an error when calling FPDFText_FindNext", func() {
			FPDFText_FindNext, err := PdfiumInstance.FPDFText_FindNext(&requests.FPDFText_FindNext{})
			Expect(err).To(MatchError("search not given"))
			Expect(FPDFText_FindNext).To(BeNil())
		})

		It("returns an error when calling FPDFText_FindPrev", func() {
			FPDFText_FindPrev, err := PdfiumInstance.FPDFText_FindPrev(&requests.FPDFText_FindPrev{})
			Expect(err).To(MatchError("search not given"))
			Expect(FPDFText_FindPrev).To(BeNil())
		})

		It("returns an error when calling FPDFText_GetSchResultIndex", func() {
			FPDFText_GetSchResultIndex, err := PdfiumInstance.FPDFText_GetSchResultIndex(&requests.FPDFText_GetSchResultIndex{})
			Expect(err).To(MatchError("search not given"))
			Expect(FPDFText_GetSchResultIndex).To(BeNil())
		})

		It("returns an error when calling FPDFText_GetSchCount", func() {
			FPDFText_GetSchCount, err := PdfiumInstance.FPDFText_GetSchCount(&requests.FPDFText_GetSchCount{})
			Expect(err).To(MatchError("search not given"))
			Expect(FPDFText_GetSchCount).To(BeNil())
		})

		It("returns an error when calling FPDFText_FindClose", func() {
			FPDFText_FindClose, err := PdfiumInstance.FPDFText_FindClose(&requests.FPDFText_FindClose{})
			Expect(err).To(MatchError("search not given"))
			Expect(FPDFText_FindClose).To(BeNil())
		})
	})

	Context("no page link is given", func() {
		It("returns an error when calling FPDFLink_CountWebLinks", func() {
			FPDFLink_CountWebLinks, err := PdfiumInstance.FPDFLink_CountWebLinks(&requests.FPDFLink_CountWebLinks{})
			Expect(err).To(MatchError("pageLink not given"))
			Expect(FPDFLink_CountWebLinks).To(BeNil())
		})

		It("returns an error when calling FPDFLink_GetURL", func() {
			FPDFLink_GetURL, err := PdfiumInstance.FPDFLink_GetURL(&requests.FPDFLink_GetURL{})
			Expect(err).To(MatchError("pageLink not given"))
			Expect(FPDFLink_GetURL).To(BeNil())
		})

		It("returns an error when calling FPDFLink_CountRects", func() {
			FPDFLink_CountRects, err := PdfiumInstance.FPDFLink_CountRects(&requests.FPDFLink_CountRects{})
			Expect(err).To(MatchError("pageLink not given"))
			Expect(FPDFLink_CountRects).To(BeNil())
		})

		It("returns an error when calling FPDFLink_GetRect", func() {
			FPDFLink_GetRect, err := PdfiumInstance.FPDFLink_GetRect(&requests.FPDFLink_GetRect{})
			Expect(err).To(MatchError("pageLink not given"))
			Expect(FPDFLink_GetRect).To(BeNil())
		})

		It("returns an error when calling FPDFLink_CloseWebLinks", func() {
			FPDFLink_CloseWebLinks, err := PdfiumInstance.FPDFLink_CloseWebLinks(&requests.FPDFLink_CloseWebLinks{})
			Expect(err).To(MatchError("pageLink not given"))
			Expect(FPDFLink_CloseWebLinks).To(BeNil())
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

			AfterEach(func() {
				resp, err := PdfiumInstance.FPDFText_ClosePage(&requests.FPDFText_ClosePage{
					TextPage: textPage,
				})
				Expect(err).To(BeNil())
				Expect(resp).To(Not(BeNil()))
			})

			It("returns the correct character count", func() {
				FPDFText_CountChars, err := PdfiumInstance.FPDFText_CountChars(&requests.FPDFText_CountChars{
					TextPage: textPage,
				})
				Expect(err).To(BeNil())
				Expect(FPDFText_CountChars).To(Equal(&responses.FPDFText_CountChars{
					Count: 57,
				}))
			})

			It("returns the correct unicode for char 0", func() {
				FPDFText_GetUnicode, err := PdfiumInstance.FPDFText_GetUnicode(&requests.FPDFText_GetUnicode{
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
				FPDFText_GetFontSize, err := PdfiumInstance.FPDFText_GetFontSize(&requests.FPDFText_GetFontSize{
					TextPage: textPage,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(FPDFText_GetFontSize).To(Equal(&responses.FPDFText_GetFontSize{
					Index:    0,
					FontSize: 1,
				}))
			})

			It("returns the correct char box for char 0", func() {
				FPDFText_GetCharBox, err := PdfiumInstance.FPDFText_GetCharBox(&requests.FPDFText_GetCharBox{
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
				FPDFText_GetCharBox, err := PdfiumInstance.FPDFText_GetCharBox(&requests.FPDFText_GetCharBox{
					TextPage: textPage,
					Index:    -1,
				})
				Expect(err).To(MatchError("could not get char box"))
				Expect(FPDFText_GetCharBox).To(BeNil())
			})

			It("returns the correct char origin for char 0", func() {
				FPDFText_GetCharOrigin, err := PdfiumInstance.FPDFText_GetCharOrigin(&requests.FPDFText_GetCharOrigin{
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
				FPDFText_GetCharOrigin, err := PdfiumInstance.FPDFText_GetCharOrigin(&requests.FPDFText_GetCharOrigin{
					TextPage: textPage,
					Index:    -1,
				})
				Expect(err).To(MatchError("could not get char origin"))
				Expect(FPDFText_GetCharOrigin).To(BeNil())
			})

			It("returns the correct char index for the given position", func() {
				FPDFText_GetCharIndexAtPos, err := PdfiumInstance.FPDFText_GetCharIndexAtPos(&requests.FPDFText_GetCharIndexAtPos{
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
				FPDFText_GetCharIndexAtPos, err := PdfiumInstance.FPDFText_GetCharIndexAtPos(&requests.FPDFText_GetCharIndexAtPos{
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
				FPDFText_GetText, err := PdfiumInstance.FPDFText_GetText(&requests.FPDFText_GetText{
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
				FPDFText_CountRects, err := PdfiumInstance.FPDFText_CountRects(&requests.FPDFText_CountRects{
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
				FPDFText_GetRect, err := PdfiumInstance.FPDFText_GetRect(&requests.FPDFText_GetRect{
					TextPage: textPage,
					Index:    0,
				})
				Expect(err).To(MatchError("could not get rect"))
				Expect(FPDFText_GetRect).To(BeNil())
			})

			It("returns the correct position for rect 0", func() {
				FPDFText_CountRects, err := PdfiumInstance.FPDFText_CountRects(&requests.FPDFText_CountRects{
					TextPage:   textPage,
					StartIndex: 0,
					Count:      57,
				})
				Expect(err).To(BeNil())
				Expect(FPDFText_CountRects).To(Equal(&responses.FPDFText_CountRects{
					Count: 3,
				}))

				// This only works if FPDFText_CountRects is called first.
				FPDFText_GetRect, err := PdfiumInstance.FPDFText_GetRect(&requests.FPDFText_GetRect{
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
				FPDFText_GetBoundedText, err := PdfiumInstance.FPDFText_GetBoundedText(&requests.FPDFText_GetBoundedText{
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
						FPDFText_FindStart, err := PdfiumInstance.FPDFText_FindStart(&requests.FPDFText_FindStart{
							TextPage: textPage,
							Find:     "e",
						})
						Expect(err).To(BeNil())
						Expect(FPDFText_FindStart).To(Not(BeNil()))

						searchHandle = FPDFText_FindStart.Search
					})

					AfterEach(func() {
						resp, err := PdfiumInstance.FPDFText_FindClose(&requests.FPDFText_FindClose{
							Search: searchHandle,
						})
						Expect(err).To(BeNil())
						Expect(resp).To(Not(BeNil()))
					})

					It("returns the correct amount of characters matches", func() {
						// We need to call FindNext first to move to the first result.
						FPDFText_FindNext, err := PdfiumInstance.FPDFText_FindNext(&requests.FPDFText_FindNext{
							Search: searchHandle,
						})

						Expect(err).To(BeNil())
						Expect(FPDFText_FindNext).To(Equal(&responses.FPDFText_FindNext{
							GotMatch: true,
						}))

						FPDFText_GetSchCount, err := PdfiumInstance.FPDFText_GetSchCount(&requests.FPDFText_GetSchCount{
							Search: searchHandle,
						})
						Expect(err).To(BeNil())
						Expect(FPDFText_GetSchCount).To(Equal(&responses.FPDFText_GetSchCount{
							Count: 1,
						}))
					})

					It("returns the correct char position of the first match", func() {
						// We need to call FindNext first to move to the first result.
						FPDFText_FindNext, err := PdfiumInstance.FPDFText_FindNext(&requests.FPDFText_FindNext{
							Search: searchHandle,
						})

						Expect(err).To(BeNil())
						Expect(FPDFText_FindNext).To(Equal(&responses.FPDFText_FindNext{
							GotMatch: true,
						}))

						FPDFText_GetSchResultIndex, err := PdfiumInstance.FPDFText_GetSchResultIndex(&requests.FPDFText_GetSchResultIndex{
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
						FPDFText_FindStart, err := PdfiumInstance.FPDFText_FindStart(&requests.FPDFText_FindStart{
							TextPage:   textPage,
							Find:       "e",
							StartIndex: -1,
						})
						Expect(err).To(BeNil())
						Expect(FPDFText_FindStart).To(Not(BeNil()))

						searchHandle = FPDFText_FindStart.Search
					})

					AfterEach(func() {
						resp, err := PdfiumInstance.FPDFText_FindClose(&requests.FPDFText_FindClose{
							Search: searchHandle,
						})
						Expect(err).To(BeNil())
						Expect(resp).To(Not(BeNil()))
					})

					It("returns the correct amount of characters matches", func() {
						// We need to call FindNext first to move to the first result.
						FPDFText_FindNext, err := PdfiumInstance.FPDFText_FindPrev(&requests.FPDFText_FindPrev{
							Search: searchHandle,
						})

						Expect(err).To(BeNil())
						Expect(FPDFText_FindNext).To(Equal(&responses.FPDFText_FindPrev{
							GotMatch: true,
						}))

						FPDFText_GetSchCount, err := PdfiumInstance.FPDFText_GetSchCount(&requests.FPDFText_GetSchCount{
							Search: searchHandle,
						})
						Expect(err).To(BeNil())
						Expect(FPDFText_GetSchCount).To(Equal(&responses.FPDFText_GetSchCount{
							Count: 1,
						}))
					})

					It("returns the correct char position of the first match", func() {
						// We need to call FindNext first to move to the first result.
						FPDFText_FindNext, err := PdfiumInstance.FPDFText_FindPrev(&requests.FPDFText_FindPrev{
							Search: searchHandle,
						})

						Expect(err).To(BeNil())
						Expect(FPDFText_FindNext).To(Equal(&responses.FPDFText_FindPrev{
							GotMatch: true,
						}))

						FPDFText_GetSchResultIndex, err := PdfiumInstance.FPDFText_GetSchResultIndex(&requests.FPDFText_GetSchResultIndex{
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
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/weblinks.pdf")
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

			AfterEach(func() {
				resp, err := PdfiumInstance.FPDFText_ClosePage(&requests.FPDFText_ClosePage{
					TextPage: textPage,
				})
				Expect(err).To(BeNil())
				Expect(resp).To(Not(BeNil()))
			})

			When("a web links are loaded", func() {
				var pageLink references.FPDF_PAGELINK

				BeforeEach(func() {
					FPDFLink_LoadWebLinks, err := PdfiumInstance.FPDFLink_LoadWebLinks(&requests.FPDFLink_LoadWebLinks{
						TextPage: textPage,
					})
					Expect(err).To(BeNil())
					Expect(FPDFLink_LoadWebLinks).To(Not(BeNil()))

					pageLink = FPDFLink_LoadWebLinks.PageLink
				})

				AfterEach(func() {
					resp, err := PdfiumInstance.FPDFLink_CloseWebLinks(&requests.FPDFLink_CloseWebLinks{
						PageLink: pageLink,
					})
					Expect(err).To(BeNil())
					Expect(resp).To(Not(BeNil()))
				})

				It("returns the correct web link count", func() {
					FPDFLink_CountWebLinks, err := PdfiumInstance.FPDFLink_CountWebLinks(&requests.FPDFLink_CountWebLinks{
						PageLink: pageLink,
					})
					Expect(err).To(BeNil())
					Expect(FPDFLink_CountWebLinks).To(Equal(&responses.FPDFLink_CountWebLinks{
						Count: 2,
					}))
				})

				It("returns the correct URL for link 0", func() {
					FPDFLink_GetURL, err := PdfiumInstance.FPDFLink_GetURL(&requests.FPDFLink_GetURL{
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
					FPDFLink_CountRects, err := PdfiumInstance.FPDFLink_CountRects(&requests.FPDFLink_CountRects{
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
					FPDFLink_GetRect, err := PdfiumInstance.FPDFLink_GetRect(&requests.FPDFLink_GetRect{
						PageLink:  pageLink,
						Index:     0,
						RectIndex: 0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFLink_GetRect).To(Not(BeNil()))
					Expect(FPDFLink_GetRect.Index).To(Equal(0))
					Expect(FPDFLink_GetRect.RectIndex).To(Equal(0))

					// Rect can be a little different depending on the platform.
					Expect(FPDFLink_GetRect.Left).To(BeNumerically("~", 50, 1))
					Expect(FPDFLink_GetRect.Top).To(BeNumerically("~", 108, 1))
					Expect(FPDFLink_GetRect.Right).To(BeNumerically("~", 187, 1))
					Expect(FPDFLink_GetRect.Bottom).To(BeNumerically("~", 97, 1))
				})
			})
		})
	})
})
