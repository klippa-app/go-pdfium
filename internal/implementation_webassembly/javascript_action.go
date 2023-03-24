package implementation_webassembly

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerJavaScriptAction(javaScriptAction *uint64, documentHandle *DocumentHandle) *JavaScriptActionHandle {
	ref := uuid.New()
	handle := &JavaScriptActionHandle{
		handle:      javaScriptAction,
		nativeRef:   references.FPDF_JAVASCRIPT_ACTION(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.javaScriptActionRefs[handle.nativeRef] = handle
	p.javaScriptActionRefs[handle.nativeRef] = handle

	return handle
}

// GetJavaScriptActions returns all the JavaScript Actions of a document.
// Experimental API.
func (p *PdfiumImplementation) GetJavaScriptActions(request *requests.GetJavaScriptActions) (*responses.GetJavaScriptActions, error) {
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

	cJavaScriptActionCount := *(*int32)(unsafe.Pointer(&res[0]))
	javaScriptActionCount := int(cJavaScriptActionCount)
	if int(javaScriptActionCount) == -1 {
		return nil, errors.New("could not get JavaScript Action count")
	}

	javaScriptActions := []responses.JavaScriptAction{}
	for i := 0; i < javaScriptActionCount; i++ {
		res, err = p.Module.ExportedFunction("FPDFDoc_GetJavaScriptAction").Call(p.Context, *documentHandle.handle, uint64(i))
		if err != nil {
			return nil, err
		}

		javaScriptAction := res[0]
		if javaScriptAction == 0 {
			continue
		}
		defer p.Module.ExportedFunction("FPDFDoc_CloseJavaScriptAction").Call(p.Context, javaScriptAction)

		// First get the name value length.
		res, err = p.Module.ExportedFunction("FPDFJavaScriptAction_GetName").Call(p.Context, javaScriptAction, 0, 0)
		if err != nil {
			return nil, err
		}

		nameSize := *(*int32)(unsafe.Pointer(&res[0]))
		if nameSize == 0 {
			return nil, errors.New("Could not get name")
		}

		charDataPointer, err := p.ByteArrayPointer(uint64(nameSize), nil)
		defer charDataPointer.Free()

		_, err = p.Module.ExportedFunction("FPDFJavaScriptAction_GetName").Call(p.Context, javaScriptAction, charDataPointer.Pointer, uint64(nameSize))
		if err != nil {
			return nil, err
		}

		charData, err := charDataPointer.Value(false)
		if err != nil {
			return nil, err
		}

		transformedName, err := p.transformUTF16LEToUTF8(charData)
		if err != nil {
			return nil, err
		}

		// First get the script value length.
		res, err = p.Module.ExportedFunction("FPDFJavaScriptAction_GetScript").Call(p.Context, javaScriptAction, 0, 0)
		if err != nil {
			return nil, err
		}

		scriptSize := *(*int32)(unsafe.Pointer(&res[0]))
		if scriptSize == 0 {
			return nil, errors.New("Could not get script")
		}

		charDataPointer, err = p.ByteArrayPointer(uint64(scriptSize), nil)
		defer charDataPointer.Free()

		_, err = p.Module.ExportedFunction("FPDFJavaScriptAction_GetScript").Call(p.Context, javaScriptAction, charDataPointer.Pointer, uint64(scriptSize))
		if err != nil {
			return nil, err
		}

		charData, err = charDataPointer.Value(false)
		if err != nil {
			return nil, err
		}

		transformedScript, err := p.transformUTF16LEToUTF8(charData)
		if err != nil {
			return nil, err
		}

		javaScriptActions = append(javaScriptActions, responses.JavaScriptAction{
			Name:   transformedName,
			Script: transformedScript,
		})
	}

	return &responses.GetJavaScriptActions{
		JavaScriptActions: javaScriptActions,
	}, nil
}
