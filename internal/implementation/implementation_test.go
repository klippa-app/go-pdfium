package implementation_test

import (
	"github.com/klippa-app/go-pdfium/internal/implementation"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Implementation", func() {
	implementation.InitLibrary()
	pdfium := implementation.Pdfium{}
	AfterEach(func() {
		pdfium.Close()
	})

	When("pinged", func() {
		It("pongs", func() {
			resp, err := pdfium.Ping()
			Expect(err).To(BeNil())
			Expect(resp).To(Equal("Pong"))
		})
	})
})
