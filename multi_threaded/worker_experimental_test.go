//go:build pdfium_experimental

package multi_threaded_test

// workerBuildTags are the build tags that the worker has to be compiled with.
// The worker runs in a separate process, so it has to be built with the same
// tags as this test binary, otherwise it offers a different set of methods than
// the tests expect. Deriving the tags from the build tags of this file means
// that the two can never disagree.
var workerBuildTags = []string{"-tags", "pdfium_experimental"}
