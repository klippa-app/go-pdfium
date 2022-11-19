package implementation_webassembly

import "C"
import (
	"bytes"
	"io/ioutil"
	"math"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/responses"

	"github.com/google/uuid"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func (p *PdfiumImplementation) registerTextPage(attachment *uint64, documentHandle *DocumentHandle) *TextPageHandle {
	ref := uuid.New()
	handle := &TextPageHandle{
		handle:      attachment,
		nativeRef:   references.FPDF_TEXTPAGE(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.textPageRefs[handle.nativeRef] = handle
	p.textPageRefs[handle.nativeRef] = handle

	return handle
}

func (p *PdfiumImplementation) registerPageLink(pageLink *uint64, documentHandle *DocumentHandle) *PageLinkHandle {
	ref := uuid.New()
	handle := &PageLinkHandle{
		handle:      pageLink,
		nativeRef:   references.FPDF_PAGELINK(ref.String()),
		documentRef: documentHandle.nativeRef,
	}

	documentHandle.pageLinkRefs[handle.nativeRef] = handle
	p.pageLinkRefs[handle.nativeRef] = handle

	return handle
}

func (p *PdfiumImplementation) transformUTF16LEToUTF8(charData []byte) (string, error) {
	pdf16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	utf16bom := unicode.BOMOverride(pdf16le.NewDecoder())
	unicodeReader := transform.NewReader(bytes.NewReader(charData), utf16bom)

	decoded, err := ioutil.ReadAll(unicodeReader)
	if err != nil {
		return "", err
	}

	// Remove NULL terminator.
	decoded = bytes.TrimSuffix(decoded, []byte("\x00"))

	return string(decoded), nil
}

func (p *PdfiumImplementation) transformUTF8ToUTF16LE(text string) ([]byte, error) {
	pdf16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	utf16bom := unicode.BOMOverride(pdf16le.NewEncoder())

	output := &bytes.Buffer{}
	unicodeWriter := transform.NewWriter(output, utf16bom)
	unicodeWriter.Write([]byte(text))
	unicodeWriter.Close()

	return output.Bytes(), nil
}

func convertPointPositions(pointPositions responses.CharPosition, ratio float64) *responses.CharPosition {
	return &responses.CharPosition{
		Left:   math.Round(pointPositions.Left * ratio),
		Top:    math.Round(pointPositions.Top * ratio),
		Right:  math.Round(pointPositions.Right * ratio),
		Bottom: math.Round(pointPositions.Bottom * ratio),
	}
}
