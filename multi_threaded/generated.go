// Code generated by tool. DO NOT EDIT.
// See the code_generation package.

package multi_threaded

import (
    "errors"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

func (i *pdfiumInstance) FPDFDoc_GetPageMode(request *requests.FPDFDoc_GetPageMode) (*responses.FPDFDoc_GetPageMode, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDFDoc_GetPageMode(request)
}

func (i *pdfiumInstance) FPDFPage_Flatten(request *requests.FPDFPage_Flatten) (*responses.FPDFPage_Flatten, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDFPage_Flatten(request)
}

func (i *pdfiumInstance) FPDFPage_GetRotation(request *requests.FPDFPage_GetRotation) (*responses.FPDFPage_GetRotation, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDFPage_GetRotation(request)
}

func (i *pdfiumInstance) FPDFPage_HasTransparency(request *requests.FPDFPage_HasTransparency) (*responses.FPDFPage_HasTransparency, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDFPage_HasTransparency(request)
}

func (i *pdfiumInstance) FPDFPage_SetRotation(request *requests.FPDFPage_SetRotation) (*responses.FPDFPage_SetRotation, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDFPage_SetRotation(request)
}

func (i *pdfiumInstance) FPDF_ClosePage(request *requests.FPDF_ClosePage) (*responses.FPDF_ClosePage, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_ClosePage(request)
}

func (i *pdfiumInstance) FPDF_CopyViewerPreferences(request *requests.FPDF_CopyViewerPreferences) (*responses.FPDF_CopyViewerPreferences, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_CopyViewerPreferences(request)
}

func (i *pdfiumInstance) FPDF_CreateNewDocument(request *requests.FPDF_CreateNewDocument) (*responses.FPDF_CreateNewDocument, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_CreateNewDocument(request)
}

func (i *pdfiumInstance) FPDF_GetDocPermissions(request *requests.FPDF_GetDocPermissions) (*responses.FPDF_GetDocPermissions, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_GetDocPermissions(request)
}

func (i *pdfiumInstance) FPDF_GetFileVersion(request *requests.FPDF_GetFileVersion) (*responses.FPDF_GetFileVersion, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_GetFileVersion(request)
}

func (i *pdfiumInstance) FPDF_GetLastError(request *requests.FPDF_GetLastError) (*responses.FPDF_GetLastError, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_GetLastError(request)
}

func (i *pdfiumInstance) FPDF_GetMetaText(request *requests.FPDF_GetMetaText) (*responses.FPDF_GetMetaText, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_GetMetaText(request)
}

func (i *pdfiumInstance) FPDF_GetPageCount(request *requests.FPDF_GetPageCount) (*responses.FPDF_GetPageCount, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_GetPageCount(request)
}

func (i *pdfiumInstance) FPDF_GetPageHeight(request *requests.FPDF_GetPageHeight) (*responses.FPDF_GetPageHeight, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_GetPageHeight(request)
}

func (i *pdfiumInstance) FPDF_GetPageSizeByIndex(request *requests.FPDF_GetPageSizeByIndex) (*responses.FPDF_GetPageSizeByIndex, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_GetPageSizeByIndex(request)
}

func (i *pdfiumInstance) FPDF_GetPageWidth(request *requests.FPDF_GetPageWidth) (*responses.FPDF_GetPageWidth, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_GetPageWidth(request)
}

func (i *pdfiumInstance) FPDF_GetSecurityHandlerRevision(request *requests.FPDF_GetSecurityHandlerRevision) (*responses.FPDF_GetSecurityHandlerRevision, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_GetSecurityHandlerRevision(request)
}

func (i *pdfiumInstance) FPDF_ImportPages(request *requests.FPDF_ImportPages) (*responses.FPDF_ImportPages, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_ImportPages(request)
}

func (i *pdfiumInstance) FPDF_LoadPage(request *requests.FPDF_LoadPage) (*responses.FPDF_LoadPage, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_LoadPage(request)
}

func (i *pdfiumInstance) FPDF_SaveAsCopy(request *requests.FPDF_SaveAsCopy) (*responses.FPDF_SaveAsCopy, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_SaveAsCopy(request)
}

func (i *pdfiumInstance) FPDF_SaveWithVersion(request *requests.FPDF_SaveWithVersion) (*responses.FPDF_SaveWithVersion, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_SaveWithVersion(request)
}

func (i *pdfiumInstance) FPDF_SetSandBoxPolicy(request *requests.FPDF_SetSandBoxPolicy) (*responses.FPDF_SetSandBoxPolicy, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.FPDF_SetSandBoxPolicy(request)
}

func (i *pdfiumInstance) GetMetaData(request *requests.GetMetaData) (*responses.GetMetaData, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.GetMetaData(request)
}

func (i *pdfiumInstance) GetPageSize(request *requests.GetPageSize) (*responses.GetPageSize, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.GetPageSize(request)
}

func (i *pdfiumInstance) GetPageSizeInPixels(request *requests.GetPageSizeInPixels) (*responses.GetPageSizeInPixels, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.GetPageSizeInPixels(request)
}

func (i *pdfiumInstance) GetPageText(request *requests.GetPageText) (*responses.GetPageText, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.GetPageText(request)
}

func (i *pdfiumInstance) GetPageTextStructured(request *requests.GetPageTextStructured) (*responses.GetPageTextStructured, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.GetPageTextStructured(request)
}

func (i *pdfiumInstance) OpenDocument(request *requests.OpenDocument) (*responses.OpenDocument, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.OpenDocument(request)
}

func (i *pdfiumInstance) RenderPageInDPI(request *requests.RenderPageInDPI) (*responses.RenderPage, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.RenderPageInDPI(request)
}

func (i *pdfiumInstance) RenderPageInPixels(request *requests.RenderPageInPixels) (*responses.RenderPage, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.RenderPageInPixels(request)
}

func (i *pdfiumInstance) RenderPagesInDPI(request *requests.RenderPagesInDPI) (*responses.RenderPages, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.RenderPagesInDPI(request)
}

func (i *pdfiumInstance) RenderPagesInPixels(request *requests.RenderPagesInPixels) (*responses.RenderPages, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.RenderPagesInPixels(request)
}

func (i *pdfiumInstance) RenderToFile(request *requests.RenderToFile) (*responses.RenderToFile, error) {
	if i.closed {
		return nil, errors.New("instance is closed")
	}
	return i.worker.plugin.RenderToFile(request)
}
