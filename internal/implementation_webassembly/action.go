package implementation_webassembly

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerAction(action *uint64) *ActionHandle {
	ref := uuid.New()
	handle := &ActionHandle{
		handle:    action,
		nativeRef: references.FPDF_ACTION(ref.String()),
	}

	p.actionRefs[handle.nativeRef] = handle

	return handle
}

func (p *PdfiumImplementation) getActionInfo(actionHandle *ActionHandle, documentHandle *DocumentHandle) (*responses.ActionInfo, error) {
	res, err := p.Module.ExportedFunction("FPDFAction_GetType").Call(p.Context, *actionHandle.handle)
	if err != nil {
		return nil, err
	}

	actionType := *(*uint32)(unsafe.Pointer(&res[0]))

	actionInfo := &responses.ActionInfo{
		Reference: actionHandle.nativeRef,
		Type:      enums.FPDF_ACTION_ACTION(actionType),
	}

	if actionInfo.Type == enums.FPDF_ACTION_ACTION_GOTO {
		res, err = p.Module.ExportedFunction("FPDFAction_GetDest").Call(p.Context, *documentHandle.handle, *actionHandle.handle)
		if err != nil {
			return nil, err
		}

		dest := res[0]
		if dest != 0 {
			destHandle := p.registerDest(&dest, documentHandle)

			destInfo, err := p.getDestInfo(destHandle, documentHandle)
			if err != nil {
				return nil, err
			}

			actionInfo.DestInfo = destInfo
		}
	} else if actionInfo.Type == enums.FPDF_ACTION_ACTION_LAUNCH || actionInfo.Type == enums.FPDF_ACTION_ACTION_REMOTEGOTO {
		// First get the file path length.
		res, err = p.Module.ExportedFunction("FPDFAction_GetFilePath").Call(p.Context, *actionHandle.handle, 0, 0)
		if err != nil {
			return nil, err
		}

		filePathLength := *(*int32)(unsafe.Pointer(&res[0]))
		if filePathLength == 0 {
			return nil, errors.New("Could not get file path")
		}

		charDataPointer, err := p.ByteArrayPointer(uint64(filePathLength), nil)
		defer charDataPointer.Free()

		_, err = p.Module.ExportedFunction("FPDFAction_GetFilePath").Call(p.Context, *actionHandle.handle, charDataPointer.Pointer, uint64(filePathLength))
		if err != nil {
			return nil, err
		}

		charData, err := charDataPointer.Value(false)
		if err != nil {
			return nil, err
		}

		filePathString := string(charData[:filePathLength-1]) // Take of NULL terminator
		actionInfo.FilePath = &filePathString
	} else if actionInfo.Type == enums.FPDF_ACTION_ACTION_URI {
		// First get the uri path length.
		res, err = p.Module.ExportedFunction("FPDFAction_GetURIPath").Call(p.Context, *actionHandle.handle, 0, 0)
		if err != nil {
			return nil, err
		}

		uriPathLength := *(*int32)(unsafe.Pointer(&res[0]))
		if uriPathLength == 0 {
			return nil, errors.New("Could not get uri path")
		}

		charDataPointer, err := p.ByteArrayPointer(uint64(uriPathLength), nil)
		defer charDataPointer.Free()

		_, err = p.Module.ExportedFunction("FPDFAction_GetURIPath").Call(p.Context, *actionHandle.handle, charDataPointer.Pointer, uint64(uriPathLength))
		if err != nil {
			return nil, err
		}

		charData, err := charDataPointer.Value(false)
		if err != nil {
			return nil, err
		}

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
