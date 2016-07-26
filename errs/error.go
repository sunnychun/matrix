package errs

import "fmt"

type Errno int

func (e Errno) Error() string {
	if s, ok := errors[e]; ok {
		return s
	}
	return fmt.Sprintf("errno %d", e)
}

func GetErrno(err error) Errno {
	if err == nil {
		return Errno(0)
	}
	err = cause(err)
	if e, ok := err.(Errno); ok {
		return e
	}
	return ErrUnknown
}

type causeError struct {
	err error
	msg string
}

func (e *causeError) Error() string {
	return e.msg + ": " + e.err.Error()
}

func New(err error, msg string) error {
	return &causeError{err: err, msg: msg}
}

func Errorf(err error, format string, a ...interface{}) error {
	return New(err, fmt.Sprintf(format, a...))
}

func Equal(a, b error) bool {
	a = cause(a)
	b = cause(b)
	return a == b
}

func NotEqual(a, b error) bool {
	return !Equal(a, b)
}

func cause(err error) error {
	for err != nil {
		c, ok := err.(*causeError)
		if !ok {
			break
		}
		err = c
	}
	return err
}
