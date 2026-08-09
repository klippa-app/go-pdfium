package main

import (
	"image/png"
	"log"
	"os"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/single_threaded"
)

var (
	pool     pdfium.Pool
	instance pdfium.Pdfium
)

func init() {
	pool = single_threaded.Init(single_threaded.Config{})

	var err error
	if instance, err = pool.GetInstance(time.Second * 30); err != nil {
		log.Fatalf("Failed to get pdfium instance: %v", err)
	}
}

func main() {
	data, err := os.ReadFile("file.pdf")
	if err != nil {
		log.Fatalf("Failed to read PDF file: %v", err)
	}

	doc, err := instance.OpenDocument(&requests.OpenDocument{File: &data})
	if err != nil {
		log.Fatalf("Failed to open PDF document: %v", err)
	}
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})

	render(doc, 0)
}

func render(doc *responses.OpenDocument, pageIdx int) {
	render, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
		DPI: 300,
		Page: requests.Page{
			ByIndex: &requests.PageByIndex{Document: doc.Document, Index: pageIdx},
		},
	})
	if err != nil {
		log.Fatalf("Failed to render page: %v", err)
	}

	f, err := os.Create("output.png")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, render.Result.Image); err != nil {
		log.Fatalf("Failed to encode PNG: %v", err)
	}
}
