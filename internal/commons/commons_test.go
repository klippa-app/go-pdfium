package commons_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCommons(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commons Suite")
}

var _ = Describe("canary", func() {
	It("passes", func() {
		Expect(true).To(BeTrue(), "This is good. Canary test passing")
	})
})
