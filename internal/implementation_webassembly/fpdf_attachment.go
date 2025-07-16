package implementation_webassembly

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

	res, err := p.Module.ExportedFunction("FPDFDoc_GetAttachmentCount").Call(p.Context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

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

	namePointer, err := p.CFPDF_WIDESTRING(request.Name)
	if err != nil {
		return nil, err
	}
	defer namePointer.Free()

	res, err := p.Module.ExportedFunction("FPDFDoc_AddAttachment").Call(p.Context, *documentHandle.handle, namePointer.Pointer)
	if err != nil {
		return nil, err
	}

	handle := res[0]
	if handle == 0 {
		return nil, errors.New("could not create attachment object")
	}

	attachmentHandle := p.registerAttachment(&handle, documentHandle)

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

	res, err := p.Module.ExportedFunction("FPDFDoc_GetAttachment").Call(p.Context, *documentHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	handle := res[0]
	if handle == 0 {
		return nil, errors.New("could not get attachment object")
	}

	attachmentHandle := p.registerAttachment(&handle, documentHandle)

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

	res, err := p.Module.ExportedFunction("FPDFDoc_DeleteAttachment").Call(p.Context, *documentHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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
	res, err := p.Module.ExportedFunction("FPDFAttachment_GetName").Call(p.Context, *attachmentHandle.handle, 0, 0)
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

	res, err = p.Module.ExportedFunction("FPDFAttachment_GetName").Call(p.Context, *attachmentHandle.handle, charDataPointer.Pointer, nameSize)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAttachment_HasKey").Call(p.Context, *attachmentHandle.handle, keyPointer.Pointer)
	if err != nil {
		return nil, err
	}

	hasKey := *(*int32)(unsafe.Pointer(&res[0]))

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

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAttachment_GetValueType").Call(p.Context, *attachmentHandle.handle, keyPointer.Pointer)
	if err != nil {
		return nil, err
	}

	valueType := *(*int32)(unsafe.Pointer(&res[0]))

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

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	valuePointer, err := p.CFPDF_WIDESTRING(request.Value)
	if err != nil {
		return nil, err
	}
	defer valuePointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAttachment_SetStringValue").Call(p.Context, *attachmentHandle.handle, keyPointer.Pointer, valuePointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	// First get the string value length.
	res, err := p.Module.ExportedFunction("FPDFAttachment_GetStringValue").Call(p.Context, *attachmentHandle.handle, keyPointer.Pointer, 0, 0)
	if err != nil {
		return nil, err
	}

	stringValueSize := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if stringValueSize == 0 {
		return nil, errors.New("Could not get string value")
	}

	charDataPointer, err := p.ByteArrayPointer(stringValueSize, nil)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDFAttachment_GetStringValue").Call(p.Context, *attachmentHandle.handle, keyPointer.Pointer, charDataPointer.Pointer, stringValueSize)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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

	fileSize := uint64(len(request.Contents))
	fileDataPointer, err := p.ByteArrayPointer(fileSize, request.Contents)
	if err != nil {
		return nil, err
	}
	defer fileDataPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAttachment_SetFile").Call(p.Context, *attachmentHandle.handle, *documentHandle.handle, fileDataPointer.Pointer, fileSize)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
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

	bufLenPointer, err := p.ULongPointer()
	if err != nil {
		return nil, err
	}
	defer bufLenPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFAttachment_GetFile").Call(p.Context, *attachmentHandle.handle, 0, 0, bufLenPointer.Pointer)
	if err != nil {
		return nil, err
	}

	getFileSizeSuccess := *(*int32)(unsafe.Pointer(&res[0]))
	if int(getFileSizeSuccess) == 0 {
		return &responses.FPDFAttachment_GetFile{
			Contents: nil,
		}, nil
	}

	bufLen, err := bufLenPointer.Value()
	if err != nil {
		return nil, err
	}

	fileDataPointer, err := p.ByteArrayPointer(bufLen, nil)
	if err != nil {
		return nil, err
	}
	defer fileDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDFAttachment_GetFile").Call(p.Context, *attachmentHandle.handle, fileDataPointer.Pointer, bufLen, bufLenPointer.Pointer)
	if err != nil {
		return nil, err
	}

	getFileSuccess := *(*int32)(unsafe.Pointer(&res[0]))
	if int(getFileSuccess) == 0 {
		return nil, errors.New("could not get file")
	}

	fileData, err := fileDataPointer.Value(true)
	if err != nil {
		return nil, err
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
	res, err := p.Module.ExportedFunction("FPDFAttachment_GetSubtype").Call(p.Context, *attachmentHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	stringValueSize := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if stringValueSize == 0 {
		return nil, errors.New("Could not get string value")
	}

	charDataPointer, err := p.ByteArrayPointer(stringValueSize, nil)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDFAttachment_GetSubtype").Call(p.Context, *attachmentHandle.handle, charDataPointer.Pointer, stringValueSize)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

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
