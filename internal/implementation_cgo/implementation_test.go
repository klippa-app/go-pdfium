package implementation_cgo_test

import (
	"os"
	"time"

	"github.com/klippa-app/go-pdfium/internal/implementation_cgo"
	"github.com/klippa-app/go-pdfium/shared_tests"
	"github.com/klippa-app/go-pdfium/single_threaded"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
)

var _ = BeforeSuite(func() {
	implementation_cgo.InitLibrary(nil)

	// Set ENV to ensure resulting values.
	err := os.Setenv("TZ", "UTC")
	Expect(err).To(BeNil())

	pool := single_threaded.Init(single_threaded.Config{})
	shared_tests.PdfiumPool = pool

	instance, err := pool.GetInstance(time.Second * 30)
	Expect(err).To(BeNil())
	shared_tests.PdfiumInstance = instance
	shared_tests.TestDataPath = "../../shared_tests"
	shared_tests.TestType = "internal"
})

var _ = AfterSuite(func() {
	err := shared_tests.PdfiumInstance.Close()
	Expect(err).To(BeNil())

	err = shared_tests.PdfiumPool.Close()
	Expect(err).To(BeNil())
})

var _ = Describe("Implementation", func() {
	When("pinged", func() {
		It("pongs", func() {
			pdfium := implementation_cgo.Pdfium.GetInstance()
			resp, err := pdfium.Ping()
			Expect(err).To(BeNil())
			Expect(resp).To(Equal("Pong"))
			pdfium.Close()
		})
	})

	shared_tests.Import()
})

var _ = AfterEach(func() {
	Eventually(Goroutines).ShouldNot(HaveLeaked())
})
