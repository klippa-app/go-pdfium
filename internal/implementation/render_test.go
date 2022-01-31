package implementation_test

import (
	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/requests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Render", func() {
	pdfium := implementation.Pdfium.GetInstance()

	Context("no document", func() {
		When("is opened", func() {
			Context("GetPageSize()", func() {
				It("returns an error", func() {
					pageSize, err := pdfium.GetPageSize(&requests.GetPageSize{
						Page: 0,
					})
					Expect(err).To(MatchError("Document.Ref not given"))
					Expect(pageSize).To(BeNil())
				})
			})

			Context("GetPageSizeInPixels()", func() {
				It("returns an error", func() {
					pageSize, err := pdfium.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
						Page: 0,
						DPI:  100,
					})
					Expect(err).To(MatchError("Document.Ref not given"))
					Expect(pageSize).To(BeNil())
				})
			})

			Context("RenderPageInDPI()", func() {
				It("returns an error", func() {
					renderedPage, err := pdfium.RenderPageInDPI(&requests.RenderPageInDPI{
						Page: 0,
						DPI:  300,
					})
					Expect(err).To(MatchError("Document.Ref not given"))
					Expect(renderedPage).To(BeNil())
				})
			})

			Context("RenderPageInDPI()", func() {
				It("returns an error", func() {
					renderedPage, err := pdfium.RenderPagesInDPI(&requests.RenderPagesInDPI{
						Pages: []requests.RenderPageInDPI{
							{
								Page: 0,
								DPI:  300,
							},
						},
						Padding: 50,
					})
					Expect(err).To(MatchError("Document.Ref not given"))
					Expect(renderedPage).To(BeNil())
				})
			})

			Context("RenderPageInPixels()", func() {
				It("returns an error", func() {
					renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
						Page:   0,
						Width:  2000,
						Height: 2000,
					})
					Expect(err).To(MatchError("Document.Ref not given"))
					Expect(renderedPage).To(BeNil())
				})
			})

			Context("RenderPagesInPixels()", func() {
				It("returns an error", func() {
					renderedPage, err := pdfium.RenderPagesInPixels(&requests.RenderPagesInPixels{
						Pages: []requests.RenderPageInPixels{
							{
								Page:   0,
								Width:  2000,
								Height: 2000,
							},
						},
						Padding: 50,
					})
					Expect(err).To(MatchError("Document.Ref not given"))
					Expect(renderedPage).To(BeNil())
				})
			})
		})
	})
})
