package responses

import (
	"github.com/klippa-app/go-pdfium/references"
)

type OpenDocument struct {
	Document references.FPDF_DOCUMENT
}

type NewPage struct {
	Page references.FPDF_PAGE
}

type PageRotation int

const (
	PageRotationNone  PageRotation = 0 // 0: no rotation.
	PageRotation90CW  PageRotation = 1 // 1: rotate 90 degrees in clockwise direction.
	PageRotation180CW PageRotation = 2 // 2: rotate 180 degrees in clockwise direction.
	PageRotation270CW PageRotation = 3 // 3: rotate 270 degrees in clockwise direction.
)

type GetPageSize struct {
	Page   int     // The page this size came from (0-index based).
	Width  float64 // The width of the page in points. One point is 1/72 inch (around 0.3528 mm).
	Height float64 // The height of the page in points. One point is 1/72 inch (around 0.3528 mm).
}

type GetPageSizeInPixels struct {
	Page              int     // The page this size came from (0-index based).
	Width             int     // The width of the page in pixels.
	Height            int     // The height of the page in pixels.
	PointToPixelRatio float64 // The point to pixel ratio for the rendered image. How many points is 1 pixel in this image.
}
