package structs

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
)

type FPDF_FS_RECTF struct {
	Left   float32
	Top    float32
	Right  float32
	Bottom float32
}

type FPDF_FS_QUADPOINTSF struct {
	X1 float32
	Y1 float32
	X2 float32
	Y2 float32
	X3 float32
	Y3 float32
	X4 float32
	Y4 float32
}

// FPDF_FS_MATRIX is a matrix that is composed as:
//   | A C E |
//   | B D F |
// and can be used to scale, rotate, shear and translate.
type FPDF_FS_MATRIX struct {
	A float32
	B float32
	C float32
	D float32
	E float32
	F float32
}

type FPDF_FS_SIZEF struct {
	Width  float32
	Height float32
}

type FPDF_COLORSCHEME struct {
	PathFillColor   uint64
	PathStrokeColor uint64
	TextFillColor   uint64
	TextStrokeColor uint64
}

type FPDF_COLOR struct {
	R uint
	G uint
	B uint
	A uint
}

type FPDF_IMAGEOBJ_METADATA struct {
	Width           uint
	Height          uint
	HorizontalDPI   float32
	VerticalDPI     float32
	BitsPerPixel    uint
	Colorspace      enums.FPDF_COLORSPACE
	MarkedContentID int
}

type FPDF_FS_POINTF struct {
	X float32
	Y float32
}

type FPDF_SYSTEMTIME struct {
	Year         uint16 // Years since 1900
	Month        uint16 // Months since January - [0,11]
	DayOfWeek    uint16 // Days since Sunday - [0,6]
	Day          uint16 // Day of the month - [1,31]
	Hour         uint16 // Hours since midnight - [0,23]
	Minute       uint16 // Minutes after the hour - [0,59]
	Second       uint16 // Seconds after the minute - [0,59]
	Milliseconds uint16 // Milliseconds after the second - [0,999]
}

// FPDF_FORMFILLINFO are the call
type FPDF_FORMFILLINFO struct {
	// Give the implementation a chance to release any resources after the
	// interface is no longer used.
	// Called by PDFium during the final cleanup process.
	Release func()

	// Invalidate the client area within the specified rectangle.
	// All positions are measured in PDF "user space".
	// Implementation should call FPDF_RenderPageBitmap() for repainting
	// the specified page area.
	//
	// Implementation required!
	FFI_Invalidate func(page references.FPDF_PAGE, left, top, right, bottom float64)

	// When the user selects text in form fields with the mouse, this
	// callback function will be invoked with the selected areas.
	// This callback function is useful for implementing special text
	// selection effects. An implementation should first record the
	// returned rectangles, then draw them one by one during the next
	// painting period. Lastly, it should remove all the recorded
	// rectangles when finished painting.
	FFI_OutputSelectedRect func(page references.FPDF_PAGE, left, top, right, bottom float64)

	// Set the Cursor shape.
	//
	// Implementation required!
	FFI_SetCursor func(cursorType enums.FXCT)

	// This method installs a system timer. An interval value is specified,
	// and every time that interval elapses, the system must call into the
	// callback function with the timer ID as returned by this function.
	//
	// Should return the timer identifier of the new timer if the function is successful.
	// An application passes this value to the FFI_KillTimer method to kill
	// the timer. Nonzero if it is successful; otherwise, it is zero.
	//
	// Implementation required!
	FFI_SetTimer func(elapse int, timerFunc func(idEvent int)) int

	// This method uninstalls a system timer, as set by an earlier call to
	// FFI_SetTimer.
	//
	// Implementation required!
	FFI_KillTimer func(timerID int)

	// This method receives the current local time on the system.
	// Note: Unused.
	//
	// Implementation required!
	FFI_GetLocalTime func() FPDF_SYSTEMTIME

	// This method will be invoked to notify the implementation when the
	// value of any FormField on the document had been changed.
	FFI_OnChange func()

	// This method receives the page handle associated with a specified
	// page index. Use FPDF_LoadPage to load the page.
	// The implementation is expected to keep track of the page handles it
	// receives from PDFium, and their mappings to page numbers. In some
	// cases, the document-level JavaScript action may refer to a page
	// which hadn't been loaded yet. To successfully run the Javascript
	// action, the implementation needs to load the page.
	//
	// Implementation required!
	FFI_GetPage func(document references.FPDF_DOCUMENT, index int) *references.FPDF_PAGE

	// This method receives the handle to the current page.
	// PDFium doesn't keep keep track of the "current page" (e.g. the one
	// that is most visible on screen), so it must ask the embedder for
	// this information.
	//
	// Implementation required when V8 support is present, otherwise unused.
	FFI_GetCurrentPage func(document references.FPDF_DOCUMENT) *references.FPDF_PAGE

	// This method receives currently rotation of the page view.
	//
	// Implementation required!
	FFI_GetRotation func(page references.FPDF_PAGE) enums.FPDF_PAGE_ROTATION

	// This method will execute a named action.
	// See ISO 32000-1:2008, section 12.6.4.11 for descriptions of the
	// standard named actions, but note that a document may supply any
	// name of its choosing.
	//
	// Implementation required!
	FFI_ExecuteNamedAction func(namedAction string)

	// Called when a text field is getting or losing focus.
	// Only supports text fields and combobox fields.
	FFI_SetTextFieldFocus func(value string, isFocus bool)

	// Ask the implementation to navigate to a uniform resource identifier.
	// If the embedder is version 2 or higher and have implementation for
	// FFI_DoURIActionWithKeyboardModifier, then
	// FFI_DoURIActionWithKeyboardModifier takes precedence over
	// FFI_DoURIAction.
	// See the URI actions description of <<PDF Reference, version 1.7>>
	// for more details.
	FFI_DoURIAction func(bsURI string)

	// This action changes the view to a specified destination.
	// See the Destinations description of <<PDF Reference, version 1.7>>
	// in 8.2.1 for more details.
	FFI_DoGoToAction func(pageIndex int, zoomMode enums.FPDF_ZOOM_MODE, pos []float32)

	// Note: the XFA methods are not implemmeted because we do not support XFA for now.
}
