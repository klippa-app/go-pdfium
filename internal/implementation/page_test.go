package implementation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/requests"
)

var _ = Describe("Page", func() {
	pdfium := implementation.Pdfium.GetInstance()

	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when getting the pdf page rotation", func() {
				pageCount, err := pdfium.GetPageRotation(&requests.GetPageRotation{
					Page: requests.Page{
						Index: 0,
					},
				})
				Expect(err).To(MatchError("document not given"))
				Expect(pageCount).To(BeNil())
			})

			It("returns an error when getting the pdf page transparency", func() {
				pageCount, err := pdfium.GetPageTransparency(&requests.GetPageTransparency{
					Page: requests.Page{
						Index: 0,
					},
				})
				Expect(err).To(MatchError("document not given"))
				Expect(pageCount).To(BeNil())
			})

			It("returns an error when flattening a pdf page", func() {
				pageCount, err := pdfium.FlattenPage(&requests.FlattenPage{
					Page: requests.Page{
						Index: 0,
					},
				})
				Expect(err).To(MatchError("document not given"))
				Expect(pageCount).To(BeNil())
			})
		})
	})
})
