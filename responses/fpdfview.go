package responses

import "github.com/klippa-app/go-pdfium/references"

type FPDF_LoadPage struct {
	Page references.FPDF_PAGE
}

type FPDF_ClosePage struct{}

type FPDF_GetFileVersion struct {
	FileVersion int // The numeric version of the file: 14 for 1.4, 15 for 1.5, ...
}

type FPDF_GetDocPermissions struct {
	DocPermissions                      uint32 // A 32-bit integer which indicates the permission flags. Please refer to "TABLE 3.20 User access permissions" in PDF Reference 1.7 P123 for detailed description. If the document is not protected, 0xffffffff (4294967295) will be returned.
	PrintDocument                       bool   // Bit position 3: (Security handlers of revision 2) Print the document, (Security handlers of revision 3 or greater) Print the document (possibly not at the highest quality level, depending on whether PrintDocumentAsFaithfulDigitalCopy (bit 12) is also set).
	ModifyContents                      bool   // Bit position 4: Modify the contents of the document by operations other than those controlled by AddOrModifyTextAnnotations (bit 6), FillInExistingInteractiveFormFields (bit 9), and AssembleDocument (bit 11).
	CopyOrExtractText                   bool   // Bit position 5: (Security handlers of revision 2) Copy or otherwise extract  text and graphics from the document, including extracting text and graphics (in support of accessibility to users with disabilities or for other purposes). (Security handlers of revision 3 or greater) Copy or otherwise extract text and graphics from the document by operations other than that controlled by ExtractTextAndGraphics (bit 10).
	AddOrModifyTextAnnotations          bool   // Bit position 6: Add or modify text annotations
	FillInInteractiveFormFields         bool   // Bit position 6: fill in interactive form fields
	CreateOrModifyInteractiveFormFields bool   // Bit position 6 & 4: create or modify interactive form fields (including signature fields).
	FillInExistingInteractiveFormFields bool   // Bit position 9: (Security handlers of revision 3 or greater) Fill in existing interactive form fields (including signature fields), even if FillInInteractiveFormFields (bit 6) is clear.
	ExtractTextAndGraphics              bool   // Bit position 10: (Security handlers of revision 3 or greater) Extract text and graphics (in support of accessibility to users with disabilities or for other purposes).
	AssembleDocument                    bool   // Bit position 11: (Security handlers of revision 3 or greater) Assemble the  document (insert, rotate, or delete pages and create bookmarks or thumbnail images), even if ModifyContents (bit 4) is clear.
	PrintDocumentAsFaithfulDigitalCopy  bool   // Bit position 12: (Security handlers of revision 3 or greater) Print the document to a representation from which a faithful digital copy of the PDF content could be generated. When this bit is clear (and PrintDocument (bit 3) is set), printing is limited to a low-level representation of the appearance, possibly of degraded quality.
}

type FPDF_GetSecurityHandlerRevision struct {
	SecurityHandlerRevision int // The revision number of security handler. Please refer to key "R" in "TABLE 3.19 Additional encryption dictionary entries for the standard security handler" in PDF Reference 1.7 P122 for detailed description. If the document is not protected, -1 will be returned.
}

type FPDF_GetPageCount struct {
	PageCount int // The amount of pages of the document.
}
