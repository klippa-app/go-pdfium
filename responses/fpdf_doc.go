package responses

type FPDF_GetMetaText struct {
	Tag   string // The requested metadata tag.
	Value string // The value of the tag if found, string is empty if the value is not found.
}
