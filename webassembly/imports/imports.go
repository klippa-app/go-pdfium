package imports

import (
	"context"

	"github.com/klippa-app/go-pdfium/internal/implementation_webassembly"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/emscripten"
)

// Instantiate instantiates the "env" module used by Emscripten into the
// runtime default namespace.
//
// # Notes
//
//   - Closing the wazero.Runtime has the same effect as closing the result.
//   - To add more functions to the "env" module, use FunctionExporter.
//   - To instantiate into another wazero.Namespace, use FunctionExporter.
func Instantiate(ctx context.Context, r wazero.Runtime, mod wazero.CompiledModule) (api.Closer, error) {
	builder := r.NewHostModuleBuilder("env")
	exporter, err := emscripten.NewFunctionExporterForModule(mod)
	if err != nil {
		return nil, err
	}
	exporter.ExportFunctions(builder)
	NewFunctionExporter().ExportFunctions(builder)
	return builder.Instantiate(ctx)
}

// FunctionExporter configures the functions in the "env" module used by
// Emscripten.
type FunctionExporter interface {
	// ExportFunctions builds functions to export with a wazero.HostModuleBuilder
	// named "env".
	ExportFunctions(builder wazero.HostModuleBuilder)
}

// NewFunctionExporter returns a FunctionExporter object with trace disabled.
func NewFunctionExporter() FunctionExporter {
	return &functionExporter{}
}

type functionExporter struct{}

// ExportFunctions implements FunctionExporter.ExportFunctions
func (e *functionExporter) ExportFunctions(b wazero.HostModuleBuilder) {
	for i := range implementation_webassembly.HostFunctions {
		hostFunction := implementation_webassembly.HostFunctions[i]
		b.NewFunctionBuilder().
			WithGoModuleFunction(goModuleFunction(hostFunction.Call), valueTypes(hostFunction.Params), valueTypes(hostFunction.Results)).
			Export(hostFunction.Name)
	}
}

// goModuleFunction adapts a runtime independent host function to wazero.
type goModuleFunction func(ctx context.Context, mod implementation_webassembly.Module, stack []uint64)

func (fn goModuleFunction) Call(ctx context.Context, mod api.Module, stack []uint64) {
	fn(ctx, NewModule(mod), stack)
}

func valueTypes(types []implementation_webassembly.ValueType) []api.ValueType {
	result := make([]api.ValueType, len(types))
	for i, valueType := range types {
		switch valueType {
		case implementation_webassembly.ValueTypeI32:
			result[i] = api.ValueTypeI32
		case implementation_webassembly.ValueTypeI64:
			result[i] = api.ValueTypeI64
		case implementation_webassembly.ValueTypeF32:
			result[i] = api.ValueTypeF32
		case implementation_webassembly.ValueTypeF64:
			result[i] = api.ValueTypeF64
		}
	}
	return result
}
