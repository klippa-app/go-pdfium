package implementation_webassembly

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
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

	res, err := p.Module.ExportedFunction("FPDFDoc_GetJavaScriptActionCount").Call(p.Context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	javaScriptActionCount := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDFDoc_GetJavaScriptAction").Call(p.Context, *documentHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	javaScriptAction := res[0]
	if javaScriptAction == 0 {
		return nil, errors.New("could not get JavaScript Action")
	}

	javaScriptActionHandle := p.registerJavaScriptAction(&javaScriptAction, documentHandle)

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

	_, err = p.Module.ExportedFunction("FPDFDoc_CloseJavaScriptAction").Call(p.Context, *javaScriptActionHandle.handle)
	if err != nil {
		return nil, err
	}

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
	res, err := p.Module.ExportedFunction("FPDFJavaScriptAction_GetName").Call(p.Context, *javaScriptActionHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	nameSize := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if nameSize == 0 {
		return nil, errors.New("Could not get name")
	}

	charDataPointer, err := p.ByteArrayPointer(nameSize, nil)
	if err != nil {
		return nil, err
	}

	defer charDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDFJavaScriptAction_GetName").Call(p.Context, *javaScriptActionHandle.handle, charDataPointer.Pointer, nameSize)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)

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
	res, err := p.Module.ExportedFunction("FPDFJavaScriptAction_GetScript").Call(p.Context, *javaScriptActionHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	scriptSize := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if scriptSize == 0 {
		return nil, errors.New("Could not get script")
	}

	charDataPointer, err := p.ByteArrayPointer(scriptSize, nil)
	if err != nil {
		return nil, err
	}

	defer charDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDFJavaScriptAction_GetScript").Call(p.Context, *javaScriptActionHandle.handle, charDataPointer.Pointer, scriptSize)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFJavaScriptAction_GetScript{
		Script: transformedText,
	}, nil
}
