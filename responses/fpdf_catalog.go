package responses

type FPDFCatalog_IsTagged struct {
	IsTagged bool
}

type FPDFCatalog_SetLanguage struct{}

type FPDFCatalog_GetLanguage struct {
	Language string
}
