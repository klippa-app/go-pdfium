package responses

type FPDF_SaveAsCopy struct {
	FileBytes *[]byte // The byte array if no path or writer was given.
	FilePath  *string
}

type FPDF_SaveWithVersion struct {
	FileBytes *[]byte // The byte array if no path or writer was given.
	FilePath  *string
}
