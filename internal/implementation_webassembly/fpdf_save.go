package implementation_webassembly

import (
	"bytes"
	"errors"
	"io"
	"os"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

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

	res, err := p.Module.ExportedFunction("FPDF_FILEWRITE_Create").Call(p.Context)
	if err != nil {
		return nil, err
	}

	fileWriterPointer := res[0]

	// Cleanup the file writer afterwards.
	defer p.Free(fileWriterPointer)

	var fileBuf *bytes.Buffer
	var curFile *os.File
	var currentWriter io.Writer
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

	fileWriterRef := &FileWriterRef{
		Writer:    currentWriter,
		FileWrite: &fileWriterPointer,
	}

	FileWriters.Mutex.Lock()
	FileWriters.Refs[uint32(fileWriterPointer)] = fileWriterRef
	FileWriters.Mutex.Unlock()

	defer func() {
		// Always remove writer.
		currentWriter = nil

		FileWriters.Mutex.Lock()
		delete(FileWriters.Refs, uint32(fileWriterPointer))
		FileWriters.Mutex.Unlock()

		if curFile != nil {
			curFile.Close()
		}
	}()

	var success int32
	if request.FileVersion == 0 {
		res, err = p.Module.ExportedFunction("FPDF_SaveAsCopy").Call(p.Context, *documentHandle.handle, fileWriterPointer, *(*uint64)(unsafe.Pointer(&request.Flags)))
		if err != nil {
			return nil, err
		}
		success = *(*int32)(unsafe.Pointer(&res[0]))
	} else {
		res, err = p.Module.ExportedFunction("FPDF_SaveWithVersion").Call(p.Context, *documentHandle.handle, fileWriterPointer, *(*uint64)(unsafe.Pointer(&request.Flags)), *(*uint64)(unsafe.Pointer(&request.FileVersion)))
		if err != nil {
			return nil, err
		}
		success = *(*int32)(unsafe.Pointer(&res[0]))
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
