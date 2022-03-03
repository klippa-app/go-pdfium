package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_doc.h"
import "C"
import (
	"errors"
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"unsafe"
)

func (p *PdfiumImplementation) registerAction(action C.FPDF_ACTION) *ActionHandle {
	ref := uuid.New()
	handle := &ActionHandle{
		handle:    action,
		nativeRef: references.FPDF_ACTION(ref.String()),
	}

	p.actionRefs[handle.nativeRef] = handle

	return handle
}

func (p *PdfiumImplementation) getActionInfo(actionHandle *ActionHandle, documentHandle *DocumentHandle) (*responses.ActionInfo, error) {
	actionType := C.FPDFAction_GetType(actionHandle.handle)

	actionInfo := &responses.ActionInfo{
		Reference: actionHandle.nativeRef,
		Type:      enums.FPDF_ACTION_ACTION(actionType),
	}

	if actionInfo.Type == enums.FPDF_ACTION_ACTION_GOTO {
		dest := C.FPDFAction_GetDest(documentHandle.handle, actionHandle.handle)
		if dest != nil {
			destHandle := p.registerDest(dest, documentHandle)

			destInfo, err := p.getDestInfo(destHandle, documentHandle)
			if err != nil {
				return nil, err
			}

			actionInfo.DestInfo = destInfo
		}
	} else if actionInfo.Type == enums.FPDF_ACTION_ACTION_LAUNCH || actionInfo.Type == enums.FPDF_ACTION_ACTION_REMOTEGOTO {
		// First get the file path length.
		filePathLength := C.FPDFAction_GetFilePath(actionHandle.handle, C.NULL, 0)
		if filePathLength == 0 {
			return nil, errors.New("Could not get file path")
		}

		charData := make([]byte, filePathLength)
		// FPDFAction_GetFilePath returns the data in UTF-8, no conversion needed.
		C.FPDFAction_GetFilePath(actionHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

		filePathString := string(charData[:filePathLength-1]) // Take of NULL terminator
		actionInfo.FilePath = &filePathString
	} else if actionInfo.Type == enums.FPDF_ACTION_ACTION_URI {
		// First get the uri path length.
		uriPathLength := C.FPDFAction_GetURIPath(documentHandle.handle, actionHandle.handle, C.NULL, 0)
		if uriPathLength == 0 {
			return nil, errors.New("Could not get uri path")
		}

		charData := make([]byte, uriPathLength)
		C.FPDFAction_GetURIPath(documentHandle.handle, actionHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

		uriPathString := string(charData[:uriPathLength-1]) // Take of NULL terminator
		actionInfo.URIPath = &uriPathString
	}

	return actionInfo, nil
}

func (p *PdfiumImplementation) GetActionInfo(request *requests.GetActionInfo) (*responses.GetActionInfo, error) {
	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	actionHandle, err := p.getActionHandle(request.Action)
	if err != nil {
		return nil, err
	}

	actionInfo, err := p.getActionInfo(actionHandle, documentHandle)
	if err != nil {
		return nil, err
	}

	return &responses.GetActionInfo{
		ActionInfo: *actionInfo,
	}, nil
}
