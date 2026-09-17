package wago

import (
	"context"
	"io"

	wagort "github.com/wago-org/wago"
	"github.com/wago-org/wasi/p1"
)

// wasiImports returns the wasi_snapshot_preview1 module for one instance,
// backed by wago's WASI plugin. Every instance gets its own set because the
// plugin keeps the file descriptor table in it.
//
// The plugin enforces WASI rights strictly. The bundled PDFium module is built
// with an Emscripten that requests the fd_read and fd_write rights that match
// the open() access mode, stock Emscripten 6.0.9 standalone builds do not
// (they request neither), which makes every read fail with ENOTCAPABLE and
// PDFium report an incorrect format. Custom modules have to be built with the
// same fix.
func wasiImports(ctx context.Context, stdout, stderr io.Writer, random io.Reader, mounts []Preopen) *wagort.Imports {
	return p1.Imports(p1.Config{
		Stdout:  stdout,
		Stderr:  stderr,
		Rand:    random,
		Context: ctx,
		Mounts:  mounts,
		Args:    []string{"pdfium"},
	})
}
