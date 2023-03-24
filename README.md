# go-pdfium

[![Go Reference](https://pkg.go.dev/badge/github.com/klippa-app/go-pdfium/pdfium.svg)](https://pkg.go.dev/github.com/klippa-app/go-pdfium)
[![Build Status][build-status]][build-url]
[![codecov](https://codecov.io/gh/klippa-app/go-pdfium/branch/main/graph/badge.svg?token=WoIlW9RbfH)](https://codecov.io/gh/klippa-app/go-pdfium)

[build-status]:https://github.com/klippa-app/go-pdfium/workflows/Go/badge.svg

[build-url]:https://github.com/klippa-app/go-pdfium/actions

:rocket: *Easy to use PDF library using Go and PDFium* :rocket:

**A fast, multi-threaded and easy to use PDF library for Go applications.**

## Features

* Option between single-threaded, multi-threaded (through subprocesses) and WebAssembly (which can do multithreading
  within go), while keeping the same interface
* This library will handle all complicated cgo/WebAssembly gymnastics for you, no direct cgo/WebAssembly usage/knowledge
  required
* Implementation of all PDFium public API methods (including methods that are marked experimental), with some exceptions
* PDFium has methods to do the following:
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
    * Windows features (`FPDF_SetPrintMode`, `FPDF_RenderPage`)
    * Transformations (page boxes, clip paths)
    * Progressive rendering
    * Document loading through data availability (loading data as needed)
    * Struct trees
    * Page/Page object editing
    * Annotations
    * Form filling
* Methods that won't be implemented for now:
    * fpdf_sysfontinfo.h (probably too complicated)
    * Skia methods ([not in pre-built binaries](https://github.com/bblanchon/pdfium-binaries/issues/29))
    * XFA/v8 JS
      methods ([not in pre-built binaries due to build issues](https://github.com/bblanchon/pdfium-binaries/issues/62))
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
documents. Therefor this project could also be called a binding.

Please be aware that PDFium comes with the `Apache License 2.0` license.

### Single/Multi-threading

Since PDFium is [not a multithreaded C++ library](https://groups.google.com/g/pdfium/c/HeZSsM_KEUk), we can not directly
make it multithreaded by calling it from Go's subroutines.

To solve this, this library has 3 different implementations that you can use to call PDFium:

- single_threaded: call PDFium from the same process with CGO
- multi_threaded: call PDFium using multiple workers with CGO, implemented using go-plugin
- webassembly: call PDFium using WebAssembly with [Wazero runtime](https://wazero.io/), you can start multiple runtimes
  to get multi-threaded behaviour

Both `single_threaded` and `multi_threaded` requires PDFium to be installed on your machine for your platform during
compilation and runtime, it also requires CGO to work on the platform you're compiling to.

`webassembly` does not need any external dependencies and also does not require CGO to work. However, Wazero currently
only has compiler support for amd64 and arm64, meaning it will be using the interpreter on other architectures which
will be much, much slower.

All implementations use exactly the same interface, so there won't be any code changes for you to switch between them.

All implementations are thread/subroutine safe, this has been guaranteed by locking the instance that's doing your work
while it's doing PDFium operations. New operations will wait until the lock becomes available again.

**Be aware that PDFium could use quite some memory depending on the size of the PDF and the requests that you do, so be
aware of the amount of workers that you configure.**

## Implementations

### Single/Multi-threading through CGO

Single-threading in CGO works by directly calling the PDFium library from the same process. Single-threaded might be
preferred if the caller is managing the workers themselves and does not want the overhead of another process. Be aware
that since PDFium is C++, we can't handle segfaults caused by PDFium, which may cause your process to be killed. So
using this library in the multi-threaded way, with only 1 worker, can still have some benefits, since it can
automatically recover from things like segfaults.

For CGO we have implemented multi-threading
using [HashiCorp's Go Plugin System](https://github.com/hashicorp/go-plugin),
which allows us to launch separate PDFium worker processes, and then route the requests through the different workers.
This also makes it a bit more safe to use PDFium, as it's less likely to segfaults or corrupt your main Go application.
The Plugin system provides the communication between the processes using gRPC, however, when implementing this library,
you won't really see anything of that. From the outside it will look like normal Go code. The inter-process
communication does come with a cost as it has to serialize/deserialize input/output as it moves between the main process
and the PDFium workers.

#### Prerequisites (CGO)

To use this Go library, you will need the actual PDFium library to run it and have it available through pkgconfig.

##### Get the library

You can try to compile PDFium yourself, but you can also use pre-compiled binaries, for example
from: https://github.com/bblanchon/pdfium-binaries/releases

If you use a pre-compiled library, make sure to extract it somewhere logical, for example /opt/pdfium.

##### Configure pkg-config

Create/edit file `/usr/lib/pkgconfig/pdfium.pc`

```
prefix={path}
libdir={path}/lib
includedir={path}/include

Name: PDFium
Description: PDFium
Version: 5664
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

#### Getting started (CGO)

To get started, make sure that you create a separate package in your application that will start the worker.

The examples below can also be found in the examples folder.

##### Single-threaded through CGO

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

##### Multi-threaded through CGO

###### Worker package

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

###### Worker configuration

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

###### Get page count

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

###### Render a page

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

#### Experimental (CGO)

Some newer API's by PDFium are marked as experimental. We do have support for these functions, but because they are
prone to change
we do not compile with support for it by default. This is to keep support for most PDFium versions by default.

**If you call any API methods marked as experimental, the call will result in an error**

Adding support for the experimental API is quite easy. All you have to do is give the build tag `pdfium_experimental`
when running `go build` or `go run`, like this:
`go run -tags pdfium_experimental {package}/{file.go}` or `go build -tags pdfium_experimental {package}/{file.go}`.

We actively monitor PDFium API additions/changes/deletions and apply them in the code-base. When API methods become
non-experimental, we will make them available in the default configuration.

### WebAssembly

Recently we have added support for a non-cgo implementation using WebAssembly, we use
the [Wazero runtime](https://wazero.io/)
for running WebAssembly within Go. The comes with quite some advantages:

- WebAssembly is one binary for every platform, which means we can embed the WebAssembly version of PDFium in this
  repository
- Because we have `go:embed`, we can embed the WebAssembly binary inside the Go binary
- Because of this, you won't have to download and distribute PDFium yourself, making deployments simpler
- Wazero is pure Go, and thus runs on any platform where you can run Go on (a lot), which also allows you to run
  go-pdfium and PDFium on all of those platforms
- Because it's running in Go, we can just start up multiple Wazero runtimes to support processing multiple PDF files at
  once
- Because it's running in Go, you can directly access the memory in Wazero, for example to render a PDFium bitmap
  directly to a Go Image
- Because it's running in Go, method calls and responses don't have to travel over gRPC like with go-plugin, saving
  quite some time in encoding/decoding, which makes it about 2x as fast as the go-pugin multithreading approach
- Since PDFium is compiled to WebAssembly and runs inside the Wazero runtime, it basically runs in a sandbox:
    - No chance of crashing the Go process like with cgo
    - No access to other local resources in case of attacks on PDFium (disk, network, memory)
    - Full control over file access (you decide which folders Wazero exposes to PDFium, by default it exposes the whole
      disk)

Of course there are also some disadvantages:

- It's about 2x as slow as the full native cgo implementation (but about 2x as fast as the multi-threaded cgo go-plugin
  implementation)
- WebAssembly doesn't have an option to give memory back to the system, once it has been claimed it will there until you
  close the instance, this could be solved by not re-using instances
- Some platform specific quirks that have been implemented in PDFium (for example for Windows and MacOS) won't work
  because the WebAssembly build is compiled as Linux

Please be aware that Wazero comes with the `Apache License 2.0` license.

#### Path handling (WebAssembly)

Because you can tell Wazero which folders have to be mounted in WebAssembly, you have full control over the filesystem.

By default, go-pdfium will mount the full root disk in Wazero on non-Windows environments.
On Windows environments, go-pdfium will get the volume of the current working directory and mount that as the root.

You can change this behaviour by overwriting FSConfig in the pool setup.

All paths given to go-pdfium in WebAssembly mode have to be in POSIX style and have to be absolute, so for
example: `/home/user/Downloads/file.pdf`. If you have mounted `/home/user/`on the root, then the path you would have to
give is `/Downloads/file.pdf`, this is the same on Windows, so no backward slashes or volume names in paths.

You can set your own mounts by overwriting FSConfig in the pool setup.

#### Getting started (WebAssembly)

The examples below can also be found in the examples folder.

To start go-pdfium workers, you will have to init the go-pdfium worker pool somewhere, this also allows you to
dynamically start
workers when needed. The best location to add this is in the `init()` of a package that is going to call the PDFium
library. Example:

`pdfium/renderer/renderer.go`

```go
package renderer

import (
	"log"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/webassembly"
)

// Be sure to close pools/instances when you're done with them.
var pool pdfium.Pool
var instance pdfium.Pdfium

func init() {
	// Init the PDFium library and return the instance to open documents.
	// You can tweak these configs to your need. Be aware that workers can use quite some memory.
	pool, err = webassembly.Init(webassembly.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
	})

	var err error
	instance, err = pool.GetInstance(time.Second * 30)
	if err != nil {
		log.Fatal(err)
	}
}
```

##### Get page count

```go
package renderer

import (
	"io/ioutil"
	"log"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
)

// Insert the webassembly init() here.

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

#### Render a page

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

// Insert the webassembly init() here.

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

	// The Render* methods return a cleanup function that has to be called when
	// using webassembly to make sure resources are cleaned up. Do this after
	// you are done with the returned image object.
	defer pageRender.Cleanup()

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

#### Experimental (WebAssembly)

Some newer API's by PDFium are marked as experimental. The WebAssembly build has support for all of them.

We actively monitor PDFium API additions/changes/deletions and apply them in the code-base.

The WebAssembly build will always be the latest PDFium version that we added support for.

## `io.ReadSeeker` and `io.Writer`

Document loading allows you to load a document with a `io.ReadSeeker`. Please be aware that this only works efficiently
when using the single-threaded or WebAssembly usage, as that lives in the same process. For multi-threaded usage this
will just load in
the complete file and pass the bytes through the gRPC interface.

Document/image saving allows you to save using a `io.Writer`. Please be aware this only works when using the
single-threaded or WebAssembly usage. It's not possible to encode the `io.Writer` with gRPC. Or share it between
processes for that
matter.

## About Klippa

Founded in 2015, [Klippa](https://www.klippa.com/en)'s goal is to digitize & automate administrative processes with
modern technologies. We help clients enhance the effectiveness of their organization by using machine learning and OCR.
Since 2015, more than a thousand happy clients have used Klippa's software solutions. Klippa currently has an
international team of 50 people, with offices in Groningen, Amsterdam and Brasov.

## License

The MIT License (MIT)
