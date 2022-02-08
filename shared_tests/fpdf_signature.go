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

var _ = Describe("fpdf_signature", func() {
	Context("no document is given", func() {
		It("returns an error when getting the document signature count", func() {
			FPDF_GetSignatureCount, err := PdfiumInstance.FPDF_GetSignatureCount(&requests.FPDF_GetSignatureCount{})
			Expect(err).To(MatchError("document not given"))
			Expect(FPDF_GetSignatureCount).To(BeNil())
		})

		It("returns an error when getting the document signature object", func() {
			FPDF_GetSignatureObject, err := PdfiumInstance.FPDF_GetSignatureObject(&requests.FPDF_GetSignatureObject{})
			Expect(err).To(MatchError("document not given"))
			Expect(FPDF_GetSignatureObject).To(BeNil())
		})
	})

	Context("no signature is given", func() {
		It("returns an error when getting the signature content", func() {
			FPDFSignatureObj_GetContents, err := PdfiumInstance.FPDFSignatureObj_GetContents(&requests.FPDFSignatureObj_GetContents{})
			Expect(err).To(MatchError("signature not given"))
			Expect(FPDFSignatureObj_GetContents).To(BeNil())
		})

		It("returns an error when getting the signature byte range", func() {
			FPDFSignatureObj_GetByteRange, err := PdfiumInstance.FPDFSignatureObj_GetByteRange(&requests.FPDFSignatureObj_GetByteRange{})
			Expect(err).To(MatchError("signature not given"))
			Expect(FPDFSignatureObj_GetByteRange).To(BeNil())
		})

		It("returns an error when getting the signature sub filter", func() {
			FPDFSignatureObj_GetSubFilter, err := PdfiumInstance.FPDFSignatureObj_GetSubFilter(&requests.FPDFSignatureObj_GetSubFilter{})
			Expect(err).To(MatchError("signature not given"))
			Expect(FPDFSignatureObj_GetSubFilter).To(BeNil())
		})

		It("returns an error when getting the signature reason", func() {
			FPDFSignatureObj_GetReason, err := PdfiumInstance.FPDFSignatureObj_GetReason(&requests.FPDFSignatureObj_GetReason{})
			Expect(err).To(MatchError("signature not given"))
			Expect(FPDFSignatureObj_GetReason).To(BeNil())
		})

		It("returns an error when getting the signature time", func() {
			FPDFSignatureObj_GetTime, err := PdfiumInstance.FPDFSignatureObj_GetTime(&requests.FPDFSignatureObj_GetTime{})
			Expect(err).To(MatchError("signature not given"))
			Expect(FPDFSignatureObj_GetTime).To(BeNil())
		})

		It("returns an error when getting the signature DocMDPPermission", func() {
			FPDFSignatureObj_GetDocMDPPermission, err := PdfiumInstance.FPDFSignatureObj_GetDocMDPPermission(&requests.FPDFSignatureObj_GetDocMDPPermission{})
			Expect(err).To(MatchError("signature not given"))
			Expect(FPDFSignatureObj_GetDocMDPPermission).To(BeNil())
		})
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
			It("returns that it has no signatures", func() {
				signatureCount, err := PdfiumInstance.FPDF_GetSignatureCount(&requests.FPDF_GetSignatureCount{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(signatureCount).To(Equal(&responses.FPDF_GetSignatureCount{
					Count: 0,
				}))
			})
		})
	})

	Context("a PDF file with two signatures", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/two_signatures.pdf")
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

		When("the document count is requested", func() {
			It("returns that it has two signatures", func() {
				signatureCount, err := PdfiumInstance.FPDF_GetSignatureCount(&requests.FPDF_GetSignatureCount{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(signatureCount).To(Equal(&responses.FPDF_GetSignatureCount{
					Count: 2,
				}))
			})
		})

		When("the first signature object is opened", func() {
			var signature references.FPDF_SIGNATURE

			BeforeEach(func() {
				signatureResp, err := PdfiumInstance.FPDF_GetSignatureObject(&requests.FPDF_GetSignatureObject{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(signatureResp).To(Not(BeNil()))
				signature = signatureResp.Signature
			})

			It("returns the correct signature content", func() {
				signatureContent, err := PdfiumInstance.FPDFSignatureObj_GetContents(&requests.FPDFSignatureObj_GetContents{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				Expect(signatureContent).To(Equal(&responses.FPDFSignatureObj_GetContents{
					Contents: []byte{0x30, 0x80, 0x06, 0x09, 0x2A, 0x86, 0x48,
						0x86, 0xF7, 0x0D, 0x01, 0x07, 0x02, 0xA0,
						0x80, 0x30, 0x80, 0x02, 0x01, 0x01},
				}))
			})

			It("returns the correct signature byte range", func() {
				signatureByteRange, err := PdfiumInstance.FPDFSignatureObj_GetByteRange(&requests.FPDFSignatureObj_GetByteRange{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				Expect(signatureByteRange).To(Equal(&responses.FPDFSignatureObj_GetByteRange{
					ByteRange: []int{0, 10, 30, 10},
				}))
			})

			It("returns the correct signature sub filter", func() {
				subFilter, err := PdfiumInstance.FPDFSignatureObj_GetSubFilter(&requests.FPDFSignatureObj_GetSubFilter{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				expectedSubfilter := "ETSI.CAdES.detached"
				Expect(subFilter).To(Equal(&responses.FPDFSignatureObj_GetSubFilter{
					SubFilter: &expectedSubfilter,
				}))
			})

			It("returns no signature reason", func() {
				reason, err := PdfiumInstance.FPDFSignatureObj_GetReason(&requests.FPDFSignatureObj_GetReason{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				Expect(reason).To(Equal(&responses.FPDFSignatureObj_GetReason{}))
			})

			It("returns the correct signature time", func() {
				signatureTime, err := PdfiumInstance.FPDFSignatureObj_GetTime(&requests.FPDFSignatureObj_GetTime{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				expectedTime := "D:20200624093114+02'00'"
				Expect(signatureTime).To(Equal(&responses.FPDFSignatureObj_GetTime{
					Time: &expectedTime,
				}))
			})

			It("returns no DocMDPPermission", func() {
				docMDPPermission, err := PdfiumInstance.FPDFSignatureObj_GetDocMDPPermission(&requests.FPDFSignatureObj_GetDocMDPPermission{
					Signature: signature,
				})
				Expect(err).To(MatchError("could not get DocMDPPermission"))
				Expect(docMDPPermission).To(BeNil())
			})
		})
	})

	Context("a PDF file with no signature sub filter", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/signature_no_sub_filter.pdf")
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

		When("the document count is requested", func() {
			It("returns that it has one signatures", func() {
				signatureCount, err := PdfiumInstance.FPDF_GetSignatureCount(&requests.FPDF_GetSignatureCount{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(signatureCount).To(Equal(&responses.FPDF_GetSignatureCount{
					Count: 1,
				}))
			})
		})

		When("the first signature object is opened", func() {
			var signature references.FPDF_SIGNATURE

			BeforeEach(func() {
				signatureResp, err := PdfiumInstance.FPDF_GetSignatureObject(&requests.FPDF_GetSignatureObject{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(signatureResp).To(Not(BeNil()))
				signature = signatureResp.Signature
			})

			It("returns the correct signature content", func() {
				signatureContent, err := PdfiumInstance.FPDFSignatureObj_GetContents(&requests.FPDFSignatureObj_GetContents{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				Expect(signatureContent).To(Equal(&responses.FPDFSignatureObj_GetContents{
					Contents: []byte{0x30, 0x80, 0x06, 0x09, 0x2A, 0x86, 0x48,
						0x86, 0xF7, 0x0D, 0x01, 0x07, 0x02, 0xA0,
						0x80, 0x30, 0x80, 0x02, 0x01, 0x01},
				}))
			})

			It("returns the correct signature byte range", func() {
				signatureByteRange, err := PdfiumInstance.FPDFSignatureObj_GetByteRange(&requests.FPDFSignatureObj_GetByteRange{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				Expect(signatureByteRange).To(Equal(&responses.FPDFSignatureObj_GetByteRange{
					ByteRange: []int{0, 10, 30, 10},
				}))
			})

			It("returns the correct signature sub filter", func() {
				subFilter, err := PdfiumInstance.FPDFSignatureObj_GetSubFilter(&requests.FPDFSignatureObj_GetSubFilter{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				Expect(subFilter).To(Equal(&responses.FPDFSignatureObj_GetSubFilter{}))
			})

			It("returns no signature reason", func() {
				reason, err := PdfiumInstance.FPDFSignatureObj_GetReason(&requests.FPDFSignatureObj_GetReason{
					Signature: signature,
				})

				Expect(err).To(BeNil())
				Expect(reason).To(Equal(&responses.FPDFSignatureObj_GetReason{}))
			})

			It("returns the correct signature time", func() {
				signatureTime, err := PdfiumInstance.FPDFSignatureObj_GetTime(&requests.FPDFSignatureObj_GetTime{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				expectedTime := "D:20200624093114+02'00'"
				Expect(signatureTime).To(Equal(&responses.FPDFSignatureObj_GetTime{
					Time: &expectedTime,
				}))
			})

			It("returns no DocMDPPermission", func() {
				docMDPPermission, err := PdfiumInstance.FPDFSignatureObj_GetDocMDPPermission(&requests.FPDFSignatureObj_GetDocMDPPermission{
					Signature: signature,
				})
				Expect(err).To(MatchError("could not get DocMDPPermission"))
				Expect(docMDPPermission).To(BeNil())
			})
		})
	})

	Context("a PDF file with signature reason", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/signature_reason.pdf")
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

		When("the document count is requested", func() {
			It("returns that it has one signatures", func() {
				signatureCount, err := PdfiumInstance.FPDF_GetSignatureCount(&requests.FPDF_GetSignatureCount{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(signatureCount).To(Equal(&responses.FPDF_GetSignatureCount{
					Count: 1,
				}))
			})
		})

		When("the first signature object is opened", func() {
			var signature references.FPDF_SIGNATURE

			BeforeEach(func() {
				signatureResp, err := PdfiumInstance.FPDF_GetSignatureObject(&requests.FPDF_GetSignatureObject{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(signatureResp).To(Not(BeNil()))
				signature = signatureResp.Signature
			})

			It("returns the correct signature content", func() {
				signatureContent, err := PdfiumInstance.FPDFSignatureObj_GetContents(&requests.FPDFSignatureObj_GetContents{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				Expect(signatureContent).To(Equal(&responses.FPDFSignatureObj_GetContents{
					Contents: []byte{0x30, 0x80, 0x06, 0x09, 0x2A, 0x86, 0x48,
						0x86, 0xF7, 0x0D, 0x01, 0x07, 0x02, 0xA0,
						0x80, 0x30, 0x80, 0x02, 0x01, 0x01},
				}))
			})

			It("returns the correct signature byte range", func() {
				signatureByteRange, err := PdfiumInstance.FPDFSignatureObj_GetByteRange(&requests.FPDFSignatureObj_GetByteRange{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				Expect(signatureByteRange).To(Equal(&responses.FPDFSignatureObj_GetByteRange{
					ByteRange: []int{0, 10, 30, 10},
				}))
			})

			It("returns the correct signature sub filter", func() {
				subFilter, err := PdfiumInstance.FPDFSignatureObj_GetSubFilter(&requests.FPDFSignatureObj_GetSubFilter{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				expectedSubfilter := "ETSI.CAdES.detached"
				Expect(subFilter).To(Equal(&responses.FPDFSignatureObj_GetSubFilter{
					SubFilter: &expectedSubfilter,
				}))
			})

			It("returns no signature reason", func() {
				reason, err := PdfiumInstance.FPDFSignatureObj_GetReason(&requests.FPDFSignatureObj_GetReason{
					Signature: signature,
				})

				expectedReason := "test reason"
				Expect(err).To(BeNil())
				Expect(reason).To(Equal(&responses.FPDFSignatureObj_GetReason{
					Reason: &expectedReason,
				}))
			})

			It("returns the correct signature time", func() {
				signatureTime, err := PdfiumInstance.FPDFSignatureObj_GetTime(&requests.FPDFSignatureObj_GetTime{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				expectedTime := "D:20200624093114+02'00'"
				Expect(signatureTime).To(Equal(&responses.FPDFSignatureObj_GetTime{
					Time: &expectedTime,
				}))
			})

			It("returns no DocMDPPermission", func() {
				docMDPPermission, err := PdfiumInstance.FPDFSignatureObj_GetDocMDPPermission(&requests.FPDFSignatureObj_GetDocMDPPermission{
					Signature: signature,
				})
				Expect(err).To(MatchError("could not get DocMDPPermission"))
				Expect(docMDPPermission).To(BeNil())
			})
		})
	})

	Context("a PDF file with signature DocMDPPermission", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/docmdp.pdf")
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

		When("the document count is requested", func() {
			It("returns that it has one signatures", func() {
				signatureCount, err := PdfiumInstance.FPDF_GetSignatureCount(&requests.FPDF_GetSignatureCount{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(signatureCount).To(Equal(&responses.FPDF_GetSignatureCount{
					Count: 1,
				}))
			})
		})

		When("the first signature object is opened", func() {
			var signature references.FPDF_SIGNATURE

			BeforeEach(func() {
				signatureResp, err := PdfiumInstance.FPDF_GetSignatureObject(&requests.FPDF_GetSignatureObject{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(signatureResp).To(Not(BeNil()))
				signature = signatureResp.Signature
			})

			It("returns the correct signature content", func() {
				signatureContent, err := PdfiumInstance.FPDFSignatureObj_GetContents(&requests.FPDFSignatureObj_GetContents{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				Expect(signatureContent).To(Equal(&responses.FPDFSignatureObj_GetContents{}))
			})

			It("returns the correct signature byte range", func() {
				signatureByteRange, err := PdfiumInstance.FPDFSignatureObj_GetByteRange(&requests.FPDFSignatureObj_GetByteRange{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				Expect(signatureByteRange).To(Equal(&responses.FPDFSignatureObj_GetByteRange{}))
			})

			It("returns the correct signature sub filter", func() {
				subFilter, err := PdfiumInstance.FPDFSignatureObj_GetSubFilter(&requests.FPDFSignatureObj_GetSubFilter{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				Expect(subFilter).To(Equal(&responses.FPDFSignatureObj_GetSubFilter{}))
			})

			It("returns no signature reason", func() {
				reason, err := PdfiumInstance.FPDFSignatureObj_GetReason(&requests.FPDFSignatureObj_GetReason{
					Signature: signature,
				})

				Expect(err).To(BeNil())
				Expect(reason).To(Equal(&responses.FPDFSignatureObj_GetReason{}))
			})

			It("returns the correct signature time", func() {
				signatureTime, err := PdfiumInstance.FPDFSignatureObj_GetTime(&requests.FPDFSignatureObj_GetTime{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				Expect(signatureTime).To(Equal(&responses.FPDFSignatureObj_GetTime{}))
			})

			It("returns the correct DocMDPPermission", func() {
				docMDPPermission, err := PdfiumInstance.FPDFSignatureObj_GetDocMDPPermission(&requests.FPDFSignatureObj_GetDocMDPPermission{
					Signature: signature,
				})
				Expect(err).To(BeNil())
				Expect(docMDPPermission).To(Equal(&responses.FPDFSignatureObj_GetDocMDPPermission{
					DocMDPPermission: 1,
				}))
			})
		})
	})
})
