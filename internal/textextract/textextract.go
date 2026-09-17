// Package textextract reimplements PDFium's CPDF_TextPage::GetTextByRect()
// on the Go side. PDFium's FPDFText_GetBoundedText scans every character on
// the page for every call, which makes extracting the text of all rects on a
// page O(chars × rects). By collecting the per-character data once (char box,
// origin and unicode, all lossless float32 round-trips through the C API) and
// indexing the boxes by their vertical position, the text of each rect is
// computed in O(matched chars) instead.
//
// The output is bit-identical to FPDFText_GetBoundedText: the intersection
// test mirrors CFX_FloatRect::Intersect/IsEmpty (strict, on normalized rects,
// in float32 math) and the line-feed/space state machine is a direct port of
// CPDF_TextPage::GetTextByPredicate.
package textextract

import "slices"

// Char is the per-character data GetTextByRect operates on, in PDFium's own
// float32 precision. Box coordinates may be un-normalized, as PDFium stores
// them.
type Char struct {
	Left, Bottom, Right, Top float32
	OriginY                  float32
	Unicode                  rune
}

const maxBuckets = 2048

// Extractor answers GetTextByRect queries over a fixed list of page chars.
type Extractor struct {
	chars []Char

	// normalized boxes, indexed like chars
	nl, nb, nr, nt []float32

	// nonSpacePrefix[i] = number of chars in [0, i) whose Unicode != ' '.
	nonSpacePrefix []int32

	// y-bucket index over the normalized boxes
	buckets  [][]int32
	minY     float32
	bucketH  float32
	seen     []int32
	epoch    int32
	matchBuf []int32
}

func New(chars []Char) *Extractor {
	e := &Extractor{
		chars:          chars,
		nl:             make([]float32, len(chars)),
		nb:             make([]float32, len(chars)),
		nr:             make([]float32, len(chars)),
		nt:             make([]float32, len(chars)),
		nonSpacePrefix: make([]int32, len(chars)+1),
		seen:           make([]int32, len(chars)),
	}

	minY, maxY := float32(0), float32(0)
	for i := range chars {
		l, r := chars[i].Left, chars[i].Right
		if l > r {
			l, r = r, l
		}
		b, t := chars[i].Bottom, chars[i].Top
		if b > t {
			b, t = t, b
		}
		e.nl[i], e.nr[i], e.nb[i], e.nt[i] = l, r, b, t

		if i == 0 {
			minY, maxY = b, t
		} else {
			if b < minY {
				minY = b
			}
			if t > maxY {
				maxY = t
			}
		}

		e.nonSpacePrefix[i+1] = e.nonSpacePrefix[i]
		if chars[i].Unicode != ' ' {
			e.nonSpacePrefix[i+1]++
		}
	}

	numBuckets := min(max(len(chars)/4, 1), maxBuckets)
	if maxY <= minY {
		numBuckets = 1
	}

	e.minY = minY
	e.buckets = make([][]int32, numBuckets)
	if numBuckets > 1 {
		e.bucketH = (maxY - minY) / float32(numBuckets)
	}

	for i := range chars {
		lo, hi := e.bucketRange(e.nb[i], e.nt[i])
		for bkt := lo; bkt <= hi; bkt++ {
			e.buckets[bkt] = append(e.buckets[bkt], int32(i))
		}
	}

	return e
}

func (e *Extractor) bucketRange(b, t float32) (int, int) {
	if e.bucketH <= 0 {
		return 0, 0
	}
	lo := int((b - e.minY) / e.bucketH)
	hi := int((t - e.minY) / e.bucketH)
	if lo < 0 {
		lo = 0
	}
	// Both ends must be clamped: for a rect entirely above the page content lo
	// alone can land past the last bucket, and the hi < lo fixup below would
	// then hand that out-of-range index straight to the bucket loop.
	if lo >= len(e.buckets) {
		lo = len(e.buckets) - 1
	}
	if hi >= len(e.buckets) {
		hi = len(e.buckets) - 1
	}
	if hi < lo {
		hi = lo
	}
	return lo, hi
}

// TextInRect returns exactly what FPDFText_GetBoundedText(left, top, right,
// bottom) returns for this page, as a Go string.
func (e *Extractor) TextInRect(left, top, right, bottom float32) string {
	// CFX_FloatRect::Intersect normalizes both rects first.
	if left > right {
		left, right = right, left
	}
	if bottom > top {
		top, bottom = bottom, top
	}

	// Collect the chars whose box intersects the rect (PDFium's IsRectIntersect:
	// strictly positive overlap in both axes).
	matched := e.matchBuf[:0]
	e.epoch++
	lo, hi := e.bucketRange(bottom, top)
	for bkt := lo; bkt <= hi; bkt++ {
		for _, i := range e.buckets[bkt] {
			if e.seen[i] == e.epoch {
				continue
			}
			e.seen[i] = e.epoch
			il := max32(left, e.nl[i])
			ir := min32(right, e.nr[i])
			if il >= ir {
				continue
			}
			ib := max32(bottom, e.nb[i])
			it := min32(top, e.nt[i])
			if ib >= it {
				continue
			}
			matched = append(matched, i)
		}
	}
	e.matchBuf = matched
	if len(matched) == 0 {
		return ""
	}
	slices.Sort(matched)

	// Direct port of CPDF_TextPage::GetTextByPredicate. The flag state between
	// two matched chars only depends on the unmatched chars in the gap:
	//   - IsContainPreChar survives only across an empty gap.
	//   - The first gap char being a space appends one ' '.
	//   - IsAddLineFeed ends up true iff the gap contains any non-space char.
	out := make([]rune, 0, len(matched)+8)
	posy := float32(0)
	containPre := false
	addLineFeed := false
	prev := int32(-1)

	for _, m := range matched {
		gapStart := prev + 1
		if prev >= 0 {
			if m > gapStart {
				if e.chars[gapStart].Unicode == ' ' {
					out = append(out, ' ')
				}
				containPre = false
				addLineFeed = e.nonSpacePrefix[m]-e.nonSpacePrefix[gapStart] > 0
			}
		} else {
			// Chars before the first match: IsContainPreChar starts false and
			// only the else-branch (non-space) can set IsAddLineFeed.
			containPre = false
			addLineFeed = e.nonSpacePrefix[m] > 0
		}

		c := &e.chars[m]
		if posy-c.OriginY != 0 && !containPre && addLineFeed {
			posy = c.OriginY
			if len(out) > 0 {
				out = append(out, '\r', '\n')
			}
		}
		containPre = true
		addLineFeed = false
		if c.Unicode != 0 {
			out = append(out, c.Unicode)
		}
		prev = m
	}

	// Trailing unmatched space directly after the last matched char.
	if int(prev+1) < len(e.chars) && e.chars[prev+1].Unicode == ' ' {
		out = append(out, ' ')
	}

	return string(out)
}

// kSizeEpsilon mirrors CPDF_TextPage's constant: chars with a box narrower or
// shorter than this (in particular PDFium's generated space/newline chars,
// which have a zero-size box) never contribute to a rect.
const kSizeEpsilon = 0.01

// Rect is a text rect as FPDFText_GetRect returns it.
type Rect struct {
	Left, Top, Right, Bottom float32
}

// FirstCharIndices returns, for every rect in rects, the index of the first
// char that makes up that rect, or -1 when no char lies in it. rects must be
// all rects of the page in FPDFText_GetRect order.
//
// PDFium builds its rects (CPDF_TextPage::GetRectArray) by walking the chars
// in order and taking the union of the char boxes of each run of chars that
// share a text object, skipping generated and degenerate chars. Every char of
// a rect therefore lies inside it, and the rects come out in char order. This
// method replays that walk from the geometry alone: a run starts at the first
// unassigned char inside the rect and continues while chars stay inside it.
// Text objects drawn over each other make the pure geometric walk ambiguous,
// which is resolved with the next rect: once the run's union has grown to the
// rect itself, the first char that lies inside the next rect starts the next
// run, even when it also lies inside the current rect.
//
// Unlike FPDFText_GetCharIndexAtPos at the rect's corner this needs no
// tolerance and cannot pick an unrelated overlapping char: the corner of a
// rect is not on the first char whenever a later char in the rect is taller,
// and the position lookup returns the first char in page order whose box
// contains the point, whichever text object it belongs to.
func (e *Extractor) FirstCharIndices(rects []Rect) []int {
	result := make([]int, len(rects))
	cursor := 0
	prevFirst := -1
	for ri := range rects {
		r := normalizeRect(rects[ri])
		var next *Rect
		if ri+1 < len(rects) {
			n := normalizeRect(rects[ri+1])
			next = &n
		}

		first := e.scanRect(cursor, r)
		if first < 0 && cursor > prevFirst+1 {
			// The previous run swallowed this rect's chars, which happens
			// when this rect lies geometrically inside the previous one.
			// Re-scan the swallowed range.
			first = e.scanRect(prevFirst+1, r)
		}
		result[ri] = first
		if first < 0 {
			continue
		}

		// Consume the run: chars that are degenerate (generated) or inside
		// this rect. The union of the run's boxes grows towards the rect; once
		// it equals the rect, a char inside the next rect belongs to it.
		ul, ut, ur, ub := e.nl[first], e.nt[first], e.nr[first], e.nb[first]
		cursor = first + 1
		for cursor < len(e.chars) {
			if e.degenerate(cursor) {
				cursor++
				continue
			}
			if !e.contained(cursor, r) {
				break
			}
			if next != nil && ul == r.Left && ut == r.Top && ur == r.Right && ub == r.Bottom && e.contained(cursor, *next) {
				break
			}
			ul, ut = min32(ul, e.nl[cursor]), max32(ut, e.nt[cursor])
			ur, ub = max32(ur, e.nr[cursor]), min32(ub, e.nb[cursor])
			cursor++
		}
		prevFirst = first
	}
	return result
}

func normalizeRect(r Rect) Rect {
	if r.Left > r.Right {
		r.Left, r.Right = r.Right, r.Left
	}
	if r.Bottom > r.Top {
		r.Top, r.Bottom = r.Bottom, r.Top
	}
	return r
}

func (e *Extractor) scanRect(from int, r Rect) int {
	for i := from; i < len(e.chars); i++ {
		if e.degenerate(i) {
			continue
		}
		if e.contained(i, r) {
			return i
		}
	}
	return -1
}

func (e *Extractor) degenerate(i int) bool {
	return e.nr[i]-e.nl[i] < kSizeEpsilon || e.nt[i]-e.nb[i] < kSizeEpsilon
}

// contained reports whether the normalized box of char i lies within the
// normalized rect. The rect is a float32 union of exactly these boxes, so the
// comparison is exact, like PDFium's own min/max.
func (e *Extractor) contained(i int, r Rect) bool {
	return e.nl[i] >= r.Left && e.nr[i] <= r.Right && e.nb[i] >= r.Bottom && e.nt[i] <= r.Top
}

func max32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func min32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
