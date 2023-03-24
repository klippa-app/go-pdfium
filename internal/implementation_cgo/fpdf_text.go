package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdf_text.h"
import "C"
import (
	"errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"unsafe"
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

	textPage := C.FPDFText_LoadPage(pageHandle.handle)
	if textPage == nil {
		return nil, errors.New("could not load text page")
	}

	documentHandle, err := p.getDocumentHandle(pageHandle.documentRef)
	if err != nil {
		return nil, err
	}

	textPageHandle := p.registerTextPage(textPage, documentHandle)

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

	C.FPDFText_ClosePage(textPageHandle.handle)

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

	charCount := C.FPDFText_CountChars(textPageHandle.handle)
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

	charUnicode := C.FPDFText_GetUnicode(textPageHandle.handle, C.int(request.Index))
	return &responses.FPDFText_GetUnicode{
		Index:   request.Index,
		Unicode: uint(charUnicode),
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

	fontSize := C.FPDFText_GetFontSize(textPageHandle.handle, C.int(request.Index))
	return &responses.FPDFText_GetFontSize{
		Index:    request.Index,
		FontSize: float64(fontSize),
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

	left := C.double(0)
	right := C.double(0)
	bottom := C.double(0)
	top := C.double(0)
	success := C.FPDFText_GetCharBox(textPageHandle.handle, C.int(request.Index), &left, &right, &bottom, &top)
	if int(success) == 0 {
		return nil, errors.New("could not get char box")
	}

	return &responses.FPDFText_GetCharBox{
		Index:  request.Index,
		Left:   float64(left),
		Right:  float64(right),
		Bottom: float64(bottom),
		Top:    float64(top),
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

	x := C.double(0)
	y := C.double(0)
	success := C.FPDFText_GetCharOrigin(textPageHandle.handle, C.int(request.Index), &x, &y)
	if int(success) == 0 {
		return nil, errors.New("could not get char origin")
	}

	return &responses.FPDFText_GetCharOrigin{
		Index: request.Index,
		X:     float64(x),
		Y:     float64(y),
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

	cCharIndex := C.FPDFText_GetCharIndexAtPos(textPageHandle.handle, C.double(request.X), C.double(request.Y), C.double(request.XTolerance), C.double(request.YTolerance))
	charIndex := int(cCharIndex)
	if charIndex == -3 {
		return nil, errors.New("could not get char index at pos")
	}

	return &responses.FPDFText_GetCharIndexAtPos{
		CharIndex: charIndex,
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

	charData := make([]byte, (request.Count+1)*2) // UTF16-LE max 2 bytes per char, add 1 char for terminator.
	charsWritten := C.FPDFText_GetText(textPageHandle.handle, C.int(request.StartIndex), C.int(request.Count), (*C.ushort)(unsafe.Pointer(&charData[0])))

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

	count := C.FPDFText_CountRects(textPageHandle.handle, C.int(request.StartIndex), C.int(request.Count))
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

	left := C.double(0)
	top := C.double(0)
	right := C.double(0)
	bottom := C.double(0)
	success := C.FPDFText_GetRect(textPageHandle.handle, C.int(request.Index), &left, &top, &right, &bottom)
	if int(success) == 0 {
		return nil, errors.New("could not get rect")
	}

	return &responses.FPDFText_GetRect{
		Left:   float64(left),
		Top:    float64(top),
		Right:  float64(right),
		Bottom: float64(bottom),
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

	charCount := C.FPDFText_GetBoundedText(textPageHandle.handle, C.double(request.Left), C.double(request.Top), C.double(request.Right), C.double(request.Bottom), nil, C.int(0))
	if int(charCount) == 0 {
		return &responses.FPDFText_GetBoundedText{}, nil
	}

	charData := make([]byte, (int(charCount)+1)*2) // UTF16-LE max 2 bytes per char, add 1 char for terminator.
	charsWritten := C.FPDFText_GetBoundedText(textPageHandle.handle, C.double(request.Left), C.double(request.Top), C.double(request.Right), C.double(request.Bottom), (*C.ushort)(unsafe.Pointer(&charData[0])), C.int(len(charData)))

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

	transformedText, err := p.transformUTF8ToUTF16LE(request.Find)
	if err != nil {
		return nil, err
	}

	search := C.FPDFText_FindStart(textPageHandle.handle, (C.FPDF_WIDESTRING)(unsafe.Pointer(&transformedText[0])), C.ulong(request.Flags), C.int(request.StartIndex))
	if search == nil {
		return nil, errors.New("could not start search")
	}

	searchHandle := p.registerSearch(search, documentHandle)

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

	match := C.FPDFText_FindNext(searchHandle.handle)
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

	match := C.FPDFText_FindPrev(searchHandle.handle)
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

	index := C.FPDFText_GetSchResultIndex(searchHandle.handle)
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

	count := C.FPDFText_GetSchCount(searchHandle.handle)
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

	C.FPDFText_FindClose(searchHandle.handle)

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

	pageLink := C.FPDFLink_LoadWebLinks(textPageHandle.handle)
	if pageLink == nil {
		return nil, errors.New("could not load web links")
	}

	documentHandle, err := p.getDocumentHandle(textPageHandle.documentRef)
	if err != nil {
		return nil, err
	}

	pageLinkHandle := p.registerPageLink(pageLink, documentHandle)

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

	count := C.FPDFLink_CountWebLinks(pageLinkhandle.handle)

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

	charCount := C.FPDFLink_GetURL(pageLinkhandle.handle, C.int(request.Index), nil, C.int(0))
	if int(charCount) == 0 {
		return &responses.FPDFLink_GetURL{}, nil
	}

	charData := make([]byte, (int(charCount)+1)*2) // UTF16-LE max 2 bytes per char, add 1 char for terminator.
	charsWritten := C.FPDFLink_GetURL(pageLinkhandle.handle, C.int(request.Index), (*C.ushort)(unsafe.Pointer(&charData[0])), C.int(len(charData)))

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

	count := C.FPDFLink_CountRects(pageLinkhandle.handle, C.int(request.Index))

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

	left := C.double(0)
	top := C.double(0)
	right := C.double(0)
	bottom := C.double(0)

	success := C.FPDFLink_GetRect(pageLinkhandle.handle, C.int(request.Index), C.int(request.RectIndex), &left, &top, &right, &bottom)
	if int(success) == 0 {
		return nil, errors.New("could not get rect")
	}

	return &responses.FPDFLink_GetRect{
		Index:     request.Index,
		RectIndex: request.RectIndex,
		Left:      float64(left),
		Top:       float64(top),
		Right:     float64(right),
		Bottom:    float64(bottom),
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

	C.FPDFLink_CloseWebLinks(pageLinkhandle.handle)

	// Cleanup refs
	delete(p.pageLinkRefs, pageLinkhandle.nativeRef)
	delete(documentHandle.pageLinkRefs, pageLinkhandle.nativeRef)

	return &responses.FPDFLink_CloseWebLinks{}, nil
}
