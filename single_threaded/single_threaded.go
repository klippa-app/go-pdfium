package single_threaded

import (
	"errors"
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"io"
	"sync"
	"time"
)

var singleThreadedMutex = &sync.Mutex{}

var poolRefs = map[string]*pdfiumPool{}

// Init will return a single-threaded pool.
// Every pool will keep track of its own instances and the documents that
// belong to those instances. When you close it, it will clean up the resources
// of that pool. Underwater every pool/instance uses the same mutex to ensure
// thread safety in pdfium across pools/instances/documents.
func Init() pdfium.Pool {
	singleThreadedMutex.Lock()
	defer singleThreadedMutex.Unlock()

	// Init the pdfium library.
	implementation.InitLibrary()

	poolRef := uuid.New()

	// Create a new pdfium pool
	pool := &pdfiumPool{
		poolRef:      poolRef.String(),
		instanceRefs: map[string]*pdfiumInstance{},
		lock:         &sync.Mutex{},
	}

	poolRefs[pool.poolRef] = pool

	return pool
}

type pdfiumPool struct {
	instanceRefs map[string]*pdfiumInstance
	poolRef      string
	closed       bool
	lock         *sync.Mutex
}

// GetInstance will return a unique pdfium instance that keeps track of its
// own documents. When you close it, it will clean up all resources of this
// instance.
func (p *pdfiumPool) GetInstance(timeout time.Duration) (pdfium.Pdfium, error) {
	if p.closed {
		return nil, errors.New("pool is closed")
	}

	newInstance := &pdfiumInstance{
		pdfium: implementation.Pdfium.GetInstance(),
		lock:   &sync.Mutex{},
		pool:   p,
	}

	instanceRef := uuid.New()
	newInstance.instanceRef = instanceRef.String()
	p.lock.Lock()
	p.instanceRefs[newInstance.instanceRef] = newInstance
	p.lock.Unlock()

	return newInstance, nil
}

// Close will close the pool and all instances in it.
func (p *pdfiumPool) Close() error {
	if p.closed {
		return errors.New("pool is already closed")
	}

	// Close all instances
	for i := range p.instanceRefs {
		p.instanceRefs[i].Close()
	}

	singleThreadedMutex.Lock()
	delete(poolRefs, p.poolRef)

	// Unload library if this was the last pool.
	if len(poolRefs) == 0 {
		implementation.DestroyLibrary()
	}

	singleThreadedMutex.Unlock()

	return nil
}

type pdfiumInstance struct {
	pdfium      *implementation.PdfiumImplementation
	instanceRef string
	closed      bool
	pool        *pdfiumPool
	lock        *sync.Mutex
}

// NewDocumentFromBytes creates a new pdfium references from a byte array.
func (i *pdfiumInstance) NewDocumentFromBytes(file *[]byte, opts ...pdfium.NewDocumentOption) (*references.FPDF_DOCUMENT, error) {
	i.lock.Lock()
	if i.closed {
		i.lock.Unlock()
		return nil, errors.New("instance is closed")
	}
	i.lock.Unlock()

	openDocRequest := &requests.OpenDocument{File: file}
	for _, opt := range opts {
		opt.AlterOpenDocumentRequest(openDocRequest)
	}

	doc, err := i.pdfium.OpenDocument(openDocRequest)
	if err != nil {
		return nil, err
	}

	return &doc.Document, nil
}

// NewDocumentFromFilePath creates a new pdfium references from a file path.
func (i *pdfiumInstance) NewDocumentFromFilePath(filePath string, opts ...pdfium.NewDocumentOption) (*references.FPDF_DOCUMENT, error) {
	i.lock.Lock()
	if i.closed {
		i.lock.Unlock()
		return nil, errors.New("instance is closed")
	}
	i.lock.Unlock()

	openDocRequest := &requests.OpenDocument{FilePath: &filePath}
	for _, opt := range opts {
		opt.AlterOpenDocumentRequest(openDocRequest)
	}

	doc, err := i.pdfium.OpenDocument(openDocRequest)
	if err != nil {
		return nil, err
	}

	return &doc.Document, nil
}

// NewDocumentFromReader creates a new pdfium references from a reader.
func (i *pdfiumInstance) NewDocumentFromReader(reader io.ReadSeeker, size int, opts ...pdfium.NewDocumentOption) (*references.FPDF_DOCUMENT, error) {
	i.lock.Lock()
	if i.closed {
		i.lock.Unlock()
		return nil, errors.New("instance is closed")
	}
	i.lock.Unlock()

	openDocRequest := &requests.OpenDocument{FileReader: reader, FileReaderSize: size}
	for _, opt := range opts {
		opt.AlterOpenDocumentRequest(openDocRequest)
	}

	doc, err := i.pdfium.OpenDocument(openDocRequest)
	if err != nil {
		return nil, err
	}

	return &doc.Document, nil
}

// Close will close the instance and will clean up the underlying pdfium resources
// by calling i.pdfium.Close().
func (i *pdfiumInstance) Close() error {
	i.lock.Lock()
	defer i.lock.Unlock()

	if i.closed {
		return errors.New("instance is already closed")
	}

	// Close underlying instance. That will close all docs.
	err := i.pdfium.Close()
	if err != nil {
		return err
	}

	i.pool.lock.Lock()
	delete(i.pool.instanceRefs, i.instanceRef)
	i.pool.lock.Unlock()

	// Remove references.
	i.pool = nil
	i.pdfium = nil
	i.closed = true

	return nil
}

// FPDF_CloseDocument closes a single Document and it's resources.
func (i *pdfiumInstance) FPDF_CloseDocument(document references.FPDF_DOCUMENT) error {
	if i.closed {
		return errors.New("instance is closed")
	}

	return i.pdfium.FPDF_CloseDocument(document)
}
