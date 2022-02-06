# go-pdfium

[![Go Reference](https://pkg.go.dev/badge/github.com/klippa-app/go-pdfium/pdfium.svg)](https://pkg.go.dev/github.com/klippa-app/go-pdfium)
[![Build Status][build-status]][build-url]
[![codecov](https://codecov.io/gh/klippa-app/go-pdfium/branch/main/graph/badge.svg?token=WoIlW9RbfH)](https://codecov.io/gh/klippa-app/go-pdfium)

[build-status]:https://github.com/klippa-app/go-pdfium/workflows/Go/badge.svg

[build-url]:https://github.com/klippa-app/go-pdfium/actions

:rocket: *Easy to use PDF library using Go and PDFium* :rocket:

**A fast, multi-threaded and easy to use PDF library for Go applications.**

## Features

* Option between single-threaded and multi-threaded (through subprocesses), while keeping the same interface
* This library will handle all complicated cgo gymnastics for you
* The goal is to implement all PDFium public API methods (including experimental), current progress: 40%
* Current PDFium methods exposed, no cgo required
    * PDFium instance configuration (sandbox policy, fonts)
    * Document loading (from bytes, path or io.ReadSeeker)
    * Document info (metadata, page count, render mode, PDF version, permissions, security handler revision)
    * Page info (size, transparency)
    * Rendering (through bitmap)
    * Bitmap handling
    * Named destinations
    * Text handling (extract, search, text size/color/font information)
    * Creation (create new documents and pages)
    * Editing (rotating, import pages from another document, copy view preferences from another document, flattening)
    * Bookmarks / Links / Weblinks
    * Document saving (to bytes, path or io.Writer)
    * JavaScript actions
    * Thumbnails
    * Attachments
    * XFA packet handling
    * ViewerRef (print settings)
* Methods to be implemented:
    * Form filling
    * Transformations (page boxes, clip paths)
    * Annotations
    * Document loading through data availability
    * Progressive rendering
    * Struct trees
* Methods that won't be implemented for now:
    * Win32-only methods
    * fpdf_sysfontinfo.h (probably too complicated)
* Useful helpers to make your life easier:
    * Get all document metadata
    * Get all document bookmarks
    * Get all document attachments
    * Get all document JavaScript actions
    * Get plain text of a page
    * Get structured text of a page (text, angle, position, size, font information)
    * Render 1 or multiple pages from 1 or multiple documents into a Go `image.Image` using either DPI or pixel size
    * Use the same render instructions to render the image directly as a jpeg or png into a file path or byte array
    * Get page size in either points or pixel size (when rendered in a specific DPI)
    * Get the point to pixel ratio when rendering or extracting text (to determine the positions when rendering into an
      image)

## PDFium

This project uses the PDFium C++ library by Google (https://pdfium.googlesource.com/pdfium/) to process the PDF
documents.

## Single/Multi-threading

Since PDFium is [not a multithreaded C++ library](https://groups.google.com/g/pdfium/c/HeZSsM_KEUk), we can not directly
make it multithreaded by calling it from Go's subroutines.

This library allows you to call PDFium in a single or multi-threaded way.

We have implemented multi-threading this using [HashiCorp's Go Plugin System](https://github.com/hashicorp/go-plugin),
which allows us launch separate PDFium worker processes, and then route the requests through the different workers. This
also makes it a bit more safe to use PDFium, as it's less likely to segfaults or corrupt your main Go application. The
Plugin system provides the communication between the processes using gRPC, however, when implementing this library, you
won't really see anything of that. From the outside it will look like normal Go code. The inter-process communication
does come with a cost as it has to serialize/deserialize input/output as it moves between the main process and the
PDFium workers.

Single-threading works by directly calling the PDFium library from the same process. Single-threaded might be preferred
if the caller is managing the workers themselves and does not want the overhead of another process. Be aware that since
PDFium is C++, we can't handle segfaults caused by PDFium, which may cause your process to be killed. So using this
library in the multi-threaded way, with only 1 worker, can still have some benefits, since it can automatically recover
from things like segfaults.

Both the single-threaded and multi-threaded implementation are thread/subroutine safe, this has been guaranteed by
locking the instance that's doing your work while it's doing PDFium operations. New operations will wait until the lock
becomes available again.

**Be aware that PDFium could use quite some memory depending on the size of the PDF and the requests that you do, so be
aware of the amount of workers that you configure.**

### `io.ReadSeeker` and `io.Writer`

Document loading allows you to load a document with a `io.ReadSeeker`. Please be aware that this only works efficiently
when using the single-threaded usage, as that lives in the same process. For multi-threaded usage this will just load in
the complete file and pass the bytes through the gRPC interface.

Document/image saving allows you to save using a `io.Writer`. Please be aware this only works when using the
single-threaded usage. It's not possible to encode the `io.Writer` with gRPC. Or share it between processes for that
matter.

## Prerequisites

To use this Go library, you will need the actual PDFium library to run it and have it available through pkgconfig.

### Get the library

You can try to compile PDFium yourself, but you can also use pre-compiled binaries, for example
from: https://github.com/bblanchon/pdfium-binaries/releases

If you use a pre-compiled library, make sure to extract it somewhere logical, for example /opt/pdfium.

### Configure pkg-config

Create/edit file `/usr/lib/pkgconfig/pdfium.pc`

```
prefix={path}
libdir={path}/lib
includedir={path}/include

Name: PDFium
Description: PDFium
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
	"log"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/single_threaded"
)

// Be sure to close pools/instances when you're done with them.
var pool pdfium.Pool
var instance pdfium.Pdfium

func init() {
	// Init the PDFium library and return the instance to open documents.
	pool = single_threaded.Init(single_threaded.Config{})

	var err error
	instance, err = pool.GetInstance(time.Second * 30)
	if err != nil {
		log.Fatal(err)
	}
}
```

### Multi-threaded

#### Worker package

This package has to be named main to make it available as a binary. The plugin system will use this to start new PDFium
workers. Example:

`pdfium/worker/main.go`

```go
package main

import (
	"github.com/klippa-app/go-pdfium/multi_threaded/worker"
)

func main() {
	worker.StartWorker()
}
```

#### Worker configuration

To actually start workers, you will have to init the PDFium library somewhere, this also allows you to dynamically start
workers when needed. The best location to add this is in the `init()` of a package that is going to call the PDFium
library. Example:

`pdfium/renderer/renderer.go`

```go
package renderer

import (
	"log"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/multi_threaded"
)

// Be sure to close pools/instances when you're done with them.
var pool pdfium.Pool
var instance pdfium.Pdfium

func init() {
	// Init the PDFium library and return the instance to open documents.
	// You can tweak these configs to your need. Be aware that workers can use quite some memory.
	pool = multi_threaded.Init(multi_threaded.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
		Command: multi_threaded.Command{
			BinPath: "go",                                     // Only do this while developing, on production put the actual binary path in here. You should not want the Go runtime on production.
			Args:    []string{"run", "pdfium/worker/main.go"}, // This is a reference to the worker package, this can be left empty when using a direct binary path.
		},
	})

	var err error
	instance, err = pool.GetInstance(time.Second * 30)
	if err != nil {
		log.Fatal(err)
	}
}
```

### Get page count

```go
package renderer

import (
	"io/ioutil"
	"log"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
)

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

	// Open the PDF using PDFium (and claim a worker)
	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &pdfBytes,
	})
	if err != nil {
		return 0, err
	}

	// Always close the document, this will release its resources.
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	pageCount, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: doc,
	})
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

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
)

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

	// Open the PDF using PDFium (and claim a worker)
	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &pdfBytes,
	})
	if err != nil {
		return err
	}

	// Always close the document, this will release its resources.
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	// Render the page in DPI 200.
	pageRender, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
		DPI: 200, // The DPI to render the page in.
		Page: requests.Page{
			ByIndex: &requests.PageByIndex{
				Document: doc,
				Index:    0,
			},
		}, // The page to render, 0-indexed.
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
