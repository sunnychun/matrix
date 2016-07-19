package rest

import (
	"ac-common-go/net/context"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/bmizerany/pat"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

type Rest struct {
	log      Logger
	logProto bool

	mux          *pat.PatternServeMux
	codecs       map[string]Codec
	defaultCodec string
}

func New() *Rest {
	r := &Rest{
		log:      noOpLogger{},
		logProto: false,

		mux:    pat.New(),
		codecs: make(map[string]Codec),
	}
	r.defaultCodec = "text/json"
	r.AddCodec("text/json", JsonCodec{})
	return r
}

func (r *Rest) SetLogger(log Logger) {
	if log == nil {
		r.log = noOpLogger{}
	}
	r.log = log
}

func (r *Rest) SetLogProto(isLog bool) {
	r.logProto = isLog
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
		r.serve(meth, pat, m, w, req)
	}
	r.mux.Add(meth, pat, http.HandlerFunc(fn))
	return nil
}

func (r *Rest) serve(meth, pat string, m *method, w http.ResponseWriter, req *http.Request) {
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
		r.log.Printf("(%q %q) unsupport Content-Type: %q", meth, pat, contentType)
		setErr(w, fmt.Errorf("unsupport Content-Type: %q", contentType))
		return
	}

	var err error
	argv := m.NewArg()
	replyv := m.NewReply()

	//Decode
	if !m.ArgIsNullInterface() {
		if req.Body == nil {
			r.log.Printf("(%q %q) request body is nil", meth, pat)
			setErr(w, errors.New("request body is nil"))
			return
		}
		if err = codec.Decode(req.Body, argv.Interface()); err != nil {
			r.log.Printf("(%q %q) decode: %v", meth, pat, err)
			setErr(w, fmt.Errorf("decode: %v", err))
			return
		}
	}

	//Call
	vars := Vars(req.URL.Query())
	ctx := context.Background()
	if err = m.Call(ctx, vars, argv, replyv); err != nil {
		r.log.Printf("(%q %q) call: %v", meth, pat, err)
		setErr(w, fmt.Errorf("call: %v", err))
		return
	}

	//Encode
	if !m.ReplyIsNullInterface() {
		if err = codec.Encode(w, replyv.Interface()); err != nil {
			r.log.Printf("(%q %q) encode: %v", meth, pat, err)
			setErr(w, fmt.Errorf("encode: %v", err))
			return
		}
	}

	setMsgName(w.Header(), "Ack")
}

func (r *Rest) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.logProto {
		_, ok := r.log.(noOpLogger)
		if !ok {
			lrw := &logResponseWriter{ResponseWriter: w, proto: req.Proto, status: http.StatusOK}
			w = lrw

			start := time.Now()
			printRequest(r.log, start, req)
			defer printResponse(r.log, start, lrw)
		}
	}

	r.mux.ServeHTTP(w, req)
}

func setErr(w http.ResponseWriter, err error) {
	setMsgName(w.Header(), "Err")
	setContentType(w.Header(), "text/plain")
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

var sep = []byte("\r\n" + strings.Repeat("-", 80) + "\r\n")

func printRequest(logger Logger, start time.Time, req *http.Request) {
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dump request: %v\r\n", err)
		return
	}

	var w bytes.Buffer
	fmt.Fprintf(&w, "[%s] Request:\r\n", start)
	w.Write(dump)
	w.Write(sep)
	logger.Printf(w.String())
}

func printResponse(logger Logger, start time.Time, lrw *logResponseWriter) {
	end := time.Now()
	dump := lrw.DumpResponse()

	var w bytes.Buffer
	fmt.Fprintf(&w, "[%s][%s] Response:\r\n", end, end.Sub(start))
	w.Write(dump)
	w.Write(sep)
	logger.Printf(w.String())
}

type logResponseWriter struct {
	http.ResponseWriter
	proto  string
	status int
	buffer bytes.Buffer
}

func (w *logResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *logResponseWriter) Write(p []byte) (int, error) {
	w.buffer.Write(p)
	return w.ResponseWriter.Write(p)
}

func (w *logResponseWriter) DumpResponse() []byte {
	var out bytes.Buffer
	fmt.Fprintf(&out, "%s %d %s\r\n", w.proto, w.status, http.StatusText(w.status))
	if len(w.Header()) > 0 {
		w.Header().Write(&out)
	}
	fmt.Fprintf(&out, "\r\n")
	if w.buffer.Len() > 0 {
		io.Copy(&out, &w.buffer)
	}
	return out.Bytes()
}
