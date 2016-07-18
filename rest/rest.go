package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/bmizerany/pat"
	"golang.org/x/net/context"
)

type Rest struct {
	mux          *pat.PatternServeMux
	codecs       map[string]Codec
	defaultCodec string
}

func New() *Rest {
	r := &Rest{
		mux:    pat.New(),
		codecs: make(map[string]Codec),
	}
	r.defaultCodec = "text/json"
	r.AddCodec("text/json", JsonCodec{})
	return r
}

func (rest *Rest) AddCodec(name string, codec Codec) {
	rest.codecs[name] = codec
}

func (rest *Rest) lookupCodec(name string) (Codec, bool) {
	c, ok := rest.codecs[name]
	return c, ok
}

func (rest *Rest) Head(pat string, api interface{}) error {
	return rest.add("HEAD", pat, api)
}

func (rest *Rest) Get(pat string, api interface{}) error {
	return rest.add("GET", pat, api)
}

func (rest *Rest) Put(pat string, api interface{}) error {
	return rest.add("PUT", pat, api)
}

func (rest *Rest) Post(pat string, api interface{}) error {
	return rest.add("POST", pat, api)
}

func (rest *Rest) Delete(pat string, api interface{}) error {
	return rest.add("DELETE", pat, api)
}

func (rest *Rest) Options(pat string, api interface{}) error {
	return rest.add("OPTIONS", pat, api)
}

func (rest *Rest) Patch(pat string, api interface{}) error {
	return rest.add("PATCH", pat, api)
}

func (rest *Rest) add(meth, pat string, api interface{}) error {
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
	sequence := getSequence(r.Header)
	contentType := getContentType(r.Header)
	if contentType == "" {
		contentType = rest.defaultCodec
	}

	setSequence(w.Header(), sequence)
	setContentType(w.Header(), contentType)

	//lookup codec
	codec, ok := rest.lookupCodec(contentType)
	if !ok {
		setErr(w, http.StatusBadRequest, fmt.Errorf("unsupport Content-Type: %q", contentType))
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
	if err = m.Call(ctx, argv, replyv); err != nil {
		setErr(w, http.StatusInternalServerError, err)
		return
	}

	//Encode
	if !m.ReplyIsNullInterface() {
		if err = codec.Encode(w, replyv.Interface()); err != nil {
			setErr(w, http.StatusInternalServerError, err)
			return
		}
	}

	setMsgName(w.Header(), "Ack")
}

func setErr(w http.ResponseWriter, status int, err error) {
	setMsgName(w.Header(), "Err")
	setContentType(w.Header(), "text/plain")
	w.WriteHeader(status)
	fmt.Fprintf(w, err.Error())
}

func setMsgName(h http.Header, v string) {
	if v != "" {
		h.Set("X-Msg-Name", v)
	}
}

func getSequence(h http.Header) string {
	return h.Get("X-Sequence")
}

func setSequence(h http.Header, v string) {
	if v != "" {
		h.Set("X-Sequence", v)
	}
}

func getContentType(h http.Header) string {
	return h.Get("Content-Type")
}

func setContentType(h http.Header, v string) {
	if v != "" {
		h.Set("Content-Type", v)
	}
}
