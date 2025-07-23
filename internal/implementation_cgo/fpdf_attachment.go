//go:build pdfium_experimental
// +build pdfium_experimental

package implementation_cgo

/*
#cgo pkg-config: pdfium
#include "fpdf_attachment.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFDoc_GetAttachmentCount returns the number of embedded files in the given document.
// Experimental API.
func (p *PdfiumImplementation) FPDFDoc_GetAttachmentCount(request *requests.FPDFDoc_GetAttachmentCount) (*responses.FPDFDoc_GetAttachmentCount, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	count := C.FPDFDoc_GetAttachmentCount(documentHandle.handle)

	return &responses.FPDFDoc_GetAttachmentCount{
		AttachmentCount: int(count),
	}, nil
}

// FPDFDoc_AddAttachment adds an embedded file with the given name in the given document. If the name is empty, or if
// the name is the name of an existing embedded file in the document, or if
// the document's embedded file name tree is too deep (i.e. the document has too
// many embedded files already), then a new attachment will not be added.
// Experimental API.
func (p *PdfiumImplementation) FPDFDoc_AddAttachment(request *requests.FPDFDoc_AddAttachment) (*responses.FPDFDoc_AddAttachment, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF8ToUTF16LE(request.Name)
	if err != nil {
		return nil, err
	}

	handle := C.FPDFDoc_AddAttachment(documentHandle.handle, (C.FPDF_WIDESTRING)(unsafe.Pointer(&transformedText[0])))
	if handle == nil {
		return nil, errors.New("could not create attachment object")
	}

	attachmentHandle := p.registerAttachment(handle, documentHandle)

	return &responses.FPDFDoc_AddAttachment{
		Attachment: attachmentHandle.nativeRef,
	}, nil
}

// FPDFDoc_GetAttachment returns the embedded attachment at the given index in the given document. Note that the returned
// attachment handle is only valid while the document is open.
// Experimental API.
func (p *PdfiumImplementation) FPDFDoc_GetAttachment(request *requests.FPDFDoc_GetAttachment) (*responses.FPDFDoc_GetAttachment, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	handle := C.FPDFDoc_GetAttachment(documentHandle.handle, C.int(request.Index))
	if handle == nil {
		return nil, errors.New("could not get attachment object")
	}

	attachmentHandle := p.registerAttachment(handle, documentHandle)

	return &responses.FPDFDoc_GetAttachment{
		Index:      request.Index,
		Attachment: attachmentHandle.nativeRef,
	}, nil
}

// FPDFDoc_DeleteAttachment deletes the embedded attachment at the given index in the given document. Note that this does
// not remove the attachment data from the PDF file; it simply removes the
// file's entry in the embedded files name tree so that it does not appear in
// the attachment list. This behavior may change in the future.
// Experimental API.
func (p *PdfiumImplementation) FPDFDoc_DeleteAttachment(request *requests.FPDFDoc_DeleteAttachment) (*responses.FPDFDoc_DeleteAttachment, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	success := C.FPDFDoc_DeleteAttachment(documentHandle.handle, C.int(request.Index))
	if int(success) == 0 {
		return nil, errors.New("could not delete attachment object")
	}

	return &responses.FPDFDoc_DeleteAttachment{
		Index: request.Index,
	}, nil
}

// FPDFAttachment_GetName returns the name of the attachment file.
// Experimental API.
func (p *PdfiumImplementation) FPDFAttachment_GetName(request *requests.FPDFAttachment_GetName) (*responses.FPDFAttachment_GetName, error) {
	p.Lock()
	defer p.Unlock()

	attachmentHandle, err := p.getAttachmentHandle(request.Attachment)
	if err != nil {
		return nil, err
	}

	// First get the title length.
	nameSize := C.FPDFAttachment_GetName(attachmentHandle.handle, nil, 0)
	if nameSize == 0 {
		return nil, errors.New("Could not get name")
	}

	charData := make([]byte, nameSize)
	C.FPDFAttachment_GetName(attachmentHandle.handle, (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAttachment_GetName{
		Name: transformedText,
	}, nil
}

// FPDFAttachment_HasKey check if the params dictionary of the given attachment has the given key as a key.
// Experimental API.
func (p *PdfiumImplementation) FPDFAttachment_HasKey(request *requests.FPDFAttachment_HasKey) (*responses.FPDFAttachment_HasKey, error) {
	p.Lock()
	defer p.Unlock()

	attachmentHandle, err := p.getAttachmentHandle(request.Attachment)
	if err != nil {
		return nil, err
	}

	keyStr := C.CString(request.Key)
	defer C.free(unsafe.Pointer(keyStr))

	hasKey := C.FPDFAttachment_HasKey(attachmentHandle.handle, keyStr)

	return &responses.FPDFAttachment_HasKey{
		Key:    request.Key,
		HasKey: int(hasKey) == 1,
	}, nil
}

// FPDFAttachment_GetValueType returns the type of the value corresponding to the given key in the params dictionary of
// the embedded attachment.
// Experimental API.
func (p *PdfiumImplementation) FPDFAttachment_GetValueType(request *requests.FPDFAttachment_GetValueType) (*responses.FPDFAttachment_GetValueType, error) {
	p.Lock()
	defer p.Unlock()

	attachmentHandle, err := p.getAttachmentHandle(request.Attachment)
	if err != nil {
		return nil, err
	}

	keyStr := C.CString(request.Key)
	defer C.free(unsafe.Pointer(keyStr))

	valueType := C.FPDFAttachment_GetValueType(attachmentHandle.handle, keyStr)
	return &responses.FPDFAttachment_GetValueType{
		Key:       request.Key,
		ValueType: enums.FPDF_OBJECT_TYPE(valueType),
	}, nil
}

// FPDFAttachment_SetStringValue sets the string value corresponding to the given key in the params dictionary of the
// embedded file attachment, overwriting the existing value if any.
// Experimental API.
func (p *PdfiumImplementation) FPDFAttachment_SetStringValue(request *requests.FPDFAttachment_SetStringValue) (*responses.FPDFAttachment_SetStringValue, error) {
	p.Lock()
	defer p.Unlock()

	attachmentHandle, err := p.getAttachmentHandle(request.Attachment)
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF8ToUTF16LE(request.Value)
	if err != nil {
		return nil, err
	}

	keyStr := C.CString(request.Key)
	defer C.free(unsafe.Pointer(keyStr))

	success := C.FPDFAttachment_SetStringValue(attachmentHandle.handle, keyStr, (C.FPDF_WIDESTRING)(unsafe.Pointer(&transformedText[0])))
	if int(success) == 0 {
		return nil, errors.New("could not set attachment string value")
	}

	return &responses.FPDFAttachment_SetStringValue{
		Key:   request.Key,
		Value: request.Value,
	}, nil
}

// FPDFAttachment_GetStringValue gets the string value corresponding to the given key in the params dictionary of the
// embedded file attachment.
// Experimental API.
func (p *PdfiumImplementation) FPDFAttachment_GetStringValue(request *requests.FPDFAttachment_GetStringValue) (*responses.FPDFAttachment_GetStringValue, error) {
	p.Lock()
	defer p.Unlock()

	attachmentHandle, err := p.getAttachmentHandle(request.Attachment)
	if err != nil {
		return nil, err
	}

	keyStr := C.CString(request.Key)
	defer C.free(unsafe.Pointer(keyStr))

	// First get the string value length.
	stringValueSize := C.FPDFAttachment_GetStringValue(attachmentHandle.handle, keyStr, nil, 0)
	if stringValueSize == 0 {
		return nil, errors.New("Could not get string value")
	}

	charData := make([]byte, stringValueSize)
	C.FPDFAttachment_GetStringValue(attachmentHandle.handle, keyStr, (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFAttachment_GetStringValue{
		Key:   request.Key,
		Value: transformedText,
	}, nil
}

// FPDFAttachment_SetFile set the file data of the given attachment, overwriting the existing file data if any.
// The creation date and checksum will be updated, while all other dictionary
// entries will be deleted. Note that only contents with a length smaller than
// INT_MAX is supported.
// Experimental API.
func (p *PdfiumImplementation) FPDFAttachment_SetFile(request *requests.FPDFAttachment_SetFile) (*responses.FPDFAttachment_SetFile, error) {
	p.Lock()
	defer p.Unlock()

	attachmentHandle, err := p.getAttachmentHandle(request.Attachment)
	if err != nil {
		return nil, err
	}

	documentHandle, err := p.getDocumentHandle(attachmentHandle.documentRef)
	if err != nil {
		return nil, err
	}

	success := C.FPDFAttachment_SetFile(attachmentHandle.handle, documentHandle.handle, unsafe.Pointer(&request.Contents[0]), C.ulong(len(request.Contents)))
	if int(success) == 0 {
		return nil, errors.New("Could not get set file")
	}

	return &responses.FPDFAttachment_SetFile{}, nil
}

// FPDFAttachment_GetFile gets the file data of the given attachment.
// Experimental API.
func (p *PdfiumImplementation) FPDFAttachment_GetFile(request *requests.FPDFAttachment_GetFile) (*responses.FPDFAttachment_GetFile, error) {
	p.Lock()
	defer p.Unlock()

	attachmentHandle, err := p.getAttachmentHandle(request.Attachment)
	if err != nil {
		return nil, err
	}

	bufLen := C.ulong(0)
	getFileSizeSuccess := C.FPDFAttachment_GetFile(attachmentHandle.handle, C.NULL, 0, &bufLen)
	if int(getFileSizeSuccess) == 0 {
		return &responses.FPDFAttachment_GetFile{
			Contents: nil,
		}, nil
	}

	fileData := make([]byte, bufLen)
	getFileSuccess := C.FPDFAttachment_GetFile(attachmentHandle.handle, unsafe.Pointer(&fileData[0]), C.ulong(bufLen), &bufLen)
	if int(getFileSuccess) == 0 {
		return nil, errors.New("could not get file")
	}

	return &responses.FPDFAttachment_GetFile{
		Contents: fileData,
	}, nil
}

// FPDFAttachment_GetSubtype gets the MIME type (Subtype) of the embedded file attachment.
// Experimental API.
func (p *PdfiumImplementation) FPDFAttachment_GetSubtype(request *requests.FPDFAttachment_GetSubtype) (*responses.FPDFAttachment_GetSubtype, error) {
	p.Lock()
	defer p.Unlock()

	attachmentHandle, err := p.getAttachmentHandle(request.Attachment)
	if err != nil {
		return nil, err
	}

	// First get the string value length.
	stringValueSize := C.FPDFAttachment_GetSubtype(attachmentHandle.handle, nil, 0)
	if stringValueSize == 0 {
		return nil, errors.New("Could not get string value")
	}

	charData := make([]byte, stringValueSize)
	C.FPDFAttachment_GetSubtype(attachmentHandle.handle, (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	if transformedText == "" {
		return &responses.FPDFAttachment_GetSubtype{
			Subtype: nil,
		}, nil
	}

	return &responses.FPDFAttachment_GetSubtype{
		Subtype: &transformedText,
	}, nil
}
