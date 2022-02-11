package implementation

/*
#cgo pkg-config: pdfium
#include "fpdf_dataavail.h"
#include <stdlib.h>

extern int go_dataavail_is_data_avail_cb(struct _FX_FILEAVAIL *me, size_t offset, size_t size);
extern void go_dataavail_add_segment_cb(struct _FX_DOWNLOADHINTS *me, size_t offset, size_t size);

extern int go_read_seeker_cb(void *param, unsigned long position, unsigned char *pBuf, unsigned long size);

static inline void FPDF_FX_FILEAVAIL_CB(FX_FILEAVAIL *p) {
	p->IsDataAvail = &go_dataavail_is_data_avail_cb;
}

static inline void FPDF_FX_DOWNLOADHINTS_CB(FX_DOWNLOADHINTS *h) {
	h->AddSegment = &go_dataavail_add_segment_cb;
}

static inline void FPDF_FILEAVAIL_FILEACCESS_SET_GET_BLOCK(FPDF_FILEACCESS *fs, char *id) {
	fs->m_GetBlock = &go_read_seeker_cb;
	fs->m_Param = id;
}
*/
import "C"
import (
	"errors"
	"log"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/google/uuid"
)

// go_progressive_render_pause_cb is the Go implementation of FX_FILEAVAIL::IsDataAvail.
// It is exported through cgo so that we can use the reference to it and set
// it on FX_FILEAVAIL structs.
//export go_dataavail_is_data_avail_cb
func go_dataavail_is_data_avail_cb(me *C.FX_FILEAVAIL, offset C.size_t, size C.size_t) C.FPDF_BOOL {
	log.Println("Checking if data is available")
	return C.FPDF_BOOL(0)
}

// go_dataavail_add_segment_cb is the Go implementation of FX_DOWNLOADHINTS::AddSegment.
// It is exported through cgo so that we can use the reference to it and set
// it on FX_FILEAVAIL structs.
//export go_dataavail_add_segment_cb
func go_dataavail_add_segment_cb(me *C.FX_DOWNLOADHINTS, offset C.size_t, size C.size_t) {
	log.Println("please add segment")
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

	var hints *C.FX_DOWNLOADHINTS
	if request.AddSegmentCallback != nil {
		hints = &C.FX_DOWNLOADHINTS{}
		C.FPDF_FX_DOWNLOADHINTS_CB(hints)
	}

	// Create a PDFium file access struct.
	readerStruct := C.FPDF_FILEACCESS{}
	readerStruct.m_FileLen = C.ulong(request.Size)

	readerRef := uuid.New()
	readerRefString := readerRef.String()
	cReaderRef := C.CString(readerRefString)

	// Set the Go callback through cgo.
	C.FPDF_FILEAVAIL_FILEACCESS_SET_GET_BLOCK(&readerStruct, cReaderRef)

	fileReaderRef := &fileReaderRef{
		stringRef:  unsafe.Pointer(cReaderRef),
		reader:     request.Reader,
		fileAccess: &readerStruct,
	}

	Pdfium.fileReaders[readerRefString] = fileReaderRef

	availStruct := C.FX_FILEAVAIL{}

	// Set the Go callback through cgo.
	C.FPDF_FX_FILEAVAIL_CB(&availStruct)

	dataAvail := C.FPDFAvail_Create(&availStruct, &readerStruct)
	dataAvailHandle := p.registerDataAvail(dataAvail, readerRefString, availStruct, hints)

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

	C.FPDFAvail_Destroy(dataAvailHandler.handle)

	delete(p.dataAvailRefs, dataAvailHandler.nativeRef)
	C.free(Pdfium.fileReaders[dataAvailHandler.fileHandleRef].stringRef)
	delete(Pdfium.fileReaders, dataAvailHandler.fileHandleRef)

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

	isDocAvail := C.FPDFAvail_IsDocAvail(dataAvailHandler.handle, dataAvailHandler.hints)

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

	var cPassword *C.char
	if request.Password != nil {
		cPassword = C.CString(*request.Password)
		defer C.free(unsafe.Pointer(cPassword))
	}

	doc := C.FPDFAvail_GetDocument(dataAvailHandler.handle, cPassword)

	documentHandle := &DocumentHandle{}
	documentHandle.handle = doc
	documentRef := uuid.New()
	documentHandle.nativeRef = references.FPDF_DOCUMENT(documentRef.String())
	Pdfium.documentRefs[documentHandle.nativeRef] = documentHandle
	p.documentRefs[documentHandle.nativeRef] = documentHandle

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

	firstPageNum := C.FPDFAvail_GetFirstPageNum(documentHandle.handle)

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

	isPageAvail := C.FPDFAvail_IsPageAvail(dataAvailHandler.handle, C.int(request.PageIndex), dataAvailHandler.hints)

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

	isFormAvail := C.FPDFAvail_IsFormAvail(dataAvailHandler.handle, dataAvailHandler.hints)

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

	isLinearized := C.FPDFAvail_IsLinearized(dataAvailHandler.handle)

	return &responses.FPDFAvail_IsLinearized{
		IsLinearized: enums.PDF_FILEAVAIL_LINEARIZATION(isLinearized),
	}, nil
}
