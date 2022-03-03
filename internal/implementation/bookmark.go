package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_doc.h"
import "C"
import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerBookmark(bookmark C.FPDF_BOOKMARK, documentHandle *DocumentHandle) *BookmarkHandle {
	bookmarkRef := uuid.New()
	bookmarkHandle := &BookmarkHandle{
		handle:      bookmark,
		nativeRef:   references.FPDF_BOOKMARK(bookmarkRef.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.bookmarkRefs[bookmarkHandle.nativeRef] = bookmarkHandle
	p.bookmarkRefs[bookmarkHandle.nativeRef] = bookmarkHandle

	return bookmarkHandle
}

// GetBookmarks returns all the bookmarks of a document.
func (p *PdfiumImplementation) GetBookmarks(request *requests.GetBookmarks) (*responses.GetBookmarks, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	bookmark := C.FPDFBookmark_GetFirstChild(documentHandle.handle, nil)
	if bookmark == nil {
		return &responses.GetBookmarks{}, nil
	}

	bookMarks, err := p.getBookMarkChildren(documentHandle, bookmark)
	if err != nil {
		return nil, err
	}

	return &responses.GetBookmarks{
		Bookmarks: bookMarks,
	}, nil
}

func (p *PdfiumImplementation) getBookMarkChildren(documentHandle *DocumentHandle, bookmark C.FPDF_BOOKMARK) ([]responses.GetBookmarksBookmark, error) {
	bookmarks := []C.FPDF_BOOKMARK{bookmark}

	currentSibling := bookmark
	for {
		newSibling := C.FPDFBookmark_GetNextSibling(documentHandle.handle, currentSibling)
		if newSibling == nil {
			break
		}

		currentSibling = newSibling
		bookmarks = append(bookmarks, newSibling)
	}

	resp := []responses.GetBookmarksBookmark{}

	for i := range bookmarks {
		respBookmark := responses.GetBookmarksBookmark{
			Children: []responses.GetBookmarksBookmark{},
		}
		child := C.FPDFBookmark_GetFirstChild(documentHandle.handle, bookmarks[i])
		if child != nil {
			myChildren, err := p.getBookMarkChildren(documentHandle, child)
			if err != nil {
				return nil, err
			}

			respBookmark.Children = myChildren
		}

		bookmarkHandle := p.registerBookmark(bookmarks[i], documentHandle)

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

		action := C.FPDFBookmark_GetAction(bookmarkHandle.handle)
		if action == nil {
			actionHandle := p.registerAction(action)
			actionInfo, err := p.getActionInfo(actionHandle, documentHandle)
			if err != nil {
				return nil, err
			}

			respBookmark.ActionInfo = actionInfo
		}

		dest := C.FPDFBookmark_GetDest(documentHandle.handle, bookmarkHandle.handle)
		if dest != nil {
			destHandle := p.registerDest(dest, documentHandle)

			destInfo, err := p.getDestInfo(destHandle, documentHandle)
			if err != nil {
				return nil, err
			}

			respBookmark.DestInfo = destInfo
		}

		respBookmark.Title = transformedText
		respBookmark.Reference = bookmarkHandle.nativeRef

		resp = append(resp, respBookmark)
	}

	return resp, nil
}
