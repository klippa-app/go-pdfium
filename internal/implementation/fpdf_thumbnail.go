package implementation

/*
#cgo pkg-config: pdfium
#include "fpdf_thumbnail.h"
*/
import "C"
import (
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"unsafe"
)

// FPDFPage_GetDecodedThumbnailData returns the decoded data from the thumbnail of the given page if it exists.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_GetDecodedThumbnailData(request *requests.FPDFPage_GetDecodedThumbnailData) (*responses.FPDFPage_GetDecodedThumbnailData, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	// First get the thumbnail length.
	thumbnailLength := C.FPDFPage_GetDecodedThumbnailData(pageHandle.handle, C.NULL, 0)
	if thumbnailLength == 0 {
		return &responses.FPDFPage_GetDecodedThumbnailData{}, nil
	}

	thumbnailData := make([]byte, thumbnailLength)
	C.FPDFPage_GetDecodedThumbnailData(pageHandle.handle, unsafe.Pointer(&thumbnailData[0]), C.ulong(len(thumbnailData)))

	return &responses.FPDFPage_GetDecodedThumbnailData{
		Thumbnail: thumbnailData,
	}, nil
}

// FPDFPage_GetRawThumbnailData returns the raw data from the thumbnail of the given page if it exists.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_GetRawThumbnailData(request *requests.FPDFPage_GetRawThumbnailData) (*responses.FPDFPage_GetRawThumbnailData, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	// First get the thumbnail length.
	rawThumbnailLength := C.FPDFPage_GetRawThumbnailData(pageHandle.handle, C.NULL, 0)
	if rawThumbnailLength == 0 {
		return &responses.FPDFPage_GetRawThumbnailData{}, nil
	}

	rawThumbnailData := make([]byte, rawThumbnailLength)
	C.FPDFPage_GetRawThumbnailData(pageHandle.handle, unsafe.Pointer(&rawThumbnailData[0]), C.ulong(len(rawThumbnailData)))

	return &responses.FPDFPage_GetRawThumbnailData{
		RawThumbnail: rawThumbnailData,
	}, nil
}

// FPDFPage_GetThumbnailAsBitmap returns the thumbnail of the given page as a FPDF_BITMAP.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_GetThumbnailAsBitmap(request *requests.FPDFPage_GetThumbnailAsBitmap) (*responses.FPDFPage_GetThumbnailAsBitmap, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	handle := C.FPDFPage_GetThumbnailAsBitmap(pageHandle.handle)
	if handle == nil {
		return &responses.FPDFPage_GetThumbnailAsBitmap{}, nil
	}

	bitmapHandle := p.registerBitmap(handle)

	return &responses.FPDFPage_GetThumbnailAsBitmap{
		Bitmap: &bitmapHandle.nativeRef,
	}, nil
}
