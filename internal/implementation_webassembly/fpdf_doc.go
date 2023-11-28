package implementation_webassembly

import (
	"bytes"
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"
)

// FPDFBookmark_GetFirstChild returns the first child of a bookmark item, or the first top level bookmark item.
// Note that another name for the bookmarks is the document outline, as
// described in ISO 32000-1:2008, section 12.3.3.
func (p *PdfiumImplementation) FPDFBookmark_GetFirstChild(request *requests.FPDFBookmark_GetFirstChild) (*responses.FPDFBookmark_GetFirstChild, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	var parentBookMark uint64
	if request.Bookmark != nil {
		bookmarkHandle, err := p.getBookmarkHandle(*request.Bookmark)
		if err != nil {
			return nil, err
		}

		parentBookMark = *bookmarkHandle.handle
	}

	res, err := p.Module.ExportedFunction("FPDFBookmark_GetFirstChild").Call(p.Context, *documentHandle.handle, parentBookMark)
	if err != nil {
		return nil, err
	}

	bookmark := res[0]
	if bookmark == 0 {
		return &responses.FPDFBookmark_GetFirstChild{}, nil
	}

	newNativeBookmark := p.registerBookmark(&bookmark, documentHandle)

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

	res, err := p.Module.ExportedFunction("FPDFBookmark_GetNextSibling").Call(p.Context, *documentHandle.handle, *bookmarkHandle.handle)
	if err != nil {
		return nil, err
	}

	bookmark := res[0]
	if bookmark == 0 {
		return &responses.FPDFBookmark_GetNextSibling{}, nil
	}

	newNativeBookmark := p.registerBookmark(&bookmark, documentHandle)

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
	res, err := p.Module.ExportedFunction("FPDFBookmark_GetTitle").Call(p.Context, *bookmarkHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	titleSize := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if titleSize == 0 {
		return nil, errors.New("Could not get title")
	}

	charDataPointer, err := p.ByteArrayPointer(titleSize, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDFBookmark_GetTitle").Call(p.Context, *bookmarkHandle.handle, charDataPointer.Pointer, titleSize)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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

	titlePointer, err := p.CFPDF_WIDESTRING(request.Title)
	if err != nil {
		return nil, err
	}
	defer titlePointer.Free()

	res, err := p.Module.ExportedFunction("FPDFBookmark_Find").Call(p.Context, *documentHandle.handle, titlePointer.Pointer)
	if err != nil {
		return nil, err
	}

	bookmark := res[0]
	if bookmark == 0 {
		return &responses.FPDFBookmark_Find{}, nil
	}

	newNativeBookmark := p.registerBookmark(&bookmark, documentHandle)

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

	res, err := p.Module.ExportedFunction("FPDFBookmark_Find").Call(p.Context, *documentHandle.handle, *bookmarkHandle.handle)
	if err != nil {
		return nil, err
	}

	dest := res[0]
	if dest == 0 {
		return &responses.FPDFBookmark_GetDest{}, nil
	}

	destHandle := p.registerDest(&dest, documentHandle)

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

	res, err := p.Module.ExportedFunction("FPDFBookmark_GetAction").Call(p.Context, *bookmarkHandle.handle)
	if err != nil {
		return nil, err
	}

	action := res[0]
	if action == 0 {
		return &responses.FPDFBookmark_GetAction{}, nil
	}

	actionHandle := p.registerAction(&action)

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

	res, err := p.Module.ExportedFunction("FPDFAction_GetType").Call(p.Context, *actionHandle.handle)
	if err != nil {
		return nil, err
	}

	actionType := *(*int32)(unsafe.Pointer(&res[0]))

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

	res, err := p.Module.ExportedFunction("FPDFAction_GetDest").Call(p.Context, *documentHandle.handle, *actionHandle.handle)
	if err != nil {
		return nil, err
	}

	dest := res[0]
	if dest == 0 {
		return &responses.FPDFAction_GetDest{}, nil
	}

	destHandle := p.registerDest(&dest, documentHandle)

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
	res, err := p.Module.ExportedFunction("FPDFAction_GetFilePath").Call(p.Context, *actionHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	filePathLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if filePathLength == 0 {
		return &responses.FPDFAction_GetFilePath{}, nil
	}

	charDataPointer, err := p.ByteArrayPointer(filePathLength, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDFAction_GetFilePath").Call(p.Context, *actionHandle.handle, charDataPointer.Pointer, filePathLength)
	if err != nil {
		return nil, err
	}

	// FPDFAction_GetFilePath returns the data in UTF-8, no conversion needed.
	charData, err := charDataPointer.Value(true)
	if err != nil {
		return nil, err
	}

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
	res, err := p.Module.ExportedFunction("FPDFAction_GetURIPath").Call(p.Context, *documentHandle.handle, *actionHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	uriPathLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if uriPathLength == 0 {
		return &responses.FPDFAction_GetURIPath{}, nil
	}

	charDataPointer, err := p.ByteArrayPointer(uriPathLength, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDFAction_GetURIPath").Call(p.Context, *documentHandle.handle, *actionHandle.handle, charDataPointer.Pointer, uriPathLength)
	if err != nil {
		return nil, err
	}

	// FPDFAction_GetURIPath returns the data in UTF-8, no conversion needed.
	charData, err := charDataPointer.Value(true)
	if err != nil {
		return nil, err
	}

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

	res, err := p.Module.ExportedFunction("FPDFDest_GetDestPageIndex").Call(p.Context, *documentHandle.handle, *destHandle.handle)
	if err != nil {
		return nil, err
	}

	pageIndex := *(*int32)(unsafe.Pointer(&res[0]))
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

	cstr, err := p.CString(request.Tag)
	if err != nil {
		return nil, err
	}
	defer cstr.Free()

	// First get the metadata length.
	res, err := p.Module.ExportedFunction("FPDF_GetMetaText").Call(p.Context, *documentHandle.handle, cstr.Pointer, 0, 0)
	if err != nil {
		return nil, err
	}

	metaSize := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if metaSize == 0 {
		return nil, errors.New("Could not get metadata")
	}

	charDataPointer, err := p.ByteArrayPointer(metaSize, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDF_GetMetaText").Call(p.Context, *documentHandle.handle, cstr.Pointer, charDataPointer.Pointer, metaSize)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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
	res, err := p.Module.ExportedFunction("FPDF_GetPageLabel").Call(p.Context, *documentHandle.handle, *(*uint64)(unsafe.Pointer(&request.Page)), 0, 0)
	if err != nil {
		return nil, err
	}

	labelSize := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if labelSize == 0 {
		return nil, errors.New("Could not get label")
	}

	charDataPointer, err := p.ByteArrayPointer(labelSize, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDF_GetPageLabel").Call(p.Context, *documentHandle.handle, *(*uint64)(unsafe.Pointer(&request.Page)), charDataPointer.Pointer, labelSize)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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

	hasXValPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer hasXValPointer.Free()

	hasYValPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer hasYValPointer.Free()

	hasZoomValPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer hasZoomValPointer.Free()

	xPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer xPointer.Free()

	yPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer yPointer.Free()

	zoomPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer zoomPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFDest_GetLocationInPage").Call(p.Context, *destHandle.handle, hasXValPointer.Pointer, hasYValPointer.Pointer, hasZoomValPointer.Pointer, xPointer.Pointer, yPointer.Pointer, zoomPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not successfully read the /XYZ value")
	}

	resp := &responses.FPDFDest_GetLocationInPage{}

	hasXVal, err := hasXValPointer.Value()
	if err != nil {
		return nil, err
	}

	if int(hasXVal) == 1 {
		x, err := xPointer.Value()
		if err != nil {
			return nil, err
		}
		xVal := float32(x)
		resp.X = &xVal
	}

	hasYVal, err := hasYValPointer.Value()
	if err != nil {
		return nil, err
	}

	if int(hasYVal) == 1 {
		y, err := yPointer.Value()
		if err != nil {
			return nil, err
		}
		yVal := float32(y)
		resp.Y = &yVal
	}

	hasZoomVal, err := hasZoomValPointer.Value()
	if err != nil {
		return nil, err
	}

	if int(hasZoomVal) == 1 {
		zoom, err := zoomPointer.Value()
		if err != nil {
			return nil, err
		}
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

	res, err := p.Module.ExportedFunction("FPDFLink_GetLinkAtPoint").Call(p.Context, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.X)), *(*uint64)(unsafe.Pointer(&request.Y)))
	if err != nil {
		return nil, err
	}

	link := res[0]
	if link == 0 {
		return &responses.FPDFLink_GetLinkAtPoint{}, nil
	}

	linkHandle := p.registerLink(&link)

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

	res, err := p.Module.ExportedFunction("FPDFLink_GetLinkZOrderAtPoint").Call(p.Context, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.X)), *(*uint64)(unsafe.Pointer(&request.Y)))
	if err != nil {
		return nil, err
	}

	zOrder := *(*int32)(unsafe.Pointer(&res[0]))

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

	res, err := p.Module.ExportedFunction("FPDFLink_GetDest").Call(p.Context, *documentHandle.handle, *linkHandle.handle)
	if err != nil {
		return nil, err
	}

	dest := res[0]
	if dest == 0 {
		return &responses.FPDFLink_GetDest{}, nil
	}

	destHandle := p.registerDest(&dest, documentHandle)

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

	res, err := p.Module.ExportedFunction("FPDFLink_GetAction").Call(p.Context, *linkHandle.handle)
	if err != nil {
		return nil, err
	}

	action := res[0]
	if action == 0 {
		return &responses.FPDFLink_GetAction{}, nil
	}

	actionHandle := p.registerAction(&action)

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

	startPosPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer startPosPointer.Free()

	writeSuccess := p.Module.Memory().WriteUint32Le(uint32(startPosPointer.Pointer), uint32(request.StartPos))
	if !writeSuccess {
		return nil, errors.New("could not write startpos to memory")
	}

	linkPointer, err := p.ULongPointer()
	if err != nil {
		return nil, err
	}
	defer linkPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFLink_Enumerate").Call(p.Context, *pageHandle.handle, startPosPointer.Pointer, linkPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return &responses.FPDFLink_Enumerate{}, nil
	}

	link, err := linkPointer.Value()
	if err != nil {
		return nil, err
	}

	linkHandle := p.registerLink(&link)

	startPos, err := startPosPointer.Value()
	if err != nil {
		return nil, err
	}

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

	rectPointer, rectValue, err := p.CStructFS_RECTF(nil)
	if err != nil {
		return nil, err
	}
	defer p.Free(rectPointer)

	res, err := p.Module.ExportedFunction("FPDFLink_GetAnnotRect").Call(p.Context, *linkHandle.handle, rectPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return &responses.FPDFLink_GetAnnotRect{}, nil
	}

	rect, err := rectValue()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFLink_GetAnnotRect{
		Rect: &structs.FPDF_FS_RECTF{
			Left:   float32(rect.Left),
			Top:    float32(rect.Top),
			Right:  float32(rect.Right),
			Bottom: float32(rect.Bottom),
		},
	}, nil
}

// FPDFLink_CountQuadPoints returns the count of quadrilateral points to the link.
func (p *PdfiumImplementation) FPDFLink_CountQuadPoints(request *requests.FPDFLink_CountQuadPoints) (*responses.FPDFLink_CountQuadPoints, error) {
	linkHandle, err := p.getLinkHandle(request.Link)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFLink_CountQuadPoints").Call(p.Context, *linkHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

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

	quadPointsPointer, quadPointsValue, err := p.CStructFS_QUADPOINTSF(nil)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFLink_GetQuadPoints").Call(p.Context, *linkHandle.handle, *(*uint64)(unsafe.Pointer(&request.QuadIndex)), quadPointsPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return &responses.FPDFLink_GetQuadPoints{}, nil
	}

	quadPoints, err := quadPointsValue()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFLink_GetQuadPoints{
		Points: &structs.FPDF_FS_QUADPOINTSF{
			X1: float32(quadPoints.X1),
			Y1: float32(quadPoints.Y1),
			X2: float32(quadPoints.X2),
			Y2: float32(quadPoints.Y2),
			X3: float32(quadPoints.X3),
			Y3: float32(quadPoints.Y3),
			X4: float32(quadPoints.X4),
			Y4: float32(quadPoints.Y4),
		},
	}, nil
}

// FPDFDest_GetView returns the view (fit type) for a given dest.
// Experimental API.
func (p *PdfiumImplementation) FPDFDest_GetView(request *requests.FPDFDest_GetView) (*responses.FPDFDest_GetView, error) {
	p.Lock()
	defer p.Unlock()

	destHandle, err := p.getDestHandle(request.Dest)
	if err != nil {
		return nil, err
	}

	numParamsPointer, err := p.ULongPointer()
	if err != nil {
		return nil, err
	}
	defer numParamsPointer.Free()

	paramsPointer, err := p.FloatArrayPointer(4)
	if err != nil {
		return nil, err
	}
	defer paramsPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFDest_GetView").Call(p.Context, *destHandle.handle, numParamsPointer.Pointer, paramsPointer.Pointer)
	if err != nil {
		return nil, err
	}

	destView := *(*int32)(unsafe.Pointer(&res[0]))

	numParams, err := numParamsPointer.Value()
	if err != nil {
		return nil, err
	}

	params, err := paramsPointer.Value()
	if err != nil {
		return nil, err
	}

	resParams := make([]float32, int(numParams), int(numParams))
	if int(numParams) > 0 {
		for i := range resParams {
			resParams[i] = float32(params[i])
		}
	}

	return &responses.FPDFDest_GetView{
		DestView: enums.FPDF_PDFDEST_VIEW(destView),
		Params:   resParams,
	}, nil
}

// FPDFLink_GetAnnot returns a FPDF_ANNOTATION object for a link.
// Experimental API.
func (p *PdfiumImplementation) FPDFLink_GetAnnot(request *requests.FPDFLink_GetAnnot) (*responses.FPDFLink_GetAnnot, error) {
	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	linkHandle, err := p.getLinkHandle(request.Link)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFLink_GetAnnot").Call(p.Context, *pageHandle.handle, *linkHandle.handle)
	if err != nil {
		return nil, err
	}

	annotation := res[0]
	if annotation == 0 {
		return &responses.FPDFLink_GetAnnot{}, nil
	}

	annotationHandle := p.registerAnnotation(&annotation)

	return &responses.FPDFLink_GetAnnot{
		Annotation: &annotationHandle.nativeRef,
	}, nil
}

// FPDF_GetPageAAction returns an additional-action from page.
// If this function returns a valid handle, it is valid as long as the page is
// valid.
// Experimental API
func (p *PdfiumImplementation) FPDF_GetPageAAction(request *requests.FPDF_GetPageAAction) (*responses.FPDF_GetPageAAction, error) {
	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDF_GetPageAAction").Call(p.Context, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.AAType)))
	if err != nil {
		return nil, err
	}

	action := res[0]
	if action == 0 {
		return &responses.FPDF_GetPageAAction{}, nil
	}

	actionHandle := p.registerAction(&action)

	return &responses.FPDF_GetPageAAction{
		AAType: &request.AAType,
		Action: &actionHandle.nativeRef,
	}, nil
}

// FPDF_GetFileIdentifier Get the file identifier defined in the trailer of a document.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetFileIdentifier(request *requests.FPDF_GetFileIdentifier) (*responses.FPDF_GetFileIdentifier, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	if request.FileIdType != enums.FPDF_FILEIDTYPE_PERMANENT && request.FileIdType != enums.FPDF_FILEIDTYPE_CHANGING {
		return nil, errors.New("invalid file id type given")
	}

	// First get the identifier length.
	res, err := p.Module.ExportedFunction("FPDF_GetFileIdentifier").Call(p.Context, *documentHandle.handle, *(*uint64)(unsafe.Pointer(&request.FileIdType)), 0, 0)
	if err != nil {
		return nil, err
	}

	identifierSize := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if identifierSize == 0 {
		return &responses.FPDF_GetFileIdentifier{
			FileIdType: request.FileIdType,
			Identifier: nil,
		}, nil
	}

	charDataPointer, err := p.ByteArrayPointer(identifierSize, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDF_GetFileIdentifier").Call(p.Context, *documentHandle.handle, *(*uint64)(unsafe.Pointer(&request.FileIdType)), charDataPointer.Pointer, identifierSize)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(true)
	if err != nil {
		return nil, err
	}

	// Remove NULL terminator.
	charData = bytes.TrimSuffix(charData, []byte("\x00"))

	return &responses.FPDF_GetFileIdentifier{
		FileIdType: request.FileIdType,
		Identifier: charData,
	}, nil
}

// FPDFBookmark_GetCount returns the number of children of a bookmark.
// Experimental API.
func (p *PdfiumImplementation) FPDFBookmark_GetCount(request *requests.FPDFBookmark_GetCount) (*responses.FPDFBookmark_GetCount, error) {
	p.Lock()
	defer p.Unlock()

	bookmarkHandle, err := p.getBookmarkHandle(request.Bookmark)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFLink_CountQuadPoints").Call(p.Context, *bookmarkHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFBookmark_GetCount{
		Count: int(count),
	}, nil
}
