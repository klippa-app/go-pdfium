package implementation_webassembly

import (
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFDoc_GetPageMode returns the document's page mode, which describes how the document should be displayed when opened.
func (p *PdfiumImplementation) FPDFDoc_GetPageMode(request *requests.FPDFDoc_GetPageMode) (*responses.FPDFDoc_GetPageMode, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFDoc_GetPageMode").Call(p.Context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	pageMode := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFDoc_GetPageMode{
		PageMode: responses.FPDFDoc_GetPageModeMode(pageMode),
	}, nil
}

var CurrentUnsupportedObjectHandler requests.UnSpObjProcessHandler

// We need to keep around a reference so that it won't get GC'ed.
var currentUnsupportedObjectHandlerPointer *uint64

// FSDK_SetUnSpObjProcessHandler set ups an unsupported object handler.
// Since callbacks can't be transferred between processes in gRPC, you can only use this in single-threaded mode.
func (p *PdfiumImplementation) FSDK_SetUnSpObjProcessHandler(request *requests.FSDK_SetUnSpObjProcessHandler) (*responses.FSDK_SetUnSpObjProcessHandler, error) {
	p.Lock()
	defer p.Unlock()

	if currentUnsupportedObjectHandlerPointer == nil {
		res, err := p.Module.ExportedFunction("UNSUPPORT_INFO_Create").Call(p.Context)
		if err != nil {
			return nil, err
		}

		currentUnsupportedObjectHandlerPointer = &res[0]
	}

	CurrentUnsupportedObjectHandler = request.UnSpObjProcessHandler

	// Set the Go callback through cgo.
	_, err := p.Module.ExportedFunction("FSDK_SetUnSpObjProcessHandler").Call(p.Context, *currentUnsupportedObjectHandlerPointer)
	if err != nil {
		return nil, err
	}

	return &responses.FSDK_SetUnSpObjProcessHandler{}, nil
}

var CurrentTimeHandler requests.SetTimeFunction

// FSDK_SetTimeFunction sets a replacement function for calls to time().
// This API is intended to be used only for testing, thus may cause PDFium to behave poorly in production environments.
// Since callbacks can't be transferred between processes in gRPC, you can only use this in single-threaded mode.
func (p *PdfiumImplementation) FSDK_SetTimeFunction(request *requests.FSDK_SetTimeFunction) (*responses.FSDK_SetTimeFunction, error) {
	p.Lock()
	defer p.Unlock()

	if request.Function == nil {
		CurrentTimeHandler = nil
		_, err := p.Module.ExportedFunction("FSDK_SetTimeFunction").Call(p.Context, 0)
		if err != nil {
			return nil, err
		}
	} else {
		CurrentTimeHandler = request.Function
		_, err := p.Module.ExportedFunction("FSDK_SetTimeFunction_SET_CB").Call(p.Context)
		if err != nil {
			return nil, err
		}
	}

	return &responses.FSDK_SetTimeFunction{}, nil
}

var CurrentLocalTimeHandler requests.SetLocaltimeFunction

// FSDK_SetLocaltimeFunction sets a replacement function for calls to localtime().
// This API is intended to be used only for testing, thus may cause PDFium to behave poorly in production environments.
// Since callbacks can't be transferred between processes in gRPC, you can only use this in single-threaded mode.
func (p *PdfiumImplementation) FSDK_SetLocaltimeFunction(request *requests.FSDK_SetLocaltimeFunction) (*responses.FSDK_SetLocaltimeFunction, error) {
	p.Lock()
	defer p.Unlock()

	if request.Function == nil {
		CurrentLocalTimeHandler = nil
		_, err := p.Module.ExportedFunction("FSDK_SetLocaltimeFunction").Call(p.Context, 0)
		if err != nil {
			return nil, err
		}
	} else {
		CurrentLocalTimeHandler = request.Function
		_, err := p.Module.ExportedFunction("FSDK_SetLocaltimeFunction_SET_CB").Call(p.Context)
		if err != nil {
			return nil, err
		}
	}

	return &responses.FSDK_SetLocaltimeFunction{}, nil
}
