//go:build pdfium_xfa
// +build pdfium_xfa

package implementation_cgo

// When XFA is enabled in the build, version MUST be 2.
func getFormFillVersion() int {
	return 2
}
