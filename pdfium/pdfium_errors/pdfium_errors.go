package pdfium_errors

import "errors"

var (
	ErrSuccess = errors.New("success")
	ErrUnknown = errors.New("unknown error")
	ErrFile = errors.New("unable to read file")
	ErrFormat = errors.New("incorrect format")
	ErrPassword = errors.New("invalid password")
	ErrSecurity = errors.New("invalid encryption")
	ErrPage = errors.New("incorrect page")
	ErrUnexpected = errors.New("unexpected error")
)
