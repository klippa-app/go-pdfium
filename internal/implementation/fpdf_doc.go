package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_doc.h"
// #include <stdlib.h>
import "C"

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

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

	newNativeBookmark := p.registerBookMark(bookmark, documentHandle)

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

	newNativeBookmark := p.registerBookMark(bookmark, documentHandle)

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

	newNativeBookmark := p.registerBookMark(bookmark, documentHandle)

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

	bookmark := C.FPDFBookmark_GetDest(documentHandle.handle, bookmarkHandle.handle)
	if bookmark == nil {
		return &responses.FPDFBookmark_GetDest{}, nil
	}

	return &responses.FPDFBookmark_GetDest{
		//Bookmark: &newNativeBookmark.nativeRef,
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

	bookmark := C.FPDFBookmark_GetAction(bookmarkHandle.handle)
	if bookmark == nil {
		return &responses.FPDFBookmark_GetAction{}, nil
	}

	return &responses.FPDFBookmark_GetAction{
		//Bookmark: &newNativeBookmark.nativeRef,
	}, nil
}
