package options

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

type values struct {
	b    bool
	i    int
	i8   int8
	i16  int16
	i32  int32
	i64  int64
	u    uint
	u8   uint8
	u16  uint16
	u32  uint32
	u64  uint64
	uptr uintptr
	f32  float32
	f64  float64
	s    string
}

func (v *values) Setup(f *flag.FlagSet) {
	f.Var(newBoolValue(reflect.ValueOf(&v.b).Elem()), "bool", "usage of bool")
	f.Var(newIntValue(reflect.ValueOf(&v.i).Elem()), "int", "usage of int")
	f.Var(newIntValue(reflect.ValueOf(&v.i8).Elem()), "int8", "usage of int8")
	f.Var(newIntValue(reflect.ValueOf(&v.i16).Elem()), "int16", "usage of int16")
	f.Var(newIntValue(reflect.ValueOf(&v.i32).Elem()), "int32", "usage of int32")
	f.Var(newIntValue(reflect.ValueOf(&v.i64).Elem()), "int64", "usage of int64")
	f.Var(newUintValue(reflect.ValueOf(&v.u).Elem()), "uint", "usage of uint")
	f.Var(newUintValue(reflect.ValueOf(&v.u8).Elem()), "uint8", "usage of uint8")
	f.Var(newUintValue(reflect.ValueOf(&v.u16).Elem()), "uint16", "usage of uint16")
	f.Var(newUintValue(reflect.ValueOf(&v.u32).Elem()), "uint32", "usage of uint32")
	f.Var(newUintValue(reflect.ValueOf(&v.u64).Elem()), "uint64", "usage of uint64")
	f.Var(newUintValue(reflect.ValueOf(&v.uptr).Elem()), "uintptr", "usage of uintptr")
	f.Var(newFloatValue(reflect.ValueOf(&v.f32).Elem()), "float32", "usage of float32")
	f.Var(newFloatValue(reflect.ValueOf(&v.f64).Elem()), "float64", "usage of float64")
	f.Var(newStringValue(reflect.ValueOf(&v.s).Elem()), "string", "usage of string")
}

func ExampleUsage() {
	v := values{
		b:   true,
		i:   -1,
		i8:  -8,
		i16: -16,
		i32: -32,
		i64: -64,
		u:   1,
		u8:  8,
		u16: 16,
		u32: 32,
		u64: 64,
		f32: 32.1,
		f64: 64.1,
		s:   "1",
	}
	f := flag.NewFlagSet("TestValue", flag.ContinueOnError)
	f.SetOutput(os.Stdout)
	v.Setup(f)
	f.Usage()

	// output:
	// Usage of TestValue:
	//   -bool value
	//     	usage of bool (default true)
	//   -float32 value
	//     	usage of float32 (default 32.099998474121094)
	//   -float64 value
	//     	usage of float64 (default 64.1)
	//   -int value
	//     	usage of int (default -1)
	//   -int16 value
	//     	usage of int16 (default -16)
	//   -int32 value
	//     	usage of int32 (default -32)
	//   -int64 value
	//     	usage of int64 (default -64)
	//   -int8 value
	//     	usage of int8 (default -8)
	//   -string value
	//     	usage of string (default 1)
	//   -uint value
	//     	usage of uint (default 1)
	//   -uint16 value
	//     	usage of uint16 (default 16)
	//   -uint32 value
	//     	usage of uint32 (default 32)
	//   -uint64 value
	//     	usage of uint64 (default 64)
	//   -uint8 value
	//     	usage of uint8 (default 8)
	//   -uintptr value
	//     	usage of uintptr
}

func TestValues(t *testing.T) {
	var v values
	f := flag.NewFlagSet("TestValue", flag.ContinueOnError)
	v.Setup(f)

	args := []string{
		"-bool", "true",
		"-int", "-1",
		"-int8", "-8",
		"-int16", "-16",
		"-int32", "-32",
		"-int64", "-64",
		"-uint", "1",
		"-uint8", "8",
		"-uint16", "16",
		"-uint32", "32",
		"-uint64", "64",
		"-uintptr", "1",
		"-float32", "32.1",
		"-float64", "64.1",
		"-string", "1",
	}
	f.Parse(args)

	if got, want := v.b, true; got != want {
		t.Errorf("bool: %v != %v", got, want)
	}
	if got, want := v.i, int(-1); got != want {
		t.Errorf("int: %v != %v", got, want)
	}
	if got, want := v.i8, int8(-8); got != want {
		t.Errorf("int8: %v != %v", got, want)
	}
	if got, want := v.i16, int16(-16); got != want {
		t.Errorf("int16: %v != %v", got, want)
	}
	if got, want := v.i32, int32(-32); got != want {
		t.Errorf("int32: %v != %v", got, want)
	}
	if got, want := v.i64, int64(-64); got != want {
		t.Errorf("int64: %v != %v", got, want)
	}
	if got, want := v.u, uint(1); got != want {
		t.Errorf("uint: %v != %v", got, want)
	}
	if got, want := v.u8, uint8(8); got != want {
		t.Errorf("uint8: %v != %v", got, want)
	}
	if got, want := v.u16, uint16(16); got != want {
		t.Errorf("uint16: %v != %v", got, want)
	}
	if got, want := v.u32, uint32(32); got != want {
		t.Errorf("uint32: %v != %v", got, want)
	}
	if got, want := v.u64, uint64(64); got != want {
		t.Errorf("uint64: %v != %v", got, want)
	}
	if got, want := v.uptr, uintptr(1); got != want {
		t.Errorf("uintptr: %v != %v", got, want)
	}
	if got, want := v.f32, float32(32.1); got != want {
		t.Errorf("float32: %v != %v", got, want)
	}
	if got, want := v.f64, float64(64.1); got != want {
		t.Errorf("float64: %v != %v", got, want)
	}
	if got, want := v.s, "1"; got != want {
		t.Errorf("string: %v != %v", got, want)
	}
}
