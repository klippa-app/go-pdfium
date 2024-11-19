package implementation_webassembly

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"
)

// FPDFText_LoadPage returns a handle to the text page information structure.
// Application must call FPDFText_ClosePage to release the text page
func (p *PdfiumImplementation) FPDFText_LoadPage(request *requests.FPDFText_LoadPage) (*responses.FPDFText_LoadPage, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_LoadPage").Call(p.Context, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	textPage := res[0]
	if textPage == 0 {
		return nil, errors.New("could not load text page")
	}

	documentHandle, err := p.getDocumentHandle(pageHandle.documentRef)
	if err != nil {
		return nil, err
	}

	textPageHandle := p.registerTextPage(&textPage, documentHandle)

	return &responses.FPDFText_LoadPage{
		TextPage: textPageHandle.nativeRef,
	}, nil
}

// FPDFText_ClosePage Release all resources allocated for a text page information structure.
func (p *PdfiumImplementation) FPDFText_ClosePage(request *requests.FPDFText_ClosePage) (*responses.FPDFText_ClosePage, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	documentHandle, err := p.getDocumentHandle(textPageHandle.documentRef)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDFText_ClosePage").Call(p.Context, *textPageHandle.handle)
	if err != nil {
		return nil, err
	}

	// Cleanup refs
	delete(p.textPageRefs, textPageHandle.nativeRef)
	delete(documentHandle.textPageRefs, textPageHandle.nativeRef)

	return &responses.FPDFText_ClosePage{}, nil
}

// FPDFText_CountChars returns the number of characters in a page.
// Characters in a page form a "stream", inside the stream, each character has an index.
// We will use the index parameters in many of FPDFTEXT functions. The first
// character in the page has an index value of zero.
func (p *PdfiumImplementation) FPDFText_CountChars(request *requests.FPDFText_CountChars) (*responses.FPDFText_CountChars, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_CountChars").Call(p.Context, *textPageHandle.handle)
	if err != nil {
		return nil, err
	}

	charCount := *(*int32)(unsafe.Pointer(&res[0]))
	if int(charCount) == -1 {
		return nil, errors.New("could not get char count")
	}

	return &responses.FPDFText_CountChars{
		Count: int(charCount),
	}, nil
}

// FPDFText_GetUnicode returns the unicode of a character in a page.
func (p *PdfiumImplementation) FPDFText_GetUnicode(request *requests.FPDFText_GetUnicode) (*responses.FPDFText_GetUnicode, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_GetUnicode").Call(p.Context, *textPageHandle.handle, uint64(request.Index))
	if err != nil {
		return nil, err
	}

	charUnicode := *(*uint)(unsafe.Pointer(&res[0]))

	return &responses.FPDFText_GetUnicode{
		Index:   request.Index,
		Unicode: charUnicode,
	}, nil
}

// FPDFText_GetFontSize returns the font size of a particular character.
func (p *PdfiumImplementation) FPDFText_GetFontSize(request *requests.FPDFText_GetFontSize) (*responses.FPDFText_GetFontSize, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_GetFontSize").Call(p.Context, *textPageHandle.handle, uint64(request.Index))
	if err != nil {
		return nil, err
	}

	fontSize := *(*float64)(unsafe.Pointer(&res[0]))

	return &responses.FPDFText_GetFontSize{
		Index:    request.Index,
		FontSize: fontSize,
	}, nil
}

// FPDFText_GetCharBox returns the bounding box of a particular character.
// All positions are measured in PDF "user space".
func (p *PdfiumImplementation) FPDFText_GetCharBox(request *requests.FPDFText_GetCharBox) (*responses.FPDFText_GetCharBox, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	leftPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer leftPointer.Free()

	topPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer topPointer.Free()

	rightPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer rightPointer.Free()

	bottomPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer bottomPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFText_GetCharBox").Call(p.Context, *textPageHandle.handle, uint64(request.Index), leftPointer.Pointer, rightPointer.Pointer, bottomPointer.Pointer, topPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get char box")
	}

	left, err := leftPointer.Value()
	if err != nil {
		return nil, err
	}

	top, err := topPointer.Value()
	if err != nil {
		return nil, err
	}

	right, err := rightPointer.Value()
	if err != nil {
		return nil, err
	}

	bottom, err := bottomPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFText_GetCharBox{
		Index:  request.Index,
		Left:   left,
		Right:  right,
		Bottom: bottom,
		Top:    top,
	}, nil
}

// FPDFText_GetCharOrigin returns origin of a particular character.
// All positions are measured in PDF "user space".
func (p *PdfiumImplementation) FPDFText_GetCharOrigin(request *requests.FPDFText_GetCharOrigin) (*responses.FPDFText_GetCharOrigin, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	xPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer xPointer.Free()

	yPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer yPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFText_GetCharOrigin").Call(p.Context, *textPageHandle.handle, uint64(request.Index), xPointer.Pointer, yPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get char origin")
	}

	x, err := xPointer.Value()
	if err != nil {
		return nil, err
	}

	y, err := yPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFText_GetCharOrigin{
		Index: request.Index,
		X:     x,
		Y:     y,
	}, nil
}

// FPDFText_GetCharIndexAtPos returns the index of a character at or nearby a certain position on the page.
func (p *PdfiumImplementation) FPDFText_GetCharIndexAtPos(request *requests.FPDFText_GetCharIndexAtPos) (*responses.FPDFText_GetCharIndexAtPos, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_GetCharIndexAtPos").Call(p.Context, *textPageHandle.handle, *(*uint64)(unsafe.Pointer(&request.X)), *(*uint64)(unsafe.Pointer(&request.Y)), *(*uint64)(unsafe.Pointer(&request.XTolerance)), *(*uint64)(unsafe.Pointer(&request.YTolerance)))
	if err != nil {
		return nil, err
	}

	charIndex := *(*int32)(unsafe.Pointer(&res[0]))
	if int(charIndex) == -3 {
		return nil, errors.New("could not get char index at pos")
	}

	return &responses.FPDFText_GetCharIndexAtPos{
		CharIndex: int(charIndex),
	}, nil
}

// FPDFText_GetText extracts unicode text string from the page.
// This function ignores characters without unicode information.
// It returns all characters on the page, even those that are not
// visible when the page has a cropbox. To filter out the characters
// outside of the cropbox, use FPDF_GetPageBoundingBox() and
// FPDFText_GetCharBox().
func (p *PdfiumImplementation) FPDFText_GetText(request *requests.FPDFText_GetText) (*responses.FPDFText_GetText, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	charDataPointer, err := p.ByteArrayPointer(uint64((request.Count+1)*2), nil) // UTF16-LE max 2 bytes per char, add 1 char for terminator.
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFText_GetText").Call(p.Context, *textPageHandle.handle, uint64(request.StartIndex), uint64(request.Count), charDataPointer.Pointer)
	if err != nil {
		return nil, err
	}
	charsWritten := *(*int32)(unsafe.Pointer(&res[0]))

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF16LEToUTF8(charData[0 : charsWritten*2])
	if err != nil {
		return nil, err
	}

	return &responses.FPDFText_GetText{
		Text: transformedText,
	}, nil
}

// FPDFText_CountRects returns the count of rectangular areas occupied by
// a segment of text, and caches the result for subsequent FPDFText_GetRect() calls.
// This function, along with FPDFText_GetRect can be used by
// applications to detect the position on the page for a text segment,
// so proper areas can be highlighted. The FPDFText_* functions will
// automatically merge small character boxes into bigger one if those
// characters are on the same line and use same font settings.
func (p *PdfiumImplementation) FPDFText_CountRects(request *requests.FPDFText_CountRects) (*responses.FPDFText_CountRects, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_CountRects").Call(p.Context, *textPageHandle.handle, uint64(request.StartIndex), uint64(request.Count))
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFText_CountRects{
		Count: int(count),
	}, nil
}

// FPDFText_GetRect returns a rectangular area from the result generated by FPDFText_CountRects
// Note: this method only works if you called FPDFText_CountRects first.
func (p *PdfiumImplementation) FPDFText_GetRect(request *requests.FPDFText_GetRect) (*responses.FPDFText_GetRect, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	leftPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer leftPointer.Free()

	topPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer topPointer.Free()

	rightPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer rightPointer.Free()

	bottomPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer bottomPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFText_GetRect").Call(p.Context, *textPageHandle.handle, uint64(request.Index), leftPointer.Pointer, topPointer.Pointer, rightPointer.Pointer, bottomPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))

	if int(success) == 0 {
		return nil, errors.New("could not get rect")
	}

	left, err := leftPointer.Value()
	if err != nil {
		return nil, err
	}

	top, err := topPointer.Value()
	if err != nil {
		return nil, err
	}

	right, err := rightPointer.Value()
	if err != nil {
		return nil, err
	}

	bottom, err := bottomPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFText_GetRect{
		Left:   left,
		Top:    top,
		Right:  right,
		Bottom: bottom,
	}, nil
}

// FPDFText_GetBoundedText extract unicode text within a rectangular boundary on the page.
func (p *PdfiumImplementation) FPDFText_GetBoundedText(request *requests.FPDFText_GetBoundedText) (*responses.FPDFText_GetBoundedText, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_GetBoundedText").Call(p.Context, *textPageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Left)), *(*uint64)(unsafe.Pointer(&request.Top)), *(*uint64)(unsafe.Pointer(&request.Right)), *(*uint64)(unsafe.Pointer(&request.Bottom)), 0, 0)
	if err != nil {
		return nil, err
	}

	charCount := *(*int32)(unsafe.Pointer(&res[0]))
	if int(charCount) == 0 {
		return &responses.FPDFText_GetBoundedText{}, nil
	}

	charDataPointer, err := p.ByteArrayPointer(uint64((charCount+1)*2), nil) // UTF16-LE max 2 bytes per char, add 1 char for terminator.
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDFText_GetBoundedText").Call(p.Context, *textPageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Left)), *(*uint64)(unsafe.Pointer(&request.Top)), *(*uint64)(unsafe.Pointer(&request.Right)), *(*uint64)(unsafe.Pointer(&request.Bottom)), charDataPointer.Pointer, uint64(charCount))
	if err != nil {
		return nil, err
	}

	charsWritten := *(*int32)(unsafe.Pointer(&res[0]))

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF16LEToUTF8(charData[0 : charsWritten*2])
	if err != nil {
		return nil, err
	}

	return &responses.FPDFText_GetBoundedText{
		Text: transformedText,
	}, nil
}

// FPDFText_FindStart returns a handle to search a page.
func (p *PdfiumImplementation) FPDFText_FindStart(request *requests.FPDFText_FindStart) (*responses.FPDFText_FindStart, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	documentHandle, err := p.getDocumentHandle(textPageHandle.documentRef)
	if err != nil {
		return nil, err
	}

	transformedTextPointer, err := p.CFPDF_WIDESTRING(request.Find)
	if err != nil {
		return nil, err
	}
	defer transformedTextPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFText_FindStart").Call(p.Context, *textPageHandle.handle, transformedTextPointer.Pointer, uint64(request.Flags), uint64(request.StartIndex))
	if err != nil {
		return nil, err
	}

	search := res[0]
	if search == 0 {
		return nil, errors.New("could not start search")
	}

	searchHandle := p.registerSearch(&search, documentHandle)

	return &responses.FPDFText_FindStart{
		Search: searchHandle.nativeRef,
	}, nil
}

// FPDFText_FindNext searches in the direction from page start to end.
func (p *PdfiumImplementation) FPDFText_FindNext(request *requests.FPDFText_FindNext) (*responses.FPDFText_FindNext, error) {
	p.Lock()
	defer p.Unlock()

	searchHandle, err := p.getSearchHandle(request.Search)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_FindNext").Call(p.Context, *searchHandle.handle)
	if err != nil {
		return nil, err
	}

	match := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFText_FindNext{
		GotMatch: int(match) == 1,
	}, nil
}

// FPDFText_FindPrev searches in the direction from page end to start.
func (p *PdfiumImplementation) FPDFText_FindPrev(request *requests.FPDFText_FindPrev) (*responses.FPDFText_FindPrev, error) {
	p.Lock()
	defer p.Unlock()

	searchHandle, err := p.getSearchHandle(request.Search)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_FindPrev").Call(p.Context, *searchHandle.handle)
	if err != nil {
		return nil, err
	}

	match := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFText_FindPrev{
		GotMatch: int(match) == 1,
	}, nil
}

// FPDFText_GetSchResultIndex returns the starting character index of the search result.
func (p *PdfiumImplementation) FPDFText_GetSchResultIndex(request *requests.FPDFText_GetSchResultIndex) (*responses.FPDFText_GetSchResultIndex, error) {
	p.Lock()
	defer p.Unlock()

	searchHandle, err := p.getSearchHandle(request.Search)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_GetSchResultIndex").Call(p.Context, *searchHandle.handle)
	if err != nil {
		return nil, err
	}

	index := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFText_GetSchResultIndex{
		Index: int(index),
	}, nil
}

// FPDFText_GetSchCount returns the number of matched characters in the search result.
func (p *PdfiumImplementation) FPDFText_GetSchCount(request *requests.FPDFText_GetSchCount) (*responses.FPDFText_GetSchCount, error) {
	p.Lock()
	defer p.Unlock()

	searchHandle, err := p.getSearchHandle(request.Search)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_GetSchCount").Call(p.Context, *searchHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFText_GetSchCount{
		Count: int(count),
	}, nil
}

// FPDFText_FindClose releases a search context.
func (p *PdfiumImplementation) FPDFText_FindClose(request *requests.FPDFText_FindClose) (*responses.FPDFText_FindClose, error) {
	p.Lock()
	defer p.Unlock()

	searchHandle, err := p.getSearchHandle(request.Search)
	if err != nil {
		return nil, err
	}

	documentHandle, err := p.getDocumentHandle(searchHandle.documentRef)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDFText_FindClose").Call(p.Context, *searchHandle.handle)
	if err != nil {
		return nil, err
	}

	// Cleanup refs
	delete(p.searchRefs, searchHandle.nativeRef)
	delete(documentHandle.searchRefs, searchHandle.nativeRef)

	return &responses.FPDFText_FindClose{}, nil
}

// FPDFLink_LoadWebLinks prepares information about weblinks in a page.
// Weblinks are those links implicitly embedded in PDF pages. PDF also
// has a type of annotation called "link" (FPDFTEXT doesn't deal with
// that kind of link). FPDFTEXT weblink feature is useful for
// automatically detecting links in the page contents. For example,
// things like "https://www.example.com" will be detected, so
// applications can allow user to click on those characters to activate
// the link, even the PDF doesn't come with link annotations.
//
// FPDFLink_CloseWebLinks must be called to release resources.
func (p *PdfiumImplementation) FPDFLink_LoadWebLinks(request *requests.FPDFLink_LoadWebLinks) (*responses.FPDFLink_LoadWebLinks, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFLink_LoadWebLinks").Call(p.Context, *textPageHandle.handle)
	if err != nil {
		return nil, err
	}

	pageLink := res[0]
	if pageLink == 0 {
		return nil, errors.New("could not load web links")
	}

	documentHandle, err := p.getDocumentHandle(textPageHandle.documentRef)
	if err != nil {
		return nil, err
	}

	pageLinkHandle := p.registerPageLink(&pageLink, documentHandle)

	return &responses.FPDFLink_LoadWebLinks{
		PageLink: pageLinkHandle.nativeRef,
	}, nil
}

// FPDFLink_CountWebLinks returns the count of detected web links.
func (p *PdfiumImplementation) FPDFLink_CountWebLinks(request *requests.FPDFLink_CountWebLinks) (*responses.FPDFLink_CountWebLinks, error) {
	p.Lock()
	defer p.Unlock()

	pageLinkhandle, err := p.getPageLinkHandle(request.PageLink)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFLink_CountWebLinks").Call(p.Context, *pageLinkhandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFLink_CountWebLinks{
		Count: int(count),
	}, nil
}

// FPDFLink_GetURL returns the URL information for a detected web link.
func (p *PdfiumImplementation) FPDFLink_GetURL(request *requests.FPDFLink_GetURL) (*responses.FPDFLink_GetURL, error) {
	p.Lock()
	defer p.Unlock()

	pageLinkhandle, err := p.getPageLinkHandle(request.PageLink)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFLink_GetURL").Call(p.Context, *pageLinkhandle.handle, uint64(request.Index), 0, 0)
	if err != nil {
		return nil, err
	}

	charCount := *(*int32)(unsafe.Pointer(&res[0]))

	charDataPointer, err := p.ByteArrayPointer(uint64((int(charCount)+1)*2), nil) // UTF16-LE max 2 bytes per char, add 1 char for terminator.
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDFLink_GetURL").Call(p.Context, *pageLinkhandle.handle, uint64(request.Index), charDataPointer.Pointer, uint64(charCount))
	if err != nil {
		return nil, err
	}
	charsWritten := *(*int32)(unsafe.Pointer(&res[0]))

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF16LEToUTF8(charData[0 : charsWritten*2])
	if err != nil {
		return nil, err
	}

	return &responses.FPDFLink_GetURL{
		Index: request.Index,
		URL:   transformedText,
	}, nil
}

// FPDFLink_CountRects returns the count of rectangular areas for the link.
func (p *PdfiumImplementation) FPDFLink_CountRects(request *requests.FPDFLink_CountRects) (*responses.FPDFLink_CountRects, error) {
	p.Lock()
	defer p.Unlock()

	pageLinkhandle, err := p.getPageLinkHandle(request.PageLink)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFLink_CountRects").Call(p.Context, *pageLinkhandle.handle, uint64(request.Index))
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFLink_CountRects{
		Count: int(count),
	}, nil
}

// FPDFLink_GetRect returns the boundaries of a rectangle for a link.
func (p *PdfiumImplementation) FPDFLink_GetRect(request *requests.FPDFLink_GetRect) (*responses.FPDFLink_GetRect, error) {
	p.Lock()
	defer p.Unlock()

	pageLinkhandle, err := p.getPageLinkHandle(request.PageLink)
	if err != nil {
		return nil, err
	}

	leftPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer leftPointer.Free()

	topPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer topPointer.Free()

	rightPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer rightPointer.Free()

	bottomPointer, err := p.DoublePointer(nil)
	if err != nil {
		return nil, err
	}
	defer bottomPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFLink_GetRect").Call(p.Context, *pageLinkhandle.handle, uint64(request.Index), uint64(request.RectIndex), leftPointer.Pointer, topPointer.Pointer, rightPointer.Pointer, bottomPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))

	if int(success) == 0 {
		return nil, errors.New("could not get rect")
	}

	left, err := leftPointer.Value()
	if err != nil {
		return nil, err
	}

	top, err := topPointer.Value()
	if err != nil {
		return nil, err
	}

	right, err := rightPointer.Value()
	if err != nil {
		return nil, err
	}

	bottom, err := bottomPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFLink_GetRect{
		Index:     request.Index,
		RectIndex: request.RectIndex,
		Left:      left,
		Top:       top,
		Right:     right,
		Bottom:    bottom,
	}, nil
}

// FPDFLink_CloseWebLinks releases resources used by weblink feature.
func (p *PdfiumImplementation) FPDFLink_CloseWebLinks(request *requests.FPDFLink_CloseWebLinks) (*responses.FPDFLink_CloseWebLinks, error) {
	p.Lock()
	defer p.Unlock()

	pageLinkhandle, err := p.getPageLinkHandle(request.PageLink)
	if err != nil {
		return nil, err
	}

	documentHandle, err := p.getDocumentHandle(pageLinkhandle.documentRef)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDFLink_CloseWebLinks").Call(p.Context, *pageLinkhandle.handle)
	if err != nil {
		return nil, err
	}

	// Cleanup refs
	delete(p.pageLinkRefs, pageLinkhandle.nativeRef)
	delete(documentHandle.pageLinkRefs, pageLinkhandle.nativeRef)

	return &responses.FPDFLink_CloseWebLinks{}, nil
}

// FPDFText_GetFontInfo returns the font name and flags of a particular character.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetFontInfo(request *requests.FPDFText_GetFontInfo) (*responses.FPDFText_GetFontInfo, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	// First get the font name length.
	res, err := p.Module.ExportedFunction("FPDFText_GetFontInfo").Call(p.Context, *textPageHandle.handle, uint64(request.Index), 0, 0, 0)
	if err != nil {
		return nil, err
	}

	fontNameSize := *(*int32)(unsafe.Pointer(&res[0]))
	if fontNameSize == 0 {
		return nil, errors.New("could not get font name")
	}

	fontFlagsPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer fontFlagsPointer.Free()

	charDataPointer, err := p.ByteArrayPointer(uint64(fontNameSize), nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFText_GetFontInfo").Call(p.Context, *textPageHandle.handle, uint64(request.Index), charDataPointer.Pointer, uint64(fontNameSize), fontFlagsPointer.Pointer)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

	fontFlags, err := fontFlagsPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFText_GetFontInfo{
		Index:    request.Index,
		FontName: string(charData[:fontNameSize-1]), // Remove NULL terminator
		Flags:    int(fontFlags),
	}, nil
}

// FPDFText_GetFontWeight returns the font weight of a particular character.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetFontWeight(request *requests.FPDFText_GetFontWeight) (*responses.FPDFText_GetFontWeight, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_GetFontWeight").Call(p.Context, *textPageHandle.handle, uint64(request.Index))
	if err != nil {
		return nil, err
	}

	fontWeight := *(*int32)(unsafe.Pointer(&res[0]))
	if int(fontWeight) == -1 {
		return nil, errors.New("could not get font weight")
	}

	return &responses.FPDFText_GetFontWeight{
		Index:      request.Index,
		FontWeight: int(fontWeight),
	}, nil
}

// FPDFText_GetFillColor returns the fill color of a particular character.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetFillColor(request *requests.FPDFText_GetFillColor) (*responses.FPDFText_GetFillColor, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	rPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer rPointer.Free()

	gPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer gPointer.Free()

	bPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer bPointer.Free()

	aPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer aPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFText_GetFillColor").Call(p.Context, *textPageHandle.handle, uint64(request.Index), rPointer.Pointer, gPointer.Pointer, bPointer.Pointer, aPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get fill color")
	}

	r, err := rPointer.Value()
	if err != nil {
		return nil, err
	}

	g, err := gPointer.Value()
	if err != nil {
		return nil, err
	}

	b, err := bPointer.Value()
	if err != nil {
		return nil, err
	}

	a, err := aPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFText_GetFillColor{
		Index: request.Index,
		R:     r,
		G:     g,
		B:     b,
		A:     a,
	}, nil
}

// FPDFText_GetStrokeColor returns the stroke color of a particular character.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetStrokeColor(request *requests.FPDFText_GetStrokeColor) (*responses.FPDFText_GetStrokeColor, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	rPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer rPointer.Free()

	gPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer gPointer.Free()

	bPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer bPointer.Free()

	aPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer aPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFText_GetStrokeColor").Call(p.Context, *textPageHandle.handle, uint64(request.Index), rPointer.Pointer, gPointer.Pointer, bPointer.Pointer, aPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get stroke color")
	}

	r, err := rPointer.Value()
	if err != nil {
		return nil, err
	}

	g, err := gPointer.Value()
	if err != nil {
		return nil, err
	}

	b, err := bPointer.Value()
	if err != nil {
		return nil, err
	}

	a, err := aPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFText_GetStrokeColor{
		Index: request.Index,
		R:     r,
		G:     g,
		B:     b,
		A:     a,
	}, nil
}

// FPDFText_GetCharAngle returns the character rotation angle.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetCharAngle(request *requests.FPDFText_GetCharAngle) (*responses.FPDFText_GetCharAngle, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_GetCharAngle").Call(p.Context, *textPageHandle.handle, uint64(request.Index))
	if err != nil {
		return nil, err
	}

	charAngle := *(*float32)(unsafe.Pointer(&res[0]))
	if float64(charAngle) == -1 {
		return nil, errors.New("could not get char angle")
	}

	return &responses.FPDFText_GetCharAngle{
		Index:     request.Index,
		CharAngle: float32(charAngle),
	}, nil
}

// FPDFText_GetLooseCharBox returns a "loose" bounding box of a particular character, i.e., covering
// the entire glyph bounds, without taking the actual glyph shape into
// account. All positions are measured in PDF "user space".
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetLooseCharBox(request *requests.FPDFText_GetLooseCharBox) (*responses.FPDFText_GetLooseCharBox, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	rectPointer, rectValue, err := p.CStructFS_RECTF(nil)
	if err != nil {
		return nil, err
	}
	defer p.Free(rectPointer)

	res, err := p.Module.ExportedFunction("FPDFText_GetLooseCharBox").Call(p.Context, *textPageHandle.handle, uint64(request.Index), rectPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get loose char box")
	}

	rect, err := rectValue()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFText_GetLooseCharBox{
		Rect: structs.FPDF_FS_RECTF{
			Left:   float32(rect.Left),
			Top:    float32(rect.Top),
			Right:  float32(rect.Right),
			Bottom: float32(rect.Bottom),
		},
	}, nil
}

// FPDFText_GetMatrix returns the effective transformation matrix for a particular character.
// All positions are measured in PDF "user space".
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetMatrix(request *requests.FPDFText_GetMatrix) (*responses.FPDFText_GetMatrix, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	matrixPointer, matrixValue, err := p.CStructFS_MATRIX(nil)
	if err != nil {
		return nil, err
	}
	defer p.Free(matrixPointer)

	res, err := p.Module.ExportedFunction("FPDFText_GetMatrix").Call(p.Context, *textPageHandle.handle, uint64(request.Index), matrixPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get char matrix")
	}

	matrix, err := matrixValue()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFText_GetMatrix{
		Matrix: structs.FPDF_FS_MATRIX{
			A: float32(matrix.A),
			B: float32(matrix.B),
			C: float32(matrix.C),
			D: float32(matrix.D),
			E: float32(matrix.E),
			F: float32(matrix.F),
		},
	}, nil
}

// FPDFLink_GetTextRange returns the start char index and char count for a link.
// Experimental API.
func (p *PdfiumImplementation) FPDFLink_GetTextRange(request *requests.FPDFLink_GetTextRange) (*responses.FPDFLink_GetTextRange, error) {
	p.Lock()
	defer p.Unlock()

	pageLinkhandle, err := p.getPageLinkHandle(request.PageLink)
	if err != nil {
		return nil, err
	}

	startCharIndexPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer startCharIndexPointer.Free()

	charCountPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer charCountPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFLink_GetTextRange").Call(p.Context, *pageLinkhandle.handle, uint64(request.Index), startCharIndexPointer.Pointer, charCountPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get text range")
	}

	startCharIndex, err := startCharIndexPointer.Value()
	if err != nil {
		return nil, err
	}

	charCount, err := charCountPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFLink_GetTextRange{
		Index:          request.Index,
		StartCharIndex: int(startCharIndex),
		CharCount:      int(charCount),
	}, nil
}

// FPDFText_GetTextObject returns the FPDF_PAGEOBJECT associated with a given character.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_GetTextObject(request *requests.FPDFText_GetTextObject) (*responses.FPDFText_GetTextObject, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_GetTextObject").Call(p.Context, *textPageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	if res[0] == 0 {
		return nil, errors.New("could not get object")
	}

	pageObject := &res[0]
	pageObjectHandle := p.registerPageObject(pageObject)

	return &responses.FPDFText_GetTextObject{
		TextObject: pageObjectHandle.nativeRef,
		Index:      request.Index,
	}, nil
}

// FPDFText_IsGenerated returns whether a character in a page is generated by PDFium.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_IsGenerated(request *requests.FPDFText_IsGenerated) (*responses.FPDFText_IsGenerated, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_IsGenerated").Call(p.Context, *textPageHandle.handle, uint64(request.Index))
	if err != nil {
		return nil, err
	}

	isGenerated := *(*int32)(unsafe.Pointer(&res[0]))
	if int(isGenerated) == -1 {
		return nil, errors.New("could not get whether text is generated")
	}

	return &responses.FPDFText_IsGenerated{
		Index:       request.Index,
		IsGenerated: int(isGenerated) == 1,
	}, nil
}

// FPDFText_IsHyphen returns whether a character in a page is a hyphen.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_IsHyphen(request *requests.FPDFText_IsHyphen) (*responses.FPDFText_IsHyphen, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_IsHyphen").Call(p.Context, *textPageHandle.handle, uint64(request.Index))
	if err != nil {
		return nil, err
	}

	isHyphen := *(*int32)(unsafe.Pointer(&res[0]))
	if int(isHyphen) == -1 {
		return nil, errors.New("could not get whether text is a hyphen")
	}

	return &responses.FPDFText_IsHyphen{
		Index:    request.Index,
		IsHyphen: int(isHyphen) == 1,
	}, nil
}

// FPDFText_HasUnicodeMapError a character in a page has an invalid unicode mapping.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_HasUnicodeMapError(request *requests.FPDFText_HasUnicodeMapError) (*responses.FPDFText_HasUnicodeMapError, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_HasUnicodeMapError").Call(p.Context, *textPageHandle.handle, uint64(request.Index))
	if err != nil {
		return nil, err
	}

	hasUnicodeMapError := *(*int32)(unsafe.Pointer(&res[0]))
	if int(hasUnicodeMapError) == -1 {
		return nil, errors.New("could not get whether text has a unicode map error")
	}

	return &responses.FPDFText_HasUnicodeMapError{
		Index:              request.Index,
		HasUnicodeMapError: int(hasUnicodeMapError) == 1,
	}, nil
}
