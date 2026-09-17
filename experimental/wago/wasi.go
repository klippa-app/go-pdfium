package wago

import (
	"context"
	"fmt"
	"io"

	wagort "github.com/wago-org/wago"
	"github.com/wago-org/wasi/p1"
)

// WASI rights bits, see the wasi_snapshot_preview1 rights type.
const (
	wasiRightFDRead  = 1 << 1
	wasiRightFDWrite = 1 << 6

	wasiErrnoNotCapable = 76
)

// wasiImports returns the wasi_snapshot_preview1 functions the module imports,
// backed by wago's WASI plugin, for one instance.
//
// Emscripten's open() does not request the fd_read and fd_write rights when it
// calls path_open, and wago's WASI plugin enforces rights strictly, so every
// read on a file PDFium opened would fail with ENOTCAPABLE. wazero ignores
// rights altogether. path_open is wrapped here to request the file rights the
// mount grants: first read and write, then read only, then what Emscripten
// asked for.
func wasiImports(ctx context.Context, compiledModule *wagort.Module, stdout, stderr io.Writer, random io.Reader, mounts []Preopen) (*wagort.Imports, error) {
	plugin := p1.Imports(p1.Config{
		Stdout:  stdout,
		Stderr:  stderr,
		Rand:    random,
		Context: ctx,
		Mounts:  mounts,
		Args:    []string{"pdfium"},
	})

	imports := wagort.NewImports()
	for _, spec := range compiledModule.Imports() {
		if spec.Module != "wasi_snapshot_preview1" {
			continue
		}

		value, ok := plugin.Lookup(spec.Module, spec.Name)
		if !ok {
			return nil, fmt.Errorf("the WASI plugin does not provide %s.%s", spec.Module, spec.Name)
		}

		fn, ok := value.(wagort.CallerHostCallFunc)
		if !ok {
			return nil, fmt.Errorf("unexpected WASI plugin function type %T for %s", value, spec.Name)
		}

		if spec.Name == "path_open" {
			fn = pathOpenWithFileRights(fn)
		}

		imports.HostFunc(spec.Module, spec.Name, fn).Params(spec.Params...).Results(spec.Results...)
	}

	return imports, nil
}

func pathOpenWithFileRights(original wagort.CallerHostCallFunc) wagort.CallerHostCallFunc {
	return func(caller wagort.Caller, call wagort.HostCall) {
		params := call.ParamSlots()
		requestedRights := params[5]

		for _, extraRights := range []uint64{wasiRightFDRead | wasiRightFDWrite, wasiRightFDRead, 0} {
			params[5] = requestedRights | extraRights
			original(caller, call)
			if call.ResultSlots()[0] != wasiErrnoNotCapable {
				return
			}
		}
	}
}
