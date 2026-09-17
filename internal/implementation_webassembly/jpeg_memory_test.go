package implementation_webassembly_test

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// The synthetic test page is 3 by 2 inches, one point is 1/72 inch.
const (
	jpegTestPageWidthPt  = 216
	jpegTestPageHeightPt = 144
)

// jpegTestPageLeftHalf is a content stream that fills the left half of the
// synthetic test page, prefixed with the given graphics state operators.
func jpegTestPageLeftHalf(operators string) string {
	return fmt.Sprintf("%s 0 0 %d %d re f\n", operators, jpegTestPageWidthPt/2, jpegTestPageHeightPt)
}

var _ = Describe("JPEG encoding", func() {
	// The specs render within the capped guest memory of memoryLimitedPool:
	// the ~23 MiB bitmaps fit, but duplicating one exhausts the budget.
	DescribeTable("reuses the bitmap within the guest memory limit", func(format requests.RenderImageFormat, width, height int) {
		instance, document := openJPEGTestPage(jpegTestPageLeftHalf("0 g"), "<< >>")

		page := requests.Page{ByIndex: &requests.PageByIndex{Document: document, Index: 0}}
		pixels := requests.RenderPageInPixels{
			Page: page, Width: width, Height: height, ImageFormat: format,
		}
		dpi := requests.RenderPageInDPI{
			Page: page, DPI: width * 72 / jpegTestPageWidthPt, ImageFormat: format,
		}
		for _, render := range []struct {
			name    string
			request requests.RenderToFile
		}{
			{"page in pixels", requests.RenderToFile{RenderPageInPixels: &pixels}},
			{"pages in pixels", requests.RenderToFile{RenderPagesInPixels: &requests.RenderPagesInPixels{Pages: []requests.RenderPageInPixels{pixels}}}},
			{"page in DPI", requests.RenderToFile{RenderPageInDPI: &dpi}},
			{"pages in DPI", requests.RenderToFile{RenderPagesInDPI: &requests.RenderPagesInDPI{Pages: []requests.RenderPageInDPI{dpi}}}},
		} {
			By(render.name)
			render.request.OutputFormat = requests.RenderToFileOutputFormatJPG
			render.request.OutputTarget = requests.RenderToFileOutputTargetBytes
			// Reuse the instance across all four paths to exercise cleanup.
			resp, err := instance.RenderToFile(&render.request)
			Expect(err).To(Succeed())
			Expect(resp.ImageBytes).NotTo(BeNil())
			checkJPEGPixels(*resp.ImageBytes, width, height, color.RGBA{A: 255})
		}
	},
		Entry("RGBA", requests.RenderImageFormatRGBA, 3000, 2000),
		Entry("grayscale", requests.RenderImageFormatGrayscale, 6000, 4000),
	)

	It("flattens a transparent page onto white within the guest memory limit", func() {
		// Screen blending makes PDFium request a transparent background, so
		// the pixels have to be flattened onto white before encoding. Staying
		// within the memory limit proves that this does not need a second
		// full-size buffer.
		instance, document := openJPEGTestPage(jpegTestPageLeftHalf("/GS1 gs 1 0 0 rg"),
			"<< /ExtGState << /GS1 << /Type /ExtGState /BM /Screen /ca 0.5 >> >> >>")

		// Render twice to make sure the bitmap of the first render was released.
		for attempt := range 2 {
			By(fmt.Sprintf("render %d", attempt+1))
			resp, err := instance.RenderToFile(&requests.RenderToFile{
				RenderPageInPixels: &requests.RenderPageInPixels{
					Page:  requests.Page{ByIndex: &requests.PageByIndex{Document: document, Index: 0}},
					Width: 3000, Height: 2000,
				},
				OutputFormat: requests.RenderToFileOutputFormatJPG,
				OutputTarget: requests.RenderToFileOutputTargetBytes,
			})
			Expect(err).To(Succeed())
			Expect(resp.Pages).To(HaveLen(1))
			Expect(resp.Pages[0].HasTransparency).To(BeTrue())
			Expect(resp.ImageBytes).NotTo(BeNil())
			checkJPEGPixels(*resp.ImageBytes, 3000, 2000, color.RGBA{R: 255, G: 128, B: 128, A: 255})
		}
	})
})

// openJPEGTestPage opens a synthetic test page in an instance of the memory
// limited pool. See openMemoryLimitedDocument for the cleanup.
func openJPEGTestPage(contents, resources string) (pdfium.Pdfium, references.FPDF_DOCUMENT) {
	GinkgoHelper()
	return openMemoryLimitedDocument(jpegTestPage(contents, resources))
}

// openMemoryLimitedDocument takes the instance of memoryLimitedPool and opens
// the given PDF in it. The document and the instance are released when the
// spec ends, in that order.
func openMemoryLimitedDocument(data []byte) (pdfium.Pdfium, references.FPDF_DOCUMENT) {
	GinkgoHelper()

	instance, err := memoryLimitedPool.GetInstance(30 * time.Second)
	Expect(err).To(Succeed())
	DeferCleanup(instance.Close)

	doc, err := instance.OpenDocument(&requests.OpenDocument{File: &data})
	Expect(err).To(Succeed())
	DeferCleanup(func() error {
		_, err := instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})
		return err
	})

	return instance, doc.Document
}

func checkJPEGPixels(data []byte, width, height int, left color.RGBA) {
	GinkgoHelper()
	img, err := jpeg.Decode(bytes.NewReader(data))
	Expect(err).To(Succeed())
	Expect(img.Bounds().Size()).To(Equal(image.Pt(width, height)))
	for _, sample := range []struct {
		x    int
		want color.RGBA
	}{
		{width / 4, left},
		{3 * width / 4, color.RGBA{R: 255, G: 255, B: 255, A: 255}},
	} {
		got := color.RGBAModel.Convert(img.At(sample.x, height/2)).(color.RGBA)
		for _, delta := range []int{int(got.R) - int(sample.want.R), int(got.G) - int(sample.want.G), int(got.B) - int(sample.want.B)} {
			Expect(delta).To(BeNumerically("~", 0, 3),
				"JPEG pixel at (%d, %d) = %v; want approximately %v", sample.x, height/2, got, sample.want)
		}
	}
}

// jpegTestPage builds a single page PDF with the given content stream and
// resources dictionary.
func jpegTestPage(contents, resources string) []byte {
	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %d %d] /Resources %s /Contents 4 0 R >>", jpegTestPageWidthPt, jpegTestPageHeightPt, resources),
		fmt.Sprintf("<< /Length %d >>\nstream\n%sendstream", len(contents), contents),
	}
	var pdf bytes.Buffer
	pdf.WriteString("%PDF-1.4\n")
	offsets := make([]int, len(objects)+1)
	for i, object := range objects {
		offsets[i+1] = pdf.Len()
		fmt.Fprintf(&pdf, "%d 0 obj\n%s\nendobj\n", i+1, object)
	}
	xref := pdf.Len()
	fmt.Fprintf(&pdf, "xref\n0 %d\n0000000000 65535 f \n", len(offsets))
	for _, offset := range offsets[1:] {
		fmt.Fprintf(&pdf, "%010d 00000 n \n", offset)
	}
	fmt.Fprintf(&pdf, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(offsets), xref)
	return pdf.Bytes()
}
