package shared_tests

import (
	"image"
	"os"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_flatten", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when flattening a pdf page", func() {
				pageCount, err := PdfiumInstance.FPDFPage_Flatten(&requests.FPDFPage_Flatten{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Index: 0,
						},
					},
				})
				Expect(err).To(MatchError("document not given"))
				Expect(pageCount).To(BeNil())
			})
		})
	})

	// PDFium parses page content when a page is loaded and never again, so
	// flattening only becomes visible on a freshly loaded page. go-pdfium
	// reloads the page behind the reference after a successful flatten, see
	// https://github.com/klippa-app/go-pdfium/issues/199.
	Context("a PDF file with annotations", func() {
		var doc references.FPDF_DOCUMENT

		renderPage := func(page requests.Page, flags enums.FPDF_RENDER_FLAG) *image.RGBA {
			render, err := PdfiumInstance.RenderPageInDPI(&requests.RenderPageInDPI{
				Page:        page,
				DPI:         72,
				RenderFlags: flags,
			})
			Expect(err).To(BeNil())
			// Copy the pixels so that the WebAssembly render buffer can be released.
			img := image.NewRGBA(render.Result.Image.Rect)
			copy(img.Pix, render.Result.Image.Pix)
			render.Cleanup()
			return img
		}

		BeforeEach(func() {
			pdfData, err := os.ReadFile(TestDataPath + "/testdata/annotation_highlight_square_with_ap.pdf")
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

		It("renders the flattened annotations through the same page reference", func() {
			page, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{
				Document: doc,
				Index:    0,
			})
			Expect(err).To(BeNil())
			pageRef := requests.Page{ByReference: &page.Page}

			// The annotations are drawn by the annotation renderer before flattening.
			withAnnotations := renderPage(pageRef, enums.FPDF_RENDER_FLAG_ANNOT)
			withoutAnnotations := renderPage(pageRef, 0)
			Expect(withAnnotations.Pix).ToNot(Equal(withoutAnnotations.Pix))

			flatten, err := PdfiumInstance.FPDFPage_Flatten(&requests.FPDFPage_Flatten{
				Page:  pageRef,
				Usage: requests.FPDFPage_FlattenUsageNormalDisplay,
			})
			Expect(err).To(BeNil())
			Expect(flatten).To(Equal(&responses.FPDFPage_Flatten{
				Page:   0,
				Result: responses.FPDFPage_FlattenResultSuccess,
			}))

			// After flattening they are part of the page content, so they show
			// up without the annotation flag, through the very same reference.
			Expect(renderPage(pageRef, 0).Pix).To(Equal(withAnnotations.Pix))

			// And a freshly loaded page gives the same result.
			freshPage, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{
				Document: doc,
				Index:    0,
			})
			Expect(err).To(BeNil())
			Expect(renderPage(requests.Page{ByReference: &freshPage.Page}, 0).Pix).To(Equal(withAnnotations.Pix))

			_, err = PdfiumInstance.FPDF_ClosePage(&requests.FPDF_ClosePage{Page: freshPage.Page})
			Expect(err).To(BeNil())
			_, err = PdfiumInstance.FPDF_ClosePage(&requests.FPDF_ClosePage{Page: page.Page})
			Expect(err).To(BeNil())
		})

		It("renders the flattened annotations through the same page index", func() {
			pageByIndex := requests.Page{ByIndex: &requests.PageByIndex{Document: doc, Index: 0}}
			withAnnotations := renderPage(pageByIndex, enums.FPDF_RENDER_FLAG_ANNOT)
			Expect(renderPage(pageByIndex, 0).Pix).ToNot(Equal(withAnnotations.Pix))

			flatten, err := PdfiumInstance.FPDFPage_Flatten(&requests.FPDFPage_Flatten{
				Page:  pageByIndex,
				Usage: requests.FPDFPage_FlattenUsageNormalDisplay,
			})
			Expect(err).To(BeNil())
			Expect(flatten.Result).To(Equal(responses.FPDFPage_FlattenResultSuccess))

			Expect(renderPage(pageByIndex, 0).Pix).To(Equal(withAnnotations.Pix))
		})

		It("reports nothing to do when flattening the page again", func() {
			pageByIndex := requests.Page{ByIndex: &requests.PageByIndex{Document: doc, Index: 0}}
			flatten, err := PdfiumInstance.FPDFPage_Flatten(&requests.FPDFPage_Flatten{
				Page:  pageByIndex,
				Usage: requests.FPDFPage_FlattenUsageNormalDisplay,
			})
			Expect(err).To(BeNil())
			Expect(flatten.Result).To(Equal(responses.FPDFPage_FlattenResultSuccess))

			flatten, err = PdfiumInstance.FPDFPage_Flatten(&requests.FPDFPage_Flatten{
				Page:  pageByIndex,
				Usage: requests.FPDFPage_FlattenUsageNormalDisplay,
			})
			Expect(err).To(BeNil())
			Expect(flatten.Result).To(Equal(responses.FPDFPage_FlattenResultNothingToDo))
		})

		It("keeps text pages loaded before flattening usable", func() {
			page, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{
				Document: doc,
				Index:    0,
			})
			Expect(err).To(BeNil())
			pageRef := requests.Page{ByReference: &page.Page}

			textPage, err := PdfiumInstance.FPDFText_LoadPage(&requests.FPDFText_LoadPage{Page: pageRef})
			Expect(err).To(BeNil())

			flatten, err := PdfiumInstance.FPDFPage_Flatten(&requests.FPDFPage_Flatten{
				Page:  pageRef,
				Usage: requests.FPDFPage_FlattenUsageNormalDisplay,
			})
			Expect(err).To(BeNil())
			Expect(flatten.Result).To(Equal(responses.FPDFPage_FlattenResultSuccess))

			// The old text page still belongs to a live (stale) page object.
			countChars, err := PdfiumInstance.FPDFText_CountChars(&requests.FPDFText_CountChars{TextPage: textPage.TextPage})
			Expect(err).To(BeNil())
			Expect(countChars.Count).To(BeNumerically(">", 0))

			_, err = PdfiumInstance.FPDFText_ClosePage(&requests.FPDFText_ClosePage{TextPage: textPage.TextPage})
			Expect(err).To(BeNil())
			_, err = PdfiumInstance.FPDF_ClosePage(&requests.FPDF_ClosePage{Page: page.Page})
			Expect(err).To(BeNil())
		})
	})

	Context("a PDF file with a filled form and a form fill environment", func() {
		var doc references.FPDF_DOCUMENT
		var formHandle references.FPDF_FORMHANDLE

		BeforeEach(func() {
			if TestType == "multi" {
				Skip("Form filling is not supported on multi-threaded usage")
			}

			pdfData, err := os.ReadFile(TestDataPath + "/testdata/text_form_filled.pdf")
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
					FFI_SetTimer:           func(elapse int, timerFunc func(idEvent int)) int { return 0 },
					FFI_KillTimer:          func(timerID int) {},
					FFI_GetLocalTime:       func() structs.FPDF_SYSTEMTIME { return structs.FPDF_SYSTEMTIME{} },
					FFI_OnChange:           func() {},
					FFI_GetPage: func(document references.FPDF_DOCUMENT, index int) *references.FPDF_PAGE {
						return nil
					},
					FFI_GetCurrentPage: func(document references.FPDF_DOCUMENT) *references.FPDF_PAGE {
						return nil
					},
					FFI_GetRotation:        func(page references.FPDF_PAGE) enums.FPDF_PAGE_ROTATION { return 0 },
					FFI_ExecuteNamedAction: func(namedAction string) {},
				},
			})
			Expect(err).To(BeNil())
			formHandle = FPDFDOC_InitFormFillEnvironment.FormHandle
		})

		AfterEach(func() {
			if TestType == "multi" {
				Skip("Form filling is not supported on multi-threaded usage")
			}

			_, err := PdfiumInstance.FPDFDOC_ExitFormFillEnvironment(&requests.FPDFDOC_ExitFormFillEnvironment{FormHandle: formHandle})
			Expect(err).To(BeNil())
			_, err = PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc})
			Expect(err).To(BeNil())
		})

		It("flattens a page that is loaded in the form fill environment", func() {
			page, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{Document: doc, Index: 0})
			Expect(err).To(BeNil())
			pageRef := requests.Page{ByReference: &page.Page}

			_, err = PdfiumInstance.FORM_OnAfterLoadPage(&requests.FORM_OnAfterLoadPage{Page: pageRef, FormHandle: formHandle})
			Expect(err).To(BeNil())

			textBefore, err := PdfiumInstance.GetPageText(&requests.GetPageText{Page: pageRef})
			Expect(err).To(BeNil())

			flatten, err := PdfiumInstance.FPDFPage_Flatten(&requests.FPDFPage_Flatten{
				Page:  pageRef,
				Usage: requests.FPDFPage_FlattenUsageNormalDisplay,
			})
			Expect(err).To(BeNil())
			Expect(flatten.Result).To(Equal(responses.FPDFPage_FlattenResultSuccess))

			// The field value became page content and is now part of the page text.
			textAfter, err := PdfiumInstance.GetPageText(&requests.GetPageText{Page: pageRef})
			Expect(err).To(BeNil())
			Expect(textAfter.Text).ToNot(Equal(textBefore.Text))

			// The reloaded page is registered with the form fill environment, so
			// the regular teardown sequence works on it.
			_, err = PdfiumInstance.FORM_OnBeforeClosePage(&requests.FORM_OnBeforeClosePage{Page: pageRef, FormHandle: formHandle})
			Expect(err).To(BeNil())
			_, err = PdfiumInstance.FPDF_ClosePage(&requests.FPDF_ClosePage{Page: page.Page})
			Expect(err).To(BeNil())
		})
	})
})
