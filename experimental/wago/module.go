package wago

import (
	"context"
	"encoding/binary"
	"errors"
	"math"

	"github.com/klippa-app/go-pdfium/internal/implementation_webassembly"

	wagort "github.com/wago-org/wago"
)

// module adapts a wago instance to the runtime independent Module interface
// of the shared implementation. One module belongs to exactly one worker and
// all access to it is serialized by the instance mutex, or happens from a
// host function that the instance itself is executing.
type module struct {
	ctx      context.Context
	instance *wagort.Instance
	compiled *wagort.Compiled
	memory   *wagort.Memory

	// interruptible makes calls go through InvokeContext with the worker
	// context, so that cancelling it (Kill) interrupts a running call.
	interruptible bool

	// caller is the host caller of the host function that is currently
	// executing, if any. While it is set, memory has to be accessed through
	// it and calls into the guest have to re-enter through it, wago rejects
	// direct access while an invocation is active.
	caller       wagort.Caller
	callerActive bool
}

// enterHostCall marks the module as executing a host function on behalf of
// caller and returns a function that restores the previous state.
func (m *module) enterHostCall(caller wagort.Caller) func() {
	previousCaller, previousActive := m.caller, m.callerActive
	m.caller, m.callerActive = caller, true
	return func() {
		m.caller, m.callerActive = previousCaller, previousActive
	}
}

func (m *module) Memory() implementation_webassembly.Memory {
	return memory{m: m}
}

func (m *module) ExportedFunction(name string) implementation_webassembly.Function {
	params, results, err := m.compiled.Signature(name)
	if err != nil {
		return nil
	}

	fn, err := m.instance.WasmFunc(name)
	if err != nil {
		return nil
	}

	return &function{
		m:       m,
		name:    name,
		fn:      fn,
		params:  len(params),
		results: len(results),
	}
}

// bytes returns the current linear memory. The slice is only valid until the
// guest grows its memory, callers must not retain it across guest calls.
func (m *module) bytes() []byte {
	if m.callerActive {
		return m.caller.Memory()
	}

	return m.memory.UnsafeBytes()
}

type memory struct {
	m *module
}

func (mem memory) slice(offset, byteCount uint32) ([]byte, bool) {
	b := mem.m.bytes()
	end := uint64(offset) + uint64(byteCount)
	if end > uint64(len(b)) {
		return nil, false
	}

	return b[offset:end:end], true
}

func (mem memory) Size() uint32 {
	return uint32(len(mem.m.bytes()))
}

func (mem memory) ReadUint16Le(offset uint32) (uint16, bool) {
	b, ok := mem.slice(offset, 2)
	if !ok {
		return 0, false
	}
	return binary.LittleEndian.Uint16(b), true
}

func (mem memory) ReadUint32Le(offset uint32) (uint32, bool) {
	b, ok := mem.slice(offset, 4)
	if !ok {
		return 0, false
	}
	return binary.LittleEndian.Uint32(b), true
}

func (mem memory) ReadFloat32Le(offset uint32) (float32, bool) {
	v, ok := mem.ReadUint32Le(offset)
	if !ok {
		return 0, false
	}
	return math.Float32frombits(v), true
}

func (mem memory) ReadUint64Le(offset uint32) (uint64, bool) {
	b, ok := mem.slice(offset, 8)
	if !ok {
		return 0, false
	}
	return binary.LittleEndian.Uint64(b), true
}

func (mem memory) ReadFloat64Le(offset uint32) (float64, bool) {
	v, ok := mem.ReadUint64Le(offset)
	if !ok {
		return 0, false
	}
	return math.Float64frombits(v), true
}

func (mem memory) Read(offset, byteCount uint32) ([]byte, bool) {
	return mem.slice(offset, byteCount)
}

func (mem memory) WriteUint16Le(offset uint32, v uint16) bool {
	b, ok := mem.slice(offset, 2)
	if !ok {
		return false
	}
	binary.LittleEndian.PutUint16(b, v)
	return true
}

func (mem memory) WriteUint32Le(offset, v uint32) bool {
	b, ok := mem.slice(offset, 4)
	if !ok {
		return false
	}
	binary.LittleEndian.PutUint32(b, v)
	return true
}

func (mem memory) WriteFloat32Le(offset uint32, v float32) bool {
	return mem.WriteUint32Le(offset, math.Float32bits(v))
}

func (mem memory) WriteUint64Le(offset uint32, v uint64) bool {
	b, ok := mem.slice(offset, 8)
	if !ok {
		return false
	}
	binary.LittleEndian.PutUint64(b, v)
	return true
}

func (mem memory) WriteFloat64Le(offset uint32, v float64) bool {
	return mem.WriteUint64Le(offset, math.Float64bits(v))
}

func (mem memory) Write(offset uint32, v []byte) bool {
	b, ok := mem.slice(offset, uint32(len(v)))
	if !ok {
		return false
	}
	copy(b, v)
	return true
}

func (mem memory) WriteString(offset uint32, v string) bool {
	b, ok := mem.slice(offset, uint32(len(v)))
	if !ok {
		return false
	}
	copy(b, v)
	return true
}

// function is an exported guest function.
type function struct {
	m       *module
	name    string
	fn      *wagort.WasmFunc
	params  int
	results int
}

func (f *function) ParamCount() int {
	return f.params
}

// invoke calls the function. The returned slice is borrowed from wago and is
// only valid until the next call on the instance.
func (f *function) invoke(ctx context.Context, params []uint64) ([]uint64, error) {
	if len(params) < f.params {
		return nil, errors.New("not enough parameters for " + f.name)
	}
	params = params[:f.params]

	// A host function that is running on behalf of the guest has to re-enter
	// the guest through its caller.
	if f.m.callerActive {
		return f.m.instance.InvokeFromHost(ctx, f.m.caller, f.name, params...)
	}

	if f.m.interruptible && ctx != nil && ctx.Done() != nil {
		return f.m.instance.InvokeContext(ctx, f.name, params...)
	}

	return f.fn.Invoke(params...)
}

func (f *function) Call(ctx context.Context, params ...uint64) ([]uint64, error) {
	results, err := f.invoke(ctx, params)
	if err != nil {
		return nil, err
	}

	return append([]uint64(nil), results...), nil
}

func (f *function) CallWithStack(ctx context.Context, stack []uint64) error {
	results, err := f.invoke(ctx, stack)
	if err != nil {
		return err
	}

	if f.results > 0 && len(results) > 0 {
		stack[0] = results[0]
	}

	return nil
}
