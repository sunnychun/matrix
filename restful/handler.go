package restful

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()
var typeOfNilInterface = reflect.TypeOf((*interface{})(nil)).Elem()

func isExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return isExported(t.Name()) || t.PkgPath() == ""
}

func checkIns(ftype reflect.Type) (in0, in1, in2 reflect.Type, err error) {
	if ftype.NumIn() != 3 {
		err = fmt.Errorf("func has wrong number of ins: %d", ftype.NumIn())
		return
	}
	in0 = ftype.In(0)
	if !in0.Implements(typeOfContext) {
		err = fmt.Errorf("in0 type not implements context: %s", in0)
		return
	}
	in1 = ftype.In(1)
	if !isExportedOrBuiltinType(in1) {
		err = fmt.Errorf("in1 type not exported: %s", in1)
		return
	}
	in2 = ftype.In(2)
	if in2.Kind() != reflect.Ptr && in2 != typeOfNilInterface {
		err = fmt.Errorf("in2 type not a pointer or interface{}: %s", in2)
		return
	}
	if !isExportedOrBuiltinType(in2) {
		err = fmt.Errorf("in2 type not exported: %s", in2)
		return
	}
	return
}

func checkOuts(ftype reflect.Type) error {
	if ftype.NumOut() != 1 {
		return fmt.Errorf("func has wrong number of outs: %d", ftype.NumOut())
	}
	if out0 := ftype.Out(0); out0 != typeOfError {
		return fmt.Errorf("func returns %s not error", out0)
	}
	return nil
}

type handler struct {
	value   reflect.Value
	in1Type reflect.Type
	in2Type reflect.Type
}

func parseHandler(i interface{}) (*handler, error) {
	value := reflect.ValueOf(i)
	ftype := reflect.TypeOf(i)
	if ftype.Kind() != reflect.Func {
		return nil, errors.New("handler kind not func")
	}
	_, in1Type, in2Type, err := checkIns(ftype)
	if err != nil {
		return nil, err
	}
	if err = checkOuts(ftype); err != nil {
		return nil, err
	}
	return &handler{value: value, in1Type: in1Type, in2Type: in2Type}, nil
}

func (h *handler) Handle(ctx context.Context, in1, in2 reflect.Value) error {
	in := []reflect.Value{reflect.ValueOf(ctx), in1, in2}
	out := h.value.Call(in)
	ret := out[0].Interface()
	if err, ok := ret.(error); ok {
		return err
	}
	panic("handler returns not error")
}

func reflectNewValue(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		return reflect.New(t.Elem())
	}
	return reflect.New(t)
}

func reflectIsNilInterface(t reflect.Type) bool {
	return t == typeOfNilInterface
}
