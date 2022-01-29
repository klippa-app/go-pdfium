package pdfium_multi_threaded_test

import (
	"github.com/klippa-app/go-pdfium/pdfium_multi_threaded"
	"github.com/klippa-app/go-pdfium/shared_tests"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Multi Threaded", func() {
	Pdfium := pdfium_multi_threaded.Init(pdfium_multi_threaded.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
		Command: pdfium_multi_threaded.Command{
			BinPath: "go",                                                         // Only do this while developing, on production put the actual binary path in here. You should not want the Go runtime on production.
			Args:    []string{"run", "../examples/multi_threaded/worker/main.go"}, // This is a reference to the worker package, this can be left empty when using a direct binary path.
		},
	})
	shared_tests.RunRenderTests(Pdfium, "../shared_tests", "multi")
	shared_tests.RunDocumentTests(Pdfium, "../shared_tests", "multi")
	shared_tests.RunTextTests(Pdfium, "../shared_tests", "multi")
})
