package textextract

import "testing"

// box places a 10-wide, 10-high char at x = 10*i on line y (baseline/origin y,
// box from y to y+10).
func box(i int, y float32, r rune) Char {
	x := float32(10 * i)
	return Char{Left: x, Right: x + 10, Bottom: y, Top: y + 10, OriginY: y, Unicode: r}
}

func TestTextInRect(t *testing.T) {
	// Two lines: "AB CD" at y=100, "EF" at y=80.
	chars := []Char{
		box(0, 100, 'A'),
		box(1, 100, 'B'),
		box(2, 100, ' '),
		box(3, 100, 'C'),
		box(4, 100, 'D'),
		box(0, 80, 'E'),
		box(1, 80, 'F'),
	}
	e := New(chars)

	cases := []struct {
		name                     string
		left, top, right, bottom float32
		want                     string
	}{
		// All chars matched contiguously: IsAddLineFeed is only ever set by
		// unmatched non-space chars, so no \r\n appears. (On real pages,
		// PDFium's generated newline chars have empty boxes, never match any
		// rect, and drive that branch.)
		{"everything", 0, 120, 100, 0, "AB CDEF"},
		{"first line", 0, 115, 100, 95, "AB CD"},
		{"second line", 0, 95, 100, 75, "EF"},
		{"AB with trailing unmatched space", 0, 115, 20, 95, "AB "},
		{"CD only, gap starts with space", 30, 115, 100, 95, "CD"},
		{"A only, B touches at zero width", 0, 115, 10, 95, "A"},
		{"touching edge is not intersecting", 0, 100, 100, 100, ""},
		{"empty rect", 50, 50, 50, 50, ""},
		{"un-normalized rect input", 100, 95, 0, 115, "AB CD"},
	}
	for _, c := range cases {
		if got := e.TextInRect(c.left, c.top, c.right, c.bottom); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}

func TestLineFeedBetweenLines(t *testing.T) {
	// A rect covering chars on two lines with a non-space, unmatched char in
	// between must produce \r\n (PDFium's IsAddLineFeed path).
	chars := []Char{
		box(0, 100, 'A'),
		box(9, 100, 'X'), // outside the query rect, non-space
		box(0, 80, 'B'),
	}
	e := New(chars)
	if got := e.TextInRect(0, 115, 20, 75); got != "A\r\nB" {
		t.Errorf("got %q, want %q", got, "A\r\nB")
	}
}

func TestUnicodeZeroCharsAppendNothing(t *testing.T) {
	chars := []Char{
		box(0, 100, 'A'),
		box(1, 100, 0), // generated char without unicode
		box(2, 100, 'B'),
	}
	e := New(chars)
	if got := e.TextInRect(0, 115, 100, 95); got != "AB" {
		t.Errorf("got %q, want %q", got, "AB")
	}
}

func TestEmptyPage(t *testing.T) {
	e := New(nil)
	if got := e.TextInRect(0, 100, 100, 0); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

// A query rect that lies entirely outside the page content must return "" and
// must not panic. Rects entirely above the content used to index one bucket
// past the end: bucketRange clamped hi but not lo, and its hi < lo fixup then
// handed lo straight to the bucket loop. Not reachable through
// GetPageTextStructured (PDFium's own rects always sit inside the char bounds),
// so only reachable by calling TextInRect directly.
func TestQueryRectOutsideContent(t *testing.T) {
	// Enough chars that New() allocates more than one y-bucket.
	chars := make([]Char, 0, 64)
	for i := range 16 {
		for line := range 4 {
			chars = append(chars, box(i, float32(100-line*20), 'A'+rune(i%26)))
		}
	}
	e := New(chars)

	cases := []struct {
		name                     string
		left, top, right, bottom float32
	}{
		{"far above content", 0, 10000, 1000, 9000},
		{"just above content", 0, 200, 1000, 111},
		{"far below content", 0, -9000, 1000, -10000},
		{"just below content", 0, 39, 1000, -100},
		{"left of content", -1000, 120, -500, 0},
		{"right of content", 5000, 120, 6000, 0},
	}
	for _, c := range cases {
		if got := e.TextInRect(c.left, c.top, c.right, c.bottom); got != "" {
			t.Errorf("%s: got %q, want empty", c.name, got)
		}
	}
}

func TestFirstCharIndices(t *testing.T) {
	// Line 1 at y=100: "ab" in one text object where 'b' is taller than 'a',
	// then a generated space (zero box), then "CD" in a second text object.
	// Line 2 at y=80: "ef".
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

	// Rects as PDFium's GetRectArray would produce them, in order. The corner
	// (0,115) of the first rect is above 'a', which is why a position lookup
	// at the corner fails; the last rect has no chars at all.
	rects := []Rect{{0, 115, 20, 100}, {30, 110, 50, 100}, {0, 90, 20, 80}, {0, 70, 20, 60}}
	want := []int{0, 3, 5, -1}
	got := New(chars).FirstCharIndices(rects)
	if !equalInts(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	// Un-normalized rect coordinates are accepted.
	got = New(chars).FirstCharIndices([]Rect{{20, 100, 0, 115}})
	if !equalInts(got, []int{0}) {
		t.Errorf("un-normalized rect: got %v, want [0]", got)
	}
}

func TestFirstCharIndicesNestedRects(t *testing.T) {
	// A second text object drawn inside the box of the first one (e.g. text
	// re-drawn over an existing, larger word). The first rect's run swallows
	// the inner chars; the inner rect must still resolve to its own first char.
	chars := []Char{
		{Left: 0, Right: 100, Bottom: 0, Top: 20, Unicode: 'W'},
		{Left: 100, Right: 200, Bottom: 0, Top: 20, Unicode: 'X'},
		{Left: 10, Right: 20, Bottom: 5, Top: 15, Unicode: 'i'},
		{Left: 20, Right: 30, Bottom: 5, Top: 15, Unicode: 'n'},
		{Left: 0, Right: 10, Bottom: 40, Top: 50, Unicode: 'z'},
	}
	got := New(chars).FirstCharIndices([]Rect{{0, 20, 200, 0}, {10, 15, 30, 5}, {0, 50, 10, 40}})
	if !equalInts(got, []int{0, 2, 4}) {
		t.Errorf("got %v, want [0 2 4]", got)
	}
}

func TestFirstCharIndicesOverdrawnRun(t *testing.T) {
	// Real-world pattern: "werk" is drawn, then a second text object redraws
	// "werk deel" starting inside the first rect. The redrawn "werk" lies
	// inside the first rect and would be swallowed by a plain walk; the union
	// of the first rect is complete after its own "werk", so the redrawn 'w',
	// which also lies in the second rect, starts the second run.
	chars := []Char{
		box(0, 100, 'w'), box(1, 100, 'e'), box(2, 100, 'r'), box(3, 100, 'k'), // rect 0: x 0..40
		box(1, 100, 'w'), box(2, 100, 'e'), box(3, 100, 'r'), box(4, 100, 'k'), // rect 1: x 10..70
		box(5, 100, 'd'), box(6, 100, 'e'),
	}
	got := New(chars).FirstCharIndices([]Rect{{0, 110, 40, 100}, {10, 110, 70, 100}})
	if !equalInts(got, []int{0, 4}) {
		t.Errorf("got %v, want [0 4]", got)
	}
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
