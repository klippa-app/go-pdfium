package main

import (
	"fmt"

	"github.com/klippa-app/go-pdfium/pdfium"
)

func main() {
	pdfium.InitLibrary(pdfium.Config{
		MinIdle:  1,
		MaxIdle:  1,
		MaxTotal: 1,
		LogCallback: func(s string) {
			fmt.Println("[PDFIUM ERROR]: " + s)
		},
		Command: pdfium.Command{
			BinPath: "go",
			Args:    []string{"run", "./subprocess"},
		},
	})
}
