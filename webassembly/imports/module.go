package imports

import (
	"github.com/klippa-app/go-pdfium/internal/implementation_webassembly"

	"github.com/tetratelabs/wazero/api"
)

// NewModule adapts a wazero module to the runtime independent Module
// interface of the implementation. The returned value compares equal for the
// same wazero module, so it can be used as a map key.
func NewModule(mod api.Module) implementation_webassembly.Module {
	return module{mod: mod}
}

type module struct {
	mod api.Module
}

func (m module) Memory() implementation_webassembly.Memory {
	return m.mod.Memory()
}

func (m module) ExportedFunction(name string) implementation_webassembly.Function {
	fn := m.mod.ExportedFunction(name)
	if fn == nil {
		return nil
	}
	return function{Function: fn}
}

// function adds ParamCount to a wazero function.
type function struct {
	api.Function
}

func (f function) ParamCount() int {
	return len(f.Definition().ParamTypes())
}
