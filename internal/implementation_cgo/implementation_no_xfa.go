//go:build !pdfium_xfa
// +build !pdfium_xfa

package implementation_cgo

import "C"

// These values will never be used, but we need this do exist.
var (
	CFPDF_ERR_XFALOAD   = C.ulong(7)
	CFPDF_ERR_XFALAYOUT = C.ulong(8)
)
