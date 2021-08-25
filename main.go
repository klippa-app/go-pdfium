package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/klippa-app/go-pdfium/pdfium/requests"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"

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

	file, err := os.Open("/home/jeroen/Downloads/Coolblue_Factuur_647776503.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := pdfium.NewDocument(&fileData)
	if err != nil {
		log.Fatal(err)
	}

	defer doc.Close()
	pageCount, err := doc.GetPageCount(&requests.GetPageCount{})
	if err != nil {
		log.Fatal(err)
	}

	if false {
		for i := 1; i <= pageCount.PageCount; i++ {
			size, err := doc.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
				Page: i - 1,
				DPI:  300,
			})
			if err != nil {
				log.Fatal(err)
			}

			log.Println(size)
		}
	}

	if false {
		for i := 1; i <= pageCount.PageCount; i++ {
			text, err := doc.GetPageText(&requests.GetPageText{
				Page: i - 1,
			})
			if err != nil {
				log.Fatal(err)
			}

			if text != nil {
				log.Println(*text)
			}
		}
	}

	if false {
		for i := 1; i <= pageCount.PageCount; i++ {
			renderedImage, err := doc.RenderPageInPixels(&requests.RenderPageInPixels{
				Page:   i - 1,
				Width:  2000,
				Height: 2000,
			})
			if err != nil {
				log.Fatal(err)
			}

			var opt jpeg.Options
			opt.Quality = 95

			var imgBuf bytes.Buffer
			err = jpeg.Encode(&imgBuf, renderedImage.Image, &opt)
			if err != nil {
				log.Fatal(err)
			}

			ioutil.WriteFile(fmt.Sprintf("tmp/image-%d.jpg", i), imgBuf.Bytes(), 0777)
		}
	}

	if true {
		for i := 1; i <= pageCount.PageCount; i++ {
			structuredText, err := doc.GetPageTextStructured(&requests.GetPageTextStructured{
				Page: i - 1,
				PixelPositions: requests.GetPageTextStructuredPixelPositions{
					Calculate: true,
					Width:     2000,
					Height:    2000,
				},
				CollectFontInformation: true,
			})
			if err != nil {
				log.Fatal(err)
			}

			/*
				if structuredText != nil {
					for _, rect := range structuredText.Rects {
						log.Println(rect.Text)
					}
				}*/

			jsonData, _ := json.MarshalIndent(structuredText, "", "    ")
			ioutil.WriteFile(fmt.Sprintf("tmp/bounds-%d.json", i), jsonData, 0777)
		}
	}
}
