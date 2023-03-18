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
var lastLocalTimePointer *uint64

func (cb FSDK_SetLocaltimeFunction_CB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	timestamp := uint32(stack[0])

	currentLocalTime := uint64(0)
	if implementation_webassembly.CurrentLocalTimeHandler != nil {
		localTime := implementation_webassembly.CurrentLocalTimeHandler(int64(timestamp))

		if lastLocalTimePointer == nil {
			// 9 int fields in localtime.
			results, err := mod.ExportedFunction("malloc").Call(ctx, 4*9)
			if err != nil {
				log.Printf("Could not allocate memory")
				stack[0] = currentLocalTime
				return
			}
			lastLocalTimePointer = &results[0]
		}

		mod.Memory().WriteUint64Le(uint32(*lastLocalTimePointer), api.EncodeI32(int32(localTime.TmSec)))
		mod.Memory().WriteUint64Le(uint32(*lastLocalTimePointer+4), api.EncodeI32(int32(localTime.TmMin)))
		mod.Memory().WriteUint64Le(uint32(*lastLocalTimePointer+8), api.EncodeI32(int32(localTime.TmHour)))
		mod.Memory().WriteUint64Le(uint32(*lastLocalTimePointer+12), api.EncodeI32(int32(localTime.TmMday)))
		mod.Memory().WriteUint64Le(uint32(*lastLocalTimePointer+16), api.EncodeI32(int32(localTime.TmMon)))
		mod.Memory().WriteUint64Le(uint32(*lastLocalTimePointer+20), api.EncodeI32(int32(localTime.TmYear)))
		mod.Memory().WriteUint64Le(uint32(*lastLocalTimePointer+24), api.EncodeI32(int32(localTime.TmWday)))
		mod.Memory().WriteUint64Le(uint32(*lastLocalTimePointer+28), api.EncodeI32(int32(localTime.TmYday)))
		mod.Memory().WriteUint64Le(uint32(*lastLocalTimePointer+32), api.EncodeI32(int32(localTime.TmIsdst)))

		currentLocalTime = *lastLocalTimePointer
	}

	stack[0] = currentLocalTime
	return
}
