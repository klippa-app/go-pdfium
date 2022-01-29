package implementation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/requests"
)

var _ = Describe("Document", func() {
	pdfium := implementation.Pdfium{}

	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when getting the page count", func() {
				pageCount, err := pdfium.GetPageCount(&requests.GetPageCount{})
				Expect(err).To(MatchError("no current document"))
				Expect(pageCount).To(BeNil())
			})
		})
	})
})
