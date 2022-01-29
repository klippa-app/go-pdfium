package implementation

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"

import (
	"sync"
	"unsafe"

	"github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"

	"github.com/hashicorp/go-hclog"
)

// Here is the real implementation of Pdfium
type Pdfium struct {
	// C data
	currentDoc  C.FPDF_DOCUMENT
	currentPage C.FPDF_PAGE

	logger            hclog.Logger
	mutex             sync.Mutex
	currentPageNumber *int    // Remember which page is currently loaded in the page variable.
	data              *[]byte // Keep a reference to the data otherwise weird stuff happens
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

	var cPassword *C.char
	if request.Password != nil {
		cPassword = C.CString(*request.Password)
	}

	doc := C.FPDF_LoadMemDocument(
		unsafe.Pointer(&((*request.File)[0])),
		C.int(len(*request.File)),
		cPassword)

	if doc == nil {
		var pdfiumError error

		errorCode := C.FPDF_GetLastError()
		switch errorCode {
		case C.FPDF_ERR_SUCCESS:
			pdfiumError = errors.ErrSuccess
		case C.FPDF_ERR_UNKNOWN:
			pdfiumError = errors.ErrUnknown
		case C.FPDF_ERR_FILE:
			pdfiumError = errors.ErrFile
		case C.FPDF_ERR_FORMAT:
			pdfiumError = errors.ErrFormat
		case C.FPDF_ERR_PASSWORD:
			pdfiumError = errors.ErrPassword
		case C.FPDF_ERR_SECURITY:
			pdfiumError = errors.ErrSecurity
		case C.FPDF_ERR_PAGE:
			pdfiumError = errors.ErrPage
		default:
			pdfiumError = errors.ErrUnexpected
		}
		return pdfiumError
	}

	p.currentDoc = doc
	p.data = request.File
	return nil
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
