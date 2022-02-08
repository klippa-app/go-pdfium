//go:build pdfium_experimental
// +build pdfium_experimental

package implementation

/*
#cgo pkg-config: pdfium
#include "fpdf_javascript.h"
*/
import "C"
import (
	"errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"unsafe"
)

// FPDFDoc_GetJavaScriptActionCount returns the number of JavaScript actions in the given document.
// Experimental API.
func (p *PdfiumImplementation) FPDFDoc_GetJavaScriptActionCount(request *requests.FPDFDoc_GetJavaScriptActionCount) (*responses.FPDFDoc_GetJavaScriptActionCount, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	javaScriptActionCount := C.FPDFDoc_GetJavaScriptActionCount(documentHandle.handle)
	if int(javaScriptActionCount) == -1 {
		return nil, errors.New("could not get JavaScript Action count")
	}

	return &responses.FPDFDoc_GetJavaScriptActionCount{
		JavaScriptActionCount: int(javaScriptActionCount),
	}, nil
}

// FPDFDoc_GetJavaScriptAction returns the JavaScript action at the given index in the given document.
// Experimental API.
func (p *PdfiumImplementation) FPDFDoc_GetJavaScriptAction(request *requests.FPDFDoc_GetJavaScriptAction) (*responses.FPDFDoc_GetJavaScriptAction, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	javaScriptAction := C.FPDFDoc_GetJavaScriptAction(documentHandle.handle, C.int(request.Index))
	if javaScriptAction == nil {
		return nil, errors.New("could not get JavaScript Action")
	}

	javaScriptActionHandle := p.registerJavaScriptAction(javaScriptAction, documentHandle)

	return &responses.FPDFDoc_GetJavaScriptAction{
		Index:            request.Index,
		JavaScriptAction: javaScriptActionHandle.nativeRef,
	}, nil
}

// FPDFDoc_CloseJavaScriptAction closes a loaded FPDF_JAVASCRIPT_ACTION object.
// Experimental API.
func (p *PdfiumImplementation) FPDFDoc_CloseJavaScriptAction(request *requests.FPDFDoc_CloseJavaScriptAction) (*responses.FPDFDoc_CloseJavaScriptAction, error) {
	p.Lock()
	defer p.Unlock()

	javaScriptActionHandle, err := p.getJavaScriptActionHandle(request.JavaScriptAction)
	if err != nil {
		return nil, err
	}

	C.FPDFDoc_CloseJavaScriptAction(javaScriptActionHandle.handle)
	delete(p.javaScriptActionRefs, javaScriptActionHandle.nativeRef)

	documentHandle, err := p.getDocumentHandle(javaScriptActionHandle.documentRef)
	if err != nil {
		return nil, err
	}

	delete(documentHandle.javaScriptActionRefs, javaScriptActionHandle.nativeRef)

	return &responses.FPDFDoc_CloseJavaScriptAction{}, nil
}

// FPDFJavaScriptAction_GetName returns the name from the javascript handle.
// Experimental API.
func (p *PdfiumImplementation) FPDFJavaScriptAction_GetName(request *requests.FPDFJavaScriptAction_GetName) (*responses.FPDFJavaScriptAction_GetName, error) {
	p.Lock()
	defer p.Unlock()

	javaScriptActionHandle, err := p.getJavaScriptActionHandle(request.JavaScriptAction)
	if err != nil {
		return nil, err
	}

	// First get the name value length.
	nameSize := C.FPDFJavaScriptAction_GetName(javaScriptActionHandle.handle, nil, 0)
	if nameSize == 0 {
		return nil, errors.New("Could not get name")
	}

	charData := make([]byte, nameSize)
	C.FPDFJavaScriptAction_GetName(javaScriptActionHandle.handle, (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFJavaScriptAction_GetName{
		Name: transformedText,
	}, nil
}

// FPDFJavaScriptAction_GetScript returns the script from the javascript handle
// Experimental API.
func (p *PdfiumImplementation) FPDFJavaScriptAction_GetScript(request *requests.FPDFJavaScriptAction_GetScript) (*responses.FPDFJavaScriptAction_GetScript, error) {
	p.Lock()
	defer p.Unlock()

	javaScriptActionHandle, err := p.getJavaScriptActionHandle(request.JavaScriptAction)
	if err != nil {
		return nil, err
	}

	// First get the script value length.
	scriptSize := C.FPDFJavaScriptAction_GetScript(javaScriptActionHandle.handle, nil, 0)
	if scriptSize == 0 {
		return nil, errors.New("Could not get script")
	}

	charData := make([]byte, scriptSize)
	C.FPDFJavaScriptAction_GetScript(javaScriptActionHandle.handle, (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFJavaScriptAction_GetScript{
		Script: transformedText,
	}, nil
}
