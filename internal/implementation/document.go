package implementation

// #cgo pkg-config: pdfium
// #include "fpdf_doc.h"
import "C"
import (
	"github.com/google/uuid"
	"github.com/klippa-app/go-pdfium/references"
)

func (p *PdfiumImplementation) registerDocument(document C.FPDF_DOCUMENT) *DocumentHandle {
	documentRef := uuid.New()
	documentHandle := &DocumentHandle{
		handle:               document,
		nativeRef:            references.FPDF_DOCUMENT(documentRef.String()),
		pageRefs:             map[references.FPDF_PAGE]*PageHandle{},
		bookmarkRefs:         map[references.FPDF_BOOKMARK]*BookmarkHandle{},
		destRefs:             map[references.FPDF_DEST]*DestHandle{},
		actionRefs:           map[references.FPDF_ACTION]*ActionHandle{},
		linkRefs:             map[references.FPDF_LINK]*LinkHandle{},
		pageLinkRefs:         map[references.FPDF_PAGELINK]*PageLinkHandle{},
		schHandleRefs:        map[references.FPDF_SCHHANDLE]*SchHandleHandle{},
		textPageRefs:         map[references.FPDF_TEXTPAGE]*TextPageHandle{},
		pageRangeRefs:        map[references.FPDF_PAGERANGE]*PageRangeHandle{},
		formHandleRefs:       map[references.FPDF_FORMHANDLE]*FormHandleHandle{},
		annotationRefs:       map[references.FPDF_ANNOTATION]*AnnotationHandle{},
		signatureRefs:        map[references.FPDF_SIGNATURE]*SignatureHandle{},
		attachmentRefs:       map[references.FPDF_ATTACHMENT]*AttachmentHandle{},
		javaScriptActionRefs: map[references.FPDF_JAVASCRIPT_ACTION]*JavaScriptActionHandle{},
		searchRefs:           map[references.FPDF_SCHHANDLE]*SearchHandle{},
		structTreeRefs:       map[references.FPDF_STRUCTTREE]*StructTreeHandle{},
		structElementRefs:    map[references.FPDF_STRUCTELEMENT]*StructElementHandle{},
	}

	Pdfium.documentRefs[documentHandle.nativeRef] = documentHandle
	p.documentRefs[documentHandle.nativeRef] = documentHandle

	return documentHandle
}
