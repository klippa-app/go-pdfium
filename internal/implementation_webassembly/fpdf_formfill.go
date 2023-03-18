package implementation_webassembly

import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"

	"github.com/tetratelabs/wazero/api"
)

type FormFillInfo struct {
	Struct           *uint64
	FormFillInfo     *structs.FPDF_FORMFILLINFO
	FormHandleHandle *FormHandleHandle
	Instance         *PdfiumImplementation
}

func (f *FormFillInfo) Release() {
	if f.FormFillInfo.Release != nil {
		f.FormFillInfo.Release()
	}
}

func (f *FormFillInfo) FFI_Invalidate_CB(page uint32, left, top, right, bottom float64) {
	var pageRef references.FPDF_PAGE
	if pointerPageRef, ok := f.FormHandleHandle.pagePointers[uint64(page)]; ok {
		pageRef = pointerPageRef
	} else {
		pageHandle := f.Instance.registerPage(uint64(page), 0, nil)
		pageRef = pageHandle.nativeRef
	}

	f.FormFillInfo.FFI_Invalidate(pageRef, left, top, right, bottom)
}

func (f *FormFillInfo) FFI_OutputSelectedRect(page uint32, left, top, right, bottom float64) {
	var pageRef references.FPDF_PAGE
	if pointerPageRef, ok := f.FormHandleHandle.pagePointers[uint64(page)]; ok {
		pageRef = pointerPageRef
	} else {
		pageHandle := f.Instance.registerPage(uint64(page), 0, nil)
		pageRef = pageHandle.nativeRef
	}

	f.FormFillInfo.FFI_OutputSelectedRect(pageRef, left, top, right, bottom)
}

func (f *FormFillInfo) FFI_SetCursor(cursor uint32) {
	f.FormFillInfo.FFI_SetCursor(enums.FXCT(cursor))
}

func (f *FormFillInfo) FFI_SetTimer(uElapse, lpTimerFunc uint32) int {
	timerFunc := func(idEvent int) {
		f.Instance.Module.ExportedFunction("FPDF_FORMFILLINFO_CALL_TIMER").Call(f.Instance.Context, *(*uint64)(unsafe.Pointer(&lpTimerFunc)), *(*uint64)(unsafe.Pointer(&idEvent)))
	}

	return f.FormFillInfo.FFI_SetTimer(int(uElapse), timerFunc)
}

func (f *FormFillInfo) FFI_KillTimer(nTimerID int) {
	f.FormFillInfo.FFI_KillTimer(nTimerID)
}

func (f *FormFillInfo) FFI_GetLocalTime() structs.FPDF_SYSTEMTIME {
	return f.FormFillInfo.FFI_GetLocalTime()
}

func (f *FormFillInfo) FFI_OnChange() {
	f.FormFillInfo.FFI_OnChange()
}

func (f *FormFillInfo) FFI_GetPage(document uint64, pageIndex int) uint64 {
	var documentRef references.FPDF_DOCUMENT
	if pointerDocumentRef, ok := f.FormHandleHandle.documentPointers[document]; ok {
		documentRef = pointerDocumentRef
	} else {
		documentHandle := f.Instance.registerDocument(&document)
		documentRef = documentHandle.nativeRef
	}

	page := f.FormFillInfo.FFI_GetPage(documentRef, int(pageIndex))
	if page == nil {
		return 0
	}

	pageHandle, err := f.Instance.getPageHandle(*page)
	if err != nil {
		return 0
	}

	return *pageHandle.handle
}

func (f *FormFillInfo) FFI_GetCurrentPage(document uint64) uint64 {
	if f.FormFillInfo.FFI_GetCurrentPage == nil {
		return 0
	}

	documentHandle := f.Instance.registerDocument(&document)
	page := f.FormFillInfo.FFI_GetCurrentPage(documentHandle.nativeRef)
	if page == nil {
		return 0
	}

	pageHandle, err := f.Instance.getPageHandle(*page)
	if err != nil {
		return 0
	}

	return *pageHandle.handle
}

func (f *FormFillInfo) FFI_GetRotation(page uint64) int {
	var pageRef references.FPDF_PAGE
	if pointerPageRef, ok := f.FormHandleHandle.pagePointers[page]; ok {
		pageRef = pointerPageRef
	} else {
		pageHandle := f.Instance.registerPage(page, 0, nil)
		pageRef = pageHandle.nativeRef
	}

	rotation := f.FormFillInfo.FFI_GetRotation(pageRef)
	return int(rotation)
}

func (f *FormFillInfo) FFI_ExecuteNamedAction(namedAction uint64) {
	// @todo: extract string.
	f.FormFillInfo.FFI_ExecuteNamedAction("")
}

func (f *FormFillInfo) FFI_SetTextFieldFocus(value uint64, valueLen uint64, isFocus uint64) {
	// @todo: extract string.
	f.FormFillInfo.FFI_SetTextFieldFocus("", isFocus == 1)
}

func (f *FormFillInfo) FFI_DoURIAction(bsURI uint64) {
	// @todo: extract string.
	f.FormFillInfo.FFI_DoURIAction("")
}

func (f *FormFillInfo) FFI_DoGoToAction(nPageIndex uint64, zoomMode uint64, fPosArray uint64, sizeofArray uint64) {
	// @todo: extract string.
	if f.FormFillInfo.FFI_DoGoToAction == nil {
		return
	}

	// @todo: extract data
	/*
		target := (*[1<<25 - 1]float32)(unsafe.Pointer(fPosArray))[:sizeofArray:sizeofArray]
		pos := make([]float32, int(sizeofArray))
		for i := range pos {
			pos[i] = float32(target[i])
		}*/

	pos := make([]float32, int(sizeofArray))
	f.FormFillInfo.FFI_DoGoToAction(int(nPageIndex), enums.FPDF_ZOOM_MODE(zoomMode), pos)
}

var FormFillInfoHandles = map[uint32]*FormFillInfo{}

// FPDFDOC_InitFormFillEnvironment initializes form fill environment
// This function should be called before any form fill operation.
// Not supported on multi-threaded usage due to its bidirectional nature.
func (p *PdfiumImplementation) FPDFDOC_InitFormFillEnvironment(request *requests.FPDFDOC_InitFormFillEnvironment) (*responses.FPDFDOC_InitFormFillEnvironment, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	if request.FormFillInfo.FFI_Invalidate == nil {
		return nil, errors.New("FormFillInfo callback FFI_Invalidate is required")
	}

	if request.FormFillInfo.FFI_SetCursor == nil {
		return nil, errors.New("FormFillInfo callback FFI_SetCursor is required")
	}

	if request.FormFillInfo.FFI_SetTimer == nil {
		return nil, errors.New("FormFillInfo callback FFI_SetTimer is required")
	}

	if request.FormFillInfo.FFI_KillTimer == nil {
		return nil, errors.New("FormFillInfo callback FFI_KillTimer is required")
	}

	if request.FormFillInfo.FFI_GetLocalTime == nil {
		return nil, errors.New("FormFillInfo callback FFI_GetLocalTime is required")
	}

	if request.FormFillInfo.FFI_GetPage == nil {
		return nil, errors.New("FormFillInfo callback FFI_GetPage is required")
	}

	if request.FormFillInfo.FFI_GetRotation == nil {
		return nil, errors.New("FormFillInfo callback FFI_GetRotation is required")
	}

	if request.FormFillInfo.FFI_ExecuteNamedAction == nil {
		return nil, errors.New("FormFillInfo callback FFI_ExecuteNamedAction is required")
	}

	res, err := p.Module.ExportedFunction("FPDF_FORMFILLINFO_Create").Call(p.Context)
	if err != nil {
		return nil, err
	}

	formInfoStruct := res[0]
	if formInfoStruct == 0 {
		return nil, errors.New("could not init form fill environment")
	}

	res, err = p.Module.ExportedFunction("FPDFDOC_InitFormFillEnvironment").Call(p.Context, *documentHandle.handle, formInfoStruct)
	if err != nil {
		return nil, err
	}

	formHandle := res[0]
	if formHandle == 0 {
		return nil, errors.New("could not init form fill environment")
	}

	formHandleHandle := p.registerFormHandle(&formHandle, &formInfoStruct)

	formFillInfo := &FormFillInfo{
		Struct:           &formInfoStruct,
		FormFillInfo:     &request.FormFillInfo,
		FormHandleHandle: formHandleHandle,
		Instance:         p,
	}

	FormFillInfoHandles[uint32(formInfoStruct)] = formFillInfo

	return &responses.FPDFDOC_InitFormFillEnvironment{
		FormHandle: formHandleHandle.nativeRef,
	}, nil
}

// FPDFDOC_ExitFormFillEnvironment takes ownership of the handle and exits form fill environment.
func (p *PdfiumImplementation) FPDFDOC_ExitFormFillEnvironment(request *requests.FPDFDOC_ExitFormFillEnvironment) (*responses.FPDFDOC_ExitFormFillEnvironment, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDFDOC_ExitFormFillEnvironment").Call(p.Context, *formHandleHandle.handle)
	if err != nil {
		return nil, err
	}

	if _, ok := FormFillInfoHandles[uint32(*formHandleHandle.formInfo)]; ok {
		delete(FormFillInfoHandles, uint32(*formHandleHandle.formInfo))
	}

	delete(p.formHandleRefs, request.FormHandle)

	return &responses.FPDFDOC_ExitFormFillEnvironment{}, nil
}

// FORM_OnAfterLoadPage
// This method is required for implementing all the form related
// functions. Should be invoked after user successfully loaded a
// PDF page, and FPDFDOC_InitFormFillEnvironment() has been invoked.
func (p *PdfiumImplementation) FORM_OnAfterLoadPage(request *requests.FORM_OnAfterLoadPage) (*responses.FORM_OnAfterLoadPage, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	documentHandle, err := p.getDocumentHandle(pageHandle.documentRef)
	if err != nil {
		return nil, err
	}

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FORM_OnAfterLoadPage").Call(p.Context, *pageHandle.handle, *formHandleHandle.handle)
	if err != nil {
		return nil, err
	}

	// Store pointers so that we can reference them in the events to prevent
	// leaving a lot of references to the same page/document.
	formHandleHandle.pagePointers[*pageHandle.handle] = pageHandle.nativeRef
	formHandleHandle.documentPointers[*documentHandle.handle] = documentHandle.nativeRef

	return &responses.FORM_OnAfterLoadPage{}, nil
}

// FORM_OnBeforeClosePage
// This method is required for implementing all the form related
// functions. Should be invoked before user closes the PDF page.
func (p *PdfiumImplementation) FORM_OnBeforeClosePage(request *requests.FORM_OnBeforeClosePage) (*responses.FORM_OnBeforeClosePage, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FORM_OnBeforeClosePage").Call(p.Context, *pageHandle.handle, *formHandleHandle.handle)
	if err != nil {
		return nil, err
	}

	// Remove pointer reference.
	if _, ok := formHandleHandle.pagePointers[*pageHandle.handle]; ok {
		delete(formHandleHandle.pagePointers, *pageHandle.handle)
	}

	return &responses.FORM_OnBeforeClosePage{}, nil
}

// FORM_DoDocumentJSAction
// This method is required for performing document-level JavaScript
// actions. It should be invoked after the PDF document has been loaded.
// If there is document-level JavaScript action embedded in the
// document, this method will execute the JavaScript action. Otherwise,
// the method will do nothing.
func (p *PdfiumImplementation) FORM_DoDocumentJSAction(request *requests.FORM_DoDocumentJSAction) (*responses.FORM_DoDocumentJSAction, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FORM_DoDocumentJSAction").Call(p.Context, *formHandleHandle.handle)
	if err != nil {
		return nil, err
	}

	return &responses.FORM_DoDocumentJSAction{}, nil
}

// FORM_DoDocumentOpenAction
// This method is required for performing open-action when the document
// is opened.
// This method will do nothing if there are no open-actions embedded
// in the document.
func (p *PdfiumImplementation) FORM_DoDocumentOpenAction(request *requests.FORM_DoDocumentOpenAction) (*responses.FORM_DoDocumentOpenAction, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FORM_DoDocumentOpenAction").Call(p.Context, *formHandleHandle.handle)
	if err != nil {
		return nil, err
	}

	return &responses.FORM_DoDocumentOpenAction{}, nil
}

// FORM_DoDocumentAAction
// This method is required for performing the document's
// additional-action.
// This method will do nothing if there is no document
// additional-action corresponding to the specified type.
func (p *PdfiumImplementation) FORM_DoDocumentAAction(request *requests.FORM_DoDocumentAAction) (*responses.FORM_DoDocumentAAction, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FORM_DoDocumentAAction").Call(p.Context, *formHandleHandle.handle, *(*uint64)(unsafe.Pointer(&request.AAType)))
	if err != nil {
		return nil, err
	}

	return &responses.FORM_DoDocumentAAction{}, nil
}

// FORM_DoPageAAction
// This method is required for performing the page object's
// additional-action when opened or closed.
// This method will do nothing if no additional-action corresponding
// to the specified type exists.
func (p *PdfiumImplementation) FORM_DoPageAAction(request *requests.FORM_DoPageAAction) (*responses.FORM_DoPageAAction, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FORM_DoPageAAction").Call(p.Context, *pageHandle.handle, *formHandleHandle.handle, *(*uint64)(unsafe.Pointer(&request.AAType)))
	if err != nil {
		return nil, err
	}

	return &responses.FORM_DoPageAAction{}, nil
}

// FORM_OnMouseMove
// Call this member function when the mouse cursor moves.
func (p *PdfiumImplementation) FORM_OnMouseMove(request *requests.FORM_OnMouseMove) (*responses.FORM_OnMouseMove, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_OnMouseMove").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Modifier)), *(*uint64)(unsafe.Pointer(&request.PageX)), *(*uint64)(unsafe.Pointer(&request.PageY)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not do mouse move")
	}

	return &responses.FORM_OnMouseMove{}, nil
}

// FORM_OnFocus
// This function focuses the form annotation at a given point. If the
// annotation at the point already has focus, nothing happens. If there
// is no annotation at the point, removes form focus.
func (p *PdfiumImplementation) FORM_OnFocus(request *requests.FORM_OnFocus) (*responses.FORM_OnFocus, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_OnFocus").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Modifier)), *(*uint64)(unsafe.Pointer(&request.PageX)), *(*uint64)(unsafe.Pointer(&request.PageY)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FORM_OnFocus{
		HasFocus: int(success) == 1,
	}, nil
}

// FORM_OnLButtonDown
// Call this member function when the user presses the left
// mouse button.
func (p *PdfiumImplementation) FORM_OnLButtonDown(request *requests.FORM_OnLButtonDown) (*responses.FORM_OnLButtonDown, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_OnLButtonDown").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Modifier)), *(*uint64)(unsafe.Pointer(&request.PageX)), *(*uint64)(unsafe.Pointer(&request.PageY)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not do l button down")
	}

	return &responses.FORM_OnLButtonDown{}, nil
}

// FORM_OnRButtonDown
// Call this member function when the user presses the right
// mouse button.
// At the present time, has no effect except in XFA builds, but is
// included for the sake of symmetry.
func (p *PdfiumImplementation) FORM_OnRButtonDown(request *requests.FORM_OnRButtonDown) (*responses.FORM_OnRButtonDown, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_OnRButtonDown").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Modifier)), *(*uint64)(unsafe.Pointer(&request.PageX)), *(*uint64)(unsafe.Pointer(&request.PageY)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not do r button down")
	}

	return &responses.FORM_OnRButtonDown{}, nil
}

// FORM_OnLButtonUp
// Call this member function when the user releases the left
// mouse button.
func (p *PdfiumImplementation) FORM_OnLButtonUp(request *requests.FORM_OnLButtonUp) (*responses.FORM_OnLButtonUp, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_OnLButtonUp").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Modifier)), *(*uint64)(unsafe.Pointer(&request.PageX)), *(*uint64)(unsafe.Pointer(&request.PageY)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not do l button up")
	}

	return &responses.FORM_OnLButtonUp{}, nil
}

// FORM_OnRButtonUp
// Call this member function when the user releases the right
// mouse button.
// At the present time, has no effect except in XFA builds, but is
// included for the sake of symmetry.
func (p *PdfiumImplementation) FORM_OnRButtonUp(request *requests.FORM_OnRButtonUp) (*responses.FORM_OnRButtonUp, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_OnRButtonUp").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Modifier)), *(*uint64)(unsafe.Pointer(&request.PageX)), *(*uint64)(unsafe.Pointer(&request.PageY)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not do r button up")
	}

	return &responses.FORM_OnRButtonUp{}, nil
}

// FORM_OnLButtonDoubleClick
// Call this member function when the user double clicks the
// left mouse button.
func (p *PdfiumImplementation) FORM_OnLButtonDoubleClick(request *requests.FORM_OnLButtonDoubleClick) (*responses.FORM_OnLButtonDoubleClick, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_OnLButtonDoubleClick").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Modifier)), *(*uint64)(unsafe.Pointer(&request.PageX)), *(*uint64)(unsafe.Pointer(&request.PageY)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not do l button double click")
	}

	return &responses.FORM_OnLButtonDoubleClick{}, nil
}

// FORM_OnKeyDown
// Call this member function when a nonsystem key is pressed.
func (p *PdfiumImplementation) FORM_OnKeyDown(request *requests.FORM_OnKeyDown) (*responses.FORM_OnKeyDown, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_OnKeyDown").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.NKeyCode)), *(*uint64)(unsafe.Pointer(&request.Modifier)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not do key down")
	}

	return &responses.FORM_OnKeyDown{}, nil
}

// FORM_OnKeyUp
// Call this member function when a nonsystem key is released.
// Currently unimplemented and always returns false. PDFium reserves this
// API and may implement it in the future on an as-needed basis.
func (p *PdfiumImplementation) FORM_OnKeyUp(request *requests.FORM_OnKeyUp) (*responses.FORM_OnKeyUp, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_OnKeyUp").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.NKeyCode)), *(*uint64)(unsafe.Pointer(&request.Modifier)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not do key up")
	}

	return &responses.FORM_OnKeyUp{}, nil
}

// FORM_OnChar
// Call this member function when a keystroke translates to a
// nonsystem character.
func (p *PdfiumImplementation) FORM_OnChar(request *requests.FORM_OnChar) (*responses.FORM_OnChar, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_OnChar").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.NChar)), *(*uint64)(unsafe.Pointer(&request.Modifier)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not do char")
	}

	return &responses.FORM_OnChar{}, nil
}

// FORM_GetSelectedText
// Call this function to obtain selected text within a form text
// field or form combobox text field.
func (p *PdfiumImplementation) FORM_GetSelectedText(request *requests.FORM_GetSelectedText) (*responses.FORM_GetSelectedText, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	// First get the text length
	res, err := p.Module.ExportedFunction("FORM_GetSelectedText").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	length := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if length == 0 {
		return nil, errors.New("could not get selected text length")
	}

	charDataPointer, err := p.ByteArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	res, err = p.Module.ExportedFunction("FORM_GetSelectedText").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, charDataPointer.Pointer, length)
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

	return &responses.FORM_GetSelectedText{
		SelectedText: transformedText,
	}, nil
}

// FORM_ReplaceSelection
// Call this function to replace the selected text in a form
// text field or user-editable form combobox text field with another
// text string (which can be empty or non-empty). If there is no
// selected text, this function will append the replacement text after
// the current caret position.
func (p *PdfiumImplementation) FORM_ReplaceSelection(request *requests.FORM_ReplaceSelection) (*responses.FORM_ReplaceSelection, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	text, err := p.CFPDF_WIDESTRING(request.Text)
	if err != nil {
		return nil, err
	}
	defer text.Free()

	_, err = p.Module.ExportedFunction("FORM_ReplaceSelection").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, text.Pointer)
	if err != nil {
		return nil, err
	}

	return &responses.FORM_ReplaceSelection{}, nil
}

// FORM_CanUndo
// Find out if it is possible for the current focused widget in a given
// form to perform an undo operation.
func (p *PdfiumImplementation) FORM_CanUndo(request *requests.FORM_CanUndo) (*responses.FORM_CanUndo, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_CanUndo").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	canUndo := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FORM_CanUndo{
		CanUndo: int(canUndo) == 1,
	}, nil
}

// FORM_CanRedo
// Find out if it is possible for the current focused widget in a given
// form to perform a redo operation.
func (p *PdfiumImplementation) FORM_CanRedo(request *requests.FORM_CanRedo) (*responses.FORM_CanRedo, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_CanRedo").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	canRedo := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FORM_CanRedo{
		CanRedo: int(canRedo) == 1,
	}, nil
}

// FORM_Undo
// Make the current focussed widget perform an undo operation.
func (p *PdfiumImplementation) FORM_Undo(request *requests.FORM_Undo) (*responses.FORM_Undo, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_Undo").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not undo")
	}

	return &responses.FORM_Undo{}, nil
}

// FORM_Redo
// Make the current focussed widget perform a redo operation.
func (p *PdfiumImplementation) FORM_Redo(request *requests.FORM_Redo) (*responses.FORM_Redo, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_Redo").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not redo")
	}

	return &responses.FORM_Redo{}, nil
}

// FORM_ForceToKillFocus
// Call this member function to force to kill the focus of the form
// field which has focus. If it would kill the focus of a form field,
// save the value of form field if was changed by theuser.
func (p *PdfiumImplementation) FORM_ForceToKillFocus(request *requests.FORM_ForceToKillFocus) (*responses.FORM_ForceToKillFocus, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_ForceToKillFocus").Call(p.Context, *formHandleHandle.handle)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not kill focus")
	}

	return &responses.FORM_ForceToKillFocus{}, nil
}

// FPDFPage_HasFormFieldAtPoint returns the form field type by point.
func (p *PdfiumImplementation) FPDFPage_HasFormFieldAtPoint(request *requests.FPDFPage_HasFormFieldAtPoint) (*responses.FPDFPage_HasFormFieldAtPoint, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPage_HasFormFieldAtPoint").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.PageX)), *(*uint64)(unsafe.Pointer(&request.PageY)))
	if err != nil {
		return nil, err
	}

	fieldType := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFPage_HasFormFieldAtPoint{
		FieldType: enums.FPDF_FORMFIELD(fieldType),
	}, nil
}

// FPDFPage_FormFieldZOrderAtPoint returns the form field z-order by point.
func (p *PdfiumImplementation) FPDFPage_FormFieldZOrderAtPoint(request *requests.FPDFPage_FormFieldZOrderAtPoint) (*responses.FPDFPage_FormFieldZOrderAtPoint, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDFPage_FormFieldZOrderAtPoint").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.PageX)), *(*uint64)(unsafe.Pointer(&request.PageY)))
	if err != nil {
		return nil, err
	}

	zOrder := *(*int32)(unsafe.Pointer(&res[0]))

	return &responses.FPDFPage_FormFieldZOrderAtPoint{
		ZOrder: int(zOrder),
	}, nil
}

// FPDF_SetFormFieldHighlightColor sets the highlight color of the specified (or all) form fields
// in the document.
func (p *PdfiumImplementation) FPDF_SetFormFieldHighlightColor(request *requests.FPDF_SetFormFieldHighlightColor) (*responses.FPDF_SetFormFieldHighlightColor, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDF_SetFormFieldHighlightColor").Call(p.Context, *formHandleHandle.handle, *(*uint64)(unsafe.Pointer(&request.FieldType)), *(*uint64)(unsafe.Pointer(&request.Color)))
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_SetFormFieldHighlightColor{}, nil
}

// FPDF_SetFormFieldHighlightAlpha sets the transparency of the form field highlight color in the
// document.
func (p *PdfiumImplementation) FPDF_SetFormFieldHighlightAlpha(request *requests.FPDF_SetFormFieldHighlightAlpha) (*responses.FPDF_SetFormFieldHighlightAlpha, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDF_SetFormFieldHighlightAlpha").Call(p.Context, *formHandleHandle.handle, *(*uint64)(unsafe.Pointer(&request.Alpha)))
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_SetFormFieldHighlightAlpha{}, nil
}

// FPDF_RemoveFormFieldHighlight removes the form field highlight color in the document.
func (p *PdfiumImplementation) FPDF_RemoveFormFieldHighlight(request *requests.FPDF_RemoveFormFieldHighlight) (*responses.FPDF_RemoveFormFieldHighlight, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDF_RemoveFormFieldHighlight").Call(p.Context, *formHandleHandle.handle)
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_RemoveFormFieldHighlight{}, nil
}

// FPDF_FFLDraw renders FormFields and popup window on a page to a device independent
// bitmap.
// This function is designed to render annotations that are
// user-interactive, which are widget annotations (for FormFields) and
// popup annotations.
// With the FPDF_ANNOT flag, this function will render a popup annotation
// when users mouse-hover on a non-widget annotation. Regardless of
// FPDF_ANNOT flag, this function will always render widget annotations
// for FormFields.
// In order to implement the FormFill functions, implementation should
// call this function after rendering functions, such as
// FPDF_RenderPageBitmap() or FPDF_RenderPageBitmap_Start(), have
// finished rendering the page contents.
func (p *PdfiumImplementation) FPDF_FFLDraw(request *requests.FPDF_FFLDraw) (*responses.FPDF_FFLDraw, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	bitmapHandle, err := p.getBitmapHandle(request.Bitmap)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	_, err = p.Module.ExportedFunction("FPDF_FFLDraw").Call(p.Context, *formHandleHandle.handle, *bitmapHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.StartX)), *(*uint64)(unsafe.Pointer(&request.StartY)), *(*uint64)(unsafe.Pointer(&request.SizeX)), *(*uint64)(unsafe.Pointer(&request.SizeY)), *(*uint64)(unsafe.Pointer(&request.Rotate)), *(*uint64)(unsafe.Pointer(&request.Flags)))
	if err != nil {
		return nil, err
	}

	return &responses.FPDF_FFLDraw{}, nil
}

// FPDF_LoadXFA load XFA fields of the document if it consists of XFA fields.
func (p *PdfiumImplementation) FPDF_LoadXFA(request *requests.FPDF_LoadXFA) (*responses.FPDF_LoadXFA, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDF_LoadXFA").Call(p.Context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not load XFA")
	}

	return &responses.FPDF_LoadXFA{}, nil
}

// FORM_OnMouseWheel
// Call this member function when the user scrolls the mouse wheel.
// For X and Y delta, the caller must normalize
// platform-specific wheel deltas. e.g. On Windows, a delta value of 240
// for a WM_MOUSEWHEEL event normalizes to 2, since Windows defines
// WHEEL_DELTA as 120.
// Experimental API
func (p *PdfiumImplementation) FORM_OnMouseWheel(request *requests.FORM_OnMouseWheel) (*responses.FORM_OnMouseWheel, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	pageCoordPointer, err := p.FS_POINTFPointer(&request.PageCoord)
	if err != nil {
		return nil, err
	}
	defer pageCoordPointer.Free()

	res, err := p.Module.ExportedFunction("FPDF_LoadXFA").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Modifier)), pageCoordPointer.Pointer, *(*uint64)(unsafe.Pointer(&request.DeltaX)), *(*uint64)(unsafe.Pointer(&request.DeltaY)))
	if err != nil {
		return nil, err
	}

	success := *(*int32)(unsafe.Pointer(&res[0]))
	if int(success) == 0 {
		return nil, errors.New("could not do mouse wheel")
	}

	return &responses.FORM_OnMouseWheel{}, nil
}

// FORM_GetFocusedText
// Call this function to obtain the text within the current focused
// field, if any.
// Experimental API
func (p *PdfiumImplementation) FORM_GetFocusedText(request *requests.FORM_GetFocusedText) (*responses.FORM_GetFocusedText, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	// First get the text length
	res, err := p.Module.ExportedFunction("FORM_GetFocusedText").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, 0, 0)
	if err != nil {
		return nil, err
	}

	length := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if length == 0 {
		return nil, errors.New("could not get focused text length")
	}

	charDataPointer, err := p.ByteArrayPointer(length, nil)
	if err != nil {
		return nil, err
	}
	defer charDataPointer.Free()

	_, err = p.Module.ExportedFunction("FORM_GetFocusedText").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, charDataPointer.Pointer, length)
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

	return &responses.FORM_GetFocusedText{
		FocusedText: transformedText,
	}, nil
}

// FORM_SelectAllText
// Call this function to select all the text within the currently focused
// form text field or form combobox text field.
// Experimental API
func (p *PdfiumImplementation) FORM_SelectAllText(request *requests.FORM_SelectAllText) (*responses.FORM_SelectAllText, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_GetFocusedText").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle)
	if err != nil {
		return nil, err
	}

	success := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if int(success) == 0 {
		return nil, errors.New("could not select all text")
	}

	return &responses.FORM_SelectAllText{}, nil
}

// FORM_GetFocusedAnnot
// Call this member function to get the currently focused annotation.
// Not currently supported for XFA forms - will report no focused
// annotation. Must call FPDFPage_CloseAnnot() when the annotation returned
// by this function is no longer needed.
// Experimental API.
func (p *PdfiumImplementation) FORM_GetFocusedAnnot(request *requests.FORM_GetFocusedAnnot) (*responses.FORM_GetFocusedAnnot, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageIndexPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer pageIndexPointer.Free()

	annotationPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}
	defer annotationPointer.Free()

	res, err := p.Module.ExportedFunction("FORM_GetFocusedAnnot").Call(p.Context, *formHandleHandle.handle, pageIndexPointer.Pointer, annotationPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if int(success) == 0 {
		return nil, errors.New("could not get focused annotation")
	}

	annotation, err := annotationPointer.Value()
	if err != nil {
		return nil, err
	}

	annotationRef := uint64(annotation)
	annotationHandle := p.registerAnnotation(&annotationRef)

	pageIndex, err := pageIndexPointer.Value()
	if err != nil {
		return nil, err
	}

	return &responses.FORM_GetFocusedAnnot{
		PageIndex:  int(pageIndex),
		Annotation: annotationHandle.nativeRef,
	}, nil
}

// FORM_SetFocusedAnnot
// Call this member function to set the currently focused annotation.
// The annotation can't be nil. To kill focus, use FORM_ForceToKillFocus() instead.
// Experimental API.
func (p *PdfiumImplementation) FORM_SetFocusedAnnot(request *requests.FORM_SetFocusedAnnot) (*responses.FORM_SetFocusedAnnot, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	annotationHandle, err := p.getAnnotationHandle(request.Annotation)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_SetFocusedAnnot").Call(p.Context, *formHandleHandle.handle, *annotationHandle.handle)
	if err != nil {
		return nil, err
	}

	success := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if int(success) == 0 {
		return nil, errors.New("could not set focused annotation")
	}

	return &responses.FORM_SetFocusedAnnot{}, nil
}

// FPDF_GetFormType returns the type of form contained in the PDF document.
// If document is nil, then the return value is FORMTYPE_NONE.
// Experimental API
func (p *PdfiumImplementation) FPDF_GetFormType(request *requests.FPDF_GetFormType) (*responses.FPDF_GetFormType, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FPDF_GetFormType").Call(p.Context, *documentHandle.handle)
	if err != nil {
		return nil, err
	}

	formType := uint64(*(*int32)(unsafe.Pointer(&res[0])))

	return &responses.FPDF_GetFormType{
		FormType: enums.FPDF_FORMTYPE(formType),
	}, nil
}

// FORM_SetIndexSelected selects/deselects the value at the given index of the focused
// annotation.
// Intended for use with listbox/combobox widget types. Comboboxes
// have at most a single value selected at a time which cannot be
// deselected. Deselect on a combobox is a no-op that returns false.
// Default implementation is a no-op that will return false for
// other types.
// Not currently supported for XFA forms - will return false.
// Experimental API
func (p *PdfiumImplementation) FORM_SetIndexSelected(request *requests.FORM_SetIndexSelected) (*responses.FORM_SetIndexSelected, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	selected := int64(0)
	if request.Selected {
		selected = int64(1)
	}

	selectedPointer, err := p.IntPointer()
	if err != nil {
		return nil, err
	}

	p.Module.Memory().WriteUint64Le(uint32(selectedPointer.Pointer), api.EncodeI64(selected))
	res, err := p.Module.ExportedFunction("FORM_SetIndexSelected").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)), selectedPointer.Pointer)
	if err != nil {
		return nil, err
	}

	success := uint64(*(*int32)(unsafe.Pointer(&res[0])))
	if int(success) == 0 {
		return nil, errors.New("could not set index selected")
	}

	return &responses.FORM_SetIndexSelected{}, nil
}

// FORM_IsIndexSelected returns whether or not the value at index of the focused
// annotation is currently selected.
// Intended for use with listbox/combobox widget types. Default
// implementation is a no-op that will return false for other types.
// Not currently supported for XFA forms - will return false.
// Experimental API
func (p *PdfiumImplementation) FORM_IsIndexSelected(request *requests.FORM_IsIndexSelected) (*responses.FORM_IsIndexSelected, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	res, err := p.Module.ExportedFunction("FORM_IsIndexSelected").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, *(*uint64)(unsafe.Pointer(&request.Index)))
	if err != nil {
		return nil, err
	}

	isIndexSelected := uint64(*(*int32)(unsafe.Pointer(&res[0])))

	return &responses.FORM_IsIndexSelected{
		IsIndexSelected: int(isIndexSelected) == 1,
	}, nil
}

// FORM_ReplaceAndKeepSelection
// Call this function to replace the selected text in a form text field or
// user-editable form combobox text field with another text string (which
// can be empty or non-empty). If there is no selected text, this function
// will append the replacement text after the current caret position. After
// the insertion, the inserted text will be selected.
// Experimental API
func (p *PdfiumImplementation) FORM_ReplaceAndKeepSelection(request *requests.FORM_ReplaceAndKeepSelection) (*responses.FORM_ReplaceAndKeepSelection, error) {
	p.Lock()
	defer p.Unlock()

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	text, err := p.CFPDF_WIDESTRING(request.Text)
	if err != nil {
		return nil, err
	}
	defer text.Free()

	_, err = p.Module.ExportedFunction("FORM_ReplaceAndKeepSelection").Call(p.Context, *formHandleHandle.handle, *pageHandle.handle, text.Pointer)
	if err != nil {
		return nil, err
	}

	return &responses.FORM_ReplaceAndKeepSelection{}, nil
}
