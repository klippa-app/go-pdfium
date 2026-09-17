package implementation_webassembly_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/shared_tests"
	"github.com/klippa-app/go-pdfium/webassembly"
	"github.com/tetratelabs/wazero"
	wazeroapi "github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
)

// interpreterMode is set at init time from the WAZERO_INTERPRETER environment
// variable. When true the whole suite uses wazero's interpreter engine instead
// of the default compiler engine.
var interpreterMode = os.Getenv("WAZERO_INTERPRETER") == "1"

// runtimeConfig returns the wazero RuntimeConfig that matches the current
// mode (interpreter or compiler). Every pool in this suite should be created
// with it so that the interpreter CI job really runs all specs interpreted.
func runtimeConfig() wazero.RuntimeConfig {
	// The bundled PDFium module requires the exception handling core feature.
	features := wazeroapi.CoreFeaturesV2 | experimental.CoreFeaturesExceptionHandling
	if interpreterMode {
		return wazero.NewRuntimeConfigInterpreter().WithCoreFeatures(features)
	}
	return wazero.NewRuntimeConfig().WithCoreFeatures(features)
}

// memoryLimitedPool is shared by the specs that need a capped guest memory.
// It has a single instance, specs take it with GetInstance and give it back
// through DeferCleanup. Compiling the module is expensive, so the pool is
// created once for the whole suite.
var memoryLimitedPool pdfium.Pool

// memoryLimitPages caps the guest memory of memoryLimitedPool. 768 pages of
// 64 KiB is 48 MiB: enough to initialize PDFium, open a document and render or
// encode a ~23 MiB bitmap, but not enough to hold that bitmap twice or to
// allocate the ~139 MiB bitmap of an A4 page at 600 DPI.
const memoryLimitPages = 768

var _ = BeforeSuite(func() {
	// Set ENV to ensure resulting values.
	err := os.Setenv("TZ", "UTC")
	Expect(err).To(BeNil())

	pool, err := webassembly.Init(webassembly.Config{
		MinIdle:       1,
		MaxIdle:       1,
		MaxTotal:      1,
		RuntimeConfig: runtimeConfig(),
	})
	Expect(err).To(BeNil())
	shared_tests.PdfiumPool = pool

	memoryLimitedPool, err = webassembly.Init(webassembly.Config{
		MinIdle:       1,
		MaxIdle:       1,
		MaxTotal:      1,
		RuntimeConfig: runtimeConfig().WithMemoryLimitPages(memoryLimitPages),
	})
	Expect(err).To(BeNil())

	instance, err := pool.GetInstance(time.Second * 30)
	Expect(err).To(BeNil())
	shared_tests.PdfiumInstance = instance
	shared_tests.TestDataPath = "../../shared_tests"

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

	err = memoryLimitedPool.Close()
	Expect(err).To(BeNil())
})

var _ = Describe("Implementation", func() {
	shared_tests.Import()
})

var _ = AfterEach(func() {
	Eventually(Goroutines).ShouldNot(HaveLeaked())
})
