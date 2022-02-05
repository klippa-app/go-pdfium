package responses

import "github.com/klippa-app/go-pdfium/references"

type FPDF_GetSignatureCount struct {
	Count int
}

type FPDF_GetSignatureObject struct {
	Index     int
	Signature references.FPDF_SIGNATURE
}

type FPDFSignatureObj_GetContents struct {
	Contents []byte // For public-key signatures, Contents is either a DER-encoded PKCS#1 binary or a DER-encoded PKCS#7 binary. nil when no content.
}

type FPDFSignatureObj_GetByteRange struct {
	ByteRange []int // ByteRange is an array of pairs of integers (starting byte offset, length in bytes) that describes the exact byte range for the digest calculation. nil when no byte range.
}

type FPDFSignatureObj_GetSubFilter struct {
	SubFilter *string // The encoding of the value of a signature object. nil when no sub filter.
}

type FPDFSignatureObj_GetReason struct {
	Reason *string // The reason (comment) of the signature object. nil when no reason.
}

type FPDFSignatureObj_GetTime struct {
	Time *string // The time of signing of a signature object. The format of time is expected to be D:YYYYMMDDHHMMSS+XX'YY', i.e. it's precision is seconds, with timezone information. This value should be used only when the time of signing is not available in the (PKCS#7 binary) signature. nil when no time.
}

type FPDFSignatureObj_GetDocMDPPermission struct {
	DocMDPPermission int
}
