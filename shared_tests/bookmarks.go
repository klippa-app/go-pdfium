package shared_tests

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"io/ioutil"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/references"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

func RunBookmarksTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("bookmarks", func() {
		Context("no document", func() {
			When("is opened", func() {
				It("returns an error when calling GetBookmarks", func() {
					GetBookmarks, err := pdfiumContainer.GetBookmarks(&requests.GetBookmarks{})
					Expect(err).To(MatchError("document not given"))
					Expect(GetBookmarks).To(BeNil())
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

			When("GetBookmarks is called", func() {
				It("returns the correct bookmarks", func() {
					metadata, err := pdfiumContainer.GetBookmarks(&requests.GetBookmarks{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(metadata).To(Equal(&responses.GetBookmarks{}))
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

			When("GetBookmarks is called", func() {
				It("returns the correct bookmarks", func() {
					bookmarks, err := pdfiumContainer.GetBookmarks(&requests.GetBookmarks{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(bookmarks).To(Not(BeNil()))

					if bookmarks.Bookmarks != nil {
						Expect(bookmarks.Bookmarks).To(HaveLen(2))

						if len(bookmarks.Bookmarks) == 2 {
							Expect(bookmarks.Bookmarks).To(ContainElement(MatchAllFields(Fields{
								"Reference": Not(BeNil()),
								"Title":     Equal("A Good Beginning"),
								"Action":    Not(BeNil()),
								"Children":  HaveLen(0),
								"Dest":      BeNil(),
							})))

							if bookmarks.Bookmarks[0].Action != nil {
								Expect(*bookmarks.Bookmarks[0].Action).To(MatchAllFields(Fields{
									"Reference": Not(BeNil()),
									"Type":      Equal(enums.FPDF_ACTION_ACTION_UNSUPPORTED),
									"Dest":      BeNil(),
									"FilePath":  BeNil(),
									"URIPath":   BeNil(),
								}))
							}

							Expect(bookmarks.Bookmarks).To(ContainElement(MatchAllFields(Fields{
								"Reference": Not(BeNil()),
								"Title":     Equal("A Good Ending"),
								"Children":  HaveLen(0),
								"Action":    Not(BeNil()),
								"Dest":      BeNil(),
							})))

							if bookmarks.Bookmarks[1].Action != nil {
								Expect(*bookmarks.Bookmarks[1].Action).To(MatchAllFields(Fields{
									"Reference": Not(BeNil()),
									"Type":      Equal(enums.FPDF_ACTION_ACTION_UNSUPPORTED),
									"Dest":      BeNil(),
									"FilePath":  BeNil(),
									"URIPath":   BeNil(),
								}))
							}
						}
					}
				})
			})
		})
	})
}
