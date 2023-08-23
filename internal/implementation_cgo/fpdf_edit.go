package implementation_cgo

/*
#cgo pkg-config: pdfium
#include "fpdf_edit.h"
#include <stdlib.h>

extern int go_read_seeker_cb(void *param, unsigned long position, unsigned char *pBuf, unsigned long size);

static inline void FPDF_FILEACCESS_SET_GET_BLOCK(FPDF_FILEACCESS *fs, char *id) {
	fs->m_GetBlock = &go_read_seeker_cb;
	fs->m_Param = id;
}

*/
import "C"
import (
	"errors"
	"io"
	"os"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"

	"github.com/google/uuid"
)

// FPDF_CreateNewDocument returns a new document.
func (p *PdfiumImplementation) FPDF_CreateNewDocument(request *requests.FPDF_CreateNewDocument) (*responses.FPDF_CreateNewDocument, error) {
	p.Lock()
	defer p.Unlock()

	doc := C.FPDF_CreateNewDocument()
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

	page := C.FPDFPage_New(documentHandle.handle, C.int(request.PageIndex), C.double(request.Width), C.double(request.Height))
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

	C.FPDFPage_Delete(documentHandle.handle, C.int(request.PageIndex))

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

	rotation := C.FPDFPage_GetRotation(pageHandle.handle)

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

	C.FPDFPage_SetRotation(pageHandle.handle, C.int(request.Rotate))

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

	C.FPDFPage_InsertObject(pageHandle.handle, pageObjectHandle.handle)

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

	count := C.FPDFPage_CountObjects(pageHandle.handle)

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

	pageObject := C.FPDFPage_GetObject(pageHandle.handle, C.int(request.Index))
	if pageObject == nil {
		return nil, errors.New("could not get object")
	}

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

	alpha := C.FPDFPage_HasTransparency(pageHandle.handle)

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

	success := C.FPDFPage_GenerateContent(pageHandle.handle)
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

	C.FPDFPageObj_Destroy(pageObjectHandle.handle)

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

	hasTransparency := C.FPDFPageObj_HasTransparency(pageObjectHandle.handle)

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

	pageObjectType := C.FPDFPageObj_GetType(pageObjectHandle.handle)

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

	C.FPDFPageObj_Transform(pageObjectHandle.handle, C.double(request.Transform.A), C.double(request.Transform.B), C.double(request.Transform.C), C.double(request.Transform.D), C.double(request.Transform.E), C.double(request.Transform.F))

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

	C.FPDFPage_TransformAnnots(pageHandle.handle, C.double(request.Transform.A), C.double(request.Transform.B), C.double(request.Transform.C), C.double(request.Transform.D), C.double(request.Transform.E), C.double(request.Transform.F))

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

	imageObject := C.FPDFPageObj_NewImageObj(documentHandle.handle)
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

	var pageHandle C.FPDF_PAGE
	if request.Page != nil {
		loadedPage, err := p.loadPage(*request.Page)
		if err != nil {
			return nil, err
		}

		pageHandle = loadedPage.handle
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

	// Create a PDFium file access struct.
	readerStruct := C.FPDF_FILEACCESS{}
	readerStruct.m_FileLen = C.ulong(request.FileReaderSize)

	readerRef := uuid.New()
	readerRefString := readerRef.String()
	cReaderRef := C.CString(readerRefString)

	// Set the Go callback through cgo.
	C.FPDF_FILEACCESS_SET_GET_BLOCK(&readerStruct, cReaderRef)

	fileReaderRef := &fileReaderRef{
		stringRef:  unsafe.Pointer(cReaderRef),
		reader:     fileReader,
		fileAccess: &readerStruct,
	}

	Pdfium.fileReaders[readerRef.String()] = fileReaderRef
	p.fileReaders[readerRef.String()] = fileReaderRef

	result := C.FPDFImageObj_LoadJpegFile(&pageHandle, C.int(request.Count), pageObjectHandle.handle, &readerStruct)
	if int(result) == 0 {
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

	var pageHandle C.FPDF_PAGE
	if request.Page != nil {
		loadedPage, err := p.loadPage(*request.Page)
		if err != nil {
			return nil, err
		}

		pageHandle = loadedPage.handle
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

	// Create a PDFium file access struct.
	readerStruct := C.FPDF_FILEACCESS{}
	readerStruct.m_FileLen = C.ulong(request.FileReaderSize)

	readerRef := uuid.New()
	readerRefString := readerRef.String()
	cReaderRef := C.CString(readerRefString)

	// Set the Go callback through cgo.
	C.FPDF_FILEACCESS_SET_GET_BLOCK(&readerStruct, cReaderRef)

	fileReaderRef := &fileReaderRef{
		stringRef:  unsafe.Pointer(cReaderRef),
		reader:     fileReader,
		fileAccess: &readerStruct,
	}

	Pdfium.fileReaders[readerRef.String()] = fileReaderRef
	p.fileReaders[readerRef.String()] = fileReaderRef

	result := C.FPDFImageObj_LoadJpegFileInline(&pageHandle, C.int(request.Count), pageObjectHandle.handle, &readerStruct)
	if int(result) == 0 {
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

	success := C.FPDFImageObj_SetMatrix(imageObjectHandle.handle, C.double(request.Transform.A), C.double(request.Transform.B), C.double(request.Transform.C), C.double(request.Transform.D), C.double(request.Transform.E), C.double(request.Transform.F))
	if int(success) == 0 {
		return nil, errors.New("could not set image object matrix")
	}

	return &responses.FPDFImageObj_SetMatrix{}, nil
}

// FPDFImageObj_SetBitmap sets the given bitmap to the given image object.
func (p *PdfiumImplementation) FPDFImageObj_SetBitmap(request *requests.FPDFImageObj_SetBitmap) (*responses.FPDFImageObj_SetBitmap, error) {
	p.Lock()
	defer p.Unlock()

	var pageHandle C.FPDF_PAGE
	if request.Page != nil {
		loadedPage, err := p.loadPage(*request.Page)
		if err != nil {
			return nil, err
		}

		pageHandle = loadedPage.handle
	}

	imageObjectHandle, err := p.getPageObjectHandle(request.ImageObject)
	if err != nil {
		return nil, err
	}

	bitmapHandle, err := p.getBitmapHandle(request.Bitmap)
	if err != nil {
		return nil, err
	}

	success := C.FPDFImageObj_SetBitmap(&pageHandle, C.int(request.Count), imageObjectHandle.handle, bitmapHandle.handle)
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

	bitmap := C.FPDFImageObj_GetBitmap(imageObjectHandle.handle)
	if bitmap == nil {
		return nil, errors.New("could not get bitmap")
	}

	bitmapHandle := p.registerBitmap(bitmap)

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

	imageDataLength := C.FPDFImageObj_GetImageDataDecoded(imageObjectHandle.handle, nil, 0)
	if int(imageDataLength) == 0 {
		return nil, errors.New("could not get decoded image data")
	}

	valueData := make([]byte, uint64(imageDataLength))
	C.FPDFImageObj_GetImageDataDecoded(imageObjectHandle.handle, unsafe.Pointer(&valueData[0]), C.ulong(len(valueData)))

	return &responses.FPDFImageObj_GetImageDataDecoded{
		Data: valueData,
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

	imageDataLength := C.FPDFImageObj_GetImageDataRaw(imageObjectHandle.handle, nil, 0)
	if int(imageDataLength) == 0 {
		return nil, errors.New("could not get raw image data")
	}

	valueData := make([]byte, uint64(imageDataLength))
	C.FPDFImageObj_GetImageDataRaw(imageObjectHandle.handle, unsafe.Pointer(&valueData[0]), C.ulong(len(valueData)))

	return &responses.FPDFImageObj_GetImageDataRaw{
		Data: valueData,
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

	count := C.FPDFImageObj_GetImageFilterCount(imageObjectHandle.handle)

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

	imageFilterLength := C.FPDFImageObj_GetImageFilter(imageObjectHandle.handle, C.int(request.Index), nil, 0)
	if int(imageFilterLength) == 0 {
		return nil, errors.New("could not get image filter")
	}

	charData := make([]byte, uint64(imageFilterLength))
	C.FPDFImageObj_GetImageFilter(imageObjectHandle.handle, C.int(request.Index), unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

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

	metadata := C.FPDF_IMAGEOBJ_METADATA{}

	success := C.FPDFImageObj_GetImageMetadata(imageObjectHandle.handle, pageHandle.handle, &metadata)
	if int(success) == 0 {
		return nil, errors.New("could not load image metadata")
	}

	return &responses.FPDFImageObj_GetImageMetadata{
		ImageMetadata: structs.FPDF_IMAGEOBJ_METADATA{
			Width:           uint(metadata.width),
			Height:          uint(metadata.height),
			HorizontalDPI:   float32(metadata.horizontal_dpi),
			VerticalDPI:     float32(metadata.vertical_dpi),
			BitsPerPixel:    uint(metadata.bits_per_pixel),
			Colorspace:      enums.FPDF_COLORSPACE(metadata.colorspace),
			MarkedContentID: int(metadata.marked_content_id),
		},
	}, nil
}

// FPDFPageObj_CreateNewPath creates a new path object at an initial position.
func (p *PdfiumImplementation) FPDFPageObj_CreateNewPath(request *requests.FPDFPageObj_CreateNewPath) (*responses.FPDFPageObj_CreateNewPath, error) {
	p.Lock()
	defer p.Unlock()

	pageObject := C.FPDFPageObj_CreateNewPath(C.float(request.X), C.float(request.Y))
	pageObjectHandle := p.registerPageObject(pageObject)

	return &responses.FPDFPageObj_CreateNewPath{
		PageObject: pageObjectHandle.nativeRef,
	}, nil
}

// FPDFPageObj_CreateNewRect creates a closed path consisting of a rectangle.
func (p *PdfiumImplementation) FPDFPageObj_CreateNewRect(request *requests.FPDFPageObj_CreateNewRect) (*responses.FPDFPageObj_CreateNewRect, error) {
	p.Lock()
	defer p.Unlock()

	pageObject := C.FPDFPageObj_CreateNewRect(C.float(request.X), C.float(request.Y), C.float(request.W), C.float(request.H))
	pageObjectHandle := p.registerPageObject(pageObject)

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

	left := C.float(0)
	bottom := C.float(0)
	right := C.float(0)
	top := C.float(0)
	success := C.FPDFPageObj_GetBounds(pageObjectHandle.handle, &left, &bottom, &right, &top)
	if int(success) == 0 {
		return nil, errors.New("could not get page object bounds")
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

	blendMode := C.CString(string(request.BlendMode))
	defer C.free(unsafe.Pointer(blendMode))

	C.FPDFPageObj_SetBlendMode(pageObjectHandle.handle, blendMode)

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

	success := C.FPDFPageObj_SetStrokeColor(pageObjectHandle.handle, C.uint(request.StrokeColor.R), C.uint(request.StrokeColor.G), C.uint(request.StrokeColor.B), C.uint(request.StrokeColor.A))
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

	r := C.uint(0)
	g := C.uint(0)
	b := C.uint(0)
	a := C.uint(0)

	success := C.FPDFPageObj_GetStrokeColor(pageObjectHandle.handle, &r, &g, &b, &a)
	if int(success) == 0 {
		return nil, errors.New("could not get page object stroke color")
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

	success := C.FPDFPageObj_SetStrokeWidth(pageObjectHandle.handle, C.float(request.StrokeWidth))
	if int(success) == 0 {
		return nil, errors.New("could not set page object stroke width")
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

	strokeWidth := C.float(0)
	success := C.FPDFPageObj_GetStrokeWidth(pageObjectHandle.handle, &strokeWidth)
	if int(success) == 0 {
		return nil, errors.New("could not get page object stroke width")
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

	lineJoin := C.FPDFPageObj_GetLineJoin(pageObjectHandle.handle)
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

	success := C.FPDFPageObj_SetLineJoin(pageObjectHandle.handle, C.int(request.LineJoin))
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

	lineCap := C.FPDFPageObj_GetLineCap(pageObjectHandle.handle)
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

	success := C.FPDFPageObj_SetLineJoin(pageObjectHandle.handle, C.int(request.LineCap))
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

	success := C.FPDFPageObj_SetFillColor(pageObjectHandle.handle, C.uint(request.FillColor.R), C.uint(request.FillColor.G), C.uint(request.FillColor.B), C.uint(request.FillColor.A))
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

	r := C.uint(0)
	g := C.uint(0)
	b := C.uint(0)
	a := C.uint(0)

	success := C.FPDFPageObj_GetFillColor(pageObjectHandle.handle, &r, &g, &b, &a)
	if int(success) == 0 {
		return nil, errors.New("could not get page object fill color")
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

	count := C.FPDFPath_CountSegments(pageObjectHandle.handle)
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

	segment := C.FPDFPath_GetPathSegment(pageObjectHandle.handle, C.int(request.Index))
	if segment == nil {
		return nil, errors.New("could not get path segment")
	}

	pathSegmentHandle := p.registerPathSegment(segment)

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

	x := C.float(0)
	y := C.float(0)
	success := C.FPDFPathSegment_GetPoint(pathSegmentHandle.handle, &x, &y)
	if int(success) == 0 {
		return nil, errors.New("could not get path segment point")
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

	segmentType := C.FPDFPathSegment_GetType(pathSegmentHandle.handle)

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

	getClose := C.FPDFPathSegment_GetClose(pathSegmentHandle.handle)

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

	result := C.FPDFPath_MoveTo(pageObjectHandle.handle, C.float(request.X), C.float(request.Y))
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

	result := C.FPDFPath_LineTo(pageObjectHandle.handle, C.float(request.X), C.float(request.Y))
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

	result := C.FPDFPath_BezierTo(pageObjectHandle.handle, C.float(request.X1), C.float(request.Y1), C.float(request.X2), C.float(request.Y2), C.float(request.X3), C.float(request.Y3))
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

	result := C.FPDFPath_Close(pageObjectHandle.handle)
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

	result := C.FPDFPath_SetDrawMode(pageObjectHandle.handle, C.int(request.FillMode), C.FPDF_BOOL(stroke))
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

	fillMode := C.int(0)
	stroke := C.FPDF_BOOL(0)

	result := C.FPDFPath_GetDrawMode(pageObjectHandle.handle, &fillMode, &stroke)
	if int(result) == 0 {
		return nil, errors.New("could not get draw mode")
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

	font := C.CString(request.Font)
	defer C.free(unsafe.Pointer(font))

	textObject := C.FPDFPageObj_NewTextObj(documentHandle.handle, font, C.float(request.FontSize))
	textObjectHandle := p.registerPageObject(textObject)

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

	transformedText, err := p.transformUTF8ToUTF16LE(request.Text)
	if err != nil {
		return nil, err
	}

	result := C.FPDFText_SetText(pageObjectHandle.handle, (C.FPDF_WIDESTRING)(unsafe.Pointer(&transformedText[0])))
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

	charCodes := make([]C.uint32_t, len(request.CharCodes))
	for i := range request.CharCodes {
		charCodes[i] = C.uint32_t(request.CharCodes[i])
	}

	result := C.FPDFText_SetCharcodes(pageObjectHandle.handle, (*C.uint32_t)(unsafe.Pointer(&charCodes[0])), C.size_t(len(request.CharCodes)))
	if int(result) == 0 {
		return nil, errors.New("could not set charcodes")
	}

	return &responses.FPDFText_SetCharcodes{}, nil
}

// FPDFText_LoadFont returns a font object loaded from a stream of data. The font is loaded
// into the document.
// The loaded font can be closed using FPDFFont_Close.
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

	font := C.FPDFText_LoadFont(documentHandle.handle, (*C.uchar)(unsafe.Pointer(&request.Data[0])), C.uint32_t(len(request.Data)), C.int(request.FontType), C.FPDF_BOOL(cid))
	fontHandle := p.registerFont(font)

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

	fontSize := C.float(0)
	result := C.FPDFTextObj_GetFontSize(pageObjectHandle.handle, &fontSize)
	if int(result) == 0 {
		return nil, errors.New("could not get font size")
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

	C.FPDFFont_Close(fontHandle.handle)

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

	textObject := C.FPDFPageObj_CreateTextObj(documentHandle.handle, fontHandle.handle, C.float(request.FontSize))
	textObjectHandle := p.registerPageObject(textObject)

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

	textRenderMode := C.FPDFTextObj_GetTextRenderMode(pageObjectHandle.handle)

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
	textSize := C.FPDFTextObj_GetText(pageObjectHandle.handle, textPageHandle.handle, nil, 0)
	if textSize == 0 {
		return nil, errors.New("could not get text")
	}

	charData := make([]byte, textSize)
	C.FPDFTextObj_GetText(pageObjectHandle.handle, textPageHandle.handle, (*C.ushort)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

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

	count := C.FPDFFormObj_CountObjects(pageObjectHandle.handle)

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

	formObject := C.FPDFFormObj_GetObject(pageObjectHandle.handle, C.ulong(request.Index))
	if formObject == nil {
		return nil, errors.New("could not get form object")
	}

	formObjectHandle := p.registerPageObject(formObject)

	return &responses.FPDFFormObj_GetObject{
		PageObject: formObjectHandle.nativeRef,
	}, nil
}
