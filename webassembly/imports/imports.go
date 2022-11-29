package imports

import (
	"context"
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
func Instantiate(ctx context.Context, r wazero.Runtime) (api.Closer, error) {
	builder := r.NewHostModuleBuilder("env")
	emscripten.NewFunctionExporter().ExportFunctions(builder)
	NewFunctionExporter().ExportFunctions(builder)
	return builder.Instantiate(ctx, r)
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
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FILEACCESS_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("FPDF_FILEACCESS_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FILEWRITE_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("FPDF_FILEWRITE_CB")
}
