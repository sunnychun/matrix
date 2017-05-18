package restful

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	if err = m.Add("POST", "/arith/add", a.Add); err != nil {
		return nil, err
	}
	if err = m.Add("POST", "/arith/sub", a.Sub); err != nil {
		return nil, err
	}
	if err = m.Add("POST", "/arith/mul", a.Mul); err != nil {
		return nil, err
	}
	if err = m.Add("POST", "/arith/div", a.Div); err != nil {
		return nil, err
	}
	return m, nil
}

func ServeHTTP(h http.Handler, method, path string, b []byte) (*httptest.ResponseRecorder, error) {
	r, err := http.NewRequest(method, path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w, nil
}

func TestServeMux(t *testing.T) {
	m, err := NewArithServeMux()
	if err != nil {
		t.Fatal(err)
	}

	var args Args
	var reply Reply
	args.A = 1
	args.B = 2
	buf, err := json.Marshal(args)
	if err != nil {
		t.Fatal(err)
	}

	r, err := ServeHTTP(m, "POST", "/arith/add", buf)
	if err != nil {
		t.Fatal(err)
	}
	if r.Code != http.StatusOK {
		t.Fatalf("%s: %d\t%s", http.StatusText(r.Code), r.Code, r.Body.String())
	}
	if r.HeaderMap.Get("Content-Type") != contentType {
		t.Fatalf("Content-Type: %s != %s", r.HeaderMap.Get("Content-Type"), contentType)
	}
	if err = json.Unmarshal(r.Body.Bytes(), &reply); err != nil {
		t.Fatal(err)
	}
	if reply.C != args.A+args.B {
		t.Errorf("C(%d) != A(%d) + B(%d)", reply.C, args.A, args.B)
	} else {
		fmt.Printf("C(%d) == A(%d) + B(%d)\n", reply.C, args.A, args.B)
	}
}
