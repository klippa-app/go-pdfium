// This go.mod keeps shared_tests/testdata out of the module zip that users
// download: Go leaves out directories that contain their own go.mod.
// The tests read these files from the repository checkout and are not affected.
// See https://go.dev/ref/mod#zip-path-size-constraints and issue #187.
module github.com/klippa-app/go-pdfium/shared_tests/testdata
