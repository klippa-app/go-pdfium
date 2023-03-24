package worker

import (
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/internal/implementation_cgo"
)

func StartWorker(config *pdfium.LibraryConfig) {
	implementation_cgo.StartPlugin(config)
}
