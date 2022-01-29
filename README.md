# go-pdfium

[![Go Reference](https://pkg.go.dev/badge/github.com/klippa-app/go-pdfium/pdfium.svg)](https://pkg.go.dev/github.com/klippa-app/go-pdfium/pdfium)
[![Build Status][build-status]][build-url]

[build-status]:https://github.com/klippa-app/go-pdfium/workflows/Go/badge.svg

[build-url]:https://github.com/klippa-app/go-pdfium/actions

:rocket: *Easy PDF rendering and text extraction using Go and pdfium* :rocket:

**A fast, multi-threaded and easy to use PDF renderer / text extractor for Go applications.**

## Features

* Option between single-threaded and multi-threaded
* Get page count
* Get plain text of a page
* Get structured text of a page (text, angle, position, size, font information)
* Render 1 or multiple pages into a Go `image.Image` using either DPI or pixel size
* Render the image above directly as a jpeg or png into a file path or byte array
* Get page size in either points or pixel size (when rendered in a specific DPI)
* High test coverage ⭐

## pdfium

This project uses the pdfium C++ library by Google (https://pdfium.googlesource.com/pdfium/) to process the PDF
documents.

## Single/Multi-threading

Since pdfium is not a multithreaded C++ library, we can not directly make it multithreaded by calling it from Go's
subroutines.

This library allows you to call pdfium in a single or multi-threaded way.

We have implemented multi-threading this using [HashiCorp's Go Plugin System](https://github.com/hashicorp/go-plugin),
which allows us launch separate pdfium worker processes, and then route the requests through the different workers. This
also makes it a bit more safe to use pdfium, as it's less likely to segfaults or corrupt your main Go application. The
Plugin system provides the communication between the processes using GRPc, however, when implementing this library, you
won't really see anything of that. From the outside it will look like normal Go code.

Single-threading works by directly calling the pdfium library from the same process. Single-threaded might be preferred
if the caller is managing the workers themselves and does not want the overhead of another process. Be aware that since
pdfium is C++, we can't handle segfaults caused by pdfium, which may cause your process to be killed.

**Be aware that pdfium could use quite some memory depending on the size of the PDF and the requests that you do, so be
aware of the amount of workers that you configure.**

## Prerequisites

To use this Go library, you will need the actual pdfium library to run it and have it available through pkgconfig.

### Get the library

You can try to compile pdfium yourself, but you can also use pre-compiled binaries, for example
from: https://github.com/bblanchon/pdfium-binaries/releases

If you use a pre-compiled library, make sure to extract it somewhere logical, for example /opt/pdfium.

### Configure pkg-config

Create/edit file `/usr/lib/pkgconfig/pdfium.pc`

```
prefix={path}
libdir={path}/lib
includedir={path}/include

Name: pdfium
Description: pdfium
Version: 4849
Requires:

Libs: -L${libdir} -lpdfium
Cflags: -I${includedir}
```

Replace `{path}` with the path you extracted/compiled pdfium in.

Make sure you extend your library path when running:

`export LD_LIBRARY_PATH={path}/lib`

You can do this globally or just in your editor.

this can globally be done on ubuntu by editing `~/.profile`
and adding the line in this file. reloading for bash can be done by relogging or running `source ~/.profile` can be used
to test the change for a terminal

## Getting started

To get started, make sure that you create a separate package in your application that will start the worker.

The examples below can also be found in the examples folder.

### Single-threaded

For single threaded implementations we just have to initialize the library.

`pdfium/renderer/renderer.go`

```go
package renderer

import (
	"github.com/klippa-app/go-pdfium/pdfium"
	"github.com/klippa-app/go-pdfium/pdfium/pdfium_single_threaded"
)

var Pdfium pdfium.Pdfium

func init() {
	// Init the pdfium library and return the instance to open documents.
	Pdfium = pdfium_single_threaded.Init()
}
```

### Multi-threaded

#### Worker package

This package has to be named main to make it available as a binary. The plugin system will use this to start new pdfium
workers. Example:

`pdfium/worker/main.go`

```go
package main

import (
	"github.com/klippa-app/go-pdfium/pdfium/pdfium_multi_threaded/worker"
)

func main() {
	worker.StartWorker()
}
```

#### Worker configuration

To actually start workers, you will have to init the pdfium library somewhere, this also allows you to dynamically start
workers when needed. The best location to add this is in the `init()` of a package that is going to call the pdfium
library. Example:

`pdfium/renderer/renderer.go`

```go
package renderer

import (
	"github.com/klippa-app/go-pdfium/pdfium"
	"github.com/klippa-app/go-pdfium/pdfium/pdfium_multi_threaded"
)

var Pdfium pdfium.Pdfium

func init() {
	// Init the pdfium library and return the instance to open documents.
	// You can tweak these configs to your need. Be aware that workers can use quite some memory.
	Pdfium = pdfium_multi_threaded.Init(pdfium_multi_threaded.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
		Command: pdfium_multi_threaded.Command{
			BinPath: "go",                                     // Only do this while developing, on production put the actual binary path in here. You should not want the Go runtime on production.
			Args:    []string{"run", "pdfium/worker/main.go"}, // This is a reference to the worker package, this can be left empty when using a direct binary path.
		},
	})
}
```

### Get page count

```go
package renderer

import (
	"io/ioutil"
	"log"

	"github.com/klippa-app/go-pdfium/pdfium"
	"github.com/klippa-app/go-pdfium/pdfium/requests"
)

var Pdfium pdfium.Pdfium

// Insert the single/multi-threaded init() here.

func main() {
	filePath := "example.pdf"
	pageCount, err := getPageCount(filePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("The PDF %s has %d pages", filePath, pageCount)
}

func getPageCount(filePath string) (int, error) {
	// Load the PDF file into a byte array.
	pdfBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	// Open the PDF using pdfium (and claim a worker)
	doc, err := Pdfium.NewDocument(&pdfBytes)
	if err != nil {
		return 0, err
	}

	// Always close the document, this will release the worker and it's resources
	defer doc.Close()

	pageCount, err := doc.GetPageCount(&requests.GetPageCount{})
	if err != nil {
		return 0, err
	}

	return pageCount.PageCount, nil
}
```

### Render a page

```go
package renderer

import (
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/klippa-app/go-pdfium/pdfium"
	"github.com/klippa-app/go-pdfium/pdfium/requests"
)

var Pdfium pdfium.Pdfium

// Insert the single/multi-threaded init() here.

func main() {
	filePath := "example.pdf"
	output := "example.pdf.png"
	err := renderPage(filePath, 1, output)
	if err != nil {
		log.Fatal(err)
	}
}

func renderPage(filePath string, page int, output string) error {
	// Load the PDF file into a byte array.
	pdfBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Open the PDF using pdfium (and claim a worker)
	doc, err := Pdfium.NewDocument(&pdfBytes)
	if err != nil {
		return err
	}

	// Always close the document, this will release the worker and it's resources
	defer doc.Close()

	// Render the page in DPI 200.
	pageRender, err := doc.RenderPageInDPI(&requests.RenderPageInDPI{
		DPI:  200,      // The DPI to render the page in.
		Page: page - 1, // The page to render, 0-indexed.
	})
	if err != nil {
		return err
	}

	// Write the output to a file.
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, pageRender.Image)
	if err != nil {
		return err
	}

	return nil
}
```

## About Klippa

Founded in 2015, [Klippa](https://www.klippa.com/en)'s goal is to digitize & automate administrative processes with
modern technologies. We help clients enhance the effectiveness of their organization by using machine learning and OCR.
Since 2015, more than a thousand happy clients have used Klippa's software solutions. Klippa currently has an
international team of 50 people, with offices in Groningen, Amsterdam and Brasov.

## License

The MIT License (MIT)
