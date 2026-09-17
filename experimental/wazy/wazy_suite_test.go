package wazy_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPdfiumWazy(t *testing.T) {
	RegisterFailHandler(Fail)

	suiteConfig, reporterConfig := GinkgoConfiguration()

	// Kill() closes the module while the interrupted call may still be
	// unwinding through Go code that reads guest memory. wazy's Module.Close
	// clears the memory buffer under a mutex that its reads do not take, so
	// the race detector reports that (harmless, the read fails) as a data
	// race inside wazy. Skip the one test that provokes it under -race until
	// wazy tolerates reads of a closed memory the way wazero does.
	if raceDetectorEnabled {
		suiteConfig.SkipStrings = append(suiteConfig.SkipStrings, "interrupts a stuck WASM execution")
	}

	suiteDescription := "Wazy Suite"
	if interpreterMode {
		suiteDescription = "Wazy Interpreter Suite"
	}
	RunSpecs(t, suiteDescription, suiteConfig, reporterConfig)
}
