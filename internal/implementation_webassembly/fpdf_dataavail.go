package implementation_webassembly

import (
	"errors"
	"sync"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

var FileAvailables = struct {
	Refs  map[uint32]*DataAvailHandle
	Mutex *sync.Mutex
}{
	Refs:  map[uint32]*DataAvailHandle{},
	Mutex: &sync.Mutex{},
}

var FileHints = struct {
	Refs  map[uint32]*DataAvailHandle
	Mutex *sync.Mutex
}{
	Refs:  map[uint32]*DataAvailHandle{},
	Mutex: &sync.Mutex{},
}

// FPDFAvail_Create creates a document availability provider.
// FPDFAvail_Destroy() must be called when done with the availability provider.
func (p *PdfiumImplementation) FPDFAvail_Create(request *requests.FPDFAvail_Create) (*responses.FPDFAvail_Create, error) {
	p.Lock()
	defer p.Unlock()

	if request.IsDataAvailableCallback == nil {
		return nil, errors.New("IsDataAvailableCallback can't be nil")
	}

	if request.Reader == nil {
		return nil, errors.New("Reader can't be nil")
	}

	if request.Size == 0 {
		return nil, errors.New("Size should be set")
	}

	fileReaderPointer, fileReaderIndex, err := p.CreateFileAccessReader(request.Size, request.Reader)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FX_FILEAVAIL_Create").Call(p.Context)
	if err != nil {
		return nil, err
	}

	fXFileAvail := res[0]

	res, err = p.Module.ExportedFunction("FPDFAvail_Create").Call(p.Context, fXFileAvail, *fileReaderPointer)
	if err != nil {
		return nil, err
	}

	fPDFAvail := res[0]

	hints := uint64(0)
	if request.AddSegmentCallback != nil {
		res, err = p.Module.ExportedFunction("FX_DOWNLOADHINTS_Create").Call(p.Context)
		if err != nil {
			return nil, err
		}

		hints = res[0]
	}

	dataAvailHandle := p.registerDataAvail(&fPDFAvail, &fXFileAvail, &hints, fileReaderIndex, request.IsDataAvailableCallback, request.AddSegmentCallback)

	FileAvailables.Mutex.Lock()
	FileAvailables.Refs[uint32(fXFileAvail)] = dataAvailHandle
	FileAvailables.Mutex.Unlock()

	if hints != 0 {
		FileHints.Mutex.Lock()
		FileHints.Refs[uint32(hints)] = dataAvailHandle
		FileHints.Mutex.Unlock()
	}

	return &responses.FPDFAvail_Create{
		AvailabilityProvider: dataAvailHandle.nativeRef,
	}, nil
}

// FPDFAvail_Destroy destroys the given document availability provider.
func (p *PdfiumImplementation) FPDFAvail_Destroy(request *requests.FPDFAvail_Destroy) (*responses.FPDFAvail_Destroy, error) {
	p.Lock()
	defer p.Unlock()

	dataAvailHandler, err := p.getDataAvailHandle(request.AvailabilityProvider)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDFAvail_Destroy").Call(p.Context, *dataAvailHandler.handle)
	if err != nil {
		return nil, err
	}

	p.Free(*p.fileReaders[*dataAvailHandler.reader].FileAccess)
	p.Free(*dataAvailHandler.fileAvail)
	p.Free(*dataAvailHandler.handle)
	delete(p.fileReaders, *dataAvailHandler.reader)
	FileReaders.Mutex.Lock()
	delete(FileReaders.Refs, *dataAvailHandler.reader)
	FileReaders.Mutex.Unlock()
	delete(p.dataAvailRefs, dataAvailHandler.nativeRef)

	FileAvailables.Mutex.Lock()
	delete(FileAvailables.Refs, uint32(*dataAvailHandler.fileAvail))
	FileAvailables.Mutex.Unlock()

	if dataAvailHandler.hints != nil {
		p.Free(*dataAvailHandler.hints)
		FileHints.Mutex.Lock()
		delete(FileHints.Refs, uint32(*dataAvailHandler.hints))
		FileHints.Mutex.Unlock()
	}

	return &responses.FPDFAvail_Destroy{}, nil
}

// FPDFAvail_IsDocAvail checks if the document is ready for loading, if not, gets download hints.
// Applications should call this function whenever new data arrives, and process
// all the generated download hints, if any, until the function returns
// enums.PDF_FILEAVAIL_DATA_ERROR or enums.PDF_FILEAVAIL_DATA_AVAIL.
// if hints is nil, the function just check current document availability.
//
// Once all data is available, call FPDFAvail_GetDocument() to get a document
// handle.
func (p *PdfiumImplementation) FPDFAvail_IsDocAvail(request *requests.FPDFAvail_IsDocAvail) (*responses.FPDFAvail_IsDocAvail, error) {
	p.Lock()
	defer p.Unlock()

	dataAvailHandler, err := p.getDataAvailHandle(request.AvailabilityProvider)
	if err != nil {
		return nil, err
	}

	hints := uint64(0)
	if dataAvailHandler.hints != nil {
		hints = *dataAvailHandler.hints
	}

	res, err := p.Module.ExportedFunction("FPDFAvail_IsDocAvail").Call(p.Context, *dataAvailHandler.handle, hints)
	if err != nil {
		return nil, err
	}

	isDocAvail := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFAvail_IsDocAvail{
		IsDocAvail: enums.PDF_FILEAVAIL_DATA(isDocAvail),
	}, nil
}

// FPDFAvail_GetDocument returns the document from the availability provider.
// When FPDFAvail_IsDocAvail() returns TRUE, call FPDFAvail_GetDocument() to
// retrieve the document handle.
func (p *PdfiumImplementation) FPDFAvail_GetDocument(request *requests.FPDFAvail_GetDocument) (*responses.FPDFAvail_GetDocument, error) {
	p.Lock()
	defer p.Unlock()

	dataAvailHandler, err := p.getDataAvailHandle(request.AvailabilityProvider)
	if err != nil {
		return nil, err
	}

	var cPassword uint64
	if request.Password != nil {
		cPasswordPointer, err := p.CString(*request.Password)
		if err != nil {
			return nil, err
		}

		defer cPasswordPointer.Free()

		cPassword = cPasswordPointer.Pointer
	}

	res, err := p.Module.ExportedFunction("FPDFAvail_GetDocument").Call(p.Context, *dataAvailHandler.handle, cPassword)
	if err != nil {
		return nil, err
	}

	doc := res[0]
	documentHandle := p.registerDocument(&doc)

	return &responses.FPDFAvail_GetDocument{
		Document: documentHandle.nativeRef,
	}, nil
}

// FPDFAvail_GetFirstPageNum returns the page number for the first available page in a linearized PDF.
// For most linearized PDFs, the first available page will be the first page,
// however, some PDFs might make another page the first available page.
// For non-linearized PDFs, this function will always return zero.
func (p *PdfiumImplementation) FPDFAvail_GetFirstPageNum(request *requests.FPDFAvail_GetFirstPageNum) (*responses.FPDFAvail_GetFirstPageNum, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFAvail_GetFirstPageNum").Call(p.Context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	firstPageNum := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFAvail_GetFirstPageNum{
		FirstPageNum: int(firstPageNum),
	}, nil
}

// FPDFAvail_IsPageAvail checks if the given page index is ready for loading, if not, it will
// call the hints to fetch more data.
func (p *PdfiumImplementation) FPDFAvail_IsPageAvail(request *requests.FPDFAvail_IsPageAvail) (*responses.FPDFAvail_IsPageAvail, error) {
	p.Lock()
	defer p.Unlock()

	dataAvailHandler, err := p.getDataAvailHandle(request.AvailabilityProvider)
	if err != nil {
		return nil, err
	}

	hints := uint64(0)
	if dataAvailHandler.hints != nil {
		hints = *dataAvailHandler.hints
	}

	res, err := p.Module.ExportedFunction("FPDFAvail_IsPageAvail").Call(p.Context, *dataAvailHandler.handle, *(*uint64)(unsafe.Pointer(&request.PageIndex)), hints)
	if err != nil {
		return nil, err
	}

	isPageAvail := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFAvail_IsPageAvail{
		IsPageAvail: enums.PDF_FILEAVAIL_DATA(isPageAvail),
	}, nil
}

// FPDFAvail_IsFormAvail
// This function can be called only after FPDFAvail_GetDocument() is called.
// Applications should call this function whenever new data arrives and process
// all the generated download hints, if any, until this function returns
// enums.PDF_FILEAVAIL_DATA_ERROR or enums.PDF_FILEAVAIL_DATA_AVAIL. Applications can then perform page
// loading.
// if hints is nil, the function just check current availability of
// specified page.
func (p *PdfiumImplementation) FPDFAvail_IsFormAvail(request *requests.FPDFAvail_IsFormAvail) (*responses.FPDFAvail_IsFormAvail, error) {
	p.Lock()
	defer p.Unlock()

	dataAvailHandler, err := p.getDataAvailHandle(request.AvailabilityProvider)
	if err != nil {
		return nil, err
	}

	hints := uint64(0)
	if dataAvailHandler.hints != nil {
		hints = *dataAvailHandler.hints
	}

	res, err := p.Module.ExportedFunction("FPDFAvail_IsFormAvail").Call(p.Context, *dataAvailHandler.handle, hints)
	if err != nil {
		return nil, err
	}

	isFormAvail := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFAvail_IsFormAvail{
		IsFormAvail: enums.PDF_FILEAVAIL_FORM(isFormAvail),
	}, nil
}

// FPDFAvail_IsLinearized Check whether a document is a linearized PDF.
// FPDFAvail_IsLinearized() will return enums.PDF_FILEAVAIL_LINEARIZED or enums.PDF_FILEAVAIL_NOT_LINEARIZED
// when we have 1k  of data. If the files size less than 1k, it returns
// enums.PDF_FILEAVAIL_LINEARIZATION_UNKNOWN as there is insufficient information to determine
// if the PDF is linearlized.
func (p *PdfiumImplementation) FPDFAvail_IsLinearized(request *requests.FPDFAvail_IsLinearized) (*responses.FPDFAvail_IsLinearized, error) {
	p.Lock()
	defer p.Unlock()

	dataAvailHandler, err := p.getDataAvailHandle(request.AvailabilityProvider)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFAvail_IsLinearized").Call(p.Context, *dataAvailHandler.handle)
	if err != nil {
		return nil, err
	}

	isLinearized := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFAvail_IsLinearized{
		IsLinearized: enums.PDF_FILEAVAIL_LINEARIZATION(isLinearized),
	}, nil
}
