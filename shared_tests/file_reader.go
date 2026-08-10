package shared_tests

import (
	"bytes"
	"io"
	"os"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// eofWithFinalBytesReader returns io.EOF together with the last bytes of the
// underlying data instead of on a following zero-byte read. Some io.ReadSeeker
// implementations in the wild behave this way, and it is allowed by the
// io.Reader contract.
type eofWithFinalBytesReader struct {
	*bytes.Reader
}

func (r *eofWithFinalBytesReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if err == nil && r.Len() == 0 {
		err = io.EOF
	}
	return n, err
}

// chunkedReader never returns more than chunkSize bytes per Read call, without
// returning an error. A single Read call leaves the requested range only
// partially filled.
type chunkedReader struct {
	*bytes.Reader
	chunkSize int
}

func (r *chunkedReader) Read(p []byte) (int, error) {
	if len(p) > r.chunkSize {
		p = p[:r.chunkSize]
	}
	return r.Reader.Read(p)
}

// The file access callbacks have to hand pdfium a completely filled buffer for
// every block it requests. They do that with io.ReadFull, which keeps reading
// until the buffer is full and treats an io.EOF that comes back together with
// the final bytes as a successful read. These tests make sure documents still
// open for readers that trigger both of those cases, and that a reader which
// genuinely cannot deliver the requested range is still reported as a failure.
var _ = Describe("file readers", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	Context("a normal PDF file with 2 pages", func() {
		var pdfData []byte

		BeforeEach(func() {
			var err error
			pdfData, err = os.ReadFile(TestDataPath + "/testdata/test_multipage.pdf")
			Expect(err).To(BeNil())
		})

		openAndCountPages := func(reader io.ReadSeeker, size int64) {
			doc, err := PdfiumInstance.FPDF_LoadCustomDocument(&requests.FPDF_LoadCustomDocument{
				Reader: reader,
				Size:   size,
			})
			Expect(err).To(BeNil())
			Expect(doc).To(Not(BeNil()))

			var document references.FPDF_DOCUMENT
			if doc != nil {
				document = doc.Document
			}

			defer func() {
				_, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
					Document: document,
				})
				Expect(err).To(BeNil())
			}()

			pageCount, err := PdfiumInstance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
				Document: document,
			})
			Expect(err).To(BeNil())
			Expect(pageCount).To(Not(BeNil()))
			Expect(pageCount.PageCount).To(Equal(2))
		}

		When("is opened with a reader that returns io.EOF together with the final bytes", func() {
			It("returns the correct page count", func() {
				reader := &eofWithFinalBytesReader{Reader: bytes.NewReader(pdfData)}
				openAndCountPages(reader, int64(len(pdfData)))
			})
		})

		When("is opened with a reader that returns short reads", func() {
			It("returns the correct page count", func() {
				reader := &chunkedReader{Reader: bytes.NewReader(pdfData), chunkSize: 7}
				openAndCountPages(reader, int64(len(pdfData)))
			})
		})

		When("is opened with a reader that can't deliver the whole requested range", func() {
			It("returns an error instead of handing pdfium a partial buffer", func() {
				if TestType == "multi" {
					Skip("Multi-threaded usage reads the file reader into memory, so it never issues partial reads")
				}

				// Only the first half of the file is available, while pdfium
				// is told the file has its full size. Every read past the
				// halfway point can't be filled completely, which has to be
				// reported to pdfium as a failed read.
				reader := bytes.NewReader(pdfData[:len(pdfData)/2])

				doc, err := PdfiumInstance.FPDF_LoadCustomDocument(&requests.FPDF_LoadCustomDocument{
					Reader: reader,
					Size:   int64(len(pdfData)),
				})
				if doc != nil {
					defer PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
						Document: doc.Document,
					})
				}
				Expect(err).To(HaveOccurred())
				Expect(doc).To(BeNil())
			})
		})
	})
})
