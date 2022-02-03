package implementation

/*
#cgo pkg-config: pdfium
#include "fpdf_signature.h"
*/
import "C"
import (
	"errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"unsafe"
)

// FPDF_GetSignatureCount returns the total number of signatures in the document.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetSignatureCount(request *requests.FPDF_GetSignatureCount) (*responses.FPDF_GetSignatureCount, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	count := C.FPDF_GetSignatureCount(documentHandle.handle)
	return &responses.FPDF_GetSignatureCount{
		Count: int(count),
	}, nil
}

// FPDF_GetSignatureObject returns the Nth signature of the document.
// Experimental API.
func (p *PdfiumImplementation) FPDF_GetSignatureObject(request *requests.FPDF_GetSignatureObject) (*responses.FPDF_GetSignatureObject, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	handle := C.FPDF_GetSignatureObject(documentHandle.handle, C.int(request.Index))
	if handle == nil {
		return nil, errors.New("could not load signature object")
	}

	signatureHandle := p.registerSignature(handle, documentHandle)

	return &responses.FPDF_GetSignatureObject{
		Index:     request.Index,
		Signature: signatureHandle.nativeRef,
	}, nil
}

// FPDFSignatureObj_GetContents returns the contents of a signature object.
// Experimental API.
func (p *PdfiumImplementation) FPDFSignatureObj_GetContents(request *requests.FPDFSignatureObj_GetContents) (*responses.FPDFSignatureObj_GetContents, error) {
	p.Lock()
	defer p.Unlock()

	signatureHandle, err := p.getSignatureHandle(request.Signature)
	if err != nil {
		return nil, err
	}

	// First get the signature length.
	signatureSize := C.FPDFSignatureObj_GetContents(signatureHandle.handle, C.NULL, 0)
	if signatureSize == 0 {
		return &responses.FPDFSignatureObj_GetContents{}, nil
	}

	signatureData := make([]byte, signatureSize)
	C.FPDFSignatureObj_GetContents(signatureHandle.handle, unsafe.Pointer(&signatureData[0]), C.ulong(len(signatureData)))

	return &responses.FPDFSignatureObj_GetContents{
		Contents: signatureData,
	}, nil
}

// FPDFSignatureObj_GetByteRange returns the byte range of a signature object.
// Experimental API.
func (p *PdfiumImplementation) FPDFSignatureObj_GetByteRange(request *requests.FPDFSignatureObj_GetByteRange) (*responses.FPDFSignatureObj_GetByteRange, error) {
	p.Lock()
	defer p.Unlock()

	signatureHandle, err := p.getSignatureHandle(request.Signature)
	if err != nil {
		return nil, err
	}

	var nullBuffer *C.int

	// First get the signature length.
	byteRangeSize := C.FPDFSignatureObj_GetByteRange(signatureHandle.handle, nullBuffer, 0)
	if byteRangeSize == 0 {
		return &responses.FPDFSignatureObj_GetByteRange{}, nil
	}

	byteRangeData := make([]C.int, byteRangeSize)
	C.FPDFSignatureObj_GetByteRange(signatureHandle.handle, (*C.int)(unsafe.Pointer(&byteRangeData[0])), C.ulong(len(byteRangeData)))

	byteRangeValues := make([]int, byteRangeSize)
	for i := range byteRangeData {
		byteRangeValues[i] = int(byteRangeData[i])
	}

	return &responses.FPDFSignatureObj_GetByteRange{
		ByteRange: byteRangeValues,
	}, nil
}

// FPDFSignatureObj_GetSubFilter returns the encoding of the value of a signature object.
// Experimental API.
func (p *PdfiumImplementation) FPDFSignatureObj_GetSubFilter(request *requests.FPDFSignatureObj_GetSubFilter) (*responses.FPDFSignatureObj_GetSubFilter, error) {
	p.Lock()
	defer p.Unlock()

	signatureHandle, err := p.getSignatureHandle(request.Signature)
	if err != nil {
		return nil, err
	}

	var nullBuffer *C.char

	// First get the signature length.
	subFilterLength := C.FPDFSignatureObj_GetSubFilter(signatureHandle.handle, nullBuffer, 0)
	if subFilterLength == 0 {
		return &responses.FPDFSignatureObj_GetSubFilter{}, nil
	}

	subFilterData := make([]byte, subFilterLength)
	C.FPDFSignatureObj_GetSubFilter(signatureHandle.handle, (*C.char)(unsafe.Pointer(&subFilterData[0])), C.ulong(len(subFilterData)))

	subFilterString := string(subFilterData[:subFilterLength-1]) // Remove NULL terminator.
	return &responses.FPDFSignatureObj_GetSubFilter{
		SubFilter: &subFilterString,
	}, nil
}

// FPDFSignatureObj_GetReason returns the reason (comment) of the signature object.
// Experimental API.
func (p *PdfiumImplementation) FPDFSignatureObj_GetReason(request *requests.FPDFSignatureObj_GetReason) (*responses.FPDFSignatureObj_GetReason, error) {
	p.Lock()
	defer p.Unlock()

	signatureHandle, err := p.getSignatureHandle(request.Signature)
	if err != nil {
		return nil, err
	}

	// First get the reason length.
	reasonLength := C.FPDFSignatureObj_GetReason(signatureHandle.handle, C.NULL, 0)
	if reasonLength == 0 {
		return &responses.FPDFSignatureObj_GetReason{}, nil
	}

	reasonData := make([]byte, reasonLength)
	C.FPDFSignatureObj_GetReason(signatureHandle.handle, unsafe.Pointer(&reasonData[0]), C.ulong(len(reasonData)))

	transformedText, err := p.transformUTF16LEToUTF8(reasonData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFSignatureObj_GetReason{
		Reason: &transformedText,
	}, nil
}

// FPDFSignatureObj_GetTime returns the time of signing of a signature object.
// Experimental API.
func (p *PdfiumImplementation) FPDFSignatureObj_GetTime(request *requests.FPDFSignatureObj_GetTime) (*responses.FPDFSignatureObj_GetTime, error) {
	p.Lock()
	defer p.Unlock()

	signatureHandle, err := p.getSignatureHandle(request.Signature)
	if err != nil {
		return nil, err
	}

	var nullBuffer *C.char

	// First get the time length.
	timeLength := C.FPDFSignatureObj_GetTime(signatureHandle.handle, nullBuffer, 0)
	if timeLength == 0 {
		return &responses.FPDFSignatureObj_GetTime{}, nil
	}

	timeData := make([]byte, timeLength)
	C.FPDFSignatureObj_GetTime(signatureHandle.handle, (*C.char)(unsafe.Pointer(&timeData[0])), C.ulong(len(timeData)))

	timeString := string(timeData[:timeLength-1]) // Remove NULL terminator.

	return &responses.FPDFSignatureObj_GetTime{
		Time: &timeString,
	}, nil
}

// FPDFSignatureObj_GetDocMDPPermission returns the DocMDP permission of a signature object.
// Experimental API.
func (p *PdfiumImplementation) FPDFSignatureObj_GetDocMDPPermission(request *requests.FPDFSignatureObj_GetDocMDPPermission) (*responses.FPDFSignatureObj_GetDocMDPPermission, error) {
	p.Lock()
	defer p.Unlock()

	signatureHandle, err := p.getSignatureHandle(request.Signature)
	if err != nil {
		return nil, err
	}

	permission := C.FPDFSignatureObj_GetDocMDPPermission(signatureHandle.handle)
	if permission == 0 {
		return nil, errors.New("could not get DocMDPPermission")
	}

	return &responses.FPDFSignatureObj_GetDocMDPPermission{
		DocMDPPermission: int(permission),
	}, nil
}
