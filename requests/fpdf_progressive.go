package requests

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/structs"
)

type FPDF_RenderPageBitmapWithColorScheme_Start struct {
	Bitmap                 references.FPDF_BITMAP
	Page                   Page
	StartX                 int
	StartY                 int
	SizeX                  int
	SizeY                  int
	Rotate                 enums.FPDF_PAGE_ROTATION
	Flags                  enums.FPDF_RENDER_FLAG
	ColorScheme            *structs.FPDF_COLORSCHEME
	NeedToPauseNowCallback func() bool
}

type FPDF_RenderPageBitmap_Start struct {
	Bitmap                 references.FPDF_BITMAP
	Page                   Page
	StartX                 int
	StartY                 int
	SizeX                  int
	SizeY                  int
	Rotate                 enums.FPDF_PAGE_ROTATION
	Flags                  enums.FPDF_RENDER_FLAG
	NeedToPauseNowCallback func() bool // A callback mechanism allowing the page rendering process to pause.
}

type FPDF_RenderPage_Continue struct {
	Page                   Page
	NeedToPauseNowCallback func() bool // A callback mechanism allowing the page rendering process to pause.This can be nil if you don't want to pause.
}

type FPDF_RenderPage_Close struct {
	Page Page
}
