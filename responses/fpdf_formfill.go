package responses

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
)

type FPDFDOC_InitFormFillEnvironment struct {
	FormHandle references.FPDF_FORMHANDLE
}

type FPDFDOC_ExitFormFillEnvironment struct{}

type FORM_OnAfterLoadPage struct{}

type FORM_OnBeforeClosePage struct{}

type FORM_DoDocumentJSAction struct{}

type FORM_DoDocumentOpenAction struct{}

type FORM_DoDocumentAAction struct{}

type FORM_DoPageAAction struct{}

type FORM_OnMouseMove struct{}

type FORM_OnMouseWheel struct{}

type FORM_OnFocus struct{}

type FORM_OnLButtonDown struct{}

type FORM_OnRButtonDown struct{}

type FORM_OnLButtonUp struct{}

type FORM_OnRButtonUp struct{}

type FORM_OnLButtonDoubleClick struct{}

type FORM_OnKeyDown struct{}

type FORM_OnKeyUp struct{}

type FORM_OnChar struct{}

type FORM_GetFocusedText struct {
	FocusedText string
}

type FORM_GetSelectedText struct {
	SelectedText string
}

type FORM_ReplaceSelection struct{}

type FORM_SelectAllText struct{}

type FORM_CanUndo struct {
	CanUndo bool
}

type FORM_CanRedo struct {
	CanRedo bool
}

type FORM_Undo struct{}

type FORM_Redo struct{}

type FORM_ForceToKillFocus struct{}

type FORM_GetFocusedAnnot struct {
	PageIndex  int
	Annotation references.FPDF_ANNOTATION
}

type FORM_SetFocusedAnnot struct{}

type FPDFPage_HasFormFieldAtPoint struct {
	FieldType enums.FPDF_FORMFIELD // The type of the form field; -1 indicates no field.
}

type FPDFPage_FormFieldZOrderAtPoint struct {
	ZOrder int // The z-order of the form field; -1 indicates no field. Higher numbers are closer to the front.
}

type FPDF_SetFormFieldHighlightColor struct{}

type FPDF_SetFormFieldHighlightAlpha struct{}

type FPDF_RemoveFormFieldHighlight struct{}

type FPDF_FFLDraw struct{}

type FPDF_GetFormType struct {
	FormType enums.FPDF_FORMTYPE
}

type FORM_SetIndexSelected struct{}

type FORM_IsIndexSelected struct {
	IsIndexSelected bool
}

type FPDF_LoadXFA struct{}
