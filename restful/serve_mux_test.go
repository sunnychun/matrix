package restful

import (
	"context"
	"errors"
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
	m = NewServeMux()
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

func ServeHTTP(h http.Handler, method, path string) (*httptest.ResponseRecorder, error) {
	r, err := http.NewRequest(method, path, nil)
	if err != nil {
		return nil, err
	}
	//r.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w, nil
}

func TestServeMux(t *testing.T) {
	m, err := NewArithServeMux()
	if err != nil {
		t.Fatal(err)
	}

	r, err := ServeHTTP(m, "XPOST", "/arith/add")
	if err != nil {
		t.Fatal(err)
	}
	if r.Code != http.StatusOK {
		t.Errorf("%s: %d != %s: %v\t%s", http.StatusText(r.Code), r.Code, http.StatusText(http.StatusOK), http.StatusOK, r.Body.String())
	}
}
