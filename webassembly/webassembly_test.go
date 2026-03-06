package webassembly_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/shared_tests"
	"github.com/klippa-app/go-pdfium/webassembly"
	"github.com/tetratelabs/wazero"

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

	Context("Kill", func() {
		// Kill() previously set i.pool = nil before calling
		// i.pool.workerPool.InvalidateObject(), causing a nil pointer
		// dereference on every call. The deferred recover() caught the
		// panic silently, but the module was never actually invalidated.
		It("does not panic when called on an idle instance", func() {
			pool, err := webassembly.Init(webassembly.Config{
				MinIdle:  0,
				MaxIdle:  1,
				MaxTotal: 1,
			})
			Expect(err).To(BeNil())
			defer pool.Close()

			instance, err := pool.GetInstance(time.Second * 30)
			Expect(err).To(BeNil())

			err = instance.Kill()
			Expect(err).To(BeNil())

			// The pool should still be usable after Kill.
			instance2, err := pool.GetInstance(time.Second * 30)
			Expect(err).To(BeNil())
			err = instance2.Close()
			Expect(err).To(BeNil())
		})

		// bug_451265.pdf is a malformed PDF from the PDFium test corpus
		// that causes pdfium to hang indefinitely during rendering.
		// Kill() can interrupt stuck execution when the caller enables
		// WithCloseOnContextDone on the RuntimeConfig. This works because
		// Kill() cancels the worker's context, and WithCloseOnContextDone
		// makes wazero respect that cancellation during execution.
		It("interrupts a stuck WASM execution and recovers the pool", func() {
			pool, err := webassembly.Init(webassembly.Config{
				MinIdle:       0,
				MaxIdle:       1,
				MaxTotal:      1,
				RuntimeConfig: wazero.NewRuntimeConfig().WithCloseOnContextDone(true),
			})
			Expect(err).To(BeNil())
			defer pool.Close()

			instance, err := pool.GetInstance(time.Second * 30)
			Expect(err).To(BeNil())

			pdfData, err := os.ReadFile("../shared_tests/testdata/bug_451265.pdf")
			Expect(err).To(BeNil())

			doc, err := instance.OpenDocument(&requests.OpenDocument{
				File: &pdfData,
			})
			Expect(err).To(BeNil())

			// Start rendering in a goroutine. This hangs indefinitely
			// because the PDF triggers a stuck state in pdfium.
			renderDone := make(chan error, 1)
			go func() {
				_, renderErr := instance.RenderToFile(&requests.RenderToFile{
					RenderPageInDPI: &requests.RenderPageInDPI{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc.Document,
								Index:    0,
							},
						},
						DPI: 72,
					},
					OutputFormat:  requests.RenderToFileOutputFormatJPG,
					OutputTarget:  requests.RenderToFileOutputTargetBytes,
					OutputQuality: 50,
				})
				renderDone <- renderErr
			}()

			// Give the render a moment to get stuck.
			time.Sleep(500 * time.Millisecond)

			// Kill must complete promptly — not block forever.
			killDone := make(chan error, 1)
			go func() {
				killDone <- instance.Kill()
			}()

			select {
			case err := <-killDone:
				// Kill may return an error (e.g. context canceled) but
				// must not hang or panic.
				_ = err
			case <-time.After(5 * time.Second):
				Fail("Kill() did not return within 5 seconds — WASM execution was not interrupted")
			}

			// The pool should still be usable: get a fresh instance.
			instance2, err := pool.GetInstance(time.Second * 30)
			Expect(err).To(BeNil())
			err = instance2.Close()
			Expect(err).To(BeNil())
		})
	})
})

var _ = AfterEach(func() {
	Eventually(Goroutines).ShouldNot(HaveLeaked())
})
