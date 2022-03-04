package implementation

/*
#cgo pkg-config: pdfium
#include "fpdf_formfill.h"
#include <stdlib.h>

extern void go_formfill_Release_cb(struct _FPDF_FORMFILLINFO *this);
extern void go_formfill_FFI_Invalidate_cb(struct _FPDF_FORMFILLINFO *this, FPDF_PAGE page, double left, double top, double right, double bottom);
extern void go_formfill_FFI_OutputSelectedRect_cb(struct _FPDF_FORMFILLINFO *this, FPDF_PAGE page, double left, double top, double right, double bottom);
extern void go_formfill_FFI_SetCursor_cb(struct _FPDF_FORMFILLINFO *this, int nCursorType);
extern int go_formfill_FFI_SetTimer_cb(struct _FPDF_FORMFILLINFO *this, int uElapse, TimerCallback lpTimerFunc);
extern void go_formfill_FFI_KillTimer_cb(struct _FPDF_FORMFILLINFO *this, int nTimerID);
extern FPDF_SYSTEMTIME go_formfill_FFI_GetLocalTime_cb(struct _FPDF_FORMFILLINFO *this);
extern void go_formfill_FFI_OnChange_cb(struct _FPDF_FORMFILLINFO *this);
extern FPDF_PAGE go_formfill_FFI_GetPage_cb(struct _FPDF_FORMFILLINFO *this, FPDF_DOCUMENT document, int nPageIndex);
extern FPDF_PAGE go_formfill_FFI_GetCurrentPage_cb(struct _FPDF_FORMFILLINFO *this, FPDF_DOCUMENT document);
extern int go_formfill_FFI_GetRotation_cb(struct _FPDF_FORMFILLINFO *this, FPDF_PAGE page);
extern void go_formfill_FFI_ExecuteNamedAction_cb(struct _FPDF_FORMFILLINFO *this, FPDF_BYTESTRING namedAction);
extern void go_formfill_FFI_SetTextFieldFocus_cb(struct _FPDF_FORMFILLINFO *this, FPDF_WIDESTRING value, FPDF_DWORD valueLen, FPDF_BOOL is_focus);
extern void go_formfill_FFI_DoURIAction_cb(struct _FPDF_FORMFILLINFO *this, FPDF_BYTESTRING bsURI);
extern void go_formfill_FFI_DoGoToAction_cb(struct _FPDF_FORMFILLINFO *this, int nPageIndex, int zoomMode, float* fPosArray, int sizeofArray);

static inline void FPDF_FORMFILLINFO_SET_CB(FPDF_FORMFILLINFO *f) {
	f->Release = &go_formfill_Release_cb;
	f->FFI_Invalidate = &go_formfill_FFI_Invalidate_cb;
	f->FFI_OutputSelectedRect = &go_formfill_FFI_OutputSelectedRect_cb;
	f->FFI_SetCursor = &go_formfill_FFI_SetCursor_cb;
	f->FFI_SetTimer = &go_formfill_FFI_SetTimer_cb;
	f->FFI_KillTimer = &go_formfill_FFI_KillTimer_cb;
	f->FFI_GetLocalTime = &go_formfill_FFI_GetLocalTime_cb;
	f->FFI_OnChange = &go_formfill_FFI_OnChange_cb;
	f->FFI_GetPage = &go_formfill_FFI_GetPage_cb;
	f->FFI_GetCurrentPage = &go_formfill_FFI_GetCurrentPage_cb;
	f->FFI_GetRotation = &go_formfill_FFI_GetRotation_cb;
	f->FFI_ExecuteNamedAction = &go_formfill_FFI_ExecuteNamedAction_cb;
	f->FFI_SetTextFieldFocus = &go_formfill_FFI_SetTextFieldFocus_cb;
	f->FFI_DoURIAction = &go_formfill_FFI_DoURIAction_cb;
	f->FFI_DoGoToAction = &go_formfill_FFI_DoGoToAction_cb;
}

*/
import "C"
import (
	"errors"
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"
)

// go_formfill_Release_cb is the Go implementation of FPDF_FORMFILLINFO::Release.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//export go_formfill_Release_cb
func go_formfill_Release_cb(me *C.FPDF_FORMFILLINFO) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	// @todo: do I have anything to cleanup for myself?

	if formFillInfoHandle.FormFillInfo.Release != nil {
		formFillInfoHandle.FormFillInfo.Release()
	}
}

// go_formfill_FFI_Invalidate_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_Invalidate.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//export go_formfill_FFI_Invalidate_cb
func go_formfill_FFI_Invalidate_cb(me *C.FPDF_FORMFILLINFO, page C.FPDF_PAGE, left C.double, top C.double, right C.double, bottom C.double) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	// @todo: is this the best way? Maybe handles should be based on the pointer to prevent duplicate handles.
	pageHandle := formFillInfoHandle.Instance.registerPage(page, 0, nil)
	formFillInfoHandle.FormFillInfo.FFI_Invalidate(pageHandle.nativeRef, float64(left), float64(top), float64(right), float64(bottom))
}

// go_formfill_FFI_OutputSelectedRect_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_OutputSelectedRect.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//export go_formfill_FFI_OutputSelectedRect_cb
func go_formfill_FFI_OutputSelectedRect_cb(me *C.FPDF_FORMFILLINFO, page C.FPDF_PAGE, left C.double, top C.double, right C.double, bottom C.double) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_OutputSelectedRect != nil {
		// @todo: is this the best way? Maybe handles should be based on the pointer to prevent duplicate handles.
		pageHandle := formFillInfoHandle.Instance.registerPage(page, 0, nil)
		formFillInfoHandle.FormFillInfo.FFI_OutputSelectedRect(pageHandle.nativeRef, float64(left), float64(top), float64(right), float64(bottom))
	}
}

// go_formfill_FFI_SetCursor_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_SetCursor.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//export go_formfill_FFI_SetCursor_cb
func go_formfill_FFI_SetCursor_cb(me *C.FPDF_FORMFILLINFO, nCursorType C.int) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	formFillInfoHandle.FormFillInfo.FFI_SetCursor(enums.FXCT(nCursorType))
}

// go_formfill_FFI_SetTimer_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_SetTimer.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//export go_formfill_FFI_SetTimer_cb
func go_formfill_FFI_SetTimer_cb(me *C.FPDF_FORMFILLINFO, uElapse C.int, lpTimerFunc C.TimerCallback) C.int {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return 0
	}

	timerFunc := func(idEvent int) {
		// @todo: implement TimerCallback
	}

	timerID := formFillInfoHandle.FormFillInfo.FFI_SetTimer(int(uElapse), timerFunc)
	return C.int(timerID)
}

// go_formfill_FFI_KillTimer_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_KillTimer.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//export go_formfill_FFI_KillTimer_cb
func go_formfill_FFI_KillTimer_cb(me *C.FPDF_FORMFILLINFO, nTimerID C.int) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	formFillInfoHandle.FormFillInfo.FFI_KillTimer(int(nTimerID))
}

// go_formfill_FFI_GetLocalTime_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_GetLocalTime.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//export go_formfill_FFI_GetLocalTime_cb
func go_formfill_FFI_GetLocalTime_cb(me *C.FPDF_FORMFILLINFO) C.FPDF_SYSTEMTIME {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return C.FPDF_SYSTEMTIME{}
	}

	localTime := formFillInfoHandle.FormFillInfo.FFI_GetLocalTime()

	return C.FPDF_SYSTEMTIME{
		wYear:         C.ushort(localTime.Year),
		wMonth:        C.ushort(localTime.Month),
		wDayOfWeek:    C.ushort(localTime.DayOfWeek),
		wDay:          C.ushort(localTime.Day),
		wHour:         C.ushort(localTime.Hour),
		wMinute:       C.ushort(localTime.Minute),
		wSecond:       C.ushort(localTime.Second),
		wMilliseconds: C.ushort(localTime.Milliseconds),
	}
}

// go_formfill_FFI_OnChange_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_OnChange.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//export go_formfill_FFI_OnChange_cb
func go_formfill_FFI_OnChange_cb(me *C.FPDF_FORMFILLINFO) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_OnChange != nil {
		formFillInfoHandle.FormFillInfo.FFI_OnChange()
	}
}

// go_formfill_FFI_GetPage_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_GetPage.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//export go_formfill_FFI_GetPage_cb
func go_formfill_FFI_GetPage_cb(me *C.FPDF_FORMFILLINFO, document C.FPDF_DOCUMENT, nPageIndex C.int) C.FPDF_PAGE {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return nil
	}

	documentHandle := formFillInfoHandle.Instance.registerDocument(document)

	page := formFillInfoHandle.FormFillInfo.FFI_GetPage(documentHandle.nativeRef, int(nPageIndex))
	if page == nil {
		return nil
	}

	pageHandle, err := formFillInfoHandle.Instance.getPageHandle(*page)
	if err != nil {
		return nil
	}

	return pageHandle.handle
}

// go_formfill_FFI_GetCurrentPage_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_GetCurrentPage.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//export go_formfill_FFI_GetCurrentPage_cb
func go_formfill_FFI_GetCurrentPage_cb(me *C.FPDF_FORMFILLINFO, document C.FPDF_DOCUMENT) C.FPDF_PAGE {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return nil
	}

	if formFillInfoHandle.FormFillInfo.FFI_GetCurrentPage == nil {
		return nil
	}

	documentHandle := formFillInfoHandle.Instance.registerDocument(document)
	page := formFillInfoHandle.FormFillInfo.FFI_GetCurrentPage(documentHandle.nativeRef)
	if page == nil {
		return nil
	}

	pageHandle, err := formFillInfoHandle.Instance.getPageHandle(*page)
	if err != nil {
		return nil
	}

	return pageHandle.handle
}

// go_formfill_FFI_GetRotation_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_GetRotation.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//export go_formfill_FFI_GetRotation_cb
func go_formfill_FFI_GetRotation_cb(me *C.FPDF_FORMFILLINFO, page C.FPDF_PAGE) C.int {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return C.int(0)
	}

	// @todo: is this the best way? Maybe handles should be based on the pointer to prevent duplicate handles.
	pageHandle := formFillInfoHandle.Instance.registerPage(page, 0, nil)
	rotation := formFillInfoHandle.FormFillInfo.FFI_GetRotation(pageHandle.nativeRef)

	return C.int(rotation)
}

// go_formfill_FFI_ExecuteNamedAction_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_ExecuteNamedAction.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//export go_formfill_FFI_ExecuteNamedAction_cb
func go_formfill_FFI_ExecuteNamedAction_cb(me *C.FPDF_FORMFILLINFO, namedAction C.FPDF_BYTESTRING) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	formFillInfoHandle.FormFillInfo.FFI_ExecuteNamedAction(C.GoString((*C.char)(namedAction)))

	return
}

// go_formfill_FFI_SetTextFieldFocus_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_SetTextFieldFocus.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//export go_formfill_FFI_SetTextFieldFocus_cb
func go_formfill_FFI_SetTextFieldFocus_cb(me *C.FPDF_FORMFILLINFO, value C.FPDF_WIDESTRING, valueLen C.FPDF_DWORD, is_focus C.FPDF_BOOL) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_SetTextFieldFocus == nil {
		return
	}

	size := uint64(valueLen) * 2
	target := (*[1<<50 - 1]byte)(unsafe.Pointer(value))[:size:size]
	decodedValue, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(target)
	if err != nil {
		return
	}

	formFillInfoHandle.FormFillInfo.FFI_SetTextFieldFocus(decodedValue, int(is_focus) == 1)

	return
}

// go_formfill_FFI_DoURIAction_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_DoURIAction.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//export go_formfill_FFI_DoURIAction_cb
func go_formfill_FFI_DoURIAction_cb(me *C.FPDF_FORMFILLINFO, bsURI C.FPDF_BYTESTRING) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_DoURIAction == nil {
		return
	}

	formFillInfoHandle.FormFillInfo.FFI_DoURIAction(C.GoString((*C.char)(bsURI)))

	return
}

// go_formfill_FFI_DoGoToAction_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_DoGoToAction.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//export go_formfill_FFI_DoGoToAction_cb
func go_formfill_FFI_DoGoToAction_cb(me *C.FPDF_FORMFILLINFO, nPageIndex C.int, zoomMode C.int, fPosArray *C.float, sizeofArray C.int) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_DoGoToAction == nil {
		return
	}

	target := (*[1<<25 - 1]float32)(unsafe.Pointer(fPosArray))[:sizeofArray:sizeofArray]
	pos := make([]float32, int(sizeofArray))
	for i := range pos {
		pos[i] = float32(target[i])
	}

	formFillInfoHandle.FormFillInfo.FFI_DoGoToAction(int(nPageIndex), enums.FPDF_ZOOM_MODE(zoomMode), pos)

	return
}

type FormFillInfo struct {
	Struct       *C.FPDF_FORMFILLINFO
	FormFillInfo *structs.FPDF_FORMFILLINFO
	NativeRef    references.FPDF_FORMHANDLE
	Instance     *PdfiumImplementation
}

var formFillInfoHandles = map[unsafe.Pointer]*FormFillInfo{}

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

	formInfoStruct := &C.FPDF_FORMFILLINFO{}
	formInfoStruct.version = 1
	C.FPDF_FORMFILLINFO_SET_CB(formInfoStruct)

	formHandle := C.FPDFDOC_InitFormFillEnvironment(documentHandle.handle, formInfoStruct)
	if formHandle == nil {
		return nil, errors.New("could not init form fill environment")
	}

	formHandleHandle := p.registerFormHandle(formHandle, unsafe.Pointer(formInfoStruct))

	formFillInfo := &FormFillInfo{
		Struct:       formInfoStruct,
		FormFillInfo: &request.FormFillInfo,
		NativeRef:    formHandleHandle.nativeRef,
		Instance:     p,
	}

	formFillInfoHandles[unsafe.Pointer(formInfoStruct)] = formFillInfo

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

	C.FPDFDOC_ExitFormFillEnvironment(formHandleHandle.handle)

	if _, ok := formFillInfoHandles[formHandleHandle.formInfo]; ok {
		delete(formFillInfoHandles, formHandleHandle.formInfo)
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

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	C.FORM_OnAfterLoadPage(pageHandle.handle, formHandleHandle.handle)

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

	C.FORM_OnBeforeClosePage(pageHandle.handle, formHandleHandle.handle)

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

	C.FORM_DoDocumentJSAction(formHandleHandle.handle)

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

	C.FORM_DoDocumentOpenAction(formHandleHandle.handle)

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

	C.FORM_DoDocumentAAction(formHandleHandle.handle, C.int(request.AAType))

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

	C.FORM_DoPageAAction(pageHandle.handle, formHandleHandle.handle, C.int(request.AAType))

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

	success := C.FORM_OnMouseMove(formHandleHandle.handle, pageHandle.handle, C.int(request.Modifier), C.double(request.PageX), C.double(request.PageY))
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

	success := C.FORM_OnFocus(formHandleHandle.handle, pageHandle.handle, C.int(request.Modifier), C.double(request.PageX), C.double(request.PageY))
	if int(success) == 0 {
		return nil, errors.New("could not do focus")
	}

	return &responses.FORM_OnFocus{}, nil
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

	success := C.FORM_OnLButtonDown(formHandleHandle.handle, pageHandle.handle, C.int(request.Modifier), C.double(request.PageX), C.double(request.PageY))
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

	success := C.FORM_OnRButtonDown(formHandleHandle.handle, pageHandle.handle, C.int(request.Modifier), C.double(request.PageX), C.double(request.PageY))
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

	success := C.FORM_OnLButtonUp(formHandleHandle.handle, pageHandle.handle, C.int(request.Modifier), C.double(request.PageX), C.double(request.PageY))
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

	success := C.FORM_OnRButtonUp(formHandleHandle.handle, pageHandle.handle, C.int(request.Modifier), C.double(request.PageX), C.double(request.PageY))
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

	success := C.FORM_OnLButtonDoubleClick(formHandleHandle.handle, pageHandle.handle, C.int(request.Modifier), C.double(request.PageX), C.double(request.PageY))
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

	success := C.FORM_OnKeyDown(formHandleHandle.handle, pageHandle.handle, C.int(request.NKeyCode), C.int(request.Modifier))
	if int(success) == 0 {
		return nil, errors.New("could not do key down")
	}

	return &responses.FORM_OnKeyDown{}, nil
}

// FORM_OnKeyUp
// Call this member function when a nonsystem key is released.
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

	success := C.FORM_OnKeyUp(formHandleHandle.handle, pageHandle.handle, C.int(request.NKeyCode), C.int(request.Modifier))
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

	success := C.FORM_OnChar(formHandleHandle.handle, pageHandle.handle, C.int(request.NChar), C.int(request.Modifier))
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
	length := C.FORM_GetSelectedText(formHandleHandle.handle, pageHandle.handle, nil, 0)
	if uint64(length) == 0 {
		return nil, errors.New("could not get selected text length")
	}

	charData := make([]byte, length)
	C.FORM_GetSelectedText(formHandleHandle.handle, pageHandle.handle, unsafe.Pointer(&charData[0]), C.ulong(len(charData)))

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

	transformedText, err := p.transformUTF8ToUTF16LE(request.Text)
	if err != nil {
		return nil, err
	}

	C.FORM_ReplaceSelection(formHandleHandle.handle, pageHandle.handle, (C.FPDF_WIDESTRING)(unsafe.Pointer(&transformedText[0])))

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

	canUndo := C.FORM_CanUndo(formHandleHandle.handle, pageHandle.handle)

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

	canRedo := C.FORM_CanRedo(formHandleHandle.handle, pageHandle.handle)

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

	success := C.FORM_Undo(formHandleHandle.handle, pageHandle.handle)
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

	success := C.FORM_Redo(formHandleHandle.handle, pageHandle.handle)
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

	success := C.FORM_ForceToKillFocus(formHandleHandle.handle)
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

	fieldType := C.FPDFPage_HasFormFieldAtPoint(formHandleHandle.handle, pageHandle.handle, C.double(request.PageX), C.double(request.PageY))

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

	zOrder := C.FPDFPage_FormFieldZOrderAtPoint(formHandleHandle.handle, pageHandle.handle, C.double(request.PageX), C.double(request.PageY))

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

	C.FPDF_SetFormFieldHighlightColor(formHandleHandle.handle, C.int(request.FieldType), C.ulong(request.Color))

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

	C.FPDF_SetFormFieldHighlightAlpha(formHandleHandle.handle, C.uchar(request.Alpha))

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

	C.FPDF_RemoveFormFieldHighlight(formHandleHandle.handle)

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

	C.FPDF_FFLDraw(formHandleHandle.handle, bitmapHandle.handle, pageHandle.handle, C.int(request.StartX), C.int(request.StartY), C.int(request.SizeX), C.int(request.SizeY), C.int(request.Rotate), C.int(request.Flags))

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

	success := C.FPDF_LoadXFA(documentHandle.handle)
	if int(success) == 0 {
		return nil, errors.New("could not load XFA")
	}

	return &responses.FPDF_LoadXFA{}, nil
}
