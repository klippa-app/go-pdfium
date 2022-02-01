package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_doc.h"
import "C"
import (
	"errors"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"unsafe"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerBookMark(bookmark C.FPDF_BOOKMARK, nativeDocument *NativeDocument) *NativeBookmark {
	bookmarkRef := uuid.New()
	newNativeBookmark := &NativeBookmark{
		bookmark:    bookmark,
		nativeRef:   references.FPDF_BOOKMARK(bookmarkRef.String()),
		documentRef: nativeDocument.nativeRef,
	}

	nativeDocument.bookmarkRefs[newNativeBookmark.nativeRef] = newNativeBookmark
	p.bookmarkRefs[newNativeBookmark.nativeRef] = newNativeBookmark

	return newNativeBookmark
}

// GetBookmarks returns all the bookmarks of a document.
func (p *PdfiumImplementation) GetBookmarks(request *requests.GetBookmarks) (*responses.GetBookmarks, error) {
	p.Lock()
	defer p.Unlock()

	nativeDoc, err := p.getNativeDocument(request.Document)
	if err != nil {
		return nil, err
	}

	bookmark := C.FPDFBookmark_GetFirstChild(nativeDoc.doc, nil)
	if bookmark == nil {
		return &responses.GetBookmarks{}, nil
	}

	bookMarks, err := getBookMarkChildren(p, nativeDoc, bookmark)
	if err != nil {
		return nil, err
	}

	return &responses.GetBookmarks{
		Bookmarks: bookMarks,
	}, nil
}

func getBookMarkChildren(p *PdfiumImplementation, nativeDoc *NativeDocument, bookmark C.FPDF_BOOKMARK) ([]responses.GetBookmarksBookmark, error) {
	bookmarks := []C.FPDF_BOOKMARK{bookmark}

	currentSibling := bookmark
	for {
		newSibling := C.FPDFBookmark_GetNextSibling(nativeDoc.doc, currentSibling)
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
		child := C.FPDFBookmark_GetFirstChild(nativeDoc.doc, bookmarks[i])
		if child != nil {
			myChildren, err := getBookMarkChildren(p, nativeDoc, child)
			if err != nil {
				return nil, err
			}

			respBookmark.Children = myChildren
		}

		nativeBookmark := p.registerBookMark(bookmarks[i], nativeDoc)

		// First get the title length.
		titleSize := C.FPDFBookmark_GetTitle(nativeBookmark.bookmark, C.NULL, 0)
		if titleSize == 0 {
			return nil, errors.New("Could not get title")
		}

		charData := make([]byte, titleSize)
		C.FPDFBookmark_GetTitle(nativeBookmark.bookmark, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

		transformedText, err := p.transformUTF16LEToUTF8(charData)
		if err != nil {
			return nil, err
		}

		respBookmark.Title = transformedText
		respBookmark.Reference = nativeBookmark.nativeRef

		resp = append(resp, respBookmark)
	}

	return resp, nil
}
