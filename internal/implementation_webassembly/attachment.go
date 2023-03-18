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

func (p *PdfiumImplementation) registerAttachment(attachment *uint64, documentHandle *DocumentHandle) *AttachmentHandle {
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

	res, err := p.Module.ExportedFunction("FPDFDoc_GetAttachmentCount").Call(p.Context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	cAttachmentCount := *(*int32)(unsafe.Pointer(&res[0]))
	attachmentCount := int(cAttachmentCount)
	if int(attachmentCount) == -1 {
		return nil, errors.New("could not get attachment count")
	}

	attachments := []responses.Attachment{}
	for i := 0; i < attachmentCount; i++ {
		res, err := p.Module.ExportedFunction("FPDFDoc_GetAttachment").Call(p.Context, *documentHandle.handle, uint64(i))
		if err != nil {
			return nil, err
		}

		attachment := res[0]
		if attachment == 0 {
			continue
		}

		// First get the name value length.
		res, err = p.Module.ExportedFunction("FPDFAttachment_GetName").Call(p.Context, attachment, 0, 0)
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

		defer charDataPointer.Free()

		_, err = p.Module.ExportedFunction("FPDFAttachment_GetName").Call(p.Context, attachment, charDataPointer.Pointer, nameSize)
		if err != nil {
			return nil, err
		}

		charData, err := charDataPointer.Value(false)
		if err != nil {
			return nil, err
		}

		transformedName, err := p.transformUTF16LEToUTF8(charData)
		if err != nil {
			return nil, err
		}

		newAttachment := responses.Attachment{
			Name:   transformedName,
			Values: []responses.AttachmentValue{},
		}

		bufLenPointer, err := p.ULongPointer()
		if err != nil {
			return nil, err
		}
		defer bufLenPointer.Free()

		res, err = p.Module.ExportedFunction("FPDFAttachment_GetFile").Call(p.Context, attachment, 0, 0, bufLenPointer.Pointer)
		if err != nil {
			return nil, err
		}

		getFileSizeSuccess := *(*int32)(unsafe.Pointer(&res[0]))
		if int(getFileSizeSuccess) != 0 {
			bufLen, err := bufLenPointer.Value()
			if err != nil {
				return nil, err
			}

			fileDataPointer, err := p.ByteArrayPointer(bufLen, nil)
			defer fileDataPointer.Free()

			res, err = p.Module.ExportedFunction("FPDFAttachment_GetFile").Call(p.Context, attachment, fileDataPointer.Pointer, bufLen, bufLenPointer.Pointer)
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

			newAttachment.Content = fileData
		}

		requestKeys := []string{"Size", "CreationDate", "CheckSum"}
		for _, requestKey := range requestKeys {
			newValue := responses.AttachmentValue{
				Key: requestKey,
			}

			keyStr, err := p.CString(requestKey)
			if err != nil {
				return nil, err
			}

			defer keyStr.Free()

			res, err = p.Module.ExportedFunction("FPDFAttachment_GetValueType").Call(p.Context, attachment, keyStr.Pointer)
			if err != nil {
				return nil, err
			}

			valueType := *(*int32)(unsafe.Pointer(&res[0]))
			newValue.ValueType = enums.FPDF_OBJECT_TYPE(valueType)

			// Only strings supported for now.
			if newValue.ValueType == enums.FPDF_OBJECT_TYPE_STRING {
				// First get the string value length.
				res, err = p.Module.ExportedFunction("FPDFAttachment_GetStringValue").Call(p.Context, attachment, keyStr.Pointer, 0, 0)
				if err != nil {
					return nil, err
				}

				stringValueSize := *(*int32)(unsafe.Pointer(&res[0]))

				if stringValueSize == 0 {
					return nil, errors.New("Could not get string value")
				}

				charDataPointer, err = p.ByteArrayPointer(uint64(stringValueSize), nil)
				defer charDataPointer.Free()

				res, err = p.Module.ExportedFunction("FPDFAttachment_GetStringValue").Call(p.Context, attachment, keyStr.Pointer, charDataPointer.Pointer, uint64(stringValueSize))
				if err != nil {
					return nil, err
				}

				charData, err = charDataPointer.Value(false)
				if err != nil {
					return nil, err
				}

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
