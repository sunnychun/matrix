package values

import (
	"fmt"
	"reflect"
)

type Value struct {
	ptr interface{}
}

func (v *Value) Get() interface{} {
	return v.ptr
}

func (v *Value) Set(x interface{}) error {
	return SetValue(v.ptr, x)
}

type Values map[string]*Value

func (values Values) Register(name string, ptr interface{}) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("value invalid: %s, %T", name, ptr)
	}
	if _, ok := values[name]; ok {
		return fmt.Errorf("value duplicate: %s", name)
	}
	values[name] = &Value{ptr: ptr}
	return nil
}
