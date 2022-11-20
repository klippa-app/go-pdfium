package shared_tests

import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("page", func() {
	BeforeEach(func() {
		Locker.Lock()

		if TestType == "webassembly" {
			// @todo: remove me when implemented.
			Skip("This test is skipped on Webassembly")
		}
	})

	AfterEach(func() {
		Locker.Unlock()

		if TestType == "webassembly" {
			// @todo: remove me when implemented.
			Skip("This test is skipped on Webassembly")
		}
	})

	Context("a normal PDF file", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/test.pdf")
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

		When("is opened", func() {
			Context("when an invalid page is given", func() {
				Context("GetPageRotation()", func() {
					It("returns an error", func() {
						pageRotation, err := PdfiumInstance.FPDFPage_GetRotation(&requests.FPDFPage_GetRotation{
							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    1,
								},
							},
						})
						Expect(err).To(MatchError(errors.ErrPage.Error()))
						Expect(pageRotation).To(BeNil())
					})
				})

				Context("GetPageTransparency()", func() {
					It("returns an error", func() {
						pageTransparency, err := PdfiumInstance.FPDFPage_HasTransparency(&requests.FPDFPage_HasTransparency{
							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    1,
								},
							},
						})
						Expect(err).To(MatchError(errors.ErrPage.Error()))
						Expect(pageTransparency).To(BeNil())
					})
				})

				Context("FlattenPage()", func() {
					It("returns an error", func() {
						flattenedPage, err := PdfiumInstance.FPDFPage_Flatten(&requests.FPDFPage_Flatten{
							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    1,
								},
							},
							Usage: requests.FPDFPage_FlattenUsageNormalDisplay,
						})
						Expect(err).To(MatchError(errors.ErrPage.Error()))
						Expect(flattenedPage).To(BeNil())
					})
				})
			})

			Context("when the page rotation is requested", func() {
				It("returns the correct rotation", func() {
					rotation, err := PdfiumInstance.FPDFPage_GetRotation(&requests.FPDFPage_GetRotation{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})
					Expect(err).To(BeNil())
					Expect(rotation).To(Equal(&responses.FPDFPage_GetRotation{}))
				})
			})

			Context("when the page transparency is requested", func() {
				It("returns the correct transparency", func() {
					pageTransparency, err := PdfiumInstance.FPDFPage_HasTransparency(&requests.FPDFPage_HasTransparency{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})
					Expect(err).To(BeNil())
					Expect(pageTransparency).To(Equal(&responses.FPDFPage_HasTransparency{
						HasTransparency: false,
					}))
				})
			})

			Context("when the page flattening is requested", func() {
				It("returns that the page does not need to be flattened", func() {
					pageFlattenResult, err := PdfiumInstance.FPDFPage_Flatten(&requests.FPDFPage_Flatten{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
						Usage: requests.FPDFPage_FlattenUsageNormalDisplay,
					})
					Expect(err).To(BeNil())
					Expect(pageFlattenResult).To(Equal(&responses.FPDFPage_Flatten{
						Result: responses.FPDFPage_FlattenResultNothingToDo,
					}))
				})
			})
		})
	})

	Context("a PDF file that uses an alpha channel", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/alpha_channel.pdf")
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

		When("the page transparency is requested", func() {
			It("returns the correct transparency", func() {
				pageTransparency, err := PdfiumInstance.FPDFPage_HasTransparency(&requests.FPDFPage_HasTransparency{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(pageTransparency).To(Equal(&responses.FPDFPage_HasTransparency{
					HasTransparency: true,
				}))
			})
		})

		When("the page transparency is requested by a wrong page reference", func() {
			It("loadPage returns an error", func() {
				fakeReference := references.FPDF_PAGE("")
				pageTransparency, err := PdfiumInstance.FPDFPage_HasTransparency(&requests.FPDFPage_HasTransparency{
					Page: requests.Page{
						ByReference: &fakeReference,
					},
				})
				Expect(err).To(Not(BeNil()))
				Expect(pageTransparency).To(BeNil())
			})
		})

		When("the page transparency is requested by page reference", func() {
			It("returns the correct transparency", func() {
				FPDF_LoadPage, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_LoadPage).To(Not(BeNil()))
				Expect(FPDF_LoadPage.Page).To(Not(BeNil()))

				pageTransparency, err := PdfiumInstance.FPDFPage_HasTransparency(&requests.FPDFPage_HasTransparency{
					Page: requests.Page{
						ByReference: &FPDF_LoadPage.Page,
					},
				})
				Expect(err).To(BeNil())
				Expect(pageTransparency).To(Equal(&responses.FPDFPage_HasTransparency{
					HasTransparency: true,
				}))
			})
		})
	})
})
