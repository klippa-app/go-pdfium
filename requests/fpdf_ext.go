package requests

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
)

type FPDFDoc_GetPageMode struct {
	Document references.FPDF_DOCUMENT
}

type UnSpObjProcessHandler func(enums.FPDF_UNSP)

type FSDK_SetUnSpObjProcessHandler struct {
	UnSpObjProcessHandler func(enums.FPDF_UNSP)
}

type SetTimeFunction func() int64

type FSDK_SetTimeFunction struct {
	Function SetTimeFunction // Alternate implementation of time(), or nil to restore to actual time() call itself.
}

type SetLocaltime struct {
	TmSec   int /* seconds */
	TmMin   int /* minutes */
	TmHour  int /* hours */
	TmMday  int /* day of the month */
	TmMon   int /* month */
	TmYear  int /* year */
	TmWday  int /* day of the week */
	TmYday  int /* day in the year */
	TmIsdst int /* daylight saving time */
}

type SetLocaltimeFunction func(int64) SetLocaltime

type FSDK_SetLocaltimeFunction struct {
	Function SetLocaltimeFunction // Alternate implementation of localtime(), or nil to restore to actual localtime() call itself.
}
