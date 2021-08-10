package pkg

import (
	goctx "context"
	"errors"
	"fmt"
	"image"
	"os"
	"os/exec"
	"time"

	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/klippa-app/go-pdfium/pkg/internal/commons"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
)

type worker struct {
	plugin       commons.Pdfium
	pluginClient *plugin.Client
	rpcClient    plugin.ClientProtocol
}

var workerPool *pool.ObjectPool

type Config struct {
	MinIdle  int
	MaxIdle  int
	MaxTotal int
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

			binPath := "go"
			args := []string{"run", "./pkg/internal/subprocess"}

			client := plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig: handshakeConfig,
				Plugins:         pluginMap,
				Cmd:             exec.Command(binPath, args...),
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
				log.Println("[pdfium error] Worker exited")
				sentry.CaptureException(fmt.Errorf("[pdfium error] Worker exited"))
				return false
			}

			err := worker.rpcClient.Ping()
			if err != nil {
				log.Printf("[pdfium error] Error on RPC ping: %s", err.Error())
				sentry.CaptureException(fmt.Errorf("[pdfium error] Error on RPC ping: %s", err.Error()))
				return false
			}

			pong, err := worker.plugin.Ping()
			if err != nil {
				log.Printf("[pdfium error] Error on plugin ping: %s", err.Error())
				sentry.CaptureException(fmt.Errorf("[pdfium error] Error on plugin ping: %s", err.Error()))
				return false
			}

			if pong != "Pong" {
				err = errors.New("Wrong ping/pong result")
				log.Printf("[pdfium error] Error on plugin ping: %s", err.Error())
				sentry.CaptureException(fmt.Errorf("[pdfium error] Error on plugin ping: %s", err.Error()))
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

func NewDocument(file []byte) (Document, error) {
	selectedWorker, err := getWorker()
	if err != nil {
		return nil, fmt.Errorf("Could not get worker: %s", err.Error())
	}

	newDocument := pdfiumDocument{}
	newDocument.worker = selectedWorker

	err = newDocument.worker.plugin.OpenDocument(&commons.OpenDocumentRequest{File: &file})
	if err != nil {
		newDocument.Close()
		return nil, err
	}

	return &newDocument, nil
}

type Document interface {
	GetPageCount() (int, error)
	GetText(i int) string
	RenderPage(i int, dpi int) (*image.RGBA, error)
	GetPageSize(i int, dpi int) (int, int, error)
	Close()
}

type pdfiumDocument struct {
	worker *worker
}

func (d *pdfiumDocument) GetPageCount() (int, error) {
	return d.worker.plugin.GetPageCount()
}

func (d *pdfiumDocument) GetText(i int) string {
	return ""
}

func (d *pdfiumDocument) RenderPage(i int, dpi int) (*image.RGBA, error) {
	renderedPage, err := d.worker.plugin.RenderPage(&commons.RenderPageRequest{Page: i, DPI: dpi})
	if renderedPage.Image == nil {
		return nil, errors.New("Did not receive an image")
	}
	return renderedPage.Image, err
}

func (d *pdfiumDocument) GetPageSize(i int, dpi int) (int, int, error) {
	pageSize, err := d.worker.plugin.GetPageSize(&commons.GetPageSizeRequest{Page: i, DPI: dpi})
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
