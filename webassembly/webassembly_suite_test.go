package webassembly_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPdfiumSingleThreaded(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webassembly Suite")
}
