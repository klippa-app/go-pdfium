package wazy_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPdfiumWazy(t *testing.T) {
	RegisterFailHandler(Fail)
	suiteDescription := "Wazy Suite"
	if interpreterMode {
		suiteDescription = "Wazy Interpreter Suite"
	}
	RunSpecs(t, suiteDescription)
}
