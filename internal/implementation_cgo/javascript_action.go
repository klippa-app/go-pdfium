package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdf_javascript.h"
import "C"
import (
	"errors"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"unsafe"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerJavaScriptAction(javaScriptAction C.FPDF_JAVASCRIPT_ACTION, documentHandle *DocumentHandle) *JavaScriptActionHandle {
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

	cJavaScriptActionCount := C.FPDFDoc_GetJavaScriptActionCount(documentHandle.handle)
	javaScriptActionCount := int(cJavaScriptActionCount)
	if int(javaScriptActionCount) == -1 {
		return nil, errors.New("could not get JavaScript Action count")
	}

	javaScriptActions := []responses.JavaScriptAction{}
	for i := 0; i < javaScriptActionCount; i++ {
		javaScriptAction := C.FPDFDoc_GetJavaScriptAction(documentHandle.handle, C.int(i))
		if javaScriptAction == nil {
			continue
		}
		defer C.FPDFDoc_CloseJavaScriptAction(javaScriptAction)

		// First get the name value length.
		nameSize := C.FPDFJavaScriptAction_GetName(javaScriptAction, nil, 0)
		if nameSize == 0 {
			return nil, errors.New("Could not get name")
		}

		charData := make([]byte, nameSize)
		C.FPDFJavaScriptAction_GetName(javaScriptAction, (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

		transformedName, err := p.transformUTF16LEToUTF8(charData)
		if err != nil {
			return nil, err
		}

		// First get the script value length.
		scriptSize := C.FPDFJavaScriptAction_GetScript(javaScriptAction, nil, 0)
		if scriptSize == 0 {
			return nil, errors.New("Could not get script")
		}

		charData = make([]byte, scriptSize)
		C.FPDFJavaScriptAction_GetScript(javaScriptAction, (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

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
