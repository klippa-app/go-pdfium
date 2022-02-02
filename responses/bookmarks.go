package responses

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
)

type GetBookmarksDest struct {
	Reference references.FPDF_DEST
	PageIndex int
}

type GetBookmarksAction struct {
	Reference references.FPDF_ACTION
	Type      enums.FPDF_ACTION_ACTION
	Dest      *GetBookmarksDest // Is set when the action is GOTO. When the action is REMOTEGOTO, we will not fetch the destination.
	FilePath  *string           // When action is LAUNCH or REMOTEGOTO.
	URIPath   *string           // When action is URI.
}

type GetBookmarksBookmark struct {
	Title     string
	Reference references.FPDF_BOOKMARK
	Action    *GetBookmarksAction
	Dest      *GetBookmarksDest
	Children  []GetBookmarksBookmark
}

type GetBookmarks struct {
	Bookmarks []GetBookmarksBookmark
}
