package webassembly

import (
	goctx "context"
	"crypto/rand"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sync"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/internal/implementation_webassembly"
	"github.com/klippa-app/go-pdfium/webassembly/imports"

	"github.com/google/uuid"
	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"golang.org/x/net/context"
)

//go:embed pdfium.wasm
var pdfiumWasm []byte

type worker struct {
	Context   context.Context
	Functions map[string]api.Function
	Module    api.Module
	Instance  *implementation_webassembly.PdfiumImplementation
}

type Config struct {
	MinIdle      int
	MaxIdle      int
	MaxTotal     int
	WASM         []byte
	FS           fs.FS
	Stdout       io.Writer
	Stderr       io.Writer
	RandomSource io.Reader
}

type pdfiumPool struct {
	runtime        wazero.Runtime
	compiledModule wazero.CompiledModule
	workerPool     *pool.ObjectPool
	instanceRefs   map[string]*pdfiumInstance
	poolRef        string
	closed         bool
	lock           *sync.Mutex
}

var poolRefs = map[string]*pdfiumPool{}
var multiThreadedMutex = &sync.Mutex{}

// Init will return a multithreaded webassembly pool.
// It will launch a new worker for every requested instance as long as the limits
// allow it. If the pool has been exhausted. It will wait until a worker becomes
// available. So it's important that you close instances when you're done with them.
func Init(config Config) (pdfium.Pool, error) {
	poolContext := context.Background()
	runtime := wazero.NewRuntimeWithConfig(poolContext, wazero.NewRuntimeConfig())

	// Import WASI features.
	if _, err := wasi_snapshot_preview1.Instantiate(poolContext, runtime); err != nil {
		runtime.Close(poolContext)
		return nil, fmt.Errorf("could not instantiate webassembly wasi_snapshot_preview1 module: %w", err)
	}

	// Add missing emscripten and syscalls.
	if _, err := imports.Instantiate(poolContext, runtime); err != nil {
		runtime.Close(poolContext)
		return nil, fmt.Errorf("could not instantiate webassembly emscripten/env module: %w", err)
	}

	compiledModule, err := runtime.CompileModule(poolContext, pdfiumWasm)
	if err != nil {
		runtime.Close(poolContext)
		return nil, fmt.Errorf("could not compile webassembly module: %w", err)
	}

	factory := pool.NewPooledObjectFactory(
		func(goctx.Context) (interface{}, error) {
			if config.WASM == nil {
				config.WASM = pdfiumWasm
			}

			if config.FS == nil {
				config.FS = os.DirFS("")
			}

			if config.Stderr == nil {
				config.Stderr = os.Stderr
			}

			if config.Stdout == nil {
				config.Stdout = os.Stdout
			}

			if config.RandomSource == nil {
				config.RandomSource = rand.Reader
			}

			newWorker := &worker{
				Context: poolContext,
			}

			moduleConfig := wazero.NewModuleConfig().
				WithStartFunctions("_initialize").
				WithStdout(config.Stdout).
				WithStderr(config.Stderr).
				WithRandSource(config.RandomSource).
				WithFS(config.FS)

			mod, err := runtime.InstantiateModule(newWorker.Context, compiledModule, moduleConfig)
			if err != nil {
				return nil, fmt.Errorf("could not instantiate webassembly module: %w", err)
			}

			newWorker.Module = mod

			malloc := mod.ExportedFunction("malloc")
			free := mod.ExportedFunction("free")

			newWorker.Functions = map[string]api.Function{
				"malloc": malloc,
				"free":   free,
			}

			_, err = mod.ExportedFunction("FPDF_InitLibrary").Call(newWorker.Context)
			if err != nil {
				return nil, fmt.Errorf("could not call FPDF_InitLibrary: %w", err)
			}

			newWorker.Instance = implementation_webassembly.GetInstance(newWorker.Context, newWorker.Functions, newWorker.Module)

			return newWorker, nil
		}, func(ctx goctx.Context, object *pool.PooledObject) error {
			worker := object.Object.(*worker)
			err := worker.Module.Close(worker.Context)
			if err != nil {
				return err
			}

			worker = nil
			return err
		}, func(ctx goctx.Context, object *pool.PooledObject) bool {
			worker := object.Object.(*worker)
			// @todo: how do to alive check?
			// @todo: do we need to do an alive check?

			pong, err := worker.Instance.Ping()
			if err != nil {
				return false
			}

			if pong != "Pong" {
				err = errors.New("Wrong ping/pong result")
				return false
			}

			return true
		}, nil, nil)
	p := pool.NewObjectPoolWithDefaultConfig(goctx.Background(), factory)
	p.Config = &pool.ObjectPoolConfig{
		BlockWhenExhausted: true,
		MinIdle:            config.MinIdle,
		MaxIdle:            config.MaxIdle,
		MaxTotal:           config.MaxTotal,
		TestOnBorrow:       true,
		TestOnReturn:       true,
		TestOnCreate:       true,
	}

	p.PreparePool(goctx.Background())

	multiThreadedMutex.Lock()
	defer multiThreadedMutex.Unlock()

	poolRef := uuid.New()

	// Create a new PDFium pool.
	newPool := &pdfiumPool{
		runtime:        runtime,
		compiledModule: compiledModule,
		poolRef:        poolRef.String(),
		instanceRefs:   map[string]*pdfiumInstance{},
		lock:           &sync.Mutex{},
		workerPool:     p,
	}

	poolRefs[newPool.poolRef] = newPool

	return newPool, nil
}

func (p *pdfiumPool) GetInstance(timeout time.Duration) (pdfium.Pdfium, error) {
	if p.closed {
		return nil, errors.New("pool is closed")
	}

	timeoutCtx, cancel := goctx.WithTimeout(goctx.Background(), timeout)
	defer cancel()
	workerObject, err := p.workerPool.BorrowObject(timeoutCtx)
	if err != nil {
		return nil, err
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	newInstance := &pdfiumInstance{
		worker: workerObject.(*worker),
		lock:   &sync.Mutex{},
	}

	instanceRef := uuid.New()
	newInstance.instanceRef = instanceRef.String()
	newInstance.pool = p
	p.instanceRefs[newInstance.instanceRef] = newInstance

	return newInstance, nil
}

func (p *pdfiumPool) Close() (err error) {
	if p.closed {
		return errors.New("pool is already closed")
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	defer func() {
		if panicError := recover(); panicError != nil {
			err = fmt.Errorf("panic occurred in %s: %v", "Close", panicError)
		}
	}()

	// Close all instances
	for i := range p.instanceRefs {
		p.instanceRefs[i].worker = nil
		p.instanceRefs[i].pool = nil
		p.instanceRefs[i].closed = true

		delete(p.instanceRefs, i)
	}

	multiThreadedMutex.Lock()
	delete(poolRefs, p.poolRef)
	multiThreadedMutex.Unlock()

	// Close the underlying pool and destroy workers.
	p.workerPool.Close(goctx.Background())

	p.runtime.Close(goctx.Background())

	return nil
}

type pdfiumInstance struct {
	worker      *worker
	pool        *pdfiumPool
	instanceRef string
	closed      bool
	lock        *sync.Mutex
}

// Close will close the instance and will clean up the underlying PDFium resources
// by calling i.worker.plugin.Close().
func (i *pdfiumInstance) Close() (err error) {
	i.lock.Lock()

	if i.closed {
		i.lock.Unlock()
		return errors.New("instance is already closed")
	}

	defer func() {
		if panicError := recover(); panicError != nil {
			err = fmt.Errorf("panic occurred in %s: %v", "Close", panicError)
		}
	}()

	defer func() {
		i.pool.workerPool.ReturnObject(goctx.Background(), i.worker)
		i.worker = nil
		delete(i.pool.instanceRefs, i.instanceRef)
		i.pool = nil
		i.closed = true
		i.lock.Unlock()
	}()

	return i.worker.Instance.Close()
}

// Kill will kill the actual subprocess and return the worker to the pool
// so that the pool system can re-create the process.
func (i *pdfiumInstance) Kill() (err error) {
	// Kill should not be protected by a lock, since Kill is a last-effort
	// to "recover" a broken process.
	if i.closed {
		return errors.New("instance is already closed")
	}

	defer func() {
		if panicError := recover(); panicError != nil {
			err = fmt.Errorf("panic occurred in %s: %v", "Close", panicError)
		}
	}()

	delete(i.pool.instanceRefs, i.instanceRef)
	i.pool = nil
	i.closed = true

	// Invalidate will close the module.
	return i.pool.workerPool.InvalidateObject(goctx.Background(), i.worker)
}
