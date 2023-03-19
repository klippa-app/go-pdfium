//go:build pdfium_experimental
// +build pdfium_experimental

package shared_tests

import "C"
import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_transformpage_experimental", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	Context("no page object", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFPageObj_GetClipPath", func() {
				FPDFPageObj_GetClipPath, err := PdfiumInstance.FPDFPageObj_GetClipPath(&requests.FPDFPageObj_GetClipPath{})
				Expect(err).To(MatchError("pageObject not given"))
				Expect(FPDFPageObj_GetClipPath).To(BeNil())
			})
		})
	})

	Context("no clippath", func() {
		When("is opened", func() {
			It("returns an error when calling FPDFClipPath_CountPaths", func() {
				FPDFClipPath_CountPaths, err := PdfiumInstance.FPDFClipPath_CountPaths(&requests.FPDFClipPath_CountPaths{})
				Expect(err).To(MatchError("clipPath not given"))
				Expect(FPDFClipPath_CountPaths).To(BeNil())
			})

			It("returns an error when calling FPDFClipPath_CountPathSegments", func() {
				FPDFClipPath_CountPathSegments, err := PdfiumInstance.FPDFClipPath_CountPathSegments(&requests.FPDFClipPath_CountPathSegments{})
				Expect(err).To(MatchError("clipPath not given"))
				Expect(FPDFClipPath_CountPathSegments).To(BeNil())
			})

			It("returns an error when calling FPDFClipPath_GetPathSegment", func() {
				FPDFClipPath_GetPathSegment, err := PdfiumInstance.FPDFClipPath_GetPathSegment(&requests.FPDFClipPath_GetPathSegment{})
				Expect(err).To(MatchError("clipPath not given"))
				Expect(FPDFClipPath_GetPathSegment).To(BeNil())
			})
		})
	})

	Context("a pdf file with clip paths", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/clip_path.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
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
			When("a page object is opened", func() {
				var pageObject references.FPDF_PAGEOBJECT
				BeforeEach(func() {
					FPDFPage_GetObject, err := PdfiumInstance.FPDFPage_GetObject(&requests.FPDFPage_GetObject{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
						Index: 0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPage_GetObject).To(Not(BeNil()))
					Expect(FPDFPage_GetObject.PageObject).To(Not(BeEmpty()))
					pageObject = FPDFPage_GetObject.PageObject
				})

				It("allows to get the clip path", func() {
					FPDFPageObj_GetClipPath, err := PdfiumInstance.FPDFPageObj_GetClipPath(&requests.FPDFPageObj_GetClipPath{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetClipPath).ToNot(BeNil())
					Expect(FPDFPageObj_GetClipPath.ClipPath).ToNot(BeEmpty())
				})

				It("returns the correct clip path info", func() {
					By("loading the clip path")
					FPDFPageObj_GetClipPath, err := PdfiumInstance.FPDFPageObj_GetClipPath(&requests.FPDFPageObj_GetClipPath{
						PageObject: pageObject,
					})
					Expect(err).To(BeNil())
					Expect(FPDFPageObj_GetClipPath).ToNot(BeNil())
					Expect(FPDFPageObj_GetClipPath.ClipPath).ToNot(BeEmpty())

					By("counting the paths in the clip path")
					FPDFClipPath_CountPaths, err := PdfiumInstance.FPDFClipPath_CountPaths(&requests.FPDFClipPath_CountPaths{
						ClipPath: FPDFPageObj_GetClipPath.ClipPath,
					})
					Expect(err).To(BeNil())
					Expect(FPDFClipPath_CountPaths).To(Equal(&responses.FPDFClipPath_CountPaths{
						Count: 1,
					}))

					By("receiving an error when giving an invalid index")
					FPDFClipPath_CountPathSegments, err := PdfiumInstance.FPDFClipPath_CountPathSegments(&requests.FPDFClipPath_CountPathSegments{
						ClipPath:  FPDFPageObj_GetClipPath.ClipPath,
						PathIndex: 25,
					})
					Expect(err).To(MatchError("could not get clip path path segment count"))
					Expect(FPDFClipPath_CountPathSegments).To(BeNil())

					By("counting the path segments in the clip path")
					FPDFClipPath_CountPathSegments, err = PdfiumInstance.FPDFClipPath_CountPathSegments(&requests.FPDFClipPath_CountPathSegments{
						ClipPath:  FPDFPageObj_GetClipPath.ClipPath,
						PathIndex: 0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFClipPath_CountPathSegments).To(Equal(&responses.FPDFClipPath_CountPathSegments{
						Count: 4,
					}))

					By("getting an error when requesting an invalid path segment in the clip path")
					FPDFClipPath_GetPathSegment, err := PdfiumInstance.FPDFClipPath_GetPathSegment(&requests.FPDFClipPath_GetPathSegment{
						ClipPath:     FPDFPageObj_GetClipPath.ClipPath,
						PathIndex:    25,
						SegmentIndex: 25,
					})
					Expect(err).To(MatchError("could not get clip path segment"))
					Expect(FPDFClipPath_GetPathSegment).To(BeNil())

					By("getting a path segment in the clip path")
					FPDFClipPath_GetPathSegment, err = PdfiumInstance.FPDFClipPath_GetPathSegment(&requests.FPDFClipPath_GetPathSegment{
						ClipPath:     FPDFPageObj_GetClipPath.ClipPath,
						PathIndex:    0,
						SegmentIndex: 0,
					})
					Expect(err).To(BeNil())
					Expect(FPDFClipPath_GetPathSegment).ToNot(BeNil())
					Expect(FPDFClipPath_GetPathSegment.PathSegment).ToNot(BeEmpty())
				})
			})
		})
	})
})
