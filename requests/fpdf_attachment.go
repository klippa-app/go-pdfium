package requests

import "github.com/klippa-app/go-pdfium/references"

type FPDFDoc_GetAttachmentCount struct {
	Document references.FPDF_DOCUMENT
}

type FPDFDoc_AddAttachment struct {
	Document references.FPDF_DOCUMENT
	Name     string
}

type FPDFDoc_GetAttachment struct {
	Document references.FPDF_DOCUMENT
	Index    int
}

type FPDFDoc_DeleteAttachment struct {
	Document references.FPDF_DOCUMENT
	Index    int
}

type FPDFAttachment_GetName struct {
	Attachment references.FPDF_ATTACHMENT
}

type FPDFAttachment_HasKey struct {
	Attachment references.FPDF_ATTACHMENT
	Key        string
}

type FPDFAttachment_GetValueType struct {
	Attachment references.FPDF_ATTACHMENT
	Key        string
}

type FPDFAttachment_SetStringValue struct {
	Attachment references.FPDF_ATTACHMENT
	Key        string
	Value      string
}

type FPDFAttachment_GetStringValue struct {
	Attachment references.FPDF_ATTACHMENT
	Key        string
}

type FPDFAttachment_SetFile struct {
	Attachment references.FPDF_ATTACHMENT
	Contents   []byte
}

type FPDFAttachment_GetFile struct {
	Attachment references.FPDF_ATTACHMENT
}

type FPDFAttachment_GetSubtype struct {
	Attachment references.FPDF_ATTACHMENT
}
