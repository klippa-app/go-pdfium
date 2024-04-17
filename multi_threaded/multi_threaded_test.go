package multi_threaded_test

import (
	"os"
	"strings"
	"time"

	"github.com/klippa-app/go-pdfium/multi_threaded"
	"github.com/klippa-app/go-pdfium/shared_tests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {
	// Set ENV to ensure resulting values.
	err := os.Setenv("TZ", "UTC")
	Expect(err).To(BeNil())

	args := []string{"run", "-exec", "env DYLD_LIBRARY_PATH=/opt/pdfium/lib"}
	experimental := os.Getenv("IS_EXPERIMENTAL") == "1"
	xfa := os.Getenv("IS_XFA") == "1"

	if experimental || xfa {
		tags := []string{}
		if experimental {
			tags = append(tags, "pdfium_experimental")
		}
		if xfa {
			tags = append(tags, "pdfium_xfa")
		}
		args = append(args, "-tags", strings.Join(tags, ","))
	}

	args = append(args, "../examples/multi_threaded/worker/main.go")

	pool := multi_threaded.Init(multi_threaded.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
		Command: multi_threaded.Command{
			BinPath:      "go",             // Only do this while developing, on production put the actual binary path in here. You should not want the Go runtime on production.
			Args:         args,             // This is a reference to the worker package, this can be left empty when using a direct binary path.
			StartTimeout: time.Minute * 15, // Some test environments are real slow.
		},
	})

	shared_tests.PdfiumPool = pool

	instance, err := pool.GetInstance(time.Second * 30)
	Expect(err).To(BeNil())

	shared_tests.PdfiumInstance = instance
	shared_tests.TestDataPath = "../shared_tests"
	shared_tests.TestType = "multi"
})

var _ = AfterSuite(func() {
	err := shared_tests.PdfiumInstance.Close()
	Expect(err).To(BeNil())

	err = shared_tests.PdfiumPool.Close()
	Expect(err).To(BeNil())
})

var _ = Describe("Multi Threaded", func() {
	shared_tests.Import()
})
