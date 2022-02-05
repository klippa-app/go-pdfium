package implementation_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestImplementation(t *testing.T) {
	// Set ENV to ensure resulting values.
	os.Setenv("TZ", "UTC")

	RegisterFailHandler(Fail)
	RunSpecs(t, "Implementation Suite")
}
