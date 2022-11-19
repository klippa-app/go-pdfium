package implementation_webassembly

import (
	"errors"
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

	// @todo: fix me, what if we get return -1? Rolls over to 4294967295 because of uint64.
	securityHandlerRevision := res[0]

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
