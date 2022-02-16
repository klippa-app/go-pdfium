package implementation

/*
#cgo pkg-config: pdfium
#include "fpdfview.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"github.com/klippa-app/go-pdfium/references"
)

type DocumentHandle struct {
	handle        C.FPDF_DOCUMENT
	currentPage   *PageHandle
	data          *[]byte                  // Keep a reference to the data otherwise weird stuff happens
	nativeRef     references.FPDF_DOCUMENT // A string that is our reference inside the process. We need this to close the documents in DestroyLibrary.
	fileHandleRef *string

	// lookup tables keeps track of the opened handles for this instance.
	// we need this for handle lookups and in case of closing the document

	pageRefs             map[references.FPDF_PAGE]*PageHandle
	bookmarkRefs         map[references.FPDF_BOOKMARK]*BookmarkHandle
	destRefs             map[references.FPDF_DEST]*DestHandle
	actionRefs           map[references.FPDF_ACTION]*ActionHandle
	linkRefs             map[references.FPDF_LINK]*LinkHandle
	pageLinkRefs         map[references.FPDF_PAGELINK]*PageLinkHandle
	schHandleRefs        map[references.FPDF_SCHHANDLE]*SchHandleHandle
	textPageRefs         map[references.FPDF_TEXTPAGE]*TextPageHandle
	pageRangeRefs        map[references.FPDF_PAGERANGE]*PageRangeHandle
	formHandleRefs       map[references.FPDF_FORMHANDLE]*FormHandleHandle
	annotationRefs       map[references.FPDF_ANNOTATION]*AnnotationHandle
	signatureRefs        map[references.FPDF_SIGNATURE]*SignatureHandle
	attachmentRefs       map[references.FPDF_ATTACHMENT]*AttachmentHandle
	javaScriptActionRefs map[references.FPDF_JAVASCRIPT_ACTION]*JavaScriptActionHandle
	searchRefs           map[references.FPDF_SCHHANDLE]*SearchHandle
	structTreeRefs       map[references.FPDF_STRUCTTREE]*StructTreeHandle
	structElementRefs    map[references.FPDF_STRUCTELEMENT]*StructElementHandle
}

func (d *DocumentHandle) getPageHandle(pageRef references.FPDF_PAGE) (*PageHandle, error) {
	if pageRef == "" {
		return nil, errors.New("page not given")
	}

	if val, ok := d.pageRefs[pageRef]; ok {
		return val, nil
	}

	return nil, errors.New("could not find page handle, perhaps the page was already closed or you tried to share pages between instances or documents")
}

func (d *DocumentHandle) getBookmarkHandle(bookmarkRef references.FPDF_BOOKMARK) (*BookmarkHandle, error) {
	if bookmarkRef == "" {
		return nil, errors.New("bookmark not given")
	}

	if val, ok := d.bookmarkRefs[bookmarkRef]; ok {
		return val, nil
	}

	return nil, errors.New("could not find bookmark handle, perhaps the bookmark was already closed or you tried to share bookmarks between instances or documents")
}

// Close closes the internal references in FPDF
func (d *DocumentHandle) Close() error {
	if d.handle == nil {
		return errors.New("no current document")
	}

	if d.currentPage != nil {
		d.currentPage.Close()
		d.currentPage = nil
	}

	for i := range d.pageRefs {
		d.pageRefs[i].Close()
		delete(d.pageRefs, i)
	}

	// Remove refs, they don't have a close method.
	for i := range d.bookmarkRefs {
		delete(d.bookmarkRefs, i)
	}

	for i := range d.destRefs {
		delete(d.destRefs, i)
	}

	for i := range d.actionRefs {
		delete(d.actionRefs, i)
	}

	for i := range d.linkRefs {
		delete(d.linkRefs, i)
	}

	for i := range d.pageLinkRefs {
		delete(d.pageLinkRefs, i)
	}

	for i := range d.schHandleRefs {
		delete(d.schHandleRefs, i)
	}

	for i := range d.textPageRefs {
		delete(d.textPageRefs, i)
	}

	for i := range d.pageRangeRefs {
		delete(d.pageRangeRefs, i)
	}

	for i := range d.formHandleRefs {
		delete(d.formHandleRefs, i)
	}

	for i := range d.annotationRefs {
		delete(d.annotationRefs, i)
	}

	for i := range d.signatureRefs {
		delete(d.signatureRefs, i)
	}

	for i := range d.attachmentRefs {
		delete(d.attachmentRefs, i)
	}

	for i := range d.javaScriptActionRefs {
		delete(d.javaScriptActionRefs, i)
	}

	for i := range d.searchRefs {
		delete(d.searchRefs, i)
	}

	for i := range d.structTreeRefs {
		delete(d.structTreeRefs, i)
	}

	for i := range d.structElementRefs {
		delete(d.structElementRefs, i)
	}

	C.FPDF_CloseDocument(d.handle)
	d.handle = nil

	// Remove reference to data.
	if d.data != nil {
		d.data = nil
	}

	// Cleanup file handle.
	if d.fileHandleRef != nil {
		Pdfium.fileReaders[*d.fileHandleRef].fileAccess = nil
		C.free(Pdfium.fileReaders[*d.fileHandleRef].stringRef)
		delete(Pdfium.fileReaders, *d.fileHandleRef)
	}

	delete(Pdfium.documentRefs, d.nativeRef)

	return nil
}

type PageHandle struct {
	handle      C.FPDF_PAGE
	index       int // -1 when unknown.
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_PAGE // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

// Close closes the internal references in FPDF
func (p *PageHandle) Close() {
	if p.handle != nil {
		C.FPDF_ClosePage(p.handle)
		p.handle = nil
	}
}

type BookmarkHandle struct {
	handle      C.FPDF_BOOKMARK
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_BOOKMARK // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type DestHandle struct {
	handle      C.FPDF_DEST
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_DEST // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type ActionHandle struct {
	handle      C.FPDF_ACTION
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_ACTION // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type LinkHandle struct {
	handle      C.FPDF_LINK
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_LINK // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type PageLinkHandle struct {
	handle      C.FPDF_PAGELINK
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_PAGELINK // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type SchHandleHandle struct {
	handle      C.FPDF_SCHHANDLE
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_SCHHANDLE // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type BitmapHandle struct {
	handle    C.FPDF_BITMAP
	nativeRef references.FPDF_BITMAP // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type TextPageHandle struct {
	handle      C.FPDF_TEXTPAGE
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_TEXTPAGE // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type PageRangeHandle struct {
	handle      C.FPDF_PAGERANGE
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_PAGERANGE // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type PageObjectHandle struct {
	handle    C.FPDF_PAGEOBJECT
	nativeRef references.FPDF_PAGEOBJECT // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type ClipPathHandle struct {
	handle    C.FPDF_CLIPPATH
	nativeRef references.FPDF_CLIPPATH // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type FormHandleHandle struct {
	handle      C.FPDF_FORMHANDLE
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_FORMHANDLE // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type AnnotationHandle struct {
	handle      C.FPDF_ANNOTATION
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_ANNOTATION // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type XObjectHandle struct {
	handle    C.FPDF_XOBJECT
	nativeRef references.FPDF_XOBJECT // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type SignatureHandle struct {
	handle      C.FPDF_SIGNATURE
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_SIGNATURE // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type AttachmentHandle struct {
	handle      C.FPDF_ATTACHMENT
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_ATTACHMENT // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type JavaScriptActionHandle struct {
	handle      C.FPDF_JAVASCRIPT_ACTION
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_JAVASCRIPT_ACTION // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type SearchHandle struct {
	handle      C.FPDF_SCHHANDLE
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_SCHHANDLE // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type PathSegmentHandle struct {
	handle    C.FPDF_PATHSEGMENT
	nativeRef references.FPDF_PATHSEGMENT // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type StructTreeHandle struct {
	handle      C.FPDF_STRUCTTREE
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_STRUCTTREE // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type StructElementHandle struct {
	handle      C.FPDF_STRUCTELEMENT
	documentRef references.FPDF_DOCUMENT
	nativeRef   references.FPDF_STRUCTELEMENT // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}

type PageObjectMarkHandle struct {
	handle    C.FPDF_PAGEOBJECTMARK
	nativeRef references.FPDF_PAGEOBJECTMARK // A string that is our reference inside the process. We need this to close the references in DestroyLibrary.
}
