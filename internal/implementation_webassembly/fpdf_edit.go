package implementation_webassembly

import (
	"errors"
	"io"
	"os"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"

	"github.com/tetratelabs/wazero/api"
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

// FPDF_MovePages Move the given pages to a new index position.
// When this call fails, the document may be left in an indeterminate state.
// Experimental API.
func (p *PdfiumImplementation) FPDF_MovePages(request *requests.FPDF_MovePages) (*responses.FPDF_MovePages, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	if len(request.PageIndices) == 0 {
		return nil, errors.New("no page indices were given")
	}

	pageIndicesSize := uint64(uint64(len(request.PageIndices)))

	// Create an array that's big enough.
	valueDataPointer, err := p.IntArrayPointer(pageIndicesSize)
	if err != nil {
		return nil, err
	}
	defer valueDataPointer.Free()

	// Put the values in the array.
	for i := range request.PageIndices {
		p.Module.Memory().WriteUint32Le(uint32(valueDataPointer.Pointer+(p.CSizeInt()*uint64(i))), uint32(request.PageIndices[i]))
	}

	res, err := p.Module.ExportedFunction("FPDF_MovePages").Call(p.Context, *documentHandle.handle, valueDataPointer.Pointer, pageIndicesSize, *(*uint64)(unsafe.Pointer(&request.DestPageIndex)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not move pages")
	}

	return &responses.FPDF_MovePages{}, nil
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

// FPDFPage_InsertObjectAtIndex inserts the given object into a page at a specific index.
func (p *PdfiumImplementation) FPDFPage_InsertObjectAtIndex(request *requests.FPDFPage_InsertObjectAtIndex) (*responses.FPDFPage_InsertObjectAtIndex, error) {
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

	_, err = p.Module.ExportedFunction("FPDFPage_InsertObjectAtIndex").Call(p.Context, *pageHandle.handle, *pageObjectHandle.handle, uint64(request.Index))
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPage_InsertObjectAtIndex{}, nil
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

// FPDFPageObj_GetIsActive returns the active state for the given page
// object within the page.
// For page objects where active is filled with false, the page object is
// treated as if it wasn't in the document even though it is still held
// internally.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetIsActive(request *requests.FPDFPageObj_GetIsActive) (*responses.FPDFPageObj_GetIsActive, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	activePointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}

	defer activePointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObj_GetIsActive").Call(p.Context, *pageObjectHandle.handle, activePointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get active state")
	}

	isActive, err := activePointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObj_GetIsActive{
		Active: int(isActive) == 1,
	}, nil
}

// FPDFPageObj_SetIsActive sets the active state for the given page object
// within the page.
// Page objects all start in the active state by default, and remain in that
// state unless this function is called.
// When active is false, this makes the page_object be treated as if it
// wasn't in the document even though it is still held internally.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_SetIsActive(request *requests.FPDFPageObj_SetIsActive) (*responses.FPDFPageObj_SetIsActive, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	isActive := 0
	if request.Active {
		isActive = 1
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_SetIsActive").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&isActive)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set active state")
	}

	return &responses.FPDFPageObj_SetIsActive{}, nil
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

	_, err = p.Module.ExportedFunction("FPDFPageObj_Transform").Call(p.Context, *pageObjectHandle.handle, api.EncodeF64(float64(request.Transform.A)), api.EncodeF64(float64(request.Transform.B)), api.EncodeF64(float64(request.Transform.C)), api.EncodeF64(float64(request.Transform.D)), api.EncodeF64(float64(request.Transform.E)), api.EncodeF64(float64(request.Transform.F)))
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

	_, err = p.Module.ExportedFunction("FPDFPage_TransformAnnots").Call(p.Context, *pageHandle.handle, api.EncodeF64(float64(request.Transform.A)), api.EncodeF64(float64(request.Transform.B)), api.EncodeF64(float64(request.Transform.C)), api.EncodeF64(float64(request.Transform.D)), api.EncodeF64(float64(request.Transform.E)), api.EncodeF64(float64(request.Transform.F)))
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
		request.FileReaderSize = int64(len(request.FileData))
	} else {
		openedFile, err := os.Open(request.FilePath)
		if err != nil {
			return nil, err
		}

		stat, err := openedFile.Stat()
		if err != nil {
			return nil, err
		}

		request.FileReaderSize = stat.Size()
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
		request.FileReaderSize = int64(len(request.FileData))
	} else {
		openedFile, err := os.Open(request.FilePath)
		if err != nil {
			return nil, err
		}

		stat, err := openedFile.Stat()
		if err != nil {
			return nil, err
		}

		request.FileReaderSize = stat.Size()
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

	res, err := p.Module.ExportedFunction("FPDFImageObj_SetMatrix").Call(p.Context, *imageObjectHandle.handle, api.EncodeF64(float64(request.Transform.A)), api.EncodeF64(float64(request.Transform.B)), api.EncodeF64(float64(request.Transform.C)), api.EncodeF64(float64(request.Transform.D)), api.EncodeF64(float64(request.Transform.E)), api.EncodeF64(float64(request.Transform.F)))
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

// FPDFImageObj_GetImagePixelSize get the image size in pixels. Faster method to get only image size.
// Experimental API.
func (p *PdfiumImplementation) FPDFImageObj_GetImagePixelSize(request *requests.FPDFImageObj_GetImagePixelSize) (*responses.FPDFImageObj_GetImagePixelSize, error) {
	p.Lock()
	defer p.Unlock()

	imageObjectHandle, err := p.getPageObjectHandle(request.ImageObject)
	if err != nil {
		return nil, err
	}

	widthPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer widthPointer.Free()

	heightPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer heightPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFImageObj_GetImagePixelSize").Call(p.Context, *imageObjectHandle.handle, widthPointer.Pointer, heightPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get image pixel size")
	}

	widthValue, err := widthPointer.Value()
	if err != nil {
		return nil, err
	}

	heightValue, err := heightPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFImageObj_GetImagePixelSize{
		Width:  widthValue,
		Height: heightValue,
	}, nil
}

// FPDFImageObj_GetIccProfileDataDecoded returns the ICC profile decoded
// data of the given image object. If the image object is not an image
// object or if it does not have an image, then the return value will
// be nil. It also returns nil if the image object has no ICC profile.
// Experimental API.
func (p *PdfiumImplementation) FPDFImageObj_GetIccProfileDataDecoded(request *requests.FPDFImageObj_GetIccProfileDataDecoded) (*responses.FPDFImageObj_GetIccProfileDataDecoded, error) {
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

	iccProfileDataLengthPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer iccProfileDataLengthPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFImageObj_GetIccProfileDataDecoded").Call(p.Context, *imageObjectHandle.handle, *pageHandle.handle, 0, 0, iccProfileDataLengthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get ICC Profile Data length")
	}

	iccProfileDataLength, err := iccProfileDataLengthPointer.Value()
	if err != nil {
		return nil, err
	}

	length := uint64(iccProfileDataLength)

	valueData, err := p.ByteArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer valueData.Free()

	res, err = p.Module.ExportedFunction("FPDFImageObj_GetIccProfileDataDecoded").Call(p.Context, *imageObjectHandle.handle, *pageHandle.handle, valueData.Pointer, *(*uint64)(unsafe.Pointer(&length)), iccProfileDataLengthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success = *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get ICC Profile Data")
	}

	data, err := valueData.Value(true)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFImageObj_GetIccProfileDataDecoded{
		Data: data,
	}, nil
}

// FPDFPageObj_CreateNewPath creates a new path object at an initial position.
func (p *PdfiumImplementation) FPDFPageObj_CreateNewPath(request *requests.FPDFPageObj_CreateNewPath) (*responses.FPDFPageObj_CreateNewPath, error) {
	p.Lock()
	defer p.Unlock()

	res, err := p.Module.ExportedFunction("FPDFPageObj_CreateNewPath").Call(p.Context, *(*uint64)(unsafe.Pointer(&request.X)), *(*uint64)(unsafe.Pointer(&request.Y)))
	if err != nil {
		return nil, err
	}

	pageObject := res[0]
	pageObjectHandle := p.registerPageObject(&pageObject)

	return &responses.FPDFPageObj_CreateNewPath{
		PageObject: pageObjectHandle.nativeRef,
	}, nil
}

// FPDFPageObj_CreateNewRect creates a closed path consisting of a rectangle.
func (p *PdfiumImplementation) FPDFPageObj_CreateNewRect(request *requests.FPDFPageObj_CreateNewRect) (*responses.FPDFPageObj_CreateNewRect, error) {
	p.Lock()
	defer p.Unlock()

	res, err := p.Module.ExportedFunction("FPDFPageObj_CreateNewRect").Call(p.Context, *(*uint64)(unsafe.Pointer(&request.X)), *(*uint64)(unsafe.Pointer(&request.Y)), *(*uint64)(unsafe.Pointer(&request.W)), *(*uint64)(unsafe.Pointer(&request.H)))
	if err != nil {
		return nil, err
	}

	pageObject := res[0]
	pageObjectHandle := p.registerPageObject(&pageObject)
	return &responses.FPDFPageObj_CreateNewRect{
		PageObject: pageObjectHandle.nativeRef,
	}, nil
}

// FPDFPageObj_GetBounds returns the bounding box of the given page object.
func (p *PdfiumImplementation) FPDFPageObj_GetBounds(request *requests.FPDFPageObj_GetBounds) (*responses.FPDFPageObj_GetBounds, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	leftPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer leftPointer.Free()

	bottomPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer bottomPointer.Free()

	rightPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer rightPointer.Free()

	topPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer topPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObj_GetBounds").Call(p.Context, *pageObjectHandle.handle, leftPointer.Pointer, bottomPointer.Pointer, rightPointer.Pointer, topPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get page object bounds")
	}

	left, err := leftPointer.Value()
	if err != nil {
		return nil, err
	}

	bottom, err := bottomPointer.Value()
	if err != nil {
		return nil, err
	}

	right, err := rightPointer.Value()
	if err != nil {
		return nil, err
	}

	top, err := topPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObj_GetBounds{
		Left:   float32(left),
		Bottom: float32(bottom),
		Right:  float32(right),
		Top:    float32(top),
	}, nil
}

// FPDFPageObj_SetBlendMode sets the blend mode of the page object.
func (p *PdfiumImplementation) FPDFPageObj_SetBlendMode(request *requests.FPDFPageObj_SetBlendMode) (*responses.FPDFPageObj_SetBlendMode, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	blendMode, err := p.CString(string(request.BlendMode))
	if err != nil {
		return nil, err
	}

	defer blendMode.Free()

	_, err = p.Module.ExportedFunction("FPDFPageObj_SetBlendMode").Call(p.Context, *pageObjectHandle.handle, blendMode.Pointer)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObj_SetBlendMode{}, nil
}

// FPDFPageObj_SetStrokeColor sets the stroke RGBA of a page object.
func (p *PdfiumImplementation) FPDFPageObj_SetStrokeColor(request *requests.FPDFPageObj_SetStrokeColor) (*responses.FPDFPageObj_SetStrokeColor, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_SetStrokeColor").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.StrokeColor.R)), *(*uint64)(unsafe.Pointer(&request.StrokeColor.G)), *(*uint64)(unsafe.Pointer(&request.StrokeColor.B)), *(*uint64)(unsafe.Pointer(&request.StrokeColor.A)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set page object stroke color")
	}

	return &responses.FPDFPageObj_SetStrokeColor{}, nil
}

// FPDFPageObj_GetStrokeColor returns the stroke RGBA of a page object
func (p *PdfiumImplementation) FPDFPageObj_GetStrokeColor(request *requests.FPDFPageObj_GetStrokeColor) (*responses.FPDFPageObj_GetStrokeColor, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	rPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer rPointer.Free()

	gPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer gPointer.Free()

	bPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer bPointer.Free()

	aPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer aPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObj_GetStrokeColor").Call(p.Context, *pageObjectHandle.handle, rPointer.Pointer, gPointer.Pointer, bPointer.Pointer, aPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get page object stroke color")
	}

	r, err := rPointer.Value()
	if err != nil {
		return nil, err
	}

	g, err := gPointer.Value()
	if err != nil {
		return nil, err
	}

	b, err := bPointer.Value()
	if err != nil {
		return nil, err
	}

	a, err := aPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObj_GetStrokeColor{
		StrokeColor: structs.FPDF_COLOR{
			R: uint(r),
			G: uint(g),
			B: uint(b),
			A: uint(a),
		},
	}, nil
}

// FPDFPageObj_SetStrokeWidth sets the stroke width of a page object
func (p *PdfiumImplementation) FPDFPageObj_SetStrokeWidth(request *requests.FPDFPageObj_SetStrokeWidth) (*responses.FPDFPageObj_SetStrokeWidth, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_SetStrokeWidth").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.StrokeWidth)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set page object stroke color")
	}

	return &responses.FPDFPageObj_SetStrokeWidth{}, nil
}

// FPDFPageObj_GetStrokeWidth returns the stroke width of a page object.
func (p *PdfiumImplementation) FPDFPageObj_GetStrokeWidth(request *requests.FPDFPageObj_GetStrokeWidth) (*responses.FPDFPageObj_GetStrokeWidth, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	strokeWidthPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer strokeWidthPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObj_GetStrokeWidth").Call(p.Context, *pageObjectHandle.handle, strokeWidthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get page object stroke width")
	}

	strokeWidth, err := strokeWidthPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObj_GetStrokeWidth{
		StrokeWidth: float32(strokeWidth),
	}, nil
}

// FPDFPageObj_GetLineJoin returns the line join of the page object.
func (p *PdfiumImplementation) FPDFPageObj_GetLineJoin(request *requests.FPDFPageObj_GetLineJoin) (*responses.FPDFPageObj_GetLineJoin, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_GetLineJoin").Call(p.Context, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	lineJoin := *(*int32)(unsafe.Pointer(&res[0]))
	if int(lineJoin) == -1 {
		return nil, errors.New("could not get page object line join")
	}

	return &responses.FPDFPageObj_GetLineJoin{
		LineJoin: enums.FPDF_LINEJOIN(lineJoin),
	}, nil
}

// FPDFPageObj_SetLineJoin sets the line join of the page object.
func (p *PdfiumImplementation) FPDFPageObj_SetLineJoin(request *requests.FPDFPageObj_SetLineJoin) (*responses.FPDFPageObj_SetLineJoin, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_SetLineJoin").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.LineJoin)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set page object line join")
	}

	return &responses.FPDFPageObj_SetLineJoin{}, nil
}

// FPDFPageObj_GetLineCap returns the line cap of the page object.
func (p *PdfiumImplementation) FPDFPageObj_GetLineCap(request *requests.FPDFPageObj_GetLineCap) (*responses.FPDFPageObj_GetLineCap, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_GetLineCap").Call(p.Context, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	lineCap := *(*int32)(unsafe.Pointer(&res[0]))
	if int(lineCap) == -1 {
		return nil, errors.New("could not get page object line cap")
	}

	return &responses.FPDFPageObj_GetLineCap{
		LineCap: enums.FPDF_LINECAP(lineCap),
	}, nil
}

// FPDFPageObj_SetLineCap sets the line cap of the page object.
func (p *PdfiumImplementation) FPDFPageObj_SetLineCap(request *requests.FPDFPageObj_SetLineCap) (*responses.FPDFPageObj_SetLineCap, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_SetLineJoin").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.LineCap)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set page object line cap")
	}

	return &responses.FPDFPageObj_SetLineCap{}, nil
}

// FPDFPageObj_SetFillColor sets the fill RGBA of a page object
func (p *PdfiumImplementation) FPDFPageObj_SetFillColor(request *requests.FPDFPageObj_SetFillColor) (*responses.FPDFPageObj_SetFillColor, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_SetFillColor").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.FillColor.R)), *(*uint64)(unsafe.Pointer(&request.FillColor.G)), *(*uint64)(unsafe.Pointer(&request.FillColor.B)), *(*uint64)(unsafe.Pointer(&request.FillColor.A)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set page object fill color")
	}

	return &responses.FPDFPageObj_SetFillColor{}, nil
}

// FPDFPageObj_GetFillColor returns the fill RGBA of a page object
func (p *PdfiumImplementation) FPDFPageObj_GetFillColor(request *requests.FPDFPageObj_GetFillColor) (*responses.FPDFPageObj_GetFillColor, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	rPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer rPointer.Free()

	gPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer gPointer.Free()

	bPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer bPointer.Free()

	aPointer, err := p.UIntPointer()
	if err != nil {
		return nil, err
	}
	defer aPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObj_GetFillColor").Call(p.Context, *pageObjectHandle.handle, rPointer.Pointer, gPointer.Pointer, bPointer.Pointer, aPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get page object fill color")
	}

	r, err := rPointer.Value()
	if err != nil {
		return nil, err
	}

	g, err := gPointer.Value()
	if err != nil {
		return nil, err
	}

	b, err := bPointer.Value()
	if err != nil {
		return nil, err
	}

	a, err := aPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObj_GetFillColor{
		FillColor: structs.FPDF_COLOR{
			R: uint(r),
			G: uint(g),
			B: uint(b),
			A: uint(a),
		},
	}, nil
}

// FPDFPath_CountSegments returns the number of segments inside the given path.
// A segment is a command, created by e.g. FPDFPath_MoveTo(),
// FPDFPath_LineTo() or FPDFPath_BezierTo().
func (p *PdfiumImplementation) FPDFPath_CountSegments(request *requests.FPDFPath_CountSegments) (*responses.FPDFPath_CountSegments, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPath_CountSegments").Call(p.Context, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))
	if int(count) == -1 {
		return nil, errors.New("could not get path segment count")
	}

	return &responses.FPDFPath_CountSegments{
		Count: int(count),
	}, nil
}

// FPDFPath_GetPathSegment returns the segment in the given path at the given index.
func (p *PdfiumImplementation) FPDFPath_GetPathSegment(request *requests.FPDFPath_GetPathSegment) (*responses.FPDFPath_GetPathSegment, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPath_GetPathSegment").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	pathSegment := res[0]
	if pathSegment == 0 {
		return nil, errors.New("could not get path segment")
	}

	pathSegmentHandle := p.registerPathSegment(&pathSegment)

	return &responses.FPDFPath_GetPathSegment{
		PathSegment: pathSegmentHandle.nativeRef,
	}, nil
}

// FPDFPathSegment_GetPoint returns the coordinates of the given segment.
func (p *PdfiumImplementation) FPDFPathSegment_GetPoint(request *requests.FPDFPathSegment_GetPoint) (*responses.FPDFPathSegment_GetPoint, error) {
	p.Lock()
	defer p.Unlock()

	pathSegmentHandle, err := p.getPathSegmentHandle(request.PathSegment)
	if err != nil {
		return nil, err
	}

	xPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer xPointer.Free()

	yPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}
	defer yPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPathSegment_GetPoint").Call(p.Context, *pathSegmentHandle.handle, xPointer.Pointer, yPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get path segment point")
	}

	x, err := xPointer.Value()
	if err != nil {
		return nil, err
	}

	y, err := yPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPathSegment_GetPoint{
		X: float32(x),
		Y: float32(y),
	}, nil
}

// FPDFPathSegment_GetType returns the type of the given segment.
func (p *PdfiumImplementation) FPDFPathSegment_GetType(request *requests.FPDFPathSegment_GetType) (*responses.FPDFPathSegment_GetType, error) {
	p.Lock()
	defer p.Unlock()

	pathSegmentHandle, err := p.getPathSegmentHandle(request.PathSegment)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPathSegment_GetType").Call(p.Context, *pathSegmentHandle.handle)
	if err != nil {
		return nil, err
	}

	segmentType := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFPathSegment_GetType{
		Type: enums.FPDF_SEGMENT(segmentType),
	}, nil
}

// FPDFPathSegment_GetClose returns whether the segment closes the current subpath of a given path.
func (p *PdfiumImplementation) FPDFPathSegment_GetClose(request *requests.FPDFPathSegment_GetClose) (*responses.FPDFPathSegment_GetClose, error) {
	p.Lock()
	defer p.Unlock()

	pathSegmentHandle, err := p.getPathSegmentHandle(request.PathSegment)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPathSegment_GetClose").Call(p.Context, *pathSegmentHandle.handle)
	if err != nil {
		return nil, err
	}

	getClose := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFPathSegment_GetClose{
		IsClose: int(getClose) == 1,
	}, nil
}

// FPDFPath_MoveTo moves a path's current point.
// Note that no line will be created between the previous current point and the
// new one.
func (p *PdfiumImplementation) FPDFPath_MoveTo(request *requests.FPDFPath_MoveTo) (*responses.FPDFPath_MoveTo, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPath_MoveTo").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.X)), *(*uint64)(unsafe.Pointer(&request.Y)))
	if err != nil {
		return nil, err
	}

	result := *(*int32)(unsafe.Pointer(&res[0]))
	if int(result) == 0 {
		return nil, errors.New("could not move path")
	}

	return &responses.FPDFPath_MoveTo{}, nil
}

// FPDFPath_LineTo adds a line between the current point and a new point in the path.
func (p *PdfiumImplementation) FPDFPath_LineTo(request *requests.FPDFPath_LineTo) (*responses.FPDFPath_LineTo, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPath_LineTo").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.X)), *(*uint64)(unsafe.Pointer(&request.Y)))
	if err != nil {
		return nil, err
	}

	result := *(*int32)(unsafe.Pointer(&res[0]))
	if int(result) == 0 {
		return nil, errors.New("could not add line")
	}

	return &responses.FPDFPath_LineTo{}, nil
}

// FPDFPath_BezierTo adds a cubic Bezier curve to the given path, starting at the current point.
func (p *PdfiumImplementation) FPDFPath_BezierTo(request *requests.FPDFPath_BezierTo) (*responses.FPDFPath_BezierTo, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPath_BezierTo").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.X1)), *(*uint64)(unsafe.Pointer(&request.Y1)), *(*uint64)(unsafe.Pointer(&request.X2)), *(*uint64)(unsafe.Pointer(&request.Y2)), *(*uint64)(unsafe.Pointer(&request.X3)), *(*uint64)(unsafe.Pointer(&request.Y3)))
	if err != nil {
		return nil, err
	}

	result := *(*int32)(unsafe.Pointer(&res[0]))
	if int(result) == 0 {
		return nil, errors.New("could not add bezier")
	}

	return &responses.FPDFPath_BezierTo{}, nil
}

// FPDFPath_Close closes the current subpath of a given path.
func (p *PdfiumImplementation) FPDFPath_Close(request *requests.FPDFPath_Close) (*responses.FPDFPath_Close, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPath_Close").Call(p.Context, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	result := *(*int32)(unsafe.Pointer(&res[0]))
	if int(result) == 0 {
		return nil, errors.New("could not close path")
	}

	return &responses.FPDFPath_Close{}, nil
}

// FPDFPath_SetDrawMode sets the drawing mode of a path.
func (p *PdfiumImplementation) FPDFPath_SetDrawMode(request *requests.FPDFPath_SetDrawMode) (*responses.FPDFPath_SetDrawMode, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	stroke := 0
	if request.Stroke {
		stroke = 1
	}

	res, err := p.Module.ExportedFunction("FPDFPath_SetDrawMode").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.FillMode)), *(*uint64)(unsafe.Pointer(&stroke)))
	if err != nil {
		return nil, err
	}

	result := *(*int32)(unsafe.Pointer(&res[0]))
	if int(result) == 0 {
		return nil, errors.New("could not set draw mode")
	}

	return &responses.FPDFPath_SetDrawMode{}, nil
}

// FPDFPath_GetDrawMode returns the drawing mode of a path.
func (p *PdfiumImplementation) FPDFPath_GetDrawMode(request *requests.FPDFPath_GetDrawMode) (*responses.FPDFPath_GetDrawMode, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	fillModePointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer fillModePointer.Free()

	strokePointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer strokePointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPath_GetDrawMode").Call(p.Context, *pageObjectHandle.handle, fillModePointer.Pointer, strokePointer.Pointer)
	if err != nil {
		return nil, err
	}

	result := *(*int32)(unsafe.Pointer(&res[0]))
	if int(result) == 0 {
		return nil, errors.New("could not get draw mode")
	}

	fillMode, err := fillModePointer.Value()
	if err != nil {
		return nil, err
	}

	stroke, err := strokePointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPath_GetDrawMode{
		FillMode: enums.FPDF_FILLMODE(fillMode),
		Stroke:   int(stroke) == 1,
	}, nil
}

// FPDFPageObj_NewTextObj creates a new text object using one of the standard PDF fonts.
func (p *PdfiumImplementation) FPDFPageObj_NewTextObj(request *requests.FPDFPageObj_NewTextObj) (*responses.FPDFPageObj_NewTextObj, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	font, err := p.CString(request.Font)
	defer font.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObj_NewTextObj").Call(p.Context, *documentHandle.handle, font.Pointer, *(*uint64)(unsafe.Pointer(&request.FontSize)))
	if err != nil {
		return nil, err
	}

	textObject := res[0]
	textObjectHandle := p.registerPageObject(&textObject)

	return &responses.FPDFPageObj_NewTextObj{
		PageObject: textObjectHandle.nativeRef,
	}, nil
}

// FPDFText_SetText sets the text for a text object. If it had text, it will be replaced.
func (p *PdfiumImplementation) FPDFText_SetText(request *requests.FPDFText_SetText) (*responses.FPDFText_SetText, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	transformedTextPointer, err := p.CFPDF_WIDESTRING(request.Text)
	if err != nil {
		return nil, err
	}

	defer transformedTextPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFText_SetText").Call(p.Context, *pageObjectHandle.handle, transformedTextPointer.Pointer)
	if err != nil {
		return nil, err
	}

	result := *(*int32)(unsafe.Pointer(&res[0]))
	if int(result) == 0 {
		return nil, errors.New("could not set text")
	}

	return &responses.FPDFText_SetText{}, nil
}

// FPDFText_SetCharcodes sets the text using charcodes for a text object. If it had text, it will be
// replaced.
func (p *PdfiumImplementation) FPDFText_SetCharcodes(request *requests.FPDFText_SetCharcodes) (*responses.FPDFText_SetCharcodes, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	length := uint64(len(request.CharCodes))
	charCodes, err := p.UIntArrayPointer(length)
	if err != nil {
		return nil, err
	}

	defer charCodes.Free()
	for i := range request.CharCodes {
		success := p.Module.Memory().WriteUint32Le(uint32(charCodes.Pointer+(p.CSizeUInt()*uint64(i))), request.CharCodes[i])
		if !success {
			return nil, errors.New("could not write uint array data to memory")
		}
	}

	res, err := p.Module.ExportedFunction("FPDFText_SetCharcodes").Call(p.Context, *pageObjectHandle.handle, charCodes.Pointer, *(*uint64)(unsafe.Pointer(&length)))
	if err != nil {
		return nil, err
	}

	result := *(*int32)(unsafe.Pointer(&res[0]))
	if int(result) == 0 {
		return nil, errors.New("could not set charcodes")
	}

	return &responses.FPDFText_SetCharcodes{}, nil
}

// FPDFText_LoadFont returns a font object loaded from a stream of data. The font is loaded
// into the document. Various font data structures, such as the ToUnicode data, are auto-generated based
// on the inputs.
// The loaded font can be closed using FPDFFont_Close().
func (p *PdfiumImplementation) FPDFText_LoadFont(request *requests.FPDFText_LoadFont) (*responses.FPDFText_LoadFont, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	cid := 0
	if request.CID {
		cid = 1
	}

	dataLength := uint32(len(request.Data))
	fontData, err := p.ByteArrayPointer(uint64(len(request.Data)), request.Data)
	if err != nil {
		return nil, err
	}

	defer fontData.Free()

	res, err := p.Module.ExportedFunction("FPDFText_LoadFont").Call(p.Context, *documentHandle.handle, fontData.Pointer, *(*uint64)(unsafe.Pointer(&dataLength)), *(*uint64)(unsafe.Pointer(&request.FontType)), *(*uint64)(unsafe.Pointer(&cid)))
	if err != nil {
		return nil, err
	}

	font := res[0]
	fontHandle := p.registerFont(&font)

	return &responses.FPDFText_LoadFont{
		Font: fontHandle.nativeRef,
	}, nil
}

// FPDFTextObj_GetFontSize returns the font size of a text object.
func (p *PdfiumImplementation) FPDFTextObj_GetFontSize(request *requests.FPDFTextObj_GetFontSize) (*responses.FPDFTextObj_GetFontSize, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	fontSizePointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}

	defer fontSizePointer.Free()
	res, err := p.Module.ExportedFunction("FPDFTextObj_GetFontSize").Call(p.Context, *pageObjectHandle.handle, fontSizePointer.Pointer)
	if err != nil {
		return nil, err
	}

	result := *(*int32)(unsafe.Pointer(&res[0]))
	if int(result) == 0 {
		return nil, errors.New("could not get font size")
	}

	fontSize, err := fontSizePointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFTextObj_GetFontSize{
		FontSize: float32(fontSize),
	}, nil
}

// FPDFFont_Close closes a loaded PDF font
func (p *PdfiumImplementation) FPDFFont_Close(request *requests.FPDFFont_Close) (*responses.FPDFFont_Close, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDFFont_Close").Call(p.Context, *fontHandle.handle)
	if err != nil {
		return nil, err
	}

	delete(p.fontRefs, fontHandle.nativeRef)

	return &responses.FPDFFont_Close{}, nil
}

// FPDFPageObj_CreateTextObj creates a new text object using a loaded font.
func (p *PdfiumImplementation) FPDFPageObj_CreateTextObj(request *requests.FPDFPageObj_CreateTextObj) (*responses.FPDFPageObj_CreateTextObj, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_CreateTextObj").Call(p.Context, *documentHandle.handle, *fontHandle.handle, *(*uint64)(unsafe.Pointer(&request.FontSize)))
	if err != nil {
		return nil, err
	}

	textObject := res[0]
	textObjectHandle := p.registerPageObject(&textObject)

	return &responses.FPDFPageObj_CreateTextObj{
		PageObject: textObjectHandle.nativeRef,
	}, nil
}

// FPDFTextObj_GetTextRenderMode returns the text rendering mode of a text object.
func (p *PdfiumImplementation) FPDFTextObj_GetTextRenderMode(request *requests.FPDFTextObj_GetTextRenderMode) (*responses.FPDFTextObj_GetTextRenderMode, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFTextObj_GetTextRenderMode").Call(p.Context, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	textRenderMode := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFTextObj_GetTextRenderMode{
		TextRenderMode: enums.FPDF_TEXT_RENDERMODE(textRenderMode),
	}, nil
}

// FPDFTextObj_GetText returns the text of a text object.
func (p *PdfiumImplementation) FPDFTextObj_GetText(request *requests.FPDFTextObj_GetText) (*responses.FPDFTextObj_GetText, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	textPageHandle, err := p.getTextPageHandle(request.TextPage)
	if err != nil {
		return nil, err
	}

	// First get the text value length.
	res, err := p.Module.ExportedFunction("FPDFTextObj_GetText").Call(p.Context, *pageObjectHandle.handle, *textPageHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	textSize := *(*int32)(unsafe.Pointer(&res[0]))
	if textSize == 0 {
		return nil, errors.New("could not get text")
	}

	charDataPointer, err := p.ByteArrayPointer(uint64(textSize), nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDFTextObj_GetText").Call(p.Context, *pageObjectHandle.handle, *textPageHandle.handle, charDataPointer.Pointer, *(*uint64)(unsafe.Pointer(&textSize)))
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

	return &responses.FPDFTextObj_GetText{
		Text: transformedName,
	}, nil
}

// FPDFFormObj_CountObjects returns the number of page objects inside the given form object.
func (p *PdfiumImplementation) FPDFFormObj_CountObjects(request *requests.FPDFFormObj_CountObjects) (*responses.FPDFFormObj_CountObjects, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFFormObj_CountObjects").Call(p.Context, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFFormObj_CountObjects{
		Count: int(count),
	}, nil
}

// FPDFFormObj_GetObject returns the page object in the given form object at the given index.
func (p *PdfiumImplementation) FPDFFormObj_GetObject(request *requests.FPDFFormObj_GetObject) (*responses.FPDFFormObj_GetObject, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFFormObj_GetObject").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	formObject := res[0]
	if formObject == 0 {
		return nil, errors.New("could not get form object")
	}

	formObjectHandle := p.registerPageObject(&formObject)

	return &responses.FPDFFormObj_GetObject{
		PageObject: formObjectHandle.nativeRef,
	}, nil
}

// FPDFFormObj_GetObject returns the page object in the given form object at the given index.
func (p *PdfiumImplementation) FPDFFormObj_RemoveObject(request *requests.FPDFFormObj_RemoveObject) (*responses.FPDFFormObj_RemoveObject, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	formObjectHandle, err := p.getPageObjectHandle(request.FormObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFFormObj_RemoveObject").Call(p.Context, *pageObjectHandle.handle, *formObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not remove form object")
	}

	return &responses.FPDFFormObj_RemoveObject{}, nil
}

// FPDFPage_RemoveObject removes an object from a page.
// Ownership is transferred to the caller. Call FPDFPageObj_Destroy() to free
// it.
// Experimental API.
func (p *PdfiumImplementation) FPDFPage_RemoveObject(request *requests.FPDFPage_RemoveObject) (*responses.FPDFPage_RemoveObject, error) {
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

	res, err := p.Module.ExportedFunction("FPDFPage_RemoveObject").Call(p.Context, *pageHandle.handle, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not remove object")
	}

	return &responses.FPDFPage_RemoveObject{}, nil
}

// FPDFPageObj_TransformF transforms the page object by the given matrix.
// The matrix is composed as:
//
//	|a c e|
//	|b d f|
//
// and can be used to scale, rotate, shear and translate the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_TransformF(request *requests.FPDFPageObj_TransformF) (*responses.FPDFPageObj_TransformF, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	matrixPointer, _, err := p.CStructFS_MATRIX(&request.Transform)
	if err != nil {
		return nil, err
	}
	defer p.Free(matrixPointer)

	res, err := p.Module.ExportedFunction("FPDFPageObj_TransformF").Call(p.Context, *pageObjectHandle.handle, matrixPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not transform object")
	}

	return &responses.FPDFPageObj_TransformF{}, nil
}

// FPDFPageObj_GetMatrix returns the transform matrix of a page object.
// The matrix is composed as:
//
//	|a c e|
//	|b d f|
//
// and can be used to scale, rotate, shear and translate the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetMatrix(request *requests.FPDFPageObj_GetMatrix) (*responses.FPDFPageObj_GetMatrix, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	matrixPointer, matrixValue, err := p.CStructFS_MATRIX(nil)
	if err != nil {
		return nil, err
	}
	defer p.Free(matrixPointer)

	res, err := p.Module.ExportedFunction("FPDFPageObj_GetMatrix").Call(p.Context, *pageObjectHandle.handle, matrixPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get page object matrix")
	}

	matrix, err := matrixValue()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObj_GetMatrix{
		Matrix: structs.FPDF_FS_MATRIX{
			A: float32(matrix.A),
			B: float32(matrix.B),
			C: float32(matrix.C),
			D: float32(matrix.D),
			E: float32(matrix.E),
			F: float32(matrix.F),
		},
	}, nil
}

// FPDFPageObj_SetMatrix sets the transform matrix on a page object.
// The matrix is composed as:
//
//	|a c e|
//	|b d f|
//
// and can be used to scale, rotate, shear and translate the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_SetMatrix(request *requests.FPDFPageObj_SetMatrix) (*responses.FPDFPageObj_SetMatrix, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	matrixPointer, _, err := p.CStructFS_MATRIX(&request.Transform)
	if err != nil {
		return nil, err
	}
	defer p.Free(matrixPointer)

	res, err := p.Module.ExportedFunction("FPDFPageObj_SetMatrix").Call(p.Context, *pageObjectHandle.handle, matrixPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set page object matrix")
	}

	return &responses.FPDFPageObj_SetMatrix{}, nil
}

// FPDFPageObj_GetMarkedContentID returns the marked content ID of a page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetMarkedContentID(request *requests.FPDFPageObj_GetMarkedContentID) (*responses.FPDFPageObj_GetMarkedContentID, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_GetMarkedContentID").Call(p.Context, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	markedContentID := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFPageObj_GetMarkedContentID{
		MarkedContentID: int(markedContentID),
	}, nil
}

// FPDFPageObj_CountMarks returns the count of content marks in a page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_CountMarks(request *requests.FPDFPageObj_CountMarks) (*responses.FPDFPageObj_CountMarks, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_CountMarks").Call(p.Context, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFPageObj_CountMarks{
		Count: int(count),
	}, nil
}

// FPDFPageObj_GetMark returns the content mark of a page object at the given index.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetMark(request *requests.FPDFPageObj_GetMark) (*responses.FPDFPageObj_GetMark, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_GetMark").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	mark := res[0]
	if mark == 0 {
		return nil, errors.New("could not get mark")
	}

	markHandle := p.registerPageObjectMark(&mark)

	return &responses.FPDFPageObj_GetMark{
		Mark: markHandle.nativeRef,
	}, nil
}

// FPDFPageObj_AddMark adds a new content mark to the given page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_AddMark(request *requests.FPDFPageObj_AddMark) (*responses.FPDFPageObj_AddMark, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	name, err := p.CString(request.Name)
	if err != nil {
		return nil, err
	}

	defer name.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObj_AddMark").Call(p.Context, *pageObjectHandle.handle, name.Pointer)
	if err != nil {
		return nil, err
	}

	mark := res[0]
	if mark == 0 {
		return nil, errors.New("could not add mark")
	}

	markHandle := p.registerPageObjectMark(&mark)

	return &responses.FPDFPageObj_AddMark{
		Mark: markHandle.nativeRef,
	}, nil
}

// FPDFPageObj_RemoveMark removes the given content mark from the given page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_RemoveMark(request *requests.FPDFPageObj_RemoveMark) (*responses.FPDFPageObj_RemoveMark, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_RemoveMark").Call(p.Context, *pageObjectHandle.handle, *pageObjectMarkHandle.handle)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not remove mark")
	}

	return &responses.FPDFPageObj_RemoveMark{}, nil
}

// FPDFPageObjMark_GetName returns the name of a content mark.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetName(request *requests.FPDFPageObjMark_GetName) (*responses.FPDFPageObjMark_GetName, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	nameLengthPointer, err := p.ULongPointer()
	if err != nil {
		return nil, err
	}
	defer nameLengthPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObjMark_GetName").Call(p.Context, *pageObjectMarkHandle.handle, 0, 0, nameLengthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get name")
	}

	nameLengthValue, err := nameLengthPointer.Value()
	if err != nil {
		return nil, err
	}

	nameLength := uint64(nameLengthValue)
	if nameLength == 0 {
		return nil, errors.New("could not get name")
	}

	charDataPointer, err := p.ByteArrayPointer(nameLength, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFPageObjMark_GetName").Call(p.Context, *pageObjectMarkHandle.handle, charDataPointer.Pointer, nameLength, nameLengthPointer.Pointer)
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

	return &responses.FPDFPageObjMark_GetName{
		Name: transformedText,
	}, nil
}

// FPDFPageObjMark_CountParams returns the number of key/value pair parameters in the given mark.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_CountParams(request *requests.FPDFPageObjMark_CountParams) (*responses.FPDFPageObjMark_CountParams, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObjMark_CountParams").Call(p.Context, *pageObjectMarkHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFPageObjMark_CountParams{
		Count: int(count),
	}, nil
}

// FPDFPageObjMark_GetParamKey returns the key of a property in a content mark.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetParamKey(request *requests.FPDFPageObjMark_GetParamKey) (*responses.FPDFPageObjMark_GetParamKey, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	keyLengthPointer, err := p.ULongPointer()
	if err != nil {
		return nil, err
	}
	defer keyLengthPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObjMark_GetParamKey").Call(p.Context, *pageObjectMarkHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)), 0, 0, keyLengthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get key")
	}

	keyLengthValue, err := keyLengthPointer.Value()
	if err != nil {
		return nil, err
	}

	keyLength := uint64(keyLengthValue)
	if keyLength == 0 {
		return nil, errors.New("could not get key")
	}

	charDataPointer, err := p.ByteArrayPointer(keyLength, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFPageObjMark_GetParamKey").Call(p.Context, *pageObjectMarkHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)), charDataPointer.Pointer, keyLength, keyLengthPointer.Pointer)
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

	return &responses.FPDFPageObjMark_GetParamKey{
		Key: transformedText,
	}, nil
}

// FPDFPageObjMark_GetParamValueType returns the type of the value of a property in a content mark by key.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetParamValueType(request *requests.FPDFPageObjMark_GetParamValueType) (*responses.FPDFPageObjMark_GetParamValueType, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObjMark_GetParamValueType").Call(p.Context, *pageObjectMarkHandle.handle, keyPointer.Pointer)
	if err != nil {
		return nil, err
	}

	valueType := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFPageObjMark_GetParamValueType{
		ValueType: enums.FPDF_OBJECT_TYPE(valueType),
	}, nil
}

// FPDFPageObjMark_GetParamIntValue returns the value of a number property in a content mark by key as int.
// FPDFPageObjMark_GetParamValueType() should have returned FPDF_OBJECT_NUMBER
// for this property.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetParamIntValue(request *requests.FPDFPageObjMark_GetParamIntValue) (*responses.FPDFPageObjMark_GetParamIntValue, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	intValuePointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer intValuePointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObjMark_GetParamIntValue").Call(p.Context, *pageObjectMarkHandle.handle, keyPointer.Pointer, intValuePointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get value")
	}

	intValue, err := intValuePointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObjMark_GetParamIntValue{
		Value: int(intValue),
	}, nil
}

// FPDFPageObjMark_GetParamStringValue returns the value of a string property in a content mark by key.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetParamStringValue(request *requests.FPDFPageObjMark_GetParamStringValue) (*responses.FPDFPageObjMark_GetParamStringValue, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	valueLengthPointer, err := p.ULongPointer()
	if err != nil {
		return nil, err
	}
	defer valueLengthPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObjMark_GetParamStringValue").Call(p.Context, *pageObjectMarkHandle.handle, keyPointer.Pointer, 0, 0, valueLengthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get value")
	}

	valueLengthValue, err := valueLengthPointer.Value()
	if err != nil {
		return nil, err
	}

	valueLength := uint64(valueLengthValue)
	if valueLength == 0 {
		return nil, errors.New("could not get value")
	}

	charDataPointer, err := p.ByteArrayPointer(valueLength, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFPageObjMark_GetParamStringValue").Call(p.Context, *pageObjectMarkHandle.handle, keyPointer.Pointer, charDataPointer.Pointer, valueLength, valueLengthPointer.Pointer)
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

	return &responses.FPDFPageObjMark_GetParamStringValue{
		Value: transformedText,
	}, nil
}

// FPDFPageObjMark_GetParamBlobValue returns the value of a blob property in a content mark by key.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_GetParamBlobValue(request *requests.FPDFPageObjMark_GetParamBlobValue) (*responses.FPDFPageObjMark_GetParamBlobValue, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	valueLengthPointer, err := p.ULongPointer()
	if err != nil {
		return nil, err
	}
	defer valueLengthPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObjMark_GetParamBlobValue").Call(p.Context, *pageObjectMarkHandle.handle, keyPointer.Pointer, 0, 0, valueLengthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get value")
	}

	valueLengthValue, err := valueLengthPointer.Value()
	if err != nil {
		return nil, err
	}

	valueLength := uint64(valueLengthValue)
	if valueLength == 0 {
		return nil, errors.New("could not get value")
	}

	paramDataPointer, err := p.ByteArrayPointer(valueLength, nil)
	if err != nil {
		return nil, err
	}
	defer paramDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFPageObjMark_GetParamBlobValue").Call(p.Context, *pageObjectMarkHandle.handle, keyPointer.Pointer, paramDataPointer.Pointer, valueLength, valueLengthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	paramData, err := paramDataPointer.Value(true)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObjMark_GetParamBlobValue{
		Value: paramData,
	}, nil
}

// FPDFPageObjMark_SetIntParam sets the value of an int property in a content mark by key. If a parameter
// with the given key exists, its value is set to the given value. Otherwise, it is added as
// a new parameter.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_SetIntParam(request *requests.FPDFPageObjMark_SetIntParam) (*responses.FPDFPageObjMark_SetIntParam, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObjMark_SetIntParam").Call(p.Context, *documentHandle.handle, *pageObjectHandle.handle, *pageObjectMarkHandle.handle, keyPointer.Pointer, *(*uint64)(unsafe.Pointer(&request.Value)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set value")
	}

	return &responses.FPDFPageObjMark_SetIntParam{}, nil
}

// FPDFPageObjMark_SetStringParam sets the value of a string property in a content mark by key. If a parameter
// with the given key exists, its value is set to the given value. Otherwise, it is added as
// a new parameter.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_SetStringParam(request *requests.FPDFPageObjMark_SetStringParam) (*responses.FPDFPageObjMark_SetStringParam, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	valuePointer, err := p.CString(request.Value)
	if err != nil {
		return nil, err
	}
	defer valuePointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObjMark_SetStringParam").Call(p.Context, *documentHandle.handle, *pageObjectHandle.handle, *pageObjectMarkHandle.handle, keyPointer.Pointer, valuePointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set value")
	}

	return &responses.FPDFPageObjMark_SetStringParam{}, nil
}

// FPDFPageObjMark_SetBlobParam sets the value of a blob property in a content mark by key. If a parameter
// with the given key exists, its value is set to the given value. Otherwise, it is added as
// a new parameter.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_SetBlobParam(request *requests.FPDFPageObjMark_SetBlobParam) (*responses.FPDFPageObjMark_SetBlobParam, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	if request.Value == nil || len(request.Value) == 0 {
		return nil, errors.New("blob value cant be empty")
	}

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	dataLength := uint64(len(request.Value))
	valuePointer, err := p.ByteArrayPointer(dataLength, request.Value)
	if err != nil {
		return nil, err
	}
	defer valuePointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObjMark_SetBlobParam").Call(p.Context, *documentHandle.handle, *pageObjectHandle.handle, *pageObjectMarkHandle.handle, keyPointer.Pointer, valuePointer.Pointer, dataLength)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set value")
	}

	return &responses.FPDFPageObjMark_SetBlobParam{}, nil
}

// FPDFPageObjMark_RemoveParam removes a property from a content mark by key.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObjMark_RemoveParam(request *requests.FPDFPageObjMark_RemoveParam) (*responses.FPDFPageObjMark_RemoveParam, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	pageObjectMarkHandle, err := p.getPageObjectMarkHandle(request.PageObjectMark)
	if err != nil {
		return nil, err
	}

	keyPointer, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}
	defer keyPointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObjMark_RemoveParam").Call(p.Context, *pageObjectHandle.handle, *pageObjectMarkHandle.handle, keyPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set value")
	}

	return &responses.FPDFPageObjMark_RemoveParam{}, nil
}

// FPDFImageObj_GetRenderedBitmap returns a bitmap rasterization of the given image object that takes the image mask and
// image matrix into account. To render correctly, the caller must provide the
// document associated with the image object. If there is a page associated
// with the image object the caller should provide that as well.
// The returned bitmap will be owned by the caller, and FPDFBitmap_Destroy()
// must be called on the returned bitmap when it is no longer needed.
// Experimental API.
func (p *PdfiumImplementation) FPDFImageObj_GetRenderedBitmap(request *requests.FPDFImageObj_GetRenderedBitmap) (*responses.FPDFImageObj_GetRenderedBitmap, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	imageObjectHandle, err := p.getPageObjectHandle(request.ImageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFImageObj_GetRenderedBitmap").Call(p.Context, *documentHandle.handle, *pageHandle.handle, *imageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	bitmap := res[0]
	if bitmap == 0 {
		return nil, errors.New("could not get bitmap")
	}

	bitmapHandle := p.registerBitmap(&bitmap)

	return &responses.FPDFImageObj_GetRenderedBitmap{
		Bitmap: bitmapHandle.nativeRef,
	}, nil
}

// FPDFPageObj_GetRotatedBounds Get the quad points that bounds the page object.
// Similar to FPDFPageObj_GetBounds(), this returns the bounds of a page
// object. When the object is rotated by a non-multiple of 90 degrees, this API
// returns a tighter bound that cannot be represented with just the 4 sides of
// a rectangle.
//
// Currently only works the following page object types: FPDF_PAGEOBJ_TEXT and
// FPDF_PAGEOBJ_IMAGE.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetRotatedBounds(request *requests.FPDFPageObj_GetRotatedBounds) (*responses.FPDFPageObj_GetRotatedBounds, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	quadPointsPointer, quadPointsValue, err := p.CStructFS_QUADPOINTSF(nil)
	if err != nil {
		return nil, err
	}
	defer p.Free(quadPointsPointer)

	res, err := p.Module.ExportedFunction("FPDFPageObj_GetRotatedBounds").Call(p.Context, *pageObjectHandle.handle, quadPointsPointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get rotated bounds for page object")
	}

	quadPoints, err := quadPointsValue()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObj_GetRotatedBounds{
		QuadPoints: *quadPoints,
	}, nil
}

// FPDFPageObj_GetDashPhase returns the line dash phase of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetDashPhase(request *requests.FPDFPageObj_GetDashPhase) (*responses.FPDFPageObj_GetDashPhase, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	dashPhasePointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}

	defer dashPhasePointer.Free()

	res, err := p.Module.ExportedFunction("FPDFPageObj_GetDashPhase").Call(p.Context, *pageObjectHandle.handle, dashPhasePointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not get dash phase")
	}

	dashPhase, err := dashPhasePointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObj_GetDashPhase{
		DashPhase: float32(dashPhase),
	}, nil
}

// FPDFPageObj_SetDashPhase sets the line dash phase of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_SetDashPhase(request *requests.FPDFPageObj_SetDashPhase) (*responses.FPDFPageObj_SetDashPhase, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_SetDashPhase").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.DashPhase)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not set dash phase")
	}

	return &responses.FPDFPageObj_SetDashPhase{}, nil
}

// FPDFPageObj_GetDashCount returns the line dash array size of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetDashCount(request *requests.FPDFPageObj_GetDashCount) (*responses.FPDFPageObj_GetDashCount, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPageObj_GetDashCount").Call(p.Context, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	dashCount := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFPageObj_GetDashCount{
		DashCount: int(dashCount),
	}, nil
}

// FPDFPageObj_GetDashArray returns the line dash array of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_GetDashArray(request *requests.FPDFPageObj_GetDashArray) (*responses.FPDFPageObj_GetDashArray, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	// First get the Dash Count.
	res, err := p.Module.ExportedFunction("FPDFPageObj_GetDashCount").Call(p.Context, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	dashCount := *(*int32)(unsafe.Pointer(&res[0]))
	dashCountSize := uint64(dashCount)

	convertedData := make([]float32, 0)
	if int(dashCount) > 0 {
		// Create an array that's big enough.
		valueDataPointer, err := p.FloatArrayPointer(dashCountSize)
		if err != nil {
			return nil, err
		}
		defer valueDataPointer.Free()

		_, err = p.Module.ExportedFunction("FPDFPageObj_GetDashArray").Call(p.Context, *pageObjectHandle.handle, valueDataPointer.Pointer, dashCountSize)
		if err != nil {
			return nil, err
		}

		valueData, err := valueDataPointer.Value()
		if err != nil {
			return nil, err
		}

		convertedData = valueData
	}

	return &responses.FPDFPageObj_GetDashArray{
		DashArray: convertedData,
	}, nil
}

// FPDFPageObj_SetDashArray sets the line dash array of the page object.
// Experimental API.
func (p *PdfiumImplementation) FPDFPageObj_SetDashArray(request *requests.FPDFPageObj_SetDashArray) (*responses.FPDFPageObj_SetDashArray, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	dashCountSize := uint64(uint64(len(request.DashArray)))

	// Create an array that's big enough.
	valueDataPointer, err := p.FloatArrayPointer(dashCountSize)
	if err != nil {
		return nil, err
	}
	defer valueDataPointer.Free()

	// Put the values in the array.
	for i := range request.DashArray {
		p.Module.Memory().WriteFloat32Le(uint32(valueDataPointer.Pointer+(p.CSizeFloat()*uint64(i))), request.DashArray[i])
	}

	_, err = p.Module.ExportedFunction("FPDFPageObj_SetDashArray").Call(p.Context, *pageObjectHandle.handle, valueDataPointer.Pointer, dashCountSize, *(*uint64)(unsafe.Pointer(&request.DashPhase)))
	if err != nil {
		return nil, err
	}

	return &responses.FPDFPageObj_SetDashArray{}, nil
}

// FPDFText_LoadStandardFont loads one of the standard 14 fonts per PDF spec 1.7 page 416. The preferred
// way of using font style is using a dash to separate the name from the style,
// for example 'Helvetica-BoldItalic'.
// The loaded font can be closed using FPDFFont_Close.
// Experimental API.
func (p *PdfiumImplementation) FPDFText_LoadStandardFont(request *requests.FPDFText_LoadStandardFont) (*responses.FPDFText_LoadStandardFont, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	fontNamePointer, err := p.CString(request.Font)
	if err != nil {
		return nil, err
	}
	defer fontNamePointer.Free()

	res, err := p.Module.ExportedFunction("FPDFText_LoadStandardFont").Call(p.Context, *documentHandle.handle, fontNamePointer.Pointer)
	if err != nil {
		return nil, err
	}

	font := res[0]
	if font == 0 {
		return nil, errors.New("could not load standard font")
	}

	fontHandle := p.registerFont(&font)

	return &responses.FPDFText_LoadStandardFont{
		Font: fontHandle.nativeRef,
	}, nil
}

// FPDFText_LoadCidType2Font returns a font object loaded from a stream of data for a type 2 CID font.
// The font is loaded into the document. Unlike FPDFText_LoadFont(), the ToUnicode data and the CIDToGIDMap
// data are caller provided, instead of auto-generated.
// The loaded font can be closed using FPDFFont_Close().
// Experimental API.
func (p *PdfiumImplementation) FPDFText_LoadCidType2Font(request *requests.FPDFText_LoadCidType2Font) (*responses.FPDFText_LoadCidType2Font, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	fontDataLength := uint32(len(request.FontData))
	fontData, err := p.ByteArrayPointer(uint64(len(request.FontData)), request.FontData)
	if err != nil {
		return nil, err
	}

	defer fontData.Free()

	toUnicodeCmapPointer, err := p.CString(request.ToUnicodeCmap)
	if err != nil {
		return nil, err
	}

	defer toUnicodeCmapPointer.Free()

	cidToGIDMapDataLength := uint32(len(request.CIDToGIDMapData))
	cidToGIDMapData, err := p.ByteArrayPointer(uint64(len(request.CIDToGIDMapData)), request.CIDToGIDMapData)
	if err != nil {
		return nil, err
	}

	defer cidToGIDMapData.Free()

	res, err := p.Module.ExportedFunction("FPDFText_LoadCidType2Font").Call(p.Context, *documentHandle.handle, fontData.Pointer, *(*uint64)(unsafe.Pointer(&fontDataLength)), toUnicodeCmapPointer.Pointer, cidToGIDMapData.Pointer, *(*uint64)(unsafe.Pointer(&cidToGIDMapDataLength)))
	if err != nil {
		return nil, err
	}

	font := res[0]
	if font == 0 {
		return nil, errors.New("could not load CID Type2 font")
	}

	fontHandle := p.registerFont(&font)

	return &responses.FPDFText_LoadCidType2Font{
		Font: fontHandle.nativeRef,
	}, nil
}

// FPDFTextObj_SetTextRenderMode sets the text rendering mode of a text object.
// Experimental API.
func (p *PdfiumImplementation) FPDFTextObj_SetTextRenderMode(request *requests.FPDFTextObj_SetTextRenderMode) (*responses.FPDFTextObj_SetTextRenderMode, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFTextObj_SetTextRenderMode").Call(p.Context, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.TextRenderMode)))
	if err != nil {
		return nil, err
	}

	result := *(*int32)(unsafe.Pointer(&res[0]))
	if int(result) == 0 {
		return nil, errors.New("could not set text render mode")
	}

	return &responses.FPDFTextObj_SetTextRenderMode{}, nil
}

// FPDFTextObj_GetRenderedBitmap returns a bitmap rasterization of the given text object.
// To render correctly, the caller must provide the document associated with the text object.
// If there is a page associated with text object, the caller should provide that as well.
// The returned bitmap will be owned by the caller, and FPDFBitmap_Destroy()
// must be called on the returned bitmap when it is no longer needed.
// Experimental API.
func (p *PdfiumImplementation) FPDFTextObj_GetRenderedBitmap(request *requests.FPDFTextObj_GetRenderedBitmap) (*responses.FPDFTextObj_GetRenderedBitmap, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	var pageHandle uint64
	if request.Page != "" {
		pageHandleReference, err := p.getPageHandle(request.Page)
		if err != nil {
			return nil, err
		}

		pageHandle = *pageHandleReference.handle
	}

	res, err := p.Module.ExportedFunction("FPDFTextObj_GetRenderedBitmap").Call(p.Context, *documentHandle.handle, pageHandle, *pageObjectHandle.handle, *(*uint64)(unsafe.Pointer(&request.Scale)))
	if err != nil {
		return nil, err
	}

	bitmap := res[0]
	if bitmap == 0 {
		return nil, errors.New("could not render text object as bitmap")
	}

	bitmapHandle := p.registerBitmap(&bitmap)

	return &responses.FPDFTextObj_GetRenderedBitmap{
		Bitmap: bitmapHandle.nativeRef,
	}, nil
}

// FPDFTextObj_GetFont returns the font of a text object.
// Experimental API.
func (p *PdfiumImplementation) FPDFTextObj_GetFont(request *requests.FPDFTextObj_GetFont) (*responses.FPDFTextObj_GetFont, error) {
	p.Lock()
	defer p.Unlock()

	pageObjectHandle, err := p.getPageObjectHandle(request.PageObject)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFTextObj_GetFont").Call(p.Context, *pageObjectHandle.handle)
	if err != nil {
		return nil, err
	}

	font := res[0]
	if font == 0 {
		return nil, errors.New("could not load standard font")
	}

	fontHandle := p.registerFont(&font)

	return &responses.FPDFTextObj_GetFont{
		Font: fontHandle.nativeRef,
	}, nil
}

// FPDFFont_GetBaseFontName returns the base font name of a font.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetBaseFontName(request *requests.FPDFFont_GetBaseFontName) (*responses.FPDFFont_GetBaseFontName, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	// First get the text value length.
	res, err := p.Module.ExportedFunction("FPDFFont_GetBaseFontName").Call(p.Context, *fontHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	nameSize := *(*int32)(unsafe.Pointer(&res[0]))
	if nameSize == 0 {
		return nil, errors.New("could not get font name")
	}

	charDataSize := uint64(nameSize)
	charDataPointer, err := p.ByteArrayPointer(charDataSize, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFFont_GetBaseFontName").Call(p.Context, *fontHandle.handle, charDataPointer.Pointer, charDataSize)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(true)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFFont_GetBaseFontName{
		BaseFontName: string(charData[:len(charData)-1]), // Remove NULL-terminator
	}, nil
}

// FPDFFont_GetFamilyName returns the family name of a font.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetFamilyName(request *requests.FPDFFont_GetFamilyName) (*responses.FPDFFont_GetFamilyName, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	// First get the text value length.
	res, err := p.Module.ExportedFunction("FPDFFont_GetFamilyName").Call(p.Context, *fontHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	nameSize := *(*int32)(unsafe.Pointer(&res[0]))
	if nameSize == 0 {
		return nil, errors.New("could not get font name")
	}

	charDataSize := uint64(nameSize)
	charDataPointer, err := p.ByteArrayPointer(charDataSize, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FPDFFont_GetFamilyName").Call(p.Context, *fontHandle.handle, charDataPointer.Pointer, charDataSize)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value(true)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFFont_GetFamilyName{
		FamilyName: string(charData[:len(charData)-1]), // Remove NULL-terminator
	}, nil
}

// FPDFFont_GetFontData returns the decoded data from the given font.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetFontData(request *requests.FPDFFont_GetFontData) (*responses.FPDFFont_GetFontData, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	outBufLenPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer outBufLenPointer.Free()

	// First get the font data length.
	res, err := p.Module.ExportedFunction("FPDFFont_GetFontData").Call(p.Context, *fontHandle.handle, 0, 0, outBufLenPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) != 1 {
		return nil, errors.New("could not get font data")
	}

	outBufLenValue, err := outBufLenPointer.Value()
	if err != nil {
		return nil, err
	}

	outBufLen := uint64(outBufLenValue)
	if int(outBufLen) == 0 {
		return nil, errors.New("could not get font data")
	}

	fontDataPointer, err := p.ByteArrayPointer(outBufLen, nil)
	if err != nil {
		return nil, err
	}
	defer fontDataPointer.Free()

	res, err = p.Module.ExportedFunction("FPDFFont_GetFontData").Call(p.Context, *fontHandle.handle, fontDataPointer.Pointer, outBufLen, outBufLenPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success = *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) != 1 {
		return nil, errors.New("could not get font data")
	}

	fontData, err := fontDataPointer.Value(true)
	if err != nil {
		return nil, err
	}

	return &responses.FPDFFont_GetFontData{
		FontData: fontData,
	}, nil
}

// FPDFFont_GetIsEmbedded returns whether the given font is embedded or not.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetIsEmbedded(request *requests.FPDFFont_GetIsEmbedded) (*responses.FPDFFont_GetIsEmbedded, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFFont_GetIsEmbedded").Call(p.Context, *fontHandle.handle)
	if err != nil {
		return nil, err
	}

	isEmbedded := *(*int32)(unsafe.Pointer(&res[0]))
	return &responses.FPDFFont_GetIsEmbedded{
		IsEmbedded: int(isEmbedded) == 1,
	}, nil
}

// FPDFFont_GetFlags returns the descriptor flags of a font.
// Returns the bit flags specifying various characteristics of the font as
// defined in ISO 32000-1:2008, table 123.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetFlags(request *requests.FPDFFont_GetFlags) (*responses.FPDFFont_GetFlags, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFFont_GetFlags").Call(p.Context, *fontHandle.handle)
	if err != nil {
		return nil, err
	}

	flags := *(*int32)(unsafe.Pointer(&res[0]))
	if int(flags) == -1 {
		return nil, errors.New("could not get font flags")
	}

	fontFlags := &responses.FPDFFont_GetFlags{
		Flags: uint32(flags),
	}

	FixedPitch := uint32(1 << 1)
	Serif := uint32(1 << 2)
	Symbolic := uint32(1 << 3)
	Script := uint32(1 << 4)
	Nonsymbolic := uint32(1 << 6)
	Italic := uint32(1 << 7)
	AllCap := uint32(1 << 17)
	SmallCap := uint32(1 << 18)
	ForceBold := uint32(1 << 19)

	hasFlag := func(flag uint32) bool {
		if fontFlags.Flags&flag > 0 {
			return true
		}

		return false
	}

	fontFlags.FixedPitch = hasFlag(FixedPitch)
	fontFlags.Serif = hasFlag(Serif)
	fontFlags.Symbolic = hasFlag(Symbolic)
	fontFlags.Script = hasFlag(Script)
	fontFlags.Nonsymbolic = hasFlag(Nonsymbolic)
	fontFlags.Italic = hasFlag(Italic)
	fontFlags.AllCap = hasFlag(AllCap)
	fontFlags.SmallCap = hasFlag(SmallCap)
	fontFlags.ForceBold = hasFlag(ForceBold)

	return fontFlags, nil
}

// FPDFFont_GetWeight returns the font weight of a font.
// Typical values are 400 (normal) and 700 (bold).
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetWeight(request *requests.FPDFFont_GetWeight) (*responses.FPDFFont_GetWeight, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFFont_GetWeight").Call(p.Context, *fontHandle.handle)
	if err != nil {
		return nil, err
	}

	fontWeight := *(*int32)(unsafe.Pointer(&res[0]))
	if int(fontWeight) == -1 {
		return nil, errors.New("could not get font weight")
	}

	return &responses.FPDFFont_GetWeight{
		Weight: int(fontWeight),
	}, nil
}

// FPDFFont_GetItalicAngle returns the italic angle of a font.
// The italic angle of a font is defined as degrees counterclockwise
// from vertical. For a font that slopes to the right, this will be negative.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetItalicAngle(request *requests.FPDFFont_GetItalicAngle) (*responses.FPDFFont_GetItalicAngle, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	anglePointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFFont_GetItalicAngle").Call(p.Context, *fontHandle.handle, anglePointer.Pointer)
	if err != nil {
		return nil, err
	}

	result := *(*int32)(unsafe.Pointer(&res[0]))
	if int(result) == 0 {
		return nil, errors.New("could not get italic angle")
	}

	angle, err := anglePointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFFont_GetItalicAngle{
		ItalicAngle: int(angle),
	}, nil
}

// FPDFFont_GetAscent returns ascent distance of a font.
// Ascent is the maximum distance in points above the baseline reached by the
// glyphs of the font. One point is 1/72 inch (around 0.3528 mm).
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetAscent(request *requests.FPDFFont_GetAscent) (*responses.FPDFFont_GetAscent, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	ascentPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFFont_GetAscent").Call(p.Context, *fontHandle.handle, *(*uint64)(unsafe.Pointer(&request.FontSize)), ascentPointer.Pointer)
	if err != nil {
		return nil, err
	}

	result := *(*int32)(unsafe.Pointer(&res[0]))
	if int(result) == 0 {
		return nil, errors.New("could not get ascent")
	}

	ascent, err := ascentPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFFont_GetAscent{
		Ascent: float32(ascent),
	}, nil
}

// FPDFFont_GetDescent returns the descent distance of a font.
// Descent is the maximum distance in points below the baseline reached by the
// glyphs of the font. One point is 1/72 inch (around 0.3528 mm).
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetDescent(request *requests.FPDFFont_GetDescent) (*responses.FPDFFont_GetDescent, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	descentPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFFont_GetDescent").Call(p.Context, *fontHandle.handle, *(*uint64)(unsafe.Pointer(&request.FontSize)), descentPointer.Pointer)
	if err != nil {
		return nil, err
	}

	result := *(*int32)(unsafe.Pointer(&res[0]))
	if int(result) == 0 {
		return nil, errors.New("could not get descent")
	}

	descent, err := descentPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFFont_GetDescent{
		Descent: float32(descent),
	}, nil
}

// FPDFFont_GetGlyphWidth returns the width of a glyph in a font.
// Glyph width is the distance from the end of the prior glyph to the next
// glyph. This will be the vertical distance for vertical writing.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetGlyphWidth(request *requests.FPDFFont_GetGlyphWidth) (*responses.FPDFFont_GetGlyphWidth, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	glyphWidthPointer, err := p.FloatPointer(nil)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFFont_GetGlyphWidth").Call(p.Context, *fontHandle.handle, *(*uint64)(unsafe.Pointer(&request.Glyph)), *(*uint64)(unsafe.Pointer(&request.FontSize)), glyphWidthPointer.Pointer)
	if err != nil {
		return nil, err
	}

	result := *(*int32)(unsafe.Pointer(&res[0]))
	if int(result) == 0 {
		return nil, errors.New("could not get glyph width")
	}

	glyphWidth, err := glyphWidthPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDFFont_GetGlyphWidth{
		GlyphWidth: float32(glyphWidth),
	}, nil
}

// FPDFFont_GetGlyphPath returns the glyphpath describing how to draw a font glyph.
// Experimental API.
func (p *PdfiumImplementation) FPDFFont_GetGlyphPath(request *requests.FPDFFont_GetGlyphPath) (*responses.FPDFFont_GetGlyphPath, error) {
	p.Lock()
	defer p.Unlock()

	fontHandle, err := p.getFontHandle(request.Font)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFFont_GetGlyphPath").Call(p.Context, *fontHandle.handle, *(*uint64)(unsafe.Pointer(&request.Glyph)), *(*uint64)(unsafe.Pointer(&request.FontSize)))
	if err != nil {
		return nil, err
	}

	glyphPath := res[0]
	if glyphPath == 0 {
		return nil, errors.New("could not get glyph path")
	}

	glyphPathHandle := p.registerGlyphPath(&glyphPath)

	return &responses.FPDFFont_GetGlyphPath{
		GlyphPath: glyphPathHandle.nativeRef,
	}, nil
}

// FPDFGlyphPath_CountGlyphSegments returns the number of segments inside the given glyphpath.
// Experimental API.
func (p *PdfiumImplementation) FPDFGlyphPath_CountGlyphSegments(request *requests.FPDFGlyphPath_CountGlyphSegments) (*responses.FPDFGlyphPath_CountGlyphSegments, error) {
	p.Lock()
	defer p.Unlock()

	glyphPathHandle, err := p.getGlyphPathHandle(request.GlyphPath)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFGlyphPath_CountGlyphSegments").Call(p.Context, *glyphPathHandle.handle)
	if err != nil {
		return nil, err
	}

	count := *(*int32)(unsafe.Pointer(&res[0]))
	if int(count) == -1 {
		return nil, errors.New("could not get glyph path segment count")
	}

	return &responses.FPDFGlyphPath_CountGlyphSegments{
		Count: int(count),
	}, nil
}

// FPDFGlyphPath_GetGlyphPathSegment returns the segment in glyphpath at the given index.
// Experimental API.
func (p *PdfiumImplementation) FPDFGlyphPath_GetGlyphPathSegment(request *requests.FPDFGlyphPath_GetGlyphPathSegment) (*responses.FPDFGlyphPath_GetGlyphPathSegment, error) {
	p.Lock()
	defer p.Unlock()

	glyphPathHandle, err := p.getGlyphPathHandle(request.GlyphPath)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFGlyphPath_GetGlyphPathSegment").Call(p.Context, *glyphPathHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	segment := res[0]
	if segment == 0 {
		return nil, errors.New("could not get glyph path segment")
	}

	segmentHandle := p.registerPathSegment(&segment)

	return &responses.FPDFGlyphPath_GetGlyphPathSegment{
		GlyphPathSegment: segmentHandle.nativeRef,
	}, nil
}
