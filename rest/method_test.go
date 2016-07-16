package rest

import (
	"ac"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

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
