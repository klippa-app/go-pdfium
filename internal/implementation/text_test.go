package implementation_test

import (
	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/requests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Text", func() {
	pdfium := implementation.Pdfium{}

	Context("no document", func() {
		When("is opened", func() {
			Context("GetPageText()", func() {
				It("returns an error", func() {
					pageText, err := pdfium.GetPageText(&requests.GetPageText{
						Page: 0,
					})
					Expect(err).To(MatchError("no current document"))
					Expect(pageText).To(BeNil())
				})
			})

			Context("GetPageTextStructured()", func() {
				It("returns an error", func() {
					pageTextStructured, err := pdfium.GetPageTextStructured(&requests.GetPageTextStructured{
						Page: 0,
					})
					Expect(err).To(MatchError("no current document"))
					Expect(pageTextStructured).To(BeNil())
				})
			})
		})
	})
})
