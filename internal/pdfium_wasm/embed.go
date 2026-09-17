// Package pdfium_wasm embeds the PDFium WebAssembly module that ships with
// go-pdfium, so that every WebAssembly runtime backend can use the same
// module without embedding it more than once.
package pdfium_wasm

import _ "embed"

// Module is the PDFium WebAssembly module. It is only linked into a binary
// when a backend's Init function that references it is linked in.
//
//go:embed pdfium.wasm
var Module []byte
