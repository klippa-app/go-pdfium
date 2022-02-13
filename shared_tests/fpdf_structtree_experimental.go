//go:build pdfium_experimental
// +build pdfium_experimental

package shared_tests

import (
	"github.com/klippa-app/go-pdfium/responses"
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_structtree_experimental", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	Context("no struct element", func() {
		When("is opened", func() {
			It("returns an error when calling FPDF_StructElement_GetID", func() {
				FPDF_StructElement_GetID, err := PdfiumInstance.FPDF_StructElement_GetID(&requests.FPDF_StructElement_GetID{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetID).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_GetLang", func() {
				FPDF_StructElement_GetLang, err := PdfiumInstance.FPDF_StructElement_GetLang(&requests.FPDF_StructElement_GetLang{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetLang).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_GetStringAttribute", func() {
				FPDF_StructElement_GetStringAttribute, err := PdfiumInstance.FPDF_StructElement_GetStringAttribute(&requests.FPDF_StructElement_GetStringAttribute{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetStringAttribute).To(BeNil())
			})
		})
	})

	Context("a normal PDF file with a tagged table", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/tagged_table.pdf")
			Expect(err).To(BeNil())

			newDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())

			doc = newDoc.Document
		})

		AfterEach(func() {
			FPDF_CloseDocument, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: doc,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))
		})

		When("is opened", func() {
			When("a page structtree is opened", func() {
				var structTree references.FPDF_STRUCTTREE

				BeforeEach(func() {
					FPDF_StructTree_GetForPage, err := PdfiumInstance.FPDF_StructTree_GetForPage(&requests.FPDF_StructTree_GetForPage{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})
					Expect(err).To(BeNil())
					Expect(FPDF_StructTree_GetForPage).To(Not(BeNil()))

					structTree = FPDF_StructTree_GetForPage.StructTree
				})

				AfterEach(func() {
					FPDF_StructTree_Close, err := PdfiumInstance.FPDF_StructTree_Close(&requests.FPDF_StructTree_Close{
						StructTree: structTree,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_StructTree_Close).To(Not(BeNil()))
				})

				It("returns the correct struct tree children count", func() {
					FPDF_StructTree_CountChildren, err := PdfiumInstance.FPDF_StructTree_CountChildren(&requests.FPDF_StructTree_CountChildren{
						StructTree: structTree,
					})
					Expect(err).To(BeNil())
					Expect(FPDF_StructTree_CountChildren).To(Equal(&responses.FPDF_StructTree_CountChildren{
						Count: 1,
					}))
				})

				When("a struct tree struct element is opened", func() {
					var structElement references.FPDF_STRUCTELEMENT

					BeforeEach(func() {
						FPDF_StructTree_GetChildAtIndex, err := PdfiumInstance.FPDF_StructTree_GetChildAtIndex(&requests.FPDF_StructTree_GetChildAtIndex{
							StructTree: structTree,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructTree_GetChildAtIndex).To(Not(BeNil()))

						structElement = FPDF_StructTree_GetChildAtIndex.StructElement
					})

					It("returns the correct struct element children count", func() {
						FPDF_StructElement_CountChildren, err := PdfiumInstance.FPDF_StructElement_CountChildren(&requests.FPDF_StructElement_CountChildren{
							StructElement: structElement,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_CountChildren).To(Equal(&responses.FPDF_StructElement_CountChildren{
							Count: 1,
						}))
					})

					It("returns the correct element type", func() {
						FPDF_StructElement_GetType, err := PdfiumInstance.FPDF_StructElement_GetType(&requests.FPDF_StructElement_GetType{
							StructElement: structElement,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetType).To(Equal(&responses.FPDF_StructElement_GetType{
							Type: "Document",
						}))
					})

					It("returns an error when not giving a name when requesting an attribute name", func() {
						table, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: structElement,
						})
						Expect(err).To(BeNil())
						Expect(table).To(Not(BeNil()))

						row, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: table.StructElement,
						})
						Expect(err).To(BeNil())
						Expect(row).To(Not(BeNil()))

						header_cell, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: row.StructElement,
						})
						Expect(err).To(BeNil())
						Expect(header_cell).To(Not(BeNil()))

						FPDF_StructElement_GetStringAttribute, err := PdfiumInstance.FPDF_StructElement_GetStringAttribute(&requests.FPDF_StructElement_GetStringAttribute{
							StructElement: header_cell.StructElement,
						})
						Expect(err).To(MatchError("could not get attribute"))
						Expect(FPDF_StructElement_GetStringAttribute).To(BeNil())
					})

					It("returns the correct string attribute", func() {
						table, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: structElement,
						})
						Expect(err).To(BeNil())
						Expect(table).To(Not(BeNil()))

						row, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: table.StructElement,
						})
						Expect(err).To(BeNil())
						Expect(row).To(Not(BeNil()))

						header_cell, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: row.StructElement,
						})
						Expect(err).To(BeNil())
						Expect(header_cell).To(Not(BeNil()))

						FPDF_StructElement_GetStringAttribute, err := PdfiumInstance.FPDF_StructElement_GetStringAttribute(&requests.FPDF_StructElement_GetStringAttribute{
							StructElement: header_cell.StructElement,
							AttributeName: "Scope",
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetStringAttribute).To(Equal(&responses.FPDF_StructElement_GetStringAttribute{
							Attribute: "Scope",
							Value:     "Row",
						}))
					})

					It("returns the correct ID", func() {
						table, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: structElement,
						})
						Expect(err).To(BeNil())
						Expect(table).To(Not(BeNil()))

						FPDF_StructElement_GetID, err := PdfiumInstance.FPDF_StructElement_GetID(&requests.FPDF_StructElement_GetID{
							StructElement: table.StructElement,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetID).To(Equal(&responses.FPDF_StructElement_GetID{
							ID: "node12",
						}))
					})

					It("returns the correct lang", func() {
						table, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: structElement,
						})
						Expect(err).To(BeNil())
						Expect(table).To(Not(BeNil()))

						FPDF_StructElement_GetLang, err := PdfiumInstance.FPDF_StructElement_GetLang(&requests.FPDF_StructElement_GetLang{
							StructElement: table.StructElement,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetLang).To(Equal(&responses.FPDF_StructElement_GetLang{
							Lang: "hu",
						}))
					})
				})
			})
		})
	})
})
