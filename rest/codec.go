package rest

import (
	"encoding/json"
	"io"
)

type Codec interface {
	Encode(io.Writer, interface{}) error
	Decode(io.Reader, interface{}) error
}

type JsonCodec struct{}

func (c JsonCodec) Encode(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(v)
}

func (c JsonCodec) Decode(r io.Reader, v interface{}) error {
	dec := json.NewDecoder(r)
	return dec.Decode(v)
}

var _ Codec = JsonCodec{}
