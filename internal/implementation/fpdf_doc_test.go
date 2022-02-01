package implementation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/requests"
)

var _ = Describe("fpdf_doc", func() {
	pdfium := implementation.Pdfium.GetInstance()

	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when getting the page metadata", func() {
				FPDF_GetMetaText, err := pdfium.FPDF_GetMetaText(&requests.FPDF_GetMetaText{
					Tag: "Creator",
				})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_GetMetaText).To(BeNil())
			})

			It("returns an error when calling FPDFBookmark_GetFirstChild", func() {
				FPDFBookmark_GetFirstChild, err := pdfium.FPDFBookmark_GetFirstChild(&requests.FPDFBookmark_GetFirstChild{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFBookmark_GetFirstChild).To(BeNil())
			})

			It("returns an error when calling FPDFBookmark_GetNextSibling", func() {
				FPDFBookmark_GetNextSibling, err := pdfium.FPDFBookmark_GetNextSibling(&requests.FPDFBookmark_GetNextSibling{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFBookmark_GetNextSibling).To(BeNil())
			})

			It("returns an error when calling FPDFBookmark_Find", func() {
				FPDFBookmark_Find, err := pdfium.FPDFBookmark_Find(&requests.FPDFBookmark_Find{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFBookmark_Find).To(BeNil())
			})
		})
	})

	Context("no bookmark", func() {
		When("is given", func() {
			It("returns an error when calling FPDFBookmark_GetTitle", func() {
				FPDFBookmark_GetTitle, err := pdfium.FPDFBookmark_GetTitle(&requests.FPDFBookmark_GetTitle{})
				Expect(err).To(MatchError("bookmark not given"))
				Expect(FPDFBookmark_GetTitle).To(BeNil())
			})
		})
	})
})
