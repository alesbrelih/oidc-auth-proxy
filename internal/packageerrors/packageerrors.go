// This is meant for package errors
package packageerrors

import "net/http"

type Error struct {
	Code    int
	Message string
	Err     error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) WithErr(err error) *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Err:     err,
	}
}

func New(code int) *Error {
	return &Error{
		Code:    code,
		Message: http.StatusText(code),
	}
}

var ErrInternal = New(500)
