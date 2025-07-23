package multi_threaded_test

import (
	"context"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/multi_threaded"
	"github.com/klippa-app/go-pdfium/shared_tests"
)

var _ = BeforeSuite(func() {
	// Set ENV to ensure resulting values.
	err := os.Setenv("TZ", "UTC")
	Expect(err).To(BeNil())

	args := []string{"run", "-exec", "env DYLD_LIBRARY_PATH=/opt/pdfium/lib"}
	experimental := os.Getenv("IS_EXPERIMENTAL")
	if experimental == "1" {
		args = append(args, []string{"-tags", "pdfium_experimental"}...)
	}

	args = append(args, "../examples/multi_threaded/worker/main.go")

	pool := multi_threaded.Init(multi_threaded.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // The maximum number of workers in total, allows the number of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
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

	Eventually(Goroutines).ShouldNot(HaveLeaked())
})

var _ = Describe("Multi Threaded", func() {
	shared_tests.Import()

	Context("pooling", func() {
		When("a pool is opened", func() {
			var TestPool pdfium.Pool

			BeforeEach(func() {
				args := []string{"run", "-exec", "env DYLD_LIBRARY_PATH=/opt/pdfium/lib"}
				experimental := os.Getenv("IS_EXPERIMENTAL")
				if experimental == "1" {
					args = append(args, []string{"-tags", "pdfium_experimental"}...)
				}

				args = append(args, "../examples/multi_threaded/worker/main.go")

				pool := multi_threaded.Init(multi_threaded.Config{
					MinIdle:  1, // Makes sure that at least x workers are always available
					MaxIdle:  1, // Makes sure that at most x workers are ever available
					MaxTotal: 1, // The maximum number of workers in total, allows the number of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
					Command: multi_threaded.Command{
						BinPath:      "go",             // Only do this while developing, on production put the actual binary path in here. You should not want the Go runtime on production.
						Args:         args,             // This is a reference to the worker package, this can be left empty when using a direct binary path.
						StartTimeout: time.Minute * 15, // Some test environments are real slow.
					},
				})
				TestPool = pool
			})

			When("an instance is retrieved", func() {
				var TestInstance pdfium.Pdfium

				BeforeEach(func() {
					instance, err := TestPool.GetInstance(time.Second * 30)
					Expect(err).To(BeNil())
					TestInstance = instance
				})

				It("allows the pool to be closed when all the instances are closed", func() {
					err := TestInstance.Close()
					Expect(err).To(BeNil())
				})

				It("allows the pool to be closed when there are still open instances", func() {
					// Do nothing here, we're testing closing the pool without closing the instance.
				})
			})

			AfterEach(func(ctx context.Context) {
				err := TestPool.Close()
				Expect(err).To(BeNil())
			}, NodeTimeout(time.Second))
		})
	})
})
