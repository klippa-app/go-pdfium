package main

import (
	"io/ioutil"
	"log"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/pdfium_single_threaded"
	"github.com/klippa-app/go-pdfium/requests"
)

var Pdfium pdfium.Pdfium

func init() {
	// Init the pdfium library and return the instance to open documents.
	Pdfium = pdfium_single_threaded.Init()
}

func main() {
	filePath := "internal/implementation/testdata/test.pdf"
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
