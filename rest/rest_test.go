package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"
)

func TestRest(t *testing.T) {
	f := func(ctx context.Context, vars Vars, req interface{}, resp interface{}) error {
		fmt.Printf("uid=%q, phone=%q\n", vars.Get(":uid"), vars.Get(":phone"))
		return nil
	}

	r := New()
	Must(r.Get("/account/uid/:uid", f))
	Must(r.Get("/account/phone/:phone", f))

	testfunc := func(method, urlstr string) {
		req, err := http.NewRequest(method, urlstr, nil)
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

	testfunc("GET", "/account/uid/1")
	testfunc("GET", "/account/phone/13564171399")
}
