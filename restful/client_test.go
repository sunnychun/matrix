package restful

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ironzhang/matrix/codes"
	"github.com/ironzhang/matrix/context-value"
	"github.com/ironzhang/matrix/tlog"
)

func TestClientDoContext(t *testing.T) {
	m, err := NewArithServeMux()
	if err != nil {
		t.Fatal(err)
	}
	m.SetVerbose(0)
	s := httptest.NewServer(m)

	var c Client
	var args Args
	var reply Reply
	c.Verbose = 1
	ctx := context_value.WithVerbose(context.Background(), false)

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
	}
	for i, tt := range tests {
		args.A, args.B = tt.a, tt.b
		if err = c.DoContext(ctx, tt.method, s.URL+tt.path, args, &reply); err != nil {
			t.Errorf("tests[%d]: do context: %v", i, err)
			continue
		}
		if reply.C != tt.c {
			t.Errorf("tests[%d]: got(%d) != want(%d)", i, reply.C, tt.c)
			continue
		}
	}
}

func TestClientDoContextReturnErr(t *testing.T) {
	tlog.Init(tlog.Config{DisableStderr: true})

	m, err := NewTesterServeMux()
	if err != nil {
		t.Fatal(err)
	}
	m.SetVerbose(0)
	s := httptest.NewServer(m)

	var c Client
	c.Verbose = 1
	ctx := context_value.WithVerbose(context.Background(), false)

	tests := []struct {
		method string
		path   string
		status int
		code   codes.Code
		cause  string
	}{
		{"POST", "/NotFound", http.StatusNotFound, codes.NotFound, "page(/NotFound) not found"},
		{"POST", "/ReturnNil", http.StatusMethodNotAllowed, codes.NotAllowed, "method(POST) not allowed"},
		{"GET", "/ReturnDecodeFailError", http.StatusBadRequest, codes.DecodeFail, "EOF"},
		{"GET", "/ReturnInternalError", http.StatusInternalServerError, codes.Internal, "internal error"},
		{"GET", "/ReturnInvalidParamError", http.StatusBadRequest, codes.InvalidParam, ""},
		{"GET", "/ReturnOutOfRangeErrorWithCause", http.StatusInternalServerError, codes.OutOfRange, "out of range"},
	}
	for i, tt := range tests {
		err = c.DoContext(ctx, tt.method, s.URL+tt.path, nil, nil)
		if err == nil {
			t.Errorf("tests[%d]: client do context expect error buf not", i)
			continue
		}
		e, ok := err.(Error)
		if !ok {
			t.Errorf("tests[%d]: return error type not %T", i, Error{})
			continue
		}

		if got, want := e.Status, tt.status; got != want {
			t.Errorf("tests[%d]: error status: %v != %v", i, got, want)
		}
		if got, want := e.Code, tt.code; got != want {
			t.Errorf("tests[%d]: error code: %v != %v", i, got, want)
		}
		if got, want := e.Cause, tt.cause; got != want {
			t.Errorf("tests[%d]: error cause: %q != %q", i, got, want)
		}
	}
}
