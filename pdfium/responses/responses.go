package responses

import "image"

type GetPageCount struct {
	PageCount int
}

type RenderPage struct {
	Image             *image.RGBA
	PointToPixelRatio float64
}

type GetPageSize struct {
	Width  float64
	Height float64
}

type GetPageSizeInPixels struct {
	Width             int
	Height            int
	PointToPixelRatio float64
}

type GetPageText struct {
	Text string
}

type CharPosition struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

type GetPageTextStructuredChar struct {
	Text          string
	Angle         float64
	PointPosition CharPosition
	PixelPosition *CharPosition
}

type GetPageTextStructuredRect struct {
	Text          string
	PointPosition CharPosition
	PixelPosition *CharPosition
}

type GetPageTextStructured struct {
	Chars             []*GetPageTextStructuredChar
	Rects             []*GetPageTextStructuredRect
	PointToPixelRatio float64
}
