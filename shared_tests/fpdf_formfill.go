package shared_tests

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_formfill", func() {
	BeforeEach(func() {
		if TestType == "multi" {
			Skip("Form filling bitmap is not supported on multi-threaded usage")
		}
		Locker.Lock()
	})

	AfterEach(func() {
		if TestType == "multi" {
			Skip("Form filling is not supported on multi-threaded usage")
		}
		Locker.Unlock()
	})

	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFDOC_InitFormFillEnvironment", func() {
				FPDFDOC_InitFormFillEnvironment, err := PdfiumInstance.FPDFDOC_InitFormFillEnvironment(&requests.FPDFDOC_InitFormFillEnvironment{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFDOC_InitFormFillEnvironment).To(BeNil())
			})
			It("returns an error when calling FPDF_LoadXFA", func() {
				FPDF_LoadXFA, err := PdfiumInstance.FPDF_LoadXFA(&requests.FPDF_LoadXFA{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_LoadXFA).To(BeNil())
			})
		})
	})

	Context("no form handle", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFDOC_ExitFormFillEnvironment", func() {
				FPDFDOC_ExitFormFillEnvironment, err := PdfiumInstance.FPDFDOC_ExitFormFillEnvironment(&requests.FPDFDOC_ExitFormFillEnvironment{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FPDFDOC_ExitFormFillEnvironment).To(BeNil())
			})
			It("returns an error when calling FORM_DoDocumentJSAction", func() {
				FORM_DoDocumentJSAction, err := PdfiumInstance.FORM_DoDocumentJSAction(&requests.FORM_DoDocumentJSAction{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_DoDocumentJSAction).To(BeNil())
			})
			It("returns an error when calling FORM_DoDocumentOpenAction", func() {
				FORM_DoDocumentOpenAction, err := PdfiumInstance.FORM_DoDocumentOpenAction(&requests.FORM_DoDocumentOpenAction{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_DoDocumentOpenAction).To(BeNil())
			})
			It("returns an error when calling FORM_DoDocumentAAction", func() {
				FORM_DoDocumentAAction, err := PdfiumInstance.FORM_DoDocumentAAction(&requests.FORM_DoDocumentAAction{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_DoDocumentAAction).To(BeNil())
			})
			It("returns an error when calling FORM_OnMouseMove", func() {
				FORM_OnMouseMove, err := PdfiumInstance.FORM_OnMouseMove(&requests.FORM_OnMouseMove{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_OnMouseMove).To(BeNil())
			})
			It("returns an error when calling FORM_OnFocus", func() {
				FORM_OnFocus, err := PdfiumInstance.FORM_OnFocus(&requests.FORM_OnFocus{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_OnFocus).To(BeNil())
			})
			It("returns an error when calling FORM_OnLButtonDown", func() {
				FORM_OnLButtonDown, err := PdfiumInstance.FORM_OnLButtonDown(&requests.FORM_OnLButtonDown{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_OnLButtonDown).To(BeNil())
			})
			It("returns an error when calling FORM_OnRButtonDown", func() {
				FORM_OnRButtonDown, err := PdfiumInstance.FORM_OnRButtonDown(&requests.FORM_OnRButtonDown{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_OnRButtonDown).To(BeNil())
			})
			It("returns an error when calling FORM_OnLButtonUp", func() {
				FORM_OnLButtonUp, err := PdfiumInstance.FORM_OnLButtonUp(&requests.FORM_OnLButtonUp{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_OnLButtonUp).To(BeNil())
			})
			It("returns an error when calling FORM_OnRButtonUp", func() {
				FORM_OnRButtonUp, err := PdfiumInstance.FORM_OnRButtonUp(&requests.FORM_OnRButtonUp{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_OnRButtonUp).To(BeNil())
			})
			It("returns an error when calling FORM_OnLButtonDoubleClick", func() {
				FORM_OnLButtonDoubleClick, err := PdfiumInstance.FORM_OnLButtonDoubleClick(&requests.FORM_OnLButtonDoubleClick{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_OnLButtonDoubleClick).To(BeNil())
			})
			It("returns an error when calling FORM_OnKeyDown", func() {
				FORM_OnKeyDown, err := PdfiumInstance.FORM_OnKeyDown(&requests.FORM_OnKeyDown{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_OnKeyDown).To(BeNil())
			})
			It("returns an error when calling FORM_OnKeyUp", func() {
				FORM_OnKeyUp, err := PdfiumInstance.FORM_OnKeyUp(&requests.FORM_OnKeyUp{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_OnKeyUp).To(BeNil())
			})
			It("returns an error when calling FORM_OnChar", func() {
				FORM_OnChar, err := PdfiumInstance.FORM_OnChar(&requests.FORM_OnChar{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_OnChar).To(BeNil())
			})
			It("returns an error when calling FORM_GetSelectedText", func() {
				FORM_GetSelectedText, err := PdfiumInstance.FORM_GetSelectedText(&requests.FORM_GetSelectedText{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_GetSelectedText).To(BeNil())
			})
			It("returns an error when calling FORM_ReplaceSelection", func() {
				FORM_ReplaceSelection, err := PdfiumInstance.FORM_ReplaceSelection(&requests.FORM_ReplaceSelection{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_ReplaceSelection).To(BeNil())
			})
			It("returns an error when calling FORM_CanUndo", func() {
				FORM_CanUndo, err := PdfiumInstance.FORM_CanUndo(&requests.FORM_CanUndo{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_CanUndo).To(BeNil())
			})
			It("returns an error when calling FORM_CanRedo", func() {
				FORM_CanRedo, err := PdfiumInstance.FORM_CanRedo(&requests.FORM_CanRedo{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_CanRedo).To(BeNil())
			})
			It("returns an error when calling FORM_Undo", func() {
				FORM_Undo, err := PdfiumInstance.FORM_Undo(&requests.FORM_Undo{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_Undo).To(BeNil())
			})
			It("returns an error when calling FORM_Redo", func() {
				FORM_Redo, err := PdfiumInstance.FORM_Redo(&requests.FORM_Redo{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_Redo).To(BeNil())
			})
			It("returns an error when calling FORM_ForceToKillFocus", func() {
				FORM_ForceToKillFocus, err := PdfiumInstance.FORM_ForceToKillFocus(&requests.FORM_ForceToKillFocus{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_ForceToKillFocus).To(BeNil())
			})
			It("returns an error when calling FPDFPage_HasFormFieldAtPoint", func() {
				FPDFPage_HasFormFieldAtPoint, err := PdfiumInstance.FPDFPage_HasFormFieldAtPoint(&requests.FPDFPage_HasFormFieldAtPoint{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FPDFPage_HasFormFieldAtPoint).To(BeNil())
			})
			It("returns an error when calling FPDFPage_FormFieldZOrderAtPoint", func() {
				FPDFPage_FormFieldZOrderAtPoint, err := PdfiumInstance.FPDFPage_FormFieldZOrderAtPoint(&requests.FPDFPage_FormFieldZOrderAtPoint{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FPDFPage_FormFieldZOrderAtPoint).To(BeNil())
			})
			It("returns an error when calling FPDF_SetFormFieldHighlightColor", func() {
				FPDF_SetFormFieldHighlightColor, err := PdfiumInstance.FPDF_SetFormFieldHighlightColor(&requests.FPDF_SetFormFieldHighlightColor{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FPDF_SetFormFieldHighlightColor).To(BeNil())
			})
			It("returns an error when calling FPDF_SetFormFieldHighlightAlpha", func() {
				FPDF_SetFormFieldHighlightAlpha, err := PdfiumInstance.FPDF_SetFormFieldHighlightAlpha(&requests.FPDF_SetFormFieldHighlightAlpha{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FPDF_SetFormFieldHighlightAlpha).To(BeNil())
			})
			It("returns an error when calling FPDF_RemoveFormFieldHighlight", func() {
				FPDF_RemoveFormFieldHighlight, err := PdfiumInstance.FPDF_RemoveFormFieldHighlight(&requests.FPDF_RemoveFormFieldHighlight{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FPDF_RemoveFormFieldHighlight).To(BeNil())
			})
			It("returns an error when calling FPDF_FFLDraw", func() {
				FPDF_FFLDraw, err := PdfiumInstance.FPDF_FFLDraw(&requests.FPDF_FFLDraw{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FPDF_FFLDraw).To(BeNil())
			})
		})
	})

	Context("no page", func() {
		When("is opened", func() {
			It("returns an error when calling FORM_OnAfterLoadPage", func() {
				FORM_OnAfterLoadPage, err := PdfiumInstance.FORM_OnAfterLoadPage(&requests.FORM_OnAfterLoadPage{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FORM_OnAfterLoadPage).To(BeNil())
			})
			It("returns an error when calling FORM_OnBeforeClosePage", func() {
				FORM_OnBeforeClosePage, err := PdfiumInstance.FORM_OnBeforeClosePage(&requests.FORM_OnBeforeClosePage{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FORM_OnBeforeClosePage).To(BeNil())
			})
			It("returns an error when calling FORM_DoPageAAction", func() {
				FORM_DoPageAAction, err := PdfiumInstance.FORM_DoPageAAction(&requests.FORM_DoPageAAction{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FORM_DoPageAAction).To(BeNil())
			})
		})
	})

	type FormHistory struct {
		Name string
		Args []interface{}
	}

	type FormTicker struct {
		Timer *time.Ticker
		Done  chan bool
	}

	Context("a normal PDF file", func() {
		var doc references.FPDF_DOCUMENT
		var formHandle references.FPDF_FORMHANDLE
		formHistory := []FormHistory{}
		timers := map[int]*FormTicker{}

		BeforeEach(func() {
			formHistory = []FormHistory{}
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/text_form.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document

			FPDFDOC_InitFormFillEnvironment, err := PdfiumInstance.FPDFDOC_InitFormFillEnvironment(&requests.FPDFDOC_InitFormFillEnvironment{
				Document: doc,
				FormFillInfo: structs.FPDF_FORMFILLINFO{
					Release: func() {
						formHistory = append(formHistory, FormHistory{
							Name: "Release",
						})
					},
					FFI_Invalidate: func(page references.FPDF_PAGE, left, top, right, bottom float64) {
						formHistory = append(formHistory, FormHistory{
							Name: "FFI_Invalidate",
							Args: []interface{}{page, left, top, right, bottom},
						})
					},
					FFI_OutputSelectedRect: func(page references.FPDF_PAGE, left, top, right, bottom float64) {
						formHistory = append(formHistory, FormHistory{
							Name: "FFI_OutputSelectedRect",
							Args: []interface{}{page, left, top, right, bottom},
						})
					},
					FFI_SetCursor: func(cursorType enums.FXCT) {
						formHistory = append(formHistory, FormHistory{
							Name: "FFI_SetCursor",
							Args: []interface{}{cursorType},
						})
					},
					FFI_SetTimer: func(elapse int, timerFunc func(idEvent int)) int {
						formHistory = append(formHistory, FormHistory{
							Name: "FFI_SetTimer",
							Args: []interface{}{elapse},
						})

						ticker := time.NewTicker(time.Duration(elapse) * time.Millisecond)
						formTicker := &FormTicker{
							Timer: ticker,
							Done:  make(chan bool),
						}

						id := len(timers) + 1 // ID can't be 0
						timers[id] = formTicker

						go func() {
							for {
								select {
								case <-formTicker.Done:
									return
								case <-ticker.C:
									timerFunc(id)
								}
							}
						}()

						return id
					},
					FFI_KillTimer: func(timerID int) {
						formHistory = append(formHistory, FormHistory{
							Name: "FFI_KillTimer",
							Args: []interface{}{timerID},
						})

						_, ok := timers[timerID]
						if !ok {
							return
						}

						timers[timerID].Timer.Stop()
						timers[timerID].Done <- true
					},
					FFI_GetLocalTime: func() structs.FPDF_SYSTEMTIME {
						formHistory = append(formHistory, FormHistory{
							Name: "FFI_GetLocalTime",
						})
						return structs.FPDF_SYSTEMTIME{}
					},
					FFI_OnChange: func() {
						formHistory = append(formHistory, FormHistory{
							Name: "FFI_OnChange",
						})
					},
					FFI_GetPage: func(document references.FPDF_DOCUMENT, index int) *references.FPDF_PAGE {
						formHistory = append(formHistory, FormHistory{
							Name: "FFI_GetPage",
							Args: []interface{}{document, index},
						})

						return nil
					},
					FFI_GetCurrentPage: func(document references.FPDF_DOCUMENT) *references.FPDF_PAGE {
						formHistory = append(formHistory, FormHistory{
							Name: "FFI_GetCurrentPage",
							Args: []interface{}{document},
						})
						return nil
					},
					FFI_GetRotation: func(page references.FPDF_PAGE) enums.FPDF_PAGE_ROTATION {
						formHistory = append(formHistory, FormHistory{
							Name: "FFI_GetRotation",
							Args: []interface{}{page},
						})
						return enums.FPDF_PAGE_ROTATION_NONE
					},
					FFI_ExecuteNamedAction: func(namedAction string) {
						formHistory = append(formHistory, FormHistory{
							Name: "FFI_ExecuteNamedAction",
							Args: []interface{}{namedAction},
						})
					},
					FFI_SetTextFieldFocus: func(value string, isFocus bool) {
						formHistory = append(formHistory, FormHistory{
							Name: "FFI_SetTextFieldFocus",
							Args: []interface{}{value, isFocus},
						})
					},
					FFI_DoURIAction: func(bsURI string) {
						formHistory = append(formHistory, FormHistory{
							Name: "FFI_DoURIAction",
							Args: []interface{}{bsURI},
						})
					},
					FFI_DoGoToAction: func(pageIndex int, zoomMode enums.FPDF_ZOOM_MODE, pos []float32) {
						formHistory = append(formHistory, FormHistory{
							Name: "FFI_DoGoToAction",
							Args: []interface{}{pageIndex, zoomMode, pos},
						})
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(FPDFDOC_InitFormFillEnvironment).ToNot(BeNil())
			Expect(FPDFDOC_InitFormFillEnvironment.FormHandle).ToNot(BeEmpty())
			formHandle = FPDFDOC_InitFormFillEnvironment.FormHandle
		})

		AfterEach(func() {
			FPDFDOC_ExitFormFillEnvironment, err := PdfiumInstance.FPDFDOC_ExitFormFillEnvironment(&requests.FPDFDOC_ExitFormFillEnvironment{
				FormHandle: formHandle,
			})
			Expect(err).To(BeNil())
			Expect(FPDFDOC_ExitFormFillEnvironment).To(Not(BeNil()))

			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))

			log.Println(formHistory)
		})

		When("required callbacks are missing", func() {
			It("returns an error when calling FPDFDOC_InitFormFillEnvironment", func() {
				FPDFDOC_InitFormFillEnvironment, err := PdfiumInstance.FPDFDOC_InitFormFillEnvironment(&requests.FPDFDOC_InitFormFillEnvironment{
					Document: doc,
				})
				Expect(err).To(MatchError("FormFillInfo callback FFI_Invalidate is required"))
				Expect(FPDFDOC_InitFormFillEnvironment).To(BeNil())
			})

			It("returns an error when calling FPDFDOC_InitFormFillEnvironment", func() {
				FPDFDOC_InitFormFillEnvironment, err := PdfiumInstance.FPDFDOC_InitFormFillEnvironment(&requests.FPDFDOC_InitFormFillEnvironment{
					Document: doc,
					FormFillInfo: structs.FPDF_FORMFILLINFO{
						FFI_Invalidate: func(page references.FPDF_PAGE, left, top, right, bottom float64) {},
					},
				})
				Expect(err).To(MatchError("FormFillInfo callback FFI_SetCursor is required"))
				Expect(FPDFDOC_InitFormFillEnvironment).To(BeNil())
			})

			It("returns an error when calling FPDFDOC_InitFormFillEnvironment", func() {
				FPDFDOC_InitFormFillEnvironment, err := PdfiumInstance.FPDFDOC_InitFormFillEnvironment(&requests.FPDFDOC_InitFormFillEnvironment{
					Document: doc,
					FormFillInfo: structs.FPDF_FORMFILLINFO{
						FFI_Invalidate: func(page references.FPDF_PAGE, left, top, right, bottom float64) {},
						FFI_SetCursor:  func(cursorType enums.FXCT) {},
					},
				})
				Expect(err).To(MatchError("FormFillInfo callback FFI_SetTimer is required"))
				Expect(FPDFDOC_InitFormFillEnvironment).To(BeNil())
			})

			It("returns an error when calling FPDFDOC_InitFormFillEnvironment", func() {
				FPDFDOC_InitFormFillEnvironment, err := PdfiumInstance.FPDFDOC_InitFormFillEnvironment(&requests.FPDFDOC_InitFormFillEnvironment{
					Document: doc,
					FormFillInfo: structs.FPDF_FORMFILLINFO{
						FFI_Invalidate:         func(page references.FPDF_PAGE, left, top, right, bottom float64) {},
						FFI_OutputSelectedRect: func(page references.FPDF_PAGE, left, top, right, bottom float64) {},
						FFI_SetCursor:          func(cursorType enums.FXCT) {},
						FFI_SetTimer: func(elapse int, timerFunc func(idEvent int)) int {
							return 0
						},
					},
				})
				Expect(err).To(MatchError("FormFillInfo callback FFI_KillTimer is required"))
				Expect(FPDFDOC_InitFormFillEnvironment).To(BeNil())
			})

			It("returns an error when calling FPDFDOC_InitFormFillEnvironment", func() {
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
					},
				})
				Expect(err).To(MatchError("FormFillInfo callback FFI_GetLocalTime is required"))
				Expect(FPDFDOC_InitFormFillEnvironment).To(BeNil())
			})

			It("returns an error when calling FPDFDOC_InitFormFillEnvironment", func() {
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
					},
				})
				Expect(err).To(MatchError("FormFillInfo callback FFI_GetPage is required"))
				Expect(FPDFDOC_InitFormFillEnvironment).To(BeNil())
			})

			It("returns an error when calling FPDFDOC_InitFormFillEnvironment", func() {
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
					},
				})
				Expect(err).To(MatchError("FormFillInfo callback FFI_GetRotation is required"))
				Expect(FPDFDOC_InitFormFillEnvironment).To(BeNil())
			})

			It("returns an error when calling FPDFDOC_InitFormFillEnvironment", func() {
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
					},
				})
				Expect(err).To(MatchError("FormFillInfo callback FFI_ExecuteNamedAction is required"))
				Expect(FPDFDOC_InitFormFillEnvironment).To(BeNil())
			})
		})

		When("is opened", func() {
			When("no form handle is given", func() {
				It("returns an error when calling FORM_OnAfterLoadPage", func() {
					FORM_OnAfterLoadPage, err := PdfiumInstance.FORM_OnAfterLoadPage(&requests.FORM_OnAfterLoadPage{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})
					Expect(err).To(MatchError("formHandle not given"))
					Expect(FORM_OnAfterLoadPage).To(BeNil())
				})
				It("returns an error when calling FORM_OnBeforeClosePage", func() {
					FORM_OnBeforeClosePage, err := PdfiumInstance.FORM_OnBeforeClosePage(&requests.FORM_OnBeforeClosePage{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})
					Expect(err).To(MatchError("formHandle not given"))
					Expect(FORM_OnBeforeClosePage).To(BeNil())
				})

				It("returns an error when calling FORM_DoPageAAction", func() {
					FORM_DoPageAAction, err := PdfiumInstance.FORM_DoPageAAction(&requests.FORM_DoPageAAction{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})
					Expect(err).To(MatchError("formHandle not given"))
					Expect(FORM_DoPageAAction).To(BeNil())
				})
			})

			When("no page is given", func() {
				It("returns an error when calling FORM_OnMouseMove", func() {
					FORM_OnMouseMove, err := PdfiumInstance.FORM_OnMouseMove(&requests.FORM_OnMouseMove{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_OnMouseMove).To(BeNil())
				})

				It("returns an error when calling FORM_OnFocus", func() {
					FORM_OnFocus, err := PdfiumInstance.FORM_OnFocus(&requests.FORM_OnFocus{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_OnFocus).To(BeNil())
				})

				It("returns an error when calling FORM_OnLButtonDown", func() {
					FORM_OnLButtonDown, err := PdfiumInstance.FORM_OnLButtonDown(&requests.FORM_OnLButtonDown{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_OnLButtonDown).To(BeNil())
				})

				It("returns an error when calling FORM_OnRButtonDown", func() {
					FORM_OnRButtonDown, err := PdfiumInstance.FORM_OnRButtonDown(&requests.FORM_OnRButtonDown{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_OnRButtonDown).To(BeNil())
				})

				It("returns an error when calling FORM_OnLButtonUp", func() {
					FORM_OnLButtonUp, err := PdfiumInstance.FORM_OnLButtonUp(&requests.FORM_OnLButtonUp{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_OnLButtonUp).To(BeNil())
				})
				It("returns an error when calling FORM_OnRButtonUp", func() {
					FORM_OnRButtonUp, err := PdfiumInstance.FORM_OnRButtonUp(&requests.FORM_OnRButtonUp{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_OnRButtonUp).To(BeNil())
				})
				It("returns an error when calling FORM_OnLButtonDoubleClick", func() {
					FORM_OnLButtonDoubleClick, err := PdfiumInstance.FORM_OnLButtonDoubleClick(&requests.FORM_OnLButtonDoubleClick{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_OnLButtonDoubleClick).To(BeNil())
				})
				It("returns an error when calling FORM_OnKeyDown", func() {
					FORM_OnKeyDown, err := PdfiumInstance.FORM_OnKeyDown(&requests.FORM_OnKeyDown{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_OnKeyDown).To(BeNil())
				})
				It("returns an error when calling FORM_OnKeyUp", func() {
					FORM_OnKeyUp, err := PdfiumInstance.FORM_OnKeyUp(&requests.FORM_OnKeyUp{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_OnKeyUp).To(BeNil())
				})
				It("returns an error when calling FORM_OnChar", func() {
					FORM_OnChar, err := PdfiumInstance.FORM_OnChar(&requests.FORM_OnChar{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_OnChar).To(BeNil())
				})
				It("returns an error when calling FORM_GetSelectedText", func() {
					FORM_GetSelectedText, err := PdfiumInstance.FORM_GetSelectedText(&requests.FORM_GetSelectedText{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_GetSelectedText).To(BeNil())
				})
				It("returns an error when calling FORM_ReplaceSelection", func() {
					FORM_ReplaceSelection, err := PdfiumInstance.FORM_ReplaceSelection(&requests.FORM_ReplaceSelection{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_ReplaceSelection).To(BeNil())
				})
				It("returns an error when calling FORM_CanUndo", func() {
					FORM_CanUndo, err := PdfiumInstance.FORM_CanUndo(&requests.FORM_CanUndo{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_CanUndo).To(BeNil())
				})
				It("returns an error when calling FORM_CanRedo", func() {
					FORM_CanRedo, err := PdfiumInstance.FORM_CanRedo(&requests.FORM_CanRedo{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_CanRedo).To(BeNil())
				})
				It("returns an error when calling FORM_Undo", func() {
					FORM_Undo, err := PdfiumInstance.FORM_Undo(&requests.FORM_Undo{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_Undo).To(BeNil())
				})
				It("returns an error when calling FORM_Redo", func() {
					FORM_Redo, err := PdfiumInstance.FORM_Redo(&requests.FORM_Redo{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_Redo).To(BeNil())
				})
				It("returns an error when calling FPDFPage_HasFormFieldAtPoint", func() {
					FPDFPage_HasFormFieldAtPoint, err := PdfiumInstance.FPDFPage_HasFormFieldAtPoint(&requests.FPDFPage_HasFormFieldAtPoint{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FPDFPage_HasFormFieldAtPoint).To(BeNil())
				})
				It("returns an error when calling FPDFPage_FormFieldZOrderAtPoint", func() {
					FPDFPage_FormFieldZOrderAtPoint, err := PdfiumInstance.FPDFPage_FormFieldZOrderAtPoint(&requests.FPDFPage_FormFieldZOrderAtPoint{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FPDFPage_FormFieldZOrderAtPoint).To(BeNil())
				})
				It("returns an error when calling FPDF_FFLDraw", func() {
					FPDFBitmap_Create, err := PdfiumInstance.FPDFBitmap_Create(&requests.FPDFBitmap_Create{
						Width:  100,
						Height: 100,
						Alpha:  1,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_Create).ToNot(BeNil())
					Expect(FPDFBitmap_Create.Bitmap).ToNot(BeEmpty())

					FPDF_FFLDraw, err := PdfiumInstance.FPDF_FFLDraw(&requests.FPDF_FFLDraw{
						FormHandle: formHandle,
						Bitmap:     FPDFBitmap_Create.Bitmap,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FPDF_FFLDraw).To(BeNil())
				})
			})
			When("no bitmap is given", func() {
				It("returns an error when calling FPDF_FFLDraw", func() {
					FPDF_FFLDraw, err := PdfiumInstance.FPDF_FFLDraw(&requests.FPDF_FFLDraw{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("bitmap not given"))
					Expect(FPDF_FFLDraw).To(BeNil())
				})
			})

			It("allows loading/closing pages", func() {
				FPDF_LoadPage, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_LoadPage).ToNot(BeNil())
				Expect(FPDF_LoadPage.Page).ToNot(BeEmpty())

				FORM_OnAfterLoadPage, err := PdfiumInstance.FORM_OnAfterLoadPage(&requests.FORM_OnAfterLoadPage{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
				})
				Expect(err).To(BeNil())
				Expect(FORM_OnAfterLoadPage).To(Equal(&responses.FORM_OnAfterLoadPage{}))

				FORM_OnBeforeClosePage, err := PdfiumInstance.FORM_OnBeforeClosePage(&requests.FORM_OnBeforeClosePage{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
				})
				Expect(err).To(BeNil())
				Expect(FORM_OnBeforeClosePage).To(Equal(&responses.FORM_OnBeforeClosePage{}))

				FPDF_ClosePage, err := PdfiumInstance.FPDF_ClosePage(&requests.FPDF_ClosePage{
					Page: FPDF_LoadPage.Page,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_ClosePage).ToNot(BeNil())
				Expect(FPDF_ClosePage).To(Equal(&responses.FPDF_ClosePage{}))
			})

			It("allows calling actions", func() {
				FPDF_LoadPage, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_LoadPage).ToNot(BeNil())
				Expect(FPDF_LoadPage.Page).ToNot(BeEmpty())

				FORM_DoDocumentJSAction, err := PdfiumInstance.FORM_DoDocumentJSAction(&requests.FORM_DoDocumentJSAction{
					FormHandle: formHandle,
				})
				Expect(err).To(BeNil())
				Expect(FORM_DoDocumentJSAction).To(Equal(&responses.FORM_DoDocumentJSAction{}))

				FORM_DoDocumentOpenAction, err := PdfiumInstance.FORM_DoDocumentOpenAction(&requests.FORM_DoDocumentOpenAction{
					FormHandle: formHandle,
				})
				Expect(err).To(BeNil())
				Expect(FORM_DoDocumentOpenAction).To(Equal(&responses.FORM_DoDocumentOpenAction{}))

				FORM_DoDocumentAAction, err := PdfiumInstance.FORM_DoDocumentAAction(&requests.FORM_DoDocumentAAction{
					FormHandle: formHandle,
					AAType:     enums.FPDFDOC_AACTION_WC,
				})
				Expect(err).To(BeNil())
				Expect(FORM_DoDocumentAAction).To(Equal(&responses.FORM_DoDocumentAAction{}))

				FORM_DoPageAAction, err := PdfiumInstance.FORM_DoPageAAction(&requests.FORM_DoPageAAction{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					AAType:     enums.FPDFPAGE_AACTION_OPEN,
				})
				Expect(err).To(BeNil())
				Expect(FORM_DoPageAAction).To(Equal(&responses.FORM_DoPageAAction{}))
			})

			It("allows the mouse to be moved", func() {
				FPDF_LoadPage, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_LoadPage).ToNot(BeNil())
				Expect(FPDF_LoadPage.Page).ToNot(BeEmpty())

				FORM_OnMouseMove, err := PdfiumInstance.FORM_OnMouseMove(&requests.FORM_OnMouseMove{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					PageX:      120,
					PageY:      120,
					Modifier:   0,
				})
				Expect(err).To(BeNil())
				Expect(FORM_OnMouseMove).To(Equal(&responses.FORM_OnMouseMove{}))
			})

			It("allows to detect a form field on a position", func() {
				FPDF_LoadPage, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_LoadPage).ToNot(BeNil())
				Expect(FPDF_LoadPage.Page).ToNot(BeEmpty())

				FPDFPage_HasFormFieldAtPoint, err := PdfiumInstance.FPDFPage_HasFormFieldAtPoint(&requests.FPDFPage_HasFormFieldAtPoint{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					PageX:      120,
					PageY:      120,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_HasFormFieldAtPoint).To(Equal(&responses.FPDFPage_HasFormFieldAtPoint{
					FieldType: enums.FPDF_FORMFIELD_TEXTFIELD,
				}))
			})

			It("allows to get the z-order a form field on a position", func() {
				FPDF_LoadPage, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_LoadPage).ToNot(BeNil())
				Expect(FPDF_LoadPage.Page).ToNot(BeEmpty())

				FPDFPage_FormFieldZOrderAtPoint, err := PdfiumInstance.FPDFPage_FormFieldZOrderAtPoint(&requests.FPDFPage_FormFieldZOrderAtPoint{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					PageX:      120,
					PageY:      120,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_FormFieldZOrderAtPoint).To(Equal(&responses.FPDFPage_FormFieldZOrderAtPoint{
					ZOrder: 0,
				}))
			})

			It("allows to focus on a text field and type something", func() {
				FPDF_LoadPage, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_LoadPage).ToNot(BeNil())
				Expect(FPDF_LoadPage.Page).ToNot(BeEmpty())

				FORM_OnMouseMove, err := PdfiumInstance.FORM_OnMouseMove(&requests.FORM_OnMouseMove{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					PageX:      120,
					PageY:      120,
					Modifier:   0,
				})
				Expect(err).To(BeNil())
				Expect(FORM_OnMouseMove).To(Equal(&responses.FORM_OnMouseMove{}))

				FORM_OnLButtonDown, err := PdfiumInstance.FORM_OnLButtonDown(&requests.FORM_OnLButtonDown{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					PageX:      120,
					PageY:      120,
					Modifier:   0,
				})
				Expect(err).To(BeNil())
				Expect(FORM_OnLButtonDown).To(Equal(&responses.FORM_OnLButtonDown{}))

				FORM_OnLButtonUp, err := PdfiumInstance.FORM_OnLButtonUp(&requests.FORM_OnLButtonUp{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					PageX:      120,
					PageY:      120,
					Modifier:   0,
				})
				Expect(err).To(BeNil())
				Expect(FORM_OnLButtonUp).To(Equal(&responses.FORM_OnLButtonUp{}))

				FORM_OnChar, err := PdfiumInstance.FORM_OnChar(&requests.FORM_OnChar{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					NChar:      'A',
					Modifier:   0,
				})
				Expect(err).To(BeNil())
				Expect(FORM_OnChar).To(Equal(&responses.FORM_OnChar{}))

				FORM_OnChar, err = PdfiumInstance.FORM_OnChar(&requests.FORM_OnChar{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					NChar:      'B',
					Modifier:   0,
				})
				Expect(err).To(BeNil())
				Expect(FORM_OnChar).To(Equal(&responses.FORM_OnChar{}))

				FORM_OnChar, err = PdfiumInstance.FORM_OnChar(&requests.FORM_OnChar{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					NChar:      'C',
					Modifier:   0,
				})
				Expect(err).To(BeNil())
				Expect(FORM_OnChar).To(Equal(&responses.FORM_OnChar{}))
			})
		})
	})
})
