package renderutil_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRenderutil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Renderutil Suite")
}
