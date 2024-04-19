package implementation_cgo

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

// XFA methods.
extern void go_formfill_FFI_DisplayCaret_cb(struct _FPDF_FORMFILLINFO *this, FPDF_PAGE page, FPDF_BOOL bVisible, double left, double top, double right, double bottom);
extern int go_formfill_FFI_GetCurrentPageIndex_cb(struct _FPDF_FORMFILLINFO *this, FPDF_DOCUMENT document);
extern void go_formfill_FFI_SetCurrentPage_cb(struct _FPDF_FORMFILLINFO *this, FPDF_DOCUMENT document, int iCurPage);
extern void go_formfill_FFI_GotoURL_cb(struct _FPDF_FORMFILLINFO *this, FPDF_DOCUMENT document, FPDF_WIDESTRING wsURL);
extern void go_formfill_FFI_GetPageViewRect_cb(struct _FPDF_FORMFILLINFO *this, FPDF_PAGE page, double* left, double* top, double* right, double* bottom);
extern void go_formfill_FFI_PageEvent_cb(struct _FPDF_FORMFILLINFO *this, int page_count, FPDF_DWORD event_type);
extern FPDF_BOOL go_formfill_FFI_PopupMenu_cb(struct _FPDF_FORMFILLINFO *this, FPDF_PAGE page, FPDF_WIDGET hWidget, int menuFlag, float x, float y);
typedef const char* ccharp;
extern FPDF_FILEHANDLER* go_formfill_FFI_OpenFile_cb(struct _FPDF_FORMFILLINFO *this, int fileFlag, FPDF_WIDESTRING wsURL, ccharp mode);
extern void go_formfill_FFI_EmailTo_cb(struct _FPDF_FORMFILLINFO *this, FPDF_FILEHANDLER* fileHandler, FPDF_WIDESTRING pTo, FPDF_WIDESTRING pSubject, FPDF_WIDESTRING pCC, FPDF_WIDESTRING pBcc, FPDF_WIDESTRING pMsg);
extern void go_formfill_FFI_UploadTo_cb(struct _FPDF_FORMFILLINFO *this, FPDF_FILEHANDLER* fileHandler, int fileFlag, FPDF_WIDESTRING uploadTo);
extern int go_formfill_FFI_GetPlatform_cb(struct _FPDF_FORMFILLINFO *this, void* platform, int length);
extern int go_formfill_FFI_GetLanguage_cb(struct _FPDF_FORMFILLINFO *this, void* language, int length);
extern FPDF_FILEHANDLER* go_formfill_FFI_DownloadFromURL_cb(struct _FPDF_FORMFILLINFO *this, FPDF_WIDESTRING URL);
extern FPDF_BOOL go_formfill_FFI_PostRequestURL_cb(struct _FPDF_FORMFILLINFO *this, FPDF_WIDESTRING wsURL, FPDF_WIDESTRING wsData, FPDF_WIDESTRING wsContentType, FPDF_WIDESTRING wsEncode, FPDF_WIDESTRING wsHeader, FPDF_BSTR* response);
extern FPDF_BOOL go_formfill_FFI_PutRequestURL_cb(struct _FPDF_FORMFILLINFO *this, FPDF_WIDESTRING wsURL, FPDF_WIDESTRING wsData, FPDF_WIDESTRING wsEncode);
extern void go_formfill_FFI_OnFocusChange_cb(struct _FPDF_FORMFILLINFO *this, FPDF_ANNOTATION annot, int page_index);
extern void go_formfill_FFI_DoURIActionWithKeyboardModifier_cb(struct _FPDF_FORMFILLINFO *this, FPDF_BYTESTRING uri, int modifiers);

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

    // XFA methods.
	if (f->version > 1) {
     	f->FFI_DisplayCaret = &go_formfill_FFI_DisplayCaret_cb;
     	f->FFI_GetCurrentPageIndex = &go_formfill_FFI_GetCurrentPageIndex_cb;
     	f->FFI_SetCurrentPage = &go_formfill_FFI_SetCurrentPage_cb;
     	f->FFI_GotoURL = &go_formfill_FFI_GotoURL_cb;
     	f->FFI_GetPageViewRect = &go_formfill_FFI_GetPageViewRect_cb;
     	f->FFI_PageEvent = &go_formfill_FFI_PageEvent_cb;
     	f->FFI_PopupMenu = &go_formfill_FFI_PopupMenu_cb;
     	f->FFI_OpenFile = &go_formfill_FFI_OpenFile_cb;
     	f->FFI_EmailTo = &go_formfill_FFI_EmailTo_cb;
     	f->FFI_UploadTo = &go_formfill_FFI_UploadTo_cb;
     	f->FFI_GetPlatform = &go_formfill_FFI_GetPlatform_cb;
     	f->FFI_GetLanguage = &go_formfill_FFI_GetLanguage_cb;
     	f->FFI_DownloadFromURL = &go_formfill_FFI_DownloadFromURL_cb;
     	f->FFI_PostRequestURL = &go_formfill_FFI_PostRequestURL_cb;
     	f->FFI_PutRequestURL = &go_formfill_FFI_PutRequestURL_cb;
     	f->FFI_OnFocusChange = &go_formfill_FFI_OnFocusChange_cb;
     	f->FFI_DoURIActionWithKeyboardModifier = &go_formfill_FFI_DoURIActionWithKeyboardModifier_cb;
		IPDF_JSPLATFORM *jsPlatform;
		jsPlatform = malloc(sizeof(IPDF_JSPLATFORM));
		jsPlatform->version = 3;
		jsPlatform->m_pFormfillinfo = f;
		f->m_pJsPlatform = jsPlatform;
    }
}

static inline void FPDF_FORMFILLINFO_CALL_TIMER(TimerCallback t, int id) {
	t(id);
}

extern int go_jsplatform_app_alert_cb(struct _IPDF_JsPlatform *this, FPDF_WIDESTRING Msg, FPDF_WIDESTRING Title, int Type, int Icon);
extern void go_jsplatform_app_beep_cb(struct _IPDF_JsPlatform *this, int nType);
extern int go_jsplatform_app_response_cb(struct _IPDF_JsPlatform *this, FPDF_WIDESTRING Question, FPDF_WIDESTRING Title, FPDF_WIDESTRING Default, FPDF_WIDESTRING cLabel, FPDF_BOOL bPassword, void* response, int length);
extern int go_jsplatform_Doc_getFilePath_cb(struct _IPDF_JsPlatform *this, void* filePath, int length);
extern void go_jsplatform_Doc_mail_cb(struct _IPDF_JsPlatform *this, void* mailData, int length, FPDF_BOOL bUI, FPDF_WIDESTRING To, FPDF_WIDESTRING Subject, FPDF_WIDESTRING CC, FPDF_WIDESTRING BCC, FPDF_WIDESTRING Msg);
extern void go_jsplatform_Doc_print_cb(struct _IPDF_JsPlatform *this, FPDF_BOOL bUI, int nStart, int nEnd, FPDF_BOOL bSilent, FPDF_BOOL bShrinkToFit, FPDF_BOOL bPrintAsImage, FPDF_BOOL bReverse, FPDF_BOOL bAnnotations);
extern void go_jsplatform_Doc_submitForm_cb(struct _IPDF_JsPlatform *this, void* formData, int length, FPDF_WIDESTRING URL);
extern void go_jsplatform_Doc_gotoPage_cb(struct _IPDF_JsPlatform *this, int nPageNum);
extern int go_jsplatform_Field_browse_cb(struct _IPDF_JsPlatform *this, void* filePath, int length);

static inline void IPDF_JSPLATFORM_SET_CB(IPDF_JSPLATFORM *jsP) {
     	jsP->app_alert = &go_jsplatform_app_alert_cb;
     	jsP->app_beep = &go_jsplatform_app_beep_cb;
     	jsP->app_response = &go_jsplatform_app_response_cb;
     	jsP->Doc_getFilePath = &go_jsplatform_Doc_getFilePath_cb;
     	jsP->Doc_mail = &go_jsplatform_Doc_mail_cb;
     	jsP->Doc_print = &go_jsplatform_Doc_print_cb;
     	jsP->Doc_submitForm = &go_jsplatform_Doc_submitForm_cb;
     	jsP->Doc_gotoPage = &go_jsplatform_Doc_gotoPage_cb;
     	jsP->Field_browse = &go_jsplatform_Field_browse_cb;
}

extern void go_filehandler_Release_cb(void* clientData);
extern FPDF_DWORD go_filehandler_GetSize_cb(void* clientData);
extern FPDF_RESULT go_filehandler_ReadBlock_cb(void* clientData, FPDF_DWORD offset, void* buffer, FPDF_DWORD size);
typedef const void* cvoidp;
extern FPDF_RESULT go_filehandler_WriteBlock_cb(void* clientData, FPDF_DWORD offset, cvoidp buffer, FPDF_DWORD size);
extern FPDF_RESULT go_filehandler_Flush_cb(void* clientData);
extern FPDF_RESULT go_filehandler_Truncate_cb(void* clientData, FPDF_DWORD size);

static inline void FPDF_FILEHANDLER_SET_CB(FPDF_FILEHANDLER *f, char *id) {
    f->clientData = id;
	f->Release = &go_filehandler_Release_cb;
	f->GetSize = &go_filehandler_GetSize_cb;
	f->ReadBlock = &go_filehandler_ReadBlock_cb;
	f->WriteBlock = &go_filehandler_WriteBlock_cb;
	f->Flush = &go_filehandler_Flush_cb;
	f->Truncate = &go_filehandler_Truncate_cb;
}

*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"

	"github.com/google/uuid"
)

// go_formfill_Release_cb is the Go implementation of FPDF_FORMFILLINFO::Release.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_Release_cb
func go_formfill_Release_cb(me *C.FPDF_FORMFILLINFO) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.Release != nil {
		formFillInfoHandle.FormFillInfo.Release()
	}

	delete(formFillInfoHandles, pointer)
}

// go_formfill_FFI_Invalidate_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_Invalidate.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_Invalidate_cb
func go_formfill_FFI_Invalidate_cb(me *C.FPDF_FORMFILLINFO, page C.FPDF_PAGE, left C.double, top C.double, right C.double, bottom C.double) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	var pageRef references.FPDF_PAGE
	if pointerPageRef, ok := formFillInfoHandle.FormHandleHandle.pagePointers[unsafe.Pointer(page)]; ok {
		pageRef = pointerPageRef
	} else {
		pageHandle := formFillInfoHandle.Instance.registerPage(page, 0, nil)
		pageRef = pageHandle.nativeRef
	}

	formFillInfoHandle.FormFillInfo.FFI_Invalidate(pageRef, float64(left), float64(top), float64(right), float64(bottom))
}

// go_formfill_FFI_OutputSelectedRect_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_OutputSelectedRect.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_OutputSelectedRect_cb
func go_formfill_FFI_OutputSelectedRect_cb(me *C.FPDF_FORMFILLINFO, page C.FPDF_PAGE, left C.double, top C.double, right C.double, bottom C.double) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_OutputSelectedRect != nil {
		var pageRef references.FPDF_PAGE
		if pointerPageRef, ok := formFillInfoHandle.FormHandleHandle.pagePointers[unsafe.Pointer(page)]; ok {
			pageRef = pointerPageRef
		} else {
			pageHandle := formFillInfoHandle.Instance.registerPage(page, 0, nil)
			pageRef = pageHandle.nativeRef
		}

		formFillInfoHandle.FormFillInfo.FFI_OutputSelectedRect(pageRef, float64(left), float64(top), float64(right), float64(bottom))
	}
}

// go_formfill_FFI_SetCursor_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_SetCursor.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
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
//
//export go_formfill_FFI_SetTimer_cb
func go_formfill_FFI_SetTimer_cb(me *C.FPDF_FORMFILLINFO, uElapse C.int, lpTimerFunc C.TimerCallback) C.int {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return 0
	}

	timerFunc := func(idEvent int) {
		C.FPDF_FORMFILLINFO_CALL_TIMER(lpTimerFunc, C.int(idEvent))
	}

	timerID := formFillInfoHandle.FormFillInfo.FFI_SetTimer(int(uElapse), timerFunc)
	return C.int(timerID)
}

// go_formfill_FFI_KillTimer_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_KillTimer.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
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
//
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
//
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
//
//export go_formfill_FFI_GetPage_cb
func go_formfill_FFI_GetPage_cb(me *C.FPDF_FORMFILLINFO, document C.FPDF_DOCUMENT, nPageIndex C.int) C.FPDF_PAGE {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return nil
	}

	var documentRef references.FPDF_DOCUMENT
	if pointerDocumentRef, ok := formFillInfoHandle.FormHandleHandle.documentPointers[unsafe.Pointer(document)]; ok {
		documentRef = pointerDocumentRef
	} else {
		documentHandle := formFillInfoHandle.Instance.registerDocument(document)
		documentRef = documentHandle.nativeRef
	}

	page := formFillInfoHandle.FormFillInfo.FFI_GetPage(documentRef, int(nPageIndex))
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
//
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
//
//export go_formfill_FFI_GetRotation_cb
func go_formfill_FFI_GetRotation_cb(me *C.FPDF_FORMFILLINFO, page C.FPDF_PAGE) C.int {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return C.int(0)
	}

	var pageRef references.FPDF_PAGE
	if pointerPageRef, ok := formFillInfoHandle.FormHandleHandle.pagePointers[unsafe.Pointer(page)]; ok {
		pageRef = pointerPageRef
	} else {
		pageHandle := formFillInfoHandle.Instance.registerPage(page, 0, nil)
		pageRef = pageHandle.nativeRef
	}

	rotation := formFillInfoHandle.FormFillInfo.FFI_GetRotation(pageRef)

	return C.int(rotation)
}

// go_formfill_FFI_ExecuteNamedAction_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_ExecuteNamedAction.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
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
//
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
	// We create a Go slice backed by a C array (without copying the original data).
	target := unsafe.Slice((*byte)(unsafe.Pointer(value)), uint64(size))

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
//
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
//
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

// XFA methods

// go_formfill_FFI_DisplayCaret_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_DisplayCaret.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_DisplayCaret_cb
func go_formfill_FFI_DisplayCaret_cb(me *C.FPDF_FORMFILLINFO, page C.FPDF_PAGE, bVisible C.FPDF_BOOL, left, top, right, bottom C.double) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_DisplayCaret == nil {
		return
	}

	var pageRef references.FPDF_PAGE
	if pointerPageRef, ok := formFillInfoHandle.FormHandleHandle.pagePointers[unsafe.Pointer(page)]; ok {
		pageRef = pointerPageRef
	} else {
		pageHandle := formFillInfoHandle.Instance.registerPage(page, 0, nil)
		pageRef = pageHandle.nativeRef
	}

	formFillInfoHandle.FormFillInfo.FFI_DisplayCaret(pageRef, int(bVisible) == 1, float64(left), float64(top), float64(right), float64(bottom))

	return
}

// go_formfill_FFI_GetCurrentPageIndex_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_GetCurrentPageIndex.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_GetCurrentPageIndex_cb
func go_formfill_FFI_GetCurrentPageIndex_cb(me *C.FPDF_FORMFILLINFO, document C.FPDF_DOCUMENT) C.int {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return C.int(0)
	}

	if formFillInfoHandle.FormFillInfo.FFI_GetCurrentPageIndex == nil {
		return C.int(0)
	}

	var docRef references.FPDF_DOCUMENT
	if pointerDocRef, ok := formFillInfoHandle.FormHandleHandle.documentPointers[unsafe.Pointer(document)]; ok {
		docRef = pointerDocRef
	} else {
		docHandle := formFillInfoHandle.Instance.registerDocument(document)
		docRef = docHandle.nativeRef
	}

	return C.int(formFillInfoHandle.FormFillInfo.FFI_GetCurrentPageIndex(docRef))
}

// go_formfill_FFI_SetCurrentPage_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_SetCurrentPage.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_SetCurrentPage_cb
func go_formfill_FFI_SetCurrentPage_cb(me *C.FPDF_FORMFILLINFO, document C.FPDF_DOCUMENT, iCurPage C.int) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_SetCurrentPage == nil {
		return
	}

	var docRef references.FPDF_DOCUMENT
	if pointerDocRef, ok := formFillInfoHandle.FormHandleHandle.documentPointers[unsafe.Pointer(document)]; ok {
		docRef = pointerDocRef
	} else {
		docHandle := formFillInfoHandle.Instance.registerDocument(document)
		docRef = docHandle.nativeRef
	}

	formFillInfoHandle.FormFillInfo.FFI_SetCurrentPage(docRef, int(iCurPage))

	return
}

// go_formfill_FFI_GotoURL_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_GotoURL.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_GotoURL_cb
func go_formfill_FFI_GotoURL_cb(me *C.FPDF_FORMFILLINFO, document C.FPDF_DOCUMENT, wsURL C.FPDF_WIDESTRING) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_GotoURL == nil {
		return
	}

	var docRef references.FPDF_DOCUMENT
	if pointerDocRef, ok := formFillInfoHandle.FormHandleHandle.documentPointers[unsafe.Pointer(document)]; ok {
		docRef = pointerDocRef
	} else {
		docHandle := formFillInfoHandle.Instance.registerDocument(document)
		docRef = docHandle.nativeRef
	}

	url := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(wsURL))
	decodedValue, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(url)
	if err != nil {
		return
	}

	formFillInfoHandle.FormFillInfo.FFI_GotoURL(docRef, decodedValue)

	return
}

// go_formfill_FFI_GetPageViewRect_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_GetPageViewRect.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_GetPageViewRect_cb
func go_formfill_FFI_GetPageViewRect_cb(me *C.FPDF_FORMFILLINFO, page C.FPDF_PAGE, left, top, right, bottom *C.double) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_GetPageViewRect == nil {
		return
	}

	var pageRef references.FPDF_PAGE
	if pointerPageRef, ok := formFillInfoHandle.FormHandleHandle.pagePointers[unsafe.Pointer(page)]; ok {
		pageRef = pointerPageRef
	} else {
		pageHandle := formFillInfoHandle.Instance.registerPage(page, 0, nil)
		pageRef = pageHandle.nativeRef
	}

	pageViewRectLeft, pageViewRectTop, pageViewRectRight, pageViewRectBottom := formFillInfoHandle.FormFillInfo.FFI_GetPageViewRect(pageRef)
	cgoPageViewRectLeft := C.double(pageViewRectLeft)
	cgoPageViewRectTop := C.double(pageViewRectTop)
	cgoPageViewRectRight := C.double(pageViewRectRight)
	cgoPageViewRectBottom := C.double(pageViewRectBottom)

	left = &cgoPageViewRectLeft
	top = &cgoPageViewRectTop
	right = &cgoPageViewRectRight
	bottom = &cgoPageViewRectBottom

	return
}

// go_formfill_FFI_PageEvent_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_PageEvent.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_PageEvent_cb
func go_formfill_FFI_PageEvent_cb(me *C.FPDF_FORMFILLINFO, page_count C.int, event_type C.FPDF_DWORD) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_PageEvent == nil {
		return
	}

	formFillInfoHandle.FormFillInfo.FFI_PageEvent(int(page_count), enums.FXFA_PAGEVIEWEVENT(event_type))

	return
}

// go_formfill_FFI_PopupMenu_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_PopupMenu.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_PopupMenu_cb
func go_formfill_FFI_PopupMenu_cb(me *C.FPDF_FORMFILLINFO, page C.FPDF_PAGE, hWidget C.FPDF_WIDGET, menuFlag C.int, x, y C.float) C.FPDF_BOOL {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return C.FPDF_BOOL(0)
	}

	if formFillInfoHandle.FormFillInfo.FFI_PopupMenu == nil {
		return C.FPDF_BOOL(0)
	}

	var pageRef references.FPDF_PAGE
	if pointerPageRef, ok := formFillInfoHandle.FormHandleHandle.pagePointers[unsafe.Pointer(page)]; ok {
		pageRef = pointerPageRef
	} else {
		pageHandle := formFillInfoHandle.Instance.registerPage(page, 0, nil)
		pageRef = pageHandle.nativeRef
	}

	result := formFillInfoHandle.FormFillInfo.FFI_PopupMenu(pageRef, int(menuFlag), float32(x), float32(y))
	if result {
		return C.FPDF_BOOL(1)
	}

	return C.FPDF_BOOL(0)
}

// go_formfill_FFI_OpenFile_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_OpenFile.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_OpenFile_cb
func go_formfill_FFI_OpenFile_cb(me *C.FPDF_FORMFILLINFO, fileFlag C.int, wsURL C.FPDF_WIDESTRING, mode C.ccharp) *C.FPDF_FILEHANDLER {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return nil
	}

	if formFillInfoHandle.FormFillInfo.FFI_OpenFile == nil {
		return nil
	}

	url := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(wsURL))
	decodedURL, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(url)
	if err != nil {
		return nil
	}

	resp := formFillInfoHandle.FormFillInfo.FFI_OpenFile(enums.FXFA_SAVEAS(fileFlag), decodedURL, C.GoString(mode))
	if resp != nil {
		fileHandlerStruct := &C.FPDF_FILEHANDLER{}
		fileHandlerStruct.clientData = pointer

		fileHandlerRef := uuid.New()
		fileHandlerRefString := fileHandlerRef.String()
		cFileHandlerRef := C.CString(fileHandlerRefString)

		C.FPDF_FILEHANDLER_SET_CB(fileHandlerStruct, cFileHandlerRef)

		fileHandlerInstance := &fileHandler{
			Struct:      fileHandlerStruct,
			FileHandler: resp,
			stringRef:   unsafe.Pointer(cFileHandlerRef),
		}

		fileHandlerHandles[fileHandlerRefString] = fileHandlerInstance
		fileHandlerPointers[unsafe.Pointer(fileHandlerStruct)] = fileHandlerInstance

		return fileHandlerStruct
	}

	return nil
}

// go_formfill_FFI_EmailTo_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_EmailTo.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_EmailTo_cb
func go_formfill_FFI_EmailTo_cb(me *C.FPDF_FORMFILLINFO, fileHandler *C.FPDF_FILEHANDLER, pTo, pSubject, pCC, pBcc, pMsg C.FPDF_WIDESTRING) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	// Check if we still have the file handler.
	fileHandlerHandler, ok := fileHandlerPointers[unsafe.Pointer(fileHandler)]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_EmailTo == nil {
		return
	}

	to := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(pTo))
	decodedTo, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(to)
	if err != nil {
		return
	}

	subject := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(pSubject))
	decodedSubject, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(subject)
	if err != nil {
		return
	}

	cc := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(pCC))
	decodedCC, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(cc)
	if err != nil {
		return
	}

	bcc := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(pBcc))
	decodedBcc, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(bcc)
	if err != nil {
		return
	}

	msg := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(pMsg))
	decodedMsg, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(msg)
	if err != nil {
		return
	}

	formFillInfoHandle.FormFillInfo.FFI_EmailTo(fileHandlerHandler.FileHandler, decodedTo, decodedSubject, decodedCC, decodedBcc, decodedMsg)

	return
}

// go_formfill_FFI_UploadTo_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_UploadTo.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_UploadTo_cb
func go_formfill_FFI_UploadTo_cb(me *C.FPDF_FORMFILLINFO, fileHandler *C.FPDF_FILEHANDLER, fileFlag C.int, uploadTo C.FPDF_WIDESTRING) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	// Check if we still have the file handler.
	fileHandlerHandler, ok := fileHandlerPointers[unsafe.Pointer(fileHandler)]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_UploadTo == nil {
		return
	}

	uploadToGo := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(uploadTo))
	decodedUploadTo, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(uploadToGo)
	if err != nil {
		return
	}

	formFillInfoHandle.FormFillInfo.FFI_UploadTo(fileHandlerHandler.FileHandler, enums.FXFA_SAVEAS(fileFlag), decodedUploadTo)

	return
}

// go_formfill_FFI_GetPlatform_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_GetPlatform.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_GetPlatform_cb
func go_formfill_FFI_GetPlatform_cb(me *C.FPDF_FORMFILLINFO, platform unsafe.Pointer, length C.int) C.int {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return C.int(0)
	}

	if formFillInfoHandle.FormFillInfo.FFI_GetPlatform == nil {
		return C.int(0)
	}

	goPlatform := formFillInfoHandle.FormFillInfo.FFI_GetPlatform()
	if platform != nil && length > 0 {
		target := unsafe.Slice((*byte)(platform), uint64(length))
		platformUTF16, err := formFillInfoHandle.Instance.transformUTF8ToUTF16LE(goPlatform)
		if err != nil {
			return C.int(0)
		}
		copy(target, platformUTF16)

		// Set last byte to NULL terminator.
		target[uint64(length)-1] = 0x00
	}

	return C.int((len(goPlatform) * 2) + 1)
}

// go_formfill_FFI_GetLanguage_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_GetLanguage.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_GetLanguage_cb
func go_formfill_FFI_GetLanguage_cb(me *C.FPDF_FORMFILLINFO, language unsafe.Pointer, length C.int) C.int {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return C.int(0)
	}

	if formFillInfoHandle.FormFillInfo.FFI_GetLanguage == nil {
		return C.int(0)
	}

	goLanguage := formFillInfoHandle.FormFillInfo.FFI_GetLanguage()
	if language != nil && length > 0 {
		target := unsafe.Slice((*byte)(language), uint64(length))
		languageUTF16, err := formFillInfoHandle.Instance.transformUTF8ToUTF16LE(goLanguage)
		if err != nil {
			return C.int(0)
		}
		copy(target, languageUTF16)
		// Set last byte to NULL terminator.
		target[uint64(length)-1] = 0x00
	}

	return C.int((len(goLanguage) * 2) + 1)
}

// go_formfill_FFI_DownloadFromURL_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_DownloadFromURL.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_DownloadFromURL_cb
func go_formfill_FFI_DownloadFromURL_cb(me *C.FPDF_FORMFILLINFO, URL C.FPDF_WIDESTRING) *C.FPDF_FILEHANDLER {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return nil
	}

	if formFillInfoHandle.FormFillInfo.FFI_DownloadFromURL == nil {
		return nil
	}

	url := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(URL))
	decodedURL, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(url)
	if err != nil {
		return nil
	}

	formFillInfoHandle.FormFillInfo.FFI_DownloadFromURL(decodedURL)

	return nil
}

// go_formfill_FFI_PostRequestURL_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_PostRequestURL.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_PostRequestURL_cb
func go_formfill_FFI_PostRequestURL_cb(me *C.FPDF_FORMFILLINFO, wsURL, wsData, wsContentType, wsEncode, wsHeader C.FPDF_WIDESTRING, response *C.FPDF_BSTR) C.FPDF_BOOL {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return C.FPDF_BOOL(0)
	}

	if formFillInfoHandle.FormFillInfo.FFI_PostRequestURL == nil {
		return C.FPDF_BOOL(0)
	}

	url := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(wsURL))
	decodedURL, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(url)
	if err != nil {
		return C.FPDF_BOOL(0)
	}

	data := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(wsData))
	decodedData, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(data)
	if err != nil {
		return C.FPDF_BOOL(0)
	}

	contentType := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(wsContentType))
	decodedContentType, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(contentType)
	if err != nil {
		return C.FPDF_BOOL(0)
	}

	encode := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(wsEncode))
	decodedEncode, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(encode)
	if err != nil {
		return C.FPDF_BOOL(0)
	}

	header := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(wsHeader))
	decodedHeader, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(header)
	if err != nil {
		return C.FPDF_BOOL(0)
	}

	var bStrRef references.FPDF_BSTR
	if pointerBStrRef, ok := formFillInfoHandle.FormHandleHandle.bStrPointers[unsafe.Pointer(response)]; ok {
		bStrRef = pointerBStrRef
	} else {
		bStrHandle := formFillInfoHandle.Instance.registerBStr(*response)
		bStrRef = bStrHandle.nativeRef
		formFillInfoHandle.FormHandleHandle.bStrPointers[unsafe.Pointer(response)] = bStrHandle.nativeRef
	}

	result := formFillInfoHandle.FormFillInfo.FFI_PostRequestURL(decodedURL, decodedData, decodedContentType, decodedEncode, decodedHeader, bStrRef)
	if result {
		return C.FPDF_BOOL(1)
	}

	return C.FPDF_BOOL(0)
}

// go_formfill_FFI_PutRequestURL_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_PutRequestURL.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_PutRequestURL_cb
func go_formfill_FFI_PutRequestURL_cb(me *C.FPDF_FORMFILLINFO, wsURL, wsData, wsEncode C.FPDF_WIDESTRING) C.FPDF_BOOL {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return C.FPDF_BOOL(0)
	}

	if formFillInfoHandle.FormFillInfo.FFI_PutRequestURL == nil {
		return C.FPDF_BOOL(0)
	}

	URL := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(wsURL))
	decodedURL, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(URL)
	if err != nil {
		return C.FPDF_BOOL(0)
	}

	data := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(wsData))
	decodedData, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(data)
	if err != nil {
		return C.FPDF_BOOL(0)
	}

	encode := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(wsEncode))
	decodedEncode, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(encode)
	if err != nil {
		return C.FPDF_BOOL(0)
	}

	result := formFillInfoHandle.FormFillInfo.FFI_PutRequestURL(decodedURL, decodedData, decodedEncode)
	if result {
		return C.FPDF_BOOL(1)
	}

	return C.FPDF_BOOL(0)
}

// go_formfill_FFI_OnFocusChange_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_OnFocusChange.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_OnFocusChange_cb
func go_formfill_FFI_OnFocusChange_cb(me *C.FPDF_FORMFILLINFO, annot C.FPDF_ANNOTATION, page_index C.int) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_OnFocusChange == nil {
		return
	}

	var annotationRef references.FPDF_ANNOTATION
	if pointerAnnotationRef, ok := formFillInfoHandle.FormHandleHandle.annotationPointers[unsafe.Pointer(annot)]; ok {
		annotationRef = pointerAnnotationRef
	} else {
		annotationHandle := formFillInfoHandle.Instance.registerAnnotation(annot)
		annotationRef = annotationHandle.nativeRef
		formFillInfoHandle.FormHandleHandle.annotationPointers[unsafe.Pointer(annotationHandle.handle)] = annotationHandle.nativeRef
	}

	formFillInfoHandle.FormFillInfo.FFI_OnFocusChange(annotationRef, int(page_index))

	return
}

// go_formfill_FFI_DoURIActionWithKeyboardModifier_cb is the Go implementation of FPDF_FORMFILLINFO::FFI_DoURIActionWithKeyboardModifier.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FORMFILLINFO structs.
//
//export go_formfill_FFI_DoURIActionWithKeyboardModifier_cb
func go_formfill_FFI_DoURIActionWithKeyboardModifier_cb(me *C.FPDF_FORMFILLINFO, uri C.FPDF_BYTESTRING, modifiers C.int) {
	pointer := unsafe.Pointer(me)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_DoURIActionWithKeyboardModifier == nil {
		return
	}

	formFillInfoHandle.FormFillInfo.FFI_DoURIActionWithKeyboardModifier(C.GoString((*C.char)(uri)), int(modifiers))

	return
}

// go_jsplatform_app_alert_cb is the Go implementation of IPDF_JSPLATFORM::app_alert.
// It is exported through cgo so that we can use the reference to it and set
// it on IPDF_JSPLATFORM structs.
//
//export go_jsplatform_app_alert_cb
func go_jsplatform_app_alert_cb(me *C.IPDF_JSPLATFORM, msg, title C.FPDF_WIDESTRING, alertType, icon C.int) C.int {
	pointer := unsafe.Pointer(me.m_pFormfillinfo)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return C.int(0)
	}

	if formFillInfoHandle.JSPlatform == nil || formFillInfoHandle.JSPlatform.App_alert == nil {
		return C.int(0)
	}

	wsMsg := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(msg))
	decodedMsg, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(wsMsg)
	if err != nil {
		return C.FPDF_BOOL(0)
	}

	wsTitle := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(title))
	decodedTitle, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(wsTitle)
	if err != nil {
		return C.FPDF_BOOL(0)
	}

	return C.int(formFillInfoHandle.JSPlatform.App_alert(decodedMsg, decodedTitle, enums.JSPLATFORM_ALERT_BUTTON(alertType), enums.JSPLATFORM_ALERT_ICON(icon)))
}

// go_jsplatform_app_beep_cb is the Go implementation of IPDF_JSPLATFORM::app_beep.
// It is exported through cgo so that we can use the reference to it and set
// it on IPDF_JSPLATFORM structs.
//
//export go_jsplatform_app_beep_cb
func go_jsplatform_app_beep_cb(me *C.IPDF_JSPLATFORM, nType C.int) {
	pointer := unsafe.Pointer(me.m_pFormfillinfo)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.JSPlatform == nil || formFillInfoHandle.JSPlatform.App_beep == nil {
		return
	}

	formFillInfoHandle.JSPlatform.App_beep(enums.JSPLATFORM_BEEP(nType))

	return
}

// go_jsplatform_app_response_cb is the Go implementation of IPDF_JSPLATFORM::app_response.
// It is exported through cgo so that we can use the reference to it and set
// it on IPDF_JSPLATFORM structs.
//
//export go_jsplatform_app_response_cb
func go_jsplatform_app_response_cb(me *C.IPDF_JSPLATFORM, question, title, defaultValue, cLabel C.FPDF_WIDESTRING, bPassword C.FPDF_BOOL, response unsafe.Pointer, length C.int) C.int {
	pointer := unsafe.Pointer(me.m_pFormfillinfo)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return C.int(0)
	}

	if formFillInfoHandle.JSPlatform == nil || formFillInfoHandle.JSPlatform.App_response == nil {
		return C.int(0)
	}

	wsQuestion := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(question))
	decodedQuestion, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(wsQuestion)
	if err != nil {
		return C.int(0)
	}

	wsTitle := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(title))
	decodedTitle, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(wsTitle)
	if err != nil {
		return C.int(0)
	}

	wsDefaultValue := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(defaultValue))
	decodedDefaultValue, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(wsDefaultValue)
	if err != nil {
		return C.int(0)
	}

	wsCLabel := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(cLabel))
	decodedCLabel, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(wsCLabel)
	if err != nil {
		return C.int(0)
	}

	returnedResponse := formFillInfoHandle.JSPlatform.App_response(decodedQuestion, decodedTitle, decodedDefaultValue, decodedCLabel, bPassword == C.FPDF_BOOL(1))
	if int(length) > 0 {
		target := unsafe.Slice((*byte)(response), uint64(length))
		copy(target, append([]byte(returnedResponse), 0x00))
	}

	return C.int(len(returnedResponse))
}

// go_jsplatform_Doc_getFilePath_cb is the Go implementation of IPDF_JSPLATFORM::Doc_getFilePath.
// It is exported through cgo so that we can use the reference to it and set
// it on IPDF_JSPLATFORM structs.
//
//export go_jsplatform_Doc_getFilePath_cb
func go_jsplatform_Doc_getFilePath_cb(me *C.IPDF_JSPLATFORM, filePath unsafe.Pointer, length C.int) C.int {
	pointer := unsafe.Pointer(me.m_pFormfillinfo)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return C.int(0)
	}

	if formFillInfoHandle.JSPlatform == nil || formFillInfoHandle.JSPlatform.Doc_getFilePath == nil {
		return C.int(0)
	}

	returnedFilePath := formFillInfoHandle.JSPlatform.Doc_getFilePath()
	neededLength := len(returnedFilePath) + 1

	if int(length) > 0 && int(length) >= neededLength {
		target := unsafe.Slice((*byte)(filePath), uint64(length))
		copy(target, append([]byte(returnedFilePath), 0x00))
	}

	return C.int(neededLength)
}

// go_jsplatform_Doc_mail_cb is the Go implementation of IPDF_JSPLATFORM::Doc_mail.
// It is exported through cgo so that we can use the reference to it and set
// it on IPDF_JSPLATFORM structs.
//
//export go_jsplatform_Doc_mail_cb
func go_jsplatform_Doc_mail_cb(me *C.IPDF_JSPLATFORM, mailData unsafe.Pointer, length C.int, bUI C.FPDF_BOOL, to, subject, cc, bcc, msg C.FPDF_WIDESTRING) {
	pointer := unsafe.Pointer(me.m_pFormfillinfo)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.JSPlatform == nil || formFillInfoHandle.JSPlatform.Doc_mail == nil {
		return
	}

	data := unsafe.Slice((*byte)(mailData), uint64(length))

	wsTo := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(to))
	decodedTo, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(wsTo)
	if err != nil {
		return
	}

	wsSubject := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(subject))
	decodedSubject, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(wsSubject)
	if err != nil {
		return
	}

	wsCc := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(cc))
	decodedCc, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(wsCc)
	if err != nil {
		return
	}

	wsBcc := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(bcc))
	decodedBcc, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(wsBcc)
	if err != nil {
		return
	}

	wsMsg := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(msg))
	decodedMsg, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(wsMsg)
	if err != nil {
		return
	}

	formFillInfoHandle.JSPlatform.Doc_mail(data, bUI == C.FPDF_BOOL(1), decodedTo, decodedSubject, decodedCc, decodedBcc, decodedMsg)

	return
}

// go_jsplatform_Doc_print_cb is the Go implementation of IPDF_JSPLATFORM::Doc_print.
// It is exported through cgo so that we can use the reference to it and set
// it on IPDF_JSPLATFORM structs.
//
//export go_jsplatform_Doc_print_cb
func go_jsplatform_Doc_print_cb(me *C.IPDF_JSPLATFORM, bUI C.FPDF_BOOL, nStart, nEnd C.int, bSilent, bShrinkToFit, bPrintAsImage, bReverse, bAnnotations C.FPDF_BOOL) {
	pointer := unsafe.Pointer(me.m_pFormfillinfo)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.JSPlatform == nil || formFillInfoHandle.JSPlatform.Doc_print == nil {
		return
	}

	formFillInfoHandle.JSPlatform.Doc_print(bUI == C.FPDF_BOOL(1), int(nStart), int(nEnd), bSilent == C.FPDF_BOOL(1), bShrinkToFit == C.FPDF_BOOL(1), bPrintAsImage == C.FPDF_BOOL(1), bReverse == C.FPDF_BOOL(1), bAnnotations == C.FPDF_BOOL(1))

	return
}

// go_jsplatform_Doc_submitForm_cb is the Go implementation of IPDF_JSPLATFORM::Doc_submitForm.
// It is exported through cgo so that we can use the reference to it and set
// it on IPDF_JSPLATFORM structs.
//
//export go_jsplatform_Doc_submitForm_cb
func go_jsplatform_Doc_submitForm_cb(me *C.IPDF_JSPLATFORM, formData unsafe.Pointer, length C.int, url C.FPDF_WIDESTRING) {
	pointer := unsafe.Pointer(me.m_pFormfillinfo)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.JSPlatform == nil || formFillInfoHandle.JSPlatform.Doc_submitForm == nil {
		return
	}

	data := unsafe.Slice((*byte)(formData), uint64(length))

	wsURL := formFillInfoHandle.Instance.readBytesUntilTerminator(unsafe.Pointer(url))
	decodedURL, err := formFillInfoHandle.Instance.transformUTF16LEToUTF8(wsURL)
	if err != nil {
		return
	}

	formFillInfoHandle.JSPlatform.Doc_submitForm(data, decodedURL)

	return
}

// go_jsplatform_Doc_gotoPage_cb is the Go implementation of IPDF_JSPLATFORM::Doc_gotoPage.
// It is exported through cgo so that we can use the reference to it and set
// it on IPDF_JSPLATFORM structs.
//
//export go_jsplatform_Doc_gotoPage_cb
func go_jsplatform_Doc_gotoPage_cb(me *C.IPDF_JSPLATFORM, nPageNum C.int) {
	pointer := unsafe.Pointer(me.m_pFormfillinfo)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return
	}

	if formFillInfoHandle.JSPlatform == nil || formFillInfoHandle.JSPlatform.Doc_gotoPage == nil {
		return
	}

	formFillInfoHandle.JSPlatform.Doc_gotoPage(int(nPageNum))

	return
}

// go_jsplatform_Field_browse_cb is the Go implementation of IPDF_JSPLATFORM::Field_browse.
// It is exported through cgo so that we can use the reference to it and set
// it on IPDF_JSPLATFORM structs.
//
//export go_jsplatform_Field_browse_cb
func go_jsplatform_Field_browse_cb(me *C.IPDF_JSPLATFORM, filePath unsafe.Pointer, length C.int) C.int {
	pointer := unsafe.Pointer(me.m_pFormfillinfo)

	// Check if we still have the callback.
	formFillInfoHandle, ok := formFillInfoHandles[pointer]
	if !ok {
		return C.int(0)
	}

	if formFillInfoHandle.JSPlatform == nil || formFillInfoHandle.JSPlatform.Field_browse == nil {
		return C.int(0)
	}

	returnedFilePath := formFillInfoHandle.JSPlatform.Field_browse()
	neededLength := len(returnedFilePath) + 1

	if int(length) > 0 && int(length) >= neededLength {
		target := unsafe.Slice((*byte)(filePath), uint64(length))
		copy(target, append([]byte(returnedFilePath), 0x00))
	}

	return C.int(neededLength)
}

// go_filehandler_Release_cb is the Go implementation of FPDF_FILEHANDLER::Release.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FILEHANDLER structs.
//
//export go_filehandler_Release_cb
func go_filehandler_Release_cb(clientData unsafe.Pointer) {
	fileHandlerRef := C.GoString((*C.char)(clientData))

	// Check if we still have the callback.
	fileHandlerHandle, ok := fileHandlerHandles[fileHandlerRef]
	if !ok {
		return
	}

	if fileHandlerHandle.FileHandler.Release != nil {
		fileHandlerHandle.FileHandler.Release()
	}

	C.free(fileHandlerHandle.stringRef)

	delete(fileHandlerHandles, fileHandlerRef)
	delete(fileHandlerPointers, unsafe.Pointer(fileHandlerHandle.FileHandler))
}

// go_filehandler_GetSize_cb is the Go implementation of FPDF_FILEHANDLER::GetSize.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FILEHANDLER structs.
//
//export go_filehandler_GetSize_cb
func go_filehandler_GetSize_cb(clientData unsafe.Pointer) C.FPDF_DWORD {
	fileHandlerRef := C.GoString((*C.char)(clientData))

	// Check if we still have the callback.
	fileHandlerHandle, ok := fileHandlerHandles[fileHandlerRef]
	if !ok {
		return C.FPDF_DWORD(0)
	}

	if fileHandlerHandle.FileHandler.GetSize != nil {
		return C.FPDF_DWORD(fileHandlerHandle.FileHandler.GetSize())
	}

	return C.FPDF_DWORD(0)
}

// go_filehandler_ReadBlock_cb is the Go implementation of FPDF_FILEHANDLER::ReadBlock.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FILEHANDLER structs.
//
//export go_filehandler_ReadBlock_cb
func go_filehandler_ReadBlock_cb(clientData unsafe.Pointer, offset C.FPDF_DWORD, buffer unsafe.Pointer, size C.FPDF_DWORD) C.FPDF_RESULT {
	fileHandlerRef := C.GoString((*C.char)(clientData))

	// Check if we still have the callback.
	fileHandlerHandle, ok := fileHandlerHandles[fileHandlerRef]
	if !ok {
		return C.FPDF_RESULT(0)
	}

	if fileHandlerHandle.FileHandler.ReadBlock != nil {
		data, err := fileHandlerHandle.FileHandler.ReadBlock(uint64(offset), uint64(size))
		if err != nil {
			return C.FPDF_RESULT(1)
		}

		if int(size) > 0 && int(size) >= len(data) {
			target := unsafe.Slice((*byte)(buffer), uint64(size))
			copy(target, data)
		}
	}

	return C.FPDF_RESULT(0)
}

// go_filehandler_WriteBlock_cb is the Go implementation of FPDF_FILEHANDLER::WriteBlock.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FILEHANDLER structs.
//
//export go_filehandler_WriteBlock_cb
func go_filehandler_WriteBlock_cb(clientData unsafe.Pointer, offset C.FPDF_DWORD, buffer C.cvoidp, size C.FPDF_DWORD) C.FPDF_RESULT {
	fileHandlerRef := C.GoString((*C.char)(clientData))

	// Check if we still have the callback.
	fileHandlerHandle, ok := fileHandlerHandles[fileHandlerRef]
	if !ok {
		return C.FPDF_RESULT(0)
	}

	if fileHandlerHandle.FileHandler.WriteBlock != nil {
		data := unsafe.Slice((*byte)(buffer), uint64(size))
		err := fileHandlerHandle.FileHandler.WriteBlock(uint64(offset), data)
		if err != nil {
			return C.FPDF_RESULT(1)
		}
	}

	return C.FPDF_RESULT(0)
}

// go_filehandler_Flush_cb is the Go implementation of FPDF_FILEHANDLER::Flush.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FILEHANDLER structs.
//
//export go_filehandler_Flush_cb
func go_filehandler_Flush_cb(clientData unsafe.Pointer) C.FPDF_RESULT {
	fileHandlerRef := C.GoString((*C.char)(clientData))

	// Check if we still have the callback.
	fileHandlerHandle, ok := fileHandlerHandles[fileHandlerRef]
	if !ok {
		return C.FPDF_RESULT(0)
	}

	if fileHandlerHandle.FileHandler.Flush != nil {
		err := fileHandlerHandle.FileHandler.Flush()
		if err != nil {
			return C.FPDF_RESULT(1)
		}
	}

	return C.FPDF_RESULT(0)
}

// go_filehandler_Truncate_cb is the Go implementation of FPDF_FILEHANDLER::Truncate.
// It is exported through cgo so that we can use the reference to it and set
// it on FPDF_FILEHANDLER structs.
//
//export go_filehandler_Truncate_cb
func go_filehandler_Truncate_cb(clientData unsafe.Pointer, size C.FPDF_DWORD) C.FPDF_RESULT {
	fileHandlerRef := C.GoString((*C.char)(clientData))

	// Check if we still have the callback.
	fileHandlerHandle, ok := fileHandlerHandles[fileHandlerRef]
	if !ok {
		return C.FPDF_RESULT(0)
	}

	if fileHandlerHandle.FileHandler.Truncate != nil {
		err := fileHandlerHandle.FileHandler.Truncate(uint64(size))
		if err != nil {
			return C.FPDF_RESULT(1)
		}
	}

	return C.FPDF_RESULT(0)
}

type FormFillInfo struct {
	Struct           *C.FPDF_FORMFILLINFO
	JSPlatformStruct *C.IPDF_JSPLATFORM
	FormFillInfo     *structs.FPDF_FORMFILLINFO
	JSPlatform       *structs.IPDF_JSPLATFORM
	FormHandleHandle *FormHandleHandle
	Instance         *PdfiumImplementation
}

var formFillInfoHandles = map[unsafe.Pointer]*FormFillInfo{}

type fileHandler struct {
	Struct      *C.FPDF_FILEHANDLER
	FileHandler *structs.FPDF_FILEHANDLER
	stringRef   unsafe.Pointer
}

var fileHandlerHandles = map[string]*fileHandler{}
var fileHandlerPointers = map[unsafe.Pointer]*fileHandler{}

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

	formFillVersion := getFormFillVersion()

	formInfoStruct := &C.FPDF_FORMFILLINFO{}
	formInfoStruct.version = C.int(formFillVersion)
	if request.FormFillInfo.XFA_disabled {
		formInfoStruct.xfa_disabled = C.FPDF_BOOL(1)
	} else {
		formInfoStruct.xfa_disabled = C.FPDF_BOOL(0)
	}

	C.FPDF_FORMFILLINFO_SET_CB(formInfoStruct)

	if formFillVersion > 1 && request.FormFillInfo.JsPlatform != nil {
		C.IPDF_JSPLATFORM_SET_CB(formInfoStruct.m_pJsPlatform)
	}

	formHandle := C.FPDFDOC_InitFormFillEnvironment(documentHandle.handle, formInfoStruct)
	if formHandle == nil {
		return nil, errors.New("could not init form fill environment")
	}

	formHandleHandle := p.registerFormHandle(formHandle, unsafe.Pointer(formInfoStruct))

	formFillInfo := &FormFillInfo{
		Struct:           formInfoStruct,
		JSPlatformStruct: formInfoStruct.m_pJsPlatform,
		FormFillInfo:     &request.FormFillInfo,
		JSPlatform:       request.FormFillInfo.JsPlatform,
		FormHandleHandle: formHandleHandle,
		Instance:         p,
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

	documentHandle, err := p.getDocumentHandle(pageHandle.documentRef)
	if err != nil {
		return nil, err
	}

	formHandleHandle, err := p.getFormHandleHandle(request.FormHandle)
	if err != nil {
		return nil, err
	}

	C.FORM_OnAfterLoadPage(pageHandle.handle, formHandleHandle.handle)

	// Store pointers so that we can reference them in the events to prevent
	// leaving a lot of references to the same page/document.
	formHandleHandle.pagePointers[unsafe.Pointer(pageHandle.handle)] = pageHandle.nativeRef
	formHandleHandle.documentPointers[unsafe.Pointer(documentHandle.handle)] = documentHandle.nativeRef

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

	// Remove pointer reference.
	if _, ok := formHandleHandle.pagePointers[unsafe.Pointer(pageHandle.handle)]; ok {
		delete(formHandleHandle.pagePointers, unsafe.Pointer(pageHandle.handle))
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
