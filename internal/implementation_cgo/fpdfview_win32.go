//go:build windows
// +build windows

package implementation_cgo

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"
import (
	"errors"
	"syscall"
	"unsafe"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// FPDF_RenderPage renders contents of a page to a device (screen, bitmap, or printer).
// This feature does not work on multi-threaded usage as you will need to give a device handle.
// Windows only!
func (p *PdfiumImplementation) FPDF_RenderPage(request *requests.FPDF_RenderPage) (*responses.FPDF_RenderPage, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	hdc, err := deviceContext(request.DC)
	if err != nil {
		return nil, err
	}

	C.FPDF_RenderPage(hdc, pageHandle.handle, C.int(request.StartX), C.int(request.StartY), C.int(request.SizeX), C.int(request.SizeY), C.int(request.Rotate), C.int(request.Flags))

	return &responses.FPDF_RenderPage{}, nil
}

// deviceContext converts the device context handle a caller gives us into
// the C type. cgo types are unexported and private to the package that
// imports "C", so callers outside this package can never produce a C.HDC;
// they hand in the handle as it comes from the Windows API: a uintptr or
// syscall.Handle (for example from GetDC, CreateCompatibleDC or
// golang.org/x/sys/windows), or an unsafe.Pointer.
func deviceContext(dc any) (C.HDC, error) {
	var handle uintptr
	switch v := dc.(type) {
	case C.HDC:
		return v, nil
	case uintptr:
		handle = v
	case syscall.Handle:
		handle = uintptr(v)
	case unsafe.Pointer:
		handle = uintptr(v)
	default:
		return nil, errors.New("DC must be a Windows device context handle given as uintptr, syscall.Handle or unsafe.Pointer")
	}

	if handle == 0 {
		return nil, errors.New("DC is a null device context handle")
	}

	// An HDC is an opaque Windows handle, not a pointer into Go memory, so
	// reinterpreting the handle value as the pointer typed C.HDC is safe.
	return *(*C.HDC)(unsafe.Pointer(&handle)), nil
}
