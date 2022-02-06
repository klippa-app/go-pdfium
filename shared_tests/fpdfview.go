package shared_tests

import "C"
import (
	"io/ioutil"
	"os"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdfview", func() {
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

			It("returns an error when calling FPDF_DocumentHasValidCrossReferenceTable", func() {
				FPDF_DocumentHasValidCrossReferenceTable, err := PdfiumInstance.FPDF_DocumentHasValidCrossReferenceTable(&requests.FPDF_DocumentHasValidCrossReferenceTable{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_DocumentHasValidCrossReferenceTable).To(BeNil())
			})

			It("returns an error when calling FPDF_GetTrailerEnds", func() {
				FPDF_GetTrailerEnds, err := PdfiumInstance.FPDF_GetTrailerEnds(&requests.FPDF_GetTrailerEnds{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_GetTrailerEnds).To(BeNil())
			})

			It("returns an error when calling FPDF_GetPageSizeByIndexF", func() {
				FPDF_GetPageSizeByIndexF, err := PdfiumInstance.FPDF_GetPageSizeByIndexF(&requests.FPDF_GetPageSizeByIndexF{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_GetPageSizeByIndexF).To(BeNil())
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

			It("returns an error when calling FPDF_GetXFAPacketCount", func() {
				FPDF_GetXFAPacketCount, err := PdfiumInstance.FPDF_GetXFAPacketCount(&requests.FPDF_GetXFAPacketCount{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_GetXFAPacketCount).To(BeNil())
			})

			It("returns an error when calling FPDF_GetXFAPacketName", func() {
				FPDF_GetXFAPacketName, err := PdfiumInstance.FPDF_GetXFAPacketName(&requests.FPDF_GetXFAPacketName{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_GetXFAPacketName).To(BeNil())
			})

			It("returns an error when calling FPDF_GetXFAPacketContent", func() {
				FPDF_GetXFAPacketContent, err := PdfiumInstance.FPDF_GetXFAPacketContent(&requests.FPDF_GetXFAPacketContent{})
				Expect(err).To(MatchError("document not given"))
				Expect(FPDF_GetXFAPacketContent).To(BeNil())
			})
		})
	})

	Context("no page", func() {
		When("is opened", func() {
			It("returns an error when calling FPDF_GetPageWidthF", func() {
				FPDF_GetPageWidthF, err := PdfiumInstance.FPDF_GetPageWidthF(&requests.FPDF_GetPageWidthF{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDF_GetPageWidthF).To(BeNil())
			})

			It("returns an error when calling FPDF_GetPageHeightF", func() {
				FPDF_GetPageHeightF, err := PdfiumInstance.FPDF_GetPageHeightF(&requests.FPDF_GetPageHeightF{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDF_GetPageHeightF).To(BeNil())
			})

			It("returns an error when calling FPDF_GetPageBoundingBox", func() {
				FPDF_GetPageBoundingBox, err := PdfiumInstance.FPDF_GetPageBoundingBox(&requests.FPDF_GetPageBoundingBox{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDF_GetPageBoundingBox).To(BeNil())
			})

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

	Context("no page range", func() {
		When("is opened", func() {
			It("returns an error when calling FPDF_VIEWERREF_GetPrintPageRangeCount", func() {
				FPDF_VIEWERREF_GetPrintPageRangeCount, err := PdfiumInstance.FPDF_VIEWERREF_GetPrintPageRangeCount(&requests.FPDF_VIEWERREF_GetPrintPageRangeCount{})
				Expect(err).To(MatchError("pageRange not given"))
				Expect(FPDF_VIEWERREF_GetPrintPageRangeCount).To(BeNil())
			})

			It("returns an error when calling FPDF_VIEWERREF_GetPrintPageRangeElement", func() {
				FPDF_VIEWERREF_GetPrintPageRangeElement, err := PdfiumInstance.FPDF_VIEWERREF_GetPrintPageRangeElement(&requests.FPDF_VIEWERREF_GetPrintPageRangeElement{})
				Expect(err).To(MatchError("pageRange not given"))
				Expect(FPDF_VIEWERREF_GetPrintPageRangeElement).To(BeNil())
			})
		})
	})

	Context("a PDF file with an invalid cross reference table", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			file, err := os.Open(TestDataPath + "/testdata/bug_664284.pdf")
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
			It("returns that it has an invalid cross reference table", func() {
				FPDF_DocumentHasValidCrossReferenceTable, err := PdfiumInstance.FPDF_DocumentHasValidCrossReferenceTable(&requests.FPDF_DocumentHasValidCrossReferenceTable{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_DocumentHasValidCrossReferenceTable).To(Equal(&responses.FPDF_DocumentHasValidCrossReferenceTable{
					DocumentHasValidCrossReferenceTable: false,
				}))
			})
		})
	})

	Context("a PDF file with multiple trailer ends", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			file, err := os.Open(TestDataPath + "/testdata/annotation_stamp_with_ap.pdf")
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
			It("returns the correct trailer ends", func() {
				FPDF_GetTrailerEnds, err := PdfiumInstance.FPDF_GetTrailerEnds(&requests.FPDF_GetTrailerEnds{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetTrailerEnds).To(Equal(&responses.FPDF_GetTrailerEnds{
					TrailerEnds: []int{441, 7945, 101719},
				}))
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

			It("returns that it has a valid cross reference table", func() {
				FPDF_DocumentHasValidCrossReferenceTable, err := PdfiumInstance.FPDF_DocumentHasValidCrossReferenceTable(&requests.FPDF_DocumentHasValidCrossReferenceTable{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_DocumentHasValidCrossReferenceTable).To(Equal(&responses.FPDF_DocumentHasValidCrossReferenceTable{
					DocumentHasValidCrossReferenceTable: true,
				}))
			})

			It("returns the correct trailer ends", func() {
				FPDF_GetTrailerEnds, err := PdfiumInstance.FPDF_GetTrailerEnds(&requests.FPDF_GetTrailerEnds{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetTrailerEnds).To(Equal(&responses.FPDF_GetTrailerEnds{
					TrailerEnds: []int{11616},
				}))
			})

			It("returns the correct page width in float32", func() {
				FPDF_GetPageWidthF, err := PdfiumInstance.FPDF_GetPageWidthF(&requests.FPDF_GetPageWidthF{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetPageWidthF).To(Equal(&responses.FPDF_GetPageWidthF{
					PageWidth: 595.2755737304688,
				}))
			})

			It("returns the correct page height in float32", func() {
				FPDF_GetPageHeightF, err := PdfiumInstance.FPDF_GetPageHeightF(&requests.FPDF_GetPageHeightF{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetPageHeightF).To(Equal(&responses.FPDF_GetPageHeightF{
					PageHeight: 841.8897094726562,
				}))
			})

			It("returns the correct page bounding box", func() {
				FPDF_GetPageBoundingBox, err := PdfiumInstance.FPDF_GetPageBoundingBox(&requests.FPDF_GetPageBoundingBox{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetPageBoundingBox).To(Equal(&responses.FPDF_GetPageBoundingBox{
					Rect: structs.FPDF_FS_RECTF{
						Left:   0,
						Top:    841.8897094726562,
						Right:  595.2755737304688,
						Bottom: 0,
					},
				}))
			})

			It("returns the correct page size by index in float32", func() {
				FPDF_GetPageSizeByIndexF, err := PdfiumInstance.FPDF_GetPageSizeByIndexF(&requests.FPDF_GetPageSizeByIndexF{
					Document: doc,
					Index:    0,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetPageSizeByIndexF).To(Equal(&responses.FPDF_GetPageSizeByIndexF{
					Size: structs.FPDF_FS_SIZEF{
						Width:  595.2755737304688,
						Height: 841.8897094726562,
					},
				}))
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

			When("a print page range has been loaded", func() {
				var pageRange references.FPDF_PAGERANGE
				BeforeEach(func() {
					FPDF_VIEWERREF_GetPrintPageRange, err := PdfiumInstance.FPDF_VIEWERREF_GetPrintPageRange(&requests.FPDF_VIEWERREF_GetPrintPageRange{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_VIEWERREF_GetPrintPageRange).To(Not(BeNil()))
					pageRange = FPDF_VIEWERREF_GetPrintPageRange.PageRange
				})

				It("returns the correct print page range count", func() {
					FPDF_VIEWERREF_GetPrintPageRangeCount, err := PdfiumInstance.FPDF_VIEWERREF_GetPrintPageRangeCount(&requests.FPDF_VIEWERREF_GetPrintPageRangeCount{
						PageRange: pageRange,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_VIEWERREF_GetPrintPageRangeCount).To(Equal(&responses.FPDF_VIEWERREF_GetPrintPageRangeCount{}))
				})
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

			When("an external bitmap has been created", func() {
				if TestType == "multi" {
					Skip("External bitmap is not supported on multi-threaded usage")
				}
				var bitmap references.FPDF_BITMAP
				var buffer []byte
				BeforeEach(func() {
					buffer = make([]byte, (1000*4)*1000) // 1000 pixels in width * 4 bytes per pixel * 1000 pixels in height

					FPDFBitmap_CreateEx, err := PdfiumInstance.FPDFBitmap_CreateEx(&requests.FPDFBitmap_CreateEx{
						Width:  1000,
						Height: 1000,
						Format: enums.FPDF_BITMAP_FORMAT_BGRA,
						Buffer: buffer,
						Stride: 4000,
					})
					Expect(err).To(BeNil())
					Expect(FPDFBitmap_CreateEx).To(Not(BeNil()))
					bitmap = FPDFBitmap_CreateEx.Bitmap
				})

				AfterEach(func() {
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

	Context("a normal PDF file", func() {
		It("can be loaded with FPDF_LoadMemDocument", func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/test.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument64(&requests.FPDF_LoadMemDocument64{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: newDoc.Document,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})
	})

	Context("a PDF file with XFA data", func() {
		var doc references.FPDF_DOCUMENT
		var file *os.File

		BeforeEach(func() {
			pdfFile, err := os.Open(TestDataPath + "/testdata/simple_xfa.pdf")
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
			It("returns the correct XFA packet count", func() {
				FPDF_GetXFAPacketCount, err := PdfiumInstance.FPDF_GetXFAPacketCount(&requests.FPDF_GetXFAPacketCount{
					Document: doc,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetXFAPacketCount).To(Equal(&responses.FPDF_GetXFAPacketCount{
					Count: 5,
				}))
			})

			It("returns the correct XFA packet name", func() {
				FPDF_GetXFAPacketName, err := PdfiumInstance.FPDF_GetXFAPacketName(&requests.FPDF_GetXFAPacketName{
					Document: doc,
					Index:    1,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetXFAPacketName).To(Equal(&responses.FPDF_GetXFAPacketName{
					Index: 1,
					Name:  "config",
				}))
			})

			It("returns an error when requesting an incorrect XFA packet name", func() {
				FPDF_GetXFAPacketName, err := PdfiumInstance.FPDF_GetXFAPacketName(&requests.FPDF_GetXFAPacketName{
					Document: doc,
					Index:    25,
				})
				Expect(err).To(MatchError("could not get name of the XFA packet"))
				Expect(FPDF_GetXFAPacketName).To(BeNil())
			})

			It("returns the correct XFA packet content", func() {
				FPDF_GetXFAPacketContent, err := PdfiumInstance.FPDF_GetXFAPacketContent(&requests.FPDF_GetXFAPacketContent{
					Document: doc,
					Index:    1,
				})
				Expect(err).To(BeNil())
				Expect(FPDF_GetXFAPacketContent).To(Equal(&responses.FPDF_GetXFAPacketContent{
					Index:   1,
					Content: []byte{60, 99, 111, 110, 102, 105, 103, 32, 120, 109, 108, 110, 115, 61, 34, 104, 116, 116, 112, 58, 47, 47, 119, 119, 119, 46, 120, 102, 97, 46, 111, 114, 103, 47, 115, 99, 104, 101, 109, 97, 47, 120, 99, 105, 47, 51, 46, 48, 47, 34, 62, 10, 60, 97, 103, 101, 110, 116, 32, 110, 97, 109, 101, 61, 34, 100, 101, 115, 105, 103, 110, 101, 114, 34, 62, 10, 32, 32, 60, 100, 101, 115, 116, 105, 110, 97, 116, 105, 111, 110, 62, 112, 100, 102, 60, 47, 100, 101, 115, 116, 105, 110, 97, 116, 105, 111, 110, 62, 10, 32, 32, 60, 112, 100, 102, 62, 10, 32, 32, 32, 32, 60, 102, 111, 110, 116, 73, 110, 102, 111, 47, 62, 10, 32, 32, 60, 47, 112, 100, 102, 62, 10, 60, 47, 97, 103, 101, 110, 116, 62, 10, 60, 112, 114, 101, 115, 101, 110, 116, 62, 10, 32, 32, 60, 112, 100, 102, 62, 10, 32, 32, 32, 32, 60, 118, 101, 114, 115, 105, 111, 110, 62, 49, 46, 55, 60, 47, 118, 101, 114, 115, 105, 111, 110, 62, 10, 32, 32, 32, 32, 60, 97, 100, 111, 98, 101, 69, 120, 116, 101, 110, 115, 105, 111, 110, 76, 101, 118, 101, 108, 62, 56, 60, 47, 97, 100, 111, 98, 101, 69, 120, 116, 101, 110, 115, 105, 111, 110, 76, 101, 118, 101, 108, 62, 10, 32, 32, 32, 32, 60, 114, 101, 110, 100, 101, 114, 80, 111, 108, 105, 99, 121, 62, 99, 108, 105, 101, 110, 116, 60, 47, 114, 101, 110, 100, 101, 114, 80, 111, 108, 105, 99, 121, 62, 10, 32, 32, 32, 32, 60, 115, 99, 114, 105, 112, 116, 77, 111, 100, 101, 108, 62, 88, 70, 65, 60, 47, 115, 99, 114, 105, 112, 116, 77, 111, 100, 101, 108, 62, 10, 32, 32, 32, 32, 60, 105, 110, 116, 101, 114, 97, 99, 116, 105, 118, 101, 62, 49, 60, 47, 105, 110, 116, 101, 114, 97, 99, 116, 105, 118, 101, 62, 10, 32, 32, 60, 47, 112, 100, 102, 62, 10, 32, 32, 60, 120, 100, 112, 62, 10, 32, 32, 32, 32, 60, 112, 97, 99, 107, 101, 116, 115, 62, 42, 60, 47, 112, 97, 99, 107, 101, 116, 115, 62, 10, 32, 32, 60, 47, 120, 100, 112, 62, 10, 32, 32, 60, 100, 101, 115, 116, 105, 110, 97, 116, 105, 111, 110, 62, 112, 100, 102, 60, 47, 100, 101, 115, 116, 105, 110, 97, 116, 105, 111, 110, 62, 10, 32, 32, 60, 115, 99, 114, 105, 112, 116, 62, 10, 32, 32, 32, 32, 60, 114, 117, 110, 83, 99, 114, 105, 112, 116, 115, 62, 115, 101, 114, 118, 101, 114, 60, 47, 114, 117, 110, 83, 99, 114, 105, 112, 116, 115, 62, 10, 32, 32, 60, 47, 115, 99, 114, 105, 112, 116, 62, 10, 60, 47, 112, 114, 101, 115, 101, 110, 116, 62, 10, 60, 97, 99, 114, 111, 98, 97, 116, 62, 10, 32, 32, 60, 97, 99, 114, 111, 98, 97, 116, 55, 62, 10, 32, 32, 32, 32, 60, 100, 121, 110, 97, 109, 105, 99, 82, 101, 110, 100, 101, 114, 62, 114, 101, 113, 117, 105, 114, 101, 100, 60, 47, 100, 121, 110, 97, 109, 105, 99, 82, 101, 110, 100, 101, 114, 62, 10, 32, 32, 60, 47, 97, 99, 114, 111, 98, 97, 116, 55, 62, 10, 32, 32, 60, 118, 97, 108, 105, 100, 97, 116, 101, 62, 112, 114, 101, 83, 117, 98, 109, 105, 116, 60, 47, 118, 97, 108, 105, 100, 97, 116, 101, 62, 10, 60, 47, 97, 99, 114, 111, 98, 97, 116, 62, 10, 60, 47, 99, 111, 110, 102, 105, 103, 62, 10},
				}))
			})

			It("returns an error when requesting an incorrect XFA packet content", func() {
				FPDF_GetXFAPacketContent, err := PdfiumInstance.FPDF_GetXFAPacketContent(&requests.FPDF_GetXFAPacketContent{
					Document: doc,
					Index:    25,
				})
				Expect(err).To(MatchError("could not get content of the XFA packet"))
				Expect(FPDF_GetXFAPacketContent).To(BeNil())
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

			When("a print page range has been loaded", func() {
				var pageRange references.FPDF_PAGERANGE
				BeforeEach(func() {
					FPDF_VIEWERREF_GetPrintPageRange, err := PdfiumInstance.FPDF_VIEWERREF_GetPrintPageRange(&requests.FPDF_VIEWERREF_GetPrintPageRange{
						Document: doc,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_VIEWERREF_GetPrintPageRange).To(Not(BeNil()))
					pageRange = FPDF_VIEWERREF_GetPrintPageRange.PageRange
				})

				It("returns the correct print page range count", func() {
					FPDF_VIEWERREF_GetPrintPageRangeCount, err := PdfiumInstance.FPDF_VIEWERREF_GetPrintPageRangeCount(&requests.FPDF_VIEWERREF_GetPrintPageRangeCount{
						PageRange: pageRange,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_VIEWERREF_GetPrintPageRangeCount).To(Equal(&responses.FPDF_VIEWERREF_GetPrintPageRangeCount{
						Count: 4,
					}))
				})

				It("returns an error when requesting an invalid print page range element", func() {
					FPDF_VIEWERREF_GetPrintPageRangeElement, err := PdfiumInstance.FPDF_VIEWERREF_GetPrintPageRangeElement(&requests.FPDF_VIEWERREF_GetPrintPageRangeElement{
						PageRange: pageRange,
						Index:     25,
					})
					Expect(err).To(MatchError("could not load page range element"))
					Expect(FPDF_VIEWERREF_GetPrintPageRangeElement).To(BeNil())
				})

				It("returns the correct print page range element", func() {
					FPDF_VIEWERREF_GetPrintPageRangeElement, err := PdfiumInstance.FPDF_VIEWERREF_GetPrintPageRangeElement(&requests.FPDF_VIEWERREF_GetPrintPageRangeElement{
						PageRange: pageRange,
						Index:     1,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_VIEWERREF_GetPrintPageRangeElement).To(Equal(&responses.FPDF_VIEWERREF_GetPrintPageRangeElement{
						Value: 2,
					}))
				})
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
