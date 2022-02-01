package implementation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/requests"
)

var _ = Describe("fpdfview", func() {
	pdfium := implementation.Pdfium.GetInstance()

	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when getting the pdf version", func() {
				pageCount, err := pdfium.FPDF_GetFileVersion(&requests.FPDF_GetFileVersion{})
				Expect(err).To(MatchError("document not given"))
				Expect(pageCount).To(BeNil())
			})

			It("returns an error when getting the doc permissions", func() {
				pageCount, err := pdfium.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{})
				Expect(err).To(MatchError("document not given"))
				Expect(pageCount).To(BeNil())
			})

			It("returns an error when getting the doc revision number of security handler", func() {
				pageCount, err := pdfium.FPDF_GetSecurityHandlerRevision(&requests.FPDF_GetSecurityHandlerRevision{})
				Expect(err).To(MatchError("document not given"))
				Expect(pageCount).To(BeNil())
			})

			It("returns an error when getting the page count", func() {
				pageCount, err := pdfium.FPDF_GetPageCount(&requests.FPDF_GetPageCount{})
				Expect(err).To(MatchError("document not given"))
				Expect(pageCount).To(BeNil())
			})
		})
	})
})
