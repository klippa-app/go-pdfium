package shared_tests

import (
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"io/ioutil"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/references"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func RunfpdfDocTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("fpdf_doc", func() {
		Context("a normal PDF file", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				pdfData, err := ioutil.ReadFile(testsPath + "/testdata/test.pdf")
				Expect(err).To(BeNil())
				if err != nil {
					return
				}

				newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData)
				if err != nil {
					return
				}

				doc = *newDoc
			})

			AfterEach(func() {
				err := pdfiumContainer.FPDF_CloseDocument(doc)
				Expect(err).To(BeNil())
			})

			When("is opened", func() {
				It("returns the correct metadata text", func() {
					metadata, err := pdfiumContainer.FPDF_GetMetaText(&requests.FPDF_GetMetaText{
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
					metadata, err := pdfiumContainer.GetMetaData(&requests.GetMetaData{
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
					metadata, err := pdfiumContainer.GetMetaData(&requests.GetMetaData{
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
		})

		Context("a PDF file with no bookmarks", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				pdfData, err := ioutil.ReadFile(testsPath + "/testdata/test.pdf")
				Expect(err).To(BeNil())
				if err != nil {
					return
				}

				newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData)
				if err != nil {
					return
				}

				doc = *newDoc
			})

			AfterEach(func() {
				err := pdfiumContainer.FPDF_CloseDocument(doc)
				Expect(err).To(BeNil())
			})

			When("FPDFBookmark_GetFirstChild is called", func() {
				It("returns no bookmark", func() {
					fistChild, err := pdfiumContainer.FPDFBookmark_GetFirstChild(&requests.FPDFBookmark_GetFirstChild{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(fistChild).To(Equal(&responses.FPDFBookmark_GetFirstChild{}))
				})
			})

			When("FPDFBookmark_GetNextSibling is called without a bookmark", func() {
				It("returns an error", func() {
					nextSibling, err := pdfiumContainer.FPDFBookmark_GetNextSibling(&requests.FPDFBookmark_GetNextSibling{
						Document: doc,
					})
					Expect(err).To(MatchError("bookmark not given"))
					Expect(nextSibling).To(BeNil())
				})
			})

			When("FPDFBookmark_Find is called without a title", func() {
				It("returns an error", func() {
					bookmark, err := pdfiumContainer.FPDFBookmark_Find(&requests.FPDFBookmark_Find{
						Document: doc,
					})
					Expect(err).To(MatchError("no title given"))
					Expect(bookmark).To(BeNil())
				})
			})

			When("FPDFBookmark_Find is called with a title", func() {
				It("returns no bookmark", func() {
					bookmark, err := pdfiumContainer.FPDFBookmark_Find(&requests.FPDFBookmark_Find{
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
				pdfData, err := ioutil.ReadFile(testsPath + "/testdata/bookmarks.pdf")
				Expect(err).To(BeNil())
				if err != nil {
					return
				}

				newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData)
				if err != nil {
					return
				}

				doc = *newDoc
			})

			AfterEach(func() {
				err := pdfiumContainer.FPDF_CloseDocument(doc)
				Expect(err).To(BeNil())
			})

			When("FPDFBookmark_GetFirstChild is called", func() {
				It("returns a bookmark with a matching tile, no children and 1 sibling", func() {
					topLevelBookmark, err := pdfiumContainer.FPDFBookmark_GetFirstChild(&requests.FPDFBookmark_GetFirstChild{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(topLevelBookmark).To(Not(BeNil()))

					if topLevelBookmark != nil {
						Expect(topLevelBookmark.Bookmark).To(Not(BeNil()))
						// Check title of bookmark.
						topLevelBookmarkTitle, err := pdfiumContainer.FPDFBookmark_GetTitle(&requests.FPDFBookmark_GetTitle{
							Bookmark: *topLevelBookmark.Bookmark,
						})
						Expect(err).To(BeNil())
						Expect(topLevelBookmarkTitle).To(Equal(&responses.FPDFBookmark_GetTitle{
							Title: "A Good Beginning",
						}))

						// Check that we have no children
						topLevelBookmarkSibling, err := pdfiumContainer.FPDFBookmark_GetFirstChild(&requests.FPDFBookmark_GetFirstChild{
							Document: doc,
							Bookmark: topLevelBookmark.Bookmark,
						})
						Expect(err).To(BeNil())
						Expect(topLevelBookmarkSibling).To(Not(BeNil()))
						if topLevelBookmarkSibling != nil {
							Expect(topLevelBookmarkSibling.Bookmark).To(BeNil())
						}

						// Check that we have a sibling
						sibling, err := pdfiumContainer.FPDFBookmark_GetNextSibling(&requests.FPDFBookmark_GetNextSibling{
							Document: doc,
							Bookmark: *topLevelBookmark.Bookmark,
						})
						Expect(err).To(BeNil())
						Expect(sibling).To(Not(BeNil()))
						if sibling != nil {
							Expect(sibling.Bookmark).To(Not(BeNil()))

							// Check title of bookmark.
							siblingTitle, err := pdfiumContainer.FPDFBookmark_GetTitle(&requests.FPDFBookmark_GetTitle{
								Bookmark: *sibling.Bookmark,
							})
							Expect(err).To(BeNil())
							Expect(siblingTitle).To(Equal(&responses.FPDFBookmark_GetTitle{
								Title: "A Good Ending",
							}))

							// Check that we have no children
							siblingChildren, err := pdfiumContainer.FPDFBookmark_GetFirstChild(&requests.FPDFBookmark_GetFirstChild{
								Document: doc,
								Bookmark: sibling.Bookmark,
							})
							Expect(err).To(BeNil())
							Expect(siblingChildren).To(Not(BeNil()))
							if siblingChildren != nil {
								Expect(siblingChildren.Bookmark).To(BeNil())
							}

							// Check that we have no sibling
							siblingSibling, err := pdfiumContainer.FPDFBookmark_GetNextSibling(&requests.FPDFBookmark_GetNextSibling{
								Document: doc,
								Bookmark: *sibling.Bookmark,
							})
							Expect(err).To(BeNil())
							Expect(siblingSibling).To(Not(BeNil()))
							if sibling != nil {
								Expect(siblingSibling.Bookmark).To(BeNil())
							}
						}
					}
				})
			})

			When("FPDFBookmark_Find is called", func() {
				It("it returns the correct bookmark when there is a match", func() {
					bookmark, err := pdfiumContainer.FPDFBookmark_Find(&requests.FPDFBookmark_Find{
						Document: doc,
						Title:    "A Good Beginning",
					})
					Expect(err).To(BeNil())
					Expect(bookmark).To(Not(BeNil()))
					if bookmark != nil {
						Expect(bookmark.Bookmark).To(Not(BeNil()))
					}
				})

				It("it returns the no bookmark when there is no match", func() {
					bookmark, err := pdfiumContainer.FPDFBookmark_Find(&requests.FPDFBookmark_Find{
						Document: doc,
						Title:    "No Good Beginning",
					})
					Expect(err).To(BeNil())
					Expect(bookmark).To(Equal(&responses.FPDFBookmark_Find{}))
				})
			})
		})
	})
}
