package experimental

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/ironzhang/matrix/errs"
)

var errInvalidValue = errors.New("invalid value")

func indirect(v reflect.Value) reflect.Value {
	if !v.IsValid() {
		panic(errs.ErrorAt("indirect", errInvalidValue))
	}

	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}

L:
	for {
		switch k := v.Kind(); k {
		case reflect.Ptr:
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		case reflect.Interface:
			if v.IsNil() {
				break L
			}
			v = v.Elem()
		default:
			break L
		}
	}
	return v
}

func indirectUnmarshaler(v reflect.Value) (json.Unmarshaler, encoding.TextUnmarshaler, reflect.Value) {
	if !v.IsValid() {
		panic(errs.ErrorAt("indirectUnmarshaler", errInvalidValue))
	}

	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}

L:
	for {
		if v.Type().NumMethod() > 0 {
			if u, ok := v.Interface().(json.Unmarshaler); ok {
				return u, nil, reflect.Value{}
			}
			if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
				return nil, u, reflect.Value{}
			}
		}

		switch k := v.Kind(); k {
		case reflect.Ptr:
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		case reflect.Interface:
			if v.IsNil() {
				break L
			}
			v = v.Elem()
		default:
			break L
		}
	}
	return nil, nil, v
}

type setError struct {
	src reflect.Type
	dst reflect.Type
	err error
}

func (e setError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("can not set value from %s[%s] to %s[%s]: %v", e.src.Kind(), e.src, e.dst.Kind(), e.dst, e.err)
	}
	return fmt.Sprintf("can not set value from %s[%s] to %s[%s]", e.src.Kind(), e.src, e.dst.Kind(), e.dst)
}

type setState struct {
}

// x = y
func (s *setState) set(x, y reflect.Value) {
}

func (s *setState) setValue(x, y reflect.Value) {
}

func (s *setState) setString(x, y reflect.Value) {
}

func (s *setState) setMap(x, y reflect.Value) {
}

func (s *setState) setObject(x, y reflect.Value) {
}
