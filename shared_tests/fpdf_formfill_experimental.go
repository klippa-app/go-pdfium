//go:build pdfium_experimental
// +build pdfium_experimental

package shared_tests

import "C"
import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"time"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_formfill_experimental", func() {
	BeforeEach(func() {
		if TestType == "multi" {
			Skip("Form filling is not supported on multi-threaded usage")
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
			It("returns an error when calling FPDF_GetFormType", func() {
				FPDF_GetFormType, err := PdfiumInstance.FPDF_GetFormType(&requests.FPDF_GetFormType{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_GetFormType).To(BeNil())
			})
		})
	})

	Context("no form handle", func() {
		When("is opened", func() {
			It("returns an error when calling FORM_OnMouseWheel", func() {
				FORM_OnMouseWheel, err := PdfiumInstance.FORM_OnMouseWheel(&requests.FORM_OnMouseWheel{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_OnMouseWheel).To(BeNil())
			})
			It("returns an error when calling FORM_GetFocusedText", func() {
				FORM_GetFocusedText, err := PdfiumInstance.FORM_GetFocusedText(&requests.FORM_GetFocusedText{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_GetFocusedText).To(BeNil())
			})
			It("returns an error when calling FORM_SelectAllText", func() {
				FORM_SelectAllText, err := PdfiumInstance.FORM_SelectAllText(&requests.FORM_SelectAllText{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_SelectAllText).To(BeNil())
			})
			It("returns an error when calling FORM_GetFocusedAnnot", func() {
				FORM_GetFocusedAnnot, err := PdfiumInstance.FORM_GetFocusedAnnot(&requests.FORM_GetFocusedAnnot{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_GetFocusedAnnot).To(BeNil())
			})
			It("returns an error when calling FORM_SetFocusedAnnot", func() {
				FORM_SetFocusedAnnot, err := PdfiumInstance.FORM_SetFocusedAnnot(&requests.FORM_SetFocusedAnnot{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_SetFocusedAnnot).To(BeNil())
			})
			It("returns an error when calling FORM_SetIndexSelected", func() {
				FORM_SetIndexSelected, err := PdfiumInstance.FORM_SetIndexSelected(&requests.FORM_SetIndexSelected{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_SetIndexSelected).To(BeNil())
			})
			It("returns an error when calling FORM_IsIndexSelected", func() {
				FORM_IsIndexSelected, err := PdfiumInstance.FORM_IsIndexSelected(&requests.FORM_IsIndexSelected{})
				Expect(err).To(MatchError("formHandle not given"))
				Expect(FORM_IsIndexSelected).To(BeNil())
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

	Context("a PDF file with a text form", func() {
		var doc references.FPDF_DOCUMENT
		var formHandle references.FPDF_FORMHANDLE
		formHistory := []FormHistory{}
		timers := map[int]*FormTicker{}
		var bitmap references.FPDF_BITMAP
		var img *image.RGBA
		renderCount := 0

		addToHistory := func(history FormHistory) {
			formHistory = append(formHistory, history)
			//log.Printf("New history: %s: %v", history.Name, history.Args)
		}

		renderFormImage := func(page references.FPDF_PAGE) {
			FPDFBitmap_FillRect, err := PdfiumInstance.FPDFBitmap_FillRect(&requests.FPDFBitmap_FillRect{
				Bitmap: bitmap,
				Color:  0xFFFFFFFF,
				Left:   0,
				Top:    0,
				Width:  300,
				Height: 300,
			})
			Expect(err).To(BeNil())
			Expect(FPDFBitmap_FillRect).To(Equal(&responses.FPDFBitmap_FillRect{}))

			FPDF_RenderPageBitmap, err := PdfiumInstance.FPDF_RenderPageBitmap(&requests.FPDF_RenderPageBitmap{
				Bitmap: bitmap,
				Page: requests.Page{
					ByReference: &page,
				},
				StartX: 0,
				StartY: 0,
				SizeX:  300,
				SizeY:  300,
				Rotate: enums.FPDF_PAGE_ROTATION_NONE,
				Flags:  enums.FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER,
			})

			Expect(err).To(BeNil())
			Expect(FPDF_RenderPageBitmap).To(Equal(&responses.FPDF_RenderPageBitmap{}))

			FPDF_FFLDraw, err := PdfiumInstance.FPDF_FFLDraw(&requests.FPDF_FFLDraw{
				FormHandle: formHandle,
				Bitmap:     bitmap,
				Page: requests.Page{
					ByReference: &page,
				},
				StartX: 0,
				StartY: 0,
				SizeX:  300,
				SizeY:  300,
				Rotate: enums.FPDF_PAGE_ROTATION_NONE,
				Flags:  enums.FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER,
			})

			Expect(err).To(BeNil())
			Expect(FPDF_FFLDraw).To(Equal(&responses.FPDF_FFLDraw{}))

			var opt jpeg.Options
			opt.Quality = 95

			var imgBuf bytes.Buffer
			err = jpeg.Encode(&imgBuf, img, &opt)
			if err != nil {
				return
			}

			ioutil.WriteFile(TestDataPath+"/testdata/"+fmt.Sprintf("render_fpdf_formfill_experimental-%d.jpg", renderCount), imgBuf.Bytes(), 0777)
			renderCount++
			//log.Println("did render")
		}

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
						addToHistory(FormHistory{
							Name: "Release",
						})
					},
					FFI_Invalidate: func(page references.FPDF_PAGE, left, top, right, bottom float64) {
						addToHistory(FormHistory{
							Name: "FFI_Invalidate",
							Args: []interface{}{page, left, top, right, bottom},
						})
					},
					FFI_OutputSelectedRect: func(page references.FPDF_PAGE, left, top, right, bottom float64) {
						addToHistory(FormHistory{
							Name: "FFI_OutputSelectedRect",
							Args: []interface{}{page, left, top, right, bottom},
						})
					},
					FFI_SetCursor: func(cursorType enums.FXCT) {
						addToHistory(FormHistory{
							Name: "FFI_SetCursor",
							Args: []interface{}{cursorType},
						})
					},
					FFI_SetTimer: func(elapse int, timerFunc func(idEvent int)) int {
						addToHistory(FormHistory{
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
						addToHistory(FormHistory{
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
						addToHistory(FormHistory{
							Name: "FFI_GetLocalTime",
						})
						return structs.FPDF_SYSTEMTIME{}
					},
					FFI_OnChange: func() {
						addToHistory(FormHistory{
							Name: "FFI_OnChange",
						})
					},
					FFI_GetPage: func(document references.FPDF_DOCUMENT, index int) *references.FPDF_PAGE {
						addToHistory(FormHistory{
							Name: "FFI_GetPage",
							Args: []interface{}{document, index},
						})

						return nil
					},
					FFI_GetCurrentPage: func(document references.FPDF_DOCUMENT) *references.FPDF_PAGE {
						addToHistory(FormHistory{
							Name: "FFI_GetCurrentPage",
							Args: []interface{}{document},
						})
						return nil
					},
					FFI_GetRotation: func(page references.FPDF_PAGE) enums.FPDF_PAGE_ROTATION {
						addToHistory(FormHistory{
							Name: "FFI_GetRotation",
							Args: []interface{}{page},
						})
						return enums.FPDF_PAGE_ROTATION_NONE
					},
					FFI_ExecuteNamedAction: func(namedAction string) {
						addToHistory(FormHistory{
							Name: "FFI_ExecuteNamedAction",
							Args: []interface{}{namedAction},
						})
					},
					FFI_SetTextFieldFocus: func(value string, isFocus bool) {
						addToHistory(FormHistory{
							Name: "FFI_SetTextFieldFocus",
							Args: []interface{}{value, isFocus},
						})
					},
					FFI_DoURIAction: func(bsURI string) {
						addToHistory(FormHistory{
							Name: "FFI_DoURIAction",
							Args: []interface{}{bsURI},
						})
					},
					FFI_DoGoToAction: func(pageIndex int, zoomMode enums.FPDF_ZOOM_MODE, pos []float32) {
						addToHistory(FormHistory{
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

			renderCount = 0
			img = image.NewRGBA(image.Rect(0, 0, 300, 300))
			FPDFBitmap_CreateEx, err := PdfiumInstance.FPDFBitmap_CreateEx(&requests.FPDFBitmap_CreateEx{
				Width:  300,
				Height: 300,
				Format: enums.FPDF_BITMAP_FORMAT_BGRA,
				Buffer: img.Pix,
				Stride: img.Stride,
			})
			Expect(err).To(BeNil())
			Expect(FPDFBitmap_CreateEx).To(Not(BeNil()))
			bitmap = FPDFBitmap_CreateEx.Bitmap
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

			FPDFBitmap_Destroy, err := PdfiumInstance.FPDFBitmap_Destroy(&requests.FPDFBitmap_Destroy{
				Bitmap: bitmap,
			})
			Expect(err).To(BeNil())
			Expect(FPDFBitmap_Destroy).To(Equal(&responses.FPDFBitmap_Destroy{}))

			//formattedHistory, _ := json.MarshalIndent(formHistory, "", "  ")

			//log.Println(formHistory)
			//log.Printf(string(formattedHistory))
		})

		When("is opened", func() {
			When("no page is given", func() {
				It("returns an error when calling FORM_OnMouseWheel", func() {
					FORM_OnMouseWheel, err := PdfiumInstance.FORM_OnMouseWheel(&requests.FORM_OnMouseWheel{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_OnMouseWheel).To(BeNil())
				})

				It("returns an error when calling FORM_GetFocusedText", func() {
					FORM_GetFocusedText, err := PdfiumInstance.FORM_GetFocusedText(&requests.FORM_GetFocusedText{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_GetFocusedText).To(BeNil())
				})

				It("returns an error when calling FORM_SelectAllText", func() {
					FORM_SelectAllText, err := PdfiumInstance.FORM_SelectAllText(&requests.FORM_SelectAllText{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_SelectAllText).To(BeNil())
				})

				It("returns an error when calling FORM_SetIndexSelected", func() {
					FORM_SetIndexSelected, err := PdfiumInstance.FORM_SetIndexSelected(&requests.FORM_SetIndexSelected{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_SetIndexSelected).To(BeNil())
				})

				It("returns an error when calling FORM_IsIndexSelected", func() {
					FORM_IsIndexSelected, err := PdfiumInstance.FORM_IsIndexSelected(&requests.FORM_IsIndexSelected{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FORM_IsIndexSelected).To(BeNil())
				})

				It("returns an error when calling FPDFAnnot_GetFormFieldAtPoint", func() {
					FPDFAnnot_GetFormFieldAtPoint, err := PdfiumInstance.FPDFAnnot_GetFormFieldAtPoint(&requests.FPDFAnnot_GetFormFieldAtPoint{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(FPDFAnnot_GetFormFieldAtPoint).To(BeNil())
				})
			})

			When("no annotation is given", func() {
				It("returns an error when calling FORM_SetFocusedAnnot", func() {
					FORM_SetFocusedAnnot, err := PdfiumInstance.FORM_SetFocusedAnnot(&requests.FORM_SetFocusedAnnot{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("annotation not given"))
					Expect(FORM_SetFocusedAnnot).To(BeNil())
				})

				It("returns an error when calling FPDFAnnot_GetFormFieldFlags", func() {
					FPDFAnnot_GetFormFieldFlags, err := PdfiumInstance.FPDFAnnot_GetFormFieldFlags(&requests.FPDFAnnot_GetFormFieldFlags{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("annotation not given"))
					Expect(FPDFAnnot_GetFormFieldFlags).To(BeNil())
				})

				It("returns an error when calling FPDFAnnot_GetFormFieldName", func() {
					FPDFAnnot_GetFormFieldName, err := PdfiumInstance.FPDFAnnot_GetFormFieldName(&requests.FPDFAnnot_GetFormFieldName{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("annotation not given"))
					Expect(FPDFAnnot_GetFormFieldName).To(BeNil())
				})

				It("returns an error when calling FPDFAnnot_GetFormFieldType", func() {
					FPDFAnnot_GetFormFieldType, err := PdfiumInstance.FPDFAnnot_GetFormFieldType(&requests.FPDFAnnot_GetFormFieldType{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("annotation not given"))
					Expect(FPDFAnnot_GetFormFieldType).To(BeNil())
				})

				It("returns an error when calling FPDFAnnot_GetFormFieldValue", func() {
					FPDFAnnot_GetFormFieldValue, err := PdfiumInstance.FPDFAnnot_GetFormFieldValue(&requests.FPDFAnnot_GetFormFieldValue{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("annotation not given"))
					Expect(FPDFAnnot_GetFormFieldValue).To(BeNil())
				})

				It("returns an error when calling FPDFAnnot_GetOptionCount", func() {
					FPDFAnnot_GetOptionCount, err := PdfiumInstance.FPDFAnnot_GetOptionCount(&requests.FPDFAnnot_GetOptionCount{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("annotation not given"))
					Expect(FPDFAnnot_GetOptionCount).To(BeNil())
				})

				It("returns an error when calling FPDFAnnot_GetOptionLabel", func() {
					FPDFAnnot_GetOptionLabel, err := PdfiumInstance.FPDFAnnot_GetOptionLabel(&requests.FPDFAnnot_GetOptionLabel{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("annotation not given"))
					Expect(FPDFAnnot_GetOptionLabel).To(BeNil())
				})

				It("returns an error when calling FPDFAnnot_IsOptionSelected", func() {
					FPDFAnnot_IsOptionSelected, err := PdfiumInstance.FPDFAnnot_IsOptionSelected(&requests.FPDFAnnot_IsOptionSelected{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("annotation not given"))
					Expect(FPDFAnnot_IsOptionSelected).To(BeNil())
				})

				It("returns an error when calling FPDFAnnot_GetOptionLabel", func() {
					FPDFAnnot_GetFontSize, err := PdfiumInstance.FPDFAnnot_GetFontSize(&requests.FPDFAnnot_GetFontSize{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("annotation not given"))
					Expect(FPDFAnnot_GetFontSize).To(BeNil())
				})

				It("returns an error when calling FPDFAnnot_IsChecked", func() {
					FPDFAnnot_IsChecked, err := PdfiumInstance.FPDFAnnot_IsChecked(&requests.FPDFAnnot_IsChecked{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("annotation not given"))
					Expect(FPDFAnnot_IsChecked).To(BeNil())
				})

				It("returns an error when calling FPDFAnnot_GetFormControlCount", func() {
					FPDFAnnot_GetFormControlCount, err := PdfiumInstance.FPDFAnnot_GetFormControlCount(&requests.FPDFAnnot_GetFormControlCount{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("annotation not given"))
					Expect(FPDFAnnot_GetFormControlCount).To(BeNil())
				})

				It("returns an error when calling FPDFAnnot_GetFormControlIndex", func() {
					FPDFAnnot_GetFormControlIndex, err := PdfiumInstance.FPDFAnnot_GetFormControlIndex(&requests.FPDFAnnot_GetFormControlIndex{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("annotation not given"))
					Expect(FPDFAnnot_GetFormControlIndex).To(BeNil())
				})

				It("returns an error when calling FPDFAnnot_GetFormFieldExportValue", func() {
					FPDFAnnot_GetFormFieldExportValue, err := PdfiumInstance.FPDFAnnot_GetFormFieldExportValue(&requests.FPDFAnnot_GetFormFieldExportValue{
						FormHandle: formHandle,
					})
					Expect(err).To(MatchError("annotation not given"))
					Expect(FPDFAnnot_GetFormFieldExportValue).To(BeNil())
				})
			})

			It("allows to get the form type", func() {
				FPDF_GetFormType, err := PdfiumInstance.FPDF_GetFormType(&requests.FPDF_GetFormType{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetFormType).To(Equal(&responses.FPDF_GetFormType{
					FormType: enums.FPDF_FORMTYPE_ACRO_FORM,
				}))
			})

			It("allows experimental form methods to be called", func() {
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

				FORM_ReplaceSelection, err := PdfiumInstance.FORM_ReplaceSelection(&requests.FORM_ReplaceSelection{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					Text:       "Jeroen",
				})
				Expect(err).To(BeNil())
				Expect(FORM_ReplaceSelection).To(Equal(&responses.FORM_ReplaceSelection{}))
				renderFormImage(FPDF_LoadPage.Page)

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
				renderFormImage(FPDF_LoadPage.Page)

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
				renderFormImage(FPDF_LoadPage.Page)

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
				renderFormImage(FPDF_LoadPage.Page)

				FORM_GetFocusedText, err := PdfiumInstance.FORM_GetFocusedText(&requests.FORM_GetFocusedText{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
				})
				Expect(err).To(BeNil())
				Expect(FORM_GetFocusedText).To(Equal(&responses.FORM_GetFocusedText{
					FocusedText: "JeroenABC",
				}))

				FORM_SelectAllText, err := PdfiumInstance.FORM_SelectAllText(&requests.FORM_SelectAllText{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
				})
				Expect(err).To(BeNil())
				Expect(FORM_SelectAllText).To(Equal(&responses.FORM_SelectAllText{}))

				FORM_GetSelectedText, err := PdfiumInstance.FORM_GetSelectedText(&requests.FORM_GetSelectedText{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
				})
				Expect(err).To(BeNil())
				Expect(FORM_GetSelectedText).To(Equal(&responses.FORM_GetSelectedText{
					SelectedText: "JeroenABC",
				}))

				PdfiumInstance.FORM_OnMouseWheel(&requests.FORM_OnMouseWheel{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					PageCoord: structs.FPDF_FS_POINTF{
						X: 30,
						Y: 30,
					},
					DeltaX:   2,
					DeltaY:   2,
					Modifier: 0,
				})
				// For some reason mousewheel errors here. Perhaps one returns success on scroll lists?

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
		})
	})

	Context("a PDF file with a form annotation", func() {
		var doc references.FPDF_DOCUMENT
		var formHandle references.FPDF_FORMHANDLE
		formHistory := []FormHistory{}
		timers := map[int]*FormTicker{}
		var bitmap references.FPDF_BITMAP
		var img *image.RGBA
		renderCount := 0

		addToHistory := func(history FormHistory) {
			formHistory = append(formHistory, history)
			//log.Printf("New history: %s: %v", history.Name, history.Args)
		}

		renderFormImage := func(page references.FPDF_PAGE) {
			FPDFBitmap_FillRect, err := PdfiumInstance.FPDFBitmap_FillRect(&requests.FPDFBitmap_FillRect{
				Bitmap: bitmap,
				Color:  0xFFFFFFFF,
				Left:   0,
				Top:    0,
				Width:  300,
				Height: 300,
			})
			Expect(err).To(BeNil())
			Expect(FPDFBitmap_FillRect).To(Equal(&responses.FPDFBitmap_FillRect{}))

			FPDF_RenderPageBitmap, err := PdfiumInstance.FPDF_RenderPageBitmap(&requests.FPDF_RenderPageBitmap{
				Bitmap: bitmap,
				Page: requests.Page{
					ByReference: &page,
				},
				StartX: 0,
				StartY: 0,
				SizeX:  300,
				SizeY:  300,
				Rotate: enums.FPDF_PAGE_ROTATION_NONE,
				Flags:  enums.FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER,
			})

			Expect(err).To(BeNil())
			Expect(FPDF_RenderPageBitmap).To(Equal(&responses.FPDF_RenderPageBitmap{}))

			FPDF_FFLDraw, err := PdfiumInstance.FPDF_FFLDraw(&requests.FPDF_FFLDraw{
				FormHandle: formHandle,
				Bitmap:     bitmap,
				Page: requests.Page{
					ByReference: &page,
				},
				StartX: 0,
				StartY: 0,
				SizeX:  300,
				SizeY:  300,
				Rotate: enums.FPDF_PAGE_ROTATION_NONE,
				Flags:  enums.FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER,
			})

			Expect(err).To(BeNil())
			Expect(FPDF_FFLDraw).To(Equal(&responses.FPDF_FFLDraw{}))

			var opt jpeg.Options
			opt.Quality = 95

			var imgBuf bytes.Buffer
			err = jpeg.Encode(&imgBuf, img, &opt)
			if err != nil {
				return
			}

			ioutil.WriteFile(TestDataPath+"/testdata/"+fmt.Sprintf("render_fpdf_formfill_experimental_annotation-%d.jpg", renderCount), imgBuf.Bytes(), 0777)
			renderCount++
			//log.Println("did render")
		}

		BeforeEach(func() {
			formHistory = []FormHistory{}
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/annotiter.pdf")
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
						addToHistory(FormHistory{
							Name: "Release",
						})
					},
					FFI_Invalidate: func(page references.FPDF_PAGE, left, top, right, bottom float64) {
						addToHistory(FormHistory{
							Name: "FFI_Invalidate",
							Args: []interface{}{page, left, top, right, bottom},
						})
					},
					FFI_OutputSelectedRect: func(page references.FPDF_PAGE, left, top, right, bottom float64) {
						addToHistory(FormHistory{
							Name: "FFI_OutputSelectedRect",
							Args: []interface{}{page, left, top, right, bottom},
						})
					},
					FFI_SetCursor: func(cursorType enums.FXCT) {
						addToHistory(FormHistory{
							Name: "FFI_SetCursor",
							Args: []interface{}{cursorType},
						})
					},
					FFI_SetTimer: func(elapse int, timerFunc func(idEvent int)) int {
						addToHistory(FormHistory{
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
						addToHistory(FormHistory{
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
						addToHistory(FormHistory{
							Name: "FFI_GetLocalTime",
						})
						return structs.FPDF_SYSTEMTIME{}
					},
					FFI_OnChange: func() {
						addToHistory(FormHistory{
							Name: "FFI_OnChange",
						})
					},
					FFI_GetPage: func(document references.FPDF_DOCUMENT, index int) *references.FPDF_PAGE {
						addToHistory(FormHistory{
							Name: "FFI_GetPage",
							Args: []interface{}{document, index},
						})

						return nil
					},
					FFI_GetCurrentPage: func(document references.FPDF_DOCUMENT) *references.FPDF_PAGE {
						addToHistory(FormHistory{
							Name: "FFI_GetCurrentPage",
							Args: []interface{}{document},
						})
						return nil
					},
					FFI_GetRotation: func(page references.FPDF_PAGE) enums.FPDF_PAGE_ROTATION {
						addToHistory(FormHistory{
							Name: "FFI_GetRotation",
							Args: []interface{}{page},
						})
						return enums.FPDF_PAGE_ROTATION_NONE
					},
					FFI_ExecuteNamedAction: func(namedAction string) {
						addToHistory(FormHistory{
							Name: "FFI_ExecuteNamedAction",
							Args: []interface{}{namedAction},
						})
					},
					FFI_SetTextFieldFocus: func(value string, isFocus bool) {
						addToHistory(FormHistory{
							Name: "FFI_SetTextFieldFocus",
							Args: []interface{}{value, isFocus},
						})
					},
					FFI_DoURIAction: func(bsURI string) {
						addToHistory(FormHistory{
							Name: "FFI_DoURIAction",
							Args: []interface{}{bsURI},
						})
					},
					FFI_DoGoToAction: func(pageIndex int, zoomMode enums.FPDF_ZOOM_MODE, pos []float32) {
						addToHistory(FormHistory{
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

			renderCount = 0
			img = image.NewRGBA(image.Rect(0, 0, 300, 300))
			FPDFBitmap_CreateEx, err := PdfiumInstance.FPDFBitmap_CreateEx(&requests.FPDFBitmap_CreateEx{
				Width:  300,
				Height: 300,
				Format: enums.FPDF_BITMAP_FORMAT_BGRA,
				Buffer: img.Pix,
				Stride: img.Stride,
			})
			Expect(err).To(BeNil())
			Expect(FPDFBitmap_CreateEx).To(Not(BeNil()))
			bitmap = FPDFBitmap_CreateEx.Bitmap
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

			FPDFBitmap_Destroy, err := PdfiumInstance.FPDFBitmap_Destroy(&requests.FPDFBitmap_Destroy{
				Bitmap: bitmap,
			})
			Expect(err).To(BeNil())
			Expect(FPDFBitmap_Destroy).To(Equal(&responses.FPDFBitmap_Destroy{}))

			//formattedHistory, _ := json.MarshalIndent(formHistory, "", "  ")

			//log.Println(formHistory)
			//log.Printf(string(formattedHistory))
		})

		When("is opened", func() {
			It("allows experimental form methods to be called", func() {
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

				FORM_OnMouseMove, err := PdfiumInstance.FORM_OnMouseMove(&requests.FORM_OnMouseMove{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					PageX:      410,
					PageY:      210,
					Modifier:   0,
				})
				Expect(err).To(BeNil())
				Expect(FORM_OnMouseMove).To(Equal(&responses.FORM_OnMouseMove{}))

				FORM_OnLButtonDown, err := PdfiumInstance.FORM_OnLButtonDown(&requests.FORM_OnLButtonDown{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					PageX:      410,
					PageY:      210,
					Modifier:   0,
				})
				Expect(err).To(BeNil())
				Expect(FORM_OnLButtonDown).To(Equal(&responses.FORM_OnLButtonDown{}))

				FORM_OnLButtonUp, err := PdfiumInstance.FORM_OnLButtonUp(&requests.FORM_OnLButtonUp{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					PageX:      410,
					PageY:      210,
					Modifier:   0,
				})
				Expect(err).To(BeNil())
				Expect(FORM_OnLButtonUp).To(Equal(&responses.FORM_OnLButtonUp{}))
				renderFormImage(FPDF_LoadPage.Page)

				FORM_GetFocusedAnnot, err := PdfiumInstance.FORM_GetFocusedAnnot(&requests.FORM_GetFocusedAnnot{
					FormHandle: formHandle,
				})
				Expect(err).To(BeNil())
				Expect(FORM_GetFocusedAnnot).ToNot(BeNil())
				Expect(FORM_GetFocusedAnnot.Annotation).ToNot(BeEmpty())
				Expect(FORM_GetFocusedAnnot.PageIndex).To(Equal(0))
				renderFormImage(FPDF_LoadPage.Page)

				FORM_SetFocusedAnnot, err := PdfiumInstance.FORM_SetFocusedAnnot(&requests.FORM_SetFocusedAnnot{
					FormHandle: formHandle,
					Annotation: FORM_GetFocusedAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FORM_SetFocusedAnnot).To(Equal(&responses.FORM_SetFocusedAnnot{}))

				FPDFAnnot_GetFormFieldFlags, err := PdfiumInstance.FPDFAnnot_GetFormFieldFlags(&requests.FPDFAnnot_GetFormFieldFlags{
					FormHandle: formHandle,
					Annotation: FORM_GetFocusedAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFormFieldFlags).To(Equal(&responses.FPDFAnnot_GetFormFieldFlags{
					Flags: enums.FPDF_FORMFLAG_NONE,
				}))

				FPDFAnnot_GetFormFieldAtPoint, err := PdfiumInstance.FPDFAnnot_GetFormFieldAtPoint(&requests.FPDFAnnot_GetFormFieldAtPoint{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					Point: structs.FPDF_FS_POINTF{
						X: 410,
						Y: 210,
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFormFieldAtPoint).ToNot(BeNil())
				Expect(FPDFAnnot_GetFormFieldAtPoint.Annotation).ToNot(BeEmpty())

				FPDFPage_CloseAnnot, err := PdfiumInstance.FPDFPage_CloseAnnot(&requests.FPDFPage_CloseAnnot{
					Annotation: FPDFAnnot_GetFormFieldAtPoint.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CloseAnnot).To(Equal(&responses.FPDFPage_CloseAnnot{}))

				FPDFAnnot_GetFormFieldName, err := PdfiumInstance.FPDFAnnot_GetFormFieldName(&requests.FPDFAnnot_GetFormFieldName{
					FormHandle: formHandle,
					Annotation: FORM_GetFocusedAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFormFieldName).To(Equal(&responses.FPDFAnnot_GetFormFieldName{
					FormFieldName: "Sub_RightBottom",
				}))

				FPDFAnnot_GetFormFieldType, err := PdfiumInstance.FPDFAnnot_GetFormFieldType(&requests.FPDFAnnot_GetFormFieldType{
					FormHandle: formHandle,
					Annotation: FORM_GetFocusedAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFormFieldType).To(Equal(&responses.FPDFAnnot_GetFormFieldType{
					FormFieldType: enums.FPDF_FORMFIELD_TYPE_TEXTFIELD,
				}))

				FPDFAnnot_GetFormFieldValue, err := PdfiumInstance.FPDFAnnot_GetFormFieldValue(&requests.FPDFAnnot_GetFormFieldValue{
					FormHandle: formHandle,
					Annotation: FORM_GetFocusedAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFormFieldValue).To(Equal(&responses.FPDFAnnot_GetFormFieldValue{}))

				FPDFAnnot_GetFontSize, err := PdfiumInstance.FPDFAnnot_GetFontSize(&requests.FPDFAnnot_GetFontSize{
					FormHandle: formHandle,
					Annotation: FORM_GetFocusedAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFontSize).To(Equal(&responses.FPDFAnnot_GetFontSize{
					FontSize: 0,
				}))

				FPDFPage_CloseAnnot, err = PdfiumInstance.FPDFPage_CloseAnnot(&requests.FPDFPage_CloseAnnot{
					Annotation: FORM_GetFocusedAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CloseAnnot).To(Equal(&responses.FPDFPage_CloseAnnot{}))

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
		})
	})

	Context("a PDF file with a form combobox", func() {
		var doc references.FPDF_DOCUMENT
		var formHandle references.FPDF_FORMHANDLE
		formHistory := []FormHistory{}
		timers := map[int]*FormTicker{}
		var bitmap references.FPDF_BITMAP
		var img *image.RGBA
		renderCount := 0

		addToHistory := func(history FormHistory) {
			formHistory = append(formHistory, history)
			//log.Printf("New history: %s: %v", history.Name, history.Args)
		}

		renderFormImage := func(page references.FPDF_PAGE) {
			FPDFBitmap_FillRect, err := PdfiumInstance.FPDFBitmap_FillRect(&requests.FPDFBitmap_FillRect{
				Bitmap: bitmap,
				Color:  0xFFFFFFFF,
				Left:   0,
				Top:    0,
				Width:  300,
				Height: 300,
			})
			Expect(err).To(BeNil())
			Expect(FPDFBitmap_FillRect).To(Equal(&responses.FPDFBitmap_FillRect{}))

			FPDF_RenderPageBitmap, err := PdfiumInstance.FPDF_RenderPageBitmap(&requests.FPDF_RenderPageBitmap{
				Bitmap: bitmap,
				Page: requests.Page{
					ByReference: &page,
				},
				StartX: 0,
				StartY: 0,
				SizeX:  300,
				SizeY:  300,
				Rotate: enums.FPDF_PAGE_ROTATION_NONE,
				Flags:  enums.FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER,
			})

			Expect(err).To(BeNil())
			Expect(FPDF_RenderPageBitmap).To(Equal(&responses.FPDF_RenderPageBitmap{}))

			FPDF_FFLDraw, err := PdfiumInstance.FPDF_FFLDraw(&requests.FPDF_FFLDraw{
				FormHandle: formHandle,
				Bitmap:     bitmap,
				Page: requests.Page{
					ByReference: &page,
				},
				StartX: 0,
				StartY: 0,
				SizeX:  300,
				SizeY:  300,
				Rotate: enums.FPDF_PAGE_ROTATION_NONE,
				Flags:  enums.FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER,
			})

			Expect(err).To(BeNil())
			Expect(FPDF_FFLDraw).To(Equal(&responses.FPDF_FFLDraw{}))

			var opt jpeg.Options
			opt.Quality = 95

			var imgBuf bytes.Buffer
			err = jpeg.Encode(&imgBuf, img, &opt)
			if err != nil {
				return
			}

			ioutil.WriteFile(TestDataPath+"/testdata/"+fmt.Sprintf("render_fpdf_formfill_experimental_combobox-%d.jpg", renderCount), imgBuf.Bytes(), 0777)
			renderCount++
			//log.Println("did render")
		}

		BeforeEach(func() {
			formHistory = []FormHistory{}
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/combobox_form.pdf")
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
						addToHistory(FormHistory{
							Name: "Release",
						})
					},
					FFI_Invalidate: func(page references.FPDF_PAGE, left, top, right, bottom float64) {
						addToHistory(FormHistory{
							Name: "FFI_Invalidate",
							Args: []interface{}{page, left, top, right, bottom},
						})
					},
					FFI_OutputSelectedRect: func(page references.FPDF_PAGE, left, top, right, bottom float64) {
						addToHistory(FormHistory{
							Name: "FFI_OutputSelectedRect",
							Args: []interface{}{page, left, top, right, bottom},
						})
					},
					FFI_SetCursor: func(cursorType enums.FXCT) {
						addToHistory(FormHistory{
							Name: "FFI_SetCursor",
							Args: []interface{}{cursorType},
						})
					},
					FFI_SetTimer: func(elapse int, timerFunc func(idEvent int)) int {
						addToHistory(FormHistory{
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
						addToHistory(FormHistory{
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
						addToHistory(FormHistory{
							Name: "FFI_GetLocalTime",
						})
						return structs.FPDF_SYSTEMTIME{}
					},
					FFI_OnChange: func() {
						addToHistory(FormHistory{
							Name: "FFI_OnChange",
						})
					},
					FFI_GetPage: func(document references.FPDF_DOCUMENT, index int) *references.FPDF_PAGE {
						addToHistory(FormHistory{
							Name: "FFI_GetPage",
							Args: []interface{}{document, index},
						})

						return nil
					},
					FFI_GetCurrentPage: func(document references.FPDF_DOCUMENT) *references.FPDF_PAGE {
						addToHistory(FormHistory{
							Name: "FFI_GetCurrentPage",
							Args: []interface{}{document},
						})
						return nil
					},
					FFI_GetRotation: func(page references.FPDF_PAGE) enums.FPDF_PAGE_ROTATION {
						addToHistory(FormHistory{
							Name: "FFI_GetRotation",
							Args: []interface{}{page},
						})
						return enums.FPDF_PAGE_ROTATION_NONE
					},
					FFI_ExecuteNamedAction: func(namedAction string) {
						addToHistory(FormHistory{
							Name: "FFI_ExecuteNamedAction",
							Args: []interface{}{namedAction},
						})
					},
					FFI_SetTextFieldFocus: func(value string, isFocus bool) {
						addToHistory(FormHistory{
							Name: "FFI_SetTextFieldFocus",
							Args: []interface{}{value, isFocus},
						})
					},
					FFI_DoURIAction: func(bsURI string) {
						addToHistory(FormHistory{
							Name: "FFI_DoURIAction",
							Args: []interface{}{bsURI},
						})
					},
					FFI_DoGoToAction: func(pageIndex int, zoomMode enums.FPDF_ZOOM_MODE, pos []float32) {
						addToHistory(FormHistory{
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

			renderCount = 0
			img = image.NewRGBA(image.Rect(0, 0, 300, 300))
			FPDFBitmap_CreateEx, err := PdfiumInstance.FPDFBitmap_CreateEx(&requests.FPDFBitmap_CreateEx{
				Width:  300,
				Height: 300,
				Format: enums.FPDF_BITMAP_FORMAT_BGRA,
				Buffer: img.Pix,
				Stride: img.Stride,
			})
			Expect(err).To(BeNil())
			Expect(FPDFBitmap_CreateEx).To(Not(BeNil()))
			bitmap = FPDFBitmap_CreateEx.Bitmap
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

			FPDFBitmap_Destroy, err := PdfiumInstance.FPDFBitmap_Destroy(&requests.FPDFBitmap_Destroy{
				Bitmap: bitmap,
			})
			Expect(err).To(BeNil())
			Expect(FPDFBitmap_Destroy).To(Equal(&responses.FPDFBitmap_Destroy{}))

			//formattedHistory, _ := json.MarshalIndent(formHistory, "", "  ")

			//log.Println(formHistory)
			//log.Printf(string(formattedHistory))
		})

		When("is opened", func() {
			It("allows experimental form methods to be called", func() {
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

				renderFormImage(FPDF_LoadPage.Page)

				FORM_OnMouseMove, err := PdfiumInstance.FORM_OnMouseMove(&requests.FORM_OnMouseMove{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					PageX:      192,
					PageY:      410,
					Modifier:   0,
				})
				Expect(err).To(BeNil())
				Expect(FORM_OnMouseMove).To(Equal(&responses.FORM_OnMouseMove{}))

				FORM_OnLButtonDown, err := PdfiumInstance.FORM_OnLButtonDown(&requests.FORM_OnLButtonDown{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					PageX:      192,
					PageY:      410,
					Modifier:   0,
				})
				Expect(err).To(BeNil())
				Expect(FORM_OnLButtonDown).To(Equal(&responses.FORM_OnLButtonDown{}))

				FORM_OnLButtonUp, err := PdfiumInstance.FORM_OnLButtonUp(&requests.FORM_OnLButtonUp{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					FormHandle: formHandle,
					PageX:      192,
					PageY:      410,
					Modifier:   0,
				})
				Expect(err).To(BeNil())
				Expect(FORM_OnLButtonUp).To(Equal(&responses.FORM_OnLButtonUp{}))
				renderFormImage(FPDF_LoadPage.Page)

				FORM_SetIndexSelected, err := PdfiumInstance.FORM_SetIndexSelected(&requests.FORM_SetIndexSelected{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					Index:    2,
					Selected: true,
				})
				Expect(err).To(BeNil())
				Expect(FORM_SetIndexSelected).To(Equal(&responses.FORM_SetIndexSelected{}))
				renderFormImage(FPDF_LoadPage.Page)

				FORM_IsIndexSelected, err := PdfiumInstance.FORM_IsIndexSelected(&requests.FORM_IsIndexSelected{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					Index: 2,
				})
				Expect(err).To(BeNil())
				Expect(FORM_IsIndexSelected).To(Equal(&responses.FORM_IsIndexSelected{
					IsIndexSelected: true,
				}))
				renderFormImage(FPDF_LoadPage.Page)

				FORM_SetIndexSelected, err = PdfiumInstance.FORM_SetIndexSelected(&requests.FORM_SetIndexSelected{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					Index:    -25,
					Selected: true,
				})
				Expect(err).To(MatchError("could not set index selected"))
				Expect(FORM_SetIndexSelected).To(BeNil())

				FORM_GetFocusedAnnot, err := PdfiumInstance.FORM_GetFocusedAnnot(&requests.FORM_GetFocusedAnnot{
					FormHandle: formHandle,
				})
				Expect(err).To(BeNil())
				Expect(FORM_GetFocusedAnnot).ToNot(BeNil())
				Expect(FORM_GetFocusedAnnot.Annotation).ToNot(BeEmpty())
				Expect(FORM_GetFocusedAnnot.PageIndex).To(Equal(0))
				renderFormImage(FPDF_LoadPage.Page)

				FPDFAnnot_GetOptionCount, err := PdfiumInstance.FPDFAnnot_GetOptionCount(&requests.FPDFAnnot_GetOptionCount{
					FormHandle: formHandle,
					Annotation: FORM_GetFocusedAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetOptionCount).To(Equal(&responses.FPDFAnnot_GetOptionCount{
					OptionCount: 26,
				}))

				FPDFAnnot_GetOptionLabel, err := PdfiumInstance.FPDFAnnot_GetOptionLabel(&requests.FPDFAnnot_GetOptionLabel{
					FormHandle: formHandle,
					Annotation: FORM_GetFocusedAnnot.Annotation,
					Index:      1,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetOptionLabel).To(Equal(&responses.FPDFAnnot_GetOptionLabel{
					OptionLabel: "Banana",
				}))

				FPDFAnnot_IsOptionSelected, err := PdfiumInstance.FPDFAnnot_IsOptionSelected(&requests.FPDFAnnot_IsOptionSelected{
					FormHandle: formHandle,
					Annotation: FORM_GetFocusedAnnot.Annotation,
					Index:      1,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_IsOptionSelected).To(Equal(&responses.FPDFAnnot_IsOptionSelected{
					IsOptionSelected: true,
				}))

				FPDFPage_CloseAnnot, err := PdfiumInstance.FPDFPage_CloseAnnot(&requests.FPDFPage_CloseAnnot{
					Annotation: FORM_GetFocusedAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CloseAnnot).To(Equal(&responses.FPDFPage_CloseAnnot{}))

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
		})
	})

	Context("a PDF file with a form checkbox/radiobutton", func() {
		var doc references.FPDF_DOCUMENT
		var formHandle references.FPDF_FORMHANDLE
		formHistory := []FormHistory{}
		timers := map[int]*FormTicker{}
		var bitmap references.FPDF_BITMAP
		var img *image.RGBA
		renderCount := 0

		addToHistory := func(history FormHistory) {
			formHistory = append(formHistory, history)
			//log.Printf("New history: %s: %v", history.Name, history.Args)
		}

		renderFormImage := func(page references.FPDF_PAGE) {
			FPDFBitmap_FillRect, err := PdfiumInstance.FPDFBitmap_FillRect(&requests.FPDFBitmap_FillRect{
				Bitmap: bitmap,
				Color:  0xFFFFFFFF,
				Left:   0,
				Top:    0,
				Width:  300,
				Height: 300,
			})
			Expect(err).To(BeNil())
			Expect(FPDFBitmap_FillRect).To(Equal(&responses.FPDFBitmap_FillRect{}))

			FPDF_RenderPageBitmap, err := PdfiumInstance.FPDF_RenderPageBitmap(&requests.FPDF_RenderPageBitmap{
				Bitmap: bitmap,
				Page: requests.Page{
					ByReference: &page,
				},
				StartX: 0,
				StartY: 0,
				SizeX:  300,
				SizeY:  300,
				Rotate: enums.FPDF_PAGE_ROTATION_NONE,
				Flags:  enums.FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER,
			})

			Expect(err).To(BeNil())
			Expect(FPDF_RenderPageBitmap).To(Equal(&responses.FPDF_RenderPageBitmap{}))

			FPDF_FFLDraw, err := PdfiumInstance.FPDF_FFLDraw(&requests.FPDF_FFLDraw{
				FormHandle: formHandle,
				Bitmap:     bitmap,
				Page: requests.Page{
					ByReference: &page,
				},
				StartX: 0,
				StartY: 0,
				SizeX:  300,
				SizeY:  300,
				Rotate: enums.FPDF_PAGE_ROTATION_NONE,
				Flags:  enums.FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER,
			})

			Expect(err).To(BeNil())
			Expect(FPDF_FFLDraw).To(Equal(&responses.FPDF_FFLDraw{}))

			var opt jpeg.Options
			opt.Quality = 95

			var imgBuf bytes.Buffer
			err = jpeg.Encode(&imgBuf, img, &opt)
			if err != nil {
				return
			}

			ioutil.WriteFile(TestDataPath+"/testdata/"+fmt.Sprintf("render_fpdf_formfill_experimental_combobox-%d.jpg", renderCount), imgBuf.Bytes(), 0777)
			renderCount++
			//log.Println("did render")
		}

		BeforeEach(func() {
			formHistory = []FormHistory{}
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
					Release: func() {
						addToHistory(FormHistory{
							Name: "Release",
						})
					},
					FFI_Invalidate: func(page references.FPDF_PAGE, left, top, right, bottom float64) {
						addToHistory(FormHistory{
							Name: "FFI_Invalidate",
							Args: []interface{}{page, left, top, right, bottom},
						})
					},
					FFI_OutputSelectedRect: func(page references.FPDF_PAGE, left, top, right, bottom float64) {
						addToHistory(FormHistory{
							Name: "FFI_OutputSelectedRect",
							Args: []interface{}{page, left, top, right, bottom},
						})
					},
					FFI_SetCursor: func(cursorType enums.FXCT) {
						addToHistory(FormHistory{
							Name: "FFI_SetCursor",
							Args: []interface{}{cursorType},
						})
					},
					FFI_SetTimer: func(elapse int, timerFunc func(idEvent int)) int {
						addToHistory(FormHistory{
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
						addToHistory(FormHistory{
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
						addToHistory(FormHistory{
							Name: "FFI_GetLocalTime",
						})
						return structs.FPDF_SYSTEMTIME{}
					},
					FFI_OnChange: func() {
						addToHistory(FormHistory{
							Name: "FFI_OnChange",
						})
					},
					FFI_GetPage: func(document references.FPDF_DOCUMENT, index int) *references.FPDF_PAGE {
						addToHistory(FormHistory{
							Name: "FFI_GetPage",
							Args: []interface{}{document, index},
						})

						return nil
					},
					FFI_GetCurrentPage: func(document references.FPDF_DOCUMENT) *references.FPDF_PAGE {
						addToHistory(FormHistory{
							Name: "FFI_GetCurrentPage",
							Args: []interface{}{document},
						})
						return nil
					},
					FFI_GetRotation: func(page references.FPDF_PAGE) enums.FPDF_PAGE_ROTATION {
						addToHistory(FormHistory{
							Name: "FFI_GetRotation",
							Args: []interface{}{page},
						})
						return enums.FPDF_PAGE_ROTATION_NONE
					},
					FFI_ExecuteNamedAction: func(namedAction string) {
						addToHistory(FormHistory{
							Name: "FFI_ExecuteNamedAction",
							Args: []interface{}{namedAction},
						})
					},
					FFI_SetTextFieldFocus: func(value string, isFocus bool) {
						addToHistory(FormHistory{
							Name: "FFI_SetTextFieldFocus",
							Args: []interface{}{value, isFocus},
						})
					},
					FFI_DoURIAction: func(bsURI string) {
						addToHistory(FormHistory{
							Name: "FFI_DoURIAction",
							Args: []interface{}{bsURI},
						})
					},
					FFI_DoGoToAction: func(pageIndex int, zoomMode enums.FPDF_ZOOM_MODE, pos []float32) {
						addToHistory(FormHistory{
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

			renderCount = 0
			img = image.NewRGBA(image.Rect(0, 0, 300, 300))
			FPDFBitmap_CreateEx, err := PdfiumInstance.FPDFBitmap_CreateEx(&requests.FPDFBitmap_CreateEx{
				Width:  300,
				Height: 300,
				Format: enums.FPDF_BITMAP_FORMAT_BGRA,
				Buffer: img.Pix,
				Stride: img.Stride,
			})
			Expect(err).To(BeNil())
			Expect(FPDFBitmap_CreateEx).To(Not(BeNil()))
			bitmap = FPDFBitmap_CreateEx.Bitmap
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

			FPDFBitmap_Destroy, err := PdfiumInstance.FPDFBitmap_Destroy(&requests.FPDFBitmap_Destroy{
				Bitmap: bitmap,
			})
			Expect(err).To(BeNil())
			Expect(FPDFBitmap_Destroy).To(Equal(&responses.FPDFBitmap_Destroy{}))

			//formattedHistory, _ := json.MarshalIndent(formHistory, "", "  ")

			//log.Println(formHistory)
			//log.Printf(string(formattedHistory))
		})

		When("is opened", func() {
			It("allows experimental form methods to be called", func() {
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

				renderFormImage(FPDF_LoadPage.Page)

				FPDFPage_GetAnnot, err := PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					Index: 1,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnot).ToNot(BeNil())
				Expect(FPDFPage_GetAnnot.Annotation).ToNot(BeEmpty())
				renderFormImage(FPDF_LoadPage.Page)

				FPDFAnnot_IsChecked, err := PdfiumInstance.FPDFAnnot_IsChecked(&requests.FPDFAnnot_IsChecked{
					FormHandle: formHandle,
					Annotation: FPDFPage_GetAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_IsChecked).To(Equal(&responses.FPDFAnnot_IsChecked{}))

				FPDFPage_CloseAnnot, err := PdfiumInstance.FPDFPage_CloseAnnot(&requests.FPDFPage_CloseAnnot{
					Annotation: FPDFPage_GetAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CloseAnnot).To(Equal(&responses.FPDFPage_CloseAnnot{}))

				FPDFPage_GetAnnot, err = PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					Index: 3,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnot).ToNot(BeNil())
				Expect(FPDFPage_GetAnnot.Annotation).ToNot(BeEmpty())
				renderFormImage(FPDF_LoadPage.Page)

				FPDFAnnot_GetFormControlCount, err := PdfiumInstance.FPDFAnnot_GetFormControlCount(&requests.FPDFAnnot_GetFormControlCount{
					FormHandle: formHandle,
					Annotation: FPDFPage_GetAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFormControlCount).To(Equal(&responses.FPDFAnnot_GetFormControlCount{
					FormControlCount: 3,
				}))

				FPDFAnnot_GetFormControlIndex, err := PdfiumInstance.FPDFAnnot_GetFormControlIndex(&requests.FPDFAnnot_GetFormControlIndex{
					FormHandle: formHandle,
					Annotation: FPDFPage_GetAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFormControlIndex).To(Equal(&responses.FPDFAnnot_GetFormControlIndex{
					FormControlIndex: 1,
				}))

				FPDFPage_CloseAnnot, err = PdfiumInstance.FPDFPage_CloseAnnot(&requests.FPDFPage_CloseAnnot{
					Annotation: FPDFPage_GetAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CloseAnnot).To(Equal(&responses.FPDFPage_CloseAnnot{}))

				FPDFPage_GetAnnot, err = PdfiumInstance.FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
					Index: 6,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_GetAnnot).ToNot(BeNil())
				Expect(FPDFPage_GetAnnot.Annotation).ToNot(BeEmpty())
				renderFormImage(FPDF_LoadPage.Page)

				FPDFAnnot_GetFormFieldExportValue, err := PdfiumInstance.FPDFAnnot_GetFormFieldExportValue(&requests.FPDFAnnot_GetFormFieldExportValue{
					FormHandle: formHandle,
					Annotation: FPDFPage_GetAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFormFieldExportValue).To(Equal(&responses.FPDFAnnot_GetFormFieldExportValue{
					Value: "value2",
				}))

				FPDFPage_CloseAnnot, err = PdfiumInstance.FPDFPage_CloseAnnot(&requests.FPDFPage_CloseAnnot{
					Annotation: FPDFPage_GetAnnot.Annotation,
				})
				Expect(err).To(BeNil())
				Expect(FPDFPage_CloseAnnot).To(Equal(&responses.FPDFPage_CloseAnnot{}))

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
		})
	})

	Context("a PDF file with a form and annotations", func() {
		var doc references.FPDF_DOCUMENT
		var formHandle references.FPDF_FORMHANDLE
		formHistory := []FormHistory{}
		timers := map[int]*FormTicker{}

		addToHistory := func(history FormHistory) {
			formHistory = append(formHistory, history)
			//log.Printf("New history: %s: %v", history.Name, history.Args)
		}

		BeforeEach(func() {
			formHistory = []FormHistory{}
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/annots.pdf")
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
						addToHistory(FormHistory{
							Name: "Release",
						})
					},
					FFI_Invalidate: func(page references.FPDF_PAGE, left, top, right, bottom float64) {
						addToHistory(FormHistory{
							Name: "FFI_Invalidate",
							Args: []interface{}{page, left, top, right, bottom},
						})
					},
					FFI_OutputSelectedRect: func(page references.FPDF_PAGE, left, top, right, bottom float64) {
						addToHistory(FormHistory{
							Name: "FFI_OutputSelectedRect",
							Args: []interface{}{page, left, top, right, bottom},
						})
					},
					FFI_SetCursor: func(cursorType enums.FXCT) {
						addToHistory(FormHistory{
							Name: "FFI_SetCursor",
							Args: []interface{}{cursorType},
						})
					},
					FFI_SetTimer: func(elapse int, timerFunc func(idEvent int)) int {
						addToHistory(FormHistory{
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
						addToHistory(FormHistory{
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
						addToHistory(FormHistory{
							Name: "FFI_GetLocalTime",
						})
						return structs.FPDF_SYSTEMTIME{}
					},
					FFI_OnChange: func() {
						addToHistory(FormHistory{
							Name: "FFI_OnChange",
						})
					},
					FFI_GetPage: func(document references.FPDF_DOCUMENT, index int) *references.FPDF_PAGE {
						addToHistory(FormHistory{
							Name: "FFI_GetPage",
							Args: []interface{}{document, index},
						})

						return nil
					},
					FFI_GetCurrentPage: func(document references.FPDF_DOCUMENT) *references.FPDF_PAGE {
						addToHistory(FormHistory{
							Name: "FFI_GetCurrentPage",
							Args: []interface{}{document},
						})
						return nil
					},
					FFI_GetRotation: func(page references.FPDF_PAGE) enums.FPDF_PAGE_ROTATION {
						addToHistory(FormHistory{
							Name: "FFI_GetRotation",
							Args: []interface{}{page},
						})
						return enums.FPDF_PAGE_ROTATION_NONE
					},
					FFI_ExecuteNamedAction: func(namedAction string) {
						addToHistory(FormHistory{
							Name: "FFI_ExecuteNamedAction",
							Args: []interface{}{namedAction},
						})
					},
					FFI_SetTextFieldFocus: func(value string, isFocus bool) {
						addToHistory(FormHistory{
							Name: "FFI_SetTextFieldFocus",
							Args: []interface{}{value, isFocus},
						})
					},
					FFI_DoURIAction: func(bsURI string) {
						addToHistory(FormHistory{
							Name: "FFI_DoURIAction",
							Args: []interface{}{bsURI},
						})
					},
					FFI_DoGoToAction: func(pageIndex int, zoomMode enums.FPDF_ZOOM_MODE, pos []float32) {
						addToHistory(FormHistory{
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
			//formattedHistory, _ := json.MarshalIndent(formHistory, "", "  ")

			//log.Println(formHistory)
			//log.Printf(string(formattedHistory))
		})

		When("is opened", func() {
			It("allows experimental form methods to be called", func() {
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

				FPDFAnnot_SetFocusableSubtypes, err := PdfiumInstance.FPDFAnnot_SetFocusableSubtypes(&requests.FPDFAnnot_SetFocusableSubtypes{
					FormHandle: formHandle,
					Subtypes:   []enums.FPDF_ANNOTATION_SUBTYPE{enums.FPDF_ANNOT_SUBTYPE_TEXT, enums.FPDF_ANNOT_SUBTYPE_LINE},
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_SetFocusableSubtypes).To(Equal(&responses.FPDFAnnot_SetFocusableSubtypes{}))

				FPDFAnnot_GetFocusableSubtypesCount, err := PdfiumInstance.FPDFAnnot_GetFocusableSubtypesCount(&requests.FPDFAnnot_GetFocusableSubtypesCount{
					FormHandle: formHandle,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFocusableSubtypesCount).To(Equal(&responses.FPDFAnnot_GetFocusableSubtypesCount{
					FocusableSubtypesCount: 2,
				}))

				FPDFAnnot_GetFocusableSubtypes, err := PdfiumInstance.FPDFAnnot_GetFocusableSubtypes(&requests.FPDFAnnot_GetFocusableSubtypes{
					FormHandle: formHandle,
				})
				Expect(err).To(BeNil())
				Expect(FPDFAnnot_GetFocusableSubtypes).To(Equal(&responses.FPDFAnnot_GetFocusableSubtypes{
					FocusableSubtypes: []enums.FPDF_ANNOTATION_SUBTYPE{enums.FPDF_ANNOT_SUBTYPE_TEXT, enums.FPDF_ANNOT_SUBTYPE_LINE},
				}))

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
		})
	})
})
