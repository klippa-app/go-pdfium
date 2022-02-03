package implementation

/*
#cgo pkg-config: pdfium
#include "fpdf_sysfontinfo.h"
*/
import "C"
import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"unsafe"
)

// FPDF_GetDefaultTTFMap returns the the default character set to TT Font name map. The map is an array of FPDF_CharsetFontMap structs
func (p *PdfiumImplementation) FPDF_GetDefaultTTFMap(request *requests.FPDF_GetDefaultTTFMap) (*responses.FPDF_GetDefaultTTFMap, error) {
	p.Lock()
	defer p.Unlock()

	ttfMap := map[enums.FPDF_FXFONT_CHARSET]string{}

	var maps *C.FPDF_CharsetFontMap = C.FPDF_GetDefaultTTFMap()
	for {
		// Add font to the map.
		ttfMap[enums.FPDF_FXFONT_CHARSET(maps.charset)] = C.GoString(maps.fontname)

		// Go to next font map.
		maps = (*C.FPDF_CharsetFontMap)(unsafe.Pointer(uintptr(unsafe.Pointer(maps)) + unsafe.Sizeof(*maps)))

		// If last one, break out.
		if maps.charset == -1 {
			break
		}
	}

	return &responses.FPDF_GetDefaultTTFMap{
		TTFMap: ttfMap,
	}, nil
}

// FPDF_AddInstalledFont add a system font to the list in PDFium.
func (p *PdfiumImplementation) FPDF_AddInstalledFont(request *requests.FPDF_AddInstalledFont) (*responses.FPDF_AddInstalledFont, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDF_AddInstalledFont{}, nil
}

// FPDF_SetSystemFontInfo set the system font info interface into PDFium.
func (p *PdfiumImplementation) FPDF_SetSystemFontInfo(request *requests.FPDF_SetSystemFontInfo) (*responses.FPDF_SetSystemFontInfo, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDF_SetSystemFontInfo{}, nil
}

// FPDF_GetDefaultSystemFontInfo gets the default system font info interface for current platform.
func (p *PdfiumImplementation) FPDF_GetDefaultSystemFontInfo(request *requests.FPDF_GetDefaultSystemFontInfo) (*responses.FPDF_GetDefaultSystemFontInfo, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDF_GetDefaultSystemFontInfo{}, nil
}

// FPDF_FreeDefaultSystemFontInfo frees a default system font info interface.
func (p *PdfiumImplementation) FPDF_FreeDefaultSystemFontInfo(request *requests.FPDF_FreeDefaultSystemFontInfo) (*responses.FPDF_FreeDefaultSystemFontInfo, error) {
	p.Lock()
	defer p.Unlock()

	return &responses.FPDF_FreeDefaultSystemFontInfo{}, nil
}
