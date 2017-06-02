package options

import (
	"flag"
	"os"
	"testing"
)

func ExampleUsage1() {
	f := flag.NewFlagSet("", flag.ContinueOnError)
	f.SetOutput(os.Stdout)

	v := values{
		B:   true,
		I:   -1,
		I8:  -8,
		I16: -16,
		I32: -32,
		I64: -64,
		U:   1,
		U8:  8,
		U16: 16,
		U32: 32,
		U64: 64,
		F32: 32.1,
		F64: 64.1,
		S:   "1",
	}
	v.Setup(f)
	f.Usage()

	// output:
	// Usage:
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

func ExampleUsage2() {
	f := flag.NewFlagSet("", flag.ContinueOnError)
	f.SetOutput(os.Stdout)

	v := values{
		B:   true,
		I:   -1,
		I8:  -8,
		I16: -16,
		I32: -32,
		I64: -64,
		U:   1,
		U8:  8,
		U16: 16,
		U32: 32,
		U64: 64,
		F32: 32.1,
		F64: 64.1,
		S:   "1",
	}
	Setup(f, "", "", &v)
	f.Usage()

	// output:
	// Usage:
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

func TestSetupByValues(t *testing.T) {
	f := flag.NewFlagSet("TestSetupByValues", flag.ContinueOnError)

	var v values
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

	if got, want := v.B, true; got != want {
		t.Errorf("bool: %v != %v", got, want)
	}
	if got, want := v.I, int(-1); got != want {
		t.Errorf("int: %v != %v", got, want)
	}
	if got, want := v.I8, int8(-8); got != want {
		t.Errorf("int8: %v != %v", got, want)
	}
	if got, want := v.I16, int16(-16); got != want {
		t.Errorf("int16: %v != %v", got, want)
	}
	if got, want := v.I32, int32(-32); got != want {
		t.Errorf("int32: %v != %v", got, want)
	}
	if got, want := v.I64, int64(-64); got != want {
		t.Errorf("int64: %v != %v", got, want)
	}
	if got, want := v.U, uint(1); got != want {
		t.Errorf("uint: %v != %v", got, want)
	}
	if got, want := v.U8, uint8(8); got != want {
		t.Errorf("uint8: %v != %v", got, want)
	}
	if got, want := v.U16, uint16(16); got != want {
		t.Errorf("uint16: %v != %v", got, want)
	}
	if got, want := v.U32, uint32(32); got != want {
		t.Errorf("uint32: %v != %v", got, want)
	}
	if got, want := v.U64, uint64(64); got != want {
		t.Errorf("uint64: %v != %v", got, want)
	}
	if got, want := v.Uptr, uintptr(1); got != want {
		t.Errorf("uintptr: %v != %v", got, want)
	}
	if got, want := v.F32, float32(32.1); got != want {
		t.Errorf("float32: %v != %v", got, want)
	}
	if got, want := v.F64, float64(64.1); got != want {
		t.Errorf("float64: %v != %v", got, want)
	}
	if got, want := v.S, "1"; got != want {
		t.Errorf("string: %v != %v", got, want)
	}
}

func TestSetupByStructs(t *testing.T) {
	f := flag.NewFlagSet("TestSetupByValues", flag.ContinueOnError)

	var v values
	if err := Setup(f, "", "", &v); err != nil {
		t.Fatalf("setup: %v", err)
	}

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

	if got, want := v.B, true; got != want {
		t.Errorf("bool: %v != %v", got, want)
	}
	if got, want := v.I, int(-1); got != want {
		t.Errorf("int: %v != %v", got, want)
	}
	if got, want := v.I8, int8(-8); got != want {
		t.Errorf("int8: %v != %v", got, want)
	}
	if got, want := v.I16, int16(-16); got != want {
		t.Errorf("int16: %v != %v", got, want)
	}
	if got, want := v.I32, int32(-32); got != want {
		t.Errorf("int32: %v != %v", got, want)
	}
	if got, want := v.I64, int64(-64); got != want {
		t.Errorf("int64: %v != %v", got, want)
	}
	if got, want := v.U, uint(1); got != want {
		t.Errorf("uint: %v != %v", got, want)
	}
	if got, want := v.U8, uint8(8); got != want {
		t.Errorf("uint8: %v != %v", got, want)
	}
	if got, want := v.U16, uint16(16); got != want {
		t.Errorf("uint16: %v != %v", got, want)
	}
	if got, want := v.U32, uint32(32); got != want {
		t.Errorf("uint32: %v != %v", got, want)
	}
	if got, want := v.U64, uint64(64); got != want {
		t.Errorf("uint64: %v != %v", got, want)
	}
	if got, want := v.Uptr, uintptr(1); got != want {
		t.Errorf("uintptr: %v != %v", got, want)
	}
	if got, want := v.F32, float32(32.1); got != want {
		t.Errorf("float32: %v != %v", got, want)
	}
	if got, want := v.F64, float64(64.1); got != want {
		t.Errorf("float64: %v != %v", got, want)
	}
	if got, want := v.S, "1"; got != want {
		t.Errorf("string: %v != %v", got, want)
	}
}
