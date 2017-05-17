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

type Code interface {
	Code() codes.Code
}

type Cause interface {
	Cause() string
}

type Error struct {
	Status int
	Code   codes.Code
	Cause  string
}

func toError(err error) error {
	if e, ok := err.(Error); ok {
		return e
	}
	status := http.StatusInternalServerError
	if e, ok := err.(HTTPStatus); ok {
		status = e.HTTPStatus()
	}
	code := codes.Internal
	if e, ok := err.(Code); ok {
		code = e.Code()
	}
	cause := err.Error()
	if e, ok := err.(Cause); ok {
		cause = e.Cause()
	}
	return NewError(status, code, cause)
}

func NewError(status int, code codes.Code, cause string) Error {
	return Error{Status: status, Code: code, Cause: cause}
}

func Errorf(status int, code codes.Code, format string, a ...interface{}) Error {
	return NewError(status, code, fmt.Sprintf(format, a...))
}

func (e Error) Error() string {
	if e.Cause != "" {
		return fmt.Sprintf("[%d: %s] %d: %s (%s)", e.Status, http.StatusText(e.Status), e.Code, e.Code.String(), e.Cause)
	}
	return fmt.Sprintf("[%d: %s] %d: %s", e.Status, http.StatusText(e.Status), e.Code, e.Code.String())
}

type jsonError struct {
	Code  int    `json:"code"`
	Desc  string `json:"desc"`
	Cause string `json:"cause,omitempty"`
}

func (e jsonError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}
