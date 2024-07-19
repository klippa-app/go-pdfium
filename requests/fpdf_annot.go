package requests

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/structs"
)

type FPDFAnnot_IsSupportedSubtype struct {
	Subtype enums.FPDF_ANNOTATION_SUBTYPE // The subtype to check.
}

type FPDFPage_CreateAnnot struct {
	Page    Page
	Subtype enums.FPDF_ANNOTATION_SUBTYPE
}

type FPDFPage_GetAnnotCount struct {
	Page Page
}

type FPDFPage_GetAnnot struct {
	Page  Page
	Index int
}

type FPDFPage_GetAnnotIndex struct {
	Page       Page
	Annotation references.FPDF_ANNOTATION
}

type FPDFPage_CloseAnnot struct {
	Annotation references.FPDF_ANNOTATION
}

type FPDFPage_RemoveAnnot struct {
	Page  Page
	Index int
}

type FPDFAnnot_GetSubtype struct {
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_IsObjectSupportedSubtype struct {
	Subtype enums.FPDF_ANNOTATION_SUBTYPE
}

type FPDFAnnot_UpdateObject struct {
	Annotation references.FPDF_ANNOTATION
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFAnnot_AddInkStroke struct {
	Annotation references.FPDF_ANNOTATION
	Points     []structs.FPDF_FS_POINTF
}

type FPDFAnnot_RemoveInkList struct {
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_AppendObject struct {
	Annotation references.FPDF_ANNOTATION
	PageObject references.FPDF_PAGEOBJECT
}

type FPDFAnnot_GetObjectCount struct {
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetObject struct {
	Annotation references.FPDF_ANNOTATION
	Index      int
}

type FPDFAnnot_RemoveObject struct {
	Annotation references.FPDF_ANNOTATION
	Index      int
}
type FPDFAnnot_SetColor struct {
	Annotation references.FPDF_ANNOTATION
	ColorType  enums.FPDFANNOT_COLORTYPE
	R          uint
	G          uint
	B          uint
	A          uint
}

type FPDFAnnot_GetColor struct {
	Annotation references.FPDF_ANNOTATION
	ColorType  enums.FPDFANNOT_COLORTYPE
}

type FPDFAnnot_HasAttachmentPoints struct {
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_SetAttachmentPoints struct {
	Annotation       references.FPDF_ANNOTATION
	Index            uint64 // Index of the set of quadpoints.
	AttachmentPoints structs.FPDF_FS_QUADPOINTSF
}

type FPDFAnnot_AppendAttachmentPoints struct {
	Annotation       references.FPDF_ANNOTATION
	AttachmentPoints structs.FPDF_FS_QUADPOINTSF
}

type FPDFAnnot_CountAttachmentPoints struct {
	Annotation references.FPDF_ANNOTATION
	Count      uint64
}
type FPDFAnnot_GetAttachmentPoints struct {
	Annotation references.FPDF_ANNOTATION
	Index      uint64
}

type FPDFAnnot_SetRect struct {
	Annotation references.FPDF_ANNOTATION
	Rect       structs.FPDF_FS_RECTF
}

type FPDFAnnot_GetRect struct {
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetVertices struct {
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetInkListCount struct {
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetInkListPath struct {
	Annotation references.FPDF_ANNOTATION
	Index      uint64
}

type FPDFAnnot_GetLine struct {
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_SetBorder struct {
	Annotation       references.FPDF_ANNOTATION
	HorizontalRadius float32
	VerticalRadius   float32
	BorderWidth      float32
}

type FPDFAnnot_GetBorder struct {
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_HasKey struct {
	Annotation references.FPDF_ANNOTATION
	Key        string
}

type FPDFAnnot_GetValueType struct {
	Annotation references.FPDF_ANNOTATION
	Key        string
}

type FPDFAnnot_SetStringValue struct {
	Annotation references.FPDF_ANNOTATION
	Key        string
	Value      string
}

type FPDFAnnot_GetStringValue struct {
	Annotation references.FPDF_ANNOTATION
	Key        string
}

type FPDFAnnot_GetNumberValue struct {
	Annotation references.FPDF_ANNOTATION
	Key        string
}

type FPDFAnnot_SetAP struct {
	Annotation     references.FPDF_ANNOTATION
	AppearanceMode enums.FPDF_ANNOT_APPEARANCEMODE
	Value          *string // If nil is passed, the AP is cleared for that mode. If the mode is Normal, APs for all modes are cleared.
}

type FPDFAnnot_GetAP struct {
	Annotation     references.FPDF_ANNOTATION
	AppearanceMode enums.FPDF_ANNOT_APPEARANCEMODE
}

type FPDFAnnot_GetLinkedAnnot struct {
	Annotation references.FPDF_ANNOTATION
	Key        string
}

type FPDFAnnot_GetFlags struct {
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_SetFlags struct {
	Annotation references.FPDF_ANNOTATION
	Flags      enums.FPDF_ANNOT_FLAG
}

type FPDFAnnot_GetFormFieldFlags struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetFormFieldAtPoint struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	Point      structs.FPDF_FS_POINTF
}

type FPDFAnnot_GetFormAdditionalActionJavaScript struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
	Event      enums.FPDF_ANNOT_AACTION
}

type FPDFAnnot_GetFormFieldName struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetFormFieldAlternateName struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetFormFieldType struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetFormFieldValue struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetOptionCount struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetOptionLabel struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
	Index      int
}

type FPDFAnnot_IsOptionSelected struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
	Index      int
}

type FPDFAnnot_GetFontSize struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetFontColor struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_IsChecked struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_SetFocusableSubtypes struct {
	FormHandle references.FPDF_FORMHANDLE
	Subtypes   []enums.FPDF_ANNOTATION_SUBTYPE
}

type FPDFAnnot_GetFocusableSubtypesCount struct {
	FormHandle references.FPDF_FORMHANDLE
}

type FPDFAnnot_GetFocusableSubtypes struct {
	FormHandle references.FPDF_FORMHANDLE
}

type FPDFAnnot_GetLink struct {
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetFormControlCount struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetFormControlIndex struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_GetFormFieldExportValue struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_SetURI struct {
	Annotation references.FPDF_ANNOTATION
	URI        string
}

type FPDFAnnot_GetFileAttachment struct {
	Document   references.FPDF_DOCUMENT
	Annotation references.FPDF_ANNOTATION
}

type FPDFAnnot_AddFileAttachment struct {
	Document   references.FPDF_DOCUMENT
	Annotation references.FPDF_ANNOTATION
	Name       string
}
