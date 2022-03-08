package requests

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/structs"
)

type FPDFDOC_InitFormFillEnvironment struct {
	Document     references.FPDF_DOCUMENT
	FormFillInfo structs.FPDF_FORMFILLINFO
}

type FPDFDOC_ExitFormFillEnvironment struct {
	FormHandle references.FPDF_FORMHANDLE
}

type FORM_OnAfterLoadPage struct {
	Page       Page
	FormHandle references.FPDF_FORMHANDLE
}

type FORM_OnBeforeClosePage struct {
	Page       Page
	FormHandle references.FPDF_FORMHANDLE
}

type FORM_DoDocumentJSAction struct {
	FormHandle references.FPDF_FORMHANDLE
}

type FORM_DoDocumentOpenAction struct {
	FormHandle references.FPDF_FORMHANDLE
}

type FORM_DoDocumentAAction struct {
	FormHandle references.FPDF_FORMHANDLE
	AAType     enums.FPDFDOC_AACTION
}

type FORM_DoPageAAction struct {
	Page       Page
	FormHandle references.FPDF_FORMHANDLE
	AAType     enums.FPDFPAGE_AACTION
}

type FORM_OnMouseMove struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	Modifier   int     // Indicates whether various virtual keys are down.
	PageX      float64 // Specifies the x-coordinate of the cursor in PDF user space.
	PageY      float64 // Specifies the y-coordinate of the cursor in PDF user space.
}

type FORM_OnMouseWheel struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	Modifier   int                    // Indicates whether various virtual keys are down.
	PageCoord  structs.FPDF_FS_POINTF // Specifies the coordinates of the cursor in PDF user space.
	DeltaX     int                    // Specifies the amount of wheel movement on the x-axis, in units of platform-agnostic wheel deltas. Negative values mean left.
	DeltaY     int                    // Specifies the amount of wheel movement on the y-axis, in units of platform-agnostic wheel deltas. Negative values mean down.

}

type FORM_OnFocus struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	Modifier   int     // Indicates whether various virtual keys are down.
	PageX      float64 // Specifies the x-coordinate of the cursor in PDF user space.
	PageY      float64 // Specifies the y-coordinate of the cursor in PDF user space.
}

type FORM_OnLButtonDown struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	Modifier   int     // Indicates whether various virtual keys are down.
	PageX      float64 // Specifies the x-coordinate of the cursor in PDF user space.
	PageY      float64 // Specifies the y-coordinate of the cursor in PDF user space.
}

type FORM_OnRButtonDown struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	Modifier   int     // Indicates whether various virtual keys are down.
	PageX      float64 // Specifies the x-coordinate of the cursor in PDF user space.
	PageY      float64 // Specifies the y-coordinate of the cursor in PDF user space.
}

type FORM_OnLButtonUp struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	Modifier   int     // Indicates whether various virtual keys are down.
	PageX      float64 // Specifies the x-coordinate of the cursor in PDF user space.
	PageY      float64 // Specifies the y-coordinate of the cursor in PDF user space.
}

type FORM_OnRButtonUp struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	Modifier   int     // Indicates whether various virtual keys are down.
	PageX      float64 // Specifies the x-coordinate of the cursor in PDF user space.
	PageY      float64 // Specifies the y-coordinate of the cursor in PDF user space.
}

type FORM_OnLButtonDoubleClick struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	Modifier   int     // Indicates whether various virtual keys are down.
	PageX      float64 // Specifies the x-coordinate of the cursor in PDF user space.
	PageY      float64 // Specifies the y-coordinate of the cursor in PDF user space.
}

type FORM_OnKeyDown struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	NKeyCode   enums.FWL_VKEYCODE  // The virtual-key code of the given key (see fpdf_fwlevent.h for virtual key codes).
	Modifier   enums.FWL_EVENTFLAG // Mask of key flags (see fpdf_fwlevent.h for key flag values).
}

type FORM_OnKeyUp struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	NKeyCode   enums.FWL_VKEYCODE  // The virtual-key code of the given key (see fpdf_fwlevent.h for virtual key codes).
	Modifier   enums.FWL_EVENTFLAG // Mask of key flags (see fpdf_fwlevent.h for key flag values).
}

type FORM_OnChar struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	NChar      int                 // The character code value itself.
	Modifier   enums.FWL_EVENTFLAG // Mask of key flags (see fpdf_fwlevent.h for key flag values).
}

type FORM_GetFocusedText struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
}

type FORM_GetSelectedText struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
}

type FORM_ReplaceSelection struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	Text       string
}

type FORM_SelectAllText struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
}

type FORM_CanUndo struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
}

type FORM_CanRedo struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
}

type FORM_Undo struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
}

type FORM_Redo struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
}

type FORM_ForceToKillFocus struct {
	FormHandle references.FPDF_FORMHANDLE
}

type FORM_GetFocusedAnnot struct {
	FormHandle references.FPDF_FORMHANDLE
}

type FORM_SetFocusedAnnot struct {
	FormHandle references.FPDF_FORMHANDLE
	Annotation references.FPDF_ANNOTATION
}

type FPDFPage_HasFormFieldAtPoint struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	PageX      float64 // X position in PDF "user space".
	PageY      float64 // Y position in PDF "user space".
}

type FPDFPage_FormFieldZOrderAtPoint struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	PageX      float64 // X position in PDF "user space".
	PageY      float64 // Y position in PDF "user space".
}

type FPDF_SetFormFieldHighlightColor struct {
	FormHandle references.FPDF_FORMHANDLE
	FieldType  enums.FPDF_FORMFIELD
	Color      uint64 // The highlight color of the form field. Constructed by 0xxxrrggbb
}

type FPDF_SetFormFieldHighlightAlpha struct {
	FormHandle references.FPDF_FORMHANDLE
	Alpha      uint8
}

type FPDF_RemoveFormFieldHighlight struct {
	FormHandle references.FPDF_FORMHANDLE
}

type FPDF_FFLDraw struct {
	FormHandle references.FPDF_FORMHANDLE
	Bitmap     references.FPDF_BITMAP
	Page       Page
	StartX     int
	StartY     int
	SizeX      int
	SizeY      int
	Rotate     enums.FPDF_PAGE_ROTATION
	Flags      enums.FPDF_RENDER_FLAG
}

type FPDF_GetFormType struct {
	Document references.FPDF_DOCUMENT
}

type FORM_SetIndexSelected struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	Index      int
	Selected   bool
}

type FORM_IsIndexSelected struct {
	FormHandle references.FPDF_FORMHANDLE
	Page       Page
	Index      int
}

type FPDF_LoadXFA struct {
	Document references.FPDF_DOCUMENT
}
