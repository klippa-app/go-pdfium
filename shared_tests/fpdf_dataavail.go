package shared_tests

import (
	"fmt"
	"os"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type FakeReadSeeker struct {
	FileData       []byte
	OriginalReader *os.File
	Size           int64
	CurrentPos     int64
	LoadedBytes    []bool
}

func (f *FakeReadSeeker) Read(p []byte) (n int, err error) {
	size := len(p)
	start := f.CurrentPos
	bytesRead := 0

	for i := start; i < start+int64(size); i++ {
		if !f.LoadedBytes[i] {
			return 0, fmt.Errorf("byte %d i not loaded yet", i)
		}

		p[bytesRead] = f.FileData[i]

		f.CurrentPos++
		bytesRead++
	}

	return bytesRead, nil
}

func (f *FakeReadSeeker) Seek(offset int64, whence int) (int64, error) {
	f.CurrentPos = offset
	return offset, nil
}

func (f *FakeReadSeeker) IsDataAvailableCallback(offset, size uint64) bool {
	for i := offset; i < offset+size; i++ {
		if !f.LoadedBytes[i] {
			return false
		}
	}

	return true
}

func (f *FakeReadSeeker) AddSegmentCallback(offset, size uint64) {
	res := make([]byte, size)
	amountRead, err := f.OriginalReader.ReadAt(res, int64(offset))
	if err != nil {
		return
	}

	for i := 0; i < amountRead; i++ {
		f.LoadedBytes[int(offset)+i] = true
		f.FileData[int(offset)+i] = res[i]
	}
}

var _ = Describe("fpdf_dataavail", func() {
	BeforeEach(func() {
		if TestType == "multi" {
			Skip("Multi-threaded usage does not support setting callbacks")
		}

		if TestType == "webassembly" {
			// @todo: remove me when implemented.
			Skip("This test is skipped on Webassembly")
		}
		Locker.Lock()
	})

	AfterEach(func() {
		if TestType == "multi" {
			Skip("Multi-threaded usage does not support setting callbacks")
		}

		if TestType == "webassembly" {
			// @todo: remove me when implemented.
			Skip("This test is skipped on Webassembly")
		}
		Locker.Unlock()
	})

	Context("no reader", func() {
		When("is set", func() {
			It("returns an error when calling FPDFAvail_Create", func() {
				FPDFAvail_Create, err := PdfiumInstance.FPDFAvail_Create(&requests.FPDFAvail_Create{
					IsDataAvailableCallback: func(offset, size uint64) bool {
						return false
					},
				})
				Expect(err).To(MatchError("Reader can't be nil"))
				Expect(FPDFAvail_Create).To(BeNil())
			})
		})
	})

	Context("no size", func() {
		When("is set", func() {
			It("returns an error when calling FPDFAvail_Create", func() {
				FPDFAvail_Create, err := PdfiumInstance.FPDFAvail_Create(&requests.FPDFAvail_Create{
					IsDataAvailableCallback: func(offset, size uint64) bool {
						return false
					},
					Reader: &FakeReadSeeker{},
				})
				Expect(err).To(MatchError("Size should be set"))
				Expect(FPDFAvail_Create).To(BeNil())
			})
		})
	})

	Context("no document", func() {
		When("is set", func() {
			It("returns an error when calling FPDFAvail_GetFirstPageNum", func() {
				FPDFAvail_GetFirstPageNum, err := PdfiumInstance.FPDFAvail_GetFirstPageNum(&requests.FPDFAvail_GetFirstPageNum{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDFAvail_GetFirstPageNum).To(BeNil())
			})
		})
	})

	Context("no callback", func() {
		When("is set", func() {
			It("returns an error when calling FPDFAvail_Create", func() {
				FPDFAvail_Create, err := PdfiumInstance.FPDFAvail_Create(&requests.FPDFAvail_Create{})
				Expect(err).To(MatchError("IsDataAvailableCallback can't be nil"))
				Expect(FPDFAvail_Create).To(BeNil())
			})
		})
	})

	Context("no availability provider", func() {
		When("is set", func() {
			It("returns an error when calling FPDFAvail_IsDocAvail", func() {
				FPDFAvail_IsDocAvail, err := PdfiumInstance.FPDFAvail_IsDocAvail(&requests.FPDFAvail_IsDocAvail{})
				Expect(err).To(MatchError("dataAvail not given"))
				Expect(FPDFAvail_IsDocAvail).To(BeNil())
			})

			It("returns an error when calling FPDFAvail_GetDocument", func() {
				FPDFAvail_GetDocument, err := PdfiumInstance.FPDFAvail_GetDocument(&requests.FPDFAvail_GetDocument{})
				Expect(err).To(MatchError("dataAvail not given"))
				Expect(FPDFAvail_GetDocument).To(BeNil())
			})

			It("returns an error when calling FPDFAvail_IsPageAvail", func() {
				FPDFAvail_IsPageAvail, err := PdfiumInstance.FPDFAvail_IsPageAvail(&requests.FPDFAvail_IsPageAvail{})
				Expect(err).To(MatchError("dataAvail not given"))
				Expect(FPDFAvail_IsPageAvail).To(BeNil())
			})

			It("returns an error when calling FPDFAvail_Destroy", func() {
				FPDFAvail_Destroy, err := PdfiumInstance.FPDFAvail_Destroy(&requests.FPDFAvail_Destroy{})
				Expect(err).To(MatchError("dataAvail not given"))
				Expect(FPDFAvail_Destroy).To(BeNil())
			})

			It("returns an error when calling FPDFAvail_IsFormAvail", func() {
				FPDFAvail_IsFormAvail, err := PdfiumInstance.FPDFAvail_IsFormAvail(&requests.FPDFAvail_IsFormAvail{})
				Expect(err).To(MatchError("dataAvail not given"))
				Expect(FPDFAvail_IsFormAvail).To(BeNil())
			})

			It("returns an error when calling FPDFAvail_IsLinearized", func() {
				FPDFAvail_IsLinearized, err := PdfiumInstance.FPDFAvail_IsLinearized(&requests.FPDFAvail_IsLinearized{})
				Expect(err).To(MatchError("dataAvail not given"))
				Expect(FPDFAvail_IsLinearized).To(BeNil())
			})
		})
	})

	Context("a normal PDF file", func() {
		var avail references.FPDF_AVAIL
		var fakeReadSeeker *FakeReadSeeker

		BeforeEach(func() {
			pdfFile, err := os.Open(TestDataPath + "/testdata/test.pdf")
			Expect(err).To(BeNil())

			stat, err := pdfFile.Stat()
			Expect(err).To(BeNil())

			fakeReadSeeker = &FakeReadSeeker{
				Size:           stat.Size(),
				OriginalReader: pdfFile,
				FileData:       make([]byte, stat.Size()),
				CurrentPos:     0,
				LoadedBytes:    make([]bool, stat.Size()),
			}

			newAvail, err := PdfiumInstance.FPDFAvail_Create(&requests.FPDFAvail_Create{
				Reader:                  fakeReadSeeker,
				Size:                    fakeReadSeeker.Size,
				IsDataAvailableCallback: fakeReadSeeker.IsDataAvailableCallback,
				AddSegmentCallback:      fakeReadSeeker.AddSegmentCallback,
			})
			Expect(err).To(BeNil())

			avail = newAvail.AvailabilityProvider
		})

		AfterEach(func() {
			FPDFAvail_Destroy, err := PdfiumInstance.FPDFAvail_Destroy(&requests.FPDFAvail_Destroy{
				AvailabilityProvider: avail,
			})
			Expect(err).To(BeNil())
			Expect(FPDFAvail_Destroy).To(Not(BeNil()))

			err = fakeReadSeeker.OriginalReader.Close()
			Expect(err).To(BeNil())
		})

		When("is opened with a data availability provider", func() {
			When("no data has been loaded yet", func() {
				It("returns that the document is not available yet", func() {
					FPDFAvail_IsDocAvail, err := PdfiumInstance.FPDFAvail_IsDocAvail(&requests.FPDFAvail_IsDocAvail{
						AvailabilityProvider: avail,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAvail_IsDocAvail).To(Equal(&responses.FPDFAvail_IsDocAvail{
						IsDocAvail: enums.PDF_FILEAVAIL_DATA_NOTAVAIL,
					}))
				})

				It("returns an error the page is not available yet (because the document has not been loaded yet)", func() {
					FPDFAvail_IsPageAvail, err := PdfiumInstance.FPDFAvail_IsPageAvail(&requests.FPDFAvail_IsPageAvail{
						AvailabilityProvider: avail,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAvail_IsPageAvail).To(Equal(&responses.FPDFAvail_IsPageAvail{
						IsPageAvail: enums.PDF_FILEAVAIL_DATA_ERROR,
					}))
				})
			})

			When("no data has been loaded yet", func() {
				It("returns that the document is ready to be loaded after a few hints", func() {
					By("We don't know if the document is linearized")
					FPDFAvail_IsLinearized, err := PdfiumInstance.FPDFAvail_IsLinearized(&requests.FPDFAvail_IsLinearized{
						AvailabilityProvider: avail,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAvail_IsLinearized).To(Equal(&responses.FPDFAvail_IsLinearized{
						IsLinearized: enums.PDF_FILEAVAIL_LINEARIZATION_UNKNOWN,
					}))

					By("Avail check 1, not ready")
					FPDFAvail_IsDocAvail, err := PdfiumInstance.FPDFAvail_IsDocAvail(&requests.FPDFAvail_IsDocAvail{
						AvailabilityProvider: avail,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAvail_IsDocAvail).To(Equal(&responses.FPDFAvail_IsDocAvail{
						IsDocAvail: enums.PDF_FILEAVAIL_DATA_NOTAVAIL,
					}))

					By("Avail check 2, not ready")
					FPDFAvail_IsDocAvail, err = PdfiumInstance.FPDFAvail_IsDocAvail(&requests.FPDFAvail_IsDocAvail{
						AvailabilityProvider: avail,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAvail_IsDocAvail).To(Equal(&responses.FPDFAvail_IsDocAvail{
						IsDocAvail: enums.PDF_FILEAVAIL_DATA_NOTAVAIL,
					}))

					By("Avail check 3, not ready")
					FPDFAvail_IsDocAvail, err = PdfiumInstance.FPDFAvail_IsDocAvail(&requests.FPDFAvail_IsDocAvail{
						AvailabilityProvider: avail,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAvail_IsDocAvail).To(Equal(&responses.FPDFAvail_IsDocAvail{
						IsDocAvail: enums.PDF_FILEAVAIL_DATA_NOTAVAIL,
					}))

					By("Avail check 4 ready")
					FPDFAvail_IsDocAvail, err = PdfiumInstance.FPDFAvail_IsDocAvail(&requests.FPDFAvail_IsDocAvail{
						AvailabilityProvider: avail,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAvail_IsDocAvail).To(Equal(&responses.FPDFAvail_IsDocAvail{
						IsDocAvail: enums.PDF_FILEAVAIL_DATA_AVAIL,
					}))

					By("Avail check 4 ready")
					FPDFAvail_IsDocAvail, err = PdfiumInstance.FPDFAvail_IsDocAvail(&requests.FPDFAvail_IsDocAvail{
						AvailabilityProvider: avail,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAvail_IsDocAvail).To(Equal(&responses.FPDFAvail_IsDocAvail{
						IsDocAvail: enums.PDF_FILEAVAIL_DATA_AVAIL,
					}))

					By("Document is not linearized")
					FPDFAvail_IsLinearized, err = PdfiumInstance.FPDFAvail_IsLinearized(&requests.FPDFAvail_IsLinearized{
						AvailabilityProvider: avail,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAvail_IsLinearized).To(Equal(&responses.FPDFAvail_IsLinearized{
						IsLinearized: enums.PDF_FILEAVAIL_LINEARIZATION_NOT_LINEARIZED,
					}))

					By("Document can be loaded")
					FPDFAvail_GetDocument, err := PdfiumInstance.FPDFAvail_GetDocument(&requests.FPDFAvail_GetDocument{
						AvailabilityProvider: avail,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAvail_GetDocument).To(Not(BeNil()))

					By("returning the correct first page num")
					FPDFAvail_GetFirstPageNum, err := PdfiumInstance.FPDFAvail_GetFirstPageNum(&requests.FPDFAvail_GetFirstPageNum{
						Document: FPDFAvail_GetDocument.Document,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAvail_GetFirstPageNum).To(Equal(&responses.FPDFAvail_GetFirstPageNum{
						FirstPageNum: 0,
					}))

					By("returning that we don't have a form")
					FPDFAvail_IsFormAvail, err := PdfiumInstance.FPDFAvail_IsFormAvail(&requests.FPDFAvail_IsFormAvail{
						AvailabilityProvider: avail,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAvail_IsFormAvail).To(Equal(&responses.FPDFAvail_IsFormAvail{
						IsFormAvail: enums.PDF_FILEAVAIL_FORM_NOTEXIST,
					}))

					By("returning that the page is not ready, attempt 1")
					FPDFAvail_IsPageAvail, err := PdfiumInstance.FPDFAvail_IsPageAvail(&requests.FPDFAvail_IsPageAvail{
						AvailabilityProvider: avail,
						PageIndex:            0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAvail_IsPageAvail).To(Equal(&responses.FPDFAvail_IsPageAvail{
						IsPageAvail: enums.PDF_FILEAVAIL_DATA_NOTAVAIL,
					}))

					By("returning that the page is not ready, attempt 2")
					FPDFAvail_IsPageAvail, err = PdfiumInstance.FPDFAvail_IsPageAvail(&requests.FPDFAvail_IsPageAvail{
						AvailabilityProvider: avail,
						PageIndex:            0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAvail_IsPageAvail).To(Equal(&responses.FPDFAvail_IsPageAvail{
						IsPageAvail: enums.PDF_FILEAVAIL_DATA_NOTAVAIL,
					}))

					By("returning that the page is ready, attempt 3")
					FPDFAvail_IsPageAvail, err = PdfiumInstance.FPDFAvail_IsPageAvail(&requests.FPDFAvail_IsPageAvail{
						AvailabilityProvider: avail,
						PageIndex:            0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFAvail_IsPageAvail).To(Equal(&responses.FPDFAvail_IsPageAvail{
						IsPageAvail: enums.PDF_FILEAVAIL_DATA_AVAIL,
					}))

					By("confirming that we can actually get page information")
					FPDFPage_GetRotation, err := PdfiumInstance.FPDFPage_GetRotation(&requests.FPDFPage_GetRotation{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: FPDFAvail_GetDocument.Document,
								Index:    0,
							},
						},
					})
					Expect(err).To(BeNil())
					Expect(FPDFPage_GetRotation).To(Equal(&responses.FPDFPage_GetRotation{
						Page:         0,
						PageRotation: enums.FPDF_PAGE_ROTATION_NONE,
					}))

					By("Document can be closed")
					FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
						Document: FPDFAvail_GetDocument.Document,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_CloseDocument).To(Not(BeNil()))
				})
			})
		})
	})
})
