package requests

import "github.com/klippa-app/go-pdfium/references"

type FPDF_GetSignatureCount struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetSignatureObject struct {
	Document references.FPDF_DOCUMENT
	Index    int
}

type FPDFSignatureObj_GetContents struct {
	Signature references.FPDF_SIGNATURE
}

type FPDFSignatureObj_GetByteRange struct {
	Signature references.FPDF_SIGNATURE
}

type FPDFSignatureObj_GetSubFilter struct {
	Signature references.FPDF_SIGNATURE
}

type FPDFSignatureObj_GetReason struct {
	Signature references.FPDF_SIGNATURE
}

type FPDFSignatureObj_GetTime struct {
	Signature references.FPDF_SIGNATURE
}

type FPDFSignatureObj_GetDocMDPPermission struct {
	Signature references.FPDF_SIGNATURE
}
