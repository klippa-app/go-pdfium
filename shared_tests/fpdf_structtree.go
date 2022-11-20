package shared_tests

import (
	"github.com/klippa-app/go-pdfium/responses"
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fpdf_structtree", func() {
	BeforeEach(func() {
		Locker.Lock()

		if TestType == "webassembly" {
			// @todo: remove me when implemented.
			Skip("This test is skipped on Webassembly")
		}
	})

	AfterEach(func() {
		Locker.Unlock()

		if TestType == "webassembly" {
			// @todo: remove me when implemented.
			Skip("This test is skipped on Webassembly")
		}
	})

	Context("no page", func() {
		When("is opened", func() {
			It("returns an error when calling FPDF_StructTree_GetForPage", func() {
				FPDF_StructTree_GetForPage, err := PdfiumInstance.FPDF_StructTree_GetForPage(&requests.FPDF_StructTree_GetForPage{})
				Expect(err).To(MatchError("either page reference or index should be given"))
				Expect(FPDF_StructTree_GetForPage).To(BeNil())
			})
		})
	})

	Context("no struct tree", func() {
		When("is opened", func() {
			It("returns an error when calling FPDF_StructTree_Close", func() {
				FPDF_StructTree_Close, err := PdfiumInstance.FPDF_StructTree_Close(&requests.FPDF_StructTree_Close{})
				Expect(err).To(MatchError("structTree not given"))
				Expect(FPDF_StructTree_Close).To(BeNil())
			})

			It("returns an error when calling FPDF_StructTree_CountChildren", func() {
				FPDF_StructTree_CountChildren, err := PdfiumInstance.FPDF_StructTree_CountChildren(&requests.FPDF_StructTree_CountChildren{})
				Expect(err).To(MatchError("structTree not given"))
				Expect(FPDF_StructTree_CountChildren).To(BeNil())
			})

			It("returns an error when calling FPDF_StructTree_GetChildAtIndex", func() {
				FPDF_StructTree_GetChildAtIndex, err := PdfiumInstance.FPDF_StructTree_GetChildAtIndex(&requests.FPDF_StructTree_GetChildAtIndex{})
				Expect(err).To(MatchError("structTree not given"))
				Expect(FPDF_StructTree_GetChildAtIndex).To(BeNil())
			})
		})
	})

	Context("no struct element", func() {
		When("is opened", func() {
			It("returns an error when calling FPDF_StructElement_GetAltText", func() {
				FPDF_StructElement_GetAltText, err := PdfiumInstance.FPDF_StructElement_GetAltText(&requests.FPDF_StructElement_GetAltText{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetAltText).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_GetMarkedContentID", func() {
				FPDF_StructElement_GetMarkedContentID, err := PdfiumInstance.FPDF_StructElement_GetMarkedContentID(&requests.FPDF_StructElement_GetMarkedContentID{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetMarkedContentID).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_GetType", func() {
				FPDF_StructElement_GetType, err := PdfiumInstance.FPDF_StructElement_GetType(&requests.FPDF_StructElement_GetType{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetType).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_GetTitle", func() {
				FPDF_StructElement_GetTitle, err := PdfiumInstance.FPDF_StructElement_GetTitle(&requests.FPDF_StructElement_GetTitle{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetTitle).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_CountChildren", func() {
				FPDF_StructElement_CountChildren, err := PdfiumInstance.FPDF_StructElement_CountChildren(&requests.FPDF_StructElement_CountChildren{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_CountChildren).To(BeNil())
			})

			It("returns an error when calling FPDF_StructElement_GetChildAtIndex", func() {
				FPDF_StructElement_GetChildAtIndex, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{})
				Expect(err).To(MatchError("structElement not given"))
				Expect(FPDF_StructElement_GetChildAtIndex).To(BeNil())
			})
		})
	})

	Context("a normal PDF file without a struct tree", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/test.pdf")
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
			It("returns an error when calling FPDF_StructTree_GetForPage", func() {
				FPDF_StructTree_GetForPage, err := PdfiumInstance.FPDF_StructTree_GetForPage(&requests.FPDF_StructTree_GetForPage{
					Page: requests.Page{
						ByIndex: &requests.PageByIndex{
							Document: doc,
							Index:    0,
						},
					},
				})
				Expect(err).To(MatchError("could not load struct tree"))
				Expect(FPDF_StructTree_GetForPage).To(BeNil())
			})
		})
	})

	Context("a normal PDF file with a struct tree", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/tagged_alt_text.pdf")
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

				It("returns an error when loading an invalid struct element", func() {
					FPDF_StructTree_GetChildAtIndex, err := PdfiumInstance.FPDF_StructTree_GetChildAtIndex(&requests.FPDF_StructTree_GetChildAtIndex{
						StructTree: structTree,
						Index:      30,
					})
					Expect(err).To(MatchError("could not load struct tree child"))
					Expect(FPDF_StructTree_GetChildAtIndex).To(BeNil())
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

					It("returns the correct marked content ID", func() {
						FPDF_StructElement_GetMarkedContentID, err := PdfiumInstance.FPDF_StructElement_GetMarkedContentID(&requests.FPDF_StructElement_GetMarkedContentID{
							StructElement: structElement,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetMarkedContentID).To(Equal(&responses.FPDF_StructElement_GetMarkedContentID{
							MarkedContentID: -1,
						}))
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

					It("returns an error when loading an invalid subelement", func() {
						FPDF_StructElement_GetChildAtIndex, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
							StructElement: structElement,
							Index:         30,
						})
						Expect(err).To(MatchError("could not load struct element child"))
						Expect(FPDF_StructElement_GetChildAtIndex).To(BeNil())
					})

					When("a struct tree struct subelement is opened", func() {
						var structSubElement references.FPDF_STRUCTELEMENT

						BeforeEach(func() {
							FPDF_StructElement_GetChildAtIndex, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
								StructElement: structElement,
							})
							Expect(err).To(BeNil())
							Expect(FPDF_StructElement_GetChildAtIndex).To(Not(BeNil()))

							structSubElement = FPDF_StructElement_GetChildAtIndex.StructElement
						})

						It("returns the correct struct element children count", func() {
							FPDF_StructElement_CountChildren, err := PdfiumInstance.FPDF_StructElement_CountChildren(&requests.FPDF_StructElement_CountChildren{
								StructElement: structSubElement,
							})
							Expect(err).To(BeNil())
							Expect(FPDF_StructElement_CountChildren).To(Equal(&responses.FPDF_StructElement_CountChildren{
								Count: 1,
							}))
						})

						It("returns the correct struct element children count", func() {
							FPDF_StructElement_CountChildren, err := PdfiumInstance.FPDF_StructElement_CountChildren(&requests.FPDF_StructElement_CountChildren{
								StructElement: structSubElement,
							})
							Expect(err).To(BeNil())
							Expect(FPDF_StructElement_CountChildren).To(Equal(&responses.FPDF_StructElement_CountChildren{
								Count: 1,
							}))
						})

						It("returns an error when the struct element has no alt text", func() {
							FPDF_StructElement_GetAltText, err := PdfiumInstance.FPDF_StructElement_GetAltText(&requests.FPDF_StructElement_GetAltText{
								StructElement: structSubElement,
							})
							Expect(err).To(MatchError("Could not get alt text"))
							Expect(FPDF_StructElement_GetAltText).To(BeNil())
						})

						When("a struct tree struct subsubelement is opened", func() {
							var structSubSubElement references.FPDF_STRUCTELEMENT

							BeforeEach(func() {
								FPDF_StructElement_GetChildAtIndex, err := PdfiumInstance.FPDF_StructElement_GetChildAtIndex(&requests.FPDF_StructElement_GetChildAtIndex{
									StructElement: structSubElement,
								})
								Expect(err).To(BeNil())
								Expect(FPDF_StructElement_GetChildAtIndex).To(Not(BeNil()))

								structSubSubElement = FPDF_StructElement_GetChildAtIndex.StructElement
							})

							It("returns the correct alt text", func() {
								FPDF_StructElement_GetAltText, err := PdfiumInstance.FPDF_StructElement_GetAltText(&requests.FPDF_StructElement_GetAltText{
									StructElement: structSubSubElement,
								})
								Expect(err).To(BeNil())
								Expect(FPDF_StructElement_GetAltText).To(Equal(&responses.FPDF_StructElement_GetAltText{
									AltText: "Black Image",
								}))
							})

							It("returns an error when requesting element title", func() {
								FPDF_StructElement_GetTitle, err := PdfiumInstance.FPDF_StructElement_GetTitle(&requests.FPDF_StructElement_GetTitle{
									StructElement: structSubSubElement,
								})
								Expect(err).To(MatchError("Could not get title"))
								Expect(FPDF_StructElement_GetTitle).To(BeNil())
							})

							It("returns an error when requesting element type", func() {
								FPDF_StructElement_GetType, err := PdfiumInstance.FPDF_StructElement_GetType(&requests.FPDF_StructElement_GetType{
									StructElement: structSubSubElement,
								})
								Expect(err).To(BeNil())
								Expect(FPDF_StructElement_GetType).To(Equal(&responses.FPDF_StructElement_GetType{
									Type: "Figure",
								}))
							})
						})
					})
				})
			})
		})
	})

	Context("a normal PDF file with a marked content ID", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/marked_content_id.pdf")
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

					It("returns the correct marked content ID", func() {
						FPDF_StructElement_GetMarkedContentID, err := PdfiumInstance.FPDF_StructElement_GetMarkedContentID(&requests.FPDF_StructElement_GetMarkedContentID{
							StructElement: structElement,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetMarkedContentID).To(Equal(&responses.FPDF_StructElement_GetMarkedContentID{
							MarkedContentID: 0,
						}))
					})
				})
			})
		})
	})

	Context("a normal PDF file with a tagged alt text", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/tagged_alt_text.pdf")
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

					It("returns the correct element title", func() {
						FPDF_StructElement_GetTitle, err := PdfiumInstance.FPDF_StructElement_GetTitle(&requests.FPDF_StructElement_GetTitle{
							StructElement: structElement,
						})
						Expect(err).To(BeNil())
						Expect(FPDF_StructElement_GetTitle).To(Equal(&responses.FPDF_StructElement_GetTitle{
							Title: "TitleText",
						}))
					})
				})
			})
		})
	})
})
