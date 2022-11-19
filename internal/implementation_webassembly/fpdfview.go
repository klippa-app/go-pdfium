package implementation_webassembly

import "C"
import (
	"errors"
	"github.com/klippa-app/go-pdfium/enums"
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"unsafe"
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

// FPDF_LoadMemDocument64 opens and load a PDF document from memory.
// Loaded document can be closed by FPDF_CloseDocument().
// If this function fails, you can use FPDF_GetLastError() to retrieve
// the reason why it failed.
// Experimental API.
func (p *PdfiumImplementation) FPDF_LoadMemDocument64(request *requests.FPDF_LoadMemDocument64) (*responses.FPDF_LoadMemDocument64, error) {
	// Don't lock, OpenDocument will do that.
	doc, err := p.OpenDocument(&requests.OpenDocument{
		File:     request.Data,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_LoadMemDocument64{
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

	err = nativeDocument.Close(p)
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

	errorCode, err := p.module.ExportedFunction("FPDF_GetLastError").Call(p.context)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_GetLastError{
		Error: responses.FPDF_GetLastErrorError(int(errorCode[0])),
	}, nil
}

// FPDF_SetSandBoxPolicy set the policy for the sandbox environment.
func (p *PdfiumImplementation) FPDF_SetSandBoxPolicy(request *requests.FPDF_SetSandBoxPolicy) (*responses.FPDF_SetSandBoxPolicy, error) {
	p.Lock()
	defer p.Unlock()

	enable := uint64(0)
	if request.Enable {
		enable = uint64(1)
	}

	_, err := p.module.ExportedFunction("FPDF_SetSandBoxPolicy").Call(p.context, uint64(request.Policy), enable)
	if err != nil {
		return nil, err
	}

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

	pageObject, err := p.module.ExportedFunction("FPDF_LoadPage").Call(p.context, *documentHandle.handle, uint64(request.Index))
	if err != nil {
		return nil, err
	}

	if len(pageObject) == 0 || pageObject[0] == 0 {
		return nil, pdfium_errors.ErrPage
	}

	pageHandle := p.registerPage(pageObject[0], request.Index, documentHandle)

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

	pageRef.Close(p)
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

	fileVersion, err := p.IntPointer()
	if err != nil {
		return nil, err
	}

	defer fileVersion.Free()

	success, err := p.module.ExportedFunction("FPDF_GetFileVersion").Call(p.context, *documentHandle.handle, fileVersion.Pointer)
	if err != nil {
		return nil, err
	}

	if len(success) == 0 || success[0] == 0 {
		return nil, errors.New("could not get file version")
	}

	val, err := fileVersion.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_GetFileVersion{
		FileVersion: val,
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

	res, err := p.module.ExportedFunction("FPDF_GetDocPermissions").Call(p.context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	permissions := res[0]

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

	res, err := p.module.ExportedFunction("FPDF_GetSecurityHandlerRevision").Call(p.context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	securityHandlerRevision := *(*int32)(unsafe.Pointer(&res[0]))

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

	res, err := p.module.ExportedFunction("FPDF_GetPageCount").Call(p.context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_GetPageCount{
		PageCount: int(res[0]),
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

	res, err := p.module.ExportedFunction("FPDF_GetPageWidth").Call(p.context, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	width := *(*float64)(unsafe.Pointer(&res[0]))

	return &responses.FPDF_GetPageWidth{
		Page:  pageHandle.index,
		Width: width,
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

	res, err := p.module.ExportedFunction("FPDF_GetPageHeight").Call(p.context, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	height := *(*float64)(unsafe.Pointer(&res[0]))

	return &responses.FPDF_GetPageHeight{
		Page:   pageHandle.index,
		Height: height,
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

	widthPointer, err := p.DoublePointer()
	if err != nil {
		return nil, err
	}
	defer widthPointer.Free()

	heightPointer, err := p.DoublePointer()
	if err != nil {
		return nil, err
	}
	defer heightPointer.Free()

	res, err := p.module.ExportedFunction("FPDF_GetPageSizeByIndex").Call(p.context, *documentHandle.handle, uint64(request.Index), widthPointer.Pointer, heightPointer.Pointer)
	if err != nil {
		return nil, err
	}

	if int(res[0]) == 0 {
		return nil, errors.New("could not load page size by index")
	}

	width, err := widthPointer.Value()
	if err != nil {
		return nil, err
	}

	height, err := heightPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_GetPageSizeByIndex{
		Page:   request.Index,
		Width:  width,
		Height: height,
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

	_, err = p.module.ExportedFunction("FPDF_RenderPageBitmap").Call(p.context, *bitmapHandle.handle, *pageHandle.handle, uint64(request.StartX), uint64(request.StartY), uint64(request.SizeX), uint64(request.SizeY), uint64(request.Rotate), uint64(request.Flags))
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_RenderPageBitmap{}, nil
}

// FPDF_RenderPageBitmapWithMatrix renders contents of a page to a device independent bitmap.
func (p *PdfiumImplementation) FPDF_RenderPageBitmapWithMatrix(request *requests.FPDF_RenderPageBitmapWithMatrix) (*responses.FPDF_RenderPageBitmapWithMatrix, error) {
	//TODO implement me
	// @todo: how to handle structs?
	return nil, pdfium_errors.ErrUnsupportedOnWebassembly
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

	pageXPointer, err := p.DoublePointer()
	if err != nil {
		return nil, err
	}
	defer pageXPointer.Free()

	pageYPointer, err := p.DoublePointer()
	if err != nil {
		return nil, err
	}
	defer pageYPointer.Free()

	res, err := p.module.ExportedFunction("FPDF_DeviceToPage").Call(p.context, *pageHandle.handle, uint64(request.StartX), uint64(request.StartY), uint64(request.SizeX), uint64(request.SizeY), uint64(request.Rotate), uint64(request.DeviceX), uint64(request.DeviceY), pageXPointer.Pointer, pageYPointer.Pointer)
	if err != nil {
		return nil, err
	}

	if int(res[0]) == 0 {
		return nil, errors.New("could not calculate from device to page")
	}

	pageX, err := pageXPointer.Value()
	if err != nil {
		return nil, err
	}

	pageY, err := pageYPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_DeviceToPage{
		PageX: pageX,
		PageY: pageY,
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

	deviceXPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer deviceXPointer.Free()

	deviceYPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer deviceYPointer.Free()

	res, err := p.module.ExportedFunction("FPDF_PageToDevice").Call(p.context, *pageHandle.handle, uint64(request.StartX), uint64(request.StartY), uint64(request.SizeX), uint64(request.SizeY), uint64(request.Rotate), uint64(request.PageX), uint64(request.PageY), deviceXPointer.Pointer, deviceYPointer.Pointer)
	if err != nil {
		return nil, err
	}

	if int(res[0]) == 0 {
		return nil, errors.New("could not calculate from page to device")
	}

	deviceX, err := deviceXPointer.Value()
	if err != nil {
		return nil, err
	}

	deviceY, err := deviceYPointer.Value()
	if err != nil {
		return nil, err
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

	res, err := p.module.ExportedFunction("FPDFBitmap_Create").Call(p.context, uint64(request.Width), uint64(request.Height), uint64(request.Alpha))
	if err != nil {
		return nil, err
	}

	bitmapHandle := p.registerBitmap(&res[0])

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
	return nil, pdfium_errors.ErrUnsupportedOnWebassembly

	p.Lock()
	defer p.Unlock()

	// @todo: figure out how we are going to do this.
	// @todo: we can't use a pointer to the byte array since this would be
	// @todo: the address in Go land, but we need the address in WASM land.
	res, err := p.module.ExportedFunction("FPDFBitmap_CreateEx").Call(p.context, uint64(request.Width), uint64(request.Height), uint64(request.Format), 0, uint64(request.Stride))
	if err != nil {
		return nil, err
	}

	bitmapHandle := p.registerBitmap(&res[0])

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

	res, err := p.module.ExportedFunction("FPDFBitmap_GetFormat").Call(p.context, *bitmapHandle.handle)
	if err != nil {
		return nil, err
	}

	format := *(*int32)(unsafe.Pointer(&res[0]))

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

	_, err = p.module.ExportedFunction("FPDFBitmap_FillRect").Call(p.context, *bitmapHandle.handle, uint64(request.Left), uint64(request.Top), uint64(request.Width), uint64(request.Height), uint64(request.Color))
	if err != nil {
		return nil, err
	}

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
	res, err := p.module.ExportedFunction("FPDFBitmap_GetStride").Call(p.context, *bitmapHandle.handle)
	if err != nil {
		return nil, err
	}

	stride := *(*int32)(unsafe.Pointer(&res[0]))

	res, err = p.module.ExportedFunction("FPDFBitmap_GetHeight").Call(p.context, *bitmapHandle.handle)
	if err != nil {
		return nil, err
	}

	height := *(*int32)(unsafe.Pointer(&res[0]))

	size := int(stride * height)

	// The pointer to the first byte of the bitmap buffer.
	res, err = p.module.ExportedFunction("FPDFBitmap_GetBuffer").Call(p.context, *bitmapHandle.handle)
	if err != nil {
		return nil, err
	}

	// Create a view of the underlying memory, not a copy.
	data, success := p.module.Memory().Read(p.context, uint32(res[0]), uint32(size))
	if !success {
		return nil, errors.New("could not get bitmap buffer")
	}

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

	res, err := p.module.ExportedFunction("FPDFBitmap_GetWidth").Call(p.context, *bitmapHandle.handle)
	if err != nil {
		return nil, err
	}

	width := *(*int32)(unsafe.Pointer(&res[0]))

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

	res, err := p.module.ExportedFunction("FPDFBitmap_GetHeight").Call(p.context, *bitmapHandle.handle)
	if err != nil {
		return nil, err
	}

	height := *(*int32)(unsafe.Pointer(&res[0]))

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

	res, err := p.module.ExportedFunction("FPDFBitmap_GetStride").Call(p.context, *bitmapHandle.handle)
	if err != nil {
		return nil, err
	}

	stride := *(*int32)(unsafe.Pointer(&res[0]))

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

	_, err = p.module.ExportedFunction("FPDFBitmap_Destroy").Call(p.context, *bitmapHandle.handle)
	if err != nil {
		return nil, err
	}

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

	res, err := p.module.ExportedFunction("FPDF_VIEWERREF_GetPrintScaling").Call(p.context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	printScaling := *(*int32)(unsafe.Pointer(&res[0]))

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

	res, err := p.module.ExportedFunction("FPDF_VIEWERREF_GetNumCopies").Call(p.context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}
	numCopies := *(*int32)(unsafe.Pointer(&res[0]))

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

	res, err := p.module.ExportedFunction("FPDF_VIEWERREF_GetPrintPageRange").Call(p.context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	pageRangeHandle := p.registerPageRange(&res[0], documentHandle)

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

	res, err := p.module.ExportedFunction("FPDF_VIEWERREF_GetDuplex").Call(p.context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}
	duplexType := *(*int32)(unsafe.Pointer(&res[0]))

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

	cstr, err := p.CString(request.Key)
	if err != nil {
		return nil, err
	}

	defer cstr.Free()

	// First get the metadata length.
	res, err := p.module.ExportedFunction("FPDF_VIEWERREF_GetName").Call(p.context, *documentHandle.handle, cstr.Pointer, 0, 0)
	if err != nil {
		return nil, err
	}

	nameSize := res[0]
	if nameSize == 0 {
		return nil, errors.New("could not get name")
	}

	charDataPointer, err := p.ByteArrayPointer(nameSize)
	if err != nil {
		return nil, err
	}

	defer charDataPointer.Free()

	_, err = p.module.ExportedFunction("FPDF_VIEWERREF_GetName").Call(p.context, *documentHandle.handle, cstr.Pointer, charDataPointer.Pointer, nameSize)
	if err != nil {
		return nil, err
	}

	charData, err := charDataPointer.Value()
	if err != nil {
		return nil, err
	}

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

	res, err := p.module.ExportedFunction("FPDF_CountNamedDests").Call(p.context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	count := res[0]

	return &responses.FPDF_CountNamedDests{
		Count: count,
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

	cstr, err := p.CString(request.Name)
	if err != nil {
		return nil, err
	}

	defer cstr.Free()

	res, err := p.module.ExportedFunction("FPDF_CountNamedDests").Call(p.context, *documentHandle.handle, cstr.Pointer)
	if err != nil {
		return nil, err
	}

	dest := res[0]
	if dest == 0 {
		return nil, errors.New("could not get named dest by name")
	}

	destHandle := p.registerDest(&dest, documentHandle)
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

	bufLenPointer, err := p.LongPointer()
	if err != nil {
		return nil, err
	}
	defer bufLenPointer.Free()

	// First get the name length.
	res, err := p.module.ExportedFunction("FPDF_VIEWERREF_GetName").Call(p.context, *documentHandle.handle, uint64(request.Index), 0, bufLenPointer.Pointer)
	if err != nil {
		return nil, err
	}

	bufLen := res[0]
	if int64(bufLen) <= 0 {
		return nil, errors.New("could not get name of named dest")
	}

	charDataPointer, err := p.ByteArrayPointer(bufLen)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	res, err = p.module.ExportedFunction("FPDF_GetNamedDest").Call(p.context, *documentHandle.handle, uint64(request.Index), charDataPointer.Pointer, bufLenPointer.Pointer)
	if err != nil {
		return nil, err
	}

	dest := res[0]

	charData, err := charDataPointer.Value()
	if err != nil {
		return nil, err
	}

	transformedText, err := p.transformUTF16LEToUTF8(charData)
	if err != nil {
		return nil, err
	}

	destHandle := p.registerDest(&dest, documentHandle)
	return &responses.FPDF_GetNamedDest{
		Dest: destHandle.nativeRef,
		Name: transformedText,
	}, nil
}
