package shared_tests

import "github.com/klippa-app/go-pdfium"

func RunTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	RunBookmarksTests(pdfiumContainer, testsPath, prefix)
	RunfpdfDocTests(pdfiumContainer, testsPath, prefix)
	RunfpdfCatalogTests(pdfiumContainer, testsPath, prefix)
	RunRenderTests(pdfiumContainer, testsPath, prefix)
	RunDocumentTests(pdfiumContainer, testsPath, prefix)
	RunTextTests(pdfiumContainer, testsPath, prefix)
	RunPageTests(pdfiumContainer, testsPath, prefix)
}
