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

func toPointer[T any](input T) *T {
	return &input
}

var _ = Describe("fpdf_attachment", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	Context("no document is given", func() {
		It("returns an error when GetForm is called", func() {
			GetForm, err := PdfiumInstance.GetForm(&requests.GetForm{})
			Expect(err).To(MatchError("either page reference or index should be given"))
			Expect(GetForm).To(BeNil())
		})
	})

	Context("a normal PDF file without form fields", func() {
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

		It("returns no form fields", func() {
			GetForm, err := PdfiumInstance.GetForm(&requests.GetForm{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(GetForm).To(Equal(&responses.GetForm{
				Page:   0,
				Fields: []responses.FormField{},
			}))
		})
	})

	Context("a normal PDF file with multiple form text fields", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/text_form_multiple.pdf")
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

		It("returns the correct form fields", func() {
			GetForm, err := PdfiumInstance.GetForm(&requests.GetForm{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(GetForm).To(Equal(&responses.GetForm{
				Page: 0,
				Fields: []responses.FormField{
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_TEXTFIELD,
						Name:      "Text Box",
						Value:     toPointer(""),
						IsChecked: nil,
						ToolTip:   "",
						Options:   nil,
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_TEXTFIELD,
						Name:      "ReadOnly",
						Value:     toPointer(""),
						IsChecked: nil,
						ToolTip:   "",
						Options:   nil,
						Flags:     responses.FormFieldFlags{ReadOnly: true, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_TEXTFIELD,
						Name:      "CharLimit",
						Value:     toPointer("Elephant"),
						IsChecked: nil,
						ToolTip:   "",
						Options:   nil,
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_TEXTFIELD,
						Name:      "Password",
						Value:     toPointer(""),
						IsChecked: nil,
						ToolTip:   "",
						Options:   nil,
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
				},
			}))
		})
	})

	Context("a normal PDF file with multiple click fields (radio/checkbox)", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/click_form.pdf")
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

		It("returns the correct form fields", func() {
			GetForm, err := PdfiumInstance.GetForm(&requests.GetForm{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
			})

			Expect(err).To(BeNil())
			Expect(GetForm).To(Equal(&responses.GetForm{
				Page: 0,
				Fields: []responses.FormField{
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_CHECKBOX,
						Name:      "readOnlyCheckbox",
						Values:    nil,
						IsChecked: toPointer(true),
						ToolTip:   "readOnlyCheckbox",
						Options:   nil,
						Flags:     responses.FormFieldFlags{ReadOnly: true, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_CHECKBOX,
						Name:      "checkbox",
						Values:    nil,
						IsChecked: toPointer(false),
						ToolTip:   "checkbox",
						Options:   nil,
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: true, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_RADIOBUTTON,
						Name:      "readOnlyRadioButton",
						Value:     toPointer("value3"),
						IsChecked: toPointer(true),
						ToolTip:   "readOnlyRadioButton1",
						Options:   []string{"value1", "value2", "value3"},
						Flags:     responses.FormFieldFlags{ReadOnly: true, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_RADIOBUTTON,
						Name:      "radioButton",
						Value:     toPointer("value3"),
						IsChecked: toPointer(true),
						ToolTip:   "radioButton1",
						Options:   []string{"value1", "value2", "value3"},
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: true, NoExport: false},
					},
				},
			}))
		})
	})

	Context("a normal PDF file with listbox form fields", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/listbox_form.pdf")
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

		It("returns the correct form fields", func() {
			GetForm, err := PdfiumInstance.GetForm(&requests.GetForm{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
			})

			Expect(err).To(BeNil())
			Expect(GetForm).To(Equal(&responses.GetForm{
				Page: 0,
				Fields: []responses.FormField{
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_LISTBOX,
						Name:      "Listbox_SingleSelect",
						Values:    []string{},
						IsChecked: nil,
						ToolTip:   "",
						Options:   []string{"Foo", "Bar", "Qux"},
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_LISTBOX,
						Name:      "Listbox_MultiSelect",
						Values:    []string{"Banana"},
						IsChecked: nil,
						ToolTip:   "",
						Options:   []string{"Apple", "Banana", "Cherry", "Date", "Elderberry", "Fig", "Guava", "Honeydew", "Indian Fig", "Jackfruit", "Kiwi", "Lemon", "Mango", "Nectarine", "Orange", "Persimmon", "Quince", "Raspberry", "Strawberry", "Tamarind", "Ugli Fruit", "Voavanga", "Wolfberry", "Xigua", "Yangmei", "Zucchini"},
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_LISTBOX,
						Name:      "Listbox_ReadOnly",
						Values:    []string{},
						IsChecked: nil,
						ToolTip:   "",
						Options:   []string{"Dog", "Elephant", "Frog"},
						Flags:     responses.FormFieldFlags{ReadOnly: true, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_LISTBOX,
						Name:      "Listbox_MultiSelectMultipleIndices",
						Values:    []string{"Belgium", "Denmark"},
						IsChecked: nil,
						ToolTip:   "",
						Options:   []string{"Albania", "Belgium", "Croatia", "Denmark", "Estonia"},
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_LISTBOX,
						Name:      "Listbox_MultiSelectMultipleValues",
						Values:    []string{"Gamma", "Epsilon"},
						IsChecked: nil,
						ToolTip:   "",
						Options:   []string{"Alpha", "Beta", "Gamma", "Delta", "Epsilon"},
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_LISTBOX,
						Name:      "Listbox_MultiSelectMultipleMismatch",
						Values:    []string{"Alligator", "Cougar"},
						IsChecked: nil,
						ToolTip:   "",
						Options:   []string{"Alligator", "Bear", "Cougar", "Deer", "Echidna"},
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_LISTBOX,
						Name:      "Listbox_SingleSelectLastSelected",
						Values:    []string{"Saskatchewan"},
						IsChecked: nil,
						ToolTip:   "",
						Options: []string{
							"Alberta",
							"British Columbia",
							"Manitoba",
							"New Brunswick",
							"Newfoundland and Labrador",
							"Nova Scotia",
							"Ontario",
							"Prince Edward Island",
							"Quebec",
							"Saskatchewan",
						},
						Flags: responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
				},
			}))
		})
	})

	Context("a normal PDF file with combobox form fields", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/combobox_form.pdf")
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

		It("returns the correct form fields", func() {
			GetForm, err := PdfiumInstance.GetForm(&requests.GetForm{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
			})

			Expect(err).To(BeNil())
			Expect(GetForm).To(Equal(&responses.GetForm{
				Page: 0,
				Fields: []responses.FormField{
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_COMBOBOX,
						Name:      "Combo_Editable",
						Value:     nil,
						Values:    []string{},
						IsChecked: nil,
						ToolTip:   "",
						Options:   []string{"Foo", "Bar", "Qux"},
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_COMBOBOX,
						Name:      "Combo1",
						Value:     nil,
						Values:    []string{"Banana"},
						IsChecked: nil,
						ToolTip:   "",
						Options:   []string{"Apple", "Banana", "Cherry", "Date", "Elderberry", "Fig", "Guava", "Honeydew", "Indian Fig", "Jackfruit", "Kiwi", "Lemon", "Mango", "Nectarine", "Orange", "Persimmon", "Quince", "Raspberry", "Strawberry", "Tamarind", "Ugli Fruit", "Voavanga", "Wolfberry", "Xigua", "Yangmei", "Zucchini"},
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_COMBOBOX,
						Name:      "Combo_ReadOnly",
						Value:     nil,
						Values:    []string{},
						IsChecked: nil,
						ToolTip:   "",
						Options:   []string{"Dog", "Elephant", "Frog"},
						Flags:     responses.FormFieldFlags{ReadOnly: true, Required: false, NoExport: false},
					},
				},
			}))
		})
	})

	Context("a normal PDF file with multiple form field types", func() {
		var doc references.FPDF_DOCUMENT

		BeforeEach(func() {
			pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/multiple_form_types.pdf")
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

		It("returns the correct form fields", func() {
			GetForm, err := PdfiumInstance.GetForm(&requests.GetForm{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc,
						Index:    0,
					},
				},
			})

			Expect(err).To(BeNil())
			Expect(GetForm).To(Equal(&responses.GetForm{
				Page: 0,
				Fields: []responses.FormField{
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_COMBOBOX,
						Name:      "Combo_Editable",
						Value:     nil,
						Values:    []string{},
						IsChecked: nil,
						ToolTip:   "",
						Options:   []string{"Foo", "Bar", "Qux"},
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_LISTBOX,
						Name:      "Listbox_MultiSelectMultipleSelected",
						Value:     nil,
						Values:    []string{"Gamma", "Epsilon"},
						IsChecked: nil,
						ToolTip:   "",
						Options:   []string{"Alpha", "Beta", "Gamma", "Delta", "Epsilon"},
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_TEXTFIELD,
						Name:      "Text Box",
						Value:     toPointer(""),
						Values:    nil,
						IsChecked: nil,
						ToolTip:   "",
						Options:   nil,
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_CHECKBOX,
						Name:      "Checkbox",
						Value:     nil,
						Values:    nil,
						IsChecked: toPointer(false),
						ToolTip:   "",
						Options:   nil,
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
					{
						Type:      enums.FPDF_FORMFIELD_TYPE_RADIOBUTTON,
						Name:      "radioButton",
						Value:     nil,
						Values:    nil,
						IsChecked: toPointer(false),
						ToolTip:   "radioButton1",
						Options:   []string{"Yes"},
						Flags:     responses.FormFieldFlags{ReadOnly: false, Required: false, NoExport: false},
					},
				},
			}))
		})
	})
})
