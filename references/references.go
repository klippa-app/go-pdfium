package references

// All the references are just an alias for string, but it's useful to keep the alias,
// so that we can change it later. They all should contain a UUID.

// FPDF_DOCUMENT is an internal reference to a C.FPDF_DOCUMENT object.
type FPDF_DOCUMENT string

// FPDF_PAGE is an internal reference to a C.FPDF_PAGE object.
type FPDF_PAGE string

// FPDF_BOOKMARK is an internal reference to a C.FPDF_BOOKMARK object.
type FPDF_BOOKMARK string
