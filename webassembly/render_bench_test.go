package webassembly_test

import (
	"os"
	"testing"
	"time"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
)

func BenchmarkRenderPage(b *testing.B) {
	pool, err := webassembly.Init(webassembly.Config{
		MinIdle:  1,
		MaxIdle:  1,
		MaxTotal: 1,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer pool.Close()

	instance, err := pool.GetInstance(time.Second * 30)
	if err != nil {
		b.Fatal(err)
	}
	defer instance.Close()

	pdfBytes, err := os.ReadFile("../shared_tests/testdata/rect-wrong.pdf")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc, err := instance.OpenDocument(&requests.OpenDocument{File: &pdfBytes})
		if err != nil {
			b.Fatal(err)
		}

		_, err = instance.RenderPageInDPI(&requests.RenderPageInDPI{
			DPI: 200,
			Page: requests.Page{
				ByIndex: &requests.PageByIndex{
					Document: doc.Document,
					Index:    0,
				},
			},
		})
		if err != nil {
			b.Fatal(err)
		}

		instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})
	}
}
