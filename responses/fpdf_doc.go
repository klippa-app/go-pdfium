package responses

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
)

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
	Bookmark *references.FPDF_BOOKMARK // Reference to the found bookmark. nil if the title can't be found.
}

type FPDFBookmark_GetDest struct {
	Dest *references.FPDF_DEST // Reference to the bookmark dest. nil if not found.
}

type FPDFBookmark_GetAction struct {
	Action *references.FPDF_ACTION // Reference to the bookmark action. nil if not found.
}

type FPDFAction_GetType struct {
	Type enums.FPDF_ACTION_ACTION
}

type FPDFAction_GetDest struct {
	Dest *references.FPDF_DEST // Reference to the bookmark dest. nil if not found.
}

type FPDFAction_GetFilePath struct {
	FilePath string
}

type FPDFAction_GetURIPath struct {
	URIPath string
}

type FPDFDest_GetDestPageIndex struct {
	Index int
}

type FPDF_GetFileIdentifier struct {
	FileIdType enums.FPDF_FILEIDTYPE
	Identifier []byte // Can be nil if no identifier was found
}

type FPDF_GetMetaText struct {
	Tag   string // The requested metadata tag.
	Value string // The value of the tag if found, string is empty if the value is not found.
}

type FPDF_GetPageLabel struct {
	Page  int
	Label string
}
