// Package wago is a WebAssembly backend for go-pdfium that runs the PDFium
// WebAssembly module with the Wago runtime (github.com/wago-org/wago), a
// pure Go ahead-of-time compiler. It is an alternative to the webassembly
// package, which uses wazero, and exposes the same pool API.
//
// Wago is still in beta and this backend is experimental. It passes
// go-pdfium's complete test suite on amd64 and arm64 with the wago commit
// pinned in go.mod; earlier wago versions could not compile or miscompiled
// the PDFium module, so do not downgrade the dependency. One known problem
// remains on arm64: wago's memory.fill clobbers a live register, which makes
// PDFium's std::fill_n on a std::vector<bool> corrupt memory in rare cases (2
// of 1,000 real world documents in one test sequence rendered differently).
// This has been reported to wago. See the README for the current state.
package wago

import (
	goctx "context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/internal/implementation_webassembly"
	"github.com/klippa-app/go-pdfium/internal/pdfium_wasm"

	"github.com/google/uuid"
	pool "github.com/jolestar/go-commons-pool/v2"
	wagort "github.com/wago-org/wago"
	"github.com/wago-org/wasi/p1"
)

// Preopen is a directory of the host filesystem that is made available to
// PDFium under GuestPath. Paths given to go-pdfium have to be absolute POSIX
// paths inside a preopen.
type Preopen = p1.Preopen

type worker struct {
	Context  goctx.Context
	Cancel   goctx.CancelFunc
	Module   *module
	Instance *implementation_webassembly.PdfiumImplementation
}

type Config struct {
	// Context is used as the parent context for the workers. It must remain
	// valid for the lifetime of the pool. If nil, context.Background is used.
	Context  goctx.Context
	MinIdle  int
	MaxIdle  int
	MaxTotal int

	// WASM is the PDFium WebAssembly module to run. Init uses the embedded
	// module when nil.
	WASM []byte

	// RuntimeConfig configures the wago compiler and runtime. The default
	// configuration works for the bundled module.
	RuntimeConfig *wagort.RuntimeConfig

	// Mounts are the host directories PDFium can access. By default the full
	// root disk is mounted read/write on non-Windows systems. On Windows the
	// volume of the current working directory is mounted as root.
	Mounts []Preopen

	Stdout       io.Writer
	Stderr       io.Writer
	RandomSource io.Reader

	// ReuseWorkers keeps workers alive when an instance is closed. By default
	// a worker is destroyed on close because creating a new one is cheap and
	// it is the only way to give memory back to the system.
	ReuseWorkers bool

	// InterruptibleCalls makes Kill able to interrupt a call that is running
	// inside PDFium, for example on a malformed document that never finishes
	// rendering. It has a small cost on every call into PDFium.
	InterruptibleCalls bool
}

type pdfiumPool struct {
	runtime      *wagort.Runtime
	module       *wagort.Module
	workerPool   *pool.ObjectPool
	instanceRefs map[string]*pdfiumInstance
	poolRef      string
	closed       bool
	lock         *sync.Mutex
	reuseWorkers bool
}

var poolRefs = map[string]*pdfiumPool{}
var multiThreadedMutex = &sync.Mutex{}

// Init will return a multithreaded wago pool running the embedded PDFium
// module. It will launch a new worker for every requested instance as long as
// the limits allow it. If the pool has been exhausted it will wait until a
// worker becomes available, so it's important that you close instances when
// you're done with them.
func Init(config Config) (pdfium.Pool, error) {
	if config.WASM == nil {
		config.WASM = pdfium_wasm.Module
	}
	return initWithConfig(config)
}

// InitWithWASM will return a multithreaded wago pool using the module in
// Config.WASM. Unlike Init, it does not reference the embedded default module,
// allowing applications that provide their own module to omit it at link time.
// This causes Go not to embed the default WASM file, saving about ~5MB in the
// resulting binary file.
func InitWithWASM(config Config) (pdfium.Pool, error) {
	if config.WASM == nil {
		return nil, errors.New("webassembly module must be provided")
	}
	return initWithConfig(config)
}

func initWithConfig(config Config) (pdfium.Pool, error) {
	// Mount the full root by default.
	if config.Mounts == nil {
		hostRoot := "/"

		// On Windows we mount the volume of the current working directory as
		// root. On Linux we mount / as root.
		if runtime.GOOS == "windows" {
			cwdDir, err := os.Getwd()
			if err != nil {
				return nil, err
			}

			volumeName := filepath.VolumeName(cwdDir)
			if volumeName != "" {
				hostRoot = fmt.Sprintf("%s\\", volumeName)
			}
		}

		config.Mounts = []Preopen{{
			GuestPath:       "/",
			HostPath:        hostRoot,
			Read:            true,
			Write:           true,
			MutateDirectory: true,
		}}
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

	poolContext := config.Context
	if poolContext == nil {
		poolContext = goctx.Background()
	}

	runtimeOptions := []wagort.RuntimeOption{}
	if config.RuntimeConfig != nil {
		runtimeOptions = append(runtimeOptions, wagort.WithRuntimeConfig(config.RuntimeConfig))
	}

	wagoRuntime := wagort.NewRuntime(runtimeOptions...)

	compiledModule, err := wagoRuntime.Compile(config.WASM)
	if err != nil {
		wagoRuntime.Close()
		return nil, fmt.Errorf("could not compile webassembly module: %w", err)
	}

	factory := pool.NewPooledObjectFactory(
		func(goctx.Context) (any, error) {
			workerCtx, cancel := goctx.WithCancel(poolContext)
			newWorker := &worker{
				Context: workerCtx,
				Cancel:  cancel,
			}

			newModule := &module{
				ctx:           workerCtx,
				compiled:      compiledModule.Compiled(),
				interruptible: config.InterruptibleCalls,
			}

			// Every instance gets its own WASI state (file descriptors) and
			// its own env module, because the env functions have to know
			// which instance they run for.
			wasi := wasiImports(workerCtx, config.Stdout, config.Stderr, config.RandomSource, config.Mounts)

			instance, err := wagoRuntime.Instantiate(workerCtx, compiledModule, wagort.WithImports(hostImports(newModule)), wagort.WithImports(wasi))
			if err != nil {
				cancel()
				return nil, fmt.Errorf("could not instantiate webassembly module: %w", err)
			}

			newModule.instance = instance
			newModule.memory = instance.Memory()
			if newModule.memory == nil {
				instance.Close()
				cancel()
				return nil, errors.New("webassembly module does not export a memory")
			}

			newWorker.Module = newModule

			fail := func(err error) (any, error) {
				instance.Close()
				cancel()
				return nil, err
			}

			// The module is a WASI reactor, _initialize runs the static
			// constructors and has to be called before anything else.
			initialize := newModule.ExportedFunction("_initialize")
			if initialize == nil {
				return fail(errors.New("could not find _initialize in exported methods"))
			}

			if _, err := initialize.Call(workerCtx); err != nil {
				return fail(fmt.Errorf("could not call _initialize: %w", err))
			}

			malloc := newModule.ExportedFunction("malloc")
			if malloc == nil {
				return fail(errors.New("could not find malloc in exported methods"))
			}

			free := newModule.ExportedFunction("free")
			if free == nil {
				return fail(errors.New("could not find free in exported methods"))
			}

			initLibrary := newModule.ExportedFunction("FPDF_InitLibrary")
			if initLibrary == nil {
				return fail(errors.New("could not find FPDF_InitLibrary in exported methods"))
			}

			if _, err := initLibrary.Call(workerCtx); err != nil {
				return fail(fmt.Errorf("could not call FPDF_InitLibrary: %w", err))
			}

			functions := map[string]implementation_webassembly.Function{
				"malloc": malloc,
				"free":   free,
			}

			newWorker.Instance = implementation_webassembly.GetInstance(workerCtx, functions, newModule)

			return newWorker, nil
		}, func(ctx goctx.Context, object *pool.PooledObject) error {
			worker := object.Object.(*worker)
			err := worker.Module.instance.Close()
			worker.Cancel()
			return err
		}, func(ctx goctx.Context, object *pool.PooledObject) bool {
			worker := object.Object.(*worker)

			pong, err := worker.Instance.Ping()
			if err != nil {
				return false
			}

			if pong != "Pong" {
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
		runtime:      wagoRuntime,
		module:       compiledModule,
		poolRef:      poolRef.String(),
		instanceRefs: map[string]*pdfiumInstance{},
		lock:         &sync.Mutex{},
		workerPool:   p,
		reuseWorkers: config.ReuseWorkers,
	}

	poolRefs[newPool.poolRef] = newPool

	return newPool, nil
}

func (p *pdfiumPool) GetInstance(timeout time.Duration) (pdfium.Pdfium, error) {
	timeoutCtx, cancel := goctx.WithTimeout(goctx.Background(), timeout)
	defer cancel()

	return p.GetInstanceWithContext(timeoutCtx)
}

func (p *pdfiumPool) GetInstanceWithContext(ctx goctx.Context) (pdfium.Pdfium, error) {
	p.lock.Lock()

	if p.closed {
		p.lock.Unlock()
		return nil, errors.New("pool is closed")
	}

	p.lock.Unlock()

	workerObject, err := p.workerPool.BorrowObject(ctx)
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
	p.lock.Lock()

	if p.closed {
		p.lock.Unlock()
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

	multiThreadedMutex.Lock()
	delete(poolRefs, p.poolRef)
	multiThreadedMutex.Unlock()

	// Close the underlying pool and destroy workers.
	p.workerPool.Close(goctx.Background())

	p.module.Close()
	p.runtime.Close()

	return nil
}

type pdfiumInstance struct {
	worker      *worker
	pool        *pdfiumPool
	instanceRef string
	// closed is atomic because Kill deliberately does not take the instance
	// lock (a stuck call may hold it) while the generated methods read it.
	closed atomic.Bool
	lock   *sync.Mutex
}

// Close will close the instance and will clean up the underlying PDFium resources.
func (i *pdfiumInstance) Close() (err error) {
	i.lock.Lock()

	if i.closed.Load() {
		i.lock.Unlock()
		return errors.New("instance is already closed")
	}

	defer func() {
		if panicError := recover(); panicError != nil {
			err = fmt.Errorf("panic occurred in %s: %v", "Close", panicError)
		}
	}()

	defer func() {
		if i.pool.reuseWorkers {
			i.pool.workerPool.ReturnObject(goctx.Background(), i.worker)
		} else {
			i.pool.workerPool.InvalidateObject(goctx.Background(), i.worker)
		}

		i.worker = nil
		i.pool.lock.Lock()
		delete(i.pool.instanceRefs, i.instanceRef)
		i.pool.lock.Unlock()
		i.pool = nil
		i.closed.Store(true)
		i.lock.Unlock()
	}()

	return i.worker.Instance.Close()
}

// Kill will destroy the underlying WebAssembly instance and remove the worker
// from the pool so that the pool can create a new one. When
// Config.InterruptibleCalls is enabled, a call that is running inside PDFium
// is interrupted.
func (i *pdfiumInstance) Kill() (err error) {
	// Kill should not be protected by a lock, since Kill is a last-effort
	// to "recover" a broken instance.
	if i.closed.Load() {
		return errors.New("instance is already closed")
	}

	defer func() {
		if panicError := recover(); panicError != nil {
			err = fmt.Errorf("panic occurred in %s: %v", "Kill", panicError)
		}
	}()

	// Cancel the worker context to interrupt any in-flight call.
	i.worker.Cancel()

	i.pool.lock.Lock()
	delete(i.pool.instanceRefs, i.instanceRef)
	i.pool.lock.Unlock()

	// Invalidate will close the instance.
	err = i.pool.workerPool.InvalidateObject(goctx.Background(), i.worker)

	i.pool = nil
	i.closed.Store(true)

	return err
}

func (i *pdfiumInstance) GetImplementation() any {
	return i.worker.Instance
}
