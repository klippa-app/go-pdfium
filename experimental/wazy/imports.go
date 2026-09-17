package wazy

import (
	"context"

	"github.com/klippa-app/go-pdfium/internal/implementation_webassembly"

	wazyrt "github.com/samyfodil/wazy"
	"github.com/samyfodil/wazy/api"
	"github.com/samyfodil/wazy/imports/emscripten"
)

// instantiateEnv instantiates the "env" module the PDFium module imports: the
// Emscripten functions it needs plus go-pdfium's host callbacks.
func instantiateEnv(ctx context.Context, r wazyrt.Runtime, mod wazyrt.CompiledModule) (api.Closer, error) {
	builder := r.NewHostModuleBuilder("env")
	exporter, err := emscripten.NewFunctionExporterForModule(mod)
	if err != nil {
		return nil, err
	}
	exporter.ExportFunctions(builder)

	for i := range implementation_webassembly.HostFunctions {
		hostFunction := implementation_webassembly.HostFunctions[i]
		builder.NewFunctionBuilder().
			WithGoModuleFunction(goModuleFunction(hostFunction.Call), valueTypes(hostFunction.Params), valueTypes(hostFunction.Results)).
			Export(hostFunction.Name)
	}

	return builder.Instantiate(ctx)
}

// goModuleFunction adapts a runtime independent host function to wazy.
type goModuleFunction func(ctx context.Context, mod implementation_webassembly.Module, stack []uint64)

func (fn goModuleFunction) Call(ctx context.Context, mod api.Module, stack []uint64) {
	fn(ctx, newModule(mod), stack)
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
