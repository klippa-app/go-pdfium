package requests

import "github.com/klippa-app/go-pdfium/references"

type FPDFDoc_GetJavaScriptActionCount struct {
	Document references.FPDF_DOCUMENT
}

type FPDFDoc_GetJavaScriptAction struct {
	Document references.FPDF_DOCUMENT
	Index    int
}

type FPDFDoc_CloseJavaScriptAction struct {
	JavaScriptAction references.FPDF_JAVASCRIPT_ACTION
}

type FPDFJavaScriptAction_GetName struct {
	JavaScriptAction references.FPDF_JAVASCRIPT_ACTION
}

type FPDFJavaScriptAction_GetScript struct {
	JavaScriptAction references.FPDF_JAVASCRIPT_ACTION
}
