package pdfium

import (
	goctx "context"
	"errors"
	"fmt"
	"image"
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
	AlterOpenDocumentRequest(*commons.OpenDocumentRequest)
}

type openDocumentWithPassword struct{ password string }

func (p openDocumentWithPassword) AlterOpenDocumentRequest(r *commons.OpenDocumentRequest) {
	r.Password = &p.password
}

func OpenDocumentWithPasswordOption(password string) NewDocumentOption {
	return openDocumentWithPassword{
		password: password,
	}
}

func NewDocument(file *[]byte, opts ...NewDocumentOption) (Document, error) {
	selectedWorker, err := getWorker()
	if err != nil {
		return nil, fmt.Errorf("Could not get worker: %s", err.Error())
	}

	newDocument := pdfiumDocument{}
	newDocument.worker = selectedWorker

	openDocRequest := &commons.OpenDocumentRequest{File: file}
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
	GetPageCount() (int, error)
	GetPageText(page int) (*string, error)
	GetPageTextStructured(page int) (*commons.GetPageTextStructuredResponse, error)
	RenderPageInDPI(page, dpi int) (*image.RGBA, error)
	RenderPageInPixels(page, width, height int) (*image.RGBA, error)
	GetPageSize(page int) (float64, float64, error)
	GetPageSizeInPixels(page, dpi int) (int, int, error)
	Close()
}

type pdfiumDocument struct {
	worker *worker
}

func (d *pdfiumDocument) GetPageCount() (int, error) {
	return d.worker.plugin.GetPageCount()
}

func (d *pdfiumDocument) GetPageText(page int) (*string, error) {
	pageText, err := d.worker.plugin.GetPageText(&commons.GetPageTextRequest{Page: page})
	if err != nil {
		return nil, err
	}
	if pageText.Text == nil {
		return nil, errors.New("did not receive text")
	}
	return pageText.Text, err
}

func (d *pdfiumDocument) GetPageTextStructured(page int) (*commons.GetPageTextStructuredResponse, error) {
	pageText, err := d.worker.plugin.GetPageTextStructured(&commons.GetPageTextStructuredRequest{Page: page})
	if err != nil {
		return nil, err
	}

	if pageText.Chars == nil && pageText.Rects == nil {
		return nil, errors.New("did not receive structured text")
	}
	return &commons.GetPageTextStructuredResponse{
		Chars: pageText.Chars,
		Rects: pageText.Rects,
	}, err
}

func (d *pdfiumDocument) RenderPageInDPI(page, dpi int) (*image.RGBA, error) {
	renderedPage, err := d.worker.plugin.RenderPageInDPI(&commons.RenderPageInDPIRequest{Page: page, DPI: dpi})
	if err != nil {
		return nil, err
	}
	if renderedPage.Image == nil {
		return nil, errors.New("did not receive an image")
	}
	return renderedPage.Image, err
}

func (d *pdfiumDocument) RenderPageInPixels(page, width, height int) (*image.RGBA, error) {
	renderedPage, err := d.worker.plugin.RenderPageInPixels(&commons.RenderPageInPixelsRequest{Page: page, Width: width, Height: height})
	if err != nil {
		return nil, err
	}
	if renderedPage.Image == nil {
		return nil, errors.New("did not receive an image")
	}
	return renderedPage.Image, err
}

func (d *pdfiumDocument) GetPageSize(page int) (float64, float64, error) {
	pageSize, err := d.worker.plugin.GetPageSize(&commons.GetPageSizeRequest{Page: page})
	if err != nil {
		return 0, 0, err
	}

	return pageSize.Width, pageSize.Height, nil
}

func (d *pdfiumDocument) GetPageSizeInPixels(page, dpi int) (int, int, error) {
	pageSize, err := d.worker.plugin.GetPageSizeInPixels(&commons.GetPageSizeInPixelsRequest{Page: page, DPI: dpi})
	if err != nil {
		return 0, 0, err
	}

	return pageSize.Width, pageSize.Height, nil
}

func (d *pdfiumDocument) Close() {
	defer func() {
		workerPool.ReturnObject(goctx.Background(), d.worker)
		d.worker = nil
	}()

	d.worker.plugin.Close()

	return
}
