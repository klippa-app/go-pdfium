package implementation_webassembly_test

import (
	"os"
	"path/filepath"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/shared_tests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Render bitmap allocation failure", func() {
	// A render whose bitmap cannot be allocated must return an error instead
	// of silently rendering into a NULL bitmap and returning a garbage image.
	// The guest memory of memoryLimitedPool is capped so the bitmap
	// allocation reliably fails.
	It("returns an error when the bitmap can not be allocated", func() {
		pdfData, err := os.ReadFile(filepath.Join(shared_tests.TestDataPath, "testdata", "test.pdf"))
		Expect(err).To(Succeed())
		instance, document := openMemoryLimitedDocument(pdfData)

		// A4 at 600 DPI is a ~139 MB BGRA bitmap, above the memory cap.
		resp, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
			DPI:  600,
			Page: requests.Page{ByIndex: &requests.PageByIndex{Document: document, Index: 0}},
		})
		Expect(err).To(MatchError(ContainSubstring("could not create bitmap")))
		Expect(resp).To(BeNil())

		// The instance must still be usable for renders that do fit.
		smallResp, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
			DPI:  72,
			Page: requests.Page{ByIndex: &requests.PageByIndex{Document: document, Index: 0}},
		})
		Expect(err).To(Succeed())
		Expect(smallResp).ToNot(BeNil())
		smallResp.Cleanup()
	})
})
