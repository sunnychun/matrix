package codes

import "fmt"

type Code int

const (
	OK Code = 0

	Internal Code = -1
	Unknown  Code = -2
	Aborted  Code = -3

	NotFound   Code = -101
	NotAllowed Code = -102

	InvalidHeader Code = -201
	InvalidParam  Code = -202

	OutOfRange Code = -301
	EncodeFail Code = -302
	DecodeFail Code = -303
)

var codes = map[Code]string{
	OK: "ok",

	Internal: "internal",
	Unknown:  "unknown",
	Aborted:  "aborted",

	NotFound:   "not found",
	NotAllowed: "not allowed",

	InvalidHeader: "invalid header",
	InvalidParam:  "invalid param",

	OutOfRange: "out of range",
	EncodeFail: "encode fail",
	DecodeFail: "decode fail",
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
