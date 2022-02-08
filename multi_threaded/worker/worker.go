package worker

import (
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/internal/implementation"
)

func StartWorker(config *pdfium.LibraryConfig) {
	implementation.StartPlugin(config)
}
