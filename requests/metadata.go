package requests

import "github.com/klippa-app/go-pdfium/references"

type GetMetaData struct {
	Document references.FPDF_DOCUMENT
	Tags     *[]string // A list of metadata tags. If nil, it will return: Title, Author, Subject, Keywords, Creator, Producer, CreationDate, ModDate. For detailed explanation of these tags and their respective values, please refer to section 10.2.1 "Document Information Dictionary" in PDF Reference 1.7.
}
