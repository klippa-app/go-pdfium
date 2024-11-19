package responses

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/structs"
)

type FPDFAnnot_IsSupportedSubtype struct {
	IsSupported bool
}

type FPDFPage_CreateAnnot struct {
	Annotation references.FPDF_ANNOTATION
}

type FPDFPage_GetAnnotCount struct {
	Count int
}

type FPDFPage_GetAnnot struct {
	Annotation references.FPDF_ANNOTATION
}

type FPDFPage_GetAnnotIndex struct {
	Index int
}

type FPDFPage_CloseAnnot struct{}

type FPDFPage_RemoveAnnot struct{}

type FPDFAnnot_GetSubtype struct {
	Subtype enums.FPDF_ANNOTATION_SUBTYPE
}

type FPDFAnnot_IsObjectSupportedSubtype struct {
	IsObjectSupportedSubtype bool
}

type FPDFAnnot_UpdateObject struct{}

type FPDFAnnot_AddInkStroke struct {
	Index int
}

type FPDFAnnot_RemoveInkList struct{}

type FPDFAnnot_AppendObject struct{}

type FPDFAnnot_GetObjectCount struct {
	Count int
}

type FPDFAnnot_GetObject struct {
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFAnnot_RemoveObject struct{}

type FPDFAnnot_SetColor struct{}

type FPDFAnnot_GetColor struct {
	R uint
	G uint
	B uint
	A uint
}

type FPDFAnnot_HasAttachmentPoints struct {
	HasAttachmentPoints bool
}

type FPDFAnnot_SetAttachmentPoints struct{}

type FPDFAnnot_AppendAttachmentPoints struct{}

type FPDFAnnot_CountAttachmentPoints struct {
	Count uint64
}

type FPDFAnnot_GetAttachmentPoints struct {
	QuadPoints structs.FPDF_FS_QUADPOINTSF
}

type FPDFAnnot_SetRect struct{}

type FPDFAnnot_GetRect struct {
	Rect structs.FPDF_FS_RECTF
}

type FPDFAnnot_GetVertices struct {
	Vertices []structs.FPDF_FS_POINTF
}

type FPDFAnnot_GetInkListCount struct {
	Count uint64
}

type FPDFAnnot_GetInkListPath struct {
	Path []structs.FPDF_FS_POINTF
}

type FPDFAnnot_GetLine struct {
	Start structs.FPDF_FS_POINTF
	End   structs.FPDF_FS_POINTF
}

type FPDFAnnot_SetBorder struct{}

type FPDFAnnot_GetBorder struct {
	HorizontalRadius float32
	VerticalRadius   float32
	BorderWidth      float32
}

type FPDFAnnot_HasKey struct {
	HasKey bool
}

type FPDFAnnot_GetValueType struct {
	ValueType enums.FPDF_OBJECT_TYPE
}

type FPDFAnnot_SetStringValue struct{}

type FPDFAnnot_GetStringValue struct {
	Value string
}

type FPDFAnnot_GetNumberValue struct {
	Value float32
}

type FPDFAnnot_SetAP struct{}

type FPDFAnnot_GetAP struct {
	Value string
}

type FPDFAnnot_GetLinkedAnnot struct {
	LinkedAnnotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetFlags struct {
	Flags enums.FPDF_ANNOT_FLAG
}

type FPDFAnnot_SetFlags struct{}

type FPDFAnnot_GetFormFieldFlags struct {
	Flags enums.FPDF_FORMFLAG
}

type FPDFAnnot_GetFormFieldAtPoint struct {
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetFormAdditionalActionJavaScript struct {
	FormAdditionalActionJavaScript string
}

type FPDFAnnot_GetFormFieldName struct {
	FormFieldName string
}

type FPDFAnnot_GetFormFieldAlternateName struct {
	FormFieldAlternateName string
}

type FPDFAnnot_GetFormFieldType struct {
	FormFieldType enums.FPDF_FORMFIELD_TYPE
}

type FPDFAnnot_GetFormFieldValue struct {
	FormFieldValue string
}

type FPDFAnnot_GetOptionCount struct {
	OptionCount int
}

type FPDFAnnot_GetOptionLabel struct {
	OptionLabel string
}

type FPDFAnnot_IsOptionSelected struct {
	IsOptionSelected bool
}

type FPDFAnnot_GetFontSize struct {
	FontSize float32
}

type FPDFAnnot_GetFontColor struct {
	R uint
	G uint
	B uint
}

type FPDFAnnot_IsChecked struct {
	IsChecked bool
}

type FPDFAnnot_SetFocusableSubtypes struct{}

type FPDFAnnot_GetFocusableSubtypesCount struct {
	FocusableSubtypesCount int
}

type FPDFAnnot_GetFocusableSubtypes struct {
	FocusableSubtypes []enums.FPDF_ANNOTATION_SUBTYPE
}

type FPDFAnnot_GetLink struct {
	Link references.FPDF_LINK
}

type FPDFAnnot_GetFormControlCount struct {
	FormControlCount int
}

type FPDFAnnot_GetFormControlIndex struct {
	FormControlIndex int
}

type FPDFAnnot_GetFormFieldExportValue struct {
	Value string
}

type FPDFAnnot_SetURI struct{}

type FPDFAnnot_GetFileAttachment struct {
	Attachment references.FPDF_ATTACHMENT
}

type FPDFAnnot_AddFileAttachment struct {
	Attachment references.FPDF_ATTACHMENT
}
