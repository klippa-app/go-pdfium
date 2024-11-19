package shared_tests

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("fpdf_save", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when calling FPDF_SaveAsCopy", func() {
				FPDF_SaveAsCopy, err := PdfiumInstance.FPDF_SaveAsCopy(&requests.FPDF_SaveAsCopy{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_SaveAsCopy).To(BeNil())
			})

			It("returns an error when calling FPDF_SaveWithVersion", func() {
				FPDF_SaveWithVersion, err := PdfiumInstance.FPDF_SaveWithVersion(&requests.FPDF_SaveWithVersion{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_SaveWithVersion).To(BeNil())
			})
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
			Context("and saved to a byte array", func() {
				It("it returns the correct bytes", func() {
					FPDF_SaveAsCopy, err := PdfiumInstance.FPDF_SaveAsCopy(&requests.FPDF_SaveAsCopy{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_SaveAsCopy).To(Not(BeNil()))
					Expect(FPDF_SaveAsCopy.FileBytes).To(Not(BeNil()))
					Expect(FPDF_SaveAsCopy.FileBytes).To(SatisfyAny(PointTo(HaveLen(11375)), PointTo(HaveLen(11183)), PointTo(HaveLen(11188)))) // 11375 < Pdfium 5854, 11183 >= Pdfium 5854, 11188 => Pdfium 6721
				})
			})

			Context("and saved to a file path", func() {
				It("it returns the correct bytes", func() {
					tempFile, err := ioutil.TempFile("", "")
					Expect(err).To(BeNil())
					defer tempFile.Close()
					defer os.Remove(tempFile.Name())

					tempFileName := tempFile.Name()
					FPDF_SaveAsCopy, err := PdfiumInstance.FPDF_SaveAsCopy(&requests.FPDF_SaveAsCopy{
						Document: doc,
						FilePath: &tempFileName,
					})

					Expect(err).To(BeNil())
					fileStat, err := tempFile.Stat()
					Expect(err).To(BeNil())
					Expect(FPDF_SaveAsCopy).To(Not(BeNil()))
					Expect(FPDF_SaveAsCopy.FileBytes).To(BeNil())
					Expect(fileStat.Size()).To(SatisfyAny(Equal(int64(11375)), Equal(int64(11183)), Equal(int64(11188)))) // 11375 < Pdfium 5854, 11183 >= Pdfium 5854, 11188 => Pdfium 6721
				})
			})

			Context("and saved to a file path that does not work", func() {
				It("it returns an error", func() {
					fakeFilePath := "/path/that/will/never/work"
					FPDF_SaveAsCopy, err := PdfiumInstance.FPDF_SaveAsCopy(&requests.FPDF_SaveAsCopy{
						Document: doc,
						FilePath: &fakeFilePath,
					})

					Expect(err).To(Not(BeNil()))
					Expect(FPDF_SaveAsCopy).To(BeNil())
				})
			})

			Context("and saved to a io.Writer", func() {
				BeforeEach(func() {
					if TestType == "multi" {
						Skip("Multi-threaded usage does not support io.Writer")
					}
				})

				It("it returns the correct bytes", func() {
					buffer := bytes.Buffer{}
					FPDF_SaveAsCopy, err := PdfiumInstance.FPDF_SaveAsCopy(&requests.FPDF_SaveAsCopy{
						Document:   doc,
						FileWriter: &buffer,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_SaveAsCopy).To(Not(BeNil()))
					Expect(FPDF_SaveAsCopy.FileBytes).To(BeNil())
					Expect(buffer.Len()).To(SatisfyAny(Equal(11375), Equal(11183), Equal(11188))) // 11375 < Pdfium 5854, 11183 >= Pdfium 5854, 11188 => Pdfium 6721
				})
			})

			Context("and saved with another PDF version", func() {
				It("it returns the correct byte array and the result has the correct version", func() {
					FPDF_SaveWithVersion, err := PdfiumInstance.FPDF_SaveWithVersion(&requests.FPDF_SaveWithVersion{
						Document:    doc,
						FileVersion: 13,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_SaveWithVersion).To(Not(BeNil()))
					Expect(FPDF_SaveWithVersion.FileBytes).To(Not(BeNil()))
					Expect(FPDF_SaveWithVersion.FileBytes).To(SatisfyAny(PointTo(HaveLen(11375)), PointTo(HaveLen(11183)), PointTo(HaveLen(11188)))) // 11375 < Pdfium 5854, 11183 >= Pdfium 5854, 11188 => Pdfium 6721

					savedDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
						Data: FPDF_SaveWithVersion.FileBytes,
					})
					Expect(err).To(BeNil())

					fileVersion, err := PdfiumInstance.FPDF_GetFileVersion(&requests.FPDF_GetFileVersion{
						Document: savedDoc.Document,
					})
					Expect(err).To(BeNil())
					Expect(fileVersion).To(Equal(&responses.FPDF_GetFileVersion{
						FileVersion: 13,
					}))

					FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
						Document: savedDoc.Document,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_CloseDocument).To(Not(BeNil()))
				})
			})

			Context("and saved with incremental flag", func() {
				It("it returns the correct byte array", func() {
					FPDF_SaveWithVersion, err := PdfiumInstance.FPDF_SaveWithVersion(&requests.FPDF_SaveWithVersion{
						Document:    doc,
						FileVersion: 13,
						Flags:       requests.SaveFlagIncremental,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_SaveWithVersion).To(Not(BeNil()))
					Expect(FPDF_SaveWithVersion.FileBytes).To(Not(BeNil()))
					Expect(FPDF_SaveWithVersion.FileBytes).To(PointTo(HaveLen(11780)))
				})
			})
		})
	})
})
