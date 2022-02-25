package pdfium

import (
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
	// If an external buffer is used, then the application should destroy
	// the buffer by itself. FPDFBitmap_Destroy function will not destroy
	// the buffer.
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
	// The data is in BGRA format. Where the A maybe unused if alpha was
	// not specified.
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

	// FPDFPage_SetRotation sets the page rotation for a given page.
	FPDFPage_SetRotation(request *requests.FPDFPage_SetRotation) (*responses.FPDFPage_SetRotation, error)

	// FPDFPage_GetRotation returns the rotation of the given page.
	FPDFPage_GetRotation(request *requests.FPDFPage_GetRotation) (*responses.FPDFPage_GetRotation, error)

	// FPDFPage_InsertObject inserts the given object into a page.
	FPDFPage_InsertObject(request *requests.FPDFPage_InsertObject) (*responses.FPDFPage_InsertObject, error)

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

	// FPDFPageObj_Transform transforms the page object by the given matrix.
	// The matrix is composed as:
	//   |a c e|
	//   |b d f|
	// and can be used to scale, rotate, shear and translate the page object.
	FPDFPageObj_Transform(request *requests.FPDFPageObj_Transform) (*responses.FPDFPageObj_Transform, error)

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
	// into the document.
	// The loaded font can be closed using FPDFFont_Close.
	FPDFText_LoadFont(request *requests.FPDFText_LoadFont) (*responses.FPDFText_LoadFont, error)

	// FPDFText_LoadStandardFont loads one of the standard 14 fonts per PDF spec 1.7 page 416. The preferred
	// way of using font style is using a dash to separate the name from the style,
	// for example 'Helvetica-BoldItalic'.
	// The loaded font can be closed using FPDFFont_Close.
	// Experimental API.
	FPDFText_LoadStandardFont(request *requests.FPDFText_LoadStandardFont) (*responses.FPDFText_LoadStandardFont, error)

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

	// FPDFTextObj_GetFont returns the font of a text object.
	// Experimental API.
	FPDFTextObj_GetFont(request *requests.FPDFTextObj_GetFont) (*responses.FPDFTextObj_GetFont, error)

	// FPDFFont_GetFontName returns the font name of a font.
	// Experimental API.
	FPDFFont_GetFontName(request *requests.FPDFFont_GetFontName) (*responses.FPDFFont_GetFontName, error)

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
	FPDFBookmark_GetFirstChild(request *requests.FPDFBookmark_GetFirstChild) (*responses.FPDFBookmark_GetFirstChild, error)

	// FPDFBookmark_GetNextSibling returns the next bookmark item at the same level.
	FPDFBookmark_GetNextSibling(request *requests.FPDFBookmark_GetNextSibling) (*responses.FPDFBookmark_GetNextSibling, error)

	// FPDFBookmark_GetTitle returns the title of a bookmark.
	FPDFBookmark_GetTitle(request *requests.FPDFBookmark_GetTitle) (*responses.FPDFBookmark_GetTitle, error)

	// FPDFBookmark_Find finds a bookmark in the document, using the bookmark title.
	FPDFBookmark_Find(request *requests.FPDFBookmark_Find) (*responses.FPDFBookmark_Find, error)

	// FPDFBookmark_GetDest returns the destination associated with a bookmark item.
	// If the returned destination is nil, none is associated to the bookmark item.
	FPDFBookmark_GetDest(request *requests.FPDFBookmark_GetDest) (*responses.FPDFBookmark_GetDest, error)

	// FPDFBookmark_GetAction returns the action associated with a bookmark item.
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

	// FPDFText_GetFontSize returns the font size of a particular character.
	FPDFText_GetFontSize(request *requests.FPDFText_GetFontSize) (*responses.FPDFText_GetFontSize, error)

	// FPDFText_GetFontInfo returns the font name and flags of a particular character.
	// Experimental API.
	FPDFText_GetFontInfo(request *requests.FPDFText_GetFontInfo) (*responses.FPDFText_GetFontInfo, error)

	// FPDFText_GetFontWeight returns the font weight of a particular character.
	// Experimental API.
	FPDFText_GetFontWeight(request *requests.FPDFText_GetFontWeight) (*responses.FPDFText_GetFontWeight, error)

	// FPDFText_GetTextRenderMode returns the text rendering mode of character.
	// Experimental API.
	FPDFText_GetTextRenderMode(request *requests.FPDFText_GetTextRenderMode) (*responses.FPDFText_GetTextRenderMode, error)

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
	FPDFText_GetText(request *requests.FPDFText_GetText) (*responses.FPDFText_GetText, error)

	// FPDFText_CountRects returns the count of rectangular areas occupied by a segment of texts.
	// This function, along with FPDFText_GetRect can be used by
	// applications to detect the position on the page for a text segment,
	// so proper areas can be highlighted. FPDFTEXT will automatically
	// merge small character boxes into bigger one if those characters
	// are on the same line and use same font settings.
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

	// FPDF_StructElement_GetTitle returns the title (/T) for a given element.
	FPDF_StructElement_GetTitle(request *requests.FPDF_StructElement_GetTitle) (*responses.FPDF_StructElement_GetTitle, error)

	// FPDF_StructElement_CountChildren counts the number of children for the structure element.
	FPDF_StructElement_CountChildren(request *requests.FPDF_StructElement_CountChildren) (*responses.FPDF_StructElement_CountChildren, error)

	// FPDF_StructElement_GetChildAtIndex returns a child in the structure element.
	// If the child exists but is not an element, then this function will
	// return an error. This will also return an error for out of bounds indices.
	FPDF_StructElement_GetChildAtIndex(request *requests.FPDF_StructElement_GetChildAtIndex) (*responses.FPDF_StructElement_GetChildAtIndex, error)

	// End fpdf_structtree.h
}
