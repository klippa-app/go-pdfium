package imports

import (
	"context"
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/internal/implementation_webassembly"
	"github.com/tetratelabs/wazero/api"
	"log"
)

type FPDF_FILEACCESS_CB struct {
}

func (cb FPDF_FILEACCESS_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	paramPointer := uint32(stack[0])
	position := uint32(stack[1])
	pBufPointer := uint32(stack[2])
	size := uint32(stack[3])

	mem := mod.Memory()

	param, ok := mem.ReadUint32Le(paramPointer)
	if !ok {
		stack[0] = uint64(0)
		return
	}

	// Check if we have the file referenced in param.
	openFile, ok := implementation_webassembly.FileReaders[param]
	if !ok {
		stack[0] = uint64(0)
		return
	}

	// Seek to the right position.
	_, err := openFile.Reader.Seek(int64(position), 0)
	if err != nil {
		stack[0] = uint64(0)
		return
	}

	// Read the requested data into a buffer.
	readBuffer := make([]byte, size)
	n, err := openFile.Reader.Read(readBuffer)
	if n == 0 || err != nil {
		stack[0] = uint64(0)
		return
	}

	ok = mem.Write(pBufPointer, readBuffer)
	if !ok {
		stack[0] = uint64(0)
		return
	}

	stack[0] = uint64(n)
	return
}

type FPDF_FILEWRITE_CB struct {
}

func (cb FPDF_FILEWRITE_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	fileWritePointer := uint32(stack[0])
	pDataPointer := uint32(stack[1])
	size := uint32(stack[2])

	mem := mod.Memory()

	// Check if we have the file referenced in param.
	openWriter, ok := implementation_webassembly.FileWriters[fileWritePointer]
	if !ok {
		stack[0] = uint64(0)
		return
	}

	pBuf, ok := mem.Read(pDataPointer, size)
	if !ok {
		stack[0] = uint64(0)
		return
	}

	n, err := openWriter.Writer.Write(pBuf)
	if err != nil {
		stack[0] = uint64(0)
		return
	}

	stack[0] = uint64(n)
	return
}

type FX_FILEAVAIL_IS_DATA_AVAILABLE_CB struct {
}

func (cb FX_FILEAVAIL_IS_DATA_AVAILABLE_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])
	offset := uint32(stack[1])
	size := uint32(stack[2])

	fileAvail, ok := implementation_webassembly.FileAvailables[me]
	if !ok {
		stack[0] = uint64(0)
		return
	}

	if fileAvail.DataAvailableCallback(uint64(offset), uint64(size)) {
		stack[0] = uint64(1)
		return
	}

	stack[0] = uint64(0)
	return
}

type FX_DOWNLOADHINTS_ADD_SEGMENT_CB struct {
}

func (cb FX_DOWNLOADHINTS_ADD_SEGMENT_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])
	offset := uint32(stack[1])
	size := uint32(stack[2])

	fileHint, ok := implementation_webassembly.FileHints[me]
	if !ok {
		return
	}

	fileHint.AddSegmentCallback(uint64(offset), uint64(size))
}

type UNSUPPORT_INFO_HANDLER_CB struct {
}

func (cb UNSUPPORT_INFO_HANDLER_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	ntype := uint32(stack[1])

	if implementation_webassembly.CurrentUnsupportedObjectHandler != nil {
		implementation_webassembly.CurrentUnsupportedObjectHandler(enums.FPDF_UNSP(ntype))
	}
}

type FSDK_SetTimeFunction_CB struct {
}

func (cb FSDK_SetTimeFunction_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	currentTime := uint64(0)
	if implementation_webassembly.CurrentTimeHandler != nil {
		currentTime = api.EncodeI64(implementation_webassembly.CurrentTimeHandler())
	}

	stack[0] = currentTime
	return
}

type FSDK_SetLocaltimeFunction_CB struct {
}

// re-use memory to prevent allocating more than necessary.
var lastFSDK_SetLocaltimeFunctionPointer *uint64

func (cb FSDK_SetLocaltimeFunction_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	timestamp := uint32(stack[0])

	currentLocalTime := uint64(0)
	if implementation_webassembly.CurrentLocalTimeHandler != nil {
		localTime := implementation_webassembly.CurrentLocalTimeHandler(int64(timestamp))

		if lastFSDK_SetLocaltimeFunctionPointer == nil {
			// 9 int fields in localtime.
			results, err := mod.ExportedFunction("malloc").Call(ctx, 4*9)
			if err != nil {
				log.Printf("Could not allocate memory")
				stack[0] = currentLocalTime
				return
			}
			lastFSDK_SetLocaltimeFunctionPointer = &results[0]
		}

		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer), api.EncodeI32(int32(localTime.TmSec)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+4), api.EncodeI32(int32(localTime.TmMin)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+8), api.EncodeI32(int32(localTime.TmHour)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+12), api.EncodeI32(int32(localTime.TmMday)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+16), api.EncodeI32(int32(localTime.TmMon)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+20), api.EncodeI32(int32(localTime.TmYear)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+24), api.EncodeI32(int32(localTime.TmWday)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+28), api.EncodeI32(int32(localTime.TmYday)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+32), api.EncodeI32(int32(localTime.TmIsdst)))

		currentLocalTime = *lastFSDK_SetLocaltimeFunctionPointer
	}

	stack[0] = currentLocalTime
	return
}

type FPDF_FORMFILLINFO_Release_CB struct {
}

func (cb FPDF_FORMFILLINFO_Release_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])

	// Check if we still have the callback.
	formFillInfoHandle, ok := implementation_webassembly.FormFillInfoHandles[me]
	if !ok {
		return
	}

	// @todo: do I have anything to cleanup for myself?

	formFillInfoHandle.Release()
}

type FPDF_FORMFILLINFO_FFI_Invalidate_CB struct {
}

func (cb FPDF_FORMFILLINFO_FFI_Invalidate_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])
	page := uint32(stack[1])
	left := uint64(stack[2])
	top := uint64(stack[3])
	right := uint64(stack[4])
	bottom := uint64(stack[5])

	// Check if we still have the callback.
	formFillInfoHandle, ok := implementation_webassembly.FormFillInfoHandles[me]
	if !ok {
		return
	}

	formFillInfoHandle.FFI_Invalidate_CB(page, left, top, right, bottom)
}

type FPDF_FORMFILLINFO_FFI_OutputSelectedRect_CB struct {
}

func (cb FPDF_FORMFILLINFO_FFI_OutputSelectedRect_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])
	page := uint32(stack[1])
	left := uint64(stack[2])
	top := uint64(stack[3])
	right := uint64(stack[4])
	bottom := uint64(stack[5])

	// Check if we still have the callback.
	formFillInfoHandle, ok := implementation_webassembly.FormFillInfoHandles[me]
	if !ok {
		return
	}

	formFillInfoHandle.FFI_OutputSelectedRect(page, left, top, right, bottom)
}

type FPDF_FORMFILLINFO_FFI_SetCursor_CB struct {
}

func (cb FPDF_FORMFILLINFO_FFI_SetCursor_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])
	cursor := uint32(stack[1])

	// Check if we still have the callback.
	formFillInfoHandle, ok := implementation_webassembly.FormFillInfoHandles[me]
	if !ok {
		return
	}

	formFillInfoHandle.FFI_SetCursor(cursor)
}

type FPDF_FORMFILLINFO_FFI_SetTimer_CB struct {
}

func (cb FPDF_FORMFILLINFO_FFI_SetTimer_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])
	uElapse := uint32(stack[1])
	lpTimerFunc := uint32(stack[2])

	// Check if we still have the callback.
	formFillInfoHandle, ok := implementation_webassembly.FormFillInfoHandles[me]
	if !ok {
		return
	}

	id := formFillInfoHandle.FFI_SetTimer(uElapse, lpTimerFunc)
	stack[0] = uint64(id)
}

type FPDF_FORMFILLINFO_FFI_KillTimer_CB struct {
}

func (cb FPDF_FORMFILLINFO_FFI_KillTimer_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])
	nTimerID := uint32(stack[1])

	// Check if we still have the callback.
	formFillInfoHandle, ok := implementation_webassembly.FormFillInfoHandles[me]
	if !ok {
		return
	}

	formFillInfoHandle.FFI_KillTimer(int(nTimerID))
}

type FPDF_FORMFILLINFO_FFI_GetLocalTime_CB struct {
}

// re-use memory to prevent allocating more than necessary.
var lastFPDF_FORMFILLINFO_FFI_GetLocalTimePointer *uint64

func (cb FPDF_FORMFILLINFO_FFI_GetLocalTime_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])

	// Check if we still have the callback.
	formFillInfoHandle, ok := implementation_webassembly.FormFillInfoHandles[me]
	if !ok {
		return
	}

	localTime := formFillInfoHandle.FFI_GetLocalTime()

	currentLocalTime := uint64(0)

	if lastFPDF_FORMFILLINFO_FFI_GetLocalTimePointer == nil {
		// 9 int fields in localtime.
		results, err := mod.ExportedFunction("malloc").Call(ctx, 4*9)
		if err != nil {
			log.Printf("Could not allocate memory")
			stack[0] = currentLocalTime
			return
		}
		lastFPDF_FORMFILLINFO_FFI_GetLocalTimePointer = &results[0]
	}

	// 8 * ushort (2 bytes)
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer), api.EncodeI32(int32(localTime.Year)))
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+2), api.EncodeI32(int32(localTime.Month)))
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+4), api.EncodeI32(int32(localTime.DayOfWeek)))
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+6), api.EncodeI32(int32(localTime.Day)))
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+8), api.EncodeI32(int32(localTime.Hour)))
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+10), api.EncodeI32(int32(localTime.Minute)))
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+12), api.EncodeI32(int32(localTime.Second)))
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+14), api.EncodeI32(int32(localTime.Milliseconds)))

	stack[0] = uint64(*lastFSDK_SetLocaltimeFunctionPointer)
}

type FPDF_FORMFILLINFO_FFI_OnChange_CB struct {
}

func (cb FPDF_FORMFILLINFO_FFI_OnChange_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])

	// Check if we still have the callback.
	formFillInfoHandle, ok := implementation_webassembly.FormFillInfoHandles[me]
	if !ok {
		return
	}

	formFillInfoHandle.FFI_OnChange()
}

type FPDF_FORMFILLINFO_FFI_GetPage_CB struct {
}

func (cb FPDF_FORMFILLINFO_FFI_GetPage_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])
	document := uint32(stack[1])
	pageIndex := uint32(stack[2])

	// Check if we still have the callback.
	formFillInfoHandle, ok := implementation_webassembly.FormFillInfoHandles[me]
	if !ok {
		return
	}

	stack[0] = formFillInfoHandle.FFI_GetPage(uint64(document), int(pageIndex))
}

type FPDF_FORMFILLINFO_FFI_GetCurrentPage_CB struct {
}

func (cb FPDF_FORMFILLINFO_FFI_GetCurrentPage_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])
	document := uint32(stack[1])

	// Check if we still have the callback.
	formFillInfoHandle, ok := implementation_webassembly.FormFillInfoHandles[me]
	if !ok {
		return
	}

	stack[0] = formFillInfoHandle.FFI_GetCurrentPage(uint64(document))
}

type FPDF_FORMFILLINFO_FFI_GetRotation_CB struct {
}

func (cb FPDF_FORMFILLINFO_FFI_GetRotation_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])
	page := uint32(stack[1])

	// Check if we still have the callback.
	formFillInfoHandle, ok := implementation_webassembly.FormFillInfoHandles[me]
	if !ok {
		return
	}

	stack[0] = api.EncodeI32(int32(formFillInfoHandle.FFI_GetRotation(uint64(page))))
}

type FPDF_FORMFILLINFO_FFI_ExecuteNamedAction_CB struct {
}

func (cb FPDF_FORMFILLINFO_FFI_ExecuteNamedAction_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])
	namedAction := uint32(stack[1])

	// Check if we still have the callback.
	formFillInfoHandle, ok := implementation_webassembly.FormFillInfoHandles[me]
	if !ok {
		return
	}

	formFillInfoHandle.FFI_ExecuteNamedAction(uint64(namedAction))
}

type FPDF_FORMFILLINFO_FFI_SetTextFieldFocus_CB struct {
}

func (cb FPDF_FORMFILLINFO_FFI_SetTextFieldFocus_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])
	value := uint32(stack[1])
	valueLen := uint32(stack[2])
	isFocus := uint32(stack[3])

	// Check if we still have the callback.
	formFillInfoHandle, ok := implementation_webassembly.FormFillInfoHandles[me]
	if !ok {
		return
	}

	formFillInfoHandle.FFI_SetTextFieldFocus(value, valueLen, isFocus)
}

type FPDF_FORMFILLINFO_FFI_DoURIAction_CB struct {
}

func (cb FPDF_FORMFILLINFO_FFI_DoURIAction_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])
	bsURI := uint32(stack[1])

	// Check if we still have the callback.
	formFillInfoHandle, ok := implementation_webassembly.FormFillInfoHandles[me]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_DoURIAction == nil {
		return
	}

	formFillInfoHandle.FFI_DoURIAction(bsURI)
}

type FPDF_FORMFILLINFO_FFI_DoGoToAction_CB struct {
}

func (cb FPDF_FORMFILLINFO_FFI_DoGoToAction_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	me := uint32(stack[0])
	nPageIndex := uint32(stack[1])
	zoomMode := uint32(stack[2])
	fPosArray := uint32(stack[3])
	sizeofArray := uint32(stack[4])

	// Check if we still have the callback.
	formFillInfoHandle, ok := implementation_webassembly.FormFillInfoHandles[me]
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_DoGoToAction == nil {
		return
	}

	formFillInfoHandle.FFI_DoGoToAction(nPageIndex, zoomMode, fPosArray, sizeofArray)
}
