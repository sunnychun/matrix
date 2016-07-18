package rest

import (
	"fmt"
	"net/url"
	"testing"

	"golang.org/x/net/context"
)

func TestParseMethod(t *testing.T) {
	f00 := func(_ context.Context, a1 Vars, a2 *int, a3 *int) error { return nil }
	f01 := func(_ context.Context, a1 Vars, a2 interface{}, a3 interface{}) error { return nil }
	f02 := func(_ context.Context, a1 Vars, a2 *interface{}, a3 *interface{}) error { return nil }

	f10 := func(_ interface{}, a1 string, a2 int, a3 int) error { return nil }
	f20 := func(_ context.Context, a1 string, a2 int, a3 int) error { return nil }
	f21 := func(_ context.Context, a1 url.Values, a2 *int, a3 int) error { return nil }
	f22 := func(_ context.Context, a1 map[string][]string, a2 *int, a3 int) error { return nil }
	f30 := func(_ context.Context, a1 Vars, a2 int, a3 int) error { return nil }
	f40 := func(_ context.Context, a1 Vars, a2 *int, a3 int) error { return nil }

	type testcase struct {
		api interface{}
		err bool
	}

	testcases := []testcase{
		{f00, false},
		{f01, false},
		{f02, false},

		{f10, true},
		{f20, true},
		{f21, true},
		{f22, true},
		{f30, true},
		{f40, true},
	}
	for i, tc := range testcases {
		_, err := parseMethod(tc.api)
		if err != nil && tc.err {
			t.Logf("case%d expect: %v", i, err)
		} else if err == nil && !tc.err {
			t.Logf("case%d expect: success", i)
		} else {
			t.Errorf("case%d unexpect: err[%v, %t]", i, err, tc.err)
		}
	}
}

func TestCallMethod(t *testing.T) {
	f00 := func(_ context.Context, a1 Vars, a2 *int, a3 *int) error { return nil }
	f01 := func(_ context.Context, a1 Vars, a2 *int, a3 *int) error { return nil }
	f02 := func(_ context.Context, a1 Vars, a2 interface{}, a3 interface{}) error { return nil }
	f03 := func(_ context.Context, a1 Vars, a2 *interface{}, a3 *interface{}) error { return nil }
	f04 := func(_ context.Context, a1 Vars, a2 *interface{}, a3 *interface{}) error {
		return fmt.Errorf("f04 error: a=%s, b=%s", a1.Get("a"), a1.Get("b"))
	}

	type testcase struct {
		api interface{}
		err bool
	}

	testcases := []testcase{
		{f00, false},
		{f01, false},
		{f02, false},
		{f03, false},
		{f04, true},
	}

	ctx := context.Background()
	values, err := url.ParseQuery("a=3&a=2&a=banana;b=c")
	if err != nil {
		t.Fatalf("parse query: %v", err)
	}

	for i, tc := range testcases {
		m, err := parseMethod(tc.api)
		if err != nil {
			t.Fatalf("case%d parse method: %v", i, err)
		}
		err = m.Call(ctx, Vars(values), m.NewArg(), m.NewReply())
		if err != nil && tc.err {
			t.Logf("case%d expect: %v", i, err)
		} else if err == nil && !tc.err {
			t.Logf("case%d expect: success", i)
		} else {
			t.Errorf("case%d unexpect: err[%v, %t]", i, err, tc.err)
		}
	}
}

/*
func testfunc0(a, b int, c *int) error {
	fmt.Printf("%d, %d, %d\n", a, b, *c)
	return nil
}

type handler struct {
	n int
}

func (h *handler) testfunc(a, b int, c *int) error {
	fmt.Printf("[%d] %d, %d, %d\n", h.n, a, b, *c)
	return nil
}

func TestMethod1(t *testing.T) {
	testcases := []interface{}{
		testfunc0,
		//&handler{},
		(&handler{1}).testfunc,
	}

	for i, tc := range testcases {
		m, err := parseMethod(tc)
		if err != nil {
			t.Errorf("%d parse failed: err[%v]", i, err)
			return
		}
		if err = m.Call(1, reflect.ValueOf(2), reflect.ValueOf(&i)); err != nil {
			t.Errorf("%d call failed: err[%v]", i, err)
			return
		}
	}
}

type Test1Req struct {
	A, B int
}

type Test1Resp struct {
	A, B string
}

func handleTest1(ctx *ac.ZContext, req *Test1Req, resp *Test1Resp) error {
	resp.A = strconv.Itoa(req.A)
	resp.B = strconv.Itoa(req.B)
	fmt.Println(req.A, req.B, resp.A, resp.B)
	return nil
}

func TestMethod2(t *testing.T) {
	m, err := parseMethod(handleTest1)
	if err != nil {
		t.Errorf("parse failed: err[%v]", err)
		return
	}

	argv := m.NewArg()
	err = json.Unmarshal([]byte(`{"A":1, "B":2}`), argv.Interface())
	if err != nil {
		t.Errorf("json unmarshal failed: err[%v]", err)
		return
	}

	replyv := m.NewReply()
	ctx := (*ac.ZContext)(nil)
	err = m.Call(ctx, argv, replyv)
	if err != nil {
		t.Errorf("call failed: err[%v]", err)
		return
	}

	data, err := json.Marshal(replyv.Interface())
	if err != nil {
		t.Errorf("json marshal failed: err[%v]", err)
		return
	}

	t.Log(string(data))
}

func TestNullInterface(t *testing.T) {
	f := func(a int, b interface{}, c *interface{}) error {
		return nil
	}

	m, err := parseMethod(f)
	if err != nil {
		t.Errorf("parse failed: err[%v]", err)
		return
	}

	if m.ArgIsNullInterface() {
		t.Logf("arg is a null interface, arg[%v]", m.argType)
	} else {
		t.Errorf("arg is not a null interface, arg[%v]", m.argType)
	}

	if !m.ReplyIsNullInterface() {
		t.Logf("reply is not a null interface, reply[%v]", m.replyType)
	} else {
		t.Errorf("reply is a null interface, reply[%v]", m.replyType)
	}
}
*/
