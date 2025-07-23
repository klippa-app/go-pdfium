package pdfium

import (
	goctx "context"
	"time"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

type LibraryConfig struct {
	UserFontPaths []string // Array of paths to scan in place of the defaults when using built-in FXGE font loading code. The Array may be nil or empty itself to use the default paths. May be ignored entirely depending upon the platform.
}

// Pool describes a PDFium worker pool. Every instance in the pool manages
// its own resources.
type Pool interface {
	// GetInstance returns an instance to the pool.
	// For single-threaded this is thread safe, but you can only do one PDFium action at the same time.
	// For multi-threaded it will try to get a worker from the pool for the length of timeout
	// It is important to Close instances when you are done with them. To either return them to the pool
	// or clear it's resources.
	GetInstance(timeout time.Duration) (Pdfium, error)

	// Same as GetInstance but with a go context.
	// If the context provided here does not contain any timeouts or cancellations, this function
	// might block forever. It is the user's responsibility to handle this context.
	GetInstanceWithContext(ctx goctx.Context) (Pdfium, error)

	// Close closes the pool.
	// It will close any unclosed instances.
	// For single-threaded it will unload the library if it's the last pool.
	// For multi-threaded it will stop all the pool workers.
	Close() error
}

// Pdfium describes a Pdfium worker instance. Documents and handles can't be
// shared between different instances. WHen a worker is closed, all resources
// and open documents are released.
type Pdfium interface {
	// Start instance functions.

	// OpenDocument returns a PDFium references for the given file data.
	// This is a gateway to FPDF_LoadMemDocument, FPDF_LoadMemDocument64,
	// FPDF_LoadDocument and FPDF_LoadCustomDocument. Please note that
	// FPDF_LoadCustomDocument will only work efficiently on single-threaded
	// usage, on multi-threaded this will just fully read from the reader
	// into a byte array before it's being sent over to PDFium.
	// This method already checks FPDF_GetLastError internally for the result.
	OpenDocument(request *requests.OpenDocument) (*responses.OpenDocument, error)

	// Close closes the instance.
	// It will close any unclosed documents.
	// For multi-threaded it will give back the worker to the pool.
	Close() error

	// Kill kills the instance.
	// On multi-thread this will kill the subprocess.
	// On single-threaded this is the same as Close().
	// Use this when you detected that your process has hung.
	Kill() error

	// GetImplementation returns the specific runtime implementation.
	GetImplementation() interface{}

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
	RenderPageInDPI(request *requests.RenderPageInDPI) (*responses.RenderPageInDPI, error)

	// RenderPagesInDPI renders the given pages in the given DPI.
	RenderPagesInDPI(request *requests.RenderPagesInDPI) (*responses.RenderPagesInDPI, error)

	// RenderPageInPixels renders a given page in the given pixel size.
	RenderPageInPixels(request *requests.RenderPageInPixels) (*responses.RenderPageInPixels, error)

	// RenderPagesInPixels renders the given pages in the given pixel sizes.
	RenderPagesInPixels(request *requests.RenderPagesInPixels) (*responses.RenderPagesInPixels, error)

	// GetPageSize returns the size of the page in points.
	GetPageSize(request *requests.GetPageSize) (*responses.GetPageSize, error)

	// GetPageSizeInPixels returns the size of a page in pixels when rendered in the given DPI.
	GetPageSizeInPixels(request *requests.GetPageSizeInPixels) (*responses.GetPageSizeInPixels, error)

	// RenderToFile allows you to call one of the other render functions
	// and output the resulting image into a file.
	RenderToFile(request *requests.RenderToFile) (*responses.RenderToFile, error)

	// End render

	// Start bookmark: bookmark helpers

	// GetBookmarks returns all the bookmarks of a document.
	GetBookmarks(request *requests.GetBookmarks) (*responses.GetBookmarks, error)

	// End bookmark

	// Start action: action helpers

	// GetActionInfo returns all the information of an action.
	GetActionInfo(request *requests.GetActionInfo) (*responses.GetActionInfo, error)

	// End action

	// Start action: dest helpers

	// GetDestInfo returns all the information of a dest.
	GetDestInfo(request *requests.GetDestInfo) (*responses.GetDestInfo, error)

	// End dest

	// Start fpdfview.h

	// FPDF_LoadDocument opens and load a PDF document from a file path.
	// Loaded document can be closed by FPDF_CloseDocument().
	// This method already checks FPDF_GetLastError internally for the result.
	FPDF_LoadDocument(request *requests.FPDF_LoadDocument) (*responses.FPDF_LoadDocument, error)

	// FPDF_LoadMemDocument opens and load a PDF document from memory.
	// Loaded document can be closed by FPDF_CloseDocument().
	// This method already checks FPDF_GetLastError internally for the result.
	FPDF_LoadMemDocument(request *requests.FPDF_LoadMemDocument) (*responses.FPDF_LoadMemDocument, error)

	// FPDF_LoadMemDocument64 opens and load a PDF document from memory.
	// Loaded document can be closed by FPDF_CloseDocument().
	// This method already checks FPDF_GetLastError internally for the result.
	// Experimental API.
	FPDF_LoadMemDocument64(request *requests.FPDF_LoadMemDocument64) (*responses.FPDF_LoadMemDocument64, error)

	// FPDF_LoadCustomDocument loads a PDF document from a custom access descriptor.
	// This is implemented as an io.ReadSeeker in go-pdfium.
	// This is only really efficient for single threaded usage, the multi-threaded
	// usage will just load the file in memory because it can't transfer readers
	// over gRPC. The single-threaded usage will actually efficiently walk over
	// the PDF as it's being used by PDFium.
	// Loaded document can be closed by FPDF_CloseDocument().
	// This method already checks FPDF_GetLastError internally for the result.
	FPDF_LoadCustomDocument(request *requests.FPDF_LoadCustomDocument) (*responses.FPDF_LoadCustomDocument, error)

	// FPDF_CloseDocument closes the references, releases the resources.
	FPDF_CloseDocument(request *requests.FPDF_CloseDocument) (*responses.FPDF_CloseDocument, error)

	// FPDF_GetLastError returns the last error code of a PDFium function, which is just called.
	// Usually, this function is called after a PDFium function returns, in order to check the error code of the previous PDFium function.
	// If the previous SDK call succeeded, the return value of this function is not defined. This function only works in conjunction
	// with APIs that mention FPDF_GetLastError() in their documentation.
	// Please note that when using go-pdfium from the same instance (on single-threaded any instance)
	// from different subroutines, FPDF_GetLastError might already be reset from
	// executing another PDFium method.
	FPDF_GetLastError(request *requests.FPDF_GetLastError) (*responses.FPDF_GetLastError, error)

	// FPDF_SetSandBoxPolicy set the policy for the sandbox environment.
	FPDF_SetSandBoxPolicy(request *requests.FPDF_SetSandBoxPolicy) (*responses.FPDF_SetSandBoxPolicy, error)

	// FPDF_LoadPage loads a page and returns a reference.
	FPDF_LoadPage(request *requests.FPDF_LoadPage) (*responses.FPDF_LoadPage, error)

	// FPDF_ClosePage closes a page that was loaded by LoadPage.
	FPDF_ClosePage(request *requests.FPDF_ClosePage) (*responses.FPDF_ClosePage, error)

	// FPDF_GetFileVersion returns the numeric version of the file:  14 for 1.4, 15 for 1.5, ...
	FPDF_GetFileVersion(request *requests.FPDF_GetFileVersion) (*responses.FPDF_GetFileVersion, error)

	// FPDF_GetDocPermissions returns the permission flags of the file.
	FPDF_GetDocPermissions(request *requests.FPDF_GetDocPermissions) (*responses.FPDF_GetDocPermissions, error)

	// FPDF_GetDocUserPermissions returns the user permission flags of the file.
	// Always returns user permissions, even if the document was unlocked by the owner.
	// Experimental API.
	FPDF_GetDocUserPermissions(request *requests.FPDF_GetDocUserPermissions) (*responses.FPDF_GetDocUserPermissions, error)

	// FPDF_GetSecurityHandlerRevision returns the revision number of security handlers of the file.
	FPDF_GetSecurityHandlerRevision(request *requests.FPDF_GetSecurityHandlerRevision) (*responses.FPDF_GetSecurityHandlerRevision, error)

	// FPDF_GetPageCount returns the amount of pages for the references.
	FPDF_GetPageCount(request *requests.FPDF_GetPageCount) (*responses.FPDF_GetPageCount, error)

	// FPDF_GetPageWidth returns the width of a page.
	// Prefer FPDF_GetPageWidthF(). This will be deprecated in the future.
	FPDF_GetPageWidth(request *requests.FPDF_GetPageWidth) (*responses.FPDF_GetPageWidth, error)

	// FPDF_GetPageHeight returns the height of a page.
	// Prefer FPDF_GetPageHeightF(). This will be deprecated in the future.
	FPDF_GetPageHeight(request *requests.FPDF_GetPageHeight) (*responses.FPDF_GetPageHeight, error)

	// FPDF_GetPageSizeByIndex returns the size of a page by the page index.
	FPDF_GetPageSizeByIndex(request *requests.FPDF_GetPageSizeByIndex) (*responses.FPDF_GetPageSizeByIndex, error)

	// FPDF_DocumentHasValidCrossReferenceTable returns whether the document's cross reference table is valid or not.
	// Experimental API.
	FPDF_DocumentHasValidCrossReferenceTable(request *requests.FPDF_DocumentHasValidCrossReferenceTable) (*responses.FPDF_DocumentHasValidCrossReferenceTable, error)

	// FPDF_GetTrailerEnds returns the byte offsets of trailer ends.
	// Experimental API.
	FPDF_GetTrailerEnds(request *requests.FPDF_GetTrailerEnds) (*responses.FPDF_GetTrailerEnds, error)

	// FPDF_GetPageWidthF returns the page width in float32.
	// Experimental API.
	FPDF_GetPageWidthF(request *requests.FPDF_GetPageWidthF) (*responses.FPDF_GetPageWidthF, error)

	// FPDF_GetPageHeightF returns the page height in float32.
	// Experimental API.
	FPDF_GetPageHeightF(request *requests.FPDF_GetPageHeightF) (*responses.FPDF_GetPageHeightF, error)

	// FPDF_GetPageBoundingBox returns the bounding box of the page. This is the intersection between
	// its media box and its crop box.
	// Experimental API.
	FPDF_GetPageBoundingBox(request *requests.FPDF_GetPageBoundingBox) (*responses.FPDF_GetPageBoundingBox, error)

	// FPDF_GetPageSizeByIndexF returns the size of the page at the given index.
	// Prefer FPDF_GetPageSizeByIndexF(). This will be deprecated in the future.
	// Experimental API.
	FPDF_GetPageSizeByIndexF(request *requests.FPDF_GetPageSizeByIndexF) (*responses.FPDF_GetPageSizeByIndexF, error)

	// FPDF_RenderPageBitmap renders contents of a page to a device independent bitmap.
	FPDF_RenderPageBitmap(request *requests.FPDF_RenderPageBitmap) (*responses.FPDF_RenderPageBitmap, error)

	// FPDF_RenderPageBitmapWithMatrix renders contents of a page to a device independent bitmap.
	FPDF_RenderPageBitmapWithMatrix(request *requests.FPDF_RenderPageBitmapWithMatrix) (*responses.FPDF_RenderPageBitmapWithMatrix, error)

	// FPDF_DeviceToPage converts the screen coordinates of a point to page coordinates.
	// The page coordinate system has its origin at the left-bottom corner
	// of the page, with the X-axis on the bottom going to the right, and
	// the Y-axis on the left side going up.
	//
	// NOTE: this coordinate system can be altered when you zoom, scroll,
	// or rotate a page, however, a point on the page should always have
	// the same coordinate values in the page coordinate system.
	//
	// The device coordinate system is device dependent. For screen device,
	// its origin is at the left-top corner of the window. However this
	// origin can be altered by the Windows coordinate transformation
	// utilities.
	//
	// You must make sure the start_x, start_y, size_x, size_y
	// and rotate parameters have exactly same values as you used in
	// the FPDF_RenderPage() function call.
	FPDF_DeviceToPage(request *requests.FPDF_DeviceToPage) (*responses.FPDF_DeviceToPage, error)

	// FPDF_PageToDevice converts the page coordinates of a point to screen coordinates.
	// See comments for FPDF_DeviceToPage().
	FPDF_PageToDevice(request *requests.FPDF_PageToDevice) (*responses.FPDF_PageToDevice, error)

	// FPDFBitmap_Create Create a device independent bitmap (FXDIB).
	FPDFBitmap_Create(request *requests.FPDFBitmap_Create) (*responses.FPDFBitmap_Create, error)

	// FPDFBitmap_CreateEx Create a device independent bitmap (FXDIB) with an
	// external buffer.
	// Similar to FPDFBitmap_Create function, but allows for more formats
	// and an external buffer is supported. The bitmap created by this
	// function can be used in any place that a FPDF_BITMAP handle is
	// required.
	//
	// If an external buffer is used, then the caller should destroy the
	// buffer. FPDFBitmap_Destroy() will not destroy the buffer.
	//
	// It is recommended to use FPDFBitmap_GetStride() to get the stride
	// value.
	//
	// Not supported on multi-threaded usage.
	FPDFBitmap_CreateEx(request *requests.FPDFBitmap_CreateEx) (*responses.FPDFBitmap_CreateEx, error)

	// FPDFBitmap_GetFormat returns the format of the bitmap.
	// Only formats supported by FPDFBitmap_CreateEx are supported by this function.
	FPDFBitmap_GetFormat(request *requests.FPDFBitmap_GetFormat) (*responses.FPDFBitmap_GetFormat, error)

	// FPDFBitmap_FillRect fills a rectangle in a bitmap.
	// This function sets the color and (optionally) alpha value in the
	// specified region of the bitmap.
	//
	// NOTE: If the alpha channel is used, this function does NOT
	// composite the background with the source color, instead the
	// background will be replaced by the source color and the alpha.
	//
	// If the alpha channel is not used, the alpha parameter is ignored.
	FPDFBitmap_FillRect(request *requests.FPDFBitmap_FillRect) (*responses.FPDFBitmap_FillRect, error)

	// FPDFBitmap_GetBuffer returns the data buffer of a bitmap.
	// The stride may be more than width * number of bytes per pixel
	//
	// Applications can use this function to get the bitmap buffer pointer,
	// then manipulate any color and/or alpha values for any pixels in the
	// bitmap.
	//
	// Use FPDFBitmap_GetFormat() to find out the format of the data.
	FPDFBitmap_GetBuffer(request *requests.FPDFBitmap_GetBuffer) (*responses.FPDFBitmap_GetBuffer, error)

	// FPDFBitmap_GetWidth returns the width of a bitmap.
	FPDFBitmap_GetWidth(request *requests.FPDFBitmap_GetWidth) (*responses.FPDFBitmap_GetWidth, error)

	// FPDFBitmap_GetHeight returns the height of a bitmap.
	FPDFBitmap_GetHeight(request *requests.FPDFBitmap_GetHeight) (*responses.FPDFBitmap_GetHeight, error)

	// FPDFBitmap_GetStride returns the number of bytes for each line in the bitmap buffer.
	FPDFBitmap_GetStride(request *requests.FPDFBitmap_GetStride) (*responses.FPDFBitmap_GetStride, error)

	// FPDFBitmap_Destroy destroys a bitmap and release all related buffers.
	// This function will not destroy any external buffers provided when
	// the bitmap was created.
	FPDFBitmap_Destroy(request *requests.FPDFBitmap_Destroy) (*responses.FPDFBitmap_Destroy, error)

	// FPDF_VIEWERREF_GetPrintScaling returns whether the PDF document prefers to be scaled or not.
	FPDF_VIEWERREF_GetPrintScaling(request *requests.FPDF_VIEWERREF_GetPrintScaling) (*responses.FPDF_VIEWERREF_GetPrintScaling, error)

	// FPDF_VIEWERREF_GetNumCopies returns the number of copies to be printed.
	FPDF_VIEWERREF_GetNumCopies(request *requests.FPDF_VIEWERREF_GetNumCopies) (*responses.FPDF_VIEWERREF_GetNumCopies, error)

	// FPDF_VIEWERREF_GetPrintPageRange returns the page numbers to initialize print dialog box when file is printed.
	FPDF_VIEWERREF_GetPrintPageRange(request *requests.FPDF_VIEWERREF_GetPrintPageRange) (*responses.FPDF_VIEWERREF_GetPrintPageRange, error)

	// FPDF_VIEWERREF_GetPrintPageRangeCount returns the number of elements in a FPDF_PAGERANGE.
	// Experimental API.
	FPDF_VIEWERREF_GetPrintPageRangeCount(request *requests.FPDF_VIEWERREF_GetPrintPageRangeCount) (*responses.FPDF_VIEWERREF_GetPrintPageRangeCount, error)

	// FPDF_VIEWERREF_GetPrintPageRangeElement returns an element from a FPDF_PAGERANGE.
	// Experimental API.
	FPDF_VIEWERREF_GetPrintPageRangeElement(request *requests.FPDF_VIEWERREF_GetPrintPageRangeElement) (*responses.FPDF_VIEWERREF_GetPrintPageRangeElement, error)

	// FPDF_VIEWERREF_GetDuplex returns the paper handling option to be used when printing from the print dialog.
	FPDF_VIEWERREF_GetDuplex(request *requests.FPDF_VIEWERREF_GetDuplex) (*responses.FPDF_VIEWERREF_GetDuplex, error)

	// FPDF_VIEWERREF_GetName returns the contents for a viewer ref, with a given key. The value must be of type "name".
	FPDF_VIEWERREF_GetName(request *requests.FPDF_VIEWERREF_GetName) (*responses.FPDF_VIEWERREF_GetName, error)

	// FPDF_CountNamedDests returns the count of named destinations in the PDF document.
	FPDF_CountNamedDests(request *requests.FPDF_CountNamedDests) (*responses.FPDF_CountNamedDests, error)

	// FPDF_GetNamedDestByName returns the destination handle for the given name.
	FPDF_GetNamedDestByName(request *requests.FPDF_GetNamedDestByName) (*responses.FPDF_GetNamedDestByName, error)

	// FPDF_GetNamedDest returns the named destination by index.
	FPDF_GetNamedDest(request *requests.FPDF_GetNamedDest) (*responses.FPDF_GetNamedDest, error)

	// FPDF_GetXFAPacketCount returns the number of valid packets in the XFA entry.
	// Experimental API.
	FPDF_GetXFAPacketCount(request *requests.FPDF_GetXFAPacketCount) (*responses.FPDF_GetXFAPacketCount, error)

	// FPDF_GetXFAPacketName returns the name of a packet in the XFA array.
	// Experimental API.
	FPDF_GetXFAPacketName(request *requests.FPDF_GetXFAPacketName) (*responses.FPDF_GetXFAPacketName, error)

	// FPDF_GetXFAPacketContent returns the content of a packet in the XFA array.
	FPDF_GetXFAPacketContent(request *requests.FPDF_GetXFAPacketContent) (*responses.FPDF_GetXFAPacketContent, error)

	// FPDF_SetPrintMode sets printing mode when printing on Windows.
	// Experimental API.
	// Windows only!
	FPDF_SetPrintMode(request *requests.FPDF_SetPrintMode) (*responses.FPDF_SetPrintMode, error)

	// FPDF_RenderPage renders contents of a page to a device (screen, bitmap, or printer).
	// This feature does not work on multi-threaded usage as you will need to give a device handle.
	// Windows only!
	FPDF_RenderPage(request *requests.FPDF_RenderPage) (*responses.FPDF_RenderPage, error)

	// End fpdfview.h

	// Start fpdf_edit.h

	// FPDF_CreateNewDocument returns a new document.
	FPDF_CreateNewDocument(request *requests.FPDF_CreateNewDocument) (*responses.FPDF_CreateNewDocument, error)

	// FPDFPage_New creates a new PDF page.
	// The page should be closed with FPDF_ClosePage() when finished as
	// with any other page in the document.
	FPDFPage_New(request *requests.FPDFPage_New) (*responses.FPDFPage_New, error)

	// FPDFPage_Delete deletes the page at the given index.
	FPDFPage_Delete(request *requests.FPDFPage_Delete) (*responses.FPDFPage_Delete, error)

	// FPDF_MovePages Move the given pages to a new index position.
	// When this call fails, the document may be left in an indeterminate state.
	// Experimental API.
	FPDF_MovePages(request *requests.FPDF_MovePages) (*responses.FPDF_MovePages, error)

	// FPDFPage_SetRotation sets the page rotation for a given page.
	FPDFPage_SetRotation(request *requests.FPDFPage_SetRotation) (*responses.FPDFPage_SetRotation, error)

	// FPDFPage_GetRotation returns the rotation of the given page.
	FPDFPage_GetRotation(request *requests.FPDFPage_GetRotation) (*responses.FPDFPage_GetRotation, error)

	// FPDFPage_InsertObject inserts the given object into a page.
	FPDFPage_InsertObject(request *requests.FPDFPage_InsertObject) (*responses.FPDFPage_InsertObject, error)

	// FPDFPage_InsertObjectAtIndex inserts the given object into a page at a specific index.
	// While technically this is not an experimental API function, in go-pdfium
	// this has been implemented as an experimental API function to ensure
	// backwards compatibility to older pdfium versions.
	FPDFPage_InsertObjectAtIndex(request *requests.FPDFPage_InsertObjectAtIndex) (*responses.FPDFPage_InsertObjectAtIndex, error)

	// FPDFPage_RemoveObject removes an object from a page.
	// Ownership is transferred to the caller. Call FPDFPageObj_Destroy() to free
	// it.
	// Experimental API.
	FPDFPage_RemoveObject(request *requests.FPDFPage_RemoveObject) (*responses.FPDFPage_RemoveObject, error)

	// FPDFPage_CountObjects returns the number of page objects inside the given page.
	FPDFPage_CountObjects(request *requests.FPDFPage_CountObjects) (*responses.FPDFPage_CountObjects, error)

	// FPDFPage_GetObject returns the object at the given index.
	FPDFPage_GetObject(request *requests.FPDFPage_GetObject) (*responses.FPDFPage_GetObject, error)

	// FPDFPage_HasTransparency returns whether a page has transparency.
	FPDFPage_HasTransparency(request *requests.FPDFPage_HasTransparency) (*responses.FPDFPage_HasTransparency, error)

	// FPDFPage_GenerateContent generates the contents of the page.
	FPDFPage_GenerateContent(request *requests.FPDFPage_GenerateContent) (*responses.FPDFPage_GenerateContent, error)

	// FPDFPageObj_Destroy destroys the page object by releasing its resources. The page object must have been
	// created by FPDFPageObj_CreateNew{Path|Rect}() or
	// FPDFPageObj_New{Text|Image}Obj(). This function must be called on
	// newly-created objects if they are not added to a page through
	// FPDFPage_InsertObject() or to an annotation through FPDFAnnot_AppendObject().
	FPDFPageObj_Destroy(request *requests.FPDFPageObj_Destroy) (*responses.FPDFPageObj_Destroy, error)

	// FPDFPageObj_HasTransparency returns whether the given page object contains transparency.
	FPDFPageObj_HasTransparency(request *requests.FPDFPageObj_HasTransparency) (*responses.FPDFPageObj_HasTransparency, error)

	// FPDFPageObj_GetType returns the type of the given page object.
	FPDFPageObj_GetType(request *requests.FPDFPageObj_GetType) (*responses.FPDFPageObj_GetType, error)

	// FPDFPageObj_GetIsActive returns the active state for the given page
	// object within the page.
	// For page objects where active is filled with false, the page object is
	// treated as if it wasn't in the document even though it is still held
	// internally.
	// Experimental API.
	FPDFPageObj_GetIsActive(request *requests.FPDFPageObj_GetIsActive) (*responses.FPDFPageObj_GetIsActive, error)

	// FPDFPageObj_SetIsActive sets the active state for the given page object
	// within the page.
	// Page objects all start in the active state by default, and remain in that
	// state unless this function is called.
	// When active is false, this makes the page_object be treated as if it
	// wasn't in the document even though it is still held internally.
	// Experimental API.
	FPDFPageObj_SetIsActive(request *requests.FPDFPageObj_SetIsActive) (*responses.FPDFPageObj_SetIsActive, error)

	// FPDFPageObj_Transform transforms the page object by the given matrix.
	// The matrix is composed as:
	//   |a c e|
	//   |b d f|
	// and can be used to scale, rotate, shear and translate the page object.
	FPDFPageObj_Transform(request *requests.FPDFPageObj_Transform) (*responses.FPDFPageObj_Transform, error)

	// FPDFPageObj_TransformF transforms the page object by the given matrix.
	// The matrix is composed as:
	//   |a c e|
	//   |b d f|
	// and can be used to scale, rotate, shear and translate the page object.
	// Experimental API.
	FPDFPageObj_TransformF(request *requests.FPDFPageObj_TransformF) (*responses.FPDFPageObj_TransformF, error)

	// FPDFPageObj_GetMatrix returns the transform matrix of a page object.
	// The matrix is composed as:
	//   |a c e|
	//   |b d f|
	// and can be used to scale, rotate, shear and translate the page object.
	// Experimental API.
	FPDFPageObj_GetMatrix(request *requests.FPDFPageObj_GetMatrix) (*responses.FPDFPageObj_GetMatrix, error)

	// FPDFPageObj_SetMatrix sets the transform matrix on a page object.
	// The matrix is composed as:
	//   |a c e|
	//   |b d f|
	// and can be used to scale, rotate, shear and translate the page object.
	// Experimental API.
	FPDFPageObj_SetMatrix(request *requests.FPDFPageObj_SetMatrix) (*responses.FPDFPageObj_SetMatrix, error)

	// FPDFPage_TransformAnnots transforms all annotations in the given page.
	// The matrix is composed as:
	//   |a c e|
	//   |b d f|
	// and can be used to scale, rotate, shear and translate the page annotations.
	FPDFPage_TransformAnnots(request *requests.FPDFPage_TransformAnnots) (*responses.FPDFPage_TransformAnnots, error)

	// FPDFPageObj_NewImageObj creates a new image object.
	FPDFPageObj_NewImageObj(request *requests.FPDFPageObj_NewImageObj) (*responses.FPDFPageObj_NewImageObj, error)

	// FPDFPageObj_GetMarkedContentID returns the marked content ID of a page object.
	// Experimental API.
	FPDFPageObj_GetMarkedContentID(request *requests.FPDFPageObj_GetMarkedContentID) (*responses.FPDFPageObj_GetMarkedContentID, error)

	// FPDFPageObj_CountMarks returns the count of content marks in a page object.
	// Experimental API.
	FPDFPageObj_CountMarks(request *requests.FPDFPageObj_CountMarks) (*responses.FPDFPageObj_CountMarks, error)

	// FPDFPageObj_GetMark returns the content mark of a page object at the given index.
	// Experimental API.
	FPDFPageObj_GetMark(request *requests.FPDFPageObj_GetMark) (*responses.FPDFPageObj_GetMark, error)

	// FPDFPageObj_AddMark adds a new content mark to the given page object.
	// Experimental API.
	FPDFPageObj_AddMark(request *requests.FPDFPageObj_AddMark) (*responses.FPDFPageObj_AddMark, error)

	// FPDFPageObj_RemoveMark removes the given content mark from the given page object.
	// Experimental API.
	FPDFPageObj_RemoveMark(request *requests.FPDFPageObj_RemoveMark) (*responses.FPDFPageObj_RemoveMark, error)

	// FPDFPageObjMark_GetName returns the name of a content mark.
	// Experimental API.
	FPDFPageObjMark_GetName(request *requests.FPDFPageObjMark_GetName) (*responses.FPDFPageObjMark_GetName, error)

	// FPDFPageObjMark_CountParams returns the number of key/value pair parameters in the given mark.
	// Experimental API.
	FPDFPageObjMark_CountParams(request *requests.FPDFPageObjMark_CountParams) (*responses.FPDFPageObjMark_CountParams, error)

	// FPDFPageObjMark_GetParamKey returns the key of a property in a content mark.
	// Experimental API.
	FPDFPageObjMark_GetParamKey(request *requests.FPDFPageObjMark_GetParamKey) (*responses.FPDFPageObjMark_GetParamKey, error)

	// FPDFPageObjMark_GetParamValueType returns the type of the value of a property in a content mark by key.
	// Experimental API.
	FPDFPageObjMark_GetParamValueType(request *requests.FPDFPageObjMark_GetParamValueType) (*responses.FPDFPageObjMark_GetParamValueType, error)

	// FPDFPageObjMark_GetParamIntValue returns the value of a number property in a content mark by key as int.
	// FPDFPageObjMark_GetParamValueType() should have returned FPDF_OBJECT_NUMBER
	// for this property.
	// Experimental API.
	FPDFPageObjMark_GetParamIntValue(request *requests.FPDFPageObjMark_GetParamIntValue) (*responses.FPDFPageObjMark_GetParamIntValue, error)

	// FPDFPageObjMark_GetParamStringValue returns the value of a string property in a content mark by key.
	// Experimental API.
	FPDFPageObjMark_GetParamStringValue(request *requests.FPDFPageObjMark_GetParamStringValue) (*responses.FPDFPageObjMark_GetParamStringValue, error)

	// FPDFPageObjMark_GetParamBlobValue returns the value of a blob property in a content mark by key.
	// Experimental API.
	FPDFPageObjMark_GetParamBlobValue(request *requests.FPDFPageObjMark_GetParamBlobValue) (*responses.FPDFPageObjMark_GetParamBlobValue, error)

	// FPDFPageObjMark_SetIntParam sets the value of an int property in a content mark by key. If a parameter
	// with the given key exists, its value is set to the given value. Otherwise, it is added as
	// a new parameter.
	// Experimental API.
	FPDFPageObjMark_SetIntParam(request *requests.FPDFPageObjMark_SetIntParam) (*responses.FPDFPageObjMark_SetIntParam, error)

	// FPDFPageObjMark_SetStringParam sets the value of a string property in a content mark by key. If a parameter
	// with the given key exists, its value is set to the given value. Otherwise, it is added as
	// a new parameter.
	// Experimental API.
	FPDFPageObjMark_SetStringParam(request *requests.FPDFPageObjMark_SetStringParam) (*responses.FPDFPageObjMark_SetStringParam, error)

	// FPDFPageObjMark_SetBlobParam sets the value of a blob property in a content mark by key. If a parameter
	// with the given key exists, its value is set to the given value. Otherwise, it is added as
	// a new parameter.
	// Experimental API.
	FPDFPageObjMark_SetBlobParam(request *requests.FPDFPageObjMark_SetBlobParam) (*responses.FPDFPageObjMark_SetBlobParam, error)

	// FPDFPageObjMark_RemoveParam removes a property from a content mark by key.
	// Experimental API.
	FPDFPageObjMark_RemoveParam(request *requests.FPDFPageObjMark_RemoveParam) (*responses.FPDFPageObjMark_RemoveParam, error)

	// FPDFImageObj_LoadJpegFile loads an image from a JPEG image file and then set it into the given image object.
	// The image object might already have an associated image, which is shared and
	// cached by the loaded pages. In that case, we need to clear the cached image
	// for all the loaded pages. Pass the pages and page count to this API
	// to clear the image cache. If the image is not previously shared, nil is a
	// valid pages value.
	FPDFImageObj_LoadJpegFile(request *requests.FPDFImageObj_LoadJpegFile) (*responses.FPDFImageObj_LoadJpegFile, error)

	// FPDFImageObj_LoadJpegFileInline
	// The image object might already have an associated image, which is shared and
	// cached by the loaded pages. In that case, we need to clear the cached image
	// for all the loaded pages. Pass the pages and page count to this API
	// to clear the image cache. If the image is not previously shared, nil is a
	// valid pages value. This function loads the JPEG image inline, so the image
	// content is copied to the file. This allows the file access and its associated
	// data to be deleted after this function returns.
	FPDFImageObj_LoadJpegFileInline(request *requests.FPDFImageObj_LoadJpegFileInline) (*responses.FPDFImageObj_LoadJpegFileInline, error)

	// FPDFImageObj_SetMatrix sets the transform matrix of the given image object.
	// The matrix is composed as:
	//   |a c e|
	//   |b d f|
	// and can be used to scale, rotate, shear and translate the image object.
	// Will be deprecated once FPDFPageObj_SetMatrix() is stable.
	FPDFImageObj_SetMatrix(request *requests.FPDFImageObj_SetMatrix) (*responses.FPDFImageObj_SetMatrix, error)

	// FPDFImageObj_SetBitmap sets the given bitmap to the given image object.
	FPDFImageObj_SetBitmap(request *requests.FPDFImageObj_SetBitmap) (*responses.FPDFImageObj_SetBitmap, error)

	// FPDFImageObj_GetBitmap returns a bitmap rasterization of the given image object. FPDFImageObj_GetBitmap() only
	// operates on the image object and does not take the associated image mask into
	// account. It also ignores the matrix for the image object.
	// The returned bitmap will be owned by the caller, and FPDFBitmap_Destroy()
	// must be called on the returned bitmap when it is no longer needed.
	FPDFImageObj_GetBitmap(request *requests.FPDFImageObj_GetBitmap) (*responses.FPDFImageObj_GetBitmap, error)

	// FPDFImageObj_GetRenderedBitmap returns a bitmap rasterization of the given image object that takes the image mask and
	// image matrix into account. To render correctly, the caller must provide the
	// document associated with the image object. If there is a page associated
	// with the image object the caller should provide that as well.
	// The returned bitmap will be owned by the caller, and FPDFBitmap_Destroy()
	// must be called on the returned bitmap when it is no longer needed.
	// Experimental API.
	FPDFImageObj_GetRenderedBitmap(request *requests.FPDFImageObj_GetRenderedBitmap) (*responses.FPDFImageObj_GetRenderedBitmap, error)

	// FPDFImageObj_GetImageDataDecoded returns the decoded image data of the image object. The decoded data is the
	// uncompressed image data, i.e. the raw image data after having all filters
	// applied.
	FPDFImageObj_GetImageDataDecoded(request *requests.FPDFImageObj_GetImageDataDecoded) (*responses.FPDFImageObj_GetImageDataDecoded, error)

	// FPDFImageObj_GetImageDataRaw returns the raw image data of the image object. The raw data is the image data as
	// stored in the PDF without applying any filters.
	FPDFImageObj_GetImageDataRaw(request *requests.FPDFImageObj_GetImageDataRaw) (*responses.FPDFImageObj_GetImageDataRaw, error)

	// FPDFImageObj_GetImageFilterCount returns the number of filters (i.e. decoders) of the image in image object.
	FPDFImageObj_GetImageFilterCount(request *requests.FPDFImageObj_GetImageFilterCount) (*responses.FPDFImageObj_GetImageFilterCount, error)

	// FPDFImageObj_GetImageFilter returns the filter at index of the image object's list of filters. Note that the
	// filters need to be applied in order, i.e. the first filter should be applied
	// first, then the second, etc.
	FPDFImageObj_GetImageFilter(request *requests.FPDFImageObj_GetImageFilter) (*responses.FPDFImageObj_GetImageFilter, error)

	// FPDFImageObj_GetImageMetadata returns the image metadata of the image object, including dimension, DPI, bits per
	// pixel, and colorspace. If the image object is not an image object or if it
	// does not have an image, then the return value will be false. Otherwise,
	// failure to retrieve any specific parameter would result in its value being 0.
	FPDFImageObj_GetImageMetadata(request *requests.FPDFImageObj_GetImageMetadata) (*responses.FPDFImageObj_GetImageMetadata, error)

	// FPDFImageObj_GetImagePixelSize get the image size in pixels. Faster method to get only image size.
	// Experimental API.
	FPDFImageObj_GetImagePixelSize(request *requests.FPDFImageObj_GetImagePixelSize) (*responses.FPDFImageObj_GetImagePixelSize, error)

	// FPDFImageObj_GetIccProfileDataDecoded returns the ICC profile decoded
	// data of the given image object. If the image object is not an image
	// object or if it does not have an image, then the return value will
	// be nil. It also returns nil if the image object has no ICC profile.
	// Experimental API
	FPDFImageObj_GetIccProfileDataDecoded(request *requests.FPDFImageObj_GetIccProfileDataDecoded) (*responses.FPDFImageObj_GetIccProfileDataDecoded, error)

	// FPDFPageObj_CreateNewPath creates a new path object at an initial position.
	FPDFPageObj_CreateNewPath(request *requests.FPDFPageObj_CreateNewPath) (*responses.FPDFPageObj_CreateNewPath, error)

	// FPDFPageObj_CreateNewRect creates a closed path consisting of a rectangle.
	FPDFPageObj_CreateNewRect(request *requests.FPDFPageObj_CreateNewRect) (*responses.FPDFPageObj_CreateNewRect, error)

	// FPDFPageObj_GetBounds returns the bounding box of the given page object.
	FPDFPageObj_GetBounds(request *requests.FPDFPageObj_GetBounds) (*responses.FPDFPageObj_GetBounds, error)

	// FPDFPageObj_SetBlendMode sets the blend mode of the page object.
	FPDFPageObj_SetBlendMode(request *requests.FPDFPageObj_SetBlendMode) (*responses.FPDFPageObj_SetBlendMode, error)

	// FPDFPageObj_SetStrokeColor sets the stroke RGBA of a page object.
	FPDFPageObj_SetStrokeColor(request *requests.FPDFPageObj_SetStrokeColor) (*responses.FPDFPageObj_SetStrokeColor, error)

	// FPDFPageObj_GetStrokeColor returns the stroke RGBA of a page object
	FPDFPageObj_GetStrokeColor(request *requests.FPDFPageObj_GetStrokeColor) (*responses.FPDFPageObj_GetStrokeColor, error)

	// FPDFPageObj_SetStrokeWidth sets the stroke width of a page object
	FPDFPageObj_SetStrokeWidth(request *requests.FPDFPageObj_SetStrokeWidth) (*responses.FPDFPageObj_SetStrokeWidth, error)

	// FPDFPageObj_GetStrokeWidth returns the stroke width of a page object.
	FPDFPageObj_GetStrokeWidth(request *requests.FPDFPageObj_GetStrokeWidth) (*responses.FPDFPageObj_GetStrokeWidth, error)

	// FPDFPageObj_GetLineJoin returns the line join of the page object.
	FPDFPageObj_GetLineJoin(request *requests.FPDFPageObj_GetLineJoin) (*responses.FPDFPageObj_GetLineJoin, error)

	// FPDFPageObj_SetLineJoin sets the line join of the page object.
	FPDFPageObj_SetLineJoin(request *requests.FPDFPageObj_SetLineJoin) (*responses.FPDFPageObj_SetLineJoin, error)

	// FPDFPageObj_GetLineCap returns the line cap of the page object.
	FPDFPageObj_GetLineCap(request *requests.FPDFPageObj_GetLineCap) (*responses.FPDFPageObj_GetLineCap, error)

	// FPDFPageObj_SetLineCap sets the line cap of the page object.
	FPDFPageObj_SetLineCap(request *requests.FPDFPageObj_SetLineCap) (*responses.FPDFPageObj_SetLineCap, error)

	// FPDFPageObj_SetFillColor sets the fill RGBA of a page object
	FPDFPageObj_SetFillColor(request *requests.FPDFPageObj_SetFillColor) (*responses.FPDFPageObj_SetFillColor, error)

	// FPDFPageObj_GetFillColor returns the fill RGBA of a page object
	FPDFPageObj_GetFillColor(request *requests.FPDFPageObj_GetFillColor) (*responses.FPDFPageObj_GetFillColor, error)

	// FPDFPageObj_GetRotatedBounds Get the quad points that bounds the page object.
	// Similar to FPDFPageObj_GetBounds(), this returns the bounds of a page
	// object. When the object is rotated by a non-multiple of 90 degrees, this API
	// returns a tighter bound that cannot be represented with just the 4 sides of
	// a rectangle.
	//
	// Currently only works the following page object types: FPDF_PAGEOBJ_TEXT and
	// FPDF_PAGEOBJ_IMAGE.
	// Experimental API.
	FPDFPageObj_GetRotatedBounds(request *requests.FPDFPageObj_GetRotatedBounds) (*responses.FPDFPageObj_GetRotatedBounds, error)

	// FPDFPageObj_GetDashPhase returns the line dash phase of the page object.
	// Experimental API.
	FPDFPageObj_GetDashPhase(request *requests.FPDFPageObj_GetDashPhase) (*responses.FPDFPageObj_GetDashPhase, error)

	// FPDFPageObj_SetDashPhase sets the line dash phase of the page object.
	// Experimental API.
	FPDFPageObj_SetDashPhase(request *requests.FPDFPageObj_SetDashPhase) (*responses.FPDFPageObj_SetDashPhase, error)

	// FPDFPageObj_GetDashCount returns the line dash array size of the page object.
	// Experimental API.
	FPDFPageObj_GetDashCount(request *requests.FPDFPageObj_GetDashCount) (*responses.FPDFPageObj_GetDashCount, error)

	// FPDFPageObj_GetDashArray returns the line dash array of the page object.
	// Experimental API.
	FPDFPageObj_GetDashArray(request *requests.FPDFPageObj_GetDashArray) (*responses.FPDFPageObj_GetDashArray, error)

	// FPDFPageObj_SetDashArray sets the line dash array of the page object.
	// Experimental API.
	FPDFPageObj_SetDashArray(request *requests.FPDFPageObj_SetDashArray) (*responses.FPDFPageObj_SetDashArray, error)

	// FPDFPath_CountSegments returns the number of segments inside the given path.
	// A segment is a command, created by e.g. FPDFPath_MoveTo(),
	// FPDFPath_LineTo() or FPDFPath_BezierTo().
	FPDFPath_CountSegments(request *requests.FPDFPath_CountSegments) (*responses.FPDFPath_CountSegments, error)

	// FPDFPath_GetPathSegment returns the segment in the given path at the given index.
	FPDFPath_GetPathSegment(request *requests.FPDFPath_GetPathSegment) (*responses.FPDFPath_GetPathSegment, error)

	// FPDFPathSegment_GetPoint returns the coordinates of the given segment.
	FPDFPathSegment_GetPoint(request *requests.FPDFPathSegment_GetPoint) (*responses.FPDFPathSegment_GetPoint, error)

	// FPDFPathSegment_GetType returns the type of the given segment.
	FPDFPathSegment_GetType(request *requests.FPDFPathSegment_GetType) (*responses.FPDFPathSegment_GetType, error)

	// FPDFPathSegment_GetClose returns whether the segment closes the current subpath of a given path.
	FPDFPathSegment_GetClose(request *requests.FPDFPathSegment_GetClose) (*responses.FPDFPathSegment_GetClose, error)

	// FPDFPath_MoveTo moves a path's current point.
	// Note that no line will be created between the previous current point and the
	// new one.
	FPDFPath_MoveTo(request *requests.FPDFPath_MoveTo) (*responses.FPDFPath_MoveTo, error)

	// FPDFPath_LineTo adds a line between the current point and a new point in the path.
	FPDFPath_LineTo(request *requests.FPDFPath_LineTo) (*responses.FPDFPath_LineTo, error)

	// FPDFPath_BezierTo adds a cubic Bezier curve to the given path, starting at the current point.
	FPDFPath_BezierTo(request *requests.FPDFPath_BezierTo) (*responses.FPDFPath_BezierTo, error)

	// FPDFPath_Close closes the current subpath of a given path.
	FPDFPath_Close(request *requests.FPDFPath_Close) (*responses.FPDFPath_Close, error)

	// FPDFPath_SetDrawMode sets the drawing mode of a path.
	FPDFPath_SetDrawMode(request *requests.FPDFPath_SetDrawMode) (*responses.FPDFPath_SetDrawMode, error)

	// FPDFPath_GetDrawMode returns the drawing mode of a path.
	FPDFPath_GetDrawMode(request *requests.FPDFPath_GetDrawMode) (*responses.FPDFPath_GetDrawMode, error)

	// FPDFPageObj_NewTextObj creates a new text object using one of the standard PDF fonts.
	FPDFPageObj_NewTextObj(request *requests.FPDFPageObj_NewTextObj) (*responses.FPDFPageObj_NewTextObj, error)

	// FPDFText_SetText sets the text for a text object. If it had text, it will be replaced.
	FPDFText_SetText(request *requests.FPDFText_SetText) (*responses.FPDFText_SetText, error)

	// FPDFText_SetCharcodes sets the text using charcodes for a text object. If it had text, it will be
	// replaced.
	FPDFText_SetCharcodes(request *requests.FPDFText_SetCharcodes) (*responses.FPDFText_SetCharcodes, error)

	// FPDFText_LoadFont returns a font object loaded from a stream of data. The font is loaded
	// into the document. Various font data structures, such as the ToUnicode data, are auto-generated based
	// on the inputs
	// The loaded font can be closed using FPDFFont_Close().
	FPDFText_LoadFont(request *requests.FPDFText_LoadFont) (*responses.FPDFText_LoadFont, error)

	// FPDFText_LoadStandardFont loads one of the standard 14 fonts per PDF spec 1.7 page 416. The preferred
	// way of using font style is using a dash to separate the name from the style,
	// for example 'Helvetica-BoldItalic'.
	// The loaded font can be closed using FPDFFont_Close().
	// Experimental API.
	FPDFText_LoadStandardFont(request *requests.FPDFText_LoadStandardFont) (*responses.FPDFText_LoadStandardFont, error)

	// FPDFText_LoadCidType2Font returns a font object loaded from a stream of data for a type 2 CID font.
	// The font is loaded into the document. Unlike FPDFText_LoadFont(), the ToUnicode data and the CIDToGIDMap
	// data are caller provided, instead of auto-generated.
	// The loaded font can be closed using FPDFFont_Close().
	// Experimental API.
	FPDFText_LoadCidType2Font(request *requests.FPDFText_LoadCidType2Font) (*responses.FPDFText_LoadCidType2Font, error)

	// FPDFTextObj_GetFontSize returns the font size of a text object.
	FPDFTextObj_GetFontSize(request *requests.FPDFTextObj_GetFontSize) (*responses.FPDFTextObj_GetFontSize, error)

	// FPDFFont_Close closes a loaded PDF font
	FPDFFont_Close(request *requests.FPDFFont_Close) (*responses.FPDFFont_Close, error)

	// FPDFPageObj_CreateTextObj creates a new text object using a loaded font.
	FPDFPageObj_CreateTextObj(request *requests.FPDFPageObj_CreateTextObj) (*responses.FPDFPageObj_CreateTextObj, error)

	// FPDFTextObj_GetTextRenderMode returns the text rendering mode of a text object.
	FPDFTextObj_GetTextRenderMode(request *requests.FPDFTextObj_GetTextRenderMode) (*responses.FPDFTextObj_GetTextRenderMode, error)

	// FPDFTextObj_SetTextRenderMode sets the text rendering mode of a text object.
	// Experimental API.
	FPDFTextObj_SetTextRenderMode(request *requests.FPDFTextObj_SetTextRenderMode) (*responses.FPDFTextObj_SetTextRenderMode, error)

	// FPDFTextObj_GetText returns the text of a text object.
	FPDFTextObj_GetText(request *requests.FPDFTextObj_GetText) (*responses.FPDFTextObj_GetText, error)

	// FPDFTextObj_GetRenderedBitmap returns a bitmap rasterization of the given text object.
	// To render correctly, the caller must provide the document associated with the text object.
	// If there is a page associated with text object, the caller should provide that as well.
	// The returned bitmap will be owned by the caller, and FPDFBitmap_Destroy()
	// must be called on the returned bitmap when it is no longer needed.
	// Experimental API.
	FPDFTextObj_GetRenderedBitmap(request *requests.FPDFTextObj_GetRenderedBitmap) (*responses.FPDFTextObj_GetRenderedBitmap, error)

	// FPDFTextObj_GetFont returns the font of a text object.
	// Experimental API.
	FPDFTextObj_GetFont(request *requests.FPDFTextObj_GetFont) (*responses.FPDFTextObj_GetFont, error)

	// FPDFFont_GetBaseFontName returns the base font name of a font.
	// Experimental API.
	FPDFFont_GetBaseFontName(request *requests.FPDFFont_GetBaseFontName) (*responses.FPDFFont_GetBaseFontName, error)

	// FPDFFont_GetFamilyName returns the family name of a font.
	// Experimental API.
	FPDFFont_GetFamilyName(request *requests.FPDFFont_GetFamilyName) (*responses.FPDFFont_GetFamilyName, error)

	// FPDFFont_GetFontData returns the decoded data from the given font.
	// Experimental API.
	FPDFFont_GetFontData(request *requests.FPDFFont_GetFontData) (*responses.FPDFFont_GetFontData, error)

	// FPDFFont_GetIsEmbedded returns whether the given font is embedded or not.
	// Experimental API.
	FPDFFont_GetIsEmbedded(request *requests.FPDFFont_GetIsEmbedded) (*responses.FPDFFont_GetIsEmbedded, error)

	// FPDFFont_GetFlags returns the descriptor flags of a font.
	// Returns the bit flags specifying various characteristics of the font as
	// defined in ISO 32000-1:2008, table 123.
	// Experimental API.
	FPDFFont_GetFlags(request *requests.FPDFFont_GetFlags) (*responses.FPDFFont_GetFlags, error)

	// FPDFFont_GetWeight returns the font weight of a font.
	// Typical values are 400 (normal) and 700 (bold).
	// Experimental API.
	FPDFFont_GetWeight(request *requests.FPDFFont_GetWeight) (*responses.FPDFFont_GetWeight, error)

	// FPDFFont_GetItalicAngle returns the italic angle of a font.
	// The italic angle of a font is defined as degrees counterclockwise
	// from vertical. For a font that slopes to the right, this will be negative.
	// Experimental API.
	FPDFFont_GetItalicAngle(request *requests.FPDFFont_GetItalicAngle) (*responses.FPDFFont_GetItalicAngle, error)

	// FPDFFont_GetAscent returns ascent distance of a font.
	// Ascent is the maximum distance in points above the baseline reached by the
	// glyphs of the font. One point is 1/72 inch (around 0.3528 mm).
	// Experimental API.
	FPDFFont_GetAscent(request *requests.FPDFFont_GetAscent) (*responses.FPDFFont_GetAscent, error)

	// FPDFFont_GetDescent returns the descent distance of a font.
	// Descent is the maximum distance in points below the baseline reached by the
	// glyphs of the font. One point is 1/72 inch (around 0.3528 mm).
	// Experimental API.
	FPDFFont_GetDescent(request *requests.FPDFFont_GetDescent) (*responses.FPDFFont_GetDescent, error)

	// FPDFFont_GetGlyphWidth returns the width of a glyph in a font.
	// Glyph width is the distance from the end of the prior glyph to the next
	// glyph. This will be the vertical distance for vertical writing.
	// Experimental API.
	FPDFFont_GetGlyphWidth(request *requests.FPDFFont_GetGlyphWidth) (*responses.FPDFFont_GetGlyphWidth, error)

	// FPDFFont_GetGlyphPath returns the glyphpath describing how to draw a font glyph.
	// Experimental API.
	FPDFFont_GetGlyphPath(request *requests.FPDFFont_GetGlyphPath) (*responses.FPDFFont_GetGlyphPath, error)

	// FPDFGlyphPath_CountGlyphSegments returns the number of segments inside the given glyphpath.
	// Experimental API.
	FPDFGlyphPath_CountGlyphSegments(request *requests.FPDFGlyphPath_CountGlyphSegments) (*responses.FPDFGlyphPath_CountGlyphSegments, error)

	// FPDFGlyphPath_GetGlyphPathSegment returns the segment in glyphpath at the given index.
	// Experimental API.
	FPDFGlyphPath_GetGlyphPathSegment(request *requests.FPDFGlyphPath_GetGlyphPathSegment) (*responses.FPDFGlyphPath_GetGlyphPathSegment, error)

	// FPDFFormObj_CountObjects returns the number of page objects inside the given form object.
	FPDFFormObj_CountObjects(request *requests.FPDFFormObj_CountObjects) (*responses.FPDFFormObj_CountObjects, error)

	// FPDFFormObj_GetObject returns the page object in the given form object at the given index.
	FPDFFormObj_GetObject(request *requests.FPDFFormObj_GetObject) (*responses.FPDFFormObj_GetObject, error)

	// FPDFFormObj_RemoveObject removes the page object in the given form object.
	// Ownership of the removed page object is transferred to the caller, call FPDFPageObj_Destroy() on the
	// removed page_object to free it.
	// Experimental API.
	FPDFFormObj_RemoveObject(request *requests.FPDFFormObj_RemoveObject) (*responses.FPDFFormObj_RemoveObject, error)

	// End fpdf_edit.h

	// Start fpdf_ppo.h

	// FPDF_ImportPages imports some pages from one PDF document to another one.
	FPDF_ImportPages(request *requests.FPDF_ImportPages) (*responses.FPDF_ImportPages, error)

	// FPDF_CopyViewerPreferences copies the viewer preferences from one PDF document to another
	FPDF_CopyViewerPreferences(request *requests.FPDF_CopyViewerPreferences) (*responses.FPDF_CopyViewerPreferences, error)

	// FPDF_ImportPagesByIndex imports pages to a FPDF_DOCUMENT.
	// Experimental API.
	FPDF_ImportPagesByIndex(request *requests.FPDF_ImportPagesByIndex) (*responses.FPDF_ImportPagesByIndex, error)

	// FPDF_ImportNPagesToOne creates a new document from source document. The pages of source document will be
	// combined to provide NumPagesOnXAxis x NumPagesOnYAxis pages per page of the output document.
	// Experimental API.
	FPDF_ImportNPagesToOne(request *requests.FPDF_ImportNPagesToOne) (*responses.FPDF_ImportNPagesToOne, error)

	// FPDF_NewXObjectFromPage creates a template to generate form xobjects from the source document's page at
	// the given index, for use in the destination document.
	// Experimental API.
	FPDF_NewXObjectFromPage(request *requests.FPDF_NewXObjectFromPage) (*responses.FPDF_NewXObjectFromPage, error)

	// FPDF_CloseXObject closes an FPDF_XOBJECT handle created by FPDF_NewXObjectFromPage().
	// Experimental API.
	FPDF_CloseXObject(request *requests.FPDF_CloseXObject) (*responses.FPDF_CloseXObject, error)

	// FPDF_NewFormObjectFromXObject creates a new form object from an FPDF_XOBJECT object.
	// Experimental API.
	FPDF_NewFormObjectFromXObject(request *requests.FPDF_NewFormObjectFromXObject) (*responses.FPDF_NewFormObjectFromXObject, error)

	// End fpdf_ppo.h

	// Start fpdf_flatten.h

	// FPDFPage_Flatten makes annotations and form fields become part of the page contents itself
	FPDFPage_Flatten(request *requests.FPDFPage_Flatten) (*responses.FPDFPage_Flatten, error)

	// End fpdf_flatten.h

	// Start fpdf_ext.h

	// FPDFDoc_GetPageMode returns the document's page mode, which describes how the references should be displayed when opened.
	FPDFDoc_GetPageMode(request *requests.FPDFDoc_GetPageMode) (*responses.FPDFDoc_GetPageMode, error)

	// FSDK_SetUnSpObjProcessHandler set ups an unsupported object handler.
	// Since callbacks can't be transferred between processes in gRPC, you can only use this in single-threaded mode.
	FSDK_SetUnSpObjProcessHandler(request *requests.FSDK_SetUnSpObjProcessHandler) (*responses.FSDK_SetUnSpObjProcessHandler, error)

	// FSDK_SetTimeFunction sets a replacement function for calls to time().
	// This API is intended to be used only for testing, thus may cause PDFium to behave poorly in production environments.
	// Since callbacks can't be transferred between processes in gRPC, you can only use this in single-threaded mode.
	FSDK_SetTimeFunction(request *requests.FSDK_SetTimeFunction) (*responses.FSDK_SetTimeFunction, error)

	// FSDK_SetLocaltimeFunction sets a replacement function for calls to localtime().
	// This API is intended to be used only for testing, thus may cause PDFium to behave poorly in production environments.
	// Since callbacks can't be transferred between processes in gRPC, you can only use this in single-threaded mode.
	FSDK_SetLocaltimeFunction(request *requests.FSDK_SetLocaltimeFunction) (*responses.FSDK_SetLocaltimeFunction, error)

	// End fpdf_ext.h

	// Start fpdf_doc.h

	// FPDFBookmark_GetFirstChild returns the first child of a bookmark item, or the first top level bookmark item.
	// Note that another name for the bookmarks is the document outline, as
	// described in ISO 32000-1:2008, section 12.3.3.
	FPDFBookmark_GetFirstChild(request *requests.FPDFBookmark_GetFirstChild) (*responses.FPDFBookmark_GetFirstChild, error)

	// FPDFBookmark_GetNextSibling returns the next bookmark item at the same level.
	// Note that the caller is responsible for handling circular bookmark
	// references, as may arise from malformed documents.
	FPDFBookmark_GetNextSibling(request *requests.FPDFBookmark_GetNextSibling) (*responses.FPDFBookmark_GetNextSibling, error)

	// FPDFBookmark_GetTitle returns the title of a bookmark.
	FPDFBookmark_GetTitle(request *requests.FPDFBookmark_GetTitle) (*responses.FPDFBookmark_GetTitle, error)

	// FPDFBookmark_GetCount returns the number of children of a bookmark.
	// Experimental API.
	FPDFBookmark_GetCount(request *requests.FPDFBookmark_GetCount) (*responses.FPDFBookmark_GetCount, error)

	// FPDFBookmark_Find finds a bookmark in the document, using the bookmark title.
	FPDFBookmark_Find(request *requests.FPDFBookmark_Find) (*responses.FPDFBookmark_Find, error)

	// FPDFBookmark_GetDest returns the destination associated with a bookmark item.
	// If the returned destination is nil, none is associated to the bookmark item.
	FPDFBookmark_GetDest(request *requests.FPDFBookmark_GetDest) (*responses.FPDFBookmark_GetDest, error)

	// FPDFBookmark_GetAction returns the action associated with a bookmark item.
	// If this function returns a valid handle, it is valid as long as the bookmark is
	// valid.
	// If the returned action is nil, you should try FPDFBookmark_GetDest.
	FPDFBookmark_GetAction(request *requests.FPDFBookmark_GetAction) (*responses.FPDFBookmark_GetAction, error)

	// FPDFAction_GetType returns the action associated with a bookmark item.
	FPDFAction_GetType(request *requests.FPDFAction_GetType) (*responses.FPDFAction_GetType, error)

	// FPDFAction_GetDest returns the destination of a specific go-to or remote-goto action.
	// Only action with type PDF_ACTION_ACTION_GOTO and PDF_ACTION_ACTION_REMOTEGOTO can have destination data.
	// In case of remote goto action, the application should first use function FPDFAction_GetFilePath to get file path, then load that particular document, and use its document handle to call this function.
	FPDFAction_GetDest(request *requests.FPDFAction_GetDest) (*responses.FPDFAction_GetDest, error)

	// FPDFAction_GetFilePath returns the file path from a remote goto or launch action.
	// Only works on actions that have the type FPDF_ACTION_ACTION_REMOTEGOTO or FPDF_ACTION_ACTION_LAUNCH.
	FPDFAction_GetFilePath(request *requests.FPDFAction_GetFilePath) (*responses.FPDFAction_GetFilePath, error)

	// FPDFAction_GetURIPath returns the URI path from a URI action.
	FPDFAction_GetURIPath(request *requests.FPDFAction_GetURIPath) (*responses.FPDFAction_GetURIPath, error)

	// FPDFDest_GetDestPageIndex returns the page index from destination data.
	FPDFDest_GetDestPageIndex(request *requests.FPDFDest_GetDestPageIndex) (*responses.FPDFDest_GetDestPageIndex, error)

	// FPDF_GetFileIdentifier Get the file identifier defined in the trailer of a document.
	// Experimental API.
	FPDF_GetFileIdentifier(request *requests.FPDF_GetFileIdentifier) (*responses.FPDF_GetFileIdentifier, error)

	// FPDF_GetMetaText returns the requested metadata.
	FPDF_GetMetaText(request *requests.FPDF_GetMetaText) (*responses.FPDF_GetMetaText, error)

	// FPDF_GetPageLabel returns the label for the given page.
	FPDF_GetPageLabel(request *requests.FPDF_GetPageLabel) (*responses.FPDF_GetPageLabel, error)

	// FPDFDest_GetView returns the view (fit type) for a given dest.
	// Experimental API.
	FPDFDest_GetView(request *requests.FPDFDest_GetView) (*responses.FPDFDest_GetView, error)

	// FPDFDest_GetLocationInPage returns the (x, y, zoom) location of dest in the destination page, if the
	// destination is in [page /XYZ x y zoom] syntax.
	FPDFDest_GetLocationInPage(request *requests.FPDFDest_GetLocationInPage) (*responses.FPDFDest_GetLocationInPage, error)

	// FPDFLink_GetLinkAtPoint finds a link at a point on a page.
	// You can convert coordinates from screen coordinates to page coordinates using
	// FPDF_DeviceToPage().
	FPDFLink_GetLinkAtPoint(request *requests.FPDFLink_GetLinkAtPoint) (*responses.FPDFLink_GetLinkAtPoint, error)

	// FPDFLink_GetLinkZOrderAtPoint finds the Z-order of link at a point on a page.
	// You can convert coordinates from screen coordinates to page coordinates using
	// FPDF_DeviceToPage().
	FPDFLink_GetLinkZOrderAtPoint(request *requests.FPDFLink_GetLinkZOrderAtPoint) (*responses.FPDFLink_GetLinkZOrderAtPoint, error)

	// FPDFLink_GetDest returns the destination info for a link.
	FPDFLink_GetDest(request *requests.FPDFLink_GetDest) (*responses.FPDFLink_GetDest, error)

	// FPDFLink_GetAction returns the action info for a link
	// If this function returns a valid handle, it is valid as long as the link is
	// valid.
	FPDFLink_GetAction(request *requests.FPDFLink_GetAction) (*responses.FPDFLink_GetAction, error)

	// FPDFLink_Enumerate Enumerates all the link annotations in a page.
	FPDFLink_Enumerate(request *requests.FPDFLink_Enumerate) (*responses.FPDFLink_Enumerate, error)

	// FPDFLink_GetAnnot returns a FPDF_ANNOTATION object for a link.
	// Experimental API.
	FPDFLink_GetAnnot(request *requests.FPDFLink_GetAnnot) (*responses.FPDFLink_GetAnnot, error)

	// FPDFLink_GetAnnotRect returns the count of quadrilateral points to the link.
	FPDFLink_GetAnnotRect(request *requests.FPDFLink_GetAnnotRect) (*responses.FPDFLink_GetAnnotRect, error)

	// FPDFLink_CountQuadPoints returns the count of quadrilateral points to the link.
	FPDFLink_CountQuadPoints(request *requests.FPDFLink_CountQuadPoints) (*responses.FPDFLink_CountQuadPoints, error)

	// FPDFLink_GetQuadPoints returns the quadrilateral points for the specified quad index in the link.
	FPDFLink_GetQuadPoints(request *requests.FPDFLink_GetQuadPoints) (*responses.FPDFLink_GetQuadPoints, error)

	// FPDF_GetPageAAction returns an additional-action from page.
	// If this function returns a valid handle, it is valid as long as the page is
	// valid.
	// Experimental API
	FPDF_GetPageAAction(request *requests.FPDF_GetPageAAction) (*responses.FPDF_GetPageAAction, error)

	// End fpdf_doc.h

	// Start fpdf_save.h

	// FPDF_SaveAsCopy saves the document to a copy.
	// If no path or writer is given, it will return the saved file as a byte array.
	// Note that using a fileWriter only works when using the single-threaded version,
	// the reason for that is that a fileWriter can't be transferred over gRPC
	// (or between processes at all).
	FPDF_SaveAsCopy(request *requests.FPDF_SaveAsCopy) (*responses.FPDF_SaveAsCopy, error)

	// FPDF_SaveWithVersion save the document to a copy, with a specific file version.
	// If no path or writer is given, it will return the saved file as a byte array.
	// Note that using a fileWriter only works when using the single-threaded version,
	// the reason for that is that a fileWriter can't be transferred over gRPC
	// (or between processes at all).
	FPDF_SaveWithVersion(request *requests.FPDF_SaveWithVersion) (*responses.FPDF_SaveWithVersion, error)

	// End fpdf_save.h

	// Start fpdf_catalog.h

	// FPDFCatalog_IsTagged determines if the given document represents a tagged PDF.
	// For the definition of tagged PDF, See (see 10.7 "Tagged PDF" in PDF Reference 1.7).
	// Experimental API.
	FPDFCatalog_IsTagged(request *requests.FPDFCatalog_IsTagged) (*responses.FPDFCatalog_IsTagged, error)

	// FPDFCatalog_SetLanguage sets the language of a document.
	// Experimental API.
	FPDFCatalog_SetLanguage(request *requests.FPDFCatalog_SetLanguage) (*responses.FPDFCatalog_SetLanguage, error)

	// End fpdf_catalog.h

	// Start fpdf_signature.h

	// FPDF_GetSignatureCount returns the total number of signatures in the document.
	// Experimental API.
	FPDF_GetSignatureCount(request *requests.FPDF_GetSignatureCount) (*responses.FPDF_GetSignatureCount, error)

	// FPDF_GetSignatureObject returns the Nth signature of the document.
	// Experimental API.
	FPDF_GetSignatureObject(request *requests.FPDF_GetSignatureObject) (*responses.FPDF_GetSignatureObject, error)

	// FPDFSignatureObj_GetContents returns the contents of a signature object.
	// Experimental API.
	FPDFSignatureObj_GetContents(request *requests.FPDFSignatureObj_GetContents) (*responses.FPDFSignatureObj_GetContents, error)

	// FPDFSignatureObj_GetByteRange returns the byte range of a signature object.
	// Experimental API.
	FPDFSignatureObj_GetByteRange(request *requests.FPDFSignatureObj_GetByteRange) (*responses.FPDFSignatureObj_GetByteRange, error)

	// FPDFSignatureObj_GetSubFilter returns the encoding of the value of a signature object.
	// Experimental API.
	FPDFSignatureObj_GetSubFilter(request *requests.FPDFSignatureObj_GetSubFilter) (*responses.FPDFSignatureObj_GetSubFilter, error)

	// FPDFSignatureObj_GetReason returns the reason (comment) of the signature object.
	// Experimental API.
	FPDFSignatureObj_GetReason(request *requests.FPDFSignatureObj_GetReason) (*responses.FPDFSignatureObj_GetReason, error)

	// FPDFSignatureObj_GetTime returns the time of signing of a signature object.
	// Experimental API.
	FPDFSignatureObj_GetTime(request *requests.FPDFSignatureObj_GetTime) (*responses.FPDFSignatureObj_GetTime, error)

	// FPDFSignatureObj_GetDocMDPPermission returns the DocMDP permission of a signature object.
	// Experimental API.
	FPDFSignatureObj_GetDocMDPPermission(request *requests.FPDFSignatureObj_GetDocMDPPermission) (*responses.FPDFSignatureObj_GetDocMDPPermission, error)

	// End fpdf_signature.h

	// Start fpdf_thumbnail.h

	// FPDFPage_GetDecodedThumbnailData returns the decoded data from the thumbnail of the given page if it exists.
	// Experimental API.
	FPDFPage_GetDecodedThumbnailData(request *requests.FPDFPage_GetDecodedThumbnailData) (*responses.FPDFPage_GetDecodedThumbnailData, error)

	// FPDFPage_GetRawThumbnailData returns the raw data from the thumbnail of the given page if it exists.
	// Experimental API.
	FPDFPage_GetRawThumbnailData(request *requests.FPDFPage_GetRawThumbnailData) (*responses.FPDFPage_GetRawThumbnailData, error)

	// FPDFPage_GetThumbnailAsBitmap returns the thumbnail of the given page as a FPDF_BITMAP.
	// Experimental API.
	FPDFPage_GetThumbnailAsBitmap(request *requests.FPDFPage_GetThumbnailAsBitmap) (*responses.FPDFPage_GetThumbnailAsBitmap, error)

	// End fpdf_thumbnail.h

	// Start fpdf_attachment.h

	// FPDFDoc_GetAttachmentCount returns the number of embedded files in the given document.
	// Experimental API.
	FPDFDoc_GetAttachmentCount(request *requests.FPDFDoc_GetAttachmentCount) (*responses.FPDFDoc_GetAttachmentCount, error)

	// FPDFDoc_AddAttachment adds an embedded file with the given name in the given document. If the name is empty, or if
	// the name is the name of an existing embedded file in the document, or if
	// the document's embedded file name tree is too deep (i.e. the document has too
	// many embedded files already), then a new attachment will not be added.
	// Experimental API.
	FPDFDoc_AddAttachment(request *requests.FPDFDoc_AddAttachment) (*responses.FPDFDoc_AddAttachment, error)

	// FPDFDoc_GetAttachment returns the embedded attachment at the given index in the given document. Note that the returned
	// attachment handle is only valid while the document is open.
	// Experimental API.
	FPDFDoc_GetAttachment(request *requests.FPDFDoc_GetAttachment) (*responses.FPDFDoc_GetAttachment, error)

	// FPDFDoc_DeleteAttachment deletes the embedded attachment at the given index in the given document. Note that this does
	// not remove the attachment data from the PDF file; it simply removes the
	// file's entry in the embedded files name tree so that it does not appear in
	// the attachment list. This behavior may change in the future.
	// Experimental API.
	FPDFDoc_DeleteAttachment(request *requests.FPDFDoc_DeleteAttachment) (*responses.FPDFDoc_DeleteAttachment, error)

	// FPDFAttachment_GetName returns the name of the attachment file.
	// Experimental API.
	FPDFAttachment_GetName(request *requests.FPDFAttachment_GetName) (*responses.FPDFAttachment_GetName, error)

	// FPDFAttachment_HasKey check if the params dictionary of the given attachment has the given key as a key.
	// Experimental API.
	FPDFAttachment_HasKey(request *requests.FPDFAttachment_HasKey) (*responses.FPDFAttachment_HasKey, error)

	// FPDFAttachment_GetValueType returns the type of the value corresponding to the given key in the params dictionary of
	// the embedded attachment.
	// Experimental API.
	FPDFAttachment_GetValueType(request *requests.FPDFAttachment_GetValueType) (*responses.FPDFAttachment_GetValueType, error)

	// FPDFAttachment_SetStringValue sets the string value corresponding to the given key in the params dictionary of the
	// embedded file attachment, overwriting the existing value if any.
	// Experimental API.
	FPDFAttachment_SetStringValue(request *requests.FPDFAttachment_SetStringValue) (*responses.FPDFAttachment_SetStringValue, error)

	// FPDFAttachment_GetStringValue gets the string value corresponding to the given key in the params dictionary of the
	// embedded file attachment.
	// Experimental API.
	FPDFAttachment_GetStringValue(request *requests.FPDFAttachment_GetStringValue) (*responses.FPDFAttachment_GetStringValue, error)

	// FPDFAttachment_SetFile set the file data of the given attachment, overwriting the existing file data if any.
	// The creation date and checksum will be updated, while all other dictionary
	// entries will be deleted. Note that only contents with a length smaller than
	// INT_MAX is supported.
	// Experimental API.
	FPDFAttachment_SetFile(request *requests.FPDFAttachment_SetFile) (*responses.FPDFAttachment_SetFile, error)

	// FPDFAttachment_GetFile gets the file data of the given attachment.
	// Experimental API.
	FPDFAttachment_GetFile(request *requests.FPDFAttachment_GetFile) (*responses.FPDFAttachment_GetFile, error)

	// FPDFAttachment_GetSubtype gets the MIME type (Subtype) of the embedded file attachment.
	// Experimental API.
	FPDFAttachment_GetSubtype(request *requests.FPDFAttachment_GetSubtype) (*responses.FPDFAttachment_GetSubtype, error)

	// End fpdf_attachment.h

	// Start attachment: attachment helpers

	// GetAttachments returns all the attachments of a document.
	// Experimental API.
	GetAttachments(request *requests.GetAttachments) (*responses.GetAttachments, error)

	// End attachment

	// Start fpdf_javascript.h

	// FPDFDoc_GetJavaScriptActionCount returns the number of JavaScript actions in the given document.
	// Experimental API.
	FPDFDoc_GetJavaScriptActionCount(request *requests.FPDFDoc_GetJavaScriptActionCount) (*responses.FPDFDoc_GetJavaScriptActionCount, error)

	// FPDFDoc_GetJavaScriptAction returns the JavaScript action at the given index in the given document.
	// Experimental API.
	FPDFDoc_GetJavaScriptAction(request *requests.FPDFDoc_GetJavaScriptAction) (*responses.FPDFDoc_GetJavaScriptAction, error)

	// FPDFDoc_CloseJavaScriptAction closes a loaded FPDF_JAVASCRIPT_ACTION object.
	// Experimental API.
	FPDFDoc_CloseJavaScriptAction(request *requests.FPDFDoc_CloseJavaScriptAction) (*responses.FPDFDoc_CloseJavaScriptAction, error)

	// FPDFJavaScriptAction_GetName returns the name from the javascript handle.
	// Experimental API.
	FPDFJavaScriptAction_GetName(request *requests.FPDFJavaScriptAction_GetName) (*responses.FPDFJavaScriptAction_GetName, error)

	// FPDFJavaScriptAction_GetScript returns the script from the javascript handle
	// Experimental API.
	FPDFJavaScriptAction_GetScript(request *requests.FPDFJavaScriptAction_GetScript) (*responses.FPDFJavaScriptAction_GetScript, error)

	// End fpdf_javascript.h

	// start javascript_action: javascript action helper

	// GetJavaScriptActions returns all the JavaScript Actions of a document.
	// Experimental API.
	GetJavaScriptActions(request *requests.GetJavaScriptActions) (*responses.GetJavaScriptActions, error)

	// End javascript_action

	// Start fpdf_text.h

	// FPDFText_LoadPage returns a handle to the text page information structure.
	// Application must call FPDFText_ClosePage to release the text page
	FPDFText_LoadPage(request *requests.FPDFText_LoadPage) (*responses.FPDFText_LoadPage, error)

	// FPDFText_ClosePage Release all resources allocated for a text page information structure.
	FPDFText_ClosePage(request *requests.FPDFText_ClosePage) (*responses.FPDFText_ClosePage, error)

	// FPDFText_CountChars returns the number of characters in a page.
	// Characters in a page form a "stream", inside the stream, each character has an index.
	// We will use the index parameters in many of FPDFTEXT functions. The first
	// character in the page has an index value of zero.
	FPDFText_CountChars(request *requests.FPDFText_CountChars) (*responses.FPDFText_CountChars, error)

	// FPDFText_GetUnicode returns the unicode of a character in a page.
	FPDFText_GetUnicode(request *requests.FPDFText_GetUnicode) (*responses.FPDFText_GetUnicode, error)

	// FPDFText_GetTextObject returns the FPDF_PAGEOBJECT associated with a given character.
	// Experimental API.
	FPDFText_GetTextObject(request *requests.FPDFText_GetTextObject) (*responses.FPDFText_GetTextObject, error)

	// FPDFText_IsGenerated returns whether a character in a page is generated by PDFium.
	// Experimental API.
	FPDFText_IsGenerated(request *requests.FPDFText_IsGenerated) (*responses.FPDFText_IsGenerated, error)

	// FPDFText_IsHyphen returns whether a character in a page is a hyphen.
	// Experimental API.
	FPDFText_IsHyphen(request *requests.FPDFText_IsHyphen) (*responses.FPDFText_IsHyphen, error)

	// FPDFText_HasUnicodeMapError a character in a page has an invalid unicode mapping.
	// Experimental API.
	FPDFText_HasUnicodeMapError(request *requests.FPDFText_HasUnicodeMapError) (*responses.FPDFText_HasUnicodeMapError, error)

	// FPDFText_GetFontSize returns the font size of a particular character.
	FPDFText_GetFontSize(request *requests.FPDFText_GetFontSize) (*responses.FPDFText_GetFontSize, error)

	// FPDFText_GetFontInfo returns the font name and flags of a particular character.
	// Experimental API.
	FPDFText_GetFontInfo(request *requests.FPDFText_GetFontInfo) (*responses.FPDFText_GetFontInfo, error)

	// FPDFText_GetFontWeight returns the font weight of a particular character.
	// Experimental API.
	FPDFText_GetFontWeight(request *requests.FPDFText_GetFontWeight) (*responses.FPDFText_GetFontWeight, error)

	// FPDFText_GetFillColor returns the fill color of a particular character.
	// Experimental API.
	FPDFText_GetFillColor(request *requests.FPDFText_GetFillColor) (*responses.FPDFText_GetFillColor, error)

	// FPDFText_GetStrokeColor returns the stroke color of a particular character.
	// Experimental API.
	FPDFText_GetStrokeColor(request *requests.FPDFText_GetStrokeColor) (*responses.FPDFText_GetStrokeColor, error)

	// FPDFText_GetCharAngle returns the character rotation angle.
	// Experimental API.
	FPDFText_GetCharAngle(request *requests.FPDFText_GetCharAngle) (*responses.FPDFText_GetCharAngle, error)

	// FPDFText_GetCharBox returns the bounding box of a particular character.
	// All positions are measured in PDF "user space".
	FPDFText_GetCharBox(request *requests.FPDFText_GetCharBox) (*responses.FPDFText_GetCharBox, error)

	// FPDFText_GetLooseCharBox returns a "loose" bounding box of a particular character, i.e., covering
	// the entire glyph bounds, without taking the actual glyph shape into
	// account. All positions are measured in PDF "user space".
	// Experimental API.
	FPDFText_GetLooseCharBox(request *requests.FPDFText_GetLooseCharBox) (*responses.FPDFText_GetLooseCharBox, error)

	// FPDFText_GetMatrix returns the effective transformation matrix for a particular character.
	// All positions are measured in PDF "user space".
	// Experimental API.
	FPDFText_GetMatrix(request *requests.FPDFText_GetMatrix) (*responses.FPDFText_GetMatrix, error)

	// FPDFText_GetCharOrigin returns origin of a particular character.
	// All positions are measured in PDF "user space".
	FPDFText_GetCharOrigin(request *requests.FPDFText_GetCharOrigin) (*responses.FPDFText_GetCharOrigin, error)

	// FPDFText_GetCharIndexAtPos returns the index of a character at or nearby a certain position on the page.
	FPDFText_GetCharIndexAtPos(request *requests.FPDFText_GetCharIndexAtPos) (*responses.FPDFText_GetCharIndexAtPos, error)

	// FPDFText_GetText extracts unicode text string from the page.
	// This function ignores characters without unicode information.
	// It returns all characters on the page, even those that are not
	// visible when the page has a cropbox. To filter out the characters
	// outside of the cropbox, use FPDF_GetPageBoundingBox() and
	// FPDFText_GetCharBox().
	FPDFText_GetText(request *requests.FPDFText_GetText) (*responses.FPDFText_GetText, error)

	// FPDFText_CountRects returns the count of rectangular areas occupied by
	// a segment of text, and caches the result for subsequent FPDFText_GetRect() calls.
	// This function, along with FPDFText_GetRect can be used by
	// applications to detect the position on the page for a text segment,
	// so proper areas can be highlighted. The FPDFText_* functions will
	// automatically merge small character boxes into bigger one if those
	// characters are on the same line and use same font settings.
	FPDFText_CountRects(request *requests.FPDFText_CountRects) (*responses.FPDFText_CountRects, error)

	// FPDFText_GetRect returns a rectangular area from the result generated by FPDFText_CountRects.
	// Note: this method only works if you called FPDFText_CountRects first.
	FPDFText_GetRect(request *requests.FPDFText_GetRect) (*responses.FPDFText_GetRect, error)

	// FPDFText_GetBoundedText extract unicode text within a rectangular boundary on the page.
	FPDFText_GetBoundedText(request *requests.FPDFText_GetBoundedText) (*responses.FPDFText_GetBoundedText, error)

	// FPDFText_FindStart returns a handle to search a page.
	FPDFText_FindStart(request *requests.FPDFText_FindStart) (*responses.FPDFText_FindStart, error)

	// FPDFText_FindNext searches in the direction from page start to end.
	FPDFText_FindNext(request *requests.FPDFText_FindNext) (*responses.FPDFText_FindNext, error)

	// FPDFText_FindPrev searches in the direction from page end to start.
	FPDFText_FindPrev(request *requests.FPDFText_FindPrev) (*responses.FPDFText_FindPrev, error)

	// FPDFText_GetSchResultIndex returns the starting character index of the search result.
	FPDFText_GetSchResultIndex(request *requests.FPDFText_GetSchResultIndex) (*responses.FPDFText_GetSchResultIndex, error)

	// FPDFText_GetSchCount returns the number of matched characters in the search result.
	FPDFText_GetSchCount(request *requests.FPDFText_GetSchCount) (*responses.FPDFText_GetSchCount, error)

	// FPDFText_FindClose releases a search context.
	FPDFText_FindClose(request *requests.FPDFText_FindClose) (*responses.FPDFText_FindClose, error)

	// FPDFLink_LoadWebLinks prepares information about weblinks in a page.
	FPDFLink_LoadWebLinks(request *requests.FPDFLink_LoadWebLinks) (*responses.FPDFLink_LoadWebLinks, error)

	// FPDFLink_CountWebLinks returns the count of detected web links.
	FPDFLink_CountWebLinks(request *requests.FPDFLink_CountWebLinks) (*responses.FPDFLink_CountWebLinks, error)

	// FPDFLink_GetURL returns the URL information for a detected web link.
	FPDFLink_GetURL(request *requests.FPDFLink_GetURL) (*responses.FPDFLink_GetURL, error)

	// FPDFLink_CountRects returns the count of rectangular areas for the link.
	FPDFLink_CountRects(request *requests.FPDFLink_CountRects) (*responses.FPDFLink_CountRects, error)

	// FPDFLink_GetRect returns the boundaries of a rectangle for a link.
	FPDFLink_GetRect(request *requests.FPDFLink_GetRect) (*responses.FPDFLink_GetRect, error)

	// FPDFLink_GetTextRange returns the start char index and char count for a link.
	// Experimental API.
	FPDFLink_GetTextRange(request *requests.FPDFLink_GetTextRange) (*responses.FPDFLink_GetTextRange, error)

	// FPDFLink_CloseWebLinks releases resources used by weblink feature.
	FPDFLink_CloseWebLinks(request *requests.FPDFLink_CloseWebLinks) (*responses.FPDFLink_CloseWebLinks, error)

	// End fpdf_text.h

	// Start fpdf_searchex.h

	// FPDFText_GetCharIndexFromTextIndex returns the character index in the text page internal character list.
	// Where the character index is an index of the text returned from FPDFText_GetText().
	FPDFText_GetCharIndexFromTextIndex(request *requests.FPDFText_GetCharIndexFromTextIndex) (*responses.FPDFText_GetCharIndexFromTextIndex, error)

	// FPDFText_GetTextIndexFromCharIndex returns the text index in the text page internal character list.
	// Where the text index is an index of the character in the internal character list.
	FPDFText_GetTextIndexFromCharIndex(request *requests.FPDFText_GetTextIndexFromCharIndex) (*responses.FPDFText_GetTextIndexFromCharIndex, error)

	// End fpdf_searchex.h

	// Start fpdf_transformpage.h

	// FPDFPage_SetMediaBox sets the "MediaBox" entry to the page dictionary.
	FPDFPage_SetMediaBox(request *requests.FPDFPage_SetMediaBox) (*responses.FPDFPage_SetMediaBox, error)

	// FPDFPage_SetCropBox sets the "CropBox" entry to the page dictionary.
	FPDFPage_SetCropBox(request *requests.FPDFPage_SetCropBox) (*responses.FPDFPage_SetCropBox, error)

	// FPDFPage_SetBleedBox sets the "BleedBox" entry to the page dictionary.
	FPDFPage_SetBleedBox(request *requests.FPDFPage_SetBleedBox) (*responses.FPDFPage_SetBleedBox, error)

	// FPDFPage_SetTrimBox sets the "TrimBox" entry to the page dictionary.
	FPDFPage_SetTrimBox(request *requests.FPDFPage_SetTrimBox) (*responses.FPDFPage_SetTrimBox, error)

	// FPDFPage_SetArtBox sets the "ArtBox" entry to the page dictionary.
	FPDFPage_SetArtBox(request *requests.FPDFPage_SetArtBox) (*responses.FPDFPage_SetArtBox, error)

	// FPDFPage_GetMediaBox gets the "MediaBox" entry from the page dictionary
	FPDFPage_GetMediaBox(request *requests.FPDFPage_GetMediaBox) (*responses.FPDFPage_GetMediaBox, error)

	// FPDFPage_GetCropBox gets the "CropBox" entry from the page dictionary.
	FPDFPage_GetCropBox(request *requests.FPDFPage_GetCropBox) (*responses.FPDFPage_GetCropBox, error)

	// FPDFPage_GetBleedBox gets the "BleedBox" entry from the page dictionary.
	FPDFPage_GetBleedBox(request *requests.FPDFPage_GetBleedBox) (*responses.FPDFPage_GetBleedBox, error)

	// FPDFPage_GetTrimBox gets the "TrimBox" entry from the page dictionary.
	FPDFPage_GetTrimBox(request *requests.FPDFPage_GetTrimBox) (*responses.FPDFPage_GetTrimBox, error)

	// FPDFPage_GetArtBox gets the "ArtBox" entry from the page dictionary.
	FPDFPage_GetArtBox(request *requests.FPDFPage_GetArtBox) (*responses.FPDFPage_GetArtBox, error)

	// FPDFPage_TransFormWithClip applies the transforms to the page.
	FPDFPage_TransFormWithClip(request *requests.FPDFPage_TransFormWithClip) (*responses.FPDFPage_TransFormWithClip, error)

	// FPDFPageObj_TransformClipPath transform (scale, rotate, shear, move) the clip path of page object.
	FPDFPageObj_TransformClipPath(request *requests.FPDFPageObj_TransformClipPath) (*responses.FPDFPageObj_TransformClipPath, error)

	// FPDFPageObj_GetClipPath Get the clip path of the page object.
	// Experimental API.
	FPDFPageObj_GetClipPath(request *requests.FPDFPageObj_GetClipPath) (*responses.FPDFPageObj_GetClipPath, error)

	// FPDFClipPath_CountPaths returns the number of paths inside the given clip path.
	// Experimental API.
	FPDFClipPath_CountPaths(request *requests.FPDFClipPath_CountPaths) (*responses.FPDFClipPath_CountPaths, error)

	// FPDFClipPath_CountPathSegments returns the number of segments inside one path of the given clip path.
	// Experimental API.
	FPDFClipPath_CountPathSegments(request *requests.FPDFClipPath_CountPathSegments) (*responses.FPDFClipPath_CountPathSegments, error)

	// FPDFClipPath_GetPathSegment returns the segment in one specific path of the given clip path at index.
	// Experimental API.
	FPDFClipPath_GetPathSegment(request *requests.FPDFClipPath_GetPathSegment) (*responses.FPDFClipPath_GetPathSegment, error)

	// FPDF_CreateClipPath creates a new clip path, with a rectangle inserted.
	FPDF_CreateClipPath(request *requests.FPDF_CreateClipPath) (*responses.FPDF_CreateClipPath, error)

	// FPDF_DestroyClipPath destroys the clip path.
	FPDF_DestroyClipPath(request *requests.FPDF_DestroyClipPath) (*responses.FPDF_DestroyClipPath, error)

	// FPDFPage_InsertClipPath Clip the page content, the page content that outside the clipping region become invisible.
	FPDFPage_InsertClipPath(request *requests.FPDFPage_InsertClipPath) (*responses.FPDFPage_InsertClipPath, error)

	// End fpdf_transformpage.h

	// Start fpdf_progressive.h

	// FPDF_RenderPageBitmapWithColorScheme_Start starts to render page contents to a device independent bitmap progressively with a specified color scheme for the content.
	// Not supported on multi-threaded usage.
	// Experimental API.
	FPDF_RenderPageBitmapWithColorScheme_Start(request *requests.FPDF_RenderPageBitmapWithColorScheme_Start) (*responses.FPDF_RenderPageBitmapWithColorScheme_Start, error)

	// FPDF_RenderPageBitmap_Start starts to render page contents to a device independent bitmap progressively.
	// Not supported on multi-threaded usage.
	FPDF_RenderPageBitmap_Start(request *requests.FPDF_RenderPageBitmap_Start) (*responses.FPDF_RenderPageBitmap_Start, error)

	// FPDF_RenderPage_Continue continues rendering a PDF page.
	// Not supported on multi-threaded usage.
	FPDF_RenderPage_Continue(request *requests.FPDF_RenderPage_Continue) (*responses.FPDF_RenderPage_Continue, error)

	// FPDF_RenderPage_Close Release the resource allocate during page rendering. Need to be called after finishing rendering or cancel the rendering.
	// Not supported on multi-threaded usage.
	FPDF_RenderPage_Close(request *requests.FPDF_RenderPage_Close) (*responses.FPDF_RenderPage_Close, error)

	// End fpdf_progressive.h

	// Start fpdf_dataavail.h

	// FPDFAvail_Create creates a document availability provider.
	// FPDFAvail_Destroy() must be called when done with the availability provider.
	FPDFAvail_Create(request *requests.FPDFAvail_Create) (*responses.FPDFAvail_Create, error)

	// FPDFAvail_Destroy destroys the given document availability provider.
	FPDFAvail_Destroy(request *requests.FPDFAvail_Destroy) (*responses.FPDFAvail_Destroy, error)

	// FPDFAvail_IsDocAvail checks if the document is ready for loading, if not, gets download hints.
	// Applications should call this function whenever new data arrives, and process
	// all the generated download hints, if any, until the function returns
	// enums.PDF_FILEAVAIL_DATA_ERROR or enums.PDF_FILEAVAIL_DATA_AVAIL.
	// if hints is nil, the function just check current document availability.
	//
	// Once all data is available, call FPDFAvail_GetDocument() to get a document
	// handle.
	FPDFAvail_IsDocAvail(request *requests.FPDFAvail_IsDocAvail) (*responses.FPDFAvail_IsDocAvail, error)

	// FPDFAvail_GetDocument returns the document from the availability provider.
	// When FPDFAvail_IsDocAvail() returns TRUE, call FPDFAvail_GetDocument() to
	// retrieve the document handle.
	FPDFAvail_GetDocument(request *requests.FPDFAvail_GetDocument) (*responses.FPDFAvail_GetDocument, error)

	// FPDFAvail_GetFirstPageNum returns the page number for the first available page in a linearized PDF.
	// For most linearized PDFs, the first available page will be the first page,
	// however, some PDFs might make another page the first available page.
	// For non-linearized PDFs, this function will always return zero.
	FPDFAvail_GetFirstPageNum(request *requests.FPDFAvail_GetFirstPageNum) (*responses.FPDFAvail_GetFirstPageNum, error)

	// FPDFAvail_IsPageAvail checks if the given page index is ready for loading, if not, it will
	// call the hints to fetch more data.
	FPDFAvail_IsPageAvail(request *requests.FPDFAvail_IsPageAvail) (*responses.FPDFAvail_IsPageAvail, error)

	// FPDFAvail_IsFormAvail
	// This function can be called only after FPDFAvail_GetDocument() is called.
	// Applications should call this function whenever new data arrives and process
	// all the generated download hints, if any, until this function returns
	// enums.PDF_FILEAVAIL_DATA_ERROR or enums.PDF_FILEAVAIL_DATA_AVAIL. Applications can then perform page
	// loading.
	// if hints is nil, the function just check current availability of
	// specified page.
	FPDFAvail_IsFormAvail(request *requests.FPDFAvail_IsFormAvail) (*responses.FPDFAvail_IsFormAvail, error)

	// FPDFAvail_IsLinearized Check whether a document is a linearized PDF.
	// FPDFAvail_IsLinearized() will return enums.PDF_FILEAVAIL_LINEARIZED or enums.PDF_FILEAVAIL_NOT_LINEARIZED
	// when we have 1k  of data. If the files size less than 1k, it returns
	// enums.PDF_FILEAVAIL_LINEARIZATION_UNKNOWN as there is insufficient information to determine
	// if the PDF is linearlized.
	FPDFAvail_IsLinearized(request *requests.FPDFAvail_IsLinearized) (*responses.FPDFAvail_IsLinearized, error)

	// End fpdf_dataavail.h

	// Start fpdf_structtree.h

	// FPDF_StructTree_GetForPage returns the structure tree for a page.
	FPDF_StructTree_GetForPage(request *requests.FPDF_StructTree_GetForPage) (*responses.FPDF_StructTree_GetForPage, error)

	// FPDF_StructTree_Close releases a resource allocated by FPDF_StructTree_GetForPage().
	FPDF_StructTree_Close(request *requests.FPDF_StructTree_Close) (*responses.FPDF_StructTree_Close, error)

	// FPDF_StructTree_CountChildren counts the number of children for the structure tree.
	FPDF_StructTree_CountChildren(request *requests.FPDF_StructTree_CountChildren) (*responses.FPDF_StructTree_CountChildren, error)

	// FPDF_StructTree_GetChildAtIndex returns a child in the structure tree.
	FPDF_StructTree_GetChildAtIndex(request *requests.FPDF_StructTree_GetChildAtIndex) (*responses.FPDF_StructTree_GetChildAtIndex, error)

	// FPDF_StructElement_GetAltText returns the alt text for a given element.
	FPDF_StructElement_GetAltText(request *requests.FPDF_StructElement_GetAltText) (*responses.FPDF_StructElement_GetAltText, error)

	// FPDF_StructElement_GetActualText returns the actual text for a given element.
	// Experimental API.
	FPDF_StructElement_GetActualText(request *requests.FPDF_StructElement_GetActualText) (*responses.FPDF_StructElement_GetActualText, error)

	// FPDF_StructElement_GetID returns the ID for a given element.
	// Experimental API.
	FPDF_StructElement_GetID(request *requests.FPDF_StructElement_GetID) (*responses.FPDF_StructElement_GetID, error)

	// FPDF_StructElement_GetLang returns the case-insensitive IETF BCP 47 language code for an element.
	// Experimental API.
	FPDF_StructElement_GetLang(request *requests.FPDF_StructElement_GetLang) (*responses.FPDF_StructElement_GetLang, error)

	// FPDF_StructElement_GetStringAttribute returns a struct element attribute of type "name" or "string"
	// Experimental API.
	FPDF_StructElement_GetStringAttribute(request *requests.FPDF_StructElement_GetStringAttribute) (*responses.FPDF_StructElement_GetStringAttribute, error)

	// FPDF_StructElement_GetMarkedContentID returns the marked content ID for a given element.
	FPDF_StructElement_GetMarkedContentID(request *requests.FPDF_StructElement_GetMarkedContentID) (*responses.FPDF_StructElement_GetMarkedContentID, error)

	// FPDF_StructElement_GetType returns the type (/S) for a given element.
	FPDF_StructElement_GetType(request *requests.FPDF_StructElement_GetType) (*responses.FPDF_StructElement_GetType, error)

	// FPDF_StructElement_GetObjType returns the object type (/Type) for a given element.
	// Experimental API.
	FPDF_StructElement_GetObjType(request *requests.FPDF_StructElement_GetObjType) (*responses.FPDF_StructElement_GetObjType, error)

	// FPDF_StructElement_GetTitle returns the title (/T) for a given element.
	FPDF_StructElement_GetTitle(request *requests.FPDF_StructElement_GetTitle) (*responses.FPDF_StructElement_GetTitle, error)

	// FPDF_StructElement_CountChildren counts the number of children for the structure element.
	FPDF_StructElement_CountChildren(request *requests.FPDF_StructElement_CountChildren) (*responses.FPDF_StructElement_CountChildren, error)

	// FPDF_StructElement_GetChildAtIndex returns a child in the structure element.
	// If the child exists but is not an element, then this function will
	// return an error. This will also return an error for out of bounds indices.
	FPDF_StructElement_GetChildAtIndex(request *requests.FPDF_StructElement_GetChildAtIndex) (*responses.FPDF_StructElement_GetChildAtIndex, error)

	// FPDF_StructElement_GetChildMarkedContentID returns the child's content id.
	// If the child exists but is not a stream or object, then this
	// function will return an error. This will also return an error for out of bounds
	// indices. Compared to FPDF_StructElement_GetMarkedContentIdAtIndex,
	// it is scoped to the current page.
	// Experimental API.
	FPDF_StructElement_GetChildMarkedContentID(request *requests.FPDF_StructElement_GetChildMarkedContentID) (*responses.FPDF_StructElement_GetChildMarkedContentID, error)

	// FPDF_StructElement_GetParent returns the parent of the structure element.
	// If structure element is StructTreeRoot, then this function will return an error.
	// Experimental API.
	FPDF_StructElement_GetParent(request *requests.FPDF_StructElement_GetParent) (*responses.FPDF_StructElement_GetParent, error)

	// FPDF_StructElement_GetAttributeCount returns the number of attributes for the structure element.
	// Experimental API.
	FPDF_StructElement_GetAttributeCount(request *requests.FPDF_StructElement_GetAttributeCount) (*responses.FPDF_StructElement_GetAttributeCount, error)

	// FPDF_StructElement_GetAttributeAtIndex returns an attribute object in the structure element.
	// If the attribute object exists but is not a dict, then this
	// function will return an error. This will also return an error for out of
	// bounds indices.
	// Experimental API.
	FPDF_StructElement_GetAttributeAtIndex(request *requests.FPDF_StructElement_GetAttributeAtIndex) (*responses.FPDF_StructElement_GetAttributeAtIndex, error)

	// FPDF_StructElement_Attr_GetCount returns the number of attributes in a structure element attribute map.
	// Experimental API.
	FPDF_StructElement_Attr_GetCount(request *requests.FPDF_StructElement_Attr_GetCount) (*responses.FPDF_StructElement_Attr_GetCount, error)

	// FPDF_StructElement_Attr_GetName returns the name of an attribute in a structure element attribute map.
	// Experimental API.
	FPDF_StructElement_Attr_GetName(request *requests.FPDF_StructElement_Attr_GetName) (*responses.FPDF_StructElement_Attr_GetName, error)

	// FPDF_StructElement_Attr_GetValue returns a handle to a value for an attribute in a structure element
	// attribute map. The caller does not own the handle. The handle remains valid as long as the
	// struct_attribute, remains valid.
	// Experimental API.
	FPDF_StructElement_Attr_GetValue(request *requests.FPDF_StructElement_Attr_GetValue) (*responses.FPDF_StructElement_Attr_GetValue, error)

	// FPDF_StructElement_Attr_GetType returns the type of an attribute in a structure element attribute map.
	// Experimental API.
	FPDF_StructElement_Attr_GetType(request *requests.FPDF_StructElement_Attr_GetType) (*responses.FPDF_StructElement_Attr_GetType, error)

	// FPDF_StructElement_Attr_GetBooleanValue returns the value of a boolean attribute in an attribute map by name as
	// FPDF_BOOL. FPDF_StructElement_Attr_GetType() should have returned
	// FPDF_OBJECT_BOOLEAN for this property.
	// Experimental API.
	FPDF_StructElement_Attr_GetBooleanValue(request *requests.FPDF_StructElement_Attr_GetBooleanValue) (*responses.FPDF_StructElement_Attr_GetBooleanValue, error)

	// FPDF_StructElement_Attr_GetNumberValue returns the value of a number attribute in an attribute map by name as
	// float. FPDF_StructElement_Attr_GetType() should have returned
	// FPDF_OBJECT_NUMBER for this property.
	// Experimental API.
	FPDF_StructElement_Attr_GetNumberValue(request *requests.FPDF_StructElement_Attr_GetNumberValue) (*responses.FPDF_StructElement_Attr_GetNumberValue, error)

	// FPDF_StructElement_Attr_GetStringValue returns the value of a string attribute in an attribute map by name as
	// string. FPDF_StructElement_Attr_GetType() should have returned
	// FPDF_OBJECT_STRING or FPDF_OBJECT_NAME for this property.
	// Experimental API.
	FPDF_StructElement_Attr_GetStringValue(request *requests.FPDF_StructElement_Attr_GetStringValue) (*responses.FPDF_StructElement_Attr_GetStringValue, error)

	// FPDF_StructElement_Attr_GetBlobValue returns the value of a blob attribute in an attribute map by name as
	// string.
	// Experimental API.
	FPDF_StructElement_Attr_GetBlobValue(request *requests.FPDF_StructElement_Attr_GetBlobValue) (*responses.FPDF_StructElement_Attr_GetBlobValue, error)

	// FPDF_StructElement_Attr_CountChildren returns the count of the number of children values in an attribute.
	// Experimental API.
	FPDF_StructElement_Attr_CountChildren(request *requests.FPDF_StructElement_Attr_CountChildren) (*responses.FPDF_StructElement_Attr_CountChildren, error)

	// FPDF_StructElement_Attr_GetChildAtIndex returns a child from an attribute at the given index.
	// The index must be less than the result of FPDF_StructElement_Attr_CountChildren().
	// Experimental API.
	FPDF_StructElement_Attr_GetChildAtIndex(request *requests.FPDF_StructElement_Attr_GetChildAtIndex) (*responses.FPDF_StructElement_Attr_GetChildAtIndex, error)

	// FPDF_StructElement_GetMarkedContentIdCount returns the count of marked content ids for a given element.
	// Experimental API.
	FPDF_StructElement_GetMarkedContentIdCount(request *requests.FPDF_StructElement_GetMarkedContentIdCount) (*responses.FPDF_StructElement_GetMarkedContentIdCount, error)

	// FPDF_StructElement_GetMarkedContentIdAtIndex returns the marked content id at a given index for a given element.
	// Experimental API.
	FPDF_StructElement_GetMarkedContentIdAtIndex(request *requests.FPDF_StructElement_GetMarkedContentIdAtIndex) (*responses.FPDF_StructElement_GetMarkedContentIdAtIndex, error)

	// End fpdf_structtree.h

	// Start fpdf_annot.h

	// FPDFAnnot_IsSupportedSubtype returns whether an annotation subtype is currently supported for creation.
	// Experimental API.
	FPDFAnnot_IsSupportedSubtype(request *requests.FPDFAnnot_IsSupportedSubtype) (*responses.FPDFAnnot_IsSupportedSubtype, error)

	// FPDFPage_CreateAnnot creates an annotation in the given page of the given subtype. If the specified
	// subtype is illegal or unsupported, then a new annotation will not be created.
	// Must call FPDFPage_CloseAnnot() when the annotation returned by this
	// function is no longer needed.
	// Experimental API.
	FPDFPage_CreateAnnot(request *requests.FPDFPage_CreateAnnot) (*responses.FPDFPage_CreateAnnot, error)

	// FPDFPage_GetAnnotCount returns the number of annotations in a given page.
	// Experimental API.
	FPDFPage_GetAnnotCount(request *requests.FPDFPage_GetAnnotCount) (*responses.FPDFPage_GetAnnotCount, error)

	// FPDFPage_GetAnnot returns annotation at the given page and index. Must call FPDFPage_CloseAnnot() when the
	// annotation returned by this function is no longer needed.
	// Experimental API.
	FPDFPage_GetAnnot(request *requests.FPDFPage_GetAnnot) (*responses.FPDFPage_GetAnnot, error)

	// FPDFPage_GetAnnotIndex returns the index of the given annotation in the given page. This is the opposite of
	// FPDFPage_GetAnnot().
	// Experimental API.
	FPDFPage_GetAnnotIndex(request *requests.FPDFPage_GetAnnotIndex) (*responses.FPDFPage_GetAnnotIndex, error)

	// FPDFPage_CloseAnnot closes an annotation. Must be called when the annotation returned by
	// FPDFPage_CreateAnnot() or FPDFPage_GetAnnot() is no longer needed. This
	// function does not remove the annotation from the document.
	// Experimental API.
	FPDFPage_CloseAnnot(request *requests.FPDFPage_CloseAnnot) (*responses.FPDFPage_CloseAnnot, error)

	// FPDFPage_RemoveAnnot removes the annotation in the given page at the given index.
	// Experimental API.
	FPDFPage_RemoveAnnot(request *requests.FPDFPage_RemoveAnnot) (*responses.FPDFPage_RemoveAnnot, error)

	// FPDFAnnot_GetSubtype returns the subtype of an annotation.
	// Experimental API.
	FPDFAnnot_GetSubtype(request *requests.FPDFAnnot_GetSubtype) (*responses.FPDFAnnot_GetSubtype, error)

	// FPDFAnnot_IsObjectSupportedSubtype checks whether an annotation subtype is currently supported for object extraction,
	// update, and removal.
	// Experimental API.
	FPDFAnnot_IsObjectSupportedSubtype(request *requests.FPDFAnnot_IsObjectSupportedSubtype) (*responses.FPDFAnnot_IsObjectSupportedSubtype, error)

	// FPDFAnnot_UpdateObject updates the given object in the given annotation. The object must be in the annotation already and must have
	// been retrieved by FPDFAnnot_GetObject(). Currently, only ink and stamp
	// annotations are supported by this API. Also note that only path, image, and
	///text objects have APIs for modification; see FPDFPath_*(), FPDFText_*(), and
	// FPDFImageObj_*().
	// Experimental API.
	FPDFAnnot_UpdateObject(request *requests.FPDFAnnot_UpdateObject) (*responses.FPDFAnnot_UpdateObject, error)

	// FPDFAnnot_AddInkStroke adds a new InkStroke, represented by an array of points, to the InkList of
	// the annotation. The API creates an InkList if one doesn't already exist in the annotation.
	// This API works only for ink annotations. Please refer to ISO 32000-1:2008
	// spec, section 12.5.6.13.
	// Experimental API.
	FPDFAnnot_AddInkStroke(request *requests.FPDFAnnot_AddInkStroke) (*responses.FPDFAnnot_AddInkStroke, error)

	// FPDFAnnot_RemoveInkList removes an InkList in the given annotation.
	// This API works only for ink annotations.
	// Experimental API.
	FPDFAnnot_RemoveInkList(request *requests.FPDFAnnot_RemoveInkList) (*responses.FPDFAnnot_RemoveInkList, error)

	// FPDFAnnot_AppendObject adds the given object to the given annotation. The object must have been created by
	// FPDFPageObj_CreateNew{Path|Rect}() or FPDFPageObj_New{Text|Image}Obj(), and
	// will be owned by the annotation. Note that an object cannot belong to more than one
	// annotation. Currently, only ink and stamp annotations are supported by this API.
	// Also note that only path, image, and text objects have APIs for creation.
	// Experimental API.
	FPDFAnnot_AppendObject(request *requests.FPDFAnnot_AppendObject) (*responses.FPDFAnnot_AppendObject, error)

	// FPDFAnnot_GetObjectCount returns the total number of objects in the given annotation, including path objects, text
	// objects, external objects, image objects, and shading objects.
	// Experimental API.
	FPDFAnnot_GetObjectCount(request *requests.FPDFAnnot_GetObjectCount) (*responses.FPDFAnnot_GetObjectCount, error)

	// FPDFAnnot_GetObject returns the object in the given annotation at the given index.
	// Experimental API.
	FPDFAnnot_GetObject(request *requests.FPDFAnnot_GetObject) (*responses.FPDFAnnot_GetObject, error)

	// FPDFAnnot_RemoveObject removes the object in the given annotation at the given index.
	// Experimental API.
	FPDFAnnot_RemoveObject(request *requests.FPDFAnnot_RemoveObject) (*responses.FPDFAnnot_RemoveObject, error)

	// FPDFAnnot_SetColor sets the color of an annotation. Fails when called on annotations with
	// appearance streams already defined; instead use
	// FPDFPath_Set{Stroke|Fill}Color().
	// Experimental API.
	FPDFAnnot_SetColor(request *requests.FPDFAnnot_SetColor) (*responses.FPDFAnnot_SetColor, error)

	// FPDFAnnot_GetColor returns the color of an annotation. If no color is specified, default to yellow
	// for highlight annotation, black for all else. Fails when called on
	// annotations with appearance streams already defined; instead use
	// FPDFPath_Get{Stroke|Fill}Color().
	// Experimental API.
	FPDFAnnot_GetColor(request *requests.FPDFAnnot_GetColor) (*responses.FPDFAnnot_GetColor, error)

	// FPDFAnnot_HasAttachmentPoints returns whether the annotation is of a type that has attachment points
	// (i.e. quadpoints). Quadpoints are the vertices of the rectangle that
	// encompasses the texts affected by the annotation. They provide the
	// coordinates in the page where the annotation is attached. Only text markup
	// annotations (i.e. highlight, strikeout, squiggly, and underline) and link
	// annotations have quadpoints.
	// Experimental API.
	FPDFAnnot_HasAttachmentPoints(request *requests.FPDFAnnot_HasAttachmentPoints) (*responses.FPDFAnnot_HasAttachmentPoints, error)

	// FPDFAnnot_SetAttachmentPoints replaces the attachment points (i.e. quadpoints) set of an annotation at
	// the given quad index. This index needs to be within the result of
	// FPDFAnnot_CountAttachmentPoints().
	// If the annotation's appearance stream is defined and this annotation is of a
	// type with quadpoints, then update the bounding box too if the new quadpoints
	// define a bigger one.
	// Experimental API.
	FPDFAnnot_SetAttachmentPoints(request *requests.FPDFAnnot_SetAttachmentPoints) (*responses.FPDFAnnot_SetAttachmentPoints, error)

	// FPDFAnnot_AppendAttachmentPoints appends to the list of attachment points (i.e. quadpoints) of an annotation.
	// If the annotation's appearance stream is defined and this annotation is of a
	// type with quadpoints, then update the bounding box too if the new quadpoints
	// define a bigger one.
	// Experimental API.
	FPDFAnnot_AppendAttachmentPoints(request *requests.FPDFAnnot_AppendAttachmentPoints) (*responses.FPDFAnnot_AppendAttachmentPoints, error)

	// FPDFAnnot_CountAttachmentPoints returns the number of sets of quadpoints of an annotation.
	// Experimental API.
	FPDFAnnot_CountAttachmentPoints(request *requests.FPDFAnnot_CountAttachmentPoints) (*responses.FPDFAnnot_CountAttachmentPoints, error)

	// FPDFAnnot_GetAttachmentPoints returns the attachment points (i.e. quadpoints) of an annotation.
	// Experimental API.
	FPDFAnnot_GetAttachmentPoints(request *requests.FPDFAnnot_GetAttachmentPoints) (*responses.FPDFAnnot_GetAttachmentPoints, error)

	// FPDFAnnot_SetRect sets the annotation rectangle defining the location of the annotation. If the
	// annotation's appearance stream is defined and this annotation is of a type
	// without quadpoints, then update the bounding box too if the new rectangle
	// defines a bigger one.
	// Experimental API.
	FPDFAnnot_SetRect(request *requests.FPDFAnnot_SetRect) (*responses.FPDFAnnot_SetRect, error)

	// FPDFAnnot_GetRect returns the annotation rectangle defining the location of the annotation.
	// Experimental API.
	FPDFAnnot_GetRect(request *requests.FPDFAnnot_GetRect) (*responses.FPDFAnnot_GetRect, error)

	// FPDFAnnot_GetVertices returns the vertices of a polygon or polyline annotation.
	// Experimental API.
	FPDFAnnot_GetVertices(request *requests.FPDFAnnot_GetVertices) (*responses.FPDFAnnot_GetVertices, error)

	// FPDFAnnot_GetInkListCount returns the number of paths in the ink list of an ink annotation.
	// Experimental API.
	FPDFAnnot_GetInkListCount(request *requests.FPDFAnnot_GetInkListCount) (*responses.FPDFAnnot_GetInkListCount, error)

	// FPDFAnnot_GetInkListPath returns a path in the ink list of an ink annotation.
	// Experimental API.
	FPDFAnnot_GetInkListPath(request *requests.FPDFAnnot_GetInkListPath) (*responses.FPDFAnnot_GetInkListPath, error)

	// FPDFAnnot_GetLine returns the starting and ending coordinates of a line annotation.
	// Experimental API.
	FPDFAnnot_GetLine(request *requests.FPDFAnnot_GetLine) (*responses.FPDFAnnot_GetLine, error)

	// FPDFAnnot_SetBorder sets the characteristics of the annotation's border (rounded rectangle).
	// Experimental API.
	FPDFAnnot_SetBorder(request *requests.FPDFAnnot_SetBorder) (*responses.FPDFAnnot_SetBorder, error)

	// FPDFAnnot_GetBorder returns the characteristics of the annotation's border (rounded rectangle).
	// Experimental API.
	FPDFAnnot_GetBorder(request *requests.FPDFAnnot_GetBorder) (*responses.FPDFAnnot_GetBorder, error)

	// FPDFAnnot_HasKey checks whether the given annotation's dictionary has the given key as a key.
	// Experimental API.
	FPDFAnnot_HasKey(request *requests.FPDFAnnot_HasKey) (*responses.FPDFAnnot_HasKey, error)

	// FPDFAnnot_GetValueType returns the type of the value corresponding to the given key the annotation's dictionary.
	// Experimental API.
	FPDFAnnot_GetValueType(request *requests.FPDFAnnot_GetValueType) (*responses.FPDFAnnot_GetValueType, error)

	// FPDFAnnot_SetStringValue sets the string value corresponding to the given key in the annotations's dictionary,
	// overwriting the existing value if any. The value type would be
	// FPDF_OBJECT_STRING after this function call succeeds.
	// Experimental API.
	FPDFAnnot_SetStringValue(request *requests.FPDFAnnot_SetStringValue) (*responses.FPDFAnnot_SetStringValue, error)

	// FPDFAnnot_GetStringValue returns the string value corresponding to the given key in the annotations's dictionary.
	// Experimental API.
	FPDFAnnot_GetStringValue(request *requests.FPDFAnnot_GetStringValue) (*responses.FPDFAnnot_GetStringValue, error)

	// FPDFAnnot_GetNumberValue returns the float value corresponding to the given key in the annotations's dictionary.
	// Experimental API.
	FPDFAnnot_GetNumberValue(request *requests.FPDFAnnot_GetNumberValue) (*responses.FPDFAnnot_GetNumberValue, error)

	// FPDFAnnot_SetAP sets the AP (appearance string) in annotations's dictionary for a given appearance mode.
	// Experimental API.
	FPDFAnnot_SetAP(request *requests.FPDFAnnot_SetAP) (*responses.FPDFAnnot_SetAP, error)

	// FPDFAnnot_GetAP returns the AP (appearance string) from annotation's dictionary for a given
	// appearance mode.
	// Experimental API.
	FPDFAnnot_GetAP(request *requests.FPDFAnnot_GetAP) (*responses.FPDFAnnot_GetAP, error)

	// FPDFAnnot_GetLinkedAnnot returns the annotation corresponding to the given key in the annotations's dictionary. Common
	// keys for linking annotations include "IRT" and "Popup". Must call
	// FPDFPage_CloseAnnot() when the annotation returned by this function is no
	// longer needed.
	// Experimental API.
	FPDFAnnot_GetLinkedAnnot(request *requests.FPDFAnnot_GetLinkedAnnot) (*responses.FPDFAnnot_GetLinkedAnnot, error)

	// FPDFAnnot_GetFlags returns the annotation flags of the given annotation.
	// Experimental API.
	FPDFAnnot_GetFlags(request *requests.FPDFAnnot_GetFlags) (*responses.FPDFAnnot_GetFlags, error)

	// FPDFAnnot_SetFlags sets the annotation flags of the given annotation.
	// Experimental API.
	FPDFAnnot_SetFlags(request *requests.FPDFAnnot_SetFlags) (*responses.FPDFAnnot_SetFlags, error)

	// FPDFAnnot_GetFormFieldFlags returns the form field annotation flags of the given annotation.
	// Experimental API.
	FPDFAnnot_GetFormFieldFlags(request *requests.FPDFAnnot_GetFormFieldFlags) (*responses.FPDFAnnot_GetFormFieldFlags, error)

	// FPDFAnnot_SetFormFieldFlags sets the form field flags for an interactive form annotation.
	// Experimental API.
	FPDFAnnot_SetFormFieldFlags(request *requests.FPDFAnnot_SetFormFieldFlags) (*responses.FPDFAnnot_SetFormFieldFlags, error)

	// FPDFAnnot_GetFormFieldAtPoint returns an interactive form annotation whose rectangle contains a given
	// point on a page. Must call FPDFPage_CloseAnnot() when the annotation returned
	// is no longer needed.
	// Experimental API.
	FPDFAnnot_GetFormFieldAtPoint(request *requests.FPDFAnnot_GetFormFieldAtPoint) (*responses.FPDFAnnot_GetFormFieldAtPoint, error)

	// FPDFAnnot_GetFormAdditionalActionJavaScript returns the JavaScript of an event of the annotation's additional actions.
	// Experimental API.
	FPDFAnnot_GetFormAdditionalActionJavaScript(request *requests.FPDFAnnot_GetFormAdditionalActionJavaScript) (*responses.FPDFAnnot_GetFormAdditionalActionJavaScript, error)

	// FPDFAnnot_GetFormFieldName returns the name of the given annotation, which is an interactive form annotation.
	// Experimental API.
	FPDFAnnot_GetFormFieldName(request *requests.FPDFAnnot_GetFormFieldName) (*responses.FPDFAnnot_GetFormFieldName, error)

	// FPDFAnnot_GetFormFieldAlternateName returns the alternate name of an annotation, which is an interactive form annotation.
	// Experimental API.
	FPDFAnnot_GetFormFieldAlternateName(request *requests.FPDFAnnot_GetFormFieldAlternateName) (*responses.FPDFAnnot_GetFormFieldAlternateName, error)

	// FPDFAnnot_GetFormFieldType returns the form field type of the given annotation, which is an interactive form annotation.
	// Experimental API.
	FPDFAnnot_GetFormFieldType(request *requests.FPDFAnnot_GetFormFieldType) (*responses.FPDFAnnot_GetFormFieldType, error)

	// FPDFAnnot_GetFormFieldValue returns the value of the given annotation, which is an interactive form annotation.
	// Experimental API.
	FPDFAnnot_GetFormFieldValue(request *requests.FPDFAnnot_GetFormFieldValue) (*responses.FPDFAnnot_GetFormFieldValue, error)

	// FPDFAnnot_GetOptionCount returns the number of options in the annotation's "Opt" dictionary. Intended for
	// use with listbox and combobox widget annotations.
	// Experimental API.
	FPDFAnnot_GetOptionCount(request *requests.FPDFAnnot_GetOptionCount) (*responses.FPDFAnnot_GetOptionCount, error)

	// FPDFAnnot_GetOptionLabel returns the string value for the label of the option at the given index in annotation's
	// "Opt" dictionary. Intended for use with listbox and combobox widget
	// annotations.
	// Experimental API.
	FPDFAnnot_GetOptionLabel(request *requests.FPDFAnnot_GetOptionLabel) (*responses.FPDFAnnot_GetOptionLabel, error)

	// FPDFAnnot_IsOptionSelected returns whether or not the option at the given index in annotation's "Opt" dictionary
	// is selected. Intended for use with listbox and combobox widget annotations.
	// Experimental API.
	FPDFAnnot_IsOptionSelected(request *requests.FPDFAnnot_IsOptionSelected) (*responses.FPDFAnnot_IsOptionSelected, error)

	// FPDFAnnot_GetFontSize returns the float value of the font size for an annotation with variable text.
	// If 0, the font is to be auto-sized: its size is computed as a function of
	// the height of the annotation rectangle.
	// Experimental API.
	FPDFAnnot_GetFontSize(request *requests.FPDFAnnot_GetFontSize) (*responses.FPDFAnnot_GetFontSize, error)

	// FPDFAnnot_SetFontColor Set the text color of an annotation.
	// Currently supported subtypes: freetext.
	// The range for the color components is 0 to 255.
	// Experimental API.
	FPDFAnnot_SetFontColor(request *requests.FPDFAnnot_SetFontColor) (*responses.FPDFAnnot_SetFontColor, error)

	// FPDFAnnot_GetFontColor returns the RGB value of the font color for an annotation with variable text.
	// Experimental API.
	FPDFAnnot_GetFontColor(request *requests.FPDFAnnot_GetFontColor) (*responses.FPDFAnnot_GetFontColor, error)

	// FPDFAnnot_IsChecked returns whether the given annotation is a form widget that is checked. Intended for use with
	// checkbox and radio button widgets.
	// Experimental API.
	FPDFAnnot_IsChecked(request *requests.FPDFAnnot_IsChecked) (*responses.FPDFAnnot_IsChecked, error)

	// FPDFAnnot_SetFocusableSubtypes sets the list of focusable annotation subtypes. Annotations of subtype
	// FPDF_ANNOT_WIDGET are by default focusable. New subtypes set using this API
	// will override the existing subtypes.
	// Experimental API.
	FPDFAnnot_SetFocusableSubtypes(request *requests.FPDFAnnot_SetFocusableSubtypes) (*responses.FPDFAnnot_SetFocusableSubtypes, error)

	// FPDFAnnot_GetFocusableSubtypesCount returns the count of focusable annotation subtypes as set by host.
	// Experimental API.
	FPDFAnnot_GetFocusableSubtypesCount(request *requests.FPDFAnnot_GetFocusableSubtypesCount) (*responses.FPDFAnnot_GetFocusableSubtypesCount, error)

	// FPDFAnnot_GetFocusableSubtypes returns the list of focusable annotation subtype as set by host.
	// Experimental API.
	FPDFAnnot_GetFocusableSubtypes(request *requests.FPDFAnnot_GetFocusableSubtypes) (*responses.FPDFAnnot_GetFocusableSubtypes, error)

	// FPDFAnnot_GetLink returns FPDF_LINK object for the given annotation. Intended to use for link annotations.
	// Experimental API.
	FPDFAnnot_GetLink(request *requests.FPDFAnnot_GetLink) (*responses.FPDFAnnot_GetLink, error)

	// FPDFAnnot_GetFormControlCount returns the count of annotations in the annotation's control group.
	// A group of interactive form annotations is collectively called a form
	// control group. Here, annotation, an interactive form annotation, should be
	// either a radio button or a checkbox.
	// Experimental API.
	FPDFAnnot_GetFormControlCount(request *requests.FPDFAnnot_GetFormControlCount) (*responses.FPDFAnnot_GetFormControlCount, error)

	// FPDFAnnot_GetFormControlIndex returns the index of the given annotation it's control group.
	// A group of interactive form annotations is collectively called a form
	// control group. Here, the annotation, an interactive form annotation, should be
	// either a radio button or a checkbox.
	// Experimental API.
	FPDFAnnot_GetFormControlIndex(request *requests.FPDFAnnot_GetFormControlIndex) (*responses.FPDFAnnot_GetFormControlIndex, error)

	// FPDFAnnot_GetFormFieldExportValue returns the export value of the given annotation which is an interactive form annotation.
	// Intended for use with radio button and checkbox widget annotations.
	// Experimental API.
	FPDFAnnot_GetFormFieldExportValue(request *requests.FPDFAnnot_GetFormFieldExportValue) (*responses.FPDFAnnot_GetFormFieldExportValue, error)

	// FPDFAnnot_SetURI adds a URI action to the given annotation, overwriting the existing action, if any.
	// Experimental API.
	FPDFAnnot_SetURI(request *requests.FPDFAnnot_SetURI) (*responses.FPDFAnnot_SetURI, error)

	// FPDFAnnot_GetFileAttachment get the attachment from the given annotation.
	// Experimental API.
	FPDFAnnot_GetFileAttachment(request *requests.FPDFAnnot_GetFileAttachment) (*responses.FPDFAnnot_GetFileAttachment, error)

	// FPDFAnnot_AddFileAttachment Add an embedded file to the given annotation.
	// Experimental API.
	FPDFAnnot_AddFileAttachment(request *requests.FPDFAnnot_AddFileAttachment) (*responses.FPDFAnnot_AddFileAttachment, error)

	// End fpdf_annot.h

	// Start fpdf_formfill.h

	// FPDFDOC_InitFormFillEnvironment initializes form fill environment
	// This function should be called before any form fill operation.
	FPDFDOC_InitFormFillEnvironment(request *requests.FPDFDOC_InitFormFillEnvironment) (*responses.FPDFDOC_InitFormFillEnvironment, error)

	// FPDFDOC_ExitFormFillEnvironment takes ownership of the handle and exits form fill environment.
	FPDFDOC_ExitFormFillEnvironment(request *requests.FPDFDOC_ExitFormFillEnvironment) (*responses.FPDFDOC_ExitFormFillEnvironment, error)

	// FORM_OnAfterLoadPage
	// This method is required for implementing all the form related
	// functions. Should be invoked after user successfully loaded a
	// PDF page, and FPDFDOC_InitFormFillEnvironment() has been invoked.
	FORM_OnAfterLoadPage(request *requests.FORM_OnAfterLoadPage) (*responses.FORM_OnAfterLoadPage, error)

	// FORM_OnBeforeClosePage
	// This method is required for implementing all the form related
	// functions. Should be invoked before user closes the PDF page.
	FORM_OnBeforeClosePage(request *requests.FORM_OnBeforeClosePage) (*responses.FORM_OnBeforeClosePage, error)

	// FORM_DoDocumentJSAction
	// This method is required for performing document-level JavaScript
	// actions. It should be invoked after the PDF document has been loaded.
	// If there is document-level JavaScript action embedded in the
	// document, this method will execute the JavaScript action. Otherwise,
	// the method will do nothing.
	FORM_DoDocumentJSAction(request *requests.FORM_DoDocumentJSAction) (*responses.FORM_DoDocumentJSAction, error)

	// FORM_DoDocumentOpenAction
	// This method is required for performing open-action when the document
	// is opened.
	// This method will do nothing if there are no open-actions embedded
	// in the document.
	FORM_DoDocumentOpenAction(request *requests.FORM_DoDocumentOpenAction) (*responses.FORM_DoDocumentOpenAction, error)

	// FORM_DoDocumentAAction
	// This method is required for performing the document's
	// additional-action.
	// This method will do nothing if there is no document
	// additional-action corresponding to the specified type.
	FORM_DoDocumentAAction(request *requests.FORM_DoDocumentAAction) (*responses.FORM_DoDocumentAAction, error)

	// FORM_DoPageAAction
	// This method is required for performing the page object's
	// additional-action when opened or closed.
	// This method will do nothing if no additional-action corresponding
	// to the specified type exists.
	FORM_DoPageAAction(request *requests.FORM_DoPageAAction) (*responses.FORM_DoPageAAction, error)

	// FORM_OnMouseMove
	// Call this member function when the mouse cursor moves.
	FORM_OnMouseMove(request *requests.FORM_OnMouseMove) (*responses.FORM_OnMouseMove, error)

	// FORM_OnMouseWheel
	// Call this member function when the user scrolls the mouse wheel.
	// For X and Y delta, the caller must normalize
	// platform-specific wheel deltas. e.g. On Windows, a delta value of 240
	// for a WM_MOUSEWHEEL event normalizes to 2, since Windows defines
	// WHEEL_DELTA as 120.
	// Experimental API
	FORM_OnMouseWheel(request *requests.FORM_OnMouseWheel) (*responses.FORM_OnMouseWheel, error)

	// FORM_OnFocus
	// This function focuses the form annotation at a given point. If the
	// annotation at the point already has focus, nothing happens. If there
	// is no annotation at the point, removes form focus.
	FORM_OnFocus(request *requests.FORM_OnFocus) (*responses.FORM_OnFocus, error)

	// FORM_OnLButtonDown
	// Call this member function when the user presses the left
	// mouse button.
	FORM_OnLButtonDown(request *requests.FORM_OnLButtonDown) (*responses.FORM_OnLButtonDown, error)

	// FORM_OnRButtonDown
	// Call this member function when the user presses the right
	// mouse button.
	// At the present time, has no effect except in XFA builds, but is
	// included for the sake of symmetry.
	FORM_OnRButtonDown(request *requests.FORM_OnRButtonDown) (*responses.FORM_OnRButtonDown, error)

	// FORM_OnLButtonUp
	// Call this member function when the user releases the left
	// mouse button.
	FORM_OnLButtonUp(request *requests.FORM_OnLButtonUp) (*responses.FORM_OnLButtonUp, error)

	// FORM_OnRButtonUp
	// Call this member function when the user releases the right
	// mouse button.
	// At the present time, has no effect except in XFA builds, but is
	// included for the sake of symmetry.
	FORM_OnRButtonUp(request *requests.FORM_OnRButtonUp) (*responses.FORM_OnRButtonUp, error)

	// FORM_OnLButtonDoubleClick
	// Call this member function when the user double clicks the
	// left mouse button.
	FORM_OnLButtonDoubleClick(request *requests.FORM_OnLButtonDoubleClick) (*responses.FORM_OnLButtonDoubleClick, error)

	// FORM_OnKeyDown
	// Call this member function when a nonsystem key is pressed.
	FORM_OnKeyDown(request *requests.FORM_OnKeyDown) (*responses.FORM_OnKeyDown, error)

	// FORM_OnKeyUp
	// Call this member function when a nonsystem key is released.
	// Currently unimplemented and always returns false. PDFium reserves this
	// API and may implement it in the future on an as-needed basis.
	FORM_OnKeyUp(request *requests.FORM_OnKeyUp) (*responses.FORM_OnKeyUp, error)

	// FORM_OnChar
	// Call this member function when a keystroke translates to a
	// nonsystem character.
	FORM_OnChar(request *requests.FORM_OnChar) (*responses.FORM_OnChar, error)

	// FORM_GetFocusedText
	// Call this function to obtain the text within the current focused
	// field, if any.
	// Experimental API
	FORM_GetFocusedText(request *requests.FORM_GetFocusedText) (*responses.FORM_GetFocusedText, error)

	// FORM_GetSelectedText
	// Call this function to obtain selected text within a form text
	// field or form combobox text field.
	FORM_GetSelectedText(request *requests.FORM_GetSelectedText) (*responses.FORM_GetSelectedText, error)

	// FORM_ReplaceAndKeepSelection
	// Call this function to replace the selected text in a form text field or
	// user-editable form combobox text field with another text string (which
	// can be empty or non-empty). If there is no selected text, this function
	// will append the replacement text after the current caret position. After
	// the insertion, the inserted text will be selected.
	// Experimental API
	FORM_ReplaceAndKeepSelection(request *requests.FORM_ReplaceAndKeepSelection) (*responses.FORM_ReplaceAndKeepSelection, error)

	// FORM_ReplaceSelection
	// Call this function to replace the selected text in a form
	// text field or user-editable form combobox text field with another
	// text string (which can be empty or non-empty). If there is no
	// selected text, this function will append the replacement text after
	// the current caret position.
	FORM_ReplaceSelection(request *requests.FORM_ReplaceSelection) (*responses.FORM_ReplaceSelection, error)

	// FORM_SelectAllText
	// Call this function to select all the text within the currently focused
	// form text field or form combobox text field.
	// Experimental API
	FORM_SelectAllText(request *requests.FORM_SelectAllText) (*responses.FORM_SelectAllText, error)

	// FORM_CanUndo
	// Find out if it is possible for the current focused widget in a given
	// form to perform an undo operation.
	FORM_CanUndo(request *requests.FORM_CanUndo) (*responses.FORM_CanUndo, error)

	// FORM_CanRedo
	// Find out if it is possible for the current focused widget in a given
	// form to perform a redo operation.
	FORM_CanRedo(request *requests.FORM_CanRedo) (*responses.FORM_CanRedo, error)

	// FORM_Undo
	// Make the current focussed widget perform an undo operation.
	FORM_Undo(request *requests.FORM_Undo) (*responses.FORM_Undo, error)

	// FORM_Redo
	// Make the current focussed widget perform a redo operation.
	FORM_Redo(request *requests.FORM_Redo) (*responses.FORM_Redo, error)

	// FORM_ForceToKillFocus
	// Call this member function to force to kill the focus of the form
	// field which has focus. If it would kill the focus of a form field,
	// save the value of form field if was changed by theuser.
	FORM_ForceToKillFocus(request *requests.FORM_ForceToKillFocus) (*responses.FORM_ForceToKillFocus, error)

	// FORM_GetFocusedAnnot
	// Call this member function to get the currently focused annotation.
	// Not currently supported for XFA forms - will report no focused
	// annotation. Must call FPDFPage_CloseAnnot() when the annotation returned
	// by this function is no longer needed.
	// Experimental API.
	FORM_GetFocusedAnnot(request *requests.FORM_GetFocusedAnnot) (*responses.FORM_GetFocusedAnnot, error)

	// FORM_SetFocusedAnnot
	// Call this member function to set the currently focused annotation.
	// The annotation can't be nil. To kill focus, use FORM_ForceToKillFocus() instead.
	// Experimental API.
	FORM_SetFocusedAnnot(request *requests.FORM_SetFocusedAnnot) (*responses.FORM_SetFocusedAnnot, error)

	// FPDFPage_HasFormFieldAtPoint returns the form field type by point.
	FPDFPage_HasFormFieldAtPoint(request *requests.FPDFPage_HasFormFieldAtPoint) (*responses.FPDFPage_HasFormFieldAtPoint, error)

	// FPDFPage_FormFieldZOrderAtPoint returns the form field z-order by point.
	FPDFPage_FormFieldZOrderAtPoint(request *requests.FPDFPage_FormFieldZOrderAtPoint) (*responses.FPDFPage_FormFieldZOrderAtPoint, error)

	// FPDF_SetFormFieldHighlightColor sets the highlight color of the specified (or all) form fields
	// in the document.
	FPDF_SetFormFieldHighlightColor(request *requests.FPDF_SetFormFieldHighlightColor) (*responses.FPDF_SetFormFieldHighlightColor, error)

	// FPDF_SetFormFieldHighlightAlpha sets the transparency of the form field highlight color in the
	// document.
	FPDF_SetFormFieldHighlightAlpha(request *requests.FPDF_SetFormFieldHighlightAlpha) (*responses.FPDF_SetFormFieldHighlightAlpha, error)

	// FPDF_RemoveFormFieldHighlight removes the form field highlight color in the document.
	FPDF_RemoveFormFieldHighlight(request *requests.FPDF_RemoveFormFieldHighlight) (*responses.FPDF_RemoveFormFieldHighlight, error)

	// FPDF_FFLDraw renders FormFields and popup window on a page to a device independent
	// bitmap.
	// This function is designed to render annotations that are
	// user-interactive, which are widget annotations (for FormFields) and
	// popup annotations.
	// With the FPDF_ANNOT flag, this function will render a popup annotation
	// when users mouse-hover on a non-widget annotation. Regardless of
	// FPDF_ANNOT flag, this function will always render widget annotations
	// for FormFields.
	// In order to implement the FormFill functions, implementation should
	// call this function after rendering functions, such as
	// FPDF_RenderPageBitmap() or FPDF_RenderPageBitmap_Start(), have
	// finished rendering the page contents.
	FPDF_FFLDraw(request *requests.FPDF_FFLDraw) (*responses.FPDF_FFLDraw, error)

	// FPDF_GetFormType returns the type of form contained in the PDF document.
	// If document is nil, then the return value is FORMTYPE_NONE.
	// Experimental API
	FPDF_GetFormType(request *requests.FPDF_GetFormType) (*responses.FPDF_GetFormType, error)

	// FORM_SetIndexSelected selects/deselects the value at the given index of the focused
	// annotation.
	// Intended for use with listbox/combobox widget types. Comboboxes
	// have at most a single value selected at a time which cannot be
	// deselected. Deselect on a combobox is a no-op that returns false.
	// Default implementation is a no-op that will return false for
	// other types.
	// Not currently supported for XFA forms - will return false.
	// Experimental API
	FORM_SetIndexSelected(request *requests.FORM_SetIndexSelected) (*responses.FORM_SetIndexSelected, error)

	// FORM_IsIndexSelected returns whether or not the value at index of the focused
	// annotation is currently selected.
	// Intended for use with listbox/combobox widget types. Default
	// implementation is a no-op that will return false for other types.
	// Not currently supported for XFA forms - will return false.
	// Experimental API
	FORM_IsIndexSelected(request *requests.FORM_IsIndexSelected) (*responses.FORM_IsIndexSelected, error)

	// FPDF_LoadXFA load XFA fields of the document if it consists of XFA fields.
	FPDF_LoadXFA(request *requests.FPDF_LoadXFA) (*responses.FPDF_LoadXFA, error)

	// End fpdf_formfill.h
}
