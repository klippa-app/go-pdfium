package implementation_cgo

import (
	"errors"
	"io"
	"unicode/utf8"
)

// A Reader implements the io.Reader, io.ReaderAt, io.WriterTo, io.Seeker,
// io.ByteScanner, and io.RuneScanner interfaces by reading from
// a byte slice.
// Unlike a Buffer, a Reader is read-only and supports seeking.
type BytesReaderCloser struct {
	s        []byte
	i        int64 // current reading index
	prevRune int   // index of previous rune; or < 0
}

// Len returns the number of bytes of the unread portion of the
// slice.
func (r *BytesReaderCloser) Len() int {
	if r.i >= int64(len(r.s)) {
		return 0
	}
	return int(int64(len(r.s)) - r.i)
}

// Size returns the original length of the underlying byte slice.
// Size is the number of bytes available for reading via ReadAt.
// The returned value is always the same and is not affected by calls
// to any other method.
func (r *BytesReaderCloser) Size() int64 { return int64(len(r.s)) }

// Read implements the io.Reader interface.
func (r *BytesReaderCloser) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	r.prevRune = -1
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}

// ReadAt implements the io.ReaderAt interface.
func (r *BytesReaderCloser) ReadAt(b []byte, off int64) (n int, err error) {
	// cannot modify state - see io.ReaderAt
	if off < 0 {
		return 0, errors.New("BytesReaderCloser.BytesReaderCloser.ReadAt: negative offset")
	}
	if off >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(b, r.s[off:])
	if n < len(b) {
		err = io.EOF
	}
	return
}

// ReadByte implements the io.ByteReader interface.
func (r *BytesReaderCloser) ReadByte() (byte, error) {
	r.prevRune = -1
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	b := r.s[r.i]
	r.i++
	return b, nil
}

// UnreadByte complements ReadByte in implementing the io.ByteScanner interface.
func (r *BytesReaderCloser) UnreadByte() error {
	r.prevRune = -1
	if r.i <= 0 {
		return errors.New("BytesReaderCloser.BytesReaderCloser.UnreadByte: at beginning of slice")
	}
	r.i--
	return nil
}

// ReadRune implements the io.RuneReader interface.
func (r *BytesReaderCloser) ReadRune() (ch rune, size int, err error) {
	if r.i >= int64(len(r.s)) {
		r.prevRune = -1
		return 0, 0, io.EOF
	}
	r.prevRune = int(r.i)
	if c := r.s[r.i]; c < utf8.RuneSelf {
		r.i++
		return rune(c), 1, nil
	}
	ch, size = utf8.DecodeRune(r.s[r.i:])
	r.i += int64(size)
	return
}

// UnreadRune complements ReadRune in implementing the io.RuneScanner interface.
func (r *BytesReaderCloser) UnreadRune() error {
	if r.prevRune < 0 {
		return errors.New("BytesReaderCloser.BytesReaderCloser.UnreadRune: previous operation was not ReadRune")
	}
	r.i = int64(r.prevRune)
	r.prevRune = -1
	return nil
}

// Seek implements the io.Seeker interface.
func (r *BytesReaderCloser) Seek(offset int64, whence int) (int64, error) {
	r.prevRune = -1
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.i + offset
	case io.SeekEnd:
		abs = int64(len(r.s)) + offset
	default:
		return 0, errors.New("BytesReaderCloser.BytesReaderCloser.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("BytesReaderCloser.BytesReaderCloser.Seek: negative position")
	}
	r.i = abs
	return abs, nil
}

// WriteTo implements the io.WriterTo interface.
func (r *BytesReaderCloser) WriteTo(w io.Writer) (n int64, err error) {
	r.prevRune = -1
	if r.i >= int64(len(r.s)) {
		return 0, nil
	}
	b := r.s[r.i:]
	m, err := w.Write(b)
	if m > len(b) {
		panic("BytesReaderCloser.BytesReaderCloser.WriteTo: invalid Write count")
	}
	r.i += int64(m)
	n = int64(m)
	if m != len(b) && err == nil {
		err = io.ErrShortWrite
	}
	return
}

// Reset resets the Reader to be reading from b.
func (r *BytesReaderCloser) Reset(b []byte) { *r = BytesReaderCloser{b, 0, -1} }

func (r *BytesReaderCloser) Close() error {
	// no-op for this reader.
	return nil
}

// BytesReaderCloser returns a new Reader reading from b.
func NewBytesReaderCloser(b []byte) *BytesReaderCloser { return &BytesReaderCloser{b, 0, -1} }
