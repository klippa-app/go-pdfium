package requests

import "github.com/klippa-app/go-pdfium/references"

type FPDF_GetMetaText struct {
	Document references.FPDF_DOCUMENT
	Tag      string // A metadata tag. Title, Author, Subject, Keywords, Creator, Producer, CreationDate, ModDate. For detailed explanation of these tags and their respective values, please refer to section 10.2.1 "Document Information Dictionary" in PDF Reference 1.7.
}

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
