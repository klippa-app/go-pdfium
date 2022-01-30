package single_threaded_test

import (
	"github.com/klippa-app/go-pdfium/shared_tests"
	"github.com/klippa-app/go-pdfium/single_threaded"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Single Threaded", func() {
	Pdfium := single_threaded.Init()
	shared_tests.RunTests(Pdfium, "../shared_tests", "single")
})
