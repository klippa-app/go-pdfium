package shared_tests

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

// Regression tests for the rect font lookup in GetPageTextStructured. The rect
// font used to be looked up with FPDFText_GetCharIndexAtPos at the rect's
// top-left corner with a fixed tolerance, which returned no char (zero-valued
// font information) or the font of an unrelated overlapping char for about one
// in ten rects. All chars of a rect share one text object, so the rect's font
// must equal the font of every char whose box lies inside the rect.
var _ = Describe("text rect font information", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	// text_font.pdf and text_render_mode.pdf have a single rect whose corner
	// the position lookup could not resolve at all; alpha_channel.pdf has
	// overlapping text objects where the lookup picked another object's char;
	// rect-wrong.pdf is a real-world page with many rects of both kinds.
	for _, file := range []string{"text_font.pdf", "text_render_mode.pdf", "alpha_channel.pdf", "rect-wrong.pdf"} {
		file := file
		It("returns the font of the rect's own chars for every rect in "+file, func() {
			pdfData, err := os.ReadFile(TestDataPath + "/testdata/" + file)
			Expect(err).To(BeNil())

			doc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{
				Data: &pdfData,
			})
			Expect(err).To(BeNil())
			defer PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})

			pageText, err := PdfiumInstance.GetPageTextStructured(&requests.GetPageTextStructured{
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc.Document,
						Index:    0,
					},
				},
				Mode:                   requests.GetPageTextStructuredModeBoth,
				CollectFontInformation: true,
			})
			Expect(err).To(BeNil())
			Expect(pageText.Rects).ToNot(BeEmpty())

			for i, rect := range pageText.Rects {
				Expect(rect.FontInformation).ToNot(BeNil(), "rect %d has no font information", i)

				// Every non-degenerate char inside the rect belongs to the rect
				// and must carry the same font.
				matched := 0
				for _, char := range charsInRect(pageText.Chars, rect.PointPosition) {
					matched++
					Expect(char.FontInformation).To(Equal(rect.FontInformation), "rect %d %q: char %q has a different font", i, rect.Text, char.Text)
				}
				Expect(matched).ToNot(BeZero(), "rect %d %q contains no chars", i, rect.Text)
			}
		})
	}
})

func charsInRect(chars []*responses.GetPageTextStructuredChar, rect responses.CharPosition) []*responses.GetPageTextStructuredChar {
	var result []*responses.GetPageTextStructuredChar
	for _, char := range chars {
		box := char.PointPosition
		if box.Left > box.Right {
			box.Left, box.Right = box.Right, box.Left
		}
		if box.Bottom > box.Top {
			box.Bottom, box.Top = box.Top, box.Bottom
		}
		// PDFium skips chars with a box smaller than 0.01 in either direction
		// when building rects, generated spaces and newlines among them.
		if box.Right-box.Left < 0.01 || box.Top-box.Bottom < 0.01 {
			continue
		}
		if box.Left >= rect.Left && box.Right <= rect.Right && box.Bottom >= rect.Bottom && box.Top <= rect.Top {
			result = append(result, char)
		}
	}
	return result
}
