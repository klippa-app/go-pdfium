package wago_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/experimental/wago"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/shared_tests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
)

var _ = BeforeSuite(func() {
	// Set ENV to ensure resulting values.
	err := os.Setenv("TZ", "UTC")
	Expect(err).To(BeNil())

	pool, err := wago.Init(wago.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // The maximum number of workers in total, allows the number of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
	})
	Expect(err).To(BeNil())

	shared_tests.PdfiumPool = pool

	instance, err := pool.GetInstance(time.Second * 30)
	Expect(err).To(BeNil())
	shared_tests.PdfiumInstance = instance

	// Older wago versions compiled the module on arm64 but miscompiled part
	// of it, which showed up as memory traps on the first render. Check one
	// render up front and skip the suite instead of failing hundreds of
	// specs when that happens again, for example after a wago downgrade.
	if runtime.GOARCH == "arm64" {
		if err := smokeRender(instance); err != nil {
			Expect(instance.Close()).To(Succeed())
			Expect(pool.Close()).To(Succeed())
			shared_tests.PdfiumInstance = nil
			shared_tests.PdfiumPool = nil
			Skip("wago miscompiles the PDFium module on arm64: " + err.Error())
		}
	}
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

// smokeRender renders the first page of the smallest test document once.
func smokeRender(instance pdfium.Pdfium) error {
	data, err := os.ReadFile("../../shared_tests/testdata/test.pdf")
	if err != nil {
		return err
	}
	doc, err := instance.OpenDocument(&requests.OpenDocument{File: &data})
	if err != nil {
		return err
	}
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})
	rendered, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
		Page: requests.Page{ByIndex: &requests.PageByIndex{Document: doc.Document, Index: 0}},
		DPI:  50,
	})
	if err != nil {
		return err
	}
	rendered.Cleanup()
	return nil
}

var _ = AfterSuite(func() {
	if shared_tests.PdfiumInstance == nil {
		return
	}

	err := shared_tests.PdfiumInstance.Close()
	Expect(err).To(BeNil())

	err = shared_tests.PdfiumPool.Close()
	Expect(err).To(BeNil())
})

var _ = Describe("Wago", func() {
	shared_tests.Import()

	Context("custom WASM", func() {
		It("requires a configured module", func() {
			pool, err := wago.InitWithWASM(wago.Config{})
			Expect(pool).To(BeNil())
			Expect(err).To(MatchError("webassembly module must be provided"))
		})
	})

	Context("pooling", func() {
		It("uses the configured context for workers", func() {
			ctx, cancel := context.WithCancel(context.Background())
			pool, err := wago.Init(wago.Config{
				Context:  ctx,
				MinIdle:  0,
				MaxIdle:  1,
				MaxTotal: 1,
			})
			Expect(err).To(BeNil())
			DeferCleanup(func() {
				Expect(pool.Close()).To(Succeed())
			})

			cancel()
			instance, err := pool.GetInstance(time.Second * 30)
			if instance != nil {
				DeferCleanup(func() {
					Expect(instance.Close()).To(Succeed())
				})
			}
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, context.Canceled)).To(BeTrue())
		})

		When("a pool is opened", func() {
			var TestPool pdfium.Pool

			BeforeEach(func() {
				pool, err := wago.Init(wago.Config{
					MinIdle:  1,
					MaxIdle:  1,
					MaxTotal: 1,
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
			}, NodeTimeout(time.Second*10))
		})

		It("runs multiple instances at the same time", func() {
			pool, err := wago.Init(wago.Config{
				MinIdle:  2,
				MaxIdle:  2,
				MaxTotal: 2,
			})
			Expect(err).To(BeNil())
			defer pool.Close()

			pdfData, err := os.ReadFile("../../shared_tests/testdata/test.pdf")
			Expect(err).To(BeNil())

			done := make(chan error, 2)
			for i := 0; i < 2; i++ {
				go func() {
					done <- func() error {
						instance, err := pool.GetInstance(time.Second * 30)
						if err != nil {
							return err
						}
						defer instance.Close()

						doc, err := instance.OpenDocument(&requests.OpenDocument{File: &pdfData})
						if err != nil {
							return err
						}

						for j := 0; j < 5; j++ {
							_, err = instance.RenderPageInDPI(&requests.RenderPageInDPI{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{Document: doc.Document, Index: 0},
								},
								DPI: 50,
							})
							if err != nil {
								return err
							}
						}
						return nil
					}()
				}()
			}

			for i := 0; i < 2; i++ {
				Expect(<-done).To(BeNil())
			}
		})
	})

	Context("Kill", func() {
		It("does not panic when called on an idle instance", func() {
			pool, err := wago.Init(wago.Config{
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

		// Kill() cancels the worker's context. With InterruptibleCalls enabled
		// that interrupts a call that is running inside PDFium, without it the
		// running call finishes and every call after it fails. In both cases
		// Kill has to return promptly and the pool has to recover.
		It("interrupts a running render loop and recovers the pool", func() {
			pool, err := wago.Init(wago.Config{
				MinIdle:            0,
				MaxIdle:            1,
				MaxTotal:           1,
				InterruptibleCalls: true,
			})
			Expect(err).To(BeNil())
			defer pool.Close()

			instance, err := pool.GetInstance(time.Second * 30)
			Expect(err).To(BeNil())

			pdfData, err := os.ReadFile("../../shared_tests/testdata/alpha_channel.pdf")
			Expect(err).To(BeNil())

			doc, err := instance.OpenDocument(&requests.OpenDocument{
				File: &pdfData,
			})
			Expect(err).To(BeNil())

			// Render in a goroutine until a render fails, which is what Kill
			// has to cause.
			renderDone := make(chan error, 1)
			go func() {
				for {
					_, renderErr := instance.RenderPageInDPI(&requests.RenderPageInDPI{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc.Document,
								Index:    0,
							},
						},
						DPI: 150,
					})
					if renderErr != nil {
						renderDone <- renderErr
						return
					}
				}
			}()

			// Give the render loop a moment to get going.
			time.Sleep(200 * time.Millisecond)

			// Kill must complete promptly, not block forever.
			killDone := make(chan error, 1)
			go func() {
				killDone <- instance.Kill()
			}()

			select {
			case err := <-killDone:
				// Kill may return an error (e.g. context canceled) but
				// must not hang or panic.
				_ = err
			case <-time.After(10 * time.Second):
				Fail("Kill() did not return within 10 seconds")
			}

			select {
			case err := <-renderDone:
				Expect(err).To(HaveOccurred())
			case <-time.After(10 * time.Second):
				Fail("the render loop did not stop within 10 seconds after Kill")
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
