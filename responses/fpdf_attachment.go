package responses

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
)

type FPDFDoc_GetAttachmentCount struct {
	AttachmentCount int
}

type FPDFDoc_AddAttachment struct {
	Attachment references.FPDF_ATTACHMENT
}

type FPDFDoc_GetAttachment struct {
	Index      int
	Attachment references.FPDF_ATTACHMENT
}

type FPDFDoc_DeleteAttachment struct {
	Index int
}

type FPDFAttachment_GetName struct {
	Name string
}

type FPDFAttachment_HasKey struct {
	Key    string
	HasKey bool
}

type FPDFAttachment_GetValueType struct {
	Key       string
	ValueType enums.FPDF_OBJECT_TYPE
}

type FPDFAttachment_SetStringValue struct {
	Key   string
	Value string
}

type FPDFAttachment_GetStringValue struct {
	Key   string
	Value string
}

type FPDFAttachment_SetFile struct{}

type FPDFAttachment_GetFile struct {
	Contents []byte // nil when not found.
}
