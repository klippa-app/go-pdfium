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

	param, ok := mem.ReadUint32Le(ctx, paramPointer)
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

	ok = mem.Write(ctx, pBufPointer, readBuffer)
	if err != nil {
		stack[0] = uint64(0)
		return
	}

	stack[0] = uint64(n)
	return
}
