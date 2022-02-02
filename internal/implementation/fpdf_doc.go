package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_doc.h"
// #include <stdlib.h>
import "C"

import (
	"bytes"
	"encoding/ascii85"
	"errors"
	"github.com/klippa-app/go-pdfium/enums"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFBookmark_GetFirstChild returns the first child of a bookmark item, or the first top level bookmark item.
func (p *PdfiumImplementation) FPDFBookmark_GetFirstChild(request *requests.FPDFBookmark_GetFirstChild) (*responses.FPDFBookmark_GetFirstChild, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	var parentBookMark C.FPDF_BOOKMARK
	if request.Bookmark != nil {
		bookmarkHandle, err := p.getBookmarkHandle(*request.Bookmark)
		if err != nil {
			return nil, err
		}

		parentBookMark = bookmarkHandle.handle
	}

	bookmark := C.FPDFBookmark_GetFirstChild(documentHandle.handle, parentBookMark)
	if bookmark == nil {
		return &responses.FPDFBookmark_GetFirstChild{}, nil
	}

	newNativeBookmark := p.registerBookmark(bookmark, documentHandle)

	return &responses.FPDFBookmark_GetFirstChild{
		Bookmark: &newNativeBookmark.nativeRef,
	}, nil
}

// FPDFBookmark_GetNextSibling returns the next bookmark item at the same level.
func (p *PdfiumImplementation) FPDFBookmark_GetNextSibling(request *requests.FPDFBookmark_GetNextSibling) (*responses.FPDFBookmark_GetNextSibling, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	bookmarkHandle, err := p.getBookmarkHandle(request.Bookmark)
	if err != nil {
		return nil, err
	}

	bookmark := C.FPDFBookmark_GetNextSibling(documentHandle.handle, bookmarkHandle.handle)
	if bookmark == nil {
		return &responses.FPDFBookmark_GetNextSibling{}, nil
	}

	newNativeBookmark := p.registerBookmark(bookmark, documentHandle)

	return &responses.FPDFBookmark_GetNextSibling{
		Bookmark: &newNativeBookmark.nativeRef,
	}, nil
}

// FPDFBookmark_GetTitle returns the title of a bookmark.
func (p *PdfiumImplementation) FPDFBookmark_GetTitle(request *requests.FPDFBookmark_GetTitle) (*responses.FPDFBookmark_GetTitle, error) {
	p.Lock()
	defer p.Unlock()

	bookmarkHandle, err := p.getBookmarkHandle(request.Bookmark)
	if err != nil {
		return nil, err
	}

	// First get the title length.
	titleSize := C.FPDFBookmark_GetTitle(bookmarkHandle.handle, C.NULL, 0)
	if titleSize == 0 {
		return nil, errors.New("Could not get title")
	}

	charData := make([]byte, titleSize)
	C.FPDFBookmark_GetTitle(bookmarkHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFBookmark_GetTitle{
		Title: transformedText,
	}, nil
}

// FPDFBookmark_Find finds a bookmark in the document, using the bookmark title.
func (p *PdfiumImplementation) FPDFBookmark_Find(request *requests.FPDFBookmark_Find) (*responses.FPDFBookmark_Find, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	if request.Title == "" {
		return nil, errors.New("no title given")
	}

	transformedText, err := p.transformUTF8ToUTF16LE(request.Title)
	if err != nil {
		return nil, err
	}

	bookmark := C.FPDFBookmark_Find(documentHandle.handle, (C.FPDF_WIDESTRING)(unsafe.Pointer(&transformedText[0])))
	if bookmark == nil {
		return &responses.FPDFBookmark_Find{}, nil
	}

	newNativeBookmark := p.registerBookmark(bookmark, documentHandle)

	return &responses.FPDFBookmark_Find{
		Bookmark: &newNativeBookmark.nativeRef,
	}, nil
}

// FPDFBookmark_GetDest returns the destination associated with a bookmark item.
// If the returned destination is nil, none is associated to the bookmark item.
func (p *PdfiumImplementation) FPDFBookmark_GetDest(request *requests.FPDFBookmark_GetDest) (*responses.FPDFBookmark_GetDest, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	bookmarkHandle, err := p.getBookmarkHandle(request.Bookmark)
	if err != nil {
		return nil, err
	}

	dest := C.FPDFBookmark_GetDest(documentHandle.handle, bookmarkHandle.handle)
	if dest == nil {
		return &responses.FPDFBookmark_GetDest{}, nil
	}

	destHandle := p.registerDest(dest, documentHandle)

	return &responses.FPDFBookmark_GetDest{
		Dest: &destHandle.nativeRef,
	}, nil
}

// FPDFBookmark_GetAction returns the action associated with a bookmark item.
// If the returned action is nil, you should try FPDFBookmark_GetDest.
func (p *PdfiumImplementation) FPDFBookmark_GetAction(request *requests.FPDFBookmark_GetAction) (*responses.FPDFBookmark_GetAction, error) {
	p.Lock()
	defer p.Unlock()

	bookmarkHandle, err := p.getBookmarkHandle(request.Bookmark)
	if err != nil {
		return nil, err
	}

	action := C.FPDFBookmark_GetAction(bookmarkHandle.handle)
	if action == nil {
		return &responses.FPDFBookmark_GetAction{}, nil
	}

	documentHandle, err := p.getDocumentHandle(bookmarkHandle.documentRef)
	if err != nil {
		return nil, err
	}

	actionHandle := p.registerAction(action, documentHandle)

	return &responses.FPDFBookmark_GetAction{
		Action: &actionHandle.nativeRef,
	}, nil
}

// FPDFAction_GetType returns the action associated with a bookmark item.
func (p *PdfiumImplementation) FPDFAction_GetType(request *requests.FPDFAction_GetType) (*responses.FPDFAction_GetType, error) {
	p.Lock()
	defer p.Unlock()

	actionHandle, err := p.getActionHandle(request.Action)
	if err != nil {
		return nil, err
	}

	actionType := C.FPDFAction_GetType(actionHandle.handle)

	return &responses.FPDFAction_GetType{
		Type: enums.FPDF_ACTION_ACTION(actionType),
	}, nil
}

// FPDFAction_GetDest returns the destination of a specific go-to or remote-goto action.
// Only action with type PDF_ACTION_ACTION_GOTO and PDF_ACTION_ACTION_REMOTEGOTO can have destination data.
// In case of remote goto action, the application should first use function FPDFAction_GetFilePath to get file path, then load that particular document, and use its document handle to call this function.
func (p *PdfiumImplementation) FPDFAction_GetDest(request *requests.FPDFAction_GetDest) (*responses.FPDFAction_GetDest, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	actionHandle, err := p.getActionHandle(request.Action)
	if err != nil {
		return nil, err
	}

	dest := C.FPDFAction_GetDest(documentHandle.handle, actionHandle.handle)
	if dest == nil {
		return &responses.FPDFAction_GetDest{}, nil
	}

	destHandle := p.registerDest(dest, documentHandle)

	return &responses.FPDFAction_GetDest{
		Dest: &destHandle.nativeRef,
	}, nil
}

// FPDFAction_GetFilePath returns the file path from a remote goto or launch action.
// Only works on actions that have the type FPDF_ACTION_ACTION_REMOTEGOTO or FPDF_ACTION_ACTION_LAUNCH.
func (p *PdfiumImplementation) FPDFAction_GetFilePath(request *requests.FPDFAction_GetFilePath) (*responses.FPDFAction_GetFilePath, error) {
	p.Lock()
	defer p.Unlock()

	actionHandle, err := p.getActionHandle(request.Action)
	if err != nil {
		return nil, err
	}

	// First get the file path length.
	filePathLength := C.FPDFAction_GetFilePath(actionHandle.handle, C.NULL, 0)
	if filePathLength == 0 {
		return nil, errors.New("Could not get file path")
	}

	charData := make([]byte, filePathLength)
	// FPDFAction_GetFilePath returns the data in UTF-8, no conversion needed.
	C.FPDFAction_GetFilePath(actionHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	return &responses.FPDFAction_GetFilePath{
		FilePath: string(charData),
	}, nil
}

// FPDFAction_GetURIPath returns the URI path from a URI action.
func (p *PdfiumImplementation) FPDFAction_GetURIPath(request *requests.FPDFAction_GetURIPath) (*responses.FPDFAction_GetURIPath, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	actionHandle, err := p.getActionHandle(request.Action)
	if err != nil {
		return nil, err
	}

	// First get the uri path length.
	uriPathLength := C.FPDFAction_GetURIPath(documentHandle.handle, actionHandle.handle, C.NULL, 0)
	if uriPathLength == 0 {
		return nil, errors.New("Could not get uri path")
	}

	charData := make([]byte, uriPathLength)
	C.FPDFAction_GetURIPath(documentHandle.handle, actionHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	// Convert 7-bit ASCII to UTF-8.
	dst := make([]byte, uriPathLength, uriPathLength)
	_, _, err = ascii85.Decode(dst, charData, true)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAction_GetURIPath{
		URIPath: string(dst),
	}, nil
}

// FPDFDest_GetDestPageIndex returns the page index from destination data.
func (p *PdfiumImplementation) FPDFDest_GetDestPageIndex(request *requests.FPDFDest_GetDestPageIndex) (*responses.FPDFDest_GetDestPageIndex, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	destHandle, err := p.getDestHandle(request.Dest)
	if err != nil {
		return nil, err
	}

	pageIndex := C.FPDFDest_GetDestPageIndex(documentHandle.handle, destHandle.handle)
	return &responses.FPDFDest_GetDestPageIndex{
		Index: int(pageIndex),
	}, nil
}

// FPDF_GetFileIdentifier Get the file identifier defined in the trailer of a document.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetFileIdentifier(request *requests.FPDF_GetFileIdentifier) (*responses.FPDF_GetFileIdentifier, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	if request.FileIdType != enums.FPDF_FILEIDTYPE_PERMANENT && request.FileIdType != enums.FPDF_FILEIDTYPE_CHANGING {
		return nil, errors.New("invalid file id type given")
	}

	// First get the identifier length.
	identifierSize := C.FPDF_GetFileIdentifier(documentHandle.handle, C.FPDF_FILEIDTYPE(request.FileIdType), C.NULL, 0)
	if identifierSize == 0 {
		return &responses.FPDF_GetFileIdentifier{
			FileIdType: request.FileIdType,
			Identifier: nil,
		}, nil
	}

	charData := make([]byte, identifierSize)
	C.FPDF_GetFileIdentifier(documentHandle.handle, C.FPDF_FILEIDTYPE(request.FileIdType), unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	// Remove NULL terminator.
	charData = bytes.TrimSuffix(charData, []byte("\x00"))

	return &responses.FPDF_GetFileIdentifier{
		FileIdType: request.FileIdType,
		Identifier: charData,
	}, nil
}

// FPDF_GetMetaText returns the requested metadata.
func (p *PdfiumImplementation) FPDF_GetMetaText(request *requests.FPDF_GetMetaText) (*responses.FPDF_GetMetaText, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	cstr := C.CString(request.Tag)
	defer C.free(unsafe.Pointer(cstr))

	// First get the metadata length.
	metaSize := C.FPDF_GetMetaText(documentHandle.handle, cstr, C.NULL, 0)
	if metaSize == 0 {
		return nil, errors.New("Could not get metadata")
	}

	charData := make([]byte, metaSize)
	C.FPDF_GetMetaText(documentHandle.handle, cstr, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_GetMetaText{
		Tag:   request.Tag,
		Value: transformedText,
	}, nil
}

// FPDF_GetPageLabel returns the label for the given page.
func (p *PdfiumImplementation) FPDF_GetPageLabel(request *requests.FPDF_GetPageLabel) (*responses.FPDF_GetPageLabel, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	// First get the label length.
	labelSize := C.FPDF_GetPageLabel(documentHandle.handle, C.int(request.Page), C.NULL, 0)
	if labelSize == 0 {
		return nil, errors.New("Could not get label")
	}

	charData := make([]byte, labelSize)
	C.FPDF_GetPageLabel(documentHandle.handle, C.int(request.Page), unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_GetPageLabel{
		Page:  request.Page,
		Label: transformedText,
	}, nil
}
