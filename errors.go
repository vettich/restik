package restik

import (
	"fmt"
	"net/http"
)

var (
	// ErrNotFoundEndpoint is error
	ErrNotFoundEndpoint = NewError(404, "endpoint_not_found", "Endpoint not found")
)

// Error - rest errors
type Error struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Msg    string `json:"msg"`
}

// NewError return new Error instance
func NewError(status int, code, msg string) *Error {
	return &Error{status, code, msg}
}

// FromAnotherError wrap any error to Error
func FromAnotherError(err error) *Error {
	if e, ok := err.(*Error); ok {
		return e
	}
	return NewBadRequestError(err.Error())
}

// NewBadRequestError return error with BadRequest status
func NewBadRequestError(args ...string) *Error {
	return NewError(parseErrorArgs(http.StatusBadRequest, "bad_request", "Bad request", args...))
}

// NewNotFoundError return error with BadRequest status
func NewNotFoundError(args ...string) *Error {
	return NewError(parseErrorArgs(http.StatusNotFound, "not_found", "Not found", args...))
}

// Error implement error interface
func (err Error) Error() string {
	return fmt.Sprintf("[%d] %s", err.Status, err.Msg)
}

// GetStatus return http status
func (err Error) GetStatus() int {
	return err.Status
}

// GetCode return error code
func (err Error) GetCode() string {
	return err.Code
}

// GetMsg return error message
func (err Error) GetMsg() string {
	return err.Msg
}

func parseErrorArgs(status int, defaultCode, defaultMsg string, args ...string) (int, string, string) {
	code, msg := defaultCode, defaultMsg
	if len(args) == 1 {
		msg = args[0]
	} else if len(args) == 2 {
		code, msg = args[0], args[1]
	}
	return status, code, msg
}
