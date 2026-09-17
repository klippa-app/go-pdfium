//go:build !windows

package shared_tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
)

var _ = Describe("fpdfview FPDF_RenderPage", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	It("is reported as unsupported outside Windows", func() {
		_, err := PdfiumInstance.FPDF_RenderPage(&requests.FPDF_RenderPage{
			DC:   uintptr(1),
			Page: requests.Page{ByIndex: &requests.PageByIndex{Index: 0}},
		})
		Expect(err).To(MatchError(pdfium_errors.ErrWindowsUnsupported))
	})
})
