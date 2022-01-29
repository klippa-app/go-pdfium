package single_threaded_test

import (
	"github.com/klippa-app/go-pdfium/shared_tests"
	"github.com/klippa-app/go-pdfium/single_threaded"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Single Threaded", func() {
	Pdfium := single_threaded.Init()
	shared_tests.RunRenderTests(Pdfium, "../shared_tests", "single")
	shared_tests.RunDocumentTests(Pdfium, "../shared_tests", "single")
	shared_tests.RunTextTests(Pdfium, "../shared_tests", "single")
})
