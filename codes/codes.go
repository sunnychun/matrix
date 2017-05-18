package codes

import "fmt"

type Code int

const (
	OK Code = 0

	Internal Code = -1

	Unknown Code = -2

	Aborted Code = -3

	NotFound Code = -4

	NotAllowed Code = -5

	InvalidParam Code = -6

	OutOfRange Code = -7
)

var codes = map[Code]string{
	OK:           "ok",
	Internal:     "internal",
	Unknown:      "unknown",
	Aborted:      "aborted",
	NotFound:     "not found",
	NotAllowed:   "not allowed",
	InvalidParam: "invalid param",
	OutOfRange:   "out of range",
}

func Register(code Code, desc string) {
	_, ok := codes[code]
	if ok {
		panic(fmt.Sprintf("%d:%s code is registered", code, desc))
	}
	codes[code] = desc
}

func (c Code) String() string {
	if desc, ok := codes[c]; ok {
		return desc
	}
	return fmt.Sprintf("code(%d)", c)
}
