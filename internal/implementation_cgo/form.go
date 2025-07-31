//go:build pdfium_experimental
// +build pdfium_experimental

package implementation_cgo

import (
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/structs"
)

// GetForm returns the form elements in the given page, including option values and current values.
// Experimental API.
func (p *PdfiumImplementation) GetForm(request *requests.GetForm) (*responses.GetForm, error) {
	p.Lock()
	defer p.Unlock()

	pageHandle, err := p.loadPage(request.Page)
	if err != nil {
		return nil, err
	}

	formFillEnv, err := p.internal_FPDFDOC_InitFormFillEnvironment(&requests.FPDFDOC_InitFormFillEnvironment{
		Document: pageHandle.documentRef,
		FormFillInfo: structs.FPDF_FORMFILLINFO{
			FFI_Invalidate:         func(page references.FPDF_PAGE, left, top, right, bottom float64) {},
			FFI_SetCursor:          func(cursorType enums.FXCT) {},
			FFI_SetTimer:           func(elapse int, timerFunc func(idEvent int)) int { return 0 },
			FFI_KillTimer:          func(timerID int) {},
			FFI_GetLocalTime:       func() structs.FPDF_SYSTEMTIME { return structs.FPDF_SYSTEMTIME{} },
			FFI_GetPage:            func(document references.FPDF_DOCUMENT, index int) *references.FPDF_PAGE { return nil },
			FFI_GetRotation:        func(page references.FPDF_PAGE) enums.FPDF_PAGE_ROTATION { return enums.FPDF_PAGE_ROTATION_NONE },
			FFI_ExecuteNamedAction: func(namedAction string) {},
		},
	})
	if err != nil {
		return nil, err
	}

	defer func() {
		p.internal_FPDFDOC_ExitFormFillEnvironment(&requests.FPDFDOC_ExitFormFillEnvironment{
			FormHandle: formFillEnv.FormHandle,
		})
	}()

	results := []*responses.FormField{}
	nameReference := map[string]*responses.FormField{}

	// This is required to make sure all the information is there about the form.
	_, err = p.internal_FORM_OnAfterLoadPage(&requests.FORM_OnAfterLoadPage{
		Page: requests.Page{
			ByReference: &pageHandle.nativeRef,
		},
		FormHandle: formFillEnv.FormHandle,
	})
	if err != nil {
		return nil, err
	}

	defer func() {
		p.internal_FORM_OnBeforeClosePage(&requests.FORM_OnBeforeClosePage{
			Page: requests.Page{
				ByReference: &pageHandle.nativeRef,
			},
			FormHandle: formFillEnv.FormHandle,
		})
	}()

	annotationCount, err := p.internal_FPDFPage_GetAnnotCount(&requests.FPDFPage_GetAnnotCount{
		Page: requests.Page{
			ByReference: &pageHandle.nativeRef,
		},
	})
	if err != nil {
		return nil, err
	}

	for annotationIndex := 0; annotationIndex < annotationCount.Count; annotationIndex++ {
		annotation, err := p.internal_FPDFPage_GetAnnot(&requests.FPDFPage_GetAnnot{
			Page: requests.Page{
				ByReference: &pageHandle.nativeRef,
			},
			Index: annotationIndex,
		})
		if err != nil {
			return nil, err
		}

		subType, err := p.internal_FPDFAnnot_GetSubtype(&requests.FPDFAnnot_GetSubtype{
			Annotation: annotation.Annotation,
		})
		if err != nil {
			return nil, err
		}

		// We only want widgets (form fields).
		if subType.Subtype != enums.FPDF_ANNOT_SUBTYPE_WIDGET {
			continue
		}

		formType, err := p.internal_FPDFAnnot_GetFormFieldType(&requests.FPDFAnnot_GetFormFieldType{
			FormHandle: formFillEnv.FormHandle,
			Annotation: annotation.Annotation,
		})
		if err != nil {
			return nil, err
		}

		fieldName, err := p.internal_FPDFAnnot_GetFormFieldName(&requests.FPDFAnnot_GetFormFieldName{
			Annotation: annotation.Annotation,
			FormHandle: formFillEnv.FormHandle,
		})
		if err != nil {
			return nil, err
		}

		formField := responses.FormField{
			Type: formType.FormFieldType,
			Name: fieldName.FormFieldName,
		}

		if formField.Type == enums.FPDF_FORMFIELD_TYPE_CHECKBOX {
			// For checkboxes, we rely on the IsChecked information; since the
			// value is often just made up by pdfium, we don't load it in.
			isChecked, err := p.internal_FPDFAnnot_IsChecked(&requests.FPDFAnnot_IsChecked{
				Annotation: annotation.Annotation,
				FormHandle: formFillEnv.FormHandle,
			})
			if err != nil {
				return nil, err
			}
			formField.IsChecked = &isChecked.IsChecked
		} else if formField.Type == enums.FPDF_FORMFIELD_TYPE_RADIOBUTTON {
			// Radiobuttons allow for multiple options, each option is its own widget.
			// The form control index tells us where this widget is in the list of options.
			index, err := p.internal_FPDFAnnot_GetFormControlIndex(&requests.FPDFAnnot_GetFormControlIndex{
				Annotation: annotation.Annotation,
				FormHandle: formFillEnv.FormHandle,
			})
			if err != nil {
				return nil, err
			}

			// The export value is what the value would be if this was the selected value.
			exportValue, err := p.internal_FPDFAnnot_GetFormFieldExportValue(&requests.FPDFAnnot_GetFormFieldExportValue{
				Annotation: annotation.Annotation,
				FormHandle: formFillEnv.FormHandle,
			})
			if err != nil {
				return nil, err
			}

			isChecked, err := p.internal_FPDFAnnot_IsChecked(&requests.FPDFAnnot_IsChecked{
				Annotation: annotation.Annotation,
				FormHandle: formFillEnv.FormHandle,
			})
			if err != nil {
				return nil, err
			}

			// Since radio options have multiple widgets for the same field, we
			// keep a reference to the form element to put the other options
			// on the same form field result.
			existingFormField, formFieldExists := nameReference[formField.Name]
			if !formFieldExists {
				// The form control count lets us know the total number of
				// controls in this group (the total number of options).
				controlCount, err := p.internal_FPDFAnnot_GetFormControlCount(&requests.FPDFAnnot_GetFormControlCount{
					Annotation: annotation.Annotation,
					FormHandle: formFillEnv.FormHandle,
				})
				if err != nil {
					return nil, err
				}

				formField.IsChecked = &isChecked.IsChecked
				formField.Options = make([]string, controlCount.FormControlCount)
				formField.Options[index.FormControlIndex] = exportValue.Value
			} else {
				// We have seen this field already, add the current export
				// value to the field and don't do anything else.
				existingFormField.Options[index.FormControlIndex] = exportValue.Value
				if isChecked.IsChecked {
					existingFormField.IsChecked = &isChecked.IsChecked
				}
				continue
			}
		}

		// Listboxes and Comboboxes are special, the options need to be fetched
		// separately, and it can also have multiple selected values.
		if formType.FormFieldType == enums.FPDF_FORMFIELD_TYPE_LISTBOX || formType.FormFieldType == enums.FPDF_FORMFIELD_TYPE_COMBOBOX {
			optionCount, err := p.internal_FPDFAnnot_GetOptionCount(&requests.FPDFAnnot_GetOptionCount{
				Annotation: annotation.Annotation,
				FormHandle: formFillEnv.FormHandle,
			})
			if err != nil {
				return nil, err
			}

			formField.Options = make([]string, optionCount.OptionCount)
			formField.Values = []string{}

			// Loop through the available options.
			for optionIndex := 0; optionIndex < optionCount.OptionCount; optionIndex++ {
				optionLabel, err := p.internal_FPDFAnnot_GetOptionLabel(&requests.FPDFAnnot_GetOptionLabel{
					Annotation: annotation.Annotation,
					FormHandle: formFillEnv.FormHandle,
					Index:      optionIndex,
				})
				if err != nil {
					return nil, err
				}
				formField.Options[optionIndex] = optionLabel.OptionLabel

				isSelected, err := p.internal_FPDFAnnot_IsOptionSelected(&requests.FPDFAnnot_IsOptionSelected{
					Annotation: annotation.Annotation,
					FormHandle: formFillEnv.FormHandle,
					Index:      optionIndex,
				})
				if err != nil {
					return nil, err
				}

				// When the option is selected, add the label to the value
				// list.
				if isSelected.IsOptionSelected {
					formField.Values = append(formField.Values, optionLabel.OptionLabel)
				}
			}
		}

		// Checkboxes don't have a sensible value, we use IsChecked there.
		// For ListBox and ComboBox we use FPDFAnnot_IsOptionSelected.
		if formType.FormFieldType != enums.FPDF_FORMFIELD_TYPE_CHECKBOX &&
			formType.FormFieldType != enums.FPDF_FORMFIELD_TYPE_LISTBOX &&
			formType.FormFieldType != enums.FPDF_FORMFIELD_TYPE_COMBOBOX {
			fieldValue, err := p.internal_FPDFAnnot_GetFormFieldValue(&requests.FPDFAnnot_GetFormFieldValue{
				Annotation: annotation.Annotation,
				FormHandle: formFillEnv.FormHandle,
			})
			if err != nil {
				return nil, err
			}

			formField.Value = &fieldValue.FormFieldValue
		}

		// The alternate name is what is shown as a tooltip when hovering over
		// the field.
		alternateName, err := p.internal_FPDFAnnot_GetFormFieldAlternateName(&requests.FPDFAnnot_GetFormFieldAlternateName{
			Annotation: annotation.Annotation,
			FormHandle: formFillEnv.FormHandle,
		})
		if err != nil {
			return nil, err
		}

		formField.ToolTip = alternateName.FormFieldAlternateName

		// Form field flags contain information about whether a field is
		// read-only, required, and allowed to be exported.
		formFieldFlags, err := p.internal_FPDFAnnot_GetFormFieldFlags(&requests.FPDFAnnot_GetFormFieldFlags{
			Annotation: annotation.Annotation,
			FormHandle: formFillEnv.FormHandle,
		})
		if err != nil {
			return nil, err
		}

		formField.Flags = responses.FormFieldFlags{
			ReadOnly: formFieldFlags.Flags&enums.FPDF_FORMFLAG_READONLY != 0,
			Required: formFieldFlags.Flags&enums.FPDF_FORMFLAG_REQUIRED != 0,
			NoExport: formFieldFlags.Flags&enums.FPDF_FORMFLAG_NOEXPORT != 0,
		}

		results = append(results, &formField)

		nameReference[formField.Name] = &formField

		_, err = p.internal_FPDFPage_CloseAnnot(&requests.FPDFPage_CloseAnnot{
			Annotation: annotation.Annotation,
		})
		if err != nil {
			return nil, err
		}
	}

	fields := []responses.FormField{}
	for i := range results {
		// Clear out values for radios that are not selected.
		if results[i].Type == enums.FPDF_FORMFIELD_TYPE_RADIOBUTTON && !*results[i].IsChecked {
			results[i].Value = nil
		}

		fields = append(fields, *results[i])
	}

	return &responses.GetForm{
		Page:   pageHandle.index,
		Fields: fields,
	}, nil
}
