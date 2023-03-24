package shared_tests

import (
	"os"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/internal/implementation_webassembly"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdfview", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	It("allows for the last error to be fetched", func() {
		FPDF_GetLastError, err := PdfiumInstance.FPDF_GetLastError(&requests.FPDF_GetLastError{})
		Expect(err).To(BeNil())
		Expect(FPDF_GetLastError).To(Not(BeNil()))
	})

	It("allows for the sandbox policy to be enabled", func() {
		FPDF_SetSandBoxPolicy, err := PdfiumInstance.FPDF_SetSandBoxPolicy(&requests.FPDF_SetSandBoxPolicy{
			Policy: requests.FPDF_SetSandBoxPolicyPolicyMachinetimeAccess,
			Enable: true,
		})
		Expect(err).To(BeNil())
		Expect(FPDF_SetSandBoxPolicy).To(Equal(&responses.FPDF_SetSandBoxPolicy{}))
	})

	It("allows for the sandbox policy to be disabled", func() {
		FPDF_SetSandBoxPolicy, err := PdfiumInstance.FPDF_SetSandBoxPolicy(&requests.FPDF_SetSandBoxPolicy{
			Policy: requests.FPDF_SetSandBoxPolicyPolicyMachinetimeAccess,
			Enable: false,
		})
		Expect(err).To(BeNil())
		Expect(FPDF_SetSandBoxPolicy).To(Equal(&responses.FPDF_SetSandBoxPolicy{}))
	})

	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when getting the pdf version", func() {
				pageCount, err := PdfiumInstance.FPDF_GetFileVersion(&requests.FPDF_GetFileVersion{})
				Expect(err).To(MatchError("document not given"))
				Expect(pageCount).To(BeNil())
			})

			It("returns an error when getting the doc permissions", func() {
				pageCount, err := PdfiumInstance.FPDF_GetDocPermissions(&requests.FPDF_GetDocPermissions{})
				Expect(err).To(MatchError("document not given"))
				Expect(pageCount).To(BeNil())
			})

			It("returns an error when getting the doc revision number of security handler", func() {
				pageCount, err := PdfiumInstance.FPDF_GetSecurityHandlerRevision(&requests.FPDF_GetSecurityHandlerRevision{})
				Expect(err).To(MatchError("document not given"))
				Expect(pageCount).To(BeNil())
			})

			It("returns an error when getting the page count", func() {
				pageCount, err := PdfiumInstance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{})
				Expect(err).To(MatchError("document not given"))
				Expect(pageCount).To(BeNil())
			})

			It("returns an error when calling FPDF_GetPageSizeByIndex", func() {
				FPDF_GetPageSizeByIndex, err := PdfiumInstance.FPDF_GetPageSizeByIndex(&requests.FPDF_GetPageSizeByIndex{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_GetPageSizeByIndex).To(BeNil())
			})

			It("returns an error when calling FPDF_VIEWERREF_GetPrintScaling", func() {
				FPDF_VIEWERREF_GetPrintScaling, err := PdfiumInstance.FPDF_VIEWERREF_GetPrintScaling(&requests.FPDF_VIEWERREF_GetPrintScaling{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_VIEWERREF_GetPrintScaling).To(BeNil())
			})

			It("returns an error when calling FPDF_VIEWERREF_GetNumCopies", func() {
				FPDF_VIEWERREF_GetNumCopies, err := PdfiumInstance.FPDF_VIEWERREF_GetNumCopies(&requests.FPDF_VIEWERREF_GetNumCopies{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_VIEWERREF_GetNumCopies).To(BeNil())
			})

			It("returns an error when calling FPDF_VIEWERREF_GetPrintPageRange", func() {
				FPDF_VIEWERREF_GetPrintPageRange, err := PdfiumInstance.FPDF_VIEWERREF_GetPrintPageRange(&requests.FPDF_VIEWERREF_GetPrintPageRange{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_VIEWERREF_GetPrintPageRange).To(BeNil())
			})

			It("returns an error when calling FPDF_VIEWERREF_GetDuplex", func() {
				FPDF_VIEWERREF_GetDuplex, err := PdfiumInstance.FPDF_VIEWERREF_GetDuplex(&requests.FPDF_VIEWERREF_GetDuplex{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_VIEWERREF_GetDuplex).To(BeNil())
			})

			It("returns an error when calling FPDF_VIEWERREF_GetName", func() {
				FPDF_VIEWERREF_GetName, err := PdfiumInstance.FPDF_VIEWERREF_GetName(&requests.FPDF_VIEWERREF_GetName{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_VIEWERREF_GetName).To(BeNil())
			})

			It("returns an error when calling FPDF_CountNamedDests", func() {
				FPDF_CountNamedDests, err := PdfiumInstance.FPDF_CountNamedDests(&requests.FPDF_CountNamedDests{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_CountNamedDests).To(BeNil())
			})

			It("returns an error when calling FPDF_GetNamedDestByName", func() {
				FPDF_GetNamedDestByName, err := PdfiumInstance.FPDF_GetNamedDestByName(&requests.FPDF_GetNamedDestByName{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_GetNamedDestByName).To(BeNil())
			})

			It("returns an error when calling FPDF_GetNamedDest", func() {
				FPDF_GetNamedDest, err := PdfiumInstance.FPDF_GetNamedDest(&requests.FPDF_GetNamedDest{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_GetNamedDest).To(BeNil())
			})
		})
	})

	Context("no page", func() {
		When("is opened", func() {
			It("returns an error when calling FPDF_RenderPageBitmap", func() {
				FPDF_RenderPageBitmap, err := PdfiumInstance.FPDF_RenderPageBitmap(&requests.FPDF_RenderPageBitmap{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDF_RenderPageBitmap).To(BeNil())
			})

			It("returns an error when calling FPDF_RenderPageBitmapWithMatrix", func() {
				FPDF_RenderPageBitmapWithMatrix, err := PdfiumInstance.FPDF_RenderPageBitmapWithMatrix(&requests.FPDF_RenderPageBitmapWithMatrix{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDF_RenderPageBitmapWithMatrix).To(BeNil())
			})

			It("returns an error when calling FPDF_DeviceToPage", func() {
				FPDF_DeviceToPage, err := PdfiumInstance.FPDF_DeviceToPage(&requests.FPDF_DeviceToPage{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDF_DeviceToPage).To(BeNil())
			})

			It("returns an error when calling FPDF_PageToDevice", func() {
				FPDF_PageToDevice, err := PdfiumInstance.FPDF_PageToDevice(&requests.FPDF_PageToDevice{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDF_PageToDevice).To(BeNil())
			})
		})
	})

	Context("no bitmap", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFBitmap_GetFormat", func() {
				FPDFBitmap_GetFormat, err := PdfiumInstance.FPDFBitmap_GetFormat(&requests.FPDFBitmap_GetFormat{})
				Expect(err).To(MatchError("bitmap not given"))
				Expect(FPDFBitmap_GetFormat).To(BeNil())
			})

			It("returns an error when calling FPDFBitmap_FillRect", func() {
				FPDFBitmap_FillRect, err := PdfiumInstance.FPDFBitmap_FillRect(&requests.FPDFBitmap_FillRect{})
				Expect(err).To(MatchError("bitmap not given"))
				Expect(FPDFBitmap_FillRect).To(BeNil())
			})

			It("returns an error when calling FPDFBitmap_GetBuffer", func() {
				FPDFBitmap_GetBuffer, err := PdfiumInstance.FPDFBitmap_GetBuffer(&requests.FPDFBitmap_GetBuffer{})
				Expect(err).To(MatchError("bitmap not given"))
				Expect(FPDFBitmap_GetBuffer).To(BeNil())
			})

			It("returns an error when calling FPDFBitmap_GetWidth", func() {
				FPDFBitmap_GetWidth, err := PdfiumInstance.FPDFBitmap_GetWidth(&requests.FPDFBitmap_GetWidth{})
				Expect(err).To(MatchError("bitmap not given"))
				Expect(FPDFBitmap_GetWidth).To(BeNil())
			})

			It("returns an error when calling FPDFBitmap_GetFormat", func() {
				FPDFBitmap_GetHeight, err := PdfiumInstance.FPDFBitmap_GetHeight(&requests.FPDFBitmap_GetHeight{})
				Expect(err).To(MatchError("bitmap not given"))
				Expect(FPDFBitmap_GetHeight).To(BeNil())
			})

			It("returns an error when calling FPDFBitmap_GetStride", func() {
				FPDFBitmap_GetStride, err := PdfiumInstance.FPDFBitmap_GetStride(&requests.FPDFBitmap_GetStride{})
				Expect(err).To(MatchError("bitmap not given"))
				Expect(FPDFBitmap_GetStride).To(BeNil())
			})

			It("returns an error when calling FPDFBitmap_Destroy", func() {
				FPDFBitmap_Destroy, err := PdfiumInstance.FPDFBitmap_Destroy(&requests.FPDFBitmap_Destroy{})
				Expect(err).To(MatchError("bitmap not given"))
				Expect(FPDFBitmap_Destroy).To(BeNil())
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
			It("returns the correct page count", func() {
				FPDF_GetPageCount, err := PdfiumInstance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetPageCount).To(Equal(&responses.FPDF_GetPageCount{
					PageCount: 1,
				}))
			})

			It("returns an error when FPDF_LoadPage is called without a document", func() {
				FPDF_LoadPage, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_LoadPage).To(BeNil())
			})

			It("returns an error when FPDF_LoadPage is called with an invalid page", func() {
				FPDF_LoadPage, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{
					Document: doc,
					Index:    23,
				})
				Expect(err).To(MatchError("incorrect page"))
				Expect(FPDF_LoadPage).To(BeNil())
			})

			It("returns an error when FPDF_ClosePage is called without a page", func() {
				FPDF_ClosePage, err := PdfiumInstance.FPDF_ClosePage(&requests.FPDF_ClosePage{})
				Expect(err).To(MatchError("page not given"))
				Expect(FPDF_ClosePage).To(BeNil())
			})

			It("allows for a page to be loaded and closed", func() {
				FPDF_LoadPage, err := PdfiumInstance.FPDF_LoadPage(&requests.FPDF_LoadPage{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_LoadPage).To(Not(BeNil()))
				Expect(FPDF_LoadPage.Page).To(Not(BeNil()))

				FPDF_ClosePage, err := PdfiumInstance.FPDF_ClosePage(&requests.FPDF_ClosePage{
					Page: FPDF_LoadPage.Page,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_ClosePage).To(Not(BeNil()))
			})

			It("returns the correct page width", func() {
				pageCount, err := PdfiumInstance.FPDF_GetPageWidth(&requests.FPDF_GetPageWidth{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(pageCount).To(Equal(&responses.FPDF_GetPageWidth{
					Page:  0,
					Width: 595.2755737304688,
				}))
			})

			It("returns an error when calling FPDF_GetPageWidth without a page", func() {
				pageCount, err := PdfiumInstance.FPDF_GetPageWidth(&requests.FPDF_GetPageWidth{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(pageCount).To(BeNil())
			})

			It("returns the correct page height", func() {
				pageCount, err := PdfiumInstance.FPDF_GetPageHeight(&requests.FPDF_GetPageHeight{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(pageCount).To(Equal(&responses.FPDF_GetPageHeight{
					Page:   0,
					Height: 841.8897094726562,
				}))
			})

			It("returns an error when calling FPDF_GetPageHeight without a page", func() {
				pageCount, err := PdfiumInstance.FPDF_GetPageHeight(&requests.FPDF_GetPageHeight{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(pageCount).To(BeNil())
			})

			It("returns the correct page size by index", func() {
				pageCount, err := PdfiumInstance.FPDF_GetPageSizeByIndex(&requests.FPDF_GetPageSizeByIndex{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(pageCount).To(Equal(&responses.FPDF_GetPageSizeByIndex{
					Page:   0,
					Width:  595.2755737304688,
					Height: 841.8897094726562,
				}))
			})

			It("returns an error when calling FPDF_GetPageSizeByIndex with an invalid page", func() {
				pageCount, err := PdfiumInstance.FPDF_GetPageSizeByIndex(&requests.FPDF_GetPageSizeByIndex{
					Document: doc,
					Index:    3,
				})
				Expect(err).To(MatchError("could not load page size by index"))
				Expect(pageCount).To(BeNil())
			})

			It("returns the correct device to page calculations", func() {
				FPDF_DeviceToPage, err := PdfiumInstance.FPDF_DeviceToPage(&requests.FPDF_DeviceToPage{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					StartX:  0,
					StartY:  0,
					SizeX:   1000,
					SizeY:   1000,
					Rotate:  enums.FPDF_PAGE_ROTATION_NONE,
					DeviceX: 500,
					DeviceY: 500,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_DeviceToPage).To(Equal(&responses.FPDF_DeviceToPage{
					PageX: 297.6377868652344,
					PageY: 420.9448547363281,
				}))
			})

			It("returns the correct page to device calculations", func() {
				FPDF_PageToDevice, err := PdfiumInstance.FPDF_PageToDevice(&requests.FPDF_PageToDevice{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
					StartX: 0,
					StartY: 0,
					SizeX:  1000,
					SizeY:  1000,
					Rotate: enums.FPDF_PAGE_ROTATION_NONE,
					PageX:  297.6377868652344,
					PageY:  420.9448547363281,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_PageToDevice).To(Equal(&responses.FPDF_PageToDevice{
					DeviceX: 500,
					DeviceY: 500,
				}))
			})

			It("returns the correct print scaling", func() {
				FPDF_VIEWERREF_GetPrintScaling, err := PdfiumInstance.FPDF_VIEWERREF_GetPrintScaling(&requests.FPDF_VIEWERREF_GetPrintScaling{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_VIEWERREF_GetPrintScaling).To(Equal(&responses.FPDF_VIEWERREF_GetPrintScaling{
					PreferPrintScaling: true,
				}))
			})

			It("returns the correct print number of copies", func() {
				FPDF_VIEWERREF_GetNumCopies, err := PdfiumInstance.FPDF_VIEWERREF_GetNumCopies(&requests.FPDF_VIEWERREF_GetNumCopies{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_VIEWERREF_GetNumCopies).To(Equal(&responses.FPDF_VIEWERREF_GetNumCopies{
					NumCopies: 1,
				}))
			})

			It("returns the correct print duplex", func() {
				FPDF_VIEWERREF_GetDuplex, err := PdfiumInstance.FPDF_VIEWERREF_GetDuplex(&requests.FPDF_VIEWERREF_GetDuplex{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_VIEWERREF_GetDuplex).To(Equal(&responses.FPDF_VIEWERREF_GetDuplex{
					DuplexType: enums.FPDF_DUPLEXTYPE_UNDEFINED,
				}))
			})

			When("no bitmap is given", func() {
				It("returns an error when FPDF_RenderPageBitmap is called", func() {
					FPDF_RenderPageBitmap, err := PdfiumInstance.FPDF_RenderPageBitmap(&requests.FPDF_RenderPageBitmap{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})
					Expect(err).To(MatchError("bitmap not given"))
					Expect(FPDF_RenderPageBitmap).To(BeNil())
				})

				It("returns an error when FPDF_RenderPageBitmapWithMatrix is called", func() {
					FPDF_RenderPageBitmapWithMatrix, err := PdfiumInstance.FPDF_RenderPageBitmapWithMatrix(&requests.FPDF_RenderPageBitmapWithMatrix{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
							},
						},
					})
					Expect(err).To(MatchError("bitmap not given"))
					Expect(FPDF_RenderPageBitmapWithMatrix).To(BeNil())
				})
			})

			When("a bitmap has been created", func() {
				var bitmap references.FPDF_BITMAP
				BeforeEach(func() {
					FPDFBitmap_Create, err := PdfiumInstance.FPDFBitmap_Create(&requests.FPDFBitmap_Create{
						Width:  1000,
						Height: 1000,
						Alpha:  1,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_Create).To(Not(BeNil()))
					bitmap = FPDFBitmap_Create.Bitmap
				})

				AfterEach(func() {
					FPDFBitmap_Destroy, err := PdfiumInstance.FPDFBitmap_Destroy(&requests.FPDFBitmap_Destroy{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_Destroy).To(Equal(&responses.FPDFBitmap_Destroy{}))
				})

				It("returns the correct bitmap format", func() {
					FPDFBitmap_GetFormat, err := PdfiumInstance.FPDFBitmap_GetFormat(&requests.FPDFBitmap_GetFormat{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetFormat).To(Equal(&responses.FPDFBitmap_GetFormat{
						Format: enums.FPDF_BITMAP_FORMAT_BGRA,
					}))
				})

				It("returns the correct bitmap width", func() {
					FPDFBitmap_GetWidth, err := PdfiumInstance.FPDFBitmap_GetWidth(&requests.FPDFBitmap_GetWidth{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetWidth).To(Equal(&responses.FPDFBitmap_GetWidth{
						Width: 1000,
					}))
				})

				It("returns the correct bitmap height", func() {
					FPDFBitmap_GetHeight, err := PdfiumInstance.FPDFBitmap_GetHeight(&requests.FPDFBitmap_GetHeight{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetHeight).To(Equal(&responses.FPDFBitmap_GetHeight{
						Height: 1000,
					}))
				})

				It("returns the correct bitmap stride", func() {
					FPDFBitmap_GetStride, err := PdfiumInstance.FPDFBitmap_GetStride(&requests.FPDFBitmap_GetStride{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetStride).To(Equal(&responses.FPDFBitmap_GetStride{
						Stride: 4000,
					}))
				})

				It("returns the correct bitmap buffer", func() {
					FPDFBitmap_GetBuffer, err := PdfiumInstance.FPDFBitmap_GetBuffer(&requests.FPDFBitmap_GetBuffer{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetBuffer).To(Not(BeNil()))
					Expect(FPDFBitmap_GetBuffer.Buffer).To(Not(BeNil()))
					Expect(FPDFBitmap_GetBuffer.Buffer).To(HaveLen(4000000))
					Expect(FPDFBitmap_GetBuffer.Buffer[0]).To(Equal(uint8(0)))
				})

				It("allows the bitmap to be filled by white", func() {
					By("when the bitmap is filled")
					FPDFBitmap_FillRect, err := PdfiumInstance.FPDFBitmap_FillRect(&requests.FPDFBitmap_FillRect{
						Bitmap: bitmap,
						Color:  0xFFFFFFFF,
						Left:   0,
						Top:    0,
						Width:  1000,
						Height: 1000,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_FillRect).To(Equal(&responses.FPDFBitmap_FillRect{}))

					By("the first byte of the buffer is white")
					FPDFBitmap_GetBuffer, err := PdfiumInstance.FPDFBitmap_GetBuffer(&requests.FPDFBitmap_GetBuffer{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetBuffer).To(Not(BeNil()))
					Expect(FPDFBitmap_GetBuffer.Buffer).To(Not(BeNil()))
					Expect(FPDFBitmap_GetBuffer.Buffer[0]).To(Equal(uint8(255)))
				})

				It("allows a page to be rendered to the bitmap", func() {
					By("when the page is rendered")
					FPDF_RenderPageBitmap, err := PdfiumInstance.FPDF_RenderPageBitmap(&requests.FPDF_RenderPageBitmap{
						Bitmap: bitmap,
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
						StartX: 0,
						StartY: 0,
						SizeX:  1000,
						SizeY:  1000,
						Rotate: enums.FPDF_PAGE_ROTATION_NONE,
						Flags:  0,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_RenderPageBitmap).To(Equal(&responses.FPDF_RenderPageBitmap{}))

					By("the byte 212483 of the buffer is non-transparent")
					FPDFBitmap_GetBuffer, err := PdfiumInstance.FPDFBitmap_GetBuffer(&requests.FPDFBitmap_GetBuffer{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetBuffer).To(Not(BeNil()))
					Expect(FPDFBitmap_GetBuffer.Buffer).To(Not(BeNil()))
					Expect(FPDFBitmap_GetBuffer.Buffer[212483]).To(Equal(uint8(44)))
				})

				It("allows a page to be rendered to the bitmap by matrix and clipping", func() {
					By("when the page is rendered")
					FPDF_RenderPageBitmapWithMatrix, err := PdfiumInstance.FPDF_RenderPageBitmapWithMatrix(&requests.FPDF_RenderPageBitmapWithMatrix{
						Bitmap: bitmap,
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
						Matrix:   structs.FPDF_FS_MATRIX{0.5, 0, 0, 0.5, 0, 0}, // Half scale
						Clipping: structs.FPDF_FS_RECTF{0, 0, 1000, 1000},
						Flags:    0,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_RenderPageBitmapWithMatrix).To(Equal(&responses.FPDF_RenderPageBitmapWithMatrix{}))

					By("the byte 84143 of the buffer is non-transparent")
					FPDFBitmap_GetBuffer, err := PdfiumInstance.FPDFBitmap_GetBuffer(&requests.FPDFBitmap_GetBuffer{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetBuffer).To(Not(BeNil()))
					Expect(FPDFBitmap_GetBuffer.Buffer).To(Not(BeNil()))
					Expect(FPDFBitmap_GetBuffer.Buffer[84143]).To(Equal(uint8(2)))
				})
			})

			When("an external bitmap has been created with a buffer reference", func() {
				var bitmap references.FPDF_BITMAP
				var buffer []byte
				BeforeEach(func() {
					if TestType == "multi" {
						Skip("External bitmap is not supported on multi-threaded usage")
					}

					if TestType == "webassembly" {
						Skip("External bitmap by buffer reference is not supported on webassembly usage")
					}

					buffer = make([]byte, (1000*4)*1500) // 1000 pixels in width * 4 bytes per pixel * 1000 pixels in height

					FPDFBitmap_CreateEx, err := PdfiumInstance.FPDFBitmap_CreateEx(&requests.FPDFBitmap_CreateEx{
						Width:  1000,
						Height: 1500,
						Format: enums.FPDF_BITMAP_FORMAT_BGRA,
						Buffer: buffer,
						Stride: 4000,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_CreateEx).To(Not(BeNil()))
					bitmap = FPDFBitmap_CreateEx.Bitmap
				})

				AfterEach(func() {
					if TestType == "multi" {
						Skip("External bitmap is not supported on multi-threaded usage")
					}

					if TestType == "webassembly" {
						Skip("External bitmap by buffer reference is not supported on webassembly usage")
					}

					FPDFBitmap_Destroy, err := PdfiumInstance.FPDFBitmap_Destroy(&requests.FPDFBitmap_Destroy{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_Destroy).To(Equal(&responses.FPDFBitmap_Destroy{}))
					buffer = nil
				})

				It("returns the correct bitmap format", func() {
					FPDFBitmap_GetFormat, err := PdfiumInstance.FPDFBitmap_GetFormat(&requests.FPDFBitmap_GetFormat{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetFormat).To(Equal(&responses.FPDFBitmap_GetFormat{
						Format: enums.FPDF_BITMAP_FORMAT_BGRA,
					}))
				})

				It("returns the correct bitmap width", func() {
					FPDFBitmap_GetWidth, err := PdfiumInstance.FPDFBitmap_GetWidth(&requests.FPDFBitmap_GetWidth{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetWidth).To(Equal(&responses.FPDFBitmap_GetWidth{
						Width: 1000,
					}))
				})

				It("returns the correct bitmap height", func() {
					FPDFBitmap_GetHeight, err := PdfiumInstance.FPDFBitmap_GetHeight(&requests.FPDFBitmap_GetHeight{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetHeight).To(Equal(&responses.FPDFBitmap_GetHeight{
						Height: 1500,
					}))
				})

				It("returns the correct bitmap stride", func() {
					FPDFBitmap_GetStride, err := PdfiumInstance.FPDFBitmap_GetStride(&requests.FPDFBitmap_GetStride{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetStride).To(Equal(&responses.FPDFBitmap_GetStride{
						Stride: 4000,
					}))
				})

				It("allows the bitmap to be filled by white", func() {
					By("the first byte of the buffer is transparent")
					Expect(buffer[0]).To(Equal(uint8(0)))

					By("when the bitmap is filled")
					FPDFBitmap_FillRect, err := PdfiumInstance.FPDFBitmap_FillRect(&requests.FPDFBitmap_FillRect{
						Bitmap: bitmap,
						Color:  0xFFFFFFFF,
						Left:   0,
						Top:    0,
						Width:  1000,
						Height: 1000,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_FillRect).To(Equal(&responses.FPDFBitmap_FillRect{}))

					By("the first byte of the buffer is white")
					Expect(buffer[0]).To(Equal(uint8(255)))
				})
			})

			When("an external bitmap has been created with a pointer reference", func() {
				var bitmap references.FPDF_BITMAP
				var buffer []byte
				var pointer interface{}
				width := 1000
				height := 1500
				stride := width * 4

				BeforeEach(func() {
					if TestType == "multi" {
						Skip("External bitmap is not supported on multi-threaded usage")
					}

					// Size = 1000 pixels in width * 4 bytes per pixel * 1000 pixels in height
					fileSize := stride * height

					if TestType == "single" || TestType == "internal" {
						buffer = make([]byte, fileSize)
						pointer = unsafe.Pointer(&buffer[0])
					} else if TestType == "webassembly" {
						webassemblyImplementation := PdfiumInstance.GetImplementation().(*implementation_webassembly.PdfiumImplementation)

						// Request memory
						memoryPointer, err := webassemblyImplementation.Malloc(uint64(fileSize))
						if err != nil {
							Expect(err).To(BeNil())
							return
						}

						// Create a view of the underlying memory.
						memoryBuffer, ok := webassemblyImplementation.Module.Memory().Read(uint32(memoryPointer), uint32(fileSize))
						Expect(ok).To(BeTrue())
						buffer = memoryBuffer
						pointer = memoryPointer
					}

					FPDFBitmap_CreateEx, err := PdfiumInstance.FPDFBitmap_CreateEx(&requests.FPDFBitmap_CreateEx{
						Width:   width,
						Height:  height,
						Format:  enums.FPDF_BITMAP_FORMAT_BGRA,
						Pointer: pointer,
						Stride:  stride,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_CreateEx).To(Not(BeNil()))
					bitmap = FPDFBitmap_CreateEx.Bitmap
				})

				AfterEach(func() {
					if TestType == "multi" {
						Skip("External bitmap is not supported on multi-threaded usage")
					}

					FPDFBitmap_Destroy, err := PdfiumInstance.FPDFBitmap_Destroy(&requests.FPDFBitmap_Destroy{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_Destroy).To(Equal(&responses.FPDFBitmap_Destroy{}))
					buffer = nil

					if TestType == "webassembly" {
						webassemblyImplementation := PdfiumInstance.GetImplementation().(*implementation_webassembly.PdfiumImplementation)

						// Free memory
						err := webassemblyImplementation.Free(pointer.(uint64))
						if err != nil {
							Expect(err).To(BeNil())
						}
					}
				})

				It("returns the correct bitmap format", func() {
					FPDFBitmap_GetFormat, err := PdfiumInstance.FPDFBitmap_GetFormat(&requests.FPDFBitmap_GetFormat{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetFormat).To(Equal(&responses.FPDFBitmap_GetFormat{
						Format: enums.FPDF_BITMAP_FORMAT_BGRA,
					}))
				})

				It("returns the correct bitmap width", func() {
					FPDFBitmap_GetWidth, err := PdfiumInstance.FPDFBitmap_GetWidth(&requests.FPDFBitmap_GetWidth{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetWidth).To(Equal(&responses.FPDFBitmap_GetWidth{
						Width: width,
					}))
				})

				It("returns the correct bitmap height", func() {
					FPDFBitmap_GetHeight, err := PdfiumInstance.FPDFBitmap_GetHeight(&requests.FPDFBitmap_GetHeight{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetHeight).To(Equal(&responses.FPDFBitmap_GetHeight{
						Height: height,
					}))
				})

				It("returns the correct bitmap stride", func() {
					FPDFBitmap_GetStride, err := PdfiumInstance.FPDFBitmap_GetStride(&requests.FPDFBitmap_GetStride{
						Bitmap: bitmap,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_GetStride).To(Equal(&responses.FPDFBitmap_GetStride{
						Stride: stride,
					}))
				})

				It("allows the bitmap to be filled by white", func() {
					By("the first byte of the buffer is transparent")
					Expect(buffer[0]).To(Equal(uint8(0)))

					By("when the bitmap is filled")
					FPDFBitmap_FillRect, err := PdfiumInstance.FPDFBitmap_FillRect(&requests.FPDFBitmap_FillRect{
						Bitmap: bitmap,
						Color:  0xFFFFFFFF,
						Left:   0,
						Top:    0,
						Width:  width,
						Height: height,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_FillRect).To(Equal(&responses.FPDFBitmap_FillRect{}))

					By("the first byte of the buffer is white")
					Expect(buffer[0]).To(Equal(uint8(255)))
				})
			})
		})
	})

	Context("a PDF file with named dests", func() {
		var doc references.FPDF_DOCUMENT
		var file *os.File

		BeforeEach(func() {
			pdfFile, err := os.Open(TestDataPath + "/testdata/named_dests.pdf")
			Expect(err).To(BeNil())

			file = pdfFile
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
			file.Close()
		})

		When("is opened", func() {
			It("returns the correct dest count", func() {
				FPDF_CountNamedDests, err := PdfiumInstance.FPDF_CountNamedDests(&requests.FPDF_CountNamedDests{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_CountNamedDests).To(Equal(&responses.FPDF_CountNamedDests{
					Count: 6,
				}))
			})

			It("returns a dest by name", func() {
				FPDF_GetNamedDestByName, err := PdfiumInstance.FPDF_GetNamedDestByName(&requests.FPDF_GetNamedDestByName{
					Document: doc,
					Name:     "First",
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetNamedDestByName).To(Not(BeNil()))
				Expect(FPDF_GetNamedDestByName.Dest).To(Not(BeNil()))
			})

			It("returns an error when it can't find a dest by name", func() {
				FPDF_GetNamedDestByName, err := PdfiumInstance.FPDF_GetNamedDestByName(&requests.FPDF_GetNamedDestByName{
					Document: doc,
					Name:     "Firstfake",
				})
				Expect(err).To(MatchError("could not get named dest by name"))
				Expect(FPDF_GetNamedDestByName).To(BeNil())
			})

			It("returns the correct dest by index", func() {
				FPDF_GetNamedDest, err := PdfiumInstance.FPDF_GetNamedDest(&requests.FPDF_GetNamedDest{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetNamedDest).To(Not(BeNil()))
				Expect(FPDF_GetNamedDest.Name).To(Equal("First"))
			})

			It("returns an error when getting a dest with an invalid index", func() {
				FPDF_GetNamedDest, err := PdfiumInstance.FPDF_GetNamedDest(&requests.FPDF_GetNamedDest{
					Document: doc,
					Index:    25,
				})
				Expect(err).To(MatchError("could not get name of named dest"))
				Expect(FPDF_GetNamedDest).To(BeNil())
			})
		})
	})

	Context("a normal PDF file with multiple pages", func() {
		var doc references.FPDF_DOCUMENT
		var file *os.File

		BeforeEach(func() {
			pdfFile, err := os.Open(TestDataPath + "/testdata/viewer_ref.pdf")
			Expect(err).To(BeNil())

			file = pdfFile
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
			file.Close()
		})

		When("is opened", func() {
			It("returns the correct print duplex", func() {
				FPDF_VIEWERREF_GetDuplex, err := PdfiumInstance.FPDF_VIEWERREF_GetDuplex(&requests.FPDF_VIEWERREF_GetDuplex{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_VIEWERREF_GetDuplex).To(Equal(&responses.FPDF_VIEWERREF_GetDuplex{
					DuplexType: enums.FPDF_DUPLEXTYPE_UNDEFINED,
				}))
			})

			It("returns an error when requesting an invalid viewer ref name", func() {
				FPDF_VIEWERREF_GetName, err := PdfiumInstance.FPDF_VIEWERREF_GetName(&requests.FPDF_VIEWERREF_GetName{
					Document: doc,
					Key:      "foo",
				})
				Expect(err).To(MatchError("could not get name"))
				Expect(FPDF_VIEWERREF_GetName).To(BeNil())
			})

			It("returns the viewer ref name", func() {
				FPDF_VIEWERREF_GetName, err := PdfiumInstance.FPDF_VIEWERREF_GetName(&requests.FPDF_VIEWERREF_GetName{
					Document: doc,
					Key:      "Direction",
				})
				Expect(err).To(BeNil())
				Expect(FPDF_VIEWERREF_GetName).To(Equal(&responses.FPDF_VIEWERREF_GetName{
					Value: "R2L",
				}))
			})
		})
	})

	Context("a normal PDF file with multiple pages", func() {
		var doc references.FPDF_DOCUMENT
		var file *os.File

		BeforeEach(func() {
			pdfFile, err := os.Open(TestDataPath + "/testdata/test_multipage.pdf")
			Expect(err).To(BeNil())

			file = pdfFile
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
			file.Close()
		})

		When("is opened", func() {
			It("returns the correct file version", func() {
				pageCount, err := PdfiumInstance.FPDF_GetFileVersion(&requests.FPDF_GetFileVersion{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(pageCount).To(Equal(&responses.FPDF_GetFileVersion{
					FileVersion: 15,
				}))
			})

			It("returns the correct page count", func() {
				pageCount, err := PdfiumInstance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(pageCount).To(Equal(&responses.FPDF_GetPageCount{
					PageCount: 2,
				}))
			})
		})
	})

	Context("a normal PDF file with alpha channel", func() {
		var doc references.FPDF_DOCUMENT
		var file *os.File

		BeforeEach(func() {
			pdfFile, err := os.Open(TestDataPath + "/testdata/alpha_channel.pdf")
			Expect(err).To(BeNil())

			file = pdfFile
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
			file.Close()
		})

		When("is opened", func() {
			It("returns the correct file version", func() {
				pageCount, err := PdfiumInstance.FPDF_GetFileVersion(&requests.FPDF_GetFileVersion{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(pageCount).To(Equal(&responses.FPDF_GetFileVersion{
					FileVersion: 17,
				}))
			})

			It("returns the correct page count", func() {
				pageCount, err := PdfiumInstance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(pageCount).To(Equal(&responses.FPDF_GetPageCount{
					PageCount: 1,
				}))
			})
		})
	})
})
