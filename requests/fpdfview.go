package requests

import (
	"github.com/klippa-app/go-pdfium/references"
	"io"
)

type FPDF_LoadDocument struct {
	Path     *string // A path to a PDF file.
	Password *string // The password of the document.
}

type FPDF_LoadMemDocument struct {
	Data     *[]byte // A reference to the file data.
	Password *string // The password of the document.
}

type FPDF_LoadMemDocument64 struct {
	Data     *[]byte // A reference to the file data.
	Password *string // The password of the document.
}

type FPDF_LoadCustomDocument struct {
	Reader   io.ReadSeeker
	Size     int64
	Password *string // The password of the document.
}

type FPDF_GetLastError struct{}

type FPDF_SetSandBoxPolicyPolicy uint32

const (
	FPDF_SetSandBoxPolicyPolicyMachinetimeAccess FPDF_SetSandBoxPolicyPolicy = 1 // Policy for accessing the local machine time.
)

type FPDF_SetSandBoxPolicy struct {
	Policy FPDF_SetSandBoxPolicyPolicy
	Enable bool
}

type FPDF_CloseDocument struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_LoadPage struct {
	Document references.FPDF_DOCUMENT
	Index    int // The page number (0-index based).
}

type FPDF_ClosePage struct {
	Page references.FPDF_PAGE
}

type FPDF_GetFileVersion struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetDocPermissions struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetSecurityHandlerRevision struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetPageCount struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetPageWidth struct {
	Page Page
}

type FPDF_GetPageHeight struct {
	Page Page
}

type FPDF_GetPageSizeByIndex struct {
	Document references.FPDF_DOCUMENT
	Index    int
}
