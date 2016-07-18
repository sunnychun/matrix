package rest

import (
	"errors"
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"

	"golang.org/x/net/context"
)

var typeOfVars = reflect.TypeOf(Vars{})
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()
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
	if mtype.NumIn() != 4 {
		return nil, errors.New("expect 4 in num")
	}
	if mtype.NumOut() != 1 {
		return nil, errors.New("expect 1 out num")
	}

	arg0Type := mtype.In(0)
	if arg0Type != typeOfContext {
		return nil, fmt.Errorf("argument0 %q is not a %q", arg0Type.String(), typeOfContext.String())
	}
	arg1Type := mtype.In(1)
	if arg1Type != typeOfVars {
		return nil, fmt.Errorf("argument1 %q is not a %q", arg1Type.String(), typeOfVars.String())
	}
	arg2Type := mtype.In(2)
	if arg2Type.Kind() != reflect.Ptr && arg2Type != typeOfNullInterface {
		return nil, fmt.Errorf("argument2 %q is not a pointer or interface{}", arg2Type.String())
	}
	if !isExportedOrBuiltinType(arg2Type) {
		return nil, fmt.Errorf("argument2 %q is not exported", arg2Type.String())
	}
	arg3Type := mtype.In(3)
	if arg3Type.Kind() != reflect.Ptr && arg3Type != typeOfNullInterface {
		return nil, fmt.Errorf("argument3 %q is not a pointer or interface{}", arg3Type.String())
	}
	if !isExportedOrBuiltinType(arg3Type) {
		return nil, fmt.Errorf("argument3 %q is not exported", arg3Type.String())
	}
	retType := mtype.Out(0)
	if retType != typeOfError {
		return nil, fmt.Errorf("return argument %q is not an error", retType.String())
	}
	return &method{value: value, argType: arg2Type, replyType: arg3Type}, nil
}

func (m *method) Call(ctx context.Context, vars Vars, argValue, replyValue reflect.Value) error {
	in := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(vars), argValue, replyValue}
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
