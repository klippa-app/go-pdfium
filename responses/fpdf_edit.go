package responses

type FPDFPage_GetRotation struct {
	Page         int          // The page number (0-index based).
	PageRotation PageRotation // The page rotation.
}

type FPDFPage_SetRotation struct{}

type FPDFPage_HasTransparency struct {
	Page            int  // The page number (0-index based).
	HasTransparency bool // Whether the page has transparency.
}
