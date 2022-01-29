package main

import (
	"io/ioutil"
	"log"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/multi_threaded"
	"github.com/klippa-app/go-pdfium/requests"
)

var Pdfium pdfium.Pdfium

func init() {
	// Init the pdfium library and return the instance to open documents.
	// You can tweak these configs to your need. Be aware that workers can use quite some memory.
	Pdfium = multi_threaded.Init(multi_threaded.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
		Command: multi_threaded.Command{
			BinPath: "go",                                                      // Only do this while developing, on production put the actual binary path in here. You should not want the Go runtime on production.
			Args:    []string{"run", "examples/multi_threaded/worker/main.go"}, // This is a reference to the worker package, this can be left empty when using a direct binary path.
		},
	})
}

func main() {
	filePath := "pdfium/internal/implementation/testdata/test.pdf"
	pageCount, err := getPageCount(filePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("The PDF %s has %d page(s)", filePath, pageCount)
}

func getPageCount(filePath string) (int, error) {
	// Load the PDF file into a byte array.
	pdfBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	// Open the PDF using pdfium (and claim a worker)
	doc, err := Pdfium.NewDocument(&pdfBytes)
	if err != nil {
		return 0, err
	}

	// Always close the document, this will release the worker and it's resources
	defer doc.Close()

	pageCount, err := doc.GetPageCount(&requests.GetPageCount{})
	if err != nil {
		return 0, err
	}

	return pageCount.PageCount, nil
}
