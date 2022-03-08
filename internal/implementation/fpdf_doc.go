package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_doc.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"
)

// FPDFBookmark_GetFirstChild returns the first child of a bookmark item, or the first top level bookmark item.
func (p *PdfiumImplementation) FPDFBookmark_GetFirstChild(request *requests.FPDFBookmark_GetFirstChild) (*responses.FPDFBookmark_GetFirstChild, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	var parentBookMark C.FPDF_BOOKMARK
	if request.Bookmark != nil {
		bookmarkHandle, err := p.getBookmarkHandle(*request.Bookmark)
		if err != nil {
			return nil, err
		}

		parentBookMark = bookmarkHandle.handle
	}

	bookmark := C.FPDFBookmark_GetFirstChild(documentHandle.handle, parentBookMark)
	if bookmark == nil {
		return &responses.FPDFBookmark_GetFirstChild{}, nil
	}

	newNativeBookmark := p.registerBookmark(bookmark, documentHandle)

	return &responses.FPDFBookmark_GetFirstChild{
		Bookmark: &newNativeBookmark.nativeRef,
	}, nil
}

// FPDFBookmark_GetNextSibling returns the next bookmark item at the same level.
// Note that the caller is responsible for handling circular bookmark
// references, as may arise from malformed documents.
func (p *PdfiumImplementation) FPDFBookmark_GetNextSibling(request *requests.FPDFBookmark_GetNextSibling) (*responses.FPDFBookmark_GetNextSibling, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	bookmarkHandle, err := p.getBookmarkHandle(request.Bookmark)
	if err != nil {
		return nil, err
	}

	bookmark := C.FPDFBookmark_GetNextSibling(documentHandle.handle, bookmarkHandle.handle)
	if bookmark == nil {
		return &responses.FPDFBookmark_GetNextSibling{}, nil
	}

	newNativeBookmark := p.registerBookmark(bookmark, documentHandle)

	return &responses.FPDFBookmark_GetNextSibling{
		Bookmark: &newNativeBookmark.nativeRef,
	}, nil
}

// FPDFBookmark_GetTitle returns the title of a bookmark.
func (p *PdfiumImplementation) FPDFBookmark_GetTitle(request *requests.FPDFBookmark_GetTitle) (*responses.FPDFBookmark_GetTitle, error) {
	p.Lock()
	defer p.Unlock()

	bookmarkHandle, err := p.getBookmarkHandle(request.Bookmark)
	if err != nil {
		return nil, err
	}

	// First get the title length.
	titleSize := C.FPDFBookmark_GetTitle(bookmarkHandle.handle, C.NULL, 0)
	if titleSize == 0 {
		return nil, errors.New("Could not get title")
	}

	charData := make([]byte, titleSize)
	C.FPDFBookmark_GetTitle(bookmarkHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFBookmark_GetTitle{
		Title: transformedText,
	}, nil
}

// FPDFBookmark_Find finds a bookmark in the document, using the bookmark title.
func (p *PdfiumImplementation) FPDFBookmark_Find(request *requests.FPDFBookmark_Find) (*responses.FPDFBookmark_Find, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	if request.Title == "" {
		return nil, errors.New("no title given")
	}

	transformedText, err := p.transformUTF8ToUTF16LE(request.Title)
	if err != nil {
		return nil, err
	}

	bookmark := C.FPDFBookmark_Find(documentHandle.handle, (C.FPDF_WIDESTRING)(unsafe.Pointer(&transformedText[0])))
	if bookmark == nil {
		return &responses.FPDFBookmark_Find{}, nil
	}

	newNativeBookmark := p.registerBookmark(bookmark, documentHandle)

	return &responses.FPDFBookmark_Find{
		Bookmark: &newNativeBookmark.nativeRef,
	}, nil
}

// FPDFBookmark_GetDest returns the destination associated with a bookmark item.
// If the returned destination is nil, none is associated to the bookmark item.
func (p *PdfiumImplementation) FPDFBookmark_GetDest(request *requests.FPDFBookmark_GetDest) (*responses.FPDFBookmark_GetDest, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	bookmarkHandle, err := p.getBookmarkHandle(request.Bookmark)
	if err != nil {
		return nil, err
	}

	dest := C.FPDFBookmark_GetDest(documentHandle.handle, bookmarkHandle.handle)
	if dest == nil {
		return &responses.FPDFBookmark_GetDest{}, nil
	}

	destHandle := p.registerDest(dest, documentHandle)

	return &responses.FPDFBookmark_GetDest{
		Dest: &destHandle.nativeRef,
	}, nil
}

// FPDFBookmark_GetAction returns the action associated with a bookmark item.
// If this function returns a valid handle, it is valid as long as the bookmark is
// valid.
// If the returned action is nil, you should try FPDFBookmark_GetDest.
func (p *PdfiumImplementation) FPDFBookmark_GetAction(request *requests.FPDFBookmark_GetAction) (*responses.FPDFBookmark_GetAction, error) {
	p.Lock()
	defer p.Unlock()

	bookmarkHandle, err := p.getBookmarkHandle(request.Bookmark)
	if err != nil {
		return nil, err
	}

	action := C.FPDFBookmark_GetAction(bookmarkHandle.handle)
	if action == nil {
		return &responses.FPDFBookmark_GetAction{}, nil
	}

	actionHandle := p.registerAction(action)

	return &responses.FPDFBookmark_GetAction{
		Action: &actionHandle.nativeRef,
	}, nil
}

// FPDFAction_GetType returns the action associated with a bookmark item.
func (p *PdfiumImplementation) FPDFAction_GetType(request *requests.FPDFAction_GetType) (*responses.FPDFAction_GetType, error) {
	p.Lock()
	defer p.Unlock()

	actionHandle, err := p.getActionHandle(request.Action)
	if err != nil {
		return nil, err
	}

	actionType := C.FPDFAction_GetType(actionHandle.handle)

	return &responses.FPDFAction_GetType{
		Type: enums.FPDF_ACTION_ACTION(actionType),
	}, nil
}

// FPDFAction_GetDest returns the destination of a specific go-to or remote-goto action.
// Only action with type PDF_ACTION_ACTION_GOTO and PDF_ACTION_ACTION_REMOTEGOTO can have destination data.
// In case of remote goto action, the application should first use function FPDFAction_GetFilePath to get file path, then load that particular document, and use its document handle to call this function.
func (p *PdfiumImplementation) FPDFAction_GetDest(request *requests.FPDFAction_GetDest) (*responses.FPDFAction_GetDest, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	actionHandle, err := p.getActionHandle(request.Action)
	if err != nil {
		return nil, err
	}

	dest := C.FPDFAction_GetDest(documentHandle.handle, actionHandle.handle)
	if dest == nil {
		return &responses.FPDFAction_GetDest{}, nil
	}

	destHandle := p.registerDest(dest, documentHandle)

	return &responses.FPDFAction_GetDest{
		Dest: &destHandle.nativeRef,
	}, nil
}

// FPDFAction_GetFilePath returns the file path from a remote goto or launch action.
// Only works on actions that have the type FPDF_ACTION_ACTION_REMOTEGOTO or FPDF_ACTION_ACTION_LAUNCH.
func (p *PdfiumImplementation) FPDFAction_GetFilePath(request *requests.FPDFAction_GetFilePath) (*responses.FPDFAction_GetFilePath, error) {
	p.Lock()
	defer p.Unlock()

	actionHandle, err := p.getActionHandle(request.Action)
	if err != nil {
		return nil, err
	}

	// First get the file path length.
	filePathLength := C.FPDFAction_GetFilePath(actionHandle.handle, C.NULL, 0)
	if filePathLength == 0 {
		return &responses.FPDFAction_GetFilePath{}, nil
	}

	charData := make([]byte, filePathLength)
	// FPDFAction_GetFilePath returns the data in UTF-8, no conversion needed.
	C.FPDFAction_GetFilePath(actionHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	filePath := string(charData[:filePathLength-1]) // Remove NULL terminator

	return &responses.FPDFAction_GetFilePath{
		FilePath: &filePath,
	}, nil
}

// FPDFAction_GetURIPath returns the URI path from a URI action.
func (p *PdfiumImplementation) FPDFAction_GetURIPath(request *requests.FPDFAction_GetURIPath) (*responses.FPDFAction_GetURIPath, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	actionHandle, err := p.getActionHandle(request.Action)
	if err != nil {
		return nil, err
	}

	// First get the uri path length.
	uriPathLength := C.FPDFAction_GetURIPath(documentHandle.handle, actionHandle.handle, C.NULL, 0)
	if uriPathLength == 0 {
		return &responses.FPDFAction_GetURIPath{}, nil
	}

	charData := make([]byte, uriPathLength)
	C.FPDFAction_GetURIPath(documentHandle.handle, actionHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	uriPath := string(charData[:uriPathLength-1]) // Remove NULL terminator

	return &responses.FPDFAction_GetURIPath{
		URIPath: &uriPath,
	}, nil
}

// FPDFDest_GetDestPageIndex returns the page index from destination data.
func (p *PdfiumImplementation) FPDFDest_GetDestPageIndex(request *requests.FPDFDest_GetDestPageIndex) (*responses.FPDFDest_GetDestPageIndex, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	destHandle, err := p.getDestHandle(request.Dest)
	if err != nil {
		return nil, err
	}

	pageIndex := C.FPDFDest_GetDestPageIndex(documentHandle.handle, destHandle.handle)
	return &responses.FPDFDest_GetDestPageIndex{
		Index: int(pageIndex),
	}, nil
}

// FPDF_GetMetaText returns the requested metadata.
func (p *PdfiumImplementation) FPDF_GetMetaText(request *requests.FPDF_GetMetaText) (*responses.FPDF_GetMetaText, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	cstr := C.CString(request.Tag)
	defer C.free(unsafe.Pointer(cstr))

	// First get the metadata length.
	metaSize := C.FPDF_GetMetaText(documentHandle.handle, cstr, C.NULL, 0)
	if metaSize == 0 {
		return nil, errors.New("Could not get metadata")
	}

	charData := make([]byte, metaSize)
	C.FPDF_GetMetaText(documentHandle.handle, cstr, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_GetMetaText{
		Tag:   request.Tag,
		Value: transformedText,
	}, nil
}

// FPDF_GetPageLabel returns the label for the given page.
func (p *PdfiumImplementation) FPDF_GetPageLabel(request *requests.FPDF_GetPageLabel) (*responses.FPDF_GetPageLabel, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	// First get the label length.
	labelSize := C.FPDF_GetPageLabel(documentHandle.handle, C.int(request.Page), C.NULL, 0)
	if labelSize == 0 {
		return nil, errors.New("Could not get label")
	}

	charData := make([]byte, labelSize)
	C.FPDF_GetPageLabel(documentHandle.handle, C.int(request.Page), unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_GetPageLabel{
		Page:  request.Page,
		Label: transformedText,
	}, nil
}

// FPDFDest_GetLocationInPage returns the (x, y, zoom) location of dest in the destination page, if the
// destination is in [page /XYZ x y zoom] syntax.
func (p *PdfiumImplementation) FPDFDest_GetLocationInPage(request *requests.FPDFDest_GetLocationInPage) (*responses.FPDFDest_GetLocationInPage, error) {
	p.Lock()
	defer p.Unlock()

	destHandle, err := p.getDestHandle(request.Dest)
	if err != nil {
		return nil, err
	}

	hasXVal := C.FPDF_BOOL(0)
	hasYVal := C.FPDF_BOOL(0)
	hasZoomVal := C.FPDF_BOOL(0)

	x := C.FS_FLOAT(0)
	y := C.FS_FLOAT(0)
	zoom := C.FS_FLOAT(0)

	success := C.FPDFDest_GetLocationInPage(destHandle.handle, &hasXVal, &hasYVal, &hasZoomVal, &x, &y, &zoom)
	if int(success) == 0 {
		return nil, errors.New("could not successfully read the /XYZ value")
	}

	resp := &responses.FPDFDest_GetLocationInPage{}

	if int(hasXVal) == 1 {
		xVal := float32(x)
		resp.X = &xVal
	}

	if int(hasYVal) == 1 {
		yVal := float32(y)
		resp.Y = &yVal
	}

	if int(hasZoomVal) == 1 {
		zoomVal := float32(zoom)
		resp.Zoom = &zoomVal
	}

	return resp, nil
}

// FPDFLink_GetLinkAtPoint finds a link at a point on a page.
// You can convert coordinates from screen coordinates to page coordinates using
// FPDF_DeviceToPage().
func (p *PdfiumImplementation) FPDFLink_GetLinkAtPoint(request *requests.FPDFLink_GetLinkAtPoint) (*responses.FPDFLink_GetLinkAtPoint, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	link := C.FPDFLink_GetLinkAtPoint(pageHandle.handle, C.double(request.X), C.double(request.Y))
	if link == nil {
		return &responses.FPDFLink_GetLinkAtPoint{}, nil
	}

	linkHandle := p.registerLink(link)

	return &responses.FPDFLink_GetLinkAtPoint{
		Link: &linkHandle.nativeRef,
	}, nil
}

// FPDFLink_GetLinkZOrderAtPoint finds the Z-order of link at a point on a page.
// You can convert coordinates from screen coordinates to page coordinates using
// FPDF_DeviceToPage().
func (p *PdfiumImplementation) FPDFLink_GetLinkZOrderAtPoint(request *requests.FPDFLink_GetLinkZOrderAtPoint) (*responses.FPDFLink_GetLinkZOrderAtPoint, error) {
	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	zOrder := C.FPDFLink_GetLinkZOrderAtPoint(pageHandle.handle, C.double(request.X), C.double(request.Y))

	return &responses.FPDFLink_GetLinkZOrderAtPoint{
		ZOrder: int(zOrder),
	}, nil
}

// FPDFLink_GetDest returns the destination info for a link.
func (p *PdfiumImplementation) FPDFLink_GetDest(request *requests.FPDFLink_GetDest) (*responses.FPDFLink_GetDest, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	linkHandle, err := p.getLinkHandle(request.Link)
	if err != nil {
		return nil, err
	}

	dest := C.FPDFLink_GetDest(documentHandle.handle, linkHandle.handle)
	if dest == nil {
		return &responses.FPDFLink_GetDest{}, nil
	}

	destHandle := p.registerDest(dest, documentHandle)

	return &responses.FPDFLink_GetDest{
		Dest: &destHandle.nativeRef,
	}, nil
}

// FPDFLink_GetAction returns the action info for a link
// If this function returns a valid handle, it is valid as long as the link is
// valid.
func (p *PdfiumImplementation) FPDFLink_GetAction(request *requests.FPDFLink_GetAction) (*responses.FPDFLink_GetAction, error) {
	p.Lock()
	defer p.Unlock()

	linkHandle, err := p.getLinkHandle(request.Link)
	if err != nil {
		return nil, err
	}

	action := C.FPDFLink_GetAction(linkHandle.handle)
	if action == nil {
		return &responses.FPDFLink_GetAction{}, nil
	}

	actionHandle := p.registerAction(action)

	return &responses.FPDFLink_GetAction{
		Action: &actionHandle.nativeRef,
	}, nil
}

// FPDFLink_Enumerate Enumerates all the link annotations in a page.
func (p *PdfiumImplementation) FPDFLink_Enumerate(request *requests.FPDFLink_Enumerate) (*responses.FPDFLink_Enumerate, error) {
	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	var link C.FPDF_LINK
	startPos := C.int(request.StartPos)

	success := C.FPDFLink_Enumerate(pageHandle.handle, &startPos, &link)
	if int(success) == 0 {
		return &responses.FPDFLink_Enumerate{}, nil
	}

	linkHandle := p.registerLink(link)

	nextStartPos := int(startPos)
	return &responses.FPDFLink_Enumerate{
		NextStartPos: &nextStartPos,
		Link:         &linkHandle.nativeRef,
	}, nil
}

// FPDFLink_GetAnnotRect returns the count of quadrilateral points to the link.
func (p *PdfiumImplementation) FPDFLink_GetAnnotRect(request *requests.FPDFLink_GetAnnotRect) (*responses.FPDFLink_GetAnnotRect, error) {
	linkHandle, err := p.getLinkHandle(request.Link)
	if err != nil {
		return nil, err
	}

	var rect C.FS_RECTF

	success := C.FPDFLink_GetAnnotRect(linkHandle.handle, &rect)
	if int(success) == 0 {
		return &responses.FPDFLink_GetAnnotRect{}, nil
	}

	return &responses.FPDFLink_GetAnnotRect{
		Rect: &structs.FPDF_FS_RECTF{
			Left:   float32(rect.left),
			Top:    float32(rect.top),
			Right:  float32(rect.right),
			Bottom: float32(rect.bottom),
		},
	}, nil
}

// FPDFLink_CountQuadPoints returns the count of quadrilateral points to the link.
func (p *PdfiumImplementation) FPDFLink_CountQuadPoints(request *requests.FPDFLink_CountQuadPoints) (*responses.FPDFLink_CountQuadPoints, error) {
	linkHandle, err := p.getLinkHandle(request.Link)
	if err != nil {
		return nil, err
	}

	count := C.FPDFLink_CountQuadPoints(linkHandle.handle)

	return &responses.FPDFLink_CountQuadPoints{
		Count: int(count),
	}, nil
}

// FPDFLink_GetQuadPoints returns the quadrilateral points for the specified quad index in the link.
func (p *PdfiumImplementation) FPDFLink_GetQuadPoints(request *requests.FPDFLink_GetQuadPoints) (*responses.FPDFLink_GetQuadPoints, error) {
	linkHandle, err := p.getLinkHandle(request.Link)
	if err != nil {
		return nil, err
	}

	var quadPoints C.FS_QUADPOINTSF

	success := C.FPDFLink_GetQuadPoints(linkHandle.handle, C.int(request.QuadIndex), &quadPoints)
	if int(success) == 0 {
		return &responses.FPDFLink_GetQuadPoints{}, nil
	}

	return &responses.FPDFLink_GetQuadPoints{
		Points: &structs.FPDF_FS_QUADPOINTSF{
			X1: float32(quadPoints.x1),
			Y1: float32(quadPoints.y1),
			X2: float32(quadPoints.x2),
			Y2: float32(quadPoints.y2),
			X3: float32(quadPoints.x3),
			Y3: float32(quadPoints.y3),
			X4: float32(quadPoints.x4),
			Y4: float32(quadPoints.y4),
		},
	}, nil
}
