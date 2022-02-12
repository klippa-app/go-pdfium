//go:build windows
// +build windows

package shared_tests

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include <windows.h>
import "C"
import (
	"os"
	"syscall"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdfview_win32", func() {
	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when calling FPDF_RenderPage", func() {
				FPDF_RenderPage, err := PdfiumInstance.FPDF_RenderPage(&requests.FPDF_RenderPage{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDF_RenderPage).To(BeNil())
			})
		})
	})

	Context("a normal PDF file with 1 page", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			if TestType == "multi" {
				Skip("Multi-threaded usage does not support FPDF_RenderPage")
			}

			file, err := os.Open(TestDataPath + "/testdata/test.pdf")
			Expect(err).To(BeNil())

			fileStat, err := file.Stat()
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadCustomDocument(&requests.FPDF_LoadCustomDocument{
				Reader: file,
				Size:   fileStat.Size(),
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document
		})

		AfterEach(func() {
			if TestType == "multi" {
				Skip("Multi-threaded usage does not support FPDF_RenderPage")
			}

			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		When("is opened", func() {
			It("returns the correct page render", func() {
				dc, _, _ := procCreateEnhMetaFileA.Call(uintptr(0), uintptr(0), uintptr(0), uintptr(0))
				Expect(dc).To(Not(BeNil()))

				FPDF_RenderPage, err := PdfiumInstance.FPDF_RenderPage(&requests.FPDF_RenderPage{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					DC:     dc.(C.HDC),
					StartX: 0,
					StartY: 0,
					SizeX:  width,
					SizeY:  height,
					Rotate: enums.FPDF_PAGE_ROTATION_NONE,
					Flags:  enums.FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_RenderPage).To(Equal(&responses.FPDF_RenderPage{}))

				emf, _, _ := procCloseEnhMetaFile.Call(uintptr(dc))
				Expect(emf).To(Not(BeNil()))

				res, _, _ := procDeleteEnhMetaFile.Call(uintptr(emf))
				Expect(res).To(Not(Equal(0)))
			})
		})
	})
})

var (
	modgdi32               = syscall.NewLazyDLL("gdi32.dll")
	procCreateEnhMetaFileA = modgdi32.NewProc("CreateEnhMetaFileA")
	procCloseEnhMetaFile   = modgdi32.NewProc("CloseEnhMetaFile")
	procDeleteEnhMetaFile  = modgdi32.NewProc("DeleteEnhMetaFile")
)
