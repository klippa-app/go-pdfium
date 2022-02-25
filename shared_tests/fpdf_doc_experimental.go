//go:build pdfium_experimental
// +build pdfium_experimental

package shared_tests

import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("fpdf_doc_experimental", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when getting the page metadata", func() {
				FPDF_GetMetaText, err := PdfiumInstance.FPDF_GetMetaText(&requests.FPDF_GetMetaText{
					Tag: "Creator",
				})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_GetMetaText).To(BeNil())
			})

			It("returns an error when calling FPDFBookmark_GetFirstChild", func() {
				FPDFBookmark_GetFirstChild, err := PdfiumInstance.FPDFBookmark_GetFirstChild(&requests.FPDFBookmark_GetFirstChild{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFBookmark_GetFirstChild).To(BeNil())
			})

			It("returns an error when calling FPDFBookmark_GetNextSibling", func() {
				FPDFBookmark_GetNextSibling, err := PdfiumInstance.FPDFBookmark_GetNextSibling(&requests.FPDFBookmark_GetNextSibling{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFBookmark_GetNextSibling).To(BeNil())
			})

			It("returns an error when calling FPDFBookmark_Find", func() {
				FPDFBookmark_Find, err := PdfiumInstance.FPDFBookmark_Find(&requests.FPDFBookmark_Find{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFBookmark_Find).To(BeNil())
			})

			It("returns an error when calling FPDFAction_GetDest", func() {
				FPDFAction_GetDest, err := PdfiumInstance.FPDFAction_GetDest(&requests.FPDFAction_GetDest{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFAction_GetDest).To(BeNil())
			})

			It("returns an error when calling FPDFAction_GetURIPath", func() {
				FPDFAction_GetURIPath, err := PdfiumInstance.FPDFAction_GetURIPath(&requests.FPDFAction_GetURIPath{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFAction_GetURIPath).To(BeNil())
			})

			It("returns an error when calling FPDFDest_GetDestPageIndex", func() {
				FPDFDest_GetDestPageIndex, err := PdfiumInstance.FPDFDest_GetDestPageIndex(&requests.FPDFDest_GetDestPageIndex{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFDest_GetDestPageIndex).To(BeNil())
			})

			It("returns an error when calling FPDF_GetFileIdentifier", func() {
				FPDF_GetFileIdentifier, err := PdfiumInstance.FPDF_GetFileIdentifier(&requests.FPDF_GetFileIdentifier{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_GetFileIdentifier).To(BeNil())
			})

			It("returns an error when calling FPDF_GetPageLabel", func() {
				FPDF_GetPageLabel, err := PdfiumInstance.FPDF_GetPageLabel(&requests.FPDF_GetPageLabel{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_GetPageLabel).To(BeNil())
			})

			It("returns an error when calling FPDFLink_GetDest", func() {
				FPDFLink_GetDest, err := PdfiumInstance.FPDFLink_GetDest(&requests.FPDFLink_GetDest{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFLink_GetDest).To(BeNil())
			})

			It("returns an error when calling GetMetaData", func() {
				GetMetaData, err := PdfiumInstance.GetMetaData(&requests.GetMetaData{})
				Expect(err).To(MatchError("document not given"))
				Expect(GetMetaData).To(BeNil())
			})
		})
	})

	Context("no bookmark", func() {
		When("is given", func() {
			It("returns an error when calling FPDFBookmark_GetTitle", func() {
				FPDFBookmark_GetTitle, err := PdfiumInstance.FPDFBookmark_GetTitle(&requests.FPDFBookmark_GetTitle{})
				Expect(err).To(MatchError("bookmark not given"))
				Expect(FPDFBookmark_GetTitle).To(BeNil())
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
			Context("requesting metadata", func() {
				It("returns the correct metadata text", func() {
					metadata, err := PdfiumInstance.FPDF_GetMetaText(&requests.FPDF_GetMetaText{
						Document: doc,
						Tag:      "Producer",
					})
					Expect(err).To(BeNil())
					Expect(metadata).To(Equal(&responses.FPDF_GetMetaText{
						Tag:   "Producer",
						Value: "cairo 1.16.0 (https://cairographics.org)",
					}))
				})

				It("returns the correct metadata tag", func() {
					metadata, err := PdfiumInstance.GetMetaData(&requests.GetMetaData{
						Document: doc,
						Tags:     &[]string{"Producer"},
					})
					Expect(err).To(BeNil())
					Expect(metadata).To(Equal(&responses.GetMetaData{
						Tags: []responses.GetMetaDataTag{
							{
								Tag:   "Producer",
								Value: "cairo 1.16.0 (https://cairographics.org)",
							},
						},
					}))
				})

				It("returns the correct metadata tags when no tags were given", func() {
					metadata, err := PdfiumInstance.GetMetaData(&requests.GetMetaData{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(metadata).To(Equal(&responses.GetMetaData{
						Tags: []responses.GetMetaDataTag{
							{Tag: "Title", Value: ""},
							{Tag: "Author", Value: ""},
							{Tag: "Subject", Value: ""},
							{Tag: "Keywords", Value: ""},
							{Tag: "Creator", Value: ""},
							{
								Tag:   "Producer",
								Value: "cairo 1.16.0 (https://cairographics.org)",
							},
							{
								Tag:   "CreationDate",
								Value: "D:20210823145142+02'00",
							},
							{Tag: "ModDate", Value: ""},
						},
					}))
				})
			})

			Context("without giving an action", func() {
				It("FPDFAction_GetURIPath returns an error", func() {
					uriPath, err := PdfiumInstance.FPDFAction_GetURIPath(&requests.FPDFAction_GetURIPath{
						Document: doc,
					})
					Expect(err).To(MatchError("action not given"))
					Expect(uriPath).To(BeNil())
				})

				It("FPDFAction_GetFilePath returns an error", func() {
					filePath, err := PdfiumInstance.FPDFAction_GetFilePath(&requests.FPDFAction_GetFilePath{})
					Expect(err).To(MatchError("action not given"))
					Expect(filePath).To(BeNil())
				})

				It("FPDFAction_GetDest returns an error", func() {
					dest, err := PdfiumInstance.FPDFAction_GetDest(&requests.FPDFAction_GetDest{
						Document: doc,
					})
					Expect(err).To(MatchError("action not given"))
					Expect(dest).To(BeNil())
				})

				It("FPDFAction_GetType returns an error", func() {
					actionType, err := PdfiumInstance.FPDFAction_GetType(&requests.FPDFAction_GetType{})
					Expect(err).To(MatchError("action not given"))
					Expect(actionType).To(BeNil())
				})
			})

			Context("without giving a dest", func() {
				It("FPDFDest_GetView returns an error", func() {
					destView, err := PdfiumInstance.FPDFDest_GetView(&requests.FPDFDest_GetView{})
					Expect(err).To(MatchError("dest not given"))
					Expect(destView).To(BeNil())
				})

				It("FPDFDest_GetDestPageIndex returns an error", func() {
					destPageIndex, err := PdfiumInstance.FPDFDest_GetDestPageIndex(&requests.FPDFDest_GetDestPageIndex{
						Document: doc,
					})
					Expect(err).To(MatchError("dest not given"))
					Expect(destPageIndex).To(BeNil())
				})

				It("FPDFDest_GetLocationInPage returns an error", func() {
					location, err := PdfiumInstance.FPDFDest_GetLocationInPage(&requests.FPDFDest_GetLocationInPage{})
					Expect(err).To(MatchError("dest not given"))
					Expect(location).To(BeNil())
				})
			})

			Context("without giving a page", func() {
				It("FPDFLink_GetLinkAtPoint returns an error", func() {
					FPDFLink_GetLinkAtPoint, err := PdfiumInstance.FPDFLink_GetLinkAtPoint(&requests.FPDFLink_GetLinkAtPoint{})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FPDFLink_GetLinkAtPoint).To(BeNil())
				})

				It("FPDFLink_GetLinkZOrderAtPoint returns an error", func() {
					FPDFLink_GetLinkZOrderAtPoint, err := PdfiumInstance.FPDFLink_GetLinkZOrderAtPoint(&requests.FPDFLink_GetLinkZOrderAtPoint{})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FPDFLink_GetLinkZOrderAtPoint).To(BeNil())
				})

				It("FPDFLink_Enumerate returns an error", func() {
					FPDFLink_Enumerate, err := PdfiumInstance.FPDFLink_Enumerate(&requests.FPDFLink_Enumerate{})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FPDFLink_Enumerate).To(BeNil())
				})

				It("FPDFLink_GetAnnot returns an error", func() {
					FPDFLink_GetAnnot, err := PdfiumInstance.FPDFLink_GetAnnot(&requests.FPDFLink_GetAnnot{})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FPDFLink_GetAnnot).To(BeNil())
				})

				It("FPDF_GetPageAAction returns an error", func() {
					FPDF_GetPageAAction, err := PdfiumInstance.FPDF_GetPageAAction(&requests.FPDF_GetPageAAction{})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FPDF_GetPageAAction).To(BeNil())
				})
			})

			Context("without giving a link", func() {
				It("FPDFLink_GetDest returns an error", func() {
					FPDFLink_GetDest, err := PdfiumInstance.FPDFLink_GetDest(&requests.FPDFLink_GetDest{
						Document: doc,
					})
					Expect(err).To(MatchError("link not given"))
					Expect(FPDFLink_GetDest).To(BeNil())
				})

				It("FPDFLink_GetAction returns an error", func() {
					FPDFLink_GetAction, err := PdfiumInstance.FPDFLink_GetAction(&requests.FPDFLink_GetAction{})
					Expect(err).To(MatchError("link not given"))
					Expect(FPDFLink_GetAction).To(BeNil())
				})

				It("FPDFLink_GetAnnot returns an error", func() {
					FPDFLink_GetAnnot, err := PdfiumInstance.FPDFLink_GetAnnot(&requests.FPDFLink_GetAnnot{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})
					Expect(err).To(MatchError("link not given"))
					Expect(FPDFLink_GetAnnot).To(BeNil())
				})

				It("FPDFLink_GetAnnotRect returns an error", func() {
					FPDFLink_GetAnnotRect, err := PdfiumInstance.FPDFLink_GetAnnotRect(&requests.FPDFLink_GetAnnotRect{})
					Expect(err).To(MatchError("link not given"))
					Expect(FPDFLink_GetAnnotRect).To(BeNil())
				})

				It("FPDFLink_CountQuadPoints returns an error", func() {
					FPDFLink_CountQuadPoints, err := PdfiumInstance.FPDFLink_CountQuadPoints(&requests.FPDFLink_CountQuadPoints{})
					Expect(err).To(MatchError("link not given"))
					Expect(FPDFLink_CountQuadPoints).To(BeNil())
				})

				It("FPDFLink_GetQuadPoints returns an error", func() {
					FPDFLink_GetQuadPoints, err := PdfiumInstance.FPDFLink_GetQuadPoints(&requests.FPDFLink_GetQuadPoints{})
					Expect(err).To(MatchError("link not given"))
					Expect(FPDFLink_GetQuadPoints).To(BeNil())
				})
			})

			Context("without having a link", func() {
				It("FPDFLink_GetLinkAtPoint returns no link", func() {
					pageLink, err := PdfiumInstance.FPDFLink_GetLinkAtPoint(&requests.FPDFLink_GetLinkAtPoint{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
						X: 100,
						Y: 100,
					})
					Expect(err).To(BeNil())
					Expect(pageLink).To(Not(BeNil()))
					Expect(pageLink.Link).To(BeNil())
				})
			})
		})
	})

	Context("a PDF file with no bookmarks", func() {
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

		When("FPDFBookmark_GetFirstChild is called", func() {
			It("returns no bookmark", func() {
				fistChild, err := PdfiumInstance.FPDFBookmark_GetFirstChild(&requests.FPDFBookmark_GetFirstChild{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(fistChild).To(Equal(&responses.FPDFBookmark_GetFirstChild{}))
			})
		})

		When("FPDFBookmark_GetNextSibling is called without a bookmark", func() {
			It("returns an error", func() {
				nextSibling, err := PdfiumInstance.FPDFBookmark_GetNextSibling(&requests.FPDFBookmark_GetNextSibling{
					Document: doc,
				})
				Expect(err).To(MatchError("bookmark not given"))
				Expect(nextSibling).To(BeNil())
			})
		})

		When("FPDFBookmark_Find is called without a title", func() {
			It("returns an error", func() {
				bookmark, err := PdfiumInstance.FPDFBookmark_Find(&requests.FPDFBookmark_Find{
					Document: doc,
				})
				Expect(err).To(MatchError("no title given"))
				Expect(bookmark).To(BeNil())
			})
		})

		When("FPDFBookmark_Find is called with a title", func() {
			It("returns no bookmark", func() {
				bookmark, err := PdfiumInstance.FPDFBookmark_Find(&requests.FPDFBookmark_Find{
					Document: doc,
					Title:    "Can't find",
				})
				Expect(err).To(BeNil())
				Expect(bookmark).To(Equal(&responses.FPDFBookmark_Find{}))
			})
		})
	})

	Context("a PDF file with bookmarks", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/bookmarks.pdf")
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

		When("FPDFBookmark_GetFirstChild is called", func() {
			It("returns a bookmark with a matching tile, no children and 1 sibling", func() {
				topLevelBookmark, err := PdfiumInstance.FPDFBookmark_GetFirstChild(&requests.FPDFBookmark_GetFirstChild{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(topLevelBookmark).To(Not(BeNil()))

				Expect(topLevelBookmark.Bookmark).To(Not(BeNil()))
				// Check title of bookmark.
				topLevelBookmarkTitle, err := PdfiumInstance.FPDFBookmark_GetTitle(&requests.FPDFBookmark_GetTitle{
					Bookmark: *topLevelBookmark.Bookmark,
				})
				Expect(err).To(BeNil())
				Expect(topLevelBookmarkTitle).To(Equal(&responses.FPDFBookmark_GetTitle{
					Title: "A Good Beginning",
				}))

				// Check that we have no children
				topLevelBookmarkSibling, err := PdfiumInstance.FPDFBookmark_GetFirstChild(&requests.FPDFBookmark_GetFirstChild{
					Document: doc,
					Bookmark: topLevelBookmark.Bookmark,
				})
				Expect(err).To(BeNil())
				Expect(topLevelBookmarkSibling).To(Not(BeNil()))
				Expect(topLevelBookmarkSibling.Bookmark).To(BeNil())

				// Check that we have a sibling
				sibling, err := PdfiumInstance.FPDFBookmark_GetNextSibling(&requests.FPDFBookmark_GetNextSibling{
					Document: doc,
					Bookmark: *topLevelBookmark.Bookmark,
				})
				Expect(err).To(BeNil())
				Expect(sibling).To(Not(BeNil()))
				Expect(sibling.Bookmark).To(Not(BeNil()))

				// Check title of bookmark.
				siblingTitle, err := PdfiumInstance.FPDFBookmark_GetTitle(&requests.FPDFBookmark_GetTitle{
					Bookmark: *sibling.Bookmark,
				})
				Expect(err).To(BeNil())
				Expect(siblingTitle).To(Equal(&responses.FPDFBookmark_GetTitle{
					Title: "A Good Ending",
				}))

				// Check that we have no children
				siblingChildren, err := PdfiumInstance.FPDFBookmark_GetFirstChild(&requests.FPDFBookmark_GetFirstChild{
					Document: doc,
					Bookmark: sibling.Bookmark,
				})
				Expect(err).To(BeNil())
				Expect(siblingChildren).To(Not(BeNil()))
				Expect(siblingChildren.Bookmark).To(BeNil())

				// Check that we have no sibling
				siblingSibling, err := PdfiumInstance.FPDFBookmark_GetNextSibling(&requests.FPDFBookmark_GetNextSibling{
					Document: doc,
					Bookmark: *sibling.Bookmark,
				})
				Expect(err).To(BeNil())
				Expect(siblingSibling).To(Not(BeNil()))
				Expect(siblingSibling.Bookmark).To(BeNil())
			})
		})

		When("FPDFBookmark_Find is called", func() {
			It("it returns the correct bookmark when there is a match", func() {
				bookmark, err := PdfiumInstance.FPDFBookmark_Find(&requests.FPDFBookmark_Find{
					Document: doc,
					Title:    "A Good Beginning",
				})
				Expect(err).To(BeNil())
				Expect(bookmark).To(Not(BeNil()))
				Expect(bookmark.Bookmark).To(Not(BeNil()))
			})

			It("it returns the no bookmark when there is no match", func() {
				bookmark, err := PdfiumInstance.FPDFBookmark_Find(&requests.FPDFBookmark_Find{
					Document: doc,
					Title:    "No Good Beginning",
				})
				Expect(err).To(BeNil())
				Expect(bookmark).To(Equal(&responses.FPDFBookmark_Find{}))
			})
		})

		When("FPDFBookmark_GetDest is called", func() {
			It("it returns an error when no document is given", func() {
				action, err := PdfiumInstance.FPDFBookmark_GetDest(&requests.FPDFBookmark_GetDest{})

				Expect(err).To(MatchError("document not given"))
				Expect(action).To(BeNil())
			})
			It("it returns an error when no bookmark is given", func() {
				action, err := PdfiumInstance.FPDFBookmark_GetDest(&requests.FPDFBookmark_GetDest{
					Document: doc,
				})

				Expect(err).To(MatchError("bookmark not given"))
				Expect(action).To(BeNil())
			})
			It("it returns no dest because the bookmark has none", func() {
				bookmark, err := PdfiumInstance.FPDFBookmark_Find(&requests.FPDFBookmark_Find{
					Document: doc,
					Title:    "A Good Beginning",
				})
				Expect(err).To(BeNil())
				Expect(bookmark).To(Not(BeNil()))
				Expect(bookmark.Bookmark).To(Not(BeNil()))

				dest, err := PdfiumInstance.FPDFBookmark_GetDest(&requests.FPDFBookmark_GetDest{
					Document: doc,
					Bookmark: *bookmark.Bookmark,
				})

				Expect(err).To(BeNil())
				Expect(dest).To(Equal(&responses.FPDFBookmark_GetDest{}))
			})
		})

		When("FPDFBookmark_GetAction is called", func() {
			It("it returns an error when no bookmark is given", func() {
				action, err := PdfiumInstance.FPDFBookmark_GetAction(&requests.FPDFBookmark_GetAction{})

				Expect(err).To(MatchError("bookmark not given"))
				Expect(action).To(BeNil())
			})
			It("it returns no action because the bookmark has none", func() {
				bookmark, err := PdfiumInstance.FPDFBookmark_Find(&requests.FPDFBookmark_Find{
					Document: doc,
					Title:    "A Good Beginning",
				})
				Expect(err).To(BeNil())
				Expect(bookmark).To(Not(BeNil()))
				Expect(bookmark.Bookmark).To(Not(BeNil()))

				action, err := PdfiumInstance.FPDFBookmark_GetAction(&requests.FPDFBookmark_GetAction{
					Bookmark: *bookmark.Bookmark,
				})

				Expect(err).To(BeNil())
				Expect(action).To(Equal(&responses.FPDFBookmark_GetAction{}))
			})
		})
	})

	Context("a PDF file with no page labels", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/about_blank.pdf")
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

		When("FPDF_GetPageLabel is called", func() {
			It("returns the correct page label", func() {
				pageLabel, err := PdfiumInstance.FPDF_GetPageLabel(&requests.FPDF_GetPageLabel{
					Document: doc,
					Page:     0,
				})
				Expect(err).To(MatchError("Could not get label"))
				Expect(pageLabel).To(BeNil())
			})
		})
	})

	Context("a PDF file with page labels", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/page_labels.pdf")
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

		When("FPDF_GetPageLabel is called", func() {
			It("returns the correct page label", func() {
				pageLabel, err := PdfiumInstance.FPDF_GetPageLabel(&requests.FPDF_GetPageLabel{
					Document: doc,
					Page:     -2,
				})
				Expect(err).To(MatchError("Could not get label"))
				Expect(pageLabel).To(BeNil())

				pageLabel, err = PdfiumInstance.FPDF_GetPageLabel(&requests.FPDF_GetPageLabel{
					Document: doc,
					Page:     -1,
				})
				Expect(err).To(MatchError("Could not get label"))
				Expect(pageLabel).To(BeNil())

				pageLabel, err = PdfiumInstance.FPDF_GetPageLabel(&requests.FPDF_GetPageLabel{
					Document: doc,
					Page:     0,
				})
				Expect(err).To(BeNil())
				Expect(pageLabel).To(Equal(&responses.FPDF_GetPageLabel{
					Page:  0,
					Label: "i",
				}))

				pageLabel, err = PdfiumInstance.FPDF_GetPageLabel(&requests.FPDF_GetPageLabel{
					Document: doc,
					Page:     1,
				})
				Expect(err).To(BeNil())
				Expect(pageLabel).To(Equal(&responses.FPDF_GetPageLabel{
					Page:  1,
					Label: "ii",
				}))

				pageLabel, err = PdfiumInstance.FPDF_GetPageLabel(&requests.FPDF_GetPageLabel{
					Document: doc,
					Page:     2,
				})
				Expect(err).To(BeNil())
				Expect(pageLabel).To(Equal(&responses.FPDF_GetPageLabel{
					Page:  2,
					Label: "1",
				}))

				pageLabel, err = PdfiumInstance.FPDF_GetPageLabel(&requests.FPDF_GetPageLabel{
					Document: doc,
					Page:     3,
				})
				Expect(err).To(BeNil())
				Expect(pageLabel).To(Equal(&responses.FPDF_GetPageLabel{
					Page:  3,
					Label: "2",
				}))

				pageLabel, err = PdfiumInstance.FPDF_GetPageLabel(&requests.FPDF_GetPageLabel{
					Document: doc,
					Page:     4,
				})
				Expect(err).To(BeNil())
				Expect(pageLabel).To(Equal(&responses.FPDF_GetPageLabel{
					Page:  4,
					Label: "zzA",
				}))

				pageLabel, err = PdfiumInstance.FPDF_GetPageLabel(&requests.FPDF_GetPageLabel{
					Document: doc,
					Page:     5,
				})
				Expect(err).To(BeNil())
				Expect(pageLabel).To(Equal(&responses.FPDF_GetPageLabel{
					Page:  5,
					Label: "zzB",
				}))

				pageLabel, err = PdfiumInstance.FPDF_GetPageLabel(&requests.FPDF_GetPageLabel{
					Document: doc,
					Page:     6,
				})
				Expect(err).To(BeNil())
				Expect(pageLabel).To(Equal(&responses.FPDF_GetPageLabel{
					Page:  6,
					Label: "",
				}))

				pageLabel, err = PdfiumInstance.FPDF_GetPageLabel(&requests.FPDF_GetPageLabel{
					Document: doc,
					Page:     7,
				})
				Expect(err).To(MatchError("Could not get label"))
				Expect(pageLabel).To(BeNil())

				pageLabel, err = PdfiumInstance.FPDF_GetPageLabel(&requests.FPDF_GetPageLabel{
					Document: doc,
					Page:     8,
				})
				Expect(err).To(MatchError("Could not get label"))
				Expect(pageLabel).To(BeNil())
			})
		})
	})

	Context("a PDF file with a non-text identifier", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/split_streams.pdf")
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

		When("FPDF_GetFileIdentifier is called with an invalid type", func() {
			It("should return an error", func() {
				identifier, err := PdfiumInstance.FPDF_GetFileIdentifier(&requests.FPDF_GetFileIdentifier{
					Document:   doc,
					FileIdType: -1,
				})
				Expect(err).To(MatchError("invalid file id type given"))
				Expect(identifier).To(BeNil())

				identifier, err = PdfiumInstance.FPDF_GetFileIdentifier(&requests.FPDF_GetFileIdentifier{
					Document:   doc,
					FileIdType: 2,
				})
				Expect(err).To(MatchError("invalid file id type given"))
				Expect(identifier).To(BeNil())
			})
		})

		When("FPDF_GetFileIdentifier is called with a valid type", func() {
			It("returns the correct identifier", func() {
				identifier, err := PdfiumInstance.FPDF_GetFileIdentifier(&requests.FPDF_GetFileIdentifier{
					Document:   doc,
					FileIdType: enums.FPDF_FILEIDTYPE_PERMANENT,
				})
				Expect(err).To(BeNil())
				Expect(identifier).To(Equal(&responses.FPDF_GetFileIdentifier{
					FileIdType: enums.FPDF_FILEIDTYPE_PERMANENT,
					Identifier: []byte{243, 65, 174, 101, 74, 119, 172, 213, 6, 90, 118, 69, 229, 150, 230, 230}, // Byte identifier
				}))

				identifier, err = PdfiumInstance.FPDF_GetFileIdentifier(&requests.FPDF_GetFileIdentifier{
					Document:   doc,
					FileIdType: enums.FPDF_FILEIDTYPE_CHANGING,
				})
				Expect(err).To(BeNil())
				Expect(identifier).To(Equal(&responses.FPDF_GetFileIdentifier{
					FileIdType: enums.FPDF_FILEIDTYPE_CHANGING,
					Identifier: []byte{188, 55, 41, 138, 63, 135, 244, 121, 34, 155, 206, 153, 124, 167, 145, 247}, // Byte identifier
				}))
			})
		})
	})

	Context("a PDF file with a text identifier", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/non_hex_file_id.pdf")
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

		When("FPDF_GetFileIdentifier is called with a valid type", func() {
			It("returns the correct identifier", func() {
				identifier, err := PdfiumInstance.FPDF_GetFileIdentifier(&requests.FPDF_GetFileIdentifier{
					Document:   doc,
					FileIdType: enums.FPDF_FILEIDTYPE_PERMANENT,
				})
				Expect(err).To(BeNil())
				Expect(identifier).To(Equal(&responses.FPDF_GetFileIdentifier{
					FileIdType: enums.FPDF_FILEIDTYPE_PERMANENT,
					Identifier: []byte("permanent non-hex"), // Text identifier
				}))

				identifier, err = PdfiumInstance.FPDF_GetFileIdentifier(&requests.FPDF_GetFileIdentifier{
					Document:   doc,
					FileIdType: enums.FPDF_FILEIDTYPE_CHANGING,
				})
				Expect(err).To(BeNil())
				Expect(identifier).To(Equal(&responses.FPDF_GetFileIdentifier{
					FileIdType: enums.FPDF_FILEIDTYPE_CHANGING,
					Identifier: []byte("changing non-hex"), // Text identifier
				}))
			})
		})
	})

	Context("a PDF file without identifier", func() {
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

		When("FPDF_GetFileIdentifier is called with a valid type", func() {
			It("returns a nil identifier", func() {
				identifier, err := PdfiumInstance.FPDF_GetFileIdentifier(&requests.FPDF_GetFileIdentifier{
					Document:   doc,
					FileIdType: enums.FPDF_FILEIDTYPE_PERMANENT,
				})
				Expect(err).To(BeNil())
				Expect(identifier).To(Equal(&responses.FPDF_GetFileIdentifier{
					FileIdType: enums.FPDF_FILEIDTYPE_PERMANENT,
					Identifier: nil,
				}))

				identifier, err = PdfiumInstance.FPDF_GetFileIdentifier(&requests.FPDF_GetFileIdentifier{
					Document:   doc,
					FileIdType: enums.FPDF_FILEIDTYPE_CHANGING,
				})
				Expect(err).To(BeNil())
				Expect(identifier).To(Equal(&responses.FPDF_GetFileIdentifier{
					FileIdType: enums.FPDF_FILEIDTYPE_CHANGING,
					Identifier: nil,
				}))
			})
		})
	})

	Context("a PDF file with a link", func() {
		var doc references.FPDF_DOCUMENT
		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/launch_action.pdf")
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

		When("looking for the link", func() {
			It("FPDFLink_GetLinkAtPoint returns the link", func() {
				pageLink, err := PdfiumInstance.FPDFLink_GetLinkAtPoint(&requests.FPDFLink_GetLinkAtPoint{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					X: 100,
					Y: 100,
				})
				Expect(err).To(BeNil())
				Expect(pageLink).To(Not(BeNil()))
				Expect(pageLink.Link).To(Not(BeNil()))
			})

			It("FPDFLink_GetLinkZOrderAtPoint returns the z-order of the link", func() {
				pageLinkZOrder, err := PdfiumInstance.FPDFLink_GetLinkZOrderAtPoint(&requests.FPDFLink_GetLinkZOrderAtPoint{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					X: 100,
					Y: 100,
				})
				Expect(err).To(BeNil())
				Expect(pageLinkZOrder).To(Equal(&responses.FPDFLink_GetLinkZOrderAtPoint{
					ZOrder: 0,
				}))
			})
		})
	})

	Context("a PDF file with a launch action", func() {
		var doc references.FPDF_DOCUMENT
		var link references.FPDF_LINK
		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/launch_action.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document

			pageLink, err := PdfiumInstance.FPDFLink_GetLinkAtPoint(&requests.FPDFLink_GetLinkAtPoint{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
				X: 100,
				Y: 100,
			})
			Expect(err).To(BeNil())
			Expect(pageLink).To(Not(BeNil()))
			Expect(pageLink.Link).To(Not(BeNil()))
			link = *pageLink.Link
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		When("FPDFLink_GetAction is called", func() {
			It("returns an action", func() {
				action, err := PdfiumInstance.FPDFLink_GetAction(&requests.FPDFLink_GetAction{
					Link: link,
				})
				Expect(err).To(BeNil())
				Expect(action).To(Not(BeNil()))
				Expect(action.Action).To(Not(BeNil()))
			})
		})

		When("FPDFLink_GetAnnot is called", func() {
			It("returns an annotation", func() {
				annotation, err := PdfiumInstance.FPDFLink_GetAnnot(&requests.FPDFLink_GetAnnot{
					Link: link,
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(annotation).To(Not(BeNil()))
				Expect(annotation.Annotation).To(Not(BeNil()))
			})
		})

		When("FPDFLink_GetAnnotRect is called", func() {
			It("returns an annotation rect", func() {
				annotationRect, err := PdfiumInstance.FPDFLink_GetAnnotRect(&requests.FPDFLink_GetAnnotRect{
					Link: link,
				})
				Expect(err).To(BeNil())
				Expect(annotationRect).To(Not(BeNil()))
				Expect(annotationRect).To(Equal(&responses.FPDFLink_GetAnnotRect{
					Rect: &structs.FPDF_FS_RECTF{Left: 1, Top: 199, Right: 199, Bottom: 1},
				}))
			})
		})

		When("FPDFLink_CountQuadPoints is called", func() {
			It("returns the quad points count of 0", func() {
				quadPointsCount, err := PdfiumInstance.FPDFLink_CountQuadPoints(&requests.FPDFLink_CountQuadPoints{
					Link: link,
				})
				Expect(err).To(BeNil())
				Expect(quadPointsCount).To(Not(BeNil()))
				Expect(quadPointsCount).To(Equal(&responses.FPDFLink_CountQuadPoints{
					Count: 0,
				}))
			})
		})

		When("FPDFLink_GetQuadPoints is called", func() {
			It("returns no quad points", func() {
				quadPoints, err := PdfiumInstance.FPDFLink_GetQuadPoints(&requests.FPDFLink_GetQuadPoints{
					Link:      link,
					QuadIndex: 0,
				})
				Expect(err).To(BeNil())
				Expect(quadPoints).To(Not(BeNil()))
				Expect(quadPoints).To(Equal(&responses.FPDFLink_GetQuadPoints{}))
			})
		})

		Context("A launch action is loaded", func() {
			var action references.FPDF_ACTION

			BeforeEach(func() {
				actionResp, err := PdfiumInstance.FPDFLink_GetAction(&requests.FPDFLink_GetAction{
					Link: link,
				})
				Expect(err).To(BeNil())
				Expect(actionResp).To(Not(BeNil()))
				Expect(actionResp.Action).To(Not(BeNil()))
				action = *actionResp.Action
			})

			When("FPDFAction_GetType is called", func() {
				It("returns an type", func() {
					actionType, err := PdfiumInstance.FPDFAction_GetType(&requests.FPDFAction_GetType{
						Action: action,
					})
					Expect(err).To(BeNil())
					Expect(actionType).To(Not(BeNil()))
					Expect(actionType).To(Equal(&responses.FPDFAction_GetType{
						Type: enums.FPDF_ACTION_ACTION_LAUNCH,
					}))
				})
			})

			When("FPDFAction_GetFilePath is called", func() {
				It("returns a file path", func() {
					actionType, err := PdfiumInstance.FPDFAction_GetFilePath(&requests.FPDFAction_GetFilePath{
						Action: action,
					})
					Expect(err).To(BeNil())
					Expect(actionType).To(Not(BeNil()))
					expectedPath := "test.pdf"
					Expect(actionType).To(Equal(&responses.FPDFAction_GetFilePath{
						FilePath: &expectedPath,
					}))
				})
			})

			When("FPDFAction_GetDest is called", func() {
				It("returns no dest", func() {
					actionDest, err := PdfiumInstance.FPDFAction_GetDest(&requests.FPDFAction_GetDest{
						Document: doc,
						Action:   action,
					})
					Expect(err).To(BeNil())
					Expect(actionDest).To(Not(BeNil()))
					Expect(actionDest).To(Equal(&responses.FPDFAction_GetDest{}))
				})
			})

			When("FPDFAction_GetURIPath is called", func() {
				It("returns no uri path", func() {
					uriPath, err := PdfiumInstance.FPDFAction_GetURIPath(&requests.FPDFAction_GetURIPath{
						Document: doc,
						Action:   action,
					})
					Expect(err).To(BeNil())
					Expect(uriPath).To(Not(BeNil()))
					Expect(uriPath).To(Equal(&responses.FPDFAction_GetURIPath{}))
				})
			})
		})
	})

	Context("a PDF file with a uri action", func() {
		var doc references.FPDF_DOCUMENT
		var link references.FPDF_LINK
		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/uri_action.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document

			pageLink, err := PdfiumInstance.FPDFLink_GetLinkAtPoint(&requests.FPDFLink_GetLinkAtPoint{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
				X: 100,
				Y: 100,
			})
			Expect(err).To(BeNil())
			Expect(pageLink).To(Not(BeNil()))
			Expect(pageLink.Link).To(Not(BeNil()))
			link = *pageLink.Link
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		When("FPDFLink_GetAction is called", func() {
			It("returns an action", func() {
				action, err := PdfiumInstance.FPDFLink_GetAction(&requests.FPDFLink_GetAction{
					Link: link,
				})
				Expect(err).To(BeNil())
				Expect(action).To(Not(BeNil()))
				Expect(action.Action).To(Not(BeNil()))
			})
		})

		Context("A uri action is loaded", func() {
			var action references.FPDF_ACTION

			BeforeEach(func() {
				actionResp, err := PdfiumInstance.FPDFLink_GetAction(&requests.FPDFLink_GetAction{
					Link: link,
				})
				Expect(err).To(BeNil())
				Expect(actionResp).To(Not(BeNil()))
				Expect(actionResp.Action).To(Not(BeNil()))
				action = *actionResp.Action
			})

			When("FPDFAction_GetType is called", func() {
				It("returns an type", func() {
					actionType, err := PdfiumInstance.FPDFAction_GetType(&requests.FPDFAction_GetType{
						Action: action,
					})
					Expect(err).To(BeNil())
					Expect(actionType).To(Not(BeNil()))
					Expect(actionType).To(Equal(&responses.FPDFAction_GetType{
						Type: enums.FPDF_ACTION_ACTION_URI,
					}))
				})
			})

			When("FPDFAction_GetURIPath is called", func() {
				It("returns a uri path", func() {
					uriPath, err := PdfiumInstance.FPDFAction_GetURIPath(&requests.FPDFAction_GetURIPath{
						Document: doc,
						Action:   action,
					})
					Expect(err).To(BeNil())
					Expect(uriPath).To(Not(BeNil()))
					expectedUriPath := "https://example.com/page.html"
					Expect(uriPath).To(Equal(&responses.FPDFAction_GetURIPath{
						URIPath: &expectedUriPath,
					}))
				})
			})

			When("FPDFAction_GetFilePath is called", func() {
				It("returns no file path", func() {
					actionType, err := PdfiumInstance.FPDFAction_GetFilePath(&requests.FPDFAction_GetFilePath{
						Action: action,
					})
					Expect(err).To(BeNil())
					Expect(actionType).To(Not(BeNil()))
					Expect(actionType).To(Equal(&responses.FPDFAction_GetFilePath{}))
				})
			})

			When("FPDFAction_GetDest is called", func() {
				It("returns no dest", func() {
					actionDest, err := PdfiumInstance.FPDFAction_GetDest(&requests.FPDFAction_GetDest{
						Document: doc,
						Action:   action,
					})
					Expect(err).To(BeNil())
					Expect(actionDest).To(Not(BeNil()))
					Expect(actionDest).To(Equal(&responses.FPDFAction_GetDest{}))
				})
			})
		})
	})

	Context("a PDF file with a goto action", func() {
		var doc references.FPDF_DOCUMENT
		var link references.FPDF_LINK
		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/goto_action.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document

			pageLink, err := PdfiumInstance.FPDFLink_GetLinkAtPoint(&requests.FPDFLink_GetLinkAtPoint{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
				X: 100,
				Y: 100,
			})
			Expect(err).To(BeNil())
			Expect(pageLink).To(Not(BeNil()))
			Expect(pageLink.Link).To(Not(BeNil()))
			link = *pageLink.Link
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		When("FPDFLink_GetAction is called", func() {
			It("returns an action", func() {
				action, err := PdfiumInstance.FPDFLink_GetAction(&requests.FPDFLink_GetAction{
					Link: link,
				})
				Expect(err).To(BeNil())
				Expect(action).To(Not(BeNil()))
				Expect(action.Action).To(Not(BeNil()))
			})
		})

		Context("A goto action is loaded", func() {
			var action references.FPDF_ACTION

			BeforeEach(func() {
				actionResp, err := PdfiumInstance.FPDFLink_GetAction(&requests.FPDFLink_GetAction{
					Link: link,
				})
				Expect(err).To(BeNil())
				Expect(actionResp).To(Not(BeNil()))
				Expect(actionResp.Action).To(Not(BeNil()))
				action = *actionResp.Action
			})

			When("FPDFAction_GetType is called", func() {
				It("returns an type", func() {
					actionType, err := PdfiumInstance.FPDFAction_GetType(&requests.FPDFAction_GetType{
						Action: action,
					})
					Expect(err).To(BeNil())
					Expect(actionType).To(Not(BeNil()))
					Expect(actionType).To(Equal(&responses.FPDFAction_GetType{
						Type: enums.FPDF_ACTION_ACTION_GOTO,
					}))
				})
			})

			When("FPDFAction_GetURIPath is called", func() {
				It("returns no uri path", func() {
					uriPath, err := PdfiumInstance.FPDFAction_GetURIPath(&requests.FPDFAction_GetURIPath{
						Document: doc,
						Action:   action,
					})
					Expect(err).To(BeNil())
					Expect(uriPath).To(Not(BeNil()))
					Expect(uriPath).To(Equal(&responses.FPDFAction_GetURIPath{}))
				})
			})

			When("FPDFAction_GetFilePath is called", func() {
				It("returns no file path", func() {
					actionType, err := PdfiumInstance.FPDFAction_GetFilePath(&requests.FPDFAction_GetFilePath{
						Action: action,
					})
					Expect(err).To(BeNil())
					Expect(actionType).To(Not(BeNil()))
					Expect(actionType).To(Equal(&responses.FPDFAction_GetFilePath{}))
				})
			})

			When("FPDFAction_GetDest is called", func() {
				It("returns a dest", func() {
					actionDest, err := PdfiumInstance.FPDFAction_GetDest(&requests.FPDFAction_GetDest{
						Document: doc,
						Action:   action,
					})
					Expect(err).To(BeNil())
					Expect(actionDest).To(Not(BeNil()))
					Expect(actionDest.Dest).To(Not(BeNil()))
				})
			})

			Context("A dest action is loaded", func() {
				var dest references.FPDF_DEST

				BeforeEach(func() {
					actionDest, err := PdfiumInstance.FPDFAction_GetDest(&requests.FPDFAction_GetDest{
						Document: doc,
						Action:   action,
					})
					Expect(err).To(BeNil())
					Expect(actionDest).To(Not(BeNil()))
					Expect(actionDest.Dest).To(Not(BeNil()))
					dest = *actionDest.Dest
				})

				When("FPDFDest_GetDestPageIndex is called", func() {
					It("returns the correct page index", func() {
						destPageIndex, err := PdfiumInstance.FPDFDest_GetDestPageIndex(&requests.FPDFDest_GetDestPageIndex{
							Document: doc,
							Dest:     dest,
						})
						Expect(err).To(BeNil())
						Expect(destPageIndex).To(Not(BeNil()))
						Expect(destPageIndex).To(Equal(&responses.FPDFDest_GetDestPageIndex{
							Index: 1,
						}))
					})
				})
			})
		})
	})

	Context("a PDF file with an embedded goto action", func() {
		var doc references.FPDF_DOCUMENT
		var link references.FPDF_LINK
		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/gotoe_action.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document

			pageLink, err := PdfiumInstance.FPDFLink_GetLinkAtPoint(&requests.FPDFLink_GetLinkAtPoint{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
				X: 100,
				Y: 100,
			})
			Expect(err).To(BeNil())
			Expect(pageLink).To(Not(BeNil()))
			Expect(pageLink.Link).To(Not(BeNil()))
			link = *pageLink.Link
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		When("FPDFLink_GetAction is called", func() {
			It("returns an action", func() {
				action, err := PdfiumInstance.FPDFLink_GetAction(&requests.FPDFLink_GetAction{
					Link: link,
				})
				Expect(err).To(BeNil())
				Expect(action).To(Not(BeNil()))
				Expect(action.Action).To(Not(BeNil()))
			})
		})

		Context("An embedded goto action is loaded", func() {
			var action references.FPDF_ACTION

			BeforeEach(func() {
				actionResp, err := PdfiumInstance.FPDFLink_GetAction(&requests.FPDFLink_GetAction{
					Link: link,
				})
				Expect(err).To(BeNil())
				Expect(actionResp).To(Not(BeNil()))
				Expect(actionResp.Action).To(Not(BeNil()))
				action = *actionResp.Action
			})

			When("FPDFAction_GetType is called", func() {
				It("returns an type", func() {
					actionType, err := PdfiumInstance.FPDFAction_GetType(&requests.FPDFAction_GetType{
						Action: action,
					})
					Expect(err).To(BeNil())
					Expect(actionType).To(Not(BeNil()))
					Expect(actionType).To(Equal(&responses.FPDFAction_GetType{
						Type: enums.FPDF_ACTION_ACTION_EMBEDDEDGOTO,
					}))
				})
			})

			When("FPDFAction_GetURIPath is called", func() {
				It("returns no uri path", func() {
					uriPath, err := PdfiumInstance.FPDFAction_GetURIPath(&requests.FPDFAction_GetURIPath{
						Document: doc,
						Action:   action,
					})
					Expect(err).To(BeNil())
					Expect(uriPath).To(Not(BeNil()))
					Expect(uriPath).To(Equal(&responses.FPDFAction_GetURIPath{}))
				})
			})

			When("FPDFAction_GetFilePath is called", func() {
				It("returns the file path", func() {
					actionType, err := PdfiumInstance.FPDFAction_GetFilePath(&requests.FPDFAction_GetFilePath{
						Action: action,
					})
					Expect(err).To(BeNil())
					Expect(actionType).To(Not(BeNil()))
					expectedFilePath := "ExampleFile.pdf"
					Expect(actionType).To(Equal(&responses.FPDFAction_GetFilePath{
						FilePath: &expectedFilePath,
					}))
				})
			})

			When("FPDFAction_GetDest is called", func() {
				It("returns a dest", func() {
					actionDest, err := PdfiumInstance.FPDFAction_GetDest(&requests.FPDFAction_GetDest{
						Document: doc,
						Action:   action,
					})
					Expect(err).To(BeNil())
					Expect(actionDest).To(Not(BeNil()))
					Expect(actionDest.Dest).To(Not(BeNil()))
				})
			})

			Context("A dest action is loaded", func() {
				var dest references.FPDF_DEST

				BeforeEach(func() {
					actionDest, err := PdfiumInstance.FPDFAction_GetDest(&requests.FPDFAction_GetDest{
						Document: doc,
						Action:   action,
					})
					Expect(err).To(BeNil())
					Expect(actionDest).To(Not(BeNil()))
					Expect(actionDest.Dest).To(Not(BeNil()))
					dest = *actionDest.Dest
				})

				When("FPDFDest_GetDestPageIndex is called", func() {
					It("returns the correct page index", func() {
						destPageIndex, err := PdfiumInstance.FPDFDest_GetDestPageIndex(&requests.FPDFDest_GetDestPageIndex{
							Document: doc,
							Dest:     dest,
						})
						Expect(err).To(BeNil())
						Expect(destPageIndex).To(Not(BeNil()))
						Expect(destPageIndex).To(Equal(&responses.FPDFDest_GetDestPageIndex{
							Index: 1,
						}))
					})
				})
			})
		})
	})

	Context("a PDF file with a link dest", func() {
		var doc references.FPDF_DOCUMENT
		var link references.FPDF_LINK
		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/bug_821454.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document

			pageLink, err := PdfiumInstance.FPDFLink_GetLinkAtPoint(&requests.FPDFLink_GetLinkAtPoint{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
				X: 150,
				Y: 360,
			})
			Expect(err).To(BeNil())
			Expect(pageLink).To(Not(BeNil()))
			Expect(pageLink.Link).To(Not(BeNil()))
			link = *pageLink.Link
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		When("FPDFLink_GetAction is called", func() {
			It("returns no action", func() {
				action, err := PdfiumInstance.FPDFLink_GetAction(&requests.FPDFLink_GetAction{
					Link: link,
				})
				Expect(err).To(BeNil())
				Expect(action).To(Not(BeNil()))
				Expect(action.Action).To(BeNil())
			})
		})

		When("FPDFLink_GetDest is called", func() {
			It("returns a dest", func() {
				dest, err := PdfiumInstance.FPDFLink_GetDest(&requests.FPDFLink_GetDest{
					Document: doc,
					Link:     link,
				})
				Expect(err).To(BeNil())
				Expect(dest).To(Not(BeNil()))
				Expect(dest.Dest).To(Not(BeNil()))
			})
		})

		When("FPDFLink_Enumerate is called", func() {
			It("returns a link", func() {
				FPDFLink_Enumerate, err := PdfiumInstance.FPDFLink_Enumerate(&requests.FPDFLink_Enumerate{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDFLink_Enumerate).To(Not(BeNil()))
				Expect(FPDFLink_Enumerate.Link).To(Not(BeNil()))
				Expect(FPDFLink_Enumerate.NextStartPos).To(PointTo(Equal(1)))

				FPDFLink_Enumerate, err = PdfiumInstance.FPDFLink_Enumerate(&requests.FPDFLink_Enumerate{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					StartPos: *FPDFLink_Enumerate.NextStartPos,
				})
				Expect(err).To(BeNil())
				Expect(FPDFLink_Enumerate).To(Not(BeNil()))
				Expect(FPDFLink_Enumerate.Link).To(Not(BeNil()))
				Expect(FPDFLink_Enumerate.NextStartPos).To(PointTo(Equal(2)))

				FPDFLink_Enumerate, err = PdfiumInstance.FPDFLink_Enumerate(&requests.FPDFLink_Enumerate{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					StartPos: *FPDFLink_Enumerate.NextStartPos,
				})
				Expect(err).To(BeNil())
				Expect(FPDFLink_Enumerate).To(Not(BeNil()))
				Expect(FPDFLink_Enumerate.Link).To(BeNil())
				Expect(FPDFLink_Enumerate.NextStartPos).To(BeNil())
			})
		})

		When("a dest is loaded", func() {
			var dest references.FPDF_DEST
			BeforeEach(func() {
				destResp, err := PdfiumInstance.FPDFLink_GetDest(&requests.FPDFLink_GetDest{
					Document: doc,
					Link:     link,
				})
				Expect(err).To(BeNil())
				Expect(destResp).To(Not(BeNil()))
				Expect(destResp.Dest).To(Not(BeNil()))
				dest = *destResp.Dest
			})

			When("FPDFDest_GetDestPageIndex is called", func() {
				It("returns the right page index", func() {
					FPDFDest_GetDestPageIndex, err := PdfiumInstance.FPDFDest_GetDestPageIndex(&requests.FPDFDest_GetDestPageIndex{
						Document: doc,
						Dest:     dest,
					})
					Expect(err).To(BeNil())
					Expect(FPDFDest_GetDestPageIndex).To(Equal(&responses.FPDFDest_GetDestPageIndex{
						Index: 0,
					}))
				})
			})

			When("FPDFDest_GetLocationInPage is called", func() {
				It("returns the right location", func() {
					FPDFDest_GetLocationInPage, err := PdfiumInstance.FPDFDest_GetLocationInPage(&requests.FPDFDest_GetLocationInPage{
						Dest: dest,
					})
					Expect(err).To(BeNil())
					expectedX := float32(100)
					expectedY := float32(200)
					Expect(FPDFDest_GetLocationInPage).To(Equal(&responses.FPDFDest_GetLocationInPage{
						X:    &expectedX,
						Y:    &expectedY,
						Zoom: nil,
					}))
				})
			})

			When("FPDFDest_GetView is called", func() {
				It("returns the right view", func() {
					FPDFDest_GetView, err := PdfiumInstance.FPDFDest_GetView(&requests.FPDFDest_GetView{
						Dest: dest,
					})
					Expect(err).To(BeNil())
					Expect(FPDFDest_GetView).To(Equal(&responses.FPDFDest_GetView{
						DestView: 1,
						Params:   []float32{100, 200, 0},
					}))
				})
			})
		})
	})

	Context("a PDF file with an aaction", func() {
		var doc references.FPDF_DOCUMENT
		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/get_page_aaction.pdf")
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

		When("FPDF_GetPageAAction is called", func() {
			It("returns null for FPDFPAGE_AACTION_CLOSE", func() {
				FPDF_GetPageAAction, err := PdfiumInstance.FPDF_GetPageAAction(&requests.FPDF_GetPageAAction{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					AAType: enums.FPDF_PAGE_AACTION_CLOSE,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetPageAAction).To(Equal(&responses.FPDF_GetPageAAction{}))
			})

			It("returns null for -1", func() {
				FPDF_GetPageAAction, err := PdfiumInstance.FPDF_GetPageAAction(&requests.FPDF_GetPageAAction{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					AAType: -1,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetPageAAction).To(Equal(&responses.FPDF_GetPageAAction{}))
			})

			It("returns null for 99", func() {
				FPDF_GetPageAAction, err := PdfiumInstance.FPDF_GetPageAAction(&requests.FPDF_GetPageAAction{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					AAType: 99,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetPageAAction).To(Equal(&responses.FPDF_GetPageAAction{}))
			})

			It("returns an action for FPDFPAGE_AACTION_OPEN", func() {
				FPDF_GetPageAAction, err := PdfiumInstance.FPDF_GetPageAAction(&requests.FPDF_GetPageAAction{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					AAType: enums.FPDF_PAGE_AACTION_OPEN,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetPageAAction).To(Not(BeNil()))
				Expect(FPDF_GetPageAAction.Action).To(Not(BeNil()))
			})

			Context("An embedded goto action is loaded", func() {
				var action references.FPDF_ACTION

				BeforeEach(func() {
					FPDF_GetPageAAction, err := PdfiumInstance.FPDF_GetPageAAction(&requests.FPDF_GetPageAAction{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
						AAType: enums.FPDF_PAGE_AACTION_OPEN,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_GetPageAAction).To(Not(BeNil()))
					Expect(FPDF_GetPageAAction.Action).To(Not(BeNil()))
					action = *FPDF_GetPageAAction.Action
				})

				When("FPDFAction_GetType is called", func() {
					It("returns an type", func() {
						actionType, err := PdfiumInstance.FPDFAction_GetType(&requests.FPDFAction_GetType{
							Action: action,
						})
						Expect(err).To(BeNil())
						Expect(actionType).To(Not(BeNil()))
						Expect(actionType).To(Equal(&responses.FPDFAction_GetType{
							Type: enums.FPDF_ACTION_ACTION_EMBEDDEDGOTO,
						}))
					})
				})

				When("FPDFAction_GetURIPath is called", func() {
					It("returns no uri path", func() {
						uriPath, err := PdfiumInstance.FPDFAction_GetURIPath(&requests.FPDFAction_GetURIPath{
							Document: doc,
							Action:   action,
						})
						Expect(err).To(BeNil())
						Expect(uriPath).To(Not(BeNil()))
						Expect(uriPath).To(Equal(&responses.FPDFAction_GetURIPath{}))
					})
				})

				When("FPDFAction_GetFilePath is called", func() {
					It("returns the file path", func() {
						actionType, err := PdfiumInstance.FPDFAction_GetFilePath(&requests.FPDFAction_GetFilePath{
							Action: action,
						})
						Expect(err).To(BeNil())
						Expect(actionType).To(Not(BeNil()))
						expectedFilePath := "\\\\127.0.0.1\\c$\\Program Files\\test.exe"
						Expect(actionType).To(Equal(&responses.FPDFAction_GetFilePath{
							FilePath: &expectedFilePath,
						}))
					})
				})

				When("FPDFAction_GetDest is called", func() {
					It("returns a dest", func() {
						actionDest, err := PdfiumInstance.FPDFAction_GetDest(&requests.FPDFAction_GetDest{
							Document: doc,
							Action:   action,
						})
						Expect(err).To(BeNil())
						Expect(actionDest).To(Not(BeNil()))
						Expect(actionDest.Dest).To(Not(BeNil()))
					})
				})

				Context("A dest action is loaded", func() {
					var dest references.FPDF_DEST

					BeforeEach(func() {
						actionDest, err := PdfiumInstance.FPDFAction_GetDest(&requests.FPDFAction_GetDest{
							Document: doc,
							Action:   action,
						})
						Expect(err).To(BeNil())
						Expect(actionDest).To(Not(BeNil()))
						Expect(actionDest.Dest).To(Not(BeNil()))
						dest = *actionDest.Dest
					})

					When("FPDFDest_GetDestPageIndex is called", func() {
						It("returns the correct page index", func() {
							destPageIndex, err := PdfiumInstance.FPDFDest_GetDestPageIndex(&requests.FPDFDest_GetDestPageIndex{
								Document: doc,
								Dest:     dest,
							})
							Expect(err).To(BeNil())
							Expect(destPageIndex).To(Not(BeNil()))
							Expect(destPageIndex).To(Equal(&responses.FPDFDest_GetDestPageIndex{
								Index: 1,
							}))
						})
					})
				})
			})
		})
	})

	Context("a PDF file with quad point links", func() {
		var doc references.FPDF_DOCUMENT
		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/annots.pdf")
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

		When("the quad point link is requested", func() {
			It("returns the correct link, quad points count and quad points", func() {
				FPDFLink_Enumerate, err := PdfiumInstance.FPDFLink_Enumerate(&requests.FPDFLink_Enumerate{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					StartPos: 3,
				})
				Expect(err).To(BeNil())
				Expect(FPDFLink_Enumerate).To(Not(BeNil()))
				Expect(FPDFLink_Enumerate.Link).To(Not(BeNil()))
				Expect(FPDFLink_Enumerate.NextStartPos).To(PointTo(Equal(4)))

				FPDFLink_CountQuadPoints, err := PdfiumInstance.FPDFLink_CountQuadPoints(&requests.FPDFLink_CountQuadPoints{
					Link: *FPDFLink_Enumerate.Link,
				})

				Expect(err).To(BeNil())
				Expect(FPDFLink_CountQuadPoints).To(Equal(&responses.FPDFLink_CountQuadPoints{
					Count: 1,
				}))

				FPDFLink_GetQuadPoints, err := PdfiumInstance.FPDFLink_GetQuadPoints(&requests.FPDFLink_GetQuadPoints{
					Link:      *FPDFLink_Enumerate.Link,
					QuadIndex: 0,
				})

				Expect(err).To(BeNil())
				Expect(FPDFLink_GetQuadPoints).To(Equal(&responses.FPDFLink_GetQuadPoints{
					Points: &structs.FPDF_FS_QUADPOINTSF{X1: 83, Y1: 453, X2: 178, Y2: 453, X3: 83, Y3: 440, X4: 178, Y4: 440},
				}))
			})
		})
	})
})
