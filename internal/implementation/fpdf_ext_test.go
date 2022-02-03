package implementation_test

import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_ext", func() {
	pdfium := implementation.Pdfium.GetInstance()

	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when getting the page mode", func() {
				pageMode, err := pdfium.FPDFDoc_GetPageMode(&requests.FPDFDoc_GetPageMode{})
				Expect(err).To(MatchError("document not given"))
				Expect(pageMode).To(BeNil())
			})
		})
	})

	Context("the time function has been overwritten", func() {
		BeforeEach(func() {
			var alwaysTheSameTimeFunction requests.SetTimeFunction
			alwaysTheSameTimeFunction = func() int64 {
				return 123456
			}

			resp, err := pdfium.FSDK_SetTimeFunction(&requests.FSDK_SetTimeFunction{
				Function: alwaysTheSameTimeFunction,
			})
			Expect(err).To(BeNil())
			Expect(resp).To(Equal(&responses.FSDK_SetTimeFunction{}))
		})

		AfterEach(func() {
			resp, err := pdfium.FSDK_SetTimeFunction(&requests.FSDK_SetTimeFunction{
				Function: nil,
			})
			Expect(err).To(BeNil())
			Expect(resp).To(Equal(&responses.FSDK_SetTimeFunction{}))
		})

		When("a new document is created", func() {
			It("returns the correct CreationDate date in the metadata", func() {
				newDoc, err := pdfium.FPDF_CreateNewDocument(&requests.FPDF_CreateNewDocument{})
				Expect(err).To(BeNil())
				Expect(newDoc).To(Not(BeNil()))

				meta, err := pdfium.GetMetaData(&requests.GetMetaData{
					Document: newDoc.Document,
				})
				Expect(err).To(BeNil())
				Expect(meta).To(Equal(&responses.GetMetaData{
					Tags: []responses.GetMetaDataTag{
						{Tag: "Title", Value: ""},
						{Tag: "Author", Value: ""},
						{Tag: "Subject", Value: ""},
						{Tag: "Keywords", Value: ""},
						{Tag: "Creator", Value: "PDFium"},
						{Tag: "Producer", Value: ""},
						{
							Tag:   "CreationDate",
							Value: "D:19700102101736",
						},
						{Tag: "ModDate", Value: ""},
					},
				}))
			})
		})
	})

	Context("the localtime function has been overwritten", func() {
		BeforeEach(func() {
			var alwaysTheSameLocalTimeFunction requests.SetLocaltimeFunction
			alwaysTheSameLocalTimeFunction = func(currentUnixTime int64) requests.SetLocaltime {
				return requests.SetLocaltime{
					TmSec:   30,
					TmMin:   30,
					TmHour:  12,
					TmMday:  231,
					TmMon:   8,
					TmYear:  1992,
					TmWday:  2,
					TmYday:  18,
					TmIsdst: 0,
				}
			}

			resp, err := pdfium.FSDK_SetLocaltimeFunction(&requests.FSDK_SetLocaltimeFunction{
				Function: alwaysTheSameLocalTimeFunction,
			})
			Expect(err).To(BeNil())
			Expect(resp).To(Equal(&responses.FSDK_SetLocaltimeFunction{}))
		})

		AfterEach(func() {
			resp, err := pdfium.FSDK_SetLocaltimeFunction(&requests.FSDK_SetLocaltimeFunction{
				Function: nil,
			})
			Expect(err).To(BeNil())
			Expect(resp).To(Equal(&responses.FSDK_SetLocaltimeFunction{}))
		})

		When("a new document is created", func() {
			It("returns the correct CreationDate date in the metadata", func() {
				newDoc, err := pdfium.FPDF_CreateNewDocument(&requests.FPDF_CreateNewDocument{})
				Expect(err).To(BeNil())
				Expect(newDoc).To(Not(BeNil()))

				meta, err := pdfium.GetMetaData(&requests.GetMetaData{
					Document: newDoc.Document,
				})
				Expect(err).To(BeNil())
				Expect(meta).To(Equal(&responses.GetMetaData{
					Tags: []responses.GetMetaDataTag{
						{Tag: "Title", Value: ""},
						{Tag: "Author", Value: ""},
						{Tag: "Subject", Value: ""},
						{Tag: "Keywords", Value: ""},
						{Tag: "Creator", Value: "PDFium"},
						{Tag: "Producer", Value: ""},
						{
							Tag:   "CreationDate",
							Value: "D:389209231123030",
						},
						{Tag: "ModDate", Value: ""},
					},
				}))
			})
		})
	})

	Context("the unsupported object processor function has been overwritten", func() {
		var lastReportedUnsupportedObject enums.FPDF_UNSP
		BeforeEach(func() {
			var reportUnsupportedObj requests.UnSpObjProcessHandler
			reportUnsupportedObj = func(obj enums.FPDF_UNSP) {
				lastReportedUnsupportedObject = obj
			}

			resp, err := pdfium.FSDK_SetUnSpObjProcessHandler(&requests.FSDK_SetUnSpObjProcessHandler{
				UnSpObjProcessHandler: reportUnsupportedObj,
			})
			Expect(err).To(BeNil())
			Expect(resp).To(Equal(&responses.FSDK_SetUnSpObjProcessHandler{}))
		})

		AfterEach(func() {
			resp, err := pdfium.FSDK_SetUnSpObjProcessHandler(&requests.FSDK_SetUnSpObjProcessHandler{
				UnSpObjProcessHandler: nil,
			})
			Expect(err).To(BeNil())
			Expect(resp).To(Equal(&responses.FSDK_SetUnSpObjProcessHandler{}))
		})

		When("a document with unsupported objects is opened", func() {
			It("reports the unsupported objects to the handler", func() {
				pdfData, err := ioutil.ReadFile("../../shared_tests/testdata/unsupported_feature.pdf")
				Expect(err).To(BeNil())
				if err != nil {
					return
				}

				newDoc, err := pdfium.OpenDocument(&requests.OpenDocument{
					File: &pdfData,
				})
				if err != nil {
					return
				}
				Expect(err).To(BeNil())
				Expect(newDoc).To(Not(BeNil()))
				Expect(lastReportedUnsupportedObject).To(Equal(enums.FPDF_UNSP_DOC_PORTABLECOLLECTION))
			})
		})
	})
})
