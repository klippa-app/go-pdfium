package requests

import (
	"io"

	"github.com/klippa-app/go-pdfium/references"
)

type OpenDocument struct {
	File           *[]byte // A reference to the file data.
	FilePath       *string // A path to a PDF file.
	FileReader     io.ReadSeeker
	FileReaderSize int
	Password       *string // The password of the document.
}

// Page can either be the index of a page or a page reference.
// When you use an index. The library will always cache the last opened page.
type Page struct {
	Index     int                  // The page number (0-index based).
	Reference references.FPDF_PAGE // A reference to a page. Received by GetPage()
}

type NewPage struct {
	Document references.FPDF_DOCUMENT
	Index    int     // A zero-based index which specifies the position of the created page in PDF document. Range: 0 to (pagecount-1). If this value is below 0, the new page will be inserted to the first. If this value is above (pagecount-1), the new page will be inserted to the last.
	Width    float64 // The page width in points.
	Height   float64 // The page height in points.
}

type GetPageSize struct {
	Document references.FPDF_DOCUMENT
	Page     Page
}

type GetPageSizeInPixels struct {
	Document references.FPDF_DOCUMENT
	Page     Page
	DPI      int // The DPI to calculate the size for.
}

type PageRotation int

const (
	PageRotationNone  PageRotation = 0 // 0: no rotation.
	PageRotation90CW  PageRotation = 1 // 1: rotate 90 degrees in clockwise direction.
	PageRotation180CW PageRotation = 2 // 2: rotate 180 degrees in clockwise direction.
	PageRotation270CW PageRotation = 3 // 3: rotate 270 degrees in clockwise direction.
)
