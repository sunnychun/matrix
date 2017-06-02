package framework

import (
	"fmt"
	"reflect"
)

type values map[string]interface{}

func (values values) Register(name string, value interface{}) error {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("value invalid: %s, %T", name, value)
	}
	if _, ok := values[name]; ok {
		return fmt.Errorf("value duplicate: %s", name)
	}
	values[name] = value
	return nil
}
