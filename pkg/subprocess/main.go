package main

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_annot.h"
// #include "fpdf_edit.h"
// #include "fpdf_structtree.h"
// #include "fpdf_text.h"
import "C"

import (
	"errors"
	"image"
	"image/color"
	"os"
	"sync"
	"unsafe"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/klippa-app/go-pdfium/pkg/commons"
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

func (p *Pdfium) GetPageCount() (int, error) {
	if currentDoc == nil {
		return 0, errors.New("No current document")
	}
	pageCount := currentDoc.GetPageCount()
	return pageCount, nil
}

func (p *Pdfium) RenderPage(request *commons.RenderPageRequest) (commons.RenderPageResponse, error) {
	if currentDoc == nil {
		return commons.RenderPageResponse{}, errors.New("No current document")
	}
	renderedPage := currentDoc.RenderPage(request.Page, request.DPI)
	return commons.RenderPageResponse{
		Image: renderedPage,
	}, nil
}

func (p *Pdfium) GetPageSize(request *commons.GetPageSizeRequest) (commons.GetPageSizeResponse, error) {
	if currentDoc == nil {
		return commons.GetPageSizeResponse{}, errors.New("No current document")
	}
	width, height := currentDoc.GetPageSize(request.Page, request.DPI)
	return commons.GetPageSizeResponse{
		Width:  width,
		Height: height,
	}, nil
}

func (p *Pdfium) Close() error {
	if currentDoc == nil {
		return errors.New("No current document")
	}

	currentDoc.Close()
	currentDoc = nil

	return nil
}

func main() {
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

var currentDoc *Document

// Document is good
type Document struct {
	doc  C.FPDF_DOCUMENT
	data *[]byte // Keep a reference to the data otherwise weird stuff happens
}

var mutex = &sync.Mutex{}

// NewDocument creates a new pdfium doc from a byte array
func NewDocument(data *[]byte) (*Document, error) {
	mutex.Lock()
	defer mutex.Unlock()
	doc := C.FPDF_LoadMemDocument(
		unsafe.Pointer(&((*data)[0])),
		C.int(len(*data)),
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
		return nil, errors.New(errMsg)
	}
	return &Document{doc: doc, data: data}, nil
}

// GetPageCount counts the amount of pages
func (d *Document) GetPageCount() int {
	mutex.Lock()
	defer mutex.Unlock()
	return int(C.FPDF_GetPageCount(d.doc))
}

// GetText returns the text of a page
func (d *Document) GetText(i int) string {
	mutex.Lock()
	defer mutex.Unlock()
	page := C.FPDF_LoadPage(d.doc, C.int(i))
	textPage := C.FPDFText_LoadPage(page)
	charsInPage := int(C.FPDFText_CountChars(textPage))
	charData := make([]byte, charsInPage+20)

	charsWritten := C.FPDFText_GetText(textPage, C.int(0), C.int(charsInPage), (*C.ushort)(unsafe.Pointer(&charData[0])))
	C.FPDFText_ClosePage(textPage)
	C.FPDF_ClosePage(page)

	return string(charData[0:charsWritten])
}

// CloseDocument should have docs
func (d *Document) Close() {
	mutex.Lock()
	C.FPDF_CloseDocument(d.doc)
	mutex.Unlock()
}

// RenderPage renders a specific page in a specific dpi, the result is an image.
func (d *Document) RenderPage(page int, dpi int) *image.RGBA {
	mutex.Lock()

	pageObject := C.FPDF_LoadPage(d.doc, C.int(page))
	width, height := d.getPageSize(pageObject, dpi)
	alpha := C.FPDFPage_HasTransparency(pageObject)
	bitmap := C.FPDFBitmap_Create(width, height, alpha)

	fillColor := 4294967295
	if int(alpha) == 1 {
		fillColor = 0
	}

	C.FPDFBitmap_FillRect(bitmap, 0, 0, width, height, C.ulong(fillColor))
	C.FPDF_RenderPageBitmap(bitmap, pageObject, 0, 0, width, height, 0, C.FPDF_ANNOT)

	p := C.FPDFBitmap_GetBuffer(bitmap)
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	img.Stride = int(C.FPDFBitmap_GetStride(bitmap))
	mutex.Unlock()

	// This takes a bit of time and I *think* we can do this without the lock
	// @todo: figure out if we can do this better/faster.
	bgra := make([]byte, 4)
	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			for i := range bgra {
				bgra[i] = *((*byte)(p))
				p = unsafe.Pointer(uintptr(p) + 1)
			}
			pixelColor := color.RGBA{B: bgra[0], G: bgra[1], R: bgra[2], A: bgra[3]}
			img.SetRGBA(x, y, pixelColor)
		}
	}
	mutex.Lock()
	C.FPDFBitmap_Destroy(bitmap)
	C.FPDF_ClosePage(pageObject)
	mutex.Unlock()

	return img
}

// getPageSize returns the pixel size of a page given the pdfium page object DPI.
func (d *Document) getPageSize(page C.FPDF_PAGE, dpi int) (C.int, C.int) {
	// Warning: don't lock the mutex here, we're not calling this from outside.
	scale := float64(dpi) / 72.0
	imgWidth := C.FPDF_GetPageWidth(page) * C.double(scale)
	imgHeight := C.FPDF_GetPageHeight(page) * C.double(scale)

	width := C.int(imgWidth)
	height := C.int(imgHeight)

	return width, height
}

// GetPageSize returns the pixel size of a page given the page number and the DPI.
func (d *Document) GetPageSize(page int, dpi int) (int, int) {
	mutex.Lock()
	pageObject := C.FPDF_LoadPage(d.doc, C.int(page))
	width, height := d.getPageSize(pageObject, dpi)
	C.FPDF_ClosePage(pageObject)
	mutex.Unlock()

	return int(width), int(height)
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
