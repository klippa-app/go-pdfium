//go:build pdfium_experimental
// +build pdfium_experimental

package shared_tests

import (
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

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

			It("returns an error when calling FPDF_StructElement_GetActualText", func() {
				FPDF_StructElement_GetActualText, err := PdfiumInstance.FPDF_StructElement_GetActualText(&requests.FPDF_StructElement_GetActualText{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetActualText).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_GetObjType", func() {
				FPDF_StructElement_GetObjType, err := PdfiumInstance.FPDF_StructElement_GetObjType(&requests.FPDF_StructElement_GetObjType{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetObjType).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_GetParent", func() {
				FPDF_StructElement_GetParent, err := PdfiumInstance.FPDF_StructElement_GetParent(&requests.FPDF_StructElement_GetParent{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetParent).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_GetAttributeCount", func() {
				FPDF_StructElement_GetAttributeCount, err := PdfiumInstance.FPDF_StructElement_GetAttributeCount(&requests.FPDF_StructElement_GetAttributeCount{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetAttributeCount).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_GetAttributeAtIndex", func() {
				FPDF_StructElement_GetAttributeAtIndex, err := PdfiumInstance.FPDF_StructElement_GetAttributeAtIndex(&requests.FPDF_StructElement_GetAttributeAtIndex{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetAttributeAtIndex).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_Attr_GetCount", func() {
				FPDF_StructElement_Attr_GetCount, err := PdfiumInstance.FPDF_StructElement_Attr_GetCount(&requests.FPDF_StructElement_Attr_GetCount{})
				Expect(err).To(MatchError("structElementAttribute not given"))
				Expect(FPDF_StructElement_Attr_GetCount).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_Attr_GetName", func() {
				FPDF_StructElement_Attr_GetName, err := PdfiumInstance.FPDF_StructElement_Attr_GetName(&requests.FPDF_StructElement_Attr_GetName{})
				Expect(err).To(MatchError("structElementAttribute not given"))
				Expect(FPDF_StructElement_Attr_GetName).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_Attr_GetValue", func() {
				FPDF_StructElement_Attr_GetValue, err := PdfiumInstance.FPDF_StructElement_Attr_GetValue(&requests.FPDF_StructElement_Attr_GetValue{})
				Expect(err).To(MatchError("structElementAttribute not given"))
				Expect(FPDF_StructElement_Attr_GetValue).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_Attr_GetType", func() {
				FPDF_StructElement_Attr_GetType, err := PdfiumInstance.FPDF_StructElement_Attr_GetType(&requests.FPDF_StructElement_Attr_GetType{})
				Expect(err).To(MatchError("structElementAttributeValue not given"))
				Expect(FPDF_StructElement_Attr_GetType).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_Attr_GetBooleanValue", func() {
				FPDF_StructElement_Attr_GetBooleanValue, err := PdfiumInstance.FPDF_StructElement_Attr_GetBooleanValue(&requests.FPDF_StructElement_Attr_GetBooleanValue{})
				Expect(err).To(MatchError("structElementAttributeValue not given"))
				Expect(FPDF_StructElement_Attr_GetBooleanValue).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_Attr_GetNumberValue", func() {
				FPDF_StructElement_Attr_GetNumberValue, err := PdfiumInstance.FPDF_StructElement_Attr_GetNumberValue(&requests.FPDF_StructElement_Attr_GetNumberValue{})
				Expect(err).To(MatchError("structElementAttributeValue not given"))
				Expect(FPDF_StructElement_Attr_GetNumberValue).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_Attr_GetStringValue", func() {
				FPDF_StructElement_Attr_GetStringValue, err := PdfiumInstance.FPDF_StructElement_Attr_GetStringValue(&requests.FPDF_StructElement_Attr_GetStringValue{})
				Expect(err).To(MatchError("structElementAttributeValue not given"))
				Expect(FPDF_StructElement_Attr_GetStringValue).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_Attr_GetBlobValue", func() {
				FPDF_StructElement_Attr_GetBlobValue, err := PdfiumInstance.FPDF_StructElement_Attr_GetBlobValue(&requests.FPDF_StructElement_Attr_GetBlobValue{})
				Expect(err).To(MatchError("structElementAttributeValue not given"))
				Expect(FPDF_StructElement_Attr_GetBlobValue).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_Attr_CountChildren", func() {
				FPDF_StructElement_Attr_CountChildren, err := PdfiumInstance.FPDF_StructElement_Attr_CountChildren(&requests.FPDF_StructElement_Attr_CountChildren{})
				Expect(err).To(MatchError("structElementAttributeValue not given"))
				Expect(FPDF_StructElement_Attr_CountChildren).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_Attr_GetChildAtIndex", func() {
				FPDF_StructElement_Attr_GetChildAtIndex, err := PdfiumInstance.FPDF_StructElement_Attr_GetChildAtIndex(&requests.FPDF_StructElement_Attr_GetChildAtIndex{})
				Expect(err).To(MatchError("structElementAttributeValue not given"))
				Expect(FPDF_StructElement_Attr_GetChildAtIndex).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_GetMarkedContentIdCount", func() {
				FPDF_StructElement_GetMarkedContentIdCount, err := PdfiumInstance.FPDF_StructElement_GetMarkedContentIdCount(&requests.FPDF_StructElement_GetMarkedContentIdCount{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetMarkedContentIdCount).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_GetMarkedContentIdAtIndex", func() {
				FPDF_StructElement_GetMarkedContentIdAtIndex, err := PdfiumInstance.FPDF_StructElement_GetMarkedContentIdAtIndex(&requests.FPDF_StructElement_GetMarkedContentIdAtIndex{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetMarkedContentIdAtIndex).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_GetChildMarkedContentID", func() {
				FPDF_StructElement_GetChildMarkedContentID, err := PdfiumInstance.FPDF_StructElement_GetChildMarkedContentID(&requests.FPDF_StructElement_GetChildMarkedContentID{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetChildMarkedContentID).To(BeNil())
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

						table, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: structElement,
							Index:         0,
						})
						Expect(err).To(BeNil())
						Expect(table).To(Not(BeNil()))

						FPDF_StructElement_GetType, err = PdfiumInstance.FPDF_StructElement_GetType(&requests.FPDF_StructElement_GetType{
							StructElement: table.StructElement,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetType).To(Equal(&responses.FPDF_StructElement_GetType{
							Type: "Table",
						}))

						row, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: table.StructElement,
							Index:         0,
						})
						Expect(err).To(BeNil())
						Expect(row).To(Not(BeNil()))

						headerCell, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: row.StructElement,
							Index:         0,
						})
						Expect(err).To(BeNil())
						Expect(headerCell).To(Not(BeNil()))

						FPDF_StructElement_GetType, err = PdfiumInstance.FPDF_StructElement_GetType(&requests.FPDF_StructElement_GetType{
							StructElement: headerCell.StructElement,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetType).To(Equal(&responses.FPDF_StructElement_GetType{
							Type: "TH",
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

					It("returns the struct element attribute count", func() {
						table, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: structElement,
						})
						Expect(err).To(BeNil())
						Expect(table).To(Not(BeNil()))

						row, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: table.StructElement,
							Index:         0,
						})
						Expect(err).To(BeNil())
						Expect(row).To(Not(BeNil()))

						FPDF_StructElement_GetAttributeCount, err := PdfiumInstance.FPDF_StructElement_GetAttributeCount(&requests.FPDF_StructElement_GetAttributeCount{
							StructElement: row.StructElement,
						})
						Expect(err).To(MatchError("could not get struct element attribute count"))
						Expect(FPDF_StructElement_GetAttributeCount).To(BeNil())

						headerCell, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: row.StructElement,
							Index:         0,
						})
						Expect(err).To(BeNil())
						Expect(headerCell).To(Not(BeNil()))

						FPDF_StructElement_GetAttributeCount, err = PdfiumInstance.FPDF_StructElement_GetAttributeCount(&requests.FPDF_StructElement_GetAttributeCount{
							StructElement: headerCell.StructElement,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetAttributeCount).To(Equal(&responses.FPDF_StructElement_GetAttributeCount{
							Count: 2,
						}))
					})

					It("returns the struct element attribute", func() {
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

						headerCell, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: row.StructElement,
						})
						Expect(err).To(BeNil())
						Expect(headerCell).To(Not(BeNil()))

						FPDF_StructElement_GetAttributeAtIndex, err := PdfiumInstance.FPDF_StructElement_GetAttributeAtIndex(&requests.FPDF_StructElement_GetAttributeAtIndex{
							StructElement: headerCell.StructElement,
							Index:         3,
						})
						Expect(err).To(MatchError("could not get struct element attribute"))
						Expect(FPDF_StructElement_GetAttributeAtIndex).To(BeNil())

						FPDF_StructElement_GetAttributeAtIndex, err = PdfiumInstance.FPDF_StructElement_GetAttributeAtIndex(&requests.FPDF_StructElement_GetAttributeAtIndex{
							StructElement: headerCell.StructElement,
							Index:         1,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetAttributeAtIndex).ToNot(BeNil())
						Expect(FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute).ToNot(BeEmpty())

						FPDF_StructElement_Attr_GetCount, err := PdfiumInstance.FPDF_StructElement_Attr_GetCount(&requests.FPDF_StructElement_Attr_GetCount{
							StructElementAttribute: FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetCount).To(Equal(&responses.FPDF_StructElement_Attr_GetCount{
							Count: 2,
						}))

						FPDF_StructElement_Attr_GetName, err := PdfiumInstance.FPDF_StructElement_Attr_GetName(&requests.FPDF_StructElement_Attr_GetName{
							StructElementAttribute: FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute,
							Index:                  25,
						})
						Expect(err).To(MatchError("could not get attribute name"))
						Expect(FPDF_StructElement_Attr_GetName).To(BeNil())

						FPDF_StructElement_Attr_GetName, err = PdfiumInstance.FPDF_StructElement_Attr_GetName(&requests.FPDF_StructElement_Attr_GetName{
							StructElementAttribute: FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute,
							Index:                  0,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetName).To(Equal(&responses.FPDF_StructElement_Attr_GetName{
							Name: "ColSpan",
						}))

						FPDF_StructElement_Attr_GetValueColSpan, err := PdfiumInstance.FPDF_StructElement_Attr_GetValue(&requests.FPDF_StructElement_Attr_GetValue{
							StructElementAttribute: FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute,
							Name:                   "ColSpan",
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetValueColSpan).To(Not(BeNil()))

						FPDF_StructElement_Attr_GetNumberValue, err := PdfiumInstance.FPDF_StructElement_Attr_GetNumberValue(&requests.FPDF_StructElement_Attr_GetNumberValue{
							StructElementAttributeValue: FPDF_StructElement_Attr_GetValueColSpan.StructElementAttributeValue,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetNumberValue).To(Equal(&responses.FPDF_StructElement_Attr_GetNumberValue{
							Value: 2,
						}))

						FPDF_StructElement_Attr_GetType, err := PdfiumInstance.FPDF_StructElement_Attr_GetType(&requests.FPDF_StructElement_Attr_GetType{
							StructElementAttributeValue: FPDF_StructElement_Attr_GetValueColSpan.StructElementAttributeValue,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetType).To(Equal(&responses.FPDF_StructElement_Attr_GetType{
							ObjectType: enums.FPDF_OBJECT_TYPE_NUMBER,
						}))

						FPDF_StructElement_Attr_GetName, err = PdfiumInstance.FPDF_StructElement_Attr_GetName(&requests.FPDF_StructElement_Attr_GetName{
							StructElementAttribute: FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute,
							Index:                  1,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetName).To(Equal(&responses.FPDF_StructElement_Attr_GetName{
							Name: "O",
						}))

						// @todo: figure out why this fails.
						/*
							FPDF_StructElement_Attr_GetValue0, err := PdfiumInstance.FPDF_StructElement_Attr_GetValue(&requests.FPDF_StructElement_Attr_GetValue{
								StructElementAttribute: FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute,
								Name:                   "0",
							})
							Expect(err).To(BeNil())
							Expect(FPDF_StructElement_Attr_GetValue0).To(Not(BeNil()))
						*/
					})

					It("returns the struct element attribute values", func() {
						table, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: structElement,
							Index:         0,
						})
						Expect(err).To(BeNil())
						Expect(table.StructElement).To(Not(BeEmpty()))

						tr, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: table.StructElement,
							Index:         1,
						})
						Expect(err).To(BeNil())
						Expect(tr.StructElement).To(Not(BeEmpty()))

						td, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: tr.StructElement,
							Index:         1,
						})
						Expect(err).To(BeNil())
						Expect(td).To(Not(BeNil()))
						Expect(td.StructElement).To(Not(BeEmpty()))

						FPDF_StructElement_GetAttributeCount, err := PdfiumInstance.FPDF_StructElement_GetAttributeCount(&requests.FPDF_StructElement_GetAttributeCount{
							StructElement: td.StructElement,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetAttributeCount).To(Equal(&responses.FPDF_StructElement_GetAttributeCount{
							Count: 1,
						}))

						FPDF_StructElement_GetAttributeAtIndex, err := PdfiumInstance.FPDF_StructElement_GetAttributeAtIndex(&requests.FPDF_StructElement_GetAttributeAtIndex{
							StructElement: td.StructElement,
							Index:         0,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetAttributeAtIndex).ToNot(BeNil())
						Expect(FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute).ToNot(BeEmpty())

						FPDF_StructElement_Attr_GetCount, err := PdfiumInstance.FPDF_StructElement_Attr_GetCount(&requests.FPDF_StructElement_Attr_GetCount{
							StructElementAttribute: FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetCount).To(Equal(&responses.FPDF_StructElement_Attr_GetCount{
							Count: 4,
						}))

						FPDF_StructElement_Attr_GetName, err := PdfiumInstance.FPDF_StructElement_Attr_GetName(&requests.FPDF_StructElement_Attr_GetName{
							StructElementAttribute: FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute,
							Index:                  0,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetName).To(Equal(&responses.FPDF_StructElement_Attr_GetName{
							Name: "ColProp",
						}))

						FPDF_StructElement_Attr_GetValueColProp, err := PdfiumInstance.FPDF_StructElement_Attr_GetValue(&requests.FPDF_StructElement_Attr_GetValue{
							StructElementAttribute: FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute,
							Name:                   "ColProp",
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetValueColProp).To(Not(BeNil()))

						FPDF_StructElement_Attr_GetType, err := PdfiumInstance.FPDF_StructElement_Attr_GetType(&requests.FPDF_StructElement_Attr_GetType{
							StructElementAttributeValue: FPDF_StructElement_Attr_GetValueColProp.StructElementAttributeValue,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetType).To(Equal(&responses.FPDF_StructElement_Attr_GetType{
							ObjectType: enums.FPDF_OBJECT_TYPE_STRING,
						}))

						FPDF_StructElement_Attr_GetStringValue, err := PdfiumInstance.FPDF_StructElement_Attr_GetStringValue(&requests.FPDF_StructElement_Attr_GetStringValue{
							StructElementAttributeValue: FPDF_StructElement_Attr_GetValueColProp.StructElementAttributeValue,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetStringValue).To(Equal(&responses.FPDF_StructElement_Attr_GetStringValue{
							Value: "Sum",
						}))

						FPDF_StructElement_Attr_GetBlobValue, err := PdfiumInstance.FPDF_StructElement_Attr_GetBlobValue(&requests.FPDF_StructElement_Attr_GetBlobValue{
							StructElementAttributeValue: FPDF_StructElement_Attr_GetValueColProp.StructElementAttributeValue,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetBlobValue).To(Equal(&responses.FPDF_StructElement_Attr_GetBlobValue{
							Value: []byte{'S', 'u', 'm'},
						}))

						FPDF_StructElement_Attr_GetName, err = PdfiumInstance.FPDF_StructElement_Attr_GetName(&requests.FPDF_StructElement_Attr_GetName{
							StructElementAttribute: FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute,
							Index:                  1,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetName).To(Equal(&responses.FPDF_StructElement_Attr_GetName{
							Name: "CurUSD",
						}))

						FPDF_StructElement_Attr_GetValueCurUSD, err := PdfiumInstance.FPDF_StructElement_Attr_GetValue(&requests.FPDF_StructElement_Attr_GetValue{
							StructElementAttribute: FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute,
							Name:                   "CurUSD",
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetValueCurUSD).To(Not(BeNil()))

						FPDF_StructElement_Attr_GetType, err = PdfiumInstance.FPDF_StructElement_Attr_GetType(&requests.FPDF_StructElement_Attr_GetType{
							StructElementAttributeValue: FPDF_StructElement_Attr_GetValueCurUSD.StructElementAttributeValue,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetType).To(Equal(&responses.FPDF_StructElement_Attr_GetType{
							ObjectType: enums.FPDF_OBJECT_TYPE_BOOLEAN,
						}))

						FPDF_StructElement_Attr_GetBooleanValue, err := PdfiumInstance.FPDF_StructElement_Attr_GetBooleanValue(&requests.FPDF_StructElement_Attr_GetBooleanValue{
							StructElementAttributeValue: FPDF_StructElement_Attr_GetValueCurUSD.StructElementAttributeValue,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetBooleanValue).To(Equal(&responses.FPDF_StructElement_Attr_GetBooleanValue{
							Value: true,
						}))

						FPDF_StructElement_Attr_GetName, err = PdfiumInstance.FPDF_StructElement_Attr_GetName(&requests.FPDF_StructElement_Attr_GetName{
							StructElementAttribute: FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute,
							Index:                  3,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetName).To(Equal(&responses.FPDF_StructElement_Attr_GetName{
							Name: "RowSpan",
						}))

						FPDF_StructElement_Attr_GetValueRowSpan, err := PdfiumInstance.FPDF_StructElement_Attr_GetValue(&requests.FPDF_StructElement_Attr_GetValue{
							StructElementAttribute: FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute,
							Name:                   "RowSpan",
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetValueRowSpan).To(Not(BeNil()))

						FPDF_StructElement_Attr_GetType, err = PdfiumInstance.FPDF_StructElement_Attr_GetType(&requests.FPDF_StructElement_Attr_GetType{
							StructElementAttributeValue: FPDF_StructElement_Attr_GetValueRowSpan.StructElementAttributeValue,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetType).To(Equal(&responses.FPDF_StructElement_Attr_GetType{
							ObjectType: enums.FPDF_OBJECT_TYPE_NUMBER,
						}))

						FPDF_StructElement_Attr_GetNumberValue, err := PdfiumInstance.FPDF_StructElement_Attr_GetNumberValue(&requests.FPDF_StructElement_Attr_GetNumberValue{
							StructElementAttributeValue: FPDF_StructElement_Attr_GetValueRowSpan.StructElementAttributeValue,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_Attr_GetNumberValue).To(Equal(&responses.FPDF_StructElement_Attr_GetNumberValue{
							Value: 3,
						}))
					})
				})
			})
		})
	})

	Context("a normal PDF file with tagged text", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/tagged_actual_text.pdf")
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

				When("a struct tree struct element is opened", func() {
					var element references.FPDF_STRUCTELEMENT

					BeforeEach(func() {
						FPDF_StructTree_GetChildAtIndex, err := PdfiumInstance.FPDF_StructTree_GetChildAtIndex(&requests.FPDF_StructTree_GetChildAtIndex{
							StructTree: structTree,
							Index:      0,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructTree_GetChildAtIndex).To(Not(BeNil()))

						element = FPDF_StructTree_GetChildAtIndex.StructElement
					})

					It("returns no parent when there is none", func() {
						FPDF_StructElement_GetParent, err := PdfiumInstance.FPDF_StructElement_GetParent(&requests.FPDF_StructElement_GetParent{
							StructElement: element,
						})
						Expect(err).To(MatchError("could not get struct element parent"))
						Expect(FPDF_StructElement_GetParent).To(BeNil())
					})

					It("returns a parent when there is one", func() {
						child_element, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: element,
							Index:         0,
						})
						Expect(err).To(BeNil())
						Expect(child_element).To(Not(BeNil()))

						FPDF_StructElement_GetParent, err := PdfiumInstance.FPDF_StructElement_GetParent(&requests.FPDF_StructElement_GetParent{
							StructElement: child_element.StructElement,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetParent).ToNot(BeNil())
						Expect(FPDF_StructElement_GetParent.StructElement).ToNot(BeEmpty())
					})

					It("returns the struct element object type", func() {
						FPDF_StructElement_GetObjType, err := PdfiumInstance.FPDF_StructElement_GetObjType(&requests.FPDF_StructElement_GetObjType{
							StructElement: element,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetObjType).To(Equal(&responses.FPDF_StructElement_GetObjType{
							ObjType: "StructElem",
						}))
					})

					It("returns the correct text", func() {
						child_element, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: element,
						})
						Expect(err).To(BeNil())
						Expect(child_element).To(Not(BeNil()))

						FPDF_StructElement_GetActualText, err := PdfiumInstance.FPDF_StructElement_GetActualText(&requests.FPDF_StructElement_GetActualText{
							StructElement: child_element.StructElement,
						})
						Expect(err).To(MatchError("could not get actual text"))
						Expect(FPDF_StructElement_GetActualText).To(BeNil())

						gchild_element, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: child_element.StructElement,
						})
						Expect(err).To(BeNil())
						Expect(gchild_element).To(Not(BeNil()))

						FPDF_StructElement_GetActualText, err = PdfiumInstance.FPDF_StructElement_GetActualText(&requests.FPDF_StructElement_GetActualText{
							StructElement: gchild_element.StructElement,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetActualText).To(Equal(&responses.FPDF_StructElement_GetActualText{
							Actualtext: "Actual Text",
						}))
					})

					When("a struct tree struct element child is opened", func() {
						var child_element references.FPDF_STRUCTELEMENT

						BeforeEach(func() {
							FPDF_StructElement_GetChildAtIndex, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
								StructElement: element,
								Index:         0,
							})
							Expect(err).To(BeNil())
							Expect(FPDF_StructElement_GetChildAtIndex).To(Not(BeNil()))

							child_element = FPDF_StructElement_GetChildAtIndex.StructElement
						})

						When("a struct tree struct element child child is opened", func() {
							var gchild_element references.FPDF_STRUCTELEMENT

							BeforeEach(func() {
								FPDF_StructElement_GetChildAtIndex, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
									StructElement: child_element,
									Index:         0,
								})
								Expect(err).To(BeNil())
								Expect(FPDF_StructElement_GetChildAtIndex).To(Not(BeNil()))

								gchild_element = FPDF_StructElement_GetChildAtIndex.StructElement
							})

							It("returns the correct gchild attribute count", func() {
								FPDF_StructElement_GetAttributeCount, err := PdfiumInstance.FPDF_StructElement_GetAttributeCount(&requests.FPDF_StructElement_GetAttributeCount{
									StructElement: gchild_element,
								})
								Expect(err).To(BeNil())
								Expect(FPDF_StructElement_GetAttributeCount).To(Equal(&responses.FPDF_StructElement_GetAttributeCount{
									Count: 1,
								}))
							})

							When("a struct tree struct element child child attribute is opened", func() {
								var attribute references.FPDF_STRUCTELEMENT_ATTR

								BeforeEach(func() {
									FPDF_StructElement_GetAttributeAtIndex, err := PdfiumInstance.FPDF_StructElement_GetAttributeAtIndex(&requests.FPDF_StructElement_GetAttributeAtIndex{
										StructElement: gchild_element,
										Index:         0,
									})
									Expect(err).To(BeNil())
									Expect(FPDF_StructElement_GetAttributeAtIndex).To(Not(BeNil()))

									attribute = FPDF_StructElement_GetAttributeAtIndex.StructElementAttribute
								})

								It("returns the correct attribute count", func() {
									FPDF_StructElement_Attr_GetCount, err := PdfiumInstance.FPDF_StructElement_Attr_GetCount(&requests.FPDF_StructElement_Attr_GetCount{
										StructElementAttribute: attribute,
									})
									Expect(err).To(BeNil())
									Expect(FPDF_StructElement_Attr_GetCount).To(Equal(&responses.FPDF_StructElement_Attr_GetCount{
										Count: 5,
									}))
								})

								It("returns the correct attribute name", func() {
									FPDF_StructElement_Attr_GetName, err := PdfiumInstance.FPDF_StructElement_Attr_GetName(&requests.FPDF_StructElement_Attr_GetName{
										StructElementAttribute: attribute,
										Index:                  1,
									})
									Expect(err).To(BeNil())
									Expect(FPDF_StructElement_Attr_GetName).To(Equal(&responses.FPDF_StructElement_Attr_GetName{
										Name: "Height",
									}))
								})

								It("returns the correct attribute name and value", func() {
									FPDF_StructElement_Attr_GetName, err := PdfiumInstance.FPDF_StructElement_Attr_GetName(&requests.FPDF_StructElement_Attr_GetName{
										StructElementAttribute: attribute,
										Index:                  0,
									})
									Expect(err).To(BeNil())
									Expect(FPDF_StructElement_Attr_GetName).To(Equal(&responses.FPDF_StructElement_Attr_GetName{
										Name: "BBox",
									}))

									FPDF_StructElement_Attr_GetValue, err := PdfiumInstance.FPDF_StructElement_Attr_GetValue(&requests.FPDF_StructElement_Attr_GetValue{
										StructElementAttribute: attribute,
										Name:                   FPDF_StructElement_Attr_GetName.Name,
									})
									Expect(err).To(BeNil())
									Expect(FPDF_StructElement_Attr_GetValue).To(Not(BeNil()))

									FPDF_StructElement_Attr_GetType, err := PdfiumInstance.FPDF_StructElement_Attr_GetType(&requests.FPDF_StructElement_Attr_GetType{
										StructElementAttributeValue: FPDF_StructElement_Attr_GetValue.StructElementAttributeValue,
									})
									Expect(err).To(BeNil())
									Expect(FPDF_StructElement_Attr_GetType).To(Equal(&responses.FPDF_StructElement_Attr_GetType{
										ObjectType: enums.FPDF_OBJECT_TYPE_ARRAY,
									}))

									FPDF_StructElement_Attr_CountChildren, err := PdfiumInstance.FPDF_StructElement_Attr_CountChildren(&requests.FPDF_StructElement_Attr_CountChildren{
										StructElementAttributeValue: FPDF_StructElement_Attr_GetValue.StructElementAttributeValue,
									})
									Expect(err).To(BeNil())
									Expect(FPDF_StructElement_Attr_CountChildren).To(Equal(&responses.FPDF_StructElement_Attr_CountChildren{
										Count: 4,
									}))

									_, err = PdfiumInstance.FPDF_StructElement_Attr_GetChildAtIndex(&requests.FPDF_StructElement_Attr_GetChildAtIndex{
										StructElementAttributeValue: FPDF_StructElement_Attr_GetValue.StructElementAttributeValue,
										Index:                       -1,
									})
									Expect(err).To(MatchError("could not get struct element attribute child"))

									_, err = PdfiumInstance.FPDF_StructElement_Attr_GetChildAtIndex(&requests.FPDF_StructElement_Attr_GetChildAtIndex{
										StructElementAttributeValue: FPDF_StructElement_Attr_GetValue.StructElementAttributeValue,
										Index:                       6,
									})
									Expect(err).To(MatchError("could not get struct element attribute child"))

									nested_attr_value0, err := PdfiumInstance.FPDF_StructElement_Attr_GetChildAtIndex(&requests.FPDF_StructElement_Attr_GetChildAtIndex{
										StructElementAttributeValue: FPDF_StructElement_Attr_GetValue.StructElementAttributeValue,
										Index:                       0,
									})
									Expect(err).To(BeNil())
									Expect(nested_attr_value0).To(Not(BeNil()))

									FPDF_StructElement_Attr_GetType, err = PdfiumInstance.FPDF_StructElement_Attr_GetType(&requests.FPDF_StructElement_Attr_GetType{
										StructElementAttributeValue: nested_attr_value0.StructElementAttributeValue,
									})
									Expect(err).To(BeNil())
									Expect(FPDF_StructElement_Attr_GetType).To(Equal(&responses.FPDF_StructElement_Attr_GetType{
										ObjectType: enums.FPDF_OBJECT_TYPE_NUMBER,
									}))

									nested_attr_value3, err := PdfiumInstance.FPDF_StructElement_Attr_GetChildAtIndex(&requests.FPDF_StructElement_Attr_GetChildAtIndex{
										StructElementAttributeValue: FPDF_StructElement_Attr_GetValue.StructElementAttributeValue,
										Index:                       3,
									})
									Expect(err).To(BeNil())
									Expect(nested_attr_value3).To(Not(BeNil()))

									FPDF_StructElement_Attr_GetType, err = PdfiumInstance.FPDF_StructElement_Attr_GetType(&requests.FPDF_StructElement_Attr_GetType{
										StructElementAttributeValue: nested_attr_value3.StructElementAttributeValue,
									})
									Expect(err).To(BeNil())
									Expect(FPDF_StructElement_Attr_GetType).To(Equal(&responses.FPDF_StructElement_Attr_GetType{
										ObjectType: enums.FPDF_OBJECT_TYPE_NUMBER,
									}))
								})
							})
						})
					})
				})
			})
		})
	})

	Context("a normal PDF file with tagged marked content", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/tagged_marked_content.pdf")
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

				When("a struct tree struct element is opened", func() {
					var child1 references.FPDF_STRUCTELEMENT
					var child2 references.FPDF_STRUCTELEMENT

					BeforeEach(func() {
						FPDF_StructTree_GetChildAtIndex, err := PdfiumInstance.FPDF_StructTree_GetChildAtIndex(&requests.FPDF_StructTree_GetChildAtIndex{
							StructTree: structTree,
							Index:      0,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructTree_GetChildAtIndex).To(Not(BeNil()))

						child1 = FPDF_StructTree_GetChildAtIndex.StructElement

						FPDF_StructTree_GetChildAtIndex, err = PdfiumInstance.FPDF_StructTree_GetChildAtIndex(&requests.FPDF_StructTree_GetChildAtIndex{
							StructTree: structTree,
							Index:      1,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructTree_GetChildAtIndex).To(Not(BeNil()))

						child2 = FPDF_StructTree_GetChildAtIndex.StructElement
					})

					It("returns the correct marked content id count", func() {
						FPDF_StructElement_GetMarkedContentIdCount, err := PdfiumInstance.FPDF_StructElement_GetMarkedContentIdCount(&requests.FPDF_StructElement_GetMarkedContentIdCount{
							StructElement: child1,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetMarkedContentIdCount).To(Equal(&responses.FPDF_StructElement_GetMarkedContentIdCount{
							Count: 1,
						}))
					})

					It("returns the correct marked content id", func() {
						FPDF_StructElement_GetMarkedContentIdAtIndex, err := PdfiumInstance.FPDF_StructElement_GetMarkedContentIdAtIndex(&requests.FPDF_StructElement_GetMarkedContentIdAtIndex{
							StructElement: child1,
							Index:         0,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetMarkedContentIdAtIndex).To(Equal(&responses.FPDF_StructElement_GetMarkedContentIdAtIndex{}))
					})

					It("returns the correct marked content id count", func() {
						FPDF_StructElement_GetMarkedContentIdCount, err := PdfiumInstance.FPDF_StructElement_GetMarkedContentIdCount(&requests.FPDF_StructElement_GetMarkedContentIdCount{
							StructElement: child2,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetMarkedContentIdCount).To(Equal(&responses.FPDF_StructElement_GetMarkedContentIdCount{
							Count: 1,
						}))
					})

					It("returns the correct marked content id", func() {
						FPDF_StructElement_GetMarkedContentIdAtIndex, err := PdfiumInstance.FPDF_StructElement_GetMarkedContentIdAtIndex(&requests.FPDF_StructElement_GetMarkedContentIdAtIndex{
							StructElement: child2,
							Index:         0,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetMarkedContentIdAtIndex).To(Equal(&responses.FPDF_StructElement_GetMarkedContentIdAtIndex{
							MarkedContentID: 1,
						}))
					})
				})
			})
		})
	})

	Context("a normal PDF file with tagged marked content multipage", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/tagged_mcr_multipage.pdf")
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

		pageIndexes := []int{0, 1}
		for i := range pageIndexes {
			page_i := i
			When("is opened", func() {
				When("a page structtree is opened", func() {
					var structTree references.FPDF_STRUCTTREE

					BeforeEach(func() {
						FPDF_StructTree_GetForPage, err := PdfiumInstance.FPDF_StructTree_GetForPage(&requests.FPDF_StructTree_GetForPage{
							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    page_i,
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
						var struct_doc references.FPDF_STRUCTELEMENT

						BeforeEach(func() {
							FPDF_StructTree_GetChildAtIndex, err := PdfiumInstance.FPDF_StructTree_GetChildAtIndex(&requests.FPDF_StructTree_GetChildAtIndex{
								StructTree: structTree,
								Index:      0,
							})
							Expect(err).To(BeNil())
							Expect(FPDF_StructTree_GetChildAtIndex).To(Not(BeNil()))

							struct_doc = FPDF_StructTree_GetChildAtIndex.StructElement
						})

						It("gives an error when loading children", func() {
							FPDF_StructElement_GetChildAtIndex, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
								StructElement: struct_doc,
								Index:         0,
							})
							Expect(err).To(MatchError("could not load struct element child"))
							Expect(FPDF_StructElement_GetChildAtIndex).To(BeNil())
						})

						It("returns the correct marked content id for the struct doc", func() {
							FPDF_StructElement_GetChildAtIndex, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
								StructElement: struct_doc,
								Index:         1,
							})
							Expect(err).To(MatchError("could not load struct element child"))
							Expect(FPDF_StructElement_GetChildAtIndex).To(BeNil())
						})

						It("returns the correct children count for the struct doc", func() {
							FPDF_StructElement_CountChildren, err := PdfiumInstance.FPDF_StructElement_CountChildren(&requests.FPDF_StructElement_CountChildren{
								StructElement: struct_doc,
							})
							Expect(err).To(BeNil())
							Expect(FPDF_StructElement_CountChildren).To(Equal(&responses.FPDF_StructElement_CountChildren{
								Count: 2,
							}))
						})

						It("returns the correct marked contetn id count for the struct doc", func() {
							FPDF_StructElement_GetMarkedContentIdCount, err := PdfiumInstance.FPDF_StructElement_GetMarkedContentIdCount(&requests.FPDF_StructElement_GetMarkedContentIdCount{
								StructElement: struct_doc,
							})
							Expect(err).To(BeNil())
							Expect(FPDF_StructElement_GetMarkedContentIdCount).To(Equal(&responses.FPDF_StructElement_GetMarkedContentIdCount{
								Count: 2,
							}))
						})

						// Both MCID are returned as if part of this page, while they are not.
						// So `FPDF_StructElement_GetMarkedContentIdAtIndex(...)` does not work
						// for StructElement spanning multiple pages.

						It("returns the correct marked content id", func() {
							FPDF_StructElement_GetMarkedContentIdAtIndex, err := PdfiumInstance.FPDF_StructElement_GetMarkedContentIdAtIndex(&requests.FPDF_StructElement_GetMarkedContentIdAtIndex{
								StructElement: struct_doc,
								Index:         0,
							})
							Expect(err).To(BeNil())
							Expect(FPDF_StructElement_GetMarkedContentIdAtIndex).To(Equal(&responses.FPDF_StructElement_GetMarkedContentIdAtIndex{}))
						})

						It("returns the correct marked content id", func() {
							FPDF_StructElement_GetMarkedContentIdAtIndex, err := PdfiumInstance.FPDF_StructElement_GetMarkedContentIdAtIndex(&requests.FPDF_StructElement_GetMarkedContentIdAtIndex{
								StructElement: struct_doc,
								Index:         1,
							})
							Expect(err).To(BeNil())
							Expect(FPDF_StructElement_GetMarkedContentIdAtIndex).To(Equal(&responses.FPDF_StructElement_GetMarkedContentIdAtIndex{}))
						})

						// One MCR is pointing to page 1, another to page2, so those are different
						// for different pages.

						It("returns the correct child marked content id", func() {
							FPDF_StructElement_GetChildMarkedContentID, err := PdfiumInstance.FPDF_StructElement_GetChildMarkedContentID(&requests.FPDF_StructElement_GetChildMarkedContentID{
								StructElement: struct_doc,
								Index:         0,
							})

							if page_i == 0 {
								Expect(err).To(BeNil())
								Expect(FPDF_StructElement_GetChildMarkedContentID).To(Equal(&responses.FPDF_StructElement_GetChildMarkedContentID{
									ChildMarkedContentID: 0,
								}))
							} else {
								Expect(err).To(MatchError("could not get struct element child marked content id"))
								Expect(FPDF_StructElement_GetChildMarkedContentID).To(BeNil())
							}
						})

						It("returns the correct child marked content id", func() {
							FPDF_StructElement_GetChildMarkedContentID, err := PdfiumInstance.FPDF_StructElement_GetChildMarkedContentID(&requests.FPDF_StructElement_GetChildMarkedContentID{
								StructElement: struct_doc,
								Index:         1,
							})

							if page_i == 1 {
								Expect(err).To(BeNil())
								Expect(FPDF_StructElement_GetChildMarkedContentID).To(Equal(&responses.FPDF_StructElement_GetChildMarkedContentID{
									ChildMarkedContentID: 0,
								}))
							} else {
								Expect(err).To(MatchError("could not get struct element child marked content id"))
								Expect(FPDF_StructElement_GetChildMarkedContentID).To(BeNil())
							}
						})

						It("returns an error when giving an invalid child index", func() {
							FPDF_StructElement_GetChildMarkedContentID, err := PdfiumInstance.FPDF_StructElement_GetChildMarkedContentID(&requests.FPDF_StructElement_GetChildMarkedContentID{
								StructElement: struct_doc,
								Index:         -1,
							})
							Expect(err).To(MatchError("could not get struct element child marked content id"))
							Expect(FPDF_StructElement_GetChildMarkedContentID).To(BeNil())

							FPDF_StructElement_GetChildMarkedContentID, err = PdfiumInstance.FPDF_StructElement_GetChildMarkedContentID(&requests.FPDF_StructElement_GetChildMarkedContentID{
								StructElement: struct_doc,
								Index:         2,
							})
							Expect(err).To(MatchError("could not get struct element child marked content id"))
							Expect(FPDF_StructElement_GetChildMarkedContentID).To(BeNil())
						})

						It("returns an error when giving an invalid struct element", func() {
							FPDF_StructElement_GetChildMarkedContentID, err := PdfiumInstance.FPDF_StructElement_GetChildMarkedContentID(&requests.FPDF_StructElement_GetChildMarkedContentID{
								Index: 0,
							})
							Expect(err).To(MatchError("structElement not given"))
							Expect(FPDF_StructElement_GetChildMarkedContentID).To(BeNil())
						})
					})
				})
			})
		}
	})
})
