package requests

import "github.com/klippa-app/go-pdfium/references"

type FPDF_StructTree_GetForPage struct {
	Page Page
}

type FPDF_StructTree_Close struct {
	StructTree references.FPDF_STRUCTTREE
}

type FPDF_StructTree_CountChildren struct {
	StructTree references.FPDF_STRUCTTREE
}

type FPDF_StructTree_GetChildAtIndex struct {
	StructTree references.FPDF_STRUCTTREE
	Index      int
}

type FPDF_StructElement_GetAltText struct {
	StructElement references.FPDF_STRUCTELEMENT
}

type FPDF_StructElement_GetID struct {
	StructElement references.FPDF_STRUCTELEMENT
}

type FPDF_StructElement_GetLang struct {
	StructElement references.FPDF_STRUCTELEMENT
}

type FPDF_StructElement_GetStringAttribute struct {
	StructElement references.FPDF_STRUCTELEMENT
	AttributeName string // The name of the attribute to retrieve.
}

type FPDF_StructElement_GetMarkedContentID struct {
	StructElement references.FPDF_STRUCTELEMENT
}

type FPDF_StructElement_GetType struct {
	StructElement references.FPDF_STRUCTELEMENT
}

type FPDF_StructElement_GetTitle struct {
	StructElement references.FPDF_STRUCTELEMENT
}

type FPDF_StructElement_CountChildren struct {
	StructElement references.FPDF_STRUCTELEMENT
}

type FPDF_StructElement_GetChildAtIndex struct {
	StructElement references.FPDF_STRUCTELEMENT
	Index         int // The index for the child, 0-based.
}

type FPDF_StructElement_GetChildMarkedContentID struct {
	StructElement references.FPDF_STRUCTELEMENT
	Index         int // The index for the child, 0-based.
}

type FPDF_StructElement_GetActualText struct {
	StructElement references.FPDF_STRUCTELEMENT
}

type FPDF_StructElement_GetObjType struct {
	StructElement references.FPDF_STRUCTELEMENT
}

type FPDF_StructElement_GetParent struct {
	StructElement references.FPDF_STRUCTELEMENT
}

type FPDF_StructElement_GetAttributeCount struct {
	StructElement references.FPDF_STRUCTELEMENT
}

type FPDF_StructElement_GetAttributeAtIndex struct {
	StructElement references.FPDF_STRUCTELEMENT
	Index         int
}

type FPDF_StructElement_Attr_GetCount struct {
	StructElementAttribute references.FPDF_STRUCTELEMENT_ATTR
}

type FPDF_StructElement_Attr_GetName struct {
	StructElementAttribute references.FPDF_STRUCTELEMENT_ATTR
	Index                  int
}

type FPDF_StructElement_Attr_GetValue struct {
	StructElementAttribute references.FPDF_STRUCTELEMENT_ATTR
	Name                   string
}

type FPDF_StructElement_Attr_GetType struct {
	StructElementAttribute references.FPDF_STRUCTELEMENT_ATTR
	Name                   string
}

type FPDF_StructElement_Attr_GetBooleanValue struct {
	StructElementAttributeValue references.FPDF_STRUCTELEMENT_ATTR_VALUE
	Name                        string
}

type FPDF_StructElement_Attr_GetNumberValue struct {
	StructElementAttributeValue references.FPDF_STRUCTELEMENT_ATTR_VALUE
	Name                        string
}

type FPDF_StructElement_Attr_GetStringValue struct {
	StructElementAttributeValue references.FPDF_STRUCTELEMENT_ATTR_VALUE
	Name                        string
}

type FPDF_StructElement_Attr_GetBlobValue struct {
	StructElementAttributeValue references.FPDF_STRUCTELEMENT_ATTR_VALUE
	Name                        string
}

type FPDF_StructElement_Attr_CountChildren struct {
	StructElementAttributeValue references.FPDF_STRUCTELEMENT_ATTR_VALUE
}

type FPDF_StructElement_Attr_GetChildAtIndex struct {
	StructElementAttributeValue references.FPDF_STRUCTELEMENT_ATTR_VALUE
	Index                       int
}

type FPDF_StructElement_GetMarkedContentIdCount struct {
	StructElement references.FPDF_STRUCTELEMENT
}

type FPDF_StructElement_GetMarkedContentIdAtIndex struct {
	StructElement references.FPDF_STRUCTELEMENT
	Index         int
}
