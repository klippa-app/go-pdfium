package wago_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPdfiumWago(t *testing.T) {
	RegisterFailHandler(Fail)

	suiteConfig, reporterConfig := GinkgoConfiguration()

	// wago's amd64 backend miscompiles part of Little-CMS, which makes the page
	// with an ICC based colour space render with a black instead of a white
	// background. Skip that shared render test until wago is fixed, see the
	// package documentation.
	suiteConfig.SkipStrings = append(suiteConfig.SkipStrings, "a PDF file that uses an alpha channel")

	RunSpecs(t, "Wago Suite", suiteConfig, reporterConfig)
}
