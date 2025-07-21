//go:build pdfium_experimental
// +build pdfium_experimental

package shared_tests

import (
	"io/ioutil"
	"strconv"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_annot", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	Context("no page is given", func() {
		It("returns an error when FPDFPage_CreateAnnot is called", func() {
			FPDFPage_CreateAnnot, err := PdfiumInstance.FPDFPage_CreateAnnot(&requests.FPDFPage_CreateAnnot{})
			Expect(err).To(MatchError("either page reference or index should be given"))
			Expect(FPDFPage_CreateAnnot).To(BeNil())
		})
		It("returns an error when FPDFPage_GetAnnotCount is called", func() {
			FPDFPage_GetAnnotCount, err := PdfiumInstance.FPDFPage_GetAnnotCount(&requests.FPDFPage_GetAnnotCount{})
			Expect(err).To(MatchError("either page reference or index should be given"))
			Expect(FPDFPage_GetAnnotCount).To(BeNil())
		})
		It("returns an error when FPDFPage_GetAnnot is called", func() {
			FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{})
			Expect(err).To(MatchError("either page reference or index should be given"))
			Expect(FPDFPage_GetAnnot).To(BeNil())
		})
		It("returns an error when FPDFPage_GetAnnotIndex is called", func() {
			FPDFPage_GetAnnotIndex, err := PdfiumInstance.FPDFPage_GetAnnotIndex(&requests.FPDFPage_GetAnnotIndex{})
			Expect(err).To(MatchError("either page reference or index should be given"))
			Expect(FPDFPage_GetAnnotIndex).To(BeNil())
		})
		It("returns an error when FPDFPage_RemoveAnnot is called", func() {
			FPDFPage_RemoveAnnot, err := PdfiumInstance.FPDFPage_RemoveAnnot(&requests.FPDFPage_RemoveAnnot{})
			Expect(err).To(MatchError("either page reference or index should be given"))
			Expect(FPDFPage_RemoveAnnot).To(BeNil())
		})
	})

	Context("no annotation is given", func() {
		It("returns an error when FPDFPage_CloseAnnot is called", func() {
			FPDFPage_CloseAnnot, err := PdfiumInstance.FPDFPage_CloseAnnot(&requests.FPDFPage_CloseAnnot{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFPage_CloseAnnot).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetSubtype is called", func() {
			FPDFAnnot_GetSubtype, err := PdfiumInstance.FPDFAnnot_GetSubtype(&requests.FPDFAnnot_GetSubtype{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetSubtype).To(BeNil())
		})
		It("returns an error when FPDFAnnot_UpdateObject is called", func() {
			FPDFAnnot_UpdateObject, err := PdfiumInstance.FPDFAnnot_UpdateObject(&requests.FPDFAnnot_UpdateObject{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_UpdateObject).To(BeNil())
		})
		It("returns an error when FPDFAnnot_AddInkStroke is called", func() {
			FPDFAnnot_AddInkStroke, err := PdfiumInstance.FPDFAnnot_AddInkStroke(&requests.FPDFAnnot_AddInkStroke{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_AddInkStroke).To(BeNil())
		})
		It("returns an error when FPDFAnnot_RemoveInkList is called", func() {
			FPDFAnnot_RemoveInkList, err := PdfiumInstance.FPDFAnnot_RemoveInkList(&requests.FPDFAnnot_RemoveInkList{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_RemoveInkList).To(BeNil())
		})
		It("returns an error when FPDFAnnot_AppendObject is called", func() {
			FPDFAnnot_AppendObject, err := PdfiumInstance.FPDFAnnot_AppendObject(&requests.FPDFAnnot_AppendObject{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_AppendObject).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetObjectCount is called", func() {
			FPDFAnnot_GetObjectCount, err := PdfiumInstance.FPDFAnnot_GetObjectCount(&requests.FPDFAnnot_GetObjectCount{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetObjectCount).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetObject is called", func() {
			FPDFAnnot_GetObject, err := PdfiumInstance.FPDFAnnot_GetObject(&requests.FPDFAnnot_GetObject{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetObject).To(BeNil())
		})
		It("returns an error when FPDFAnnot_RemoveObject is called", func() {
			FPDFAnnot_RemoveObject, err := PdfiumInstance.FPDFAnnot_RemoveObject(&requests.FPDFAnnot_RemoveObject{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_RemoveObject).To(BeNil())
		})
		It("returns an error when FPDFAnnot_SetColor is called", func() {
			FPDFAnnot_SetColor, err := PdfiumInstance.FPDFAnnot_SetColor(&requests.FPDFAnnot_SetColor{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_SetColor).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetColor is called", func() {
			FPDFAnnot_GetColor, err := PdfiumInstance.FPDFAnnot_GetColor(&requests.FPDFAnnot_GetColor{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetColor).To(BeNil())
		})
		It("returns an error when FPDFAnnot_HasAttachmentPoints is called", func() {
			FPDFAnnot_HasAttachmentPoints, err := PdfiumInstance.FPDFAnnot_HasAttachmentPoints(&requests.FPDFAnnot_HasAttachmentPoints{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_HasAttachmentPoints).To(BeNil())
		})
		It("returns an error when FPDFAnnot_SetAttachmentPoints is called", func() {
			FPDFAnnot_SetAttachmentPoints, err := PdfiumInstance.FPDFAnnot_SetAttachmentPoints(&requests.FPDFAnnot_SetAttachmentPoints{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_SetAttachmentPoints).To(BeNil())
		})
		It("returns an error when FPDFAnnot_AppendAttachmentPoints is called", func() {
			FPDFAnnot_AppendAttachmentPoints, err := PdfiumInstance.FPDFAnnot_AppendAttachmentPoints(&requests.FPDFAnnot_AppendAttachmentPoints{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_AppendAttachmentPoints).To(BeNil())
		})
		It("returns an error when FPDFAnnot_CountAttachmentPoints is called", func() {
			FPDFAnnot_CountAttachmentPoints, err := PdfiumInstance.FPDFAnnot_CountAttachmentPoints(&requests.FPDFAnnot_CountAttachmentPoints{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_CountAttachmentPoints).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetAttachmentPoints is called", func() {
			FPDFAnnot_GetAttachmentPoints, err := PdfiumInstance.FPDFAnnot_GetAttachmentPoints(&requests.FPDFAnnot_GetAttachmentPoints{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetAttachmentPoints).To(BeNil())
		})
		It("returns an error when FPDFAnnot_SetRect is called", func() {
			FPDFAnnot_SetRect, err := PdfiumInstance.FPDFAnnot_SetRect(&requests.FPDFAnnot_SetRect{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_SetRect).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetRect is called", func() {
			FPDFAnnot_GetRect, err := PdfiumInstance.FPDFAnnot_GetRect(&requests.FPDFAnnot_GetRect{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetRect).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetVertices is called", func() {
			FPDFAnnot_GetVertices, err := PdfiumInstance.FPDFAnnot_GetVertices(&requests.FPDFAnnot_GetVertices{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetVertices).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetInkListCount is called", func() {
			FPDFAnnot_GetInkListCount, err := PdfiumInstance.FPDFAnnot_GetInkListCount(&requests.FPDFAnnot_GetInkListCount{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetInkListCount).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetValueType is called", func() {
			FPDFAnnot_GetValueType, err := PdfiumInstance.FPDFAnnot_GetValueType(&requests.FPDFAnnot_GetValueType{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetValueType).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetInkListPath is called", func() {
			FPDFAnnot_GetInkListPath, err := PdfiumInstance.FPDFAnnot_GetInkListPath(&requests.FPDFAnnot_GetInkListPath{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetInkListPath).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetLine is called", func() {
			FPDFAnnot_GetLine, err := PdfiumInstance.FPDFAnnot_GetLine(&requests.FPDFAnnot_GetLine{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetLine).To(BeNil())
		})
		It("returns an error when FPDFAnnot_SetBorder is called", func() {
			FPDFAnnot_SetBorder, err := PdfiumInstance.FPDFAnnot_SetBorder(&requests.FPDFAnnot_SetBorder{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_SetBorder).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetBorder is called", func() {
			FPDFAnnot_GetBorder, err := PdfiumInstance.FPDFAnnot_GetBorder(&requests.FPDFAnnot_GetBorder{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetBorder).To(BeNil())
		})
		It("returns an error when FPDFAnnot_HasKey is called", func() {
			FPDFAnnot_HasKey, err := PdfiumInstance.FPDFAnnot_HasKey(&requests.FPDFAnnot_HasKey{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_HasKey).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetValueType is called", func() {
			FPDFAnnot_GetValueType, err := PdfiumInstance.FPDFAnnot_GetValueType(&requests.FPDFAnnot_GetValueType{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetValueType).To(BeNil())
		})
		It("returns an error when FPDFAnnot_SetStringValue is called", func() {
			FPDFAnnot_SetStringValue, err := PdfiumInstance.FPDFAnnot_SetStringValue(&requests.FPDFAnnot_SetStringValue{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_SetStringValue).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetStringValue is called", func() {
			FPDFAnnot_GetStringValue, err := PdfiumInstance.FPDFAnnot_GetStringValue(&requests.FPDFAnnot_GetStringValue{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetStringValue).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetNumberValue is called", func() {
			FPDFAnnot_GetNumberValue, err := PdfiumInstance.FPDFAnnot_GetNumberValue(&requests.FPDFAnnot_GetNumberValue{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetNumberValue).To(BeNil())
		})
		It("returns an error when FPDFAnnot_SetAP is called", func() {
			FPDFAnnot_SetAP, err := PdfiumInstance.FPDFAnnot_SetAP(&requests.FPDFAnnot_SetAP{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_SetAP).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetAP is called", func() {
			FPDFAnnot_GetAP, err := PdfiumInstance.FPDFAnnot_GetAP(&requests.FPDFAnnot_GetAP{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetAP).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetLinkedAnnot is called", func() {
			FPDFAnnot_GetLinkedAnnot, err := PdfiumInstance.FPDFAnnot_GetLinkedAnnot(&requests.FPDFAnnot_GetLinkedAnnot{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetLinkedAnnot).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetFlags is called", func() {
			FPDFAnnot_GetFlags, err := PdfiumInstance.FPDFAnnot_GetFlags(&requests.FPDFAnnot_GetFlags{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetFlags).To(BeNil())
		})
		It("returns an error when FPDFAnnot_SetFlags is called", func() {
			FPDFAnnot_SetFlags, err := PdfiumInstance.FPDFAnnot_SetFlags(&requests.FPDFAnnot_SetFlags{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_SetFlags).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetLink is called", func() {
			FPDFAnnot_GetLink, err := PdfiumInstance.FPDFAnnot_GetLink(&requests.FPDFAnnot_GetLink{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetLink).To(BeNil())
		})
		It("returns an error when FPDFAnnot_SetURI is called", func() {
			FPDFAnnot_SetURI, err := PdfiumInstance.FPDFAnnot_SetURI(&requests.FPDFAnnot_SetURI{})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_SetURI).To(BeNil())
		})
	})

	Context("no form handle is given", func() {
		It("returns an error when FPDFAnnot_GetFormFieldFlags is called", func() {
			FPDFAnnot_GetFormFieldFlags, err := PdfiumInstance.FPDFAnnot_GetFormFieldFlags(&requests.FPDFAnnot_GetFormFieldFlags{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetFormFieldFlags).To(BeNil())
		})
		It("returns an error when FPDFAnnot_SetFormFieldFlags is called", func() {
			FPDFAnnot_SetFormFieldFlags, err := PdfiumInstance.FPDFAnnot_SetFormFieldFlags(&requests.FPDFAnnot_SetFormFieldFlags{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_SetFormFieldFlags).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetFormFieldAtPoint is called", func() {
			FPDFAnnot_GetFormFieldAtPoint, err := PdfiumInstance.FPDFAnnot_GetFormFieldAtPoint(&requests.FPDFAnnot_GetFormFieldAtPoint{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetFormFieldAtPoint).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetFormAdditionalActionJavaScript is called", func() {
			FPDFAnnot_GetFormAdditionalActionJavaScript, err := PdfiumInstance.FPDFAnnot_GetFormAdditionalActionJavaScript(&requests.FPDFAnnot_GetFormAdditionalActionJavaScript{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetFormAdditionalActionJavaScript).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetFormFieldName is called", func() {
			FPDFAnnot_GetFormFieldName, err := PdfiumInstance.FPDFAnnot_GetFormFieldName(&requests.FPDFAnnot_GetFormFieldName{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetFormFieldName).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetFormFieldAlternateName is called", func() {
			FPDFAnnot_GetFormFieldAlternateName, err := PdfiumInstance.FPDFAnnot_GetFormFieldAlternateName(&requests.FPDFAnnot_GetFormFieldAlternateName{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetFormFieldAlternateName).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetFormFieldType is called", func() {
			FPDFAnnot_GetFormFieldType, err := PdfiumInstance.FPDFAnnot_GetFormFieldType(&requests.FPDFAnnot_GetFormFieldType{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetFormFieldType).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetFormFieldValue is called", func() {
			FPDFAnnot_GetFormFieldValue, err := PdfiumInstance.FPDFAnnot_GetFormFieldValue(&requests.FPDFAnnot_GetFormFieldValue{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetFormFieldValue).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetOptionCount is called", func() {
			FPDFAnnot_GetOptionCount, err := PdfiumInstance.FPDFAnnot_GetOptionCount(&requests.FPDFAnnot_GetOptionCount{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetOptionCount).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetOptionLabel is called", func() {
			FPDFAnnot_GetOptionLabel, err := PdfiumInstance.FPDFAnnot_GetOptionLabel(&requests.FPDFAnnot_GetOptionLabel{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetOptionLabel).To(BeNil())
		})
		It("returns an error when FPDFAnnot_IsOptionSelected is called", func() {
			FPDFAnnot_IsOptionSelected, err := PdfiumInstance.FPDFAnnot_IsOptionSelected(&requests.FPDFAnnot_IsOptionSelected{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_IsOptionSelected).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetFontSize is called", func() {
			FPDFAnnot_GetFontSize, err := PdfiumInstance.FPDFAnnot_GetFontSize(&requests.FPDFAnnot_GetFontSize{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetFontSize).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetFontColor is called", func() {
			FPDFAnnot_GetFontColor, err := PdfiumInstance.FPDFAnnot_GetFontColor(&requests.FPDFAnnot_GetFontColor{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetFontColor).To(BeNil())
		})
		It("returns an error when FPDFAnnot_IsChecked is called", func() {
			FPDFAnnot_IsChecked, err := PdfiumInstance.FPDFAnnot_IsChecked(&requests.FPDFAnnot_IsChecked{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_IsChecked).To(BeNil())
		})
		It("returns an error when FPDFAnnot_SetFocusableSubtypes is called", func() {
			FPDFAnnot_SetFocusableSubtypes, err := PdfiumInstance.FPDFAnnot_SetFocusableSubtypes(&requests.FPDFAnnot_SetFocusableSubtypes{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_SetFocusableSubtypes).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetFocusableSubtypesCount is called", func() {
			FPDFAnnot_GetFocusableSubtypesCount, err := PdfiumInstance.FPDFAnnot_GetFocusableSubtypesCount(&requests.FPDFAnnot_GetFocusableSubtypesCount{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetFocusableSubtypesCount).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetFocusableSubtypes is called", func() {
			FPDFAnnot_GetFocusableSubtypes, err := PdfiumInstance.FPDFAnnot_GetFocusableSubtypes(&requests.FPDFAnnot_GetFocusableSubtypes{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetFocusableSubtypes).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetFormControlCount is called", func() {
			FPDFAnnot_GetFormControlCount, err := PdfiumInstance.FPDFAnnot_GetFormControlCount(&requests.FPDFAnnot_GetFormControlCount{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetFormControlCount).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetFormControlIndex is called", func() {
			FPDFAnnot_GetFormControlIndex, err := PdfiumInstance.FPDFAnnot_GetFormControlIndex(&requests.FPDFAnnot_GetFormControlIndex{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetFormControlIndex).To(BeNil())
		})
		It("returns an error when FPDFAnnot_GetFormFieldExportValue is called", func() {
			FPDFAnnot_GetFormFieldExportValue, err := PdfiumInstance.FPDFAnnot_GetFormFieldExportValue(&requests.FPDFAnnot_GetFormFieldExportValue{})
			Expect(err).To(MatchError("formHandle not given"))
			Expect(FPDFAnnot_GetFormFieldExportValue).To(BeNil())
		})
	})

	It("returns the correct supported state per subtype", func() {
		subTypes := map[enums.FPDF_ANNOTATION_SUBTYPE]bool{
			enums.FPDF_ANNOT_SUBTYPE_TEXT:           true,
			enums.FPDF_ANNOT_SUBTYPE_LINK:           true,
			enums.FPDF_ANNOT_SUBTYPE_FREETEXT:       true,
			enums.FPDF_ANNOT_SUBTYPE_LINE:           false,
			enums.FPDF_ANNOT_SUBTYPE_SQUARE:         true,
			enums.FPDF_ANNOT_SUBTYPE_CIRCLE:         true,
			enums.FPDF_ANNOT_SUBTYPE_POLYGON:        false,
			enums.FPDF_ANNOT_SUBTYPE_POLYLINE:       false,
			enums.FPDF_ANNOT_SUBTYPE_HIGHLIGHT:      true,
			enums.FPDF_ANNOT_SUBTYPE_UNDERLINE:      true,
			enums.FPDF_ANNOT_SUBTYPE_SQUIGGLY:       true,
			enums.FPDF_ANNOT_SUBTYPE_STRIKEOUT:      true,
			enums.FPDF_ANNOT_SUBTYPE_STAMP:          true,
			enums.FPDF_ANNOT_SUBTYPE_CARET:          false,
			enums.FPDF_ANNOT_SUBTYPE_INK:            true,
			enums.FPDF_ANNOT_SUBTYPE_POPUP:          true,
			enums.FPDF_ANNOT_SUBTYPE_FILEATTACHMENT: true,
			enums.FPDF_ANNOT_SUBTYPE_SOUND:          false,
			enums.FPDF_ANNOT_SUBTYPE_MOVIE:          false,
			enums.FPDF_ANNOT_SUBTYPE_WIDGET:         false,
			enums.FPDF_ANNOT_SUBTYPE_SCREEN:         false,
			enums.FPDF_ANNOT_SUBTYPE_PRINTERMARK:    false,
			enums.FPDF_ANNOT_SUBTYPE_TRAPNET:        false,
			enums.FPDF_ANNOT_SUBTYPE_WATERMARK:      false,
			enums.FPDF_ANNOT_SUBTYPE_THREED:         false,
			enums.FPDF_ANNOT_SUBTYPE_RICHMEDIA:      false,
			enums.FPDF_ANNOT_SUBTYPE_XFAWIDGET:      false,
			enums.FPDF_ANNOT_SUBTYPE_REDACT:         false,
		}

		for subType := range subTypes {
			By("testing subtype " + strconv.Itoa(int(subType)))
			FPDFAnnot_IsSupportedSubtype, err := PdfiumInstance.FPDFAnnot_IsSupportedSubtype(&requests.FPDFAnnot_IsSupportedSubtype{
				Subtype: subType,
			})
			Expect(err).To(BeNil())
			Expect(FPDFAnnot_IsSupportedSubtype).To(Equal(&responses.FPDFAnnot_IsSupportedSubtype{
				IsSupported: subTypes[subType],
			}))
		}
	})

	It("returns the correct object supported state per subtype", func() {
		subTypes := map[enums.FPDF_ANNOTATION_SUBTYPE]bool{
			enums.FPDF_ANNOT_SUBTYPE_TEXT:           false,
			enums.FPDF_ANNOT_SUBTYPE_LINK:           false,
			enums.FPDF_ANNOT_SUBTYPE_FREETEXT:       false,
			enums.FPDF_ANNOT_SUBTYPE_LINE:           false,
			enums.FPDF_ANNOT_SUBTYPE_SQUARE:         false,
			enums.FPDF_ANNOT_SUBTYPE_CIRCLE:         false,
			enums.FPDF_ANNOT_SUBTYPE_POLYGON:        false,
			enums.FPDF_ANNOT_SUBTYPE_POLYLINE:       false,
			enums.FPDF_ANNOT_SUBTYPE_HIGHLIGHT:      false,
			enums.FPDF_ANNOT_SUBTYPE_UNDERLINE:      false,
			enums.FPDF_ANNOT_SUBTYPE_SQUIGGLY:       false,
			enums.FPDF_ANNOT_SUBTYPE_STRIKEOUT:      false,
			enums.FPDF_ANNOT_SUBTYPE_STAMP:          true,
			enums.FPDF_ANNOT_SUBTYPE_CARET:          false,
			enums.FPDF_ANNOT_SUBTYPE_INK:            true,
			enums.FPDF_ANNOT_SUBTYPE_POPUP:          false,
			enums.FPDF_ANNOT_SUBTYPE_FILEATTACHMENT: false,
			enums.FPDF_ANNOT_SUBTYPE_SOUND:          false,
			enums.FPDF_ANNOT_SUBTYPE_MOVIE:          false,
			enums.FPDF_ANNOT_SUBTYPE_WIDGET:         false,
			enums.FPDF_ANNOT_SUBTYPE_SCREEN:         false,
			enums.FPDF_ANNOT_SUBTYPE_PRINTERMARK:    false,
			enums.FPDF_ANNOT_SUBTYPE_TRAPNET:        false,
			enums.FPDF_ANNOT_SUBTYPE_WATERMARK:      false,
			enums.FPDF_ANNOT_SUBTYPE_THREED:         false,
			enums.FPDF_ANNOT_SUBTYPE_RICHMEDIA:      false,
			enums.FPDF_ANNOT_SUBTYPE_XFAWIDGET:      false,
			enums.FPDF_ANNOT_SUBTYPE_REDACT:         false,
		}

		for subType := range subTypes {
			By("testing subtype " + strconv.Itoa(int(subType)))
			FPDFAnnot_IsObjectSupportedSubtype, err := PdfiumInstance.FPDFAnnot_IsObjectSupportedSubtype(&requests.FPDFAnnot_IsObjectSupportedSubtype{
				Subtype: subType,
			})
			Expect(err).To(BeNil())
			Expect(FPDFAnnot_IsObjectSupportedSubtype).To(Equal(&responses.FPDFAnnot_IsObjectSupportedSubtype{
				IsObjectSupportedSubtype: subTypes[subType],
			}))
		}
	})

	Context("a normal PDF file with annotations", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/annotation_stamp_with_ap.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		It("returns the correct annotations count", func() {
			FPDFPage_GetAnnotCount, err := PdfiumInstance.FPDFPage_GetAnnotCount(&requests.FPDFPage_GetAnnotCount{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(FPDFPage_GetAnnotCount).To(Equal(&responses.FPDFPage_GetAnnotCount{
				Count: 2,
			}))
		})

		It("returns an error when trying to get an annotation that isn't there", func() {
			FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
				Index: 3,
			})
			Expect(err).To(MatchError("could not get annotation"))
			Expect(FPDFPage_GetAnnot).To(BeNil())
		})

		It("returns an error when adding an annotation without a type", func() {
			FPDFPage_CreateAnnot, err := PdfiumInstance.FPDFPage_CreateAnnot(&requests.FPDFPage_CreateAnnot{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
			})
			Expect(err).To(MatchError("could not create annotation"))
			Expect(FPDFPage_CreateAnnot).To(BeNil())
		})

		It("allows to add a new annotation", func() {
			FPDFPage_CreateAnnot, err := PdfiumInstance.FPDFPage_CreateAnnot(&requests.FPDFPage_CreateAnnot{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
				Subtype: enums.FPDF_ANNOT_SUBTYPE_TEXT,
			})
			Expect(err).To(BeNil())
			Expect(FPDFPage_CreateAnnot).To(Not(BeNil()))
			Expect(FPDFPage_CreateAnnot.Annotation).To(Not(BeNil()))
		})

		When("an annotation has been loaded", func() {
			var annotation references.FPDF_ANNOTATION
			BeforeEach(func() {
				FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Index: 0,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnot).To(Not(BeNil()))
				Expect(FPDFPage_GetAnnot.Annotation).To(Not(BeNil()))
				annotation = FPDFPage_GetAnnot.Annotation
			})

			It("returns that it has the specified key", func() {
				FPDFAnnot_HasKey, err := PdfiumInstance.FPDFAnnot_HasKey(&requests.FPDFAnnot_HasKey{
					Annotation: annotation,
					Key:        "AAPL:Hash",
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_HasKey).To(Equal(&responses.FPDFAnnot_HasKey{
					HasKey: true,
				}))
			})

			It("returns the correct value type for the given key", func() {
				FPDFAnnot_GetValueType, err := PdfiumInstance.FPDFAnnot_GetValueType(&requests.FPDFAnnot_GetValueType{
					Annotation: annotation,
					Key:        "AAPL:Hash",
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetValueType).To(Equal(&responses.FPDFAnnot_GetValueType{
					ValueType: enums.FPDF_OBJECT_TYPE_NAME,
				}))
			})

			It("returns the correct value for the given key", func() {
				FPDFAnnot_GetStringValue, err := PdfiumInstance.FPDFAnnot_GetStringValue(&requests.FPDFAnnot_GetStringValue{
					Annotation: annotation,
					Key:        "AAPL:Hash",
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetStringValue).To(Equal(&responses.FPDFAnnot_GetStringValue{
					Value: "395fbcb98d558681742f30683a62a2ad",
				}))
			})

			It("allows setting and getting the value for the given key", func() {
				FPDFAnnot_SetStringValue, err := PdfiumInstance.FPDFAnnot_SetStringValue(&requests.FPDFAnnot_SetStringValue{
					Annotation: annotation,
					Key:        "AAPL:Hash",
					Value:      "Text",
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_SetStringValue).To(Equal(&responses.FPDFAnnot_SetStringValue{}))

				FPDFAnnot_GetStringValue, err := PdfiumInstance.FPDFAnnot_GetStringValue(&requests.FPDFAnnot_GetStringValue{
					Annotation: annotation,
					Key:        "AAPL:Hash",
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetStringValue).To(Equal(&responses.FPDFAnnot_GetStringValue{
					Value: "Text",
				}))
			})
		})

		When("an annotation has been added", func() {
			var annotation references.FPDF_ANNOTATION
			BeforeEach(func() {
				FPDFPage_CreateAnnot, err := PdfiumInstance.FPDFPage_CreateAnnot(&requests.FPDFPage_CreateAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Subtype: enums.FPDF_ANNOT_SUBTYPE_TEXT,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CreateAnnot).To(Not(BeNil()))
				Expect(FPDFPage_CreateAnnot.Annotation).To(Not(BeNil()))
				annotation = FPDFPage_CreateAnnot.Annotation
			})

			AfterEach(func() {
				FPDFPage_CloseAnnot, err := PdfiumInstance.FPDFPage_CloseAnnot(&requests.FPDFPage_CloseAnnot{
					Annotation: annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CloseAnnot).To(Not(BeNil()))
			})

			It("returns the attachment that has just been added", func() {
				FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Index: 2,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnot).To(Not(BeNil()))
				Expect(FPDFPage_GetAnnot.Annotation).To(Not(BeNil()))
			})

			It("returns an error when trying to get the index for an annotation without giving the annotation", func() {
				FPDFPage_GetAnnotIndex, err := PdfiumInstance.FPDFPage_GetAnnotIndex(&requests.FPDFPage_GetAnnotIndex{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(MatchError("annotation not given"))
				Expect(FPDFPage_GetAnnotIndex).To(BeNil())
			})

			It("gives the correct index for the annotation", func() {
				FPDFPage_GetAnnotIndex, err := PdfiumInstance.FPDFPage_GetAnnotIndex(&requests.FPDFPage_GetAnnotIndex{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Annotation: annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnotIndex).To(Equal(&responses.FPDFPage_GetAnnotIndex{
					Index: 2,
				}))
			})

			It("gives the correct subtype for the annotation", func() {
				FPDFAnnot_GetSubtype, err := PdfiumInstance.FPDFAnnot_GetSubtype(&requests.FPDFAnnot_GetSubtype{
					Annotation: annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetSubtype).To(Equal(&responses.FPDFAnnot_GetSubtype{
					Subtype: enums.FPDF_ANNOT_SUBTYPE_TEXT,
				}))
			})

			It("allows the annotation to be deleted", func() {
				FPDFPage_RemoveAnnot, err := PdfiumInstance.FPDFPage_RemoveAnnot(&requests.FPDFPage_RemoveAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Index: 2,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_RemoveAnnot).To(Equal(&responses.FPDFPage_RemoveAnnot{}))
			})

			It("does not allow the annotation to be deleted twice", func() {
				FPDFPage_RemoveAnnot, err := PdfiumInstance.FPDFPage_RemoveAnnot(&requests.FPDFPage_RemoveAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Index: 2,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_RemoveAnnot).To(Equal(&responses.FPDFPage_RemoveAnnot{}))

				FPDFPage_RemoveAnnot, err = PdfiumInstance.FPDFPage_RemoveAnnot(&requests.FPDFPage_RemoveAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Index: 2,
				})
				Expect(err).To(MatchError("could not remove annotation"))
				Expect(FPDFPage_RemoveAnnot).To(BeNil())
			})

			It("returns an error when trying to get update an annotation without giving the object", func() {
				FPDFAnnot_UpdateObject, err := PdfiumInstance.FPDFAnnot_UpdateObject(&requests.FPDFAnnot_UpdateObject{
					Annotation: annotation,
				})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFAnnot_UpdateObject).To(BeNil())
			})

			It("returns an error when trying to append an annotation without giving the object", func() {
				FPDFAnnot_AppendObject, err := PdfiumInstance.FPDFAnnot_AppendObject(&requests.FPDFAnnot_AppendObject{
					Annotation: annotation,
				})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFAnnot_AppendObject).To(BeNil())
			})

			It("allows for stamp annotations to be added and modified", func() {
				By("creating a stamp annotation and set its annotation rectangle")
				FPDFPage_CreateAnnot, err := PdfiumInstance.FPDFPage_CreateAnnot(&requests.FPDFPage_CreateAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Subtype: enums.FPDF_ANNOT_SUBTYPE_STAMP,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CreateAnnot).ToNot(BeNil())
				Expect(FPDFPage_CreateAnnot.Annotation).ToNot(BeEmpty())

				FPDFAnnot_SetRect, err := PdfiumInstance.FPDFAnnot_SetRect(&requests.FPDFAnnot_SetRect{
					Annotation: FPDFPage_CreateAnnot.Annotation,
					Rect: structs.FPDF_FS_RECTF{
						Left:   200,
						Bottom: 600,
						Right:  400,
						Top:    800,
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_SetRect).To(Equal(&responses.FPDFAnnot_SetRect{}))

				FPDFAnnot_GetRect, err := PdfiumInstance.FPDFAnnot_GetRect(&requests.FPDFAnnot_GetRect{
					Annotation: FPDFPage_CreateAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetRect).To(Equal(&responses.FPDFAnnot_GetRect{
					Rect: structs.FPDF_FS_RECTF{
						Left:   200,
						Bottom: 600,
						Right:  400,
						Top:    800,
					},
				}))

				By("adding a solid-color translucent image object to the new annotation")
				FPDFBitmap_Create, err := PdfiumInstance.FPDFBitmap_Create(&requests.FPDFBitmap_Create{
					Width:  200,
					Height: 200,
					Alpha:  1,
				})
				Expect(err).To(BeNil())
				Expect(FPDFBitmap_Create).ToNot(BeNil())
				Expect(FPDFBitmap_Create.Bitmap).ToNot(BeEmpty())

				FPDFBitmap_FillRect, err := PdfiumInstance.FPDFBitmap_FillRect(&requests.FPDFBitmap_FillRect{
					Bitmap: FPDFBitmap_Create.Bitmap,
					Left:   0,
					Top:    0,
					Width:  200,
					Height: 200,
					Color:  0xeeeecccc,
				})
				Expect(err).To(BeNil())
				Expect(FPDFBitmap_FillRect).To(Equal(&responses.FPDFBitmap_FillRect{}))

				FPDFPageObj_NewImageObj, err := PdfiumInstance.FPDFPageObj_NewImageObj(&requests.FPDFPageObj_NewImageObj{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPageObj_NewImageObj).ToNot(BeNil())
				Expect(FPDFPageObj_NewImageObj.PageObject).ToNot(BeEmpty())

				FPDFImageObj_SetBitmap, err := PdfiumInstance.FPDFImageObj_SetBitmap(&requests.FPDFImageObj_SetBitmap{
					Page: &requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Count:       0,
					Bitmap:      FPDFBitmap_Create.Bitmap,
					ImageObject: FPDFPageObj_NewImageObj.PageObject,
				})
				Expect(err).To(BeNil())
				Expect(FPDFImageObj_SetBitmap).To(Equal(&responses.FPDFImageObj_SetBitmap{}))

				FPDFPageObj_SetMatrix, err := PdfiumInstance.FPDFPageObj_SetMatrix(&requests.FPDFPageObj_SetMatrix{
					PageObject: FPDFPageObj_NewImageObj.PageObject,
					Transform:  structs.FPDF_FS_MATRIX{200, 0, 0, 200, 0, 0},
				})
				Expect(err).To(BeNil())
				Expect(FPDFPageObj_SetMatrix).To(Equal(&responses.FPDFPageObj_SetMatrix{}))

				FPDFPageObj_Transform, err := PdfiumInstance.FPDFPageObj_Transform(&requests.FPDFPageObj_Transform{
					PageObject: FPDFPageObj_NewImageObj.PageObject,
					Transform:  structs.FPDF_FS_MATRIX{1, 0, 0, 1, 200, 600},
				})
				Expect(err).To(BeNil())
				Expect(FPDFPageObj_Transform).To(Equal(&responses.FPDFPageObj_Transform{}))

				FPDFAnnot_AppendObject, err := PdfiumInstance.FPDFAnnot_AppendObject(&requests.FPDFAnnot_AppendObject{
					PageObject: FPDFPageObj_NewImageObj.PageObject,
					Annotation: FPDFPage_CreateAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_AppendObject).To(Equal(&responses.FPDFAnnot_AppendObject{}))

				By("retrieving the newly added stamp annotation and its image object")
				FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Index: 3,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnot).To(Not(BeNil()))
				Expect(FPDFPage_GetAnnot.Annotation).To(Not(BeNil()))

				FPDFAnnot_GetObjectCount, err := PdfiumInstance.FPDFAnnot_GetObjectCount(&requests.FPDFAnnot_GetObjectCount{
					Annotation: FPDFPage_CreateAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetObjectCount).To(Equal(&responses.FPDFAnnot_GetObjectCount{
					Count: 1,
				}))

				FPDFAnnot_GetObject, err := PdfiumInstance.FPDFAnnot_GetObject(&requests.FPDFAnnot_GetObject{
					Annotation: FPDFPage_CreateAnnot.Annotation,
					Index:      0,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetObject).To(Not(BeNil()))
				Expect(FPDFAnnot_GetObject.PageObject).To(Not(BeNil()))

				FPDFPageObj_GetType, err := PdfiumInstance.FPDFPageObj_GetType(&requests.FPDFPageObj_GetType{
					PageObject: FPDFAnnot_GetObject.PageObject,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPageObj_GetType).To(Equal(&responses.FPDFPageObj_GetType{
					Type: enums.FPDF_PAGEOBJ_IMAGE,
				}))

				FPDFBitmap_FillRect, err = PdfiumInstance.FPDFBitmap_FillRect(&requests.FPDFBitmap_FillRect{
					Bitmap: FPDFBitmap_Create.Bitmap,
					Left:   0,
					Top:    0,
					Width:  200,
					Height: 200,
					Color:  0xff000000,
				})
				Expect(err).To(BeNil())
				Expect(FPDFBitmap_FillRect).To(Equal(&responses.FPDFBitmap_FillRect{}))

				FPDFImageObj_SetBitmap, err = PdfiumInstance.FPDFImageObj_SetBitmap(&requests.FPDFImageObj_SetBitmap{
					Page: &requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Count:       0,
					Bitmap:      FPDFBitmap_Create.Bitmap,
					ImageObject: FPDFAnnot_GetObject.PageObject,
				})
				Expect(err).To(BeNil())
				Expect(FPDFImageObj_SetBitmap).To(Equal(&responses.FPDFImageObj_SetBitmap{}))

				FPDFAnnot_UpdateObject, err := PdfiumInstance.FPDFAnnot_UpdateObject(&requests.FPDFAnnot_UpdateObject{
					Annotation: FPDFPage_CreateAnnot.Annotation,
					PageObject: FPDFAnnot_GetObject.PageObject,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_UpdateObject).To(Equal(&responses.FPDFAnnot_UpdateObject{}))

				FPDFAnnot_RemoveObject, err := PdfiumInstance.FPDFAnnot_RemoveObject(&requests.FPDFAnnot_RemoveObject{
					Annotation: FPDFPage_CreateAnnot.Annotation,
					Index:      0,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_RemoveObject).To(Equal(&responses.FPDFAnnot_RemoveObject{}))

				FPDFAnnot_RemoveObject, err = PdfiumInstance.FPDFAnnot_RemoveObject(&requests.FPDFAnnot_RemoveObject{
					Annotation: FPDFPage_CreateAnnot.Annotation,
					Index:      0,
				})
				Expect(err).To(MatchError("could not remove object"))
				Expect(FPDFAnnot_RemoveObject).To(BeNil())

				FPDFAnnot_SetColor, err := PdfiumInstance.FPDFAnnot_SetColor(&requests.FPDFAnnot_SetColor{
					Annotation: FPDFPage_CreateAnnot.Annotation,
					ColorType:  enums.FPDFANNOT_COLORTYPE_Color,
					R:          51,
					G:          102,
					B:          153,
					A:          204,
				})
				Expect(err).To(MatchError("could not set annotation color"))
				Expect(FPDFAnnot_SetColor).To(BeNil())

				FPDFAnnot_GetColor, err := PdfiumInstance.FPDFAnnot_GetColor(&requests.FPDFAnnot_GetColor{
					Annotation: FPDFPage_CreateAnnot.Annotation,
				})
				Expect(err).To(MatchError("could not get annotation color"))
				Expect(FPDFAnnot_GetColor).To(BeNil())
			})

			It("gives an error when trying to add an ink stroke to a non-ink annotation", func() {
				FPDFPage_CreateAnnot, err := PdfiumInstance.FPDFPage_CreateAnnot(&requests.FPDFPage_CreateAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Subtype: enums.FPDF_ANNOT_SUBTYPE_HIGHLIGHT,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CreateAnnot).ToNot(BeNil())
				Expect(FPDFPage_CreateAnnot.Annotation).ToNot(BeEmpty())

				FPDFAnnot_AddInkStroke, err := PdfiumInstance.FPDFAnnot_AddInkStroke(&requests.FPDFAnnot_AddInkStroke{
					Annotation: FPDFPage_CreateAnnot.Annotation,
					Points: []structs.FPDF_FS_POINTF{
						{60.0, 90.0},
						{61.0, 91.0},
						{62.0, 92.0},
						{63.0, 93.0},
						{64.0, 94.0},
					},
				})
				Expect(err).To(MatchError("could not add ink stroke"))
				Expect(FPDFAnnot_AddInkStroke).To(BeNil())

				FPDFAnnot_RemoveInkList, err := PdfiumInstance.FPDFAnnot_RemoveInkList(&requests.FPDFAnnot_RemoveInkList{
					Annotation: FPDFPage_CreateAnnot.Annotation,
				})
				Expect(err).To(MatchError("could not remove ink list"))
				Expect(FPDFAnnot_RemoveInkList).To(BeNil())
			})

			It("gives an error when trying to add ink annotations without points", func() {
				FPDFPage_CreateAnnot, err := PdfiumInstance.FPDFPage_CreateAnnot(&requests.FPDFPage_CreateAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Subtype: enums.FPDF_ANNOT_SUBTYPE_INK,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CreateAnnot).ToNot(BeNil())
				Expect(FPDFPage_CreateAnnot.Annotation).ToNot(BeEmpty())

				FPDFAnnot_AddInkStroke, err := PdfiumInstance.FPDFAnnot_AddInkStroke(&requests.FPDFAnnot_AddInkStroke{
					Annotation: FPDFPage_CreateAnnot.Annotation,
					Points:     []structs.FPDF_FS_POINTF{},
				})
				Expect(err).To(MatchError("at least one point is required"))
				Expect(FPDFAnnot_AddInkStroke).To(BeNil())
			})

			It("allows for ink annotations to be added and removed", func() {
				FPDFPage_CreateAnnot, err := PdfiumInstance.FPDFPage_CreateAnnot(&requests.FPDFPage_CreateAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Subtype: enums.FPDF_ANNOT_SUBTYPE_INK,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CreateAnnot).ToNot(BeNil())
				Expect(FPDFPage_CreateAnnot.Annotation).ToNot(BeEmpty())

				FPDFAnnot_AddInkStroke, err := PdfiumInstance.FPDFAnnot_AddInkStroke(&requests.FPDFAnnot_AddInkStroke{
					Annotation: FPDFPage_CreateAnnot.Annotation,
					Points: []structs.FPDF_FS_POINTF{
						{60.0, 90.0},
						{61.0, 91.0},
						{62.0, 92.0},
						{63.0, 93.0},
						{64.0, 94.0},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_AddInkStroke).To(Equal(&responses.FPDFAnnot_AddInkStroke{}))

				FPDFAnnot_GetInkListCount, err := PdfiumInstance.FPDFAnnot_GetInkListCount(&requests.FPDFAnnot_GetInkListCount{
					Annotation: FPDFPage_CreateAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetInkListCount).To(Equal(&responses.FPDFAnnot_GetInkListCount{
					Count: 1,
				}))

				FPDFAnnot_GetInkListPath, err := PdfiumInstance.FPDFAnnot_GetInkListPath(&requests.FPDFAnnot_GetInkListPath{
					Annotation: FPDFPage_CreateAnnot.Annotation,
					Index:      0,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetInkListPath).To(Equal(&responses.FPDFAnnot_GetInkListPath{
					Path: []structs.FPDF_FS_POINTF{
						{60.0, 90.0},
						{61.0, 91.0},
						{62.0, 92.0},
						{63.0, 93.0},
						{64.0, 94.0},
					},
				}))

				FPDFAnnot_RemoveInkList, err := PdfiumInstance.FPDFAnnot_RemoveInkList(&requests.FPDFAnnot_RemoveInkList{
					Annotation: FPDFPage_CreateAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_RemoveInkList).To(Equal(&responses.FPDFAnnot_RemoveInkList{}))
			})

			It("allows for text annotations to be added and modified", func() {
				FPDFPage_CreateAnnot, err := PdfiumInstance.FPDFPage_CreateAnnot(&requests.FPDFPage_CreateAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Subtype: enums.FPDF_ANNOT_SUBTYPE_TEXT,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CreateAnnot).ToNot(BeNil())
				Expect(FPDFPage_CreateAnnot.Annotation).ToNot(BeEmpty())

				FPDFAnnot_SetColor, err := PdfiumInstance.FPDFAnnot_SetColor(&requests.FPDFAnnot_SetColor{
					Annotation: FPDFPage_CreateAnnot.Annotation,
					ColorType:  enums.FPDFANNOT_COLORTYPE_Color,
					R:          51,
					G:          102,
					B:          153,
					A:          204,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_SetColor).To(Equal(&responses.FPDFAnnot_SetColor{}))

				FPDFAnnot_GetColor, err := PdfiumInstance.FPDFAnnot_GetColor(&requests.FPDFAnnot_GetColor{
					Annotation: FPDFPage_CreateAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetColor).To(Equal(&responses.FPDFAnnot_GetColor{R: 51, G: 102, B: 153, A: 204}))
			})

			It("allows for annotation points to be added and retrieved", func() {
				FPDFPage_CreateAnnot, err := PdfiumInstance.FPDFPage_CreateAnnot(&requests.FPDFPage_CreateAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Subtype: enums.FPDF_ANNOT_SUBTYPE_HIGHLIGHT,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CreateAnnot).ToNot(BeNil())
				Expect(FPDFPage_CreateAnnot.Annotation).ToNot(BeEmpty())

				FPDFAnnot_HasAttachmentPoints, err := PdfiumInstance.FPDFAnnot_HasAttachmentPoints(&requests.FPDFAnnot_HasAttachmentPoints{
					Annotation: FPDFPage_CreateAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_HasAttachmentPoints).To(Equal(&responses.FPDFAnnot_HasAttachmentPoints{
					HasAttachmentPoints: true,
				}))

				FPDFAnnot_AppendAttachmentPoints, err := PdfiumInstance.FPDFAnnot_AppendAttachmentPoints(&requests.FPDFAnnot_AppendAttachmentPoints{
					Annotation:       FPDFPage_CreateAnnot.Annotation,
					AttachmentPoints: structs.FPDF_FS_QUADPOINTSF{1, 2, 3, 4, 5, 6, 7, 8},
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_AppendAttachmentPoints).To(Equal(&responses.FPDFAnnot_AppendAttachmentPoints{}))

				FPDFAnnot_CountAttachmentPoints, err := PdfiumInstance.FPDFAnnot_CountAttachmentPoints(&requests.FPDFAnnot_CountAttachmentPoints{
					Annotation: FPDFPage_CreateAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_CountAttachmentPoints).To(Equal(&responses.FPDFAnnot_CountAttachmentPoints{
					Count: 1,
				}))

				FPDFAnnot_SetAttachmentPoints, err := PdfiumInstance.FPDFAnnot_SetAttachmentPoints(&requests.FPDFAnnot_SetAttachmentPoints{
					Annotation:       FPDFPage_CreateAnnot.Annotation,
					AttachmentPoints: structs.FPDF_FS_QUADPOINTSF{1, 2, 3, 4, 5, 6, 7, 8},
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_SetAttachmentPoints).To(Equal(&responses.FPDFAnnot_SetAttachmentPoints{}))

				FPDFAnnot_GetAttachmentPoints, err := PdfiumInstance.FPDFAnnot_GetAttachmentPoints(&requests.FPDFAnnot_GetAttachmentPoints{
					Annotation: FPDFPage_CreateAnnot.Annotation,
					Index:      0,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetAttachmentPoints).To(Equal(&responses.FPDFAnnot_GetAttachmentPoints{
					QuadPoints: structs.FPDF_FS_QUADPOINTSF{1, 2, 3, 4, 5, 6, 7, 8},
				}))

				FPDFAnnot_GetLine, err := PdfiumInstance.FPDFAnnot_GetLine(&requests.FPDFAnnot_GetLine{
					Annotation: FPDFPage_CreateAnnot.Annotation,
				})
				Expect(err).To(MatchError("could not get line"))
				Expect(FPDFAnnot_GetLine).To(BeNil())
			})
		})
	})

	Context("a PDF file with polygon annotations", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/polygon_annot.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		When("an annotation has been loaded", func() {
			var annotation references.FPDF_ANNOTATION
			BeforeEach(func() {
				FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Index: 0,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnot).To(Not(BeNil()))
				Expect(FPDFPage_GetAnnot.Annotation).To(Not(BeNil()))
				annotation = FPDFPage_GetAnnot.Annotation
			})

			AfterEach(func() {
				FPDFPage_CloseAnnot, err := PdfiumInstance.FPDFPage_CloseAnnot(&requests.FPDFPage_CloseAnnot{
					Annotation: annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CloseAnnot).To(Not(BeNil()))
			})

			It("returns the correct vertices", func() {
				FPDFAnnot_GetVertices, err := PdfiumInstance.FPDFAnnot_GetVertices(&requests.FPDFAnnot_GetVertices{
					Annotation: annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetVertices).To(Equal(&responses.FPDFAnnot_GetVertices{
					Vertices: []structs.FPDF_FS_POINTF{
						{X: 159, Y: 296},
						{X: 350, Y: 411},
						{X: 472, Y: 243.4199981689453},
					},
				}))
			})
		})
	})

	Context("a PDF file with line annotations", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/line_annot.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		When("an annotation has been loaded", func() {
			var annotation references.FPDF_ANNOTATION
			BeforeEach(func() {
				FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Index: 0,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnot).To(Not(BeNil()))
				Expect(FPDFPage_GetAnnot.Annotation).To(Not(BeNil()))
				annotation = FPDFPage_GetAnnot.Annotation
			})

			AfterEach(func() {
				FPDFPage_CloseAnnot, err := PdfiumInstance.FPDFPage_CloseAnnot(&requests.FPDFPage_CloseAnnot{
					Annotation: annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CloseAnnot).To(Not(BeNil()))
			})

			It("returns the correct line", func() {
				FPDFAnnot_GetLine, err := PdfiumInstance.FPDFAnnot_GetLine(&requests.FPDFAnnot_GetLine{
					Annotation: annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetLine).To(Equal(&responses.FPDFAnnot_GetLine{
					Start: structs.FPDF_FS_POINTF{X: 159, Y: 296},
					End:   structs.FPDF_FS_POINTF{X: 472, Y: 243.4199981689453},
				}))
			})

			It("returns the correct border", func() {
				FPDFAnnot_GetBorder, err := PdfiumInstance.FPDFAnnot_GetBorder(&requests.FPDFAnnot_GetBorder{
					Annotation: annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetBorder).To(Equal(&responses.FPDFAnnot_GetBorder{
					HorizontalRadius: 0.25,
					VerticalRadius:   0.5,
					BorderWidth:      2,
				}))
			})

			It("allows to set and get the border", func() {
				FPDFAnnot_SetBorder, err := PdfiumInstance.FPDFAnnot_SetBorder(&requests.FPDFAnnot_SetBorder{
					Annotation:       annotation,
					HorizontalRadius: 1,
					VerticalRadius:   2,
					BorderWidth:      3,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_SetBorder).To(Equal(&responses.FPDFAnnot_SetBorder{}))

				FPDFAnnot_GetBorder, err := PdfiumInstance.FPDFAnnot_GetBorder(&requests.FPDFAnnot_GetBorder{
					Annotation: annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetBorder).To(Equal(&responses.FPDFAnnot_GetBorder{
					HorizontalRadius: 1,
					VerticalRadius:   2,
					BorderWidth:      3,
				}))
			})
		})
	})

	Context("a normal PDF file with form annotations", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/text_form_multiple.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		When("an annotation has been loaded", func() {
			var annotation references.FPDF_ANNOTATION
			BeforeEach(func() {
				FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Index: 2,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnot).To(Not(BeNil()))
				Expect(FPDFPage_GetAnnot.Annotation).To(Not(BeNil()))
				annotation = FPDFPage_GetAnnot.Annotation
			})

			It("returns that it has the specified key", func() {
				FPDFAnnot_HasKey, err := PdfiumInstance.FPDFAnnot_HasKey(&requests.FPDFAnnot_HasKey{
					Annotation: annotation,
					Key:        "MaxLen",
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_HasKey).To(Equal(&responses.FPDFAnnot_HasKey{
					HasKey: true,
				}))
			})

			It("returns an error when requesting a key that doesnt exist", func() {
				FPDFAnnot_GetNumberValue, err := PdfiumInstance.FPDFAnnot_GetNumberValue(&requests.FPDFAnnot_GetNumberValue{
					Annotation: annotation,
					Key:        "Nope",
				})
				Expect(err).To(MatchError("could not get number value"))
				Expect(FPDFAnnot_GetNumberValue).To(BeNil())
			})

			It("returns the correct value for the specified key", func() {
				FPDFAnnot_GetNumberValue, err := PdfiumInstance.FPDFAnnot_GetNumberValue(&requests.FPDFAnnot_GetNumberValue{
					Annotation: annotation,
					Key:        "MaxLen",
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetNumberValue).To(Equal(&responses.FPDFAnnot_GetNumberValue{
					Value: 10,
				}))
			})

			It("gives an error when using invalid APs", func() {
				value := "test"
				FPDFAnnot_SetAP, err := PdfiumInstance.FPDFAnnot_SetAP(&requests.FPDFAnnot_SetAP{
					Annotation:     annotation,
					AppearanceMode: 25,
					Value:          &value,
				})
				Expect(err).To(MatchError("could not set appearance mode"))
				Expect(FPDFAnnot_SetAP).To(BeNil())

				FPDFAnnot_GetAP, err := PdfiumInstance.FPDFAnnot_GetAP(&requests.FPDFAnnot_GetAP{
					Annotation:     annotation,
					AppearanceMode: 25,
				})
				Expect(err).To(MatchError("could not get appearance mode"))
				Expect(FPDFAnnot_GetAP).To(BeNil())

				FPDFAnnot_SetAP, err = PdfiumInstance.FPDFAnnot_SetAP(&requests.FPDFAnnot_SetAP{
					Annotation:     annotation,
					AppearanceMode: 25,
				})
				Expect(err).To(MatchError("could not set appearance mode"))
				Expect(FPDFAnnot_SetAP).To(BeNil())
			})

			It("allows setting and getting an AP", func() {
				value := "test"
				FPDFAnnot_SetAP, err := PdfiumInstance.FPDFAnnot_SetAP(&requests.FPDFAnnot_SetAP{
					Annotation:     annotation,
					AppearanceMode: enums.FPDF_ANNOT_APPEARANCEMODE_DOWN,
					Value:          &value,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_SetAP).To(Equal(&responses.FPDFAnnot_SetAP{}))

				FPDFAnnot_GetAP, err := PdfiumInstance.FPDFAnnot_GetAP(&requests.FPDFAnnot_GetAP{
					Annotation:     annotation,
					AppearanceMode: enums.FPDF_ANNOT_APPEARANCEMODE_DOWN,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetAP).To(Equal(&responses.FPDFAnnot_GetAP{
					Value: value,
				}))

				FPDFAnnot_SetAP, err = PdfiumInstance.FPDFAnnot_SetAP(&requests.FPDFAnnot_SetAP{
					Annotation:     annotation,
					AppearanceMode: enums.FPDF_ANNOT_APPEARANCEMODE_DOWN,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_SetAP).To(Equal(&responses.FPDFAnnot_SetAP{}))

				FPDFAnnot_GetAP, err = PdfiumInstance.FPDFAnnot_GetAP(&requests.FPDFAnnot_GetAP{
					Annotation:     annotation,
					AppearanceMode: enums.FPDF_ANNOT_APPEARANCEMODE_DOWN,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetAP).To(Equal(&responses.FPDFAnnot_GetAP{
					Value: "",
				}))
			})
		})
	})

	Context("a PDF file with linked annotations", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/annotation_highlight_square_with_ap.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		When("an annotation has been loaded", func() {
			var annotation references.FPDF_ANNOTATION
			BeforeEach(func() {
				FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Index: 0,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnot).To(Not(BeNil()))
				Expect(FPDFPage_GetAnnot.Annotation).To(Not(BeNil()))
				annotation = FPDFPage_GetAnnot.Annotation
			})

			It("returns an error when getting a linked annotation with unknown key", func() {
				FPDFAnnot_GetLinkedAnnot, err := PdfiumInstance.FPDFAnnot_GetLinkedAnnot(&requests.FPDFAnnot_GetLinkedAnnot{
					Annotation: annotation,
					Key:        "Fake",
				})
				Expect(err).To(MatchError("could not get linked annotation"))
				Expect(FPDFAnnot_GetLinkedAnnot).To(BeNil())
			})

			It("returns the linked annotation", func() {
				FPDFAnnot_GetLinkedAnnot, err := PdfiumInstance.FPDFAnnot_GetLinkedAnnot(&requests.FPDFAnnot_GetLinkedAnnot{
					Annotation: annotation,
					Key:        "Popup",
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetLinkedAnnot).ToNot(BeNil())
				Expect(FPDFAnnot_GetLinkedAnnot.LinkedAnnotation).ToNot(BeEmpty())
			})

			It("returns the correct annotation flags", func() {
				FPDFAnnot_GetFlags, err := PdfiumInstance.FPDFAnnot_GetFlags(&requests.FPDFAnnot_GetFlags{
					Annotation: annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFlags).To(Equal(&responses.FPDFAnnot_GetFlags{
					Flags: 4,
				}))
			})

			It("allows to set and get flags", func() {
				FPDFAnnot_GetFlags, err := PdfiumInstance.FPDFAnnot_GetFlags(&requests.FPDFAnnot_GetFlags{
					Annotation: annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFlags).To(Equal(&responses.FPDFAnnot_GetFlags{
					Flags: 4,
				}))

				FPDFAnnot_SetFlags, err := PdfiumInstance.FPDFAnnot_SetFlags(&requests.FPDFAnnot_SetFlags{
					Annotation: annotation,
					Flags:      8,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_SetFlags).To(Equal(&responses.FPDFAnnot_SetFlags{}))

				FPDFAnnot_GetFlags, err = PdfiumInstance.FPDFAnnot_GetFlags(&requests.FPDFAnnot_GetFlags{
					Annotation: annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFlags).To(Equal(&responses.FPDFAnnot_GetFlags{
					Flags: 8,
				}))
			})
		})
	})

	Context("a PDF file with link annotations", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/annots.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		It("returns an error when getting the link of a non-link annotation", func() {
			FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
				Index: 4,
			})
			Expect(err).To(BeNil())
			Expect(FPDFPage_GetAnnot).To(Not(BeNil()))
			Expect(FPDFPage_GetAnnot.Annotation).To(Not(BeNil()))

			FPDFAnnot_GetLink, err := PdfiumInstance.FPDFAnnot_GetLink(&requests.FPDFAnnot_GetLink{
				Annotation: FPDFPage_GetAnnot.Annotation,
			})
			Expect(err).To(MatchError("could not get link"))
			Expect(FPDFAnnot_GetLink).To(BeNil())
		})

		It("returns the link of an annotation", func() {
			FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
				Index: 3,
			})
			Expect(err).To(BeNil())
			Expect(FPDFPage_GetAnnot).To(Not(BeNil()))
			Expect(FPDFPage_GetAnnot.Annotation).To(Not(BeNil()))

			FPDFAnnot_GetLink, err := PdfiumInstance.FPDFAnnot_GetLink(&requests.FPDFAnnot_GetLink{
				Annotation: FPDFPage_GetAnnot.Annotation,
			})
			Expect(err).To(BeNil())
			Expect(FPDFAnnot_GetLink).ToNot(BeNil())
			Expect(FPDFAnnot_GetLink.Link).ToNot(BeEmpty())
		})

		It("allows setting the link of an annotation", func() {
			FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
				Index: 3,
			})
			Expect(err).To(BeNil())
			Expect(FPDFPage_GetAnnot).To(Not(BeNil()))
			Expect(FPDFPage_GetAnnot.Annotation).To(Not(BeNil()))

			FPDFAnnot_SetURI, err := PdfiumInstance.FPDFAnnot_SetURI(&requests.FPDFAnnot_SetURI{
				Annotation: FPDFPage_GetAnnot.Annotation,
				URI:        "https://github.com/klippa-app/go-pdfium",
			})
			Expect(err).To(BeNil())
			Expect(FPDFAnnot_SetURI).To(Equal(&responses.FPDFAnnot_SetURI{}))
		})

		It("returns an error when setting a uri on a non-link annotation", func() {
			FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
				Index: 4,
			})
			Expect(err).To(BeNil())
			Expect(FPDFPage_GetAnnot).To(Not(BeNil()))
			Expect(FPDFPage_GetAnnot.Annotation).To(Not(BeNil()))

			FPDFAnnot_SetURI, err := PdfiumInstance.FPDFAnnot_SetURI(&requests.FPDFAnnot_SetURI{
				Annotation: FPDFPage_GetAnnot.Annotation,
				URI:        "https://github.com/klippa-app/go-pdfium",
			})
			Expect(err).To(MatchError("could net set uri"))
			Expect(FPDFAnnot_SetURI).To(BeNil())
		})
	})

	Context("a PDF file with alternate form names", func() {
		var doc references.FPDF_DOCUMENT
		var formHandle references.FPDF_FORMHANDLE

		BeforeEach(func() {
			if TestType == "multi" {
				Skip("Form filling is not supported on multi-threaded usage")
			}

			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/click_form.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document

			FPDFDOC_InitFormFillEnvironment, err := PdfiumInstance.FPDFDOC_InitFormFillEnvironment(&requests.FPDFDOC_InitFormFillEnvironment{
				Document: doc,
				FormFillInfo: structs.FPDF_FORMFILLINFO{
					FFI_Invalidate:         func(page references.FPDF_PAGE, left, top, right, bottom float64) {},
					FFI_OutputSelectedRect: func(page references.FPDF_PAGE, left, top, right, bottom float64) {},
					FFI_SetCursor:          func(cursorType enums.FXCT) {},
					FFI_SetTimer: func(elapse int, timerFunc func(idEvent int)) int {
						return 0
					},
					FFI_KillTimer: func(timerID int) {},
					FFI_GetLocalTime: func() structs.FPDF_SYSTEMTIME {
						return structs.FPDF_SYSTEMTIME{}
					},
					FFI_GetPage: func(document references.FPDF_DOCUMENT, index int) *references.FPDF_PAGE {
						return nil
					},
					FFI_GetRotation: func(page references.FPDF_PAGE) enums.FPDF_PAGE_ROTATION {
						return enums.FPDF_PAGE_ROTATION_NONE
					},
					FFI_ExecuteNamedAction: func(namedAction string) {},
				},
			})
			Expect(err).To(BeNil())
			Expect(FPDFDOC_InitFormFillEnvironment).ToNot(BeNil())
			Expect(FPDFDOC_InitFormFillEnvironment.FormHandle).ToNot(BeEmpty())
			formHandle = FPDFDOC_InitFormFillEnvironment.FormHandle
		})

		AfterEach(func() {
			if TestType == "multi" {
				Skip("Form filling is not supported on multi-threaded usage")
			}

			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		It("returns an error when calling FPDFAnnot_GetFormFieldAlternateName without an annotation", func() {
			FPDFAnnot_GetFormFieldAlternateName, err := PdfiumInstance.FPDFAnnot_GetFormFieldAlternateName(&requests.FPDFAnnot_GetFormFieldAlternateName{
				FormHandle: formHandle,
			})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetFormFieldAlternateName).To(BeNil())
		})

		It("returns an error calling FPDFAnnot_GetFormFieldAlternateName without valid annotation", func() {
			FPDFAnnot_GetFormFieldAlternateName, err := PdfiumInstance.FPDFAnnot_GetFormFieldAlternateName(&requests.FPDFAnnot_GetFormFieldAlternateName{
				FormHandle: formHandle,
				Annotation: "test123",
			})
			Expect(err).To(MatchError("could not find annotation handle, perhaps the annotation was already closed or you tried to share annotations between instances"))
			Expect(FPDFAnnot_GetFormFieldAlternateName).To(BeNil())
		})

		When("an annotation has been loaded", func() {
			var annotation references.FPDF_ANNOTATION
			BeforeEach(func() {
				FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Index: 0,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnot).To(Not(BeNil()))
				Expect(FPDFPage_GetAnnot.Annotation).To(Not(BeNil()))
				annotation = FPDFPage_GetAnnot.Annotation
			})

			It("returns an alternate form name for an annotation that has one", func() {
				FPDFAnnot_GetFormFieldAlternateName, err := PdfiumInstance.FPDFAnnot_GetFormFieldAlternateName(&requests.FPDFAnnot_GetFormFieldAlternateName{
					FormHandle: formHandle,
					Annotation: annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFormFieldAlternateName).To(Equal(&responses.FPDFAnnot_GetFormFieldAlternateName{
					FormFieldAlternateName: "readOnlyCheckbox",
				}))
			})
		})
	})

	Context("a PDF file with form additional action JavaScript", func() {
		var doc references.FPDF_DOCUMENT
		var formHandle references.FPDF_FORMHANDLE

		BeforeEach(func() {
			if TestType == "multi" {
				Skip("Form filling is not supported on multi-threaded usage")
			}

			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/annot_javascript.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document

			FPDFDOC_InitFormFillEnvironment, err := PdfiumInstance.FPDFDOC_InitFormFillEnvironment(&requests.FPDFDOC_InitFormFillEnvironment{
				Document: doc,
				FormFillInfo: structs.FPDF_FORMFILLINFO{
					FFI_Invalidate:         func(page references.FPDF_PAGE, left, top, right, bottom float64) {},
					FFI_OutputSelectedRect: func(page references.FPDF_PAGE, left, top, right, bottom float64) {},
					FFI_SetCursor:          func(cursorType enums.FXCT) {},
					FFI_SetTimer: func(elapse int, timerFunc func(idEvent int)) int {
						return 0
					},
					FFI_KillTimer: func(timerID int) {},
					FFI_GetLocalTime: func() structs.FPDF_SYSTEMTIME {
						return structs.FPDF_SYSTEMTIME{}
					},
					FFI_GetPage: func(document references.FPDF_DOCUMENT, index int) *references.FPDF_PAGE {
						return nil
					},
					FFI_GetRotation: func(page references.FPDF_PAGE) enums.FPDF_PAGE_ROTATION {
						return enums.FPDF_PAGE_ROTATION_NONE
					},
					FFI_ExecuteNamedAction: func(namedAction string) {},
				},
			})
			Expect(err).To(BeNil())
			Expect(FPDFDOC_InitFormFillEnvironment).ToNot(BeNil())
			Expect(FPDFDOC_InitFormFillEnvironment.FormHandle).ToNot(BeEmpty())
			formHandle = FPDFDOC_InitFormFillEnvironment.FormHandle
		})

		AfterEach(func() {
			if TestType == "multi" {
				Skip("Form filling is not supported on multi-threaded usage")
			}

			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		It("returns an error when calling FPDFAnnot_GetFormAdditionalActionJavaScript without an annotation", func() {
			FPDFAnnot_GetFormAdditionalActionJavaScript, err := PdfiumInstance.FPDFAnnot_GetFormAdditionalActionJavaScript(&requests.FPDFAnnot_GetFormAdditionalActionJavaScript{
				FormHandle: formHandle,
			})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetFormAdditionalActionJavaScript).To(BeNil())
		})

		It("returns an error calling FPDFAnnot_GetFormAdditionalActionJavaScript without valid annotation", func() {
			FPDFAnnot_GetFormAdditionalActionJavaScript, err := PdfiumInstance.FPDFAnnot_GetFormAdditionalActionJavaScript(&requests.FPDFAnnot_GetFormAdditionalActionJavaScript{
				FormHandle: formHandle,
				Annotation: "test123",
			})
			Expect(err).To(MatchError("could not find annotation handle, perhaps the annotation was already closed or you tried to share annotations between instances"))
			Expect(FPDFAnnot_GetFormAdditionalActionJavaScript).To(BeNil())
		})

		When("an annotation has been loaded", func() {
			var annotation references.FPDF_ANNOTATION
			BeforeEach(func() {
				FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Index: 0,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnot).To(Not(BeNil()))
				Expect(FPDFPage_GetAnnot.Annotation).To(Not(BeNil()))
				annotation = FPDFPage_GetAnnot.Annotation
			})

			It("returns an error when requesting a form additional action JavaScript without an event type", func() {
				FPDFAnnot_GetFormAdditionalActionJavaScript, err := PdfiumInstance.FPDFAnnot_GetFormAdditionalActionJavaScript(&requests.FPDFAnnot_GetFormAdditionalActionJavaScript{
					FormHandle: formHandle,
					Annotation: annotation,
				})
				Expect(err).To(MatchError("could not get form additional action JavaScript"))
				Expect(FPDFAnnot_GetFormAdditionalActionJavaScript).To(BeNil())
			})

			It("returns an empty javascript when requesting a form additional action JavaScript that doesn't have one for the given type", func() {
				FPDFAnnot_GetFormAdditionalActionJavaScript, err := PdfiumInstance.FPDFAnnot_GetFormAdditionalActionJavaScript(&requests.FPDFAnnot_GetFormAdditionalActionJavaScript{
					FormHandle: formHandle,
					Annotation: annotation,
					Event:      enums.FPDF_ANNOT_AACTION_KEY_STROKE,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFormAdditionalActionJavaScript).To(Equal(&responses.FPDFAnnot_GetFormAdditionalActionJavaScript{
					FormAdditionalActionJavaScript: "",
				}))
			})

			It("returns a form additional action JavaScript for an annotation that has one of the given type", func() {
				FPDFAnnot_GetFormAdditionalActionJavaScript, err := PdfiumInstance.FPDFAnnot_GetFormAdditionalActionJavaScript(&requests.FPDFAnnot_GetFormAdditionalActionJavaScript{
					FormHandle: formHandle,
					Annotation: annotation,
					Event:      enums.FPDF_ANNOT_AACTION_FORMAT,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFormAdditionalActionJavaScript).To(Equal(&responses.FPDFAnnot_GetFormAdditionalActionJavaScript{
					FormAdditionalActionJavaScript: "AFDate_FormatEx(\"yyyy-mm-dd\");",
				}))
			})
		})
	})

	Context("a PDF file with a file attached to an annotation", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/annotation_fileattachment.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		It("returns an error when calling FPDFAnnot_GetFileAttachment without a document", func() {
			FPDFAnnot_GetFileAttachment, err := PdfiumInstance.FPDFAnnot_GetFileAttachment(&requests.FPDFAnnot_GetFileAttachment{})
			Expect(err).To(MatchError("document not given"))
			Expect(FPDFAnnot_GetFileAttachment).To(BeNil())
		})

		It("returns an error when calling FPDFAnnot_GetFileAttachment without a valid document", func() {
			FPDFAnnot_GetFileAttachment, err := PdfiumInstance.FPDFAnnot_GetFileAttachment(&requests.FPDFAnnot_GetFileAttachment{
				Document: "test123",
			})
			Expect(err).To(MatchError("could not find document handle, perhaps the doc was already closed or you tried to share documents between instances"))
			Expect(FPDFAnnot_GetFileAttachment).To(BeNil())
		})

		It("returns an error when calling FPDFAnnot_GetFileAttachment without an annotation", func() {
			FPDFAnnot_GetFileAttachment, err := PdfiumInstance.FPDFAnnot_GetFileAttachment(&requests.FPDFAnnot_GetFileAttachment{
				Document: doc,
			})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_GetFileAttachment).To(BeNil())
		})

		It("returns an error when calling FPDFAnnot_GetFileAttachment without a valid annotation", func() {
			FPDFAnnot_GetFileAttachment, err := PdfiumInstance.FPDFAnnot_GetFileAttachment(&requests.FPDFAnnot_GetFileAttachment{
				Document:   doc,
				Annotation: "test123",
			})
			Expect(err).To(MatchError("could not find annotation handle, perhaps the annotation was already closed or you tried to share annotations between instances"))
			Expect(FPDFAnnot_GetFileAttachment).To(BeNil())
		})

		It("returns an error when calling FPDFAnnot_AddFileAttachment without a document", func() {
			FPDFAnnot_AddFileAttachment, err := PdfiumInstance.FPDFAnnot_AddFileAttachment(&requests.FPDFAnnot_AddFileAttachment{})
			Expect(err).To(MatchError("document not given"))
			Expect(FPDFAnnot_AddFileAttachment).To(BeNil())
		})

		It("returns an error when calling FPDFAnnot_AddFileAttachment without a valid document", func() {
			FPDFAnnot_AddFileAttachment, err := PdfiumInstance.FPDFAnnot_AddFileAttachment(&requests.FPDFAnnot_AddFileAttachment{
				Document: "test123",
			})
			Expect(err).To(MatchError("could not find document handle, perhaps the doc was already closed or you tried to share documents between instances"))
			Expect(FPDFAnnot_AddFileAttachment).To(BeNil())
		})

		It("returns an error when calling FPDFAnnot_AddFileAttachment without an annotation", func() {
			FPDFAnnot_AddFileAttachment, err := PdfiumInstance.FPDFAnnot_AddFileAttachment(&requests.FPDFAnnot_AddFileAttachment{
				Document: doc,
			})
			Expect(err).To(MatchError("annotation not given"))
			Expect(FPDFAnnot_AddFileAttachment).To(BeNil())
		})

		It("returns an error when calling FPDFAnnot_AddFileAttachment without a valid annotation", func() {
			FPDFAnnot_AddFileAttachment, err := PdfiumInstance.FPDFAnnot_AddFileAttachment(&requests.FPDFAnnot_AddFileAttachment{
				Document:   doc,
				Annotation: "test123",
			})
			Expect(err).To(MatchError("could not find annotation handle, perhaps the annotation was already closed or you tried to share annotations between instances"))
			Expect(FPDFAnnot_AddFileAttachment).To(BeNil())
		})

		When("an annotation has been loaded", func() {
			var annotation references.FPDF_ANNOTATION
			BeforeEach(func() {
				FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Index: 0,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnot).To(Not(BeNil()))
				Expect(FPDFPage_GetAnnot.Annotation).To(Not(BeNil()))
				annotation = FPDFPage_GetAnnot.Annotation
			})

			It("returns the correct annotation count", func() {
				FPDFPage_GetAnnotCount, err := PdfiumInstance.FPDFPage_GetAnnotCount(&requests.FPDFPage_GetAnnotCount{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnotCount).To(Equal(&responses.FPDFPage_GetAnnotCount{
					Count: 1,
				}))
			})

			It("returns the correct annotation subtype", func() {
				FPDFAnnot_GetSubtype, err := PdfiumInstance.FPDFAnnot_GetSubtype(&requests.FPDFAnnot_GetSubtype{
					Annotation: annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetSubtype).To(Equal(&responses.FPDFAnnot_GetSubtype{
					Subtype: enums.FPDF_ANNOT_SUBTYPE_FILEATTACHMENT,
				}))
			})

			When("an attachment has been loaded", func() {
				var attachment references.FPDF_ATTACHMENT
				BeforeEach(func() {
					FPDFAnnot_GetFileAttachment, err := PdfiumInstance.FPDFAnnot_GetFileAttachment(&requests.FPDFAnnot_GetFileAttachment{
						Document:   doc,
						Annotation: annotation,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAnnot_GetFileAttachment).To(Not(BeNil()))
					Expect(FPDFAnnot_GetFileAttachment.Attachment).To(Not(BeNil()))
					attachment = FPDFAnnot_GetFileAttachment.Attachment
				})

				It("returns the correct filename", func() {
					FPDFAttachment_GetName, err := PdfiumInstance.FPDFAttachment_GetName(&requests.FPDFAttachment_GetName{
						Attachment: attachment,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_GetName).To(Equal(&responses.FPDFAttachment_GetName{
						Name: "test.txt",
					}))
				})

				It("returns the correct file data", func() {
					FPDFAttachment_GetFile, err := PdfiumInstance.FPDFAttachment_GetFile(&requests.FPDFAttachment_GetFile{
						Attachment: attachment,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_GetFile).To(Equal(&responses.FPDFAttachment_GetFile{
						Contents: []byte("test text"),
					}))
				})
			})
		})

		When("an annotation has been added", func() {
			var annotation references.FPDF_ANNOTATION
			BeforeEach(func() {
				FPDFPage_CreateAnnot, err := PdfiumInstance.FPDFPage_CreateAnnot(&requests.FPDFPage_CreateAnnot{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					Subtype: enums.FPDF_ANNOT_SUBTYPE_FILEATTACHMENT,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CreateAnnot).To(Not(BeNil()))
				Expect(FPDFPage_CreateAnnot.Annotation).To(Not(BeNil()))
				annotation = FPDFPage_CreateAnnot.Annotation

				FPDFAnnot_AddFileAttachment, err := PdfiumInstance.FPDFAnnot_AddFileAttachment(&requests.FPDFAnnot_AddFileAttachment{
					Document:   doc,
					Annotation: annotation,
					Name:       "0.txt",
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_AddFileAttachment).To(Not(BeNil()))
				Expect(FPDFAnnot_AddFileAttachment.Attachment).To(Not(BeNil()))
			})

			It("returns the correct annotation count", func() {
				FPDFPage_GetAnnotCount, err := PdfiumInstance.FPDFPage_GetAnnotCount(&requests.FPDFPage_GetAnnotCount{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnotCount).To(Equal(&responses.FPDFPage_GetAnnotCount{
					Count: 2,
				}))
			})

			It("returns the correct annotation subtype", func() {
				FPDFAnnot_GetSubtype, err := PdfiumInstance.FPDFAnnot_GetSubtype(&requests.FPDFAnnot_GetSubtype{
					Annotation: annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetSubtype).To(Equal(&responses.FPDFAnnot_GetSubtype{
					Subtype: enums.FPDF_ANNOT_SUBTYPE_FILEATTACHMENT,
				}))
			})

			When("an attachment has been loaded", func() {
				var attachment references.FPDF_ATTACHMENT
				BeforeEach(func() {
					FPDFAnnot_GetFileAttachment, err := PdfiumInstance.FPDFAnnot_GetFileAttachment(&requests.FPDFAnnot_GetFileAttachment{
						Document:   doc,
						Annotation: annotation,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAnnot_GetFileAttachment).To(Not(BeNil()))
					Expect(FPDFAnnot_GetFileAttachment.Attachment).To(Not(BeNil()))
					attachment = FPDFAnnot_GetFileAttachment.Attachment
				})

				It("returns the correct filename", func() {
					FPDFAttachment_GetName, err := PdfiumInstance.FPDFAttachment_GetName(&requests.FPDFAttachment_GetName{
						Attachment: attachment,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_GetName).To(Equal(&responses.FPDFAttachment_GetName{
						Name: "0.txt",
					}))
				})
			})
		})
	})
})
