package implementation_cgo_test

import (
	"runtime"
	"unsafe"

	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/shared_tests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// FPDFBitmap_CreateEx hands PDFium a raw pointer into caller memory and PDFium
// keeps writing through it for the lifetime of the bitmap. When that memory is
// a Go allocation, the bitmap handle has to keep it reachable: if the garbage
// collector frees and reuses it, the next render writes pixels straight into
// whatever Go objects now live there, and the process dies much later in an
// unrelated place ("found pointer to free object", "index out of range" inside
// the runtime).
//
// The buffer can reach the implementation either as a []byte (request.Buffer)
// or as a bare unsafe.Pointer (request.Pointer); both are covered here,
// because only the former was retained at first.
var _ = Describe("FPDFBitmap_CreateEx", func() {
	const width = 8
	const height = 8
	const stride = width * 4
	const size = stride * height

	// createBitmap allocates a buffer, hands it to PDFium, and returns the
	// bitmap plus a channel that is closed once the buffer is collected. The
	// buffer deliberately does not escape: after this returns, the only thing
	// that can keep it alive is the bitmap handle.
	createBitmap := func(useRawPointer bool) (references.FPDF_BITMAP, chan struct{}) {
		collected := make(chan struct{})
		buffer := new([size]byte)
		runtime.SetFinalizer(buffer, func(*[size]byte) { close(collected) })

		request := &requests.FPDFBitmap_CreateEx{
			Width:  width,
			Height: height,
			Format: enums.FPDF_BITMAP_FORMAT_BGRA,
			Stride: stride,
		}

		if useRawPointer {
			request.Pointer = unsafe.Pointer(&buffer[0])
		} else {
			request.Buffer = buffer[:]
		}

		resp, err := shared_tests.PdfiumInstance.FPDFBitmap_CreateEx(request)
		Expect(err).To(BeNil())
		Expect(resp).To(Not(BeNil()))

		return resp.Bitmap, collected
	}

	for _, variant := range []struct {
		name          string
		useRawPointer bool
	}{
		{name: "a buffer", useRawPointer: false},
		{name: "a raw pointer", useRawPointer: true},
	} {
		variant := variant

		When("given an external bitmap buffer as "+variant.name, func() {
			It("keeps the Go memory alive until the bitmap is destroyed", func() {
				bitmap, collected := createBitmap(variant.useRawPointer)

				// Nothing in the test holds the buffer anymore. Only the
				// bitmap handle can keep it out of the sweeper's hands.
				runtime.GC()
				runtime.GC()
				Consistently(collected).ShouldNot(BeClosed())

				// PDFium writes the whole buffer here. Without a live
				// reference this is the write that corrupts the Go heap.
				FPDFBitmap_FillRect, err := shared_tests.PdfiumInstance.FPDFBitmap_FillRect(&requests.FPDFBitmap_FillRect{
					Bitmap: bitmap,
					Color:  0xFFFFFFFF,
					Left:   0,
					Top:    0,
					Width:  width,
					Height: height,
				})
				Expect(err).To(BeNil())
				Expect(FPDFBitmap_FillRect).To(Not(BeNil()))

				FPDFBitmap_Destroy, err := shared_tests.PdfiumInstance.FPDFBitmap_Destroy(&requests.FPDFBitmap_Destroy{
					Bitmap: bitmap,
				})
				Expect(err).To(BeNil())
				Expect(FPDFBitmap_Destroy).To(Not(BeNil()))

				// And it must not be pinned for longer than that either.
				Eventually(func() bool {
					runtime.GC()
					select {
					case <-collected:
						return true
					default:
						return false
					}
				}).Should(BeTrue())
			})
		})
	}
})
