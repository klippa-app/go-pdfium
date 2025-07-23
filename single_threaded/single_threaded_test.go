package single_threaded_test

import (
	"context"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/shared_tests"
	"github.com/klippa-app/go-pdfium/single_threaded"
)

var _ = BeforeSuite(func() {
	// Set ENV to ensure resulting values.
	err := os.Setenv("TZ", "UTC")
	Expect(err).To(BeNil())

	pool := single_threaded.Init(single_threaded.Config{})
	shared_tests.PdfiumPool = pool

	instance, err := pool.GetInstance(time.Second * 30)
	Expect(err).To(BeNil())
	shared_tests.PdfiumInstance = instance
	shared_tests.TestDataPath = "../shared_tests"
	shared_tests.TestType = "single"
})

var _ = AfterSuite(func() {
	err := shared_tests.PdfiumInstance.Close()
	Expect(err).To(BeNil())

	err = shared_tests.PdfiumPool.Close()
	Expect(err).To(BeNil())
})

var _ = Describe("Single Threaded", func() {
	shared_tests.Import()

	Context("pooling", func() {
		When("a pool is opened", func() {
			var TestPool pdfium.Pool

			BeforeEach(func() {
				pool := single_threaded.Init(single_threaded.Config{})
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
