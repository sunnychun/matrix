package values

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ironzhang/matrix/errs"
	"github.com/ironzhang/matrix/framework/pkg/tags"
)

type convertError struct {
	src reflect.Type
	dst reflect.Type
	err error
}

func (e convertError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("can not convert %s@%s to %s@%s", e.src.Kind(), e.src, e.dst.Kind(), e.dst, e.err)
	}
	return fmt.Sprintf("can not convert %s@%s to %s@%s", e.src.Kind(), e.src, e.dst.Kind(), e.dst)
}

// i => x
func set(i, x reflect.Value) {
	i = indirect(i)
	switch k := i.Kind(); k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		setValue(i, x)

	case reflect.String:
		setString(i, x)

	case reflect.Array, reflect.Slice:
		setArray(i, x)

	case reflect.Map:
		setMap(i, x)

	case reflect.Struct:
		setObject(i, x)

	default:
		panic(errs.ErrorAt("set", convertError{i.Type(), x.Type(), nil}))
	}
}

func setValue(i, x reflect.Value) {
	x = indirect(x)
	x.Set(i)
}

func setString(i, x reflect.Value) {
	u, ut, v := indirectUnmarshaler(x)
	if u != nil {
		b := []byte(`"` + i.String() + `"`)
		if err := u.UnmarshalJSON(b); err != nil {
			panic(errs.ErrorAt("setString", convertError{i.Type(), x.Type(), err}))
		}
		return
	}
	if ut != nil {
		b := []byte(i.String())
		if err := ut.UnmarshalText(b); err != nil {
			panic(errs.ErrorAt("setString", convertError{i.Type(), x.Type(), err}))
		}
		return
	}
	v.Set(i)
}

func setArray(i, x reflect.Value) {
	x = indirect(x)
	switch k := x.Kind(); k {
	case reflect.Interface:
		if x.NumMethod() == 0 {
			x.Set(i)
		} else {
			panic(errs.ErrorAt("setArray", convertError{i.Type(), x.Type(), nil}))
		}

	case reflect.Array:
		for n := 0; n < x.Len(); n++ {
			v := x.Index(n)
			z := reflect.Zero(v.Type())
			v.Set(z)
		}
		for n := 0; n < i.Len() && n < x.Len(); n++ {
			set(i.Index(n), x.Index(n))
		}

	case reflect.Slice:
		s := reflect.MakeSlice(x.Type(), i.Len(), i.Len())
		for n := 0; n < i.Len(); n++ {
			set(i.Index(n), s.Index(n))
		}
		x.Set(s)

	default:
		panic(errs.ErrorAt("setArray", convertError{i.Type(), x.Type(), nil}))
	}
}

func setMap(i, x reflect.Value) {
	x = indirect(x)
	switch k := x.Kind(); k {
	case reflect.Interface:
		if x.NumMethod() == 0 {
			x.Set(i)
		} else {
			panic(errs.ErrorAt("setMap", convertError{i.Type(), x.Type(), nil}))
		}

	case reflect.Map:
		nm := reflect.MakeMap(x.Type())
		for _, k := range i.MapKeys() {
			nm.SetMapIndex(k, i.MapIndex(k))
		}
		x.Set(nm)

	case reflect.Struct:
		if i.Type().Key().Kind() != reflect.String {
			panic(errs.ErrorAt("setMap", convertError{i.Type(), x.Type(), nil}))
		}
		fs := typeFields(x.Type())
		for _, k := range i.MapKeys() {
			f, ok := fs[k.String()]
			if !ok {
				continue
			}
			if f.readonly {
				continue
			}
			set(i.MapIndex(k), x.FieldByIndex(f.index))
		}

	default:
		panic(errs.ErrorAt("setMap", convertError{i.Type(), x.Type(), nil}))
	}
}

func setObject(i, x reflect.Value) {
	x = indirect(x)
	switch k := x.Kind(); k {
	case reflect.Interface:
		if x.NumMethod() == 0 {
			x.Set(i)
		} else {
			panic(errs.ErrorAt("setObject", convertError{i.Type(), x.Type(), nil}))
		}

	case reflect.Struct:
		xfs := typeFields(x.Type())
		ifs := typeFields(i.Type())
		for k, f := range ifs {
			xf, ok := xfs[k]
			if !ok {
				continue
			}
			if xf.readonly {
				continue
			}
			set(i.FieldByIndex(f.index), x.FieldByIndex(xf.index))
		}

	case reflect.Map:
		if x.Type().Key().Kind() != reflect.String {
			panic(errs.ErrorAt("setObject", convertError{i.Type(), x.Type(), nil}))
		}
		nm := reflect.MakeMap(x.Type())
		fs := typeFields(i.Type())
		for k, f := range fs {
			nm.SetMapIndex(reflect.ValueOf(k), i.FieldByIndex(f.index))
		}
		x.Set(nm)
	}
}

func indirect(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			v = v.Elem()
			continue
		}
		if v.Kind() != reflect.Ptr {
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}

func indirectUnmarshaler(v reflect.Value) (json.Unmarshaler, encoding.TextUnmarshaler, reflect.Value) {
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() {
				v = e
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 {
			if u, ok := v.Interface().(json.Unmarshaler); ok {
				return u, nil, reflect.Value{}
			}
			if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
				return nil, u, reflect.Value{}
			}
		}
		v = v.Elem()
	}
	return nil, nil, v
}

type field struct {
	name     string
	index    []int
	typ      reflect.Type
	readonly bool
}

type fields map[string]field

func parseTag(tag reflect.StructTag) (string, bool) {
	if s := tag.Get("json"); s != "" {
		name, opts := tags.ParseTag(s)
		if !tags.IsValidTag(name) {
			return "", opts.Contains("readonly")
		}
		return name, opts.Contains("readonly")
	}
	return "", false
}

func typeFields(t reflect.Type) fields {
	fs := make(fields)
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous { // unexported
			continue
		}

		name := sf.Name
		tname, readonly := parseTag(sf.Tag)
		if tname != "" {
			name = tname
		}
		if name == "-" {
			continue
		}
		fs[name] = field{name: name, index: sf.Index, typ: sf.Type, readonly: readonly}
	}
	return fs
}
