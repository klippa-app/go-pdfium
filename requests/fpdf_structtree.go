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
