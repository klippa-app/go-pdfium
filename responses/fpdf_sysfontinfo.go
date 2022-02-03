package responses

import "github.com/klippa-app/go-pdfium/enums"

type FPDF_GetDefaultTTFMap struct {
	TTFMap map[enums.FPDF_FXFONT_CHARSET]string
}

type FPDF_AddInstalledFont struct{}
type FPDF_SetSystemFontInfo struct{}
type FPDF_GetDefaultSystemFontInfo struct{}
type FPDF_FreeDefaultSystemFontInfo struct{}
