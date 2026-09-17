package textextract

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTextExtract(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Text Extract Suite")
}

// box places a 10-wide, 10-high char at x = 10*i on line y (baseline/origin y,
// box from y to y+10).
func box(i int, y float32, r rune) Char {
	x := float32(10 * i)
	return Char{Left: x, Right: x + 10, Bottom: y, Top: y + 10, OriginY: y, Unicode: r}
}

var _ = Describe("TextInRect", func() {
	// Two lines: "AB CD" at y=100, "EF" at y=80.
	twoLines := []Char{
		box(0, 100, 'A'),
		box(1, 100, 'B'),
		box(2, 100, ' '),
		box(3, 100, 'C'),
		box(4, 100, 'D'),
		box(0, 80, 'E'),
		box(1, 80, 'F'),
	}

	DescribeTable("returns the text in the rect",
		func(left, top, right, bottom float32, want string) {
			Expect(New(twoLines).TextInRect(left, top, right, bottom)).To(Equal(want))
		},
		// All chars matched contiguously: IsAddLineFeed is only ever set by
		// unmatched non-space chars, so no \r\n appears. (On real pages,
		// PDFium's generated newline chars have empty boxes, never match any
		// rect, and drive that branch.)
		Entry("everything", float32(0), float32(120), float32(100), float32(0), "AB CDEF"),
		Entry("first line", float32(0), float32(115), float32(100), float32(95), "AB CD"),
		Entry("second line", float32(0), float32(95), float32(100), float32(75), "EF"),
		Entry("AB with trailing unmatched space", float32(0), float32(115), float32(20), float32(95), "AB "),
		Entry("CD only, gap starts with space", float32(30), float32(115), float32(100), float32(95), "CD"),
		Entry("A only, B touches at zero width", float32(0), float32(115), float32(10), float32(95), "A"),
		Entry("touching edge is not intersecting", float32(0), float32(100), float32(100), float32(100), ""),
		Entry("empty rect", float32(50), float32(50), float32(50), float32(50), ""),
		Entry("un-normalized rect input", float32(100), float32(95), float32(0), float32(115), "AB CD"),
	)

	It("adds a line feed between lines", func() {
		// A rect covering chars on two lines with a non-space, unmatched char
		// in between must produce \r\n (PDFium's IsAddLineFeed path).
		chars := []Char{
			box(0, 100, 'A'),
			box(9, 100, 'X'), // outside the query rect, non-space
			box(0, 80, 'B'),
		}
		Expect(New(chars).TextInRect(0, 115, 20, 75)).To(Equal("A\r\nB"))
	})

	It("appends nothing for chars without unicode", func() {
		chars := []Char{
			box(0, 100, 'A'),
			box(1, 100, 0), // generated char without unicode
			box(2, 100, 'B'),
		}
		Expect(New(chars).TextInRect(0, 115, 100, 95)).To(Equal("AB"))
	})

	It("returns nothing for an empty page", func() {
		Expect(New(nil).TextInRect(0, 100, 100, 0)).To(Equal(""))
	})

	// A query rect that lies entirely outside the page content must return ""
	// and must not panic. Rects entirely above the content used to index one
	// bucket past the end: bucketRange clamped hi but not lo, and its hi < lo
	// fixup then handed lo straight to the bucket loop. Not reachable through
	// GetPageTextStructured (PDFium's own rects always sit inside the char
	// bounds), so only reachable by calling TextInRect directly.
	Describe("a query rect outside the content", func() {
		var e *Extractor

		BeforeEach(func() {
			// Enough chars that New() allocates more than one y-bucket.
			chars := make([]Char, 0, 64)
			for i := range 16 {
				for line := range 4 {
					chars = append(chars, box(i, float32(100-line*20), 'A'+rune(i%26)))
				}
			}
			e = New(chars)
		})

		DescribeTable("returns nothing",
			func(left, top, right, bottom float32) {
				Expect(e.TextInRect(left, top, right, bottom)).To(Equal(""))
			},
			Entry("far above content", float32(0), float32(10000), float32(1000), float32(9000)),
			Entry("just above content", float32(0), float32(200), float32(1000), float32(111)),
			Entry("far below content", float32(0), float32(-9000), float32(1000), float32(-10000)),
			Entry("just below content", float32(0), float32(39), float32(1000), float32(-100)),
			Entry("left of content", float32(-1000), float32(120), float32(-500), float32(0)),
			Entry("right of content", float32(5000), float32(120), float32(6000), float32(0)),
		)
	})
})

var _ = Describe("FirstCharIndices", func() {
	It("finds the first char of each rect", func() {
		// Line 1 at y=100: "ab" in one text object where 'b' is taller than
		// 'a', then a generated space (zero box), then "CD" in a second text
		// object. Line 2 at y=80: "ef".
		tall := box(1, 100, 'b')
		tall.Top = 115
		chars := []Char{
			box(0, 100, 'a'),
			tall,
			{Left: 20, Right: 20, Bottom: 100, Top: 100, OriginY: 100, Unicode: ' '},
			box(3, 100, 'C'),
			box(4, 100, 'D'),
			box(0, 80, 'e'),
			box(1, 80, 'f'),
		}

		// Rects as PDFium's GetRectArray would produce them, in order. The
		// corner (0,115) of the first rect is above 'a', which is why a
		// position lookup at the corner fails; the last rect has no chars at
		// all.
		rects := []Rect{{0, 115, 20, 100}, {30, 110, 50, 100}, {0, 90, 20, 80}, {0, 70, 20, 60}}
		Expect(New(chars).FirstCharIndices(rects)).To(Equal([]int{0, 3, 5, -1}))

		// Un-normalized rect coordinates are accepted.
		Expect(New(chars).FirstCharIndices([]Rect{{20, 100, 0, 115}})).To(Equal([]int{0}))
	})

	It("resolves a rect nested inside the previous one", func() {
		// A second text object drawn inside the box of the first one (e.g.
		// text re-drawn over an existing, larger word). The first rect's run
		// swallows the inner chars; the inner rect must still resolve to its
		// own first char.
		chars := []Char{
			{Left: 0, Right: 100, Bottom: 0, Top: 20, Unicode: 'W'},
			{Left: 100, Right: 200, Bottom: 0, Top: 20, Unicode: 'X'},
			{Left: 10, Right: 20, Bottom: 5, Top: 15, Unicode: 'i'},
			{Left: 20, Right: 30, Bottom: 5, Top: 15, Unicode: 'n'},
			{Left: 0, Right: 10, Bottom: 40, Top: 50, Unicode: 'z'},
		}
		Expect(New(chars).FirstCharIndices([]Rect{{0, 20, 200, 0}, {10, 15, 30, 5}, {0, 50, 10, 40}})).To(Equal([]int{0, 2, 4}))
	})

	It("ends a run at the first char of an overdrawn next rect", func() {
		// Real-world pattern: "werk" is drawn, then a second text object
		// redraws "werk deel" starting inside the first rect. The redrawn
		// "werk" lies inside the first rect and would be swallowed by a plain
		// walk; the union of the first rect is complete after its own "werk",
		// so the redrawn 'w', which also lies in the second rect, starts the
		// second run.
		chars := []Char{
			box(0, 100, 'w'), box(1, 100, 'e'), box(2, 100, 'r'), box(3, 100, 'k'), // rect 0: x 0..40
			box(1, 100, 'w'), box(2, 100, 'e'), box(3, 100, 'r'), box(4, 100, 'k'), // rect 1: x 10..70
			box(5, 100, 'd'), box(6, 100, 'e'),
		}
		Expect(New(chars).FirstCharIndices([]Rect{{0, 110, 40, 100}, {10, 110, 70, 100}})).To(Equal([]int{0, 4}))
	})
})
