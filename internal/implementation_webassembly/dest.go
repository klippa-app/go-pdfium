package implementation_webassembly

import (
	"unsafe"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerDest(dest *uint64, documentHandle *DocumentHandle) *DestHandle {
	ref := uuid.New()
	handle := &DestHandle{
		handle:      dest,
		nativeRef:   references.FPDF_DEST(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.destRefs[handle.nativeRef] = handle
	p.destRefs[handle.nativeRef] = handle

	return handle
}

func (p *PdfiumImplementation) getDestInfo(destHandle *DestHandle, documentHandle *DocumentHandle) (*responses.DestInfo, error) {
	res, err := p.module.ExportedFunction("FPDFDest_GetDestPageIndex").Call(p.context, *documentHandle.handle, *destHandle.handle)
	if err != nil {
		return nil, err
	}

	pageIndex := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.DestInfo{
		Reference: destHandle.nativeRef,
		PageIndex: int(pageIndex),
	}, nil
}

func (p *PdfiumImplementation) GetDestInfo(request *requests.GetDestInfo) (*responses.GetDestInfo, error) {
	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	destHandle, err := p.getDestHandle(request.Dest)
	if err != nil {
		return nil, err
	}

	destInfo, err := p.getDestInfo(destHandle, documentHandle)
	if err != nil {
		return nil, err
	}

	return &responses.GetDestInfo{
		DestInfo: *destInfo,
	}, nil
}
