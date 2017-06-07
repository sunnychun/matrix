package values

import (
	"reflect"
	"testing"
	"time"

	"github.com/ironzhang/matrix/jsoncfg"
)

func TestSetString(t *testing.T) {
	{
		var i string
		var x string
		i = "hello"
		setString(reflect.ValueOf(i), reflect.ValueOf(&x))
		if got, want := x, i; got != want {
			t.Errorf("case0: %s != %s", got, want)
		} else {
			t.Logf("case0: %s == %s", got, want)
		}
	}

	{
		var i string
		var x jsoncfg.Duration
		i = "10m30s"
		setString(reflect.ValueOf(i), reflect.ValueOf(&x))
		if got, want := time.Duration(x), 10*time.Minute+30*time.Second; got != want {
			t.Errorf("case1: %s != %s", got, want)
		} else {
			t.Logf("case1: %s == %s", got, want)
		}
	}
}

func TestSetArray(t *testing.T) {
	if true {
		var i = [10]int{1, 2, 3}
		var x = [5]int{2, 4}
		setArray(reflect.ValueOf(i), reflect.ValueOf(&x))
		if got, want := x, [5]int{1, 2, 3}; got != want {
			t.Fatalf("case0: %v != %v", got, want)
		} else {
			t.Logf("case0: %v == %v", got, want)
		}
	}

	if true {
		var i = []int{1, 2, 3}
		var x = [2]int{2, 4}
		setArray(reflect.ValueOf(i), reflect.ValueOf(&x))
		if got, want := x, [2]int{1, 2}; got != want {
			t.Fatalf("case1: %v != %v", got, want)
		} else {
			t.Logf("case1: %v == %v", got, want)
		}
	}

	if true {
		var i = [10]int{1, 2, 3}
		var x = []int{2, 4}
		setArray(reflect.ValueOf(i), reflect.ValueOf(&x))
		if got, want := x, i[:]; !reflect.DeepEqual(got, want) {
			t.Fatalf("case2: %v != %v", got, want)
		} else {
			t.Logf("case2: %v == %v", got, want)
		}
	}

	if true {
		var i = []int{1, 2, 3}
		var x = []int{2, 4}
		setArray(reflect.ValueOf(i), reflect.ValueOf(&x))
		if got, want := x, i[:]; !reflect.DeepEqual(got, want) {
			t.Errorf("case3: %v != %v", got, want)
		} else {
			t.Logf("case3: %v == %v", got, want)
		}
	}

	if true {
		var i = []int{1, 2, 3}
		var x interface{}
		setArray(reflect.ValueOf(i), reflect.ValueOf(&x))
		if got, want := x, i[:]; !reflect.DeepEqual(got, want) {
			t.Errorf("case4: %v != %v", got, want)
		} else {
			t.Logf("case4: %v == %v", got, want)
		}
	}

	if true {
		var i = []int{1, 2, 3}
		var x []interface{}
		setArray(reflect.ValueOf(i), reflect.ValueOf(&x))
		if got, want := x, []interface{}{1, 2, 3}; !reflect.DeepEqual(got, want) {
			t.Errorf("case5: %v != %v", got, want)
		} else {
			t.Logf("case5: %v == %v", got, want)
		}
	}
}

func TestSetMap(t *testing.T) {
	if true {
		var i = map[int]string{1: "1", 2: "2"}
		var x interface{}
		setMap(reflect.ValueOf(i), reflect.ValueOf(&x))
		if got, want := x, i; !reflect.DeepEqual(got, want) {
			t.Errorf("case0: %v != %v", got, want)
		} else {
			t.Logf("case0: %v == %v", got, want)
		}
	}

	if true {
		var i = map[int]string{1: "1", 2: "2"}
		var x map[int]string
		setMap(reflect.ValueOf(i), reflect.ValueOf(&x))
		if got, want := x, i; !reflect.DeepEqual(got, want) {
			t.Errorf("case1: %v != %v", got, want)
		} else {
			t.Logf("case1: %v == %v", got, want)
		}
	}

	if true {
		type T struct {
			A string
			B int
		}

		var i = map[string]interface{}{"A": "1", "B": 2}
		var x T
		var t1 = T{A: "1", B: 2}
		setMap(reflect.ValueOf(i), reflect.ValueOf(&x))
		if got, want := x, t1; !reflect.DeepEqual(got, want) {
			t.Errorf("case2: %v != %v", got, want)
		} else {
			t.Logf("case2: %v == %v", got, want)
		}
	}

	if true {
		type T struct {
			A map[string]int
			B map[string]string
			C map[int]int
		}

		var i = map[string]interface{}{
			"A": map[string]int{
				"A1": 1,
				"A2": 2,
			},
			"B": map[string]string{
				"B1": "1",
				"B2": "2",
			},
			"C": map[int]int{
				1: 1,
				2: 2,
			},
		}
		var x T
		var t1 = T{
			A: map[string]int{
				"A1": 1,
				"A2": 2,
			},
			B: map[string]string{
				"B1": "1",
				"B2": "2",
			},
			C: map[int]int{
				1: 1,
				2: 2,
			},
		}
		setMap(reflect.ValueOf(i), reflect.ValueOf(&x))
		if got, want := x, t1; !reflect.DeepEqual(got, want) {
			t.Errorf("case2: %v != %v", got, want)
		} else {
			t.Logf("case2: %v == %v", got, want)
		}
	}
}

func TestSetObject(t *testing.T) {
	if true {
		type T struct {
			A int
			B string
		}
		var i = T{A: 1, B: "B"}
		var x interface{}
		setObject(reflect.ValueOf(i), reflect.ValueOf(&x))
		if got, want := x, i; !reflect.DeepEqual(got, want) {
			t.Errorf("case0: %v != %v", got, want)
		} else {
			t.Logf("case0: %v == %v", got, want)
		}
	}

	if true {
		type T1 struct {
			A int
			B string
		}

		type T2 struct {
			X int    `json:"A"`
			Y string `json:"B"`
		}

		var i = T1{A: 1, B: "B"}
		var x T2
		setObject(reflect.ValueOf(i), reflect.ValueOf(&x))
		if got, want := x.X, i.A; got != want {
			t.Errorf("case1: %v != %v", got, want)
		} else {
			t.Logf("case1: %v == %v", got, want)
		}
		if got, want := x.Y, i.B; got != want {
			t.Errorf("case2: %v != %v", got, want)
		} else {
			t.Logf("case1: %v == %v", got, want)
		}
	}

	if true {
		type T struct {
			A int
			B string
		}
		var i = T{A: 1, B: "B"}
		var x map[string]interface{}
		var t1 = map[string]interface{}{
			"A": 1,
			"B": "B",
		}
		setObject(reflect.ValueOf(i), reflect.ValueOf(&x))
		if got, want := x, t1; !reflect.DeepEqual(got, want) {
			t.Errorf("case0: %v != %v", got, want)
		} else {
			t.Logf("case0: %v == %v", got, want)
		}
	}
}
