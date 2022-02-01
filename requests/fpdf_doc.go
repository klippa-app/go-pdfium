package requests

import "github.com/klippa-app/go-pdfium/references"

type FPDF_GetMetaText struct {
	Document references.FPDF_DOCUMENT
	Tag      string // A metadata tag. Title, Author, Subject, Keywords, Creator, Producer, CreationDate, ModDate. For detailed explanation of these tags and their respective values, please refer to section 10.2.1 "Document Information Dictionary" in PDF Reference 1.7.
}
