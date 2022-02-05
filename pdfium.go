package pdfium

import (
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

// Pdfium describes a Pdfium instance.
type Pdfium interface {
	// Start instance functions.

	// OpenDocument returns a PDFium references for the given file data.
	// This is a gateway to FPDF_LoadMemDocument, FPDF_LoadMemDocument64, FPDF_LoadDocument and FPDF_LoadCustomDocument.
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
	// If this function fails, you can use FPDF_GetLastError() to retrieve
	// the reason why it failed.
	FPDF_LoadDocument(request *requests.FPDF_LoadDocument) (*responses.FPDF_LoadDocument, error)

	// FPDF_LoadMemDocument opens and load a PDF document from memory.
	// Loaded document can be closed by FPDF_CloseDocument().
	// If this function fails, you can use FPDF_GetLastError() to retrieve
	// the reason why it failed.
	FPDF_LoadMemDocument(request *requests.FPDF_LoadMemDocument) (*responses.FPDF_LoadMemDocument, error)

	// FPDF_LoadMemDocument64 opens and load a PDF document from memory.
	// Loaded document can be closed by FPDF_CloseDocument().
	// If this function fails, you can use FPDF_GetLastError() to retrieve
	// the reason why it failed.
	FPDF_LoadMemDocument64(request *requests.FPDF_LoadMemDocument64) (*responses.FPDF_LoadMemDocument64, error)

	// FPDF_LoadCustomDocument loads a PDF document from a custom access descriptor.
	// This is implemented as an io.ReadSeeker in go-pdfium.
	// This is only really efficient for single threaded usage, the multi-threaded
	// usage will just load the file in memory because it can't transfer readers
	// over gRPC. The single-threaded usage will actually efficiently walk over
	// the PDF as it's being used by PDFium.
	// Loaded document can be closed by FPDF_CloseDocument().
	// If this function fails, you can use FPDF_GetLastError() to retrieve
	// the reason why it failed.
	FPDF_LoadCustomDocument(request *requests.FPDF_LoadCustomDocument) (*responses.FPDF_LoadCustomDocument, error)

	// FPDF_CloseDocument closes the references, releases the resources.
	FPDF_CloseDocument(request *requests.FPDF_CloseDocument) (*responses.FPDF_CloseDocument, error)

	// FPDF_GetLastError returns the last error code of a PDFium function, which is just called.
	// Usually, this function is called after a PDFium function returns, in order to check the error code of the previous PDFium function.
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
	FPDF_GetPageWidth(request *requests.FPDF_GetPageWidth) (*responses.FPDF_GetPageWidth, error)

	// FPDF_GetPageHeight returns the height of a page.
	FPDF_GetPageHeight(request *requests.FPDF_GetPageHeight) (*responses.FPDF_GetPageHeight, error)

	// FPDF_GetPageSizeByIndex returns the size of a page by the page index.
	FPDF_GetPageSizeByIndex(request *requests.FPDF_GetPageSizeByIndex) (*responses.FPDF_GetPageSizeByIndex, error)

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
}
