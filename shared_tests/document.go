package shared_tests

import (
	"github.com/klippa-app/go-pdfium/errors"
	"io/ioutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

func RunDocumentTests(pdfiumContainer pdfium.Pdfium, testsPath string, prefix string) {
	Describe("Document", func() {
		Context("a normal PDF file with 1 page", func() {
			var doc pdfium.Document

			BeforeEach(func() {
				pdfData, err := ioutil.ReadFile(testsPath + "/testdata/test.pdf")
				Expect(err).To(BeNil())
				if err != nil {
					return
				}

				newDoc, err := pdfiumContainer.NewDocument(&pdfData)
				if err != nil {
					return
				}

				doc = newDoc
			})

			AfterEach(func() {
				doc.Close()
			})

			When("is opened", func() {
				It("returns the correct page count", func() {
					pageCount, err := doc.GetPageCount(&requests.GetPageCount{})
					Expect(err).To(BeNil())
					Expect(pageCount).To(Equal(&responses.GetPageCount{
						PageCount: 1,
					}))
				})
			})
		})

		Context("a normal PDF file with multiple pages", func() {
			var doc pdfium.Document

			BeforeEach(func() {
				pdfData, err := ioutil.ReadFile(testsPath + "/testdata/test_multipage.pdf")
				Expect(err).To(BeNil())
				if err != nil {
					return
				}

				newDoc, err := pdfiumContainer.NewDocument(&pdfData)
				if err != nil {
					return
				}

				doc = newDoc
			})

			AfterEach(func() {
				doc.Close()
			})

			When("is opened", func() {
				It("returns the correct page count", func() {
					pageCount, err := doc.GetPageCount(&requests.GetPageCount{})
					Expect(err).To(BeNil())
					Expect(pageCount).To(Equal(&responses.GetPageCount{
						PageCount: 2,
					}))
				})
			})
		})

		Context("a password protected PDF file", func() {
			pdfData, _ := ioutil.ReadFile(testsPath + "/testdata/password_test123.pdf")
			When("is opened with no password", func() {
				It("returns the password error", func() {
					doc, err := pdfiumContainer.NewDocument(&pdfData)
					Expect(err).To(MatchError(errors.ErrPassword.Error()))
					Expect(doc).To(BeNil())
				})
			})
			When("is opened with the wrong password", func() {
				It("returns the password error", func() {
					wrongPassword := "test"
					doc, err := pdfiumContainer.NewDocument(&pdfData, pdfium.OpenDocumentWithPasswordOption(wrongPassword))
					Expect(err).To(MatchError(errors.ErrPassword.Error()))
					Expect(doc).To(BeNil())
				})
			})
			When("is opened with the correct password", func() {
				It("does not return an error", func() {
					pdfPassword := "test123"
					doc, err := pdfiumContainer.NewDocument(&pdfData, pdfium.OpenDocumentWithPasswordOption(pdfPassword))
					Expect(err).To(BeNil())
					Expect(doc).To(Not(BeNil()))
					doc.Close()
				})
			})
		})
	})
}
