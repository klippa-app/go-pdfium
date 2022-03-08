package references

// All the references are just an alias for string, but it's useful to keep the alias,
// so that we can change it later. They all should contain a UUID.

// FPDF_DOCUMENT is an internal reference to a C.FPDF_DOCUMENT handle.
type FPDF_DOCUMENT string

// FPDF_PAGE is an internal reference to a C.FPDF_PAGE handle.
type FPDF_PAGE string

// FPDF_BOOKMARK is an internal reference to a C.FPDF_BOOKMARK handle.
type FPDF_BOOKMARK string

// FPDF_DEST is an internal reference to a C.FPDF_DEST handle.
type FPDF_DEST string

// FPDF_ACTION is an internal reference to a C.FPDF_ACTION handle.
type FPDF_ACTION string

// FPDF_LINK is an internal reference to a C.FPDF_LINK handle.
type FPDF_LINK string

// FPDF_PAGELINK is an internal reference to a C.FPDF_PAGELINK handle.
type FPDF_PAGELINK string

// FPDF_SCHHANDLE is an internal reference to a C.FPDF_SCHHANDLE handle.
type FPDF_SCHHANDLE string

// FPDF_BITMAP is an internal reference to a C.FPDF_BITMAP handle.
type FPDF_BITMAP string

// FPDF_TEXTPAGE is an internal reference to a C.FPDF_TEXTPAGE handle.
type FPDF_TEXTPAGE string

// FPDF_PAGERANGE is an internal reference to a C.FPDF_PAGERANGE handle.
type FPDF_PAGERANGE string

// FPDF_PAGEOBJECT is an internal reference to a C.FPDF_PAGEOBJECT handle.
type FPDF_PAGEOBJECT string

// FPDF_CLIPPATH is an internal reference to a C.FPDF_CLIPPATH handle.
type FPDF_CLIPPATH string

// FPDF_FORMHANDLE is an internal reference to a C.FPDF_FORMHANDLE handle.
type FPDF_FORMHANDLE string

// FPDF_ANNOTATION is an internal reference to a C.FPDF_ANNOTATION handle.
type FPDF_ANNOTATION string

// FPDF_XOBJECT is an internal reference to a C.FPDF_XOBJECT handle.
type FPDF_XOBJECT string

// FPDF_SIGNATURE is an internal reference to a C.FPDF_SIGNATURE handle.
type FPDF_SIGNATURE string

// FPDF_ATTACHMENT is an internal reference to a C.FPDF_ATTACHMENT handle.
type FPDF_ATTACHMENT string

// FPDF_JAVASCRIPT_ACTION is an internal reference to a C.FPDF_JAVASCRIPT_ACTION handle.
type FPDF_JAVASCRIPT_ACTION string

// FPDF_PATHSEGMENT is an internal reference to a C.FPDF_PATHSEGMENT handle.
type FPDF_PATHSEGMENT string

// FPDF_AVAIL is an internal reference to a C.FPDF_AVAIL handle.
type FPDF_AVAIL string

// FPDF_STRUCTTREE is an internal reference to a C.FPDF_STRUCTTREE handle.
type FPDF_STRUCTTREE string

// FPDF_STRUCTELEMENT is an internal reference to a C.FPDF_STRUCTELEMENT handle.
type FPDF_STRUCTELEMENT string

// FPDF_PAGEOBJECTMARK is an internal reference to a C.FPDF_PAGEOBJECTMARK handle.
type FPDF_PAGEOBJECTMARK string

// FPDF_FONT is an internal reference to a C.FPDF_FONT handle.
type FPDF_FONT string

// FPDF_GLYPHPATH is an internal reference to a C.FPDF_GLYPHPATH handle.
type FPDF_GLYPHPATH string

// FPDF_STRUCTELEMENT_ATTR is an internal reference to a C.FPDF_STRUCTELEMENT_ATTR handle.
type FPDF_STRUCTELEMENT_ATTR string
