package restful

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ironzhang/matrix/codes"
)

func ServeHTTP(h http.Handler, method, path string, b []byte) (*httptest.ResponseRecorder, error) {
	r, err := http.NewRequest(method, path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w, nil
}

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith int

func (t *Arith) Add(ctx context.Context, args Args, reply *Reply) error {
	reply.C = args.A + args.B
	return nil
}

func (t *Arith) Sub(ctx context.Context, args Args, reply *Reply) error {
	reply.C = args.A - args.B
	return nil
}

func (t *Arith) Mul(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func (t *Arith) Div(ctx context.Context, args Args, reply *Reply) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	reply.C = args.A / args.B
	return nil
}

func NewArithServeMux() (m *ServeMux, err error) {
	var a Arith
	m = NewServeMux(nil)
	if err = m.Add("POST", "/add", a.Add); err != nil {
		return nil, err
	}
	if err = m.Add("POST", "/sub", a.Sub); err != nil {
		return nil, err
	}
	if err = m.Add("POST", "/mul", a.Mul); err != nil {
		return nil, err
	}
	if err = m.Add("POST", "/div", a.Div); err != nil {
		return nil, err
	}
	return m, nil
}

func CallArith(h http.Handler, method, path string, a, b int) (c int, err error) {
	args := Args{a, b}
	buf, err := json.Marshal(args)
	if err != nil {
		return 0, err
	}

	r, err := ServeHTTP(h, method, path, buf)
	if err != nil {
		return 0, err
	}
	if r.Code != http.StatusOK {
		var e rpcError
		if err = json.Unmarshal(r.Body.Bytes(), &e); err != nil {
			return 0, err
		}
		return 0, Errorf(r.Code, codes.Code(e.Code), e.Cause)
	}

	var reply Reply
	if err = json.Unmarshal(r.Body.Bytes(), &reply); err != nil {
		return 0, err
	}
	return reply.C, nil
}

func TestArithServeMux(t *testing.T) {
	m, err := NewArithServeMux()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		method string
		path   string
		a      int
		b      int
		c      int
	}{
		{method: "POST", path: "/add", a: 1, b: 2, c: 3},
		{method: "Post", path: "/sub", a: 1, b: 2, c: -1},
		{method: "post", path: "/mul", a: 1, b: 2, c: 2},
		{method: "post", path: "/div", a: 1, b: 2, c: 0},
		//{method: "post", path: "/div", a: 1, b: 0, c: 0},
	}

	for i, tt := range tests {
		c, err := CallArith(m, tt.method, tt.path, tt.a, tt.b)
		if err != nil {
			t.Errorf("tests[%d]: call arith: %v", i, err)
			continue
		}
		if c != tt.c {
			t.Errorf("tests[%d]: got(%v) != want(%v)", i, c, tt.c)
		}
	}
}
