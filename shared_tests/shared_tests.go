package shared_tests

import "github.com/klippa-app/go-pdfium"

func RunTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	RunRenderTests(pdfiumContainer, testsPath, prefix)
	RunDocumentTests(pdfiumContainer, testsPath, prefix)
	RunTextTests(pdfiumContainer, testsPath, prefix)
	RunPageTests(pdfiumContainer, testsPath, prefix)
}
