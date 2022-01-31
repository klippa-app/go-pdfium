package implementation

import (
	"errors"
	"github.com/klippa-app/go-pdfium/document"
	"io"
	"sync"
	"unsafe"

	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

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

// Pdfium is a container so that we can always only have 1 instance of pdfium
// per process. We need this so that we can guarantee thread safety.
var Pdfium = &mainPdfium{
	mutex:        &sync.Mutex{},
	instanceRefs: map[int]*PdfiumImplementation{},
	documentRefs: map[document.Ref]*NativeDocument{},
}

var isInitialized = false

// InitLibrary loads the actual C++ library.
func InitLibrary() {
	Pdfium.mutex.Lock()
	defer Pdfium.mutex.Unlock()

	// Only initialize when we aren't already.
	if isInitialized {
		return
	}

	C.FPDF_InitLibrary()
	isInitialized = true
}

// DestroyLibrary unloads the actual C++ library.
// If any documents were loaded, it closes them.
func DestroyLibrary() {
	Pdfium.mutex.Lock()
	defer Pdfium.mutex.Unlock()

	// Only destroy when we're initialized.
	if !isInitialized {
		return
	}

	for i := range Pdfium.instanceRefs {
		Pdfium.instanceRefs[i].Close()
		delete(Pdfium.instanceRefs, Pdfium.instanceRefs[i].instanceRef)
	}

	C.FPDF_DestroyLibrary()
	isInitialized = false
}

// Here is the real implementation of Pdfium
type mainPdfium struct {
	// logger is for communication with the plugin.
	logger hclog.Logger

	// mutex will ensure thread safety.
	mutex *sync.Mutex

	// instance keeps track of the opened instances for this process.
	instanceRefs map[int]*PdfiumImplementation

	// documentRefs keeps track of the opened document for this process.
	// we need this for document lookups and in case of closing the instance
	documentRefs map[document.Ref]*NativeDocument
}

func (p *mainPdfium) GetInstance() *PdfiumImplementation {
	newInstance := &PdfiumImplementation{
		logger:       p.logger,
		documentRefs: map[document.Ref]*NativeDocument{},
	}

	newInstance.instanceRef = len(p.instanceRefs)
	p.instanceRefs[newInstance.instanceRef] = newInstance

	return newInstance
}

// Here is the real implementation of Pdfium
type PdfiumImplementation struct {
	// logger is for communication with the plugin.
	logger hclog.Logger

	// documentRefs keeps track of the opened document for this instance.
	// we need this for document lookups and in case of closing the instance
	documentRefs map[document.Ref]*NativeDocument

	// We need to keep track of our own instance.
	instanceRef int
}

func (p *PdfiumImplementation) Ping() (string, error) {
	return "Pong", nil
}

func (p *PdfiumImplementation) Lock() {
	Pdfium.mutex.Lock()
}

func (p *PdfiumImplementation) Unlock() {
	Pdfium.mutex.Unlock()
}

func (p *PdfiumImplementation) OpenDocument(request *requests.OpenDocument) (*responses.OpenDocument, error) {
	p.Lock()
	defer p.Unlock()

	var cPassword *C.char
	if request.Password != nil {
		cPassword = C.CString(*request.Password)
	}

	nativeDoc := &NativeDocument{}
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
			return nil, errors.New("FileReaderSize should be given when FileReader is set")
		}

		// Allocate memory on C heap. we send the io.ReadSeeker address in this pointer.
		readSeekerAlloc := C.malloc(C.size_t(unsafe.Sizeof(uintptr(0))))

		// Create array to write the address in the array.
		a := (*[1]*io.ReadSeeker)(readSeekerAlloc)

		// Save the address in index 0 of the array.
		a[0] = &(*(*io.ReadSeeker)(unsafe.Pointer(&request.FileReader)))

		// Keep track of the allocated memory to free it later on.
		nativeDoc.readSeekerRef = readSeekerAlloc

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
		return nil, errors.New("No file given")
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
		return nil, pdfiumError
	}

	nativeDoc.currentDoc = doc
	nativeDoc.data = request.File
	nativeDoc.nativeRef = document.Ref(len(Pdfium.documentRefs) + 1)
	Pdfium.documentRefs[nativeDoc.nativeRef] = nativeDoc
	p.documentRefs[nativeDoc.nativeRef] = nativeDoc

	return &responses.OpenDocument{
		Document: nativeDoc.nativeRef,
	}, nil
}

func (p *PdfiumImplementation) CloseDocument(document document.Ref) error {
	p.Lock()
	defer p.Unlock()

	nativeDocument, err := p.getNativeDocument(document)
	if err != nil {
		return err
	}

	err = nativeDocument.Close()
	if err != nil {
		return err
	}

	delete(p.documentRefs, nativeDocument.nativeRef)

	return nil
}

func (p *PdfiumImplementation) Close() error {
	p.Lock()
	defer p.Unlock()

	for i := range p.documentRefs {
		err := p.documentRefs[i].Close()
		if err != nil {
			return err
		}

		delete(p.documentRefs, p.documentRefs[i].nativeRef)
	}

	delete(Pdfium.instanceRefs, p.instanceRef)

	return nil
}

func (p *PdfiumImplementation) getNativeDocument(documentRef document.Ref) (*NativeDocument, error) {
	if documentRef == 0 {
		return nil, errors.New("Document.Ref not given")
	}

	if val, ok := p.documentRefs[documentRef]; ok {
		return val, nil
	}

	return nil, errors.New("could not find native instance of document, perhaps the document was already closed or you tried to share documents between instances")
}

type NativeDocument struct {
	currentDoc    C.FPDF_DOCUMENT
	currentPage   C.FPDF_PAGE
	readSeekerRef unsafe.Pointer

	currentPageNumber *int         // Remember which page is currently loaded in the page variable.
	data              *[]byte      // Keep a reference to the data otherwise weird stuff happens
	nativeRef         document.Ref // An integer that is our reference inside the process. We need this to close the document in DestroyLibrary.
}

// Close closes the internal references in FPDF
func (d *NativeDocument) Close() error {
	if d.currentDoc == nil {
		return errors.New("no current document")
	}

	if d.currentPageNumber != nil {
		C.FPDF_ClosePage(d.currentPage)
		d.currentPage = nil
		d.currentPageNumber = nil
	}
	C.FPDF_CloseDocument(d.currentDoc)
	d.currentDoc = nil

	if d.readSeekerRef != nil {
		C.free(d.readSeekerRef)
		d.readSeekerRef = nil
	}

	// Remove reference to data.
	if d.data != nil {
		d.data = nil
	}

	delete(Pdfium.documentRefs, d.nativeRef)

	return nil
}
