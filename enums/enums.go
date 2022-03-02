package enums

// The file identifier entry type. See section 14.4 "File Identifiers" of the
// ISO 32000-1:2008 spec.
type FPDF_FILEIDTYPE int

const (
	FPDF_FILEIDTYPE_PERMANENT FPDF_FILEIDTYPE = 0
	FPDF_FILEIDTYPE_CHANGING  FPDF_FILEIDTYPE = 1
)

type FPDF_ACTION_ACTION uint32

const (
	FPDF_ACTION_ACTION_UNSUPPORTED  FPDF_ACTION_ACTION = 0 // Action type: unsupported action type.
	FPDF_ACTION_ACTION_GOTO         FPDF_ACTION_ACTION = 1 // This action contains information which can be used to go to a destination within current document.
	FPDF_ACTION_ACTION_REMOTEGOTO   FPDF_ACTION_ACTION = 2 // This action contains information which can be used to launch an application or opens or prints a document.
	FPDF_ACTION_ACTION_URI          FPDF_ACTION_ACTION = 3 // This action contains information which can be used to go to a destination within another document.
	FPDF_ACTION_ACTION_LAUNCH       FPDF_ACTION_ACTION = 4 // This action contains information which identifies (resolves to) a resource on the Internet - such as web pages, a file that is the destination of a hypertext link, and etc.
	FPDF_ACTION_ACTION_EMBEDDEDGOTO FPDF_ACTION_ACTION = 5 // This action contains information which can be used to Go to a destination in an embedded file.
)

type FPDF_PAGE_ROTATION int

const (
	FPDF_PAGE_ROTATION_NONE   FPDF_PAGE_ROTATION = 0 // 0: no rotation.
	FPDF_PAGE_ROTATION_90_CW  FPDF_PAGE_ROTATION = 1 // 1: rotate 90 degrees in clockwise direction.
	FPDF_PAGE_ROTATION_180_CW FPDF_PAGE_ROTATION = 2 // 2: rotate 180 degrees in clockwise direction.
	FPDF_PAGE_ROTATION_270_CW FPDF_PAGE_ROTATION = 3 // 3: rotate 270 degrees in clockwise direction.
)

// View destination fit types. See pdfmark reference v9, page 48.
type FPDF_PDFDEST_VIEW uint32

const (
	FPDF_PDFDEST_VIEW_UNKNOWN_MODE FPDF_PDFDEST_VIEW = 0
	FPDF_PDFDEST_VIEW_XYZ          FPDF_PDFDEST_VIEW = 1
	FPDF_PDFDEST_VIEW_FIT          FPDF_PDFDEST_VIEW = 2
	FPDF_PDFDEST_VIEW_FITH         FPDF_PDFDEST_VIEW = 3
	FPDF_PDFDEST_VIEW_FITV         FPDF_PDFDEST_VIEW = 4
	FPDF_PDFDEST_VIEW_FITR         FPDF_PDFDEST_VIEW = 5
	FPDF_PDFDEST_VIEW_FITB         FPDF_PDFDEST_VIEW = 6
	FPDF_PDFDEST_VIEW_FITBH        FPDF_PDFDEST_VIEW = 7
	FPDF_PDFDEST_VIEW_FITBV        FPDF_PDFDEST_VIEW = 8
)

// Additional-action types of page object
type FPDF_PAGE_AACTION int

const (
	FPDF_PAGE_AACTION_OPEN  FPDF_PAGE_AACTION = 0 // OPEN (/O) -- An action to be performed when the page is opened
	FPDF_PAGE_AACTION_CLOSE FPDF_PAGE_AACTION = 1 // CLOSE (/C) -- An action to be performed when the page is closed
)

// Additional actions type of document
type FPDF_DOC_AACTION int

const (
	FPDF_DOC_AACTION_WC FPDF_DOC_AACTION = 0x10 // WC, before closing document, JavaScript action.
	FPDF_DOC_AACTION_WS FPDF_DOC_AACTION = 0x11 // WS before saving document, JavaScript action.
	FPDF_DOC_AACTION_DS FPDF_DOC_AACTION = 0x12 // DS, after saving document, JavaScript action.
	FPDF_DOC_AACTION_WP FPDF_DOC_AACTION = 0x13 // WP, before printing document, JavaScript action.
	FPDF_DOC_AACTION_DP FPDF_DOC_AACTION = 0x14 // DP, after printing document, JavaScript action.
)

type FPDF_UNSP int

const (
	FPDF_UNSP_DOC_XFAFORM               FPDF_UNSP = 1
	FPDF_UNSP_DOC_PORTABLECOLLECTION    FPDF_UNSP = 2
	FPDF_UNSP_DOC_ATTACHMENT            FPDF_UNSP = 3
	FPDF_UNSP_DOC_SECURITY              FPDF_UNSP = 4
	FPDF_UNSP_DOC_SHAREDREVIEW          FPDF_UNSP = 5
	FPDF_UNSP_DOC_SHAREDFORM_ACROBAT    FPDF_UNSP = 6
	FPDF_UNSP_DOC_SHAREDFORM_FILESYSTEM FPDF_UNSP = 7
	FPDF_UNSP_DOC_SHAREDFORM_EMAIL      FPDF_UNSP = 8
	FPDF_UNSP_ANNOT_3DANNOT             FPDF_UNSP = 11
	FPDF_UNSP_ANNOT_MOVIE               FPDF_UNSP = 12
	FPDF_UNSP_ANNOT_SOUND               FPDF_UNSP = 13
	FPDF_UNSP_ANNOT_SCREEN_MEDIA        FPDF_UNSP = 14
	FPDF_UNSP_ANNOT_SCREEN_RICHMEDIA    FPDF_UNSP = 15
	FPDF_UNSP_ANNOT_ATTACHMENT          FPDF_UNSP = 16
	FPDF_UNSP_ANNOT_SIG                 FPDF_UNSP = 17
)

type FPDF_FXFONT_CHARSET int

const (
	FPDF_FXFONT_ANSI_CHARSET            FPDF_FXFONT_CHARSET = 0
	FPDF_FXFONT_DEFAULT_CHARSET         FPDF_FXFONT_CHARSET = 1
	FPDF_FXFONT_SYMBOL_CHARSET          FPDF_FXFONT_CHARSET = 2
	FPDF_FXFONT_SHIFTJIS_CHARSET        FPDF_FXFONT_CHARSET = 128
	FPDF_FXFONT_HANGEUL_CHARSET         FPDF_FXFONT_CHARSET = 129
	FPDF_FXFONT_GB2312_CHARSET          FPDF_FXFONT_CHARSET = 134
	FPDF_FXFONT_CHINESEBIG5_CHARSET     FPDF_FXFONT_CHARSET = 136
	FPDF_FXFONT_GREEK_CHARSET           FPDF_FXFONT_CHARSET = 161
	FPDF_FXFONT_VIETNAMESE_CHARSET      FPDF_FXFONT_CHARSET = 163
	FPDF_FXFONT_HEBREW_CHARSET          FPDF_FXFONT_CHARSET = 177
	FPDF_FXFONT_ARABIC_CHARSET          FPDF_FXFONT_CHARSET = 178
	FPDF_FXFONT_CYRILLIC_CHARSET        FPDF_FXFONT_CHARSET = 204
	FPDF_FXFONT_THAI_CHARSET            FPDF_FXFONT_CHARSET = 22
	FPDF_FXFONT_EASTERNEUROPEAN_CHARSET FPDF_FXFONT_CHARSET = 238
)

type FPDF_OBJECT_TYPE int

const (
	FPDF_OBJECT_TYPE_UNKNOWN    FPDF_OBJECT_TYPE = 0
	FPDF_OBJECT_TYPE_BOOLEAN    FPDF_OBJECT_TYPE = 1
	FPDF_OBJECT_TYPE_NUMBER     FPDF_OBJECT_TYPE = 2
	FPDF_OBJECT_TYPE_STRING     FPDF_OBJECT_TYPE = 3
	FPDF_OBJECT_TYPE_NAME       FPDF_OBJECT_TYPE = 4
	FPDF_OBJECT_TYPE_ARRAY      FPDF_OBJECT_TYPE = 5
	FPDF_OBJECT_TYPE_DICTIONARY FPDF_OBJECT_TYPE = 6
	FPDF_OBJECT_TYPE_STREAM     FPDF_OBJECT_TYPE = 7
	FPDF_OBJECT_TYPE_NULLOBJ    FPDF_OBJECT_TYPE = 8
	FPDF_OBJECT_TYPE_REFERENCE  FPDF_OBJECT_TYPE = 9
)

type FPDF_TEXT_RENDERMODE int

const (
	FPDF_TEXTRENDERMODE_UNKNOWN          FPDF_TEXT_RENDERMODE = -1
	FPDF_TEXTRENDERMODE_FILL             FPDF_TEXT_RENDERMODE = 0
	FPDF_TEXTRENDERMODE_STROKE           FPDF_TEXT_RENDERMODE = 1
	FPDF_TEXTRENDERMODE_FILL_STROKE      FPDF_TEXT_RENDERMODE = 2
	FPDF_TEXTRENDERMODE_INVISIBLE        FPDF_TEXT_RENDERMODE = 3
	FPDF_TEXTRENDERMODE_FILL_CLIP        FPDF_TEXT_RENDERMODE = 4
	FPDF_TEXTRENDERMODE_STROKE_CLIP      FPDF_TEXT_RENDERMODE = 5
	FPDF_TEXTRENDERMODE_FILL_STROKE_CLIP FPDF_TEXT_RENDERMODE = 6
	FPDF_TEXTRENDERMODE_CLIP             FPDF_TEXT_RENDERMODE = 7
)

type FPDF_BITMAP_FORMAT int

const (
	FPDF_BITMAP_FORMAT_UNKNOWN FPDF_BITMAP_FORMAT = 0 // Unknown or unsupported format.
	FPDF_BITMAP_FORMAT_GRAY    FPDF_BITMAP_FORMAT = 1 // Gray scale bitmap, one byte per pixel.
	FPDF_BITMAP_FORMAT_BGR     FPDF_BITMAP_FORMAT = 2 // 3 bytes per pixel, byte order: blue, green, red.
	FPDF_BITMAP_FORMAT_BGRX    FPDF_BITMAP_FORMAT = 3 // 4 bytes per pixel, byte order: blue, green, red, unused.
	FPDF_BITMAP_FORMAT_BGRA    FPDF_BITMAP_FORMAT = 4 // 4 bytes per pixel, byte order: blue, green, red, alpha.
)

type FPDF_DUPLEXTYPE int

const (
	FPDF_DUPLEXTYPE_UNDEFINED              FPDF_DUPLEXTYPE = 0
	FPDF_DUPLEXTYPE_SIMPLEX                FPDF_DUPLEXTYPE = 1
	FPDF_DUPLEXTYPE_DUPLEX_FLIP_SHORT_EDGE FPDF_DUPLEXTYPE = 2
	FPDF_DUPLEXTYPE_DUPLEX_FLIP_LONG_EDGE  FPDF_DUPLEXTYPE = 3
)

type FPDF_RENDER_FLAG int

const (
	FPDF_RENDER_FLAG_ANNOT                    FPDF_RENDER_FLAG = 0x01   // Set if annotations are to be rendered.
	FPDF_RENDER_FLAG_LCD_TEXT                 FPDF_RENDER_FLAG = 0x02   // Set if using text rendering optimized for LCD display. This flag will only take effect if anti-aliasing is enabled for text.
	FPDF_RENDER_FLAG_NO_NATIVETEXT            FPDF_RENDER_FLAG = 0x04   // Don't use the native text output available on some platforms.
	FPDF_RENDER_FLAG_GRAYSCALE                FPDF_RENDER_FLAG = 0x08   // Grayscale output.
	FPDF_RENDER_FLAG_DEBUG_INFO               FPDF_RENDER_FLAG = 0x80   // Obsolete, has no effect, retained for compatibility.
	FPDF_RENDER_FLAG_NO_CATCH                 FPDF_RENDER_FLAG = 0x100  // Obsolete, has no effect, retained for compatibility.
	FPDF_RENDER_FLAG_RENDER_LIMITEDIMAGECACHE FPDF_RENDER_FLAG = 0x200  // Limit image cache size.
	FPDF_RENDER_FLAG_RENDER_FORCEHALFTONE     FPDF_RENDER_FLAG = 0x400  // Always use halftone for image stretching.
	FPDF_RENDER_FLAG_PRINTING                 FPDF_RENDER_FLAG = 0x800  // Render for printing.
	FPDF_RENDER_FLAG_RENDER_NO_SMOOTHTEXT     FPDF_RENDER_FLAG = 0x1000 // Set to disable anti-aliasing on text. This flag will also disable LCD optimization for text rendering.
	FPDF_RENDER_FLAG_RENDER_NO_SMOOTHIMAGE    FPDF_RENDER_FLAG = 0x2000 // Set to disable anti-aliasing on images.
	FPDF_RENDER_FLAG_RENDER_NO_SMOOTHPATH     FPDF_RENDER_FLAG = 0x4000 // Set to disable anti-aliasing on paths.
	FPDF_RENDER_FLAG_REVERSE_BYTE_ORDER       FPDF_RENDER_FLAG = 0x10   // Set whether to render in a reverse Byte order, this flag is only used when rendering to a bitmap.
	FPDF_RENDER_FLAG_CONVERT_FILL_TO_STROKE   FPDF_RENDER_FLAG = 0x20   // Set whether fill paths need to be stroked. This flag is only used when FPDF_COLORSCHEME is passed in, since with a single fill color for paths the boundaries of adjacent fill paths are less visible.
)

type FPDF_PRINTMODE int

const (
	FPDF_PRINTMODE_EMF                            FPDF_PRINTMODE = 0 // To output EMF (default)
	FPDF_PRINTMODE_TEXTONLY                       FPDF_PRINTMODE = 1 // to output text only (for charstream devices)
	FPDF_PRINTMODE_POSTSCRIPT2                    FPDF_PRINTMODE = 2 // to output level 2 PostScript into EMF as a series of GDI comments.
	FPDF_PRINTMODE_POSTSCRIPT3                    FPDF_PRINTMODE = 3 // to output level 3 PostScript into EMF as a series of GDI comments.
	FPDF_PRINTMODE_POSTSCRIPT2_PASSTHROUGH        FPDF_PRINTMODE = 4 // to output level 2 PostScript via ExtEscape() in PASSTHROUGH mode.
	FPDF_PRINTMODE_POSTSCRIPT3_PASSTHROUGH        FPDF_PRINTMODE = 5 // to output level 3 PostScript via ExtEscape() in PASSTHROUGH mode.
	FPDF_PRINTMODE_EMF_IMAGE_MASKS                FPDF_PRINTMODE = 6 // to output EMF, with more efficient processing of documents containing image masks.
	FPDF_PRINTMODE_POSTSCRIPT3_TYPE42             FPDF_PRINTMODE = 7 // to output level 3 PostScript with embedded Type 42 fonts, when applicable, into EMF as a series of GDI comments.
	FPDF_PRINTMODE_POSTSCRIPT3_TYPE42_PASSTHROUGH FPDF_PRINTMODE = 8 // to output level 3 PostScript with embedded Type 42 fonts, when applicable, via ExtEscape() in PASSTHROUGH mode.
)

type FPDF_RENDER_STATUS int

const (
	FPDF_RENDER_STATUS_READY         FPDF_RENDER_STATUS = 0
	FPDF_RENDER_STATUS_TOBECONTINUED FPDF_RENDER_STATUS = 1
	FPDF_RENDER_STATUS_DONE          FPDF_RENDER_STATUS = 2
	FPDF_RENDER_STATUS_FAILED        FPDF_RENDER_STATUS = 3
)

type PDF_FILEAVAIL_LINEARIZATION int

const (
	PDF_FILEAVAIL_LINEARIZATION_UNKNOWN        PDF_FILEAVAIL_LINEARIZATION = -1
	PDF_FILEAVAIL_LINEARIZATION_NOT_LINEARIZED PDF_FILEAVAIL_LINEARIZATION = 0
	PDF_FILEAVAIL_LINEARIZATION_LINEARIZED     PDF_FILEAVAIL_LINEARIZATION = 1
)

type PDF_FILEAVAIL_DATA int

const (
	PDF_FILEAVAIL_DATA_ERROR    PDF_FILEAVAIL_DATA = -1
	PDF_FILEAVAIL_DATA_NOTAVAIL PDF_FILEAVAIL_DATA = 0
	PDF_FILEAVAIL_DATA_AVAIL    PDF_FILEAVAIL_DATA = 1
)

type PDF_FILEAVAIL_FORM int

const (
	PDF_FILEAVAIL_FORM_ERROR    PDF_FILEAVAIL_FORM = -1
	PDF_FILEAVAIL_FORM_NOTAVAIL PDF_FILEAVAIL_FORM = 0
	PDF_FILEAVAIL_FORM_AVAIL    PDF_FILEAVAIL_FORM = 1
	PDF_FILEAVAIL_FORM_NOTEXIST PDF_FILEAVAIL_FORM = 2
)

type PDF_BLEND_MODE string

const (
	PDF_BLEND_MODE_COLOR       PDF_BLEND_MODE = "Color"
	PDF_BLEND_MODE_COLOR_BURN  PDF_BLEND_MODE = "ColorBurn"
	PDF_BLEND_MODE_COLOR_DODGE PDF_BLEND_MODE = "ColorDodge"
	PDF_BLEND_MODE_DARKEN      PDF_BLEND_MODE = "Darken"
	PDF_BLEND_MODE_DIFFERENCE  PDF_BLEND_MODE = "Difference"
	PDF_BLEND_MODE_EXCLUSION   PDF_BLEND_MODE = "Exclusion"
	PDF_BLEND_MODE_HARD_LIGHT  PDF_BLEND_MODE = "HardLight"
	PDF_BLEND_MODE_HUE         PDF_BLEND_MODE = "Hue"
	PDF_BLEND_MODE_LIGHTEN     PDF_BLEND_MODE = "Lighten"
	PDF_BLEND_MODE_LUMINOSITY  PDF_BLEND_MODE = "Luminosity"
	PDF_BLEND_MODE_MULTIPLY    PDF_BLEND_MODE = "Multiply"
	PDF_BLEND_MODE_NORMAL      PDF_BLEND_MODE = "Normal"
	PDF_BLEND_MODE_OVERLAY     PDF_BLEND_MODE = "Overlay"
	PDF_BLEND_MODE_SATURATION  PDF_BLEND_MODE = "Saturation"
	PDF_BLEND_MODE_SCREEN      PDF_BLEND_MODE = "Screen"
	PDF_BLEND_MODE_SOFT_LIGHT  PDF_BLEND_MODE = "SoftLight"
)

type FPDF_FILLMODE int

const (
	FPDF_FILLMODE_NONE      FPDF_FILLMODE = 0
	FPDF_FILLMODE_ALTERNATE FPDF_FILLMODE = 1
	FPDF_FILLMODE_WINDING   FPDF_FILLMODE = 2
)

type FPDF_FONT int

const (
	FPDF_FONT_TYPE1    FPDF_FONT = 1
	FPDF_FONT_TRUETYPE FPDF_FONT = 2
)

type FPDF_PAGEOBJ int

const (
	FPDF_PAGEOBJ_UNKNOWN FPDF_PAGEOBJ = 0
	FPDF_PAGEOBJ_TEXT    FPDF_PAGEOBJ = 1
	FPDF_PAGEOBJ_PATH    FPDF_PAGEOBJ = 2
	FPDF_PAGEOBJ_IMAGE   FPDF_PAGEOBJ = 3
	FPDF_PAGEOBJ_SHADING FPDF_PAGEOBJ = 4
	FPDF_PAGEOBJ_FORM    FPDF_PAGEOBJ = 5
)

type FPDF_LINEJOIN int

const (
	FPDF_LINEJOIN_MITER FPDF_LINEJOIN = 0
	FPDF_LINEJOIN_ROUND FPDF_LINEJOIN = 1
	FPDF_LINEJOIN_BEVEL FPDF_LINEJOIN = 2
)

type FPDF_LINECAP int

const (
	FPDF_LINECAP_BUTT             FPDF_LINECAP = 0
	FPDF_LINECAP_ROUND            FPDF_LINECAP = 1
	FPDF_LINECAP_PROJECTING_SQUAR FPDF_LINECAP = 2
)

type FPDF_SEGMENT int

const (
	FPDF_SEGMENT_UNKNOWN  FPDF_SEGMENT = -1
	FPDF_SEGMENT_LINETO   FPDF_SEGMENT = 0
	FPDF_SEGMENT_BEZIERTO FPDF_SEGMENT = 1
	FPDF_SEGMENT_MOVETO   FPDF_SEGMENT = 2
)

// Refer to PDF Reference version 1.7 table 4.12 for all color space families.
type FPDF_COLORSPACE int

const (
	FPDF_COLORSPACE_UNKNOWN    FPDF_COLORSPACE = 0
	FPDF_COLORSPACE_DEVICEGRAY FPDF_COLORSPACE = 1
	FPDF_COLORSPACE_DEVICERGB  FPDF_COLORSPACE = 2
	FPDF_COLORSPACE_DEVICECMYK FPDF_COLORSPACE = 3
	FPDF_COLORSPACE_CALGRAY    FPDF_COLORSPACE = 4
	FPDF_COLORSPACE_CALRGB     FPDF_COLORSPACE = 5
	FPDF_COLORSPACE_LAB        FPDF_COLORSPACE = 6
	FPDF_COLORSPACE_ICCBASED   FPDF_COLORSPACE = 7
	FPDF_COLORSPACE_SEPARATION FPDF_COLORSPACE = 8
	FPDF_COLORSPACE_DEVICEN    FPDF_COLORSPACE = 9
	FPDF_COLORSPACE_INDEXED    FPDF_COLORSPACE = 10
	FPDF_COLORSPACE_PATTERN    FPDF_COLORSPACE = 11
)

type FPDF_ANNOTATION_SUBTYPE int

const (
	FPDF_ANNOT_SUBTYPE_UNKNOWN        FPDF_ANNOTATION_SUBTYPE = 0
	FPDF_ANNOT_SUBTYPE_TEXT           FPDF_ANNOTATION_SUBTYPE = 1
	FPDF_ANNOT_SUBTYPE_LINK           FPDF_ANNOTATION_SUBTYPE = 2
	FPDF_ANNOT_SUBTYPE_FREETEXT       FPDF_ANNOTATION_SUBTYPE = 3
	FPDF_ANNOT_SUBTYPE_LINE           FPDF_ANNOTATION_SUBTYPE = 4
	FPDF_ANNOT_SUBTYPE_SQUARE         FPDF_ANNOTATION_SUBTYPE = 5
	FPDF_ANNOT_SUBTYPE_CIRCLE         FPDF_ANNOTATION_SUBTYPE = 6
	FPDF_ANNOT_SUBTYPE_POLYGON        FPDF_ANNOTATION_SUBTYPE = 7
	FPDF_ANNOT_SUBTYPE_POLYLINE       FPDF_ANNOTATION_SUBTYPE = 8
	FPDF_ANNOT_SUBTYPE_HIGHLIGHT      FPDF_ANNOTATION_SUBTYPE = 9
	FPDF_ANNOT_SUBTYPE_UNDERLINE      FPDF_ANNOTATION_SUBTYPE = 10
	FPDF_ANNOT_SUBTYPE_SQUIGGLY       FPDF_ANNOTATION_SUBTYPE = 11
	FPDF_ANNOT_SUBTYPE_STRIKEOUT      FPDF_ANNOTATION_SUBTYPE = 12
	FPDF_ANNOT_SUBTYPE_STAMP          FPDF_ANNOTATION_SUBTYPE = 13
	FPDF_ANNOT_SUBTYPE_CARET          FPDF_ANNOTATION_SUBTYPE = 14
	FPDF_ANNOT_SUBTYPE_INK            FPDF_ANNOTATION_SUBTYPE = 15
	FPDF_ANNOT_SUBTYPE_POPUP          FPDF_ANNOTATION_SUBTYPE = 16
	FPDF_ANNOT_SUBTYPE_FILEATTACHMENT FPDF_ANNOTATION_SUBTYPE = 17
	FPDF_ANNOT_SUBTYPE_SOUND          FPDF_ANNOTATION_SUBTYPE = 18
	FPDF_ANNOT_SUBTYPE_MOVIE          FPDF_ANNOTATION_SUBTYPE = 19
	FPDF_ANNOT_SUBTYPE_WIDGET         FPDF_ANNOTATION_SUBTYPE = 20
	FPDF_ANNOT_SUBTYPE_SCREEN         FPDF_ANNOTATION_SUBTYPE = 21
	FPDF_ANNOT_SUBTYPE_PRINTERMARK    FPDF_ANNOTATION_SUBTYPE = 22
	FPDF_ANNOT_SUBTYPE_TRAPNET        FPDF_ANNOTATION_SUBTYPE = 23
	FPDF_ANNOT_SUBTYPE_WATERMARK      FPDF_ANNOTATION_SUBTYPE = 24
	FPDF_ANNOT_SUBTYPE_THREED         FPDF_ANNOTATION_SUBTYPE = 25
	FPDF_ANNOT_SUBTYPE_RICHMEDIA      FPDF_ANNOTATION_SUBTYPE = 26
	FPDF_ANNOT_SUBTYPE_XFAWIDGET      FPDF_ANNOTATION_SUBTYPE = 27
	FPDF_ANNOT_SUBTYPE_REDACT         FPDF_ANNOTATION_SUBTYPE = 28
)

type FPDFANNOT_COLORTYPE int

const (
	FPDFANNOT_COLORTYPE_Color         FPDFANNOT_COLORTYPE = 0
	FPDFANNOT_COLORTYPE_InteriorColor FPDFANNOT_COLORTYPE = 1
)

type FPDF_ANNOT_APPEARANCEMODE int

const (
	FPDF_ANNOT_APPEARANCEMODE_NORMAL   FPDF_ANNOT_APPEARANCEMODE = 0
	FPDF_ANNOT_APPEARANCEMODE_ROLLOVER FPDF_ANNOT_APPEARANCEMODE = 1
	FPDF_ANNOT_APPEARANCEMODE_DOWN     FPDF_ANNOT_APPEARANCEMODE = 2
	FPDF_ANNOT_APPEARANCEMODE_COUNT    FPDF_ANNOT_APPEARANCEMODE = 3
)

type FPDF_FORMFIELD_TYPE int

const (
	FPDF_FORMFIELD_TYPE_UNKNOWN        FPDF_FORMFIELD_TYPE = 0
	FPDF_FORMFIELD_TYPE_PUSHBUTTON     FPDF_FORMFIELD_TYPE = 1
	FPDF_FORMFIELD_TYPE_CHECKBOX       FPDF_FORMFIELD_TYPE = 2
	FPDF_FORMFIELD_TYPE_RADIOBUTTON    FPDF_FORMFIELD_TYPE = 3
	FPDF_FORMFIELD_TYPE_COMBOBOX       FPDF_FORMFIELD_TYPE = 4
	FPDF_FORMFIELD_TYPE_LISTBOX        FPDF_FORMFIELD_TYPE = 5
	FPDF_FORMFIELD_TYPE_TEXTFIELD      FPDF_FORMFIELD_TYPE = 6
	FPDF_FORMFIELD_TYPE_SIGNATURE      FPDF_FORMFIELD_TYPE = 7
	FPDF_FORMFIELD_TYPE_XFA            FPDF_FORMFIELD_TYPE = 8
	FPDF_FORMFIELD_TYPE_XFA_CHECKBOX   FPDF_FORMFIELD_TYPE = 9
	FPDF_FORMFIELD_TYPE_XFA_COMBOBOX   FPDF_FORMFIELD_TYPE = 10
	FPDF_FORMFIELD_TYPE_XFA_IMAGEFIELD FPDF_FORMFIELD_TYPE = 11
	FPDF_FORMFIELD_TYPE_XFA_LISTBOX    FPDF_FORMFIELD_TYPE = 12
	FPDF_FORMFIELD_TYPE_XFA_PUSHBUTTON FPDF_FORMFIELD_TYPE = 13
	FPDF_FORMFIELD_TYPE_XFA_SIGNATURE  FPDF_FORMFIELD_TYPE = 14
	FPDF_FORMFIELD_TYPE_XFA_TEXTFIELD  FPDF_FORMFIELD_TYPE = 15
)
