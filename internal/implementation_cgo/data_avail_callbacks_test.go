package implementation_cgo

import (
	"os"

	"github.com/klippa-app/go-pdfium/requests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// FPDFAvail_Create registers its callbacks in package level maps, keyed by the
// address of the C struct that PDFium was given. FPDFAvail_Destroy has to
// remove them again with those same addresses, or every availability provider
// leaks its entry (and the callback it closes over) for the lifetime of the
// process.
var _ = Describe("FPDFAvail_Create", func() {
	When("the availability provider is destroyed", func() {
		It("removes the callbacks that were registered for it", func() {
			file, err := os.Open("../../shared_tests/testdata/test.pdf")
			Expect(err).To(BeNil())
			defer file.Close()

			stat, err := file.Stat()
			Expect(err).To(BeNil())

			instance := Pdfium.GetInstance()
			defer func() { Expect(instance.Close()).To(BeNil()) }()

			availCallbacks := len(dataAvailAbilityCallbacks)
			segmentCallbacks := len(addSegmentCallbacks)

			FPDFAvail_Create, err := instance.FPDFAvail_Create(&requests.FPDFAvail_Create{
				Reader:                  file,
				Size:                    stat.Size(),
				IsDataAvailableCallback: func(offset, size uint64) bool { return true },
				AddSegmentCallback:      func(offset, size uint64) {},
			})
			Expect(err).To(BeNil())
			Expect(dataAvailAbilityCallbacks).To(HaveLen(availCallbacks + 1))
			Expect(addSegmentCallbacks).To(HaveLen(segmentCallbacks + 1))

			FPDFAvail_Destroy, err := instance.FPDFAvail_Destroy(&requests.FPDFAvail_Destroy{
				AvailabilityProvider: FPDFAvail_Create.AvailabilityProvider,
			})
			Expect(err).To(BeNil())
			Expect(FPDFAvail_Destroy).To(Not(BeNil()))

			Expect(dataAvailAbilityCallbacks).To(HaveLen(availCallbacks))
			Expect(addSegmentCallbacks).To(HaveLen(segmentCallbacks))
		})
	})
})
