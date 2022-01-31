package implementation_test

import (
	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/requests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Text", func() {
	pdfium := implementation.Pdfium.GetInstance()

	Context("no document", func() {
		When("is given", func() {
			Context("GetPageText()", func() {
				It("returns an error", func() {
					pageText, err := pdfium.GetPageText(&requests.GetPageText{
						Page: 0,
					})
					Expect(err).To(MatchError("Document.Ref not given"))
					Expect(pageText).To(BeNil())
				})
			})

			Context("GetPageTextStructured()", func() {
				It("returns an error", func() {
					pageTextStructured, err := pdfium.GetPageTextStructured(&requests.GetPageTextStructured{
						Page: 0,
					})
					Expect(err).To(MatchError("Document.Ref not given"))
					Expect(pageTextStructured).To(BeNil())
				})
			})
		})
	})
})
