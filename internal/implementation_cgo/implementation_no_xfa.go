//go:build !pdfium_xfa
// +build !pdfium_xfa

package implementation_cgo

import "C"

// These values will never be true, but we need this do exist.
var (
	CFPDF_ERR_XFALOAD   = C.int(7)
	CFPDF_ERR_XFALAYOUT = C.int(8)
)
