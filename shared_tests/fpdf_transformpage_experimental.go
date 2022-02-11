//go:build pdfium_experimental
// +build pdfium_experimental

package shared_tests

import "C"
import (
	"github.com/klippa-app/go-pdfium/requests"
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

	// @todo: add extra tests when FPDFPage_GetObject has been implemented.
})
