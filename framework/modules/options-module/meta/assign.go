package meta

import (
	"fmt"
	"reflect"
	"runtime"
)

type assignError struct {
	Method string
	Reason string
}

func (e assignError) Error() string {
	return fmt.Sprintf("%s: %s", e.Method, e.Reason)
}

type field struct {
	name     string
	index    []int
	typ      reflect.Type
	readonly bool
}

type fields map[string]field

func typeFields(t reflect.Type) fields {
	var name string
	var readonly bool
	fs := make(fields)
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		name = sf.Name
		readonly = false
		tag := sf.Tag.Get("json")
		if tag != "" {
			var opts tagOptions
			name, opts = parseTag(tag)
			readonly = opts.Contains("readonly")
		}
		fs[name] = field{name: name, index: sf.Index, typ: sf.Type, readonly: readonly}
	}
	return fs
}

func floatAssign(x reflect.Value, f float64) {
	switch kind := x.Kind(); kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x.SetInt(int64(f))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		x.SetUint(uint64(f))
	case reflect.Float32, reflect.Float64:
		x.SetFloat(f)
	default:
		panic(assignError{"config.floatAssign", fmt.Sprintf("unsupport %s kind", kind)})
	}
}

func mapAssign(x reflect.Value, m map[string]interface{}) {
	switch kind := x.Kind(); kind {
	case reflect.Map:
		nm := reflect.MakeMap(x.Type())
		for k, v := range m {
			nm.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
		}
		x.Set(nm)
	case reflect.Struct:
		fs := typeFields(x.Type())
		for k, v := range m {
			f, ok := fs[k]
			if !ok {
				continue
			}
			if f.readonly {
				continue
			}
			valueAssign(x.FieldByIndex(f.index), v)
		}
	default:
		panic(assignError{"config.mapAssign", fmt.Sprintf("unsupport %s kind", kind)})
	}
}

func valueAssign(x reflect.Value, value interface{}) {
	switch v := value.(type) {
	case bool:
		x.SetBool(v)

	case int:
		x.SetInt(int64(v))
	case int8:
		x.SetInt(int64(v))
	case int16:
		x.SetInt(int64(v))
	case int32:
		x.SetInt(int64(v))
	case int64:
		x.SetInt(int64(v))

	case uint:
		x.SetUint(uint64(v))
	case uint8:
		x.SetUint(uint64(v))
	case uint16:
		x.SetUint(uint64(v))
	case uint32:
		x.SetUint(uint64(v))
	case uint64:
		x.SetUint(uint64(v))
	case uintptr:
		x.SetUint(uint64(v))

	case float32:
		floatAssign(x, float64(v))
	case float64:
		floatAssign(x, float64(v))

	case string:
		x.SetString(v)

	case map[string]interface{}:
		mapAssign(x, v)

	default:
		panic(assignError{"config.valueAssign", fmt.Sprintf("unsupport %T type", v)})
	}
}

func Assign(x interface{}, m map[string]interface{}) (err error) {
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

	mapAssign(reflect.ValueOf(x).Elem(), m)
	return
}
