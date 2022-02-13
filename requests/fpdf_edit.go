package requests

import "github.com/klippa-app/go-pdfium/enums"

type FPDF_CreateNewDocument struct{}

type FPDFPage_SetRotation struct {
	Page   Page
	Rotate enums.FPDF_PAGE_ROTATION // New value of PDF page rotation.
}

type FPDFPage_GetRotation struct {
	Page Page
}

type FPDFPage_HasTransparency struct {
	Page Page
}

type FPDFPage_New struct{}
type FPDFPage_Delete struct{}
type FPDFPage_InsertObject struct{}
type FPDFPage_RemoveObject struct{}
type FPDFPage_CountObjects struct{}
type FPDFPage_GetObject struct{}
type FPDFPage_GenerateContent struct{}
type FPDFPageObj_Destroy struct{}
type FPDFPageObj_HasTransparency struct{}
type FPDFPageObj_GetType struct{}
type FPDFPageObj_Transform struct{}
type FPDFPageObj_GetMatrix struct{}
type FPDFPageObj_SetMatrix struct{}
type FPDFPage_TransformAnnots struct{}
type FPDFPageObj_NewImageObj struct{}
type FPDFPageObj_CountMarks struct{}
type FPDFPageObj_GetMark struct{}
type FPDFPageObj_AddMark struct{}
type FPDFPageObj_RemoveMark struct{}
type FPDFPageObjMark_GetName struct{}
type FPDFPageObjMark_CountParams struct{}
type FPDFPageObjMark_GetParamKey struct{}
type FPDFPageObjMark_GetParamValueType struct{}
type FPDFPageObjMark_GetParamIntValue struct{}
type FPDFPageObjMark_GetParamStringValue struct{}
type FPDFPageObjMark_GetParamBlobValue struct{}
type FPDFPageObjMark_SetIntParam struct{}
type FPDFPageObjMark_SetStringParam struct{}
type FPDFPageObjMark_SetBlobParam struct{}
type FPDFPageObjMark_RemoveParam struct{}
type FPDFImageObj_LoadJpegFile struct{}
type FPDFImageObj_LoadJpegFileInline struct{}
type FPDFImageObj_SetMatrix struct{}
type FPDFImageObj_SetBitmap struct{}
type FPDFImageObj_GetBitmap struct{}
type FPDFImageObj_GetRenderedBitmap struct{}
type FPDFImageObj_GetImageDataDecoded struct{}
type FPDFImageObj_GetImageDataRaw struct{}
type FPDFImageObj_GetImageFilterCount struct{}
type FPDFImageObj_GetImageFilter struct{}
type FPDFImageObj_GetImageMetadata struct{}
type FPDFPageObj_CreateNewPath struct{}
type FPDFPageObj_CreateNewRect struct{}
type FPDFPageObj_GetBounds struct{}
type FPDFPageObj_SetBlendMode struct{}
type FPDFPageObj_SetStrokeColor struct{}
type FPDFPageObj_GetStrokeColor struct{}
type FPDFPageObj_SetStrokeWidth struct{}
type FPDFPageObj_GetStrokeWidth struct{}
type FPDFPageObj_GetLineJoin struct{}
type FPDFPageObj_SetLineJoin struct{}
type FPDFPageObj_GetLineCap struct{}
type FPDFPageObj_SetLineCap struct{}
type FPDFPageObj_SetFillColor struct{}
type FPDFPageObj_GetFillColor struct{}
type FPDFPageObj_GetDashPhase struct{}
type FPDFPageObj_SetDashPhase struct{}
type FPDFPageObj_GetDashCount struct{}
type FPDFPageObj_GetDashArray struct{}
type FPDFPageObj_SetDashArray struct{}
type FPDFPath_CountSegments struct{}
type FPDFPath_GetPathSegment struct{}
type FPDFPathSegment_GetPoint struct{}
type FPDFPathSegment_GetType struct{}
type FPDFPathSegment_GetClose struct{}
type FPDFPath_MoveTo struct{}
type FPDFPath_LineTo struct{}
type FPDFPath_BezierTo struct{}
type FPDFPath_Close struct{}
type FPDFPath_SetDrawMode struct{}
type FPDFPath_GetDrawMode struct{}
type FPDFPageObj_NewTextObj struct{}
type FPDFText_SetText struct{}
type FPDFText_SetCharcodes struct{}
type FPDFText_LoadFont struct{}
type FPDFText_LoadStandardFont struct{}
type FPDFTextObj_GetFontSize struct{}
type FPDFFont_Close struct{}
type FPDFPageObj_CreateTextObj struct{}
type FPDFTextObj_GetTextRenderMode struct{}
type FPDFTextObj_SetTextRenderMode struct{}
type FPDFTextObj_GetText struct{}
type FPDFTextObj_GetFont struct{}
type FPDFFont_GetFontName struct{}
type FPDFFont_GetFlags struct{}
type FPDFFont_GetWeight struct{}
type FPDFFont_GetItalicAngle struct{}
type FPDFFont_GetAscent struct{}
type FPDFFont_GetDescent struct{}
type FPDFFont_GetGlyphWidth struct{}
type FPDFFont_GetGlyphPath struct{}
type FPDFGlyphPath_CountGlyphSegments struct{}
type FPDFGlyphPath_GetGlyphPathSegment struct{}
type FPDFFormObj_CountObjects struct{}
type FPDFFormObj_GetObject struct{}
