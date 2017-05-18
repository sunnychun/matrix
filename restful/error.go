package restful

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ironzhang/matrix/codes"
)

type HTTPStatus interface {
	HTTPStatus() int
}

type ErrorCode interface {
	ErrorCode() codes.Code
}

type ErrorCause interface {
	ErrorCause() string
}

func NewError(status int, code codes.Code) Error {
	return Error{Status: status, Code: code}
}

func Errorf(status int, code codes.Code, format string, a ...interface{}) Error {
	return Error{
		Status: status,
		Code:   code,
		Cause:  fmt.Sprintf(format, a...),
	}
}

type Error struct {
	Status int
	Code   codes.Code
	Cause  string
}

func (e Error) HTTPStatus() int {
	return e.Status
}

func (e Error) ErrorCode() codes.Code {
	return e.Code
}

func (e Error) ErrorCause() string {
	return e.Cause
}

func (e Error) Error() string {
	if e.Cause != "" {
		return fmt.Sprintf("[%d: %s] %d: %s (%s)", e.Status, http.StatusText(e.Status), e.Code, e.Code.String(), e.Cause)
	}
	return fmt.Sprintf("[%d: %s] %d: %s", e.Status, http.StatusText(e.Status), e.Code, e.Code.String())
}

type rpcError struct {
	Code  int    `json:"code"`
	Desc  string `json:"desc"`
	Cause string `json:"cause,omitempty"`
}

func (e rpcError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func toRPCError(err error) error {
	if e, ok := err.(rpcError); ok {
		return e
	}
	code := codes.Internal
	if e, ok := err.(ErrorCode); ok {
		code = e.ErrorCode()
	}
	cause := err.Error()
	if e, ok := err.(ErrorCause); ok {
		cause = e.ErrorCause()
	}
	return rpcError{
		Code:  int(code),
		Desc:  code.String(),
		Cause: cause,
	}
}
