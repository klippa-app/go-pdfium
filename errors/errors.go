package errors

import "errors"

var (
	ErrSuccess                  = errors.New("0: success")
	ErrUnknown                  = errors.New("1: unknown error")
	ErrFile                     = errors.New("2: unable to read file")
	ErrFormat                   = errors.New("3: incorrect format")
	ErrPassword                 = errors.New("4: invalid password")
	ErrSecurity                 = errors.New("5: invalid encryption")
	ErrPage                     = errors.New("6: incorrect page")
	ErrXFALoad                  = errors.New("7: load XFA error")
	ErrXFALayout                = errors.New("8: layout XFA error")
	ErrUnexpected               = errors.New("unexpected error")
	ErrExperimentalUnsupported  = errors.New("this functionality is only supported when using the pdfium_experimental build flag, see https://github.com/klippa-app/go-pdfium#experimental for more information")
	ErrWindowsUnsupported       = errors.New("this functionality is Windows only")
	ErrUnsupportedOnWebassembly = errors.New("this functionality is not supported on Webassembly")
	ErrXFAUnsupported           = errors.New("this functionality is only supported when using the pdfium_xfa build flag, see https://github.com/klippa-app/go-pdfium#xfa for more information")
	ErrV8Unsupported            = errors.New("this functionality is only supported when using the pdfium_v8 build flag, see https://github.com/klippa-app/go-pdfium#v8 for more information")
)
