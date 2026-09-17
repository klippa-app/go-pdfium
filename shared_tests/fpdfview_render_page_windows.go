//go:build windows

package shared_tests

import (
	"os"
	"syscall"
	"unsafe"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
)

// FPDF_RenderPage draws to a Windows device context. The handle is passed
// the way the Windows API hands it out, so a caller never needs the C type.
var _ = Describe("fpdfview FPDF_RenderPage", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	Context("a normal PDF file", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := os.ReadFile(TestDataPath + "/testdata/test.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document
		})

		AfterEach(func() {
			_, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
		})

		It("returns an error when the DC is not a handle", func() {
			_, err := PdfiumInstance.FPDF_RenderPage(&requests.FPDF_RenderPage{
				DC:   "not a handle",
				Page: requests.Page{ByIndex: &requests.PageByIndex{Document: doc, Index: 0}},
			})
			Expect(err).To(MatchError("DC must be a Windows device context handle given as uintptr, syscall.Handle or unsafe.Pointer"))
		})

		It("returns an error when the DC is null", func() {
			_, err := PdfiumInstance.FPDF_RenderPage(&requests.FPDF_RenderPage{
				DC:   uintptr(0),
				Page: requests.Page{ByIndex: &requests.PageByIndex{Document: doc, Index: 0}},
			})
			Expect(err).To(MatchError("DC is a null device context handle"))
		})

		It("renders to a memory device context given as uintptr, syscall.Handle and unsafe.Pointer", func() {
			// A memory DC compatible with the screen, the same thing an
			// application would create to render into a bitmap.
			gdi32 := syscall.NewLazyDLL("gdi32.dll")
			createCompatibleDC := gdi32.NewProc("CreateCompatibleDC")
			deleteDC := gdi32.NewProc("DeleteDC")

			hdc, _, _ := createCompatibleDC.Call(0)
			Expect(hdc).ToNot(BeZero())
			defer deleteDC.Call(hdc)

			// The handle is not a Go pointer, so reinterpret the value rather
			// than converting it, which go vet would flag.
			hdcPointer := *(*unsafe.Pointer)(unsafe.Pointer(&hdc))
			for _, dc := range []any{hdc, syscall.Handle(hdc), hdcPointer} {
				resp, err := PdfiumInstance.FPDF_RenderPage(&requests.FPDF_RenderPage{
					DC:     dc,
					Page:   requests.Page{ByIndex: &requests.PageByIndex{Document: doc, Index: 0}},
					StartX: 0,
					StartY: 0,
					SizeX:  100,
					SizeY:  100,
					Rotate: enums.FPDF_PAGE_ROTATION_NONE,
					Flags:  enums.FPDF_RENDER_FLAG_ANNOT,
				})
				Expect(err).To(BeNil())
				Expect(resp).ToNot(BeNil())
			}
		})
	})
})
