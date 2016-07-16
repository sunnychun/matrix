package rest

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/bmizerany/pat"
	"golang.org/x/net/context"
)

type Codec interface {
	Decode(io.Reader, interface{}) error
	Encode(io.Writer, interface{}) error
}

type Rest struct {
	mux    *pat.PatternServeMux
	codecs map[string]Codec
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
	contentType := getContentType(r.Header)

	//lookup codec
	codec, ok := rest.lookupCodec(contentType)
	if !ok {
		setErr(w, http.StatusBadRequest, fmt.Errorf("unsupport Content-Type: %s", contentType))
		return
	}

	var err error
	argv := m.NewArg()
	replyv := m.NewReply()

	//Decode
	if !m.ArgIsNullInterface() {
		if r.Body == nil {
			setErr(w, http.StatusBadRequest, errors.New("body is nil"))
			return
		}
		if err = codec.Decode(r.Body, argv.Interface()); err != nil {
			setErr(w, http.StatusBadRequest, err)
			return
		}
	}

	//Call
	values := r.URL.Query()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "values", values)
	ctx = context.WithValue(ctx, "header", r.Header)
	if err = m.Call(ctx, argv, replyv); err != nil {
		setErr(w, http.StatusInternalServerError, err)
		return
	}

	setContentType(w.Header(), contentType)

	//Encode
	if !m.ReplyIsNullInterface() {
		if err = codec.Encode(w, replyv.Interface()); err != nil {
			setErr(w, http.StatusInternalServerError, err)
			return
		}
	}
}

func (rest *Rest) lookupCodec(name string) (Codec, bool) {
	c, ok := rest.codecs[name]
	return c, ok
}

func getSeq(h http.Header) string {
	return h.Get("X-Sequence")
}

func setSeq(h http.Header, v string) {
	h.Set("X-Sequence", v)
}

func getContentType(h http.Header) string {
	return h.Get("Content-Type")
}

func setContentType(h http.Header, v string) {
	h.Set("Content-Type", v)
}

func setErr(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
}
