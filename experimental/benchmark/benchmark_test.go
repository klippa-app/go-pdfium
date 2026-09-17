// Package benchmark compares the WebAssembly runtimes go-pdfium can run the
// PDFium module on: wazero (the webassembly package) and wazy (the
// experimental package). Every benchmark runs the same workload through the
// public pdfium.Pool API of each backend, so the numbers include go-pdfium's
// own overhead in the same way for all of them.
//
// The experimental wago backend is left out until wago's arm64 compile and
// amd64 miscompile issues are fixed, see the wago package documentation.
//
// Run with:
//
//	go test -run '^$' -bench . -benchmem -count 5 ./experimental/benchmark
//
// Set PDFIUM_BENCHMARK_RUNTIMES to a comma separated subset of wazero,wazy to
// benchmark fewer runtimes.
package benchmark

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/experimental/wazy"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
)

const testdata = "../../shared_tests/testdata/"

type runtimeSpec struct {
	name string
	init func(reuseWorkers bool) (pdfium.Pool, error)
}

var runtimes = []runtimeSpec{
	{"wazero", func(reuse bool) (pdfium.Pool, error) {
		return webassembly.Init(webassembly.Config{MinIdle: 1, MaxIdle: 1, MaxTotal: 1, ReuseWorkers: reuse})
	}},
	{"wazy", func(reuse bool) (pdfium.Pool, error) {
		return wazy.Init(wazy.Config{MinIdle: 1, MaxIdle: 1, MaxTotal: 1, ReuseWorkers: reuse})
	}},
}

func selectedRuntimes(b *testing.B) []runtimeSpec {
	env := os.Getenv("PDFIUM_BENCHMARK_RUNTIMES")
	if env == "" {
		return runtimes
	}
	var out []runtimeSpec
	for _, spec := range runtimes {
		if strings.Contains(","+env+",", ","+spec.name+",") {
			out = append(out, spec)
		}
	}
	return out
}

// forEachRuntime runs fn as a sub benchmark per runtime with an open pool and
// a borrowed instance. Runtimes that can not start are reported and skipped.
func forEachRuntime(b *testing.B, fn func(b *testing.B, instance pdfium.Pdfium)) {
	for _, spec := range selectedRuntimes(b) {
		spec := spec
		b.Run(spec.name, func(b *testing.B) {
			pool, err := spec.init(true)
			if err != nil {
				b.Skipf("%s: %v", spec.name, err)
			}
			defer pool.Close()
			instance, err := pool.GetInstance(30 * time.Second)
			if err != nil {
				b.Fatal(err)
			}
			defer instance.Close()
			fn(b, instance)
		})
	}
}

func readFile(b *testing.B, name string) []byte {
	data, err := os.ReadFile(testdata + name)
	if err != nil {
		b.Fatal(err)
	}
	return data
}

// BenchmarkInit measures creating a pool: compiling the 5.7 MB PDFium module
// to machine code and starting one worker.
func BenchmarkInit(b *testing.B) {
	for _, spec := range selectedRuntimes(b) {
		spec := spec
		b.Run(spec.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				pool, err := spec.init(true)
				if err != nil {
					b.Skipf("%s: %v", spec.name, err)
				}
				pool.Close()
			}
		})
	}
}

// BenchmarkNewInstance measures creating a worker: instantiating the compiled
// module, running its constructors and FPDF_InitLibrary. This is what every
// GetInstance costs when workers are not reused (the default).
func BenchmarkNewInstance(b *testing.B) {
	for _, spec := range selectedRuntimes(b) {
		spec := spec
		b.Run(spec.name, func(b *testing.B) {
			pool, err := spec.init(false)
			if err != nil {
				b.Skipf("%s: %v", spec.name, err)
			}
			defer pool.Close()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				instance, err := pool.GetInstance(30 * time.Second)
				if err != nil {
					b.Fatal(err)
				}
				instance.Close()
			}
		})
	}
}

// BenchmarkOpenDocument measures loading a document from memory and closing
// it again, which exercises parsing plus the host callback free path.
func BenchmarkOpenDocument(b *testing.B) {
	forEachRuntime(b, func(b *testing.B, instance pdfium.Pdfium) {
		data := readFile(b, "test.pdf")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			doc, err := instance.OpenDocument(&requests.OpenDocument{File: &data})
			if err != nil {
				b.Fatal(err)
			}
			if _, err := instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document}); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkRender renders the first page of a few documents at 100 DPI: a
// small text page, a page with an ICC based colour space and transparency, a
// page with embedded images and a large PNG predicted flate page.
func BenchmarkRender(b *testing.B) {
	for _, file := range []string{"test.pdf", "alpha_channel.pdf", "embedded_images.pdf", "rect-wrong.pdf"} {
		file := file
		b.Run(strings.TrimSuffix(file, ".pdf"), func(b *testing.B) {
			forEachRuntime(b, func(b *testing.B, instance pdfium.Pdfium) {
				data := readFile(b, file)
				doc, err := instance.OpenDocument(&requests.OpenDocument{File: &data})
				if err != nil {
					b.Fatal(err)
				}
				defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})
				page := requests.Page{ByIndex: &requests.PageByIndex{Document: doc.Document, Index: 0}}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					rendered, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{Page: page, DPI: 100})
					if err != nil {
						b.Fatal(err)
					}
					// The pixel buffer lives in guest memory until released.
					rendered.Cleanup()
				}
			})
		})
	}
}

// BenchmarkRenderToJPEG renders the first page of the image document at 150
// DPI and encodes it to JPEG, the encode running inside the module as well.
func BenchmarkRenderToJPEG(b *testing.B) {
	forEachRuntime(b, func(b *testing.B, instance pdfium.Pdfium) {
		data := readFile(b, "embedded_images.pdf")
		doc, err := instance.OpenDocument(&requests.OpenDocument{File: &data})
		if err != nil {
			b.Fatal(err)
		}
		defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})
		page := requests.Page{ByIndex: &requests.PageByIndex{Document: doc.Document, Index: 0}}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := instance.RenderToFile(&requests.RenderToFile{
				RenderPageInDPI: &requests.RenderPageInDPI{Page: page, DPI: 150},
				OutputFormat:    requests.RenderToFileOutputFormatJPG,
				OutputTarget:    requests.RenderToFileOutputTargetBytes,
				OutputQuality:   90,
			})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkGetPageText extracts the text of the first page, a call-heavy
// workload with many small host to guest transitions per character.
func BenchmarkGetPageText(b *testing.B) {
	forEachRuntime(b, func(b *testing.B, instance pdfium.Pdfium) {
		data := readFile(b, "test.pdf")
		doc, err := instance.OpenDocument(&requests.OpenDocument{File: &data})
		if err != nil {
			b.Fatal(err)
		}
		defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})
		page := requests.Page{ByIndex: &requests.PageByIndex{Document: doc.Document, Index: 0}}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, err := instance.GetPageTextStructured(&requests.GetPageTextStructured{Page: page, Mode: requests.GetPageTextStructuredModeBoth}); err != nil {
				b.Fatal(err)
			}
		}
	})
}
