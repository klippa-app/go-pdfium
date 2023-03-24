package requests

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/structs"
	"io"
)

type FPDF_LoadDocument struct {
	Path     *string // A path to a PDF file.
	Password *string // The password of the document.
}

type FPDF_LoadMemDocument struct {
	Data     *[]byte // A reference to the file data.
	Password *string // The password of the document.
}

type FPDF_LoadMemDocument64 struct {
	Data     *[]byte // A reference to the file data.
	Password *string // The password of the document.
}

type FPDF_LoadCustomDocument struct {
	Reader   io.ReadSeeker
	Size     int64
	Password *string // The password of the document.
}

type FPDF_GetLastError struct{}

type FPDF_SetSandBoxPolicyPolicy uint32

const (
	FPDF_SetSandBoxPolicyPolicyMachinetimeAccess FPDF_SetSandBoxPolicyPolicy = 1 // Policy for accessing the local machine time.
)

type FPDF_SetSandBoxPolicy struct {
	Policy FPDF_SetSandBoxPolicyPolicy
	Enable bool
}

type FPDF_CloseDocument struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_LoadPage struct {
	Document references.FPDF_DOCUMENT
	Index    int // The page number (0-index based).
}

type FPDF_ClosePage struct {
	Page references.FPDF_PAGE
}

type FPDF_GetFileVersion struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetDocPermissions struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetSecurityHandlerRevision struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetPageCount struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetPageWidth struct {
	Page Page
}

type FPDF_GetPageHeight struct {
	Page Page
}

type FPDF_GetPageSizeByIndex struct {
	Document references.FPDF_DOCUMENT
	Index    int
}

type FPDF_DocumentHasValidCrossReferenceTable struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetTrailerEnds struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetPageWidthF struct {
	Page Page
}

type FPDF_GetPageHeightF struct {
	Page Page
}

type FPDF_GetPageBoundingBox struct {
	Page Page
}

type FPDF_GetPageSizeByIndexF struct {
	Document references.FPDF_DOCUMENT
	Index    int
}

type FPDF_RenderPageBitmap struct {
	Bitmap references.FPDF_BITMAP
	Page   Page
	StartX int                      // Left pixel position of the display area in bitmap coordinates.
	StartY int                      // Top pixel position of the display area in bitmap coordinates.
	SizeX  int                      // Horizontal size (in pixels) for displaying the page.
	SizeY  int                      // Vertical size (in pixels) for displaying the page.
	Rotate enums.FPDF_PAGE_ROTATION // Page orientation.
	Flags  enums.FPDF_RENDER_FLAG   // 0 for normal display, or combination of enums.FPDF_RENDER_FLAG. With the enums.FPDF_RENDER_FLAG_ANNOT flag, it renders all annotations that do not require user-interaction, which are all annotations except widget and popup annotations.
}

type FPDF_RenderPageBitmapWithMatrix struct {
	Bitmap   references.FPDF_BITMAP
	Page     Page
	Matrix   structs.FPDF_FS_MATRIX // The transform matrix, which must be invertible. See PDF Reference 1.7, 4.2.2 Common Transformations.
	Clipping structs.FPDF_FS_RECTF  // The rect to clip to in device coords.
	Flags    enums.FPDF_RENDER_FLAG // 0 for normal display, or combination of enums.FPDF_RENDER_FLAG. With the enums.FPDF_RENDER_FLAG_ANNOT flag, it renders all annotations that do not require user-interaction, which are all annotations except widget and popup annotations.
}

type FPDF_DeviceToPage struct {
	Page    Page
	StartX  int                      // Left pixel position of the display area in device coordinates.
	StartY  int                      // Top pixel position of the display area in device coordinates.
	SizeX   int                      // Horizontal size (in pixels) for displaying the page.
	SizeY   int                      // Vertical size (in pixels) for displaying the page.
	Rotate  enums.FPDF_PAGE_ROTATION // Page orientation.
	DeviceX int                      // X value in device coordinates to be converted.
	DeviceY int                      // Y value in device coordinates to be converted.
}

type FPDF_PageToDevice struct {
	Page   Page
	StartX int                      // Left pixel position of the display area in device coordinates.
	StartY int                      // Top pixel position of the display area in device coordinates.
	SizeX  int                      // Horizontal size (in pixels) for displaying the page.
	SizeY  int                      // Vertical size (in pixels) for displaying the page.
	Rotate enums.FPDF_PAGE_ROTATION // Page orientation.
	PageX  float64                  // X value in page coordinates to be converted.
	PageY  float64                  // Y value in page coordinates to be converted.
}

type FPDFBitmap_Create struct {
	Width  int // The number of pixels in width for the bitmap. Must be greater than 0.
	Height int // The number of pixels in height for the bitmap. Must be greater than 0.
	Alpha  int // A flag indicating whether the alpha channel is used. Non-zero for using alpha, zero for not using.
}

type FPDFBitmap_CreateEx struct {
	Width   int // The number of pixels in width for the bitmap. Must be greater than 0.
	Height  int // The number of pixels in height for the bitmap. Must be greater than 0.
	Format  enums.FPDF_BITMAP_FORMAT
	Buffer  []byte      // DEPRECATED: use Pointer, unsupported on Webassembly runtime.
	Pointer interface{} // In the CGO runtime this must be an unsafe.Pointer to the first byte of a byte array, use unsafe.Pointer(&byteArray[0]) to get it. In the Webassembly runtime this must be uint64 with a pointer inside the Webassembly memory space created by malloc.
	Stride  int         // Number of bytes for each scan line. The value must be 0 or greater. When the value is 0, FPDFBitmap_CreateEx() will automatically calculate the appropriate value using Width and Format. When using an external buffer, it is recommended for the caller to pass in the value. When not using an external buffer, it is recommended for the caller to pass in 0.

}

type FPDFBitmap_GetFormat struct {
	Bitmap references.FPDF_BITMAP
}

type FPDFBitmap_FillRect struct {
	Bitmap references.FPDF_BITMAP
	Left   int
	Top    int
	Width  int
	Height int
	Color  uint64
}

type FPDFBitmap_GetBuffer struct {
	Bitmap references.FPDF_BITMAP
}

type FPDFBitmap_GetWidth struct {
	Bitmap references.FPDF_BITMAP
}

type FPDFBitmap_GetHeight struct {
	Bitmap references.FPDF_BITMAP
}

type FPDFBitmap_GetStride struct {
	Bitmap references.FPDF_BITMAP
}

type FPDFBitmap_Destroy struct {
	Bitmap references.FPDF_BITMAP
}

type FPDF_VIEWERREF_GetPrintScaling struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_VIEWERREF_GetNumCopies struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_VIEWERREF_GetPrintPageRange struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_VIEWERREF_GetPrintPageRangeCount struct {
	PageRange references.FPDF_PAGERANGE
}

type FPDF_VIEWERREF_GetPrintPageRangeElement struct {
	PageRange references.FPDF_PAGERANGE
	Index     uint64
}

type FPDF_VIEWERREF_GetDuplex struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_VIEWERREF_GetName struct {
	Document references.FPDF_DOCUMENT
	Key      string // Name of the key in the viewer pref dictionary.
}

type FPDF_CountNamedDests struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetNamedDestByName struct {
	Document references.FPDF_DOCUMENT
	Name     string
}

type FPDF_GetNamedDest struct {
	Document references.FPDF_DOCUMENT
	Index    int
}

type FPDF_GetXFAPacketCount struct {
	Document references.FPDF_DOCUMENT
}

type FPDF_GetXFAPacketName struct {
	Document references.FPDF_DOCUMENT
	Index    int
}

type FPDF_GetXFAPacketContent struct {
	Document references.FPDF_DOCUMENT
	Index    int
}

type FPDF_SetPrintMode struct {
	PrintMode enums.FPDF_PRINTMODE
}

type FPDF_RenderPage struct {
	DC     interface{} // Handle to the device context. This should be of type C.HDC, which is a device (screen, bitmap, or printer).
	Page   Page
	StartX int                      // Left pixel position of the display area in bitmap coordinates.
	StartY int                      // Top pixel position of the display area in bitmap coordinates.
	SizeX  int                      // Horizontal size (in pixels) for displaying the page.
	SizeY  int                      // Vertical size (in pixels) for displaying the page.
	Rotate enums.FPDF_PAGE_ROTATION // Page orientation.
	Flags  enums.FPDF_RENDER_FLAG   // 0 for normal display, or combination of enums.FPDF_RENDER_FLAG. With the enums.FPDF_RENDER_FLAG_ANNOT flag, it renders all annotations that do not require user-interaction, which are all annotations except widget and popup annotations.
}
