//go:build !pdfium_experimental
// +build !pdfium_experimental

package implementation_cgo

import (
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFDoc_GetAttachmentCount returns the number of embedded files in the given document.
// Experimental API.
func (p *PdfiumImplementation) FPDFDoc_GetAttachmentCount(request *requests.FPDFDoc_GetAttachmentCount) (*responses.FPDFDoc_GetAttachmentCount, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFDoc_AddAttachment adds an embedded file with the given name in the given document. If the name is empty, or if
// the name is the name of an existing embedded file in the document, or if
// the document's embedded file name tree is too deep (i.e. the document has too
// many embedded files already), then a new attachment will not be added.
// Experimental API.
func (p *PdfiumImplementation) FPDFDoc_AddAttachment(request *requests.FPDFDoc_AddAttachment) (*responses.FPDFDoc_AddAttachment, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFDoc_GetAttachment returns the embedded attachment at the given index in the given document. Note that the returned
// attachment handle is only valid while the document is open.
// Experimental API.
func (p *PdfiumImplementation) FPDFDoc_GetAttachment(request *requests.FPDFDoc_GetAttachment) (*responses.FPDFDoc_GetAttachment, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFDoc_DeleteAttachment deletes the embedded attachment at the given index in the given document. Note that this does
// not remove the attachment data from the PDF file; it simply removes the
// file's entry in the embedded files name tree so that it does not appear in
// the attachment list. This behavior may change in the future.
// Experimental API.
func (p *PdfiumImplementation) FPDFDoc_DeleteAttachment(request *requests.FPDFDoc_DeleteAttachment) (*responses.FPDFDoc_DeleteAttachment, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAttachment_GetName returns the name of the attachment file.
// Experimental API.
func (p *PdfiumImplementation) FPDFAttachment_GetName(request *requests.FPDFAttachment_GetName) (*responses.FPDFAttachment_GetName, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAttachment_HasKey check if the params dictionary of the given attachment has the given key as a key.
// Experimental API.
func (p *PdfiumImplementation) FPDFAttachment_HasKey(request *requests.FPDFAttachment_HasKey) (*responses.FPDFAttachment_HasKey, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAttachment_GetValueType returns the type of the value corresponding to the given key in the params dictionary of
// the embedded attachment.
// Experimental API.
func (p *PdfiumImplementation) FPDFAttachment_GetValueType(request *requests.FPDFAttachment_GetValueType) (*responses.FPDFAttachment_GetValueType, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAttachment_SetStringValue sets the string value corresponding to the given key in the params dictionary of the
// embedded file attachment, overwriting the existing value if any.
// Experimental API.
func (p *PdfiumImplementation) FPDFAttachment_SetStringValue(request *requests.FPDFAttachment_SetStringValue) (*responses.FPDFAttachment_SetStringValue, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAttachment_GetStringValue gets the string value corresponding to the given key in the params dictionary of the
// embedded file attachment.
// Experimental API.
func (p *PdfiumImplementation) FPDFAttachment_GetStringValue(request *requests.FPDFAttachment_GetStringValue) (*responses.FPDFAttachment_GetStringValue, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAttachment_SetFile set the file data of the given attachment, overwriting the existing file data if any.
// The creation date and checksum will be updated, while all other dictionary
// entries will be deleted. Note that only contents with a length smaller than
// INT_MAX is supported.
// Experimental API.
func (p *PdfiumImplementation) FPDFAttachment_SetFile(request *requests.FPDFAttachment_SetFile) (*responses.FPDFAttachment_SetFile, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}

// FPDFAttachment_GetFile gets the file data of the given attachment.
// Experimental API.
func (p *PdfiumImplementation) FPDFAttachment_GetFile(request *requests.FPDFAttachment_GetFile) (*responses.FPDFAttachment_GetFile, error) {
	return nil, pdfium_errors.ErrExperimentalUnsupported
}
