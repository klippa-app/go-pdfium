package responses

import "github.com/klippa-app/go-pdfium/enums"

type FPDF_RenderPageBitmapWithColorScheme_Start struct {
	RenderStatus enums.FPDF_RENDER_STATUS
}

type FPDF_RenderPageBitmap_Start struct {
	RenderStatus enums.FPDF_RENDER_STATUS
}

type FPDF_RenderPage_Continue struct {
	RenderStatus enums.FPDF_RENDER_STATUS
}

type FPDF_RenderPage_Close struct{}
