package implementation_test

import (
	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/pdfium_single_threaded"
	"github.com/klippa-app/go-pdfium/shared_tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Implementation", func() {
	implementation.InitLibrary()

	When("pinged", func() {
		It("pongs", func() {
			pdfium := implementation.Pdfium{}
			resp, err := pdfium.Ping()
			Expect(err).To(BeNil())
			Expect(resp).To(Equal("Pong"))
			pdfium.Close()
		})
	})

	Pdfium := pdfium_single_threaded.Init()
	shared_tests.RunRenderTests(Pdfium, "../../shared_tests", "internal")
	shared_tests.RunDocumentTests(Pdfium, "../../shared_tests", "internal")
	shared_tests.RunTextTests(Pdfium, "../../shared_tests", "internal")
})
