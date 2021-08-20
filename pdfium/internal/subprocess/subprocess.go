package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"

import (
	"errors"
	"os"
	"sync"
	"unsafe"

	"github.com/klippa-app/go-pdfium/pdfium/internal/commons"
	"github.com/klippa-app/go-pdfium/pdfium/requests"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

func init() {
	InitLibrary()
}

// Document is good
type Document struct {
	// C data
	doc  C.FPDF_DOCUMENT
	page C.FPDF_PAGE

	currentPage *int    // Remember which page is currently loaded in the page variable.
	data        *[]byte // Keep a reference to the data otherwise weird stuff happens
}

// Here is the real implementation of Pdfium
type Pdfium struct {
	logger     hclog.Logger
	currentDoc *Document
	mutex      sync.Mutex
}

func (p *Pdfium) Ping() (string, error) {
	return "Pong", nil
}

func (p *Pdfium) Lock() {
	p.mutex.Lock()
}

func (p *Pdfium) Unlock() {
	p.mutex.Unlock()
}

func (p *Pdfium) OpenDocument(request *requests.OpenDocument) error {
	p.Lock()
	defer p.Unlock()
	doc := C.FPDF_LoadMemDocument(
		unsafe.Pointer(&((*request.File)[0])),
		C.int(len(*request.File)),
		nil)

	if doc == nil {
		var errMsg string

		errorCode := C.FPDF_GetLastError()
		switch errorCode {
		case C.FPDF_ERR_SUCCESS:
			errMsg = "Success"
		case C.FPDF_ERR_UNKNOWN:
			errMsg = "Unknown error"
		case C.FPDF_ERR_FILE:
			errMsg = "Unable to read file"
		case C.FPDF_ERR_FORMAT:
			errMsg = "Incorrect format"
		case C.FPDF_ERR_PASSWORD:
			errMsg = "Invalid password"
		case C.FPDF_ERR_SECURITY:
			errMsg = "Invalid encryption"
		case C.FPDF_ERR_PAGE:
			errMsg = "Incorrect page"
		default:
			errMsg = "Unexpected error"
		}
		return errors.New(errMsg)
	}

	p.currentDoc = &Document{doc: doc, data: request.File}

	return nil
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

var globalMutex = &sync.Mutex{}

func InitLibrary() {
	globalMutex.Lock()
	C.FPDF_InitLibrary()
	globalMutex.Unlock()
}

func DestroyLibrary() {
	globalMutex.Lock()
	C.FPDF_DestroyLibrary()
	globalMutex.Unlock()
}
