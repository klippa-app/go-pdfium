package main

// A tool to ensure that all PDFium methods are implemented.
// This tool also measures implementation progress.

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/klippa-app/go-pdfium"
)

func main() {
	implementedMethods := map[string]bool{}
	docType := reflect.TypeOf((*pdfium.Pdfium)(nil)).Elem()
	numMethods := docType.NumMethod()

	fmt.Println("Currently implemented methods:")
	for i := 0; i < numMethods; i++ {
		method := docType.Method(i)
		implementedMethods[method.Name] = true
		fmt.Println(method.Name)
	}

	pdfiumFolder := "/opt/lib/pdfium"
	items, err := ioutil.ReadDir(pdfiumFolder + "/include")
	if err != nil {
		log.Fatal(err)
	}

	methodCount := 0
	implementedCount := 0
	for _, item := range items {
		if !item.IsDir() {
			if strings.HasSuffix(item.Name(), ".h") {
				// Skip fpdf_sysfontinfo.h
				if item.Name() == "fpdf_sysfontinfo.h" {
					continue
				}

				file, err := os.Open(pdfiumFolder + "/include/" + item.Name())
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("\nMethods missing in: %s\n", item.Name())

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()
					if strings.HasPrefix(line, "FPDF_EXPORT") {
						callConvSplit := strings.Split(line, "FPDF_CALLCONV")
						method := strings.TrimSpace(callConvSplit[1])
						if method == "" {
							scanner.Scan()
							nextLine := scanner.Text()
							method = nextLine
						}

						methodParts := strings.SplitN(method, "(", 2)
						method = methodParts[0]

						// Skip methods that should never be implemented.
						if method == "FPDF_InitLibrary" || method == "FPDF_InitLibraryWithConfig" || method == "FPDF_DestroyLibrary" {
							continue
						}

						methodCount++
						if _, ok := implementedMethods[method]; !ok {
							fmt.Println(method)
						} else {
							implementedCount++
						}
					}
				}

				if err := scanner.Err(); err != nil {
					log.Fatal(err)
				}

				file.Close()
			}
		}
	}

	log.Printf("Methods to implement: %d, methods implemented: %d, progress: %f", methodCount, implementedCount, (float64(implementedCount)/float64(methodCount))*100)
}
