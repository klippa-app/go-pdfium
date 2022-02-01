package requests

type FPDF_CreateNewDocument struct{}

type FPDFPage_SetRotation struct {
	Page   Page
	Rotate PageRotation // New value of PDF page rotation.
}

type FPDFPage_GetRotation struct {
	Page Page
}

type FPDFPage_HasTransparency struct {
	Page Page
}
