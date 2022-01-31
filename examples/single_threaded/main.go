package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/single_threaded"
)

func main() {
	// Init the pdfium library and return the instance to open documents.
	pool := single_threaded.Init()
	instance, err := pool.GetInstance(time.Second * 30)
	if err != nil {
		log.Fatal(err)
	}

	// Cleanup
	defer func() {
		instance.Close()
		pool.Close()
	}()

	filePath := "shared_tests/testdata/test.pdf"
	pageCount, err := getPageCount(instance, filePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("The PDF %s has %d page(s)", filePath, pageCount)
}

func getPageCount(instance pdfium.Pdfium, filePath string) (int, error) {
	// Load the PDF file into a byte array.
	pdfBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	// Open the PDF using pdfium (and claim a worker)
	doc, err := instance.NewDocumentFromBytes(&pdfBytes)
	if err != nil {
		return 0, err
	}

	// Always close the document, this will release the worker and it's resources
	defer instance.CloseDocument(*doc)

	pageCount, err := instance.GetPageCount(&requests.GetPageCount{
		Document: *doc,
	})
	if err != nil {
		return 0, err
	}

	return pageCount.PageCount, nil
}
