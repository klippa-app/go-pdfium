package pdfium

import (
	"github.com/klippa-app/go-pdfium/references"
	"io"
	"time"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

type NewDocumentOption interface {
	AlterOpenDocumentRequest(*requests.OpenDocument)
}

type openDocumentWithPassword struct{ password string }

func (p openDocumentWithPassword) AlterOpenDocumentRequest(r *requests.OpenDocument) {
	r.Password = &p.password
}

// OpenDocumentWithPasswordOption can be used as NewDocumentOption when your PDF contains a password.
func OpenDocumentWithPasswordOption(password string) NewDocumentOption {
	return openDocumentWithPassword{
		password: password,
	}
}

type Pool interface {
	// GetInstance returns an instance to the pool.
	// For single-threaded this is thread safe, but you can only do one pdfium action at the same time.
	// For multi-threaded it will try to get a worker from the pool for the length of timeout
	// It is important to Close instances when you are done with them. To either return them to the pool
	// or clear it's resources.
	GetInstance(timeout time.Duration) (Pdfium, error)

	// Close closes the pool.
	// It will close any unclosed instances.
	// For single-threaded it will unload the library if it's the last pool.
	// For multi-threaded it will stop all the pool workers.
	Close() error
}

// Pdfium describes a Pdfium instance.
type Pdfium interface {
	// Start instance functions.

	// NewDocumentFromBytes returns a pdfium Document from the given PDF bytes.
	// This is a helper around OpenDocument. This is a gateway to FPDF_LoadMemDocument.
	NewDocumentFromBytes(file *[]byte, opts ...NewDocumentOption) (*references.FPDF_DOCUMENT, error)

	// NewDocumentFromFilePath returns a pdfium Document from the given PDF file path.
	// This is a helper around OpenDocument. This is a gateway to FPDF_LoadDocument.
	NewDocumentFromFilePath(filePath string, opts ...NewDocumentOption) (*references.FPDF_DOCUMENT, error)

	// NewDocumentFromReader returns a pdfium Document from the given PDF file reader.
	// This is a helper around OpenDocument.
	// This is only really efficient for single threaded usage, the multi-threaded
	// usage will just load the file in memory because it can't transfer readers
	// over gRPC. The single-threaded usage will actually efficiently walk over
	// the PDF as it's being used by pdfium. This is a gatweway to FPDF_LoadCustomDocument.
	NewDocumentFromReader(reader io.ReadSeeker, size int, opts ...NewDocumentOption) (*references.FPDF_DOCUMENT, error)

	// OpenDocument returns a pdfium references for the given file data.
	// This is a gateway to FPDF_LoadMemDocument, FPDF_LoadDocument and FPDF_LoadCustomDocument.
	OpenDocument(request *requests.OpenDocument) (*responses.OpenDocument, error)

	// Close closes the instance.
	// It will close any unclosed documents.
	// For multi-threaded it will give back the worker to the pool.
	Close() error

	// End instance functions.

	// Start text: text helpers

	// GetPageText returns the text of a given page in plain text.
	GetPageText(request *requests.GetPageText) (*responses.GetPageText, error)

	// GetPageTextStructured returns the text of a given page in a structured way,
	// with coordinates and font information.
	GetPageTextStructured(request *requests.GetPageTextStructured) (*responses.GetPageTextStructured, error)

	// End text: text helpers

	// Start text: metadata helpers

	// GetMetaData returns the metadata values of the document.
	GetMetaData(request *requests.GetMetaData) (*responses.GetMetaData, error)

	// End text: metadata helpers

	// Start render: render helpers

	// RenderPageInDPI renders a given page in the given DPI.
	RenderPageInDPI(request *requests.RenderPageInDPI) (*responses.RenderPage, error)

	// RenderPagesInDPI renders the given pages in the given DPI.
	RenderPagesInDPI(request *requests.RenderPagesInDPI) (*responses.RenderPages, error)

	// RenderPageInPixels renders a given page in the given pixel size.
	RenderPageInPixels(request *requests.RenderPageInPixels) (*responses.RenderPage, error)

	// RenderPagesInPixels renders the given pages in the given pixel sizes.
	RenderPagesInPixels(request *requests.RenderPagesInPixels) (*responses.RenderPages, error)

	// GetPageSize returns the size of the page in points.
	GetPageSize(request *requests.GetPageSize) (*responses.GetPageSize, error)

	// GetPageSizeInPixels returns the size of a page in pixels when rendered in the given DPI.
	GetPageSizeInPixels(request *requests.GetPageSizeInPixels) (*responses.GetPageSizeInPixels, error)

	// RenderToFile allows you to call one of the other render functions
	// and output the resulting image into a file.
	RenderToFile(request *requests.RenderToFile) (*responses.RenderToFile, error)

	// End render

	// Start fpdfview.h

	// FPDF_GetLastError returns the last error code of a PDFium function, which is just called.
	// Usually, this function is called after a PDFium function returns, in order to check the error code of the previous PDFium function.
	FPDF_GetLastError(request *requests.FPDF_GetLastError) (*responses.FPDF_GetLastError, error)

	// FPDF_CloseDocument closes the references, releases the resources.
	FPDF_CloseDocument(request references.FPDF_DOCUMENT) error

	// FPDF_LoadPage loads a page and returns a reference.
	FPDF_LoadPage(request *requests.FPDF_LoadPage) (*responses.FPDF_LoadPage, error)

	// FPDF_ClosePage closes a page that was loaded by LoadPage.
	FPDF_ClosePage(request *requests.FPDF_ClosePage) (*responses.FPDF_ClosePage, error)

	// FPDF_GetFileVersion returns the numeric version of the file:  14 for 1.4, 15 for 1.5, ...
	FPDF_GetFileVersion(request *requests.FPDF_GetFileVersion) (*responses.FPDF_GetFileVersion, error)

	// FPDF_GetDocPermissions returns the permission flags of the file.
	FPDF_GetDocPermissions(request *requests.FPDF_GetDocPermissions) (*responses.FPDF_GetDocPermissions, error)

	// FPDF_GetSecurityHandlerRevision returns the revision number of security handlers of the file.
	FPDF_GetSecurityHandlerRevision(request *requests.FPDF_GetSecurityHandlerRevision) (*responses.FPDF_GetSecurityHandlerRevision, error)

	// FPDF_GetPageCount returns the amount of pages for the references.
	FPDF_GetPageCount(request *requests.FPDF_GetPageCount) (*responses.FPDF_GetPageCount, error)

	// End fpdfview.h

	// Start fpdf_edit.h

	// FPDF_CreateNewDocument returns a new document.
	FPDF_CreateNewDocument(request *requests.FPDF_CreateNewDocument) (*responses.FPDF_CreateNewDocument, error)

	// FPDFPage_SetRotation sets the page rotation for a given page.
	FPDFPage_SetRotation(request *requests.FPDFPage_SetRotation) (*responses.FPDFPage_SetRotation, error)

	// FPDFPage_GetRotation returns the rotation of the given page.
	FPDFPage_GetRotation(request *requests.FPDFPage_GetRotation) (*responses.FPDFPage_GetRotation, error)

	// FPDFPage_HasTransparency returns whether a page has transparency.
	FPDFPage_HasTransparency(request *requests.FPDFPage_HasTransparency) (*responses.FPDFPage_HasTransparency, error)

	// End fpdf_edit.h

	// Start fpdf_ppo.h

	// FPDF_ImportPages imports some pages from one PDF document to another one.
	FPDF_ImportPages(request *requests.FPDF_ImportPages) (*responses.FPDF_ImportPages, error)

	// FPDF_CopyViewerPreferences copies the viewer preferences from one PDF document to another
	FPDF_CopyViewerPreferences(request *requests.FPDF_CopyViewerPreferences) (*responses.FPDF_CopyViewerPreferences, error)

	// End fpdf_ppo.h

	// Start fpdf_flatten.h

	// FPDFPage_Flatten makes annotations and form fields become part of the page contents itself
	FPDFPage_Flatten(request *requests.FPDFPage_Flatten) (*responses.FPDFPage_Flatten, error)

	// End fpdf_flatten.h

	// Start fpdf_ext.h

	// FPDFDoc_GetPageMode returns the document's page mode, which describes how the references should be displayed when opened.
	FPDFDoc_GetPageMode(request *requests.FPDFDoc_GetPageMode) (*responses.FPDFDoc_GetPageMode, error)

	// End fpdf_ext.h

	// Start fpdf_doc.h

	// FPDF_GetMetaText returns the requested metadata.
	FPDF_GetMetaText(request *requests.FPDF_GetMetaText) (*responses.FPDF_GetMetaText, error)

	// End fpdf_doc.h
}
