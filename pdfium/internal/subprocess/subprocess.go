package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"

import (
	"errors"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"

	"github.com/klippa-app/go-pdfium/pdfium/internal/commons"
)

func init() {
	InitLibrary()
}

// Here is a real implementation of Greeter
type Pdfium struct {
	logger hclog.Logger
}

func (p *Pdfium) Ping() (string, error) {
	return "Pong", nil
}

func (p *Pdfium) OpenDocument(request *commons.OpenDocumentRequest) error {
	newDocument, err := NewDocument(request.File)
	if err != nil {
		return err
	}

	currentDoc = newDocument

	return nil
}

func (p *Pdfium) Close() error {
	if currentDoc == nil {
		return errors.New("no current document")
	}

	currentDoc.Close()
	currentDoc = nil

	return nil
}

func (p *Pdfium) GetPageText(request *commons.GetPageTextRequest) (commons.GetPageTextResponse, error) {
	if currentDoc == nil {
		return commons.GetPageTextResponse{}, errors.New("no current document")
	}
	text := currentDoc.GetPageText(request.Page)
	return commons.GetPageTextResponse{
		Text: &text,
	}, nil
}

func (p *Pdfium) GetPageSize(request *commons.GetPageSizeRequest) (commons.GetPageSizeResponse, error) {
	if currentDoc == nil {
		return commons.GetPageSizeResponse{}, errors.New("no current document")
	}
	width, height := currentDoc.GetPageSize(request.Page)
	return commons.GetPageSizeResponse{
		Width:  width,
		Height: height,
	}, nil
}

func (p *Pdfium) GetPageSizeInPixels(request *commons.GetPageSizeInPixelsRequest) (commons.GetPageSizeInPixelsResponse, error) {
	if currentDoc == nil {
		return commons.GetPageSizeInPixelsResponse{}, errors.New("no current document")
	}
	width, height := currentDoc.GetPageSizeInPixels(request.Page, request.DPI)
	return commons.GetPageSizeInPixelsResponse{
		Width:  width,
		Height: height,
	}, nil
}

func (p *Pdfium) RenderPageInDPI(request *commons.RenderPageInDPIRequest) (commons.RenderPageResponse, error) {
	if currentDoc == nil {
		return commons.RenderPageResponse{}, errors.New("no current document")
	}

	if request.DPI== 0 {
		return commons.RenderPageResponse{}, errors.New("DPI should be given")
	}

	renderedPage := currentDoc.renderPageInDPI(request.Page, request.DPI)
	return commons.RenderPageResponse{
		Image: renderedPage,
	}, nil
}

func (p *Pdfium) RenderPageInPixels(request *commons.RenderPageInPixelsRequest) (commons.RenderPageResponse, error) {
	if currentDoc == nil {
		return commons.RenderPageResponse{}, errors.New("no current document")
	}

	if request.Width == 0 && request.Height == 0 {
		return commons.RenderPageResponse{}, errors.New("either width or height should be given")
	}

	renderedPage := currentDoc.renderPageInPixels(request.Page, request.Width, request.Height)
	return commons.RenderPageResponse{
		Image: renderedPage,
	}, nil
}

func (p *Pdfium) GetPageCount() (int, error) {
	if currentDoc == nil {
		return 0, errors.New("no current document")
	}
	pageCount := currentDoc.GetPageCount()
	return pageCount, nil
}

func Main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	pdfium := &Pdfium{
		logger: logger,
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"pdfium": &commons.PdfiumPlugin{Impl: pdfium},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func InitLibrary() {
	mutex.Lock()
	C.FPDF_InitLibrary()
	mutex.Unlock()
}

func DestroyLibrary() {
	mutex.Lock()
	C.FPDF_DestroyLibrary()
	mutex.Unlock()
}
