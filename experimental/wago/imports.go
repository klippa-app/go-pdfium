package wago

import (
	"github.com/klippa-app/go-pdfium/internal/implementation_webassembly"

	wagort "github.com/wago-org/wago"
)

// hostImports returns the "env" module that PDFium imports: the go-pdfium
// callbacks plus the Emscripten function the standalone build still needs.
// The imports are bound to one module because a host function has to know
// which instance it is executing for.
func hostImports(m *module) *wagort.Imports {
	imports := wagort.NewImports()

	for i := range implementation_webassembly.HostFunctions {
		hostFunction := implementation_webassembly.HostFunctions[i]
		stackSize := max(len(hostFunction.Params), len(hostFunction.Results), 1)

		imports.HostFunc("env", hostFunction.Name, func(caller wagort.Caller, call wagort.HostCall) {
			// The host function may re-enter the guest (malloc) which may in
			// turn call a host function again, so every call gets its own
			// stack.
			stack := make([]uint64, stackSize)
			copy(stack, call.ParamSlots())

			exit := m.enterHostCall(caller)
			hostFunction.Call(m.ctx, m, stack)
			exit()

			if len(hostFunction.Results) > 0 {
				call.ResultSlots()[0] = stack[0]
			}
		}).Params(valTypes(hostFunction.Params)...).Results(valTypes(hostFunction.Results)...)
	}

	// Emscripten calls this after memory.grow. Memory is always read through
	// the runtime, so there is nothing to invalidate on the host side.
	imports.HostFunc("env", "emscripten_notify_memory_growth", func(int32) {})

	return imports
}

func valTypes(types []implementation_webassembly.ValueType) []wagort.ValType {
	result := make([]wagort.ValType, len(types))
	for i, valueType := range types {
		switch valueType {
		case implementation_webassembly.ValueTypeI32:
			result[i] = wagort.ValI32
		case implementation_webassembly.ValueTypeI64:
			result[i] = wagort.ValI64
		case implementation_webassembly.ValueTypeF32:
			result[i] = wagort.ValF32
		case implementation_webassembly.ValueTypeF64:
			result[i] = wagort.ValF64
		}
	}
	return result
}
