package restful

import (
	"encoding/json"
	"encoding/xml"
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

type XMLCodec struct{}

func (c XMLCodec) ContentType() string {
	return "application/xml"
}

func (c XMLCodec) Encode(w io.Writer, v interface{}) error {
	return xml.NewEncoder(w).Encode(v)
}

func (c XMLCodec) Decode(r io.Reader, v interface{}) error {
	return xml.NewDecoder(r).Decode(v)
}
