package framework

import (
	"encoding/json"
	"testing"
)

type mconf struct {
	A int
	B string
}

func TestLoadConfig(t *testing.T) {
	var c1 = config{}
	var c2 = config{}

	file := "./testdata/cfg.json"
	want := mconf{A: 1, B: "B"}
	c1.Register("test-module", &want)
	if err := c1.write(file); err != nil {
		t.Fatal(err)
	}

	var got mconf
	c2.Register("test-module", &got)
	if err := c2.load(file); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("mconf: %v != %v", got, want)
	} else {
		t.Logf("mconf: %v", got)
	}
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
