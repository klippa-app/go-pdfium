package implementation_webassembly

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFText_GetCharIndexFromTextIndex returns the character index in the text page internal character list.
// Where the character index is an index of the text returned from FPDFText_GetText().
func (p *PdfiumImplementation) FPDFText_GetCharIndexFromTextIndex(request *requests.FPDFText_GetCharIndexFromTextIndex) (*responses.FPDFText_GetCharIndexFromTextIndex, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_GetCharIndexFromTextIndex").Call(p.Context, *textPageHandle.handle, *(*uint64)(unsafe.Pointer(&request.NTextIndex)))
	if err != nil {
		return nil, err
	}

	charIndex := *(*int32)(unsafe.Pointer(&res[0]))
	if int(charIndex) == -1 {
		return nil, errors.New("could not get char index")
	}

	return &responses.FPDFText_GetCharIndexFromTextIndex{
		CharIndex: int(charIndex),
	}, nil
}

// FPDFText_GetTextIndexFromCharIndex returns the text index in the text page internal character list.
// Where the text index is an index of the character in the internal character list.
func (p *PdfiumImplementation) FPDFText_GetTextIndexFromCharIndex(request *requests.FPDFText_GetTextIndexFromCharIndex) (*responses.FPDFText_GetTextIndexFromCharIndex, error) {
	p.Lock()
	defer p.Unlock()

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFText_GetTextIndexFromCharIndex").Call(p.Context, *textPageHandle.handle, *(*uint64)(unsafe.Pointer(&request.NCharIndex)))
	if err != nil {
		return nil, err
	}

	textIndex := *(*int32)(unsafe.Pointer(&res[0]))
	if int(textIndex) == -1 {
		return nil, errors.New("could not get text index")
	}

	return &responses.FPDFText_GetTextIndexFromCharIndex{
		TextIndex: int(textIndex),
	}, nil
}
