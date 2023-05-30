package implementation_cgo

import "C"
import (
	"bytes"
	"errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"io"
	"os"
	"unsafe"
)

/*
#cgo pkg-config: pdfium
#include "fpdf_save.h"
#include <stdlib.h>

typedef const void cvoid_t;
extern int go_writer_cb(struct FPDF_FILEWRITE_ *pThis, cvoid_t *pData, unsigned long size);

static inline void FPDF_FILEWRITE_SET_WRITE_BLOCK(FPDF_FILEWRITE *fs) {
	fs->WriteBlock = &go_writer_cb;
}
*/
import "C"

var currentWriter io.Writer

//export go_writer_cb
func go_writer_cb(pThis *C.FPDF_FILEWRITE, pData *C.cvoid_t, size C.ulong) C.int {
	// We create a Go slice backed by a C array (without copying the original data).
	data := unsafe.Slice((*byte)(unsafe.Pointer(pData)), uint64(size))

	if currentWriter == nil {
		return C.int(0)
	}

	writtenBytes, err := currentWriter.Write(data)
	if err != nil {
		return C.int(0)
	}

	// An integer value. It should be non-zero if successful, while zero for error.
	return C.int(writtenBytes)
}

// FPDF_SaveAsCopy saves the document to a copy.
func (p *PdfiumImplementation) FPDF_SaveAsCopy(request *requests.FPDF_SaveAsCopy) (*responses.FPDF_SaveAsCopy, error) {
	// These methods are basically the same. We switch between
	// FPDF_SaveAsCopy and FPDF_SaveWithVersion in the implementation of FPDF_SaveWithVersion.
	// Don't lock here, FPDF_SaveWithVersion does it for us.
	resp, err := p.FPDF_SaveWithVersion(&requests.FPDF_SaveWithVersion{
		Flags:       request.Flags,
		Document:    request.Document,
		FilePath:    request.FilePath,
		FileWriter:  request.FileWriter,
		FileVersion: 0,
	})

	if err != nil {
		return nil, err
	}

	return &responses.FPDF_SaveAsCopy{
		FileBytes: resp.FileBytes,
		FilePath:  resp.FilePath,
	}, nil
}

// FPDF_SaveWithVersion save the document to a copy, with a specific file version.
func (p *PdfiumImplementation) FPDF_SaveWithVersion(request *requests.FPDF_SaveWithVersion) (*responses.FPDF_SaveWithVersion, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	writer := C.FPDF_FILEWRITE{}
	writer.version = 1

	// Set the Go callback through cgo.
	C.FPDF_FILEWRITE_SET_WRITE_BLOCK(&writer)

	var fileBuf *bytes.Buffer
	var curFile *os.File
	if request.FileWriter != nil {
		currentWriter = request.FileWriter
	} else if request.FilePath != nil {
		newFile, err := os.Create(*request.FilePath)
		if err != nil {
			return nil, err
		}
		currentWriter = newFile
		curFile = newFile
	} else {
		fileBuf = &bytes.Buffer{}
		currentWriter = fileBuf
	}

	defer func() {
		// Always remove writer.
		currentWriter = nil

		if curFile != nil {
			curFile.Close()
		}
	}()

	var success C.int
	if request.FileVersion == 0 {
		success = C.FPDF_SaveAsCopy(documentHandle.handle, &writer, C.ulong(request.Flags))
	} else {
		success = C.FPDF_SaveWithVersion(documentHandle.handle, &writer, C.ulong(request.Flags), C.int(request.FileVersion))
	}

	if int(success) == 0 {
		return nil, errors.New("save of document failed")
	}

	resp := &responses.FPDF_SaveWithVersion{}
	if request.FilePath != nil {
		resp.FilePath = request.FilePath
	}

	if fileBuf != nil {
		pdfContent := fileBuf.Bytes()
		resp.FileBytes = &pdfContent
	}

	return resp, nil
}
