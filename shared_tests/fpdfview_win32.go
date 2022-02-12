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
				FPDF_GetPageWidthF, err := PdfiumInstance.FPDF_GetPageWidthF(&requests.FPDF_GetPageWidthF{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetPageWidthF).To(Not(BeNil()))

				FPDF_GetPageHeightF, err := PdfiumInstance.FPDF_GetPageHeightF(&requests.FPDF_GetPageHeightF{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetPageHeightF).To(Not(BeNil()))

				dc := procCreateEnhMetaFileA.Call(nil, nil, nil, nil)
				Expect(dc).To(Not(BeNil()))

				width := int(FPDF_GetPageWidthF.PageWidth)
				height := int(FPDF_GetPageHeightF.PageHeight)
				startX := int(0)
				startY := int(0)

				rgn := procCreateRectRgn.Call(uintptr(startX), uintptr(startY), uintptr(width), uintptr(height))
				Expect(rgn).To(Not(BeNil()))

				selectClip, _, _ := procSelectClipRgn.Call(uintptr(dc), uintptr(rgn))
				Expect(deleteRGN).To(Not(Equal(0)))

				deleteRGN, _, _ := procDeleteObject.Call(uintptr(rgn))
				Expect(deleteRGN).To(Not(Equal(0)))

				WHITE_BRUSH := 0
				NULL_PEN := 8

				nullPen, _, _ := procGetStockObject.Call(uintptr(NULL_PEN))
				Expect(nullPen).To(Not(BeNil()))

				whiteBrush, _, _ := procGetStockObject.Call(uintptr(WHITE_BRUSH))
				Expect(whiteBrush).To(Not(BeNil()))

				nullPenSelect, _, _ := procSelectObject.Call(uintptr(dc), uintptr(nullPen))
				Expect(nullPenSelect).To(Not(Equal(0)))

				whiteBrushSelect, _, _ := procSelectObject.Call(uintptr(dc), uintptr(whiteBrush))
				Expect(whiteBrushSelect).To(Not(Equal(0)))

				rectAngleX1 := int(0)
				rectAngleY1 := int(0)
				rectAngleX2 := int(width + 1)
				rectAngleY2 := int(height + 1)

				// If a PS_NULL pen is used, the dimensions of the rectangle are 1 pixel less.
				rectangleResult, _, _ := procRectangle.Call(uintptr(dc), uintptr(rectAngleX1), uintptr(rectAngleY1), uintptr(rectAngleX2), uintptr(rectAngleY2))
				Expect(rectangleResult).To(Not(Equal(0)))

				FPDF_RenderPage, err := PdfiumInstance.FPDF_RenderPage(&requests.FPDF_RenderPage{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					DC:     dc,
					StartX: 0,
					StartY: 0,
					SizeX:  width,
					SizeY:  height,
					Rotate: enums.FPDF_PAGE_ROTATION_NONE,
					Flags:  enums.FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_RenderPage).To(Equal(&responses.FPDF_RenderPage{}))
			})
		})
	})
})

var (
	modgdi32               = syscall.NewLazyDLL("gdi32.dll")
	procDeleteObject       = modgdi32.NewProc("DeleteObject")
	procSelectObject       = modgdi32.NewProc("SelectObject")
	procGetStockObject     = modgdi32.NewProc("GetStockObject")
	procRectangle          = modgdi32.NewProc("Rectangle")
	procCreateRectRgn      = modgdi32.NewProc("CreateRectRgn")
	procSelectClipRgn      = modgdi32.NewProc("SelectClipRgn")
	procCreateEnhMetaFileA = modgdi32.NewProc("CreateEnhMetaFileA")
)
