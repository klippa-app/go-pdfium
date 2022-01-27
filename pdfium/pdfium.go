package pdfium

import (
	goctx "context"
	"errors"
	"fmt"
	"github.com/klippa-app/go-pdfium/pdfium/requests"
	"github.com/klippa-app/go-pdfium/pdfium/responses"
	"os"
	"os/exec"
	"time"

	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/klippa-app/go-pdfium/pdfium/internal/commons"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

type worker struct {
	plugin       commons.Pdfium
	pluginClient *plugin.Client
	rpcClient    plugin.ClientProtocol
}

var workerPool *pool.ObjectPool

type Config struct {
	MinIdle     int
	MaxIdle     int
	MaxTotal    int
	LogCallback func(string)
	Command     Command
}

type Command struct {
	BinPath string
	Args    []string
}

func InitLibrary(config Config) { // serve one thread that is "native" through cgo

	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	var handshakeConfig = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "BASIC_PLUGIN",
		MagicCookieValue: "hello",
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"pdfium": &commons.PdfiumPlugin{},
	}

	factory := pool.NewPooledObjectFactory(
		func(goctx.Context) (interface{}, error) {
			newWorker := &worker{}

			client := plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig: handshakeConfig,
				Plugins:         pluginMap,
				Cmd:             exec.Command(config.Command.BinPath, config.Command.Args...),
				Logger:          logger,
			})

			rpcClient, err := client.Client()
			if err != nil {
				return nil, err
			}

			raw, err := rpcClient.Dispense("pdfium")
			if err != nil {
				return nil, err
			}

			pdfium := raw.(commons.Pdfium)

			pong, err := pdfium.Ping()
			if err != nil {
				return nil, err
			}

			if pong != "Pong" {
				return nil, errors.New("Wrong ping/pong result")
			}

			newWorker.pluginClient = client
			newWorker.rpcClient = rpcClient
			newWorker.plugin = pdfium

			return newWorker, nil
		}, nil, func(ctx goctx.Context, object *pool.PooledObject) bool {
			worker := object.Object.(*worker)
			if worker.pluginClient.Exited() {
				config.LogCallback("Worker exited")
				return false
			}

			err := worker.rpcClient.Ping()
			if err != nil {
				config.LogCallback(fmt.Sprintf("Error on RPC ping: %s", err.Error()))
				return false
			}

			pong, err := worker.plugin.Ping()
			if err != nil {
				config.LogCallback(fmt.Sprintf("Error on plugin ping:: %s", err.Error()))
				return false
			}

			if pong != "Pong" {
				err = errors.New("Wrong ping/pong result")
				config.LogCallback(fmt.Sprintf("Error on plugin ping:: %s", err.Error()))
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

	workerPool = p
	workerPool.PreparePool(goctx.Background())
}

func getWorker() (*worker, error) {
	timeout, cancel := goctx.WithTimeout(goctx.Background(), time.Second*30)
	defer cancel()
	workerObject, err := workerPool.BorrowObject(timeout)
	if err != nil {
		return nil, err
	}

	return workerObject.(*worker), nil
}

type NewDocumentOption interface {
	AlterOpenDocumentRequest(*requests.OpenDocument)
}

type openDocumentWithPassword struct{ password string }

func (p openDocumentWithPassword) AlterOpenDocumentRequest(r *requests.OpenDocument) {
	r.Password = &p.password
}

// OpenDocumentWithPasswordOption can be used as NewDocumentOption when your PDF contains a password.
func OpenDocumentWithPasswordOption(password string) NewDocumentOption {
	return openDocumentWithPassword{
		password: password,
	}
}

// NewDocument creates a new pdfium document from a byte array.
// This will automatically select a worker and keep it for you until you execute
// the close method on the document.
func NewDocument(file *[]byte, opts ...NewDocumentOption) (Document, error) {
	selectedWorker, err := getWorker()
	if err != nil {
		return nil, fmt.Errorf("Could not get worker: %s", err.Error())
	}

	newDocument := pdfiumDocument{}
	newDocument.worker = selectedWorker

	openDocRequest := &requests.OpenDocument{File: file}
	for _, opt := range opts {
		opt.AlterOpenDocumentRequest(openDocRequest)
	}

	err = newDocument.worker.plugin.OpenDocument(openDocRequest)
	if err != nil {
		newDocument.Close()
		return nil, err
	}

	return &newDocument, nil
}

type Document interface {
	// GetPageCount returns the amount of pages for the document.
	GetPageCount(request *requests.GetPageCount) (*responses.GetPageCount, error)

	// GetPageText returns the text of a given page in plain text.
	GetPageText(request *requests.GetPageText) (*responses.GetPageText, error)

	// GetPageTextStructured returns the text of a given page in a structured way,
	// with coordinates and font information.
	GetPageTextStructured(request *requests.GetPageTextStructured) (*responses.GetPageTextStructured, error)

	// RenderPageInDPI renders a given page in the given DPI.
	RenderPageInDPI(request *requests.RenderPageInDPI) (*responses.RenderPage, error)

	// RenderPagesInDPI renders the given pages in the given DPI.
	RenderPagesInDPI(request *requests.RenderPagesInDPI) (*responses.RenderPages, error)

	// RenderPageInPixels renders a given page in the given pixel size.
	RenderPageInPixels(request *requests.RenderPageInPixels) (*responses.RenderPage, error)

	// RenderPagesInPixels renders the given pages in the given pixel sizes.
	RenderPagesInPixels(request *requests.RenderPagesInPixels) (*responses.RenderPages, error)

	// GetPageSize returns the size of the page in points.
	GetPageSize(request *requests.GetPageSize) (*responses.GetPageSize, error)

	// GetPageSizeInPixels returns the size of a page in pixels when rendered in the given DPI.
	GetPageSizeInPixels(request *requests.GetPageSizeInPixels) (*responses.GetPageSizeInPixels, error)

	// RenderToFile allows you to call one of the other render functions
	// and output the resulting image into a file.
	RenderToFile(request *requests.RenderToFile) (*responses.RenderToFile, error)

	// Close closes the document, releases the resources and gives back the worker to the pool.
	Close()
}

type pdfiumDocument struct {
	worker *worker
}

func (d *pdfiumDocument) GetPageCount(request *requests.GetPageCount) (*responses.GetPageCount, error) {
	return d.worker.plugin.GetPageCount(request)
}

func (d *pdfiumDocument) GetPageText(request *requests.GetPageText) (*responses.GetPageText, error) {
	return d.worker.plugin.GetPageText(request)
}

func (d *pdfiumDocument) GetPageTextStructured(request *requests.GetPageTextStructured) (*responses.GetPageTextStructured, error) {
	return d.worker.plugin.GetPageTextStructured(request)
}

func (d *pdfiumDocument) RenderPageInDPI(request *requests.RenderPageInDPI) (*responses.RenderPage, error) {
	return d.worker.plugin.RenderPageInDPI(request)
}

func (d *pdfiumDocument) RenderPagesInDPI(request *requests.RenderPagesInDPI) (*responses.RenderPages, error) {
	return d.worker.plugin.RenderPagesInDPI(request)
}

func (d *pdfiumDocument) RenderPageInPixels(request *requests.RenderPageInPixels) (*responses.RenderPage, error) {
	return d.worker.plugin.RenderPageInPixels(request)
}

func (d *pdfiumDocument) RenderPagesInPixels(request *requests.RenderPagesInPixels) (*responses.RenderPages, error) {
	return d.worker.plugin.RenderPagesInPixels(request)
}

func (d *pdfiumDocument) GetPageSize(request *requests.GetPageSize) (*responses.GetPageSize, error) {
	return d.worker.plugin.GetPageSize(request)
}

func (d *pdfiumDocument) GetPageSizeInPixels(request *requests.GetPageSizeInPixels) (*responses.GetPageSizeInPixels, error) {
	return d.worker.plugin.GetPageSizeInPixels(request)
}

func (d *pdfiumDocument) RenderToFile(request *requests.RenderToFile) (*responses.RenderToFile, error) {
	return d.worker.plugin.RenderToFile(request)
}

func (d *pdfiumDocument) Close() {
	defer func() {
		workerPool.ReturnObject(goctx.Background(), d.worker)
		d.worker = nil
	}()

	d.worker.plugin.Close()

	return
}
