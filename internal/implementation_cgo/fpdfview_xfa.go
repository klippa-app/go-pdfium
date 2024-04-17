//go:build pdfium_xfa
// +build pdfium_xfa

package implementation_cgo

// #cgo pkg-config: pdfium
// #cgo CFLAGS: -DPDF_ENABLE_XFA
// #include "fpdfview.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"unsafe"
)

// FPDF_BStr_Init initializes a FPDF_BSTR.
//
// Only available when using a build that includes XFA and when using the
// build flag pdfium_xfa.
func (p *PdfiumImplementation) FPDF_BStr_Init(request *requests.FPDF_BStr_Init) (*responses.FPDF_BStr_Init, error) {
	p.Lock()
	defer p.Unlock()

	FPDF_BStr := C.FPDF_BSTR{}

	result := C.FPDF_BStr_Init(&FPDF_BStr)
	if result != C.FPDF_RESULT(0) {
		return nil, errors.New("Could not init FPDF_BSTR")
	}

	handle := p.registerBStr(FPDF_BStr)

	return &responses.FPDF_BStr_Init{
		FPDF_BSTR: handle.nativeRef,
	}, nil
}

// FPDF_BStr_Set copies string data into the FPDF_BSTR.
//
// Only available when using a build that includes XFA and when using the
// build flag pdfium_xfa.
func (p *PdfiumImplementation) FPDF_BStr_Set(request *requests.FPDF_BStr_Set) (*responses.FPDF_BStr_Set, error) {
	p.Lock()
	defer p.Unlock()

	handle, err := p.getBStrHandle(request.FPDF_BSTR)
	if err != nil {
		return nil, err
	}

	valueStr := C.CString(request.Value)
	defer C.free(unsafe.Pointer(valueStr))

	result := C.FPDF_BStr_Set(&handle.handle, valueStr, C.int(len(request.Value)))
	if result != C.FPDF_RESULT(0) {
		return nil, errors.New("Could not set FPDF_BSTR value")
	}

	return &responses.FPDF_BStr_Set{}, nil
}

// FPDF_BStr_Clear clears a FPDF_BSTR.
//
// Only available when using a build that includes XFA and when using the
// build flag pdfium_xfa.
func (p *PdfiumImplementation) FPDF_BStr_Clear(request *requests.FPDF_BStr_Clear) (*responses.FPDF_BStr_Clear, error) {
	p.Lock()
	defer p.Unlock()

	handle, err := p.getBStrHandle(request.FPDF_BSTR)
	if err != nil {
		return nil, err
	}

	result := C.FPDF_BStr_Clear(&handle.handle)
	if result != C.FPDF_RESULT(0) {
		return nil, errors.New("Could not clear FPDF_BSTR")
	}

	return &responses.FPDF_BStr_Clear{}, nil
}
