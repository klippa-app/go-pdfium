package structs

import "github.com/klippa-app/go-pdfium/enums"

type FPDF_FS_RECTF struct {
	Left   float32
	Top    float32
	Right  float32
	Bottom float32
}

type FPDF_FS_QUADPOINTSF struct {
	X1 float32
	Y1 float32
	X2 float32
	Y2 float32
	X3 float32
	Y3 float32
	X4 float32
	Y4 float32
}

// FPDF_FS_MATRIX is a matrix that is composed as:
//   | A C E |
//   | B D F |
// and can be used to scale, rotate, shear and translate.
type FPDF_FS_MATRIX struct {
	A float32
	B float32
	C float32
	D float32
	E float32
	F float32
}

type FPDF_FS_SIZEF struct {
	Width  float32
	Height float32
}

type FPDF_COLORSCHEME struct {
	PathFillColor   uint64
	PathStrokeColor uint64
	TextFillColor   uint64
	TextStrokeColor uint64
}

type FPDF_COLOR struct {
	R uint
	G uint
	B uint
	A uint
}

type FPDF_IMAGEOBJ_METADATA struct {
	Width           uint
	Height          uint
	HorizontalDPI   float32
	VerticalDPI     float32
	BitsPerPixel    uint
	Colorspace      enums.FPDF_COLORSPACE
	MarkedContentID int
}
