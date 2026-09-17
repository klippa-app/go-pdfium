package image_jpeg

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestImageJPEG(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Image JPEG Suite")
}
