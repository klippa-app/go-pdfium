package pdfium_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPdfium(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pdfium Suite")
}

var _ = Describe("canary", func() {
	It("passes", func() {
		Expect(true).To(BeTrue(), "This is good. Canary test passing")
	})
})
