package options

import (
	"reflect"
	"testing"
)

func TestSetup(t *testing.T) {
	type T1 struct {
		A int
		B string `json:"b",usage:"usage of B"`
	}

	t1 := T1{
		A: 1,
		B: "B",
	}

	setup(reflect.ValueOf(&t1).Elem())
	args := []string{
		"-A", "2",
		"-b", "C",
	}
	if err := CommandLine.Parse(args); err != nil {
		t.Fatal(err)
	}

	if got, want := t1.A, 2; got != want {
		t.Errorf("A: %v != %v", got, want)
	}
	if got, want := t1.B, "C"; got != want {
		t.Errorf("B: %v != %v", got, want)
	}
}
