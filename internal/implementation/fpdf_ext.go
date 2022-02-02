package implementation

/*
#cgo pkg-config: pdfium
#include "fpdf_ext.h"
#include <time.h>
#include <stdio.h>

extern void go_un_sp_obj_cb(struct _UNSUPPORT_INFO *pThis, int nType);
static inline void UNSUPPORT_INFO_SET_CALLBACK(UNSUPPORT_INFO *ui) {
	ui->FSDK_UnSupport_Handler = &go_un_sp_obj_cb;
}

extern time_t go_time_function_cb();
static inline void FSDK_SetTimeFunction_SET_GO_METHOD() {
	FSDK_SetTimeFunction(&go_time_function_cb);
}

typedef struct GoPdfiumLocalTime {
    int tm_sec;
	int tm_min;
	int tm_hour;
	int tm_mday;
	int tm_mon;
	int tm_year;
	int tm_wday;
	int tm_yday;
	int tm_isdst;
} GoPdfiumLocalTime;

typedef const time_t ctime_t;
extern GoPdfiumLocalTime* go_local_time_function_cb(ctime_t *curTime);

static inline struct tm* local_time_function_cb(const time_t *curTime) {
	struct tm* local_time;

	// For some reason I could not get the tm struct available in Go.
	// So we have to do it like this.
	GoPdfiumLocalTime* goLocalTime = go_local_time_function_cb(curTime);
	local_time->tm_sec = goLocalTime->tm_sec;
	local_time->tm_min = goLocalTime->tm_min;
	local_time->tm_hour = goLocalTime->tm_hour;
	local_time->tm_mday = goLocalTime->tm_mday;
	local_time->tm_mon = goLocalTime->tm_mon;
	local_time->tm_year = goLocalTime->tm_year;
	local_time->tm_wday = goLocalTime->tm_wday;
	local_time->tm_yday = goLocalTime->tm_yday;
	local_time->tm_isdst = goLocalTime->tm_isdst;

	return local_time;
}

static inline void FSDK_SetLocaltimeFunction_SET_GO_METHOD() {
	FSDK_SetLocaltimeFunction(&local_time_function_cb);
}
*/
import "C"
import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDFDoc_GetPageMode returns the document's page mode, which describes how the document should be displayed when opened.
func (p *PdfiumImplementation) FPDFDoc_GetPageMode(request *requests.FPDFDoc_GetPageMode) (*responses.FPDFDoc_GetPageMode, error) {
	p.Lock()
	defer p.Unlock()

	documentHandle, err := p.getDocumentHandle(request.Document)
	if err != nil {
		return nil, err
	}

	pageMode := C.FPDFDoc_GetPageMode(documentHandle.handle)

	return &responses.FPDFDoc_GetPageMode{
		PageMode: responses.FPDFDoc_GetPageModeMode(pageMode),
	}, nil
}

var currentUnsupportedObjectHandler requests.UnSpObjProcessHandler

//export go_un_sp_obj_cb
func go_un_sp_obj_cb(pThis *C.UNSUPPORT_INFO, nType C.int) {
	if currentUnsupportedObjectHandler != nil {
		currentUnsupportedObjectHandler(enums.FPDF_UNSP(nType))
	}
}

// FSDK_SetUnSpObjProcessHandler set ups an unsupported object handler.
// Since callbacks can't be transferred between processes in gRPC, you can only use this in single-threaded mode.
func (p *PdfiumImplementation) FSDK_SetUnSpObjProcessHandler(request *requests.FSDK_SetUnSpObjProcessHandler) (*responses.FSDK_SetUnSpObjProcessHandler, error) {
	p.Lock()
	defer p.Unlock()

	currentUnsupportedObjectHandler = request.UnSpObjProcessHandler

	handler := C.UNSUPPORT_INFO{}
	handler.version = 1

	// Set the Go callback through cgo.
	C.UNSUPPORT_INFO_SET_CALLBACK(&handler)

	return &responses.FSDK_SetUnSpObjProcessHandler{}, nil
}

var currentTimeHandler requests.SetTimeFunction

//export go_time_function_cb
func go_time_function_cb() C.time_t {
	if currentTimeHandler != nil {
		return C.time_t(currentTimeHandler())
	}

	result := C.time_t(0)
	C.time(&result)

	return result
}

// FSDK_SetTimeFunction sets a replacement function for calls to time().
// This API is intended to be used only for testing, thus may cause PDFium to behave poorly in production environments.
// Since callbacks can't be transferred between processes in gRPC, you can only use this in single-threaded mode.
func (p *PdfiumImplementation) FSDK_SetTimeFunction(request *requests.FSDK_SetTimeFunction) (*responses.FSDK_SetTimeFunction, error) {
	p.Lock()
	defer p.Unlock()

	if request.Function == nil {
		C.FSDK_SetTimeFunction(nil)
	} else {
		C.FSDK_SetTimeFunction_SET_GO_METHOD()
	}

	return nil, nil
}

var currentLocalTimeHandler requests.SetLocaltimeFunction

//export go_local_time_function_cb
func go_local_time_function_cb(timer *C.ctime_t) *C.GoPdfiumLocalTime {
	timeStruct := C.GoPdfiumLocalTime{}

	if currentTimeHandler != nil {
		// Convert from C to go.
		localTime := currentLocalTimeHandler(int64(*timer))
		timeStruct.tm_sec = C.int(localTime.TmSec)
		timeStruct.tm_min = C.int(localTime.TmMin)
		timeStruct.tm_hour = C.int(localTime.TmHour)
		timeStruct.tm_mday = C.int(localTime.TmMday)
		timeStruct.tm_mon = C.int(localTime.TmMon)
		timeStruct.tm_year = C.int(localTime.TmYear)
		timeStruct.tm_wday = C.int(localTime.TmWday)
		timeStruct.tm_yday = C.int(localTime.TmYday)
		timeStruct.tm_isdst = C.int(localTime.TmIsdst)
	} else {
		result := C.time_t(0)
		C.time(&result)

		// Copy over real local time.
		realLocalTime := C.localtime(&result)
		timeStruct.tm_sec = realLocalTime.tm_sec
		timeStruct.tm_min = realLocalTime.tm_min
		timeStruct.tm_hour = realLocalTime.tm_hour
		timeStruct.tm_mday = realLocalTime.tm_mday
		timeStruct.tm_mon = realLocalTime.tm_mon
		timeStruct.tm_year = realLocalTime.tm_year
		timeStruct.tm_wday = realLocalTime.tm_wday
		timeStruct.tm_yday = realLocalTime.tm_yday
		timeStruct.tm_isdst = realLocalTime.tm_isdst
	}

	return (*C.GoPdfiumLocalTime)(&timeStruct)
}

// FSDK_SetLocaltimeFunction sets a replacement function for calls to localtime().
// This API is intended to be used only for testing, thus may cause PDFium to behave poorly in production environments.
// Since callbacks can't be transferred between processes in gRPC, you can only use this in single-threaded mode.
func (p *PdfiumImplementation) FSDK_SetLocaltimeFunction(request *requests.FSDK_SetLocaltimeFunction) (*responses.FSDK_SetLocaltimeFunction, error) {
	p.Lock()
	defer p.Unlock()

	if request.Function == nil {
		C.FSDK_SetLocaltimeFunction(nil)
	} else {
		C.FSDK_SetLocaltimeFunction_SET_GO_METHOD()
	}

	return nil, nil
}
