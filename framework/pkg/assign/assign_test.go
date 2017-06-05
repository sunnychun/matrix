package assign

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestTypeFields(t *testing.T) {
	type t1 struct {
		A int
	}
	type t2 struct {
		A t1
		B string
	}

	tests := []struct {
		v    interface{}
		num  int
		keys []string
	}{
		{
			v:    &t1{},
			num:  1,
			keys: []string{"A"},
		},
		{
			v:    &t2{},
			num:  2,
			keys: []string{"A", "B"},
		},
	}
	for i, tt := range tests {
		fs := typeFields(reflect.TypeOf(tt.v).Elem())
		if got, want := len(fs), tt.num; got != want {
			t.Errorf("tests[%d]: len: %d != %d", i, got, want)
		}
		for _, key := range tt.keys {
			if _, ok := fs[key]; !ok {
				t.Errorf("tests[%d]: key: %s not found", i, key)
			}
		}
	}
}

func TestMapAssign(t *testing.T) {
	type t1 struct {
		A int
	}
	type t2 struct {
		A t1
		B string
	}
	tests := []struct {
		m    map[string]interface{}
		got  interface{}
		want interface{}
	}{
		{
			m:    map[string]interface{}{"A": 1},
			got:  &t1{},
			want: &t1{A: 1},
		},
		{
			m: map[string]interface{}{
				"A": map[string]interface{}{"A": 1},
				"B": "B",
			},
			got: &t2{},
			want: &t2{
				A: t1{A: 1},
				B: "B",
			},
		},
		{
			m:    map[string]interface{}{"A": 1, "B": 2},
			got:  &map[string]int{"C": 1},
			want: &map[string]int{"A": 1, "B": 2},
		},
	}
	for i, tt := range tests {
		x := reflect.ValueOf(tt.got).Elem()
		mapAssign(x, tt.m)
		if got, want := tt.got, tt.want; !reflect.DeepEqual(got, want) {
			t.Errorf("tests[%d]: %v != %v", i, got, want)
		} else {
			t.Logf("tests[%d]: %v == %v", i, got, want)
		}
	}
}

func TestValueAssign(t *testing.T) {
	type t1 struct {
		A int
		B string
	}

	var b bool
	var i8 int8
	var i16 int16
	tests := []struct {
		v    interface{}
		got  interface{}
		want interface{}
	}{
		{v: bool(true), got: &b, want: bool(true)},
		{v: int8(8), got: &i8, want: int8(8)},
		{v: int16(16), got: &i16, want: int16(16)},
		{v: map[string]interface{}{"A": 1, "B": "B"}, got: &t1{}, want: t1{A: 1, B: "B"}},
	}
	for i, tt := range tests {
		x := reflect.ValueOf(tt.got).Elem()
		valueAssign(x, tt.v)
		if got, want := x.Interface(), tt.want; !reflect.DeepEqual(got, want) {
			t.Errorf("tests[%d]: %v != %v", i, got, want)
		} else {
			t.Logf("tests[%d]: %v == %v", i, got, want)
		}
	}
}

func TestAssign(t *testing.T) {
	type T1 struct {
		A int
	}
	type T2 struct {
		A int
		B int
	}
	type T3 struct {
		A int
		B int `json:"b,readonly"`
	}
	type T4 struct {
		T1 T1
		T2 T2
		T3 T3
	}

	tests := []struct {
		value interface{}
		got   interface{}
		want  interface{}
	}{
		{
			value: &T1{A: 1},
			got:   &T1{},
			want:  &T1{A: 1},
		},
		{
			value: &T2{A: 1, B: 2},
			got:   &T2{B: 1},
			want:  &T2{A: 1, B: 2},
		},
		{
			value: &T3{A: 1, B: 2},
			got:   &T3{B: 1},
			want:  &T3{A: 1, B: 1},
		},
		{
			value: &T4{T1: T1{A: 1}, T2: T2{A: 1, B: 2}, T3: T3{A: 1, B: 2}},
			got:   &T4{},
			want:  &T4{T1: T1{A: 1}, T2: T2{A: 1, B: 2}, T3: T3{A: 1, B: 0}},
		},
	}
	for i, tt := range tests {
		data, err := json.Marshal(tt.value)
		if err != nil {
			t.Fatal(err)
		}
		var m map[string]interface{}
		err = json.Unmarshal(data, &m)
		if err != nil {
			t.Fatal(err)
		}
		//fmt.Printf("%s\n", data)

		if err = Assign(tt.got, m); err != nil {
			t.Errorf("tests[%d]: assign to: %v", i, err)
			continue
		}
		if got, want := tt.got, tt.want; !reflect.DeepEqual(got, want) {
			t.Errorf("tests[%d]: %v != %v", i, got, want)
		} else {
			t.Logf("tests[%d]: %v == %v", i, got, want)
		}
	}
}
