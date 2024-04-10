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
//
//	| A C E |
//	| B D F |
//
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

// FPDF_FORMFILLINFO is the callback interface for form filling.
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

	// A IPDF_JSPLATFORM instance.
	// Unused if PDFium is built without V8 support. Otherwise, if nil, then
	// JavaScript will be prevented from executing while rendering the document.
	JsPlatform *IPDF_JSPLATFORM

	// Whether the XFA module is disabled when built with the XFA module.
	// Ignored on non-XFA builds.
	XFA_disabled bool

	// FFI_DisplayCaret shows the caret at specified position.
	// Parameters:
	//       page            -   Handle to page. Returned by FPDF_LoadPage().
	//       left            -   Left position of the client area in PDF page
	//                           coordinates.
	//       top             -   Top position of the client area in PDF page
	//                           coordinates.
	//       right           -   Right position of the client area in PDF page
	//                           coordinates.
	//       bottom          -   Bottom position of the client area in PDF page
	//                           coordinates.
	//
	// Required for XFA, otherwise set to nil.
	// Ignored on non-XFA builds.
	FFI_DisplayCaret func(page references.FPDF_PAGE, bVisible bool, left, top, right, bottom float64)

	// FFI_GetCurrentPageIndex returns the current page index.
	// Parameters:
	//       document        -   Handle to document from FPDF_LoadDocument().
	//
	// Required for XFA, otherwise set to nil.
	// Ignored on non-XFA builds.
	FFI_GetCurrentPageIndex func(document references.FPDF_DOCUMENT) int

	// FFI_SetCurrentPage sets the current page.
	// Parameters:
	//       document        -   Handle to document from FPDF_LoadDocument().
	//       iCurPage        -   The index of the PDF page.
	//
	// Required for XFA, otherwise set to nil.
	// Ignored on non-XFA builds.
	FFI_SetCurrentPage func(document references.FPDF_DOCUMENT, iCurPage int)

	// FFI_GotoURL will navigate to the specified URL.
	// Parameters:
	//       document        -   Handle to document from FPDF_LoadDocument().
	//       url             -   The URL to navigate to
	//
	// Required for XFA, otherwise set to nil.
	// Ignored on non-XFA builds.
	FFI_GotoURL func(document references.FPDF_DOCUMENT, url string)

	// FFI_GetPageViewRect will get the current page view rectangle.
	// Parameters:
	//       page            -   Handle to page. Returned by FPDF_LoadPage().
	//       left            -   The pointer to receive left position of the page
	//                           view area in PDF page coordinates.
	//       top             -   The pointer to receive top position of the page
	//                           view area in PDF page coordinates.
	//       right           -   The pointer to receive right position of the
	//                           page view area in PDF page coordinates.
	//       bottom          -   The pointer to receive bottom position of the
	//                           page view area in PDF page coordinates.
	//
	// Required for XFA, otherwise set to nil.
	// Ignored on non-XFA builds.
	FFI_GetPageViewRect func(document references.FPDF_DOCUMENT, url string) (left, top, right, bottom float64)

	// FFI_PageEvent fires when pages have been added to or deleted from
	// the XFA document.
	// Parameters:
	//       page_count      -   The number of pages to be added or deleted.
	//       event_type      -   A enums.FXFA_PAGEVIEWEVENT value.
	//
	// Required for XFA, otherwise set to nil.
	// Ignored on non-XFA builds.
	FFI_PageEvent func(page_count int, event_type enums.FXFA_PAGEVIEWEVENT)

	// FFI_PopupMenu will track the right context menu for XFA fields.
	// Parameters:
	//       page            -   Handle to page. Returned by FPDF_LoadPage().
	//       menuFlag        -   The menu flags. Please refer to macro definition
	//                           of enums.FXFA_MENU and this can be one or a
	//                           combination of these macros.
	//       x               -   X position of the client area in PDF page
	//                           coordinates.
	//       y               -   Y position of the client area in PDF page
	//                           coordinates.
	//
	// Required for XFA, otherwise set to nil.
	// Ignored on non-XFA builds.
	FFI_PopupMenu func(page references.FPDF_PAGE, menuFlag int, x, y float32) bool

	// FFI_OpenFile will open the specified file with the specified mode.
	// Parameters:
	//       fileFlag        -   The file flag. Please refer to macro definition
	//                           of enums.FXFA_SAVEAS and use one of these macros.
	//       url             -   The file url to open.
	//       mode            -   The mode for open file, e.g. "rb" or "wb".
	//
	// Required for XFA, otherwise set to nil.
	// Ignored on non-XFA builds.
	FFI_OpenFile func(fileFlag enums.FXFA_SAVEAS, url, mode string) *FPDF_FILEHANDLER

	// FFI_EmailTo will email the specified file stream to the specified
	// contact.
	// Parameters:
	//       fileHandler    -   Handle to the FPDF_FILEHANDLER.
	//       to             -   A semicolon-delimited list of recipients.
	//       subject        -   The subject of the message.
	//       cc             -   A semicolon-delimited list of CC recipients.
	//       bcc            -   A semicolon-delimited list of BCC recipients.
	//       msg            -   The message to be sent.
	//
	// Required for XFA, otherwise set to nil.
	// Ignored on non-XFA builds.
	FFI_EmailTo func(fileHandler *FPDF_FILEHANDLER, to, subject, cc, bcc, msg string)

	// FFI_UploadTo will upload the specified file stream to the
	// specified URL.
	// Parameters:
	//       fileHandler    -   Handle to the FPDF_FILEHANDLER.
	//       fileFlag       -   The file flag. Please refer to macro definition
	//                          of enums.FXFA_SAVEAS and use one of these macros.
	//       uploadTo       -   The URL to upload to
	//
	// Required for XFA, otherwise set to nil.
	// Ignored if Version is lower than 2.
	FFI_UploadTo func(fileHandler *FPDF_FILEHANDLER, fileFlag enums.FXFA_SAVEAS, uploadTo string)

	// FFI_GetPlatform will return the current platform.
	//
	// Required for XFA, otherwise set to nil.
	// Ignored on non-XFA builds.
	FFI_GetPlatform func() string

	// FFI_GetLanguage will return the current language.
	//
	// Required for XFA, otherwise set to nil.
	// Ignored on non-XFA builds.
	FFI_GetLanguage func() string

	// FFI_DownloadFromURL will download the specified file from the URL.
	// Parameters:
	//       url             -   The file url to download.
	//
	// Required for XFA, otherwise set to nil.
	// Ignored on non-XFA builds.
	FFI_DownloadFromURL func(fileFlag enums.FXFA_SAVEAS, url, mode string) *FPDF_FILEHANDLER

	// FFI_PostRequestURL will post the request to the server URL.
	// Parameters:
	//       url             -   The server URL to post to.
	//       data            -   The post data.
	//       contentType     -   The content type of the request data.
	//       encode          -   The encode type.
	//       header          -   The request header.
	//       response        -   A FPDF_BSTR to write the response to.
	//
	// Required for XFA, otherwise set to nil.
	// Ignored on non-XFA builds.
	FFI_PostRequestURL func(url, data, contentType, encode, header string, response references.FPDF_BSTR) bool

	// FFI_PostRequestURL will put the request to the server URL.
	// Parameters:
	//       url             -   The server URL to post to.
	//       data            -   The post data.
	//       encode          -   The encode type.
	//
	// Required for XFA, otherwise set to nil.
	// Ignored on non-XFA builds.
	FFI_PutRequestURL func(url, data, encode string) bool

	// FFI_PostRequestURL is called when the focused annotation is updated.
	// Parameters:
	//       annot           -   The focused annotation.
	//       page_index      -   Index number of the page which contains the
	//                           focused annotation. 0 for the first page.
	//
	// Ignored on non-XFA builds.
	FFI_OnFocusChange func(annot references.FPDF_ANNOTATION, page_index int)

	// FFI_DoURIActionWithKeyboardModifier asks the implementation to navigate
	// to a uniform resource identifier with the specified modifiers.
	// Parameters:
	//       uri             -   The uniform resource identifier to navigate to.
	//       modifiers       -   Keyboard modifier that indicates which of
	//                           the virtual keys are down, if any.
	//
	// Ignored on non-XFA builds.
	FFI_DoURIActionWithKeyboardModifier func(uri string, modifiers int)
}

// IPDF_JSPLATFORM is the callback interface for XFA JS handling.
type IPDF_JSPLATFORM struct {
	// Pop up a dialog to show warning or hint.
	// Parameters:
	//       msg            -   A string containing the message to be displayed.
	//       title          -   The title of the dialog.
	//       nButton        -   The type of button group, one of the
	//                          enums.JSPLATFORM_ALERT_BUTTON values.
	//       nIcon          -   The type of the icon, one of the
	//                      -   enums.JSPLATFORM_ALERT_ICON values.
	//
	// Implementation required!
	App_alert func(msg, title string, nButton enums.JSPLATFORM_ALERT_BUTTON, nIcon enums.JSPLATFORM_ALERT_ICON) int

	// App_beep causes the system to play a sound.
	// Parameters:
	//       nType       -   The sound type, see enums.JSPLATFORM_BEEP.
	//
	// Implementation required!
	App_beep func(nType enums.JSPLATFORM_BEEP)

	// App_response displays a dialog box containing a question and an entry
	// field for the user to reply to the question.
	// Parameters:
	//       question        -   The question to be posed to the user.
	//       title           -   The title of the dialog box.
	//       defaultValue    -   A default value for the answer to the question. If
	//                           not specified, no default value is presented.
	//       cLabel          -   A short string to appear in front of and on the
	//                           same line as the edit text field.
	//       bPassword       -   If true, indicates that the user's response should
	//                           be shown as asterisks (*) or bullets (?) to mask
	//                           the response, which might be sensitive information.
	//
	// Implementation required!
	App_response func(question, title, defaultValue, cLabel string, bPassword bool) string

	// Doc_getFilePath returns the file path of the current document.
	//
	// Implementation required!
	Doc_getFilePath func() string

	// Doc_mail mails the data buffer as an attachment to all recipients, with
	// or without user interaction.
	// Parameters:
	//       mailData    -   The data to send as attachment.
	//       bUI         -   If true, the rest of the parameters are used in a
	//                       compose-new-message window that is displayed to the
	//                       user. If false, the cTo parameter is required and
	//                       all others are optional.
	//       to          -   A semicolon-delimited list of recipients for the
	//                        message.
	//       subject     -   The subject of the message. The length limit is
	//                        64 KB.
	//       c          -   A semicolon-delimited list of CC recipients for
	//                        the message.
	//       bcc         -   A semicolon-delimited list of BCC recipients for
	//                        the message.
	//       msg         -   The content of the message. The length limit is
	//                        64 KB.
	//
	// Implementation required!
	Doc_mail func(mailData []byte, bUI bool, to, subject, cc, bcc, msg string)

	// Doc_print prints all or a specific number of pages of the document.
	// Parameters:
	//       bUI           -   If true, will cause a UI to be presented to the
	//                         user to obtain printing information and confirm
	//                         the action.
	//       nStart        -   A 0-based index that defines the start of an
	//                         inclusive range of pages.
	//       nEnd          -   A 0-based index that defines the end of an
	//                         inclusive page range.
	//       bSilent       -   If true, suppresses the cancel dialog box while
	//                         the document is printing. The default is false.
	//       bShrinkToFit  -   If true, the page is shrunk (if necessary) to
	//                         fit within the imageable area of the printed page.
	//       bPrintAsImage -   If true, print pages as an image.
	//       bReverse      -   If true, print from nEnd to nStart.
	//       bAnnotations  -   If true (the default), annotations are
	//                         printed.
	//
	// Implementation required!
	Doc_print func(bUI bool, nStart, nEnd int, bSilent, bShrinkToFit, bPrintAsImage, bReverse, bAnnotations bool)

	// Doc_submitForm sends the form data to a specified URL.
	//
	// Implementation required!
	Doc_submitForm func(formData []byte, url string)

	// Doc_gotoPage jumps to a specified page.
	//
	// Implementation required!
	Doc_gotoPage func(nPageNum int)

	// Field_browse shows a file selection dialog, and returns the selected
	// file path.
	//
	// Implementation required!
	Field_browse func() string
}

// FPDF_FILEHANDLER is a Structure for file reading or writing (I/O).
//
// Note: This is a handler and should be implemented by callers,
// and is only used from XFA.
type FPDF_FILEHANDLER struct {
	// @todo: implement me.
	//   /*
	//   * User-defined data.
	//   * Note: Callers can use this field to track controls.
	//   */
	//  void* clientData;
	//
	//  /*
	//   * Callback function to release the current file stream object.
	//   *
	//   * Parameters:
	//   *       clientData   -  Pointer to user-defined data.
	//   * Returns:
	//   *       None.
	//   */
	//  void (*Release)(void* clientData);
	//
	//  /*
	//   * Callback function to retrieve the current file stream size.
	//   *
	//   * Parameters:
	//   *       clientData   -  Pointer to user-defined data.
	//   * Returns:
	//   *       Size of file stream.
	//   */
	//  FPDF_DWORD (*GetSize)(void* clientData);
	//
	//  /*
	//   * Callback function to read data from the current file stream.
	//   *
	//   * Parameters:
	//   *       clientData   -  Pointer to user-defined data.
	//   *       offset       -  Offset position starts from the beginning of file
	//   *                       stream. This parameter indicates reading position.
	//   *       buffer       -  Memory buffer to store data which are read from
	//   *                       file stream. This parameter should not be nil.
	//   *       size         -  Size of data which should be read from file stream,
	//   *                       in bytes. The buffer indicated by |buffer| must be
	//   *                       large enough to store specified data.
	//   * Returns:
	//   *       0 for success, other value for failure.
	//   */
	//  FPDF_RESULT (*ReadBlock)(void* clientData,
	//                           FPDF_DWORD offset,
	//                           void* buffer,
	//                           FPDF_DWORD size);
	//
	//  /*
	//   * Callback function to write data into the current file stream.
	//   *
	//   * Parameters:
	//   *       clientData   -  Pointer to user-defined data.
	//   *       offset       -  Offset position starts from the beginning of file
	//   *                       stream. This parameter indicates writing position.
	//   *       buffer       -  Memory buffer contains data which is written into
	//   *                       file stream. This parameter should not be nil.
	//   *       size         -  Size of data which should be written into file
	//   *                       stream, in bytes.
	//   * Returns:
	//   *       0 for success, other value for failure.
	//   */
	//  FPDF_RESULT (*WriteBlock)(void* clientData,
	//                            FPDF_DWORD offset,
	//                            const void* buffer,
	//                            FPDF_DWORD size);
	//  /*
	//   * Callback function to flush all internal accessing buffers.
	//   *
	//   * Parameters:
	//   *       clientData   -  Pointer to user-defined data.
	//   * Returns:
	//   *       0 for success, other value for failure.
	//   */
	//  FPDF_RESULT (*Flush)(void* clientData);
	//
	//  /*
	//   * Callback function to change file size.
	//   *
	//   * Description:
	//   *       This function is called under writing mode usually. Implementer
	//   *       can determine whether to realize it based on application requests.
	//   * Parameters:
	//   *       clientData   -  Pointer to user-defined data.
	//   *       size         -  New size of file stream, in bytes.
	//   * Returns:
	//   *       0 for success, other value for failure.
	//   */
	//  FPDF_RESULT (*Truncate)(void* clientData, FPDF_DWORD size);
}
