#if defined(WIN32)
#define GO_FPDF_EXPORT __declspec(dllexport)
#define GO_FPDF_CALLCONV __stdcall
#else
#define GO_FPDF_EXPORT
#define GO_FPDF_CALLCONV
#endif
