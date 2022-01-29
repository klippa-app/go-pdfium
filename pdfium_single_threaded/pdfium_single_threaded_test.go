package pdfium_single_threaded_test

import (
	"github.com/klippa-app/go-pdfium/pdfium_single_threaded"
	"github.com/klippa-app/go-pdfium/shared_tests"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Single Threaded", func() {
	Pdfium := pdfium_single_threaded.Init()
	shared_tests.RunRenderTests(Pdfium, "single")
	shared_tests.RunDocumentTests(Pdfium, "single")
	shared_tests.RunTextTests(Pdfium, "single")
})
