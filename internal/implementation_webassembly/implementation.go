package implementation_webassembly

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"sync"

	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/google/uuid"
	"github.com/tetratelabs/wazero/api"
)

func GetInstance(ctx context.Context, functions map[string]api.Function, module api.Module) *PdfiumImplementation {
	newInstance := &PdfiumImplementation{
		mutex:                      &sync.Mutex{},
		Context:                    ctx,
		Functions:                  functions,
		Module:                     module,
		documentRefs:               map[references.FPDF_DOCUMENT]*DocumentHandle{},
		pageRefs:                   map[references.FPDF_PAGE]*PageHandle{},
		bookmarkRefs:               map[references.FPDF_BOOKMARK]*BookmarkHandle{},
		destRefs:                   map[references.FPDF_DEST]*DestHandle{},
		actionRefs:                 map[references.FPDF_ACTION]*ActionHandle{},
		linkRefs:                   map[references.FPDF_LINK]*LinkHandle{},
		pageLinkRefs:               map[references.FPDF_PAGELINK]*PageLinkHandle{},
		schHandleRefs:              map[references.FPDF_SCHHANDLE]*SchHandleHandle{},
		bitmapRefs:                 map[references.FPDF_BITMAP]*BitmapHandle{},
		textPageRefs:               map[references.FPDF_TEXTPAGE]*TextPageHandle{},
		pageRangeRefs:              map[references.FPDF_PAGERANGE]*PageRangeHandle{},
		pageObjectRefs:             map[references.FPDF_PAGEOBJECT]*PageObjectHandle{},
		clipPathRefs:               map[references.FPDF_CLIPPATH]*ClipPathHandle{},
		formHandleRefs:             map[references.FPDF_FORMHANDLE]*FormHandleHandle{},
		annotationRefs:             map[references.FPDF_ANNOTATION]*AnnotationHandle{},
		xObjectRefs:                map[references.FPDF_XOBJECT]*XObjectHandle{},
		signatureRefs:              map[references.FPDF_SIGNATURE]*SignatureHandle{},
		attachmentRefs:             map[references.FPDF_ATTACHMENT]*AttachmentHandle{},
		javaScriptActionRefs:       map[references.FPDF_JAVASCRIPT_ACTION]*JavaScriptActionHandle{},
		searchRefs:                 map[references.FPDF_SCHHANDLE]*SearchHandle{},
		pathSegmentRefs:            map[references.FPDF_PATHSEGMENT]*PathSegmentHandle{},
		dataAvailRefs:              map[references.FPDF_AVAIL]*DataAvailHandle{},
		structTreeRefs:             map[references.FPDF_STRUCTTREE]*StructTreeHandle{},
		structElementRefs:          map[references.FPDF_STRUCTELEMENT]*StructElementHandle{},
		structElementAttributeRefs: map[references.FPDF_STRUCTELEMENT_ATTR]*StructElementAttributeHandle{},
		pageObjectMarkRefs:         map[references.FPDF_PAGEOBJECTMARK]*PageObjectMarkHandle{},
		fontRefs:                   map[references.FPDF_FONT]*FontHandle{},
		glyphPathRefs:              map[references.FPDF_GLYPHPATH]*GlyphPathHandle{},
		fileReaders:                map[uint32]*FileReaderRef{},
	}

	return newInstance
}

type FunctionWrapper struct {
	function api.Function
	mutex    *sync.Mutex
}

// Definition implements the same method as documented on api.FunctionDefinition.
func (f *FunctionWrapper) Definition() api.FunctionDefinition {
	// We need to lock these because they are not thread safe.
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.function.Definition()
}

// Call implements the same method as documented on api.Function.
func (f *FunctionWrapper) Call(ctx context.Context, params ...uint64) (ret []uint64, err error) {
	// We need to lock these because they are not thread safe.
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.function.Call(ctx, params...)
}

// Here is the real implementation of Pdfium
type PdfiumImplementation struct {
	mutex *sync.Mutex

	// Wazero items
	Context   context.Context
	Functions map[string]api.Function
	Module    api.Module

	// lookup tables keeps track of the opened handles for this instance.
	// we need this for handle lookups and in case of closing the instance
	documentRefs               map[references.FPDF_DOCUMENT]*DocumentHandle
	pageRefs                   map[references.FPDF_PAGE]*PageHandle
	bookmarkRefs               map[references.FPDF_BOOKMARK]*BookmarkHandle
	destRefs                   map[references.FPDF_DEST]*DestHandle
	actionRefs                 map[references.FPDF_ACTION]*ActionHandle
	linkRefs                   map[references.FPDF_LINK]*LinkHandle
	pageLinkRefs               map[references.FPDF_PAGELINK]*PageLinkHandle
	schHandleRefs              map[references.FPDF_SCHHANDLE]*SchHandleHandle
	textPageRefs               map[references.FPDF_TEXTPAGE]*TextPageHandle
	pageRangeRefs              map[references.FPDF_PAGERANGE]*PageRangeHandle
	pageObjectRefs             map[references.FPDF_PAGEOBJECT]*PageObjectHandle
	clipPathRefs               map[references.FPDF_CLIPPATH]*ClipPathHandle
	formHandleRefs             map[references.FPDF_FORMHANDLE]*FormHandleHandle
	bitmapRefs                 map[references.FPDF_BITMAP]*BitmapHandle
	annotationRefs             map[references.FPDF_ANNOTATION]*AnnotationHandle
	xObjectRefs                map[references.FPDF_XOBJECT]*XObjectHandle
	signatureRefs              map[references.FPDF_SIGNATURE]*SignatureHandle
	attachmentRefs             map[references.FPDF_ATTACHMENT]*AttachmentHandle
	javaScriptActionRefs       map[references.FPDF_JAVASCRIPT_ACTION]*JavaScriptActionHandle
	searchRefs                 map[references.FPDF_SCHHANDLE]*SearchHandle
	pathSegmentRefs            map[references.FPDF_PATHSEGMENT]*PathSegmentHandle
	dataAvailRefs              map[references.FPDF_AVAIL]*DataAvailHandle
	structTreeRefs             map[references.FPDF_STRUCTTREE]*StructTreeHandle
	structElementRefs          map[references.FPDF_STRUCTELEMENT]*StructElementHandle
	structElementAttributeRefs map[references.FPDF_STRUCTELEMENT_ATTR]*StructElementAttributeHandle
	pageObjectMarkRefs         map[references.FPDF_PAGEOBJECTMARK]*PageObjectMarkHandle
	fontRefs                   map[references.FPDF_FONT]*FontHandle
	glyphPathRefs              map[references.FPDF_GLYPHPATH]*GlyphPathHandle
	fileReaders                map[uint32]*FileReaderRef

	// We need to keep track of our own instance.
	instanceRef int
}

func (p *PdfiumImplementation) Ping() (string, error) {
	return "Pong", nil
}

func (p *PdfiumImplementation) Lock() {
	p.mutex.Lock()
}

func (p *PdfiumImplementation) Unlock() {
	p.mutex.Unlock()
}

func (p *PdfiumImplementation) OpenDocument(request *requests.OpenDocument) (*responses.OpenDocument, error) {
	p.Lock()
	defer p.Unlock()

	var cPasswordPointer uint64
	if request.Password != nil {
		cPassword, err := p.CString(*request.Password)
		if err != nil {
			return nil, err
		}

		defer cPassword.Free()
		cPasswordPointer = cPassword.Pointer
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

	var doc *uint64
	var dataPointer *uint64
	if request.File != nil {
		fileData := *request.File

		dataPtr, err := p.Malloc(uint64(len(fileData)))
		if err != nil {
			return nil, err
		}

		dataPointer = &dataPtr

		if !p.Module.Memory().Write(uint32(dataPtr), fileData) {
			return nil, errors.New("could not write file data to memory")
		}

		// If larger than INT_MAX, use FPDF_LoadMemDocument64
		if len(fileData) > 2147483647 {
			res, err := p.Module.ExportedFunction("FPDF_LoadMemDocument64").Call(p.Context, dataPtr, uint64(len(fileData)), cPasswordPointer)
			if err != nil {
				return nil, err
			}

			// Pointer 0 = document could not be opened.
			if res[0] != 0 {
				doc = &res[0]
			}
		} else {
			res, err := p.Module.ExportedFunction("FPDF_LoadMemDocument").Call(p.Context, dataPtr, uint64(len(fileData)), cPasswordPointer)
			if err != nil {
				return nil, err
			}

			// Pointer 0 = document could not be opened.
			if res[0] != 0 {
				doc = &res[0]
			}
		}
	} else if request.FilePath != nil {
		filePath := *request.FilePath

		// Non-root file, try to absolute it to current working directory.
		// Relative paths are not supported within Webassembly.
		if !strings.HasPrefix(filePath, "/") {
			abs, err := filepath.Abs(filePath)
			if err != nil {
				return nil, err
			}

			filePath = abs
		}

		cFilePathPointer, err := p.CString(filePath)
		if err != nil {
			return nil, err
		}

		defer cFilePathPointer.Free()

		res, err := p.Module.ExportedFunction("FPDF_LoadDocument").Call(p.Context, cFilePathPointer.Pointer, cPasswordPointer)
		if err != nil {
			return nil, err
		}

		// Pointer 0 = document could not be opened.
		if res[0] != 0 {
			doc = &res[0]
		}
	} else if request.FileReader != nil {
		if request.FileReaderSize == 0 {
			return nil, errors.New("FileReaderSize should be given when FileReader is set")
		}

		fileReaderIndex := FileReadersCounter
		FileReadersCounter++

		paramPointer, err := p.Malloc(4)
		if err != nil {
			return nil, err
		}

		p.Module.Memory().WriteUint32Le(uint32(paramPointer), fileReaderIndex)

		res, err := p.Module.ExportedFunction("FPDF_FILEACCESS_Create").Call(p.Context, uint64(request.FileReaderSize), paramPointer)
		if err != nil {
			return nil, err
		}

		fileAccessPointer := res[0]

		fileReaderRef := &FileReaderRef{
			Reader:     request.FileReader,
			FileAccess: &fileAccessPointer,
		}

		FileReaders[fileReaderIndex] = fileReaderRef
		p.fileReaders[fileReaderIndex] = fileReaderRef

		nativeDoc.fileHandleRef = &fileReaderIndex

		res, err = p.Module.ExportedFunction("FPDF_LoadCustomDocument").Call(p.Context, fileAccessPointer, cPasswordPointer)
		if err != nil {
			return nil, err
		}

		// Pointer 0 = document could not be opened.
		if res[0] != 0 {
			doc = &res[0]
		}
	} else {
		return nil, errors.New("No file given")
	}

	if doc == nil {
		errorCode, err := p.Module.ExportedFunction("FPDF_GetLastError").Call(p.Context)
		if err != nil {
			return nil, err
		}

		var pdfiumError error
		switch FPDF_ERR(errorCode[0]) {
		case FPDF_ERR_SUCCESS:
			pdfiumError = pdfium_errors.ErrSuccess
		case FPDF_ERR_UNKNOWN:
			pdfiumError = pdfium_errors.ErrUnknown
		case FPDF_ERR_FILE:
			pdfiumError = pdfium_errors.ErrFile
		case FPDF_ERR_FORMAT:
			pdfiumError = pdfium_errors.ErrFormat
		case FPDF_ERR_PASSWORD:
			pdfiumError = pdfium_errors.ErrPassword
		case FPDF_ERR_SECURITY:
			pdfiumError = pdfium_errors.ErrSecurity
		case FPDF_ERR_PAGE:
			pdfiumError = pdfium_errors.ErrPage
		default:
			pdfiumError = pdfium_errors.ErrUnexpected
		}

		// Cleanup when file loading didn't work.
		if nativeDoc.fileHandleRef != nil {
			p.Free(*p.fileReaders[*nativeDoc.fileHandleRef].FileAccess)
			delete(p.fileReaders, *nativeDoc.fileHandleRef)
		}

		return nil, pdfiumError
	}

	documentRef := uuid.New()
	nativeDoc.handle = doc
	nativeDoc.data = request.File
	nativeDoc.nativeRef = references.FPDF_DOCUMENT(documentRef.String())
	nativeDoc.dataPointer = dataPointer
	p.documentRefs[nativeDoc.nativeRef] = nativeDoc

	return &responses.OpenDocument{
		Document: nativeDoc.nativeRef,
	}, nil
}

func (p *PdfiumImplementation) Close() error {
	p.Lock()
	defer p.Unlock()

	for i := range p.documentRefs {
		err := p.documentRefs[i].Close(p)
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

	for i := range p.structElementAttributeRefs {
		delete(p.structElementAttributeRefs, i)
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
		p.Free(*p.fileReaders[i].FileAccess)

		// Cleanup file handle.
		p.fileReaders[i].FileAccess = nil

		delete(p.fileReaders, i)
		delete(FileReaders, i)
	}

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

func (p *PdfiumImplementation) getStructElementAttributeHandle(structElementAttribute references.FPDF_STRUCTELEMENT_ATTR) (*StructElementAttributeHandle, error) {
	if structElementAttribute == "" {
		return nil, errors.New("structElementAttribute not given")
	}

	if val, ok := p.structElementAttributeRefs[structElementAttribute]; ok {
		return val, nil
	}

	return nil, errors.New("could not find structElementAttribute handle, perhaps the structElementAttribute was already closed or you tried to share structElementAttributes between instances")
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
