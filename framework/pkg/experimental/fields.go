package experimental

import "reflect"

func parseJSONTag(tag string) (string, bool) {
	if tag != "" {
		name, opts := parseTag(tag)
		if isValidTag(name) {
			return name, opts.Contains("readonly")
		}
		return "", opts.Contains("readonly")
	}
	return "", false
}

type field struct {
	name     string
	index    int
	typ      reflect.Type
	readonly bool
}

func typeFields(t reflect.Type) map[string]field {
	var name string
	var readonly bool
	fields := make(map[string]field)
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous { // unexported
			continue
		}
		name, readonly = parseJSONTag(sf.Tag.Get("json"))
		if name == "-" {
			continue
		}
		if name == "" {
			name = sf.Name
		}
		fields[name] = field{name: name, index: i, typ: sf.Type, readonly: readonly}
	}
	return fields
}
