package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/bmizerany/pat"
	"golang.org/x/net/context"
)

type Encoder interface {
	Encode(io.Reader, interface{}) error
}

type Decoder interface {
	Decode(io.Reader, interface{}) error
}

type Codec interface {
	Encoder
	Decoder
}

type Rest struct {
	mux *pat.PatternServeMux
}

func New() *Rest {
	return &Rest{mux: pat.New()}
}

/*
func (r *Rest) Post(pattern string, fun interface{}) error {
}
*/

func (rest *Rest) register(meth, pat string, api interface{}) error {
	m, err := parseMethod(api)
	if err != nil {
		return fmt.Errorf("parse method: %v", err)
	}
	fn := func(w http.ResponseWriter, r *http.Request) {
		rest.serve(m, w, r)
	}
	rest.mux.Add(meth, pat, http.HandlerFunc(fn))
	return nil
}

func (rest *Rest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rest.mux.ServeHTTP(w, r)
}

func (rest *Rest) serve(m *method, w http.ResponseWriter, r *http.Request) {
	argv := m.NewArg()
	replyv := m.NewReply()

	var err error
	if !m.ArgIsNullInterface() {
		if r.Body == nil {
			setErr(w, http.StatusBadRequest, errors.New("body is nil"))
			return
		}
		dec := json.NewDecoder(r.Body)
		if err = dec.Decode(argv.Interface()); err != nil {
			setErr(w, http.StatusBadRequest, err)
			return
		}
	}

	values := r.URL.Query()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "values", values)
	ctx = context.WithValue(ctx, "header", r.Header)
	if err = m.Call(ctx, argv, replyv); err != nil {
		setErr(w, http.StatusInternalServerError, err)
		return
	}

	if !m.ReplyIsNullInterface() {
		enc := json.NewEncoder(w)
		if err = enc.Encode(replyv.Interface()); err != nil {
			setErr(w, http.StatusInternalServerError, err)
			return
		}
	}
}

func setErr(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
}
