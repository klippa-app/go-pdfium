package responses

type FPDFDoc_GetPageModeMode int

const (
	FPDFDoc_GetPageModeModeUnknown        FPDFDoc_GetPageModeMode = -1 // Page mode: unknown.
	FPDFDoc_GetPageModeModeUseNone        FPDFDoc_GetPageModeMode = 0  // Page mode: use none, which means neither document outline nor thumbnail images visible.
	FPDFDoc_GetPageModeModeUseOutlines    FPDFDoc_GetPageModeMode = 1  // Page mode: document outline visible.
	FPDFDoc_GetPageModeModeUseThumbs      FPDFDoc_GetPageModeMode = 2  // Page mode: thumbnail images visible.
	FPDFDoc_GetPageModeModeFullScreen     FPDFDoc_GetPageModeMode = 3  // Page mode: full screen - with no menu bar, no windows controls and no any other windows visible.
	FPDFDoc_GetPageModeModeUseOC          FPDFDoc_GetPageModeMode = 4  // Page mode: optional content group panel visible.
	FPDFDoc_GetPageModeModeUseAttachments FPDFDoc_GetPageModeMode = 5  // Page mode: attachments panel visible.
)

type FPDFDoc_GetPageMode struct {
	PageMode FPDFDoc_GetPageModeMode // The document's page mode, which describes how the document should be displayed when opened.
}

type FSDK_SetUnSpObjProcessHandler struct {
}

type FSDK_SetTimeFunction struct {
}

type FSDK_SetLocaltimeFunction struct {
}
