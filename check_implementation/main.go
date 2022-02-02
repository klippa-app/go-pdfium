package main

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

// A tool to ensure that all pdfium methods are implemented.
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

	for _, item := range items {
		if !item.IsDir() {
			if strings.HasSuffix(item.Name(), ".h") {
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

						if _, ok := implementedMethods[method]; !ok {
							fmt.Println(method)
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
}
