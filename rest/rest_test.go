package rest

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Response struct {
	A int
}

func TestRestRegister(t *testing.T) {
	rest := New()
	err := rest.register("GET", "/test", func(_ interface{}, a interface{}, resp *Response) error {
		resp.A = 1
		return errors.New("test error")
	})
	if err != nil {
		t.Fatalf("rest register: %v", err)
	}

	r, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	w := httptest.NewRecorder()
	rest.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("status code is not ok, code[%d]", w.Code)
	}

	buf, _ := ioutil.ReadAll(w.Body)
	t.Logf("Status: %v", w.Code)
	t.Logf("Header: %v", w.Header())
	t.Logf("Body: %s", string(buf))
}
