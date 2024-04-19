//go:build pdfium_xfa
// +build pdfium_xfa

package implementation_cgo

// #cgo pkg-config: pdfium
// #cgo CFLAGS: -DPDF_ENABLE_XFA
// #include "fpdfview.h"
import "C"

var CFPDF_ERR_XFALOAD = C.ulong(C.FPDF_ERR_XFALOAD)
var CFPDF_ERR_XFALAYOUT = C.ulong(C.FPDF_ERR_XFALAYOUT)
