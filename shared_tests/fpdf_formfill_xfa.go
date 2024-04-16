//go:build pdfium_xfa
// +build pdfium_xfa

package shared_tests

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"time"
	"unsafe"

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
			Skip("Form filling is not supported on multi-threaded usage")
		}
		if TestType == "webassembly" {
			Skip("XFA Form filling is not supported on webassembly usage")
		}
		Locker.Lock()
	})

	AfterEach(func() {
		if TestType == "multi" {
			Skip("Form filling is not supported on multi-threaded usage")
		}
		if TestType == "webassembly" {
			Skip("XFA Form filling is not supported on webassembly usage")
		}
		Locker.Unlock()
	})

	type FormHistory struct {
		Name string
		Args []interface{}
	}

	type FormTicker struct {
		Timer *time.Ticker
		Done  chan bool
	}

	Context("a XFA PDF file with an email form", func() {
		var doc references.FPDF_DOCUMENT
		var page references.FPDF_PAGE
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

		renderFormImage := func(page references.FPDF_PAGE, title string) {
			FPDFBitmap_FillRect, err := PdfiumInstance.FPDFBitmap_FillRect(&requests.FPDFBitmap_FillRect{
				Bitmap: bitmap,
				Color:  0xFFFFFFFF,
				Left:   0,
				Top:    0,
				Width:  900,
				Height: 1164,
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
				SizeX:  900,
				SizeY:  1164,
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
				SizeX:  900,
				SizeY:  1164,
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

			ioutil.WriteFile(TestDataPath+"/testdata/"+fmt.Sprintf("render_fpdf_formfill_xfa_%s-%d.jpg", title, renderCount), imgBuf.Bytes(), 0777)
			renderCount++
			//log.Println("did render")
		}

		BeforeEach(func() {
			formHistory = []FormHistory{}
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/xfa/email_recommended.pdf")
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
					FFI_DisplayCaret: func(page references.FPDF_PAGE, bVisible bool, left, top, right, bottom float64) {
						addToHistory(FormHistory{
							Name: "FFI_DisplayCaret",
							Args: []interface{}{page, bVisible, left, top, right, bottom},
						})
					},
					FFI_GetCurrentPageIndex: func(document references.FPDF_DOCUMENT) int {
						addToHistory(FormHistory{
							Name: "FFI_GetCurrentPageIndex",
							Args: []interface{}{document},
						})
						return 0
					},
					FFI_SetCurrentPage: func(document references.FPDF_DOCUMENT, iCurPage int) {
						addToHistory(FormHistory{
							Name: "FFI_SetCurrentPage",
							Args: []interface{}{document, iCurPage},
						})
					},
					FFI_GotoURL: func(document references.FPDF_DOCUMENT, url string) {
						addToHistory(FormHistory{
							Name: "FFI_GotoURL",
							Args: []interface{}{document, url},
						})
					},
					FFI_GetPageViewRect: func(page references.FPDF_PAGE) (left, top, right, bottom float64) {
						addToHistory(FormHistory{
							Name: "FFI_GetPageViewRect",
							Args: []interface{}{page},
						})
						return 0, 0, 0, 0
					},
					FFI_PageEvent: func(page_count int, event_type enums.FXFA_PAGEVIEWEVENT) {
						addToHistory(FormHistory{
							Name: "FFI_PageEvent",
							Args: []interface{}{page_count, event_type},
						})
					},
					FFI_PopupMenu: func(page references.FPDF_PAGE, menuFlag int, x, y float32) bool {
						addToHistory(FormHistory{
							Name: "FFI_PopupMenu",
							Args: []interface{}{page, menuFlag, x, y},
						})
						return false
					},
					FFI_OpenFile: func(fileFlag enums.FXFA_SAVEAS, url, mode string) *structs.FPDF_FILEHANDLER {
						addToHistory(FormHistory{
							Name: "FFI_OpenFile",
							Args: []interface{}{fileFlag, url, mode},
						})
						return nil
					},
					FFI_EmailTo: func(fileHandler *structs.FPDF_FILEHANDLER, to, subject, cc, bcc, msg string) {
						addToHistory(FormHistory{
							Name: "FFI_EmailTo",
							Args: []interface{}{fileHandler, to, subject, cc, bcc, msg},
						})
					},
					FFI_UploadTo: func(fileHandler *structs.FPDF_FILEHANDLER, fileFlag enums.FXFA_SAVEAS, uploadTo string) {
						addToHistory(FormHistory{
							Name: "FFI_UploadTo",
							Args: []interface{}{fileHandler, fileFlag, uploadTo},
						})
					},
					FFI_GetPlatform: func() string {
						addToHistory(FormHistory{
							Name: "FFI_GetPlatform",
							Args: []interface{}{},
						})
						return ""
					},
					FFI_GetLanguage: func() string {
						addToHistory(FormHistory{
							Name: "FFI_GetLanguage",
							Args: []interface{}{},
						})
						return ""
					},
					FFI_DownloadFromURL: func(url string) *structs.FPDF_FILEHANDLER {
						addToHistory(FormHistory{
							Name: "FFI_DownloadFromURL",
							Args: []interface{}{url},
						})
						return nil
					},
					FFI_PostRequestURL: func(url, data, contentType, encode, header string, response references.FPDF_BSTR) bool {
						addToHistory(FormHistory{
							Name: "FFI_PostRequestURL",
							Args: []interface{}{url, data, contentType, encode, header, response},
						})
						return false
					},
					FFI_PutRequestURL: func(url, data, encode string) bool {
						addToHistory(FormHistory{
							Name: "FFI_PutRequestURL",
							Args: []interface{}{url, data, encode},
						})
						return false
					},
					FFI_OnFocusChange: func(annot references.FPDF_ANNOTATION, page_index int) {
						addToHistory(FormHistory{
							Name: "FFI_OnFocusChange",
							Args: []interface{}{annot, page_index},
						})
					},
					FFI_DoURIActionWithKeyboardModifier: func(uri string, modifiers int) {
						addToHistory(FormHistory{
							Name: "FFI_DoURIActionWithKeyboardModifier",
							Args: []interface{}{uri, modifiers},
						})
					},
					JsPlatform: &structs.IPDF_JSPLATFORM{
						App_alert: func(msg, title string, nButton enums.JSPLATFORM_ALERT_BUTTON, nIcon enums.JSPLATFORM_ALERT_ICON) int {
							addToHistory(FormHistory{
								Name: "JsPlatform_app_alert",
								Args: []interface{}{msg, title, nButton, nIcon},
							})
							return 0
						},
						App_beep: func(nType enums.JSPLATFORM_BEEP) {
							addToHistory(FormHistory{
								Name: "JsPlatform_app_beep",
								Args: []interface{}{nType},
							})
						},
						App_response: func(question, title, defaultValue, cLabel string, bPassword bool) string {
							addToHistory(FormHistory{
								Name: "JsPlatform_app_response",
								Args: []interface{}{question, title, defaultValue, cLabel, bPassword},
							})
							return ""
						},
						Doc_getFilePath: func() string {
							addToHistory(FormHistory{
								Name: "JsPlatform_Doc_getFilePath",
								Args: []interface{}{},
							})
							return ""
						},
						Doc_mail: func(mailData []byte, bUI bool, to, subject, cc, bcc, msg string) {
							addToHistory(FormHistory{
								Name: "JsPlatform_Doc_mail",
								Args: []interface{}{mailData, bUI, to, subject, cc, bcc, msg},
							})
						},
						Doc_print: func(bUI bool, nStart, nEnd int, bSilent, bShrinkToFit, bPrintAsImage, bReverse, bAnnotations bool) {
							addToHistory(FormHistory{
								Name: "JsPlatform_Doc_print",
								Args: []interface{}{bUI, nStart, nEnd, bSilent, bShrinkToFit, bPrintAsImage, bReverse, bAnnotations},
							})
						},
						Doc_submitForm: func(formData []byte, url string) {
							addToHistory(FormHistory{
								Name: "JsPlatform_Doc_submitForm",
								Args: []interface{}{formData, url},
							})
						},
						Doc_gotoPage: func(nPageNum int) {
							addToHistory(FormHistory{
								Name: "JsPlatform_Doc_gotoPage",
								Args: []interface{}{nPageNum},
							})
						},
						Field_browse: func() string {
							addToHistory(FormHistory{
								Name: "JsPlatform_Field_browse",
								Args: []interface{}{},
							})
							return ""
						},
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(FPDFDOC_InitFormFillEnvironment).ToNot(BeNil())
			Expect(FPDFDOC_InitFormFillEnvironment.FormHandle).ToNot(BeEmpty())
			formHandle = FPDFDOC_InitFormFillEnvironment.FormHandle

			_, err = PdfiumInstance.FPDF_LoadXFA(&requests.FPDF_LoadXFA{
				Document: doc,
			})
			Expect(err).To(BeNil())

			FPDF_LoadPage, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{
				Document: doc,
				Index:    0,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_LoadPage).ToNot(BeNil())
			Expect(FPDF_LoadPage.Page).ToNot(BeEmpty())
			page = FPDF_LoadPage.Page

			FORM_OnAfterLoadPage, err := PdfiumInstance.FORM_OnAfterLoadPage(&requests.FORM_OnAfterLoadPage{
				Page: requests.Page{
					ByReference: &page,
				},
				FormHandle: formHandle,
			})
			Expect(err).To(BeNil())
			Expect(FORM_OnAfterLoadPage).To(Equal(&responses.FORM_OnAfterLoadPage{}))

			width := 900
			height := 1164
			stride := width * 4

			fileSize := stride * height
			buffer := make([]byte, fileSize)
			pointer := unsafe.Pointer(&buffer[0])

			renderCount = 0
			img = image.NewRGBA(image.Rect(0, 0, 900, 1164))
			img.Pix = buffer
			FPDFBitmap_CreateEx, err := PdfiumInstance.FPDFBitmap_CreateEx(&requests.FPDFBitmap_CreateEx{
				Width:   900,
				Height:  1164,
				Format:  enums.FPDF_BITMAP_FORMAT_BGRA,
				Pointer: pointer,
				Stride:  img.Stride,
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
			It("allows to invoke the first tab on the page", func() {
				renderFormImage(page, "normal-tab")

				_, err := PdfiumInstance.FORM_OnKeyDown(&requests.FORM_OnKeyDown{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &page,
					},
					NKeyCode: enums.FWL_VKEY_Tab,
				})
				Expect(err).To(BeNil())

				renderFormImage(page, "normal-tab")
			})

			It("allows to invoke the first shift-tab on the page", func() {
				renderFormImage(page, "shift-tab")
				_, err := PdfiumInstance.FORM_OnKeyDown(&requests.FORM_OnKeyDown{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &page,
					},
					NKeyCode: enums.FWL_VKEY_Tab,
					Modifier: enums.FWL_EVENTFLAG_ShiftKey,
				})
				Expect(err).To(BeNil())

				renderFormImage(page, "shift-tab")
			})

			It("allows to continuously tab on the page", func() {
				renderFormImage(page, "continuously-tab")

				// First tab
				_, err := PdfiumInstance.FORM_OnKeyDown(&requests.FORM_OnKeyDown{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &page,
					},
					NKeyCode: enums.FWL_VKEY_Tab,
				})
				Expect(err).To(BeNil())
				renderFormImage(page, "continuously-tab")

				for i := 0; i < 9; i++ {
					_, err = PdfiumInstance.FORM_OnKeyDown(&requests.FORM_OnKeyDown{
						FormHandle: formHandle,
						Page: requests.Page{
							ByReference: &page,
						},
						NKeyCode: enums.FWL_VKEY_Tab,
					})
					Expect(err).To(BeNil())
					renderFormImage(page, "continuously-tab")
				}

				// Tab should not be handled as the last annotation of the page is in focus.
				_, err = PdfiumInstance.FORM_OnKeyDown(&requests.FORM_OnKeyDown{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &page,
					},
					NKeyCode: enums.FWL_VKEY_Tab,
				})
				Expect(err).To(Not(BeNil()))
			})

			It("allows to continuously shift tab on the page", func() {
				renderFormImage(page, "continuously-shift-tab")

				// First tab
				_, err := PdfiumInstance.FORM_OnKeyDown(&requests.FORM_OnKeyDown{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &page,
					},
					NKeyCode: enums.FWL_VKEY_Tab,
					Modifier: enums.FWL_EVENTFLAG_ShiftKey,
				})
				Expect(err).To(BeNil())
				renderFormImage(page, "continuously-shift-tab")

				for i := 0; i < 9; i++ {
					_, err = PdfiumInstance.FORM_OnKeyDown(&requests.FORM_OnKeyDown{
						FormHandle: formHandle,
						Page: requests.Page{
							ByReference: &page,
						},
						NKeyCode: enums.FWL_VKEY_Tab,
						Modifier: enums.FWL_EVENTFLAG_ShiftKey,
					})
					Expect(err).To(BeNil())
					renderFormImage(page, "continuously-shift-tab")
				}

				// Tab should not be handled as the last annotation of the page is in focus.
				_, err = PdfiumInstance.FORM_OnKeyDown(&requests.FORM_OnKeyDown{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &page,
					},
					NKeyCode: enums.FWL_VKEY_Tab,
					Modifier: enums.FWL_EVENTFLAG_ShiftKey,
				})
				Expect(err).To(Not(BeNil()))
			})

			It("allows to select text with a left mouse click", func() {
				renderFormImage(page, "select-text-left-mouse-click")

				// Focus field.
				_, err := PdfiumInstance.FORM_OnLButtonDown(&requests.FORM_OnLButtonDown{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &page,
					},
					PageX: 115,
					PageY: 58,
				})
				Expect(err).To(BeNil())
				renderFormImage(page, "select-text-left-mouse-click")

				// Write text.
				for i := 0; i < 10; i++ {
					_, err := PdfiumInstance.FORM_OnChar(&requests.FORM_OnChar{
						FormHandle: formHandle,
						Page: requests.Page{
							ByReference: &page,
						},
						NChar: 'a' + i,
					})
					Expect(err).To(BeNil())
					renderFormImage(page, "select-text-left-mouse-click")
				}

				// Set mouse position.
				_, err = PdfiumInstance.FORM_OnLButtonDown(&requests.FORM_OnLButtonDown{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &page,
					},
					PageX: 128,
					PageY: 58,
				})
				Expect(err).To(BeNil())
				renderFormImage(page, "select-text-left-mouse-click")

				// Select text.
				_, err = PdfiumInstance.FORM_OnLButtonDown(&requests.FORM_OnLButtonDown{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &page,
					},
					PageX:    152,
					PageY:    58,
					Modifier: int(enums.FWL_EVENTFLAG_ShiftKey),
				})
				Expect(err).To(BeNil())
				renderFormImage(page, "select-text-left-mouse-click")

				// Get selected text.
				FORM_GetSelectedText, err := PdfiumInstance.FORM_GetSelectedText(&requests.FORM_GetSelectedText{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &page,
					},
				})
				Expect(err).To(BeNil())
				Expect(FORM_GetSelectedText).To(Equal(&responses.FORM_GetSelectedText{
					SelectedText: "defgh",
				}))
			})

			It("allows to do a drag mouse selection", func() {
				renderFormImage(page, "drag-mouse-selection")

				// Focus field.
				_, err := PdfiumInstance.FORM_OnLButtonDown(&requests.FORM_OnLButtonDown{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &page,
					},
					PageX: 115,
					PageY: 58,
				})
				Expect(err).To(BeNil())
				renderFormImage(page, "drag-mouse-selection")

				// Write text.
				for i := 0; i < 10; i++ {
					_, err := PdfiumInstance.FORM_OnChar(&requests.FORM_OnChar{
						FormHandle: formHandle,
						Page: requests.Page{
							ByReference: &page,
						},
						NChar: 'a' + i,
					})
					Expect(err).To(BeNil())
					renderFormImage(page, "drag-mouse-selection")
				}

				// Set mouse position.
				_, err = PdfiumInstance.FORM_OnLButtonDown(&requests.FORM_OnLButtonDown{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &page,
					},
					PageX: 128,
					PageY: 58,
				})
				Expect(err).To(BeNil())
				renderFormImage(page, "drag-mouse-selection")

				// Select text.
				_, err = PdfiumInstance.FORM_OnMouseMove(&requests.FORM_OnMouseMove{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &page,
					},
					PageX:    152,
					PageY:    58,
					Modifier: int(enums.FWL_EVENTFLAG_ShiftKey),
				})
				Expect(err).To(BeNil())
				renderFormImage(page, "drag-mouse-selection")

				// Get selected text.
				FORM_GetSelectedText, err := PdfiumInstance.FORM_GetSelectedText(&requests.FORM_GetSelectedText{
					FormHandle: formHandle,
					Page: requests.Page{
						ByReference: &page,
					},
				})
				Expect(err).To(BeNil())
				Expect(FORM_GetSelectedText).To(Equal(&responses.FORM_GetSelectedText{
					SelectedText: "defgh",
				}))
			})
		})
	})
})
