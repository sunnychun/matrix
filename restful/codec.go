package restful

import (
	"encoding/json"
	"io"
)

type Codec interface {
	ContentType() string
	Encode(io.Writer, interface{}) error
	Decode(io.Reader, interface{}) error
}

type JSONCodec struct{}

func (c JSONCodec) ContentType() string {
	return "application/json"
}

func (c JSONCodec) Encode(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func (c JSONCodec) Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

var _ Codec = JSONCodec{}
