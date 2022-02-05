package single_threaded_test

import (
	"time"

	"github.com/klippa-app/go-pdfium/shared_tests"
	"github.com/klippa-app/go-pdfium/single_threaded"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Single Threaded", func() {
	pool := single_threaded.Init()
	instance, err := pool.GetInstance(time.Second * 30)
	if err != nil {
		Expect(err).To(BeNil())
		return
	}

	shared_tests.RunTests(instance, "../shared_tests", "single")
})
