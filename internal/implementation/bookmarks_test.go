package implementation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/requests"
)

var _ = Describe("bookmarks", func() {
	pdfium := implementation.Pdfium.GetInstance()

	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when calling GetBookmarks", func() {
				GetBookmarks, err := pdfium.GetBookmarks(&requests.GetBookmarks{})
				Expect(err).To(MatchError("document not given"))
				Expect(GetBookmarks).To(BeNil())
			})
		})
	})
})
