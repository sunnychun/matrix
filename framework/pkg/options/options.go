package options

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"runtime"

	"github.com/ironzhang/matrix/framework/pkg/tags"
)

var CommandLine = flag.CommandLine

func Setup(options interface{}) (err error) {
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

	setup(reflect.ValueOf(options).Elem())
	return CommandLine.Parse(os.Args[1:])
}

func setup(val reflect.Value) {
	setupValue("", "", val)
}

func setupValue(name, usage string, val reflect.Value) {
	switch k := val.Kind(); k {
	case reflect.Bool:
		CommandLine.Var(newBoolValue(val), name, usage)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		CommandLine.Var(newIntValue(val), name, usage)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		CommandLine.Var(newUintValue(val), name, usage)
	case reflect.Float32, reflect.Float64:
		CommandLine.Var(newFloatValue(val), name, usage)
	case reflect.String:
		CommandLine.Var(newStringValue(val), name, usage)
	case reflect.Struct:
		setupStruct(name, val)
	case reflect.Map:
		setupMap(name, val)
	default:
		panic(fmt.Errorf("unsupport %s kind", k))
	}
}

func setupStruct(prefix string, val reflect.Value) {
	if prefix != "" {
		prefix += "."
	}
	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		name := sf.Name
		if v := parseNameFromTag(sf.Tag); v != "" {
			name = v
		}
		if name == "-" {
			continue
		}
		usage := parseUsageFromTag(sf.Tag)
		setupValue(prefix+name, usage, val.Field(i))
	}
}

func setupMap(prefix string, val reflect.Value) {
	if prefix != "" {
		prefix += "."
	}
	keys := val.MapKeys()
	for _, key := range keys {
		setupValue(prefix+key.String(), "", val.MapIndex(key))
	}
}

func parseNameFromTag(tag reflect.StructTag) string {
	v, _ := tags.ParseTag(tag.Get("json"))
	return v
}

func parseUsageFromTag(tag reflect.StructTag) string {
	v, _ := tags.ParseTag(tag.Get("usage"))
	return v
}
