package responses

// AfterUnmarshaler is necessary to fix Gob encoding/decoding that causes
// the output to be different between implementations.
type AfterUnmarshaler interface {
	AfterUnmarshal()
}
