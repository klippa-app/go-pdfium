package structs

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
