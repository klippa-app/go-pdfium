//go:build !pdfium_xfa
// +build !pdfium_xfa

package implementation_cgo

func getFormFillVersion() int {
	return 1
}
