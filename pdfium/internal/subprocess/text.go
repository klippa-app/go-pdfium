package subprocess

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_text.h"
import "C"

import (
	"log"
	"unsafe"
)

// GetPageText returns the text of a page
func (d *Document) GetPageText(page int) string {
	d.loadPage(page)

	log.Println("Loading text")

	mutex.Lock()
	textPage := C.FPDFText_LoadPage(d.page)
	charsInPage := int(C.FPDFText_CountChars(textPage))
	charData := make([]byte, charsInPage+20)
	charsWritten := C.FPDFText_GetText(textPage, C.int(0), C.int(charsInPage), (*C.ushort)(unsafe.Pointer(&charData[0])))
	C.FPDFText_ClosePage(textPage)
	mutex.Unlock()

	return string(charData[0:charsWritten])
}