package responses

import "github.com/klippa-app/go-pdfium/references"

type FPDFDoc_GetJavaScriptActionCount struct {
	JavaScriptActionCount int
}

type FPDFDoc_GetJavaScriptAction struct {
	Index            int
	JavaScriptAction references.FPDF_JAVASCRIPT_ACTION
}

type FPDFDoc_CloseJavaScriptAction struct{}

type FPDFJavaScriptAction_GetName struct {
	Name string
}

type FPDFJavaScriptAction_GetScript struct {
	Script string
}
