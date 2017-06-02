package options

import (
	"reflect"
	"testing"
)

func TestValue(t *testing.T) {
	var i int
	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64

	CommandLine.Var(newIntValue(reflect.ValueOf(&i).Elem()), "int", "usage")
	CommandLine.Var(newIntValue(reflect.ValueOf(&i8).Elem()), "int8", "usage")
	CommandLine.Var(newIntValue(reflect.ValueOf(&i16).Elem()), "int16", "usage")
	CommandLine.Var(newIntValue(reflect.ValueOf(&i32).Elem()), "int32", "usage")
	CommandLine.Var(newIntValue(reflect.ValueOf(&i64).Elem()), "int64", "usage")

	args := []string{
		"-int", "1",
		"-int8", "8",
		"-int16", "16",
		"-int32", "32",
		"-int64", "64",
	}
	CommandLine.Parse(args)

	if got, want := i, 1; got != want {
		t.Errorf("i: %d != %d", got, want)
	}
	if got, want := i8, int8(8); got != want {
		t.Errorf("i8: %d != %d", got, want)
	}
	if got, want := i16, int16(16); got != want {
		t.Errorf("i16: %d != %d", got, want)
	}
	if got, want := i32, int32(32); got != want {
		t.Errorf("i32: %d != %d", got, want)
	}
	if got, want := i64, int64(64); got != want {
		t.Errorf("i64: %d != %d", got, want)
	}
}

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
