package shared_tests

import (
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func RunfpdfViewTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("fpdf_flatten", func() {
		Context("no document", func() {
			When("is opened", func() {
				It("returns an error when getting the pdf version", func() {
					pageCount, err := pdfiumContainer.FPDF_GetFileVersion(&requests.FPDF_GetFileVersion{})
					Expect(err).To(MatchError("document not given"))
					Expect(pageCount).To(BeNil())
				})

				It("returns an error when getting the doc permissions", func() {
					pageCount, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{})
					Expect(err).To(MatchError("document not given"))
					Expect(pageCount).To(BeNil())
				})

				It("returns an error when getting the doc revision number of security handler", func() {
					pageCount, err := pdfiumContainer.FPDF_GetSecurityHandlerRevision(&requests.FPDF_GetSecurityHandlerRevision{})
					Expect(err).To(MatchError("document not given"))
					Expect(pageCount).To(BeNil())
				})

				It("returns an error when getting the page count", func() {
					pageCount, err := pdfiumContainer.FPDF_GetPageCount(&requests.FPDF_GetPageCount{})
					Expect(err).To(MatchError("document not given"))
					Expect(pageCount).To(BeNil())
				})
			})
		})
	})
}
