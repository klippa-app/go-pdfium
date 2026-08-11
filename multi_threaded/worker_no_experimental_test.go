//go:build !pdfium_experimental

package multi_threaded_test

// workerBuildTags are the build tags that the worker has to be compiled with.
// See the experimental variant of this file for why this is derived from the
// build tags instead of being configured separately.
var workerBuildTags []string
