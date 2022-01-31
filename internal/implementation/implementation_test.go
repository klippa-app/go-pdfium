package implementation_test

import (
	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/shared_tests"
	"github.com/klippa-app/go-pdfium/single_threaded"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Implementation", func() {
	implementation.InitLibrary()

	When("pinged", func() {
		It("pongs", func() {
			pdfium := implementation.Pdfium.GetInstance()
			resp, err := pdfium.Ping()
			Expect(err).To(BeNil())
			Expect(resp).To(Equal("Pong"))
			pdfium.Close()
		})
	})

	pool := single_threaded.Init()
	instance, err := pool.GetInstance(time.Second * 30)
	if err != nil {
		Expect(err).To(BeNil())
		return
	}

	shared_tests.RunTests(instance, "../../shared_tests", "internal")
})
