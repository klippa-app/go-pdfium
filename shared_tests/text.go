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
			/*
				Context("when the page text is requested", func() {
					It("returns the correct text", func() {
						pageText, err := PdfiumInstance.GetPageText(&requests.GetPageText{
							Page: requests.Page{
								ByIndex: &requests.PageByIndex{
									Document: doc,
									Index:    1,
								},
							},
						})
						Expect(err).To(BeNil())
						Expect(pageText).To(Equal(&responses.GetPageText{
							Text: "2\r\nTABLE I: Categorization of related work.\r\nWorks D1 D2 D3 D4 D5 D6 D7\r\n[3]–[5] ISAC ✗ Continuous Unicast ✗ ✗ ✗\r\n[6] ISAC ✗ Discrete Unicast ✗ ✗ ✗\r\n[7] Non-ISAC ✓ Discrete Multicast ✗ ✗ ✗\r\n[9] ISAC ✗ — Multicast ✗ ✗ ✗\r\n[10] Non-ISAC ✗ Continuous Multicast ✓ ✗ ✗\r\n[11] ISAC ✗ Continuous Unicast ✗ ✓ ✗\r\nProposed ISAC ✓ Discrete Multicast ✔ ✔ ✔\r\nD1: System type D2: Globally optimality D3: Phase type D4: Network topology\r\nD5: Admission control D6: Angle uncertainty D7: Hierarchical prioritization\r\npose, Hermitian transpose, and trace of A are denoted by AT,\r\nAH, and Tr (A), respectively. The l-th row and i-th column of\r\nA are denoted by [A]\r\nl,:\r\nand [A]\r\n:,i, respectively, and the l-th\r\nelement of a is denoted by [a]\r\nl\r\n. C\r\nI×J\r\nand N denote the space\r\nof I × J complex-valued matrices and the natural numbers,\r\nrespectively. Also, j ,\r\n√\r\n−1 is the imaginary unit, E{·}\r\ndenotes statistical expectation, and CN \r\nυ, ξ2\r\n\x01\r\nrepresents the\r\ncomplex Gaussian distribution with mean υ and variance ξ\r\n2\r\n.\r\nII. SYSTEM MODEL AND PROBLEM FORMULATION\r\nWe consider an ISAC system comprising a base station (BS)\r\nequipped with N transmit and N receive antennas, U single\ufffeantenna users, and one target, as shown in Fig. 1.\r\nBeamforming: The BS transmits signal d = wz, where\r\nw ∈ C\r\nN×1\r\nis the multicast beamforming vector and z ∈ C\r\nis the data symbol which serves both sensing and com\ufffemunication purposes simultaneously, and follows a complex\r\nGaussian distribution with zero mean and unit variance, e.g.,\r\nE{zz∗} = 1. To account for the constant-modulus discrete\r\nphases used in the analog beamforming design, we include\r\nconstraint C1 : [w]n ∈ S, ∀n ∈ N , where N = {1, . . . , N}\r\nindexes the antenna elements and S =\r\n\b\r\nδejφ1\r\n, . . . , δejφL\r\n\t\r\nis\r\nthe set of admissible phases. In addition, δ =\r\np\r\nPtx/N is the\r\nmagnitude, L is the number of phases, φl\r\nis the l-th phase, and\r\nPtx is the BS’s transmit power. Furthermore, Q is the number\r\nof bits needed for encoding the L phases, i.e., Q = log2\r\n(L).\r\nAdmission control: To decide which users are served by the\r\nBS, we include constraint C2 : µu ∈ {0, 1} , ∀u ∈ U, where\r\nU = {1, . . . , U} indexes the users. Here, µu = 1 indicates\r\nthat user u is admitted, and µu = 0 otherwise.\r\nCommunications model: The signal received by user u\r\nis ycom,u = h\r\nH\r\nu d + ηcom,u = h\r\nH\r\nu wz + ηcom,u, where hu ∈\r\nC\r\nN×1\r\nis the channel between the BS and user u, and ηcom,u ∼\r\nCN \r\n0, σ2\r\ncom\x01\r\nis additive white Gaussian noise (AWGN). The\r\ncommunication signal-to-noise ratio (SNR) at user u is\r\nSNRcom,u (w) = wHHe uw, ∀u ∈ U, (1)\r\nwhere He u =\r\nhuh\r\nH\r\nu\r\nσ2\r\ncom\r\n. Let Γth be the minimum SNR threshold\r\nnecessary for successfully decoding the multicast data. To\r\nenforce this requirement jointly with user admission, we\r\nincorporate constraint C3 : wHHe uw ≥ µu · Γth, ∀u ∈ U, i.e.,\r\nthe SNR threshold must be satisfied for all admitted users.\r\nSensing model: We assume the target is far from the\r\nBS, thus we model it as a single point. The BS operates\r\nas a monostatic co-located radar, i.e., the angle of departure\r\n(AoD) and angle of arrival (AoA) are the same. Hence,\r\nthe response matrix between the BS and the target is given\r\nby G (θ) = αa (θ) a\r\nH (θ), where α is the reflection co\ufffeefficient, θ is the AoD/AoA of the target, and a (θ) =\r\nh\r\ne\r\njπ −N+1\r\n2\r\ncos(θ)\r\n, . . . , ejπ N−1\r\n2\r\ncos(θ)\r\niT\r\n, ∈ C\r\nN×1\r\nis the half\ufffewavelength steering vector in the direction of θ. The reflected\r\nsignal by the target at the BS is ysen = wHG (θ) d + ηsen =\r\nwHG (θ) wz + ηsen, where ηsen ∼ CN \r\n0, σ2\r\nsen\x01\r\n, Thus, the\r\nsensing SNR, measured at he BS, is given by\r\nSNRsen (w, θ) = wHGe (θ) w. (2)\r\nwhere Ge (θ) = G(θ)\r\nσ2\r\nsen\r\n. It is assumed that the transmit and\r\nreceive antenna arrays are adequately spaced to prevent self\ufffeinterference [14]. To account for potential uncertainty in the\r\nvalue of θ, e.g., caused by the target’s speed [15], we adopt\r\nthe model in [11], where an angular interval [θ − ∆, θ + ∆]\r\nis considered, with ∆ representing the uncertainty in θ. This\r\ninterval is discretized into samples, resulting in set n\r\nΘ =\r\n¯θ |\r\n¯θ = θ − ∆ + 2∆\r\nC−1\r\nc\r\no\r\n, ∀c = 0, . . . , C − 1, where C is the\r\nnumber of samples taken within the interval. To ensure that\r\nthe sensing SNR in all angular directions within Θ exceeds\r\nsome value τ, we first include constraint C4 : τ ≥ 0 and then\r\nadd constraint C5 : wHGe (θ) w ≥ τ, ∀θ ∈ Θ.\r\nObjective function: We define the tradeoff function\r\nf (µ, τ) , ρcom · fcom (µ) + ρsen · fsen (τ), (3)\r\nwhich we aim to maximize. Here, fcom (µ) , 1\r\nTµ and\r\nfsen (τ) , τ are the objective functions related to commu\ufffenications and sensing, respectively. In particular, fcom (µ)\r\nrepresents the number of admitted users, i.e., users that are\r\nserved with the desired multicast data, while fsen (τ) is the\r\nlowest sensing SNR value for the angles in Θ. In addition,\r\nµ = [µ1, . . . , µU ]\r\nT\r\n, whereas ρcom and ρsen are the weights\r\nthat control the functionality importance.\r\nProblem formulation: We formulate the joint design of\r\nadmission control and discrete-phase beamforming as\r\nP : maximize\r\nw,µ,τ\r\nf (µ, τ) s.t. C1, C2, C3, C4, C5.\r\nAs a particular case, we consider that communications has\r\nhigher hierarchy than sensing, achieved through the careful\r\ndesign of weights, as outlined in Lemma 1. We highlight that\r\nour framework can accommodate any arbitrary weights, even\r\nwhen hierarchies are not required.\r\nLemma 1. A set of weights ensuring that communications has\r\nhigher hierarchy is given by ρcom = 1 and ρsen =\r\nσ\r\n2\r\nsen\r\n2αNPtx\r\n.\r\nProof. To ensure that the functionalities have different hierar\ufffechies, we choose the weights such that ρcom · fcom (µ) and\r\nρsen · fsen (τ) span nonoverlapping intervals. In particular, we\r\nlet ρcom · fcom (µ) ∈ N handle the integer part of f (µ, τ)\r\nand let ρsen · fsen (τ) ∈ [0, 1) handle the decimal part, thereby\r\neffectively assigning a higher hierarchy to communications.\r\nNote that fcom (µ) is an integer by definition. In order for\r\nρcom · fcom (µ) to also be an integer, we can choose any\r\nρcom ∈ [1, ∞) ∩ N. For simplicity, we adopt ρcom = 1.\r\nBesides, since τ is smaller than or equal to wHGe (θ) w,\r\n∀θ ∈ Θ, as stated in C5, we can establish an upper bound for",
						}))
					})
				})*/

			Context("when the structured page text is requested", func() {
				It("returns the correct structured text", func() {
					pageTextStructured, err := PdfiumInstance.GetPageTextStructured(&requests.GetPageTextStructured{
						Page: requests.Page{
							ByIndex: &requests.PageByIndex{
								Document: doc,
								Index:    1,
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
										Index:    1,
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
										Index:    1,
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
										Index:    1,
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
