package responses

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
)

type FPDF_StructTree_GetForPage struct {
	StructTree references.FPDF_STRUCTTREE
}

type FPDF_StructTree_Close struct{}

type FPDF_StructTree_CountChildren struct {
	Count int
}

type FPDF_StructTree_GetChildAtIndex struct {
	StructElement references.FPDF_STRUCTELEMENT
}

type FPDF_StructElement_GetAltText struct {
	AltText string
}

type FPDF_StructElement_GetID struct {
	ID string
}

type FPDF_StructElement_GetLang struct {
	Lang string // The case-insensitive IETF BCP 47 language code for an element.
}

type FPDF_StructElement_GetStringAttribute struct {
	Attribute string
	Value     string
}

type FPDF_StructElement_GetMarkedContentID struct {
	MarkedContentID int // The marked content ID of the element. If no ID exists, returns -1.
}

type FPDF_StructElement_GetType struct {
	Type string
}

type FPDF_StructElement_GetTitle struct {
	Title string
}

type FPDF_StructElement_CountChildren struct {
	Count int
}

type FPDF_StructElement_GetChildAtIndex struct {
	StructElement references.FPDF_STRUCTELEMENT
}

type FPDF_StructElement_GetChildMarkedContentID struct {
	MarkedContentID int // The marked content ID of the child. If no ID exists, returns -1.
}

type FPDF_StructElement_GetActualText struct {
	Actualtext string
}

type FPDF_StructElement_GetObjType struct {
	ObjType string
}

type FPDF_StructElement_GetParent struct {
	StructElement references.FPDF_STRUCTELEMENT
}

type FPDF_StructElement_GetAttributeCount struct {
	Count int
}

type FPDF_StructElement_GetAttributeAtIndex struct {
	StructElementAttribute references.FPDF_STRUCTELEMENT_ATTR
}

type FPDF_StructElement_Attr_GetCount struct {
	Count int
}

type FPDF_StructElement_Attr_GetName struct {
	Name string
}

type FPDF_StructElement_Attr_GetType struct {
	ObjectType enums.FPDF_OBJECT_TYPE
}

type FPDF_StructElement_Attr_GetBooleanValue struct {
	Value bool
}

type FPDF_StructElement_Attr_GetNumberValue struct {
	Value float32
}

type FPDF_StructElement_Attr_GetStringValue struct {
	Value string
}

type FPDF_StructElement_Attr_GetBlobValue struct {
	Value []byte
}

type FPDF_StructElement_GetMarkedContentIdCount struct {
	Count int
}

type FPDF_StructElement_GetMarkedContentIdAtIndex struct {
	MarkedContentID int
}
