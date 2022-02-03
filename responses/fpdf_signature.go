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
	Signature []byte // For public-key signatures, Signature is either a DER-encoded PKCS#1 binary or a DER-encoded PKCS#7 binary.
}

type FPDFSignatureObj_GetByteRange struct {
	ByteRange []byte // ByteRange is an array of pairs of integers (starting byte offset, length in bytes) that describes the exact byte range for the digest calculation.
}

type FPDFSignatureObj_GetSubFilter struct {
	SubFilter string // The encoding of the value of a signature object.
}

type FPDFSignatureObj_GetReason struct {
	Reason string // The reason (comment) of the signature object.
}

type FPDFSignatureObj_GetTime struct {
	Time string // The time of signing of a signature object. The format of time is expected to be D:YYYYMMDDHHMMSS+XX'YY', i.e. it's precision is seconds, with timezone information. This value should be used only when the time of signing is not available in the (PKCS#7 binary) signature.
}

type FPDFSignatureObj_GetDocMDPPermission struct {
	DocMDPPermission int
}
