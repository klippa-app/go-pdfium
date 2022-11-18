package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDF_LoadDocument opens and load a PDF document from a file path.
// Loaded document can be closed by FPDF_CloseDocument().
// If this function fails, you can use FPDF_GetLastError() to retrieve
// the reason why it failed.
func (p *PdfiumImplementation) FPDF_LoadDocument(request *requests.FPDF_LoadDocument) (*responses.FPDF_LoadDocument, error) {
	// Don't lock, OpenDocument will do that.
	doc, err := p.OpenDocument(&requests.OpenDocument{
		FilePath: request.Path,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_LoadDocument{
		Document: doc.Document,
	}, nil
}

// FPDF_LoadMemDocument opens and load a PDF document from memory.
// Loaded document can be closed by FPDF_CloseDocument().
// If this function fails, you can use FPDF_GetLastError() to retrieve
// the reason why it failed.
func (p *PdfiumImplementation) FPDF_LoadMemDocument(request *requests.FPDF_LoadMemDocument) (*responses.FPDF_LoadMemDocument, error) {
	// Don't lock, OpenDocument will do that.
	doc, err := p.OpenDocument(&requests.OpenDocument{
		File:     request.Data,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_LoadMemDocument{
		Document: doc.Document,
	}, nil
}

// FPDF_LoadCustomDocument loads a PDF document from a custom access descriptor.
// This is implemented as an io.ReadSeeker in go-pdfium.
// This is only really efficient for single threaded usage, the multi-threaded
// usage will just load the file in memory because it can't transfer readers
// over gRPC. The single-threaded usage will actually efficiently walk over
// the PDF as it's being used by PDFium.
// Loaded document can be closed by FPDF_CloseDocument().
// If this function fails, you can use FPDF_GetLastError() to retrieve
// the reason why it failed.
func (p *PdfiumImplementation) FPDF_LoadCustomDocument(request *requests.FPDF_LoadCustomDocument) (*responses.FPDF_LoadCustomDocument, error) {
	// Don't lock, OpenDocument will do that.
	doc, err := p.OpenDocument(&requests.OpenDocument{
		FileReader:     request.Reader,
		FileReaderSize: request.Size,
		Password:       request.Password,
	})
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_LoadCustomDocument{
		Document: doc.Document,
	}, nil
}

// FPDF_CloseDocument closes the references, releases the resources.
func (p *PdfiumImplementation) FPDF_CloseDocument(request *requests.FPDF_CloseDocument) (*responses.FPDF_CloseDocument, error) {
	p.Lock()
	defer p.Unlock()

	nativeDocument, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	err = nativeDocument.Close()
	if err != nil {
		return nil, err
	}

	delete(p.documentRefs, nativeDocument.nativeRef)

	return &responses.FPDF_CloseDocument{}, nil
}

// FPDF_GetLastError returns the last error code of a PDFium function, which is just called.
// Usually, this function is called after a PDFium function returns, in order to check the error code of the previous PDFium function.
// If the previous SDK call succeeded, the return value of this function is not defined. This function only works in conjunction
// with APIs that mention FPDF_GetLastError() in their documentation.
// Please note that when using go-pdfium from the same instance (on single-threaded any instance)
// from different subroutines, FPDF_GetLastError might already be reset from
// executing another PDFium method.
func (p *PdfiumImplementation) FPDF_GetLastError(request *requests.FPDF_GetLastError) (*responses.FPDF_GetLastError, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDF_GetLastError{
		Error: responses.FPDF_GetLastErrorError(C.FPDF_GetLastError()),
	}, nil
}

// FPDF_SetSandBoxPolicy set the policy for the sandbox environment.
func (p *PdfiumImplementation) FPDF_SetSandBoxPolicy(request *requests.FPDF_SetSandBoxPolicy) (*responses.FPDF_SetSandBoxPolicy, error) {
	p.Lock()
	defer p.Unlock()

	enable := C.int(0)
	if request.Enable {
		enable = C.int(1)
	}

	C.FPDF_SetSandBoxPolicy(C.FPDF_DWORD(request.Policy), enable)

	return &responses.FPDF_SetSandBoxPolicy{}, nil
}

// FPDF_LoadPage loads a page and returns a reference.
func (p *PdfiumImplementation) FPDF_LoadPage(request *requests.FPDF_LoadPage) (*responses.FPDF_LoadPage, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	pageObject := C.FPDF_LoadPage(documentHandle.handle, C.int(request.Index))
	if pageObject == nil {
		return nil, pdfium_errors.ErrPage
	}

	pageHandle := p.registerPage(pageObject, request.Index, documentHandle)

	return &responses.FPDF_LoadPage{
		Page: pageHandle.nativeRef,
	}, nil
}

// FPDF_ClosePage unloads a page by reference.
func (p *PdfiumImplementation) FPDF_ClosePage(request *requests.FPDF_ClosePage) (*responses.FPDF_ClosePage, error) {
	p.Lock()
	defer p.Unlock()

	pageRef, err := p.getPageHandle(request.Page)
	if err != nil {
		return nil, err
	}

	pageRef.Close()
	delete(p.pageRefs, request.Page)

	// Remove page reference from document.
	documentHandle, err := p.getDocumentHandle(pageRef.documentRef)
	if err != nil {
		return nil, err
	}
	delete(documentHandle.pageRefs, request.Page)

	return &responses.FPDF_ClosePage{}, nil
}

// FPDF_GetFileVersion returns the version of the PDF file.
func (p *PdfiumImplementation) FPDF_GetFileVersion(request *requests.FPDF_GetFileVersion) (*responses.FPDF_GetFileVersion, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	fileVersion := C.int(0)

	success := C.FPDF_GetFileVersion(documentHandle.handle, &fileVersion)
	if int(success) == 0 {
		return nil, errors.New("could not get file version")
	}

	return &responses.FPDF_GetFileVersion{
		FileVersion: int(fileVersion),
	}, nil
}

// FPDF_GetDocPermissions returns the permissions of the PDF.
func (p *PdfiumImplementation) FPDF_GetDocPermissions(request *requests.FPDF_GetDocPermissions) (*responses.FPDF_GetDocPermissions, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	permissions := C.FPDF_GetDocPermissions(documentHandle.handle)

	docPermissions := &responses.FPDF_GetDocPermissions{
		DocPermissions: uint32(permissions),
	}

	PrintDocument := uint32(1 << 2)
	ModifyContents := uint32(1 << 3)
	CopyOrExtractText := uint32(1 << 4)
	AddOrModifyTextAnnotations := uint32(1 << 5)
	FillInExistingInteractiveFormFields := uint32(1 << 8)
	ExtractTextAndGraphics := uint32(1 << 9)
	AssembleDocument := uint32(1 << 10)
	PrintDocumentAsFaithfulDigitalCopy := uint32(1 << 11)

	hasPermission := func(permission uint32) bool {
		if docPermissions.DocPermissions&permission > 0 {
			return true
		}

		return false
	}

	docPermissions.PrintDocument = hasPermission(PrintDocument)
	docPermissions.ModifyContents = hasPermission(ModifyContents)
	docPermissions.CopyOrExtractText = hasPermission(CopyOrExtractText)
	docPermissions.AddOrModifyTextAnnotations = hasPermission(AddOrModifyTextAnnotations)
	docPermissions.FillInInteractiveFormFields = hasPermission(AddOrModifyTextAnnotations)
	docPermissions.FillInExistingInteractiveFormFields = hasPermission(FillInExistingInteractiveFormFields)
	docPermissions.ExtractTextAndGraphics = hasPermission(ExtractTextAndGraphics)
	docPermissions.AssembleDocument = hasPermission(AssembleDocument)
	docPermissions.PrintDocumentAsFaithfulDigitalCopy = hasPermission(PrintDocumentAsFaithfulDigitalCopy)

	// Calculated permissions
	docPermissions.CreateOrModifyInteractiveFormFields = docPermissions.ModifyContents && docPermissions.AddOrModifyTextAnnotations

	return docPermissions, nil
}

// FPDF_GetSecurityHandlerRevision returns the revision number of security handlers of the file.
func (p *PdfiumImplementation) FPDF_GetSecurityHandlerRevision(request *requests.FPDF_GetSecurityHandlerRevision) (*responses.FPDF_GetSecurityHandlerRevision, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	securityHandlerRevision := C.FPDF_GetSecurityHandlerRevision(documentHandle.handle)

	return &responses.FPDF_GetSecurityHandlerRevision{
		SecurityHandlerRevision: int(securityHandlerRevision),
	}, nil
}

// FPDF_GetPageCount counts the amount of pages.
func (p *PdfiumImplementation) FPDF_GetPageCount(request *requests.FPDF_GetPageCount) (*responses.FPDF_GetPageCount, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_GetPageCount{
		PageCount: int(C.FPDF_GetPageCount(documentHandle.handle)),
	}, nil
}

// FPDF_GetPageWidth returns the width of a page.
func (p *PdfiumImplementation) FPDF_GetPageWidth(request *requests.FPDF_GetPageWidth) (*responses.FPDF_GetPageWidth, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	width := C.FPDF_GetPageWidth(pageHandle.handle)

	return &responses.FPDF_GetPageWidth{
		Page:  pageHandle.index,
		Width: float64(width),
	}, nil
}

// FPDF_GetPageHeight returns the height of a page.
func (p *PdfiumImplementation) FPDF_GetPageHeight(request *requests.FPDF_GetPageHeight) (*responses.FPDF_GetPageHeight, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	height := C.FPDF_GetPageHeight(pageHandle.handle)

	return &responses.FPDF_GetPageHeight{
		Page:   pageHandle.index,
		Height: float64(height),
	}, nil
}

// FPDF_GetPageSizeByIndex returns the size of a page by the page index.
func (p *PdfiumImplementation) FPDF_GetPageSizeByIndex(request *requests.FPDF_GetPageSizeByIndex) (*responses.FPDF_GetPageSizeByIndex, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	width := C.double(0)
	height := C.double(0)

	result := C.FPDF_GetPageSizeByIndex(documentHandle.handle, C.int(request.Index), &width, &height)
	if int(result) == 0 {
		return nil, errors.New("could not load page size by index")
	}

	return &responses.FPDF_GetPageSizeByIndex{
		Page:   request.Index,
		Width:  float64(width),
		Height: float64(height),
	}, nil
}

// FPDF_RenderPageBitmap renders contents of a page to a device independent bitmap.
func (p *PdfiumImplementation) FPDF_RenderPageBitmap(request *requests.FPDF_RenderPageBitmap) (*responses.FPDF_RenderPageBitmap, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	bitmapHandle, err := p.getBitmapHandle(request.Bitmap)
	if err != nil {
		return nil, err
	}

	C.FPDF_RenderPageBitmap(bitmapHandle.handle, pageHandle.handle, C.int(request.StartX), C.int(request.StartY), C.int(request.SizeX), C.int(request.SizeY), C.int(request.Rotate), C.int(request.Flags))

	return &responses.FPDF_RenderPageBitmap{}, nil
}

// FPDF_RenderPageBitmapWithMatrix renders contents of a page to a device independent bitmap.
func (p *PdfiumImplementation) FPDF_RenderPageBitmapWithMatrix(request *requests.FPDF_RenderPageBitmapWithMatrix) (*responses.FPDF_RenderPageBitmapWithMatrix, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	bitmapHandle, err := p.getBitmapHandle(request.Bitmap)
	if err != nil {
		return nil, err
	}

	matrix := C.FS_MATRIX{
		a: C.float(request.Matrix.A),
		b: C.float(request.Matrix.B),
		c: C.float(request.Matrix.C),
		d: C.float(request.Matrix.D),
		e: C.float(request.Matrix.E),
		f: C.float(request.Matrix.F),
	}

	clipping := C.FS_RECTF{
		left:   C.float(request.Clipping.Left),
		top:    C.float(request.Clipping.Top),
		right:  C.float(request.Clipping.Right),
		bottom: C.float(request.Clipping.Bottom),
	}

	C.FPDF_RenderPageBitmapWithMatrix(bitmapHandle.handle, pageHandle.handle, &matrix, &clipping, C.int(request.Flags))

	return &responses.FPDF_RenderPageBitmapWithMatrix{}, nil
}

// FPDF_DeviceToPage converts the screen coordinates of a point to page coordinates.
// The page coordinate system has its origin at the left-bottom corner
// of the page, with the X-axis on the bottom going to the right, and
// the Y-axis on the left side going up.
//
// NOTE: this coordinate system can be altered when you zoom, scroll,
// or rotate a page, however, a point on the page should always have
// the same coordinate values in the page coordinate system.
//
// The device coordinate system is device dependent. For screen device,
// its origin is at the left-top corner of the window. However this
// origin can be altered by the Windows coordinate transformation
// utilities.
//
// You must make sure the start_x, start_y, size_x, size_y
// and rotate parameters have exactly same values as you used in
// the FPDF_RenderPage() function call.
func (p *PdfiumImplementation) FPDF_DeviceToPage(request *requests.FPDF_DeviceToPage) (*responses.FPDF_DeviceToPage, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	pageX := C.double(0)
	pageY := C.double(0)

	success := C.FPDF_DeviceToPage(pageHandle.handle, C.int(request.StartX), C.int(request.StartY), C.int(request.SizeX), C.int(request.SizeY), C.int(request.Rotate), C.int(request.DeviceX), C.int(request.DeviceY), &pageX, &pageY)
	if int(success) == 0 {
		return nil, errors.New("could not calculate from device to page")
	}

	return &responses.FPDF_DeviceToPage{
		PageX: float64(pageX),
		PageY: float64(pageY),
	}, nil
}

// FPDF_PageToDevice converts the page coordinates of a point to screen coordinates.
// See comments for FPDF_DeviceToPage().
func (p *PdfiumImplementation) FPDF_PageToDevice(request *requests.FPDF_PageToDevice) (*responses.FPDF_PageToDevice, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	deviceX := C.int(0)
	deviceY := C.int(0)

	success := C.FPDF_PageToDevice(pageHandle.handle, C.int(request.StartX), C.int(request.StartY), C.int(request.SizeX), C.int(request.SizeY), C.int(request.Rotate), C.double(request.PageX), C.double(request.PageY), &deviceX, &deviceY)
	if int(success) == 0 {
		return nil, errors.New("could not calculate from page to device")
	}

	return &responses.FPDF_PageToDevice{
		DeviceX: int(deviceX),
		DeviceY: int(deviceY),
	}, nil
}

// FPDFBitmap_Create Create a device independent bitmap (FXDIB).
func (p *PdfiumImplementation) FPDFBitmap_Create(request *requests.FPDFBitmap_Create) (*responses.FPDFBitmap_Create, error) {
	p.Lock()
	defer p.Unlock()

	bitmap := C.FPDFBitmap_Create(C.int(request.Width), C.int(request.Height), C.int(request.Alpha))
	bitmapHandle := p.registerBitmap(bitmap)

	return &responses.FPDFBitmap_Create{
		Bitmap: bitmapHandle.nativeRef,
	}, nil
}

// FPDFBitmap_CreateEx Create a device independent bitmap (FXDIB) with an
// external buffer.
// Similar to FPDFBitmap_Create function, but allows for more formats
// and an external buffer is supported. The bitmap created by this
// function can be used in any place that a FPDF_BITMAP handle is
// required.
//
// If an external buffer is used, then the caller should destroy the
// buffer. FPDFBitmap_Destroy() will not destroy the buffer.
//
// It is recommended to use FPDFBitmap_GetStride() to get the stride
// value.
//
// Not supported on multi-threaded usage.
func (p *PdfiumImplementation) FPDFBitmap_CreateEx(request *requests.FPDFBitmap_CreateEx) (*responses.FPDFBitmap_CreateEx, error) {
	p.Lock()
	defer p.Unlock()

	bitmap := C.FPDFBitmap_CreateEx(C.int(request.Width), C.int(request.Height), C.int(request.Format), unsafe.Pointer(&request.Buffer[0]), C.int(request.Stride))
	bitmapHandle := p.registerBitmap(bitmap)

	return &responses.FPDFBitmap_CreateEx{
		Bitmap: bitmapHandle.nativeRef,
	}, nil
}

// FPDFBitmap_GetFormat returns the format of the bitmap.
// Only formats supported by FPDFBitmap_CreateEx are supported by this function.
func (p *PdfiumImplementation) FPDFBitmap_GetFormat(request *requests.FPDFBitmap_GetFormat) (*responses.FPDFBitmap_GetFormat, error) {
	p.Lock()
	defer p.Unlock()

	bitmapHandle, err := p.getBitmapHandle(request.Bitmap)
	if err != nil {
		return nil, err
	}

	format := C.FPDFBitmap_GetFormat(bitmapHandle.handle)

	return &responses.FPDFBitmap_GetFormat{
		Format: enums.FPDF_BITMAP_FORMAT(format),
	}, nil
}

// FPDFBitmap_FillRect fills a rectangle in a bitmap.
// This function sets the color and (optionally) alpha value in the
// specified region of the bitmap.
//
// NOTE: If the alpha channel is used, this function does NOT
// composite the background with the source color, instead the
// background will be replaced by the source color and the alpha.
//
// If the alpha channel is not used, the alpha parameter is ignored.
func (p *PdfiumImplementation) FPDFBitmap_FillRect(request *requests.FPDFBitmap_FillRect) (*responses.FPDFBitmap_FillRect, error) {
	p.Lock()
	defer p.Unlock()

	bitmapHandle, err := p.getBitmapHandle(request.Bitmap)
	if err != nil {
		return nil, err
	}

	C.FPDFBitmap_FillRect(bitmapHandle.handle, C.int(request.Left), C.int(request.Top), C.int(request.Width), C.int(request.Height), C.ulong(request.Color))

	return &responses.FPDFBitmap_FillRect{}, nil
}

// FPDFBitmap_GetBuffer returns the data buffer of a bitmap.
// The stride may be more than width * number of bytes per pixel
//
// Applications can use this function to get the bitmap buffer pointer,
// then manipulate any color and/or alpha values for any pixels in the
// bitmap.
//
// The data is in BGRA format. Where the A maybe unused if alpha was
// not specified.
func (p *PdfiumImplementation) FPDFBitmap_GetBuffer(request *requests.FPDFBitmap_GetBuffer) (*responses.FPDFBitmap_GetBuffer, error) {
	p.Lock()
	defer p.Unlock()

	bitmapHandle, err := p.getBitmapHandle(request.Bitmap)
	if err != nil {
		return nil, err
	}

	// We need to calculate the buffer size, this is stride (bytes per bitmap line) * height.
	stride := C.FPDFBitmap_GetStride(bitmapHandle.handle)
	height := C.FPDFBitmap_GetHeight(bitmapHandle.handle)
	size := int(stride * height)

	// The pointer to the first byte of the bitmap buffer.
	buffer := C.FPDFBitmap_GetBuffer(bitmapHandle.handle)

	// We create a Go slice backed by a C array (without copying the original data),
	// and acquire its length at runtime and use a type conversion to a pointer to a very big array and then slice it to the length that we want.
	// Refer https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
	data := (*[1<<50 - 1]byte)(unsafe.Pointer(buffer))[:size:size] // For 64-bit machine, the max number it can go is 50 as per https://github.com/golang/go/issues/13656#issuecomment-291957684

	return &responses.FPDFBitmap_GetBuffer{
		Buffer: data,
	}, nil
}

// FPDFBitmap_GetWidth returns the width of a bitmap.
func (p *PdfiumImplementation) FPDFBitmap_GetWidth(request *requests.FPDFBitmap_GetWidth) (*responses.FPDFBitmap_GetWidth, error) {
	p.Lock()
	defer p.Unlock()

	bitmapHandle, err := p.getBitmapHandle(request.Bitmap)
	if err != nil {
		return nil, err
	}

	width := C.FPDFBitmap_GetHeight(bitmapHandle.handle)
	return &responses.FPDFBitmap_GetWidth{
		Width: int(width),
	}, nil
}

// FPDFBitmap_GetHeight returns the height of a bitmap.
func (p *PdfiumImplementation) FPDFBitmap_GetHeight(request *requests.FPDFBitmap_GetHeight) (*responses.FPDFBitmap_GetHeight, error) {
	p.Lock()
	defer p.Unlock()

	bitmapHandle, err := p.getBitmapHandle(request.Bitmap)
	if err != nil {
		return nil, err
	}

	height := C.FPDFBitmap_GetHeight(bitmapHandle.handle)
	return &responses.FPDFBitmap_GetHeight{
		Height: int(height),
	}, nil
}

// FPDFBitmap_GetStride returns the number of bytes for each line in the bitmap buffer.
func (p *PdfiumImplementation) FPDFBitmap_GetStride(request *requests.FPDFBitmap_GetStride) (*responses.FPDFBitmap_GetStride, error) {
	p.Lock()
	defer p.Unlock()

	bitmapHandle, err := p.getBitmapHandle(request.Bitmap)
	if err != nil {
		return nil, err
	}

	stride := C.FPDFBitmap_GetStride(bitmapHandle.handle)
	return &responses.FPDFBitmap_GetStride{
		Stride: int(stride),
	}, nil
}

// FPDFBitmap_Destroy destroys a bitmap and release all related buffers.
// This function will not destroy any external buffers provided when
// the bitmap was created.
func (p *PdfiumImplementation) FPDFBitmap_Destroy(request *requests.FPDFBitmap_Destroy) (*responses.FPDFBitmap_Destroy, error) {
	p.Lock()
	defer p.Unlock()

	bitmapHandle, err := p.getBitmapHandle(request.Bitmap)
	if err != nil {
		return nil, err
	}

	C.FPDFBitmap_Destroy(bitmapHandle.handle)

	delete(p.bitmapRefs, bitmapHandle.nativeRef)

	return &responses.FPDFBitmap_Destroy{}, nil
}

// FPDF_VIEWERREF_GetPrintScaling returns whether the PDF document prefers to be scaled or not.
func (p *PdfiumImplementation) FPDF_VIEWERREF_GetPrintScaling(request *requests.FPDF_VIEWERREF_GetPrintScaling) (*responses.FPDF_VIEWERREF_GetPrintScaling, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	printScaling := C.FPDF_VIEWERREF_GetPrintScaling(documentHandle.handle)
	return &responses.FPDF_VIEWERREF_GetPrintScaling{
		PreferPrintScaling: int(printScaling) == 1,
	}, nil
}

// FPDF_VIEWERREF_GetNumCopies returns the number of copies to be printed.
func (p *PdfiumImplementation) FPDF_VIEWERREF_GetNumCopies(request *requests.FPDF_VIEWERREF_GetNumCopies) (*responses.FPDF_VIEWERREF_GetNumCopies, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	numCopies := C.FPDF_VIEWERREF_GetNumCopies(documentHandle.handle)
	return &responses.FPDF_VIEWERREF_GetNumCopies{
		NumCopies: int(numCopies),
	}, nil
}

// FPDF_VIEWERREF_GetPrintPageRange returns the page numbers to initialize print dialog box when file is printed.
func (p *PdfiumImplementation) FPDF_VIEWERREF_GetPrintPageRange(request *requests.FPDF_VIEWERREF_GetPrintPageRange) (*responses.FPDF_VIEWERREF_GetPrintPageRange, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	pageRange := C.FPDF_VIEWERREF_GetPrintPageRange(documentHandle.handle)
	pageRangeHandle := p.registerPageRange(pageRange, documentHandle)

	return &responses.FPDF_VIEWERREF_GetPrintPageRange{
		PageRange: pageRangeHandle.nativeRef,
	}, nil
}

// FPDF_VIEWERREF_GetDuplex returns the paper handling option to be used when printing from the print dialog.
func (p *PdfiumImplementation) FPDF_VIEWERREF_GetDuplex(request *requests.FPDF_VIEWERREF_GetDuplex) (*responses.FPDF_VIEWERREF_GetDuplex, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	duplexType := C.FPDF_VIEWERREF_GetDuplex(documentHandle.handle)
	return &responses.FPDF_VIEWERREF_GetDuplex{
		DuplexType: enums.FPDF_DUPLEXTYPE(duplexType),
	}, nil
}

// FPDF_VIEWERREF_GetName returns the contents for a viewer ref, with a given key. The value must be of type "name".
func (p *PdfiumImplementation) FPDF_VIEWERREF_GetName(request *requests.FPDF_VIEWERREF_GetName) (*responses.FPDF_VIEWERREF_GetName, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	cstr := C.CString(request.Key)
	defer C.free(unsafe.Pointer(cstr))

	// First get the metadata length.
	nameSize := C.FPDF_VIEWERREF_GetName(documentHandle.handle, cstr, nil, 0)
	if nameSize == 0 {
		return nil, errors.New("could not get name")
	}

	charData := make([]byte, uint64(nameSize))
	C.FPDF_VIEWERREF_GetName(documentHandle.handle, cstr, (*C.char)(unsafe.Pointer(&charData[0])), C.ulong(len(charData)))

	return &responses.FPDF_VIEWERREF_GetName{
		Value: string(charData[:len(charData)-1]), // Remove nil terminator
	}, nil
}

// FPDF_CountNamedDests returns the count of named destinations in the PDF document.
func (p *PdfiumImplementation) FPDF_CountNamedDests(request *requests.FPDF_CountNamedDests) (*responses.FPDF_CountNamedDests, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	count := C.FPDF_CountNamedDests(documentHandle.handle)
	return &responses.FPDF_CountNamedDests{
		Count: uint64(count),
	}, nil
}

// FPDF_GetNamedDestByName returns the destination handle for the given name.
func (p *PdfiumImplementation) FPDF_GetNamedDestByName(request *requests.FPDF_GetNamedDestByName) (*responses.FPDF_GetNamedDestByName, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	cstr := C.CString(request.Name)
	defer C.free(unsafe.Pointer(cstr))

	dest := C.FPDF_GetNamedDestByName(documentHandle.handle, cstr)
	if dest == nil {
		return nil, errors.New("could not get named dest by name")
	}

	destHandle := p.registerDest(dest, documentHandle)
	return &responses.FPDF_GetNamedDestByName{
		Dest: destHandle.nativeRef,
	}, nil
}

// FPDF_GetNamedDest returns the named destination by index.
func (p *PdfiumImplementation) FPDF_GetNamedDest(request *requests.FPDF_GetNamedDest) (*responses.FPDF_GetNamedDest, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	bufLen := C.long(0)

	// First get the name length.
	C.FPDF_GetNamedDest(documentHandle.handle, C.int(request.Index), nil, &bufLen)
	if int64(bufLen) <= 0 {
		return nil, errors.New("could not get name of named dest")
	}

	charData := make([]byte, int64(bufLen))
	dest := C.FPDF_GetNamedDest(documentHandle.handle, C.int(request.Index), unsafe.Pointer(&charData[0]), &bufLen)

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	destHandle := p.registerDest(dest, documentHandle)
	return &responses.FPDF_GetNamedDest{
		Dest: destHandle.nativeRef,
		Name: transformedText,
	}, nil
}
