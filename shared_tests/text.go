package shared_tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/klippa-app/go-pdfium/errors"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("text", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	Context("no references", func() {
		When("is given", func() {
			Context("GetPageText()", func() {
				It("returns an error", func() {
					pageText, err := PdfiumInstance.GetPageText(&requests.GetPageText{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Index: 0,
							},
						},
					})
					Expect(err).To(MatchError("document not given"))
					Expect(pageText).To(BeNil())
				})
			})

			Context("GetPageTextStructured()", func() {
				It("returns an error", func() {
					pageTextStructured, err := PdfiumInstance.GetPageTextStructured(&requests.GetPageTextStructured{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Index: 0,
							},
						},
					})
					Expect(err).To(MatchError("document not given"))
					Expect(pageTextStructured).To(BeNil())
				})
			})
		})
	})

	Context("a normal PDF file", func() {
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
			Context("when an invalid page is given", func() {
				Context("GetPageText()", func() {
					It("returns an error", func() {
						pageText, err := PdfiumInstance.GetPageText(&requests.GetPageText{
							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    1,
								},
							},
						})
						Expect(err).To(MatchError(errors.ErrPage.Error()))
						Expect(pageText).To(BeNil())
					})
				})

				Context("GetPageTextStructured()", func() {
					It("returns an error", func() {
						pageTextStructured, err := PdfiumInstance.GetPageTextStructured(&requests.GetPageTextStructured{
							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    1,
								},
							},
						})
						Expect(err).To(MatchError(errors.ErrPage.Error()))
						Expect(pageTextStructured).To(BeNil())
					})
				})
			})

			Context("when the page text is requested", func() {
				It("returns the correct text", func() {
					pageText, err := PdfiumInstance.GetPageText(&requests.GetPageText{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})
					Expect(err).To(BeNil())
					Expect(pageText).To(Equal(&responses.GetPageText{
						Text: "File: Untitled Document 2 Page 1 of 1\r\nThis is a test PDF",
					}))
				})
			})

			Context("when the structured page text is requested", func() {
				It("returns the correct structured text", func() {
					pageTextStructured, err := PdfiumInstance.GetPageTextStructured(&requests.GetPageTextStructured{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})

					Expect(err).To(BeNil())
					Expect(pageTextStructured).To(Or(loadStructuredText(pageTextStructured, TestDataPath+"/testdata/text_"+TestType+"_testpdf_without_pixel_calculations.json", TestDataPath+"/testdata/text_"+TestType+"_testpdf_without_pixel_calculations_7019.json")...))
				})

				Context("when PixelPositions is enabled", func() {
					Context("with no DPI and no pixels", func() {
						It("returns an error", func() {
							pageTextStructured, err := PdfiumInstance.GetPageTextStructured(&requests.GetPageTextStructured{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
								PixelPositions: requests.GetPageTextStructuredPixelPositions{
									Calculate: true,
								},
							})
							Expect(err).To(MatchError("no DPI or resolution given to calculate pixel positions"))
							Expect(pageTextStructured).To(BeNil())
						})
					})

					Context("with DPI", func() {
						It("returns the correct calculations", func() {
							pageTextStructured, err := PdfiumInstance.GetPageTextStructured(&requests.GetPageTextStructured{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
								PixelPositions: requests.GetPageTextStructuredPixelPositions{
									Calculate: true,
									DPI:       300,
								},
							})
							Expect(err).To(BeNil())
							Expect(pageTextStructured).To(Or(loadStructuredText(pageTextStructured, TestDataPath+"/testdata/text_"+TestType+"_testpdf_with_dpi_pixel_calculations.json", TestDataPath+"/testdata/text_"+TestType+"_testpdf_with_dpi_pixel_calculations_7019.json")...))
						})
					})

					Context("with pixels", func() {
						It("returns the correct calculations", func() {
							pageTextStructured, err := PdfiumInstance.GetPageTextStructured(&requests.GetPageTextStructured{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
								PixelPositions: requests.GetPageTextStructuredPixelPositions{
									Calculate: true,
									Width:     3000,
									Height:    3000,
								},
							})

							Expect(err).To(BeNil())
							Expect(pageTextStructured).To(Or(loadStructuredText(pageTextStructured, TestDataPath+"/testdata/text_"+TestType+"_testpdf_with_resolution_pixel_calculations.json", TestDataPath+"/testdata/text_"+TestType+"_testpdf_with_resolution_pixel_calculations_7019.json")...))
						})
					})
				})
			})
		})
	})

	Context("a PDF file with multibyte characters", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/rect-wrong.pdf")
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
			Context("when the page text is requested", func() {
				It("returns the correct text", func() {
					pageText, err := PdfiumInstance.GetPageText(&requests.GetPageText{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})
					Expect(err).To(BeNil())
					Expect(pageText).To(Equal(&responses.GetPageText{
						Text: "arXiv:2501.00201v1 [eess.SP] 31 Dec 2024\r\n1\r\nHierarchical Functionality Prioritization in Multicast ISAC:\r\nOptimal Admission Control and Discrete-Phase Beamforming\r\nLuis F. Abanto-Leon and Setareh Maghsudi\r\nAbstract—We investigate the joint admission control and\r\ndiscrete-phase multicast beamforming design for integrated sens\ufffeing and commmunications (ISAC) systems, where sensing and\r\ncommunications functionalities have different hierarchies. Specif\ufffeically, the ISAC system first allocates resources to the higher\ufffehierarchy functionality and opportunistically uses the remaining\r\nresources to support the lower-hierarchy one. This resource allo\ufffecation problem is a nonconvex mixed-integer nonlinear program\r\n(MINLP). We propose an exact mixed-integer linear program\r\n(MILP) reformulation, leading to a globally optimal solution. In\r\naddition, we implemented three baselines for comparison, which\r\nour proposed method outperforms by more than 39%\r\n.\r\nIndex Terms—Integrated sensing and communications, multi\ufffecast, beamforming, discrete phases, admission control.\r\nI. INTRODUCTION\r\nIntegrated sensing and commmunications (ISAC) is a dis\uffferuptive advancement in wireless technology in which sensing\r\nand communications share the same radio resources, e.g.,\r\ninfrastructure, spectrum, waveform, to enhance radio resource\r\nutilization, reduce costs, and simplify system complexity [1].\r\nSensing at high frequencies is appealing since the shorter\r\nwavelengths enable finer resolution [2]. These frequencies suf\ufffefer severe path loss, which beamforming can alleviate. Highly\r\nversatile digital beamformers are expensive to manufacture for\r\nsuch high frequencies. Hence, analog beamformers lead the\r\ninitial stages of ISAC systems operating at these frequencies.\r\nAnalog beamformers can be designed with continuous or\r\ndiscrete phases. The state-of-the-art literature features beam\ufffeforming designs with both phase types, but most works\r\nfocused on continuous phases, e.g., [3]–[5], while a few\r\naccounted for discrete phases, e.g., [6]. The latter are of\r\nimmense practical interest, as they reduce system complexity\r\nand costs. To date, however, only suboptimal beamforming\r\ndesigns exist for ISAC systems that utilize discrete phases\r\n.\r\nAnother characteristic of analog beamformers is their single\r\nradio-frequency (RF) chain, which supports one signal stream,\r\nmaking them well-suited for multicasting scenarios, such a\r\ns\r\nbroadcasting live sports or concerts to several subscribed\r\nusers simultaneously. Multicast beamforming has been well\r\ninvestigated in non-ISAC systems, e.g., [7], [8], but rarely in\r\nISAC systems, with only a few studies addressing the topic,\r\ne.g., [9]. Yet, none of such studies accounted for constant\ufffemodulus discrete phases. Particularly, multicasting and ISAC\r\ncould play a key role in live events where drones are often use\r\nd\r\nfor aerial filming. Thus, ISAC could enable drone tracking\r\nwhile supporting efficient content dissemination to users.\r\nIn non-ISAC systems, admission control is crucial in pre\ufffeventing resource allocation infeasibility, especially when radio\r\nresources are limited, allowing to serve only a selected subset\r\nFig. 1: Multicast ISAC system with many users and a target.\r\nof users [8], [10], thereby enhancing resource utilization. In\r\nlight of its advantages, incorporating admission control into\r\nISAC systems holds significant promise. However, despite it\r\ns\r\npotential, this aspect has been overlooked in ISAC contexts\r\n.\r\nMoreover, angular positions of targets may not be known\r\nprecisely due to factors such as motion. Thus, accounting fo\r\nr\r\nthis aspect in the resource allocation design can help mitigate\r\npotential performance degradation in sensing, a crucial aspect\r\nexplored in only a few studies, such as [11].\r\nIn ISAC systems, one functionality may be more critical\r\nthan the other [12]. Particularly, this view aligns with indus\ufffetry’s pragmatic stance of preserving communication perfor\ufffemance, while enabling sensing opportunistically when feasi\ufffeble. While tradeoff functions can balance the importance of\r\nfunctionalities by using weights [13], changes in paramete\r\nr\r\nsettings (e.g., number of users, transmit power) can skew\r\nobjective function values, rendering preset weights ineffective\r\nand shifting the intended operating point. To address this, we\r\npropose establishing strict hierarchies through careful weight\r\ndesign. Our approach consistently prioritizes communications\r\nregardless of parameter settings, ensuring its full optimization\r\nbefore addressing the sensing requirements, thus leading to a\r\nstrictly tiered resource allocation framework.\r\nMotivated by the above discussion, we investigate the joint\r\noptimization of admission control and multicast beamforming\r\nwith discrete phases for ISAC systems, prioritizing communi\ufffecations while enabling opportunistic sensing, and accounting\r\nfor target angular uncertainty. This novel resource alloca\r\n-\r\ntion problem, distinct from existing works (see Table I), is\r\nformulated as a nonconvex mixed-integer nonlinear program\r\n(MINLP), which is challenging to solve. We propose an\r\napproach to reformulate it, leading to a mixed-integer linear\r\nprogram (MILP) that can be solved globally optimally. Our\r\napproach employs a series of transformations to convexify the\r\nnonconvex MINLP without compromising optimality, effec\ufffetively addressing the original problem’s complexity. Addition\ufffeally, we implement three baselines based on well-established\r\noptimization methods used in the resource allocation literature.\r\nNotation: Boldface capital letters\r\nA and boldface lowercase\r\nletters\r\na denote matrices and vectors, respectively. The trans-",
					}))
				})
			})

			Context("when the structured page text is requested", func() {
				It("returns the correct structured text", func() {
					pageTextStructured, err := PdfiumInstance.GetPageTextStructured(&requests.GetPageTextStructured{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    0,
							},
						},
					})

					Expect(err).To(BeNil())
					Expect(pageTextStructured).To(Or(loadStructuredText(pageTextStructured, TestDataPath+"/testdata/multibyte_text_"+TestType+"_testpdf_without_pixel_calculations.json", TestDataPath+"/testdata/multibyte_text_"+TestType+"_testpdf_without_pixel_calculations_7019.json")...))
				})

				Context("when PixelPositions is enabled", func() {
					Context("with no DPI and no pixels", func() {
						It("returns an error", func() {
							pageTextStructured, err := PdfiumInstance.GetPageTextStructured(&requests.GetPageTextStructured{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
								PixelPositions: requests.GetPageTextStructuredPixelPositions{
									Calculate: true,
								},
							})
							Expect(err).To(MatchError("no DPI or resolution given to calculate pixel positions"))
							Expect(pageTextStructured).To(BeNil())
						})
					})

					Context("with DPI", func() {
						It("returns the correct calculations", func() {
							pageTextStructured, err := PdfiumInstance.GetPageTextStructured(&requests.GetPageTextStructured{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
								PixelPositions: requests.GetPageTextStructuredPixelPositions{
									Calculate: true,
									DPI:       300,
								},
							})
							Expect(err).To(BeNil())
							Expect(pageTextStructured).To(Or(loadStructuredText(pageTextStructured, TestDataPath+"/testdata/multibyte_text_"+TestType+"_testpdf_with_dpi_pixel_calculations.json", TestDataPath+"/testdata/multibyte_text_"+TestType+"_testpdf_with_dpi_pixel_calculations_7019.json")...))
						})
					})

					Context("with pixels", func() {
						It("returns the correct calculations", func() {
							pageTextStructured, err := PdfiumInstance.GetPageTextStructured(&requests.GetPageTextStructured{
								Page: requests.Page{
									ByIndex: &requests.PageByIndex{
										Document: doc,
										Index:    0,
									},
								},
								PixelPositions: requests.GetPageTextStructuredPixelPositions{
									Calculate: true,
									Width:     3000,
									Height:    3000,
								},
							})

							Expect(err).To(BeNil())
							Expect(pageTextStructured).To(Or(loadStructuredText(pageTextStructured, TestDataPath+"/testdata/multibyte_text_"+TestType+"_testpdf_with_resolution_pixel_calculations.json", TestDataPath+"/testdata/multibyte_text_"+TestType+"_testpdf_with_resolution_pixel_calculations_7019.json")...))
						})
					})
				})
			})
		})
	})
})

func loadStructuredText(resp *responses.GetPageTextStructured, paths ...string) []types.GomegaMatcher {
	result := []types.GomegaMatcher{}

	for _, path := range paths {
		writeStructuredText(path, resp)
		preRender, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}

		buf := bytes.NewBuffer(preRender)
		dec := json.NewDecoder(buf)

		var text responses.GetPageTextStructured
		err = dec.Decode(&text)

		result = append(result, Equal(&text))
	}

	return result
}

func writeStructuredText(path string, resp *responses.GetPageTextStructured) error {
	if _, err := os.Stat(path); err == nil {
		return nil // Comment this in case of updating PDFium versions and rendering has changed.
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(resp); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, buf.Bytes(), 0777); err != nil {
		return err
	}

	return nil
}
