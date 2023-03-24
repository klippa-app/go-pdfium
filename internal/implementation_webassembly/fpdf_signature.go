package implementation_webassembly

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
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

	res, err := p.Module.ExportedFunction("FPDF_GetSignatureCount").Call(p.Context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))
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

	res, err := p.Module.ExportedFunction("FPDF_GetSignatureObject").Call(p.Context, *documentHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	handle := res[0]
	if handle == 0 {
		return nil, errors.New("could not load signature object")
	}

	signatureHandle := p.registerSignature(&handle, documentHandle)

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
	res, err := p.Module.ExportedFunction("FPDFSignatureObj_GetContents").Call(p.Context, *signatureHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	signatureSize := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if signatureSize == 0 {
		return &responses.FPDFSignatureObj_GetContents{}, nil
	}

	signatureDataPointer, err := p.ByteArrayPointer(signatureSize, nil)
	if err != nil {
		return nil, err
	}
	defer signatureDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDFSignatureObj_GetContents").Call(p.Context, *signatureHandle.handle, signatureDataPointer.Pointer, signatureSize)
	if err != nil {
		return nil, err
	}

	signatureData, err := signatureDataPointer.Value(true)
	if err != nil {
		return nil, err
	}

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

	// First get the signature length.
	res, err := p.Module.ExportedFunction("FPDFSignatureObj_GetByteRange").Call(p.Context, *signatureHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	byteRangeSize := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if byteRangeSize == 0 {
		return &responses.FPDFSignatureObj_GetByteRange{}, nil
	}

	intArrayPointer, err := p.IntArrayPointer(byteRangeSize)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDFSignatureObj_GetByteRange").Call(p.Context, *signatureHandle.handle, intArrayPointer.Pointer, byteRangeSize)
	if err != nil {
		return nil, err
	}

	byteRangeValues, err := intArrayPointer.Value()
	if err != nil {
		return nil, err
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

	// First get the signature length.
	res, err := p.Module.ExportedFunction("FPDFSignatureObj_GetSubFilter").Call(p.Context, *signatureHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	subFilterLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if subFilterLength == 0 {
		return &responses.FPDFSignatureObj_GetSubFilter{}, nil
	}

	subFilterDataPointer, err := p.ByteArrayPointer(subFilterLength, nil)
	if err != nil {
		return nil, err
	}
	defer subFilterDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDFSignatureObj_GetSubFilter").Call(p.Context, *signatureHandle.handle, subFilterDataPointer.Pointer, subFilterLength)
	if err != nil {
		return nil, err
	}

	subFilterData, err := subFilterDataPointer.Value(true)
	if err != nil {
		return nil, err
	}

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
	res, err := p.Module.ExportedFunction("FPDFSignatureObj_GetReason").Call(p.Context, *signatureHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	reasonLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if reasonLength == 0 {
		return &responses.FPDFSignatureObj_GetReason{}, nil
	}

	reasonDataPointer, err := p.ByteArrayPointer(reasonLength, nil)
	if err != nil {
		return nil, err
	}
	defer reasonDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDFSignatureObj_GetReason").Call(p.Context, *signatureHandle.handle, reasonDataPointer.Pointer, reasonLength)
	if err != nil {
		return nil, err
	}

	reasonDataa, err := reasonDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF16LEToUTF8(reasonDataa)
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

	res, err := p.Module.ExportedFunction("FPDFSignatureObj_GetTime").Call(p.Context, *signatureHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	timeLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if timeLength == 0 {
		return &responses.FPDFSignatureObj_GetTime{}, nil
	}

	timeDataPointer, err := p.ByteArrayPointer(timeLength, nil)
	if err != nil {
		return nil, err
	}
	defer timeDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDFSignatureObj_GetTime").Call(p.Context, *signatureHandle.handle, timeDataPointer.Pointer, timeLength)
	if err != nil {
		return nil, err
	}

	timeData, err := timeDataPointer.Value(true)
	if err != nil {
		return nil, err
	}

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

	res, err := p.Module.ExportedFunction("FPDFSignatureObj_GetDocMDPPermission").Call(p.Context, *signatureHandle.handle)
	if err != nil {
		return nil, err
	}

	permission := *(*int32)(unsafe.Pointer(&res[0]))
	if permission == 0 {
		return nil, errors.New("could not get DocMDPPermission")
	}

	return &responses.FPDFSignatureObj_GetDocMDPPermission{
		DocMDPPermission: int(permission),
	}, nil
}
