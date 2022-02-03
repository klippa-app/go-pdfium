package shared_tests

import (
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func RunfpdfEditTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("fpdf_edit", func() {
		Context("no document", func() {
			When("is opened", func() {
				It("returns an error when getting the pdf page rotation", func() {
					pageCount, err := pdfiumContainer.FPDFPage_GetRotation(&requests.FPDFPage_GetRotation{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Index: 0,
							},
						},
					})
					Expect(err).To(MatchError("document not given"))
					Expect(pageCount).To(BeNil())
				})

				It("returns an error when getting the pdf page transparency", func() {
					pageCount, err := pdfiumContainer.FPDFPage_HasTransparency(&requests.FPDFPage_HasTransparency{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Index: 0,
							},
						},
					})
					Expect(err).To(MatchError("document not given"))
					Expect(pageCount).To(BeNil())
				})
			})
		})
	})
}
