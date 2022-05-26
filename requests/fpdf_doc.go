package requests

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
)

type FPDFBookmark_GetFirstChild struct {
	Document references.FPDF_DOCUMENT
	Bookmark *references.FPDF_BOOKMARK // Reference to the current bookmark. Can be nil if you want to get the first top level item.
}

type FPDFBookmark_GetNextSibling struct {
	Document references.FPDF_DOCUMENT
	Bookmark references.FPDF_BOOKMARK // Reference to the current bookmark. Cannot be nil.
}

type FPDFBookmark_GetTitle struct {
	Bookmark references.FPDF_BOOKMARK // Reference to the current bookmark.
}

type FPDFBookmark_GetCount struct {
	Bookmark references.FPDF_BOOKMARK // Reference to the current bookmark.
}

type FPDFBookmark_Find struct {
	Document references.FPDF_DOCUMENT
	Title    string // The string for the bookmark title to be searched
}

type FPDFBookmark_GetDest struct {
	Document references.FPDF_DOCUMENT
	Bookmark references.FPDF_BOOKMARK
}

type FPDFBookmark_GetAction struct {
	Bookmark references.FPDF_BOOKMARK
}

type FPDFAction_GetType struct {
	Action references.FPDF_ACTION
}

type FPDFAction_GetDest struct {
	Document references.FPDF_DOCUMENT
	Action   references.FPDF_ACTION
}

type FPDFAction_GetFilePath struct {
	Action references.FPDF_ACTION
}

type FPDFAction_GetURIPath struct {
	Document references.FPDF_DOCUMENT
	Action   references.FPDF_ACTION
}

type FPDFDest_GetDestPageIndex struct {
	Document references.FPDF_DOCUMENT
	Dest     references.FPDF_DEST
}

type FPDF_GetFileIdentifier struct {
	Document   references.FPDF_DOCUMENT
	FileIdType enums.FPDF_FILEIDTYPE
}

type FPDF_GetMetaText struct {
	Document references.FPDF_DOCUMENT
	Tag      string // A metadata tag. Title, Author, Subject, Keywords, Creator, Producer, CreationDate, ModDate. For detailed explanation of these tags and their respective values, please refer to section 10.2.1 "Document Information Dictionary" in PDF Reference 1.7.
}

type FPDF_GetPageLabel struct {
	Document references.FPDF_DOCUMENT
	Page     int // The page number (0-index based).
}

type FPDFDest_GetView struct {
	Dest references.FPDF_DEST
}

type FPDFDest_GetLocationInPage struct {
	Dest references.FPDF_DEST
}

type FPDFLink_GetLinkAtPoint struct {
	Page Page
	X    float64
	Y    float64
}

type FPDFLink_GetLinkZOrderAtPoint struct {
	Page Page
	X    float64
	Y    float64
}

type FPDFLink_GetDest struct {
	Document references.FPDF_DOCUMENT
	Link     references.FPDF_LINK
}

type FPDFLink_GetAction struct {
	Link references.FPDF_LINK
}

type FPDFLink_Enumerate struct {
	Page     Page
	StartPos int
}

type FPDFLink_GetAnnot struct {
	Page Page
	Link references.FPDF_LINK
}

type FPDFLink_GetAnnotRect struct {
	Link references.FPDF_LINK
}

type FPDFLink_CountQuadPoints struct {
	Link references.FPDF_LINK
}

type FPDFLink_GetQuadPoints struct {
	Link      references.FPDF_LINK
	QuadIndex int
}

type FPDF_GetPageAAction struct {
	Page   Page
	AAType enums.FPDF_PAGE_AACTION
}
