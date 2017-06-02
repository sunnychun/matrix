package restful

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/ironzhang/matrix/codes"
	"github.com/ironzhang/matrix/context-value"
	"github.com/ironzhang/matrix/httputils"
	"github.com/ironzhang/matrix/restful/codec"
	"github.com/ironzhang/matrix/tlog"
	"github.com/ironzhang/matrix/uuid"
)

func NewServeMux(c codec.Codec) *ServeMux {
	if c == nil {
		c = codec.DefaultCodec
	}
	return &ServeMux{
		verbose:  1,
		codec:    c,
		patterns: make([]*pattern, 0),
	}
}

type ServeMux struct {
	w        io.Writer
	verbose  int
	codec    codec.Codec
	patterns []*pattern
}

// SetVerbose 设置verbose级别
//  verbose = 0, 不打印HTTP协议
//  verbose = 1, 根据请求头部中是否含有X-Verbose来决定是否打印HTTP协议，默认级别
//  verbose = 2, 打印HTTP协议
func (m *ServeMux) SetVerbose(verbose int) {
	m.verbose = verbose
}

func (m *ServeMux) SetWriter(w io.Writer) {
	m.w = w
}

func (m *ServeMux) Delete(pat string, h interface{}) error {
	return m.Add("DELETE", pat, h)
}

func (m *ServeMux) Get(pat string, h interface{}) error {
	return m.Add("GET", pat, h)
}

func (m *ServeMux) Head(pat string, h interface{}) error {
	return m.Add("HEAD", pat, h)
}

func (m *ServeMux) Options(pat string, h interface{}) error {
	return m.Add("OPTIONS", pat, h)
}

func (m *ServeMux) Patch(pat string, h interface{}) error {
	return m.Add("PATCH", pat, h)
}

func (m *ServeMux) Post(pat string, h interface{}) error {
	return m.Add("POST", pat, h)
}

func (m *ServeMux) Put(pat string, h interface{}) error {
	return m.Add("PUT", pat, h)
}

func (m *ServeMux) Add(meth, pat string, i interface{}) error {
	h, err := parseHandler(i)
	if err != nil {
		return fmt.Errorf("parse handler: %v", err)
	}

	for _, p := range m.patterns {
		if p.pat == pat {
			p.add(meth, h)
			return nil
		}
	}
	p := newPattern(pat)
	p.add(meth, h)
	m.patterns = append(m.patterns, p)
	return nil
}

func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context_value.WithTraceId(context.Background(), getTraceId(r.Header))

	// print verbose proto
	if m.getVerbose(r.Header) {
		m.printRequest(ctx, r)
		d := httputils.NewResponseDumper(w, r)
		defer m.printResponse(ctx, d)

		w = d
	}

	// serve http
	if err := m.serveHTTP(ctx, w, r); err != nil {
		m.setError(w, err)
	}
}

func (m *ServeMux) serveHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	log := tlog.WithContext(ctx).Sugar().With("method", r.Method, "path", r.URL.Path)

	var found bool
	for _, p := range m.patterns {
		values, ok := p.try(r.URL.Path)
		if !ok {
			continue
		}
		found = true
		h, ok := p.get(r.Method)
		if !ok {
			continue
		}
		r.URL.RawQuery = values.Encode() + "&" + r.URL.RawQuery
		return m.serve(ctx, h, w, r)
	}

	if found {
		log.Info(http.StatusText(http.StatusMethodNotAllowed))
		return Errorf(http.StatusMethodNotAllowed, codes.NotAllowed, "method(%s) not allowed", r.Method)
	} else {
		log.Info(http.StatusText(http.StatusNotFound))
		return Errorf(http.StatusNotFound, codes.NotFound, "page(%s) not found", r.URL.Path)
	}
}

func (m *ServeMux) serve(ctx context.Context, h *handler, w http.ResponseWriter, r *http.Request) (err error) {
	log := tlog.WithContext(ctx).Sugar().With("method", r.Method, "path", r.URL.Path)

	// check Content-Type
	if err = m.checkContentType(r.Header); err != nil {
		log.Infow("check content type", "error", err)
		return Errorf(http.StatusBadRequest, codes.InvalidHeader, err.Error())
	}

	in1 := newReflectValue(h.in1Type)
	in2 := newReflectValue(h.in2Type)

	// Decode
	if !isNilInterface(h.in1Type) {
		if err = m.codec.Decode(r.Body, in1.Interface()); err != nil {
			log.Infow("decode", "error", err)
			return Errorf(http.StatusBadRequest, codes.DecodeFail, err.Error())
		}
	}

	// Handle
	if err = h.Handle(ctx, in1, in2); err != nil {
		log.Infow("handle", "error", err)
		return err
	}

	// Encode
	if !isNilInterface(h.in2Type) {
		w.Header().Set("Content-Type", m.codec.ContentType())
		if err = m.codec.Encode(w, in2.Interface()); err != nil {
			log.Errorw("encode", "error", err)
			return Errorf(http.StatusInternalServerError, codes.EncodeFail, err.Error())
		}
	}

	return nil
}

func (m *ServeMux) checkContentType(h http.Header) error {
	if v := h.Get("Content-Type"); v != "" && v != m.codec.ContentType() {
		return fmt.Errorf("Content-Type not %s: %s", m.codec.ContentType(), v)
	}
	return nil
}

func (m *ServeMux) setError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if te, ok := err.(HTTPStatus); ok {
		status = te.HTTPStatus()
	}
	e := codec.ToError(err)
	w.Header().Set("Content-Type", m.codec.ContentType())
	w.WriteHeader(status)
	m.codec.EncodeError(w, e)
}

func (m *ServeMux) writer() io.Writer {
	if m.w == nil {
		return os.Stdout
	}
	return m.w
}

func (m *ServeMux) printRequest(ctx context.Context, r *http.Request) {
	b, err := httputil.DumpRequest(r, true)
	if err != nil {
		tlog.WithContext(ctx).Sugar().Errorw("dump request", "error", err)
		return
	}
	traceId := context_value.ParseTraceId(ctx)
	fmt.Fprintf(m.writer(), "traceId(%s) server request:\n%s\n", traceId, b)
}

func (m *ServeMux) printResponse(ctx context.Context, r *httputils.ResponseDumper) {
	b := r.Dump(true)
	traceId := context_value.ParseTraceId(ctx)
	fmt.Fprintf(m.writer(), "traceId(%s) server response:\n%s\n", traceId, b)
}

func (m *ServeMux) getVerbose(h http.Header) bool {
	switch m.verbose {
	case 0:
		return false
	case 1:
		if v := h.Get(xVerbose); v == "1" || strings.ToLower(v) == "true" {
			return true
		}
		return false
	case 2:
		return true
	}
	return false
}

func getTraceId(h http.Header) string {
	if v := h.Get(xTraceId); v != "" {
		return v
	}
	return uuid.New().String()
}
