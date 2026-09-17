package implementation_webassembly

import (
	"context"
	"math"
)

// Memory is the linear memory access the implementation needs from a
// WebAssembly runtime. It is a subset of wazero's api.Memory, which satisfies
// it as-is, so that other runtimes can be adapted with a small wrapper.
type Memory interface {
	// Size returns the size in bytes of the linear memory.
	Size() uint32
	ReadUint16Le(offset uint32) (uint16, bool)
	ReadUint32Le(offset uint32) (uint32, bool)
	ReadFloat32Le(offset uint32) (float32, bool)
	ReadUint64Le(offset uint32) (uint64, bool)
	ReadFloat64Le(offset uint32) (float64, bool)
	// Read returns a write-through view of guest memory. The view is only
	// valid until the guest grows its memory.
	Read(offset, byteCount uint32) ([]byte, bool)
	WriteUint16Le(offset uint32, v uint16) bool
	WriteUint32Le(offset, v uint32) bool
	WriteFloat32Le(offset uint32, v float32) bool
	WriteUint64Le(offset uint32, v uint64) bool
	WriteFloat64Le(offset uint32, v float64) bool
	Write(offset uint32, v []byte) bool
	WriteString(offset uint32, v string) bool
}

// Function is an exported guest function. Parameters and results are passed
// as raw 64 bit slots, encoded with the Encode/Decode helpers below. It is a
// subset of wazero's api.Function, which satisfies it as-is.
type Function interface {
	// Call invokes the function and returns a fresh result slice.
	Call(ctx context.Context, params ...uint64) ([]uint64, error)
	// CallWithStack invokes the function with the parameters at the start of
	// stack and writes the results back to the start of stack. The stack must
	// be large enough to hold both.
	CallWithStack(ctx context.Context, stack []uint64) error
	// ParamCount returns the number of parameters the function takes.
	ParamCount() int
}

// Module is an instantiated PDFium WebAssembly module. A Module value must be
// comparable and stable for the lifetime of the instance, it is used as a map
// key to associate host callbacks with the instance that registered them.
type Module interface {
	Memory() Memory
	// ExportedFunction returns the exported function with the given name, or
	// nil when the module does not export it.
	ExportedFunction(name string) Function
}

// EncodeI32 encodes the input as a 64 bit parameter or result slot.
func EncodeI32(input int32) uint64 {
	return uint64(uint32(input))
}

// DecodeI32 decodes a 64 bit slot as a signed 32 bit integer.
func DecodeI32(input uint64) int32 {
	return int32(input)
}

// EncodeU32 encodes the input as a 64 bit parameter or result slot.
func EncodeU32(input uint32) uint64 {
	return uint64(input)
}

// DecodeU32 decodes a 64 bit slot as an unsigned 32 bit integer.
func DecodeU32(input uint64) uint32 {
	return uint32(input)
}

// EncodeI64 encodes the input as a 64 bit parameter or result slot.
func EncodeI64(input int64) uint64 {
	return uint64(input)
}

// DecodeI64 decodes a 64 bit slot as a signed 64 bit integer.
func DecodeI64(input uint64) int64 {
	return int64(input)
}

// EncodeF32 encodes the input as a 64 bit parameter or result slot.
func EncodeF32(input float32) uint64 {
	return uint64(math.Float32bits(input))
}

// DecodeF32 decodes a 64 bit slot as a 32 bit float.
func DecodeF32(input uint64) float32 {
	return math.Float32frombits(uint32(input))
}

// EncodeF64 encodes the input as a 64 bit parameter or result slot.
func EncodeF64(input float64) uint64 {
	return math.Float64bits(input)
}

// DecodeF64 decodes a 64 bit slot as a 64 bit float.
func DecodeF64(input uint64) float64 {
	return math.Float64frombits(input)
}

// ValueType is a WebAssembly value type of a host function parameter or
// result.
type ValueType byte

const (
	ValueTypeI32 ValueType = iota
	ValueTypeI64
	ValueTypeF32
	ValueTypeF64
)

// HostFunction describes one function that PDFium imports from the "env"
// module and that has to be provided by the host, independent of the runtime
// that executes the module. Call receives the parameters at the start of stack
// and writes the result, if any, to stack[0].
type HostFunction struct {
	Name    string
	Params  []ValueType
	Results []ValueType
	Call    func(ctx context.Context, mod Module, stack []uint64)
}
