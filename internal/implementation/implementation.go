package implementation

import (
	"errors"
	"io"
	"sync"
	"unsafe"

	"github.com/klippa-app/go-pdfium"
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
)

/*
#cgo pkg-config: pdfium
#include "fpdfview.h"
#include <stdlib.h>

extern int go_read_seeker_cb(void *param, unsigned long position, unsigned char *pBuf, unsigned long size);

static inline void FPDF_FILEACCESS_SET_GET_BLOCK(FPDF_FILEACCESS *fs, char *id) {
	fs->m_GetBlock = &go_read_seeker_cb;
	fs->m_Param = id;
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
	fileIdentifier := C.GoString((*C.char)(param))

	// Check if we still have the reader.
	if _, ok := Pdfium.fileReaders[fileIdentifier]; !ok {
		return C.int(0)
	}

	r := Pdfium.fileReaders[fileIdentifier].reader

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

// Pdfium is a container so that we can always only have 1 instance of PDFium
// per process. We need this so that we can guarantee thread safety.
var Pdfium = &mainPdfium{
	mutex:        &sync.Mutex{},
	instanceRefs: map[int]*PdfiumImplementation{},
	documentRefs: map[references.FPDF_DOCUMENT]*DocumentHandle{},
	fileReaders:  map[string]*fileReaderRef{},
}

var isInitialized = false

var libraryConfig C.FPDF_LIBRARY_CONFIG

// InitLibrary loads the actual C++ library.
func InitLibrary(config *pdfium.LibraryConfig) {
	Pdfium.mutex.Lock()
	defer Pdfium.mutex.Unlock()

	// Only initialize when we aren't already.
	if isInitialized {
		return
	}

	if config == nil {
		C.FPDF_InitLibrary()
	} else {
		libraryConfig = C.FPDF_LIBRARY_CONFIG{}

		if config.UserFontPaths != nil && len(config.UserFontPaths) > 0 {
			// Create array of length config.UserFontPaths + 1 for the NULL terminator.
			cArray := C.malloc(C.size_t(len(config.UserFontPaths)+1) * C.size_t(unsafe.Sizeof(uintptr(0))))

			cFonts := (*[1<<30 - 1]*C.char)(cArray)
			for i := range config.UserFontPaths {
				cFonts[i] = C.CString(config.UserFontPaths[i])
			}

			libraryConfig.m_pUserFontPaths = (**C.char)(cArray)
		}

		C.FPDF_InitLibraryWithConfig(&libraryConfig)
	}

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

type fileReaderRef struct {
	reader     io.ReadSeeker
	stringRef  unsafe.Pointer
	fileAccess *C.FPDF_FILEACCESS
}

// Here is the real implementation of Pdfium
type mainPdfium struct {
	// logger is for communication with the plugin.
	logger hclog.Logger

	// mutex will ensure thread safety.
	mutex *sync.Mutex

	// instance keeps track of the opened instances for this process.
	instanceRefs map[int]*PdfiumImplementation

	// documentRefs keeps track of the opened documents for this process.
	// we need this for document lookups and in case of closing the instance
	documentRefs map[references.FPDF_DOCUMENT]*DocumentHandle

	// fileReaders keeps track of the file readers when using FPDF_LoadCustomDocument.
	// We need this to look it up from Go.
	fileReaders map[string]*fileReaderRef
}

func (p *mainPdfium) GetInstance() *PdfiumImplementation {
	newInstance := &PdfiumImplementation{
		logger:               p.logger,
		documentRefs:         map[references.FPDF_DOCUMENT]*DocumentHandle{},
		documentPointers:     map[unsafe.Pointer]references.FPDF_DOCUMENT{},
		pageRefs:             map[references.FPDF_PAGE]*PageHandle{},
		pagePointers:         map[unsafe.Pointer]references.FPDF_PAGE{},
		bookmarkRefs:         map[references.FPDF_BOOKMARK]*BookmarkHandle{},
		destRefs:             map[references.FPDF_DEST]*DestHandle{},
		actionRefs:           map[references.FPDF_ACTION]*ActionHandle{},
		linkRefs:             map[references.FPDF_LINK]*LinkHandle{},
		pageLinkRefs:         map[references.FPDF_PAGELINK]*PageLinkHandle{},
		schHandleRefs:        map[references.FPDF_SCHHANDLE]*SchHandleHandle{},
		bitmapRefs:           map[references.FPDF_BITMAP]*BitmapHandle{},
		textPageRefs:         map[references.FPDF_TEXTPAGE]*TextPageHandle{},
		pageRangeRefs:        map[references.FPDF_PAGERANGE]*PageRangeHandle{},
		pageObjectRefs:       map[references.FPDF_PAGEOBJECT]*PageObjectHandle{},
		clipPathRefs:         map[references.FPDF_CLIPPATH]*ClipPathHandle{},
		formHandleRefs:       map[references.FPDF_FORMHANDLE]*FormHandleHandle{},
		annotationRefs:       map[references.FPDF_ANNOTATION]*AnnotationHandle{},
		xObjectRefs:          map[references.FPDF_XOBJECT]*XObjectHandle{},
		signatureRefs:        map[references.FPDF_SIGNATURE]*SignatureHandle{},
		attachmentRefs:       map[references.FPDF_ATTACHMENT]*AttachmentHandle{},
		javaScriptActionRefs: map[references.FPDF_JAVASCRIPT_ACTION]*JavaScriptActionHandle{},
		searchRefs:           map[references.FPDF_SCHHANDLE]*SearchHandle{},
		pathSegmentRefs:      map[references.FPDF_PATHSEGMENT]*PathSegmentHandle{},
		dataAvailRefs:        map[references.FPDF_AVAIL]*DataAvailHandle{},
		structTreeRefs:       map[references.FPDF_STRUCTTREE]*StructTreeHandle{},
		structElementRefs:    map[references.FPDF_STRUCTELEMENT]*StructElementHandle{},
		pageObjectMarkRefs:   map[references.FPDF_PAGEOBJECTMARK]*PageObjectMarkHandle{},
		fontRefs:             map[references.FPDF_FONT]*FontHandle{},
		glyphPathRefs:        map[references.FPDF_GLYPHPATH]*GlyphPathHandle{},
		fileReaders:          map[string]*fileReaderRef{},
	}

	newInstance.instanceRef = len(p.instanceRefs)
	p.instanceRefs[newInstance.instanceRef] = newInstance

	return newInstance
}

// Here is the real implementation of Pdfium
type PdfiumImplementation struct {
	// logger is for communication with the plugin.
	logger hclog.Logger

	// lookup tables keeps track of the opened handles for this instance.
	// we need this for handle lookups and in case of closing the instance

	documentRefs         map[references.FPDF_DOCUMENT]*DocumentHandle
	documentPointers     map[unsafe.Pointer]references.FPDF_DOCUMENT
	pageRefs             map[references.FPDF_PAGE]*PageHandle
	pagePointers         map[unsafe.Pointer]references.FPDF_PAGE
	bookmarkRefs         map[references.FPDF_BOOKMARK]*BookmarkHandle
	destRefs             map[references.FPDF_DEST]*DestHandle
	actionRefs           map[references.FPDF_ACTION]*ActionHandle
	linkRefs             map[references.FPDF_LINK]*LinkHandle
	pageLinkRefs         map[references.FPDF_PAGELINK]*PageLinkHandle
	schHandleRefs        map[references.FPDF_SCHHANDLE]*SchHandleHandle
	textPageRefs         map[references.FPDF_TEXTPAGE]*TextPageHandle
	pageRangeRefs        map[references.FPDF_PAGERANGE]*PageRangeHandle
	pageObjectRefs       map[references.FPDF_PAGEOBJECT]*PageObjectHandle
	clipPathRefs         map[references.FPDF_CLIPPATH]*ClipPathHandle
	formHandleRefs       map[references.FPDF_FORMHANDLE]*FormHandleHandle
	bitmapRefs           map[references.FPDF_BITMAP]*BitmapHandle
	annotationRefs       map[references.FPDF_ANNOTATION]*AnnotationHandle
	xObjectRefs          map[references.FPDF_XOBJECT]*XObjectHandle
	signatureRefs        map[references.FPDF_SIGNATURE]*SignatureHandle
	attachmentRefs       map[references.FPDF_ATTACHMENT]*AttachmentHandle
	javaScriptActionRefs map[references.FPDF_JAVASCRIPT_ACTION]*JavaScriptActionHandle
	searchRefs           map[references.FPDF_SCHHANDLE]*SearchHandle
	pathSegmentRefs      map[references.FPDF_PATHSEGMENT]*PathSegmentHandle
	dataAvailRefs        map[references.FPDF_AVAIL]*DataAvailHandle
	structTreeRefs       map[references.FPDF_STRUCTTREE]*StructTreeHandle
	structElementRefs    map[references.FPDF_STRUCTELEMENT]*StructElementHandle
	pageObjectMarkRefs   map[references.FPDF_PAGEOBJECTMARK]*PageObjectMarkHandle
	fontRefs             map[references.FPDF_FONT]*FontHandle
	glyphPathRefs        map[references.FPDF_GLYPHPATH]*GlyphPathHandle
	fileReaders          map[string]*fileReaderRef

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
		defer C.free(unsafe.Pointer(cPassword))
	}

	nativeDoc := &DocumentHandle{
		pageRefs:             map[references.FPDF_PAGE]*PageHandle{},
		bookmarkRefs:         map[references.FPDF_BOOKMARK]*BookmarkHandle{},
		destRefs:             map[references.FPDF_DEST]*DestHandle{},
		pageLinkRefs:         map[references.FPDF_PAGELINK]*PageLinkHandle{},
		schHandleRefs:        map[references.FPDF_SCHHANDLE]*SchHandleHandle{},
		textPageRefs:         map[references.FPDF_TEXTPAGE]*TextPageHandle{},
		pageRangeRefs:        map[references.FPDF_PAGERANGE]*PageRangeHandle{},
		formHandleRefs:       map[references.FPDF_FORMHANDLE]*FormHandleHandle{},
		signatureRefs:        map[references.FPDF_SIGNATURE]*SignatureHandle{},
		attachmentRefs:       map[references.FPDF_ATTACHMENT]*AttachmentHandle{},
		javaScriptActionRefs: map[references.FPDF_JAVASCRIPT_ACTION]*JavaScriptActionHandle{},
		searchRefs:           map[references.FPDF_SCHHANDLE]*SearchHandle{},
		structTreeRefs:       map[references.FPDF_STRUCTTREE]*StructTreeHandle{},
		structElementRefs:    map[references.FPDF_STRUCTELEMENT]*StructElementHandle{},
	}
	var doc C.FPDF_DOCUMENT

	if request.File != nil {
		fileData := *request.File

		// If larger than INT_MAX, use FPDF_LoadMemDocument64
		if len(fileData) > 2147483647 {
			doc = C.FPDF_LoadMemDocument64(
				unsafe.Pointer(&(fileData[0])),
				C.size_t(len(fileData)),
				cPassword)
		} else {
			doc = C.FPDF_LoadMemDocument(
				unsafe.Pointer(&(fileData[0])),
				C.int(len(fileData)),
				cPassword)
		}
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

		// Create a PDFium file access struct.
		readerStruct := C.FPDF_FILEACCESS{}
		readerStruct.m_FileLen = C.ulong(request.FileReaderSize)

		readerRef := uuid.New()
		readerRefString := readerRef.String()
		cReaderRef := C.CString(readerRefString)

		// Set the Go callback through cgo.
		C.FPDF_FILEACCESS_SET_GET_BLOCK(&readerStruct, cReaderRef)

		fileReaderRef := &fileReaderRef{
			stringRef:  unsafe.Pointer(cReaderRef),
			reader:     request.FileReader,
			fileAccess: &readerStruct,
		}

		Pdfium.fileReaders[readerRef.String()] = fileReaderRef
		nativeDoc.fileHandleRef = &readerRefString

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

		// Cleanup when file loading didn't work.
		if nativeDoc.fileHandleRef != nil {
			C.free(Pdfium.fileReaders[*nativeDoc.fileHandleRef].stringRef)
			delete(Pdfium.fileReaders, *nativeDoc.fileHandleRef)
		}

		return nil, pdfiumError
	}

	documentRef := uuid.New()
	nativeDoc.handle = doc
	nativeDoc.data = request.File
	nativeDoc.nativeRef = references.FPDF_DOCUMENT(documentRef.String())
	Pdfium.documentRefs[nativeDoc.nativeRef] = nativeDoc
	p.documentRefs[nativeDoc.nativeRef] = nativeDoc

	return &responses.OpenDocument{
		Document: nativeDoc.nativeRef,
	}, nil
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

	for i := range p.pageRefs {
		// Already closed by the document close.
		delete(p.pageRefs, p.pageRefs[i].nativeRef)
	}

	// Remove refs, they don't have a close method.
	for i := range p.bookmarkRefs {
		delete(p.bookmarkRefs, i)
	}

	for i := range p.destRefs {
		delete(p.destRefs, i)
	}

	for i := range p.actionRefs {
		delete(p.actionRefs, i)
	}

	for i := range p.linkRefs {
		delete(p.linkRefs, i)
	}

	for i := range p.pageLinkRefs {
		delete(p.pageLinkRefs, i)
	}

	for i := range p.schHandleRefs {
		delete(p.schHandleRefs, i)
	}

	for i := range p.bitmapRefs {
		delete(p.bitmapRefs, i)
	}

	for i := range p.textPageRefs {
		delete(p.textPageRefs, i)
	}

	for i := range p.pageRangeRefs {
		delete(p.pageRangeRefs, i)
	}

	for i := range p.pageObjectRefs {
		delete(p.pageObjectRefs, i)
	}

	for i := range p.clipPathRefs {
		delete(p.clipPathRefs, i)
	}

	for i := range p.formHandleRefs {
		delete(p.formHandleRefs, i)
	}

	for i := range p.annotationRefs {
		delete(p.annotationRefs, i)
	}

	for i := range p.xObjectRefs {
		delete(p.xObjectRefs, i)
	}

	for i := range p.signatureRefs {
		delete(p.signatureRefs, i)
	}

	for i := range p.attachmentRefs {
		delete(p.attachmentRefs, i)
	}

	for i := range p.javaScriptActionRefs {
		delete(p.javaScriptActionRefs, i)
	}

	for i := range p.searchRefs {
		delete(p.searchRefs, i)
	}

	for i := range p.pathSegmentRefs {
		delete(p.pathSegmentRefs, i)
	}

	for i := range p.dataAvailRefs {
		delete(p.dataAvailRefs, i)
	}

	for i := range p.structTreeRefs {
		delete(p.structTreeRefs, i)
	}

	for i := range p.structElementRefs {
		delete(p.structElementRefs, i)
	}

	for i := range p.pageObjectMarkRefs {
		delete(p.pageObjectMarkRefs, i)
	}

	for i := range p.fontRefs {
		delete(p.fontRefs, i)
	}

	for i := range p.glyphPathRefs {
		delete(p.glyphPathRefs, i)
	}

	for i := range p.fileReaders {
		// Cleanup file handle.
		Pdfium.fileReaders[i].fileAccess = nil
		C.free(Pdfium.fileReaders[i].stringRef)
		delete(Pdfium.fileReaders, i)
		delete(p.fileReaders, i)
	}

	delete(Pdfium.instanceRefs, p.instanceRef)

	return nil
}

func (p *PdfiumImplementation) getDocumentHandle(documentRef references.FPDF_DOCUMENT) (*DocumentHandle, error) {
	if documentRef == "" {
		return nil, errors.New("document not given")
	}

	if val, ok := p.documentRefs[documentRef]; ok {
		return val, nil
	}

	return nil, errors.New("could not find document handle, perhaps the doc was already closed or you tried to share documents between instances")
}

func (d *PdfiumImplementation) getPageHandle(pageRef references.FPDF_PAGE) (*PageHandle, error) {
	if pageRef == "" {
		return nil, errors.New("page not given")
	}

	if val, ok := d.pageRefs[pageRef]; ok {
		return val, nil
	}

	return nil, errors.New("could not find page handle, perhaps the page was already closed or you tried to share pages between instances or documents")
}

func (d *PdfiumImplementation) getBookmarkHandle(bookmarkRef references.FPDF_BOOKMARK) (*BookmarkHandle, error) {
	if bookmarkRef == "" {
		return nil, errors.New("bookmark not given")
	}

	if val, ok := d.bookmarkRefs[bookmarkRef]; ok {
		return val, nil
	}

	return nil, errors.New("could not find bookmark handle, perhaps the bookmark was already closed or you tried to share bookmarks between instances or documents")
}

func (d *PdfiumImplementation) getDestHandle(destRef references.FPDF_DEST) (*DestHandle, error) {
	if destRef == "" {
		return nil, errors.New("dest not given")
	}

	if val, ok := d.destRefs[destRef]; ok {
		return val, nil
	}

	return nil, errors.New("could not find dest handle, perhaps the dest was already closed or you tried to share dests between instances or documents")
}

func (d *PdfiumImplementation) getActionHandle(actionRef references.FPDF_ACTION) (*ActionHandle, error) {
	if actionRef == "" {
		return nil, errors.New("action not given")
	}

	if val, ok := d.actionRefs[actionRef]; ok {
		return val, nil
	}

	return nil, errors.New("could not find action handle, perhaps the action was already closed or you tried to share actions between instances or documents")
}

func (d *PdfiumImplementation) getLinkHandle(linkRef references.FPDF_LINK) (*LinkHandle, error) {
	if linkRef == "" {
		return nil, errors.New("link not given")
	}

	if val, ok := d.linkRefs[linkRef]; ok {
		return val, nil
	}

	return nil, errors.New("could not find link handle, perhaps the link was already closed or you tried to share links between instances or documents")
}

func (d *PdfiumImplementation) getXObjectHandle(xObject references.FPDF_XOBJECT) (*XObjectHandle, error) {
	if xObject == "" {
		return nil, errors.New("xObject not given")
	}

	if val, ok := d.xObjectRefs[xObject]; ok {
		return val, nil
	}

	return nil, errors.New("could not find xObject handle, perhaps the xObject was already closed or you tried to share xObjects between instances or documents")
}

func (d *PdfiumImplementation) getSignatureHandle(signature references.FPDF_SIGNATURE) (*SignatureHandle, error) {
	if signature == "" {
		return nil, errors.New("signature not given")
	}

	if val, ok := d.signatureRefs[signature]; ok {
		return val, nil
	}

	return nil, errors.New("could not find signature handle, perhaps the signature was already closed or you tried to share signatures between instances or documents")
}

func (d *PdfiumImplementation) getAttachmentHandle(attachment references.FPDF_ATTACHMENT) (*AttachmentHandle, error) {
	if attachment == "" {
		return nil, errors.New("attachment not given")
	}

	if val, ok := d.attachmentRefs[attachment]; ok {
		return val, nil
	}

	return nil, errors.New("could not find attachment handle, perhaps the attachment was already closed or you tried to share attachments between instances or documents")
}

func (d *PdfiumImplementation) getJavaScriptActionHandle(javaScriptAction references.FPDF_JAVASCRIPT_ACTION) (*JavaScriptActionHandle, error) {
	if javaScriptAction == "" {
		return nil, errors.New("javaScriptAction not given")
	}

	if val, ok := d.javaScriptActionRefs[javaScriptAction]; ok {
		return val, nil
	}

	return nil, errors.New("could not find javaScriptAction handle, perhaps the javaScriptAction was already closed or you tried to share javaScriptActions between instances or documents")
}

func (p *PdfiumImplementation) getTextPageHandle(textPage references.FPDF_TEXTPAGE) (*TextPageHandle, error) {
	if textPage == "" {
		return nil, errors.New("textPage not given")
	}

	if val, ok := p.textPageRefs[textPage]; ok {
		return val, nil
	}

	return nil, errors.New("could not find textPage handle, perhaps the textPage was already closed or you tried to share textPages between instances")
}

func (p *PdfiumImplementation) getSearchHandle(search references.FPDF_SCHHANDLE) (*SearchHandle, error) {
	if search == "" {
		return nil, errors.New("search not given")
	}

	if val, ok := p.searchRefs[search]; ok {
		return val, nil
	}

	return nil, errors.New("could not find search handle, perhaps the search was already closed or you tried to share searchs between instances")
}

func (p *PdfiumImplementation) getPageLinkHandle(pageLink references.FPDF_PAGELINK) (*PageLinkHandle, error) {
	if pageLink == "" {
		return nil, errors.New("pageLink not given")
	}

	if val, ok := p.pageLinkRefs[pageLink]; ok {
		return val, nil
	}

	return nil, errors.New("could not find pageLink handle, perhaps the pageLink was already closed or you tried to share pageLinks between instances")
}

func (p *PdfiumImplementation) getBitmapHandle(bitmap references.FPDF_BITMAP) (*BitmapHandle, error) {
	if bitmap == "" {
		return nil, errors.New("bitmap not given")
	}

	if val, ok := p.bitmapRefs[bitmap]; ok {
		return val, nil
	}

	return nil, errors.New("could not find bitmap handle, perhaps the bitmap was already closed or you tried to share bitmaps between instances")
}

func (p *PdfiumImplementation) getPageRangeHandle(pageRange references.FPDF_PAGERANGE) (*PageRangeHandle, error) {
	if pageRange == "" {
		return nil, errors.New("pageRange not given")
	}

	if val, ok := p.pageRangeRefs[pageRange]; ok {
		return val, nil
	}

	return nil, errors.New("could not find pageRange handle, perhaps the pageRange was already closed or you tried to share pageRanges between instances")
}

func (p *PdfiumImplementation) getPageObjectHandle(pageObject references.FPDF_PAGEOBJECT) (*PageObjectHandle, error) {
	if pageObject == "" {
		return nil, errors.New("pageObject not given")
	}

	if val, ok := p.pageObjectRefs[pageObject]; ok {
		return val, nil
	}

	return nil, errors.New("could not find pageObject handle, perhaps the pageObject was already closed or you tried to share pageObjects between instances")
}

func (p *PdfiumImplementation) getClipPathHandle(clipPath references.FPDF_CLIPPATH) (*ClipPathHandle, error) {
	if clipPath == "" {
		return nil, errors.New("clipPath not given")
	}

	if val, ok := p.clipPathRefs[clipPath]; ok {
		return val, nil
	}

	return nil, errors.New("could not find clipPath handle, perhaps the clipPath was already closed or you tried to share clipPaths between instances")
}

func (p *PdfiumImplementation) getDataAvailHandle(dataAvail references.FPDF_AVAIL) (*DataAvailHandle, error) {
	if dataAvail == "" {
		return nil, errors.New("dataAvail not given")
	}

	if val, ok := p.dataAvailRefs[dataAvail]; ok {
		return val, nil
	}

	return nil, errors.New("could not find dataAvail handle, perhaps the dataAvail was already closed or you tried to share dataAvails between instances")
}

func (p *PdfiumImplementation) getStructTreeHandle(structTree references.FPDF_STRUCTTREE) (*StructTreeHandle, error) {
	if structTree == "" {
		return nil, errors.New("structTree not given")
	}

	if val, ok := p.structTreeRefs[structTree]; ok {
		return val, nil
	}

	return nil, errors.New("could not find structTree handle, perhaps the structTree was already closed or you tried to share structTrees between instances")
}

func (p *PdfiumImplementation) getStructElementHandle(structElement references.FPDF_STRUCTELEMENT) (*StructElementHandle, error) {
	if structElement == "" {
		return nil, errors.New("structElement not given")
	}

	if val, ok := p.structElementRefs[structElement]; ok {
		return val, nil
	}

	return nil, errors.New("could not find structElement handle, perhaps the structElement was already closed or you tried to share structElements between instances")
}

func (p *PdfiumImplementation) getPageObjectMarkHandle(pageObjectMark references.FPDF_PAGEOBJECTMARK) (*PageObjectMarkHandle, error) {
	if pageObjectMark == "" {
		return nil, errors.New("pageObjectMark not given")
	}

	if val, ok := p.pageObjectMarkRefs[pageObjectMark]; ok {
		return val, nil
	}

	return nil, errors.New("could not find pageObjectMark handle, perhaps the pageObjectMark was already closed or you tried to share pageObjectMarks between instances")
}

func (p *PdfiumImplementation) getPathSegmentHandle(pathSegment references.FPDF_PATHSEGMENT) (*PathSegmentHandle, error) {
	if pathSegment == "" {
		return nil, errors.New("pathSegment not given")
	}

	if val, ok := p.pathSegmentRefs[pathSegment]; ok {
		return val, nil
	}

	return nil, errors.New("could not find pathSegment handle, perhaps the pathSegment was already closed or you tried to share pathSegments between instances")
}

func (p *PdfiumImplementation) getFontHandle(font references.FPDF_FONT) (*FontHandle, error) {
	if font == "" {
		return nil, errors.New("font not given")
	}

	if val, ok := p.fontRefs[font]; ok {
		return val, nil
	}

	return nil, errors.New("could not find font handle, perhaps the font was already closed or you tried to share fonts between instances")
}

func (p *PdfiumImplementation) getGlyphPathHandle(glyphPath references.FPDF_GLYPHPATH) (*GlyphPathHandle, error) {
	if glyphPath == "" {
		return nil, errors.New("glyphPath not given")
	}

	if val, ok := p.glyphPathRefs[glyphPath]; ok {
		return val, nil
	}

	return nil, errors.New("could not find glyphPath handle, perhaps the glyphPath was already closed or you tried to share glyphPaths between instances")
}

func (p *PdfiumImplementation) getAnnotationHandle(annotation references.FPDF_ANNOTATION) (*AnnotationHandle, error) {
	if annotation == "" {
		return nil, errors.New("annotation not given")
	}

	if val, ok := p.annotationRefs[annotation]; ok {
		return val, nil
	}

	return nil, errors.New("could not find annotation handle, perhaps the annotation was already closed or you tried to share annotations between instances")
}

func (p *PdfiumImplementation) getFormHandleHandle(formHandle references.FPDF_FORMHANDLE) (*FormHandleHandle, error) {
	if formHandle == "" {
		return nil, errors.New("formHandle not given")
	}

	if val, ok := p.formHandleRefs[formHandle]; ok {
		return val, nil
	}

	return nil, errors.New("could not find formHandle handle, perhaps the formHandle was already closed or you tried to share formHandles between instances")
}
