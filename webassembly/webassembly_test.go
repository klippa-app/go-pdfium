package webassembly_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/shared_tests"
	"github.com/klippa-app/go-pdfium/webassembly"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
)

var _ = BeforeSuite(func() {
	// Set ENV to ensure resulting values.
	err := os.Setenv("TZ", "UTC")
	Expect(err).To(BeNil())

	pool, err := webassembly.Init(webassembly.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // The maximum number of workers in total, allows the number of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
	})
	Expect(err).To(BeNil())

	shared_tests.PdfiumPool = pool

	instance, err := pool.GetInstance(time.Second * 30)
	Expect(err).To(BeNil())
	shared_tests.PdfiumInstance = instance
	shared_tests.TestDataPath = "../shared_tests"

	if runtime.GOOS == "windows" {
		absPath, err := filepath.Abs(shared_tests.TestDataPath)
		Expect(err).To(BeNil())

		volumeName := filepath.VolumeName(absPath)
		if volumeName != "" {
			absPath = strings.TrimPrefix(absPath, volumeName)
		}

		shared_tests.TestDataPath = strings.ReplaceAll(absPath, "\\", "/")
	}

	shared_tests.TestType = "webassembly"
})

var _ = AfterSuite(func() {
	err := shared_tests.PdfiumInstance.Close()
	Expect(err).To(BeNil())

	err = shared_tests.PdfiumPool.Close()
	Expect(err).To(BeNil())
})

var _ = Describe("Webassembly", func() {
	shared_tests.Import()

	Context("pooling", func() {
		When("a pool is opened", func() {
			var TestPool pdfium.Pool

			BeforeEach(func() {
				pool, err := webassembly.Init(webassembly.Config{
					MinIdle:  1, // Makes sure that at least x workers are always available
					MaxIdle:  1, // Makes sure that at most x workers are ever available
					MaxTotal: 1, // The maximum number of workers in total, allows the number of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
				})
				Expect(err).To(BeNil())
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

var _ = AfterEach(func() {
	Eventually(Goroutines).ShouldNot(HaveLeaked())
})
