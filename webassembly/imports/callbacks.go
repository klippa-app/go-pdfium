package imports

import (
	"context"
	"github.com/klippa-app/go-pdfium/internal/implementation_webassembly"
	"github.com/tetratelabs/wazero/api"
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
