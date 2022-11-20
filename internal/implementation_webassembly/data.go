package implementation_webassembly

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/klippa-app/go-pdfium/structs"
)

type CString struct {
	Pointer uint64
	Free    func()
}

func (p *PdfiumImplementation) CString(input string) (*CString, error) {
	inputLength := uint64(len(input)) + 1

	pointer, err := p.Malloc(inputLength)
	if err != nil {
		return nil, err
	}

	// Write string + null terminator.
	if !p.Module.Memory().Write(p.Context, uint32(pointer), append([]byte(input), byte(0))) {
		return nil, errors.New("could not write CString data")
	}

	return &CString{
		Pointer: pointer,
		Free: func() {
			p.Free(pointer)
		},
	}, nil
}

type IntPointer struct {
	Pointer uint64
	Free    func()
	Value   func() (int, error)
}

func (p *PdfiumImplementation) IntPointer() (*IntPointer, error) {
	pointer, err := p.Malloc(p.CSizeInt())
	if err != nil {
		return nil, err
	}

	return &IntPointer{
		Pointer: pointer,
		Free: func() {
			p.Free(pointer)
		},
		Value: func() (int, error) {
			b, success := p.Module.Memory().Read(p.Context, uint32(pointer), uint32(4))
			if !success {
				return 0, errors.New("could not read int data from memory")
			}

			var myInt int32
			err := binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &myInt)
			if err != nil {
				return 0, err
			}

			return int(myInt), nil
		},
	}, nil
}

type DoublePointer struct {
	Pointer uint64
	Free    func()
	Value   func() (float64, error)
}

func (p *PdfiumImplementation) DoublePointer() (*DoublePointer, error) {
	pointer, err := p.Malloc(p.CSizeDouble())
	if err != nil {
		return nil, err
	}

	return &DoublePointer{
		Pointer: pointer,
		Free: func() {
			p.Free(pointer)
		},
		Value: func() (float64, error) {
			val, success := p.Module.Memory().ReadFloat64Le(p.Context, uint32(pointer))
			if !success {
				return 0, errors.New("could not read double data from memory")
			}

			return val, nil
		},
	}, nil
}

type ByteArrayPointer struct {
	Pointer uint64
	Free    func()
	Value   func() ([]byte, error)
}

func (p *PdfiumImplementation) ByteArrayPointer(size uint64) (*ByteArrayPointer, error) {
	pointer, err := p.Malloc(size)
	if err != nil {
		return nil, err
	}

	return &ByteArrayPointer{
		Pointer: pointer,
		Free: func() {
			p.Free(pointer)
		},
		Value: func() ([]byte, error) {
			b, success := p.Module.Memory().Read(p.Context, uint32(pointer), uint32(size))
			if !success {
				return nil, errors.New("could not read byte array data from memory")
			}

			return b, nil
		},
	}, nil
}

type LongPointer struct {
	Pointer uint64
	Free    func()
	Value   func() (int64, error)
}

func (p *PdfiumImplementation) LongPointer() (*LongPointer, error) {
	pointer, err := p.Malloc(p.CSizeLong())
	if err != nil {
		return nil, err
	}

	return &LongPointer{
		Pointer: pointer,
		Free: func() {
			p.Free(pointer)
		},
		Value: func() (int64, error) {
			b, success := p.Module.Memory().Read(p.Context, uint32(pointer), uint32(p.CSizeLong()))
			if !success {
				return 0, errors.New("could not read long data from memory")
			}

			var myInt int64
			err := binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &myInt)
			if err != nil {
				return 0, err
			}

			return myInt, nil
		},
	}, nil
}

func (p *PdfiumImplementation) Malloc(size uint64) (uint64, error) {
	results, err := p.Functions["malloc"].Call(p.Context, size)
	if err != nil {
		return 0, err
	}

	pointer := results[0]

	// Create a view of the underlying memory.
	memoryBuffer, ok := p.Module.Memory().Read(p.Context, uint32(results[0]), uint32(size))
	if !ok {
		return 0, errors.New("could not get a view of the memory")
	}

	// Make all values 0.
	// @todo: why do we need to this this? Slow!
	for i := range memoryBuffer {
		memoryBuffer[i] = 0
	}

	return pointer, nil
}

func (p *PdfiumImplementation) Free(pointer uint64) error {
	_, err := p.Functions["free"].Call(p.Context, pointer)
	if err != nil {
		return err
	}
	return nil
}

func (p *PdfiumImplementation) CSizeInt() uint64 {
	// @todo: implement on pdfium/emscripten side?
	return 4
}

func (p *PdfiumImplementation) CSizeFloat() uint64 {
	// @todo: implement on pdfium/emscripten side?
	return 4
}

func (p *PdfiumImplementation) CSizeDouble() uint64 {
	// @todo: implement on pdfium/emscripten side?
	return 8
}

func (p *PdfiumImplementation) CSizeLong() uint64 {
	// @todo: implement on pdfium/emscripten side?
	return 8
}

func (p *PdfiumImplementation) CSizeStructFS_MATRIX() uint64 {
	// FS_MATRIX is 6 * float (a, b, c, d, e, f).
	return p.CSizeFloat() * 6
}

func (p *PdfiumImplementation) CStructFS_MATRIX(in structs.FPDF_FS_MATRIX) (uint64, error) {
	pointer, err := p.Malloc(p.CSizeStructFS_MATRIX())
	if err != nil {
		return 0, err
	}

	if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.A) {
		p.Free(pointer)
		return 0, errors.New("could not write float data to memory")
	}

	if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.B) {
		p.Free(pointer)
		return 0, errors.New("could not write float data to memory")
	}

	if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.C) {
		p.Free(pointer)
		return 0, errors.New("could not write float data to memory")
	}

	if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.D) {
		p.Free(pointer)
		return 0, errors.New("could not write float data to memory")
	}

	if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.E) {
		p.Free(pointer)
		return 0, errors.New("could not write float data to memory")
	}

	if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.F) {
		p.Free(pointer)
		return 0, errors.New("could not write float data to memory")
	}

	return pointer, nil
}

func (p *PdfiumImplementation) CSizeStructFS_RECTF() uint64 {
	// FS_RECTF is 4 * float (left, top, right, bottom).
	return p.CSizeFloat() * 4
}

func (p *PdfiumImplementation) CStructFS_RECTF(in structs.FPDF_FS_RECTF) (uint64, error) {
	pointer, err := p.Malloc(p.CSizeStructFS_RECTF())
	if err != nil {
		return 0, err
	}

	if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.Top) {
		p.Free(pointer)
		return 0, errors.New("could not write float data to memory")
	}

	if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.Left) {
		p.Free(pointer)
		return 0, errors.New("could not write float data to memory")
	}

	if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.Right) {
		p.Free(pointer)
		return 0, errors.New("could not write float data to memory")
	}

	if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.Bottom) {
		p.Free(pointer)
		return 0, errors.New("could not write float data to memory")
	}

	return pointer, nil
}
