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
	// NewDocumentFromBytes returns a pdfium Document from the given PDF bytes.
	// This is a helper around OpenDocument.
	NewDocumentFromBytes(file *[]byte, opts ...NewDocumentOption) (*references.Document, error)

	// NewDocumentFromFilePath returns a pdfium Document from the given PDF file path.
	// This is a helper around OpenDocument.
	NewDocumentFromFilePath(filePath string, opts ...NewDocumentOption) (*references.Document, error)

	// NewDocumentFromReader returns a pdfium Document from the given PDF file reader.
	// This is a helper around OpenDocument.
	// This is only really efficient for single threaded usage, the multi-threaded
	// usage will just load the file in memory because it can't transfer readers
	// over gRPC. The single-threaded usage will actually efficiently walk over
	// the PDF as it's being used by pdfium.
	NewDocumentFromReader(reader io.ReadSeeker, size int, opts ...NewDocumentOption) (*references.Document, error)

	// OpenDocument returns a pdfium references for the given file data.
	OpenDocument(request *requests.OpenDocument) (*responses.OpenDocument, error)

	// LoadPage loads a page and returns a reference.
	LoadPage(request *requests.LoadPage) (*responses.LoadPage, error)

	// UnloadPage unloads a page by reference.
	UnloadPage(request *requests.UnloadPage) (*responses.UnloadPage, error)

	// GetFileVersion returns the numeric version of the file:  14 for 1.4, 15 for 1.5, ...
	GetFileVersion(request *requests.GetFileVersion) (*responses.GetFileVersion, error)

	// GetDocPermissions returns the permission flags of the file.
	GetDocPermissions(request *requests.GetDocPermissions) (*responses.GetDocPermissions, error)

	// GetSecurityHandlerRevision returns the revision number of security handlers of the file.
	GetSecurityHandlerRevision(request *requests.GetSecurityHandlerRevision) (*responses.GetSecurityHandlerRevision, error)

	// GetPageCount returns the amount of pages for the references.
	GetPageCount(request *requests.GetPageCount) (*responses.GetPageCount, error)

	// GetPageMode references's page mode, which describes how the references should be displayed when opened.
	GetPageMode(request *requests.GetPageMode) (*responses.GetPageMode, error)

	// GetMetadata returns the requested metadata.
	GetMetadata(request *requests.GetMetadata) (*responses.GetMetadata, error)

	// GetPageText returns the text of a given page in plain text.
	GetPageText(request *requests.GetPageText) (*responses.GetPageText, error)

	// GetPageTextStructured returns the text of a given page in a structured way,
	// with coordinates and font information.
	GetPageTextStructured(request *requests.GetPageTextStructured) (*responses.GetPageTextStructured, error)

	// GetPageRotation returns the rotation of the given page.
	GetPageRotation(request *requests.GetPageRotation) (*responses.GetPageRotation, error)

	// GetPageTransparency returns whether a page has transparency.
	GetPageTransparency(request *requests.GetPageTransparency) (*responses.GetPageTransparency, error)

	// FlattenPage makes annotations and form fields become part of the page contents itself
	FlattenPage(request *requests.FlattenPage) (*responses.FlattenPage, error)

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

	// CloseDocument closes the references, releases the resources.
	CloseDocument(request references.Document) error

	// Close closes the instance.
	// It will close any unclosed documents.
	// For multi-threaded it will give back the worker to the pool.
	Close() error
}
