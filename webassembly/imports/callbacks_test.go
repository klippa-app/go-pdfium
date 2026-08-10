package imports

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

// FPDF_FILEACCESS_CB used to contain a failsafe that cleared out an io.EOF
// error when the reader had returned the full requested amount of bytes
// alongside it, and it only ever issued a single Read call. That failsafe was
// dropped in favour of io.ReadFull, which is documented to do both: it keeps
// reading until the buffer is full and only reports an error when it could not
// fill it.
//
// The tests below pin down that assumption, since pdfium's m_GetBlock must
// always be handed a completely filled buffer.

// eofWithFinalBytesReader returns io.EOF together with the last bytes of the
// underlying data instead of on a following zero-byte read. Some io.ReadSeeker
// implementations in the wild behave this way, and it is allowed by the
// io.Reader contract.
type eofWithFinalBytesReader struct {
	*bytes.Reader
}

func (r *eofWithFinalBytesReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if err == nil && r.Len() == 0 {
		err = io.EOF
	}
	return n, err
}

// chunkedReader never returns more than chunkSize bytes per Read call, without
// returning an error. A single Read call would leave the buffer partially
// filled.
type chunkedReader struct {
	*bytes.Reader
	chunkSize int
}

func (r *chunkedReader) Read(p []byte) (int, error) {
	if len(p) > r.chunkSize {
		p = p[:r.chunkSize]
	}
	return r.Reader.Read(p)
}

func TestReadFullAcceptsEOFWithFinalBytes(t *testing.T) {
	data := []byte("0123456789")
	reader := &eofWithFinalBytesReader{Reader: bytes.NewReader(data)}

	buf := make([]byte, len(data))
	n, err := io.ReadFull(reader, buf)
	if err != nil {
		t.Fatalf("io.ReadFull returned error %v, want nil", err)
	}
	if n != len(data) {
		t.Fatalf("io.ReadFull read %d bytes, want %d", n, len(data))
	}
	if !bytes.Equal(buf, data) {
		t.Fatalf("io.ReadFull read %q, want %q", buf, data)
	}
}

// The same reader, but only reading the last part of the data, which is what
// pdfium does when it reads the trailer at the end of a document.
func TestReadFullAcceptsEOFWithFinalBytesAfterSeek(t *testing.T) {
	data := []byte("0123456789")
	reader := &eofWithFinalBytesReader{Reader: bytes.NewReader(data)}

	if _, err := reader.Seek(6, io.SeekStart); err != nil {
		t.Fatalf("Seek returned error %v, want nil", err)
	}

	buf := make([]byte, 4)
	n, err := io.ReadFull(reader, buf)
	if err != nil {
		t.Fatalf("io.ReadFull returned error %v, want nil", err)
	}
	if n != len(buf) {
		t.Fatalf("io.ReadFull read %d bytes, want %d", n, len(buf))
	}
	if !bytes.Equal(buf, data[6:]) {
		t.Fatalf("io.ReadFull read %q, want %q", buf, data[6:])
	}
}

func TestReadFullFillsBufferOnShortReads(t *testing.T) {
	data := []byte("0123456789")
	reader := &chunkedReader{Reader: bytes.NewReader(data), chunkSize: 3}

	buf := make([]byte, len(data))
	n, err := io.ReadFull(reader, buf)
	if err != nil {
		t.Fatalf("io.ReadFull returned error %v, want nil", err)
	}
	if n != len(data) {
		t.Fatalf("io.ReadFull read %d bytes, want %d", n, len(data))
	}
	if !bytes.Equal(buf, data) {
		t.Fatalf("io.ReadFull read %q, want %q", buf, data)
	}
}

// The callback must still fail when the reader cannot deliver the full
// requested range, otherwise pdfium would be handed a partially filled buffer.
func TestReadFullReportsGenuinelyShortReads(t *testing.T) {
	data := []byte("01234")
	reader := &eofWithFinalBytesReader{Reader: bytes.NewReader(data)}

	buf := make([]byte, len(data)+5)
	n, err := io.ReadFull(reader, buf)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("io.ReadFull returned error %v, want io.ErrUnexpectedEOF", err)
	}
	if n != len(data) {
		t.Fatalf("io.ReadFull read %d bytes, want %d", n, len(data))
	}
}

// An empty reader reports io.EOF rather than io.ErrUnexpectedEOF, which the
// callback also has to treat as a failure.
func TestReadFullReportsEmptyReader(t *testing.T) {
	reader := &eofWithFinalBytesReader{Reader: bytes.NewReader(nil)}

	buf := make([]byte, 4)
	n, err := io.ReadFull(reader, buf)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("io.ReadFull returned error %v, want io.EOF", err)
	}
	if n != 0 {
		t.Fatalf("io.ReadFull read %d bytes, want 0", n)
	}
}
