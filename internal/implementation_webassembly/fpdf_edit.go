package implementation_webassembly

import (
	"errors"
	"io"
	"os"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDF_CreateNewDocument returns a new document.
func (p *PdfiumImplementation) FPDF_CreateNewDocument(request *requests.FPDF_CreateNewDocument) (*responses.FPDF_CreateNewDocument, error) {
	p.Lock()
	defer p.Unlock()

	res, err := p.Module.ExportedFunction("FPDF_CreateNewDocument").Call(p.Context)
	if err != nil {
		return nil, err
	}

	doc := &res[0]
	documentHandle := p.registerDocument(doc)

	return &responses.FPDF_CreateNewDocument{
		Document: documentHandle.nativeRef,
	}, nil
}

// FPDFPage_New creates a new PDF page.
// The page should be closed with FPDF_ClosePage() when finished as
// with any other page in the document.
func (p *PdfiumImplementation) FPDFPage_New(request *requests.FPDFPage_New) (*responses.FPDFPage_New, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPage_New").Call(p.Context, *documentHandle.handle, *(*uint64)(unsafe.Pointer(&request.PageIndex)), *(*uint64)(unsafe.Pointer(&request.Width)), *(*uint64)(unsafe.Pointer(&request.Height)))
	if err != nil {
		return nil, err
	}

	page := res[0]
	pageHandle := p.registerPage(page, -1, documentHandle)

	return &responses.FPDFPage_New{
		Page: pageHandle.nativeRef,
	}, nil
}

// FPDFPage_Delete deletes the page at the given index.
func (p *PdfiumImplementation) FPDFPage_Delete(request *requests.FPDFPage_Delete) (*responses.FPDFPage_Delete, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDFPage_Delete").Call(p.Context, *documentHandle.handle, *(*uint64)(unsafe.Pointer(&request.PageIndex)))
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPage_Delete{}, nil
}

// FPDFPage_GetRotation returns the page rotation.
func (p *PdfiumImplementation) FPDFPage_GetRotation(request *requests.FPDFPage_GetRotation) (*responses.FPDFPage_GetRotation, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPage_GetRotation").Call(p.Context, *pageHandle.handle)
	if err != nil {
		return nil, err
	}
	rotation := res[0]

	return &responses.FPDFPage_GetRotation{
		Page:         pageHandle.index,
		PageRotation: enums.FPDF_PAGE_ROTATION(rotation),
	}, nil
}

// FPDFPage_SetRotation sets the page rotation for a given page.
func (p *PdfiumImplementation) FPDFPage_SetRotation(request *requests.FPDFPage_SetRotation) (*responses.FPDFPage_SetRotation, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDFPage_SetRotation").Call(p.Context, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Rotate)))
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPage_SetRotation{}, nil
}

// FPDFPage_InsertObject inserts the given object into a page.
func (p *PdfiumImplementation) FPDFPage_InsertObject(request *requests.FPDFPage_InsertObject) (*responses.FPDFPage_InsertObject, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDFPage_InsertObject").Call(p.Context, *pageHandle.handle, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPage_InsertObject{}, nil
}

// FPDFPage_CountObjects returns the number of page objects inside the given page.
func (p *PdfiumImplementation) FPDFPage_CountObjects(request *requests.FPDFPage_CountObjects) (*responses.FPDFPage_CountObjects, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPage_CountObjects").Call(p.Context, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFPage_CountObjects{
		Count: int(count),
	}, nil
}

// FPDFPage_GetObject returns the object at the given index.
func (p *PdfiumImplementation) FPDFPage_GetObject(request *requests.FPDFPage_GetObject) (*responses.FPDFPage_GetObject, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPage_GetObject").Call(p.Context, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	if res[0] == 0 {
		return nil, errors.New("could not get object")
	}

	pageObject := &res[0]
	pageObjectHandle := p.registerPageObject(pageObject)

	return &responses.FPDFPage_GetObject{
		PageObject: pageObjectHandle.nativeRef,
	}, nil
}

// FPDFPage_HasTransparency returns whether the page has transparency.
func (p *PdfiumImplementation) FPDFPage_HasTransparency(request *requests.FPDFPage_HasTransparency) (*responses.FPDFPage_HasTransparency, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPage_HasTransparency").Call(p.Context, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	alpha := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFPage_HasTransparency{
		Page:            pageHandle.index,
		HasTransparency: int(alpha) == 1,
	}, nil
}

// FPDFPage_GenerateContent generates the contents of the page.
func (p *PdfiumImplementation) FPDFPage_GenerateContent(request *requests.FPDFPage_GenerateContent) (*responses.FPDFPage_GenerateContent, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPage_GenerateContent").Call(p.Context, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not generate content")
	}

	return &responses.FPDFPage_GenerateContent{}, nil
}

// FPDFPageObj_Destroy destroys the page object by releasing its resources. The page object must have been
// created by FPDFPageObj_CreateNew{Path|Rect}() or
// FPDFPageObj_New{Text|Image}Obj(). This function must be called on
// newly-created objects if they are not added to a page through
// FPDFPage_InsertObject() or to an annotation through FPDFAnnot_AppendObject().
func (p *PdfiumImplementation) FPDFPageObj_Destroy(request *requests.FPDFPageObj_Destroy) (*responses.FPDFPageObj_Destroy, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDFPageObj_Destroy").Call(p.Context, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObj_Destroy{}, nil
}

// FPDFPageObj_HasTransparency returns whether the given page object contains transparency.
func (p *PdfiumImplementation) FPDFPageObj_HasTransparency(request *requests.FPDFPageObj_HasTransparency) (*responses.FPDFPageObj_HasTransparency, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_HasTransparency").Call(p.Context, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	hasTransparency := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFPageObj_HasTransparency{
		HasTransparency: int(hasTransparency) == 1,
	}, nil
}

// FPDFPageObj_GetType returns the type of the given page object.
func (p *PdfiumImplementation) FPDFPageObj_GetType(request *requests.FPDFPageObj_GetType) (*responses.FPDFPageObj_GetType, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_GetType").Call(p.Context, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	pageObjectType := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFPageObj_GetType{
		Type: enums.FPDF_PAGEOBJ(pageObjectType),
	}, nil
}

// FPDFPageObj_Transform transforms the page object by the given matrix.
// The matrix is composed as:
//
//	|a c e|
//	|b d f|
//
// and can be used to scale, rotate, shear and translate the page object.
func (p *PdfiumImplementation) FPDFPageObj_Transform(request *requests.FPDFPageObj_Transform) (*responses.FPDFPageObj_Transform, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDFPageObj_Transform").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.Transform.A)), *(*uint64)(unsafe.Pointer(&request.Transform.B)), *(*uint64)(unsafe.Pointer(&request.Transform.C)), *(*uint64)(unsafe.Pointer(&request.Transform.D)), *(*uint64)(unsafe.Pointer(&request.Transform.E)), *(*uint64)(unsafe.Pointer(&request.Transform.F)))
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObj_Transform{}, nil
}

// FPDFPage_TransformAnnots transforms all annotations in the given page.
// The matrix is composed as:
//
//	|a c e|
//	|b d f|
//
// and can be used to scale, rotate, shear and translate the page annotations.
func (p *PdfiumImplementation) FPDFPage_TransformAnnots(request *requests.FPDFPage_TransformAnnots) (*responses.FPDFPage_TransformAnnots, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDFPage_TransformAnnots").Call(p.Context, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Transform.A)), *(*uint64)(unsafe.Pointer(&request.Transform.B)), *(*uint64)(unsafe.Pointer(&request.Transform.C)), *(*uint64)(unsafe.Pointer(&request.Transform.D)), *(*uint64)(unsafe.Pointer(&request.Transform.E)), *(*uint64)(unsafe.Pointer(&request.Transform.F)))
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPage_TransformAnnots{}, nil
}

// FPDFPageObj_NewImageObj creates a new image object.
func (p *PdfiumImplementation) FPDFPageObj_NewImageObj(request *requests.FPDFPageObj_NewImageObj) (*responses.FPDFPageObj_NewImageObj, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_NewImageObj").Call(p.Context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	imageObject := &res[0]
	imageObjectHandle := p.registerPageObject(imageObject)

	return &responses.FPDFPageObj_NewImageObj{
		PageObject: imageObjectHandle.nativeRef,
	}, nil
}

// FPDFImageObj_LoadJpegFile loads an image from a JPEG image file and then set it into the given image object.
// The image object might already have an associated image, which is shared and
// cached by the loaded pages. In that case, we need to clear the cached image
// for all the loaded pages. Pass the pages and page count to this API
// to clear the image cache. If the image is not previously shared, nil is a
// valid pages value.
func (p *PdfiumImplementation) FPDFImageObj_LoadJpegFile(request *requests.FPDFImageObj_LoadJpegFile) (*responses.FPDFImageObj_LoadJpegFile, error) {
	p.Lock()
	defer p.Unlock()

	var pageHandle uint64
	if request.Page != nil {
		loadedPage, err := p.loadPage(*request.Page)
		if err != nil {
			return nil, err
		}

		pageHandle = *loadedPage.handle
	}

	pageObjectHandle, err := p.getPageObjectHandle(request.ImageObject)
	if err != nil {
		return nil, err
	}

	var fileReader io.ReadSeeker
	if request.FileReader != nil {
		fileReader = request.FileReader
	} else if request.FileData != nil {
		fileReader = NewBytesReaderCloser(request.FileData)
	} else {
		openedFile, err := os.Open(request.FilePath)
		if err != nil {
			return nil, err
		}

		fileReader = openedFile
	}

	fileAccessPointer, _, err := p.CreateFileAccessReader(request.FileReaderSize, fileReader)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFImageObj_LoadJpegFile").Call(p.Context, pageHandle, *(*uint64)(unsafe.Pointer(&request.Count)), *pageObjectHandle.handle, *fileAccessPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not load jpeg file")
	}

	return &responses.FPDFImageObj_LoadJpegFile{}, nil
}

// FPDFImageObj_LoadJpegFileInline
// The image object might already have an associated image, which is shared and
// cached by the loaded pages. In that case, we need to clear the cached image
// for all the loaded pages. Pass the pages and page count to this API
// to clear the image cache. If the image is not previously shared, nil is a
// valid pages value. This function loads the JPEG image inline, so the image
// content is copied to the file. This allows the file access and its associated
// data to be deleted after this function returns.
func (p *PdfiumImplementation) FPDFImageObj_LoadJpegFileInline(request *requests.FPDFImageObj_LoadJpegFileInline) (*responses.FPDFImageObj_LoadJpegFileInline, error) {
	p.Lock()
	defer p.Unlock()

	var pageHandle uint64
	if request.Page != nil {
		loadedPage, err := p.loadPage(*request.Page)
		if err != nil {
			return nil, err
		}

		pageHandle = *loadedPage.handle
	}

	pageObjectHandle, err := p.getPageObjectHandle(request.ImageObject)
	if err != nil {
		return nil, err
	}

	var fileReader io.ReadSeeker
	if request.FileReader != nil {
		fileReader = request.FileReader
	} else if request.FileData != nil {
		fileReader = NewBytesReaderCloser(request.FileData)
	} else {
		openedFile, err := os.Open(request.FilePath)
		if err != nil {
			return nil, err
		}

		fileReader = openedFile
	}

	fileAccessPointer, _, err := p.CreateFileAccessReader(request.FileReaderSize, fileReader)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFImageObj_LoadJpegFileInline").Call(p.Context, pageHandle, *(*uint64)(unsafe.Pointer(&request.Count)), *pageObjectHandle.handle, *fileAccessPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not load jpeg file")
	}

	return &responses.FPDFImageObj_LoadJpegFileInline{}, nil
}

// FPDFImageObj_SetMatrix sets the transform matrix of the given image object.
// The matrix is composed as:
//
//	|a c e|
//	|b d f|
//
// and can be used to scale, rotate, shear and translate the image object.
// Will be deprecated once FPDFPageObj_SetMatrix() is stable.
func (p *PdfiumImplementation) FPDFImageObj_SetMatrix(request *requests.FPDFImageObj_SetMatrix) (*responses.FPDFImageObj_SetMatrix, error) {
	p.Lock()
	defer p.Unlock()

	imageObjectHandle, err := p.getPageObjectHandle(request.ImageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFImageObj_SetMatrix").Call(p.Context, *imageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.Transform.A)), *(*uint64)(unsafe.Pointer(&request.Transform.B)), *(*uint64)(unsafe.Pointer(&request.Transform.C)), *(*uint64)(unsafe.Pointer(&request.Transform.D)), *(*uint64)(unsafe.Pointer(&request.Transform.E)), *(*uint64)(unsafe.Pointer(&request.Transform.F)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set image object matrix")
	}

	return &responses.FPDFImageObj_SetMatrix{}, nil
}

// FPDFImageObj_SetBitmap sets the given bitmap to the given image object.
func (p *PdfiumImplementation) FPDFImageObj_SetBitmap(request *requests.FPDFImageObj_SetBitmap) (*responses.FPDFImageObj_SetBitmap, error) {
	p.Lock()
	defer p.Unlock()

	var pageHandle uint64
	if request.Page != nil {
		loadedPage, err := p.loadPage(*request.Page)
		if err != nil {
			return nil, err
		}

		pageHandle = *loadedPage.handle
	}

	imageObjectHandle, err := p.getPageObjectHandle(request.ImageObject)
	if err != nil {
		return nil, err
	}

	bitmapHandle, err := p.getBitmapHandle(request.Bitmap)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFImageObj_SetBitmap").Call(p.Context, pageHandle, *(*uint64)(unsafe.Pointer(&request.Count)), *imageObjectHandle.handle, *bitmapHandle.handle)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set image object bitmap")
	}

	return &responses.FPDFImageObj_SetBitmap{}, nil
}

// FPDFImageObj_GetBitmap returns a bitmap rasterization of the given image object. FPDFImageObj_GetBitmap() only
// operates on the image object and does not take the associated image mask into
// account. It also ignores the matrix for the image object.
// The returned bitmap will be owned by the caller, and FPDFBitmap_Destroy()
// must be called on the returned bitmap when it is no longer needed.
func (p *PdfiumImplementation) FPDFImageObj_GetBitmap(request *requests.FPDFImageObj_GetBitmap) (*responses.FPDFImageObj_GetBitmap, error) {
	p.Lock()
	defer p.Unlock()

	imageObjectHandle, err := p.getPageObjectHandle(request.ImageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFImageObj_GetBitmap").Call(p.Context, *imageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	bitmap := res[0]
	if bitmap == 0 {
		return nil, errors.New("could not get bitmap")
	}

	bitmapHandle := p.registerBitmap(&bitmap)

	return &responses.FPDFImageObj_GetBitmap{
		Bitmap: bitmapHandle.nativeRef,
	}, nil
}

// FPDFImageObj_GetImageDataDecoded returns the decoded image data of the image object. The decoded data is the
// uncompressed image data, i.e. the raw image data after having all filters
// applied.
func (p *PdfiumImplementation) FPDFImageObj_GetImageDataDecoded(request *requests.FPDFImageObj_GetImageDataDecoded) (*responses.FPDFImageObj_GetImageDataDecoded, error) {
	p.Lock()
	defer p.Unlock()

	imageObjectHandle, err := p.getPageObjectHandle(request.ImageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFImageObj_GetImageDataDecoded").Call(p.Context, *imageObjectHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	imageDataLength := *(*int32)(unsafe.Pointer(&res[0]))
	if int(imageDataLength) == 0 {
		return nil, errors.New("could not get decoded image data")
	}

	length := uint64(imageDataLength)

	valueData, err := p.ByteArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer valueData.Free()

	_, err = p.Module.ExportedFunction("FPDFImageObj_GetImageDataDecoded").Call(p.Context, *imageObjectHandle.handle, valueData.Pointer, *(*uint64)(unsafe.Pointer(&length)))
	if err != nil {
		return nil, err
	}

	data, err := valueData.Value(true)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFImageObj_GetImageDataDecoded{
		Data: data,
	}, nil
}

// FPDFImageObj_GetImageDataRaw returns the raw image data of the image object. The raw data is the image data as
// stored in the PDF without applying any filters.
func (p *PdfiumImplementation) FPDFImageObj_GetImageDataRaw(request *requests.FPDFImageObj_GetImageDataRaw) (*responses.FPDFImageObj_GetImageDataRaw, error) {
	p.Lock()
	defer p.Unlock()

	imageObjectHandle, err := p.getPageObjectHandle(request.ImageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFImageObj_GetImageDataRaw").Call(p.Context, *imageObjectHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	imageDataLength := *(*int32)(unsafe.Pointer(&res[0]))
	if int(imageDataLength) == 0 {
		return nil, errors.New("could not get raw image data")
	}

	length := uint64(imageDataLength)

	valueData, err := p.ByteArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer valueData.Free()

	_, err = p.Module.ExportedFunction("FPDFImageObj_GetImageDataRaw").Call(p.Context, *imageObjectHandle.handle, valueData.Pointer, *(*uint64)(unsafe.Pointer(&length)))
	if err != nil {
		return nil, err
	}

	data, err := valueData.Value(true)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFImageObj_GetImageDataRaw{
		Data: data,
	}, nil
}

// FPDFImageObj_GetImageFilterCount returns the number of filters (i.e. decoders) of the image in image object.
func (p *PdfiumImplementation) FPDFImageObj_GetImageFilterCount(request *requests.FPDFImageObj_GetImageFilterCount) (*responses.FPDFImageObj_GetImageFilterCount, error) {
	p.Lock()
	defer p.Unlock()

	imageObjectHandle, err := p.getPageObjectHandle(request.ImageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFImageObj_GetImageFilterCount").Call(p.Context, *imageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFImageObj_GetImageFilterCount{
		Count: int(count),
	}, nil
}

// FPDFImageObj_GetImageFilter returns the filter at index of the image object's list of filters. Note that the
// filters need to be applied in order, i.e. the first filter should be applied
// first, then the second, etc.
func (p *PdfiumImplementation) FPDFImageObj_GetImageFilter(request *requests.FPDFImageObj_GetImageFilter) (*responses.FPDFImageObj_GetImageFilter, error) {
	p.Lock()
	defer p.Unlock()

	imageObjectHandle, err := p.getPageObjectHandle(request.ImageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFImageObj_GetImageFilter").Call(p.Context, *imageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)), 0, 0)
	if err != nil {
		return nil, err
	}

	imageFilterLength := *(*int32)(unsafe.Pointer(&res[0]))
	if int(imageFilterLength) == 0 {
		return nil, errors.New("could not get image filter")
	}

	length := uint64(imageFilterLength)
	charDataPointer, err := p.ByteArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFImageObj_GetImageFilter").Call(p.Context, *imageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)), charDataPointer.Pointer, *(*uint64)(unsafe.Pointer(&length)))
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(false)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFImageObj_GetImageFilter{
		ImageFilter: string(charData[:len(charData)-1]), // Remove NULL-terminator.
	}, nil
}

// FPDFImageObj_GetImageMetadata returns the image metadata of the image object, including dimension, DPI, bits per
// pixel, and colorspace. If the image object is not an image object or if it
// does not have an image, then the return value will be false. Otherwise,
// failure to retrieve any specific parameter would result in its value being 0.
func (p *PdfiumImplementation) FPDFImageObj_GetImageMetadata(request *requests.FPDFImageObj_GetImageMetadata) (*responses.FPDFImageObj_GetImageMetadata, error) {
	p.Lock()
	defer p.Unlock()

	imageObjectHandle, err := p.getPageObjectHandle(request.ImageObject)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	metadataPointer, metadataValue, err := p.CStructFPDF_IMAGEOBJ_METADATA(nil)
	if err != nil {
		return nil, err
	}

	defer p.Free(metadataPointer)

	res, err := p.Module.ExportedFunction("FPDFImageObj_GetImageMetadata").Call(p.Context, *imageObjectHandle.handle, *pageHandle.handle, metadataPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not load image metadata")
	}

	metadata, err := metadataValue()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFImageObj_GetImageMetadata{
		ImageMetadata: *metadata,
	}, nil
}
