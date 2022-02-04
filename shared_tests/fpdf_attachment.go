package shared_tests

import (
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

func RunfpdfAttachmentTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("fpdf_attachment", func() {
		Context("no document is given", func() {
			It("returns an error when FPDFDoc_GetAttachmentCount is called", func() {
				FPDFDoc_GetAttachmentCount, err := pdfiumContainer.FPDFDoc_GetAttachmentCount(&requests.FPDFDoc_GetAttachmentCount{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFDoc_GetAttachmentCount).To(BeNil())
			})

			It("returns an error when FPDFDoc_AddAttachment is called", func() {
				FPDFDoc_AddAttachment, err := pdfiumContainer.FPDFDoc_AddAttachment(&requests.FPDFDoc_AddAttachment{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFDoc_AddAttachment).To(BeNil())
			})

			It("returns an error when FPDFDoc_GetAttachment is called", func() {
				FPDFDoc_GetAttachment, err := pdfiumContainer.FPDFDoc_GetAttachment(&requests.FPDFDoc_GetAttachment{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFDoc_GetAttachment).To(BeNil())
			})

			It("returns an error when FPDFDoc_DeleteAttachment is called", func() {
				FPDFDoc_DeleteAttachment, err := pdfiumContainer.FPDFDoc_DeleteAttachment(&requests.FPDFDoc_DeleteAttachment{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFDoc_DeleteAttachment).To(BeNil())
			})
		})

		Context("no attachment is given", func() {
			It("returns an error when FPDFAttachment_GetName is called", func() {
				FPDFAttachment_GetName, err := pdfiumContainer.FPDFAttachment_GetName(&requests.FPDFAttachment_GetName{})
				Expect(err).To(MatchError("attachment not given"))
				Expect(FPDFAttachment_GetName).To(BeNil())
			})

			It("returns an error when FPDFAttachment_HasKey is called", func() {
				FPDFAttachment_HasKey, err := pdfiumContainer.FPDFAttachment_HasKey(&requests.FPDFAttachment_HasKey{})
				Expect(err).To(MatchError("attachment not given"))
				Expect(FPDFAttachment_HasKey).To(BeNil())
			})

			It("returns an error when FPDFAttachment_SetStringValue is called", func() {
				FPDFAttachment_SetStringValue, err := pdfiumContainer.FPDFAttachment_SetStringValue(&requests.FPDFAttachment_SetStringValue{})
				Expect(err).To(MatchError("attachment not given"))
				Expect(FPDFAttachment_SetStringValue).To(BeNil())
			})

			It("returns an error when FPDFAttachment_GetStringValue is called", func() {
				FPDFAttachment_GetStringValue, err := pdfiumContainer.FPDFAttachment_GetStringValue(&requests.FPDFAttachment_GetStringValue{})
				Expect(err).To(MatchError("attachment not given"))
				Expect(FPDFAttachment_GetStringValue).To(BeNil())
			})

			It("returns an error when FPDFAttachment_SetFile is called", func() {
				FPDFAttachment_SetFile, err := pdfiumContainer.FPDFAttachment_SetFile(&requests.FPDFAttachment_SetFile{})
				Expect(err).To(MatchError("attachment not given"))
				Expect(FPDFAttachment_SetFile).To(BeNil())
			})

			It("returns an error when FPDFAttachment_GetFile is called", func() {
				FPDFAttachment_GetFile, err := pdfiumContainer.FPDFAttachment_GetFile(&requests.FPDFAttachment_GetFile{})
				Expect(err).To(MatchError("attachment not given"))
				Expect(FPDFAttachment_GetFile).To(BeNil())
			})
		})

		Context("a normal PDF file without attachments", func() {
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

			It("returns no an attachment count of 0", func() {
				FPDFDoc_GetAttachmentCount, err := pdfiumContainer.FPDFDoc_GetAttachmentCount(&requests.FPDFDoc_GetAttachmentCount{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDFDoc_GetAttachmentCount).To(Equal(&responses.FPDFDoc_GetAttachmentCount{
					AttachmentCount: 0,
				}))
			})

			It("returns an error when trying to get an attachment that isn't there", func() {
				FPDFDoc_GetAttachment, err := pdfiumContainer.FPDFDoc_GetAttachment(&requests.FPDFDoc_GetAttachment{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(MatchError("could not get attachment object"))
				Expect(FPDFDoc_GetAttachment).To(BeNil())
			})

			It("allows to add a new attachment", func() {
				FPDFDoc_AddAttachment, err := pdfiumContainer.FPDFDoc_AddAttachment(&requests.FPDFDoc_AddAttachment{
					Document: doc,
					Name:     "Attachment",
				})
				Expect(err).To(BeNil())
				Expect(FPDFDoc_AddAttachment).To(Not(BeNil()))
				Expect(FPDFDoc_AddAttachment.Attachment).To(Not(BeNil()))
			})

			When("an attachment has been added", func() {
				var attachment references.FPDF_ATTACHMENT
				BeforeEach(func() {
					FPDFDoc_AddAttachment, err := pdfiumContainer.FPDFDoc_AddAttachment(&requests.FPDFDoc_AddAttachment{
						Document: doc,
						Name:     "Attachment",
					})
					Expect(err).To(BeNil())
					Expect(FPDFDoc_AddAttachment).To(Not(BeNil()))
					Expect(FPDFDoc_AddAttachment.Attachment).To(Not(BeNil()))
					attachment = FPDFDoc_AddAttachment.Attachment
				})

				It("returns the attachment that has just been added", func() {
					FPDFDoc_GetAttachment, err := pdfiumContainer.FPDFDoc_GetAttachment(&requests.FPDFDoc_GetAttachment{
						Document: doc,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFDoc_GetAttachment).To(Not(BeNil()))
					Expect(FPDFDoc_GetAttachment.Attachment).To(Not(BeNil()))
				})

				It("returns the name of the attachment that has just been added", func() {
					FPDFAttachment_GetName, err := pdfiumContainer.FPDFAttachment_GetName(&requests.FPDFAttachment_GetName{
						Attachment: attachment,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_GetName).To(Equal(&responses.FPDFAttachment_GetName{
						Name: "Attachment",
					}))
				})

				It("returns that a non-existent key could not be found", func() {
					FPDFAttachment_HasKey, err := pdfiumContainer.FPDFAttachment_HasKey(&requests.FPDFAttachment_HasKey{
						Attachment: attachment,
						Key:        "CreationDate",
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_HasKey).To(Equal(&responses.FPDFAttachment_HasKey{
						Key:    "CreationDate",
						HasKey: false,
					}))
				})

				It("gives a nil content array when no file content has been set and FPDFAttachment_GetFile is executed", func() {
					FPDFAttachment_GetFile, err := pdfiumContainer.FPDFAttachment_GetFile(&requests.FPDFAttachment_GetFile{
						Attachment: attachment,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_GetFile).To(Equal(&responses.FPDFAttachment_GetFile{}))
				})

				It("allows for the file content to be set and match when we retrieve it", func() {
					FPDFAttachment_SetFile, err := pdfiumContainer.FPDFAttachment_SetFile(&requests.FPDFAttachment_SetFile{
						Attachment: attachment,
						Contents:   []byte{1, 2, 3},
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_SetFile).To(Equal(&responses.FPDFAttachment_SetFile{}))

					FPDFAttachment_GetFile, err := pdfiumContainer.FPDFAttachment_GetFile(&requests.FPDFAttachment_GetFile{
						Attachment: attachment,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_GetFile).To(Equal(&responses.FPDFAttachment_GetFile{
						Contents: []byte{1, 2, 3},
					}))
				})

				It("allows the attachment to be deleted", func() {
					FPDFDoc_DeleteAttachment, err := pdfiumContainer.FPDFDoc_DeleteAttachment(&requests.FPDFDoc_DeleteAttachment{
						Document: doc,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFDoc_DeleteAttachment).To(Equal(&responses.FPDFDoc_DeleteAttachment{}))
				})
			})
		})

		Context("a PDF file with attachments", func() {
			var doc references.FPDF_DOCUMENT

			BeforeEach(func() {
				pdfData, err := ioutil.ReadFile(testsPath + "/testdata/embedded_attachments.pdf")
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

			It("returns the correct attachments count", func() {
				FPDFDoc_GetAttachmentCount, err := pdfiumContainer.FPDFDoc_GetAttachmentCount(&requests.FPDFDoc_GetAttachmentCount{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDFDoc_GetAttachmentCount).To(Equal(&responses.FPDFDoc_GetAttachmentCount{
					AttachmentCount: 2,
				}))
			})

			It("returns the correct attachments", func() {
				GetAttachments, err := pdfiumContainer.GetAttachments(&requests.GetAttachments{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(GetAttachments).To(Not(BeNil()))
				Expect(GetAttachments.Attachments).To(Not(BeNil()))
				Expect(GetAttachments.Attachments).To(HaveLen(2))
				Expect(GetAttachments.Attachments[0].Name).To(Equal("1.txt"))
				Expect(GetAttachments.Attachments[0].Content).To(Equal([]byte("test")))
				Expect(GetAttachments.Attachments[0].Values).To(Equal([]responses.AttachmentValue{
					{Key: "Size", ValueType: 2, StringValue: ""},
					{
						Key:         "CreationDate",
						ValueType:   3,
						StringValue: "D:20170712214438-07'00'",
					},
					{
						Key:         "CheckSum",
						ValueType:   3,
						StringValue: "<098F6BCD4621D373CADE4E832627B4F6>",
					},
				}))

				Expect(GetAttachments.Attachments[1].Name).To(Equal("attached.pdf"))
				Expect(GetAttachments.Attachments[1].Content).To(HaveLen(5869))
				Expect(GetAttachments.Attachments[1].Values).To(Equal([]responses.AttachmentValue{
					{Key: "Size", ValueType: 2, StringValue: ""},
					{
						Key:         "CreationDate",
						ValueType:   3,
						StringValue: "D:20170712214443-07'00'",
					},
					{
						Key:         "CheckSum",
						ValueType:   3,
						StringValue: "<72AFCDDEDF554DDA63C0C88E06F1CE18>",
					},
				}))
			})

			When("the first attachment has been loaded", func() {
				var attachment references.FPDF_ATTACHMENT
				BeforeEach(func() {
					FPDFDoc_GetAttachment, err := pdfiumContainer.FPDFDoc_GetAttachment(&requests.FPDFDoc_GetAttachment{
						Document: doc,
						Index:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFDoc_GetAttachment).To(Not(BeNil()))
					Expect(FPDFDoc_GetAttachment.Attachment).To(Not(BeNil()))
					attachment = FPDFDoc_GetAttachment.Attachment
				})

				It("returns the correct name of the attachment", func() {
					FPDFAttachment_GetName, err := pdfiumContainer.FPDFAttachment_GetName(&requests.FPDFAttachment_GetName{
						Attachment: attachment,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_GetName).To(Equal(&responses.FPDFAttachment_GetName{
						Name: "1.txt",
					}))
				})

				It("returns the correct file content", func() {
					FPDFAttachment_GetFile, err := pdfiumContainer.FPDFAttachment_GetFile(&requests.FPDFAttachment_GetFile{
						Attachment: attachment,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_GetFile).To(Equal(&responses.FPDFAttachment_GetFile{
						Contents: []byte("test"),
					}))
				})

				It("returns that it has a key Size", func() {
					FPDFAttachment_HasKey, err := pdfiumContainer.FPDFAttachment_HasKey(&requests.FPDFAttachment_HasKey{
						Attachment: attachment,
						Key:        "Size",
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_HasKey).To(Equal(&responses.FPDFAttachment_HasKey{
						Key:    "Size",
						HasKey: true,
					}))
				})

				It("returns that it has a key CreationDate", func() {
					FPDFAttachment_HasKey, err := pdfiumContainer.FPDFAttachment_HasKey(&requests.FPDFAttachment_HasKey{
						Attachment: attachment,
						Key:        "CreationDate",
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_HasKey).To(Equal(&responses.FPDFAttachment_HasKey{
						Key:    "CreationDate",
						HasKey: true,
					}))
				})

				It("returns that it has a key CheckSum", func() {
					FPDFAttachment_HasKey, err := pdfiumContainer.FPDFAttachment_HasKey(&requests.FPDFAttachment_HasKey{
						Attachment: attachment,
						Key:        "CheckSum",
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_HasKey).To(Equal(&responses.FPDFAttachment_HasKey{
						Key:    "CheckSum",
						HasKey: true,
					}))
				})

				It("returns the right value type for key Size", func() {
					FPDFAttachment_GetValueType, err := pdfiumContainer.FPDFAttachment_GetValueType(&requests.FPDFAttachment_GetValueType{
						Attachment: attachment,
						Key:        "Size",
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_GetValueType).To(Equal(&responses.FPDFAttachment_GetValueType{
						Key:       "Size",
						ValueType: enums.FPDF_OBJECT_TYPE_NUMBER,
					}))
				})

				It("returns the right value type for key CreationDate", func() {
					FPDFAttachment_GetValueType, err := pdfiumContainer.FPDFAttachment_GetValueType(&requests.FPDFAttachment_GetValueType{
						Attachment: attachment,
						Key:        "CreationDate",
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_GetValueType).To(Equal(&responses.FPDFAttachment_GetValueType{
						Key:       "CreationDate",
						ValueType: enums.FPDF_OBJECT_TYPE_STRING,
					}))
				})

				It("returns the right value type for key CheckSum", func() {
					FPDFAttachment_GetValueType, err := pdfiumContainer.FPDFAttachment_GetValueType(&requests.FPDFAttachment_GetValueType{
						Attachment: attachment,
						Key:        "CheckSum",
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_GetValueType).To(Equal(&responses.FPDFAttachment_GetValueType{
						Key:       "CheckSum",
						ValueType: enums.FPDF_OBJECT_TYPE_STRING,
					}))
				})

				It("returns the right string value for key CreationDate", func() {
					FPDFAttachment_GetStringValue, err := pdfiumContainer.FPDFAttachment_GetStringValue(&requests.FPDFAttachment_GetStringValue{
						Attachment: attachment,
						Key:        "CreationDate",
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_GetStringValue).To(Equal(&responses.FPDFAttachment_GetStringValue{
						Key:   "CreationDate",
						Value: "D:20170712214438-07'00'",
					}))
				})

				It("returns the right string value for key CheckSum", func() {
					FPDFAttachment_GetStringValue, err := pdfiumContainer.FPDFAttachment_GetStringValue(&requests.FPDFAttachment_GetStringValue{
						Attachment: attachment,
						Key:        "CheckSum",
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_GetStringValue).To(Equal(&responses.FPDFAttachment_GetStringValue{
						Key:   "CheckSum",
						Value: "<098F6BCD4621D373CADE4E832627B4F6>",
					}))
				})

				It("allows for a string value to be set, returns the correct type and value", func() {
					FPDFAttachment_SetStringValue, err := pdfiumContainer.FPDFAttachment_SetStringValue(&requests.FPDFAttachment_SetStringValue{
						Attachment: attachment,
						Key:        "RandomValue",
						Value:      "Test123",
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_SetStringValue).To(Equal(&responses.FPDFAttachment_SetStringValue{
						Key:   "RandomValue",
						Value: "Test123",
					}))

					FPDFAttachment_GetValueType, err := pdfiumContainer.FPDFAttachment_GetValueType(&requests.FPDFAttachment_GetValueType{
						Attachment: attachment,
						Key:        "RandomValue",
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_GetValueType).To(Equal(&responses.FPDFAttachment_GetValueType{
						Key:       "RandomValue",
						ValueType: enums.FPDF_OBJECT_TYPE_STRING,
					}))

					FPDFAttachment_GetStringValue, err := pdfiumContainer.FPDFAttachment_GetStringValue(&requests.FPDFAttachment_GetStringValue{
						Attachment: attachment,
						Key:        "RandomValue",
					})
					Expect(err).To(BeNil())
					Expect(FPDFAttachment_GetStringValue).To(Equal(&responses.FPDFAttachment_GetStringValue{
						Key:   "RandomValue",
						Value: "Test123",
					}))
				})
			})
		})
	})
}
