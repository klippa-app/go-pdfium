package single_threaded

import (
	"errors"
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/document"
	"github.com/klippa-app/go-pdfium/internal/implementation"
	"github.com/klippa-app/go-pdfium/requests"
	"io"
	"sync"
	"time"
)

var singleThreadedMutex = &sync.Mutex{}

var poolRefs = map[int]*pdfiumPool{}

func Init() pdfium.Pool {
	singleThreadedMutex.Lock()
	defer singleThreadedMutex.Unlock()

	// Init the pdfium library.
	implementation.InitLibrary()

	// Create a new pdfium pool
	pool := &pdfiumPool{
		poolRef:      len(poolRefs),
		instanceRefs: map[int]*pdfiumInstance{},
		lock:         &sync.Mutex{},
	}

	poolRefs[pool.poolRef] = pool

	return pool
}

type pdfiumPool struct {
	instanceRefs map[int]*pdfiumInstance
	poolRef      int
	closed       bool
	lock         *sync.Mutex
}

func (p *pdfiumPool) GetInstance(timeout time.Duration) (pdfium.Pdfium, error) {
	if p.closed {
		return nil, errors.New("pool is closed")
	}

	newInstance := &pdfiumInstance{
		pdfium: implementation.Pdfium.GetInstance(),
		lock:   &sync.Mutex{},
		pool:   p,
	}

	instanceRef := len(p.instanceRefs)
	newInstance.instanceRef = instanceRef
	p.lock.Lock()
	p.instanceRefs[instanceRef] = newInstance
	p.lock.Unlock()

	return newInstance, nil
}

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
	instanceRef int
	closed      bool
	pool        *pdfiumPool
	lock        *sync.Mutex
}

// NewDocumentFromBytes creates a new pdfium document from a byte array.
func (i *pdfiumInstance) NewDocumentFromBytes(file *[]byte, opts ...pdfium.NewDocumentOption) (*document.Ref, error) {
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

// NewDocumentFromFilePath creates a new pdfium document from a file path.
func (i *pdfiumInstance) NewDocumentFromFilePath(filePath string, opts ...pdfium.NewDocumentOption) (*document.Ref, error) {
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

// NewDocumentFromReader creates a new pdfium document from a reader.
func (i *pdfiumInstance) NewDocumentFromReader(reader io.ReadSeeker, size int, opts ...pdfium.NewDocumentOption) (*document.Ref, error) {
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

func (i *pdfiumInstance) CloseDocument(document document.Ref) error {
	if i.closed {
		return errors.New("instance is closed")
	}

	return i.pdfium.CloseDocument(document)
}
