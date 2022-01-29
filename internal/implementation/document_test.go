package implementation_test

import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Document", func() {
	pdfium := implementation.Pdfium{}

	Context("no document", func() {
		When("is opened", func() {
			It("returns an error when getting the page count", func() {
				pageCount, err := pdfium.GetPageCount(&requests.GetPageCount{})
				Expect(err).To(MatchError("no current document"))
				Expect(pageCount).To(BeNil())
			})
		})
	})

	Context("a normal PDF file with 1 page", func() {
		BeforeEach(func() {
			pdfData, _ := ioutil.ReadFile("./testdata/test.pdf")
			pdfium.OpenDocument(&requests.OpenDocument{
				File: &pdfData,
			})
		})

		AfterEach(func() {
			pdfium.Close()
		})

		When("is opened", func() {
			It("returns the correct page count", func() {
				pageCount, err := pdfium.GetPageCount(&requests.GetPageCount{})
				Expect(err).To(BeNil())
				Expect(pageCount).To(Equal(&responses.GetPageCount{
					PageCount: 1,
				}))
			})
		})
	})

	Context("a normal PDF file with multiple pages", func() {
		BeforeEach(func() {
			pdfData, _ := ioutil.ReadFile("./testdata/test_multipage.pdf")
			pdfium.OpenDocument(&requests.OpenDocument{
				File: &pdfData,
			})
		})

		AfterEach(func() {
			pdfium.Close()
		})

		When("is opened", func() {
			It("returns the correct page count", func() {
				pageCount, err := pdfium.GetPageCount(&requests.GetPageCount{})
				Expect(err).To(BeNil())
				Expect(pageCount).To(Equal(&responses.GetPageCount{
					PageCount: 2,
				}))
			})
		})
	})
})
