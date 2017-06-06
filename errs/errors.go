package errs

import "fmt"

type errorAt struct {
	method string
	err    error
}

func ErrorAt(method string, err error) error {
	return errorAt{method: method, err: err}
}

func (e errorAt) Error() string {
	return fmt.Sprintf("call of %s occur: %v", e.method, e.err)
}
