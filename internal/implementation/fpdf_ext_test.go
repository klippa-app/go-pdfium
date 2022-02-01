package implementation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/requests"
)

var _ = Describe("fpdf_ext", func() {
	pdfium := implementation.Pdfium.GetInstance()

	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when getting the page mode", func() {
				pageMode, err := pdfium.FPDFDoc_GetPageMode(&requests.FPDFDoc_GetPageMode{})
				Expect(err).To(MatchError("document not given"))
				Expect(pageMode).To(BeNil())
			})
		})
	})
})
