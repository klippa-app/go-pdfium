package main

import "github.com/klippa-app/go-pdfium/pdfium"

func main() {
	pdfium.InitLibrary(pdfium.Config{
		MinIdle:  1,
		MaxIdle:  1,
		MaxTotal: 1,
	})
}
