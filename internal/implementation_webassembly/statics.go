package implementation_webassembly

type FPDF_ERR int

const (
	FPDF_ERR_SUCCESS  FPDF_ERR = 0
	FPDF_ERR_UNKNOWN           = 1
	FPDF_ERR_FILE              = 2
	FPDF_ERR_FORMAT            = 3
	FPDF_ERR_PASSWORD          = 4
	FPDF_ERR_SECURITY          = 5
	FPDF_ERR_PAGE              = 6
)
