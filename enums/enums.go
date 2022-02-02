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
