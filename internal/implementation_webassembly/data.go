package implementation_webassembly

import (
	"bytes"
	"encoding/binary"
	"errors"
	"unsafe"
)

type CString struct {
	Pointer uint64
	Free    func()
}

func (p *PdfiumImplementation) CString(input string) (*CString, error) {
	inputLength := uint64(len(input)) + 1

	results, err := p.Functions["malloc"].Call(p.Context, inputLength)
	if err != nil {
		return nil, err
	}
	pointer := results[0]

	// Write string + null terminator.
	if !p.Module.Memory().Write(p.Context, uint32(pointer), append([]byte(input), byte(0))) {
		return nil, errors.New("could not write CString data")
	}

	return &CString{
		Pointer: pointer,
		Free: func() {
			p.Functions["free"].Call(p.Context, pointer)
		},
	}, nil
}

type IntPointer struct {
	Pointer uint64
	Free    func()
	Value   func() (int, error)
}

func (p *PdfiumImplementation) IntPointer() (*IntPointer, error) {
	results, err := p.Functions["malloc"].Call(p.Context, 4)
	if err != nil {
		return nil, err
	}
	pointer := results[0]

	return &IntPointer{
		Pointer: pointer,
		Free: func() {
			p.Functions["free"].Call(p.Context, pointer)
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
	results, err := p.Functions["malloc"].Call(p.Context, 8)
	if err != nil {
		return nil, err
	}
	pointer := results[0]

	return &DoublePointer{
		Pointer: pointer,
		Free: func() {
			p.Functions["free"].Call(p.Context, pointer)
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
	results, err := p.Functions["malloc"].Call(p.Context, size)
	if err != nil {
		return nil, err
	}
	pointer := results[0]

	return &ByteArrayPointer{
		Pointer: pointer,
		Free: func() {
			p.Functions["free"].Call(p.Context, pointer)
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
	results, err := p.Functions["malloc"].Call(p.Context, 8)
	if err != nil {
		return nil, err
	}
	pointer := results[0]

	return &LongPointer{
		Pointer: pointer,
		Free: func() {
			p.Functions["free"].Call(p.Context, pointer)
		},
		Value: func() (int64, error) {
			b, success := p.Module.Memory().Read(p.Context, uint32(pointer), uint32(8))
			if !success {
				return 0, errors.New("could not read long data from memory")
			}

			val := *(*int64)(unsafe.Pointer(&b[0]))
			return val, nil
		},
	}, nil
}
