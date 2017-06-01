package config

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestConfigSet(t *testing.T) {
	file := "cfg.json"
	type config1 struct {
		A int
		B string
	}
	type config2 struct {
		A int
		B config1
	}

	var c1 ConfigSet
	var c2 ConfigSet

	tests := []struct {
		name  string
		value interface{}
		got   interface{}
	}{
		{"config1", &config1{A: 1, B: "B"}, &config1{}},
		{"config2", &config2{A: 1, B: config1{A: 2, B: "b"}}, &config2{}},
	}

	for _, tt := range tests {
		if err := c1.Register(tt.name, tt.value); err != nil {
			t.Fatal(err)
		}
		if err := c2.Register(tt.name, tt.got); err != nil {
			t.Fatal(err)
		}
	}
	if err := c1.WriteToFile(file); err != nil {
		t.Fatal(err)
	}
	if err := c2.LoadFromFile(file); err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		if got, want := tt.got, tt.value; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %v != %v", tt.name, got, want)
		} else {
			t.Logf("%s: %v ", tt.name, got)
		}
	}

	os.Remove(file)
}

func TestByteSlice(t *testing.T) {
	var m map[string]byteSlice
	data := []byte(`{"m1": { "A": 1, "B": "B" }, "m2":{ "C": 1, "D": "D" }}`)
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		key   string
		value string
	}{
		{
			key:   "m1",
			value: `{ "A": 1, "B": "B" }`,
		},
		{
			key:   "m2",
			value: `{ "C": 1, "D": "D" }`,
		},
	}
	for _, tt := range tests {
		if got, want := string(m[tt.key]), tt.value; got != want {
			t.Errorf("%s: %v != %v", tt.key, got, want)
		} else {
			t.Logf("%s: %v\n", tt.key, got)
		}
	}
}

func TestJSONUnmarshal(t *testing.T) {
	type Value struct {
		A int
		B string
	}

	tests := []struct {
		value Value
		data  string
	}{
		{
			value: Value{A: 1},
			data:  `{"B": "B"}`,
		},
		{
			value: Value{B: "B"},
			data:  `{"A": 1}`,
		},
		{
			value: Value{A: 1, B: "B"},
			data:  `{"A": 1, "B": "B"}`,
		},
	}
	for i, tt := range tests {
		if err := json.Unmarshal([]byte(tt.data), &tt.value); err != nil {
			t.Fatal(err)
		}
		if got, want := tt.value.A, 1; got != want {
			t.Errorf("tests[%d]: A: %d != %d", i, got, want)
		}
		if got, want := tt.value.B, "B"; got != want {
			t.Errorf("tests[%d]: B: %s != %s", i, got, want)
		}
	}
}
