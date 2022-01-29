package multi_threaded

import (
	goctx "context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/internal/commons"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	pool "github.com/jolestar/go-commons-pool/v2"
)

type worker struct {
	plugin       commons.Pdfium
	pluginClient *plugin.Client
	rpcClient    plugin.ClientProtocol
}

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

type multiThreadedPdfiumContainer struct {
	workerPool *pool.ObjectPool
}

var container *multiThreadedPdfiumContainer

func Init(config Config) pdfium.Pdfium {
	if container != nil {
		return container
	}

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

	// Create a new pdfium container and make sure we only have 1 container.
	container = &multiThreadedPdfiumContainer{
		workerPool: p,
	}

	container.workerPool.PreparePool(goctx.Background())

	return container
}

// NewDocument creates a new pdfium document from a byte array.
// This will automatically select a worker and keep it for you until you execute
// the close method on the document.
func (c *multiThreadedPdfiumContainer) NewDocument(file *[]byte, opts ...pdfium.NewDocumentOption) (pdfium.Document, error) {
	selectedWorker, err := c.getWorker()
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

func (c *multiThreadedPdfiumContainer) getWorker() (*worker, error) {
	timeout, cancel := goctx.WithTimeout(goctx.Background(), time.Second*30)
	defer cancel()
	workerObject, err := c.workerPool.BorrowObject(timeout)
	if err != nil {
		return nil, err
	}

	return workerObject.(*worker), nil
}

type pdfiumDocument struct {
	worker *worker
}

func (d *pdfiumDocument) GetPageCount(request *requests.GetPageCount) (*responses.GetPageCount, error) {
	return d.worker.plugin.GetPageCount(request)
}

func (d *pdfiumDocument) GetMetadata(request *requests.GetMetadata) (*responses.GetMetadata, error) {
	return d.worker.plugin.GetMetadata(request)
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
		container.workerPool.ReturnObject(goctx.Background(), d.worker)
		d.worker = nil
	}()

	d.worker.plugin.Close()

	return
}
