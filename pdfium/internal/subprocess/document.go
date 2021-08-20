package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"

import (
	"errors"
	"sync"
	"unsafe"
)

var currentDoc *Document

// Document is good
type Document struct {
	// C data
	doc  C.FPDF_DOCUMENT
	page C.FPDF_PAGE

	currentPage *int    // Remember which page is currently loaded in the page variable.
	data        *[]byte // Keep a reference to the data otherwise weird stuff happens
}

var mutex = &sync.Mutex{}

// NewDocument creates a new pdfium doc from a byte array
func NewDocument(data *[]byte, password *string) (*Document, error) {
	var cPassword *C.char
	if password != nil {
		cPassword = C.CString(*password)
	}

	mutex.Lock()
	defer mutex.Unlock()
	doc := C.FPDF_LoadMemDocument(
		unsafe.Pointer(&((*data)[0])),
		C.int(len(*data)),
		cPassword)

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

// Close closes the internal references in FPDF
func (d *Document) Close() {
	mutex.Lock()
	if d.currentPage != nil {
		C.FPDF_ClosePage(d.page)
		d.page = nil
		d.currentPage = nil
	}
	C.FPDF_CloseDocument(d.doc)
	d.doc = nil
	mutex.Unlock()
}
