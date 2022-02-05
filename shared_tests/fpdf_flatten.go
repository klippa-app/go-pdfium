package shared_tests

import (
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func RunfpdfFlattenTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("fpdf_flatten", func() {
		Context("no document", func() {
			When("is opened", func() {
				It("returns an error when flattening a pdf page", func() {
					pageCount, err := pdfiumContainer.FPDFPage_Flatten(&requests.FPDFPage_Flatten{
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
