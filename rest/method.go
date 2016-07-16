package rest

import (
	"errors"
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfNullInterface = reflect.TypeOf((*interface{})(nil)).Elem()

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

type method struct {
	value     reflect.Value
	argType   reflect.Type
	replyType reflect.Type
}

func parseMethod(i interface{}) (*method, error) {
	value := reflect.ValueOf(i)
	mtype := reflect.TypeOf(i)
	if mtype.Kind() != reflect.Func {
		return nil, errors.New("is not a func")
	}
	if mtype.NumIn() != 3 {
		return nil, errors.New("expect 3 in num")
	}
	if mtype.NumOut() != 1 {
		return nil, errors.New("expect 1 out num")
	}
	argType := mtype.In(1)
	if !isExportedOrBuiltinType(argType) {
		return nil, fmt.Errorf("arg type is not exported: %s", argType.String())
	}
	replyType := mtype.In(2)
	if replyType.Kind() != reflect.Ptr && replyType != typeOfNullInterface {
		return nil, fmt.Errorf("reply type unsupport: %s", replyType.String())
	}
	if !isExportedOrBuiltinType(replyType) {
		return nil, fmt.Errorf("reply type is not exported: %s", replyType.String())
	}
	retType := mtype.Out(0)
	if retType != typeOfError {
		return nil, fmt.Errorf("return argument type is not an error: %s", retType.String())
	}
	return &method{value: value, argType: argType, replyType: replyType}, nil
}

func (m *method) Call(ctx interface{}, argValue, replyValue reflect.Value) error {
	in := []reflect.Value{reflect.ValueOf(ctx), argValue, replyValue}
	out := m.value.Call(in)
	ret := out[0].Interface()
	if err, ok := ret.(error); ok {
		return err
	}
	return nil
}

func (m *method) NewArg() reflect.Value {
	if m.argType.Kind() == reflect.Ptr {
		return reflect.New(m.argType.Elem())
	}
	return reflect.New(m.argType)
}

func (m *method) NewReply() reflect.Value {
	if m.replyType.Kind() == reflect.Ptr {
		return reflect.New(m.replyType.Elem())
	}
	return reflect.New(m.replyType)
}

func (m *method) ArgIsNullInterface() bool {
	if m.argType == typeOfNullInterface {
		return true
	}
	return false
}

func (m *method) ReplyIsNullInterface() bool {
	if m.replyType == typeOfNullInterface {
		return true
	}
	return false
}
