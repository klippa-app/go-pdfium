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
