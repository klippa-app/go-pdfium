package implementation_webassembly

import (
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
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
	res, err := p.Module.ExportedFunction("FPDFPage_GetDecodedThumbnailData").Call(p.Context, *pageHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	thumbnailLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if thumbnailLength == 0 {
		return &responses.FPDFPage_GetDecodedThumbnailData{}, nil
	}

	thumbnailDataPointer, err := p.ByteArrayPointer(thumbnailLength, nil)
	if err != nil {
		return nil, err
	}
	defer thumbnailDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFPage_GetDecodedThumbnailData").Call(p.Context, *pageHandle.handle, thumbnailDataPointer.Pointer, thumbnailLength)
	if err != nil {
		return nil, err
	}

	thumbnailData, err := thumbnailDataPointer.Value(true)
	if err != nil {
		return nil, err
	}

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
	res, err := p.Module.ExportedFunction("FPDFPage_GetRawThumbnailData").Call(p.Context, *pageHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	rawThumbnailLength := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if rawThumbnailLength == 0 {
		return &responses.FPDFPage_GetRawThumbnailData{}, nil
	}

	rawThumbnailDataPointer, err := p.ByteArrayPointer(rawThumbnailLength, nil)
	if err != nil {
		return nil, err
	}
	defer rawThumbnailDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFPage_GetRawThumbnailData").Call(p.Context, *pageHandle.handle, rawThumbnailDataPointer.Pointer, rawThumbnailLength)
	if err != nil {
		return nil, err
	}

	rawThumbnailData, err := rawThumbnailDataPointer.Value(true)
	if err != nil {
		return nil, err
	}

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

	res, err := p.Module.ExportedFunction("FPDFPage_GetThumbnailAsBitmap").Call(p.Context, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	handle := res[0]
	if handle == 0 {
		return &responses.FPDFPage_GetThumbnailAsBitmap{}, nil
	}

	bitmapHandle := p.registerBitmap(&handle)

	return &responses.FPDFPage_GetThumbnailAsBitmap{
		Bitmap: &bitmapHandle.nativeRef,
	}, nil
}
