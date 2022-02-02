package implementation_test

import (
	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/requests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_catalog", func() {
	pdfium := implementation.Pdfium.GetInstance()

	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when getting the document tagged status", func() {
				isTagged, err := pdfium.FPDFCatalog_IsTagged(&requests.FPDFCatalog_IsTagged{})
				Expect(err).To(MatchError("document not given"))
				Expect(isTagged).To(BeNil())
			})
		})
	})
})
