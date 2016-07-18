package rest

import (
	"ac-common-go/net/context"
	"errors"
	"fmt"
	"net/http"

	"github.com/bmizerany/pat"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

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

func (r *Rest) AddCodec(name string, codec Codec) {
	r.codecs[name] = codec
}

func (r *Rest) lookupCodec(name string) (Codec, bool) {
	c, ok := r.codecs[name]
	return c, ok
}

func (r *Rest) Head(pat string, api interface{}) error {
	return r.add("HEAD", pat, api)
}

func (r *Rest) Get(pat string, api interface{}) error {
	return r.add("GET", pat, api)
}

func (r *Rest) Put(pat string, api interface{}) error {
	return r.add("PUT", pat, api)
}

func (r *Rest) Post(pat string, api interface{}) error {
	return r.add("POST", pat, api)
}

func (r *Rest) Delete(pat string, api interface{}) error {
	return r.add("DELETE", pat, api)
}

func (r *Rest) Options(pat string, api interface{}) error {
	return r.add("OPTIONS", pat, api)
}

func (r *Rest) Patch(pat string, api interface{}) error {
	return r.add("PATCH", pat, api)
}

func (r *Rest) add(meth, pat string, api interface{}) error {
	m, err := parseMethod(api)
	if err != nil {
		return fmt.Errorf("parse method: %v", err)
	}
	fn := func(w http.ResponseWriter, req *http.Request) {
		r.serve(m, w, req)
	}
	r.mux.Add(meth, pat, http.HandlerFunc(fn))
	return nil
}

func (r *Rest) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Rest) serve(m *method, w http.ResponseWriter, req *http.Request) {
	sequence := getSequence(req.Header)
	contentType := getContentType(req.Header)
	if contentType == "" {
		contentType = r.defaultCodec
	}

	setSequence(w.Header(), sequence)
	setContentType(w.Header(), contentType)

	//lookup codec
	codec, ok := r.lookupCodec(contentType)
	if !ok {
		setErr(w, http.StatusBadRequest, fmt.Errorf("unsupport Content-Type: %q", contentType))
		return
	}

	var err error
	argv := m.NewArg()
	replyv := m.NewReply()

	//Decode
	if !m.ArgIsNullInterface() {
		if req.Body == nil {
			setErr(w, http.StatusBadRequest, errors.New("body is nil"))
			return
		}
		if err = codec.Decode(req.Body, argv.Interface()); err != nil {
			setErr(w, http.StatusBadRequest, err)
			return
		}
	}

	//Call
	vars := Vars(req.URL.Query())
	ctx := context.Background()
	if err = m.Call(ctx, vars, argv, replyv); err != nil {
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
