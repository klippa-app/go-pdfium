package shared_tests

// PDFium's FPDFText_GetBoundedText re-scans every character on the page on
// every call, which makes "give me the text of all rects on this page"
// O(chars x rects). internal/textextract avoids that by collecting the
// per-character data once and answering each rect in O(matched chars) -- but
// doing so means CPDF_TextPage::GetTextByRect(), its float32 rect intersection
// and its line-feed/space state machine are all reimplemented in Go.
//
// That reimplementation can silently drift from PDFium: a change to
// GetTextByPredicate() upstream, or to the way character boxes are generated,
// would leave our port quietly returning something else. The tests below pin
// the port to PDFium itself rather than to hand-written expectations, by using
// FPDFText_GetBoundedText -- which is a direct, unmodified PDFium call on both
// the cgo and WebAssembly backends -- as the oracle:
//
//   - "matches PDFium for every rect on the page" covers the production path
//     end to end (GetPageTextStructured, including the per-char collection in
//     each backend's text.go) over a corpus of single-line, multi-line, table,
//     rotated, vertical, control-character and multi-byte pages.
//   - "matches PDFium for arbitrary query rects" drives the Extractor directly
//     with rects PDFium's own layout never produces (partial char overlap,
//     slivers in line gaps, degenerate and inverted rects), which is where the
//     intersection test and the state machine are most likely to diverge.
//
// If PDFium changes this logic, these fail with the exact rect and both
// strings, instead of the drift reaching users.

import (
	"fmt"
	"io/ioutil"

	"github.com/klippa-app/go-pdfium/internal/textextract"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Files whose every page is checked through the production path. Between them
// they cover a single text run, multiple lines, a table, rotated and vertical
// writing, control characters, a page with no text at all, and ~6k characters
// of multi-byte text per page (rect-wrong.pdf).
//
// All five rect-wrong.pdf pages are checked deliberately. The oracle is
// O(chars) per call, so those pages alone are most of this file's runtime
// (~29s of it on the WebAssembly interpreter, against a ~48s baseline for that
// package); full coverage of the multi-byte pages was judged worth it. Trim
// pages here, not query rects in the sweep, if that ever has to change.
var textExtractParityCorpus = []string{
	"test.pdf",
	"rect-wrong.pdf",
	"hello_world.pdf",
	"bigtable_mini.pdf",
	"rotated_text.pdf",
	"vertical_text.pdf",
	"cropped_text.pdf",
	"find_text_consecutive.pdf",
	"text_render_mode.pdf",
	"control_characters.pdf",
	"tagged_table.pdf",
}

// Pages driven with arbitrary query rects. Deliberately a subset: the oracle
// is O(chars) per call, so the char-heavy rect-wrong.pdf page gets a smaller
// query budget (see textExtractSweepQueries) to keep this affordable on the
// WebAssembly interpreter.
var textExtractSweepPages = []struct {
	File string
	Page int
}{
	{"test.pdf", 0},
	{"hello_world.pdf", 0},
	{"bigtable_mini.pdf", 0},
	{"rotated_text.pdf", 0},
	{"vertical_text.pdf", 0},
	{"cropped_text.pdf", 0},
	{"find_text_consecutive.pdf", 0},
	{"text_render_mode.pdf", 0},
	{"control_characters.pdf", 0},
	{"rect-wrong.pdf", 1},
}

// textExtractRand is a fixed-increment LCG. Hand-rolled rather than math/rand
// so the generated rects are identical across Go versions and backends -- a
// failure here must be reproducible.
type textExtractRand uint64

func (r *textExtractRand) next() uint32 {
	*r = *r*6364136223846793005 + 1442695040888963407
	return uint32(*r >> 33)
}

func (r *textExtractRand) float32Between(lo, hi float32) float32 {
	return lo + (hi-lo)*float32(r.next()%100001)/100000
}

type textExtractQuery struct {
	Left, Top, Right, Bottom float32
	Origin                   string
}

// textExtractPageChars reads the per-character data the Extractor is built
// from, using only the public API, mirroring what each backend's text.go does.
func textExtractPageChars(textPage references.FPDF_TEXTPAGE, count int) []textextract.Char {
	chars := make([]textextract.Char, 0, count)
	for i := 0; i < count; i++ {
		box, err := PdfiumInstance.FPDFText_GetCharBox(&requests.FPDFText_GetCharBox{TextPage: textPage, Index: i})
		Expect(err).To(BeNil())
		origin, err := PdfiumInstance.FPDFText_GetCharOrigin(&requests.FPDFText_GetCharOrigin{TextPage: textPage, Index: i})
		Expect(err).To(BeNil())
		unicode, err := PdfiumInstance.FPDFText_GetUnicode(&requests.FPDFText_GetUnicode{TextPage: textPage, Index: i})
		Expect(err).To(BeNil())

		chars = append(chars, textextract.Char{
			Left:    float32(box.Left),
			Bottom:  float32(box.Bottom),
			Right:   float32(box.Right),
			Top:     float32(box.Top),
			OriginY: float32(origin.Y),
			Unicode: rune(unicode.Unicode),
		})
	}
	return chars
}

// textExtractSweepQueries builds the query rects for one page. Coordinates are
// generated as float32 so the Extractor (float32) and the oracle (float64) see
// the exact same value and any difference is a real logic difference.
func textExtractSweepQueries(chars []textextract.Char) []textExtractQuery {
	if len(chars) == 0 {
		return nil
	}

	minX, maxX := min(chars[0].Left, chars[0].Right), max(chars[0].Left, chars[0].Right)
	minY, maxY := min(chars[0].Bottom, chars[0].Top), max(chars[0].Bottom, chars[0].Top)
	for _, c := range chars {
		minX, maxX = min(minX, min(c.Left, c.Right)), max(maxX, max(c.Left, c.Right))
		minY, maxY = min(minY, min(c.Bottom, c.Top)), max(maxY, max(c.Bottom, c.Top))
	}

	// Keep the char-box family bounded on char-heavy pages by sampling.
	charStride := 1 + len(chars)/150
	randomQueries := 400
	bands := 60
	if len(chars) > 1000 {
		randomQueries = 60
		bands = 20
	}

	queries := make([]textExtractQuery, 0, randomQueries+bands+16)

	// Each char's own box, then inset and outset: exercises partial overlap and
	// the "touching edge does not intersect" boundary.
	for i := 0; i < len(chars); i += charStride {
		c := chars[i]
		for _, e := range []float32{0, 0.5, -0.5, 2, -2} {
			queries = append(queries, textExtractQuery{
				Left: c.Left - e, Top: c.Top + e, Right: c.Right + e, Bottom: c.Bottom - e,
				Origin: fmt.Sprintf("char %d box inset %g", i, e),
			})
		}
	}

	// Degenerate and inverted rects.
	queries = append(queries,
		textExtractQuery{minX, maxY, maxX, minY, "whole content box"},
		textExtractQuery{minX, minY, maxX, maxY, "vertically inverted"},
		textExtractQuery{maxX, maxY, minX, minY, "horizontally inverted"},
		textExtractQuery{minX, maxY, minX, minY, "zero width"},
		textExtractQuery{minX, maxY, maxX, maxY, "zero height"},
		textExtractQuery{0, 0, 0, 0, "empty at origin"},
		textExtractQuery{minX - 50, minY - 10, maxX + 50, minY - 20, "entirely below content"},
		textExtractQuery{minX - 50, maxY + 20, maxX + 50, maxY + 10, "entirely above content"},
	)

	// Random rects over the content box plus a margin.
	rand := textExtractRand(0x5eed)
	pad := float32(5)
	for i := 0; i < randomQueries; i++ {
		x0 := rand.float32Between(minX-pad, maxX+pad)
		x1 := rand.float32Between(minX-pad, maxX+pad)
		y0 := rand.float32Between(minY-pad, maxY+pad)
		y1 := rand.float32Between(minY-pad, maxY+pad)
		queries = append(queries, textExtractQuery{
			Left: x0, Top: max(y0, y1), Right: x1, Bottom: min(y0, y1),
			Origin: fmt.Sprintf("random %d", i),
		})
	}

	// Thin horizontal bands, which land in the gaps between lines and drive the
	// line-feed branch of the state machine.
	for i := 0; i <= bands; i++ {
		y := minY + (maxY-minY)*float32(i)/float32(bands)
		queries = append(queries, textExtractQuery{
			Left: minX - pad, Top: y + 0.3, Right: maxX + pad, Bottom: y - 0.3,
			Origin: fmt.Sprintf("band %d", i),
		})
	}

	return queries
}

var _ = Describe("textextract parity with PDFium", func() {
	BeforeEach(func() {
		Locker.Lock()
	})

	AfterEach(func() {
		Locker.Unlock()
	})

	// withPage loads a file and runs body once per page, with a text page
	// loaded, cleaning up as it goes.
	withPage := func(file string, maxPages int, body func(doc references.FPDF_DOCUMENT, page int, textPage references.FPDF_TEXTPAGE, charCount int)) {
		pdfData, err := ioutil.ReadFile(TestDataPath + "/testdata/" + file)
		Expect(err).To(BeNil())

		loadedDoc, err := PdfiumInstance.FPDF_LoadMemDocument(&requests.FPDF_LoadMemDocument{Data: &pdfData})
		Expect(err).To(BeNil())
		doc := loadedDoc.Document

		defer func() {
			_, err := PdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc})
			Expect(err).To(BeNil())
		}()

		pageCount, err := PdfiumInstance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{Document: doc})
		Expect(err).To(BeNil())

		pages := pageCount.PageCount
		if maxPages > 0 && maxPages < pages {
			pages = maxPages
		}

		for page := 0; page < pages; page++ {
			textPage, err := PdfiumInstance.FPDFText_LoadPage(&requests.FPDFText_LoadPage{
				Page: requests.Page{ByIndex: &requests.PageByIndex{Document: doc, Index: page}},
			})
			Expect(err).To(BeNil())

			charCount, err := PdfiumInstance.FPDFText_CountChars(&requests.FPDFText_CountChars{TextPage: textPage.TextPage})
			Expect(err).To(BeNil())

			body(doc, page, textPage.TextPage, charCount.Count)

			_, err = PdfiumInstance.FPDFText_ClosePage(&requests.FPDFText_ClosePage{TextPage: textPage.TextPage})
			Expect(err).To(BeNil())
		}
	}

	// boundedText is the oracle: a direct, unmodified PDFium call.
	boundedText := func(textPage references.FPDF_TEXTPAGE, left, top, right, bottom float32) string {
		resp, err := PdfiumInstance.FPDFText_GetBoundedText(&requests.FPDFText_GetBoundedText{
			TextPage: textPage,
			Left:     float64(left),
			Top:      float64(top),
			Right:    float64(right),
			Bottom:   float64(bottom),
		})
		Expect(err).To(BeNil())
		return resp.Text
	}

	Context("the rect text produced by GetPageTextStructured", func() {
		for _, file := range textExtractParityCorpus {
			It("matches PDFium for every rect in "+file, func() {
				compared := 0
				mismatches := []string{}

				withPage(file, 0, func(doc references.FPDF_DOCUMENT, page int, textPage references.FPDF_TEXTPAGE, charCount int) {
					// FPDFText_GetRect needs FPDFText_CountRects first, and
					// this is the exact call the implementations make.
					rectCount, err := PdfiumInstance.FPDFText_CountRects(&requests.FPDFText_CountRects{
						TextPage:   textPage,
						StartIndex: 0,
						Count:      charCount,
					})
					Expect(err).To(BeNil())

					structured, err := PdfiumInstance.GetPageTextStructured(&requests.GetPageTextStructured{
						Page: requests.Page{ByIndex: &requests.PageByIndex{Document: doc, Index: page}},
						Mode: requests.GetPageTextStructuredModeRects,
					})
					Expect(err).To(BeNil())
					Expect(len(structured.Rects)).To(
						Equal(rectCount.Count),
						fmt.Sprintf("%s page %d: rect count", file, page),
					)

					for i := 0; i < rectCount.Count; i++ {
						rect, err := PdfiumInstance.FPDFText_GetRect(&requests.FPDFText_GetRect{
							TextPage: textPage,
							Index:    i,
						})
						Expect(err).To(BeNil())

						want := boundedText(textPage,
							float32(rect.Left), float32(rect.Top), float32(rect.Right), float32(rect.Bottom))
						got := structured.Rects[i].Text

						compared++
						if got != want {
							mismatches = append(mismatches, fmt.Sprintf(
								"%s page %d rect %d (%g,%g,%g,%g):\n   PDFium: %q\n   Go:     %q",
								file, page, i, rect.Left, rect.Top, rect.Right, rect.Bottom, want, got))
						}
					}
				})

				Expect(mismatches).To(BeEmpty(), fmt.Sprintf(
					"internal/textextract has drifted from PDFium's CPDF_TextPage::GetTextByRect (%d rects compared)",
					compared))
			})
		}
	})

	Context("the Extractor driven with arbitrary query rects", func() {
		for _, sweep := range textExtractSweepPages {
			It(fmt.Sprintf("matches PDFium for arbitrary rects in %s page %d", sweep.File, sweep.Page), func() {
				ran := false

				withPage(sweep.File, sweep.Page+1, func(doc references.FPDF_DOCUMENT, page int, textPage references.FPDF_TEXTPAGE, charCount int) {
					if page != sweep.Page {
						return
					}
					ran = true

					chars := textExtractPageChars(textPage, charCount)
					extractor := textextract.New(chars)
					queries := textExtractSweepQueries(chars)
					Expect(queries).ToNot(BeEmpty())

					compared, mismatches := 0, []string{}
					for _, q := range queries {
						got := extractor.TextInRect(q.Left, q.Top, q.Right, q.Bottom)
						want := boundedText(textPage, q.Left, q.Top, q.Right, q.Bottom)
						compared++
						if got != want {
							mismatches = append(mismatches, fmt.Sprintf(
								"%s page %d, %s (%g,%g,%g,%g):\n   PDFium: %q\n   Go:     %q",
								sweep.File, page, q.Origin, q.Left, q.Top, q.Right, q.Bottom, want, got))
							// Enough to diagnose; a real divergence usually
							// trips a large share of the queries at once.
							if len(mismatches) >= 10 {
								break
							}
						}
					}

					Expect(mismatches).To(BeEmpty(), fmt.Sprintf(
						"internal/textextract has drifted from PDFium's CPDF_TextPage::GetTextByRect (%d of %d query rects compared before stopping)",
						compared, len(queries)))
				})

				Expect(ran).To(BeTrue(), "sweep page was never reached")
			})
		}
	})
})
