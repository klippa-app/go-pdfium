package responses

import (
	"image"

	"github.com/klippa-app/go-pdfium/document"
)

type OpenDocument struct {
	Document document.Ref
}

type GetFileVersion struct {
	FileVersion int // The numeric version of the file: 14 for 1.4, 15 for 1.5, ...
}

type GetDocPermissions struct {
	DocPermissions                      uint32 // A 32-bit integer which indicates the permission flags. Please refer to "TABLE 3.20 User access permissions" in PDF Reference 1.7 P123 for detailed description. If the document is not protected, 0xffffffff (4294967295) will be returned.
	PrintDocument                       bool   // Bit position 3: (Security handlers of revision 2) Print the document, (Security handlers of revision 3 or greater) Print the document (possibly not at the highest quality level, depending on whether PrintDocumentAsFaithfulDigitalCopy (bit 12) is also set).
	ModifyContents                      bool   // Bit position 4: Modify the contents of the document by operations other than those controlled by AddOrModifyTextAnnotations (bit 6), FillInExistingInteractiveFormFields (bit 9), and AssembleDocument (bit 11).
	CopyOrExtractText                   bool   // Bit position 5: (Security handlers of revision 2) Copy or otherwise extract  text and graphics from the document, including extracting text and graphics (in support of accessibility to users with disabilities or for other purposes). (Security handlers of revision 3 or greater) Copy or otherwise extract text and graphics from the document by operations other than that controlled by ExtractTextAndGraphics (bit 10).
	AddOrModifyTextAnnotations          bool   // Bit position 6: Add or modify text annotations
	FillInInteractiveFormFields         bool   // Bit position 6: fill in interactive form fields
	CreateOrModifyInteractiveFormFields bool   // Bit position 6 & 4: create or modify interactive form fields (including signature fields).
	FillInExistingInteractiveFormFields bool   // Bit position 9: (Security handlers of revision 3 or greater) Fill in existing interactive form fields (including signature fields), even if FillInInteractiveFormFields (bit 6) is clear.
	ExtractTextAndGraphics              bool   // Bit position 10: (Security handlers of revision 3 or greater) Extract text and graphics (in support of accessibility to users with disabilities or for other purposes).
	AssembleDocument                    bool   // Bit position 11: (Security handlers of revision 3 or greater) Assemble the  document (insert, rotate, or delete pages and create bookmarks or thumbnail images), even if ModifyContents (bit 4) is clear.
	PrintDocumentAsFaithfulDigitalCopy  bool   // Bit position 12: (Security handlers of revision 3 or greater) Print the document to a representation from which a faithful digital copy of the PDF content could be generated. When this bit is clear (and PrintDocument (bit 3) is set), printing is limited to a low-level representation of the appearance, possibly of degraded quality.
}

type GetSecurityHandlerRevision struct {
	SecurityHandlerRevision int // The revision number of security handler. Please refer to key "R" in "TABLE 3.19 Additional encryption dictionary entries for the standard security handler" in PDF Reference 1.7 P122 for detailed description. If the document is not protected, -1 will be returned.
}

type GetPageCount struct {
	PageCount int // The amount of pages of the document.
}

type PageMode int

const (
	PageModeUnknown        PageMode = -1 // Page mode: unknown.
	PageModeUseNone        PageMode = 0  // Page mode: use none, which means neither document outline nor thumbnail images visible.
	PageModeUseOutlines    PageMode = 1  // Page mode: document outline visible.
	PageModeUseThumbs      PageMode = 2  // Page mode: thumbnail images visible.
	PageModeFullScreen     PageMode = 3  // Page mode: full screen - with no menu bar, no windows controls and no any other windows visible.
	PageModeUseOC          PageMode = 4  // Page mode: optional content group panel visible.
	PageModeUseAttachments PageMode = 5  // Page mode: attachments panel visible.
)

type GetPageMode struct {
	PageMode PageMode // The document's page mode, which describes how the document should be displayed when opened.
}

type GetMetadata struct {
	Tag   string // The requested metadata tag.
	Value string // The value of the tag if found, string is empty if the value is not found.
}

type PageRotation int

const (
	PageRotationNone  PageRotation = 0 // 0: no rotation.
	PageRotation90CW  PageRotation = 1 // 1: rotate 90 degrees in clockwise direction.
	PageRotation180CW PageRotation = 2 // 2: rotate 180 degrees in clockwise direction.
	PageRotation270CW PageRotation = 3 // 3: rotate 270 degrees in clockwise direction.
)

type GetPageRotation struct {
	Page         int          // The page number (0-index based).
	PageRotation PageRotation // The page rotation.
}

type GetPageTransparency struct {
	Page            int  // The page number (0-index based).
	HasTransparency bool // Whether the page has transparency.
}

type FlattenPageResult int

const (
	FlattenPageResultFail        FlattenPageResult = 0 // Flatten operation failed.
	FlattenPageResultSuccess     FlattenPageResult = 1 // Flatten operation succeed.
	FlattenPageResultNothingToDo FlattenPageResult = 2 // There is nothing can be flatten.
)

type FlattenPage struct {
	Page   int               // The page number (0-index based).
	Result FlattenPageResult // The result of the flatten.
}

type RenderPage struct {
	Page              int         // The rendered page number (0-index based).
	PointToPixelRatio float64     // The point to pixel ratio for the rendered image. How many points is 1 pixel in this image.
	Image             *image.RGBA // The rendered image.
	Width             int         // The width of the rendered image.
	Height            int         // The height of the rendered image.
}

type RenderPagesPage struct {
	Page              int     // The rendered page number (0-index based).
	PointToPixelRatio float64 // The point to pixel ratio for the rendered image. How many points is 1 pixel for this page in this image.
	Width             int     // The width of the rendered page inside the image.
	Height            int     // The height of the rendered page inside the image.
	X                 int     // The X start position of this page inside the image.
	Y                 int     // The Y start position of this page inside the image.
}

type RenderPages struct {
	Pages  []RenderPagesPage // Information about the rendered pages inside this image.
	Image  *image.RGBA       // The rendered image.
	Width  int               // The width of the rendered image.
	Height int               // The height of the rendered image.
}

type RenderToFile struct {
	Pages             []RenderPagesPage // Information about the rendered pages inside this image.
	ImageBytes        *[]byte           // The byte array of the rendered file when OutputTarget is RenderToFileOutputTargetBytes.
	ImagePath         string            // The file path when OutputTarget is RenderToFileOutputTargetFile, is a tmp path when TargetFilePath was empty in the request.
	Width             int               // The width of the rendered image.
	Height            int               // The height of the rendered image.
	PointToPixelRatio float64           // The point to pixel ratio for the rendered image. How many points is 1 pixel in this image. Only set when rendering one page.
}

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

type GetPageText struct {
	Page int    // The page this text came from (0-index based).
	Text string // The plain text of a page.
}

type CharPosition struct {
	Left   float64 // The position of this char from the left.
	Top    float64 // The position of this char from the top.
	Right  float64 // The position of this char from the right.
	Bottom float64 // The position of this char from the bottom.
}

type FontInformation struct {
	Size         float64 // Font size in points (also known as em).
	SizeInPixels *int    // Font size in pixels, only available when PixelPositions is used.
	Weight       int     // The weight of the font, can be negative for spaces and newlines.
	Name         string  // The name of the font, can be empty for spaces and newlines.
	Flags        int     // Font flags, should be interpreted per PDF spec 1.7, Section 5.7.1 Font Descriptor Flags.
}

type GetPageTextStructuredChar struct {
	Text            string           // The text of this char.
	Angle           float64          // The angle this char is in.
	PointPosition   CharPosition     // The position of this char in points.
	PixelPosition   *CharPosition    // The position of this char in pixels. When PixelPositions are requested.
	FontInformation *FontInformation // The font information of this char. When CollectFontInformation is enabled.
}

type GetPageTextStructuredRect struct {
	Text            string           // The text of this rect.
	PointPosition   CharPosition     // The position of this rect in points.
	PixelPosition   *CharPosition    // The position of this rect in pixels. When PixelPositions are requested.
	FontInformation *FontInformation // The font information of this rect. When CollectFontInformation is enabled.
}

type GetPageTextStructured struct {
	Page              int                          // The page structured this text came from (0-index based).
	Chars             []*GetPageTextStructuredChar // A list of chars in a page. When Mode is GetPageTextStructuredModeChars or GetPageTextStructuredModeBoth.
	Rects             []*GetPageTextStructuredRect // A list of rects in a page. When Mode is GetPageTextStructuredModeRects or GetPageTextStructuredModeBoth.
	PointToPixelRatio float64                      // The point to pixel ratio for the calculated positions.
}
