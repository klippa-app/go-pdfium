package implementation_webassembly

import (
	"context"
	"encoding/binary"
	"math"
	"testing"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/structs"
)

type dispatchTestMemory struct {
	data []byte
}

func (m *dispatchTestMemory) Size() uint32 { return uint32(len(m.data)) }

func (m *dispatchTestMemory) ReadUint16Le(offset uint32) (uint16, bool) {
	if int(offset)+2 > len(m.data) {
		return 0, false
	}
	return binary.LittleEndian.Uint16(m.data[offset:]), true
}

func (m *dispatchTestMemory) ReadUint32Le(offset uint32) (uint32, bool) {
	if int(offset)+4 > len(m.data) {
		return 0, false
	}
	return binary.LittleEndian.Uint32(m.data[offset:]), true
}

func (m *dispatchTestMemory) ReadFloat32Le(offset uint32) (float32, bool) {
	v, ok := m.ReadUint32Le(offset)
	return math.Float32frombits(v), ok
}

func (m *dispatchTestMemory) ReadUint64Le(offset uint32) (uint64, bool) {
	if int(offset)+8 > len(m.data) {
		return 0, false
	}
	return binary.LittleEndian.Uint64(m.data[offset:]), true
}

func (m *dispatchTestMemory) ReadFloat64Le(offset uint32) (float64, bool) {
	v, ok := m.ReadUint64Le(offset)
	return math.Float64frombits(v), ok
}

func (m *dispatchTestMemory) Read(offset, byteCount uint32) ([]byte, bool) {
	if int(offset)+int(byteCount) > len(m.data) {
		return nil, false
	}
	return m.data[offset : offset+byteCount], true
}

func (m *dispatchTestMemory) WriteUint16Le(offset uint32, v uint16) bool {
	if int(offset)+2 > len(m.data) {
		return false
	}
	binary.LittleEndian.PutUint16(m.data[offset:], v)
	return true
}

func (m *dispatchTestMemory) WriteUint32Le(offset, v uint32) bool {
	if int(offset)+4 > len(m.data) {
		return false
	}
	binary.LittleEndian.PutUint32(m.data[offset:], v)
	return true
}

func (m *dispatchTestMemory) WriteFloat32Le(offset uint32, v float32) bool {
	return m.WriteUint32Le(offset, math.Float32bits(v))
}

func (m *dispatchTestMemory) WriteUint64Le(offset uint32, v uint64) bool {
	if int(offset)+8 > len(m.data) {
		return false
	}
	binary.LittleEndian.PutUint64(m.data[offset:], v)
	return true
}

func (m *dispatchTestMemory) WriteFloat64Le(offset uint32, v float64) bool {
	return m.WriteUint64Le(offset, math.Float64bits(v))
}

func (m *dispatchTestMemory) Write(offset uint32, v []byte) bool {
	if int(offset)+len(v) > len(m.data) {
		return false
	}
	copy(m.data[offset:], v)
	return true
}

func (m *dispatchTestMemory) WriteString(offset uint32, v string) bool {
	return m.Write(offset, []byte(v))
}

type dispatchTestFunction struct {
	paramCount int
	result     uint64
	calls      [][]uint64
}

func (f *dispatchTestFunction) Call(ctx context.Context, params ...uint64) ([]uint64, error) {
	f.calls = append(f.calls, append([]uint64(nil), params...))
	return []uint64{f.result}, nil
}

func (f *dispatchTestFunction) CallWithStack(ctx context.Context, stack []uint64) error {
	f.calls = append(f.calls, append([]uint64(nil), stack[:f.paramCount]...))
	stack[0] = f.result
	return nil
}

func (f *dispatchTestFunction) ParamCount() int { return f.paramCount }

type dispatchTestModule struct {
	memory    *dispatchTestMemory
	functions map[string]*dispatchTestFunction
}

func (m *dispatchTestModule) Memory() Memory { return m.memory }

func (m *dispatchTestModule) ExportedFunction(name string) Function {
	if fn, ok := m.functions[name]; ok {
		return fn
	}
	return nil
}

func TestFORMOnMouseWheelCallsFORMOnMouseWheel(t *testing.T) {
	malloc := &dispatchTestFunction{paramCount: 1, result: 16}
	free := &dispatchTestFunction{paramCount: 1}
	formOnMouseWheel := &dispatchTestFunction{paramCount: 6, result: 1}
	loadXFA := &dispatchTestFunction{paramCount: 1, result: 1}
	module := &dispatchTestModule{
		memory: &dispatchTestMemory{data: make([]byte, 128)},
		functions: map[string]*dispatchTestFunction{
			"FORM_OnMouseWheel": formOnMouseWheel,
			"FPDF_LoadXFA":      loadXFA,
		},
	}
	implementation := GetInstance(context.Background(), map[string]Function{
		"malloc": malloc,
		"free":   free,
	}, module)

	formHandle := uint64(111)
	pageHandle := uint64(222)
	formRef := references.FPDF_FORMHANDLE("form")
	pageRef := references.FPDF_PAGE("page")
	implementation.formHandleRefs[formRef] = &FormHandleHandle{handle: &formHandle}
	implementation.pageRefs[pageRef] = &PageHandle{handle: &pageHandle}

	response, err := implementation.FORM_OnMouseWheel(&requests.FORM_OnMouseWheel{
		FormHandle: formRef,
		Page:       requests.Page{ByReference: &pageRef},
		Modifier:   5,
		PageCoord:  structs.FPDF_FS_POINTF{X: 1, Y: 2},
		DeltaX:     3,
		DeltaY:     4,
	})
	if err != nil {
		t.Fatalf("FORM_OnMouseWheel returned an error: %v", err)
	}
	if response == nil {
		t.Fatal("FORM_OnMouseWheel returned no response")
	}
	if len(loadXFA.calls) != 0 {
		t.Fatalf("FORM_OnMouseWheel called FPDF_LoadXFA with %v", loadXFA.calls[0])
	}
	if len(formOnMouseWheel.calls) != 1 {
		t.Fatalf("FORM_OnMouseWheel was called %d times, want 1", len(formOnMouseWheel.calls))
	}

	want := []uint64{111, 222, 5, 16, 3, 4}
	got := formOnMouseWheel.calls[0]
	if len(got) != len(want) {
		t.Fatalf("FORM_OnMouseWheel got %d parameters, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("parameter %d = %d, want %d", i, got[i], want[i])
		}
	}
}
