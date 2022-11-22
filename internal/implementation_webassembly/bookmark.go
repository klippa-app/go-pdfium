package implementation_webassembly

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerBookmark(bookmark *uint64, documentHandle *DocumentHandle) *BookmarkHandle {
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

	res, err := p.Module.ExportedFunction("FPDFBookmark_GetFirstChild").Call(p.Context, *documentHandle.handle, 0)
	if err != nil {
		return nil, err
	}
	bookmark := res[0]
	if bookmark == 0 {
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

func (p *PdfiumImplementation) getBookMarkChildren(documentHandle *DocumentHandle, bookmark uint64) ([]responses.GetBookmarksBookmark, error) {
	bookmarks := []uint64{bookmark}

	currentSibling := bookmark
	for {
		res, err := p.Module.ExportedFunction("FPDFBookmark_GetNextSibling").Call(p.Context, *documentHandle.handle, currentSibling)
		if err != nil {
			return nil, err
		}
		newSibling := res[0]
		if newSibling == 0 {
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
		res, err := p.Module.ExportedFunction("FPDFBookmark_GetFirstChild").Call(p.Context, *documentHandle.handle, bookmarks[i])
		if err != nil {
			return nil, err
		}
		child := res[0]
		if child != 0 {
			myChildren, err := p.getBookMarkChildren(documentHandle, child)
			if err != nil {
				return nil, err
			}

			respBookmark.Children = myChildren
		}

		bookmarkHandle := p.registerBookmark(&bookmarks[i], documentHandle)

		// First get the title length.
		res, err = p.Module.ExportedFunction("FPDFBookmark_GetTitle").Call(p.Context, *bookmarkHandle.handle, 0, 0)
		if err != nil {
			return nil, err
		}

		titleSize := *(*int32)(unsafe.Pointer(&res[0]))
		if titleSize == 0 {
			return nil, errors.New("Could not get title")
		}

		charDataPointer, err := p.ByteArrayPointer(uint64(titleSize), nil)
		defer charDataPointer.Free()

		_, err = p.Module.ExportedFunction("FPDFBookmark_GetTitle").Call(p.Context, *bookmarkHandle.handle, charDataPointer.Pointer, uint64(titleSize))
		if err != nil {
			return nil, err
		}

		charData, err := charDataPointer.Value(false)
		if err != nil {
			return nil, err
		}

		transformedText, err := p.transformUTF16LEToUTF8(charData)
		if err != nil {
			return nil, err
		}

		res, err = p.Module.ExportedFunction("FPDFBookmark_GetAction").Call(p.Context, *documentHandle.handle)
		if err != nil {
			return nil, err
		}
		action := res[0]
		if action == 0 {
			actionHandle := p.registerAction(&action)
			actionInfo, err := p.getActionInfo(actionHandle, documentHandle)
			if err != nil {
				return nil, err
			}

			respBookmark.ActionInfo = actionInfo
		}

		res, err = p.Module.ExportedFunction("FPDFBookmark_GetAction").Call(p.Context, *documentHandle.handle, *bookmarkHandle.handle)
		if err != nil {
			return nil, err
		}
		dest := res[0]
		if dest != 0 {
			destHandle := p.registerDest(&dest, documentHandle)

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
