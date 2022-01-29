package pdfium_single_threaded

import (
	"sync"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
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

// NewDocument creates a new pdfium document from a byte array.
func (c *singleThreadedPdfiumContainer) NewDocument(file *[]byte, opts ...pdfium.NewDocumentOption) (pdfium.Document, error) {
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

type pdfiumDocument struct {
	pdfium *implementation.Pdfium
}

func (d *pdfiumDocument) GetPageCount(request *requests.GetPageCount) (*responses.GetPageCount, error) {
	return d.pdfium.GetPageCount(request)
}

func (d *pdfiumDocument) GetPageText(request *requests.GetPageText) (*responses.GetPageText, error) {
	return d.pdfium.GetPageText(request)
}

func (d *pdfiumDocument) GetPageTextStructured(request *requests.GetPageTextStructured) (*responses.GetPageTextStructured, error) {
	return d.pdfium.GetPageTextStructured(request)
}

func (d *pdfiumDocument) RenderPageInDPI(request *requests.RenderPageInDPI) (*responses.RenderPage, error) {
	return d.pdfium.RenderPageInDPI(request)
}

func (d *pdfiumDocument) RenderPagesInDPI(request *requests.RenderPagesInDPI) (*responses.RenderPages, error) {
	return d.pdfium.RenderPagesInDPI(request)
}

func (d *pdfiumDocument) RenderPageInPixels(request *requests.RenderPageInPixels) (*responses.RenderPage, error) {
	return d.pdfium.RenderPageInPixels(request)
}

func (d *pdfiumDocument) RenderPagesInPixels(request *requests.RenderPagesInPixels) (*responses.RenderPages, error) {
	return d.pdfium.RenderPagesInPixels(request)
}

func (d *pdfiumDocument) GetPageSize(request *requests.GetPageSize) (*responses.GetPageSize, error) {
	return d.pdfium.GetPageSize(request)
}

func (d *pdfiumDocument) GetPageSizeInPixels(request *requests.GetPageSizeInPixels) (*responses.GetPageSizeInPixels, error) {
	return d.pdfium.GetPageSizeInPixels(request)
}

func (d *pdfiumDocument) RenderToFile(request *requests.RenderToFile) (*responses.RenderToFile, error) {
	return d.pdfium.RenderToFile(request)
}

func (d *pdfiumDocument) Close() {
	d.pdfium.Close()
	container.mutex.Unlock()

	return
}
