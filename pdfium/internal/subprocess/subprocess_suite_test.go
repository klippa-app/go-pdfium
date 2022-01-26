package subprocess_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSubprocess(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Subprocess Suite")
}
