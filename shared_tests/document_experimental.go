//go:build pdfium_experimental
// +build pdfium_experimental

package shared_tests

import (
	"io/ioutil"
	"os"

	"github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("document", func() {
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

	Describe("FPDF_LoadMemDocument", func() {
		Context("a normal PDF file with 1 page", func() {
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
				It("returns the correct file version", func() {
					fileVersion, err := PdfiumInstance.FPDF_GetFileVersion(&requests.FPDF_GetFileVersion{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(fileVersion).To(Equal(&responses.FPDF_GetFileVersion{
						FileVersion: 15,
					}))
				})

				It("returns the correct references permissions", func() {
					docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
						DocPermissions:                      0xffffffff, // 0xffffffff (4294967295) = not protected
						PrintDocument:                       true,
						ModifyContents:                      true,
						CopyOrExtractText:                   true,
						AddOrModifyTextAnnotations:          true,
						FillInInteractiveFormFields:         true,
						CreateOrModifyInteractiveFormFields: true,
						FillInExistingInteractiveFormFields: true,
						ExtractTextAndGraphics:              true,
						AssembleDocument:                    true,
						PrintDocumentAsFaithfulDigitalCopy:  true,
					}))
				})

				It("returns the correct security handler revision", func() {
					securityHandlerRevision, err := PdfiumInstance.FPDF_GetSecurityHandlerRevision(&requests.FPDF_GetSecurityHandlerRevision{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(securityHandlerRevision).To(Equal(&responses.FPDF_GetSecurityHandlerRevision{
						SecurityHandlerRevision: -1, // -1 = no security handler.
					}))
				})

				It("returns the correct page count", func() {
					pageCount, err := PdfiumInstance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(pageCount).To(Equal(&responses.FPDF_GetPageCount{
						PageCount: 1,
					}))
				})

				It("returns the correct page mode", func() {
					pageMode, err := PdfiumInstance.FPDFDoc_GetPageMode(&requests.FPDFDoc_GetPageMode{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(pageMode).To(Equal(&responses.FPDFDoc_GetPageMode{
						PageMode: responses.FPDFDoc_GetPageModeModeUseNone,
					}))
				})
			})
		})

		Context("a normal PDF file with multiple pages", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/test_multipage.pdf")
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
				It("returns the correct page count", func() {
					pageCount, err := PdfiumInstance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(pageCount).To(Equal(&responses.FPDF_GetPageCount{
						PageCount: 2,
					}))
				})
			})
		})

		Context("a password protected PDF file", func() {
			var pdfData []byte
			BeforeEach(func() {
				var err error
				pdfData, err = ioutil.ReadFile(TestDataPath + "/testdata/password_test123.pdf")
				Expect(err).To(BeNil())
			})

			When("is opened with no password", func() {
				It("returns the password error", func() {
					doc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data: &pdfData,
					})
					Expect(err).To(MatchError(errors.ErrPassword.Error()))
					Expect(doc).To(BeNil())
				})

				It("returns the password error", func() {
					doc, err := PdfiumInstance.FPDF_LoadMemDocument64(&requests.FPDF_LoadMemDocument64{
						Data: &pdfData,
					})
					Expect(err).To(MatchError(errors.ErrPassword.Error()))
					Expect(doc).To(BeNil())
				})
			})

			When("is opened with the wrong password", func() {
				It("returns the password error", func() {
					wrongPassword := "test"
					doc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &wrongPassword,
					})
					Expect(err).To(MatchError(errors.ErrPassword.Error()))
					Expect(doc).To(BeNil())
				})
			})
			When("is opened with the correct password", func() {
				It("does not return an error", func() {
					pdfPassword := "test123"
					doc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
					})
					Expect(err).To(BeNil())
					Expect(doc).To(Not(BeNil()))
					FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
						Document: doc.Document,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_CloseDocument).To(Not(BeNil()))
				})
			})
		})
	})

	Describe("FPDF_LoadDocument", func() {
		Context("a normal PDF file with 1 page", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				filePath := TestDataPath + "/testdata/test.pdf"
				newDoc, err := PdfiumInstance.FPDF_LoadDocument(&requests.FPDF_LoadDocument{
					Path: &filePath,
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
				It("returns the correct page count", func() {
					pageCount, err := PdfiumInstance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(pageCount).To(Equal(&responses.FPDF_GetPageCount{
						PageCount: 1,
					}))
				})
			})
		})

		Context("a normal PDF file with multiple pages", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				filePath := TestDataPath + "/testdata/test_multipage.pdf"
				newDoc, err := PdfiumInstance.FPDF_LoadDocument(&requests.FPDF_LoadDocument{
					Path: &filePath,
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
				It("returns the correct page count", func() {
					pageCount, err := PdfiumInstance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(pageCount).To(Equal(&responses.FPDF_GetPageCount{
						PageCount: 2,
					}))
				})
			})
		})

		Context("a password protected PDF file", func() {
			var filePath string
			BeforeEach(func() {
				filePath = TestDataPath + "/testdata/password_test123.pdf"
			})
			When("is opened with no password", func() {
				It("returns the password error", func() {
					doc, err := PdfiumInstance.FPDF_LoadDocument(&requests.FPDF_LoadDocument{
						Path: &filePath,
					})
					Expect(err).To(MatchError(errors.ErrPassword.Error()))
					Expect(doc).To(BeNil())
				})
			})
			When("is opened with the wrong password", func() {
				It("returns the password error", func() {
					wrongPassword := "test"
					doc, err := PdfiumInstance.FPDF_LoadDocument(&requests.FPDF_LoadDocument{
						Path:     &filePath,
						Password: &wrongPassword,
					})
					Expect(err).To(MatchError(errors.ErrPassword.Error()))
					Expect(doc).To(BeNil())
				})
			})
			When("is opened with the correct password", func() {
				It("does not return an error", func() {
					pdfPassword := "test123"
					doc, err := PdfiumInstance.FPDF_LoadDocument(&requests.FPDF_LoadDocument{
						Path:     &filePath,
						Password: &pdfPassword,
					})
					Expect(err).To(BeNil())
					Expect(doc).To(Not(BeNil()))
					FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
						Document: doc.Document,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_CloseDocument).To(Not(BeNil()))
				})
			})
		})

		Context("a non-existent file", func() {
			filePath := TestDataPath + "/testdata/i_dont_exist.pdf"
			When("is opened", func() {
				It("returns the file error", func() {
					doc, err := PdfiumInstance.FPDF_LoadDocument(&requests.FPDF_LoadDocument{
						Path: &filePath,
					})
					Expect(err).To(MatchError(errors.ErrFile.Error()))
					Expect(doc).To(BeNil())
				})
			})
		})
	})

	Describe("FPDF_LoadCustomDocument", func() {
		Context("a password protected PDF file", func() {
			When("is opened with no password", func() {
				It("returns the password error", func() {
					file, err := os.Open(TestDataPath + "/testdata/password_test123.pdf")
					Expect(err).To(BeNil())

					fileStat, err := file.Stat()
					Expect(err).To(BeNil())

					doc, err := PdfiumInstance.FPDF_LoadCustomDocument(&requests.FPDF_LoadCustomDocument{
						Reader: file,
						Size:   fileStat.Size(),
					})

					Expect(err).To(MatchError(errors.ErrPassword.Error()))
					Expect(doc).To(BeNil())
				})
			})
			When("is opened with the wrong password", func() {
				It("returns the password error", func() {
					file, err := os.Open(TestDataPath + "/testdata/password_test123.pdf")
					Expect(err).To(BeNil())

					fileStat, err := file.Stat()
					Expect(err).To(BeNil())

					wrongPassword := "test"

					doc, err := PdfiumInstance.FPDF_LoadCustomDocument(&requests.FPDF_LoadCustomDocument{
						Reader:   file,
						Size:     fileStat.Size(),
						Password: &wrongPassword,
					})

					Expect(err).To(MatchError(errors.ErrPassword.Error()))
					Expect(doc).To(BeNil())
				})
			})
			When("is opened with the correct password", func() {
				It("does not return an error", func() {
					file, err := os.Open(TestDataPath + "/testdata/password_test123.pdf")
					Expect(err).To(BeNil())

					fileStat, err := file.Stat()
					Expect(err).To(BeNil())

					pdfPassword := "test123"
					doc, err := PdfiumInstance.FPDF_LoadCustomDocument(&requests.FPDF_LoadCustomDocument{
						Reader:   file,
						Size:     fileStat.Size(),
						Password: &pdfPassword,
					})
					Expect(err).To(BeNil())
					Expect(doc).To(Not(BeNil()))

					fileVersion, err := PdfiumInstance.FPDF_GetFileVersion(&requests.FPDF_GetFileVersion{
						Document: doc.Document,
					})
					Expect(err).To(BeNil())
					Expect(fileVersion).To(Equal(&responses.FPDF_GetFileVersion{
						FileVersion: 15,
					}))

					docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
						Document: doc.Document,
					})
					Expect(err).To(BeNil())
					Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
						DocPermissions:                      0xFFFFFFFC, // 0xFFFFFFFC (4294967292) = owner password
						PrintDocument:                       true,
						ModifyContents:                      true,
						CopyOrExtractText:                   true,
						AddOrModifyTextAnnotations:          true,
						FillInInteractiveFormFields:         true,
						CreateOrModifyInteractiveFormFields: true,
						FillInExistingInteractiveFormFields: true,
						ExtractTextAndGraphics:              true,
						AssembleDocument:                    true,
						PrintDocumentAsFaithfulDigitalCopy:  true,
					}))

					securityHandlerRevision, err := PdfiumInstance.FPDF_GetSecurityHandlerRevision(&requests.FPDF_GetSecurityHandlerRevision{
						Document: doc.Document,
					})
					Expect(err).To(BeNil())
					Expect(securityHandlerRevision).To(Equal(&responses.FPDF_GetSecurityHandlerRevision{
						SecurityHandlerRevision: 3,
					}))

					FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
						Document: doc.Document,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_CloseDocument).To(Not(BeNil()))
				})
			})
		})

		Context("a protected PDF file with no permissions", func() {
			Context("is opened with the user password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "test123"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_none.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294963392,
							PrintDocument:                       false,
							ModifyContents:                      false,
							CopyOrExtractText:                   false,
							AddOrModifyTextAnnotations:          false,
							FillInInteractiveFormFields:         false,
							CreateOrModifyInteractiveFormFields: false,
							FillInExistingInteractiveFormFields: false,
							ExtractTextAndGraphics:              false,
							AssembleDocument:                    false,
							PrintDocumentAsFaithfulDigitalCopy:  false,
						}))
					})
				})
			})

			Context("is opened with the owner password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "123test"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_none.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294967292,
							PrintDocument:                       true,
							ModifyContents:                      true,
							CopyOrExtractText:                   true,
							AddOrModifyTextAnnotations:          true,
							FillInInteractiveFormFields:         true,
							CreateOrModifyInteractiveFormFields: true,
							FillInExistingInteractiveFormFields: true,
							ExtractTextAndGraphics:              true,
							AssembleDocument:                    true,
							PrintDocumentAsFaithfulDigitalCopy:  true,
						}))
					})
				})
			})
		})

		Context("a protected PDF file with printing permissions", func() {
			Context("is opened with the user password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "test123"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_printing.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294965444,
							PrintDocument:                       true,
							ModifyContents:                      false,
							CopyOrExtractText:                   false,
							AddOrModifyTextAnnotations:          false,
							FillInInteractiveFormFields:         false,
							CreateOrModifyInteractiveFormFields: false,
							FillInExistingInteractiveFormFields: false,
							ExtractTextAndGraphics:              false,
							AssembleDocument:                    false,
							PrintDocumentAsFaithfulDigitalCopy:  true,
						}))
					})
				})
			})

			Context("is opened with the owner password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "123test"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_printing.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294967292,
							PrintDocument:                       true,
							ModifyContents:                      true,
							CopyOrExtractText:                   true,
							AddOrModifyTextAnnotations:          true,
							FillInInteractiveFormFields:         true,
							CreateOrModifyInteractiveFormFields: true,
							FillInExistingInteractiveFormFields: true,
							ExtractTextAndGraphics:              true,
							AssembleDocument:                    true,
							PrintDocumentAsFaithfulDigitalCopy:  true,
						}))
					})
				})
			})
		})

		Context("a protected PDF file with degraded printing permissions", func() {
			Context("is opened with the user password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "test123"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_degraded_printing.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294963396,
							PrintDocument:                       true,
							ModifyContents:                      false,
							CopyOrExtractText:                   false,
							AddOrModifyTextAnnotations:          false,
							FillInInteractiveFormFields:         false,
							CreateOrModifyInteractiveFormFields: false,
							FillInExistingInteractiveFormFields: false,
							ExtractTextAndGraphics:              false,
							AssembleDocument:                    false,
							PrintDocumentAsFaithfulDigitalCopy:  false,
						}))
					})
				})
			})

			Context("is opened with the owner password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "123test"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_degraded_printing.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294967292,
							PrintDocument:                       true,
							ModifyContents:                      true,
							CopyOrExtractText:                   true,
							AddOrModifyTextAnnotations:          true,
							FillInInteractiveFormFields:         true,
							CreateOrModifyInteractiveFormFields: true,
							FillInExistingInteractiveFormFields: true,
							ExtractTextAndGraphics:              true,
							AssembleDocument:                    true,
							PrintDocumentAsFaithfulDigitalCopy:  true,
						}))
					})
				})
			})
		})

		Context("a protected PDF file with modify content permissions", func() {
			Context("is opened with the user password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "test123"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_modify_contents.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294964424,
							PrintDocument:                       false,
							ModifyContents:                      true,
							CopyOrExtractText:                   false,
							AddOrModifyTextAnnotations:          false,
							FillInInteractiveFormFields:         false,
							CreateOrModifyInteractiveFormFields: false,
							FillInExistingInteractiveFormFields: false,
							ExtractTextAndGraphics:              false,
							AssembleDocument:                    true,
							PrintDocumentAsFaithfulDigitalCopy:  false,
						}))
					})
				})
			})

			Context("is opened with the owner password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "123test"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_modify_contents.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294967292,
							PrintDocument:                       true,
							ModifyContents:                      true,
							CopyOrExtractText:                   true,
							AddOrModifyTextAnnotations:          true,
							FillInInteractiveFormFields:         true,
							CreateOrModifyInteractiveFormFields: true,
							FillInExistingInteractiveFormFields: true,
							ExtractTextAndGraphics:              true,
							AssembleDocument:                    true,
							PrintDocumentAsFaithfulDigitalCopy:  true,
						}))
					})
				})
			})
		})

		Context("a protected PDF file with assembly permissions", func() {
			Context("is opened with the user password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "test123"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_assembly.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294964416,
							PrintDocument:                       false,
							ModifyContents:                      false,
							CopyOrExtractText:                   false,
							AddOrModifyTextAnnotations:          false,
							FillInInteractiveFormFields:         false,
							CreateOrModifyInteractiveFormFields: false,
							FillInExistingInteractiveFormFields: false,
							ExtractTextAndGraphics:              false,
							AssembleDocument:                    true,
							PrintDocumentAsFaithfulDigitalCopy:  false,
						}))
					})
				})
			})

			Context("is opened with the owner password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "123test"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_assembly.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294967292,
							PrintDocument:                       true,
							ModifyContents:                      true,
							CopyOrExtractText:                   true,
							AddOrModifyTextAnnotations:          true,
							FillInInteractiveFormFields:         true,
							CreateOrModifyInteractiveFormFields: true,
							FillInExistingInteractiveFormFields: true,
							ExtractTextAndGraphics:              true,
							AssembleDocument:                    true,
							PrintDocumentAsFaithfulDigitalCopy:  true,
						}))
					})
				})
			})
		})

		Context("a protected PDF file with copy contents permissions", func() {
			Context("is opened with the user password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "test123"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_copy_contents.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294963920,
							PrintDocument:                       false,
							ModifyContents:                      false,
							CopyOrExtractText:                   true,
							AddOrModifyTextAnnotations:          false,
							FillInInteractiveFormFields:         false,
							CreateOrModifyInteractiveFormFields: false,
							FillInExistingInteractiveFormFields: false,
							ExtractTextAndGraphics:              true,
							AssembleDocument:                    false,
							PrintDocumentAsFaithfulDigitalCopy:  false,
						}))
					})
				})
			})

			Context("is opened with the owner password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "123test"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_copy_contents.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294967292,
							PrintDocument:                       true,
							ModifyContents:                      true,
							CopyOrExtractText:                   true,
							AddOrModifyTextAnnotations:          true,
							FillInInteractiveFormFields:         true,
							CreateOrModifyInteractiveFormFields: true,
							FillInExistingInteractiveFormFields: true,
							ExtractTextAndGraphics:              true,
							AssembleDocument:                    true,
							PrintDocumentAsFaithfulDigitalCopy:  true,
						}))
					})
				})
			})
		})

		Context("a protected PDF file with screen readers permissions", func() {
			Context("is opened with the user password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "test123"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_screen_readers.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294963904,
							PrintDocument:                       false,
							ModifyContents:                      false,
							CopyOrExtractText:                   false,
							AddOrModifyTextAnnotations:          false,
							FillInInteractiveFormFields:         false,
							CreateOrModifyInteractiveFormFields: false,
							FillInExistingInteractiveFormFields: false,
							ExtractTextAndGraphics:              true,
							AssembleDocument:                    false,
							PrintDocumentAsFaithfulDigitalCopy:  false,
						}))
					})
				})
			})

			Context("is opened with the owner password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "123test"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_screen_readers.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294967292,
							PrintDocument:                       true,
							ModifyContents:                      true,
							CopyOrExtractText:                   true,
							AddOrModifyTextAnnotations:          true,
							FillInInteractiveFormFields:         true,
							CreateOrModifyInteractiveFormFields: true,
							FillInExistingInteractiveFormFields: true,
							ExtractTextAndGraphics:              true,
							AssembleDocument:                    true,
							PrintDocumentAsFaithfulDigitalCopy:  true,
						}))
					})
				})
			})
		})

		Context("a protected PDF file with modify annotations permissions", func() {
			Context("is opened with the user password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "test123"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_modify_annotations.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294963680,
							PrintDocument:                       false,
							ModifyContents:                      false,
							CopyOrExtractText:                   false,
							AddOrModifyTextAnnotations:          true,
							FillInInteractiveFormFields:         true,
							CreateOrModifyInteractiveFormFields: false,
							FillInExistingInteractiveFormFields: true,
							ExtractTextAndGraphics:              false,
							AssembleDocument:                    false,
							PrintDocumentAsFaithfulDigitalCopy:  false,
						}))
					})
				})
			})

			Context("is opened with the owner password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "123test"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_modify_annotations.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294967292,
							PrintDocument:                       true,
							ModifyContents:                      true,
							CopyOrExtractText:                   true,
							AddOrModifyTextAnnotations:          true,
							FillInInteractiveFormFields:         true,
							CreateOrModifyInteractiveFormFields: true,
							FillInExistingInteractiveFormFields: true,
							ExtractTextAndGraphics:              true,
							AssembleDocument:                    true,
							PrintDocumentAsFaithfulDigitalCopy:  true,
						}))
					})
				})
			})
		})

		Context("a protected PDF file with fill in permissions", func() {
			Context("is opened with the user password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "test123"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_fill_in.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294963648,
							PrintDocument:                       false,
							ModifyContents:                      false,
							CopyOrExtractText:                   false,
							AddOrModifyTextAnnotations:          false,
							FillInInteractiveFormFields:         false,
							CreateOrModifyInteractiveFormFields: false,
							FillInExistingInteractiveFormFields: true,
							ExtractTextAndGraphics:              false,
							AssembleDocument:                    false,
							PrintDocumentAsFaithfulDigitalCopy:  false,
						}))
					})
				})
			})

			Context("is opened with the owner password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "123test"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_fill_in.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294967292,
							PrintDocument:                       true,
							ModifyContents:                      true,
							CopyOrExtractText:                   true,
							AddOrModifyTextAnnotations:          true,
							FillInInteractiveFormFields:         true,
							CreateOrModifyInteractiveFormFields: true,
							FillInExistingInteractiveFormFields: true,
							ExtractTextAndGraphics:              true,
							AssembleDocument:                    true,
							PrintDocumentAsFaithfulDigitalCopy:  true,
						}))
					})
				})
			})
		})

		Context("a protected PDF file with all feature permissions", func() {
			Context("is opened with the user password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "test123"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_all_features.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294967292,
							PrintDocument:                       true,
							ModifyContents:                      true,
							CopyOrExtractText:                   true,
							AddOrModifyTextAnnotations:          true,
							FillInInteractiveFormFields:         true,
							CreateOrModifyInteractiveFormFields: true,
							FillInExistingInteractiveFormFields: true,
							ExtractTextAndGraphics:              true,
							AssembleDocument:                    true,
							PrintDocumentAsFaithfulDigitalCopy:  true,
						}))
					})
				})
			})

			Context("is opened with the owner password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "123test"
					pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/permissions_all_features.pdf")
					Expect(err).To(BeNil())

					newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data:     &pdfData,
						Password: &pdfPassword,
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
					It("returns the correct permission", func() {
						docPermissions, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(docPermissions).To(Equal(&responses.FPDF_GetDocPermissions{
							DocPermissions:                      4294967292,
							PrintDocument:                       true,
							ModifyContents:                      true,
							CopyOrExtractText:                   true,
							AddOrModifyTextAnnotations:          true,
							FillInInteractiveFormFields:         true,
							CreateOrModifyInteractiveFormFields: true,
							FillInExistingInteractiveFormFields: true,
							ExtractTextAndGraphics:              true,
							AssembleDocument:                    true,
							PrintDocumentAsFaithfulDigitalCopy:  true,
						}))
					})
				})
			})
		})
	})
})
