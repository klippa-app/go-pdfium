package responses

type GetMetaDataTag struct {
	Tag   string // The requested metadata tag.
	Value string // The value of the tag if found, string is empty if the value is not found.
}

type GetMetaData struct {
	Tags []GetMetaDataTag
}
