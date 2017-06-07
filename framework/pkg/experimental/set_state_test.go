package experimental

import (
	"io"
	"os"
	"reflect"
	"testing"
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
