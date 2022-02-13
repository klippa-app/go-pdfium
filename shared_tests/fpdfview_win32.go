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
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdfview_win32", func() {
	BeforeEach(func() {
		if TestType == "multi" {
			Skip("Multi-threaded usage does not support setting callbacks")
		}
		Locker.Lock()
	})

	AfterEach(func() {
		if TestType == "multi" {
			Skip("Multi-threaded usage does not support setting callbacks")
		}
		Locker.Unlock()
	})

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

				dcPointer := (C.HDC)(unsafe.Pointer(dc))

				FPDF_RenderPage, err := PdfiumInstance.FPDF_RenderPage(&requests.FPDF_RenderPage{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					DC:     dcPointer,
					StartX: 0,
					StartY: 0,
					SizeX:  2000,
					SizeY:  2000,
					Rotate: enums.FPDF_PAGE_ROTATION_NONE,
					Flags:  enums.FPDF_RENDER_FLAG_ANNOT | enums.FPDF_RENDER_FLAG_PRINTING,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_RenderPage).To(Equal(&responses.FPDF_RenderPage{}))

				emf, _, _ := procCloseEnhMetaFile.Call(uintptr(dc))
				Expect(emf).To(Not(BeNil()))

				fileSize, _, _ := procGetEnhMetaFileBits.Call(uintptr(dc), uintptr(0), uintptr(0))
				Expect(int(fileSize)).To(Not(Equal(int(0))))

				uintFileSize := uint(fileSize)
				buffer := make([]byte, int(uintFileSize))
				writtenFileSize, _, _ := procGetEnhMetaFileBits.Call(uintptr(dc), uintptr(uintFileSize), uintptr(unsafe.Pointer(&buffer[0])))
				Expect(int(writtenFileSize)).To(Not(Equal(int(0))))
				Expect(buffer).To(Equal([]byte{}))

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
	procGetEnhMetaFileBits = modgdi32.NewProc("GetEnhMetaFileBits")
)
