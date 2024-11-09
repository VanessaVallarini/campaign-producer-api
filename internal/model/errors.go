package model

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type ErrorDetail struct {
	Code    string `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Error struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type BaseErrorWrapper struct {
	id      uint
	cause   error
	message string
	code    string
	details []ErrorDetail
}

func (bw *BaseErrorWrapper) Unwrap() error {
	return bw.cause
}

func (bw *BaseErrorWrapper) SetCode(code string) *BaseErrorWrapper {
	bw.code = code
	return bw
}

func (bw *BaseErrorWrapper) SetDetails(details []ErrorDetail) *BaseErrorWrapper {
	bw.details = details
	return bw
}

func (bw *BaseErrorWrapper) Builder(err error, message string) *BaseErrorWrapper {
	return &BaseErrorWrapper{
		message: message,
		cause:   err,
		id:      bw.id,
		code:    bw.code,
		details: bw.details,
	}
}

func (bw *BaseErrorWrapper) WithMessage(message string) *BaseErrorWrapper {
	return &BaseErrorWrapper{
		message: message,
		cause:   bw.cause,
		id:      bw.id,
		code:    bw.code,
		details: bw.details,
	}
}

func (bw *BaseErrorWrapper) Build() error {
	fmt.Printf("» %+v\n", bw)
	return errors.WithStack(bw) // so stacktracing "works"
}

func (bw *BaseErrorWrapper) Wrap(err error, message string) error {
	return bw.Builder(err, message).Build()
}

func (bw *BaseErrorWrapper) Wrapf(err error, message string, args ...interface{}) error {
	return bw.Wrap(err, fmt.Sprintf(message, args...))
}

func (bw *BaseErrorWrapper) Error() string {
	return bw.message
}

// pkg/errors requires this for stacktace composition
func (bw *BaseErrorWrapper) Cause() error {
	return bw.cause
}

func (bw *BaseErrorWrapper) ErrorDetails() Error {
	return Error{
		Message: bw.message,
		Code:    bw.code,
		Details: bw.details,
	}
}

// This is actually the only reason to reimplement all this stuff - it allows us
// to yield errors which are "checkable", while retaining error trace.
// If there's no need to decorate with a new trait, simply use `errors.Wrap()`
// If there's no need to keep the trace, use `«trait».Wrap(nil, «message»)`
func (bw *BaseErrorWrapper) Is(target error) bool {
	if other, ok := target.(*BaseErrorWrapper); ok {
		return other.id == bw.id
	}

	return false
}

// pkg/errors requires this for stacktace composition
func (bw *BaseErrorWrapper) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", bw.Cause())
			io.WriteString(s, bw.message)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, bw.Error())
	}
}

var (
	ErrNotFound            = &BaseErrorWrapper{message: "Not found", id: 1, code: "NOT_FOUND"}
	ErrInvalid             = &BaseErrorWrapper{message: "Invalid", id: 2, code: "INVALID_INPUT"}
	ErrInternal            = &BaseErrorWrapper{message: "Internal", id: 3, code: "INTERNAL_ERROR"}
	ErrNotImplemented      = &BaseErrorWrapper{message: "Not implemented", id: 4, code: "NOT_IMPLEMENTED"}
	ErrForbidden           = &BaseErrorWrapper{message: "Forbidden", id: 5, code: "FORBIDDEN"}
	ErrExternal            = &BaseErrorWrapper{message: "Bad gateway", id: 6, code: "BAD_GATEWAY"}
	ErrUnprocessableEntity = &BaseErrorWrapper{message: "Unprocessable entity", id: 7, code: "UNPROCESSABLE_ENTITY"}
)
