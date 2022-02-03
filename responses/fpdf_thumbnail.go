package responses

import "github.com/klippa-app/go-pdfium/references"

type FPDFPage_GetDecodedThumbnailData struct {
	Thumbnail []byte // The thumbnail data, nil when it doesn't exist.
}

type FPDFPage_GetRawThumbnailData struct {
	RawThumbnail []byte // The raw thumbnail data, nil when it doesn't exist.
}

type FPDFPage_GetThumbnailAsBitmap struct {
	Bitmap *references.FPDF_BITMAP // The thumbnail as bitmap, nil if it doesn't exist.
}
