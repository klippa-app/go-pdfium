package implementation_webassembly

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// GetMetaData returns the metadata values of the document.
func (p *PdfiumImplementation) GetMetaData(request *requests.GetMetaData) (*responses.GetMetaData, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	getMetaText := func(tag string) (string, error) {
		cstr, err := p.CString(tag)
		if err != nil {
			return "", err
		}

		defer cstr.Free()

		// First get the metadata length.
		res, err := p.Module.ExportedFunction("FPDF_GetMetaText").Call(p.Context, *documentHandle.handle, cstr.Pointer, 0, 0)
		if err != nil {
			return "", err
		}

		metaSize := *(*int32)(unsafe.Pointer(&res[0]))
		if metaSize == 0 {
			return "", errors.New("Could not get metadata")
		}

		charDataPointer, err := p.ByteArrayPointer(uint64(metaSize), nil)
		defer charDataPointer.Free()

		_, err = p.Module.ExportedFunction("FPDF_GetMetaText").Call(p.Context, *documentHandle.handle, cstr.Pointer, charDataPointer.Pointer, uint64(metaSize))
		if err != nil {
			return "", err
		}

		charData, err := charDataPointer.Value(false)
		if err != nil {
			return "", err
		}

		transformedText, err := p.transformUTF16LEToUTF8(charData)
		if err != nil {
			return "", err
		}

		return transformedText, nil
	}

	if request.Tags == nil {
		request.Tags = &[]string{"Title", "Author", "Subject", "Keywords", "Creator", "Producer", "CreationDate", "ModDate"}
	}

	results := []responses.GetMetaDataTag{}

	tags := *request.Tags
	for i := range tags {
		result, err := getMetaText(tags[i])
		if err != nil {
			return nil, err
		}

		results = append(results, responses.GetMetaDataTag{
			Tag:   tags[i],
			Value: result,
		})
	}

	return &responses.GetMetaData{
		Tags: results,
	}, nil
}
