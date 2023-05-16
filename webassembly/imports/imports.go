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
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FILEACCESS_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("FPDF_FILEACCESS_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FILEWRITE_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("FPDF_FILEWRITE_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FX_FILEAVAIL_IS_DATA_AVAILABLE_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("FX_FILEAVAIL_IS_DATA_AVAILABLE_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FX_DOWNLOADHINTS_ADD_SEGMENT_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).Export("FX_DOWNLOADHINTS_ADD_SEGMENT_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(UNSUPPORT_INFO_HANDLER_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).Export("UNSUPPORT_INFO_HANDLER_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FSDK_SetTimeFunction_CB{}, []api.ValueType{}, []api.ValueType{api.ValueTypeI64}).Export("FSDK_SetTimeFunction_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FSDK_SetLocaltimeFunction_CB{}, []api.ValueType{api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("FSDK_SetLocaltimeFunction_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FORMFILLINFO_Release_CB{}, []api.ValueType{api.ValueTypeI32}, []api.ValueType{}).Export("FPDF_FORMFILLINFO_Release_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FORMFILLINFO_FFI_Invalidate_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeF64, api.ValueTypeF64, api.ValueTypeF64, api.ValueTypeF64}, []api.ValueType{}).Export("FPDF_FORMFILLINFO_FFI_Invalidate_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FORMFILLINFO_FFI_OutputSelectedRect_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeF64, api.ValueTypeF64, api.ValueTypeF64, api.ValueTypeF64}, []api.ValueType{}).Export("FPDF_FORMFILLINFO_FFI_OutputSelectedRect_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FORMFILLINFO_FFI_SetCursor_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).Export("FPDF_FORMFILLINFO_FFI_SetCursor_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FORMFILLINFO_FFI_SetTimer_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("FPDF_FORMFILLINFO_FFI_SetTimer_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FORMFILLINFO_FFI_KillTimer_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).Export("FPDF_FORMFILLINFO_FFI_KillTimer_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FORMFILLINFO_FFI_GetLocalTime_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).Export("FPDF_FORMFILLINFO_FFI_GetLocalTime_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FORMFILLINFO_FFI_OnChange_CB{}, []api.ValueType{api.ValueTypeI32}, []api.ValueType{}).Export("FPDF_FORMFILLINFO_FFI_OnChange_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FORMFILLINFO_FFI_GetPage_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("FPDF_FORMFILLINFO_FFI_GetPage_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FORMFILLINFO_FFI_GetCurrentPage_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("FPDF_FORMFILLINFO_FFI_GetCurrentPage_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FORMFILLINFO_FFI_GetRotation_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("FPDF_FORMFILLINFO_FFI_GetRotation_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FORMFILLINFO_FFI_ExecuteNamedAction_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).Export("FPDF_FORMFILLINFO_FFI_ExecuteNamedAction_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FORMFILLINFO_FFI_SetTextFieldFocus_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).Export("FPDF_FORMFILLINFO_FFI_SetTextFieldFocus_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FORMFILLINFO_FFI_DoURIAction_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).Export("FPDF_FORMFILLINFO_FFI_DoURIAction_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(FPDF_FORMFILLINFO_FFI_DoGoToAction_CB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).Export("FPDF_FORMFILLINFO_FFI_DoGoToAction_CB")
	b.NewFunctionBuilder().WithGoModuleFunction(IFSDK_PAUSE_NeedToPauseNow_CB{}, []api.ValueType{api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("IFSDK_PAUSE_NeedToPauseNow_CB")
}
