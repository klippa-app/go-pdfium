package implementation_cgo

import (
	"os"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("an availability provider", func() {
	var instance *PdfiumImplementation
	var file *os.File
	var fileSize int64

	BeforeEach(func() {
		instance = Pdfium.GetInstance()

		openedFile, err := os.Open("../../shared_tests/testdata/test.pdf")
		Expect(err).To(BeNil())
		file = openedFile

		stat, err := file.Stat()
		Expect(err).To(BeNil())
		fileSize = stat.Size()
	})

	AfterEach(func() {
		Expect(instance.Close()).To(BeNil())
		Expect(file.Close()).To(BeNil())
	})

	create := func() *requests.FPDFAvail_Destroy {
		FPDFAvail_Create, err := instance.FPDFAvail_Create(&requests.FPDFAvail_Create{
			Reader:                  file,
			Size:                    fileSize,
			IsDataAvailableCallback: func(offset, size uint64) bool { return true },
			AddSegmentCallback:      func(offset, size uint64) {},
		})
		Expect(err).To(BeNil())

		return &requests.FPDFAvail_Destroy{
			AvailabilityProvider: FPDFAvail_Create.AvailabilityProvider,
		}
	}

	// FPDFAvail_Create registers its callbacks in package level maps, keyed by
	// the address of the C struct that PDFium was given. FPDFAvail_Destroy has
	// to remove them again with those same addresses, or every availability
	// provider leaks its entry (and the callback it closes over) for the
	// lifetime of the process.
	When("it is destroyed", func() {
		It("removes the callbacks that were registered for it", func() {
			availCallbacks := len(dataAvailAbilityCallbacks)
			segmentCallbacks := len(addSegmentCallbacks)

			destroy := create()
			Expect(dataAvailAbilityCallbacks).To(HaveLen(availCallbacks + 1))
			Expect(addSegmentCallbacks).To(HaveLen(segmentCallbacks + 1))

			FPDFAvail_Destroy, err := instance.FPDFAvail_Destroy(destroy)
			Expect(err).To(BeNil())
			Expect(FPDFAvail_Destroy).To(Not(BeNil()))

			Expect(dataAvailAbilityCallbacks).To(HaveLen(availCallbacks))
			Expect(addSegmentCallbacks).To(HaveLen(segmentCallbacks))
		})
	})

	// PDFium parses the document of FPDFAvail_GetDocument on top of the
	// provider's FPDF_FILEACCESS and keeps reading through it, so the reader
	// has to survive FPDFAvail_Destroy when such a document is still open.
	// Freeing it early leaves PDFium reading a freed struct, and calling the
	// m_GetBlock function pointer it finds in it.
	When("it is destroyed while a document parsed from it is still open", func() {
		It("keeps the reader alive until the document is closed", func() {
			destroy := create()

			FPDFAvail_IsDocAvail, err := instance.FPDFAvail_IsDocAvail(&requests.FPDFAvail_IsDocAvail{
				AvailabilityProvider: destroy.AvailabilityProvider,
			})
			Expect(err).To(BeNil())
			Expect(FPDFAvail_IsDocAvail.IsDocAvail).To(Equal(enums.PDF_FILEAVAIL_DATA_AVAIL))

			FPDFAvail_GetDocument, err := instance.FPDFAvail_GetDocument(&requests.FPDFAvail_GetDocument{
				AvailabilityProvider: destroy.AvailabilityProvider,
			})
			Expect(err).To(BeNil())

			dataAvailHandle, err := instance.getDataAvailHandle(destroy.AvailabilityProvider)
			Expect(err).To(BeNil())
			readerRef := dataAvailHandle.fileHandleRef
			Expect(Pdfium.fileReaders).To(HaveKey(readerRef))

			_, err = instance.FPDFAvail_Destroy(destroy)
			Expect(err).To(BeNil())
			Expect(Pdfium.fileReaders).To(HaveKey(readerRef))

			// Loading a page reads parts of the file that parsing the
			// document did not need yet, so this goes through the reader
			// again after the provider is gone.
			FPDFPage_GetRotation, err := instance.FPDFPage_GetRotation(&requests.FPDFPage_GetRotation{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: FPDFAvail_GetDocument.Document,
						Index:    0,
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(FPDFPage_GetRotation.PageRotation).To(Equal(enums.FPDF_PAGE_ROTATION_NONE))

			FPDF_CloseDocument, err := instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
				Document: FPDFAvail_GetDocument.Document,
			})
			Expect(err).To(BeNil())
			Expect(FPDF_CloseDocument).To(Not(BeNil()))

			Expect(Pdfium.fileReaders).To(Not(HaveKey(readerRef)))
		})
	})
})
