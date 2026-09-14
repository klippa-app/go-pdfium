package implementation_cgo

import (
	"unsafe"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// The timer callback handed to FFI_SetTimer is the only callback a caller
// invokes itself, from a goroutine of its own choosing, so it is the only one
// that has to deal with PDFium being busy or already gone.
var _ = Describe("a form fill timer tick", func() {
	var instance *PdfiumImplementation

	BeforeEach(func() {
		instance = Pdfium.GetInstance()
	})

	AfterEach(func() {
		Expect(instance.Close()).To(BeNil())
	})

	// A form fill environment is identified by the address of its
	// FPDF_FORMFILLINFO; a stand-in is enough here, formFillTimerTick only
	// uses it to look the environment up.
	newFormInfo := func() (unsafe.Pointer, func()) {
		formInfo := unsafe.Pointer(new([8]byte))
		formFillInfoHandles[formInfo] = &FormFillInfo{}

		return formInfo, func() { delete(formFillInfoHandles, formInfo) }
	}

	It("runs while the environment exists", func() {
		formInfo, exit := newFormInfo()
		defer exit()

		ticked := false
		formFillTimerTick(instance, formInfo, func() { ticked = true })
		Expect(ticked).To(BeTrue())
	})

	It("does not run after the environment was destroyed", func() {
		formInfo, exit := newFormInfo()
		exit()

		ticked := false
		formFillTimerTick(instance, formInfo, func() { ticked = true })
		Expect(ticked).To(BeFalse())
	})

	It("gives up instead of waiting when PDFium is busy", func() {
		formInfo, exit := newFormInfo()

		ticked := false
		returned := make(chan struct{})

		// PDFium kills timers from inside its own API calls, so a caller
		// that stops its timer goroutine from FFI_KillTimer and waits for it
		// would deadlock if a tick could block here.
		instance.Lock()
		go func() {
			defer close(returned)
			formFillTimerTick(instance, formInfo, func() { ticked = true })
		}()

		// Unwind whatever the assertions do: the lock is process wide, so
		// leaving it held would wedge the rest of the suite, and the
		// registry entry has to outlive a tick that did block on it.
		defer func() {
			instance.Unlock()
			<-returned
			exit()
		}()

		Eventually(returned).Should(BeClosed())
		Expect(ticked).To(BeFalse())
	})
})
