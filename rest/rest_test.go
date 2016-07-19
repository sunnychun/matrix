package rest

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang.org/x/net/context"
)

func TestRest(t *testing.T) {
	type Request struct {
		Echo string
	}
	type Response struct {
		Echo string
	}
	echo := func(ctx context.Context, vars Vars, req *Request, resp *Response) error {
		resp.Echo = req.Echo
		return nil
	}
	echo1 := func(ctx context.Context, vars Vars, req interface{}, resp *Response) error {
		resp.Echo = vars.Get(":echo")
		return errors.New("echo1 error")
	}

	r := New()
	r.SetLogProto(true)
	r.SetLogger(log.New(os.Stdout, "[rest-log] ", log.LstdFlags))
	Must(r.Get("/echo", echo))
	Must(r.Get("/echo/:echo", echo1))

	testfunc := func(method, urlstr, body string) {
		req, err := http.NewRequest(method, urlstr, bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("new request: %v", err)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("status code is not ok, code[%d]", w.Code)
		}
		buf, _ := ioutil.ReadAll(w.Body)
		t.Logf("Status: %v", w.Code)
		t.Logf("Header: %v", w.Header())
		t.Logf("Body: %s", string(buf))
	}

	testfunc("GET", "/echo", `{"Echo":"hello, world"}`)
	testfunc("GET", "/echo/hello,world", ``)
}
