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
				pageCount, err := pdfium.FPDF_GetMetaText(&requests.FPDF_GetMetaText{
					Tag: "Creator",
				})
				Expect(err).To(MatchError("document not given"))
				Expect(pageCount).To(BeNil())
			})
		})
	})
})
