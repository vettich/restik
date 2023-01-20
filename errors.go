package restik

import (
	"fmt"
	"net/http"
)

var (
	// ErrNotFoundEndpoint is error
	ErrNotFoundEndpoint = NewError(404, "endpoint_not_found", "Endpoint not found")
)

type Error interface {
	Error() string
	GetStatus() int
	GetCode() string
	GetMessage() string
}

// Error - rest errors
type errorImpl struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Msg    string `json:"msg"`
}

// NewError return new Error instance
func NewError(status int, code, msg string) Error {
	return &errorImpl{status, code, msg}
}

// FromAnotherError wrap any error to Error
func FromAnotherError(err error) Error {
	if e, ok := err.(Error); ok {
		return e
	}
	return NewBadRequestError(err.Error())
}

// NewBadRequestError return error with BadRequest status
//
// Using:
//		NewBadRequestError()
// or
//		NewBadRequestError(message)
// or
//		NewBadRequestError(code, message)
func NewBadRequestError(args ...string) Error {
	return NewError(parseErrorArgs(http.StatusBadRequest, "bad_request", "Bad request", args...))
}

// NewNotFoundError return error with BadRequest status
//
// Using:
//		NewNotFoundError()
// or
//		NewNotFoundError(message)
// or
//		NewNotFoundError(code, message)
func NewNotFoundError(args ...string) Error {
	return NewError(parseErrorArgs(http.StatusNotFound, "not_found", "Not found", args...))
}

// NewInternalError return error with BadRequest status
//
// Using:
//		NewInternalError()
// or
//		NewInternalError(message)
// or
//		NewInternalError(code, message)
func NewInternalError(args ...string) Error {
	return NewError(parseErrorArgs(http.StatusInternalServerError, "internal_error", "Intertal Server Error", args...))
}

// Error implement error interface
func (err errorImpl) Error() string {
	return fmt.Sprintf("[%d] %s", err.Status, err.Msg)
}

// GetStatus return http status
func (err errorImpl) GetStatus() int {
	return err.Status
}

// GetCode return error code
func (err errorImpl) GetCode() string {
	return err.Code
}

// GetMessage return error message
func (err errorImpl) GetMessage() string {
	return err.Msg
}

// using:
//		parseErrorArgs(status, defaultCode, defaultMsg)
// or
//		parseErrorArgs(status, defaultCode, defaultMsg, message)
// or
//		parseErrorArgs(status, defaultCode, defaultMsg, code, message)
func parseErrorArgs(status int, defaultCode, defaultMsg string, args ...string) (int, string, string) {
	code, msg := defaultCode, defaultMsg
	if len(args) == 1 {
		msg = args[0]
	} else if len(args) == 2 {
		code, msg = args[0], args[1]
	}
	return status, code, msg
}
