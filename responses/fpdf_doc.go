package responses

import "github.com/klippa-app/go-pdfium/references"

type FPDF_GetMetaText struct {
	Tag   string // The requested metadata tag.
	Value string // The value of the tag if found, string is empty if the value is not found.
}

type FPDFBookmark_GetFirstChild struct {
	Bookmark *references.FPDF_BOOKMARK // Reference to the first child or top level bookmark item. nil if no child or top level bookmark found.
}

type FPDFBookmark_GetNextSibling struct {
	Bookmark *references.FPDF_BOOKMARK // Reference to the next bookmark item at the same level. nil if this is the last bookmark at this level.
}

type FPDFBookmark_GetTitle struct {
	Title string // The title of the bookmark.
}

type FPDFBookmark_Find struct {
	Bookmark *references.FPDF_BOOKMARK // Reference to the found bookmark item. nil if the title can't be found.
}
