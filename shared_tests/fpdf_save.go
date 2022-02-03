package shared_tests

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

func RunfpdfSaveTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("fpdf_save", func() {
		Context("no document", func() {
			When("is opened", func() {
				It("returns an error when calling FPDF_SaveAsCopy", func() {
					FPDF_SaveAsCopy, err := pdfiumContainer.FPDF_SaveAsCopy(&requests.FPDF_SaveAsCopy{})
					Expect(err).To(MatchError("document not given"))
					Expect(FPDF_SaveAsCopy).To(BeNil())
				})

				It("returns an error when calling FPDF_SaveWithVersion", func() {
					FPDF_SaveWithVersion, err := pdfiumContainer.FPDF_SaveWithVersion(&requests.FPDF_SaveWithVersion{})
					Expect(err).To(MatchError("document not given"))
					Expect(FPDF_SaveWithVersion).To(BeNil())
				})
			})
		})

		Context("a normal PDF file", func() {
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
				Context("and saved to a byte array", func() {
					It("it returns the correct bytes", func() {
						FPDF_SaveAsCopy, err := pdfiumContainer.FPDF_SaveAsCopy(&requests.FPDF_SaveAsCopy{
							Document: doc,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_SaveAsCopy).To(Not(BeNil()))
						Expect(FPDF_SaveAsCopy.FileBytes).To(Not(BeNil()))
						Expect(FPDF_SaveAsCopy.FileBytes).To(PointTo(HaveLen(11375)))
					})
				})

				Context("and saved to a file path", func() {
					It("it returns the correct bytes", func() {
						tempFile, err := ioutil.TempFile("", "")
						Expect(err).To(BeNil())
						defer tempFile.Close()
						defer os.Remove(tempFile.Name())

						tempFileName := tempFile.Name()
						FPDF_SaveAsCopy, err := pdfiumContainer.FPDF_SaveAsCopy(&requests.FPDF_SaveAsCopy{
							Document: doc,
							FilePath: &tempFileName,
						})

						Expect(err).To(BeNil())
						fileStat, err := tempFile.Stat()
						Expect(err).To(BeNil())
						Expect(FPDF_SaveAsCopy).To(Not(BeNil()))
						Expect(FPDF_SaveAsCopy.FileBytes).To(BeNil())
						Expect(fileStat.Size()).To(Equal(int64(11375)))
					})
				})

				Context("and saved to a file path that does not work", func() {
					It("it returns an error", func() {
						fakeFilePath := "/path/that/will/never/work"
						FPDF_SaveAsCopy, err := pdfiumContainer.FPDF_SaveAsCopy(&requests.FPDF_SaveAsCopy{
							Document: doc,
							FilePath: &fakeFilePath,
						})

						Expect(err).To(MatchError("open /path/that/will/never/work: no such file or directory"))
						Expect(FPDF_SaveAsCopy).To(BeNil())
					})
				})

				// io.Writer not supported on multi-threaded usage.
				if prefix != "multi" {
					Context("and saved to a io.Writer", func() {
						It("it returns the correct bytes", func() {
							buffer := bytes.Buffer{}
							FPDF_SaveAsCopy, err := pdfiumContainer.FPDF_SaveAsCopy(&requests.FPDF_SaveAsCopy{
								Document:   doc,
								FileWriter: &buffer,
							})
							Expect(err).To(BeNil())
							Expect(FPDF_SaveAsCopy).To(Not(BeNil()))
							Expect(FPDF_SaveAsCopy.FileBytes).To(BeNil())
							Expect(buffer.Len()).To(Equal(11375))
						})
					})
				}

				Context("and saved with another PDF version", func() {
					It("it returns the correct byte array", func() {
						FPDF_SaveWithVersion, err := pdfiumContainer.FPDF_SaveWithVersion(&requests.FPDF_SaveWithVersion{
							Document:    doc,
							FileVersion: 13,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_SaveWithVersion).To(Not(BeNil()))
						Expect(FPDF_SaveWithVersion.FileBytes).To(Not(BeNil()))
						Expect(FPDF_SaveWithVersion.FileBytes).To(PointTo(HaveLen(11375)))
					})
				})
			})
		})
	})
}
