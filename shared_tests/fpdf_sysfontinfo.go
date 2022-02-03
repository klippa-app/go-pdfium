package shared_tests

import (
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func RunfpdfSysfontinfoTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("fpdf_sysfontinfo", func() {
		It("returns an error when getting the decoded thumbnail data", func() {
			FPDF_GetDefaultTTFMap, err := pdfiumContainer.FPDF_GetDefaultTTFMap(&requests.FPDF_GetDefaultTTFMap{})
			Expect(err).To(BeNil())
			Expect(FPDF_GetDefaultTTFMap).To(Not(BeNil()))
			Expect(FPDF_GetDefaultTTFMap.TTFMap).To(Not(BeNil()))

			// Check length, actually checking the length is hard since it's
			// platform/runtime specific.
			Expect(FPDF_GetDefaultTTFMap.TTFMap).To(Not(HaveLen(0)))
		})
	})
}
