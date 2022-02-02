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
type FPDF_PAGE_AACTION uint32

const (
	FPDF_PAGE_AACTION_OPEN  FPDF_PAGE_AACTION = 0 // OPEN (/O) -- An action to be performed when the page is opened
	FPDF_PAGE_AACTION_CLOSE FPDF_PAGE_AACTION = 1 // CLOSE (/C) -- An action to be performed when the page is closed
)

// Additional actions type of document
type FPDF_DOC_AACTION uint32

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
