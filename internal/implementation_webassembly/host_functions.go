package implementation_webassembly

import (
	"context"
	"io"
	"log"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
)

func hostFPDF_FILEACCESS_CB(ctx context.Context, mod Module, stack []uint64) {
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
	FileReaders.Mutex.Lock()
	openFile, ok := FileReaders.Refs[param]
	FileReaders.Mutex.Unlock()
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

	// Memory.Read returns a write-through view of guest memory.
	guestBuffer, ok := mem.Read(pBufPointer, size)
	if !ok {
		stack[0] = uint64(0)
		return
	}

	// m_GetBlock must fill the entire requested range. ReadFull also treats an
	// EOF returned together with the final bytes as a successful read.
	n, err := io.ReadFull(openFile.Reader, guestBuffer)
	if err != nil {
		stack[0] = uint64(0)
		return
	}

	stack[0] = uint64(n)
	return
}

func hostFPDF_FILEWRITE_CB(ctx context.Context, mod Module, stack []uint64) {
	fileWritePointer := uint32(stack[0])
	pDataPointer := uint32(stack[1])
	size := uint32(stack[2])

	mem := mod.Memory()

	key := FileWriterKey{Module: mod, Pointer: fileWritePointer}

	// Check if we have the file referenced in param.
	FileWriters.Mutex.Lock()
	openWriter, ok := FileWriters.Refs[key]
	FileWriters.Mutex.Unlock()
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

func hostFX_FILEAVAIL_IS_DATA_AVAILABLE_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])
	offset := uint32(stack[1])
	size := uint32(stack[2])

	FileAvailables.Mutex.Lock()
	fileAvail, ok := FileAvailables.Refs[me]
	FileAvailables.Mutex.Unlock()
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

func hostFX_DOWNLOADHINTS_ADD_SEGMENT_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])
	offset := uint32(stack[1])
	size := uint32(stack[2])

	FileHints.Mutex.Lock()
	fileHint, ok := FileHints.Refs[me]
	FileHints.Mutex.Unlock()
	if !ok {
		return
	}

	fileHint.AddSegmentCallback(uint64(offset), uint64(size))
}

func hostUNSUPPORT_INFO_HANDLER_CB(ctx context.Context, mod Module, stack []uint64) {
	ntype := uint32(stack[1])

	if CurrentUnsupportedObjectHandler != nil {
		CurrentUnsupportedObjectHandler(enums.FPDF_UNSP(ntype))
	}
}

func hostFSDK_SetTimeFunction_CB(ctx context.Context, mod Module, stack []uint64) {
	currentTime := uint64(0)
	if CurrentTimeHandler != nil {
		currentTime = EncodeI64(CurrentTimeHandler())
	}

	stack[0] = currentTime
	return
}

// re-use memory to prevent allocating more than necessary.
var lastFSDK_SetLocaltimeFunctionPointer *uint64

func hostFSDK_SetLocaltimeFunction_CB(ctx context.Context, mod Module, stack []uint64) {
	timestamp := uint32(stack[0])

	currentLocalTime := uint64(0)
	if CurrentLocalTimeHandler != nil {
		localTime := CurrentLocalTimeHandler(int64(timestamp))

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

		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer), EncodeI32(int32(localTime.TmSec)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+4), EncodeI32(int32(localTime.TmMin)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+8), EncodeI32(int32(localTime.TmHour)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+12), EncodeI32(int32(localTime.TmMday)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+16), EncodeI32(int32(localTime.TmMon)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+20), EncodeI32(int32(localTime.TmYear)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+24), EncodeI32(int32(localTime.TmWday)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+28), EncodeI32(int32(localTime.TmYday)))
		mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+32), EncodeI32(int32(localTime.TmIsdst)))

		currentLocalTime = *lastFSDK_SetLocaltimeFunctionPointer
	}

	stack[0] = currentLocalTime
	return
}

func hostFPDF_FORMFILLINFO_Release_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])

	// Check if we still have the callback.
	FormFillInfoHandles.Mutex.Lock()
	formFillInfoHandle, ok := FormFillInfoHandles.Refs[me]
	FormFillInfoHandles.Mutex.Unlock()
	if !ok {
		return
	}

	// @todo: do I have anything to cleanup for myself?

	formFillInfoHandle.Release()
}

func hostFPDF_FORMFILLINFO_FFI_Invalidate_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])
	page := uint32(stack[1])
	left := uint64(stack[2])
	top := uint64(stack[3])
	right := uint64(stack[4])
	bottom := uint64(stack[5])

	// Check if we still have the callback.
	FormFillInfoHandles.Mutex.Lock()
	formFillInfoHandle, ok := FormFillInfoHandles.Refs[me]
	FormFillInfoHandles.Mutex.Unlock()
	if !ok {
		return
	}

	formFillInfoHandle.FFI_Invalidate_CB(page, left, top, right, bottom)
}

func hostFPDF_FORMFILLINFO_FFI_OutputSelectedRect_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])
	page := uint32(stack[1])
	left := uint64(stack[2])
	top := uint64(stack[3])
	right := uint64(stack[4])
	bottom := uint64(stack[5])

	// Check if we still have the callback.
	FormFillInfoHandles.Mutex.Lock()
	formFillInfoHandle, ok := FormFillInfoHandles.Refs[me]
	FormFillInfoHandles.Mutex.Unlock()
	if !ok {
		return
	}

	formFillInfoHandle.FFI_OutputSelectedRect(page, left, top, right, bottom)
}

func hostFPDF_FORMFILLINFO_FFI_SetCursor_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])
	cursor := uint32(stack[1])

	// Check if we still have the callback.
	FormFillInfoHandles.Mutex.Lock()
	formFillInfoHandle, ok := FormFillInfoHandles.Refs[me]
	FormFillInfoHandles.Mutex.Unlock()
	if !ok {
		return
	}

	formFillInfoHandle.FFI_SetCursor(cursor)
}

func hostFPDF_FORMFILLINFO_FFI_SetTimer_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])
	uElapse := uint32(stack[1])
	lpTimerFunc := uint32(stack[2])

	// Check if we still have the callback.
	FormFillInfoHandles.Mutex.Lock()
	formFillInfoHandle, ok := FormFillInfoHandles.Refs[me]
	FormFillInfoHandles.Mutex.Unlock()
	if !ok {
		return
	}

	id := formFillInfoHandle.FFI_SetTimer(uElapse, lpTimerFunc)
	stack[0] = uint64(id)
}

func hostFPDF_FORMFILLINFO_FFI_KillTimer_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])
	nTimerID := uint32(stack[1])

	// Check if we still have the callback.
	FormFillInfoHandles.Mutex.Lock()
	formFillInfoHandle, ok := FormFillInfoHandles.Refs[me]
	FormFillInfoHandles.Mutex.Unlock()
	if !ok {
		return
	}

	formFillInfoHandle.FFI_KillTimer(int(nTimerID))
}

// re-use memory to prevent allocating more than necessary.
var lastFPDF_FORMFILLINFO_FFI_GetLocalTimePointer *uint64

func hostFPDF_FORMFILLINFO_FFI_GetLocalTime_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])

	// Check if we still have the callback.
	FormFillInfoHandles.Mutex.Lock()
	formFillInfoHandle, ok := FormFillInfoHandles.Refs[me]
	FormFillInfoHandles.Mutex.Unlock()
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
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer), EncodeI32(int32(localTime.Year)))
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+2), EncodeI32(int32(localTime.Month)))
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+4), EncodeI32(int32(localTime.DayOfWeek)))
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+6), EncodeI32(int32(localTime.Day)))
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+8), EncodeI32(int32(localTime.Hour)))
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+10), EncodeI32(int32(localTime.Minute)))
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+12), EncodeI32(int32(localTime.Second)))
	mod.Memory().WriteUint64Le(uint32(*lastFSDK_SetLocaltimeFunctionPointer+14), EncodeI32(int32(localTime.Milliseconds)))

	stack[0] = uint64(*lastFSDK_SetLocaltimeFunctionPointer)
}

func hostFPDF_FORMFILLINFO_FFI_OnChange_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])

	// Check if we still have the callback.
	FormFillInfoHandles.Mutex.Lock()
	formFillInfoHandle, ok := FormFillInfoHandles.Refs[me]
	FormFillInfoHandles.Mutex.Unlock()
	if !ok {
		return
	}

	formFillInfoHandle.FFI_OnChange()
}

func hostFPDF_FORMFILLINFO_FFI_GetPage_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])
	document := uint32(stack[1])
	pageIndex := uint32(stack[2])

	// Check if we still have the callback.
	FormFillInfoHandles.Mutex.Lock()
	formFillInfoHandle, ok := FormFillInfoHandles.Refs[me]
	FormFillInfoHandles.Mutex.Unlock()
	if !ok {
		return
	}

	stack[0] = formFillInfoHandle.FFI_GetPage(uint64(document), int(pageIndex))
}

func hostFPDF_FORMFILLINFO_FFI_GetCurrentPage_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])
	document := uint32(stack[1])

	// Check if we still have the callback.
	FormFillInfoHandles.Mutex.Lock()
	formFillInfoHandle, ok := FormFillInfoHandles.Refs[me]
	FormFillInfoHandles.Mutex.Unlock()
	if !ok {
		return
	}

	stack[0] = formFillInfoHandle.FFI_GetCurrentPage(uint64(document))
}

func hostFPDF_FORMFILLINFO_FFI_GetRotation_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])
	page := uint32(stack[1])

	// Check if we still have the callback.
	FormFillInfoHandles.Mutex.Lock()
	formFillInfoHandle, ok := FormFillInfoHandles.Refs[me]
	FormFillInfoHandles.Mutex.Unlock()
	if !ok {
		return
	}

	stack[0] = EncodeI32(int32(formFillInfoHandle.FFI_GetRotation(uint64(page))))
}

func hostFPDF_FORMFILLINFO_FFI_ExecuteNamedAction_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])
	namedAction := uint32(stack[1])

	// Check if we still have the callback.
	FormFillInfoHandles.Mutex.Lock()
	formFillInfoHandle, ok := FormFillInfoHandles.Refs[me]
	FormFillInfoHandles.Mutex.Unlock()
	if !ok {
		return
	}

	formFillInfoHandle.FFI_ExecuteNamedAction(uint64(namedAction))
}

func hostFPDF_FORMFILLINFO_FFI_SetTextFieldFocus_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])
	value := uint32(stack[1])
	valueLen := uint32(stack[2])
	isFocus := uint32(stack[3])

	// Check if we still have the callback.
	FormFillInfoHandles.Mutex.Lock()
	formFillInfoHandle, ok := FormFillInfoHandles.Refs[me]
	FormFillInfoHandles.Mutex.Unlock()
	if !ok {
		return
	}

	formFillInfoHandle.FFI_SetTextFieldFocus(value, valueLen, isFocus)
}

func hostFPDF_FORMFILLINFO_FFI_DoURIAction_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])
	bsURI := uint32(stack[1])

	// Check if we still have the callback.
	FormFillInfoHandles.Mutex.Lock()
	formFillInfoHandle, ok := FormFillInfoHandles.Refs[me]
	FormFillInfoHandles.Mutex.Unlock()
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_DoURIAction == nil {
		return
	}

	formFillInfoHandle.FFI_DoURIAction(bsURI)
}

func hostFPDF_FORMFILLINFO_FFI_DoGoToAction_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])
	nPageIndex := uint32(stack[1])
	zoomMode := uint32(stack[2])
	fPosArray := uint32(stack[3])
	sizeofArray := uint32(stack[4])

	// Check if we still have the callback.
	FormFillInfoHandles.Mutex.Lock()
	formFillInfoHandle, ok := FormFillInfoHandles.Refs[me]
	FormFillInfoHandles.Mutex.Unlock()
	if !ok {
		return
	}

	if formFillInfoHandle.FormFillInfo.FFI_DoGoToAction == nil {
		return
	}

	formFillInfoHandle.FFI_DoGoToAction(nPageIndex, zoomMode, fPosArray, sizeofArray)
}

func hostIFSDK_PAUSE_NeedToPauseNow_CB(ctx context.Context, mod Module, stack []uint64) {
	me := uint32(stack[0])

	stringRefPointer, _ := mod.Memory().ReadUint32Le(me + 8)

	stringRef := []byte{}
	for {
		data, success := mod.Memory().Read(stringRefPointer, 1)
		if !success {
			return
		}

		if data[0] == 0x00 {
			break
		}

		stringRef = append(stringRef, data[0])
		stringRefPointer++
	}

	// Check if we still have the reference.
	PauseHandles.Mutex.RLock()
	defer PauseHandles.Mutex.RUnlock()
	if _, ok := PauseHandles.Refs[references.FPDF_PAGE(string(stringRef))]; !ok {
		stack[0] = EncodeI32(int32(1))
		return
	}

	shouldPause := PauseHandles.Refs[references.FPDF_PAGE(string(stringRef))].Callback()
	if shouldPause {
		stack[0] = EncodeI32(int32(1))
		return
	}

	stack[0] = EncodeI32(int32(0))
}

// HostFunctions lists every function PDFium imports from the "env" module,
// with the signature the module expects. Runtimes register these as host
// functions when they instantiate the module.
var HostFunctions = []HostFunction{
	{Name: "FPDF_FILEACCESS_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32, ValueTypeI32, ValueTypeI32}, Results: []ValueType{ValueTypeI32}, Call: hostFPDF_FILEACCESS_CB},
	{Name: "FPDF_FILEWRITE_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32, ValueTypeI32}, Results: []ValueType{ValueTypeI32}, Call: hostFPDF_FILEWRITE_CB},
	{Name: "FX_FILEAVAIL_IS_DATA_AVAILABLE_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32, ValueTypeI32}, Results: []ValueType{ValueTypeI32}, Call: hostFX_FILEAVAIL_IS_DATA_AVAILABLE_CB},
	{Name: "FX_DOWNLOADHINTS_ADD_SEGMENT_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32, ValueTypeI32}, Results: []ValueType{}, Call: hostFX_DOWNLOADHINTS_ADD_SEGMENT_CB},
	{Name: "UNSUPPORT_INFO_HANDLER_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32}, Results: []ValueType{}, Call: hostUNSUPPORT_INFO_HANDLER_CB},
	{Name: "FSDK_SetTimeFunction_CB", Params: []ValueType{}, Results: []ValueType{ValueTypeI64}, Call: hostFSDK_SetTimeFunction_CB},
	{Name: "FSDK_SetLocaltimeFunction_CB", Params: []ValueType{ValueTypeI32}, Results: []ValueType{ValueTypeI32}, Call: hostFSDK_SetLocaltimeFunction_CB},
	{Name: "FPDF_FORMFILLINFO_Release_CB", Params: []ValueType{ValueTypeI32}, Results: []ValueType{}, Call: hostFPDF_FORMFILLINFO_Release_CB},
	{Name: "FPDF_FORMFILLINFO_FFI_Invalidate_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32, ValueTypeF64, ValueTypeF64, ValueTypeF64, ValueTypeF64}, Results: []ValueType{}, Call: hostFPDF_FORMFILLINFO_FFI_Invalidate_CB},
	{Name: "FPDF_FORMFILLINFO_FFI_OutputSelectedRect_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32, ValueTypeF64, ValueTypeF64, ValueTypeF64, ValueTypeF64}, Results: []ValueType{}, Call: hostFPDF_FORMFILLINFO_FFI_OutputSelectedRect_CB},
	{Name: "FPDF_FORMFILLINFO_FFI_SetCursor_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32}, Results: []ValueType{}, Call: hostFPDF_FORMFILLINFO_FFI_SetCursor_CB},
	{Name: "FPDF_FORMFILLINFO_FFI_SetTimer_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32, ValueTypeI32}, Results: []ValueType{ValueTypeI32}, Call: hostFPDF_FORMFILLINFO_FFI_SetTimer_CB},
	{Name: "FPDF_FORMFILLINFO_FFI_KillTimer_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32}, Results: []ValueType{}, Call: hostFPDF_FORMFILLINFO_FFI_KillTimer_CB},
	{Name: "FPDF_FORMFILLINFO_FFI_GetLocalTime_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32}, Results: []ValueType{}, Call: hostFPDF_FORMFILLINFO_FFI_GetLocalTime_CB},
	{Name: "FPDF_FORMFILLINFO_FFI_OnChange_CB", Params: []ValueType{ValueTypeI32}, Results: []ValueType{}, Call: hostFPDF_FORMFILLINFO_FFI_OnChange_CB},
	{Name: "FPDF_FORMFILLINFO_FFI_GetPage_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32, ValueTypeI32}, Results: []ValueType{ValueTypeI32}, Call: hostFPDF_FORMFILLINFO_FFI_GetPage_CB},
	{Name: "FPDF_FORMFILLINFO_FFI_GetCurrentPage_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32}, Results: []ValueType{ValueTypeI32}, Call: hostFPDF_FORMFILLINFO_FFI_GetCurrentPage_CB},
	{Name: "FPDF_FORMFILLINFO_FFI_GetRotation_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32}, Results: []ValueType{ValueTypeI32}, Call: hostFPDF_FORMFILLINFO_FFI_GetRotation_CB},
	{Name: "FPDF_FORMFILLINFO_FFI_ExecuteNamedAction_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32}, Results: []ValueType{}, Call: hostFPDF_FORMFILLINFO_FFI_ExecuteNamedAction_CB},
	{Name: "FPDF_FORMFILLINFO_FFI_SetTextFieldFocus_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32, ValueTypeI32, ValueTypeI32}, Results: []ValueType{}, Call: hostFPDF_FORMFILLINFO_FFI_SetTextFieldFocus_CB},
	{Name: "FPDF_FORMFILLINFO_FFI_DoURIAction_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32}, Results: []ValueType{}, Call: hostFPDF_FORMFILLINFO_FFI_DoURIAction_CB},
	{Name: "FPDF_FORMFILLINFO_FFI_DoGoToAction_CB", Params: []ValueType{ValueTypeI32, ValueTypeI32, ValueTypeI32, ValueTypeI32, ValueTypeI32}, Results: []ValueType{}, Call: hostFPDF_FORMFILLINFO_FFI_DoGoToAction_CB},
	{Name: "IFSDK_PAUSE_NeedToPauseNow_CB", Params: []ValueType{ValueTypeI32}, Results: []ValueType{ValueTypeI32}, Call: hostIFSDK_PAUSE_NeedToPauseNow_CB},
}
