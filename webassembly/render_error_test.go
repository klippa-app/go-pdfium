package webassembly_test

import (
	"os"
	"time"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Render bitmap allocation failure", func() {
	// A render whose bitmap cannot be allocated must return an error instead
	// of silently rendering into a NULL bitmap and returning a garbage image.
	// The instance memory is capped so the bitmap allocation reliably fails.
	It("returns an error when the bitmap can not be allocated", func() {
		pool, err := webassembly.Init(webassembly.Config{
			MinIdle:  1,
			MaxIdle:  1,
			MaxTotal: 1,
			// 768 pages = 48 MB: enough to initialize PDFium and open the
			// document, not enough for the ~139 MB bitmap requested below.
			RuntimeConfig: runtimeConfig().WithMemoryLimitPages(768),
		})
		Expect(err).To(BeNil())
		defer pool.Close()

		instance, err := pool.GetInstance(time.Second * 30)
		Expect(err).To(BeNil())
		defer instance.Close()

		pdfData, err := os.ReadFile("../shared_tests/testdata/test.pdf")
		Expect(err).To(BeNil())

		doc, err := instance.OpenDocument(&requests.OpenDocument{File: &pdfData})
		Expect(err).To(BeNil())
		defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})

		// A4 at 600 DPI is a ~139 MB BGRA bitmap, above the memory cap.
		resp, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
			DPI:  600,
			Page: requests.Page{ByIndex: &requests.PageByIndex{Document: doc.Document, Index: 0}},
		})
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("could not create bitmap"))
		Expect(resp).To(BeNil())

		// The instance must still be usable for renders that do fit.
		smallResp, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
			DPI:  72,
			Page: requests.Page{ByIndex: &requests.PageByIndex{Document: doc.Document, Index: 0}},
		})
		Expect(err).To(BeNil())
		Expect(smallResp).ToNot(BeNil())
		smallResp.Cleanup()
	})
})
