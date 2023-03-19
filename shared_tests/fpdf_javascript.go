//go:build pdfium_experimental
// +build pdfium_experimental

package shared_tests

import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_javascript", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	Context("no document is given", func() {
		It("returns an error when FPDFDoc_GetJavaScriptActionCount is called", func() {
			FPDFDoc_GetJavaScriptActionCount, err := PdfiumInstance.FPDFDoc_GetJavaScriptActionCount(&requests.FPDFDoc_GetJavaScriptActionCount{})
			Expect(err).To(MatchError("document not given"))
			Expect(FPDFDoc_GetJavaScriptActionCount).To(BeNil())
		})

		It("returns an error when FPDFDoc_GetJavaScriptAction is called", func() {
			FPDFDoc_GetJavaScriptAction, err := PdfiumInstance.FPDFDoc_GetJavaScriptAction(&requests.FPDFDoc_GetJavaScriptAction{})
			Expect(err).To(MatchError("document not given"))
			Expect(FPDFDoc_GetJavaScriptAction).To(BeNil())
		})
	})

	Context("no javascript action is given", func() {
		It("returns an error when FPDFDoc_CloseJavaScriptAction is called", func() {
			FPDFDoc_CloseJavaScriptAction, err := PdfiumInstance.FPDFDoc_CloseJavaScriptAction(&requests.FPDFDoc_CloseJavaScriptAction{})
			Expect(err).To(MatchError("javaScriptAction not given"))
			Expect(FPDFDoc_CloseJavaScriptAction).To(BeNil())
		})

		It("returns an error when FPDFJavaScriptAction_GetName is called", func() {
			FPDFJavaScriptAction_GetName, err := PdfiumInstance.FPDFJavaScriptAction_GetName(&requests.FPDFJavaScriptAction_GetName{})
			Expect(err).To(MatchError("javaScriptAction not given"))
			Expect(FPDFJavaScriptAction_GetName).To(BeNil())
		})

		It("returns an error when FPDFJavaScriptAction_GetScript is called", func() {
			FPDFJavaScriptAction_GetScript, err := PdfiumInstance.FPDFJavaScriptAction_GetScript(&requests.FPDFJavaScriptAction_GetScript{})
			Expect(err).To(MatchError("javaScriptAction not given"))
			Expect(FPDFJavaScriptAction_GetScript).To(BeNil())
		})
	})

	Context("a normal PDF file without javascript actions", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/test.pdf")
			Expect(err).To(BeNil())
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

		It("returns no an javascript action count of 0", func() {
			FPDFDoc_GetJavaScriptActionCount, err := PdfiumInstance.FPDFDoc_GetJavaScriptActionCount(&requests.FPDFDoc_GetJavaScriptActionCount{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDFDoc_GetJavaScriptActionCount).To(Equal(&responses.FPDFDoc_GetJavaScriptActionCount{
				JavaScriptActionCount: 0,
			}))
		})

		It("returns an error when trying to get an javascript action that isn't there", func() {
			FPDFDoc_GetJavaScriptAction, err := PdfiumInstance.FPDFDoc_GetJavaScriptAction(&requests.FPDFDoc_GetJavaScriptAction{
				Document: doc,
				Index:    0,
			})
			Expect(err).To(MatchError("could not get JavaScript Action"))
			Expect(FPDFDoc_GetJavaScriptAction).To(BeNil())
		})

		It("allows to add a new attachment", func() {
			FPDFDoc_AddAttachment, err := PdfiumInstance.FPDFDoc_AddAttachment(&requests.FPDFDoc_AddAttachment{
				Document: doc,
				Name:     "Attachment",
			})
			Expect(err).To(BeNil())
			Expect(FPDFDoc_AddAttachment).To(Not(BeNil()))
			Expect(FPDFDoc_AddAttachment.Attachment).To(Not(BeNil()))
		})
	})

	Context("a PDF file with javascript actions", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/js.pdf")
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

		It("returns the correct javascript action count", func() {
			FPDFDoc_GetJavaScriptActionCount, err := PdfiumInstance.FPDFDoc_GetJavaScriptActionCount(&requests.FPDFDoc_GetJavaScriptActionCount{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDFDoc_GetJavaScriptActionCount).To(Equal(&responses.FPDFDoc_GetJavaScriptActionCount{
				JavaScriptActionCount: 5,
			}))
		})

		It("returns the correct javascript actions", func() {
			GetJavaScriptActions, err := PdfiumInstance.GetJavaScriptActions(&requests.GetJavaScriptActions{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(GetJavaScriptActions).To(Equal(&responses.GetJavaScriptActions{
				JavaScriptActions: []responses.JavaScriptAction{
					{
						Name:   "normal",
						Script: "app.alert(\"ping\");",
					},
					{
						Name:   "no_type",
						Script: "app.alert(\"pong\");",
					},
				},
			}))
		})

		When("the first javascript action has been loaded", func() {
			var javaScriptAction references.FPDF_JAVASCRIPT_ACTION
			BeforeEach(func() {
				FPDFDoc_GetJavaScriptAction, err := PdfiumInstance.FPDFDoc_GetJavaScriptAction(&requests.FPDFDoc_GetJavaScriptAction{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(FPDFDoc_GetJavaScriptAction).To(Not(BeNil()))
				Expect(FPDFDoc_GetJavaScriptAction.JavaScriptAction).To(Not(BeNil()))
				javaScriptAction = FPDFDoc_GetJavaScriptAction.JavaScriptAction
			})

			AfterEach(func() {
				FPDFDoc_CloseJavaScriptAction, err := PdfiumInstance.FPDFDoc_CloseJavaScriptAction(&requests.FPDFDoc_CloseJavaScriptAction{
					JavaScriptAction: javaScriptAction,
				})
				Expect(err).To(BeNil())
				Expect(FPDFDoc_CloseJavaScriptAction).To(Not(BeNil()))
			})

			It("returns the correct name of the javascript action", func() {
				FPDFJavaScriptAction_GetName, err := PdfiumInstance.FPDFJavaScriptAction_GetName(&requests.FPDFJavaScriptAction_GetName{
					JavaScriptAction: javaScriptAction,
				})
				Expect(err).To(BeNil())
				Expect(FPDFJavaScriptAction_GetName).To(Equal(&responses.FPDFJavaScriptAction_GetName{
					Name: "normal",
				}))
			})

			It("returns the correct script of the javascript action", func() {
				FPDFJavaScriptAction_GetScript, err := PdfiumInstance.FPDFJavaScriptAction_GetScript(&requests.FPDFJavaScriptAction_GetScript{
					JavaScriptAction: javaScriptAction,
				})
				Expect(err).To(BeNil())
				Expect(FPDFJavaScriptAction_GetScript).To(Equal(&responses.FPDFJavaScriptAction_GetScript{
					Script: "app.alert(\"ping\");",
				}))
			})
		})
	})
})
