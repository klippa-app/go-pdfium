package responses

import "github.com/klippa-app/go-pdfium/references"

type GetBookmarksBookmark struct {
	Title     string
	Reference references.FPDF_BOOKMARK
	Children  []GetBookmarksBookmark
}

type GetBookmarks struct {
	Bookmarks []GetBookmarksBookmark
}
