package webassembly_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

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
		MaxTotal: 1, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
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
})

var _ = AfterEach(func() {
	Eventually(Goroutines).ShouldNot(HaveLeaked())
})
