package responses

import (
	"github.com/klippa-app/go-pdfium/enums"
)

type FormFieldFlags struct {
	ReadOnly bool
	Required bool
	NoExport bool
}

type FormField struct {
	Type      enums.FPDF_FORMFIELD_TYPE
	Name      string
	Value     *string  // The current value of the field, not used for FPDF_FORMFIELD_TYPE_CHECKBOX, for FPDF_FORMFIELD_TYPE_COMBOBOX andFPDF_FORMFIELD_TYPE_LISTBOX check the Values property. This can be nil for FPDF_FORMFIELD_TYPE_RADIOBUTTON if no option is selected, this can be checked in the IsChecked property.
	Values    []string // The values for fields of type FPDF_FORMFIELD_TYPE_COMBOBOX, FPDF_FORMFIELD_TYPE_LISTBOX, since those can have multiple values selected.
	IsChecked *bool    // Whether the checkbox is checked or one of the radio buttons, used for FPDF_FORMFIELD_TYPE_CHECKBOX and FPDF_FORMFIELD_TYPE_RADIOBUTTON.
	ToolTip   string   // The text shown when hovering the form field.
	Options   []string // The options that are available, used for FPDF_FORMFIELD_TYPE_RADIOBUTTON, FPDF_FORMFIELD_TYPE_COMBOBOX, FPDF_FORMFIELD_TYPE_LISTBOX.
	Flags     FormFieldFlags
}

type GetForm struct {
	Page   int // The page number (0-index based).
	Fields []FormField
}

func (f *GetForm) AfterUnmarshal() {
	emptyString := ""
	emptyBool := false

	// Sadly, we have to do this since an empty slice becomes nil after gob.
	if f.Fields == nil {
		f.Fields = []FormField{}
	}

	for i := range f.Fields {
		// Sadly, we have to do this since an empty string becomes nil after gob.
		if f.Fields[i].Type != enums.FPDF_FORMFIELD_TYPE_CHECKBOX &&
			f.Fields[i].Type != enums.FPDF_FORMFIELD_TYPE_RADIOBUTTON &&
			f.Fields[i].Type != enums.FPDF_FORMFIELD_TYPE_LISTBOX &&
			f.Fields[i].Type != enums.FPDF_FORMFIELD_TYPE_COMBOBOX && f.Fields[i].Value == nil {
			f.Fields[i].Value = &emptyString
		}

		// Sadly, we have to do this since an empty slice becomes nil after gob.
		if (f.Fields[i].Type == enums.FPDF_FORMFIELD_TYPE_LISTBOX ||
			f.Fields[i].Type == enums.FPDF_FORMFIELD_TYPE_COMBOBOX) && f.Fields[i].Values == nil {
			f.Fields[i].Values = []string{}
		}

		// Sadly, we have to do this since an empty string becomes nil after gob.
		if (f.Fields[i].Type == enums.FPDF_FORMFIELD_TYPE_CHECKBOX ||
			f.Fields[i].Type == enums.FPDF_FORMFIELD_TYPE_RADIOBUTTON) && f.Fields[i].IsChecked == nil {
			f.Fields[i].IsChecked = &emptyBool
		}
	}
}
