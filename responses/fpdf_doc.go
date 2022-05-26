package responses

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/structs"
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

type FPDFBookmark_GetCount struct {
	// A signed integer that represents the number of sub-items the given
	// bookmark has. If the value is positive, child items shall be shown by default
	// (open state). If the value is negative, child items shall be hidden by
	// default (closed state). Please refer to PDF 32000-1:2008, Table 153.
	// Returns 0 if the bookmark has no children or is invalid.
	Count int
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
	FilePath *string // nil when not set
}

type FPDFAction_GetURIPath struct {
	URIPath *string // nil when not set
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

type FPDFDest_GetView struct {
	DestView enums.FPDF_PDFDEST_VIEW
	Params   []float32
}

type FPDFDest_GetLocationInPage struct {
	X    *float32
	Y    *float32
	Zoom *float32
}

type FPDFLink_GetLinkAtPoint struct {
	Link *references.FPDF_LINK // Reference to the found link. nil if not found.
}

type FPDFLink_GetLinkZOrderAtPoint struct {
	ZOrder int // the Z-order of the link, or -1 if no link found at the given point. Larger Z-order numbers are closer to the front.
}

type FPDFLink_GetDest struct {
	// Dest is a handle to the destination, or nil if there is no destination
	// associated with the link. In this case, you should call FPDFLink_GetAction()
	// to retrieve the action associated with a link.
	Dest *references.FPDF_DEST
}

type FPDFLink_GetAction struct {
	// Action is a handle to the action associated to a link, or nil if no action.
	Action *references.FPDF_ACTION
}

type FPDFLink_Enumerate struct {
	NextStartPos *int
	Link         *references.FPDF_LINK
}

type FPDFLink_GetAnnot struct {
	Annotation *references.FPDF_ANNOTATION
}

type FPDFLink_GetAnnotRect struct {
	Rect *structs.FPDF_FS_RECTF // The rectangle for the link.
}

type FPDFLink_CountQuadPoints struct {
	Count int
}

type FPDFLink_GetQuadPoints struct {
	Points *structs.FPDF_FS_QUADPOINTSF
}

type FPDF_GetPageAAction struct {
	AAType *enums.FPDF_PAGE_AACTION
	Action *references.FPDF_ACTION
}
