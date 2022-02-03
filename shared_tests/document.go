package shared_tests

import (
	"github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/references"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

func RunDocumentTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("NewDocumentFromBytes", func() {
		Context("a normal PDF file with 1 page", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				pdfData, err := ioutil.ReadFile(testsPath + "/testdata/test.pdf")
				Expect(err).To(BeNil())
				if err != nil {
					return
				}

				newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData)
				if err != nil {
					return
				}

				doc = *newDoc
			})

			AfterEach(func() {
				err := pdfiumContainer.FPDF_CloseDocument(doc)
				Expect(err).To(BeNil())
			})

			When("is opened", func() {
				It("returns the correct file version", func() {
					fileVersion, err := pdfiumContainer.FPDF_GetFileVersion(&requests.FPDF_GetFileVersion{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(fileVersion).To(Equal(&responses.FPDF_GetFileVersion{
						FileVersion: 15,
					}))
				})

				It("returns the correct references permissions", func() {
					docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					securityHandlerRevision, err := pdfiumContainer.FPDF_GetSecurityHandlerRevision(&requests.FPDF_GetSecurityHandlerRevision{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(securityHandlerRevision).To(Equal(&responses.FPDF_GetSecurityHandlerRevision{
						SecurityHandlerRevision: -1, // -1 = no security handler.
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

				It("returns the correct page mode", func() {
					pageMode, err := pdfiumContainer.FPDFDoc_GetPageMode(&requests.FPDFDoc_GetPageMode{
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
				pdfData, err := ioutil.ReadFile(testsPath + "/testdata/test_multipage.pdf")
				Expect(err).To(BeNil())
				if err != nil {
					return
				}

				newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData)
				if err != nil {
					return
				}

				doc = *newDoc
			})

			AfterEach(func() {
				err := pdfiumContainer.FPDF_CloseDocument(doc)
				Expect(err).To(BeNil())
			})

			When("is opened", func() {
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

		Context("a password protected PDF file", func() {
			pdfData, _ := ioutil.ReadFile(testsPath + "/testdata/password_test123.pdf")
			When("is opened with no password", func() {
				It("returns the password error", func() {
					doc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData)
					Expect(err).To(MatchError(errors.ErrPassword.Error()))
					Expect(doc).To(BeNil())
				})
			})
			When("is opened with the wrong password", func() {
				It("returns the password error", func() {
					wrongPassword := "test"
					doc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(wrongPassword))
					Expect(err).To(MatchError(errors.ErrPassword.Error()))
					Expect(doc).To(BeNil())
				})
			})
			When("is opened with the correct password", func() {
				It("does not return an error", func() {
					pdfPassword := "test123"
					doc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					Expect(err).To(BeNil())
					Expect(doc).To(Not(BeNil()))
					err = pdfiumContainer.FPDF_CloseDocument(*doc)
					Expect(err).To(BeNil())
				})
			})
		})
	})

	Describe("NewDocumentFromFilePath", func() {
		Context("a normal PDF file with 1 page", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				newDoc, err := pdfiumContainer.NewDocumentFromFilePath(testsPath + "/testdata/test.pdf")
				if err != nil {
					return
				}

				doc = *newDoc
			})

			AfterEach(func() {
				err := pdfiumContainer.FPDF_CloseDocument(doc)
				Expect(err).To(BeNil())
			})

			When("is opened", func() {
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

		Context("a normal PDF file with multiple pages", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				newDoc, err := pdfiumContainer.NewDocumentFromFilePath(testsPath + "/testdata/test_multipage.pdf")
				if err != nil {
					return
				}

				doc = *newDoc
			})

			AfterEach(func() {
				err := pdfiumContainer.FPDF_CloseDocument(doc)
				Expect(err).To(BeNil())
			})

			When("is opened", func() {
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

		Context("a password protected PDF file", func() {
			filePath := testsPath + "/testdata/password_test123.pdf"
			When("is opened with no password", func() {
				It("returns the password error", func() {
					doc, err := pdfiumContainer.NewDocumentFromFilePath(filePath)
					Expect(err).To(MatchError(errors.ErrPassword.Error()))
					Expect(doc).To(BeNil())
				})
			})
			When("is opened with the wrong password", func() {
				It("returns the password error", func() {
					wrongPassword := "test"
					doc, err := pdfiumContainer.NewDocumentFromFilePath(filePath, pdfium.OpenDocumentWithPasswordOption(wrongPassword))
					Expect(err).To(MatchError(errors.ErrPassword.Error()))
					Expect(doc).To(BeNil())
				})
			})
			When("is opened with the correct password", func() {
				It("does not return an error", func() {
					pdfPassword := "test123"
					doc, err := pdfiumContainer.NewDocumentFromFilePath(filePath, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					Expect(err).To(BeNil())
					Expect(doc).To(Not(BeNil()))
					err = pdfiumContainer.FPDF_CloseDocument(*doc)
					Expect(err).To(BeNil())
				})
			})
		})

		Context("a non-existent file", func() {
			filePath := testsPath + "/testdata/i_dont_exist.pdf"
			When("is opened", func() {
				It("returns the file error", func() {
					doc, err := pdfiumContainer.NewDocumentFromFilePath(filePath)
					Expect(err).To(MatchError(errors.ErrFile.Error()))
					Expect(doc).To(BeNil())
				})
			})
		})
	})

	Describe("NewDocumentFromReader", func() {
		Context("a password protected PDF file", func() {
			When("is opened with no password", func() {
				It("returns the password error", func() {
					file, err := os.Open(testsPath + "/testdata/password_test123.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}
					fileStat, err := file.Stat()
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					doc, err := pdfiumContainer.NewDocumentFromReader(file, int(fileStat.Size()))
					Expect(err).To(MatchError(errors.ErrPassword.Error()))
					Expect(doc).To(BeNil())
				})
			})
			When("is opened with the wrong password", func() {
				It("returns the password error", func() {
					file, err := os.Open(testsPath + "/testdata/password_test123.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}
					fileStat, err := file.Stat()
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					wrongPassword := "test"
					doc, err := pdfiumContainer.NewDocumentFromReader(file, int(fileStat.Size()), pdfium.OpenDocumentWithPasswordOption(wrongPassword))
					Expect(err).To(MatchError(errors.ErrPassword.Error()))
					Expect(doc).To(BeNil())
				})
			})
			When("is opened with the correct password", func() {
				It("does not return an error", func() {
					file, err := os.Open(testsPath + "/testdata/password_test123.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}
					fileStat, err := file.Stat()
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					pdfPassword := "test123"
					doc, err := pdfiumContainer.NewDocumentFromReader(file, int(fileStat.Size()), pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					Expect(err).To(BeNil())
					Expect(doc).To(Not(BeNil()))

					fileVersion, err := pdfiumContainer.FPDF_GetFileVersion(&requests.FPDF_GetFileVersion{
						Document: *doc,
					})
					Expect(err).To(BeNil())
					Expect(fileVersion).To(Equal(&responses.FPDF_GetFileVersion{
						FileVersion: 15,
					}))

					docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
						Document: *doc,
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

					securityHandlerRevision, err := pdfiumContainer.FPDF_GetSecurityHandlerRevision(&requests.FPDF_GetSecurityHandlerRevision{
						Document: *doc,
					})
					Expect(err).To(BeNil())
					Expect(securityHandlerRevision).To(Equal(&responses.FPDF_GetSecurityHandlerRevision{
						SecurityHandlerRevision: 3,
					}))

					err = pdfiumContainer.FPDF_CloseDocument(*doc)
					Expect(err).To(BeNil())
				})
			})
		})

		Context("a protected PDF file with no permissions", func() {
			Context("is opened with the user password", func() {
				var doc references.FPDF_DOCUMENT

				BeforeEach(func() {
					pdfPassword := "test123"
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_none.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_none.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_printing.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_printing.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_degraded_printing.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_degraded_printing.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_modify_contents.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_modify_contents.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_assembly.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_assembly.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_copy_contents.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_copy_contents.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_screen_readers.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_screen_readers.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_modify_annotations.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_modify_annotations.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_fill_in.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_fill_in.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_all_features.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
					pdfData, err := ioutil.ReadFile(testsPath + "/testdata/permissions_all_features.pdf")
					Expect(err).To(BeNil())
					if err != nil {
						return
					}

					newDoc, err := pdfiumContainer.NewDocumentFromBytes(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					if err != nil {
						return
					}

					doc = *newDoc
				})

				AfterEach(func() {
					err := pdfiumContainer.FPDF_CloseDocument(doc)
					Expect(err).To(BeNil())
				})

				When("is opened", func() {
					It("returns the correct permission", func() {
						docPermissions, err := pdfiumContainer.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{
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
}
