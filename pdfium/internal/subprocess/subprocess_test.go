package subprocess_test

import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/pdfium/internal/subprocess"
	"github.com/klippa-app/go-pdfium/pdfium/pdfium_errors"
	"github.com/klippa-app/go-pdfium/pdfium/requests"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Subprocess", func() {
	pdfium := subprocess.Pdfium{}
	AfterEach(func() {
		pdfium.Close()
	})

	When("pinged", func() {
		It("pongs", func() {
			resp, err := pdfium.Ping()
			Expect(err).To(BeNil())
			Expect(resp).To(Equal("Pong"))
		})
	})

	Context("a normal PDF file", func() {
		pdfData, _ := ioutil.ReadFile("./testdata/test.pdf")
		When("is opened", func() {
			It("does not return an error", func() {
				err := pdfium.OpenDocument(&requests.OpenDocument{
					File: &pdfData,
				})
				Expect(err).To(BeNil())
			})
		})
	})

	Context("a password protected PDF file", func() {
		pdfData, _ := ioutil.ReadFile("./testdata/password_test123.pdf")
		When("is opened with no password", func() {
			It("returns the password error", func() {
				err := pdfium.OpenDocument(&requests.OpenDocument{
					File: &pdfData,
				})
				Expect(err).To(Equal(pdfium_errors.ErrPassword))
			})
		})
		When("is opened with the wrong password", func() {
			It("returns the password error", func() {
				wrongPassword := "test"
				err := pdfium.OpenDocument(&requests.OpenDocument{
					File:     &pdfData,
					Password: &wrongPassword,
				})
				Expect(err).To(Equal(pdfium_errors.ErrPassword))
			})
		})
		When("is opened with the correct password", func() {
			It("does not return an error", func() {
				pdfPassword := "test123"
				err := pdfium.OpenDocument(&requests.OpenDocument{
					File:     &pdfData,
					Password: &pdfPassword,
				})
				Expect(err).To(BeNil())
			})
		})
	})
})
