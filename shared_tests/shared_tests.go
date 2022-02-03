package shared_tests

import (
	"os"

	"github.com/klippa-app/go-pdfium"
)

func RunTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	// Set ENV to ensure resulting values.
	os.Setenv("TZ", "UTC")

	RunBookmarkTests(pdfiumContainer, testsPath, prefix)
	RunActionTests(pdfiumContainer, testsPath, prefix)
	RunfpdfDocTests(pdfiumContainer, testsPath, prefix)
	RunfpdfCatalogTests(pdfiumContainer, testsPath, prefix)
	RunfpdfSignatureTests(pdfiumContainer, testsPath, prefix)
	RunfpdfThumbnailTests(pdfiumContainer, testsPath, prefix)
	RunfpdfEditTests(pdfiumContainer, testsPath, prefix)
	RunfpdfExtTests(pdfiumContainer, testsPath, prefix)
	RunfpdfFlattenTests(pdfiumContainer, testsPath, prefix)
	RunfpdfViewTests(pdfiumContainer, testsPath, prefix)
	RunRenderTests(pdfiumContainer, testsPath, prefix)
	RunDocumentTests(pdfiumContainer, testsPath, prefix)
	RunTextTests(pdfiumContainer, testsPath, prefix)
	RunPageTests(pdfiumContainer, testsPath, prefix)
}
