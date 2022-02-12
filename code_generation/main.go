package main

// This tool is to generate the go-pdfium implementations.
// The implementations follow a format for input/output which makes it easy to
// generate the implementations, saving a lot of copy-pasting time.

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
	"text/template"

	"github.com/klippa-app/go-pdfium"
)

type GenerateDataMethod struct {
	Name   string
	Input  string
	Output string
}

func (m *GenerateDataMethod) BlockForMultiThreaded() bool {
	if m.Name == "FPDFBitmap_CreateEx" ||
		m.Name == "FSDK_SetUnSpObjProcessHandler" ||
		m.Name == "FSDK_SetTimeFunction" ||
		m.Name == "FSDK_SetLocaltimeFunction" ||
		m.Name == "FPDF_RenderPage" ||
		m.Name == "FPDF_RenderPageBitmapWithColorScheme_Start" ||
		m.Name == "FPDF_RenderPageBitmap_Start" ||
		m.Name == "FPDF_RenderPage_Continue" ||
		m.Name == "FPDF_RenderPage_Close" ||
		strings.HasPrefix(m.Name, "FPDFAvail_") {
		return true
	}
	return false
}

type GenerateData struct {
	Methods []GenerateDataMethod
}

type Template struct {
	Source string
	Target string
}

func main() {
	data := GenerateData{
		Methods: []GenerateDataMethod{},
	}

	docType := reflect.TypeOf((*pdfium.Pdfium)(nil)).Elem()
	numMethods := docType.NumMethod()

	for i := 0; i < numMethods; i++ {
		method := docType.Method(i)

		// These are special, don't generate them
		if method.Name == "Close" {
			continue
		}

		dataMethod := GenerateDataMethod{
			Name:   method.Name,
			Input:  method.Name,
			Output: method.Name,
		}

		data.Methods = append(data.Methods, dataMethod)
	}

	templates := []Template{
		{
			Source: "code_generation/templates/single_threaded.go.tmpl",
			Target: "single_threaded/generated.go",
		},
		{
			Source: "code_generation/templates/multi_threaded.go.tmpl",
			Target: "multi_threaded/generated.go",
		},
		{
			Source: "code_generation/templates/grpc.go.tmpl",
			Target: "internal/commons/generated.go",
		},
	}
	for i := range templates {
		err := generateFromTemplate(templates[i], data)
		if err != nil {
			log.Fatalf("Could not generate template %s: %s", templates[i], err.Error())
		}
	}
}

func generateFromTemplate(codeTemplate Template, data GenerateData) error {
	templateContent, err := ioutil.ReadFile(codeTemplate.Source)
	if err != nil {
		return err
	}
	t, err := template.New(path.Base(codeTemplate.Source)).Parse(string(templateContent))
	if err != nil {
		return err
	}

	f, err := os.Create(codeTemplate.Target)
	if err != nil {
		return err
	}

	err = t.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}
