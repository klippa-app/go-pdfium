package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdf_attachment.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/google/uuid"
)

func (p *PdfiumImplementation) registerAttachment(attachment C.FPDF_ATTACHMENT, documentHandle *DocumentHandle) *AttachmentHandle {
	ref := uuid.New()
	handle := &AttachmentHandle{
		handle:      attachment,
		nativeRef:   references.FPDF_ATTACHMENT(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.attachmentRefs[handle.nativeRef] = handle
	p.attachmentRefs[handle.nativeRef] = handle

	return handle
}

// GetAttachments returns all the attachments of a document.
// Experimental API.
func (p *PdfiumImplementation) GetAttachments(request *requests.GetAttachments) (*responses.GetAttachments, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	cAttachmentCount := C.FPDFDoc_GetAttachmentCount(documentHandle.handle)
	attachmentCount := int(cAttachmentCount)
	if int(attachmentCount) == -1 {
		return nil, errors.New("could not get attachment count")
	}

	attachments := []responses.Attachment{}
	for i := 0; i < attachmentCount; i++ {
		attachment := C.FPDFDoc_GetAttachment(documentHandle.handle, C.int(i))
		if attachment == nil {
			continue
		}

		// First get the name value length.
		nameSize := C.FPDFAttachment_GetName(attachment, nil, 0)
		if nameSize == 0 {
			return nil, errors.New("Could not get name")
		}

		charData := make([]byte, nameSize)
		C.FPDFAttachment_GetName(attachment, (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

		transformedName, err := p.transformUTF16LEToUTF8(charData)
		if err != nil {
			return nil, err
		}

		newAttachment := responses.Attachment{
			Name:   transformedName,
			Values: []responses.AttachmentValue{},
		}

		bufLen := C.ulong(0)
		getFileSizeSuccess := C.FPDFAttachment_GetFile(attachment, C.NULL, 0, &bufLen)
		if int(getFileSizeSuccess) != 0 {
			fileData := make([]byte, bufLen)
			getFileSuccess := C.FPDFAttachment_GetFile(attachment, unsafe.Pointer(&fileData[0]), C.ulong(bufLen), &bufLen)
			if int(getFileSuccess) == 0 {
				return nil, errors.New("could not get file")
			}

			newAttachment.Content = fileData
		}

		requestKeys := []string{"Size", "CreationDate", "CheckSum"}
		for _, requestKey := range requestKeys {
			newValue := responses.AttachmentValue{
				Key: requestKey,
			}

			keyStr := C.CString(requestKey)
			defer C.free(unsafe.Pointer(keyStr))

			valueType := C.FPDFAttachment_GetValueType(attachment, keyStr)
			newValue.ValueType = enums.FPDF_OBJECT_TYPE(valueType)

			// Only strings supported for now.
			if newValue.ValueType == enums.FPDF_OBJECT_TYPE_STRING {
				// First get the string value length.
				stringValueSize := C.FPDFAttachment_GetStringValue(attachment, keyStr, nil, 0)
				if stringValueSize == 0 {
					return nil, errors.New("Could not get string value")
				}

				charData := make([]byte, stringValueSize)
				C.FPDFAttachment_GetStringValue(attachment, keyStr, (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

				transformedText, err := p.transformUTF16LEToUTF8(charData)
				if err != nil {
					return nil, err
				}

				newValue.StringValue = transformedText
			}

			newAttachment.Values = append(newAttachment.Values, newValue)
		}

		attachments = append(attachments, newAttachment)
	}

	return &responses.GetAttachments{
		Attachments: attachments,
	}, nil
}
