package options

import (
	"flag"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unicode"

	"github.com/ironzhang/matrix/errs"
	"github.com/ironzhang/matrix/framework/pkg/tags"
)

func Setup(f *flag.FlagSet, name, usage string, value interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			if e, ok := r.(error); ok {
				err = e
				return
			}
			panic(r)
		}
	}()

	options{f}.SetupValue(name, usage, reflect.ValueOf(value).Elem())
	return
}

type options struct {
	*flag.FlagSet
}

func (o options) SetupValue(name, usage string, v reflect.Value) {
	switch k := v.Kind(); k {
	case reflect.Bool:
		o.Var(newBoolValue(v), name, usage)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		o.Var(newIntValue(v), name, usage)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		o.Var(newUintValue(v), name, usage)
	case reflect.Float32, reflect.Float64:
		o.Var(newFloatValue(v), name, usage)
	case reflect.String:
		o.Var(newStringValue(v), name, usage)
	case reflect.Struct:
		if name != "" {
			name = name + "."
		}
		if usage != "" {
			usage = usage + ": "
		}
		o.SetupStruct(name, usage, v)
	default:
		panic(errs.ErrorAt("options.SetupValue", fmt.Errorf("unsupport %s kind", k)))
	}
}

func (o options) SetupStruct(prefix, usage string, v reflect.Value) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous { // unexported
			continue
		}
		name := sf.Name
		if s := parseNameFromTag(sf.Tag); s != "" {
			name = s
		}
		if name == "-" {
			continue
		}
		o.SetupValue(prefix+name, usage+sf.Tag.Get("usage"), v.Field(i))
	}
}

func parseNameFromTag(tag reflect.StructTag) string {
	name, _ := tags.ParseTag(tag.Get("json"))
	if isValidTag(name) {
		return name
	}
	return ""
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
		}
	}
	return true
}
