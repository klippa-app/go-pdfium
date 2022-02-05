package shared_tests

import (
	"os"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func RunfpdfViewTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("fpdfview", func() {
		It("allows for the last error to be fetched", func() {
			FPDF_GetLastError, err := pdfiumContainer.FPDF_GetLastError(&requests.FPDF_GetLastError{})
			Expect(err).To(BeNil())
			Expect(FPDF_GetLastError).To(Not(BeNil()))
		})

		It("allows for the sandbox policy to be enabled", func() {
			FPDF_SetSandBoxPolicy, err := pdfiumContainer.FPDF_SetSandBoxPolicy(&requests.FPDF_SetSandBoxPolicy{
				Policy: requests.FPDF_SetSandBoxPolicyPolicyMachinetimeAccess,
				Enable: true,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_SetSandBoxPolicy).To(Equal(&responses.FPDF_SetSandBoxPolicy{}))
		})

		It("allows for the sandbox policy to be disabled", func() {
			FPDF_SetSandBoxPolicy, err := pdfiumContainer.FPDF_SetSandBoxPolicy(&requests.FPDF_SetSandBoxPolicy{
				Policy: requests.FPDF_SetSandBoxPolicyPolicyMachinetimeAccess,
				Enable: false,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_SetSandBoxPolicy).To(Equal(&responses.FPDF_SetSandBoxPolicy{}))
		})

		Context("no document", func() {
			When("is opened", func() {
				It("returns an error when getting the pdf version", func() {
					pageCount, err := pdfiumContainer.FPDF_GetFileVersion(&requests.FPDF_GetFileVersion{})
					Expect(err).To(MatchError("document not given"))
					Expect(pageCount).To(BeNil())
				})

				It("returns an error when getting the doc permissions", func() {
					pageCount, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{})
					Expect(err).To(MatchError("document not given"))
					Expect(pageCount).To(BeNil())
				})

				It("returns an error when getting the doc revision number of security handler", func() {
					pageCount, err := pdfiumContainer.FPDF_GetSecurityHandlerRevision(&requests.FPDF_GetSecurityHandlerRevision{})
					Expect(err).To(MatchError("document not given"))
					Expect(pageCount).To(BeNil())
				})

				It("returns an error when getting the page count", func() {
					pageCount, err := pdfiumContainer.FPDF_GetPageCount(&requests.FPDF_GetPageCount{})
					Expect(err).To(MatchError("document not given"))
					Expect(pageCount).To(BeNil())
				})

				It("returns an error when calling FPDF_GetPageSizeByIndex", func() {
					FPDF_GetPageSizeByIndex, err := pdfiumContainer.FPDF_GetPageSizeByIndex(&requests.FPDF_GetPageSizeByIndex{})
					Expect(err).To(MatchError("document not given"))
					Expect(FPDF_GetPageSizeByIndex).To(BeNil())
				})
			})
		})

		Context("a normal PDF file with 1 page", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				file, err := os.Open(testsPath + "/testdata/test.pdf")
				Expect(err).To(BeNil())

				fileStat, err := file.Stat()
				Expect(err).To(BeNil())

				newDoc, err := pdfiumContainer.FPDF_LoadCustomDocument(&requests.FPDF_LoadCustomDocument{
					Reader: file,
					Size:   fileStat.Size(),
				})
				Expect(err).To(BeNil())

				doc = newDoc.Document
			})

			AfterEach(func() {
				FPDF_CloseDocument, err := pdfiumContainer.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_CloseDocument).To(Not(BeNil()))
			})

			When("is opened", func() {
				It("returns the correct page count", func() {
					FPDF_GetPageCount, err := pdfiumContainer.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_GetPageCount).To(Equal(&responses.FPDF_GetPageCount{
						PageCount: 1,
					}))
				})

				It("returns an error when FPDF_LoadPage is called without a document", func() {
					FPDF_LoadPage, err := pdfiumContainer.FPDF_LoadPage(&requests.FPDF_LoadPage{})
					Expect(err).To(MatchError("document not given"))
					Expect(FPDF_LoadPage).To(BeNil())
				})

				It("returns an error when FPDF_LoadPage is called with an invalid page", func() {
					FPDF_LoadPage, err := pdfiumContainer.FPDF_LoadPage(&requests.FPDF_LoadPage{
						Document: doc,
						Index:    23,
					})
					Expect(err).To(MatchError("incorrect page"))
					Expect(FPDF_LoadPage).To(BeNil())
				})

				It("returns an error when FPDF_ClosePage is called without a page", func() {
					FPDF_ClosePage, err := pdfiumContainer.FPDF_ClosePage(&requests.FPDF_ClosePage{})
					Expect(err).To(MatchError("page not given"))
					Expect(FPDF_ClosePage).To(BeNil())
				})

				It("allows for a page to be loaded and closed", func() {
					FPDF_LoadPage, err := pdfiumContainer.FPDF_LoadPage(&requests.FPDF_LoadPage{
						Document: doc,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_LoadPage).To(Not(BeNil()))
					Expect(FPDF_LoadPage.Page).To(Not(BeNil()))

					FPDF_ClosePage, err := pdfiumContainer.FPDF_ClosePage(&requests.FPDF_ClosePage{
						Page: FPDF_LoadPage.Page,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_ClosePage).To(Not(BeNil()))
				})

				It("returns the correct page width", func() {
					pageCount, err := pdfiumContainer.FPDF_GetPageWidth(&requests.FPDF_GetPageWidth{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})
					Expect(err).To(BeNil())
					Expect(pageCount).To(Equal(&responses.FPDF_GetPageWidth{
						Page:  0,
						Width: 595.2755737304688,
					}))
				})

				It("returns an error when calling FPDF_GetPageWidth without a page", func() {
					pageCount, err := pdfiumContainer.FPDF_GetPageWidth(&requests.FPDF_GetPageWidth{})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(pageCount).To(BeNil())
				})

				It("returns the correct page height", func() {
					pageCount, err := pdfiumContainer.FPDF_GetPageHeight(&requests.FPDF_GetPageHeight{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})
					Expect(err).To(BeNil())
					Expect(pageCount).To(Equal(&responses.FPDF_GetPageHeight{
						Page:   0,
						Height: 841.8897094726562,
					}))
				})

				It("returns an error when calling FPDF_GetPageHeight without a page", func() {
					pageCount, err := pdfiumContainer.FPDF_GetPageHeight(&requests.FPDF_GetPageHeight{})
					Expect(err).To(MatchError("either page reference or index should be given"))
					Expect(pageCount).To(BeNil())
				})

				It("returns the correct page size by index", func() {
					pageCount, err := pdfiumContainer.FPDF_GetPageSizeByIndex(&requests.FPDF_GetPageSizeByIndex{
						Document: doc,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(pageCount).To(Equal(&responses.FPDF_GetPageSizeByIndex{
						Page:   0,
						Width:  595.2755737304688,
						Height: 841.8897094726562,
					}))
				})

				It("returns an error when calling FPDF_GetPageSizeByIndex with an invalid page", func() {
					pageCount, err := pdfiumContainer.FPDF_GetPageSizeByIndex(&requests.FPDF_GetPageSizeByIndex{
						Document: doc,
						Index:    3,
					})
					Expect(err).To(MatchError("could not load page size by index"))
					Expect(pageCount).To(BeNil())
				})
			})
		})

		Context("a normal PDF file with multiple pages", func() {
			var doc references.FPDF_DOCUMENT
			var file *os.File

			BeforeEach(func() {
				pdfFile, err := os.Open(testsPath + "/testdata/test_multipage.pdf")
				Expect(err).To(BeNil())

				file = pdfFile
				fileStat, err := file.Stat()
				Expect(err).To(BeNil())

				newDoc, err := pdfiumContainer.FPDF_LoadCustomDocument(&requests.FPDF_LoadCustomDocument{
					Reader: file,
					Size:   fileStat.Size(),
				})
				Expect(err).To(BeNil())

				doc = newDoc.Document
			})

			AfterEach(func() {
				FPDF_CloseDocument, err := pdfiumContainer.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_CloseDocument).To(Not(BeNil()))
				file.Close()
			})

			When("is opened", func() {
				It("returns the correct file version", func() {
					pageCount, err := pdfiumContainer.FPDF_GetFileVersion(&requests.FPDF_GetFileVersion{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(pageCount).To(Equal(&responses.FPDF_GetFileVersion{
						FileVersion: 15,
					}))
				})

				It("returns the correct page count", func() {
					pageCount, err := pdfiumContainer.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(pageCount).To(Equal(&responses.FPDF_GetPageCount{
						PageCount: 2,
					}))
				})
			})
		})

		Context("a normal PDF file with alpha channel", func() {
			var doc references.FPDF_DOCUMENT
			var file *os.File

			BeforeEach(func() {
				pdfFile, err := os.Open(testsPath + "/testdata/alpha_channel.pdf")
				Expect(err).To(BeNil())

				file = pdfFile
				fileStat, err := file.Stat()
				Expect(err).To(BeNil())

				newDoc, err := pdfiumContainer.FPDF_LoadCustomDocument(&requests.FPDF_LoadCustomDocument{
					Reader: file,
					Size:   fileStat.Size(),
				})
				Expect(err).To(BeNil())

				doc = newDoc.Document
			})

			AfterEach(func() {
				FPDF_CloseDocument, err := pdfiumContainer.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_CloseDocument).To(Not(BeNil()))
				file.Close()
			})

			When("is opened", func() {
				It("returns the correct file version", func() {
					pageCount, err := pdfiumContainer.FPDF_GetFileVersion(&requests.FPDF_GetFileVersion{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(pageCount).To(Equal(&responses.FPDF_GetFileVersion{
						FileVersion: 17,
					}))
				})

				It("returns the correct page count", func() {
					pageCount, err := pdfiumContainer.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(pageCount).To(Equal(&responses.FPDF_GetPageCount{
						PageCount: 1,
					}))
				})
			})
		})

	})
}
