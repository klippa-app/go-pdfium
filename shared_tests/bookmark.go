package shared_tests

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"io/ioutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("bookmarks", func() {
	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when calling GetBookmarks", func() {
				GetBookmarks, err := PdfiumInstance.GetBookmarks(&requests.GetBookmarks{})
				Expect(err).To(MatchError("document not given"))
				Expect(GetBookmarks).To(BeNil())
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

		When("GetBookmarks is called", func() {
			It("returns the correct bookmarks", func() {
				metadata, err := PdfiumInstance.GetBookmarks(&requests.GetBookmarks{
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

		When("GetBookmarks is called", func() {
			It("returns the correct bookmarks", func() {
				bookmarks, err := PdfiumInstance.GetBookmarks(&requests.GetBookmarks{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(bookmarks).To(Not(BeNil()))
				Expect(bookmarks.Bookmarks).To(HaveLen(2))
				Expect(bookmarks.Bookmarks).To(ContainElement(MatchAllFields(Fields{
					"Reference":  Not(BeNil()),
					"Title":      Equal("A Good Beginning"),
					"ActionInfo": Not(BeNil()),
					"DestInfo":   BeNil(),
					"Children":   HaveLen(0),
				})))

				Expect(*bookmarks.Bookmarks[0].ActionInfo).To(MatchAllFields(Fields{
					"Reference": Not(BeNil()),
					"Type":      Equal(enums.FPDF_ACTION_ACTION_UNSUPPORTED),
					"DestInfo":  BeNil(),
					"FilePath":  BeNil(),
					"URIPath":   BeNil(),
				}))

				Expect(bookmarks.Bookmarks).To(ContainElement(MatchAllFields(Fields{
					"Reference":  Not(BeNil()),
					"Title":      Equal("A Good Ending"),
					"ActionInfo": Not(BeNil()),
					"DestInfo":   BeNil(),
					"Children":   HaveLen(0),
				})))

				Expect(*bookmarks.Bookmarks[1].ActionInfo).To(MatchAllFields(Fields{
					"Reference": Not(BeNil()),
					"Type":      Equal(enums.FPDF_ACTION_ACTION_UNSUPPORTED),
					"DestInfo":  BeNil(),
					"FilePath":  BeNil(),
					"URIPath":   BeNil(),
				}))
			})
		})
	})
})
