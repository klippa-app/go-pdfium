package implementation_webassembly_test

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("JPEG encoding", func() {
	DescribeTable("reuses the bitmap within the guest memory limit", func(format requests.RenderImageFormat, width, height int) {
		pool := newJPEGTestPool()
		defer pool.Close()

		instance, err := pool.GetInstance(30 * time.Second)
		Expect(err).To(Succeed())
		defer instance.Close()
		data := jpegTestPage("0 g 0 0 108 144 re f\n", "<< >>")
		doc, err := instance.OpenDocument(&requests.OpenDocument{File: &data})
		Expect(err).To(Succeed())
		defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})

		page := requests.Page{ByIndex: &requests.PageByIndex{Document: doc.Document, Index: 0}}
		pixels := requests.RenderPageInPixels{
			Page: page, Width: width, Height: height, ImageFormat: format,
		}
		dpi := requests.RenderPageInDPI{
			// The synthetic page is 3 by 2 inches.
			Page: page, DPI: width / 3, ImageFormat: format,
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

	It("releases the transparent bitmap before encoding the composited pixels", func() {
		pool := newJPEGTestPool()
		defer pool.Close()

		instance, err := pool.GetInstance(30 * time.Second)
		Expect(err).To(Succeed())
		defer instance.Close()
		// Screen blending makes PDFium request a transparent background. The
		// pixels composited over white in Go must then be copied into guest memory.
		data := jpegTestPage("/GS1 gs 1 0 0 rg 0 0 108 144 re f\n",
			"<< /ExtGState << /GS1 << /Type /ExtGState /BM /Screen /ca 0.5 >> >> >>")
		doc, err := instance.OpenDocument(&requests.OpenDocument{File: &data})
		Expect(err).To(Succeed())
		defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})

		// Render again to catch heap corruption from freeing the bitmap twice.
		for attempt := 0; attempt < 2; attempt++ {
			By(fmt.Sprintf("render %d", attempt+1))
			resp, err := instance.RenderToFile(&requests.RenderToFile{
				RenderPageInPixels: &requests.RenderPageInPixels{
					Page:  requests.Page{ByIndex: &requests.PageByIndex{Document: doc.Document, Index: 0}},
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

func newJPEGTestPool() pdfium.Pool {
	GinkgoHelper()
	config := wazero.NewRuntimeConfig()
	if os.Getenv("WAZERO_INTERPRETER") == "1" {
		config = wazero.NewRuntimeConfigInterpreter()
	}
	// The bundled PDFium module requires the exception handling core feature.
	config = config.WithCoreFeatures(api.CoreFeaturesV2 | experimental.CoreFeaturesExceptionHandling)
	pool, err := webassembly.Init(webassembly.Config{
		MinIdle: 1, MaxIdle: 1, MaxTotal: 1,
		// 768 pages = 48 MiB: enough for rendering and JPEG encoding, but
		// duplicating the ~23 MiB bitmap exhausts the guest memory budget.
		RuntimeConfig: config.WithMemoryLimitPages(768),
	})
	Expect(err).To(Succeed())
	return pool
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

func jpegTestPage(contents, resources string) []byte {
	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 216 144] /Resources %s /Contents 4 0 R >>", resources),
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
