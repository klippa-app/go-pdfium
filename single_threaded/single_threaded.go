package single_threaded

import (
	goctx "context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/internal/implementation_cgo"
)

var singleThreadedMutex = &sync.Mutex{}

var poolRefs = map[string]*pdfiumPool{}

type Config struct {
	LibraryConfig *pdfium.LibraryConfig
}

// Init will return a single-threaded pool.
// Every pool will keep track of its own instances and the documents that
// belong to those instances. When you close it, it will clean up the resources
// of that pool. Underwater every pool/instance uses the same mutex to ensure
// thread safety in PDFium across pools/instances/documents.
func Init(config Config) pdfium.Pool {
	singleThreadedMutex.Lock()
	defer singleThreadedMutex.Unlock()

	// Init the PDFium library.
	implementation_cgo.InitLibrary(config.LibraryConfig)

	poolRef := uuid.New()

	// Create a new PDFium pool
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

// GetInstance will return a unique PDFium instance that keeps track of its
// own documents. When you close it, it will clean up all resources of this
// instance.
func (p *pdfiumPool) GetInstance(timeout time.Duration) (pdfium.Pdfium, error) {
	return p.GetInstanceWithContext(goctx.Background())
}

func (p *pdfiumPool) GetInstanceWithContext(ctx goctx.Context) (pdfium.Pdfium, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.closed {
		return nil, errors.New("pool is closed")
	}

	newInstance := &pdfiumInstance{
		pdfium: implementation_cgo.Pdfium.GetInstance(),
		lock:   &sync.Mutex{},
		pool:   p,
	}

	instanceRef := uuid.New()
	newInstance.instanceRef = instanceRef.String()
	p.instanceRefs[newInstance.instanceRef] = newInstance

	return newInstance, nil
}

// Close will close the pool and all instances in it.
func (p *pdfiumPool) Close() (err error) {
	p.lock.Lock()

	if p.closed {
		return errors.New("pool is already closed")
	}

	// Once we mark the pool as closed, the user can't do anything to change
	// the pool, except closing instances, which has its own lock anyway.
	p.closed = true
	p.lock.Unlock()

	defer func() {
		if panicError := recover(); panicError != nil {
			err = fmt.Errorf("panic occurred in %s: %v", "Close", panicError)
		}
	}()

	// Close all instances
	for i := range p.instanceRefs {
		p.instanceRefs[i].Close()
	}

	singleThreadedMutex.Lock()
	delete(poolRefs, p.poolRef)

	// Unload library if this was the last pool.
	if len(poolRefs) == 0 {
		implementation_cgo.DestroyLibrary()
	}

	singleThreadedMutex.Unlock()

	return nil
}

type pdfiumInstance struct {
	pdfium      *implementation_cgo.PdfiumImplementation
	instanceRef string
	closed      bool
	pool        *pdfiumPool
	lock        *sync.Mutex
}

// Close will close the instance and will clean up the underlying PDFium resources
// by calling i.pdfium.Close().
func (i *pdfiumInstance) Close() (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if i.closed {
		return errors.New("instance is already closed")
	}

	defer func() {
		if panicError := recover(); panicError != nil {
			err = fmt.Errorf("panic occurred in %s: %v", "NewDocumentFromReader", panicError)
		}
	}()

	// Close underlying instance. That will close all docs.
	err = i.pdfium.Close()
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

// Kill is the same as Close on single-threaded usage.
func (i *pdfiumInstance) Kill() (err error) {
	return i.Close()
}

func (i *pdfiumInstance) GetImplementation() interface{} {
	return i.pdfium
}
