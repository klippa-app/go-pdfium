//go:build pdfium_v8
// +build pdfium_v8

package implementation_cgo

import (
	pdfium_errors "github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"log"
)

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include <stdlib.h>
import "C"

// FPDF_GetRecommendedV8Flags returns a space-separated string of command
// line flags that are recommended to be passed into V8 via
// V8::SetFlagsFromString() prior to initializing the PDFium library.
//
// Only available when using a build that includes V8 and when using the
// build flag pdfium_v8.
func (p *PdfiumImplementation) FPDF_GetRecommendedV8Flags(request *requests.FPDF_GetRecommendedV8Flags) (*responses.FPDF_GetRecommendedV8Flags, error) {
	p.Lock()
	defer p.Unlock()

	recommendedFlags := C.FPDF_GetRecommendedV8Flags()
	log.Println(recommendedFlags)

	// @todo: implement me.
	return nil, pdfium_errors.ErrV8Unsupported
}

// FPDF_GetArrayBufferAllocatorSharedInstance initializes V8 isolates that
// will use PDFium's internal memory management.
//
// Use is optional, but allows external creation of isolates matching the
// ones PDFium will make when none is provided via
// |FPDF_LIBRARY_CONFIG::m_pIsolate|.
//
// Can only be called when the library is in an uninitialized or destroyed
// state.
//
// Only available when using a build that includes V8 and when using the
// build flag pdfium_v8.
func (p *PdfiumImplementation) FPDF_GetArrayBufferAllocatorSharedInstance(request *requests.FPDF_GetArrayBufferAllocatorSharedInstance) (*responses.FPDF_GetArrayBufferAllocatorSharedInstance, error) {
	p.Lock()
	defer p.Unlock()

	arrayBufferAllocator := C.FPDF_GetArrayBufferAllocatorSharedInstance()
	log.Println(arrayBufferAllocator)
	
	// @todo: implement me.
	return nil, pdfium_errors.ErrV8Unsupported
}
