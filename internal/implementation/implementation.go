package implementation

import (
	"errors"
	"io"
	"sync"
	"unsafe"

	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"

	"github.com/hashicorp/go-hclog"
)

/*
#cgo pkg-config: pdfium
#include "fpdfview.h"
#include <stdlib.h>

extern int go_read_seeker_cb(void *param, unsigned long position, unsigned char *pBuf, unsigned long size);

static inline void FPDF_FILEACCESS_SET_GET_BLOCK(FPDF_FILEACCESS *fs) {
	fs->m_GetBlock = &go_read_seeker_cb;
}
*/
import "C"

// go_read_seeker_cb is the Go implementation of FPDF_FILEACCESS::m_GetBlock.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FILEACCESS structs. It contains a lot of tricks to make this work,
// it has a pointer to the original ReadSeeker, and it also converts the
// pBuf *C.uchar into a Go []byte array so that we can directly read from the
// readSeeker into the byte array.
//export go_read_seeker_cb
func go_read_seeker_cb(param unsafe.Pointer, position C.ulong, pBuf *C.uchar, size C.ulong) C.int {
	r := *(*io.ReadSeeker)((*[1]*io.ReadSeeker)(param)[0])

	_, err := r.Seek(int64(position), 0)
	if err != nil {
		return C.int(0)
	}

	// We create a Go slice backed by a C array (without copying the original data),
	// and acquire its length at runtime and use a type conversion to a pointer to a very big array and then slice it to the length that we want.
	// Refer https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
	target := (*[1<<50 - 1]byte)(unsafe.Pointer(pBuf))[:size:size] // For 64-bit machine, the max number it can go is 50 as per https://github.com/golang/go/issues/13656#issuecomment-291957684
	readBytes, err := r.Read(target)
	if err != nil {
		return C.int(0)
	}

	if readBytes == 0 {
		return C.int(0)
	}

	// A integer value: non-zero for success while zero for error.
	return C.int(readBytes)
}

// Here is the real implementation of Pdfium
type Pdfium struct {
	// C data
	currentDoc    C.FPDF_DOCUMENT
	currentPage   C.FPDF_PAGE
	readSeekerRef unsafe.Pointer

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

	var doc C.FPDF_DOCUMENT

	if request.File != nil {
		doc = C.FPDF_LoadMemDocument(
			unsafe.Pointer(&((*request.File)[0])),
			C.int(len(*request.File)),
			cPassword)
	} else if request.FilePath != nil {
		filePath := C.CString(*request.FilePath)
		defer C.free(unsafe.Pointer(filePath))
		doc = C.FPDF_LoadDocument(
			filePath,
			cPassword)
	} else if request.FileReader != nil {
		if request.FileReaderSize == 0 {
			return errors.New("FileReaderSize should be given when FileReader is set")
		}

		// Allocate memory on C heap. we send the io.ReadSeeker address in this pointer.
		readSeekerAlloc := C.malloc(C.size_t(unsafe.Sizeof(uintptr(0))))

		// Create array to write the address in the array.
		a := (*[1]*io.ReadSeeker)(readSeekerAlloc)

		// Save the address in index 0 of the array.
		a[0] = &(*(*io.ReadSeeker)(unsafe.Pointer(&request.FileReader)))

		// Keep track of the allocated memory to free it later on.
		p.readSeekerRef = readSeekerAlloc

		// Create a pdfium file access struct.
		readerStruct := C.FPDF_FILEACCESS{}
		readerStruct.m_FileLen = C.ulong(request.FileReaderSize)
		readerStruct.m_Param = readSeekerAlloc

		// Set the Go callback through cgo.
		C.FPDF_FILEACCESS_SET_GET_BLOCK(&readerStruct)

		doc = C.FPDF_LoadCustomDocument(
			&readerStruct,
			cPassword)
	} else {
		return errors.New("No file given")
	}

	if doc == nil {
		var pdfiumError error

		errorCode := C.FPDF_GetLastError()
		switch errorCode {
		case C.FPDF_ERR_SUCCESS:
			pdfiumError = pdfium_errors.ErrSuccess
		case C.FPDF_ERR_UNKNOWN:
			pdfiumError = pdfium_errors.ErrUnknown
		case C.FPDF_ERR_FILE:
			pdfiumError = pdfium_errors.ErrFile
		case C.FPDF_ERR_FORMAT:
			pdfiumError = pdfium_errors.ErrFormat
		case C.FPDF_ERR_PASSWORD:
			pdfiumError = pdfium_errors.ErrPassword
		case C.FPDF_ERR_SECURITY:
			pdfiumError = pdfium_errors.ErrSecurity
		case C.FPDF_ERR_PAGE:
			pdfiumError = pdfium_errors.ErrPage
		default:
			pdfiumError = pdfium_errors.ErrUnexpected
		}
		return pdfiumError
	}

	p.currentDoc = doc
	p.data = request.File
	return nil
}

// Close closes the internal references in FPDF
func (p *Pdfium) Close() error {
	p.Lock()
	defer p.Unlock()

	if p.currentDoc == nil {
		return errors.New("no current document")
	}

	if p.currentPageNumber != nil {
		C.FPDF_ClosePage(p.currentPage)
		p.currentPage = nil
		p.currentPageNumber = nil
	}
	C.FPDF_CloseDocument(p.currentDoc)
	p.currentDoc = nil

	if p.readSeekerRef != nil {
		C.free(p.readSeekerRef)
		p.readSeekerRef = nil
	}

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
