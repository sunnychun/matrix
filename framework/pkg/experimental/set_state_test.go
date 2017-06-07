package experimental

import (
	"encoding"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/ironzhang/matrix/jsoncfg"
)

func TestIndirect(t *testing.T) {
	var i int
	var s string
	var a interface{}
	var b interface{}
	var c io.Writer
	var d io.Writer
	var e interface{}

	b = &i
	d = os.Stdout
	e = os.Stdout

	tests := []struct {
		v reflect.Value
		p reflect.Value
	}{
		{reflect.ValueOf(&i), reflect.ValueOf(&i).Elem()},
		{reflect.ValueOf(&s), reflect.ValueOf(&s).Elem()},
		{reflect.ValueOf(&a), reflect.ValueOf(&a).Elem()},
		{reflect.ValueOf(&b), reflect.ValueOf(&i).Elem()},
		{reflect.ValueOf(b), reflect.ValueOf(&i).Elem()},
		{reflect.ValueOf(&c), reflect.ValueOf(&c).Elem()},
		{reflect.ValueOf(&d), reflect.ValueOf(os.Stdout).Elem()},
		{reflect.ValueOf(&e), reflect.ValueOf(os.Stdout).Elem()},
		{reflect.ValueOf(e), reflect.ValueOf(os.Stdout).Elem()},
	}
	for i, tt := range tests {
		p := indirect(tt.v)
		if got, want := p, tt.p; got != want {
			t.Errorf("tests[%d]: got(%v) != want(%v)", i, got, want)
		}
	}
}

func TestIndirectUnmarshaler(t *testing.T) {
	var i int
	var s string
	var a interface{}
	var b interface{}
	var d jsoncfg.Duration

	b = &i

	tests := []struct {
		v  interface{}
		u  json.Unmarshaler
		tu encoding.TextUnmarshaler
		rv reflect.Value
	}{
		{&i, nil, nil, reflect.ValueOf(&i).Elem()},
		{&s, nil, nil, reflect.ValueOf(&s).Elem()},
		{&a, nil, nil, reflect.ValueOf(&a).Elem()},
		{&b, nil, nil, reflect.ValueOf(&i).Elem()},
		{b, nil, nil, reflect.ValueOf(&i).Elem()},
		{&d, nil, &d, reflect.Value{}},
	}
	for i, tt := range tests {
		u, tu, rv := indirectUnmarshaler(reflect.ValueOf(tt.v))
		if got, want := u, tt.u; got != want {
			t.Errorf("tests[%d]: Unmarshaler: got(%v) != want(%v)", i, got, want)
		}
		if got, want := tu, tt.tu; got != want {
			t.Errorf("tests[%d]: TextUnmarshaler: got(%v) != want(%v)", i, got, want)
		}
		if got, want := rv, tt.rv; got != want {
			t.Errorf("tests[%d]: Value: got(%v) != want(%v)", i, got, want)
		}
	}
}
