package single_threaded

import (
	"io"
	"sync"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/requests"
)

type singleThreadedPdfiumContainer struct {
	pdfium *implementation.Pdfium
	mutex  *sync.Mutex
}

var container *singleThreadedPdfiumContainer

func Init() pdfium.Pdfium {
	if container != nil {
		return container
	}

	// Init the pdfium library.
	implementation.InitLibrary()

	// Create a new pdfium container and ake sure we only have 1 container.
	container = &singleThreadedPdfiumContainer{
		pdfium: &implementation.Pdfium{},
		mutex:  &sync.Mutex{},
	}

	return container
}

func Destroy() {
	if container == nil {
		return
	}

	container.pdfium.Close()
	implementation.DestroyLibrary()
}

// NewDocumentFromBytes creates a new pdfium document from a byte array.
func (c *singleThreadedPdfiumContainer) NewDocumentFromBytes(file *[]byte, opts ...pdfium.NewDocumentOption) (pdfium.Document, error) {
	// Make sure there can only be one document at the same time.
	c.mutex.Lock()

	newDocument := pdfiumDocument{
		pdfium: c.pdfium,
	}

	openDocRequest := &requests.OpenDocument{File: file}
	for _, opt := range opts {
		opt.AlterOpenDocumentRequest(openDocRequest)
	}

	err := c.pdfium.OpenDocument(openDocRequest)
	if err != nil {
		newDocument.Close()
		return nil, err
	}

	return &newDocument, nil
}

// NewDocumentFromFilePath creates a new pdfium document from a file path.
func (c *singleThreadedPdfiumContainer) NewDocumentFromFilePath(filePath string, opts ...pdfium.NewDocumentOption) (pdfium.Document, error) {
	// Make sure there can only be one document at the same time.
	c.mutex.Lock()

	newDocument := pdfiumDocument{
		pdfium: c.pdfium,
	}

	openDocRequest := &requests.OpenDocument{FilePath: &filePath}
	for _, opt := range opts {
		opt.AlterOpenDocumentRequest(openDocRequest)
	}

	err := c.pdfium.OpenDocument(openDocRequest)
	if err != nil {
		newDocument.Close()
		return nil, err
	}

	return &newDocument, nil
}

// NewDocumentFromReader creates a new pdfium document from a reader.
func (c *singleThreadedPdfiumContainer) NewDocumentFromReader(reader io.ReadSeeker, size int, opts ...pdfium.NewDocumentOption) (pdfium.Document, error) {
	// Make sure there can only be one document at the same time.
	c.mutex.Lock()

	newDocument := pdfiumDocument{
		pdfium: c.pdfium,
	}

	openDocRequest := &requests.OpenDocument{FileReader: reader, FileReaderSize: size}
	for _, opt := range opts {
		opt.AlterOpenDocumentRequest(openDocRequest)
	}

	err := c.pdfium.OpenDocument(openDocRequest)
	if err != nil {
		newDocument.Close()
		return nil, err
	}

	return &newDocument, nil
}

type pdfiumDocument struct {
	pdfium *implementation.Pdfium
}

func (d *pdfiumDocument) Close() {
	d.pdfium.Close()
	container.mutex.Unlock()

	return
}
